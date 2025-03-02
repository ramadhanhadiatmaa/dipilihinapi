package controllers

import (
	"data/models"
	"fmt"
	"strconv"

	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
)

func CreateStatus(c *fiber.Ctx) error {
	var data map[string]interface{}
	if err := c.BodyParser(&data); err != nil {
		return jsonResponse(c, fiber.StatusBadRequest, "Invalid input", err.Error())
	}

	allowedKeys := []string{"status"}

	for key := range data {
		if !contains(allowedKeys, key) {
			return jsonResponse(c, fiber.StatusBadRequest, "Inputting data is not allowed", nil)
		}
	}

	if exampleValue, exists := data["status"]; exists {
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
	return jsonResponse(c, fiber.StatusBadRequest, "key is required", nil)
}

func ShowStatus(c *fiber.Ctx) error {
	var data []models.Status

	if err := models.DB.Find(&data).Error; err != nil {
		return jsonResponse(c, fiber.StatusInternalServerError, "Failed to load data", err.Error())
	}

	if len(data) == 0 {
		return jsonResponse(c, fiber.StatusNotFound, "No data found", nil)
	}

	return c.JSON(data)
}

func IndexStatus(c *fiber.Ctx) error {
	id := c.Params("id")
	var data models.Status

	if err := models.DB.First(&data, "id = ?", id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return jsonResponse(c, fiber.StatusNotFound, "No data found", nil)
		}
		return jsonResponse(c, fiber.StatusInternalServerError, "Failed to load data", err.Error())
	}

	return c.JSON(data)
}

func UpdateStatus(c *fiber.Ctx) error {
	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		return jsonResponse(c, fiber.StatusBadRequest, "Invalid ID format", nil)
	}

	var data models.Status
	if err := models.DB.First(&data, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return jsonResponse(c, fiber.StatusNotFound, "No data found", nil)
		}
		return jsonResponse(c, fiber.StatusInternalServerError, "Failed to load data", err.Error())
	}

	var updateData models.Status
	if err := c.BodyParser(&updateData); err != nil {
		return jsonResponse(c, fiber.StatusBadRequest, "Invalid input", err.Error())
	}

	if updateData.ID != 0 && updateData.ID != id {
		if err := models.DB.First(&models.Status{}, updateData.ID).Error; err == nil {
			return jsonResponse(c, fiber.StatusBadRequest, "The updated ID is already in use", nil)
		}
	}

	if err := models.DB.Model(&data).Updates(updateData).Error; err != nil {
		return jsonResponse(c, fiber.StatusInternalServerError, "Failed to update data", err.Error())
	}

	return jsonResponse(c, fiber.StatusOK, "Data successfully updated", nil)
}

func DeleteStatus(c *fiber.Ctx) error {
	id := c.Params("id")

	if models.DB.Delete(&models.TypeUser{}, id).RowsAffected == 0 {
		return jsonResponse(c, fiber.StatusNotFound, "Data not found or already deleted", nil)
	}

	return jsonResponse(c, fiber.StatusOK, "Successfully deleted data", nil)
}
