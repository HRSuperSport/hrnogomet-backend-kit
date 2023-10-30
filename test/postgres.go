package test

import (
	"context"
	"database/sql"
	"embed"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/rs/zerolog/log"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"time"
)

const (
	DBImage = "postgres:10"
	DBPort  = "5432"
)

type PostgresDBContainer struct {
	Container           testcontainers.Container
	Context             context.Context
	URI                 string
	URIDockerCompatible string
	Client              *sql.DB
	DbName              string
}

func SetupPostgresDB(ctx context.Context, dbUser string, dbPwd string, dbName string) *PostgresDBContainer {
	// https://golang.testcontainers.org/modules/postgres/
	req := testcontainers.ContainerRequest{
		Image:        DBImage,
		ExposedPorts: []string{"5432/tcp"},
		Env: map[string]string{
			"POSTGRES_USER":     dbUser,
			"POSTGRES_PASSWORD": dbPwd,
			"POSTGRES_DB":       dbName,
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("database system is ready to accept connections").WithStartupTimeout(5*time.Second),
			wait.ForListeningPort(DBPort),
		).WithDeadline(1 * time.Minute),
	}

	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		log.Fatal().Err(err).Msg("cannot start container")
	}

	hostIP, err := container.Host(ctx)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot get host")
	}

	mappedPort, err := container.MappedPort(ctx, DBPort)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot get port")
	}

	uri := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPwd, hostIP, mappedPort.Port(), dbName)
	log.Info().Msgf("Postgresql connection string: %s", uri)
	uriDockerCompatible := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", dbUser, dbPwd, "host.docker.internal", mappedPort.Port(), dbName)
	log.Info().Msgf("Postgresql docker compatible connection string: %s", uriDockerCompatible)

	db, err := sql.Open("postgres", uri)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot open connection")
	}

	err = db.Ping()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot ping database")
	}

	return &PostgresDBContainer{
		Container:           container,
		Context:             ctx,
		Client:              db,
		URI:                 uri,
		URIDockerCompatible: uriDockerCompatible,
		DbName:              dbName,
	}
}

func (p *PostgresDBContainer) TeardownPostgresDB() {
	err := p.Client.Close()
	if err != nil {
		log.Error().Err(err).Msg("cannot close client")
	}

	err = p.Container.Terminate(p.Context)
	if err != nil {
		log.Error().Err(err).Msg("cannot close postgresql")
	}
}

func (p *PostgresDBContainer) Cleanup() {
	_, err := p.Client.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", p.DbName))
	if err != nil {
		log.Error().Err(err).Msg("cannot drop database")
	}

	_, err = p.Client.Exec(fmt.Sprintf("CREATE DATABASE %s", p.DbName))
	if err != nil {
		log.Error().Err(err).Msg("cannot create database")
	}

	_, err = p.Client.Exec(fmt.Sprintf("USE %s", p.DbName))
	if err != nil {
		log.Error().Err(err).Msg("cannot use database")
	}
}

// InitMigratePostgres is used to apply db migration scripts from embedded file system to given postgres database
func InitMigratePostgres(conn string, migrations embed.FS, migrationTable *string) {
	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot open migrations")
	}

	var uri string
	if migrationTable != nil {
		// see https://github.com/golang-migrate/migrate/blob/master/database/postgres/README.md
		uri = fmt.Sprintf("%s&x-migrations-table=%s", conn, *migrationTable)
	} else {
		uri = conn
	}

	m, err := migrate.NewWithSourceInstance("iofs", source, uri)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot open source migrations")
	}

	defer func(m *migrate.Migrate) {
		err, _ := m.Close()
		if err != nil {
			log.Fatal().Err(err).Msg("cannot close migration")
		}
	}(m)

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		log.Fatal().Err(err).Msg("cannot migrate")
	}
}
