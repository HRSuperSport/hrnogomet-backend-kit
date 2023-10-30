package test

import (
	"context"
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"os"
	"time"
)

const (
	MariaDBImage = "mariadb:10.5.2"
	MariaDBPort  = "3306"
)

type MariadbDBContainer struct {
	Container           testcontainers.Container
	Context             context.Context
	URI                 string
	URIDockerCompatible string
	Client              *sql.DB
	DbName              string
}

func init() {
	log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	log.Logger = log.Logger.With().Caller().Logger()
}

func SetupMariaDB(ctx context.Context, dbUser string, dbPwd string, dbRootPwd string, dbName string) *MariadbDBContainer {
	req := testcontainers.ContainerRequest{
		Image:        MariaDBImage,
		ExposedPorts: []string{"3306/tcp", "33060/tcp"},
		Env: map[string]string{
			"MYSQL_USER":          dbUser,
			"MYSQL_PASSWORD":      dbPwd,
			"MYSQL_ROOT_PASSWORD": dbRootPwd,
			"MYSQL_DATABASE":      dbName,
		},
		WaitingFor: wait.ForAll(
			wait.ForLog("init process done. Ready for start up").WithStartupTimeout(5*time.Minute),
			wait.ForExec([]string{"mysqladmin", "ping", "-h", "localhost"}).
				WithPollInterval(2*time.Second).
				WithExitCodeMatcher(func(exitCode int) bool {
					return exitCode == 0
				}),
			wait.ForListeningPort(MariaDBPort),
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

	mappedPort, err := container.MappedPort(ctx, MariaDBPort)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot get port")
	}

	// why ?parseTime=true ? see:
	// https://stackoverflow.com/questions/26617957/how-to-scan-a-mysql-timestamp-value-into-a-time-time-variable
	uri := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPwd, hostIP, mappedPort.Port(), dbName)
	log.Info().Str("URI", uri).Msg("MariaDB connection string")
	uriDockerCompatible := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?parseTime=true", dbUser, dbPwd, "host.docker.internal", mappedPort.Port(), dbName)
	log.Info().Str("URI", uriDockerCompatible).Msg("MariaDB docker compatible connection string")

	db, err := sql.Open("mysql", uri)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot open connection")
	}

	err = db.Ping()
	if err != nil {
		log.Fatal().Err(err).Msg("cannot ping database")
	}

	return &MariadbDBContainer{
		Container:           container,
		Context:             ctx,
		Client:              db,
		URI:                 uri,
		URIDockerCompatible: uriDockerCompatible,
		DbName:              dbName,
	}
}

func (m *MariadbDBContainer) TeardownMariaDB() {
	err := m.Client.Close()
	if err != nil {
		log.Error().Err(err).Msg("cannot close client")
	}

	err = m.Container.Terminate(m.Context)
	if err != nil {
		log.Error().Err(err).Msg("cannot close mariadb")
	}
}

func (m *MariadbDBContainer) Cleanup() {
	_, err := m.Client.Exec(fmt.Sprintf("DROP DATABASE IF EXISTS %s", m.DbName))
	if err != nil {
		log.Error().Err(err).Msg("cannot drop database")
	}

	_, err = m.Client.Exec(fmt.Sprintf("CREATE DATABASE %s", m.DbName))
	if err != nil {
		log.Error().Err(err).Msg("cannot create database")
	}

	_, err = m.Client.Exec(fmt.Sprintf("USE %s", m.DbName))
	if err != nil {
		log.Error().Err(err).Msg("cannot use database")
	}
}
