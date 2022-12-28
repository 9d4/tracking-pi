package main

import (
	"context"
	"github.com/9d4/tracking-pi/db"
	"github.com/9d4/tracking-pi/handler"
	"github.com/9d4/tracking-pi/industry"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/basicauth"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/mongo/readpref"
	"log"
	"os"
)

func main() {
	loadEnv()
	loadDB()
	loadStores()
	app := fiber.New()

	basicAuthMw := basicauth.New(basicauth.Config{
		Users: map[string]string{
			os.Getenv("ADMIN"): os.Getenv("ADMIN_PASSWORD"),
		},
	})

	app.Group("/admin", basicAuthMw)
	api := app.Group("/api", basicAuthMw)

	api.Get("/industries", handler.HandleIndex)
	api.Post("/industries", handler.HandleStore)

	app.Static("/", "views/dist/", fiber.Static{
		Compress:       os.Getenv("DEVELOPMENT") != "true",
		ByteRange:      false,
		Browse:         false,
		Download:       false,
		Index:          "index.html",
		CacheDuration:  0,
		MaxAge:         0,
		ModifyResponse: nil,
		Next:           nil,
	})

	log.Fatal(app.Listen(os.Getenv("ADDRESS")))
}

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}
}

func loadDB() {
	cli, err := db.Open(os.Getenv("MONGODB_URI"))
	if err != nil {
		log.Fatal(err)
	}

	if err := cli.Ping(context.TODO(), readpref.Primary()); err != nil {
		log.Fatal(err)
	}

	db.SetClient(cli)
	db.SetDB(os.Getenv("MONGODB_DB"))
}

func loadStores() {
	industry.SetStore(industry.NewStore(db.DB()))
}
