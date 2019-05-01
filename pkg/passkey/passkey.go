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

const (
	keyLen  = 32
	saltLen = 128
	time    = 1
	memory  = 64 * 1024
)

/*
New creates a 32-byte array (as required by nacl/secretbox) containing the
result of applying the argon2id algorithm to a given password and salt.
The underlying argon2 parameters are those recommended in the argon2 library
documentation.
*/
func New(password, salt []byte) (key [keyLen]byte) {
	cpus := uint8(runtime.NumCPU())
	keyBytes := argon2.IDKey(password, salt, time, memory, cpus, keyLen)
	copy(key[:], keyBytes)
	return key
}

// FreshSalt generates a byte slice initialized to a random value using
// crypto/rand. A length of 128 should guarantee 'sufficient' uniqieness.
func FreshSalt() ([]byte, error) {
	salt := make([]byte, saltLen)
	_, err := rand.Read(salt)
	if err != nil {
		return nil, err
	}
	return salt, err
}
