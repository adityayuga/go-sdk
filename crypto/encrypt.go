package crypto

import (
	"context"
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"

	"github.com/adityayuga/go-sdk/errors"
)

type (
	Encryptor interface {
		Encrypt(plaintext string) ([]byte, error)
		Decrypt(ciphertext string) (string, error)
	}

	AesGcmEncryption struct {
		// secret 16 chars -> AES 128
		// secret 24 chars -> AES 192
		// secret 32 chars -> AES 256
		secret string
	}
)

func NewAesGcmEncryption(ctx context.Context, secret string) (Encryptor, error) {
	switch len(secret) {
	case 16, 24, 32:
	default:
		return nil, errors.New("[pkg encrypt] secret must be 16 / 24 / 32 chars")
	}

	return AesGcmEncryption{secret: secret}, nil
}

func (m AesGcmEncryption) Encrypt(plaintext string) ([]byte, error) {
	aes, err := aes.NewCipher([]byte(m.secret))
	if err != nil {
		return nil, err
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return nil, err
	}

	// We need a 12-byte nonce for GCM (modifiable if you use cipher.NewGCMWithNonceSize())
	// A nonce should always be randomly generated for every encryption.
	nonce := make([]byte, gcm.NonceSize())
	_, err = rand.Read(nonce)
	if err != nil {
		return nil, err
	}

	// ciphertext here is actually nonce+ciphertext
	// So that when we decrypt, just knowing the nonce size
	// is enough to separate it from the ciphertext.
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)

	return ciphertext, nil
}

func (m AesGcmEncryption) Decrypt(ciphertext string) (string, error) {
	aes, err := aes.NewCipher([]byte(m.secret))
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(aes)
	if err != nil {
		return "", err
	}

	// Since we know the ciphertext is actually nonce+ciphertext
	// And len(nonce) == NonceSize(). We can separate the two.
	nonceSize := gcm.NonceSize()
	if len(ciphertext) < nonceSize {
		return "", errors.New("[pkg encrypt] ciphertext too short")
	}

	nonce, ciphertext := ciphertext[:nonceSize], ciphertext[nonceSize:]

	plaintext, err := gcm.Open(nil, []byte(nonce), []byte(ciphertext), nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}
