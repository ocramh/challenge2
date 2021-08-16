package provider

import (
	"bytes"
	"os"

	"github.com/ocramh/challenge2/pkg/content"
	"github.com/ocramh/challenge2/pkg/indexer"
	"github.com/ocramh/challenge2/pkg/storage"
)

type Provider struct {
	idx indexer.Indexer
}

func New(cap int, rootDir string) (*Provider, error) {
	err := os.MkdirAll(rootDir, 0755)
	if err != nil {
		return nil, err
	}

	return &Provider{
		idx: indexer.NewMemoryIndex(
			rootDir, cap, indexer.EvictLeastPopular{}, storage.NoopStore{},
		),
	}, nil
}

func (p *Provider) AddItem(b []byte) ([]*content.Block, error) {
	block, err := p.idx.Put(bytes.NewReader(b), string(b))
	if err != nil {
		return nil, err
	}

	return []*content.Block{block}, nil
}

func (p *Provider) GetItem(key content.BlockKey) (*content.Block, error) {
	block, err := p.idx.Get(key)
	if err != nil {
		return nil, err
	}

	return block, nil
}
