package cipher

import (
	"bytes"
	"crypto/rand"
	"reflect"
	"testing"

	"github.com/stutonk/ef/pkg/passkey"
)

// shared state is required to keep from emptying crypto/rand. If this test
// fails with nil dereference errors, that's what happened.
var (
	salt, _ = passkey.FreshSalt()
	n1, _   = FreshNonce()
	n2, _   = FreshNonce()
	key     = passkey.New(pass, salt)
	pass    = []byte("password")
	msg     = []byte("We are discovered! Flee at once!")
	encMsg  = Encrypt(msg, n1, key)
)

func TestFreshNonce(t *testing.T) {
	if reflect.DeepEqual(*n1, *n2) {
		t.Fatal("nonces are not unique between runs")
	}

	rand.Reader = bytes.NewReader(nil)
	if _, err := FreshNonce(); err == nil {
		t.Fatal("failed to return error with empty reader")
	}
}

func TestEncrypt(t *testing.T) {
	if reflect.DeepEqual(encMsg, []byte(msg)) {
		t.Fatal("ciphertext identical to plaintext")
	}

	if em2 := Encrypt(msg, n2, key); reflect.DeepEqual(em2, encMsg) {
		t.Fatal("different nonces created identical ciphertexts")
	}

	k2 := passkey.New([]byte("80085"), salt)
	if em3 := Encrypt(msg, n1, k2); reflect.DeepEqual(em3, encMsg) {
		t.Fatal("different keys created identical ciphertexts")
	}
}

func TestDecrypt(t *testing.T) {
	if dm, _ := Decrypt(encMsg, n1, key); !reflect.DeepEqual(dm, msg) {
		t.Fatal("decrypted message differs from original message")
	}

	if _, err := Decrypt(encMsg, &[nonceLen]byte{}, key); err == nil {
		t.Fatal("decrypt succeds with invalid nonce")
	}

	if _, err := Decrypt(encMsg, n1, &[keyLen]byte{}); err == nil {
		t.Fatal("decrypt succeeds with invalid key")
	}
}
