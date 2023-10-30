package test

import (
	"embed"
	"errors"
	"fmt"
	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	"github.com/golang-migrate/migrate/v4/source/iofs"
	"github.com/rs/zerolog/log"
)

func InitMigratePostgresql(conn string, migrations embed.FS, migrationTable *string) {
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
