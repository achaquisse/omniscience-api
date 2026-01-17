---
description: Repository Information Overview
alwaysApply: true
---

# Omniscience API Information

## Summary
Omniscience API (module name: skulla-api) is a REST API built with Go and Fiber framework for managing educational data including courses, classes, students, assiduity matrices, and presence/absence records.

## Structure
```
omniscience-api/
├── src/
│   ├── db/              # Database connection and models
│   │   ├── connection.go          # Database connection with GORM
│   │   ├── student_class.go       # Student class model
│   │   └── student_enrollment.go  # Student enrollment model
│   ├── rest/            # REST API endpoints
│   │   ├── endpoints.go           # Endpoint routing
│   │   └── student_class.go       # Student class handlers
│   ├── main.go          # Application entry point
│   ├── go.mod           # Go module dependencies
│   └── go.sum           # Dependency checksums
└── .dev-db/
    ├── docker-compose.yml  # MariaDB development database
    └── init.sql            # Database initialization script
```

## Language & Runtime
**Language**: Go  
**Version**: 1.25  
**Build System**: Go modules  
**Package Manager**: Go modules

## Dependencies
**Main Dependencies**:
- `github.com/gofiber/fiber/v2 v2.52.10` - Web framework
- `gorm.io/gorm v1.31.1` - ORM library
- `gorm.io/driver/mysql v1.6.0` - MySQL/MariaDB driver

**Key Indirect Dependencies**:
- `github.com/go-sql-driver/mysql v1.9.3` - MySQL driver
- `github.com/google/uuid v1.6.0` - UUID generation
- `github.com/valyala/fasthttp v1.69.0` - HTTP implementation

## Build & Installation
```bash
# Install dependencies
cd src
go mod download

# Build the application
go build -o omniscience-api

# Run the application
go run main.go

# The API will start on port 8080
```

## Configuration
The application uses environment variables for database configuration with the following defaults:

**Environment Variables**:
- `DB_HOST` - Database host (default: `localhost`)
- `DB_PORT` - Database port (default: `3306`)
- `DB_NAME` - Database name (default: `omniscience`)
- `DB_USERNAME` - Database username (default: `admin`)
- `DB_PASSWORD` - Database password (default: `admin`)

## Docker
**Docker Compose**: `.dev-db/docker-compose.yml`  
**Service**: `omniscience-db`  
**Image**: `mariadb:latest`  
**Configuration**:
- Database: `omniscience`
- Root/User Password: `admin`
- Port: `3306:3306`
- Initialization script: `init.sql` mounted to container
- Persistent volume: `mariadb_data`

**Starting Development Database**:
```bash
cd .dev-db
docker-compose up -d
```

## Main Files
**Entry Point**: `src/main.go:10` - Main application entry point that initializes database connection and Fiber web server

**Core Components**:
- `src/db/connection.go` - Database connection handler using GORM with environment-based configuration
- `src/rest/endpoints.go` - REST API endpoint initialization and routing

**API Endpoints**:
- `GET /student-class` - List student classes

## Testing
No testing framework or test files are currently configured in this repository.
