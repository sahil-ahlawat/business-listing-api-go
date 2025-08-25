package handlers

import (
	"encoding/json"
	"net/http"
	"database/sql"
	"github.com/gofiber/fiber/v2"
	"fitness/models"
	"fitness/utils"
	"fitness/services/nats"
	"fitness/services/redis"
	"fitness/services/elasticsearch"
	"time"
	"fmt"
)

func CreateListingHandler(db *sql.DB) fiber.Handler {
	return func(c *fiber.Ctx) error {
		var req models.CreateListingRequest
		if err := c.BodyParser(&req); err != nil {
			return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request"})
		}

		// Add slug server-side before queueing
		req.Slug = utils.GenerateSlug(req.Title)

		// Marshal to JSON
		data, err := json.Marshal(req)
		if err != nil {
			return c.Status(500).SendString("Failed to serialize listing data")
		}

		// Publish to NATS subject
		if err := nats.Publish("elastic-search-sync", string(data)); err != nil {
			return c.Status(500).SendString(err.Error())
		}

		return c.JSON(fiber.Map{
			"message": "Listing queued successfully",
			"slug":    req.Slug,
		})
	}
}

func SearchListings(db *sql.DB, keyword, location string) (string, error) {
	// First, check Redis
	cacheKey := fmt.Sprintf("search:%s:%s", keyword, location)
	cachedResult, err := redis.Get(cacheKey)
	if err == nil {
		return cachedResult, nil
	}

	// If not in Redis, search Elasticsearch
	query := fmt.Sprintf(`{"query": {"bool": {"must": [{"match": {"title": "%s"}}, {"match": {"location": "%s"}}]}}}`, keyword, location)
	esResult, err := elasticsearch.Search("listings", query)
	if err != nil {
		return "", err
	}

	// Add to Redis with 30-minute expiry
	redis.Set(cacheKey, esResult, 30*time.Minute)

	return string(esResult), nil
}

func GetListingBySlug(db *sql.DB, slug string) (string, error) {
	// First, check Redis
	cacheKey := fmt.Sprintf("listing:%s", slug)
	cachedResult, err := redis.Get(cacheKey)
	if err == nil {
		return cachedResult, nil
	}

	// If not in Redis, get from DB
	var listing models.Listing
	err = db.QueryRow("SELECT id, title, description, lat, lng, location, slug FROM listings WHERE slug = ?", slug).Scan(&listing.ID, &listing.Title, &listing.Description, &listing.Lat, &listing.Lng, &listing.Location, &listing.Slug)
	if err != nil {
		return "", err
	}

	// Add to Redis with 30-minute expiry
	jsonListing, _ := json.Marshal(listing)
	redis.Set(cacheKey, jsonListing, 30*time.Minute)

	return string(jsonListing), nil
}
