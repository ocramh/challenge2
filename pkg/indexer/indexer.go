package indexer

import (
	"io"

	"github.com/ocramh/challenge2/pkg/content"
)

// Indexer defines the functionalies exposed by a storage indexer, used for fast
// access of content made available by the provider
type Indexer interface {
	// Put adds content to storage and returns its block representation
	Put(src io.Reader, name string) (*content.Block, error)

	// Get retrives a block of content from storage identified byt its cid.
	// An error will be returned if the provided cid doesn't match any available block
	Get(content.BlockKey) (*content.Block, error)

	// Size returns the current storage usage
	Size() int

	// Capacity returns the total storage capacity
	Capacity() int
}
