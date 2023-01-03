package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	dingoapp "github.com/yapi-teklif/internal/pkg/dingo"
	"github.com/yapi-teklif/internal/routes"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/gofiber/helmet/v2"
	"github.com/sarulabs/dingo/v4"
	database "github.com/yapi-teklif/internal/pkg/database/connection"
	"github.com/yapi-teklif/internal/pkg/dingo/provider"
)

func init() {
	err := dingo.GenerateContainer((*provider.Provider)(nil), "../../internal/pkg/dingo/container")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}
}

func main() {

	err := godotenv.Load(os.ExpandEnv("$GOPATH/src/github.com/yapi-teklif/.env"))

	if err != nil {
		log.Fatal(".env file couldn't loaded")
	}
	dingoapp.New()
	defer dingoapp.Application.Container.Delete()

	app := fiber.New()
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		fmt.Println("Gracefully shutting down...")
		_ = app.Shutdown()
	}()

	app.Use(helmet.New())
	app.Use(func(c *fiber.Ctx) error {
		// Set some security headers:
		c.Set("X-XSS-Protection", "1; mode=block")
		c.Set("X-Content-Type-Options", "nosniff")
		c.Set("X-Download-Options", "noopen")
		c.Set("Strict-Transport-Security", "max-age=5184000")
		c.Set("X-Frame-Options", "SAMEORIGIN")
		c.Set("X-DNS-Prefetch-Control", "off")

		// Go to next middleware:
		return c.Next()
	})

	app.Use(cors.New(cors.Config{
		AllowOrigins:     "*",
		AllowMethods:     "GET,POST,PUT,DELETE,OPTIONS",
		AllowHeaders:     "Accept,Authorization,Content-Type,X-CSRF-TOKEN",
		ExposeHeaders:    "Link",
		AllowCredentials: true,
		MaxAge:           300,
	}))

	routes.InitRoutes(app)

	migrate := flag.Bool("migrate", false, "a bool")
	if migrate != nil && *migrate {
		err := database.AutoMigrate(dingoapp.Application.Container.GetExternalConnections())
		if err != nil {
			log.Fatal(err)
		}

		return
	}

	if err := app.Listen(":8080" /*+ os.Getenv("PORT")*/); err != nil {
		log.Panic(err)
	}
	fmt.Println("Running cleanup tasks...")
	// Your cleanup tasks go here
}
