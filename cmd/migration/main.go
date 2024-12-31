package main

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

const configPath = "internal/database/migrations"

func main() {
	if err := godotenv.Load(); err != nil {
		fmt.Println("Error loading .env file:", err)
		return
	}

	if len(os.Args) < 2 {
		fmt.Println("Missing command. Usage: create migration <name> or migrate <direction>")
		return
	}

	command := os.Args[1]
	switch command {
	case "create":
		if len(os.Args) < 4 || os.Args[2] != "migration" {
			fmt.Println("Usage: create migration <name>")
			return
		}
		createMigration(os.Args[3])
	case "migrate":
		if len(os.Args) < 3 {
			fmt.Println("Usage: migrate <direction>")
			return
		}
		migrateDB(os.Args[2])
	case "migrate:fresh":
		dropDB()
		migrateDB("up")
	default:
		fmt.Println("Unknown command. Usage: create migration <name> or migrate <direction>")
	}
}

func dropDB() {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("MIGRATION_DB_HOST"), os.Getenv("DB_NAME"))

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("Error opening database connection:", err)
		return
	}
	defer db.Close()

	if _, err := db.Exec("DROP SCHEMA public CASCADE; CREATE SCHEMA public;"); err != nil {
		fmt.Println("Error dropping database:", err)
		return
	}
	fmt.Println("Database dropped successfully")
}

func createMigration(name string) {
	timestamp := time.Now().Format("20060102150405")
	upMigrationPath := fmt.Sprintf("%s/%s_%s.up.sql", configPath, timestamp, name)
	downMigrationPath := fmt.Sprintf("%s/%s_%s.down.sql", configPath, timestamp, name)

	if err := createFile(upMigrationPath); err != nil {
		fmt.Println("Error creating up migration file:", err)
		return
	}
	fmt.Println("Up migration file created successfully at:", upMigrationPath)

	if err := createFile(downMigrationPath); err != nil {
		fmt.Println("Error creating down migration file:", err)
		return
	}
	fmt.Println("Down migration file created successfully at:", downMigrationPath)
}

func createFile(path string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil
}

func migrateDB(direction string) {
	dbURL := fmt.Sprintf("postgres://%s:%s@%s:5432/%s?sslmode=disable",
		os.Getenv("DB_USER"), os.Getenv("DB_PASSWORD"), os.Getenv("MIGRATION_DB_HOST"), os.Getenv("DB_NAME"))

	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		fmt.Println("Error opening database connection:", err)
		return
	}
	defer db.Close()

	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		fmt.Println("Error creating driver instance:", err)
		return
	}

	m, err := migrate.NewWithDatabaseInstance("file://"+configPath, "postgres", driver)
	if err != nil {
		fmt.Println("Error creating migration instance:", err)
		return
	}

	switch direction {
	case "up":
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			fmt.Println("Error migrating up:", err)
			return
		}
		fmt.Println("Migration up complete")
	case "down":
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			fmt.Println("Error migrating down:", err)
			return
		}
		fmt.Println("Migration down complete")
	default:
		fmt.Println("Unknown migration direction. Use 'up' or 'down'")
	}
}
