package utils

import (
	"bytes"
	"mime/multipart"
)

func GetFileRequestWriter(
	fieldName string,
	fieldValue string,
	imageObject []byte,
	formField string,
) (*multipart.Writer, *bytes.Buffer, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if err := writer.WriteField(fieldName, fieldValue); err != nil {
		return nil, nil, err
	}

	part, err := writer.CreateFormFile(formField, "upload.jpg")
	if err != nil {
		return nil, nil, err
	}

	if _, err := part.Write(imageObject); err != nil {
		return nil, nil, err
	}

	return writer, body, nil
}
