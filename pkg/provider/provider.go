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
			rootDir, cap, indexer.EvictLeastPopular{}, storage.NewBlocksStorage(),
		),
	}, nil
}

func (p *Provider) AddItem(b []byte) (*content.Block, error) {
	return p.idx.Put(bytes.NewReader(b), string(b))
}

func (p *Provider) GetItem(key string) (*content.BlockWithData, error) {
	blockCid, err := content.CidFromString(key)
	if err != nil {
		return nil, err
	}

	block, data, err := p.idx.Get(blockCid)
	if err != nil {
		return nil, err
	}

	return &content.BlockWithData{
		Data:  data,
		Block: *block,
	}, nil
}
