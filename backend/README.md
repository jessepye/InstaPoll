# InstaPoll Backend

This is the backend service for InstaPoll, built with Go and Gin.

## Getting Started

1. Install Go 1.21 or later
2. Install dependencies:
   ```bash
   go mod download
   ```
3. Run the server:
   ```bash
   go run main.go
   ```

The server will start on port 8080. You can test it by visiting:
- http://localhost:8080/

## Development

- The server uses Gin framework for routing and middleware
- Main application logic is in `main.go`
- Additional packages and routes can be added as needed 