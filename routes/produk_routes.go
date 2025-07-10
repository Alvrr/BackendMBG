package routes

import (
	"backend/controllers"

	"github.com/gofiber/fiber/v2"
)

func ProdukRoutes(app *fiber.App) {
	produk := app.Group("/produk")

	produk.Get("/", controllers.GetAllProduk)
	produk.Get("/:id", controllers.GetProdukByID)
	produk.Post("/", controllers.CreateProduk)
	produk.Put("/:id", controllers.UpdateProduk)
	produk.Delete("/:id", controllers.DeleteProduk)
}
