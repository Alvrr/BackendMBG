package main

import (
	"backend/config"
	_ "backend/docs" // Import docs for swagger
	"backend/middleware"
	"backend/repository"
	"backend/routes"
	"log"
	"os"
	"strings"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	fiberSwagger "github.com/swaggo/fiber-swagger"
)

//	@title			Backend MBG API
//	@version		1.0
//	@description	API documentation untuk Backend MBG application.
//	@description
//	@description	**Sistem Login:**
//	@description	- Admin: admin@email.com / password123
//	@description	- Kasir: kasir@email.com / password123
//	@description	- Driver: driver@email.com / password123
//	@description
//	@description	**Authentication:**
//	@description	- Semua endpoint (kecuali login) memerlukan Bearer Token
//	@description	- Token didapat dari endpoint /auth/login
//	@description	- Format: Authorization: Bearer {token}
//	@description
//	@description	**Role Permissions:**
//	@description	- Admin: Akses penuh ke semua fitur
//	@description	- Kasir: CRUD produk, pelanggan, pembayaran, lihat riwayat
//	@description	- Driver: Lihat pembayaran assigned, update status selesai, cetak surat jalan
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	API Support
//	@contact.url	http://www.swagger.io/support
//	@contact.email	support@swagger.io

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:5000
//	@BasePath	/

//	@securityDefinitions.apikey	BearerAuth
//	@in							header
//	@name						Authorization
//	@description				Type "Bearer" followed by a space and JWT token.

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
		path := c.Path()
		if path == "/auth/login" || path == "/auth/register" || strings.HasPrefix(path, "/swagger") {
			return c.Next()
		}
		return middleware.JWTMiddleware(c)
	})

	// Swagger route
	app.Get("/swagger/*", fiberSwagger.WrapHandler)

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
