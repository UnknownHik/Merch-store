package delivery

import (
	"errors"
	"net/http"

	"API-Avito-shop/internal/dto"
	e "API-Avito-shop/internal/errors"
	s "API-Avito-shop/internal/services"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	userService s.UserService
	token       s.Token
}

func NewUserHandler(userService s.UserService, token s.Token) *UserHandler {
	return &UserHandler{
		userService: userService,
		token:       token,
	}
}

// AuthHandler обрабатывает запрос на аутентификацию пользователя
func (h *UserHandler) AuthHandler(c *gin.Context) {
	var userAuthDTO dto.UserAuth

	if err := c.ShouldBindJSON(&userAuthDTO); err != nil {
		handleError(c, http.StatusBadRequest, "Invalid username or password format", err)
		return
	}

	err := h.userService.AuthUser(c.Request.Context(), &userAuthDTO)
	if err != nil {
		if errors.Is(err, e.ErrInvalidPass) {
			handleError(c, http.StatusUnauthorized, "Authorization failed", err)
			return
		}
		handleError(c, http.StatusInternalServerError, "Authentication service error", err)
		return
	}

	token, err := h.token.GenerateToken(c.Request.Context(), userAuthDTO.UserName)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to generate token", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

// InfoHandler обрабатывает запрос на получение информации о балансе и действиях пользователя
func (h *UserHandler) InfoHandler(c *gin.Context) {
	username, err := getUsername(c)
	if err != nil {
		handleError(c, http.StatusBadRequest, err.Error(), err)
		return
	}

	userInfo, err := h.userService.UserInfo(c.Request.Context(), username)
	if err != nil {
		handleError(c, http.StatusInternalServerError, "Failed to get user information", err)
		return
	}

	response := dto.InfoResponse{
		Coins:       userInfo.Coins,
		Inventory:   userInfo.Inventory,
		CoinHistory: userInfo.CoinHistory,
	}

	c.JSON(http.StatusOK, response)
}
