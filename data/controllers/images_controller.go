package controllers

import (
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync/atomic"

	"github.com/gofiber/fiber/v2"
)

// Global counter untuk generate nama file dengan increment
var imageCounter uint64

func nextImageID() uint64 {
	return atomic.AddUint64(&imageCounter, 1)
}

// validateFileType memeriksa MIME type file berdasarkan header-nya.
// Hanya file dengan MIME type image/jpeg, image/png, dan image/gif yang diizinkan.
func validateFileType(file *multipart.FileHeader) error {
	f, err := file.Open()
	if err != nil {
		return fmt.Errorf("failed to open file for type validation: %v", err)
	}
	defer f.Close()

	// Baca 512 byte pertama untuk mendeteksi MIME type
	buffer := make([]byte, 512)
	if _, err := io.ReadFull(f, buffer); err != nil && err != io.EOF {
		return fmt.Errorf("failed to read file header: %v", err)
	}
	contentType := http.DetectContentType(buffer)

	allowedTypes := map[string]bool{
		"image/jpeg": true,
		"image/png":  true,
		"image/gif":  true,
	}

	if !allowedTypes[contentType] {
		return fmt.Errorf("unsupported file type: %s", contentType)
	}

	return nil
}

// ensureDirectoryExists memastikan direktori upload ada, jika tidak ada maka akan dibuat.
func ensureDirectoryExists(dir string) error {
	if _, err := os.Stat(dir); os.IsNotExist(err) {
		if err := os.MkdirAll(dir, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %v", err)
		}
	}
	return nil
}

func UploadImage(c *fiber.Ctx) error {
	// Ambil file dari request
	file, err := c.FormFile("image")
	if err != nil {
		fmt.Println("error reading file:", err)
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "unable to read the file. Ensure 'image' key is included in the form-data request",
		})
	}

	fmt.Printf("File received: %s (Size: %d bytes)\n", file.Filename, file.Size)

	// Validasi ukuran file (maks 5 MB)
	const maxFileSize = 5 * 1024 * 1024 // 5 MB
	if file.Size > maxFileSize {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "file size exceeds 5 MB limit",
		})
	}

	// Validasi tipe file dengan MIME
	if err := validateFileType(file); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": err.Error(),
		})
	}

	// Ambil direktori upload dari environment variable, atau gunakan default
	uploadDir := os.Getenv("UPLOAD_DIR")
	if uploadDir == "" {
		uploadDir = "/var/www/html/images"
	}
	if err := ensureDirectoryExists(uploadDir); err != nil {
		fmt.Println("error ensuring directory exists:", err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to create upload directory",
		})
	}

	// Sanitasi nama file dan validasi ekstensi
	ext := strings.ToLower(filepath.Ext(file.Filename))
	allowedExtensions := map[string]bool{
		".jpg":  true,
		".jpeg": true,
		".png":  true,
		".gif":  true,
	}
	if !allowedExtensions[ext] {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"error": "unsupported file type",
		})
	}

	// Generate nama file menggunakan increment
	fileName := fmt.Sprintf("%d%s", nextImageID(), ext)
	filePath := filepath.Join(uploadDir, fileName)

	// Pastikan path tetap di dalam direktori yang diizinkan
	if !strings.HasPrefix(filepath.Clean(filePath), filepath.Clean(uploadDir)) {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"error": "invalid file path",
		})
	}

	// Simpan file
	if err := c.SaveFile(file, filePath); err != nil {
		fmt.Printf("error saving file (%s): %v\n", filePath, err)
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"error": "failed to save the file",
		})
	}

	// URL publik untuk file
	publicURL := fmt.Sprintf("https://web.ayomenjadi.com/images/%s", fileName)

	return c.Status(fiber.StatusOK).JSON(fiber.Map{
		"message":    "image uploaded successfully",
		"image_path": publicURL,
	})
}
