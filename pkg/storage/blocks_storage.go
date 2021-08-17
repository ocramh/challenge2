package storage

import (
	"context"

	bserv "github.com/ipfs/go-blockservice"
	dstore "github.com/ipfs/go-datastore"
	bstore "github.com/ipfs/go-ipfs-blockstore"
	offline "github.com/ipfs/go-ipfs-exchange-offline"
	format "github.com/ipfs/go-ipld-format"
	mdag "github.com/ipfs/go-merkledag"

	"github.com/ocramh/challenge2/pkg/content"
)

// BlocksStore is an implementation of the StorageManager interface using a
// DAG representation of the available content backed by an ipfs blockservice
type BlocksStore struct {
	blockSrv    bserv.BlockService
	dagSrv      format.DAGService
	rootNode    *mdag.ProtoNode
	maxCapacity int
}

func NewBlocksStore(maxCapacity int) *BlocksStore {
	blockStore := bstore.NewBlockstore(dstore.NewMapDatastore())
	blockService := bserv.New(blockStore, offline.Exchange(blockStore))
	return &BlocksStore{
		blockSrv:    blockService,
		dagSrv:      mdag.NewDAGService(blockService),
		rootNode:    mdag.NodeWithData(nil),
		maxCapacity: maxCapacity,
	}
}

func (b *BlocksStore) Put(data []byte, name string) (*content.Address, error) {
	if b.Size() >= b.maxCapacity {
		return nil, ErrNoStorageAvailable
	}

	nd := mdag.NodeWithData(data)
	err := b.rootNode.AddNodeLink(name, nd)
	if err != nil {
		return nil, err
	}

	err = b.dagSrv.Add(context.TODO(), nd)
	if err != nil {
		return nil, err
	}

	return &content.Address{
		Cid:      nd.Cid(),
		Path:     nd.Cid().String(),
		NodeName: name,
	}, nil
}

func (b *BlocksStore) Get(addr *content.Address) ([]byte, error) {
	node, err := b.dagSrv.Get(context.TODO(), addr.Cid)
	if err != nil {
		return nil, err
	}

	return node.RawData(), nil
}

func (b *BlocksStore) Delete(addr *content.Address) error {
	err := b.rootNode.RemoveNodeLink(addr.NodeName)
	if err != nil {
		return err
	}

	return b.dagSrv.Remove(context.TODO(), addr.Cid)
}

// Size returns the number of the root node links.
// While this is fine for a demonstrating purpose, a more correct implementation should
// return the rootNode Size() which includes the total size of the data addressed by
// the node, including the size of its references
func (b *BlocksStore) Size() int {
	return len(b.rootNode.Links())
}

func (b *BlocksStore) Capacity() int {
	return b.maxCapacity
}
