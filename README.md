# README.md

## Go Application with Database Migration and Development Workflow

---

## Project Structure

```plaintext
.
├── README.md
├── cmd
│   ├── api
│   │   └── main.go          # Entry point for the API server
│   └── migration
│       └── main.go          # Entry point for the migration tool
├── docker-compose.yml       # Docker setup
├── go.mod                   # Go modules file
├── go.sum                   # Go dependencies checksum
├── internal
│   ├── configs              # Configuration files
│   │   └── configs.go
│   ├── database
│   │   ├── migrations       # Migration files
│   │   │   ├── *.up.sql     # SQL files for "up" migrations
│   │   │   └── *.down.sql   # SQL files for "down" migrations
│   │   ├── mysql.go         # MySQL-specific utilities
│   │   └── postgresql.go    # PostgreSQL-specific utilities
│   ├── modules              # Application modules
│   │   ├── auth             # Authentication module
│   │   ├── common           # Common utilities
│   │   └── file-uploads     # File upload module
│   └── pkg
│       └── e                # Error handling utilities
│           └── error.go
```

---

## Prerequisites

1. **Go Environment**: Install Go 1.23+.
2. **Dependencies**:
   - Install required packages using `go mod tidy`.
3. **Database**:
   - Ensure your database (e.g., MySQL, PostgreSQL) is running and configured.
4. **air**: Install `air` for hot-reloading during development:

   ```bash
   go install github.com/air-verse/air@latest
   ```

---

## Development Workflow

### Create environment file

```bash
cp .env.example .env
```

### Run the Application in Development Mode

Use `air` to start the application with hot-reloading:

```bash
air
```

---

## Database Migration

### Setup Migration with Makefile

The `Makefile` includes the following migration commands:

```makefile
# Migrate database
.PHONY: migrate
migrate:
	go run ./cmd/migration/main.go migrate up

# Migrate Fresh database
.PHONY: migrate-fresh
migrate-fresh:
	go run ./cmd/migration/main.go migrate:fresh

# Rollback database
.PHONY: rollback
rollback:
	go run ./cmd/migration/main.go migrate down

# Create migration
.PHONY: create-migration
create-migration:
	go run ./cmd/migration/main.go create migration ${name}
```

---

### Usage

1. **Create a New Migration**  
   Generate new migration files with a descriptive name:

   ```bash
   make create-migration name=<migration_name>
   ```

   Example:

   ```bash
   make create-migration name=add_users_table
   ```

2. **Run Migrations**  
   Apply all pending migrations:

   ```bash
   make migrate
   ```

3. **Fresh Migrations**  
   Drop and recreate the database, then reapply all migrations:

   ```bash
   make migrate-fresh
   ```

4. **Rollback Migrations**  
   Roll back the last applied migration:

   ```bash
   make rollback
   ```

---

## Configuration

### Database Configuration

Ensure your database configuration is defined in `internal/configs/configs.go`:

```go
package configs

type Config struct {
    DatabaseURL  string
    DatabaseName string
}

var Config = &Config{
    DatabaseURL:  "user:password@tcp(localhost:3306)/your_db_name",
    DatabaseName: "your_db_name",
}

func Setup() error {
    // Add setup logic if needed
    return nil
}
```

---

## Notes

- **Environment Variables**: Use environment variables to manage sensitive information like database credentials.
- **Air Configuration**: Customize your `.air.toml` for hot-reloading as needed.
- **Database Permissions**: Ensure your database user has permissions for creating, dropping, and modifying schemas.

