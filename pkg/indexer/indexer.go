package indexer

import (
	"io"

	"github.com/ipfs/go-cid"
	"github.com/ocramh/challenge2/pkg/content"
)

// Indexer defines the functionalies for fast access to the content made available by
// the provider
type Indexer interface {
	// Put adds content to storage and returns its block representation
	Put(src io.Reader, name string) (*content.Block, error)

	// Get retrives a block of content from storage identified byt its cid.
	// An error will be returned if the provided cid doesn't match any available block
	Get(cid.Cid) (*content.Block, []byte, error)
}
