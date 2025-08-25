package main

import (
	"fitness/services/nats"
	"database/sql"
	"log"
	"os"
	"fmt"
	"github.com/gofiber/fiber/v2"
	_ "github.com/go-sql-driver/mysql"
	natsio "github.com/nats-io/nats.go"

	"fitness/config"
	"fitness/routes"
)
func main() {
	config.LoadEnv()
	redis.Init()
	elasticsearch.Init()
	app := fiber.New()
	// nats setup
	nats.Init() // 👈 This is required
	// creating stream and subject
	err := nats.CreateStream("LISTING_STREAM", "listings.create")
	if err != nil {
		log.Fatal("Failed to create NATS stream:", err)
	}
	

	mysqlHost := os.Getenv("MYSQL_HOST")
	mysqlPort := os.Getenv("MYSQL_PORT")
	mysqlUser := os.Getenv("MYSQL_USER")
	mysqlPassword := os.Getenv("MYSQL_PASSWORD")
	mysqlDB := os.Getenv("MYSQL_DB")

	//db, err := sql.Open("mysql", "root:root@tcp(host.docker.internal:3306)/fitness") // Change credentials as needed
		// Use the environment variables to construct the DSN
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", mysqlUser, mysqlPassword, mysqlHost, mysqlPort, mysqlDB)

	db, err := sql.Open("mysql", dsn)

	if err != nil {
		log.Fatal(err)
	}

	routes.SetupListingRoutes(app, db)
	routes.SetupAuthRoutes(app, db)
	// After Init(), CreateStream() and DB connection setup
	err = nats.PollMessages("listings.create", "listing-worker", func(msg *natsio.Msg) bool {
		return nats.ProcessListingInsert(db, msg)
	}, 10)
	if err != nil {
		log.Fatal("Polling error:", err)
	}
	go elasticsearch.Sync()

	port := config.GetEnv("APP_PORT", "3000")
	log.Fatal(app.Listen(":" + port))
}
