package main

import (
	"auth/models"
	"auth/routes"
	"context"
	"fmt"
	"os"
	"strings"

	"auth/controllers"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"google.golang.org/api/option"
)

var firebaseAuth *auth.Client

func getFirebaseAuth() (*auth.Client, error) {
	if firebaseAuth != nil {
		return firebaseAuth, nil
	}

	credentials := os.Getenv("FIREBASE_CREDENTIALS")
	if credentials == "" {
		return nil, fmt.Errorf("firebase credentials not set")
	}

	// Perbaiki escape karakter jika ada
	credentials = strings.ReplaceAll(credentials, "\\n", "\n")
	opt := option.WithCredentialsJSON([]byte(credentials))

	// Inisialisasi Firebase App
	app, err := firebase.NewApp(context.Background(), nil, opt)
	if err != nil {
		return nil, fmt.Errorf("failed to initialize firebase app: %v", err)
	}

	// Dapatkan Firebase Auth Client
	authClient, err := app.Auth(context.Background())
	if err != nil {
		return nil, fmt.Errorf("failed to initialize firebase auth: %v", err)
	}

	firebaseAuth = authClient
	return authClient, nil
}

func main() {

	models.ConnectDatabase()

	authClient, err := getFirebaseAuth()
	if err != nil {
		fmt.Println("Failed to initialize Firebase Auth:", err)
		os.Exit(1)
	}
	controllers.SetFirebaseAuth(authClient) // Kirim ke controller

	port := os.Getenv("PORT")
	if port == "" {
		port = "7301"
	}

	app := fiber.New()

	app.Use(cors.New())

	routes.Route(app)

	app.Listen(":" + port)
}
