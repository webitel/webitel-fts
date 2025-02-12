package model

type SearchResult struct {
	Id         any
	ObjectName string
	Text       string
}

type ObjectName struct {
	Name    string
	RoleIds []int64
}

type SearchQuery struct {
	DomainId    int64
	Limit       int
	Page        int
	Q           string
	ObjectsName []ObjectName
}

func (o *ObjectName) String() string {
	return o.Name
}

func (sq *SearchQuery) HasObject(name string) bool {
	for _, v := range sq.ObjectsName {
		if v.Name == name {
			return true
		}
	}

	return false
}
