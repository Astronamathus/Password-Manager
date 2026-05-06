package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
)

var secretKey = []byte("12345678901234567890123456789012") // 32 bytes

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