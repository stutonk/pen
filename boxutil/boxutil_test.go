package boxutil

import (
	"bytes"
	"crypto/rand"
	"errors"
	mathrand "math/rand"
	"testing"
	"time"
)

var (
	errFoo = errors.New("foo error")
	salt   = make([]byte, 128)
	key    = Passkey([]byte("password"), salt)
)

type errReadWriter struct{}

func (e errReadWriter) Read(p []byte) (int, error) {
	return -1, errFoo
}

func (e errReadWriter) Write(p []byte) (int, error) {
	return -1, errFoo
}

func TestGen(t *testing.T) {
	const str = "test string"
	in := bytes.NewBufferString(str)
	e := make(chan error)
	a := gen(in, 4, e)
	s1, s2, s3 := string(*(<-a)), string(*(<-a)), string(*(<-a))
	if s1 != "test" || s2 != " str" || s3 != "ing" {
		t.Fatalf(
			"gen did not properly read/chunk data; wanted '%v', got '%v%v%v'",
			str,
			s1,
			s2,
			s3,
		)
	}

	if sx, ok := <-a; ok {
		t.Fatalf("extraneous data in channel: %v", sx)
	}

	if b := gen(nil, 0, e); b != nil {
		t.Fatal("successfully created generator with 0 chunk size")
	}

	_ = gen(errReadWriter{}, 1, e)
	select {
	case err := <-e:
		if err != errFoo {
			t.Fatal("incorrect or no error received; wanted errFoo")
		}
	case <-time.After(1 * time.Second):
		t.Fatal("no error returned for errReadWriter")
	}

}

func TestSystem(t *testing.T) {
	origReader := rand.Reader
	rand.Reader = errReadWriter{}
	if err := SealStream(
		bytes.NewBufferString("bar"),
		bytes.NewBuffer(nil),
		key,
	); err != errFoo {
		t.Fatal("incorrect or no error received; wanted errFoo")
	}
	rand.Reader = origReader

	size := mathrand.Intn(1024*24) + 1024*16 //at least 16KiB
	rndSlice := make([]byte, size)
	if _, err := mathrand.Read(rndSlice); err != nil {
		t.Fatal("could not create random slice")
	}
	msgBuf := bytes.NewBuffer(rndSlice)
	cryptBuf := bytes.NewBuffer(nil)
	if err := SealStream(msgBuf, cryptBuf, key); err != nil {
		t.Fatal(err)
	}
	if cryptBuf.String() == msgBuf.String() {
		t.Fatal("ciphertext matches plaintext")
	}

	if err := OpenStream(
		bytes.NewBufferString(cryptBuf.String()),
		bytes.NewBuffer(nil),
		&[32]byte{},
	); err != ErrIncorrectKey {
		t.Fatal("wrong key did not generate ErrIncorrectKey")
	}

	decryptBuf := bytes.NewBuffer(nil)
	if err := OpenStream(cryptBuf, decryptBuf, key); err != nil {
		t.Fatal(err)
	}
	if decryptBuf.String() != string(rndSlice) { //msg {
		t.Fatal("decrypt does not match plaintext")
	}
}

func TestWrite(t *testing.T) {
	e := make(chan error)
	defer close(e)
	c := make(chan *[]byte)
	defer close(c)
	x := []byte("baz")
	go func() {
		err := write(c, errReadWriter{}, e)
		if err != errFoo {
			t.Fatal("incorrect or no error received; wanted errFoo")
		}
	}()
	c <- &x
}
