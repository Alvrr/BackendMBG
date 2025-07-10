package controllers

import (
	"backend/config"
	"backend/models"
	"backend/repository"
	"backend/utils"
	"context"
	"fmt"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
)

// GET /pembayaran
func GetAllPembayaran(c *fiber.Ctx) error {
	data, err := repository.GetAllPembayaran()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal ambil data pembayaran",
			"error":   err.Error(),
		})
	}
	return c.JSON(data)
}

// GET /pembayaran/:id
func GetPembayaranByID(c *fiber.Ctx) error {
	id := c.Params("id")
	data, err := repository.GetPembayaranByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Data tidak ditemukan",
			"error":   err.Error(),
		})
	}
	return c.JSON(data)
}

// fungsi untuk generate ID pembayaran dengan prefix "PM"
func generatePembayaranID() (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	filter := bson.M{"_id": "pembayaran"}
	update := bson.M{"$inc": bson.M{"seq": 1}}
	opts := options.FindOneAndUpdate().SetUpsert(true).SetReturnDocument(options.After)

	var counter struct {
		Seq int64 `bson:"seq"`
	}

	err := config.CounterCollection.FindOneAndUpdate(ctx, filter, update, opts).Decode(&counter)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			// kalau belum ada dokumen, buat baru dengan seq 1
			_, err = config.CounterCollection.InsertOne(ctx, bson.M{"_id": "pembayaran", "seq": 1})
			if err != nil {
				return "", err
			}
			counter.Seq = 1
		} else {
			return "", err
		}
	}

	return fmt.Sprintf("PM%03d", counter.Seq), nil
}

// POST /pembayaran
func CreatePembayaran(c *fiber.Ctx) error {
	var pembayaran models.Pembayaran

	if err := c.BodyParser(&pembayaran); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Request tidak valid",
			"error":   err.Error(),
		})
	}

	// ✅ Validasi data
	if err := utils.Validate.Struct(pembayaran); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "Validasi gagal",
			"error":   err.Error(),
		})
	}

	// generate ID pembayaran otomatis
	id, err := generatePembayaranID()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal generate ID pembayaran",
			"error":   err.Error(),
		})
	}
	pembayaran.ID = id

	// Ambil nama produk berdasarkan ID
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	for i, item := range pembayaran.Produk {
		var produk models.Produk
		err := config.ProdukCollection.FindOne(ctx, bson.M{"_id": item.IDProduk}).Decode(&produk)
		if err == nil {
			pembayaran.Produk[i].NamaProduk = produk.NamaProduk
		} else {
			pembayaran.Produk[i].NamaProduk = "Produk tidak ditemukan"
		}
	}

	pembayaran.Tanggal = time.Now().Format("2006-01-02 15:04:05")

	result, err := repository.CreatePembayaran(pembayaran)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menyimpan data",
			"error":   err.Error(),
		})
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Pembayaran berhasil ditambahkan",
		"data":    result.InsertedID,
	})
}

// PUT /pembayaran/:id
func UpdatePembayaran(c *fiber.Ctx) error {
	id := c.Params("id")
	var pembayaran models.Pembayaran

	if err := c.BodyParser(&pembayaran); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Request tidak valid",
			"error":   err.Error(),
		})
	}

	// ✅ Validasi data
	if err := utils.Validate.Struct(pembayaran); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "Validasi gagal",
			"error":   err.Error(),
		})
	}

	// Tetap update tanggal ke sekarang
	pembayaran.Tanggal = time.Now().Format("2006-01-02 15:04:05")

	_, err := repository.UpdatePembayaran(id, pembayaran)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal update data",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Data berhasil diupdate",
	})
}

// DELETE /pembayaran/:id
func DeletePembayaran(c *fiber.Ctx) error {
	id := c.Params("id")

	_, err := repository.DeletePembayaran(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal hapus data",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Data pembayaran berhasil dihapus",
	})
}
