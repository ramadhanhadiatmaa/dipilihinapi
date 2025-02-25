package main

import (
	"auth/models"
	"auth/controllers"
	"auth/routes"
	"context"
	"encoding/base64"
	"fmt"
	"log"
	"os"

	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"google.golang.org/api/option"
)

var firebaseAuth *auth.Client

// getFirebaseAuth menginisialisasi dan mengembalikan Firebase Auth Client
func getFirebaseAuth() (*auth.Client, error) {
	if firebaseAuth != nil {
		return firebaseAuth, nil
	}

	// Baca string Base64 kredensial dari variabel environment
	encodedCredentials := os.Getenv("FIREBASE_CREDENTIALS_BASE64")
	if encodedCredentials == "" {
		return nil, fmt.Errorf("firebase credentials not set")
	}

	// Dekode Base64 untuk mendapatkan JSON asli kredensial
	decoded, err := base64.StdEncoding.DecodeString(encodedCredentials)
	if err != nil {
		return nil, fmt.Errorf("failed to decode firebase credentials: %v", err)
	}

	// Buat opsi untuk Firebase Admin SDK
	opt := option.WithCredentialsJSON(decoded)

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
		log.Fatalf("Failed to initialize Firebase Auth: %v", err)
	}
	log.Println("Firebase Auth initialized successfully:", authClient)

	// Kirim Firebase Auth Client ke controllers
	controllers.SetFirebaseAuth(authClient)

	port := os.Getenv("PORT")
	if port == "" {
		port = "7301"
	}

	app := fiber.New()

	app.Use(cors.New())

	routes.Route(app)

	app.Listen(":" + port)
}
