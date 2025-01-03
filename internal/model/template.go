package model

import "encoding/json"

type Template struct {
	Name string          `json:"name"`
	Data json.RawMessage `json:"data"`
}
