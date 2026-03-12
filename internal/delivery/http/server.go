package server

import (
	"arena-ban/internal/delivery/http/handler"

	"github.com/gofiber/fiber/v3"
	"github.com/gofiber/fiber/v3/middleware/static" // Wajib import ini
	"github.com/sirupsen/logrus"
)

func SetupServer(authHandler *handler.AuthHandler, logger *logrus.Logger) {
	// Inisialisasi App Fiber v3
	app := fiber.New(fiber.Config{
		AppName: "Arena Ban API v1",
	})

	// Middleware Logging
	app.Use(func(c fiber.Ctx) error {
		logger.Infof("[%s] %s - IP: %s", c.Method(), c.Path(), c.IP())
		return c.Next()
	})

	app.Get("/swagger.yaml", static.New("/app/docs/swagger.yaml"))

	app.Get("/docs/*", static.New("/app/assets/swagger-ui", static.Config{
		IndexNames: []string{"index.html"},
	}))

	app.Get("/", func(c fiber.Ctx) error {
		return c.Status(200).JSON(fiber.Map{
			"status":  "success",
			"message": "Arena Ban API is running! 🚀",
		})
	})

	api := app.Group("/api/")
	auth := api.Group("/auth")

	auth.Post("/login", authHandler.Login)
	auth.Post("/register", authHandler.Register)
	auth.Post("/verify-token", authHandler.VerifyToken)
	auth.Post("/setup-password", authHandler.RegisterPassword)
	auth.Post("/forgot-password", authHandler.ForgotPassword)
	auth.Post("/reset-password", authHandler.ResetPassword)

	port := ":3000"
	logger.Infof("Server starting on port %s", port)

	if err := app.Listen(port); err != nil {
		logger.Fatalf("Gagal menjalankan server: %v", err)
	}
}
