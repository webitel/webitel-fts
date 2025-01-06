package model

type SearchResult struct {
	Id         any
	ObjectName string
	Text       string
}

type SearchQuery struct {
	DomainId    int64
	Limit       int
	Q           string
	ObjectsName []string
}
