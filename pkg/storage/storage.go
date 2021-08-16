package storage

import (
	"github.com/ocramh/challenge2/pkg/content"
)

// StorageManager is the interface used for persisting data
type StorageManager interface {
	// Put adds a block to storage and returns its address
	Put(b []byte, path string) (*content.Address, error)

	// Get retrives a block of content identified by its cid from storage.
	// An error will be returned if the provided address cannot be resolved
	Get(c *content.Address) ([]byte, error)

	// Delete removes an item from storage at the specified address
	Delete(c *content.Address) error
}
