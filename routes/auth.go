package routes

import (
	"database/sql"

	"fitness/handlers"

	"github.com/gofiber/fiber/v2"
)

func SetupAuthRoutes(app *fiber.App, db *sql.DB) {
	app.Post("/login", handlers.Login(db))
	app.Get("/verify-login", middleware.AuthMiddleware(), handlers.VerifyLogin())
}
