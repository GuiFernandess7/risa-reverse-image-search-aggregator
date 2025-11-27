package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"
)

type ImgBBResponse struct {
	Data struct {
		ID         string `json:"id"`
		Title      string `json:"title"`
		URLViewer  string `json:"url_viewer"`
		URL        string `json:"url"`
		DisplayURL string `json:"display_url"`
		Width      int    `json:"width"`
		Height     int    `json:"height"`
		Size       int    `json:"size"`
		Time       int64  `json:"time"`
		Expiration int    `json:"expiration"`

		// Image struct {
		// 	Filename  string `json:"filename"`
		// 	Name      string `json:"name"`
		// 	Mime      string `json:"mime"`
		// 	Extension string `json:"extension"`
		// 	URL       string `json:"url"`
		// } `json:"image"`

		// Thumb struct {
		// 	Filename  string `json:"filename"`
		// 	Name      string `json:"name"`
		// 	Mime      string `json:"mime"`
		// 	Extension string `json:"extension"`
		// 	URL       string `json:"url"`
		// } `json:"thumb"`

		Medium struct {
			Filename  string `json:"filename"`
			Name      string `json:"name"`
			Mime      string `json:"mime"`
			Extension string `json:"extension"`
			URL       string `json:"url"`
		} `json:"medium"`

		DeleteURL string `json:"delete_url"`
	} `json:"data"`

	Success bool `json:"success"`
	Status  int  `json:"status"`
}

type YandexSearch struct {
	Client *http.Client
}

func (ys YandexSearch) Name() string {
	return "yandex"
}

func (ys YandexSearch) RequiresImageURL() bool {
	return true
}

func NewYandexSearch() *YandexSearch {
	return &YandexSearch{
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (ys YandexSearch) Search(input SearchInput) (any, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if err := writer.WriteField("key", os.Getenv("HOST_IMAGE_KEY")); err != nil {
		return nil, err
	}

	part, err := writer.CreateFormFile("image", "upload.jpg")
	if err != nil {
		return nil, err
	}

	if _, err := part.Write(input.ImageBytes); err != nil {
		return nil, err
	}

	writer.Close()

	req, err := http.NewRequest("POST", os.Getenv("HOST_IMAGE_URL"), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var parsed ImgBBResponse
	if err := json.Unmarshal(respBytes, &parsed); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	return parsed, nil

	// Send to Yandex image reverse search
}
