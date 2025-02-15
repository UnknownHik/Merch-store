package services

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TxExecutor interface {
	RunWithTransaction(ctx context.Context, fn func(pgx.Tx) error) error
}

// DefaultTxExecutor управляет транзакциями в pgx
type DefaultTxExecutor struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

// NewTxExecutor создает новый менеджер транзакций
func NewTxExecutor(pool *pgxpool.Pool, logger *slog.Logger) *DefaultTxExecutor {
	return &DefaultTxExecutor{pool: pool, logger: logger}
}

// RunWithTransaction запускает функцию с транзакцией
func (t *DefaultTxExecutor) RunWithTransaction(ctx context.Context, fn func(tx pgx.Tx) error) (err error) {
	tx, err := t.pool.Begin(ctx)
	if err != nil {
		t.logger.Error("failed to start transaction", err)
		return fmt.Errorf("failed to start transaction: %w", err)
	}

	t.logger.Info("Transaction started")

	defer func() {
		if err != nil {
			if rollbackErr := tx.Rollback(ctx); rollbackErr != nil {
				t.logger.Error("failed to rollback transaction", "rollbackErr", rollbackErr)
			} else {
				t.logger.Info("Transaction rolled back")
			}
		}
	}()

	if err = fn(tx); err != nil {
		t.logger.Error("error transaction function", "error", err)
		return err
	}

	if commitErr := tx.Commit(ctx); commitErr != nil {
		t.logger.Error("failed to commit transaction", "error", commitErr)
		err = fmt.Errorf("failed to commit transaction: %w", commitErr)
		return err
	}

	t.logger.Info("Transaction committed")
	return nil
}
