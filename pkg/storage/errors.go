package storage

import "errors"

var (
	ErrNoStorageAvailable = errors.New("insufficient storage")
)
