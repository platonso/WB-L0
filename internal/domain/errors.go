package domain

import "errors"

var (
	ErrOrderNotFound      = errors.New("order not found")
	ErrOrderAlreadyExists = errors.New("order already exists")
	ErrValidation         = errors.New("validation error")
)
