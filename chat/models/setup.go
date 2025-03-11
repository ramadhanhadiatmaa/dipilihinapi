package models

import (
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func ConnectDatabase() {

	if err := godotenv.Load(/* "../.env" */); err != nil {
		log.Fatalf("Failed to load .env file: %v", err)
	}

	dsn := os.Getenv("DB_USER") + ":" +
		os.Getenv("DB_PASS") + "@tcp(" +
		os.Getenv("DB_HOST") + ":" +
		os.Getenv("DB_PORT") + ")/" +
		os.Getenv("DB_NAME") + "?charset=utf8mb4&parseTime=True&loc=Local"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Failed to get database instance: %v", err)
	}

	sqlDB.SetMaxOpenConns(20)              
	sqlDB.SetMaxIdleConns(10)              
	sqlDB.SetConnMaxIdleTime(5 * time.Minute) 
	sqlDB.SetConnMaxLifetime(30 * time.Minute)

	if err := db.AutoMigrate(&Chat{}); err != nil {
		log.Fatalf("Failed to migrate the database: %v", err)
	}

	DB = db
	log.Println("Successfully connected to the database.")
}