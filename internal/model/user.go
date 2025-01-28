package model

import "strings"

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
	RoleIds  []int64             `json:"role_ids"`
}

func (s *Session) ObjectPermission(name string) *SessionPermission {
	for _, v := range s.Scopes {
		if v.Class == name {
			return &v
		}
	}
	return nil
}

func (sp *SessionPermission) HasRead() bool {
	// TODO
	if !sp.Obac || strings.Index(sp.Access, "r") > -1 {
		return true
	}
	return false
}
