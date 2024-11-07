package db

import (
	"context"
	"database/sql"

	"github.com/jmoiron/sqlx"
)

// Handler - тип функции, которой нужно обернуть код для транзакции.
type Handler func(ctx context.Context) error

// Client - клиента для работы с БД.
type Client interface {
	DB() DB
	Close() error
}

// TxManager - интерфейс для использования транзакции с уровнем изоляции read committed.
type TxManager interface {
	ReadCommitted(ctx context.Context, f Handler) error
}

// Query - обёртка над запросом, хранящая имя запроса и сам запрос.
type Query struct {
	Name     string
	QueryRaw string
}

// Transactor менеджер транзакций.
type Transactor interface {
	BeginTxx(ctx context.Context, opts *sql.TxOptions) (*sqlx.Tx, error)
}

// SQLExecer комбинирует NamedExecer и QueryExecer
type SQLExecer interface {
	NamedExecer
	QueryExecer
}

// NamedExecer - интерфейс для работы с именованными запросами с помощью тегов в структурах.
type NamedExecer interface {
	NamedExecContext(ctx context.Context, query string, arg interface{}) (sql.Result, error)
	SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error
}

// QueryExecer - интерфейс для работы с обычными запросами.
type QueryExecer interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	QueryRowxContext(ctx context.Context, query string, args ...interface{}) *sqlx.Row
}

// Pinger - интерфейс для проверки соединения с БД.
type Pinger interface {
	PingContext(ctx context.Context) error
}

// DB - интерфейс для работы с БД.
type DB interface {
	SQLExecer
	Pinger
	Transactor
	Close() error
}