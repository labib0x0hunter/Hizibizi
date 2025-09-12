package utils

import (
	"github.com/cespare/xxhash/v2"
)

// func GetHash(word string) uint64 {
// 	h := xxhash.NewWithSeed(seed)
// }

type Hash struct {
	h *xxhash.Digest
}

func NewHash() *Hash {
	// os.Remove(filepath.Join(path, seedFile))
	return &Hash{
		h: xxhash.NewWithSeed(seed),
	}
}

func (h *Hash) WriteString(msg string) error {
	if _, err := h.h.WriteString(msg); err != nil {
		return err
	}
	return nil
}

func (h *Hash) Sum() uint64 {
	return h.h.Sum64()
}

func (h *Hash) Reset() {
	h.h.Reset()
}
