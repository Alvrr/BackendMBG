package routes

import (
	"backend/controllers"

	"github.com/gofiber/fiber/v2"
)

func LaporanRoutes(app *fiber.App) {
	app.Get("/laporan/export/excel", controllers.ExportLaporanExcel)
}
