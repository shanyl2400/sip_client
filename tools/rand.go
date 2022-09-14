package tools

import (
	"crypto/rand"
	"encoding/hex"
)

func RandString() string {
	buf := make([]byte, 8)
	rand.Read(buf)

	return hex.EncodeToString(buf)
}
