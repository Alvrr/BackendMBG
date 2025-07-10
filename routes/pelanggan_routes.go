package routes

import (
	"backend/controllers"

	"github.com/gofiber/fiber/v2"
)

func PelangganRoutes(app *fiber.App) {
	pelanggan := app.Group("/pelanggan")

	pelanggan.Get("/", controllers.GetAllPelanggan)
	pelanggan.Get("/:id", controllers.GetPelangganByID)
	pelanggan.Post("/", controllers.CreatePelanggan)
	pelanggan.Put("/:id", controllers.UpdatePelanggan)
	pelanggan.Delete("/:id", controllers.DeletePelanggan)
}
