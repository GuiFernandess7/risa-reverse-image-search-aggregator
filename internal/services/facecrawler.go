package service

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
)

type FaceCrawlerStartResponse struct {
	Engine string `json:"engine"`
	JobID  struct {
		Code            *string `json:"code"`
		Error           *string `json:"error"`
		FacemonLastScan *string `json:"facemon_last_scann"`
		FacemonStatus   *string `json:"facemon_status"`
		HasEmptyImages  bool    `json:"hasEmptyImages"`
		IDSearch        string  `json:"id_search"`
		Input           []struct {
			Base64    string      `json:"base64"`
			SvgAnim   interface{} `json:"svg_anim"`
			URLSource interface{} `json:"url_source"`
		} `json:"input"`
		Message    string      `json:"message"`
		NewSeen    int         `json:"new_seen_count"`
		Output     interface{} `json:"output"`
		Progress   interface{} `json:"progress"`
		WasUpdated bool        `json:"was_updated"`
	} `json:"job_id"`
}

type FaceCrawlerStatusResponse struct {
	IDSearch string `json:"id_search"`
	Message  any    `json:"message"`
	Progress int    `json:"progress"`
	Error    any    `json:"error"`
	Code     any    `json:"code"`
	Output   any    `json:"output"`
	Input    []struct {
		Base64    string `json:"base64"`
		IDPic     string `json:"id_pic"`
		URLSource any    `json:"url_source"`
		SVGAnim   any    `json:"svg_anim"`
	} `json:"input"`
	FacemonStatus   any `json:"facemon_status"`
	FacemonLastScan any `json:"facemon_last_scann"`
	NewSeenCount    int `json:"new_seen_count"`
	NewFrom         any `json:"new_from"`
}

type FaceCrawlerStartResult struct {
	Engine   string `json:"engine"`
	IDSearch string `json:"id_search"`
	Message  string `json:"message"`
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

	log.Println("[SEARCHING] Starting FaceCrawler...")
	apiKey := os.Getenv("FACECRAWLER_KEY")
	serviceUrl := os.Getenv("SITE_URL") + "api/upload_pic"
	headers := map[string]string{
		"Authorization": apiKey,
		"Accept":        "application/json",
	}

	log.Println("[SEARCHING] Sending image bytes...")
	log.Println("URL destination:", serviceUrl)
	responseBytes, statusCode, err := fc.sendImageBytes(serviceUrl, headers, input.ImageBytes)
	if err != nil {
		log.Printf("[ERROR] Service Error - FaceCrawler: %v", err)
		return nil, err
	}

	if statusCode != 200 {
		log.Printf("[ERROR] Service Error - FaceCrawler: %v", statusCode)
		return nil, fmt.Errorf("[ERROR] Service Error - FaceCrawler: %v", statusCode)
	}

	var full FaceCrawlerStartResponse
	if err := json.Unmarshal(responseBytes, &full); err != nil {
		log.Printf("[ERROR] Service Error - FaceCrawler: %v", err)
		return nil, fmt.Errorf("invalid JSON from service: %w", err)
	}

	clean := FaceCrawlerStartResult{
		Engine:   full.Engine,
		IDSearch: full.JobID.IDSearch,
		Message:  full.JobID.Message,
	}

	return clean, nil
}

func (fc FaceCrawler) Check(jobID string) (any, error) {
	payload := map[string]any{
		"id_search":     jobID,
		"with_progress": true,
		"status_only":   true,
		"demo":          false,
	}

	log.Printf("[CHECKING PROGRESS] FaceCrawler...")
	serviceUrl := os.Getenv("SITE_URL") + "api/search"

	respBytes, statusCode, err := fc.GetImageStatus(serviceUrl, payload)
	if err != nil {
		return FaceCrawlerStatusResponse{}, err
	}

	if statusCode != 200 {
		return FaceCrawlerStatusResponse{}, fmt.Errorf("FaceCrawler returned status %d", statusCode)
	}

	var result FaceCrawlerStatusResponse
	if err := json.Unmarshal(respBytes, &result); err != nil {
		return FaceCrawlerStatusResponse{}, err
	}

	return result, nil
}

func (fc FaceCrawler) GetImageStatus(serviceUrl string, payload map[string]any) ([]byte, int, error) {
	jsonBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, 0, err
	}

	req, err := http.NewRequest("POST", serviceUrl, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return nil, 0, err
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := fc.Client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}

	return responseData, resp.StatusCode, nil
}

func (fc FaceCrawler) sendImageBytes(serviceUrl string, headers map[string]string, image []byte) ([]byte, int, error) {
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("images", "upload.jpg")
	if err != nil {
		return nil, 0, err
	}

	_, err = part.Write(image)
	if err != nil {
		return nil, 0, err
	}

	writer.Close()
	req, err := http.NewRequest("POST", serviceUrl, body)
	if err != nil {
		return nil, 0, err
	}

	for k, v := range headers {
		req.Header.Set(k, v)
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())
	resp, err := fc.Client.Do(req)
	if err != nil {
		return nil, 0, err
	}
	defer resp.Body.Close()

	responseData, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, err
	}
	return responseData, resp.StatusCode, nil
}
