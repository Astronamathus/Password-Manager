package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"

	"golang.org/x/crypto/pbkdf2"
)

var secretKey []byte

func deriveKey(password string) {
	salt := []byte("fixed-salt") // simple for now

	secretKey = pbkdf2.Key(
		[]byte(password),
		salt,
		100000,
		32,
		sha256.New,
	)

	fmt.Println("Key derived successfully")
}

func encrypt(text string) string {
	block, _ := aes.NewCipher(secretKey)

	gcm, _ := cipher.NewGCM(block)

	nonce := make([]byte, gcm.NonceSize())
	io.ReadFull(rand.Reader, nonce)

	ciphertext := gcm.Seal(nonce, nonce, []byte(text), nil)

	return base64.StdEncoding.EncodeToString(ciphertext)
}

func decrypt(encText string) string {
	data, _ := base64.StdEncoding.DecodeString(encText)

	block, _ := aes.NewCipher(secretKey)
	gcm, _ := cipher.NewGCM(block)

	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	plaintext, _ := gcm.Open(nil, nonce, ciphertext, nil)

	return string(plaintext)
}