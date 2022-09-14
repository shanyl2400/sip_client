package tools

import (
	"crypto/sha1"
	"encoding/hex"
)

func Hash(text string) string {
	sha128 := sha1.New()
	sha128.Write([]byte(text))
	return hex.EncodeToString(sha128.Sum(nil))
}
