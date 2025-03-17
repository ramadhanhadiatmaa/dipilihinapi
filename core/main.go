package main

import (
	"fmt"
	"log"
	"os"

	"core/models"
	"core/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {
	// Load environment variables dari file .env
	if err := godotenv.Load(/* "../.env" */); err != nil {
		log.Fatal("Error loading .env file:", err)
	}

	// Ambil kredensial database dari environment variables
	host := os.Getenv("DB_HOST_POS")
	user := os.Getenv("DB_USER_POS")
	password := os.Getenv("DB_PASSWORD_POS")
	dbname := os.Getenv("DB_NAME_POS")
	port := os.Getenv("DB_PORT_POS")

	// Buat DSN untuk koneksi ke PostgreSQL
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=disable", host, user, password, dbname, port)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatal("Gagal terhubung ke database:", err)
	}

	models.DB = db

	db.AutoMigrate(&models.Laptop{})

	app := fiber.New()

	app.Use(cors.New())

	routes.SetupRoutes(app)

	log.Fatal(app.Listen(":7304"))
}
