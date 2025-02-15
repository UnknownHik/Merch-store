package errors

import "errors"

var (
	ErrInvalidPass        = errors.New("invalid password")
	ErrFailedExecuteQuery = errors.New("failed to execute query")
	ErrNotEnoughCoins     = errors.New("Not enough coins")
	ErrInvalidUser        = errors.New("invalid user")
)
