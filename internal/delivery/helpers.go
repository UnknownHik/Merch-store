package delivery

import (
	"fmt"
	"log/slog"

	"github.com/gin-gonic/gin"
)

// getUsername извлекает и валидирует username из контекста
func getUsername(c *gin.Context) (string, error) {
	userIDCtx, exists := c.Get("username")
	if !exists || userIDCtx == nil {
		return "", fmt.Errorf("username missing in context")
	}

	username, ok := userIDCtx.(string)
	if !ok {
		return "", fmt.Errorf("invalid username type in context")
	}

	return username, nil
}

// handleError отправляет HTTP-ответ c ошибкой
func handleError(c *gin.Context, status int, message string, err error) {
	if err != nil {
		slog.Error(message, "method", c.Request.Method, "path", c.Request.URL.Path, "client_ip", c.ClientIP(), "error", err)
	}
	c.JSON(status, gin.H{"error": message})
}
