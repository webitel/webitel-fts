package model

type SearchResult struct {
	Id    any
	Scope string
	Text  string
}

type SearchQuery struct {
	Limit int
	Q     string
	Scope []string
}
