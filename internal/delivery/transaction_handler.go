package delivery

import (
	"errors"
	"net/http"

	"API-Avito-shop/internal/dto"
	e "API-Avito-shop/internal/errors"
	s "API-Avito-shop/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type TransactionHandler struct {
	transactionService s.TransactionService
}

func NewTransactionHandler(transactionService s.TransactionService) *TransactionHandler {
	return &TransactionHandler{
		transactionService: transactionService,
	}
}

// SendCoinHandler обрабатывает запрос на отправку монет пользователю
func (h *TransactionHandler) SendCoinHandler(c *gin.Context) {
	username, err := getUsername(c)
	if err != nil {
		handleError(c, http.StatusBadRequest, "Failed to get user_id from context", err)
		return
	}

	var sendCoinDTO dto.SendCoin

	if err = c.ShouldBindJSON(&sendCoinDTO); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid request data", err)
		return
	}

	err = h.transactionService.SendCoin(c.Request.Context(), username, &sendCoinDTO)
	if err != nil {
		switch {
		case errors.Is(err, e.ErrNotEnoughCoins):
			handleError(c, http.StatusBadRequest, "Insufficient balance", err)
		case errors.Is(err, pgx.ErrNoRows) || errors.Is(err, e.ErrInvalidUser):
			handleError(c, http.StatusBadRequest, "Invalid recipient", err)
		default:
			handleError(c, http.StatusInternalServerError, "Failed to send coins", err)
		}
		return
	}

	c.Status(http.StatusOK)
}
