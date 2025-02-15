package middleware

import (
	"log/slog"
	"net/http"
	"strings"

	"API-Avito-shop/internal/services"

	"github.com/gin-gonic/gin"
)

type AuthMiddleware struct {
	token     services.Token
	secretKey []byte
	logger    *slog.Logger
}

func NewAuthMiddleware(token services.Token, secretKey string, logger *slog.Logger) *AuthMiddleware {
	return &AuthMiddleware{
		token:     token,
		secretKey: []byte(secretKey),
		logger:    logger,
	}
}

func (m *AuthMiddleware) AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" || !strings.HasPrefix(authHeader, "Bearer ") {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "missing or invalid authorization header"})
			c.Abort()
			return
		}

		token := strings.TrimPrefix(authHeader, "Bearer ")

		username, err := m.token.ValidateToken(c.Request.Context(), token)
		if err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid token"})
			c.Abort()
			return
		}

		c.Set("username", username)
		c.Next()
	}
}
