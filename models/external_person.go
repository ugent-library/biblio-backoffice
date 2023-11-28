package models

type ExternalPerson struct {
	FullName  string `json:"full_name,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
	// TODO forms are not yet aware of these new fields
	HonorificPrefix string `json:"honorific_prefix,omitempty"`
	Affiliation     string `json:"affiliation,omitempty"`
}
