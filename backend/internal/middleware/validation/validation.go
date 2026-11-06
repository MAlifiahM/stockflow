package validation

import (
	"regexp"

	"github.com/MAlifiahM/stockflow/backend/internal/domain"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v3"
)

var validate *validator.Validate

func init() {
	validate = validator.New(validator.WithRequiredStructEnabled())

	// Register custom username validation
	validate.RegisterValidation("username", func(fl validator.FieldLevel) bool {
		usernameRegex := regexp.MustCompile(`^[a-zA-Z0-9_]+$`)
		return usernameRegex.MatchString(fl.Field().String())
	})
}

func New[V any]() fiber.Handler {
	return func(c fiber.Ctx) error {
		var v V
		if err := c.Bind().Body(&v); err != nil {
			return fiber.NewError(fiber.StatusBadRequest, err.Error())
		}
		if err := validate.Struct(v); err != nil {
			var errors []string
			for _, err := range err.(validator.ValidationErrors) {
				message := err.Field() + " is " + err.Tag()
				if err.Param() != "" {
					message += " " + err.Param()
				}
				errors = append(errors, message)
			}
			return c.Status(fiber.StatusBadRequest).JSON(domain.Error{
				Code:    fiber.StatusBadRequest,
				Errors:  errors,
				Message: "validation error",
			})
		}
		c.Locals("parser", &v)
		return c.Next()
	}
}
