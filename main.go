package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
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
			var site, user, pass string

			if len(parts) >= 3 {
				site = parts[1]
				user = parts[2]

				if pm.Exists(site, user) {
					fmt.Println("Entry already exists. Use update command instead.")
					return
				}	
				// Handle password or --generate
				if len(parts) >= 4 {
					if parts[3] == "--generate" {
						pass = generatePassword(12)
						fmt.Println("Generated password:", pass)
					} else {
						pass = parts[3]
					}
				} else {
					// fallback to interactive password
					fmt.Print("Password (leave empty to generate): ")
					input, _ := reader.ReadString('\n')
					pass = strings.TrimSpace(input)

					if pass == "" {
						pass = generatePassword(12)
						fmt.Println("Generated password:", pass)
					}
				}

			} else {
				fmt.Print("Site: ")
				s, _ := reader.ReadString('\n')

				fmt.Print("Username: ")
				u, _ := reader.ReadString('\n')

				fmt.Print("Password (leave empty to generate): ")
				p, _ := reader.ReadString('\n')

				site = strings.TrimSpace(s)
				user = strings.TrimSpace(u)
				pass = strings.TrimSpace(p)

				if pass == "" {
					pass = generatePassword(12)
					fmt.Println("Generated password:", pass)
				}
			}

				// Save credential
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
			results := pm.SearchAll(site)

			if len(results) == 0 {
				fmt.Println("Not found")
				continue
			}

			for i, r := range results {
				fmt.Printf("%d. %s | %s\n", i+1, r.Username, maskPassword(decrypt(r.Password)))
			}

			fmt.Print("Select entry to reveal: ")
			choiceStr, _ := reader.ReadString('\n')
			choiceStr = strings.TrimSpace(choiceStr)

			choice, err := strconv.Atoi(choiceStr)
			if err != nil || choice < 1 || choice > len(results) {
				fmt.Println("Invalid choice")
				continue
			}

			selected := results[choice-1]
			fmt.Println("Username: ", selected.Username)
			fmt.Println("Password:", decrypt(selected.Password))
		
		case "update":
			if len(parts) < 3 {
				fmt.Println("Usage: update <site> <username>")
				continue
			}

			site := parts[1]
			user := parts[2]

			fmt.Print("New password (leave empty to generate): ")
			input, _ := reader.ReadString('\n')
			newPass := strings.TrimSpace(input)

			if newPass == "" {
				newPass = generatePassword(12)
				fmt.Println("Generated password:", newPass)
			}

			if pm.Update(site, user, encrypt(newPass)) {
				pm.SaveToFile()
				fmt.Println("Updated successfully")
			} else {
				fmt.Println("Entry not found")
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

		case "generate":
			pass := generatePassword(12)
			fmt.Println("Generated password:", pass)

		case "exit":
			pm.SaveToFile()
			fmt.Println("Saved. Goodbye.")
			return
		default:
			fmt.Println("Unknown command. Try: add, search, delete, list, exit")
		}
	}
}