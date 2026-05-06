# Password Manager CLI (Go)

A simple command-line password manager built in Go. It allows users to store, search, delete, and manage credentials locally using JSON file storage.

## Features

- Add credentials (site, username, password)
- Search credentials by site
- Delete credentials
- View all stored credentials
- Count total entries
- Persistent storage using JSON file

## Technologies Used

- Go (Golang)
- File handling (JSON storage)
- Command-line interface (CLI)

## Project Structure
- main.go - CLI interface (user interaction)
- manager.go - Core logic (add, search, delete, list)
- storage.go - Save/load data from JSON file
- model.go - Credential data structure


## How to Run

### 1. Clone the repository
`git clone https://github.com/your-username/password-manager-go.git`
`cd password-manager-go`

### 2. Run the application
`go run .`

## Usage
When you run the program, you will see a menu:
# Password Manager CLI - Commands Reference

```add / add <site> <username> <psswd/blank>``` - Adds a new credential to the password manager.
Supports manual password entry or auto-generation by leaving blank.

```search <site>``` - Searches for credentials matching a site.

```update <site> <username>``` - Updates the password for a specific site and username.

```edit-user <site> <old-username>``` - Updates the username for an existing credential.

```delete <site> <username>``` - Deletes a specific credential matching both site and username.

```generate``` - Generates a secure random password and displays it to the user.

```exit``` - Saves all data and exits the application safely.

## Future Improvements
1) Copy to clipboard
2) Fuzzy matching/autocomplete search
3) Renaming sites
5) Export to CSV
6) Undo functionality
7) Multi-user support
