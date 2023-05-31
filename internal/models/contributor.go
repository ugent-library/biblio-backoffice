package models

import (
	"github.com/ugent-library/biblio-backoffice/internal/validation"
)

type ExternalPerson struct {
	Name      string `json:"name,omitempty"`
	FirstName string `json:"first_name,omitempty"`
	LastName  string `json:"last_name,omitempty"`
}

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
			Name:      fn + " " + ln,
			FirstName: fn,
			LastName:  ln,
		},
	}
}

func (c *Contributor) Validate() (errs validation.Errors) {
	return nil
}

// func (c *Contributor) SetPerson(p *Person) {
// 	c.ExternalPerson = nil
// 	c.PersonID = p.ID
// 	c.Person = p
// }

// func (c *Contributor) SetExternalPerson(p *ExternalPerson) {
// 	c.ExternalPerson = p
// 	c.PersonID = ""
// 	c.Person = nil
// }

// func (c *Contributor) Name() string {
// 	if c.Person != nil {
// 		return c.Person.FullName
// 	}
// 	if c.ExternalPerson != nil {
// 		return c.ExternalPerson.FullName
// 	}
// 	return ""
// }

// type ContributorDepartment struct {
// 	ID   string `json:"id"`
// 	Name string `json:"name"`
// }

// TODO only name should be required (support corporate names)
// type Contributor struct {
// 	CreditRole []string                `json:"credit_role,omitempty"`
// 	FirstName  string                  `json:"first_name,omitempty"`
// 	FullName   string                  `json:"full_name,omitempty"` // TODO rename to Name
// 	ID         string                  `json:"id,omitempty"`
// 	LastName   string                  `json:"last_name,omitempty"`
// 	ORCID      string                  `json:"orcid,omitempty"`
// 	UGentID    []string                `json:"ugent_id,omitempty"`
// 	Department []ContributorDepartment `json:"department,omitempty"`
// }

// func (c *Contributor) Validate() (errs validation.Errors) {
// 	if c.FirstName == "" {
// 		errs = append(errs, &validation.Error{
// 			Pointer: "/first_name",
// 			Code:    "first_name.required",
// 		})
// 	}
// 	if c.LastName == "" {
// 		errs = append(errs, &validation.Error{
// 			Pointer: "/last_name",
// 			Code:    "last_name.required",
// 		})
// 	}
// 	for i, cr := range c.CreditRole {
// 		if !validation.InArray(vocabularies.Map["credit_roles"], cr) {
// 			errs = append(errs, &validation.Error{
// 				Pointer: fmt.Sprintf("/credit_role/%d", i),
// 				Code:    "credit_role.invalid",
// 			})
// 		}
// 	}

// 	return
// }
