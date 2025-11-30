package utils

import (
	"mime/multipart"

	"github.com/labstack/echo/v4"
)

func GetFileObject(
	c echo.Context,
	fieldName string,
) (multipart.File, error) {
	file, err := c.FormFile(fieldName)
	if err != nil {
		return nil, err
	}

	src, err := file.Open()
	if err != nil {
		return nil, err
	}
	return src, nil
}
