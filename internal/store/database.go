package store

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func Open() (*sql.DB, error) {
	//Load env
	if err := godotenv.Load(); err != nil {
		fmt.Println("No .env file found, reading from system env")
	}

	//Getting db url
	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		return nil, fmt.Errorf("DATABASE_URL not set")
	}

	//Connecting the db
	db, err := sql.Open("pgx", dbURL)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}
	fmt.Println("Connected to db")

	//Ping the db
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	return db, nil
}
