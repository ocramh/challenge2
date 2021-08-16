package dag

import (
	"context"
	"log"

	bserv "github.com/ipfs/go-blockservice"
	"github.com/ipfs/go-cid"
	dstore "github.com/ipfs/go-datastore"
	bstore "github.com/ipfs/go-ipfs-blockstore"
	offline "github.com/ipfs/go-ipfs-exchange-offline"
	format "github.com/ipfs/go-ipld-format"
	mdag "github.com/ipfs/go-merkledag"
)

type NodesManager struct {
	blockSrv bserv.BlockService
	dagSrv   format.DAGService
	rootNode *mdag.ProtoNode
}

func NewNodesManager() *NodesManager {
	blockStore := bstore.NewBlockstore(dstore.NewMapDatastore())
	blockService := bserv.New(blockStore, offline.Exchange(blockStore))
	return &NodesManager{
		blockSrv: blockService,
		dagSrv:   mdag.NewDAGService(blockService),
		rootNode: mdag.NodeWithData(nil),
	}
}

func (n *NodesManager) AddNodeLink(content []byte, name string) (cid.Cid, error) {
	nd := mdag.NodeWithData([]byte(content))
	err := n.rootNode.AddNodeLink(name, nd)
	if err != nil {
		return cid.Undef, err
	}

	err = n.dagSrv.Add(context.TODO(), nd)
	if err != nil {
		return cid.Undef, err
	}

	return nd.Cid(), nil
}

func (n *NodesManager) GetNodeLink(nodeCID cid.Cid) (format.Node, error) {
	node, err := n.dagSrv.Get(context.TODO(), nodeCID)
	if err != nil {
		return nil, err
	}

	return node, nil
}

func (n *NodesManager) RemoveNode(nodeCID cid.Cid) error {
	err := n.dagSrv.Remove(context.TODO(), nodeCID)
	log.Println(n.RootID())
	return err

}

func (n *NodesManager) RootID() string {
	return n.rootNode.String()
}
