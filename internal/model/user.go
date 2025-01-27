package model

type SessionPermission struct {
	Id     int64
	Class  string
	Obac   bool
	Rbac   bool
	Access string
}

type Session struct {
	Id       string              `json:"id"`
	Name     string              `json:"name"`
	DomainId int64               `json:"domain_id"`
	Expire   int64               `json:"expire"`
	UserId   int64               `json:"user_id"`
	Scopes   []SessionPermission `json:"scopes"`
}
