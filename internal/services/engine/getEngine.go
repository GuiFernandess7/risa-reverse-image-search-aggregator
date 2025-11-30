package engines

import (
	"fmt"

	interfaces "github.com/GuiFernandess7/risa/internal/repository/interfaces"
	"github.com/GuiFernandess7/risa/internal/services/search/facecrawler"
	"github.com/GuiFernandess7/risa/internal/services/search/yandex"
)

func GetEngine(name string) (interfaces.SearchService, interfaces.AsyncSearchService, error) {
	switch name {

	case "facecrawler":
		return nil, facecrawler.NewFaceCrawler(), nil

	case "yandex":
		return yandex.NewYandexSearch(), nil, nil
	}

	return nil, nil, fmt.Errorf("invalid engine")
}
