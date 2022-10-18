package models

import "time"

type UserDepartment struct {
	ID string `json:"_id"`
}

type User struct {
	Active      bool             `json:"active"`
	DateCreated *time.Time       `json:"date_created"`
	DateUpdated *time.Time       `json:"date_updated"`
	Department  []UserDepartment `json:"department"`
	Email       string           `json:"email"`
	FirstName   string           `json:"first_name"`
	FullName    string           `json:"full_name"`
	ID          string           `json:"_id"`
	LastName    string           `json:"last_name"`
	ORCID       string           `json:"orcid"`
	ORCIDToken  string           `json:"orcid_token"`
	Role        string           `json:"role"`
	UGentID     []string         `json:"ugent_id"`
	Username    string           `json:"username"`
}

func (u *User) CanViewPublication(p *Publication) bool {
	if p.Status == "deleted" {
		return false
	}
	if u.CanCurate() {
		return true
	}
	if p.Creator != nil && p.Creator.ID == u.ID {
		return true
	}
	for _, c := range p.Author {
		if c.ID == u.ID {
			return true
		}
	}
	for _, c := range p.Editor {
		if c.ID == u.ID {
			return true
		}
	}
	for _, c := range p.Supervisor {
		if c.ID == u.ID {
			return true
		}
	}
	return false
}

func (u *User) CanEditPublication(p *Publication) bool {
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
	if p.Creator != nil && p.Creator.ID == u.ID {
		return true
	}
	for _, c := range p.Author {
		if c.ID == u.ID {
			return true
		}
	}
	for _, c := range p.Editor {
		if c.ID == u.ID {
			return true
		}
	}
	for _, c := range p.Supervisor {
		if c.ID == u.ID {
			return true
		}
	}
	return false
}

func (u *User) CanDeletePublication(p *Publication) bool {
	if p.Status == "deleted" {
		return false
	}
	if u.CanCurate() {
		return true
	}
	if p.Locked {
		return false
	}
	if p.Status == "private" && p.Creator != nil && p.Creator.ID == u.ID {
		return true
	}
	return false
}

func (u *User) CanViewDataset(d *Dataset) bool {
	if d.Status == "deleted" {
		return false
	}
	if u.CanCurate() {
		return true
	}
	if d.Creator != nil && d.Creator.ID == u.ID {
		return true
	}
	for _, c := range d.Author {
		if c.ID == u.ID {
			return true
		}
	}
	for _, c := range d.Contributor {
		if c.ID == u.ID {
			return true
		}
	}
	return false
}

func (u *User) CanEditDataset(d *Dataset) bool {
	if d.Status == "deleted" {
		return false
	}
	if u.CanCurate() {
		return true
	}
	if d.Locked {
		return false
	}
	if d.Creator != nil && d.Creator.ID == u.ID {
		return true
	}
	for _, c := range d.Author {
		if c.ID == u.ID {
			return true
		}
	}
	for _, c := range d.Contributor {
		if c.ID == u.ID {
			return true
		}
	}
	return false
}

func (u *User) CanDeleteDataset(d *Dataset) bool {
	if d.Status == "deleted" {
		return false
	}
	if u.Role == "admin" {
		return true
	}
	if d.Locked {
		return false
	}
	if d.Status == "private" && d.Creator != nil && d.Creator.ID == u.ID {
		return true
	}
	return false
}

func (u *User) CanImpersonateUser() bool {
	return u.Role == "admin"
}

func (u *User) CanCurate() bool {
	return u.Role == "admin"
}

func (u *User) CanViewDashboard() bool {
	return u.Role == "admin"
}

func (u *User) CanChangeType(p *Publication) bool {
	return u.CanEditPublication(p) && p.Status != "public"
}
