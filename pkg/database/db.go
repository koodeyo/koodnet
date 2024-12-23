package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/koodeyo/koodnet/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func GetPostgresURL() string {
	var db_url string = os.Getenv("DATABASE_URL")

	if db_url != "" {
		return db_url
	}

	db_hostname := os.Getenv("POSTGRES_HOST")
	db_name := os.Getenv("POSTGRES_DB")
	db_user := os.Getenv("POSTGRES_USER")
	db_pass := os.Getenv("POSTGRES_PASSWORD")
	db_port := os.Getenv("POSTGRES_PORT")

	if db_hostname != "" && db_name != "" && db_user != "" && db_pass != "" && db_port != "" {
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s", db_user, db_pass, db_hostname, db_port, db_name)
	}
	return ""
}

func NewDatabase() *gorm.DB {
	var database *gorm.DB
	var err error

	dbURL := GetPostgresURL()
	if dbURL != "" {
		for i := 1; i <= 3; i++ {
			database, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{})
			if err == nil {
				break
			} else {
				log.Printf("Attempt %d: Failed to initialize PostgreSQL database. Retrying...", i)
				time.Sleep(3 * time.Second)
			}
		}
		if err != nil {
			log.Fatalf("Failed to connect to PostgreSQL after 3 attempts: %v", err)
		}
	} else {
		log.Println("PostgreSQL environment variables not defined. Falling back to SQLite.")
		database, err = gorm.Open(sqlite.Open("koodnet.db"), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to initialize SQLite database: %v", err)
		}
	}

	database.AutoMigrate(&models.Network{})

	return database
}
