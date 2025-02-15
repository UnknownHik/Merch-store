package services

import (
	"context"
	"log/slog"

	"API-Avito-shop/internal/dto"
	e "API-Avito-shop/internal/errors"
	r "API-Avito-shop/internal/repositories"

	"github.com/jackc/pgx/v5"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	AuthUser(ctx context.Context, userAuthDTO *dto.UserAuth) error
	UserInfo(ctx context.Context, username string) (dto.InfoResponse, error)
}

type DefaultUserService struct {
	userRepo        r.UserRepository
	shopRepo        r.ShopRepository
	transactionRepo r.TransactionRepository
	txExecutor      TxExecutor
	logger          *slog.Logger
}

func NewUserService(userRepo r.UserRepository, shopRepo r.ShopRepository, transactionRepo r.TransactionRepository, txHelper TxExecutor, logger *slog.Logger) *DefaultUserService {
	return &DefaultUserService{
		userRepo:        userRepo,
		shopRepo:        shopRepo,
		transactionRepo: transactionRepo,
		txExecutor:      txHelper,
		logger:          logger,
	}
}

// AuthUser выполняет авторизацию пользователя или создает его, если первый раз
func (s *DefaultUserService) AuthUser(ctx context.Context, userAuthDTO *dto.UserAuth) error {
	s.logger.Info("Start of user authorization", "username", userAuthDTO.UserName)

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(userAuthDTO.Password), bcrypt.DefaultCost)
	if err != nil {
		s.logger.Error("Failed to hash password", "error", err)
		return err
	}

	storedPassword, err := s.userRepo.GetOrCreateUser(ctx, userAuthDTO.UserName, string(hashedPassword))
	if err != nil {
		s.logger.Error("Failed to get or create user", "error", err)
		return err
	}

	if err = bcrypt.CompareHashAndPassword([]byte(storedPassword), []byte(userAuthDTO.Password)); err != nil {
		s.logger.Warn("Incorrect password", "username", userAuthDTO.UserName)
		return e.ErrInvalidPass
	}

	s.logger.Info("User successfully created")
	return nil
}

// UserInfo предоставляет информацию о пользователе: текущий баланс, приобретенные товары и историю транзакций
func (s *DefaultUserService) UserInfo(ctx context.Context, username string) (dto.InfoResponse, error) {
	s.logger.Info("Starting to get information about user", "username", username)

	var userData dto.InfoResponse

	err := s.txExecutor.RunWithTransaction(ctx, func(tx pgx.Tx) (err error) {
		balance, err := s.userRepo.GetBalance(ctx, tx, username)
		if err != nil {
			return err
		}
		userData.Coins = balance

		s.logger.Info("Received user balance", "balance", userData.Coins)

		userData.Inventory, err = s.shopRepo.GetPurchases(ctx, tx, username)
		if err != nil {
			return err
		}

		userData.CoinHistory.Received, err = s.transactionRepo.ReceivedTransaction(ctx, tx, username)
		if err != nil {
			return err
		}

		userData.CoinHistory.Sent, err = s.transactionRepo.SentTransaction(ctx, tx, username)
		if err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		s.logger.Error("Failed to get information about user", "username", username, "error", err)
		return userData, err
	}

	s.logger.Info("User information retrieved successfully", "username", username)
	return userData, nil
}
