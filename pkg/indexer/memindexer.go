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
	rootDir       string
	maxCapacity   int
	mu            sync.Mutex
	kvStore       KVStore
	evictStrategy Evictor
	storage       storage.StorageManager
}

func NewMemoryIndex(rootDir string, cap int, ev Evictor, store storage.StorageManager) *MemIndex {
	return &MemIndex{
		rootDir:       rootDir,
		maxCapacity:   cap,
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

	if len(m.kvStore) == m.maxCapacity {
		err := m.resizeStore()
		if err != nil {
			return nil, err
		}
	}

	// add new content to storage
	addr, err := m.storage.Put(srcAsBytes, name)
	if err != nil {
		return nil, err
	}

	block, err := content.NewBlock(srcAsBytes, addr)
	if err != nil {
		return nil, err
	}

	// add block to the indexer key value store
	m.kvStore[block.ID()] = block

	return block, nil
}

func (m *MemIndex) resizeStore() error {
	for len(m.kvStore) >= m.maxCapacity {
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

	return nil
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

func (m *MemIndex) Size() int {
	return len(m.kvStore)
}

func (m *MemIndex) Capacity() int {
	return m.maxCapacity
}
