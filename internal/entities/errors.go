package entities

import "errors"

var (
	ErrNilDependency = errors.New("nil dependency: ")
	ErrInvalidParams = errors.New("invalid parameters: ")
)
