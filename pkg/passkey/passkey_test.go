package passkey

import (
	"bytes"
	"crypto/rand"
	"reflect"
	"testing"
)

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
	if _, err = FreshSalt(); err == nil {
		t.Fatal("failed to return error with empty reader")
	}
}

func TestNew(t *testing.T) {
	salt, _ := FreshSalt()
	pass := []byte("password1234")

	a := New(pass, salt)
	if len(a) != 32 {
		t.Fatalf("key has incorret length; want 32, got %d", len(a))
	}

	b := New(pass, salt)
	if !reflect.DeepEqual(*a, *b) {
		t.Fatal("same password/same salt combo generated different keys")
	}

	otherSalt, _ := FreshSalt()
	c := New(pass, otherSalt)
	if reflect.DeepEqual(*b, *c) {
		t.Fatal("same password/different salt combo generated same key")
	}

	d := New([]byte("abc123"), otherSalt)
	if reflect.DeepEqual(*c, *d) {
		t.Fatal("different password/same salt combo generated same key")
	}
}
