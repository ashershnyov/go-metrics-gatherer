package hasher

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"hash"
	"sync"
)

// Hasher allows for generating sha256 hashes.
type Hasher struct {
	mu   *sync.Mutex
	hash hash.Hash
}

// NewHasher returns a new hasher.
func NewHasher(key string) *Hasher {
	return &Hasher{
		mu:   &sync.Mutex{},
		hash: hmac.New(sha256.New, []byte(key)),
	}
}

// Hash returns a hex-encoded sha256 hash of data.
func (h *Hasher) Hash(data []byte) string {
	h.mu.Lock()
	h.hash.Reset()
	h.hash.Write(data)
	s := h.hash.Sum(nil)
	h.mu.Unlock()
	return hex.EncodeToString(s)
}
