package routes

import (
	merchantcontroller "sample/merchant/controller"
	"sample/middleware"
	customercontroller "sample/customer/controller"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// SetupRoutes initializes the routes for the Fiber app
func SetupRoutes(app *fiber.App, db *gorm.DB) {

	// Group routes for persons under /api/person
	customerGroup := app.Group("/api/customer", middleware.HeadersMiddleware())
	{
		customerGroup.Post("/", customercontroller.Createcustomer(db))
		customerGroup.Get("/", customercontroller.GetAllcustomers(db))
		customerGroup.Get("/:id", customercontroller.GetcustomerByID(db))
		customerGroup.Put("/:id", customercontroller.Updatecustomer(db))
		customerGroup.Delete("/:id", customercontroller.Deletecustomer(db))
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
