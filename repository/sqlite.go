package repository

import (
	"log"
	"os"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"

	"github.com/kyungmun/otp-server/models"
)

func ConnectSqliteDB(config *Config) (*gorm.DB, error) {

	db, err := gorm.Open(sqlite.Open("otp.dbb"), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
		os.Exit(2)
	}

	err = db.AutoMigrate(&models.OtpRegistry{},
		&models.User{})

	if err != nil {
		log.Fatal("could not migrate db")
	}

	log.Println("connected sqlite db")
	return db, nil
}
