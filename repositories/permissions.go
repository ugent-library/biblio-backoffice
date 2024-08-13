package repositories

import "github.com/ugent-library/biblio-backoffice/models"

func (s *Repo) isProxyForPublication(u *models.Person, p *models.Publication) bool {
	var personIDs []string

	if p.CreatorID != "" {
		personIDs = append(personIDs, p.CreatorID)
	}
	for _, c := range p.Author {
		if c.PersonID != "" {
			personIDs = append(personIDs, c.PersonID)
		}
	}
	for _, c := range p.Editor {
		if c.PersonID != "" {
			personIDs = append(personIDs, c.PersonID)
		}
	}
	for _, c := range p.Supervisor {
		if c.PersonID != "" {
			personIDs = append(personIDs, c.PersonID)
		}
	}

	return len(personIDs) > 0 && s.isProxy(u.ID, personIDs)
}

func (s *Repo) CanViewPublication(u *models.Person, p *models.Publication) bool {
	if !u.Active {
		return false
	}
	if p.Status == "deleted" {
		return false
	}
	if s.CanCurate(u) {
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

	return s.isProxyForPublication(u, p)
}

func (s *Repo) CanWithdrawPublication(u *models.Person, p *models.Publication) bool {
	return p.Status == "public" && s.CanEditPublication(u, p)
}

func (s *Repo) CanPublishPublication(u *models.Person, p *models.Publication) bool {
	return p.Status != "public" && s.CanEditPublication(u, p)
}

func (s *Repo) CanEditPublication(u *models.Person, p *models.Publication) bool {
	if !u.Active {
		return false
	}
	if p.Status == "deleted" {
		return false
	}
	if s.CanCurate(u) {
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

	return s.isProxyForPublication(u, p)
}

func (s *Repo) CanDeletePublication(u *models.Person, p *models.Publication) bool {
	if !u.Active {
		return false
	}
	if p.Status == "deleted" {
		return false
	}
	if s.CanCurate(u) {
		return true
	}
	if p.Locked {
		return false
	}
	if p.Status == "private" {
		if p.CreatorID == u.ID {
			return true
		}
		if p.CreatorID != "" && s.isProxy(u.ID, []string{p.CreatorID}) {
			return true
		}
	}

	return false
}

func (s *Repo) isProxyForDataset(u *models.Person, d *models.Dataset) bool {
	var personIDs []string

	if d.CreatorID != "" {
		personIDs = append(personIDs, d.CreatorID)
	}
	for _, c := range d.Author {
		if c.PersonID != "" {
			personIDs = append(personIDs, c.PersonID)
		}
	}
	for _, c := range d.Contributor {
		if c.PersonID != "" {
			personIDs = append(personIDs, c.PersonID)
		}
	}

	return len(personIDs) > 0 && s.isProxy(u.ID, personIDs)
}

func (s *Repo) CanViewDataset(u *models.Person, d *models.Dataset) bool {
	if !u.Active {
		return false
	}
	if d.Status == "deleted" {
		return false
	}
	if s.CanCurate(u) {
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

	return s.isProxyForDataset(u, d)
}

func (s *Repo) CanWithdrawDataset(u *models.Person, d *models.Dataset) bool {
	return d.Status == "public" && s.CanEditDataset(u, d)
}

func (s *Repo) CanPublishDataset(u *models.Person, d *models.Dataset) bool {
	return d.Status != "public" && s.CanEditDataset(u, d)
}

func (s *Repo) CanEditDataset(u *models.Person, d *models.Dataset) bool {
	if !u.Active {
		return false
	}
	if d.Status == "deleted" {
		return false
	}
	if s.CanCurate(u) {
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

	return s.isProxyForDataset(u, d)
}

func (s *Repo) CanDeleteDataset(u *models.Person, d *models.Dataset) bool {
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
	if d.Status == "private" {
		if d.CreatorID == u.ID {
			return true
		}
		if d.CreatorID != "" && s.isProxy(u.ID, []string{d.CreatorID}) {
			return true
		}
	}

	return false
}

func (s *Repo) CanImpersonateUser(u *models.Person) bool {
	return u.Active && u.Role == "admin"
}

func (s *Repo) CanCurate(u *models.Person) bool {
	return u.Active && u.Role == "admin"
}

func (s *Repo) CanViewDashboard(u *models.Person) bool {
	return u.Active && u.Role == "admin"
}

func (s *Repo) CanChangeType(u *models.Person, p *models.Publication) bool {
	return p.Status != "public" && s.CanEditPublication(u, p)
}
