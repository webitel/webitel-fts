package model

import (
	"encoding/json"
	"fmt"
)

type MessageId string

type Message struct {
	Id         MessageId       `json:"id"`
	DomainId   int64           `json:"domain_id,omitempty"`
	ObjectName string          `json:"object_name,omitempty"`
	Date       int64           `json:"date,omitempty"`
	Body       json.RawMessage `json:"body,omitempty"`
}

func (id *MessageId) UnmarshalJSON(data []byte) error {
	var raw any
	if err := json.Unmarshal(data, &raw); err != nil {
		return fmt.Errorf("cannot unmarshal MessageId: %w", err)
	}

	switch raw.(type) {
	case float64:
		raw = int64(raw.(float64))
	}

	*id = MessageId(fmt.Sprintf("%v", raw))
	return nil
}
