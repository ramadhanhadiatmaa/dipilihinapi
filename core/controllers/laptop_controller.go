package controllers

import (
	"strconv"
	"strings"

	"core/models"

	"github.com/gofiber/fiber/v2"
)

// GetLaptops meng-handle request untuk mengambil data laptop berdasarkan query parameter ids
func GetLaptops(c *fiber.Ctx) error {
	idsParam := c.Query("ids")
	if idsParam == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "Parameter 'ids' diperlukan",
		})
	}

	// Split string id yang dipisahkan koma, lalu konversi ke []int
	idStrs := strings.Split(idsParam, ",")
	var ids []int
	for _, idStr := range idStrs {
		idStr = strings.TrimSpace(idStr)
		id, err := strconv.Atoi(idStr)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"error": "Nilai ID tidak valid",
			})
		}
		ids = append(ids, id)
	}

	laptops, err := models.GetLaptopsByIds(ids)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "Gagal mengambil data laptop",
		})
	}

	return c.JSON(laptops)
}