package pkg

import (
	"crypto/rand"
	"encoding/hex"
)

func GenerateToken() (string, error) {
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	return "hzx_" + hex.EncodeToString(bytes), nil
}
