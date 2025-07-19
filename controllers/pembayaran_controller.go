package controllers

import (
	"backend/config"
	"backend/models"
	"backend/repository"
	"backend/utils"
	"context"
	"fmt"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// GET /pembayaran
func GetAllPembayaran(c *fiber.Ctx) error {
	role := c.Locals("userRole").(string)
	id := c.Locals("userID").(string)

	filter := bson.M{}
	if role == "driver" {
		filter["id_driver"] = id
	} else if role == "kasir" {
		filter["id_kasir"] = id
	}

	data, err := repository.GetPembayaranFiltered(filter)
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
	role := c.Locals("userRole").(string)
	userID := c.Locals("userID").(string)

	data, err := repository.GetPembayaranByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Data tidak ditemukan",
			"error":   err.Error(),
		})
	}

	if role == "driver" && data.IDDriver != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Akses ditolak",
		})
	}

	return c.JSON(data)
}

// generate ID unik untuk transaksi
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
	if err != nil && err != mongo.ErrNoDocuments {
		return "", err
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

	if err := utils.Validate.Struct(pembayaran); err != nil {
		return c.Status(fiber.StatusUnprocessableEntity).JSON(fiber.Map{
			"message": "Validasi gagal",
			"error":   err.Error(),
		})
	}

	pembayaran.NamaKasir = c.Locals("userNama").(string)
	pembayaran.IDKasir = c.Locals("userID").(string)

	// Tambahkan nama produk berdasarkan ID
	for i, item := range pembayaran.Produk {
		var produk models.Produk
		err := config.ProdukCollection.FindOne(context.Background(), bson.M{"_id": item.IDProduk}).Decode(&produk)
		if err == nil {
			pembayaran.Produk[i].NamaProduk = produk.NamaProduk
		} else {
			pembayaran.Produk[i].NamaProduk = "Produk tidak ditemukan"
		}
	}

	if pembayaran.JenisPengiriman == "motor" || pembayaran.JenisPengiriman == "mobil" {
		if pembayaran.NamaDriver == "" {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
				"message": "Driver wajib dipilih untuk pengiriman motor/mobil",
			})
		}

		// Cari ID driver berdasarkan nama
		var driver models.User
		err := config.UserCollection.FindOne(context.Background(), bson.M{
			"nama":   pembayaran.NamaDriver,
			"role":   "driver",
			"status": "aktif",
		}).Decode(&driver)
		if err == nil {
			pembayaran.IDDriver = driver.ID
		}

		if pembayaran.JenisPengiriman == "motor" {
			pembayaran.Ongkir = 10000
		} else {
			pembayaran.Ongkir = 20000
		}
	} else {
		pembayaran.NamaDriver = "-"
		pembayaran.IDDriver = ""
		pembayaran.Ongkir = 0
	}

	total := 0
	for _, item := range pembayaran.Produk {
		total += item.Subtotal
	}
	pembayaran.TotalBayar = total + pembayaran.Ongkir

	id, err := generatePembayaranID()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal generate ID",
			"error":   err.Error(),
		})
	}
	pembayaran.ID = id
	pembayaran.Tanggal = time.Now().Format("2006-01-02 15:04:05")
	pembayaran.Status = "Pending" // Set status default

	result, err := repository.CreatePembayaran(pembayaran)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal simpan data",
			"error":   err.Error(),
		})
	}

	// Kurangi stok produk sesuai jumlah yang dibeli
	for _, item := range pembayaran.Produk {
		// Update stok produk: stok = stok - item.Jumlah
		filter := bson.M{"_id": item.IDProduk}
		update := bson.M{"$inc": bson.M{"stok": -item.Jumlah}}
		_, err := config.ProdukCollection.UpdateOne(context.Background(), filter, update)
		// Jika gagal update stok, log error saja, tidak batalkan transaksi
		if err != nil {
			fmt.Fprintf(os.Stderr, "Gagal update stok produk %s: %v\n", item.IDProduk, err)
		}
	}

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"message": "Berhasil ditambahkan",
		"data":    result.InsertedID,
	})
}

// PUT /pembayaran/selesaikan/:id
func SelesaikanPembayaran(c *fiber.Ctx) error {
	id := c.Params("id")
	role := c.Locals("userRole").(string)
	userID := c.Locals("userID").(string)

	pembayaran, err := repository.GetPembayaranByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Data tidak ditemukan",
		})
	}

	if role == "driver" && pembayaran.IDDriver != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Transaksi ini bukan milik driver",
		})
	}

	if pembayaran.Status == "Selesai" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{
			"message": "Transaksi sudah selesai",
		})
	}

	err = repository.SelesaikanPembayaran(id)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{
			"message": "Gagal menyelesaikan",
			"error":   err.Error(),
		})
	}

	return c.JSON(fiber.Map{
		"message": "Transaksi berhasil diselesaikan",
	})
}

// GET /pembayaran/cetak/:id
func CetakSuratJalan(c *fiber.Ctx) error {
	id := c.Params("id")
	role := c.Locals("userRole").(string)
	userID := c.Locals("userID").(string)

	data, err := repository.GetPembayaranByID(id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
			"message": "Data tidak ditemukan",
		})
	}

	if role == "driver" && data.IDDriver != userID {
		return c.Status(fiber.StatusForbidden).JSON(fiber.Map{
			"message": "Akses ditolak",
		})
	}

	// Print fisik ke printer TM-T82X (jika ada)
	go func(pembayaran *models.Pembayaran) {
		defer func() { recover() }()
		printerPath := "LPT1" // atau \\localhost\TM-T82X jika di Windows
		f, err := os.OpenFile(printerPath, os.O_RDWR, 0)
		if err == nil {
			defer f.Close()
			fmt.Fprintf(f, "\n==== STRUK PEMBAYARAN ===\n")
			fmt.Fprintf(f, "ID: %s\nTanggal: %s\n", pembayaran.ID, pembayaran.Tanggal)
			fmt.Fprintf(f, "Kasir: %s\nPelanggan: %s\n", pembayaran.NamaKasir, pembayaran.IDPelanggan)
			fmt.Fprintf(f, "Driver: %s\nPengiriman: %s\n", pembayaran.NamaDriver, pembayaran.JenisPengiriman)
			fmt.Fprintf(f, "--------------------------\n")
			for _, item := range pembayaran.Produk {
				fmt.Fprintf(f, "%s x%d @Rp%d = Rp%d\n", item.NamaProduk, item.Jumlah, item.Harga, item.Subtotal)
			}
			fmt.Fprintf(f, "--------------------------\n")
			fmt.Fprintf(f, "Ongkir: Rp%d\nTotal: Rp%d\n", pembayaran.Ongkir, pembayaran.TotalBayar)
			fmt.Fprintf(f, "\nTerima kasih!\n\n\n\n\n")
		}
	}(data)

	return c.JSON(fiber.Map{
		"message": "Surat jalan berhasil dicetak",
		"data":    data,
	})
}
