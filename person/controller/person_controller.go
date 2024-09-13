package personcontroller

import (
	"log"
	"sample/custom"
	personmodel "sample/person/model"
	"sample/script"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func CreatePerson(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		// Main model
		var person personmodel.Person
		var Addresses personmodel.Address
		var Identifications personmodel.Identification
		var Contacts personmodel.Contact

		// Parse person data from the request body
		if err := c.Bind().Body(&person); err != nil {
			log.Printf("Error parsing person data: %+v", err)
			return custom.SendErrorResponse(c, custom.NewHttpError("Invalid person data", fiber.StatusBadRequest))
		}

		// Use the generic function to create the person and related resources
		return script.CreateResource(db, &person,
			&Addresses,
			&Identifications,
			&Contacts)(c)
	}
}

func GetAllPersons(db *gorm.DB) fiber.Handler {
	return script.GetAllResources[personmodel.Person](db, []string{"Addresses", "Identifications", "Contacts", "Merchant", "Merchant.Product", "Merchant.ContactMerchant", "Merchant.AddressMerchant"} )
}



func GetPersonByID(db *gorm.DB) fiber.Handler {
	return script.GetResourceByID[personmodel.Person](db, []string{"Addresses", "Identifications", "Contacts", "Merchant", "Merchant.Product", "Merchant.ContactMerchant", "Merchant.AddressMerchant"})
}

func UpdatePerson(db *gorm.DB) fiber.Handler {
	var person personmodel.Person
	return script.UpdateResource[personmodel.Person](db, &person)
}

func DeletePerson(db *gorm.DB) fiber.Handler {
	return script.DeleteResource[personmodel.Person](db)
}
