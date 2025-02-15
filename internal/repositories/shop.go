package repositories

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"API-Avito-shop/internal/dto"
	e "API-Avito-shop/internal/errors"
	"API-Avito-shop/internal/models"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ShopRepository interface {
	GetItem(ctx context.Context, item string) (*models.Product, error)
	AddPurchase(ctx context.Context, tx pgx.Tx, item, username string, price int) error
	GetPurchases(ctx context.Context, tx pgx.Tx, username string) ([]dto.Item, error)
}

type ShopRepo struct {
	pool   *pgxpool.Pool
	logger *slog.Logger
}

func NewShopRepository(pool *pgxpool.Pool, logger *slog.Logger) *ShopRepo {
	return &ShopRepo{pool: pool, logger: logger}
}

const (
	queryGetItem      = `SELECT item, price FROM products WHERE item = $1 FOR UPDATE`
	queryAddPurchase  = `INSERT INTO purchases (username, item, price) VALUES ($1, $2, $3)`
	queryGetPurchases = `SELECT item, COUNT(*) AS total_purchased FROM purchases WHERE username = $1 GROUP BY item`
)

// GetItem получение товара по названию из доступных к приобретению
func (r *ShopRepo) GetItem(ctx context.Context, item string) (*models.Product, error) {
	var product models.Product

	r.logger.Info("Executing query", "query", queryGetItem, "item", item)
	err := r.pool.QueryRow(ctx, queryGetItem, item).Scan(&product.Item, &product.Price)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			r.logger.Info("Product not found", "item", item)
			return nil, err
		}

		r.logger.Error("Failed to execute query to get item", "item", item, "error", err)
		return nil, fmt.Errorf("GetItem: %w", err)
	}

	r.logger.Info("Product found", "item", item)
	return &product, nil
}

// AddPurchase добавление совершенной покупки
func (r *ShopRepo) AddPurchase(ctx context.Context, tx pgx.Tx, item, username string, price int) error {
	r.logger.Info("Executing query", "query", queryAddPurchase, "item", item)

	_, err := tx.Exec(ctx, queryAddPurchase, username, item, price)
	if err != nil {
		r.logger.Error("Failed to execute query to add purchase", "username", username, "item", item, "error", err)
		return fmt.Errorf("AddPurchase: %w", e.ErrFailedExecuteQuery)
	}

	r.logger.Info("Purchase added", "username", username, "item", item)
	return nil
}

// GetPurchases предоставляет список купленных товаров
func (r *ShopRepo) GetPurchases(ctx context.Context, tx pgx.Tx, username string) ([]dto.Item, error) {
	var purchases []dto.Item

	r.logger.Info("Executing query", "query", queryGetPurchases, "username", username)
	rows, err := tx.Query(ctx, queryGetPurchases, username)
	if err != nil {
		r.logger.Error("Failed to execute query to get purchases", "username", username, "error", err)
		return purchases, fmt.Errorf("GetPurchases: %w", e.ErrFailedExecuteQuery)
	}
	defer rows.Close()

	for rows.Next() {
		var purchase dto.Item
		if err = rows.Scan(&purchase.Type, &purchase.Quantity); err != nil {
			r.logger.Error("Failed to parse row", "error", err)
			return purchases, fmt.Errorf("GetPurchases: failed to parse rows: %w", err)
		}
		purchases = append(purchases, purchase)
	}
	if err = rows.Err(); err != nil {
		r.logger.Error("Error during rows iteration", "error", err)
		return purchases, fmt.Errorf("GetPurchases: error during rows iteration: %w", err)
	}

	r.logger.Info("Purchase list received")
	return purchases, nil
}
