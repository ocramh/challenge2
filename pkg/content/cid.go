package content

import (
	"github.com/ipfs/go-cid"
)

func CidFromBytes(b []byte) (cid.Cid, error) {
	format := cid.V0Builder{}

	return format.Sum(b)
}

func CidFromString(c string) (cid.Cid, error) {
	decoded, err := cid.Decode(c)
	if err != nil {
		return cid.Undef, ErrInvalidCID
	}

	return decoded, nil
}
