package models

type Organization struct {
	ID   string `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
	Tree []struct {
		ID string `json:"id,omitempty"`
	} `json:"tree,omitempty"`
}

type RelatedOrganization struct {
	OrganizationID string        `json:"organization_id,omitempty"`
	Organization   *Organization `json:"-"`
}
