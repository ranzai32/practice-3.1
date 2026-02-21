package apperrors

import "errors"

var (
	ErrNotFound = errors.New("404")
	ErrConflict = errors.New("already exists")
	ErrInternal = errors.New("500")
	ErrValidation = errors.New("400")
)