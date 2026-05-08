package main

import (
	"bufio"
	"fmt"
	"os"
	"sort"
	"strconv"
	"strings"

	"golang.design/x/clipboard"
)

func main() {

	pm := PasswordManager{}
	pm.LoadFromFile()

	reader := bufio.NewReader(os.Stdin)

	fmt.Print(Cyan + "Enter master password: " + Reset)
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	err := initAuth(password)
	if err != nil {
		errorMsg("Auth failed: " + err.Error())
		return
	}

	if _, err := os.Stat("auth.json"); os.IsNotExist(err) {
		info("No user found. Creating new account...")

		err := createAuth(password)
		if err != nil {
			errorMsg("Error creating auth: " + err.Error())
			return
		}

		success("Account created successfully.")
	} else {
		err := verifyPassword(password)
		if err != nil {
			errorMsg("Login failed: " + err.Error())
			return
		}

		success("Login successful.")
	}

	err = clipboard.Init()
	if err != nil {
		errorMsg("Clipboard initialization failed")
		return
	}

	for {
		fmt.Print(Blue + "> " + Reset)

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
					errorMsg("Entry already exists. Use update instead.")
					continue
				}

				if len(parts) >= 4 {
					if parts[3] == "--generate" {
						pass = generatePassword(12)
						info("Generated password: " + pass)
					} else {
						pass = parts[3]
					}
				} else {
					fmt.Print("Password (leave empty to generate): ")
					input, _ := reader.ReadString('\n')
					pass = strings.TrimSpace(input)

					if pass == "" {
						pass = generatePassword(12)
						info("Generated password: " + pass)
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
					info("Generated password: " + pass)
				}
			}

			pm.Add(Credential{
				Site:     site,
				Username: user,
				Password: encrypt(pass),
			})

			pm.SaveToFile()
			success("Added successfully")

		case "search":
			if len(parts) < 2 {
				errorMsg("Usage: search <site>")
				continue
			}

			site := parts[1]
			results := pm.SearchAll(site)

			if len(results) == 0 {
				errorMsg("Not found")
				continue
			}

			if len(results) == 1 {
				r := results[0]

				fmt.Println(Cyan + "Username: " + Reset + r.Username)
				fmt.Println(Cyan + "Password: " + Reset + maskPassword(decrypt(r.Password)))

				fmt.Print("Reveal Password (y/n)? ")
				res, _ := reader.ReadString('\n')
				res = strings.TrimSpace(res)

				if strings.EqualFold(res, "y") {
					fmt.Println(decrypt(r.Password))
				} else {
					info("Password masked")
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
				errorMsg("Invalid choice")
				continue
			}

			selected := results[choice-1]

			fmt.Println(Cyan + "Username: " + Reset + selected.Username)

			fmt.Print("Reveal password? (y/n): ")
			res, _ := reader.ReadString('\n')
			res = strings.TrimSpace(res)

			if strings.EqualFold(res, "y") {
				fmt.Println(decrypt(selected.Password))
			} else {
				info("Password masked")
			}

		case "copy":
			if len(parts) < 2 {
				errorMsg("Usage: copy <site> [username]")
				continue
			}

			site := parts[1]
			results := pm.SearchAll(site)

			if len(results) == 0 {
				errorMsg("No entries found")
				continue
			}

			if len(results) == 1 {
				copyToClipboard(decrypt(results[0].Password))
				success("Password copied to clipboard")
				continue
			}

			for i, r := range results {
				fmt.Printf("%d. %s\n", i+1, r.Username)
			}

			fmt.Print("Select entry: ")
			choiceStr, _ := reader.ReadString('\n')
			choiceStr = strings.TrimSpace(choiceStr)

			choice, err := strconv.Atoi(choiceStr)
			if err != nil || choice < 1 || choice > len(results) {
				errorMsg("Invalid choice")
				continue
			}

			selected := results[choice-1]
			copyToClipboard(decrypt(selected.Password))

			success("Password copied to clipboard")

		case "help", "h", "?":
			printHelp()

		case "update":
			if len(parts) < 3 {
				errorMsg("Usage: update <site> <username>")
				continue
			}

			site := parts[1]
			user := parts[2]

			fmt.Print("New password (leave empty to generate): ")
			input, _ := reader.ReadString('\n')
			newPass := strings.TrimSpace(input)

			if newPass == "" {
				newPass = generatePassword(12)
				info("Generated password: " + newPass)
			}

			if pm.Update(site, user, encrypt(newPass)) {
				pm.SaveToFile()
				success("Updated successfully")
			} else {
				errorMsg("Entry not found")
			}

		case "edit-user":
			if len(parts) < 3 {
				errorMsg("Usage: edit-user <site> <old-username>")
				continue
			}

			site := parts[1]
			oldUser := parts[2]

			fmt.Print("New username: ")
			input, _ := reader.ReadString('\n')
			newUser := strings.TrimSpace(input)

			if newUser == "" {
				errorMsg("Username cannot be empty")
				continue
			}

			if pm.UpdateUsername(site, oldUser, newUser) {
				pm.SaveToFile()
				success("Username updated successfully")
			} else {
				errorMsg("Entry not found")
			}

		case "delete":
			if len(parts) < 2 {
				errorMsg("Usage: delete <site> [username]")
				continue
			}

			site := parts[1]

			if len(parts) >= 3 {
				user := parts[2]

				if pm.DeleteExact(site, user) {
					pm.SaveToFile()
					success("Deleted successfully")
				} else {
					errorMsg("Entry not found")
				}
				continue
			}

			results := pm.SearchAll(site)

			if len(results) == 0 {
				errorMsg("No entries found")
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
					success("Deleted")
				} else {
					info("Cancelled")
				}
				continue
			}

			for i, r := range results {
				fmt.Printf("%d. %s\n", i+1, r.Username)
			}

			fmt.Print("Select entry: ")
			choiceStr, _ := reader.ReadString('\n')
			choiceStr = strings.TrimSpace(choiceStr)

			choice, err := strconv.Atoi(choiceStr)
			if err != nil || choice < 1 || choice > len(results) {
				errorMsg("Invalid choice")
				continue
			}

			selected := results[choice-1]
			pm.DeleteExact(site, selected.Username)
			pm.SaveToFile()

			success("Deleted successfully")

		case "list":
			grouped := pm.GroupBySite()

			sites := make([]string, 0, len(grouped))
			for site := range grouped {
				sites = append(sites, site)
			}
			sort.Strings(sites)

			for _, site := range sites {
				creds := grouped[site]

				fmt.Println(Bold + Cyan + site + ":" + Reset)

				for _, c := range creds {
					fmt.Printf("  - %s | %s\n", c.Username, maskPassword(decrypt(c.Password)))
				}

				fmt.Println()
			}

		case "generate":
			pass := generatePassword(12)
			fmt.Println(Green + "Generated password: " + Reset + pass)

		case "exit":
			pm.SaveToFile()
			success("Saved. Goodbye.")
			return

		default:
			errorMsg("Unknown command. Type help")
		}
	}
}