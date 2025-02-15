package repositories

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	e "API-Avito-shop/internal/errors"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	GetOrCreateUser(ctx context.Context, username, password string) (string, error)
	GetBalance(ctx context.Context, tx pgx.Tx, username string) (int, error)
	SubtractCoins(ctx context.Context, tx pgx.Tx, username string, coins int) error
	AddCoins(ctx context.Context, tx pgx.Tx, username string, coins int) error
}

type UserRepo struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

func NewUserRepository(pool *pgxpool.Pool, logger *slog.Logger) *UserRepo {
	return &UserRepo{pool: pool, logger: logger}
}

const (
	queryCheckUser      = `SELECT password FROM users WHERE username = $1`
	queryCreateUser     = `INSERT INTO users (username, password) VALUES ($1, $2) ON CONFLICT (username) DO NOTHING RETURNING password;`
	queryGetBalanceByID = `SELECT balance FROM users WHERE username = $1`
	querySubtractCoins  = `UPDATE users SET balance = balance - $1 WHERE username = $2 AND balance >= $1 RETURNING balance`
	queryAddCoins       = `UPDATE users SET balance = balance + $1 WHERE username = $2`
)

// GetOrCreateUser находит пользователя по имени или создает нового
func (r *UserRepo) GetOrCreateUser(ctx context.Context, username, password string) (string, error) {
	var hashedPassword string

	err := r.pool.QueryRow(ctx, queryCreateUser, username, password).Scan(&hashedPassword)
	if err == nil {
		r.logger.Info("User created", "username", username)
		return hashedPassword, nil
	}
	if !errors.Is(err, pgx.ErrNoRows) {
		r.logger.Error("Failed to execute query create user", "username", username, "error", err)
		return "", err
	}

	err = r.pool.QueryRow(ctx, queryCheckUser, username).Scan(&hashedPassword)
	if err != nil {
		r.logger.Error("Failed to execute query to check user", "error", err)
		return "", e.ErrFailedExecuteQuery
	}

	r.logger.Info("User found", "username", username)
	return hashedPassword, nil
}

// GetBalance получение баланса пользователя по его id
func (r *UserRepo) GetBalance(ctx context.Context, tx pgx.Tx, username string) (int, error) {
	var balance int

	r.logger.Info("Executing query", "query", queryGetBalanceByID, "username", username)
	err := tx.QueryRow(ctx, queryGetBalanceByID, username).Scan(&balance)
	if err != nil {
		r.logger.Error("Failed to execute query to get user balance", "username", username, "error", err)
		return 0, fmt.Errorf("GetBalance: %w", e.ErrFailedExecuteQuery)
	}

	r.logger.Info("User balance found", "username", username)
	return balance, nil
}

// SubtractCoins изменение баланса после покупки или отправки транзакции
func (r *UserRepo) SubtractCoins(ctx context.Context, tx pgx.Tx, username string, coins int) error {
	r.logger.Info("Executing query", "query", querySubtractCoins, "username", username)

	var newBalance int
	err := tx.QueryRow(ctx, querySubtractCoins, coins, username).Scan(&newBalance)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Error("Not enough coins", "username", username, "error", err)
			return e.ErrNotEnoughCoins
		}
		r.logger.Error("Failed to subtract coins from balance", "username", username, "error", err)
		return fmt.Errorf("SubtractCoins: %w", e.ErrFailedExecuteQuery)
	}
	r.logger.Info("Balance updated", "username", username)
	return nil
}

// AddCoins обновление баланса после получения транзакции
func (r *UserRepo) AddCoins(ctx context.Context, tx pgx.Tx, username string, coins int) error {
	r.logger.Info("Executing query", "query", queryAddCoins, "username", username)

	_, err := tx.Exec(ctx, queryAddCoins, coins, username)
	if err != nil {
		r.logger.Error("Failed to add coins from balance", "username", username, "error", err)
		return fmt.Errorf("AddCoins: %w", e.ErrFailedExecuteQuery)
	}

	r.logger.Info("Balance updated", "username", username)
	return nil
}
