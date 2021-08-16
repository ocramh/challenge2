package indexer

import (
	"github.com/ocramh/challenge2/pkg/content"
)

type Evictor interface {
	EvictBlock(store KVStore) *content.Address
}

// EvictLeastPopular is a naive implementation of the evictor interface which simply
// removes the block of content from KVStore with the lowest number of hits
type EvictLeastPopular struct{}

func (e EvictLeastPopular) EvictBlock(store KVStore) *content.Address {
	if len(store) == 0 {
		return nil
	}

	c := 0
	var toEvict content.BlockKey
	var evictedAddr *content.Address
	var lowestHitsCount = 0
	for k, v := range store {
		if c == 0 {
			toEvict = k
			lowestHitsCount = v.GetHitsCount()
			evictedAddr = &v.Address
			c++
			continue
		}

		if lowestHitsCount > v.GetHitsCount() {
			toEvict = k
			lowestHitsCount = v.GetHitsCount()
			evictedAddr = &v.Address
		}

		c++
	}

	delete(store, toEvict)

	return evictedAddr
}
