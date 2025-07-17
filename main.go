package main

import (
	"backend/config"
	"backend/middleware"
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

	// Inisialisasi Fiber
	app := fiber.New()

	// Middleware global
	app.Use(middleware.LoggerMiddleware())
	app.Use(middleware.CorsMiddleware())

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
