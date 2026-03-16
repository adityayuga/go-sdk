package crypto

import (
	"strings"
	"testing"
)

func TestNewAesGcmEncryptionRejectsInvalidKeyLengths(t *testing.T) {
	t.Parallel()

	testCases := []struct {
		name   string
		secret string
		valid  bool
	}{
		{name: "aes128", secret: strings.Repeat("a", 16), valid: true},
		{name: "aes192", secret: strings.Repeat("b", 24), valid: true},
		{name: "aes256", secret: strings.Repeat("c", 32), valid: true},
		{name: "too short", secret: strings.Repeat("d", 15), valid: false},
		{name: "invalid 64", secret: strings.Repeat("e", 64), valid: false},
	}

	for _, testCase := range testCases {
		testCase := testCase
		t.Run(testCase.name, func(t *testing.T) {
			t.Parallel()

			enc, err := NewAesGcmEncryption(nil, testCase.secret)
			if testCase.valid {
				if err != nil {
					t.Fatalf("expected key length %d to be valid: %v", len(testCase.secret), err)
				}
				if enc == nil {
					t.Fatal("expected encryptor instance")
				}
				return
			}

			if err == nil {
				t.Fatalf("expected key length %d to be rejected", len(testCase.secret))
			}
		})
	}
}

func TestAesGcmEncryptionRoundTrip(t *testing.T) {
	t.Parallel()

	enc, err := NewAesGcmEncryption(nil, strings.Repeat("k", 32))
	if err != nil {
		t.Fatalf("unexpected constructor error: %v", err)
	}

	ciphertext, err := enc.Encrypt("hello world")
	if err != nil {
		t.Fatalf("unexpected encrypt error: %v", err)
	}

	plaintext, err := enc.Decrypt(string(ciphertext))
	if err != nil {
		t.Fatalf("unexpected decrypt error: %v", err)
	}

	if plaintext != "hello world" {
		t.Fatalf("unexpected plaintext: %q", plaintext)
	}
}

func TestAesGcmDecryptRejectsShortCiphertext(t *testing.T) {
	t.Parallel()

	enc, err := NewAesGcmEncryption(nil, strings.Repeat("k", 32))
	if err != nil {
		t.Fatalf("unexpected constructor error: %v", err)
	}

	_, err = enc.Decrypt("short")
	if err == nil {
		t.Fatal("expected error for short ciphertext")
	}
}

func TestGenerateSecureTokenValidation(t *testing.T) {
	t.Parallel()

	token, err := GenerateSecureToken(0)
	if err != nil {
		t.Fatalf("unexpected error for zero length: %v", err)
	}
	if token != "" {
		t.Fatalf("expected empty token, got %q", token)
	}

	_, err = GenerateSecureToken(-1)
	if err == nil {
		t.Fatal("expected error for negative token length")
	}
}

func TestGenerateSecureTokenLengthAndAlphabet(t *testing.T) {
	t.Parallel()

	token, err := GenerateSecureToken(32)
	if err != nil {
		t.Fatalf("unexpected token error: %v", err)
	}
	if len(token) != 32 {
		t.Fatalf("expected token length 32, got %d", len(token))
	}
	if strings.Contains(token, "=") {
		t.Fatalf("expected URL-safe token without padding, got %q", token)
	}
}

func TestGenerateCode(t *testing.T) {
	t.Parallel()

	if code := GenerateCode(0); code != "" {
		t.Fatalf("expected empty code for zero length, got %q", code)
	}

	code := GenerateCode(64)
	if len(code) != 64 {
		t.Fatalf("expected code length 64, got %d", len(code))
	}

	for _, char := range code {
		if char < '0' || char > '9' {
			t.Fatalf("expected only numeric characters, got %q in %q", char, code)
		}
	}
}
