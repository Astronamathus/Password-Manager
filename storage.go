package main

import (
	"encoding/json"
	"os"
)

func (pm *PasswordManager) SaveToFile() {
	file, err := os.Create("data.json")
	if err != nil {
		return
	}
	defer file.Close()

	json.NewEncoder(file).Encode(pm.credentials)
}

func (pm *PasswordManager) LoadFromFile() {
	file, err := os.Open("data.json")
	if err != nil {
		return // file doesn't exist yet
	}
	defer file.Close()

	json.NewDecoder(file).Decode(&pm.credentials)
}