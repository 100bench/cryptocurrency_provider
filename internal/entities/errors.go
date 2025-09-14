package entities

import "errors"

var (
	ErrNilDependency = errors.New("nil dependency: ")
	//ErrProviderUnavailable = errors.New("provider unavailable: ")
	ErrInvalidParams = errors.New("invalid parameters: ")
)
