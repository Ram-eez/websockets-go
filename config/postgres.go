package config

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func ConnectDB() (*sql.DB, error) {
	if err := godotenv.Load(); err != nil {
		log.Println("no .env found")
	}
	connStr := fmt.Sprintf(
		"user=%s password=%s dbname=%s sslmode=%s",
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("DB_SSLMODE"),
	)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opnening the DB: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to the DB: %w", err)
	}
	fmt.Println("Database conn established successfully")
	return db, nil
}
