package models

type FormError struct {
	Code   string            `json:"code,omitempty"`
	Title  string            `json:"title,omitempty"`
	Source map[string]string `json:"source,omitempty"`
}
