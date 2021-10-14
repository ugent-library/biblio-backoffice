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
	Role        string           `json:"role"`
	UGentID     []string         `json:"ugent_id"`
	Username    string           `json:"username"`
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

func (u *User) CanEditDataset(p *Dataset) bool {
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
	for _, c := range p.Creator {
		if c.ID == u.ID {
			return true
		}
	}
	for _, c := range p.Contributor {
		if c.ID == u.ID {
			return true
		}
	}
	return false
}
