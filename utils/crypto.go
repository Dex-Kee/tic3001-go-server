package utils

import "crypto/sha256"

func SHA256Digest(content string) string {
	return string(sha256.New().Sum([]byte(content)))
}
