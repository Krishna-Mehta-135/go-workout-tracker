package store

import (
	"database/sql"
	"fmt"
	"io/fs"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib" // Import pgx as a database/sql driver
	"github.com/joho/godotenv"         // Loads .env files into environment variables
	"github.com/pressly/goose/v3"      // Database migration tool
)

//
// Open establishes and verifies a connection to the PostgreSQL database.
//
// Steps:
//  1. Loads environment variables from `.env` (if present).
//  2. Reads the DATABASE_URL from environment variables.
//  3. Opens a connection using the pgx driver.
//  4. Pings the database to verify it's reachable.
//
// Returns:
//  - *sql.DB connection object (caller must close it)
//  - error if connection fails
//
func Open() (*sql.DB, error) {
	// Load .env file if available (falls back to system env)
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, reading from system env")
	}

	// Get database URL from env
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL not set")
	}

	// Connect using pgx driver
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	fmt.Println("Connected to db")

	// Verify the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}

//
// MigrateFS runs database migrations from an embedded filesystem (fs.FS).
//
// This is useful when migration files are embedded in the Go binary (using `embed`).
// The function temporarily sets goose's base filesystem, runs migrations,
// and then restores the default filesystem.
//
func MigrateFS(db *sql.DB, migrationsFS fs.FS, dir string) error {
	goose.SetBaseFS(migrationsFS) // Tell goose to look inside the embedded FS
	defer goose.SetBaseFS(nil)    // Reset after migrations

	return Migrate(db, dir)
}

//
// Migrate applies all pending migrations from a given directory path.
//
// Steps:
//  1. Sets goose dialect to "postgres" so it knows SQL syntax.
//  2. Runs `goose.Up` which applies all new migration files in the folder.
//
// Example migration directory structure:
//  migrations/
//    20230814120000_create_users_table.sql
//    20230814121000_add_workouts_table.sql
//
func Migrate(db *sql.DB, dir string) error {
	// Tell goose to use PostgreSQL-specific migration behavior
	if err := goose.SetDialect("postgres"); err != nil {
		return fmt.Errorf("migration error - %w", err)
	}

	// Apply all "up" migrations
	if err := goose.Up(db, dir); err != nil {
		return fmt.Errorf("goose up error: %w", err)
	}

	return nil
}
