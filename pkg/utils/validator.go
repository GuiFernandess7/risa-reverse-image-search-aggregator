package utils

import (
	"fmt"

	"github.com/go-playground/validator/v10"
	"github.com/labstack/echo/v4"
)

type CustomValidator struct {
	Validator *validator.Validate
}

func (cv *CustomValidator) Validate(i any) error {
	return cv.Validator.Struct(i)
}

func ValidateRequestParams(c echo.Context, allowedParams []string) error {
	allowedMap := make(map[string]bool)
	for _, v := range allowedParams {
		allowedMap[v] = true
	}

	for param := range c.QueryParams() {
		if !allowedMap[param] {
			return fmt.Errorf("parâmetro inválido: %s", param)
		}
	}
	return nil
}
