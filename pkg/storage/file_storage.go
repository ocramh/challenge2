package storage

import (
	"io/ioutil"
	"os"

	"github.com/ocramh/challenge2/pkg/content"
)

// FileStore is an implementation of the Storage interface backed by file objects
type FileStore struct{}

func (fs FileStore) Put(b []byte, path string) (*content.Address, error) {
	f, err := os.Create(path)
	if err != nil {
		return nil, err
	}

	_, err = f.Write(b)
	if err != nil {
		return nil, err
	}

	return &content.Address{
		Filepath: path,
	}, nil
}

func (fs FileStore) Get(addr *content.Address) ([]byte, error) {
	f, err := os.Open(addr.Filepath)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(f)
}

func (fs FileStore) Delete(addr *content.Address) error {
	return os.RemoveAll(addr.Filepath)
}
