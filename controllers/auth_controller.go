package controllers

import (
	"backend/models"
	"backend/repository"
	"backend/utils"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"golang.org/x/crypto/bcrypt"
)

// Register handler
func Register(c *fiber.Ctx) error {
	var input models.User

	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Data tidak valid"})
	}

	// Validasi role wajib (admin/kasir/driver)
	role := strings.ToLower(input.Role)
	if role != "admin" && role != "kasir" && role != "driver" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Role harus admin, kasir, atau driver"})
	}

	// Generate ID otomatis dari counters.go (jangan ambil dari frontend)
	id, err := repository.GenerateID(role)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membuat ID user"})
	}
	input.ID = id

	// Hash password
	hashed, err := bcrypt.GenerateFromPassword([]byte(input.Password), 12)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal enkripsi password"})
	}
	input.Password = string(hashed)
	input.CreatedAt = time.Now()

	// Simpan user
	if err := repository.CreateUser(&input); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal registrasi user"})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Registrasi berhasil",
		"id":      input.ID,
	})
}

// Login handler
func Login(c *fiber.Ctx) error {
	var input models.LoginInput
	if err := c.BodyParser(&input); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Data tidak valid"})
	}

	user, err := repository.FindUserByEmail(input.Email)
	if err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Email atau password salah"})
	}

	// Cek status aktif/nonaktif
	if user.Status != "aktif" {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Akun karyawan nonaktif, tidak dapat login"})
	}

	// Bandingkan password plaintext dan hashed
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(input.Password)); err != nil {
		return c.Status(fiber.StatusUnauthorized).JSON(fiber.Map{"error": "Email atau password salah"})
	}

	// ✅ Perbaikan: tambahkan user.Nama
	token, err := utils.GenerateToken(user.ID, user.Role, user.Nama)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membuat token"})
	}

	return c.JSON(fiber.Map{
		"message": "Login berhasil",
		"token":   token,
	})
}
