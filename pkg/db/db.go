package db

import (
	"context"
	"database/sql"

	_ "modernc.org/sqlite"
)

type Database struct {
	db *sql.DB
}

func (d *Database) Close() error {
	return d.db.Close()
}

func (d *Database) Conn() *sql.DB {
	return d.db
}

func (d *Database) BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error) {
	return d.db.BeginTx(ctx, opts)
}

func (d *Database) PingContext(ctx context.Context) error {
	return d.db.PingContext(ctx)
}

type Executor interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sql.Rows, error)
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

func NewDatabase(dsn string) (*Database, error) {
	conn, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, err
	}

	return &Database{db: conn}, nil
}
