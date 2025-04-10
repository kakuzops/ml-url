package main

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"github.com/kakuzops/ml-url/internal/config"
	"gorm.io/gorm"
)

func main() {
	log.Println("Starting migration process...")

	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	requiredEnvVars := []string{
		"POSTGRES_HOST",
		"POSTGRES_USER",
		"POSTGRES_PASSWORD",
		"POSTGRES_DB",
		"POSTGRES_PORT",
	}

	for _, envVar := range requiredEnvVars {
		if os.Getenv(envVar) == "" {
			log.Fatalf("Required environment variable %s is not set", envVar)
		}
		log.Printf("Environment variable %s is set", envVar)
	}

	db, err := config.NewDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	if err := checkDatabaseConnection(db); err != nil {
		log.Fatalf("Failed to verify database connection: %v", err)
	}

	if err := config.RunMigrations(db); err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	log.Println("Migrations completed successfully")
	os.Exit(0)
}

func checkDatabaseConnection(db *gorm.DB) error {
	sqlDB, err := db.DB()
	if err != nil {
		return err
	}

	sqlDB.SetConnMaxLifetime(time.Second * 5)

	if err := sqlDB.Ping(); err != nil {
		return err
	}

	log.Println("Successfully connected to the database")
	return nil
}
