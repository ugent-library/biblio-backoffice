package models

type Contributor struct {
	ID         string   `json:"id,omitempty"`
	ORCID      string   `json:"orcid,omitempty"`
	UGentID    []string `json:"ugent_id,omitempty"`
	FirstName  string   `json:"first_name,omitempty"`
	LastName   string   `json:"last_name,omitempty"`
	FullName   string   `json:"full_name,omitempty"`
	CreditRole []string `json:"credit_role,omitempty"`
}
