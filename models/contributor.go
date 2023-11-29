package models

import (
	"fmt"

	"slices"

	"github.com/ugent-library/biblio-backoffice/vocabularies"
	"github.com/ugent-library/okay"
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

func (c *Contributor) Validate() error {
	errs := okay.NewErrors()

	if c.ExternalPerson != nil {
		if c.ExternalPerson.FullName == "" {
			errs.Add(okay.NewError("/full_name", "full_name.required"))
		}
		if c.ExternalPerson.FirstName == "" {
			errs.Add(okay.NewError("/first_name", "first_name.required"))
		}
		if c.ExternalPerson.LastName == "" {
			errs.Add(okay.NewError("/last_name", "last_name.required"))
		}
	}
	for i, cr := range c.CreditRole {
		if !slices.Contains(vocabularies.Map["credit_roles"], cr) {
			errs.Add(okay.NewError(fmt.Sprintf("/credit_role/%d", i), "credit_role.invalid"))
		}
	}

	return errs.ErrorOrNil()
}
