package models

type User struct {
	Person
	Username   string `json:"username"`
	Role       string `json:"role"`
	ORCIDToken string `json:"orcid_token"`
}

func (u *User) CanViewPublication(p *Publication) bool {
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
	if p.Status == "private" && p.CreatorID == u.ID {
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

func (u *User) CanCurate() bool {
	return u.Role == "admin"
}

func (u *User) CanViewDashboard() bool {
	return u.Role == "admin"
}

func (u *User) CanChangeType(p *Publication) bool {
	return u.CanEditPublication(p) && p.Status != "public"
}
