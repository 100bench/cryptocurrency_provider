package repository

import "errors"

var (
	ErrDbNil      = errors.New("db is nil")
	ErrRedisIsNil = errors.New("redis is nil")
)
