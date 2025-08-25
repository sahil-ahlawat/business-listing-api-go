Miniml viable product

Backend:

Need Full name, gender, email id or phone number

Create user with above details and assign listings to the user.

Each listing should have a (name, Category, tags, address, location (lat/long), rating, phone number, email, wahtsapp number, profile pic/ logo, gallery pictures and videos)

Database should have option to feature the listing based on location/city.

Create a listing api :
Listing creation api should have support for all these fields, and with owner full name, phone number, email.

Fetch listing with pagination. (locations, ratings, and keyword(name, category) based listing.)


Listing single page should give all details and should have option to give rating and add rating comment with phone number or email otp.

Should use mysql database, redis amd elastic search.

Also create api to sync data between mysql and elastic search.

Frontend:

Should be simple vanilla js and html and css.
single html page with search bar anf featured listing. now when user start searching it shows listing like google and when single listing is clicked it take to a single listing page.
Should also have listing creation form. with all required fields. 
by detault all listing should be draft and should have option to publish. (As this is a mvp I will manually approve the listing in db via phpmyadmin)
on first load single page should be generated annd stored on server and should also be updated if listing info changes in db.

nd then we should store those html and static files in cloudflare r2 and then server with cloudflare cdn









///////////////////////////////////////////////////////////////////////////////////////////////

attached is the database mysql setup, I have fiber hello world setup, I have .env: APP_PORT=3000
APP_ENV=development
. go.mod: module fitness 

go 1.22.12

require (
	github.com/andybalholm/brotli v1.1.0 // indirect
	github.com/gofiber/fiber/v2 v2.52.8 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/joho/godotenv v1.5.1 // indirect
	github.com/klauspost/compress v1.17.9 // indirect
	github.com/mattn/go-colorable v0.1.13 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	github.com/mattn/go-runewidth v0.0.16 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	github.com/valyala/bytebufferpool v1.0.0 // indirect
	github.com/valyala/fasthttp v1.51.0 // indirect
	github.com/valyala/tcplisten v1.0.0 // indirect
	golang.org/x/sys v0.28.0 // indirect
)
, main.go: package main 

import (
    "fmt"
    "log"

    "github.com/gofiber/fiber/v2"
    "fitness/config"
)

func main() {
    config.LoadEnv()

    app := fiber.New()

    app.Get("/", func(c *fiber.Ctx) error {
        return c.SendString("Hello from Fiber!")
    })

    port := config.GetEnv("APP_PORT", "3000")
    log.Fatal(app.Listen(fmt.Sprintf(":%s", port)))
}
, config/config.go package config

import (
    "log"
    "os"

    "github.com/joho/godotenv"
)

func LoadEnv() {
    if err := godotenv.Load(); err != nil {
        log.Println("No .env file found, using system environment variables")
    }
}

func GetEnv(key, fallback string) string {
    if value := os.Getenv(key); value != "" {
        return value
    }
    return fallback
} now I help me create listing api, to get list of listings based on location, keyword(title, description, tags, features) with elastic search, api to sync between mysql and elastic search,  api to create listing, it should be a single api which will have user full name, email, phone number, by default listing status should be inactive. api to get single listing full details. api to upload image, video to cloudflarer2 and return url. and I also want to host static html of listing single pages, what it the best way to do that. keeo steps short and simple, give full code, lets create 1 api at a time. I am creating a mvp. I can update db structure if required aswell, you can use combination of mysql, redis, elasticsearch. it should be highly scaleable, so the mvp I am targetting is simple page like google where users will be searching for gyms or trainers or clubs, swimming pool, sports academy, yoga center etc in any locality. and when clicked on any listing it opens up a single page and from these users can contact gyms etc through email, phone, whatsapp etc and can see photo, video, ratings etc, features, tags, certifications etc.so you get the bigger picture, help me build this mvp in simple steps 1 complete feature at a time.   