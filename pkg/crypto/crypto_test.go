package crypto

import (
	"os"
	"testing"

	"github.com/matryer/is"
)

var ciphkey string

func TestMain(m *testing.M) {
	ciphkey = "65dTxbqk7rE3IFly1hnI1234"
	os.Exit(m.Run())
}

func TestEncrypt(t *testing.T) {
	is := is.New(t)
	ciphtxt := Encrypt(ciphkey, "radiation")
	is.Equal(ciphtxt, "Q0jhx4gItMsD")
}

func TestDecrypt(t *testing.T) {
	is := is.New(t)
	plaintxt := Decrypt(ciphkey, "Q0jhx4gItMsD")
	is.Equal(plaintxt, "radiation")
}

