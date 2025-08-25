package routes

import (
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"fitness/handlers"
)

func SetupListingRoutes(app *fiber.App, db *sql.DB) {
	listing := app.Group("/api/listings")

	// Search listings by keyword and location (city, lat | long)
	listing.Get("/search", func(c *fiber.Ctx) error {
		keyword := c.Query("q")
		location := c.Query("location")
		if keyword == "" && location == "" {
			return c.SendString("Please provide a search query or location.")
}
		listings, err := handlers.SearchListings(db, keyword, location)
		if err != nil {
			return err
		}
		return c.JSON(listings)
	})

	// Get listing details by slug
	listing.Get("/:slug", func(c *fiber.Ctx) error {
		slug := c.Params("slug")
		listing, err := handlers.GetListingBySlug(db, slug)
		if err != nil {
			return err
		}
		return c.JSON(listing)
	})

	// Create new listing
	listing.Post("/", func(c *fiber.Ctx) error {
		var listing handlers.Listing
		if err := c.BodyParser(&listing); err != nil {
			return err
		}
		err = handlers.CreateListing(db, &listing)
		if err != nil {
			return err
		}
		return c.JSON(listing)
	})
}
