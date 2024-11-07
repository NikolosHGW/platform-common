package pg

import (
	"context"
	"fmt"

	"github.com/NikolosHGW/platform-common/pkg/db"
	"github.com/jmoiron/sqlx"
)

type pgClient struct {
	masterDBC db.DB
}

// New - конструктор клиента бд.
func New(ctx context.Context, dsn string) (db.Client, error) {
	db, err := sqlx.ConnectContext(ctx, "postgres", dsn)
	if err != nil {
		return nil, fmt.Errorf("не удалось настроить соединение с бд постгрес: %w", err)
	}

	return &pgClient{masterDBC: db}, nil
}

func (pgc *pgClient) DB() db.DB {
	return pgc.masterDBC
}

func (pgc *pgClient) Close() error {
	if pgc.masterDBC != nil {
		return pgc.masterDBC.Close()
	}

	return nil
}
