package content

import (
	"encoding/json"
	"time"

	"github.com/ipfs/go-cid"
)

// Address defines the location of a block of content, either identified by it CID or
// its file path
type Address struct {
	Cid      cid.Cid
	NodeName string
	Path     string
}

func (a *Address) MarshalJSON() ([]byte, error) {
	return json.Marshal(&struct {
		Cid  string `json:"cid"`
		Path string `json:"path"`
	}{
		Cid:  a.Cid.String(),
		Path: a.Path,
	})
}

// Block is a unit of content. It contains basic information necessary for retrieving
// the original content and evaluating the freshness of its data.
type Block struct {
	Address        Address   `json:"address"`
	CreatedAt      time.Time `json:"created_at"`
	size           int
	hitsCount      int
	lastAccessedAt *time.Time
}

type BlockWithData struct {
	Data []byte `json:"data"`
	Block
}

func NewBlock(data []byte, addr *Address) (*Block, error) {
	return &Block{
		size:      len(data),
		CreatedAt: time.Now(),
		Address:   *addr,
	}, nil
}

func (b *Block) ID() string {
	return b.Address.Cid.String()
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
