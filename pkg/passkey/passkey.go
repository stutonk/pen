/*
Package passkey implements secure expansion of passwords into
cryptographic keys suitable for use with nacl/secretbox. Also included is a
utility function for generating salts to use with the key function.
*/
package passkey

import (
	"crypto/rand"
	"runtime"

	"golang.org/x/crypto/argon2"
)

// keys should be len 32 to be compatible with nacl/secretbox
const (
	keyLen  = 32
	saltLen = 128
	time    = 1
	memory  = 64 * 1024
)

/*
FreshSalt generates a byte slice initialized to a random value using
crypto/rand. A length of 128 should guarantee 'sufficient' uniqieness.
Be aware that this function reads from crypto/rand's Reader; if too many
calls are made read from this reader in too short a time, this function
may return an EOF error.
*/
func FreshSalt() ([]byte, error) {
	salt := make([]byte, saltLen)
	if _, err := rand.Read(salt); err != nil {
		return nil, err
	}
	return salt, nil
}

/*
New creates a 32-byte array (as required by nacl/secretbox) containing the
result of applying the argon2id algorithm to a given password and salt.
The underlying argon2 parameters are those recommended in the argon2 library
documentation.
*/
func New(password, salt []byte) *[keyLen]byte {
	var key [keyLen]byte
	cpus := uint8(runtime.NumCPU())
	keyBytes := argon2.IDKey(password, salt, time, memory, cpus, keyLen)
	copy(key[:], keyBytes)
	return &key
}
