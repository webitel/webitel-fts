package searchengine

import "context"

type SearchEngine interface {
	Shutdown() error
	Test() error

	Insert(ctx context.Context, id string, index string, body []byte) error
	Update(ctx context.Context, id string, index string, body []byte) error
	Delete(ctx context.Context, id string, index string) error

	Template(ctx context.Context, name string, body []byte) error

	Search(ctx context.Context, IndexName []string, text string, limit int) ([]SearchResult, error)
}

type SearchResult struct {
	Index  string
	Id     string
	Text   string
	Source map[string]any `json:"_source"`
}
