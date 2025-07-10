package controllers

import (
	"backend/models"
	"backend/repository"
	"backend/utils"

	"github.com/gofiber/fiber/v2"
)

// GET /pelanggan
func GetAllPelanggan(c *fiber.Ctx) error {
	data, err := repository.GetAllPelanggan()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal mengambil data pelanggan",
			"error":   err.Error(),
		})
	}
	return c.JSON(data)
}

// GET /pelanggan/:id
func GetPelangganByID(c *fiber.Ctx) error {
	id := c.Params("id")
	pelanggan, err := repository.GetPelangganByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Pelanggan tidak ditemukan",
			"error":   err.Error(),
		})
	}
	return c.JSON(pelanggan)
}

// POST /pelanggan
func CreatePelanggan(c *fiber.Ctx) error {
	var pelanggan models.Pelanggan

	if err := c.BodyParser(&pelanggan); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Request tidak valid",
			"error":   err.Error(),
		})
	}

	// ✅ Validasi input
	if err := utils.Validate.Struct(pelanggan); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "Validasi gagal",
			"error":   err.Error(),
		})
	}

	newID, err := repository.GenerateID("pelangganid")
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal generate ID pelanggan",
			"error":   err.Error(),
		})
	}

	pelanggan.ID = newID

	result, err := repository.CreatePelanggan(pelanggan)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menambahkan pelanggan",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Pelanggan berhasil ditambahkan",
		"data":    result.InsertedID,
	})
}

// PUT /pelanggan/:id
func UpdatePelanggan(c *fiber.Ctx) error {
	id := c.Params("id")
	var pelanggan models.Pelanggan

	if err := c.BodyParser(&pelanggan); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Request tidak valid",
			"error":   err.Error(),
		})
	}

	// ✅ Validasi input
	if err := utils.Validate.Struct(pelanggan); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "Validasi gagal",
			"error":   err.Error(),
		})
	}

	_, err := repository.UpdatePelanggan(id, pelanggan)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal update pelanggan",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Pelanggan berhasil diupdate",
	})
}

// DELETE /pelanggan/:id
func DeletePelanggan(c *fiber.Ctx) error {
	id := c.Params("id")

	_, err := repository.DeletePelanggan(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal hapus pelanggan",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Pelanggan berhasil dihapus",
	})
}
