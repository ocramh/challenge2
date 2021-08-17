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

// BlocksStorage is an implementation of the StorageManager interface using a
// DAG representation of the available content backed by an ipfs blockservice
type BlocksStorage struct {
	blockSrv bserv.BlockService
	dagSrv   format.DAGService
	rootNode *mdag.ProtoNode
}

func NewBlocksStorage() *BlocksStorage {
	blockStore := bstore.NewBlockstore(dstore.NewMapDatastore())
	blockService := bserv.New(blockStore, offline.Exchange(blockStore))
	return &BlocksStorage{
		blockSrv: blockService,
		dagSrv:   mdag.NewDAGService(blockService),
		rootNode: mdag.NodeWithData(nil),
	}
}

func (b *BlocksStorage) Put(data []byte, name string) (*content.Address, error) {
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

func (b *BlocksStorage) Get(addr *content.Address) ([]byte, error) {
	node, err := b.dagSrv.Get(context.TODO(), addr.Cid)
	if err != nil {
		return nil, err
	}

	return node.RawData(), nil
}

func (b *BlocksStorage) Delete(addr *content.Address) error {
	err := b.rootNode.RemoveNodeLink(addr.NodeName)
	if err != nil {
		return err
	}

	return b.dagSrv.Remove(context.TODO(), addr.Cid)
}

func (b *BlocksStorage) RootID() string {
	return b.rootNode.String()
}
