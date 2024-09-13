package routes

import (
	merchantcontroller "sample/merchant/controller"
	"sample/middleware"
	personcontroller "sample/person/controller"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// SetupRoutes initializes the routes for the Fiber app
func SetupRoutes(app *fiber.App, db *gorm.DB) {

	// Group routes for persons under /api/person
	personGroup := app.Group("/api/person", middleware.HeadersMiddleware())
	{
		personGroup.Post("/", personcontroller.CreatePerson(db))
		personGroup.Get("/", personcontroller.GetAllPersons(db))
		personGroup.Get("/:id", personcontroller.GetPersonByID(db))
		personGroup.Put("/:id", personcontroller.UpdatePerson(db))
		personGroup.Delete("/:id", personcontroller.DeletePerson(db))
	}

	// Group routes for persons under /api/person
	merchantGroup := app.Group("/api/merchant", middleware.HeadersMiddleware())
	{
		merchantGroup.Post("/", merchantcontroller.CreateMerchant(db))
		merchantGroup.Get("/", merchantcontroller.GetAllMerchant(db))
		merchantGroup.Get("/:id", merchantcontroller.GetMerchantByID(db))
		merchantGroup.Put("/:id", merchantcontroller.UpdateMerchant(db))
		merchantGroup.Delete("/:id", merchantcontroller.DeleteMerchant(db))
	}

	// Group routes for persons under /api/person
	productGroup := app.Group("/api/product", middleware.HeadersMiddleware())
	{
		productGroup.Post("/", merchantcontroller.CreateProduct(db))
		productGroup.Get("/", merchantcontroller.GetAllProduct(db))
	}

}
