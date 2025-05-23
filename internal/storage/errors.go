package storage

import "errors"

var (
	ErrNotFound      = errors.New("not Found")
	ErrInvalidData   = errors.New("invalid data")
	ErrAlreadyExists = errors.New("already exists")
)
