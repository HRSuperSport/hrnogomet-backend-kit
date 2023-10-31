package postgres

import (
	"context"
	"github.com/jackc/pgx/v5/pgxpool"
)

// DB represents access (via underlying pool object) to PostgreSQL database
type DB struct {
	*pgxpool.Pool
}

// NewPostgresDBFromUri creates new postgres connection from provided uri
func NewPostgresDBFromUri(ctx context.Context, uri string) (*DB, error) {
	db, err := pgxpool.New(ctx, uri)
	if err != nil {
		return nil, err
	}

	err = db.Ping(ctx)
	if err != nil {
		return nil, err
	}

	return &DB{
		db,
	}, nil
}

// Close closes the database connection
func (db *DB) Close() {
	db.Pool.Close()
}
