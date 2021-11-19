package models

type Contributor struct {
	ID         string   `json:"id,omitempty" form:"-"`
	ORCID      string   `json:"orcid,omitempty" form:"-"`
	UGentID    []string `json:"ugent_id,omitempty" form:"-"`
	FirstName  string   `json:"first_name,omitempty" form:"first_name"`
	LastName   string   `json:"last_name,omitempty" form:"last_name"`
	FullName   string   `json:"full_name,omitempty" form:"-"`
	CreditRole []string `json:"credit_role,omitempty" form:"credit_role"`
}

func (p *Contributor) CreditRoleChoices() []string {
	return []string{
		"first_author",
		"last_author",
		"conceptualization",
		"data_curation",
		"formal_analysis",
		"funding_acquisition",
		"investigation",
		"methodology",
		"project_administration",
		"resources",
		"software",
		"supervision",
		"validation",
		"visualization",
		"writing_original_draft",
		"writing_review_editing",
	}
}
