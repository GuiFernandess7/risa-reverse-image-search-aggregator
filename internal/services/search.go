package services

import (
	"fmt"
)

type SearchInput struct {
	ImageBytes []byte
	ImageURL   string
}

type SearchService interface {
	Name() string
	Search(input SearchInput) (any, error)
}

type AsyncSearchService interface {
	Name() string
	Start(SearchInput) (any, error)
	Check(string) (any, error)
}

func GetEngine(name string) (SearchService, AsyncSearchService, error) {
	switch name {

	case "facecrawler":
		return nil, NewFaceCrawler(), nil

	case "yandex":
		return NewYandexSearch(), nil, nil
	}

	return nil, nil, fmt.Errorf("invalid engine")
}
