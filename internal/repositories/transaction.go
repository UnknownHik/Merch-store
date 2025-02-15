package repositories

import (
	"context"
	"fmt"
	"log/slog"

	"API-Avito-shop/internal/dto"
	e "API-Avito-shop/internal/errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TransactionRepository interface {
	TransferCoin(ctx context.Context, tx pgx.Tx, fromUser, toUser string, coin int) error
	ReceivedTransaction(ctx context.Context, tx pgx.Tx, username string) ([]dto.ReceivedCoin, error)
	SentTransaction(ctx context.Context, tx pgx.Tx, username string) ([]dto.SentCoin, error)
}

type TransactionRepo struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

func NewTransactionRepository(pool *pgxpool.Pool, logger *slog.Logger) *TransactionRepo {
	return &TransactionRepo{pool: pool, logger: logger}
}

const (
	querySaveTransaction     = `INSERT INTO transactions (from_username, to_username, amount) VALUES ($1, $2, $3)`
	queryReceivedTransaction = `SELECT from_username, amount FROM transactions WHERE to_username = $1`
	querySendTransaction     = `SELECT to_username, amount FROM transactions WHERE from_username = $1`
)

// TransferCoin сохраняет данные транзакции монет
func (r *TransactionRepo) TransferCoin(ctx context.Context, tx pgx.Tx, fromUser, toUser string, amount int) error {
	r.logger.Info("Executing query", "query", querySaveTransaction, "from_user", fromUser, "to_user", toUser)

	_, err := tx.Exec(ctx, querySaveTransaction, fromUser, toUser, amount)
	if err != nil {
		r.logger.Error("Failed to execute query to save coins transaction", "from_user", fromUser, "to_user", toUser, "error", err)
		return fmt.Errorf("TransferCoin: %w", e.ErrFailedExecuteQuery)
	}

	r.logger.Info("Coin transaction saved successfully", "from_user", fromUser, "to_user", toUser)
	return nil
}

// ReceivedTransaction предоставляет список полученных транзакций
func (r *TransactionRepo) ReceivedTransaction(ctx context.Context, tx pgx.Tx, username string) ([]dto.ReceivedCoin, error) {
	var transactions []dto.ReceivedCoin

	r.logger.Info("Executing query", "query", queryReceivedTransaction, "username", username)
	rows, err := tx.Query(ctx, queryReceivedTransaction, username)
	if err != nil {
		r.logger.Error("Failed to execute query to receive the received transactions", "username", username, "error", err)
		return transactions, fmt.Errorf("ReceivedTransaction: %w", e.ErrFailedExecuteQuery)
	}
	defer rows.Close()

	for rows.Next() {
		var transaction dto.ReceivedCoin
		if err = rows.Scan(&transaction.FromUser, &transaction.Amount); err != nil {
			r.logger.Error("Failed to parse row", "error", err)
			return transactions, fmt.Errorf("ReceivedTransaction: failed to parse rows: %w", err)
		}
		transactions = append(transactions, transaction)
	}
	if err = rows.Err(); err != nil {
		r.logger.Error("Error during rows iteration", "error", err)
		return transactions, fmt.Errorf("ReceivedTransaction: error during rows iteration: %w", err)
	}

	r.logger.Info("Transactions received")
	return transactions, nil
}

// SentTransaction предоставляет список отправленных транзакций
func (r *TransactionRepo) SentTransaction(ctx context.Context, tx pgx.Tx, username string) ([]dto.SentCoin, error) {
	var transactions []dto.SentCoin

	r.logger.Info("Executing query", "query", querySendTransaction, "username", username)
	rows, err := tx.Query(ctx, querySendTransaction, username)
	if err != nil {
		r.logger.Error("Failed to execute query to receive the send transactions", "username", username, "error", err)
		return transactions, fmt.Errorf("SendTransaction: %w", e.ErrFailedExecuteQuery)
	}
	defer rows.Close()

	for rows.Next() {
		var transaction dto.SentCoin
		if err = rows.Scan(&transaction.ToUser, &transaction.Amount); err != nil {
			r.logger.Error("Failed to parse row", "error", err)
			return transactions, fmt.Errorf("SendTransaction: failed to parse rows: %w", err)
		}
		transactions = append(transactions, transaction)
	}
	if err = rows.Err(); err != nil {
		r.logger.Error("Error during rows iteration", "error", err)
		return transactions, fmt.Errorf("SendTransaction: error during rows iteration: %w", err)
	}

	r.logger.Info("Transactions received")
	return transactions, nil
}
