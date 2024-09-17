package customercontroller

import (
	customermodel "sample/customer/model"
	"sample/response"
	"sample/script"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func Createcustomer(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		// Main model
		var customer customermodel.Customer
		var Addresses customermodel.Address
		var Identifications customermodel.Identification
		var Contacts customermodel.Contact

		// Parse person data from the request body
		if err := c.Bind().Body(&customer); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(response.ErrorModel{
				RetCode: string(response.BadRequest),
				Message: "Invalid request body",
				Data:    err,
			})
		}

		// Use the generic function to create the person and related resources
		return script.CreateResource(db, &customer,
			&Addresses,
			&Identifications,
			&Contacts)(c)
	}
}

func GetAllcustomers(db *gorm.DB) fiber.Handler {
	return script.GetAllResources[customermodel.Customer](db, []string{"Addresses", "Identifications", "Contacts", "Merchant", "Merchant.Product", "Merchant.ContactMerchant", "Merchant.AddressMerchant"} )
}


func GetcustomerByID(db *gorm.DB) fiber.Handler {
	return script.GetResourceByID[customermodel.Customer](db, []string{"Addresses", "Identifications", "Contacts", "Merchant", "Merchant.Product", "Merchant.ContactMerchant", "Merchant.AddressMerchant"})
}

func Updatecustomer(db *gorm.DB) fiber.Handler {
	var customer customermodel.Customer
	return script.UpdateResource[customermodel.Customer](db, &customer)
}

func Deletecustomer(db *gorm.DB) fiber.Handler {
	return script.DeleteResource[customermodel.Customer](db)
}
