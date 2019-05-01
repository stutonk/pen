package passkey

import (
	"bytes"
	"crypto/rand"
	"reflect"
	"testing"
)

func TestNew(t *testing.T) {
	const errFmt = "key generation failed: %v"

	salt, err := FreshSalt()
	if err != nil {
		t.Fatalf(errFmt, err)
	}

	pass := []byte("password1234")

	a := New(pass, salt)
	if len(a) != 32 {
		t.Fatalf("key has incorret length; want 32, got %d", len(a))
	}

	b := New(pass, salt)
	if !reflect.DeepEqual(a, b) {
		t.Fatal("same password/same salt combo generated different keys")
	}

	otherSalt, err := FreshSalt()
	if err != nil {
		t.Fatalf(errFmt, err)
	}

	c := New(pass, otherSalt)
	if reflect.DeepEqual(b, c) {
		t.Fatal("same password/different salt combo generated same key")
	}

	d := New([]byte("abc123"), otherSalt)
	if reflect.DeepEqual(c, d) {
		t.Fatal("different password/same salt combo generated same key")
	}

}

func TestFreshSalt(t *testing.T) {
	const errFmt = "salt generation failed: %v"
	a, err := FreshSalt()
	if err != nil {
		t.Fatalf(errFmt, err)
	}
	if len(a) != 128 {
		t.Fatalf("salt has incorrect length; want 128, got %d", len(a))
	}

	b, err := FreshSalt()
	if err != nil {
		t.Fatalf(errFmt, err)
	}

	if reflect.DeepEqual(a, b) {
		t.Fatal("salts are not unique between runs")
	}

	rand.Reader = bytes.NewReader(nil)
	_, err = FreshSalt()
	if err == nil {
		t.Fatal("failed to return error with empty reader")
	}
}
