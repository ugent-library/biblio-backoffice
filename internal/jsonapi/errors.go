package jsonapi

import "encoding/json"

type Errors []Error

type Error struct {
	ID     string      `json:"id,omitempty"`
	Status int         `json:"status,omitempty"`
	Code   string      `json:"code,omitempty"`
	Title  string      `json:"title,omitempty"`
	Detail string      `json:"detail,omitempty"`
	Source ErrorSource `json:"source,omitempty"`
}

type ErrorSource struct {
	Pointer string `json:"pointer,omitempty"`
}

func (e Errors) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}
