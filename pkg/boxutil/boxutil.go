/*
Package boxutil implements streaming versions of the functions in
nacl/secretbox as well as utility functions for generating suitable secret
keys from password strings
*/
package boxutil

import (
	"crypto/rand"
	"errors"
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
	passes    = 1
)

// ErrIncorrectKey should be used when an Open operation is attempted with
// the wrong key
var ErrIncorrectKey = errors.New("incorrect key")

type functor func(*[]byte) (*[]byte, error)

// gen is the first stage of a concurrent pipeline. It reads
// chunkSize-sized byte slices from the Reader (or however much is left)
// and sends them down the pipeline
func gen(in io.Reader, size int, e chan<- error) <-chan *[]byte {
	if size < 1 {
		return nil
	}
	out := make(chan *[]byte)
	go func() {
		defer close(out)
		for {
			chunk := make([]byte, size)
			n, err := in.Read(chunk)
			if err != nil {
				if err == io.EOF {
					break
				} else {
					e <- err
					return
				}
			} else if n < size {
				chunk = chunk[:n]
			}
			out <- &chunk
		}
	}()
	return out
}

/*
OpenStream applies secretbox.Open to an io.Reader stream of []byte chunks
and writes them to an io.Writer stream. Each chunk should begin with the
nonce used to seal the chunk. If this function returns an error, the data
written to the output stream should be considered corrupted and discarded.
*/
func OpenStream(r io.Reader, w io.Writer, key *[keyLen]byte) error {
	e := make(chan error)
	defer close(e)
	in := gen(r, nonceLen+chunkSize+secretbox.Overhead, e)
	f := func(data *[]byte) (*[]byte, error) {
		var nonce [nonceLen]byte
		copy(nonce[:], (*data)[:nonceLen])
		opened, ok := secretbox.Open(nil, (*data)[nonceLen:], &nonce, key)
		if !ok {
			return nil, ErrIncorrectKey
		}
		return &opened, nil
	}
	out := work(in, f, e)
	return write(out, w, e)
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
	keyBytes := argon2.IDKey(password, salt, passes, memory, cpus, keyLen)
	copy(key[:], keyBytes)
	return &key
}

/*
SealStream applies secretbox.Seal to an io.Reader stream of []byte chunks
and writes them to an io.Writer stream. Each resulting chunk will be
prepended with the nonce used in its creation and will additionally be
secretbox.Overhead bytes longer. If this function returns a non-nil error,
the data written to the output stream should be considered corrupted and
discarded.
*/
func SealStream(r io.Reader, w io.Writer, key *[keyLen]byte) error {
	e := make(chan error)
	defer close(e)
	in := gen(r, chunkSize, e)
	f := func(data *[]byte) (*[]byte, error) {
		nonceBytes := make([]byte, nonceLen)
		if _, err := rand.Read(nonceBytes); err != nil {
			// this may happen if the system runs out of entropy
			// it may be better to wait for more instead
			return nil, err
		}
		var nonce [nonceLen]byte
		copy(nonce[:], nonceBytes)
		sealed := secretbox.Seal(nonceBytes, *data, &nonce, key)
		return &sealed, nil
	}
	out := work(in, f, e)
	return write(out, w, e)
}

// work is the middle stage of a concurrent pipeline. It applies a
// transformation function f to incoming data and passes it down the pile.
// Errors upstream should close their send channel so this goroutine will
// terminate.
func work(in <-chan *[]byte, f functor, e chan<- error) <-chan *[]byte {

	out := make(chan *[]byte)
	go func() {
		defer close(out)
		for data := range in {
			v, err := f(data)
			if err != nil {
				e <- err
				return
			}
			out <- v
		}
	}()
	return out
}

// write is the final stage of a concurrent pipeline. It commits transformed
// data to an output stream. Errors upstream should close their send channel
// so this function can abort with an error.
func write(in <-chan *[]byte, w io.Writer, e <-chan error) error {
	for {
		select {
		case data, ok := <-in:
			if !ok {
				return nil
			}
			if _, err := w.Write(*data); err != nil {
				return err
			}
		case err := <-e:
			return err
		}
	}
}
