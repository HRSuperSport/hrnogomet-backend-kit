package test

import (
	"embed"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/rs/zerolog/log"
	"strings"
)

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

// InitMigrateMariadb is used to apply db migration scripts from embedded file system to given mariadb database
func InitMigrateMariadb(conn string, migrations embed.FS, migrationTable *string) {
	source, err := iofs.New(migrations, "migrations")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot open migrations")
	}

	var uri string
	if migrationTable != nil {
		// see https://github.com/golang-migrate/migrate/blob/master/database/mysql/README.md
		uri = fmt.Sprintf("mysql://%s?x-migrations-table=%s", conn, *migrationTable)
	} else {
		uri = fmt.Sprintf("mysql://%s", conn)
	}

	// mysql://user:test@tcp(localhost:55250)/hrnogomet?parseTime=true
	// ->
	// mysql://user:test@tcp(localhost:55250)/hrnogomet
	// we are using ?parseTime=true in conn string so that mysql parses timestamps during sqlx scans
	// but somehow this does not go well with migrate so let's remove it at this particular place
	uri = strings.Replace(uri, "?parseTime=true", "", 1)

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
