package controllers

import (
	"chat/models"
	"fmt"
	"time"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func Create(c *fiber.Ctx) error {
	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		return jsonResponse(c, fiber.StatusBadRequest, "Invalid input", err.Error())
	}

	allowedKeys := []string{"email", "message"}
	for key := range data {
		if !contains(allowedKeys, key) {
			return jsonResponse(c, fiber.StatusBadRequest, "Inputting data is not allowed", nil)
		}
	}

	// Pastikan kedua key "email" dan "message" ada dalam data input
	email, emailExists := data["email"]
	message, messageExists := data["message"]
	if !emailExists || !messageExists {
		return jsonResponse(c, fiber.StatusBadRequest, "'email' and 'message' keys are required", nil)
	}

	// Buat object Chat dengan menyimpan email dan message pada field yang sesuai
	datas := models.Chat{
		Email:     fmt.Sprintf("%v", email),
		Message:   fmt.Sprintf("%v", message),
		CreatedAt: time.Now(),
		Status:    1,
	}

	// Simpan ke database
	if err := models.DB.Create(&datas).Error; err != nil {
		return jsonResponse(c, fiber.StatusInternalServerError, "Failed to save data", err.Error())
	}

	// Return response sukses
	return jsonResponse(c, fiber.StatusCreated, "Data successfully added", datas)
}

func Show(c *fiber.Ctx) error {
	var data []models.Chat

	if err := models.DB.Find(&data).Error; err != nil {
		return jsonResponse(c, fiber.StatusInternalServerError, "Failed to load data", err.Error())
	}

	if len(data) == 0 {
		return jsonResponse(c, fiber.StatusNotFound, "No data found", nil)
	}

	return c.JSON(data)
}

func Index(c *fiber.Ctx) error {
	id := c.Params("id")
	var data models.Chat

	if err := models.DB.First(&data, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return jsonResponse(c, fiber.StatusNotFound, "No data found", nil)
		}
		return jsonResponse(c, fiber.StatusInternalServerError, "Failed to load data", err.Error())
	}

	return c.JSON(data)
}

func Delete(c *fiber.Ctx) error {
	id := c.Params("id")

	if models.DB.Delete(&models.Chat{}, id).RowsAffected == 0 {
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

func contains(slice []string, item string) bool {
	for _, a := range slice {
		if a == item {
			return true
		}
	}
	return false
}