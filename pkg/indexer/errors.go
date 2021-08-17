package indexer

import "errors"

var (
	ErrNoItemFound        = errors.New("no item found")
	ErrNoStorageAvailable = errors.New("insufficient storage")
)
