package security

import (
	"crypto/sha256"
	"encoding/hex"

	"villainrsty-ecommerce-server/internal/core/auth/ports"
)

type sha256TokenHasher struct{}

func NewSHA256TokenHasher() ports.TokenHasher {
	return &sha256TokenHasher{}
}

func (h *sha256TokenHasher) Hash(token string) (string, error) {
	sum := sha256.Sum256([]byte(token))
	return hex.EncodeToString(sum[:]), nil // 64 char
}
