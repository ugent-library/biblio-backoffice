package models

import (
	"time"
)

type Affiliation struct {
	OrganizationID string        `json:"organization_id"`
	Organization   *Organization `json:"-"`
}

type Person struct {
	Active       bool           `json:"active"`
	DateCreated  *time.Time     `json:"date_created"`
	DateUpdated  *time.Time     `json:"date_updated"`
	Affiliations []*Affiliation `json:"affiliations"`
	Email        string         `json:"email"`
	FirstName    string         `json:"first_name"`
	FullName     string         `json:"full_name"`
	ID           string         `json:"id"`
	LastName     string         `json:"last_name"`
	ORCID        string         `json:"orcid"`
	UGentID      []string       `json:"ugent_id"`
	// fields below are only relevant for Active people (users)
	Username   string `json:"username"`
	Role       string `json:"role"`
	ORCIDToken string `json:"orcid_token"`
}

func (u *Person) CanViewPublication(p *Publication) bool {
	if !u.Active {
		return false
	}
	if p.Status == "deleted" {
		return false
	}
	if u.CanCurate() {
		return true
	}
	if p.CreatorID == u.ID {
		return true
	}
	for _, c := range p.Author {
		if c.PersonID == u.ID {
			return true
		}
	}
	for _, c := range p.Editor {
		if c.PersonID == u.ID {
			return true
		}
	}
	for _, c := range p.Supervisor {
		if c.PersonID == u.ID {
			return true
		}
	}
	return false
}

func (u *Person) CanWithdrawPublication(p *Publication) bool {
	return p.Status == "public" && u.CanEditPublication(p)
}

func (u *Person) CanPublishPublication(p *Publication) bool {
	return p.Status != "public" && u.CanEditPublication(p)
}

func (u *Person) CanEditPublication(p *Publication) bool {
	if !u.Active {
		return false
	}
	if p.Status == "deleted" {
		return false
	}
	if u.CanCurate() {
		return true
	}
	if p.Legacy {
		return false
	}
	if p.Locked {
		return false
	}
	if p.CreatorID == u.ID {
		return true
	}
	for _, c := range p.Author {
		if c.PersonID == u.ID {
			return true
		}
	}
	for _, c := range p.Editor {
		if c.PersonID == u.ID {
			return true
		}
	}
	for _, c := range p.Supervisor {
		if c.PersonID == u.ID {
			return true
		}
	}
	return false
}

func (u *Person) CanDeletePublication(p *Publication) bool {
	if !u.Active {
		return false
	}
	if p.Status == "deleted" {
		return false
	}
	if u.CanCurate() {
		return true
	}
	if p.Locked {
		return false
	}
	if p.Status == "private" && p.CreatorID == u.ID {
		return true
	}
	return false
}

func (u *Person) CanViewDataset(d *Dataset) bool {
	if !u.Active {
		return false
	}
	if d.Status == "deleted" {
		return false
	}
	if u.CanCurate() {
		return true
	}
	if d.CreatorID == u.ID {
		return true
	}
	for _, c := range d.Author {
		if c.PersonID == u.ID {
			return true
		}
	}
	for _, c := range d.Contributor {
		if c.PersonID == u.ID {
			return true
		}
	}
	return false
}

func (u *Person) CanWithdrawDataset(d *Dataset) bool {
	return d.Status == "public" && u.CanEditDataset(d)
}

func (u *Person) CanPublishDataset(d *Dataset) bool {
	return d.Status != "public" && u.CanEditDataset(d)
}

func (u *Person) CanEditDataset(d *Dataset) bool {
	if !u.Active {
		return false
	}
	if d.Status == "deleted" {
		return false
	}
	if u.CanCurate() {
		return true
	}
	if d.Locked {
		return false
	}
	if d.CreatorID == u.ID {
		return true
	}
	for _, c := range d.Author {
		if c.PersonID == u.ID {
			return true
		}
	}
	for _, c := range d.Contributor {
		if c.PersonID == u.ID {
			return true
		}
	}
	return false
}

func (u *Person) CanDeleteDataset(d *Dataset) bool {
	if !u.Active {
		return false
	}
	if d.Status == "deleted" {
		return false
	}
	if u.Role == "admin" {
		return true
	}
	if d.Locked {
		return false
	}
	if d.Status == "private" && d.CreatorID == u.ID {
		return true
	}
	return false
}

func (u *Person) CanImpersonateUser() bool {
	return u.Active && u.Role == "admin"
}

func (u *Person) CanCurate() bool {
	return u.Active && u.Role == "admin"
}

func (u *Person) CanViewDashboard() bool {
	return u.Active && u.Role == "admin"
}

func (u *Person) CanChangeType(p *Publication) bool {
	return u.CanEditPublication(p) && p.Status != "public"
}
