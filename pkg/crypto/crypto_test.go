package crypto

import (
	"os"
	"testing"

	"github.com/matryer/is"
)

var cipherKey string

func TestMain(m *testing.M) {
	cipherKey = "65dTxbqk7rE3IFly1hnI1234"
	os.Exit(m.Run())
}

func TestEncrypt(t *testing.T) {
	i := is.New(t)
	encrypt := Encrypt(cipherKey, "FBA4E319-65EF-542B-8243-6598B0BDC3AF", "radiation")
	i.Equal(encrypt, "d2vEmqxP7J1A+TZIpL2TSULqaqq6leqvAis6Lbv4U1s1dmuddZ78T4nBfrgym39A0uADC6FJqICJ1ktWy9JhPuSrqVd7Gl7537YCObLEZcJ+")
}

func TestDecrypt(t *testing.T) {
	i := is.New(t)
	decrypt := Decrypt(cipherKey, "FBA4E319-65EF-542B-8243-6598B0BDC3AF", "d2vEmqxP7J1A+TZIpL2TSULqaqq6leqvAis6Lbv4U1s1dmuddZ78T4nBfrgym39A0uADC6FJqICJ1ktWy9JhPuSrqVd7Gl7537YCObLEZcJ+")
	i.Equal(decrypt, "radiation")
}
