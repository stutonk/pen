package boxutil

import (
	"bytes"
	"crypto/rand"
	"testing"

	"github.com/stutonk/ef/pkg/passkey"
)

func StreamTest(t *testing.T) {
	msgBuf := bytes.NewBufferString("secret message")
	cryptBuf := bytes.NewBuffer(nil)

	salt := make([]byte, 128)
	if _, err := rand.Read(salt); err != nil {
		t.Fatalf("couldn't create salt: %v", err)
	}

	key := passkey.New([]byte("password"), salt)
	SealStream(msgBuf, cryptBuf, key)
	if cryptBuf.String() == msgBuf.String() {
		t.Fatal("ciphertext matches plaintext")
	}

	decryptBuf := bytes.NewBuffer(nil)
	OpenStream(cryptBuf, decryptBuf, key)
	if decryptBuf.String() != msgBuf.String() {
		t.Fatalf(
			"decrypt does not match plaintext; wanted '%v', got '%v",
			msgBuf.String(),
			decryptBuf.String(),
		)
	}
}
