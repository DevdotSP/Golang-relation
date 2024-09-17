package merchantcontroller

import (
	merchantmodel "sample/merchant/model"
	"sample/response"
	"sample/script"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

func GetAllProduct(db *gorm.DB) fiber.Handler {
	return func (c fiber.Ctx) error  {
		var product []merchantmodel.Product

		if err := db.Find(&product).Error; err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(response.ErrorModel{
				RetCode: string(response.BadRequest),
				Message: "Invalid product data",
				Data: err,
			})
		} 

		if err := db.Find(&product).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(response.ErrorModel{
				RetCode: string(response.NotFound),
				Message: "Could not retrieve data",
				Data: err,
			})
		}

		if len(product) == 0 {
			return c.Status(fiber.StatusNotFound).JSON(response.ErrorModel{
				RetCode: string(response.NotFound),
				Message: "Could not retrieve data",
				Data: fiber.ErrNotFound,
			}) }

		return c.Status(fiber.StatusOK).JSON(product)
	}
}

func CreateProduct(db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		var	product     merchantmodel.Product      

		// Parse person data from the request body
		if err := c.Bind().Body(&product); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(response.ErrorModel{
				RetCode: string(response.BadRequest),
				Message: "Invalid request body",
				Data: err,
			})
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
			return c.Status(fiber.StatusBadRequest).JSON(response.ErrorModel{
				RetCode: string(response.BadRequest),
				Message: "Invalid request body",
				Data: err,
			})
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

