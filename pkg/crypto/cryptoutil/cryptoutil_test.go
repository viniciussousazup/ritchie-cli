package cryptoutil

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
	ciphtxt := Encrypt(ciphkey, "FBA4E319-65EF-542B-8243-6598B0BDC3AF", "radiation")
	is.Equal(ciphtxt, "d2vEmqxP7J1A+TZIpL2TSULqaqq6leqvAis6Lbv4U1s1dmuddZ78T4nBfrgym39A0uADC6FJqICJ1ktWy9JhPuSrqVd7Gl7537YCObLEZcJ+")
}

func TestDecrypt(t *testing.T) {
	is := is.New(t)
	plaintxt := Decrypt(ciphkey, "FBA4E319-65EF-542B-8243-6598B0BDC3AF", "d2vEmqxP7J1A+TZIpL2TSULqaqq6leqvAis6Lbv4U1s1dmuddZ78T4nBfrgym39A0uADC6FJqICJ1ktWy9JhPuSrqVd7Gl7537YCObLEZcJ+")
	is.Equal(plaintxt, "radiation")
}
