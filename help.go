package main

import "fmt"

func printHelp() {
	fmt.Println(`
Available Commands

add
  Add a new credential interactively

add <site> <username> <password>
  Add a credential directly from CLI

add <site> <username> --generate
  Generate a secure password automatically

search <site>
  Search credentials by site

update <site> <username>
  Update password for an entry

edit-user <site> <old-username>
  Change username for an existing entry

delete <site>
  Delete credential(s) for a site

delete <site> <username>
  Delete a specific credential

list
  List all saved credentials

generate
  Generate a secure random password

copy <site>
  Copy password to clipboard

copy <site> <username>
  Copy specific account password to clipboard

help
  Show this help menu

exit
  Save and exit the application
`)
}