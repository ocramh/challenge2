package indexer

import (
	"bytes"
	"io"
	"testing"

	"github.com/ocramh/challenge2/pkg/storage"
	"github.com/stretchr/testify/assert"
)

var (
	putContent1 = []byte("content1")
	putContent2 = []byte("content2")
	putContent3 = []byte("content3")
)

func TestPut(t *testing.T) {
	idx := NewMemoryIndex("/root", 2, EvictLeastPopular{}, storage.NewNoopStore("/root"))

	testcases := []struct {
		src              io.Reader
		fPath            string
		expectedPath     string
		expectedStoreLen int
	}{
		{
			src:              bytes.NewReader(putContent1),
			fPath:            "/addr1",
			expectedPath:     "/root/addr1",
			expectedStoreLen: 1,
		},
		{
			src:              bytes.NewReader(putContent1),
			fPath:            "/another/addr",
			expectedPath:     "/root/another/addr",
			expectedStoreLen: 1,
		},
		{
			src:              bytes.NewReader(putContent2),
			fPath:            "/addr2",
			expectedPath:     "/root/addr2",
			expectedStoreLen: 2,
		},
		{
			src:              bytes.NewReader(putContent3),
			fPath:            "/addr3",
			expectedPath:     "/root/addr3",
			expectedStoreLen: 2,
		},
	}

	for _, testcase := range testcases {
		got, err := idx.Put(testcase.src, testcase.fPath)
		assert.NoError(t, err)
		assert.Equal(t, testcase.expectedPath, got.Address.Path)
		assert.Len(t, idx.kvStore, testcase.expectedStoreLen)
	}
}

func TestGet(t *testing.T) {
	idx := NewMemoryIndex("/root", 2, EvictLeastPopular{}, storage.NewNoopStore("/root"))

	c1, err := idx.Put(bytes.NewReader(putContent1), "/addr1")
	assert.NoError(t, err)

	c2, err := idx.Put(bytes.NewReader(putContent2), "/addr1")
	assert.NoError(t, err)

	b1, _, err := idx.Get(c1.Address.Cid)
	assert.NoError(t, err)

	b2, _, err := idx.Get(c2.Address.Cid)
	assert.NoError(t, err)

	assert.True(t, b1.Address.Cid.Equals(c1.Address.Cid))
	assert.True(t, b2.Address.Cid.Equals(c2.Address.Cid))
}

func TestEvict(t *testing.T) {
	idx := NewMemoryIndex("/root", 2, EvictLeastPopular{}, storage.NewNoopStore("/root"))
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

	stored := []string{}
	for k := range idx.kvStore {
		stored = append(stored, k)
	}

	expected := []string{c3.ID(), c1.ID()}
	assert.ElementsMatch(t, expected, stored)
}
