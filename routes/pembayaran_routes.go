package routes

import (
	"backend/controllers"

	"github.com/gofiber/fiber/v2"
)

func PembayaranRoutes(app *fiber.App) {
	pembayaran := app.Group("/pembayaran")

	pembayaran.Get("/", controllers.GetAllPembayaran)
	pembayaran.Get("/:id", controllers.GetPembayaranByID)
	pembayaran.Post("/", controllers.CreatePembayaran)
	pembayaran.Put("/:id", controllers.UpdatePembayaran)
	pembayaran.Delete("/:id", controllers.DeletePembayaran)
}
