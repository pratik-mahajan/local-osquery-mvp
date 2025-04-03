# Local OSQuery MVP

A Go service that collects system information using OSQuery, stores it in PostgreSQL, and provides both CLI and HTTP interfaces to access the data.

## Features

- Collects system information using OSQuery
- Stores data in PostgreSQL database
- Provides CLI interface with interactive menu
- Exposes HTTP API endpoint for data retrieval
- Supports multiple query types:
  - OS Version information
  - OSQuery version details
  - Installed applications

## Prerequisites

1. **Go**
2. **PostgreSQL**
3. **OSQuery** 

## Project Setup

1. Create a `.env` file in the project root:
```bash
# Database configuration
DB_HOST=localhost
DB_PORT=<your_db_port>
DB_USER=postgres
DB_PASSWORD=<your_password>
DB_NAME=postgres
```

2. Install Go dependencies:
```bash
go mod download
```

3. Run the application:
```
go run main.go
```

If you want to run start an HTTP server, you can use the `--server` and `--port` flags:

For example for port 1234:
```
go run main.go --server --port 1234
```

Now you can hit the `/latest_data` endpoint.

```
curl http://localhost:1234/latest_data
```

The application will start and display an interactive menu with the following options:
1. Get OS and OSQuery Info
2. Get Applications
3. Run All Queries
4. Exit

## HTTP API

When the HTTP server is running, the following endpoint is available:

### GET /latest_data
Returns the latest collected data including:
- OS version information
- OSQuery version details
- List of installed applications


## Project Structure

```
.
├── cmd/
│   ├── menu.go          # CLI menu implementation
│   └── utils.go         # Utility functions
├── pkg/
│   ├── api/             # HTTP API handlers
│   ├── config/          # Load config from .env file
│   ├── db/              # Database operations
│   ├── model/           # Data models and query definitions
├── .env                 # Environment configuration
├── go.mod               # Go module definition
├── go.sum               # Go module checksums
├── main.go              # Application entry point
└── README.md            # This file
```