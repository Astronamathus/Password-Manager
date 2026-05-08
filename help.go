package main

import "fmt"

func printHelp() {
	fmt.Println(Bold + Cyan + "\nPassword Manager Commands\n" + Reset)

	fmt.Println(Green + "add" + Reset)
	fmt.Println("  Add a new credential interactively\n")

	fmt.Println(Green + "add <site> <username> <password>" + Reset)
	fmt.Println("  Add credential directly from CLI\n")

	fmt.Println(Green + "add <site> <username> --generate" + Reset)
	fmt.Println("  Generate a secure password automatically\n")

	fmt.Println(Green + "search <site>" + Reset)
	fmt.Println("  Search credentials by site\n")

	fmt.Println(Green + "update <site> <username>" + Reset)
	fmt.Println("  Update password for an entry\n")

	fmt.Println(Green + "edit-user <site> <old-username>" + Reset)
	fmt.Println("  Change username for an existing entry\n")

	fmt.Println(Green + "delete <site>" + Reset)
	fmt.Println("  Delete credentials for a site\n")

	fmt.Println(Green + "delete <site> <username>" + Reset)
	fmt.Println("  Delete a specific credential\n")

	fmt.Println(Green + "list" + Reset)
	fmt.Println("  List all saved credentials\n")

	fmt.Println(Green + "generate" + Reset)
	fmt.Println("  Generate a secure random password\n")

	fmt.Println(Green + "copy <site>" + Reset)
	fmt.Println("  Copy password to clipboard\n")

	fmt.Println(Green + "help" + Reset)
	fmt.Println("  Show this help menu\n")

	fmt.Println(Green + "exit" + Reset)
	fmt.Println("  Save and exit the application\n")
}