package utils

import (
	"bytes"
	"io"
	"mime/multipart"
	"net/http"
)

func SendRequest(
	url string,
	body *bytes.Buffer,
	apiKey string,
	writer *multipart.Writer,
	multipart bool,
) ([]byte, int, error) {
	req, err := http.NewRequest("POST", url, body)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	req.Header.Set("Authorization", apiKey)
	req.Header.Set("Accept", "application/json")
	
	if !multipart {
		req.Header.Set("Content-Type", "application/json")
	} else {
		req.Header.Set("Content-Type", writer.FormDataContentType())
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	return respBody, resp.StatusCode, nil
}

func GetFileRequestWriter(
	fieldName string,
	fieldValue string,
	imageObject []byte,
) (*multipart.Writer, *bytes.Buffer, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if err := writer.WriteField(fieldName, fieldValue); err != nil {
		return nil, nil, err
	}

	part, err := writer.CreateFormFile("images", "upload.jpg")
	if err != nil {
		return nil, nil, err
	}

	if _, err := part.Write(imageObject); err != nil {
		return nil, nil, err
	}

	return writer, body, nil
}
