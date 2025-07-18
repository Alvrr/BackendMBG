package controllers

import (
	"backend/models"
	"backend/repository"
	"backend/utils"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// GET /produk
func GetAllProduk(c *fiber.Ctx) error {
	produks, err := repository.GetAllProduk()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data produk",
			"error":   err.Error(),
		})
	}
	return c.JSON(produks)
}

// GET /produk/:id
func GetProdukByID(c *fiber.Ctx) error {
	id := c.Params("id")
	produk, err := repository.GetProdukByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Produk tidak ditemukan",
			"error":   err.Error(),
		})
	}
	return c.JSON(produk)
}

// POST /produk
func CreateProduk(c *fiber.Ctx) error {
	var produk models.Produk

	if err := c.BodyParser(&produk); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Request tidak valid",
			"error":   err.Error(),
		})
	}

	// ✅ Tambahkan validasi
	if err := utils.Validate.Struct(produk); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "Validasi gagal",
			"error":   err.Error(),
		})
	}

	// 🔢 Generate ID dan waktu
	newID, err := repository.GenerateID("produkid")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal generate ID produk",
			"error":   err.Error(),
		})
	}

	produk.ID = newID
	produk.CreatedAt = primitive.NewDateTimeFromTime(time.Now())

	result, err := repository.CreateProduk(produk)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menambahkan produk",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Produk berhasil ditambahkan",
		"data":    result.InsertedID,
	})
}

// PUT /produk/:id
func UpdateProduk(c *fiber.Ctx) error {
	id := c.Params("id")
	var produk models.Produk

	if err := c.BodyParser(&produk); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Request tidak valid",
			"error":   err.Error(),
		})
	}

	// Hapus field ID agar tidak ikut di-update (hindari error immutable _id)
	produk.ID = ""

	// ✅ Validasi input
	if err := utils.Validate.Struct(produk); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "Validasi gagal",
			"error":   err.Error(),
		})
	}

	_, err := repository.UpdateProduk(id, produk)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal update produk",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Produk berhasil diupdate",
	})
}

// DELETE /produk/:id
func DeleteProduk(c *fiber.Ctx) error {
	id := c.Params("id")

	_, err := repository.DeleteProduk(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal hapus produk",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Produk berhasil dihapus",
	})
}
