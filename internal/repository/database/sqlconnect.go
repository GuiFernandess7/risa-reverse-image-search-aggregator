package database

import (
	"fmt"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func ConnectDB() (*gorm.DB, error) {
	fmt.Println("Trying to connect to database...")

	connectionString := os.Getenv("CONNECTION_STRING")

	db, err := gorm.Open(postgres.Open(connectionString), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("failed to connect to PostgreSQL: %w", err)
	}

	fmt.Println("Connected to PostgreSQL")
	return db, nil
}
