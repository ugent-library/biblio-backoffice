package models

import (
	"fmt"

	"github.com/ugent-library/biblio-backoffice/validation"
	"github.com/ugent-library/biblio-backoffice/vocabularies"
)

type Contributor struct {
	PersonID       string          `json:"person_id,omitempty"`
	Person         *Person         `json:"-"`
	ExternalPerson *ExternalPerson `json:"external_person,omitempty"`
	CreditRole     []string        `json:"credit_role,omitempty"`
}

func ContributorFromPerson(p *Person) *Contributor {
	return &Contributor{
		PersonID: p.ID,
		Person:   p,
	}
}

func ContributorFromFirstLastName(fn, ln string) *Contributor {
	return &Contributor{
		ExternalPerson: &ExternalPerson{
			FullName:  fn + " " + ln,
			FirstName: fn,
			LastName:  ln,
		},
	}
}

func (c *Contributor) Name() string {
	if c.Person != nil {
		return c.Person.FullName
	}
	if c.ExternalPerson != nil {
		return c.ExternalPerson.FullName
	}
	return ""
}

func (c *Contributor) FirstName() string {
	if c.Person != nil {
		return c.Person.FirstName
	}
	if c.ExternalPerson != nil {
		return c.ExternalPerson.FirstName
	}
	return ""
}

func (c *Contributor) LastName() string {
	if c.Person != nil {
		return c.Person.LastName
	}
	if c.ExternalPerson != nil {
		return c.ExternalPerson.LastName
	}
	return ""
}

func (c *Contributor) ORCID() string {
	if c.Person != nil {
		return c.Person.ORCID
	}
	return ""
}

func (c *Contributor) Validate() (errs validation.Errors) {
	if c.ExternalPerson != nil {
		if c.ExternalPerson.FullName == "" {
			errs = append(errs, &validation.Error{
				Pointer: "/full_name",
				Code:    "full_name.required",
			})
		}
		if c.ExternalPerson.FirstName == "" {
			errs = append(errs, &validation.Error{
				Pointer: "/first_name",
				Code:    "first_name.required",
			})
		}
		if c.ExternalPerson.LastName == "" {
			errs = append(errs, &validation.Error{
				Pointer: "/last_name",
				Code:    "last_name.required",
			})
		}
	}
	for i, cr := range c.CreditRole {
		if !validation.InArray(vocabularies.Map["credit_roles"], cr) {
			errs = append(errs, &validation.Error{
				Pointer: fmt.Sprintf("/credit_role/%d", i),
				Code:    "credit_role.invalid",
			})
		}
	}

	return
}
