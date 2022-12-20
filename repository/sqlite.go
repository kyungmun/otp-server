package repository

import (
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectSqliteDB(config *Config) (*gorm.DB, error) {

	db, err := gorm.Open(sqlite.Open("otp.db"), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
		os.Exit(2)
	}

	log.Println("connected sqlite db")
	return db, nil
}
