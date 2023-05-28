package hash

import (
	"github.com/cespare/xxhash"
)

type Hasher interface {
	Sum64(bytes []byte) uint64
}

func NewDefaultHasher() Hasher {
	return &xxhashHasher{}
}

type xxhashHasher struct {
}

func (x *xxhashHasher) Sum64(bytes []byte) uint64 {
	return xxhash.Sum64(bytes)
}
