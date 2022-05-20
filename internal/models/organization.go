package models

type OrganizationRef struct {
	ID string `json:"id,omitempty"`
}

type Organization struct {
	ID   string            `json:"id,omitempty"`
	Name string            `json:"name,omitempty"`
	Tree []OrganizationRef `json:"tree,omitempty"`
}
