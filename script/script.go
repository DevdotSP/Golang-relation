package script

import (

	"log"
	"reflect"
	"sample/custom"
	"strings"

	"github.com/gofiber/fiber/v3"
	"gorm.io/gorm"
)

// CreateResource creates a resource in the database
func CreateProduct[T any](db *gorm.DB, input *T) fiber.Handler {
	return func(c fiber.Ctx) error {
		// Bind the request body to the main input model
		if err := c.Bind().Body(input); err != nil {
			return custom.SendErrorResponse(c, custom.NewHttpError("Invalid request body", fiber.StatusBadRequest))
		}

		// Create the main resource
		if err := db.Create(input).Error; err != nil {
			if isUniqueConstraintError(err) {
				return custom.SendErrorResponse(c, custom.NewHttpError("Duplicate entry detected", fiber.StatusConflict))
			}
			return custom.SendErrorResponse(c, custom.NewHttpError("Could not create resource", fiber.StatusInternalServerError))
		}

		// Extract the ID from the input model
		val := reflect.ValueOf(input).Elem() // Dereference the pointer to get the value
		idField := val.FieldByName("ID")
		if !idField.IsValid() {
			return custom.SendErrorResponse(c, custom.NewHttpError("ID field not found in resource", fiber.StatusInternalServerError))
		}

		id := idField.Uint() // Get the ID value

		// Respond with success message and created resource ID
		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Resource created successfully",
			"id":      id,
		})
	}
}


// CreateResource creates a resource and can optionally preload related models
func CreateResource[T any](db *gorm.DB, input *T, relatedModels ...interface{}) fiber.Handler {
	return func(c fiber.Ctx) error {
		// Bind the request body to the main input model
		if err := c.Bind().Body(input); err != nil {
			return custom.SendErrorResponse(c, custom.NewHttpError("Invalid request body", fiber.StatusBadRequest))
		}

		// Create the main resource
		if err := db.Create(input).Error; err != nil {
			if isUniqueConstraintError(err) {
				return custom.SendErrorResponse(c, custom.NewHttpError("Duplicate entry detected", fiber.StatusConflict))
			}
			return custom.SendErrorResponse(c, custom.NewHttpError("Could not create resource", fiber.StatusInternalServerError))
		}

		// Extract the ID from the input model
		val := reflect.ValueOf(input).Elem() // Dereference the pointer to get the value
		idField := val.FieldByName("ID")
		if !idField.IsValid() {
			return custom.SendErrorResponse(c, custom.NewHttpError("ID field not found in resource", fiber.StatusInternalServerError))
		}

		id := idField.Uint() // Get the ID value

		// Update related models with the ID
		for _, relatedModel := range relatedModels {
			if relatedModel != nil {
				relatedVal := reflect.ValueOf(relatedModel)
				if relatedVal.Kind() != reflect.Ptr || relatedVal.IsNil() {
					continue // Skip invalid models
				}

				relatedVal = relatedVal.Elem() // Dereference the pointer to get the value
				if relatedVal.Kind() != reflect.Slice && relatedVal.Kind() != reflect.Array {
					continue // Ensure it's a slice or array
				}

				// Iterate through the slice/array of related models
				for i := 0; i < relatedVal.Len(); i++ {
					elem := relatedVal.Index(i).Addr().Interface()
					elemVal := reflect.ValueOf(elem).Elem()

					// Set the foreign key field (PersonID)
					if field := elemVal.FieldByName("PersonID"); field.IsValid() && field.CanSet() {
						field.Set(reflect.ValueOf(id))
					}

					// Create the related resource
					if err := db.Create(elem).Error; err != nil {
						if isUniqueConstraintError(err) {
							return custom.SendErrorResponse(c, custom.NewHttpError("Duplicate entry detected in related model", fiber.StatusConflict))
						}
						return custom.SendErrorResponse(c, custom.NewHttpError("Could not create related resource", fiber.StatusInternalServerError))
					}
				}
			}
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Resource created successfully",
		})
	}
}


// Helper function to detect unique constraint violation
func isUniqueConstraintError(err error) bool {
	// Check if error contains specific keywords indicating a unique constraint violation
	return err != nil && (strings.Contains(err.Error(), "unique constraint"))
}

// Get all resources with optional preload
func GetAllResources[T any](db *gorm.DB, preloads []string) fiber.Handler {
	return func(c fiber.Ctx) error {
		var resources []T

		query := db
		for _, preload := range preloads {
			query = query.Preload(preload)
		}

		if err := query.Find(&resources).Error; err != nil {
			err := custom.NewHttpError("Could not retrieve resources", fiber.StatusInternalServerError)
			return custom.SendErrorResponse(c, err)
		}

		if len(resources) == 0 {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{
				"message": "No records found",
			})
		}

		return c.JSON(resources)
	}
}

// Get a resource by ID with optional preload
func GetResourceByID[T any](db *gorm.DB, preloads []string) fiber.Handler {
	return func(c fiber.Ctx) error {
		var resource T
		id := c.Params("id")

		resourceID, err := custom.ParseID(id)
		if err != nil {
			return custom.SendErrorResponse(c, custom.NewHttpError("Invalid ID", fiber.StatusBadRequest))
		}

		query := db
		for _, preload := range preloads {
			query = query.Preload(preload)
		}

		if err := query.First(&resource, resourceID).Error; err != nil {
			return custom.SendErrorResponse(c, custom.NewHttpError("Could not retrieve resource", fiber.StatusNotFound))
		}

		return c.JSON(resource)
	}
}

// Update a resource by ID
func UpdateResource[T any](db *gorm.DB, input *T) fiber.Handler {
	return func(c fiber.Ctx) error {
		id := c.Params("id")
		resourceID, err := custom.ParseID(id)
		if err != nil {
			return custom.SendErrorResponse(c, custom.NewHttpError("Invalid ID", fiber.StatusBadRequest))
		}

		// Parse request body into the input model
		if err := c.Bind().Body(input); err != nil {
			log.Println("Error parsing body:", err)
			return custom.SendErrorResponse(c, custom.NewHttpError("Invalid request body", fiber.StatusBadRequest))
		}

		// Check if the user exists before updating
		var existingUser T
		if err := db.First(&existingUser, resourceID).Error; err != nil {
			return custom.SendErrorResponse(c, custom.NewHttpError("Resource not found", fiber.StatusNotFound))
		}

		// Update only the fields present in the input struct
		if err := db.Model(&existingUser).Where("id = ?", resourceID).Updates(input).Error; err != nil {
			return custom.SendErrorResponse(c, custom.NewHttpError("Could not update resource", fiber.StatusInternalServerError))
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Resource updated successfully",
		})
	}
}

// DeleteResource deletes a resource by ID using GORM's cascading feature.
func DeleteResource[T any](db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		id := c.Params("id")
		resourceID, err := custom.ParseID(id) // Assuming ParseID handles ID parsing correctly
		if err != nil {
			return custom.SendErrorResponse(c, custom.NewHttpError("Invalid ID", fiber.StatusBadRequest))
		}

		// Delete the main resource
		if err := db.Delete(new(T), resourceID).Error; err != nil {
			return custom.SendErrorResponse(c, custom.NewHttpError("Could not delete resource", fiber.StatusInternalServerError))
		}

		return c.Status(fiber.StatusOK).JSON(fiber.Map{
			"message": "Resource deleted successfully",
		})
	}
}
