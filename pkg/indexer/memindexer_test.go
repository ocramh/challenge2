package indexer

import (
	"bytes"
	"io"
	"testing"

	"github.com/ocramh/challenge2/pkg/content"
	"github.com/ocramh/challenge2/pkg/storage"
	"github.com/stretchr/testify/assert"
)

var (
	putContent1 = []byte("content1")
	putContent2 = []byte("content2")
	putContent3 = []byte("content3")
)

func TestPut(t *testing.T) {
	idx := NewMemoryIndex("/root", 2, EvictLeastPopular{}, storage.NoopStore{})

	testcases := []struct {
		src              io.Reader
		fPath            string
		expectedAddr     content.Address
		expectedStoreLen int
	}{
		{
			src:              bytes.NewReader(putContent1),
			fPath:            "/addr1",
			expectedAddr:     content.Address{Filepath: "/root/addr1"},
			expectedStoreLen: 1,
		},
		{
			src:              bytes.NewReader(putContent1),
			fPath:            "/another/addr",
			expectedAddr:     content.Address{Filepath: "/root/addr1"},
			expectedStoreLen: 1,
		},
		{
			src:              bytes.NewReader(putContent2),
			fPath:            "/addr2",
			expectedAddr:     content.Address{Filepath: "/root/addr2"},
			expectedStoreLen: 2,
		},
		{
			src:              bytes.NewReader(putContent3),
			fPath:            "/addr3",
			expectedAddr:     content.Address{Filepath: "/root/addr3"},
			expectedStoreLen: 2,
		},
	}

	for _, testcase := range testcases {
		got, err := idx.Put(testcase.src, testcase.fPath)
		assert.NoError(t, err)
		assert.Equal(t, got.Address, testcase.expectedAddr)
		assert.Len(t, idx.kvStore, testcase.expectedStoreLen)
	}
}

func TestGet(t *testing.T) {
	idx := NewMemoryIndex("/root", 2, EvictLeastPopular{}, storage.NoopStore{})

	c1, err := idx.Put(bytes.NewReader(putContent1), "/addr1")
	assert.NoError(t, err)

	c2, err := idx.Put(bytes.NewReader(putContent2), "/addr1")
	assert.NoError(t, err)

	b1, err := idx.Get(content.BlockKeyFromCid(c1.ID))
	assert.NoError(t, err)

	b2, err := idx.Get(content.BlockKeyFromCid(c2.ID))
	assert.NoError(t, err)

	assert.True(t, b1.ID.Equals(c1.ID))
	assert.True(t, b2.ID.Equals(c2.ID))
}

func TestEvict(t *testing.T) {
	idx := NewMemoryIndex("/root", 2, EvictLeastPopular{}, storage.NoopStore{})
	c1, err := idx.Put(bytes.NewReader(putContent1), "/addr1")
	assert.NoError(t, err)
	c1.IncHitsCount()
	c1.IncHitsCount()
	c1.IncHitsCount()

	c2, err := idx.Put(bytes.NewReader(putContent2), "/addr3")
	assert.NoError(t, err)
	c2.IncHitsCount()
	c2.IncHitsCount()

	c3, err := idx.Put(bytes.NewReader(putContent3), "/addr4")
	assert.NoError(t, err)

	stored := []content.BlockKey{}
	for k := range idx.kvStore {
		stored = append(stored, k)
	}

	expected := []content.BlockKey{c3.Key(), c1.Key()}
	assert.ElementsMatch(t, expected, stored)
}
