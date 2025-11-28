package yandex

import (
	"net/http"
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

type YandexSearchResponse struct {
	ImageResults []ImageResult `json:"image_results,omitempty"`
}

type ImageResult struct {
	Title   string `json:"title,omitempty"`
	Snippet string `json:"snippet,omitempty"`
	Link    string `json:"link,omitempty"`
	Source  string `json:"source,omitempty"`
}

type YandexSearch struct {
	Client *http.Client
}
