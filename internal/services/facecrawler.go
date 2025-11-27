package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/GuiFernandess7/risa/pkg/utils"
)

type BaseFaceCrawlerResponse struct {
	IDSearch string      `json:"id_search"`
	Message  string      `json:"message"`
	Progress interface{} `json:"progress"`
	Error    interface{} `json:"error"`
	Code     interface{} `json:"code"`
	Output   Output      `json:"output"`
}

type Output struct {
	Items []Item `json:"items"`
}

type FaceCrawlerStartResult struct {
	IDSearch string `json:"id_search"`
	Message  string `json:"message"`
}

type Item struct {
	URL string `json:"url"`
}

type FaceCrawler struct {
	Client *http.Client
}

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

func (fc FaceCrawler) Start(input SearchInput) (any, error) {

	if len(input.ImageBytes) == 0 {
		return nil, fmt.Errorf("FaceCrawler requires bytes")
	}

	site := os.Getenv("SITE_URL")
	apiKey := os.Getenv("FACECRAWLER_KEY")
	headers := map[string]string{
		"Authorization": apiKey,
		"Accept":        "application/json",
	}

	uploadURL := site + "api/upload_pic"
	respBytes, status, err := fc.sendImageBytes(uploadURL, headers, input.ImageBytes)
	if _, failed, err := utils.Try(respBytes, err); failed {
		return nil, err
	}

	if status != 200 {
		return nil, fmt.Errorf("service returned %d", status)
	}

	var parsed BaseFaceCrawlerResponse
	if err := json.Unmarshal(respBytes, &parsed); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}

	if parsed.Error != nil {
		return nil, fmt.Errorf("%v (%v)", parsed.Error, parsed.Code)
	}

	return FaceCrawlerStartResult{
		IDSearch: parsed.IDSearch,
		Message:  parsed.Message,
	}, nil
}

func (fc FaceCrawler) Check(jobID string) (any, error) {

	site := os.Getenv("SITE_URL")
	apiKey := os.Getenv("FACECRAWLER_KEY")

	jsonPayload := map[string]any{
		"id_search":     jobID,
		"with_progress": true,
		"status_only":   false,
		"demo":          true,
	}

	b, _ := json.Marshal(jsonPayload)
	log.Printf("[CHECK] Sending: %s", string(b))

	req, err := http.NewRequest("POST", site+"api/search", bytes.NewBuffer(b))
	if _, failed, err := utils.Try(req, err); failed {
		return nil, err
	}

	req.Header.Set("Authorization", apiKey)
	req.Header.Set("Accept", "application/json")
	req.Header.Set("Content-Type", "application/json")

	resp, err := fc.Client.Do(req)
	if _, failed, err := utils.Try(resp, err); failed {
		return nil, err
	}
	defer resp.Body.Close()

	respBytes, _ := io.ReadAll(resp.Body)
	log.Printf("[CHECK] Response: %s", respBytes)

	var parsed BaseFaceCrawlerResponse
	if err := json.Unmarshal(respBytes, &parsed); err != nil {
		return nil, err
	}

	if parsed.Error != nil {
		return nil, fmt.Errorf("%v (%v)", parsed.Error, parsed.Code)
	}

	return parsed, nil
}

func (fc FaceCrawler) sendImageBytes(url string, headers map[string]string, image []byte) ([]byte, int, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	writer.WriteField("id_search", "")
	filePart, err := writer.CreateFormFile("images", "upload.jpg")

	if _, failed, err := utils.Try(filePart, err); failed {
		return nil, http.StatusInternalServerError, err
	}
	if _, err := filePart.Write(image); err != nil {
		return nil, http.StatusInternalServerError, err
	}

	writer.Close()
	req, err := http.NewRequest("POST", url, body)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := fc.Client.Do(req)

	if err != nil {
		return nil, http.StatusInternalServerError, err
	}

	defer resp.Body.Close()
	respBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, 0, err
	}

	return respBytes, resp.StatusCode, nil
}
