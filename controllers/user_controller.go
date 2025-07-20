package controllers

import (
	"backend/models"
	"backend/repository"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// GetAllDrivers godoc
//
//	@Summary		Get all drivers
//	@Description	Mengambil semua data driver
//	@Tags			Driver
//	@Produce		json
//	@Success		200	{array}		models.User
//	@Failure		500	{object}	map[string]interface{}
//	@Router			/users/drivers [get]
//
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

// GetAllKaryawan godoc
//
//	@Summary		Get all users
//	@Description	Mengambil semua data user/karyawan (admin only)
//	@Tags			Karyawan
//	@Security		BearerAuth
//	@Produce		json
//	@Success		200	{array}		models.User
//	@Failure		403	{object}	map[string]interface{}
//	@Failure		500	{object}	map[string]interface{}
//	@Router			/users/karyawan [get]
//
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

	// Validasi role hanya kasir dan driver yang bisa dibuat
	if user.Role != "kasir" && user.Role != "driver" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Role harus kasir atau driver"})
	}

	// Generate ID untuk user berdasarkan role
	newID, err := repository.GenerateID(user.Role)
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
	_, err = repository.CreateKaryawan(&user)
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

	// Validasi role hanya kasir dan driver yang bisa dibuat
	if user.Role != "kasir" && user.Role != "driver" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Role harus kasir atau driver"})
	}

	// Generate ID untuk user berdasarkan role
	newID, err := repository.GenerateID(user.Role)
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
	err = repository.CreateUser(&user)
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
	if id == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "ID karyawan tidak boleh kosong"})
	}

	var body struct {
		Status string `json:"status"`
	}
	if err := c.BodyParser(&body); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Request tidak valid"})
	}

	// Validasi status hanya boleh aktif atau nonaktif
	if body.Status != "aktif" && body.Status != "nonaktif" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"message": "Status harus 'aktif' atau 'nonaktif'"})
	}

	// Cek apakah karyawan ada terlebih dahulu
	_, err := repository.GetKaryawanByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"message": "Karyawan tidak ditemukan"})
	}

	// Update status
	if err := repository.UpdateKaryawanStatus(id, body.Status); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal update status karyawan",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Status karyawan berhasil diupdate",
		"id":      id,
		"status":  body.Status,
	})
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
