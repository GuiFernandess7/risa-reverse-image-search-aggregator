package facecrawler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	interfaces "github.com/GuiFernandess7/risa/internal/repository/interfaces"
	"github.com/GuiFernandess7/risa/pkg/utils"
)

func (fc FaceCrawler) Name() string {
	return "facecrawler"
}

func (fc FaceCrawler) RequiresImageURL() bool {
	return false
}

func NewFaceCrawler() *FaceCrawler {
	return &FaceCrawler{
		Client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

func (fc FaceCrawler) Start(input interfaces.SearchInput) (any, error) {
	// Uploads the image to run search
	if len(input.ImageBytes) == 0 {
		return nil, fmt.Errorf("FaceCrawler requires bytes")
	}

	site := os.Getenv("SITE_URL")
	uploadURL := site + "api/upload_pic"

	log.Println("[STARTING] - Running facecrawler search...")
	writer, body, err := utils.GetFileRequestWriter("id_search", "", input.ImageBytes, "images")
	if err != nil {
		return nil, err
	}

	writer.Close()
	respBytes, statusCode, err := utils.SendRequest(uploadURL, body, os.Getenv("FACECRAWLER_KEY"), writer, true)
	if _, failed, err := utils.Try(respBytes, err); failed && statusCode != http.StatusOK {
		log.Printf("[ERROR] - Error calling facecrawler service: %v", err)
		return nil, fmt.Errorf("unexpected error occured")
	}

	var parsed BaseFaceCrawlerResponse
	if err := json.Unmarshal(respBytes, &parsed); err != nil {
		log.Printf("[ERROR] - Error parsing service response: %v", err)
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	if parsed.Error != nil {
		log.Printf("[ERROR] - Error calling facecrawler service: %v", err)
		return nil, fmt.Errorf("service error: %v (%v)", parsed.Error, parsed.Code)
	}

	log.Printf("[SUCCESS] - Facecrawler search ended.")
	return FaceCrawlerStartResult{
		IDSearch: parsed.IDSearch,
		Message:  parsed.Message,
	}, nil
}

func (fc FaceCrawler) Check(jobID string) (any, error) {
	// Checks the search status
	site := os.Getenv("SITE_URL")
	apiKey := os.Getenv("FACECRAWLER_KEY")

	jsonPayload := map[string]any{
		"id_search":     jobID,
		"with_progress": true,
		"status_only":   false,
		"demo":          true,
	}

	b, _ := json.Marshal(jsonPayload)
	log.Println("[STARTING] - Running facecrawler search by ID...")

	respBytes, statusCode, err := utils.SendRequest(site+"api/search", bytes.NewBuffer(b), apiKey, nil, false)
	if err != nil && statusCode != http.StatusOK {
		return nil, err
	}

	var parsed BaseFaceCrawlerResponse
	if err := json.Unmarshal(respBytes, &parsed); err != nil {
		return nil, err
	}
	if parsed.Error != nil {
		return nil, fmt.Errorf("%v (%v)", parsed.Error, parsed.Code)
	}
	return parsed, nil
}
