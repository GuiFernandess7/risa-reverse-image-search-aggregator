package types

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
