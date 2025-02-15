package services

import (
	"context"
	"log/slog"

	"API-Avito-shop/internal/dto"
	e "API-Avito-shop/internal/errors"
	r "API-Avito-shop/internal/repositories"

	"github.com/jackc/pgx/v5"
)

type TransactionService interface {
	SendCoin(ctx context.Context, username string, sendCoinDTO *dto.SendCoin) error
}

type DefaultTransactionService struct {
	userRepo        r.UserRepository
	transactionRepo r.TransactionRepository
	txExecutor      TxExecutor
	logger          *slog.Logger
}

func NewTransactionService(userRepo r.UserRepository, transactionRepo r.TransactionRepository, txHelper TxExecutor, logger *slog.Logger) *DefaultTransactionService {
	return &DefaultTransactionService{
		userRepo:        userRepo,
		transactionRepo: transactionRepo,
		txExecutor:      txHelper,
		logger:          logger,
	}
}

// SendCoin производит транзакцию монет между пользователями
func (s *DefaultTransactionService) SendCoin(ctx context.Context, username string, sendCoinDTO *dto.SendCoin) error {
	s.logger.Info("Starting to send coins from user to user", "from_user", username, "to_user", sendCoinDTO.ToUser)

	err := s.txExecutor.RunWithTransaction(ctx, func(tx pgx.Tx) error {
		if sendCoinDTO.ToUser == username {
			s.logger.Warn("Sender matches recipient", "from_user", username, "to_user", sendCoinDTO.ToUser)
			return e.ErrInvalidUser
		}

		if err := s.userRepo.SubtractCoins(ctx, tx, username, sendCoinDTO.Amount); err != nil {
			return err
		}
		s.logger.Info("Coins left the user", "from_user", username)

		if err := s.userRepo.AddCoins(ctx, tx, sendCoinDTO.ToUser, sendCoinDTO.Amount); err != nil {
			return err
		}
		s.logger.Info("Coins reached the user", "from_user", username)

		if err := s.transactionRepo.TransferCoin(ctx, tx, username, sendCoinDTO.ToUser, sendCoinDTO.Amount); err != nil {
			return err
		}

		return nil
	})

	if err != nil {
		s.logger.Error("Failed to send coins", "from_user", username, "to_user", sendCoinDTO.ToUser, "error", err)
		return err
	}

	s.logger.Info("Coins sent successfully", "from_user", username, "to_user", sendCoinDTO.ToUser)
	return nil
}
