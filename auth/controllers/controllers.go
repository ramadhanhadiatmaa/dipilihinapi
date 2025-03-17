package controllers

import (
	"auth/models"
	"context"
	"os"
	"strconv"
	"time"

	"firebase.google.com/go/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"gorm.io/gorm"
)

var firebaseAuth *auth.Client

// SetFirebaseAuth menerima Firebase Auth Client dari main.go
func SetFirebaseAuth(authClient *auth.Client) {
	firebaseAuth = authClient
}

// Login mengautentikasi user menggunakan Firebase token
func Login(c *fiber.Ctx) error {
	var data map[string]string

	// Parse request body
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "invalid input"})
	}

	idToken := data["id_token"]
	email := data["email"]

	if idToken == "" || email == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "id token and email are required"})
	}

	// Verifikasi token Firebase
	token, err := firebaseAuth.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "invalid firebase token"})
	}

	// Pastikan email dari token cocok dengan yang dikirim
	if token.Claims["email"] != email {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "token does not match email"})
	}

	// Cari user di database
	var user models.User
	err = models.DB.Where("email = ?", email).First(&user).Error
	if err == gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "user not found"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "database error"})
	}

	// Generate JWT internal untuk backend
	secretKey := os.Getenv("SECRET_KEY")
	claims := jwt.MapClaims{
		"email": user.Email,
		"type":  user.Type,
		"exp":   time.Now().Add(time.Hour * 87600).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, err := jwtToken.SignedString([]byte(secretKey))
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "could not generate token"})
	}

	// Return user info + JWT backend
	return c.JSON(fiber.Map{
		"token":     signedToken,
		"email":     user.Email,
		"phone":     user.Phone,
		"full_name": user.FullName,
		"type":      user.Type,
		"id":        user.ID,
	})
}

func Register(c *fiber.Ctx) error {
	// Parsing input data
	var data map[string]string
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Mengecek apakah email sudah terdaftar
	var existingUser models.User
	if err := models.DB.First(&existingUser, "email = ?", data["email"]).Error; err == nil {
		return c.Status(fiber.StatusConflict).JSON(fiber.Map{"error": "Email already exists"})
	}

	// Membuat objek user baru
	user := models.User{
		Email:     data["email"],
		FullName:  data["full_name"],
		Phone:     data["phone"],
		Type:      1,
		CreatedAt: time.Now(),
	}

	// Menyimpan user baru ke dalam database
	if err := models.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not register user"})
	}

	return c.JSON(fiber.Map{
		"message": "User registered successfully",
	})
}

func Index(c *fiber.Ctx) error {
	id := c.Params("id")
	var data models.User

	if err := models.DB.First(&data, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return jsonResponse(c, fiber.StatusNotFound, "No data found", nil)
		}
		return jsonResponse(c, fiber.StatusInternalServerError, "Failed to load data", err.Error())
	}

	return c.JSON("Data found")
}

func Update(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return jsonResponse(c, fiber.StatusBadRequest, "Invalid ID format", nil)
	}

	var data models.User
	if err := models.DB.First(&data, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return jsonResponse(c, fiber.StatusNotFound, "No data found", nil)
		}
		return jsonResponse(c, fiber.StatusInternalServerError, "Failed to load data", err.Error())
	}

	var updateData models.User
	if err := c.BodyParser(&updateData); err != nil {
		return jsonResponse(c, fiber.StatusBadRequest, "Invalid", err.Error())
	}

	if updateData.ID != 0 && updateData.ID != id {
		if err := models.DB.First(&models.User{}, updateData.ID).Error; err == nil {
			return jsonResponse(c, fiber.StatusBadRequest, "The updated ID is already in use", nil)
		}
	}

	if err := models.DB.Model(&data).Updates(updateData).Error; err != nil {
		return jsonResponse(c, fiber.StatusInternalServerError, "Failed to update data", err.Error())
	}

	return jsonResponse(c, fiber.StatusOK, "Data successfully updated", nil)
}

func Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	if models.DB.Delete(&models.User{}, id).RowsAffected == 0 {
		return jsonResponse(c, fiber.StatusNotFound, "Data not found or already deleted", nil)
	}

	return jsonResponse(c, fiber.StatusOK, "Successfully deleted data", nil)
}

func jsonResponse(c *fiber.Ctx, status int, message string, data interface{}) error {
	return c.Status(status).JSON(fiber.Map{
		"message": message,
		"data":    data,
	})
}
