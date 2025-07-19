package main

import (
	"backend/config"
	"backend/middleware"
	"backend/repository"
	"backend/routes"
	"log"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load file .env
	if err := godotenv.Load(); err != nil {
		log.Fatal("❌ Gagal load file .env")
	}

	// Koneksi ke MongoDB
	config.ConnectDB()

	// Inisialisasi counters yang diperlukan
	if err := repository.InitializeCounters(); err != nil {
		log.Printf("⚠️ Peringatan: %v", err)
	} else {
		log.Println("✅ Counters berhasil diinisialisasi")
	}

	// Inisialisasi Fiber
	app := fiber.New()

	// Middleware global
	app.Use(middleware.LoggerMiddleware())
	app.Use(middleware.CorsMiddleware())

	// JWTMiddleware global, kecuali untuk /auth/login dan /auth/register
	app.Use(func(c *fiber.Ctx) error {
		if c.Path() == "/auth/login" || c.Path() == "/auth/register" {
			return c.Next()
		}
		return middleware.JWTMiddleware(c)
	})

	// Semua route (termasuk auth/login/register)
	routes.SetupRoutes(app)

	// Port server
	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	log.Println("🚀 Server jalan di http://localhost:" + port)
	log.Fatal(app.Listen(":" + port))
}
