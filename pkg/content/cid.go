package content

import (
	"github.com/ipfs/go-cid"
)

func CidFromBytes(b []byte) (cid.Cid, error) {
	format := cid.V0Builder{}

	return pref.Sum(b)
}
