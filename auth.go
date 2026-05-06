package main

import (
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"errors"
	"os"

	"golang.org/x/crypto/pbkdf2"
)

type AuthData struct {
	Salt string `json:"salt"`
	Hash string `json:"hash"`
}

var secretKey []byte

const authFile = "auth.json"

func generateSalt() []byte {
	salt := make([]byte, 16)
	rand.Read(salt)
	return salt
}

func deriveKey(password string, salt []byte) []byte {
	return pbkdf2.Key(
		[]byte(password),
		salt,
		100000,
		32,
		sha256.New,
	)
}

func createAuth(password string) error {
	salt := generateSalt()
	key := deriveKey(password, salt)

	auth := AuthData{
		Salt: base64.StdEncoding.EncodeToString(salt),
		Hash: base64.StdEncoding.EncodeToString(hashKey(key)),
	}

	file, err := os.Create(authFile)
	if err != nil {
		return err
	}
	defer file.Close()

	return json.NewEncoder(file).Encode(auth)
}

func initAuth(password string) error {
	if _, err := os.Stat("auth.json"); os.IsNotExist(err) {
		return createAuth(password)
	}

	return verifyPassword(password)
}

func verifyPassword(password string) error {
	file, err := os.Open(authFile)
	if err != nil {
		return errors.New("no auth file found, run setup first")
	}
	defer file.Close()

	var auth AuthData
	json.NewDecoder(file).Decode(&auth)

	salt, _ := base64.StdEncoding.DecodeString(auth.Salt)

	key := deriveKey(password, salt)

	if base64.StdEncoding.EncodeToString(hashKey(key)) != auth.Hash {
		return errors.New("invalid password")
	}

	secretKey = key
	return nil
}

func hashKey(key []byte) []byte {
	h := sha256.Sum256(key)
	return h[:]
}