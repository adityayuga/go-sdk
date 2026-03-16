package crypto

import (
	"crypto/rand"
	"encoding/base64"
	"strings"

	"github.com/adityayuga/go-sdk/errors"
)

// GenerateSecureToken generates a URL-safe base64 encoded random string
func GenerateSecureToken(length int) (string, error) {
	if length < 0 {
		return "", errors.New("[pkg secure token] length must be greater than or equal to 0")
	}
	if length == 0 {
		return "", nil
	}

	byteLength := (length*6 + 7) / 8 // convert bit length to byte length
	b := make([]byte, byteLength)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	token := base64.URLEncoding.EncodeToString(b)
	return strings.TrimRight(token, "=")[:length], nil
}
