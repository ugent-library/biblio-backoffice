package models

type Organization struct {
	ID   string         `json:"id,omitempty"`
	Name string         `json:"name,omitempty"`
	Tree []Organization `json:"tree,omitempty"`
}
