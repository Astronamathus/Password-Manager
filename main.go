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
	masterPass, _ := reader.ReadString('\n')
	masterPass = strings.TrimSpace(masterPass)

	deriveKey(masterPass)

	for {
		fmt.Println("\n1. Add")
		fmt.Println("2. View All")
		fmt.Println("3. Search")
		fmt.Println("4. Total")
		fmt.Println("5. Delete")
		fmt.Println("6. Exit")
		fmt.Print("Choose: ")

		var choice int
		fmt.Scanln(&choice)

		switch choice {

		case 1:
			fmt.Print("Site: ")
			site, _ := reader.ReadString('\n')

			fmt.Print("Username: ")
			user, _ := reader.ReadString('\n')

			fmt.Print("Password: ")
			pass, _ := reader.ReadString('\n')

			pm.Add(Credential{
				Site:     strings.TrimSpace(site),
				Username: strings.TrimSpace(user),
				Password: encrypt(strings.TrimSpace(pass)),
			})

		case 2:
			for _, c := range pm.GetAll() {
				fmt.Println(c.Site, "|", c.Username, "|", maskPassword(decrypt(c.Password)))
			}

		case 3:
			fmt.Print("Enter site: ")
			site, _ := reader.ReadString('\n')

			result := pm.Search(strings.TrimSpace(site))

			if result != nil {
				fmt.Println("Username:", result.Username)
				fmt.Println("Password:", maskPassword(decrypt(result.Password)))
				fmt.Print("Reveal password? (y/n): ")
				choice, _ := reader.ReadString('\n')
				choice = strings.TrimSpace(choice)

				if choice == "y" || choice == "Y" {	
					fmt.Println("Actual Password:", decrypt(result.Password))
				} else {
					fmt.Println("Password masked")
				}
				
			} else {
				fmt.Println("Not found")
			}

		case 4:
			fmt.Println("Total:", pm.Total())

		case 5:
			fmt.Print("Enter site to delete: ")
			site, _ := reader.ReadString('\n')

			val := pm.Delete(strings.TrimSpace(site))
			if val {
				fmt.Println(site, " deleted from file")
			} else {
				fmt.Println(site, " not found on file")
			}
		case 6:
			pm.SaveToFile()
			fmt.Println("Saved. Goodbye.")
			return
		}
	}
}