package models

import "github.com/ugent-library/biblio-backend/internal/validation"

type Contributor struct {
	CreditRole []string `json:"credit_role,omitempty" form:"credit_role"`
	FirstName  string   `json:"first_name,omitempty" form:"first_name"`
	FullName   string   `json:"full_name,omitempty" form:"-"` // TODO rename to Name
	ID         string   `json:"id,omitempty" form:"ID"`
	LastName   string   `json:"last_name,omitempty" form:"last_name"`
	ORCID      string   `json:"orcid,omitempty" form:"-"`
	UGentID    []string `json:"ugent_id,omitempty" form:"-"`
}

func (p *Contributor) Clone() *Contributor {
	clone := *p
	clone.CreditRole = append(clone.CreditRole, p.CreditRole...)
	clone.UGentID = append(clone.UGentID, p.UGentID...)
	return &clone
}

func (p *Contributor) HasCreditRole(role string) bool {
	for _, r := range p.CreditRole {
		if r == role {
			return true
		}
	}
	return false
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

func (c *Contributor) Validate() (errs validation.Errors) {
	if c.FirstName == "" {
		errs = append(errs, &validation.Error{
			Pointer: "/first_name",
			Code:    "required",
		})
	}
	if c.LastName == "" {
		errs = append(errs, &validation.Error{
			Pointer: "/last_name",
			Code:    "required",
		})
	}
	return
}
