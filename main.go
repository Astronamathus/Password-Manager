package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func main() {
	
	pm := PasswordManager{}
	pm.LoadFromFile()

	reader := bufio.NewReader(os.Stdin)

	fmt.Print("Enter master password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	err := initAuth(password)
	if err != nil {
		fmt.Println("Auth failed:", err)
		return
	} // Check if auth file exists
	if _, err := os.Stat("auth.json"); os.IsNotExist(err) {
		fmt.Println("No user found. Creating new account...")

		err := createAuth(password)
		if err != nil {
			fmt.Println("Error creating auth:", err)
			return
		}

	fmt.Println("Account created successfully.")
	} else {
		err := verifyPassword(password)
		if err != nil {
			fmt.Println("Login failed:", err)
			return
		}

		fmt.Println("Login successful.")
	}

	for {
	fmt.Print("> ")

	line, _ := reader.ReadString('\n')
	line = strings.TrimSpace(line)

	parts := strings.Split(line, " ")
	command := parts[0]

	switch command {

		case "add":
		if len(parts) < 4 {
			fmt.Println("Usage: add <site> <username> <password>")
			continue
		}

		site := parts[1]
		user := parts[2]
		pass := parts[3]

		pm.Add(Credential{
			Site:     site,
			Username: user,
			Password: encrypt(pass),
		})

		pm.SaveToFile()
		fmt.Println("Added")

		case "search":
			if len(parts) < 2 {
				fmt.Println("Usage: search <site>")
				continue
			}

			site := parts[1]
			result := pm.Search(site)

			if result != nil {
				fmt.Println("Username:", result.Username)
				fmt.Println("Password:", maskPassword(result.Password))

				fmt.Print("Reveal? (y/n): ")
				choice, _ := reader.ReadString('\n')

				if strings.TrimSpace(choice) == "y" {
					fmt.Println("Password:", decrypt(result.Password))
				} else {
					fmt.Println("Password Masked")
				}
				} else {
					fmt.Println("Not found")
				}

		case "delete":
			if len(parts) < 2 {
				fmt.Println("Usage: delete <site>")
				continue
			}

			if pm.Delete(parts[1]) {
				pm.SaveToFile()
				fmt.Println("Deleted")
			} else {
				fmt.Println("Not found")
			}
	
		case "list":
			for _, c := range pm.GetAll() {
				fmt.Println(c.Site, "|", c.Username, "|", maskPassword(decrypt(c.Password)))
			}

		case "exit":
			pm.SaveToFile()
			fmt.Println("Saved. Goodbye.")
			return
		default:
			fmt.Println("Unknown command. Try: add, search, delete, list, exit")
		}
	}
}