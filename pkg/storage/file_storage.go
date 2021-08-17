package storage

import (
	"io/ioutil"
	"os"
	"path"

	"github.com/ocramh/challenge2/pkg/content"
)

// FileStore is an implementation of the Storage interface backed by file objects
type FileStore struct{}

func (fs FileStore) Put(b []byte, fpath string) (*content.Address, error) {
	f, err := os.Create(fpath)
	if err != nil {
		return nil, err
	}

	_, err = f.Write(b)
	if err != nil {
		return nil, err
	}

	cid, err := content.CidFromBytes(b)
	if err != nil {
		return nil, err
	}

	return &content.Address{
		Cid:      cid,
		Path:     fpath,
		NodeName: path.Base(fpath),
	}, nil
}

func (fs FileStore) Get(addr *content.Address) ([]byte, error) {
	f, err := os.Open(addr.Path)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(f)
}

func (fs FileStore) Delete(addr *content.Address) error {
	return os.RemoveAll(addr.Path)
}
