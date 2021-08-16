package content

import (
	"encoding/json"
	"time"

	"github.com/ipfs/go-cid"
)

// BlockKey string representation of a Cid
type BlockKey string

func BlockKeyFromCid(c cid.Cid) BlockKey {
	return BlockKey(c.String())
}

// Address defines the location of a block of content
type Address struct {
	Filepath string `json:"file_path"`
}

// Block is a unit of content. It contains basic information necessary for retrieving
// the original content and evaluating the freshness of its data.
type Block struct {
	ID             cid.Cid   `json:"-"`
	Address        Address   `json:"address"`
	CreatedAt      time.Time `json:"created_at"`
	size           int
	hitsCount      int
	lastAccessedAt *time.Time
}

func NewBlock(data []byte, addr *Address) (*Block, error) {
	blockCid, err := CidFromBytes(data)
	if err != nil {
		return nil, err
	}

	return &Block{
		ID:        blockCid,
		size:      len(data),
		CreatedAt: time.Now(),
		Address:   *addr,
	}, nil
}

func (b *Block) MarshalJSON() ([]byte, error) {
	type BlockAlias Block
	return json.Marshal(&struct {
		ID        string `json:"cid"`
		HitsCount int    `json:"hits_count"`
		*BlockAlias
	}{
		ID:         b.ID.String(),
		HitsCount:  b.GetHitsCount(),
		BlockAlias: (*BlockAlias)(b),
	})
}

// Key returns a BlockKey representation of the cid
func (b *Block) Key() BlockKey {
	return BlockKey(b.ID.String())
}

func (b *Block) Size() int {
	return b.size
}

// IncHitsCount increments the block hits count by one
func (b *Block) IncHitsCount() {
	now := time.Now()
	b.lastAccessedAt = &now
	b.hitsCount = b.hitsCount + 1
}

// GetHitsCount returns the number of hits
func (b *Block) GetHitsCount() int {
	return b.hitsCount
}
