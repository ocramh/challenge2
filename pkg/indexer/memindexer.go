package indexer

import (
	"bytes"
	"io"
	"sync"

	"github.com/ipfs/go-cid"

	"github.com/ocramh/challenge2/pkg/content"
	"github.com/ocramh/challenge2/pkg/storage"
)

// KVStore is the key - value store used to track the available blocks of content
type KVStore map[string]*content.Block

// MemIndex is an implementation of the Indexer interface which keeps track of the
// available storage using an in-memory key value map.
// Internally it has access to the underlying storage implementation for adding,
// accessing, removing and syncing the actual data
type MemIndex struct {
	mu            sync.Mutex
	kvStore       KVStore
	evictStrategy Evictor
	storage       storage.StorageManager
}

func NewMemoryIndex(ev Evictor, store storage.StorageManager) *MemIndex {
	return &MemIndex{
		kvStore:       make(map[string]*content.Block),
		evictStrategy: ev,
		storage:       store,
	}
}

func (m *MemIndex) Put(src io.Reader, name string) (*content.Block, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	buf := new(bytes.Buffer)
	buf.ReadFrom(src)
	srcAsBytes := buf.Bytes()

	if m.storage.Size() >= m.storage.Capacity() {
		resizeErr := m.resizeStore()
		if resizeErr != nil {
			return nil, resizeErr
		}
	}

	// add new content to storage
	addr, err := m.storage.Put(srcAsBytes, name)

	block, err := content.NewBlock(srcAsBytes, addr)
	if err != nil {
		return nil, err
	}

	// add block to the indexer key value store
	m.kvStore[block.ID()] = block

	return block, nil
}

func (m *MemIndex) resizeStore() error {
	for m.storage.Size() >= m.storage.Capacity() {
		err := m.evictBlock()
		if err != nil {
			return err
		}
	}

	return nil
}

func (m *MemIndex) evictBlock() error {
	rmAddr := m.evictStrategy.EvictBlock(m.kvStore)
	if rmAddr != nil {
		return m.storage.Delete(rmAddr)
	}

	return ErrNoItemFound
}

func (m *MemIndex) Get(blockID cid.Cid) (*content.Block, []byte, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	block, ok := m.kvStore[blockID.String()]
	if !ok {
		return nil, nil, ErrNoItemFound
	}

	d, err := m.storage.Get(&block.Address)
	if err != nil {
		return nil, nil, err
	}
	block.IncHitsCount()

	return block, d, nil
}
