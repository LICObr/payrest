package service

import (
	"crypto/rand"
	"encoding/hex"
)

// generateID returns a random 16-byte hex identifier, shared by every
// service that needs to create new domain entities.
func generateID() string {
	buf := make([]byte, 16)
	_, _ = rand.Read(buf)
	return hex.EncodeToString(buf)
}
