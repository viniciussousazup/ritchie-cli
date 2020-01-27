package cryptoutil

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/base64"
	"fmt"
)

var iv = []byte{35, 46, 57, 24, 85, 35, 24, 74, 87, 35, 88, 98, 66, 32, 14, 05}

func encodeBase64(b []byte) string {
	return base64.StdEncoding.EncodeToString(b)
}

func decodeBase64(s string) []byte {
	data, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		panic(err)
	}
	return data
}

// Encrypt encrypt data using aes
func Encrypt(key, text string) string {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}
	plaintext := []byte(fmt.Sprint(text))
	cfb := cipher.NewCFBEncrypter(block, iv)
	cipherText := make([]byte, len(plaintext))
	cfb.XORKeyStream(cipherText, plaintext)
	return encodeBase64(cipherText)
}

// Decrypt decrypt data using aes
func Decrypt(key, text string) string {
	block, err := aes.NewCipher([]byte(key))
	if err != nil {
		panic(err)
	}
	cipherText := decodeBase64(text)
	cfb := cipher.NewCFBDecrypter(block, iv)
	plaintext := make([]byte, len(cipherText))
	cfb.XORKeyStream(plaintext, cipherText)
	return string(plaintext)
}
