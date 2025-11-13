package db

import (
	"context"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
)

type TxHandler func(ctx context.Context) error

// Sqlizer - something that can build sql query
type Sqlizer interface {
	ToSql() (sql string, args []interface{}, err error)
}

// Client incapsulates connections to
// different databases (master, slave)
type Client interface {
	DB() DB
	Close() error
}

type TxManager interface {
	ReadCommitted(ctx context.Context, fn TxHandler) error
}

// DB provides interface for working with db
type DB interface {
	QueryExecutor
	Transactor
	Close()
}

type Query struct {
	Name    string
	Sqlizer Sqlizer
}

type Transactor interface {
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

type QueryExecutor interface {
	Exec(ctx context.Context, q Query) (pgconn.CommandTag, error)
	Query(ctx context.Context, q Query) (pgx.Rows, error)
	QueryRow(ctx context.Context, q Query) (pgx.Row, error)

	Getx(ctx context.Context, dest interface{}, q Query) error
	Selectx(ctx context.Context, dest interface{}, q Query) error
}
