package pg

import (
	"context"
	"database/sql"

	"github.com/NikolosHGW/platform-common/pkg/db"
	"github.com/jmoiron/sqlx"
)

type key string

const (
	// TxKey ключ для вытаскивания объекта tx из контекста.
	TxKey key = "tx"
)

type pg struct {
	dbc *sqlx.DB
}

// NewDB - конструктор для постгрес-обёртки.
func NewDB(dbc *sqlx.DB) db.DB {
	return &pg{dbc: dbc}
}

func (p *pg) NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error) {
	tx, ok := ctx.Value(TxKey).(sqlx.Tx)
	if ok {
		return tx.NamedExecContext(ctx, query, arg)
	}

	return p.dbc.NamedExecContext(ctx, query, arg)
}

func (p *pg) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	tx, ok := ctx.Value(TxKey).(sqlx.Tx)
	if ok {
		return tx.SelectContext(ctx, dest, query, args...)
	}

	return p.dbc.SelectContext(ctx, dest, query, args...)
}

func (p *pg) ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error) {
	tx, ok := ctx.Value(TxKey).(sqlx.Tx)
	if ok {
		return tx.ExecContext(ctx, query, args...)
	}

	return p.dbc.ExecContext(ctx, query, args...)
}

func (p *pg) QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row {
	tx, ok := ctx.Value(TxKey).(sqlx.Tx)
	if ok {
		return tx.QueryRowxContext(ctx, query, args...)
	}

	return p.dbc.QueryRowxContext(ctx, query, args...)
}

func (p *pg) PingContext(ctx context.Context) error {
	return p.dbc.PingContext(ctx)
}

func (p *pg) Close() error {
	return p.dbc.Close()
}

func (p *pg) BeginTxx(ctx context.Context, txOptions *sql.TxOptions) (*sqlx.Tx, error) {

	return p.dbc.BeginTxx(ctx, txOptions)
}

// MakeContextTx устанавливает объект tx в контекст.
func MakeContextTx(ctx context.Context, tx sqlx.Tx) context.Context {
	return context.WithValue(ctx, TxKey, tx)
}
