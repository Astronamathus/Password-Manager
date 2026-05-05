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
main.go - CLI interface (user interaction)
manager.go - Core logic (add, search, delete, list)
storage.go - Save/load data from JSON file
model.go - Credential data structure


## How to Run

### 1. Clone the repository
git clone https://github.com/your-username/password-manager-go.git
cd password-manager-go

### 2. Run the application
go run .

## Usage
When you run the program, you will see a menu:

Add a new credential
View all credentials
Search by site
Delete a credential
View total count
Exit (auto-saves data)

## Future Improvements
Encrypt stored passwords
Improve CLI commands (e.g., add, search, delete)
Input validation
Export to CSV
Multi-user support
