package pg

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/chistyakoviv/logbot/internal/db"
	"github.com/georgysavva/scany/v2/pgxscan"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
)

type key string

const (
	TxKey key = "tx"
)

type pg struct {
	name   string
	dbc    *pgxpool.Pool
	logger *slog.Logger
}

func NewDB(name string, dbc *pgxpool.Pool, logger *slog.Logger) db.DB {
	return &pg{
		name:   name,
		dbc:    dbc,
		logger: logger,
	}
}

func (p *pg) Exec(ctx context.Context, q db.Query) (pgconn.CommandTag, error) {
	query, args, err := q.Sqlizer.ToSql()
	if err != nil {
		return pgconn.CommandTag{}, fmt.Errorf("postgres: to sql: %w", err)
	}

	p.logger.Debug(
		"query debug",
		slog.Attr{Key: "name", Value: slog.StringValue(q.Name)},
		slog.Attr{Key: "sql", Value: slog.StringValue(query)},
	)

	// A transaction is initiated by calling txManager.ReadCommitted(ctx, func(ctx context.Context) error).
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.Exec(ctx, query, args...)
	}

	return p.dbc.Exec(ctx, query, args...)
}

func (p *pg) Query(ctx context.Context, q db.Query) (pgx.Rows, error) {
	query, args, err := q.Sqlizer.ToSql()
	if err != nil {
		return nil, fmt.Errorf("postgres: to sql: %w", err)
	}

	p.logger.Debug(
		"query debug",
		slog.Attr{Key: "name", Value: slog.StringValue(q.Name)},
		slog.Attr{Key: "sql", Value: slog.StringValue(query)},
	)

	// A transaction is initiated by calling txManager.ReadCommitted(ctx, func(ctx context.Context) error).
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.Query(ctx, query, args...)
	}

	return p.dbc.Query(ctx, query, args...)
}

func (p *pg) QueryRow(ctx context.Context, q db.Query) (pgx.Row, error) {
	query, args, err := q.Sqlizer.ToSql()
	if err != nil {
		return nil, fmt.Errorf("postgres: to sql: %w", err)
	}

	p.logger.Debug(
		"query debug",
		slog.Attr{Key: "name", Value: slog.StringValue(q.Name)},
		slog.Attr{Key: "sql", Value: slog.StringValue(query)},
	)

	// A transaction is initiated by calling txManager.ReadCommitted(ctx, func(ctx context.Context) error).
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return tx.QueryRow(ctx, query, args...), nil
	}

	return p.dbc.QueryRow(ctx, query, args...), nil
}

func (p *pg) Getx(ctx context.Context, dest interface{}, q db.Query) error {
	query, args, err := q.Sqlizer.ToSql()
	if err != nil {
		return fmt.Errorf("postgres: to sql: %w", err)
	}

	p.logger.Debug(
		"query debug",
		slog.Attr{Key: "name", Value: slog.StringValue(q.Name)},
		slog.Attr{Key: "sql", Value: slog.StringValue(query)},
	)

	// A transaction is initiated by calling txManager.ReadCommitted(ctx, func(ctx context.Context) error).
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return pgxscan.Get(ctx, tx, dest, query, args...)
	}

	return pgxscan.Get(ctx, p.dbc, dest, query, args...)
}

func (p *pg) Selectx(ctx context.Context, dest interface{}, q db.Query) error {
	query, args, err := q.Sqlizer.ToSql()
	if err != nil {
		return fmt.Errorf("postgres: to sql: %w", err)
	}

	p.logger.Debug(
		"query debug",
		slog.Attr{Key: "name", Value: slog.StringValue(q.Name)},
		slog.Attr{Key: "sql", Value: slog.StringValue(query)},
	)

	// A transaction is initiated by calling txManager.ReadCommitted(ctx, func(ctx context.Context) error).
	tx, ok := ctx.Value(TxKey).(pgx.Tx)
	if ok {
		return pgxscan.Get(ctx, tx, dest, query, args...)
	}

	return pgxscan.Select(ctx, p.dbc, dest, query, args...)
}

/**
* Begin acquires a connection from the Pool and starts a transaction.
* Unlike database/sql, the context only affects the begin command.
* i.e. there is no auto-rollback on context cancellation.
* Begin initiates a transaction block without explicitly setting
* a transaction mode for the block (see BeginTx with TxOptions
* if transaction mode is required). *pgxpool.Tx is returned,
* which implements the pgx.Tx interface. Commit or Rollback
* must be called on the returned transaction to finalize the transaction block.
 */
func (p *pg) Begin(ctx context.Context) (pgx.Tx, error) {
	return p.dbc.Begin(ctx)
}

func (p *pg) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return p.dbc.BeginTx(ctx, txOptions)
}

func (p *pg) Close() {
	p.dbc.Close()
}

func MakeContextTx(ctx context.Context, tx pgx.Tx) context.Context {
	return context.WithValue(ctx, TxKey, tx)
}
