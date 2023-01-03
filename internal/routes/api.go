package routes

import (
	"net/http"
	"os"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/compress"
	jwtware "github.com/gofiber/jwt/v3"
	"github.com/golang-jwt/jwt/v4"
	"github.com/yapi-teklif/internal/handlers"
	dingoapp "github.com/yapi-teklif/internal/pkg/dingo"
)

const (
	Supplier    string = "/supplier"
	Constructor string = "/constructor"
	Admin       string = "/admin"
	V1          string = "/v1"
	Auth        string = "/auth"
	Public      string = "/public"
)

func InitRoutes(app *fiber.App) {

	api := app.Group("/api")
	api.Use(compress.New(compress.Config{
		Level: compress.LevelBestSpeed, // 1
	}))

	func() {
		v1 := api.Group(V1)
		v1 = v1.Group(Auth)
		InitAuthRoutes(v1)

		supplier := api.Group(V1)
		supplier.Use(jwtware.New(jwtware.Config{
			SigningKey: []byte(os.Getenv("ACCESS_SECRET")),
		}))
		supplier = supplier.Group(Supplier)
		supplier.Use(func(ctx *fiber.Ctx) error {
			err := handlers.TokenValid(ctx)
			if err != nil {
				return ctx.JSON(http.StatusUnauthorized)
			}
			return ctx.Next()
		})
		supplier.Use(func(ctx *fiber.Ctx) error {

			supplier := ctx.Locals("user").(*jwt.Token)
			claims := supplier.Claims.(jwt.MapClaims)

			if claims == nil {
				return ctx.SendStatus(http.StatusUnauthorized)
			}
			if int64(claims["expiration"].(float64)) < time.Now().UTC().Unix() {
				return ctx.SendStatus(http.StatusUnauthorized)
			}
			is_supplier, ok := claims["is_supplier"]
			if ok {
				if is_supplier != true {
					return ctx.SendStatus(http.StatusUnauthorized)
				} else {
					ctx.Locals("is_supplier", is_supplier)
				}
			} else {
				return ctx.SendStatus(http.StatusUnauthorized)
			}
			ctx.Locals("id", claims["id"])
			return ctx.Next()
		})
		InitSupplierRoutes(supplier)

		constructor := api.Group(V1)
		constructor.Use(jwtware.New(jwtware.Config{
			SigningKey: []byte(os.Getenv("ACCESS_SECRET")),
		}))
		constructor = constructor.Group(Constructor)
		constructor.Use(func(ctx *fiber.Ctx) error {
			err := handlers.TokenValid(ctx)
			if err != nil {
				return ctx.JSON(http.StatusUnauthorized)
			}
			return ctx.Next()
		})
		constructor.Use(func(ctx *fiber.Ctx) error {

			constructor := ctx.Locals("user").(*jwt.Token)
			claims := constructor.Claims.(jwt.MapClaims)

			if claims == nil {
				return ctx.SendStatus(http.StatusUnauthorized)
			}
			if int64(claims["expiration"].(float64)) < time.Now().UTC().Unix() {
				return ctx.SendStatus(http.StatusUnauthorized)
			}
			is_constructor, ok := claims["is_constructor"]
			if ok {
				if is_constructor != true {
					return ctx.SendStatus(http.StatusUnauthorized)
				} else {
					ctx.Locals("is_constructor", is_constructor)
				}
			} else {
				return ctx.SendStatus(http.StatusUnauthorized)
			}
			ctx.Locals("id", claims["id"])
			return ctx.Next()
		})
		InitConstructorRoutes(constructor)

	}()

}

func InitAuthRoutes(group fiber.Router) {
	company_controller := dingoapp.Application.Container.GetCompanyHandler()
	group.Post("/signup", company_controller.Create)
	group.Post("/login", company_controller.Login)
	group.Delete("/logout", company_controller.Logout)
	group.Post("/token/refresh", company_controller.Refresh)
}

func InitSupplierRoutes(group fiber.Router) {
	productHandler := dingoapp.Application.Container.GetProductHandler()
	group.Post("/product", productHandler.Create)
	group.Put("/product/:id", productHandler.Update)
	group.Delete("/product/:id", productHandler.Delete)
	group.Get("/product/:id", productHandler.Get)
	group.Get("/products/:id", productHandler.GetAllByCategory)

	variantHandler := dingoapp.Application.Container.GetVariantHandler()
	group.Post("/variant", variantHandler.Create)
	group.Put("/variant/:id", variantHandler.Update)
	group.Delete("/variant/:id", variantHandler.Delete)
	group.Get("/variant/:id", variantHandler.Get)

	featureHandler := dingoapp.Application.Container.GetFeatureHandler()
	group.Get("/feature/:id", featureHandler.Get)
}

func InitConstructorRoutes(group fiber.Router) {
	construction_handler := dingoapp.Application.Container.GetConstructionHandler()
	group.Post("/construction", construction_handler.Create)
	group.Put("/construction/:id", construction_handler.Update)
	group.Delete("/construction/:id", construction_handler.Delete)
	group.Get("/construction/:id", construction_handler.Get)
	group.Get("/constructions/:id", construction_handler.GetAllByCategory)
}
