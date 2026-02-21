package apperrors

import "errors"

var (
	ErrNotFound = errors.New("not found")
	ErrConflict = errors.New("already exists")
	ErrInternal = errors.New("internal error")
	ErrValidation = errors.New("name and email are required")
)