package delivery

import (
	"errors"
	"net/http"

	e "API-Avito-shop/internal/errors"
	s "API-Avito-shop/internal/services"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5"
)

type ShopHandler struct {
	shopService s.ShopService
}

func NewShopHandler(shopService s.ShopService) *ShopHandler {
	return &ShopHandler{
		shopService: shopService,
	}
}

// BuyHandler обрабатывает запрос на покупку товара
func (h *ShopHandler) BuyHandler(c *gin.Context) {
	username, err := getUsername(c)
	if err != nil {
		handleError(c, http.StatusBadRequest, "Failed to get user_id from context", err)
		return
	}

	item := c.Param("item")
	if item == "" {
		err = errors.New("item parameter is required")
		handleError(c, http.StatusBadRequest, "Item parameter is missing", err)
		return
	}

	err = h.shopService.BuyProduct(c.Request.Context(), item, username)
	if err != nil {
		switch {
		case errors.Is(err, pgx.ErrNoRows):
			handleError(c, http.StatusBadRequest, "Invalid item", err)
		case errors.Is(err, e.ErrNotEnoughCoins):
			handleError(c, http.StatusBadRequest, "Insufficient balance", err)
		default:
			handleError(c, http.StatusInternalServerError, "Failed to complete purchase", err)
		}
		return
	}

	c.Status(http.StatusOK)
}
