package model

import "encoding/json"

type Message struct {
	Id         any             `json:"id"`
	DomainId   int64           `json:"domain_id,omitempty"`
	ObjectName string          `json:"object_name,omitempty"`
	Date       int64           `json:"date,omitempty"`
	Body       json.RawMessage `json:"body,omitempty"`
}
