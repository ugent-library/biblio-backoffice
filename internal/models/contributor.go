package models

import (
	"fmt"

	"github.com/ugent-library/biblio-backend/internal/validation"
	"github.com/ugent-library/biblio-backend/internal/vocabularies"
)

// TODO only name should be required (support corporate names)
type Contributor struct {
	CreditRole []string `json:"credit_role,omitempty"`
	FirstName  string   `json:"first_name,omitempty"`
	FullName   string   `json:"full_name,omitempty"` // TODO rename to Name
	ID         string   `json:"id,omitempty"`
	LastName   string   `json:"last_name,omitempty"`
	ORCID      string   `json:"orcid,omitempty"`
	UGentID    []string `json:"ugent_id,omitempty"`
}

// TODO remove
func (p *Contributor) Clone() *Contributor {
	clone := *p
	clone.CreditRole = nil
	clone.CreditRole = append(clone.CreditRole, p.CreditRole...)
	clone.UGentID = nil
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

func IsValidCreditRole(role string) bool {
	return validation.InArray(vocabularies.Map["credit_roles"], role)
}

func (c *Contributor) Validate() (errs validation.Errors) {
	if c.FirstName == "" {
		errs = append(errs, &validation.Error{
			Pointer: "/first_name",
			Code:    "first_name.required",
		})
	}
	if c.LastName == "" {
		errs = append(errs, &validation.Error{
			Pointer: "/last_name",
			Code:    "last_name.required",
		})
	}
	for i, cr := range c.CreditRole {
		if !IsValidCreditRole(cr) {
			errs = append(errs, &validation.Error{
				Pointer: fmt.Sprintf("/credit_role/%d", i),
				Code:    "credit_role.invalid",
			})
		}
	}

	return
}
