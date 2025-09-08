package domain

import "errors"

var ErrOrderNotFound = errors.New("order not found")
var ErrOrderAlreadyExists = errors.New("order already exists")
