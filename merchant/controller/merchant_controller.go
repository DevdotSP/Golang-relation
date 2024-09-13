package merchantcontroller

import (
	"log"
	"sample/custom"
	merchantmodel "sample/merchant/model"
	"sample/script"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func GetAllProduct(db *gorm.DB) fiber.Handler {
	return func (c fiber.Ctx) error  {
		var product []merchantmodel.Product

		if err := db.Find(&product).Error; err != nil {
			return custom.SendErrorResponse(c, custom.NewHttpError("Invalid person data", fiber.StatusBadRequest))
		} 

		if err := db.Find(&product).Error; err != nil {
			return custom.SendErrorResponse(c, custom.NewHttpError("Could not retrieve resource", fiber.StatusNotFound))
		}

		if len(product) == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "No records found",
			}) }

		return c.Status(fiber.StatusOK).JSON(product)
	}
}

func CreateProduct(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		var	product     merchantmodel.Product      

		// Parse person data from the request body
		if err := c.Bind().Body(&product); err != nil {
			log.Printf("Error parsing person data: %+v", err)
			return custom.SendErrorResponse(c, custom.NewHttpError("Invalid person data", fiber.StatusBadRequest))
		}

		// Use the generic function to create the person and related resources
		return script.CreateResource(db, &product)(c)
	}
}

func CreateMerchant(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		// Main model
		var merchant merchantmodel.Merchant
		var AddressMerchant      merchantmodel.AddressMerchant      
		var	ContactMerchant merchantmodel.ContactMerchant
		var	Product     merchantmodel.Product      

		// Parse person data from the request body
		if err := c.Bind().Body(&merchant); err != nil {
			log.Printf("Error parsing person data: %+v", err)
			return custom.SendErrorResponse(c, custom.NewHttpError("Invalid person data", fiber.StatusBadRequest))
		}

		// Use the generic function to create the person and related resources
		return script.CreateResource(db, &merchant, 
			&AddressMerchant, 
			&ContactMerchant, 
			&Product)(c)
	}
}

func GetAllMerchant(db *gorm.DB) fiber.Handler {
	return script.GetAllResources[merchantmodel.Merchant](db, []string{"AddressMerchant", "ContactMerchant", "Product"})
}


func GetMerchantByID(db *gorm.DB) fiber.Handler {
	return script.GetResourceByID[merchantmodel.Merchant](db, []string{"AddressMerchant", "ContactMerchant", "Product"})
}

func UpdateMerchant(db *gorm.DB) fiber.Handler {
	var merchant merchantmodel.Merchant
	return script.UpdateResource[merchantmodel.Merchant](db,&merchant)
}

func DeleteMerchant(db *gorm.DB) fiber.Handler {
	return script.DeleteResource[merchantmodel.Merchant](db)
}

