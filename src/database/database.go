package database

import (
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Setup() {
	dsn := os.Getenv("DSN")
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("Error:", err.Error())
		return
	}

	// db.AutoMigrate(&models.User{})
	// db.AutoMigrate(&models.Income{})
	// db.AutoMigrate(&models.Expense{})

	DB = db
	return
}

func GetDB() *gorm.DB {
	return DB
}
