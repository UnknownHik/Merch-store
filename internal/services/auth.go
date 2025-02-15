package services

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/golang-jwt/jwt/v5"
)

type Token interface {
	GenerateToken(ctx context.Context, userID string) (string, error)
	ValidateToken(ctx context.Context, token string) (string, error)
}
type DefaultToken struct {
	secretKey string
	logger    *slog.Logger
}

func NewToken(secretKey string, logger *slog.Logger) *DefaultToken {
	return &DefaultToken{
		secretKey: secretKey,
		logger:    logger,
	}
}

// GenerateToken генерирует токен
func (t *DefaultToken) GenerateToken(ctx context.Context, username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
	})

	tokenString, err := token.SignedString([]byte(t.secretKey))
	if err != nil {
		t.logger.Error("Failed to generate token", "username", username, "error", err)
		return "", err
	}

	t.logger.Info("Token successfully generated", "username", username)
	return tokenString, nil
}

// ValidateToken проверяет валидность токена
func (t *DefaultToken) ValidateToken(ctx context.Context, tokenString string) (string, error) {
	parsedToken, err := jwt.Parse(tokenString, func(j *jwt.Token) (interface{}, error) {
		if _, ok := j.Method.(*jwt.SigningMethodHMAC); !ok {
			err := fmt.Errorf("unexpected signing method: %v", j.Header["alg"])
			t.logger.Error("Invalid token signing method", "token", tokenString, "error", err)
			return nil, err
		}
		return []byte(t.secretKey), nil
	})
	if err != nil || !parsedToken.Valid {
		t.logger.Error("Invalid token", "token", tokenString, "error", err)
		return "", errors.New("invalid token")
	}

	claims, ok := parsedToken.Claims.(jwt.MapClaims)
	if !ok {
		t.logger.Error("Invalid token claims", "token", tokenString)
		return "", errors.New("invalid token claims")
	}

	username, ok := claims["username"].(string)
	if !ok {
		t.logger.Error("Invalid user ID in token claims", "token", tokenString)
		return "", errors.New("invalid user ID in token")
	}

	t.logger.Info("Token validated successfully", "username", username)
	return username, nil
}
