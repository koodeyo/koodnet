package database

import (
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func getPostgresURL() string {
	db_url := os.Getenv("POSTGRES_URL")
	if db_url != "" {
		return db_url
	}

	db_hostname := os.Getenv("POSTGRES_HOST")
	db_name := os.Getenv("POSTGRES_DB")
	db_pass := os.Getenv("POSTGRES_PASSWORD")
	db_user := os.Getenv("POSTGRES_USER")
	db_port := os.Getenv("POSTGRES_PORT")

	if db_port == "" {
		db_port = "5432"
	}

	if db_hostname != "" && db_name != "" && db_user != "" {
		return fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable",
			db_hostname,
			db_user,
			db_pass,
			db_name,
			db_port,
		)
	}

	return ""
}

var Conn *gorm.DB

func Connect() {
	var err error

	dbURL := getPostgresURL()
	if dbURL != "" {
		for i := 1; i <= 3; i++ {
			Conn, err = gorm.Open(postgres.Open(dbURL), &gorm.Config{})
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
		Conn, err = gorm.Open(sqlite.Open("koodnet.db"), &gorm.Config{})
		if err != nil {
			log.Fatalf("Failed to initialize SQLite database: %v", err)
		}
	}
}
