package validator

import (
	"reflect"
	"strings"

	"villainrsty-ecommerce-server/internal/core/shared/errors"

	"github.com/go-playground/validator/v10"
)

type Validator struct {
	validate *validator.Validate
}

func NewValidate() *Validator {
	val := validator.New()

	// Key error pakai json tag: email, password, refresh_token, dst
	val.RegisterTagNameFunc(func(fld reflect.StructField) string {
		name := strings.SplitN(fld.Tag.Get("json"), ",", 2)[0]
		if name == "" || name == "-" {
			return fld.Name
		}
		return name
	})

	return &Validator{
		validate: val,
	}
}

func fieldError(field, message string) error {
	return &errors.AppError{
		Kind:    errors.ErrValidation,
		Message: "Validation failed",
		Fields:  map[string]string{field: message},
	}
}

func (v *Validator) ValidateEmail(email string) error {
	if email == "" {
		return fieldError("email", "Email is required")
	}

	if err := v.validate.Var(email, "email"); err != nil {
		return fieldError("email", "Invalid email format")
	}

	return nil
}

func (v *Validator) ValidatePassword(password string) error {
	if password == "" {
		return fieldError("password", "Password is required")
	}

	if err := v.validate.Var(password, "min=8"); err != nil {
		return fieldError("password", "Minimum 8 characters required")
	}

	hasUpper := false
	hasLower := false
	hasNumber := false

	for _, c := range password {
		if c >= 'A' && c <= 'Z' {
			hasUpper = true
		}

		if c >= 'a' && c <= 'z' {
			hasLower = true
		}

		if c >= '0' && c <= '9' {
			hasNumber = true
		}
	}

	if !hasLower || !hasNumber || !hasUpper {
		return fieldError("password", "Password must contain Uppercase, Lowercase and Number")
	}

	return nil
}

func (v *Validator) ValidatePhone(phone string) error {
	if phone == "" {
		return fieldError("phone", "Phone is required")
	}

	if err := v.validate.Var(phone, "e164"); err != nil {
		return fieldError("phone", "Invalid phone format")
	}

	return nil
}

func (v *Validator) ValidateName(name string) error {
	if name == "" {
		return fieldError("name", "Name is require")
	}

	if err := v.validate.Var(name, "min=1,max=100"); err != nil {
		return fieldError("name", "Name must be between 1 and 100 characters")
	}

	return nil
}

func (v *Validator) ValidateURL(url string) error {
	if url == "" {
		return fieldError("url", "URL is required")
	}

	if err := v.validate.Var(url, "url"); err != nil {
		return fieldError("url", "Invalid URL format")
	}

	return nil
}

func (v *Validator) ValidateRequired(field, value string) error {
	if err := v.validate.Var(value, "required"); err != nil {
		return fieldError(field, field+" is required")
	}

	return nil
}

func (v *Validator) ValidateStruct(data interface{}) error {
	if err := v.validate.Struct(data); err != nil {
		errorFields := make(map[string]string)

		if validationErrors, ok := err.(validator.ValidationErrors); ok {
			for _, fieldErr := range validationErrors {
				errorFields[fieldErr.Field()] = msgForTag(
					fieldErr.Tag(),
					fieldErr.Param(),
				)
			}
		}

		return &errors.AppError{
			Kind:    errors.ErrValidation,
			Message: "Validation failed",
			Fields:  errorFields,
		}
	}

	return nil
}

// Helper untuk mengubah pesan error validator menjadi text user-friendly
func msgForTag(tag string, param string) string {
	switch tag {
	case "required":
		return "This field is required"
	case "email":
		return "Invalid email format"
	case "min":
		return "Minimum " + param + " characters required"
	case "max":
		return "Maximum " + param + " characters allowed"
	}
	return "Invalid value"
}
