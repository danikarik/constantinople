package util

import (
	"errors"
	"fmt"
	"strings"

	"github.com/go-chi/render"
	"gopkg.in/go-playground/validator.v9"
)

// Validate checks struct against schema.
func Validate(b render.Binder, v *validator.Validate) error {
	if err := v.Struct(b); err != nil {
		if _, ok := err.(*validator.InvalidValidationError); ok {
			return err
		}
		var text string
		for _, err := range err.(validator.ValidationErrors) {
			text += fmt.Sprintf("%s -> [%s], ", strings.ToLower(err.Field()), err.Tag())
		}
		return errors.New(text[:len(text)-2])
	}
	return nil
}
