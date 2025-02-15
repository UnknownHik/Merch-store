package services

import (
	"context"
	"log/slog"

	r "API-Avito-shop/internal/repositories"

	"github.com/jackc/pgx/v5"
)

type ShopService interface {
	BuyProduct(ctx context.Context, item string, username string) error
}

type DefaultShopService struct {
	userRepo   r.UserRepository
	shopRepo   r.ShopRepository
	txExecutor TxExecutor
	logger     *slog.Logger
}

func NewShopService(userRepo r.UserRepository, shopRepo r.ShopRepository, txHelper TxExecutor, logger *slog.Logger) *DefaultShopService {
	return &DefaultShopService{
		userRepo:   userRepo,
		shopRepo:   shopRepo,
		txExecutor: txHelper,
		logger:     logger,
	}
}

// BuyProduct позволяет купить товар
func (s *DefaultShopService) BuyProduct(ctx context.Context, item string, username string) error {
	s.logger.Info("Starting to buy item", "item", item)

	existingItem, err := s.shopRepo.GetItem(ctx, item)
	if err != nil {
		s.logger.Error("Failed to find item", "item", item, "error", err)
		return err
	}

	err = s.txExecutor.RunWithTransaction(ctx, func(tx pgx.Tx) error {
		if err = s.userRepo.SubtractCoins(ctx, tx, username, existingItem.Price); err != nil {
			return err
		}
		s.logger.Info("Payment for item made", "username", username, "item", item)

		if err = s.shopRepo.AddPurchase(ctx, tx, item, username, existingItem.Price); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		s.logger.Error("Failed to buy item", "item", item, "error", err)
		return err
	}

	s.logger.Info("Purchase completed successfully", "username", username, "item", item)
	return nil
}
