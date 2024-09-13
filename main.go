package main

import (
	"sample/database"
	merchantmodel "sample/merchant/model"
	personmodel "sample/person/model"

	"sample/routes"

	"log"

	"github.com/gofiber/fiber/v3"
)

func main() {

	// Create a new Fiber app with the custom validator
	app := fiber.New()

	// Initialize the database connection
	db := database.InitDB()

	if err := db.AutoMigrate(&personmodel.Person{}, &personmodel.Address{}, &personmodel.Identification{}, &personmodel.Contact{}, &merchantmodel.Merchant{}, &merchantmodel.Product{}, &merchantmodel.ContactMerchant{}, merchantmodel.AddressMerchant{},  ); err != nil {
log.Println("migrate", err)
	}; log.Println("success")
	
	

	// Setup routes
	routes.SetupRoutes(app, db)

	// Start the Fiber app
	log.Fatal(app.Listen(":3000"))
}

