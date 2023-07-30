package auth

import (
	"crypto/sha1"
	"encoding/base64"
)

func hasher(password string) string {
	hasher := sha1.New()
	hasher.Write([]byte(password))
	sha := base64.URLEncoding.EncodeToString(hasher.Sum(nil))
	return sha
}
