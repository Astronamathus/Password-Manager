package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"io"
	"math/big"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()"

func generatePassword(length int) string {
	password := make([]byte, length)

	for i := range password {
		num, _ := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		password[i] = charset[num.Int64()]
	}

	return string(password)
}

func encrypt(text string) string {
	if secretKey == nil {
		panic("secretKey not initialized - login required")
	}

	block, err := aes.NewCipher(secretKey)
	if err != nil {
		panic(err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		panic(err)
	}

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