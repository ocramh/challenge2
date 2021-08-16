package storage

import (
	"github.com/ocramh/challenge2/pkg/content"
)

// NoopStore is an implementation of the Storage interface which doesn't produce any
// result or error
type NoopStore struct{}

func (n NoopStore) Put(b []byte, path string) (*content.Address, error) {
	return &content.Address{
		Filepath: path,
	}, nil
}

func (n NoopStore) Get(addr *content.Address) ([]byte, error) {
	return []byte{}, nil
}

func (n NoopStore) Delete(addr *content.Address) error {
	return nil
}
