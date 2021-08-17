package storage

import (
	"testing"

	"github.com/ipfs/go-cid"
	"github.com/stretchr/testify/assert"
)

func TestPutBlocks(t *testing.T) {
	store := NewBlocksStore(2)

	b1, err := store.Put([]byte("block1"), "block 1")
	assert.NoError(t, err)
	assert.False(t, b1.Cid.Equals(cid.Undef))
	assert.Equal(t, b1.Cid.String(), b1.Path)
	assert.Equal(t, "block 1", b1.NodeName)
	assert.Equal(t, store.Size(), 1)

	_, err = store.Put([]byte("block2"), "block 2")
	assert.NoError(t, err)

	_, err = store.Put([]byte("block3"), "block 3")
	assert.EqualError(t, err, ErrNoStorageAvailable.Error())
}

func TestGetBlocks(t *testing.T) {
	store := NewBlocksStore(2)

	b1, err := store.Put([]byte("block1"), "block 1")
	assert.NoError(t, err)

	got, err := store.Get(b1)
	assert.NoError(t, err)

	assert.Subset(t, got, []byte("block1"))
}

func TestDeleteBlocks(t *testing.T) {
	store := NewBlocksStore(2)

	addr, err := store.Put([]byte("block1"), "block 1")
	assert.NoError(t, err)

	err = store.Delete(addr)
	assert.NoError(t, err)

	assert.Equal(t, store.Size(), 0)
}
