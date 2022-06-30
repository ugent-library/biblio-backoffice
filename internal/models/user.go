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
	if u.Role == "admin" {
		return true
	}
	if p.CreatorID == u.ID {
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
	if u.Role == "admin" {
		return true
	}
	if p.Locked {
		return false
	}
	if p.CreatorID == u.ID {
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

func (u *User) CanPublishPublication(p *Publication) bool {
	return u.CanEditPublication(p) && p.Status != "public"
}

func (u *User) CanDeletePublication(p *Publication) bool {
	if p.Status == "deleted" {
		return false
	}
	if u.Role == "admin" {
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

func (u *User) CanViewDataset(d *Dataset) bool {
	if d.Status == "deleted" {
		return false
	}
	if u.Role == "admin" {
		return true
	}
	if d.CreatorID == u.ID {
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
	if u.Role == "admin" {
		return true
	}
	if d.Locked {
		return false
	}
	if d.CreatorID == u.ID {
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

func (u *User) CanPublishDataset(d *Dataset) bool {
	return u.CanEditDataset(d) && d.Status != "public"
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
	if d.Status == "private" && d.CreatorID == u.ID {
		return true
	}
	return false
}

func (u *User) CanImpersonateUser() bool {
	return u.Role == "admin"
}

func (u *User) CanCuratePublications() bool {
	return u.Role == "admin"
}

func (u *User) CanCurateDatasets() bool {
	return u.Role == "admin"
}
