package crypto

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/hex"

	"github.com/adityayuga/go-sdk/errors"
)

type (
	Hash struct {
		secret string
	}

	Hasher interface {
		HashSHA256(data string) string
		HashSHA512(data string) string
	}
)

func NewHasher(ctx context.Context, secret string) (Hasher, error) {
	if len(secret) != 64 && len(secret) != 32 {
		return nil, errors.New("[pkg hash] secret must be 32 / 64 chars")
	}
	return Hash{secret: secret}, nil
}

func (h Hash) HashSHA512(data string) string {
	hmac := hmac.New(sha512.New, []byte(h.secret))
	hmac.Write([]byte(data))
	dataHmac := hmac.Sum(nil)

	return hex.EncodeToString(dataHmac)
}

func (h Hash) HashSHA256(data string) string {
	hmac := hmac.New(sha256.New, []byte(h.secret))
	hmac.Write([]byte(data))
	dataHmac := hmac.Sum(nil)

	return hex.EncodeToString(dataHmac)
}
