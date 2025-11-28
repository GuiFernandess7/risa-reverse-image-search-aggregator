package yandex

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/GuiFernandess7/risa/internal/services/types"
	"github.com/GuiFernandess7/risa/pkg/utils"
	g "github.com/serpapi/google-search-results-golang"
)

func NewYandexSearch() *YandexSearch {
	return &YandexSearch{
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (ys YandexSearch) Name() string {
	return "yandex"
}

func (ys YandexSearch) Search(input types.SearchInput) (any, error) {
	hostImageURL := os.Getenv("HOST_IMAGE_URL")
	hostImageKey := os.Getenv("HOST_IMAGE_KEY")
	serpAPIKey := os.Getenv("SERPAPI_KEY")

	writer, body, err := utils.GetFileRequestWriter("key", os.Getenv("HOST_IMAGE_KEY"), input.ImageBytes)
	if err != nil {
		return nil, err
	}

	writer.Close()
	respBytes, statusCode, err := utils.SendRequest(hostImageURL, body, hostImageKey, writer, true)
	if _, failed, err := utils.Try(respBytes, err); failed && statusCode != http.StatusOK {
		log.Printf("[ERROR] - Error calling host image service: %v", err)
		return nil, fmt.Errorf("unexpected error occured")
	}

	var imgBBResponseSchema ImgBBResponse
	if err := json.Unmarshal(respBytes, &imgBBResponseSchema); err != nil {
		log.Printf("[ERROR] - Error parsing image URL response: %v", err)
		return "", fmt.Errorf("invalid JSON: %w", err)
	}

	parameter := map[string]string{
		"engine": "yandex_images",
		"url":    imgBBResponseSchema.Data.URL,
	}

	search := g.NewGoogleSearch(parameter, serpAPIKey)
	results, err := search.GetJSON()
	if err != nil {
		log.Printf("[ERROR] - Error calling yandex service: %v", err)
		return nil, fmt.Errorf("unexpected error occured")
	}

	var yandexResponseSchema YandexSearchResponse
	if err := json.Unmarshal(respBytes, &yandexResponseSchema); err != nil {
		log.Printf("[ERROR] - Error parsing service response: %v", err)
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	log.Printf("[YANDEX] - Response data: %v", yandexResponseSchema)
	imageResults := results["image_results"]
	log.Printf("[SUCCESS] - Yandex search ended.")
	return imageResults, nil
}
