package postgres

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"sorkin_bot/internal/config"
	"time"
)

type Client interface {
	Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error)
	Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	QueryRow(ctx context.Context, sql string, args ...any) pgx.Row
	Begin(ctx context.Context) (pgx.Tx, error)
	BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error)
}

type PgxClientWrapper struct {
	conn *pgx.Conn
}

func (c *PgxClientWrapper) Exec(ctx context.Context, sql string, arguments ...any) (pgconn.CommandTag, error) {
	return c.conn.Exec(ctx, sql, arguments...)
}

func (c *PgxClientWrapper) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	return c.conn.Query(ctx, sql, args)
}

func (c *PgxClientWrapper) QueryRow(ctx context.Context, sql string, args ...any) pgx.Row {
	return c.conn.QueryRow(ctx, sql, args...)
}

func (c *PgxClientWrapper) Begin(ctx context.Context) (pgx.Tx, error) {
	return c.conn.Begin(ctx)
}

func (c *PgxClientWrapper) BeginTx(ctx context.Context, txOptions pgx.TxOptions) (pgx.Tx, error) {
	return c.conn.BeginTx(ctx, txOptions)
}

func NewClient(ctx context.Context, pg config.StorageConfig) (client Client, err error) {
	const op = "referalMS.pkg.client.postgres.NewClient"
	DSN := fmt.Sprintf("postgresql://%s:%s@%s:%s/%s", pg.User, pg.Password, pg.Host, pg.Port, pg.Database)

	err = ConnectWithRetry(func() error {
		ctx, cancel := context.WithTimeout(ctx, pg.RetryTimeout)
		defer cancel()

		client, err = pgx.Connect(ctx, DSN)
		if err != nil {
			return err
		}
		return nil
	}, pg.MaxRetry, pg.RetryTimeout)

	if err != nil {
		return client, err
	}

	return client, nil
}

func ConnectWithRetry(fn func() error, maxRetry int, timeout time.Duration) (err error) {
	for maxRetry > 0 {
		if err = fn(); err != nil {
			time.Sleep(timeout)
			maxRetry--
			continue
		}
		return nil
	}
	return
}
