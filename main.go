package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
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

			if len(results) == 1 {
				r := results[0]
	
				fmt.Println("Username:", r.Username)
				fmt.Println("Password:", maskPassword(decrypt(r.Password)))
				fmt.Println("Reveal Password (y/n)? ")
				res, _ := reader.ReadString('\n')
				res = strings.TrimSpace(res)

				if strings.EqualFold(res, "y") {
					fmt.Println("Username:", r.Username)
					fmt.Println("Password:", decrypt(r.Password))
				} else {
					fmt.Println("Password masked")
				}	
					continue
			}

			for i, r := range results {
				fmt.Printf("%d. %s | %s\n", i+1, r.Username, maskPassword(decrypt(r.Password)))
			}

			fmt.Print("Select entry: ")
			choiceStr, _ := reader.ReadString('\n')
			choiceStr = strings.TrimSpace(choiceStr)
	
			choice, err := strconv.Atoi(choiceStr)
			if err != nil || choice < 1 || choice > len(results) {
				fmt.Println("Invalid choice")
				continue
			}

			selected := results[choice-1]

			fmt.Println("Username:", selected.Username)
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

		case "edit-user":
			if len(parts) < 3 {
				fmt.Println("Usage: edit-user <site> <old-username>")
				continue
			}

			site := parts[1]
			oldUser := parts[2]

			fmt.Print("New username: ")
			input, _ := reader.ReadString('\n')
			newUser := strings.TrimSpace(input)

			if newUser == "" {
				fmt.Println("Username cannot be empty")
				continue
			}

			if pm.UpdateUsername(site, oldUser, newUser) {
				pm.SaveToFile()
				fmt.Println("Username updated successfully")
			} else {
				fmt.Println("Entry not found")
			}


		case "delete":
			if len(parts) < 2 {
				fmt.Println("Usage: delete <site> [username]")
				continue
			}
	
			site := parts[1]

			if len(parts) >= 3 {
				user := parts[2]

				if pm.DeleteExact(site, user) {
					pm.SaveToFile()
					fmt.Println("Deleted successfully")
				} else {
					fmt.Println("Entry not found")
				}
				continue
			}

			results := pm.SearchAll(site)

			if len(results) == 0 {
				fmt.Println("No entries found")
				continue
			}

			if len(results) == 1 {
				r := results[0]
	
				fmt.Printf("Delete %s (%s)? (y/n): ", r.Site, r.Username)
				confirm, _ := reader.ReadString('\n')
				confirm = strings.TrimSpace(confirm)

				if strings.EqualFold(confirm, "y") {
					pm.DeleteExact(site, r.Username)
					pm.SaveToFile()
					fmt.Println("Deleted")
				} else {
					fmt.Println("Cancelled, Password saved")
				}
				continue
			}

			// psswd mangr
			for i, r := range results {
				fmt.Printf("%d. %s\n", i+1, r.Username)
			}

			fmt.Print("Select entry to delete: ")
			choiceStr, _ := reader.ReadString('\n')
			choiceStr = strings.TrimSpace(choiceStr)

			choice, err := strconv.Atoi(choiceStr)
			if err != nil || choice < 1 || choice > len(results) {
				fmt.Println("Invalid choice")
				continue
			}

			selected := results[choice-1]

			pm.DeleteExact(site, selected.Username)
			pm.SaveToFile()

			fmt.Println("Deleted successfully")
	
		case "list":
			grouped := pm.GroupBySite()

			sites := make([]string, 0, len(grouped))
			for site := range grouped {
				sites = append(sites, site)
			}
			sort.Strings(sites)

			for _, site := range sites {
				creds := grouped[site]

				if len(creds) == 1 {
					c := creds[0]
					fmt.Println(site, "|", c.Username, "|", maskPassword(decrypt(c.Password)))
					continue
				}

				fmt.Println(site + ":")

				for _, c := range creds {
					fmt.Printf("  - %s | %s\n", c.Username, maskPassword(decrypt(c.Password)))
				}

				fmt.Println()
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