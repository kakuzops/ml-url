package config

import (
	"log"

	"github.com/kakuzops/ml-url/internal/domain"
	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) error {
	log.Println("Running database migrations...")

	var tableExists bool
	db.Raw("SELECT EXISTS (SELECT FROM information_schema.tables WHERE table_name = 'shorten_url')").Scan(&tableExists)
	if tableExists {
		log.Println("Table 'shorten_url' already exists")
	} else {
		log.Println("Table 'shorten_url' will be created")
	}

	err := db.AutoMigrate(&domain.URL{})
	if err != nil {
		log.Printf("Error during migration: %v", err)
		return err
	}

	var columns []string
	db.Raw("SELECT column_name FROM information_schema.columns WHERE table_name = 'shorten_url'").Pluck("column_name", &columns)
	log.Printf("Table columns: %v", columns)

	log.Println("Database migrations completed successfully")
	return nil
}
