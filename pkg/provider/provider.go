package provider

import (
	"bytes"
	"log"
	"os"

	"github.com/ocramh/challenge2/pkg/content"
	"github.com/ocramh/challenge2/pkg/dag"
	"github.com/ocramh/challenge2/pkg/indexer"
	"github.com/ocramh/challenge2/pkg/storage"
)

type Provider struct {
	idx indexer.Indexer
	nd  *dag.NodesManager
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
		nd: dag.NewNodesManager(),
	}, nil
}

func (p *Provider) AddItem(b []byte) ([]*content.Block, error) {
	block, err := p.idx.Put(bytes.NewReader(b), string(b))
	if err != nil {
		return nil, err
	}

	ndCid, err := p.nd.AddNodeLink(b, block.ID.String())
	if err != nil {
		return nil, err
	}

	log.Println(ndCid)

	return []*content.Block{block}, nil
}

func (p *Provider) GetItem(key content.BlockKey) (*content.Block, error) {
	block, err := p.idx.Get(key)
	if err != nil {
		return nil, err
	}

	// node, err := p.nd.GetNodeLink(block.ID)
	// if err != nil {
	// 	return nil, err
	// }
	// log.Println(node.Cid())

	return block, nil
}
