package main

import (
	"auth/models"
	"auth/routes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

func loadEnv() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Error loading .env file")
	}
}

func getFirebaseAuth() (*auth.Client, error) {
	// Load environment variables
	loadEnv()
	credentials := os.Getenv("FIREBASE_CREDENTIALS")

	if credentials == "" {
		return nil, fmt.Errorf("firebase_credentials is not set")
	}

	// Convert JSON string to map
	var credMap map[string]interface{}
	if err := json.Unmarshal([]byte(credentials), &credMap); err != nil {
		return nil, fmt.Errorf("invalid firebase_credentials JSON format")
	}

	// Convert map to Firebase option
	opt := option.WithCredentialsJSON([]byte(credentials))

	// Initialize Firebase App
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase app: %v", err)
	}

	// Get Firebase Auth Client
	authClient, err := app.Auth(context.Background())
	if err != nil {
		return nil, fmt.Errorf("error initializing firebase auth: %v", err)
	}

	return authClient, nil
}

func main() {

	authClient, err := getFirebaseAuth()
	if err != nil {
		log.Fatalf("Failed to initialize Firebase Auth: %v", err)
	}
	fmt.Println("Firebase Auth initialized successfully:", authClient)

	models.ConnectDatabase()

	port := os.Getenv("PORT")
	if port == "" {
		port = "7301"
	}

	app := fiber.New()

	app.Use(cors.New())

	routes.Route(app)

	app.Listen(":" + port)
}
