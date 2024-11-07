package transaction

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/NikolosHGW/platform-common/pkg/db"
	"github.com/NikolosHGW/platform-common/pkg/db/pg"
	"github.com/jmoiron/sqlx"
)

type manager struct {
	db db.Transactor
}

// NewTransactionManager создаёт новый менеджер транзакций.
func NewTransactionManager(db db.Transactor) db.TxManager {
	return &manager{
		db: db,
	}
}

func (m *manager) transaction(ctx context.Context, opts sql.TxOptions, fn db.Handler) (err error) {
	tx, ok := ctx.Value(pg.TxKey).(*sqlx.Tx)
	if ok {
		return fn(ctx)
	}

	tx, err = m.db.BeginTx(ctx, &opts)
	if err != nil {
		err = fmt.Errorf("не удалось начать транзакцию: %w", err)
	}

	ctx = pg.MakeContextTx(ctx, *tx)

	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("panic recovered")
		}

		if err != nil {
			if errRollback := tx.Rollback(); errRollback != nil {
				err = fmt.Errorf("транзакция откатана: %w", errRollback)
			}
		}

		if err == nil {
			err = tx.Commit()
			if err != nil {
				err = fmt.Errorf("не удалось закоммитить изменения: %w", err)
			}
		}
	}()

	if err = fn(ctx); err != nil {
		err = fmt.Errorf("ошибка выполнения кода внутри транзакции: %w", err)
	}

	return err
}

func (m *manager) ReadCommitted(ctx context.Context, f db.Handler) error {
	txOptions := sql.TxOptions{Isolation: sql.LevelReadCommitted}

	return m.transaction(ctx, txOptions, f)
}
