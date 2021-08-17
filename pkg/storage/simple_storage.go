package storage

import (
	"path"

	"github.com/ocramh/challenge2/pkg/content"
)

// SimpleStore is an implementation of the Storage interface that keeps track of
// content in memeory
type SimpleStore struct {
	roodDir     string
	maxCapacity int
	store       map[string][]byte
}

func NewSimpleStore(root string, maxCapacity int) *SimpleStore {
	return &SimpleStore{
		roodDir:     root,
		maxCapacity: maxCapacity,
		store:       make(map[string][]byte),
	}
}

func (n *SimpleStore) Put(b []byte, name string) (*content.Address, error) {
	cid, err := content.CidFromBytes(b)
	if err != nil {
		return nil, err
	}

	if n.Size() == n.Capacity() {
		return nil, ErrNoStorageAvailable
	}

	n.store[cid.String()] = b

	return &content.Address{
		Cid:      cid,
		Path:     path.Join(n.roodDir, name),
		NodeName: name,
	}, nil
}

func (n *SimpleStore) Get(addr *content.Address) ([]byte, error) {
	return n.store[addr.Cid.String()], nil
}

func (n *SimpleStore) Delete(addr *content.Address) error {
	delete(n.store, addr.Cid.String())
	return nil
}

func (n *SimpleStore) Size() int {
	return len(n.store)
}

func (n *SimpleStore) Capacity() int {
	return n.maxCapacity
}
