package database

import (
	"fmt"
	"log"

	"auth.jwt.api/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

const (
	host     = "0.0.0.0"
	port     = "5432"
	user     = "auth"
	password = "password"
	database = "auth"
	timezone = "Asia/Bangkok"
)

func Setup() {
	dns := fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable TimeZone=%s", user, password, host, port, database, timezone)
	db, err := gorm.Open(postgres.Open(dns), &gorm.Config{})
	if err != nil {
		log.Fatal("Connect into database failed: ", err.Error())
	}

	db.AutoMigrate(&models.User{})

	fmt.Println("Connect into data success")

	DB = db
}
