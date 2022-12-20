package main

import (
	"fmt"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/kyungmun/otp-server/controller"
	"github.com/kyungmun/otp-server/models"
	"github.com/kyungmun/otp-server/repository"
	"github.com/kyungmun/otp-server/service"
	"gorm.io/gorm"
)

func main() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal(err)
	}

	config := &repository.Config{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		Password: os.Getenv("DB_PASS"),
		User:     os.Getenv("DB_USER"),
		SSLMode:  os.Getenv("DB_SSLMODE"),
		DBName:   os.Getenv("DB_NAME"),
		DBEngine: os.Getenv("DB_ENGINE"),
	}

	fmt.Printf("config : %v\n", config)

	//db connect
	var db *gorm.DB
	if config.DBEngine == "sqlite" {
		db, err = repository.ConnectSqliteDB(config)
	} else {
		db, err = repository.ConnectMysqlDB(config)
	}

	if err != nil {
		log.Fatal("could not load database")
	}

	err = models.MigrateOtpRegistrys(db)
	if err != nil {
		log.Fatal("could not migrate db")
	}

	//레토지토리 생성하면서 db 연결
	otpRepo, err := repository.NewOtpRepository(db)
	if err != nil {
		log.Fatal("could not otp repository create")
	}

	//서비스에 레포지토리 연결
	otpService, err := service.NewOtpServices(otpRepo)
	if err != nil {
		log.Fatal("could not otp services create")
	}

	//fiber controller engine 생성하고 서비스 연결
	fiberApp := controller.NewFiber()
	fiberApp.SetupRoutes(otpService)

	//middleware test
	fiberApp.App.Use(func(c *fiber.Ctx) error {
		fmt.Println("fiber middleware")
		return c.Next()
	})

	// Render index template
	fiberApp.App.Get("/", func(c *fiber.Ctx) error {
		return c.Render("index", fiber.Map{
			"Title": "Google OTP Server",
			"Name":  "by Kyungmun, lim",
		})
	})

	fiberApp.Listen(":8082")
}
