package entities

import "errors"

var ErrNilDependency = errors.New("nil dependency: ")
var ErrEmptyBaseURL = errors.New("empty: ")
