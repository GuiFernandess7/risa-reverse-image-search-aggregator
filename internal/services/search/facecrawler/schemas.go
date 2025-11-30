package facecrawler

import (
	"net/http"
)

type BaseFaceCrawlerResponse struct {
	IDSearch string `json:"id_search"`
	Message  string `json:"message"`
	Progress any    `json:"progress"`
	Error    any    `json:"error"`
	Code     any    `json:"code"`
	Output   Output `json:"output"`
}

type Output struct {
	Items []Item `json:"items"`
}

type FaceCrawlerStartResult struct {
	IDSearch string `json:"id_search"`
	Message  string `json:"message"`
}

type Item struct {
	URL   string `json:"url"`
	Score int    `json:"score"`
}

type FaceCrawler struct {
	Client *http.Client
}
