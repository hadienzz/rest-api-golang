package validation

import "github.com/go-playground/validator/v10"

var validate = validator.New()

func ValidateStruct(s interface{}) ([]string, error) {
	if err := validate.Struct(s); err != nil {
		ve, ok := err.(validator.ValidationErrors)
		if !ok {
			return nil, err
		}

		var errors []string
		for _, e := range ve {
			errors = append(errors, e.Error())
		}
		return errors, nil
	}
	return nil, nil
}
