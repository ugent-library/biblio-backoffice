package models

type OrganizationTreeElement struct {
	ID string `json:"id,omitempty"`
}

type Organization struct {
	ID   string                    `json:"id,omitempty"`
	Name string                    `json:"name,omitempty"`
	Tree []OrganizationTreeElement `json:"tree,omitempty"`
}

type RelatedOrganization struct {
	OrganizationID string        `json:"organization_id,omitempty"`
	Organization   *Organization `json:"-"`
}
