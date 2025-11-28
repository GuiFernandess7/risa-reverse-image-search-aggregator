package search

import (
	"fmt"

	facecrawler "github.com/GuiFernandess7/risa/internal/services/facecrawler"
	yandex "github.com/GuiFernandess7/risa/internal/services/yandex"
)

func GetEngine(name string) (SearchService, AsyncSearchService, error) {
	switch name {

	case "facecrawler":
		return nil, facecrawler.NewFaceCrawler(), nil

	case "yandex":
		return yandex.NewYandexSearch(), nil, nil
	}

	return nil, nil, fmt.Errorf("invalid engine")
}
