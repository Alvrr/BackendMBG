package middleware

import (
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

func CorsMiddleware() fiber.Handler {
	return cors.New(cors.Config{
		AllowOrigins:     "https://frontend-mbg.vercel.app,https://backendmbg-production.up.railway.app,http://localhost:5000,https://localhost:5000", // Frontend + Railway domain + localhost untuk Swagger testing
		AllowHeaders:     "Origin, Content-Type, Accept, Authorization",
		AllowMethods:     "GET,POST,PUT,DELETE",
		AllowCredentials: true,
	})
}
