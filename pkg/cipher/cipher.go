package cipher

import (
	"crypto/rand"
	"fmt"
	"io"

	"golang.org/x/crypto/nacl/secretbox"
)

const (
	keyLen   = 32
	nonceLen = 24
)

func FreshNonce() (*[nonceLen]byte, error) {
	var nonce [nonceLen]byte
	if _, err := io.ReadFull(rand.Reader, nonce[:]); err != nil {
		return nil, err
	}
	return &nonce, nil
}

func Encrypt(message []byte, nonce *[nonceLen]byte, key *[keyLen]byte) []byte {
	return secretbox.Seal(nil, message, nonce, key)
}

func Decrypt(message []byte, nonce *[nonceLen]byte, key *[keyLen]byte) ([]byte, error) {
	plaintext, ok := secretbox.Open(nil, message, nonce, key)
	if !ok {
		return nil, fmt.Errorf("incorrect nonce or key")
	}
	return plaintext, nil
}
