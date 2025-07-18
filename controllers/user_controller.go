package controllers

import (
	"backend/models"
	"backend/repository"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// GET /drivers
func GetAllDrivers(c *fiber.Ctx) error {
	drivers, err := repository.GetAllDrivers()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data driver",
			"error":   err.Error(),
		})
	}
	return c.JSON(drivers)
}

// CRUD Karyawan (admin only)
func GetAllKaryawan(c *fiber.Ctx) error {
	if c.Locals("userRole") != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Akses hanya untuk admin"})
	}
	users, err := repository.GetAllKaryawan()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data karyawan",
			"error":   err.Error(),
		})
	}
	return c.JSON(users)
}

func GetKaryawanByID(c *fiber.Ctx) error {
	if c.Locals("userRole") != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Akses hanya untuk admin"})
	}
	id := c.Params("id")
	user, err := repository.GetKaryawanByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Karyawan tidak ditemukan"})
	}
	return c.JSON(user)
}

func CreateKaryawan(c *fiber.Ctx) error {
	if c.Locals("userRole") != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Akses hanya untuk admin"})
	}
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Request tidak valid"})
	}

	// Generate ID untuk user
	newID, err := repository.GenerateID("userid")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal generate ID user",
			"error":   err.Error(),
		})
	}
	user.ID = newID

	// Set default status aktif jika kosong
	if user.Status == "" {
		user.Status = "aktif"
	}
	// Hash password
	if user.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal hash password"})
		}
		user.Password = string(hashed)
	}
	user.CreatedAt = time.Now()
	_, err = repository.CreateKaryawan(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal menambah karyawan"})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Karyawan berhasil ditambah"})
}

func UpdateKaryawan(c *fiber.Ctx) error {
	if c.Locals("userRole") != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Akses hanya untuk admin"})
	}
	id := c.Params("id")
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Request tidak valid"})
	}
	// Hash password jika diupdate
	if user.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal hash password"})
		}
		user.Password = string(hashed)
	}
	_, err := repository.UpdateKaryawan(id, user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal update karyawan"})
	}
	return c.JSON(fiber.Map{"message": "Karyawan berhasil diupdate"})
}

// Register khusus karyawan (bukan di halaman login)
func RegisterKaryawan(c *fiber.Ctx) error {
	var user models.User
	if err := c.BodyParser(&user); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Request tidak valid"})
	}

	// Generate ID untuk user
	newID, err := repository.GenerateID("userid")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal generate ID user",
			"error":   err.Error(),
		})
	}
	user.ID = newID

	// Set default status aktif
	user.Status = "aktif"
	// Hash password
	if user.Password != "" {
		hashed, err := bcrypt.GenerateFromPassword([]byte(user.Password), 12)
		if err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal hash password"})
		}
		user.Password = string(hashed)
	}
	user.CreatedAt = time.Now()
	_, err = repository.CreateKaryawan(user)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal register karyawan"})
	}
	return c.Status(fiber.StatusCreated).JSON(fiber.Map{"message": "Register karyawan berhasil"})
}

func DeleteKaryawan(c *fiber.Ctx) error {
	if c.Locals("userRole") != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Akses hanya untuk admin"})
	}
	id := c.Params("id")
	_, err := repository.DeleteKaryawan(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal hapus karyawan"})
	}
	return c.JSON(fiber.Map{"message": "Karyawan berhasil dihapus"})
}

// PATCH /users/karyawan/:id/status
func UpdateKaryawanStatus(c *fiber.Ctx) error {
	if c.Locals("userRole") != "admin" {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{"message": "Akses hanya untuk admin"})
	}
	id := c.Params("id")
	var body struct {
		Status string `json:"status"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Request tidak valid"})
	}
	if body.Status != "aktif" && body.Status != "nonaktif" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Status harus 'aktif' atau 'nonaktif'"})
	}
	if err := repository.UpdateKaryawanStatus(id, body.Status); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal update status karyawan"})
	}
	return c.JSON(fiber.Map{"message": "Status karyawan berhasil diupdate"})
}

// GET /users/karyawan/active
func GetActiveKaryawan(c *fiber.Ctx) error {
	// Kasir/driver juga boleh akses
	users, err := repository.GetActiveKaryawan()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"message": "Gagal mengambil data karyawan aktif"})
	}
	return c.JSON(users)
}
