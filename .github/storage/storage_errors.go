package storage

import "errors"

var (
	ErrOperationNotFound = errors.New("operation not found")
	ErrUserNotFound      = errors.New("user not found")
)
