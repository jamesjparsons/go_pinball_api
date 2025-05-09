# Go Pinball Server

A simple web server built with Go, using SQLite and GORM.

## Prerequisites

- Go 1.21 or later
- Git

## Setup

1. Clone the repository
2. Install dependencies:
   ```bash
   go mod tidy
   ```

## Running the Server

1. Start the server:
   ```bash
   go run main.go
   ```

2. The server will start on `http://localhost:8080`

## Database Configuration

The server uses SQLite, which creates a local database file named `pinball.db` in the project directory. No additional setup is required as SQLite is file-based and doesn't require a separate database server. # go_pinball_api
