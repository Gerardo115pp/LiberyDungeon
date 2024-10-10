package helpers

import (
	"crypto/sha1"
	"fmt"
)

func GenerateSha1ID(from string) string {
	data := []byte(from)

	hasher := sha1.New()

	hasher.Write(data)
	hash := hasher.Sum(nil)

	return fmt.Sprintf("%x", hash)
}
