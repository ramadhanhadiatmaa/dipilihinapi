package controllers

import (
	"auth/models"
	"context"
	"fmt"
	"os"
	"strconv"
	"time"

	"firebase.google.com/go/auth"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

var FirebaseAuth *auth.Client

func Login(c *fiber.Ctx) error {
	var data map[string]string

	// Ambil data dari request body
	if err := c.BodyParser(&data); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid input"})
	}

	// Ambil Firebase ID Token dari request
	idToken := data["id_token"]

	if idToken == "" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Token required"})
	}

	// **Verifikasi Firebase ID Token**
	token, err := FirebaseAuth.VerifyIDToken(context.Background(), idToken)
	if err != nil {
		fmt.Println("Invalid Firebase Token:", err)
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Invalid token"})
	}

	// Ambil email dari token yang sudah diverifikasi
	email := token.Claims["email"].(string)

	// Cari user berdasarkan email di database
	var user models.User
	err = models.DB.Where("email = ?", email).First(&user).Error

	// Jika user tidak ditemukan, tolak akses
	if err == gorm.ErrRecordNotFound {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "User not found"})
	} else if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Database error"})
	}

	// Generate JWT internal backend
	secretKey := os.Getenv("SECRET_KEY")
	claims := jwt.MapClaims{
		"email": user.Email,
		"type":  user.Type,
		"exp":   time.Now().Add(time.Hour * 240).Unix(),
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signedToken, _ := jwtToken.SignedString([]byte(secretKey))

	// Return user info + JWT backend
	return c.JSON(fiber.Map{
		"token":     signedToken,
		"email":     user.Email,
		"phone":     user.Phone,
		"full_name": user.FullName,
		"type":      user.Type,
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

	// Men-generate hash password
	password, err := bcrypt.GenerateFromPassword([]byte(data["password"]), bcrypt.DefaultCost)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not hash password"})
	}

	// Membuat objek user baru
	user := models.User{
		Email:     data["email"],
		Password:  string(password),
		FullName:  data["full_name"],
		Phone:     data["phone"],
		Type:      1,
		CreatedAt: time.Now(),
	}

	// Menyimpan user baru ke dalam database
	if err := models.DB.Create(&user).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not register user"})
	}

	// Mengambil informasi tipe user (TypeInfo) yang sudah terpreload
	var newUser models.User
	if err := models.DB.Preload("TypeInfo").First(&newUser, user.ID).Error; err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Could not retrieve user type information"})
	}
	return c.JSON(fiber.Map{
		"message": "User registered successfully",
	})
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