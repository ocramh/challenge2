package content

import (
	"github.com/ipfs/go-cid"
	"github.com/multiformats/go-multihash"
)

func CidFromBytes(b []byte) (cid.Cid, error) {
	pref := cid.Prefix{
		Version:  1,
		Codec:    cid.Raw,
		MhType:   multihash.SHA2_256,
		MhLength: -1,
	}

	return pref.Sum(b)
}
