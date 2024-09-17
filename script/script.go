package script

import (
	"log"
	"reflect"
	"sample/custom"
	"sample/response"
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
				return c.Status(fiber.StatusForbidden).JSON(response.ErrorModel{
					RetCode: string(response.Forbidden),
					Message: "Could not create resource",
					Data: err,
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorModel{
				RetCode: string(response.InternalServerError),
				Message: "Could not create resource",
				Data: err,
			})
		}

		// Extract the ID from the input model
		val := reflect.ValueOf(input).Elem() // Dereference the pointer to get the value
		idField := val.FieldByName("ID")
		if !idField.IsValid() {
			return c.Status(fiber.StatusForbidden).JSON(response.ErrorModel{
				RetCode: string(response.Forbidden),
				Message: "id field not found",
				Data: fiber.ErrForbidden,
			})
		}

		id := idField.Uint() // Get the ID value

		return c.Status(fiber.StatusOK).JSON(response.ErrorModel{
			RetCode: string(response.SuccessOK),
			Message: "Success",
			Data: id,
		})
	}
}


// CreateResource creates a resource and can optionally preload related models
func CreateResource[T any](db *gorm.DB, input *T, relatedModels ...interface{}) fiber.Handler {
	return func(c fiber.Ctx) error {
		// Bind the request body to the main input model
		if err := c.Bind().Body(input); err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(response.ErrorModel{
				RetCode: string(response.BadRequest),
				Message: "Invalid request body",
				Data: err,
			})
		}

		// Create the main resource
		if err := db.Create(input).Error; err != nil {
			if isUniqueConstraintError(err) {
				return c.Status(fiber.StatusForbidden).JSON(response.ErrorModel{
					RetCode: string(response.Forbidden),
					Message: "Duplicate",
					Data: err,
				})
			}
			return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorModel{
				RetCode: string(response.InternalServerError),
				Message: "Could not create resource",
				Data: err,
			})
		}

		// Extract the ID from the input model
		val := reflect.ValueOf(input).Elem() // Dereference the pointer to get the value
		idField := val.FieldByName("ID")
		if !idField.IsValid() {
			return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorModel{
				RetCode: string(response.InternalServerError),
				Message: "Invalid ID",
				Data: fiber.StatusInternalServerError,
			})
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
					if field := elemVal.FieldByName("CustomerID"); field.IsValid() && field.CanSet() {
						field.Set(reflect.ValueOf(id))
					}

					// Create the related resource
					if err := db.Create(elem).Error; err != nil {
						if isUniqueConstraintError(err) {
							return c.Status(fiber.StatusForbidden).JSON(response.ErrorModel{
								RetCode: string(response.Forbidden),
								Message: "Duplicate data",
								Data: err,
							})
						}
						return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorModel{
							RetCode: string(response.InternalServerError),
							Message: "Could not create related resource",
							Data: err,
						})
					}
				}
			}
		}

		return c.Status(fiber.StatusOK).JSON(response.ErrorModel{
			RetCode: string(response.SuccessOK),
			Message: "Success Insert",
			Data: "Success",
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
			return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorModel{
				RetCode: string(response.InternalServerError),
				Message: "Could not retrieve resource",
				Data:    err,
			})
		}

		if len(resources) == 0 {
			return c.Status(fiber.StatusNotFound).JSON(response.ErrorModel{
				RetCode: string(response.NotFound),
				Message: "No resource found",
				Data:    resources,
			})
		}

		return c.Status(fiber.StatusOK).JSON(response.ErrorModel{
			RetCode: string(response.SuccessOK),
			Message: "success",
			Data:    resources,
		})
	}
}

// Get a resource by ID with optional preload
func GetResourceByID[T any](db *gorm.DB, preloads []string) fiber.Handler {
	return func(c fiber.Ctx) error {
		var resource T
		id := c.Params("id")

		resourceID, err := custom.ParseID(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(response.ErrorModel{
				RetCode: string(response.BadRequest),
				Message: "invalid id",
				Data:    err,
			})
		}

		query := db
		for _, preload := range preloads {
			query = query.Preload(preload)
		}

		if err := query.First(&resource, resourceID).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(response.ErrorModel{
				RetCode: string(response.NotFound),
				Message: "Could not find update resource",
				Data:    err,
			})
		}

		return c.Status(fiber.StatusOK).JSON(response.ErrorModel{
			RetCode: string(response.SuccessOK),
			Message: "Success",
			Data:    resource,
		})
	}
}

// Update a resource by ID
func UpdateResource[T any](db *gorm.DB, input *T) fiber.Handler {
	return func(c fiber.Ctx) error {
		id := c.Params("id")
		resourceID, err := custom.ParseID(id)
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(response.ErrorModel{
				RetCode: string(response.BadRequest),
				Message: "Invalid ID",
				Data:    resourceID,
			})
		}

		// Parse request body into the input model
		if err := c.Bind().Body(input); err != nil {
			log.Println("Error parsing body:", err)
			return c.Status(fiber.StatusBadRequest).JSON(response.ErrorModel{
				RetCode: string(response.BadRequest),
				Message: "Could not find update resource",
				Data:  err,
			})
		}

		// Check if the user exists before updating
		var existingUser T
		if err := db.First(&existingUser, resourceID).Error; err != nil {
			return c.Status(fiber.StatusNotFound).JSON(response.ErrorModel{
				RetCode: string(response.NotFound),
				Message: "Could not find update resource",
				Data:    existingUser,
			})
		}

		// Update only the fields present in the input struct
		if err := db.Model(&existingUser).Where("id = ?", resourceID).Updates(input).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorModel{
				RetCode: string(response.InternalServerError),
				Message: "Could not find update resource",
				Data:    existingUser,
			})
		}

		return c.Status(fiber.StatusOK).JSON(response.ErrorModel{
			RetCode: string(response.SuccessOK),
	    Message: "Update success",
	    Data: existingUser,
		})
	}
}

// DeleteResource deletes a resource by ID using GORM's cascading feature.
func DeleteResource[T any](db *gorm.DB) fiber.Handler {
	return func(c fiber.Ctx) error {
		id := c.Params("id")
		resourceID, err := custom.ParseID(id) // Assuming ParseID handles ID parsing correctly
		if err != nil {
			return c.Status(fiber.StatusBadRequest).JSON(response.ErrorModel{
				RetCode: string(response.BadRequest),
				Message: "Invalid ID",
				Data:    err,
			})
		}

		// Delete the main resource
		if err := db.Delete(new(T), resourceID).Error; err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(response.ErrorModel{
				RetCode: string(response.InternalServerError),
				Message: "Server Error",
				Data:    err.Error(),
			})
		}

		return c.Status(fiber.StatusOK).JSON(response.ErrorModel{
			RetCode: "200",
			Message: "Deleted Successfully",
			Data:    resourceID,
		})
	}
}
