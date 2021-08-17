package storage

import (
	"path"

	"github.com/ocramh/challenge2/pkg/content"
)

// NoopStore is an implementation of the Storage interface which doesn't produce any
// result or error
type NoopStore struct {
	roodDir string
}

func NewNoopStore(r string) NoopStore {
	return NoopStore{r}
}

func (n NoopStore) Put(b []byte, name string) (*content.Address, error) {
	cid, err := content.CidFromBytes(b)
	if err != nil {
		return nil, err
	}

	return &content.Address{
		Cid:      cid,
		Path:     path.Join(n.roodDir, name),
		NodeName: name,
	}, nil
}

func (n NoopStore) Get(addr *content.Address) ([]byte, error) {
	return []byte{}, nil
}

func (n NoopStore) Delete(addr *content.Address) error {
	return nil
}
