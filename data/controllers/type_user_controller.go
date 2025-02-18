package controllers

import (
	"data/models"
	"fmt"

	"github.com/gofiber/fiber/v2"
)

/* func CreateType(c *fiber.Ctx) error {
	var data models.TypeUser

	if err := c.BodyParser(&data); err != nil {
		return jsonResponse(c, fiber.StatusBadRequest, "Invalid input", err.Error())
	}

	if err := models.DB.Create(&data).Error; err != nil {
		return jsonResponse(c, fiber.StatusInternalServerError, "Failed to save data", err.Error())
	}

	return jsonResponse(c, fiber.StatusCreated, "Data successfully added", data)
} */

func CreateType(c *fiber.Ctx) error {
	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		return jsonResponse(c, fiber.StatusBadRequest, "Invalid input", err.Error())
	}

	allowedKeys := []string{"type"}

	for key := range data {
		if !contains(allowedKeys, key) {
			return jsonResponse(c, fiber.StatusBadRequest, "Inputting data is not allowed", nil)
		}
	}

	if exampleValue, exists := data["type"]; exists {
		typeUser := models.TypeUser{
			Type: fmt.Sprintf("%v", exampleValue), // Menyimpan value yang diterima dalam Type
		}

		// Simpan ke database
		if err := models.DB.Create(&typeUser).Error; err != nil {
			return jsonResponse(c, fiber.StatusInternalServerError, "Failed to save data", err.Error())
		}

		// Return response sukses
		return jsonResponse(c, fiber.StatusCreated, "Data successfully added", typeUser)
	}

	// Jika key "example" tidak ada
	return jsonResponse(c, fiber.StatusBadRequest, "'example' key is required", nil)
}

func DeleteType(c *fiber.Ctx) error {
	id := c.Params("id")

	if models.DB.Delete(&models.TypeUser{}, id).RowsAffected == 0 {
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