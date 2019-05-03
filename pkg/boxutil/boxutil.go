/*
Package boxutil implements streaming versions of the functions in
nacl/secretbox as well as utility functions for generating suitable secret
keys from password strings
*/
package boxutil

import (
	"crypto/rand"
	"io"
	"runtime"

	"golang.org/x/crypto/argon2"
	"golang.org/x/crypto/nacl/secretbox"
)

const (
	chunkSize = 16 * 1024
	keyLen    = 32
	memory    = 64 * 1024
	nonceLen  = 24
	time      = 1
)

func gen(in io.Reader, size int) <-chan *[]byte {
	out := make(chan *[]byte)
	go func() {
		for more := true; more; /**/ {
			chunk := make([]byte, size)
			if n, err := in.Read(chunk); err != nil {
				if err == io.EOF {
					chunk = chunk[:n]
					more = false
				} else {
					panic(err)
				}
			}
			out <- &chunk
		}
		close(out)
	}()
	return out
}

func open(in <-chan *[]byte, key *[keyLen]byte) <-chan *[]byte {
	out := make(chan *[]byte)
	go func() {
		for data := range in {
			var nonce [nonceLen]byte
			copy(nonce[:], (*data)[:nonceLen])
			opened, ok := secretbox.Open(nil, (*data)[nonceLen:], &nonce, key)
			if !ok {
				panic("incorrect key")
			}
			out <- &opened
		}
		close(out)
	}()
	return out
}

/*
OpenStream applies secretbox.Open to an io.Reader stream of []byte chunks
and writes them to an io.Writer stream. Each chunk should begin with the
nonce used to seal the chunk.
*/
func OpenStream(r io.Reader, w io.Writer, key *[keyLen]byte) {
	in := gen(r, nonceLen+chunkSize+secretbox.Overhead)
	out := open(in, key)
	write(out, w)
}

/*
Passkey creates a 32-byte array (as required by nacl/secretbox) containing the
result of applying the argon2id algorithm to a given password and salt.
The underlying argon2 parameters are those recommended in the argon2 library
documentation.
*/
func Passkey(password, salt []byte) *[keyLen]byte {
	var key [keyLen]byte
	cpus := uint8(runtime.NumCPU())
	keyBytes := argon2.IDKey(password, salt, time, memory, cpus, keyLen)
	copy(key[:], keyBytes)
	return &key
}

func seal(in <-chan *[]byte, key *[keyLen]byte) <-chan *[]byte {
	out := make(chan *[]byte)
	go func() {
		for data := range in {
			nonceBytes := make([]byte, nonceLen)
			if _, err := rand.Read(nonceBytes); err != nil {
				panic(err)
			}
			var nonce [nonceLen]byte
			copy(nonce[:], nonceBytes)
			sealed := secretbox.Seal(nonceBytes, *data, &nonce, key)
			out <- &sealed
		}
		close(out)
	}()
	return out
}

/*
SealStream applies secretbox.Seal to an io.Reader stream of []byte chunks
and writes them to an io.Writer stream. Each resulting chunk will be
prepended with the nonce used in its creation and will additionally be
secretbox.Overhead bytes longer.
*/
func SealStream(r io.Reader, w io.Writer, key *[keyLen]byte) {
	in := gen(r, chunkSize)
	out := seal(in, key)
	write(out, w)
}

func write(in <-chan *[]byte, w io.Writer) {
	for data := range in {
		if _, err := w.Write(*data); err != nil {
			panic(err)
		}
	}
}
