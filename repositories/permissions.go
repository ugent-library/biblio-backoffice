package repositories

import "github.com/ugent-library/biblio-backoffice/models"

func (s *Repo) isProxyForPublication(u *models.Person, p *models.Publication) bool {
	var personIDs []string

	if p.Creator != nil {
		personIDs = append(personIDs, p.Creator.IDs...)
	}
	for _, c := range p.Author {
		if c.Person != nil {
			personIDs = append(personIDs, c.Person.IDs...)
		}
	}
	for _, c := range p.Editor {
		if c.Person != nil {
			personIDs = append(personIDs, c.Person.IDs...)
		}
	}
	for _, c := range p.Supervisor {
		if c.Person != nil {
			personIDs = append(personIDs, c.Person.IDs...)
		}
	}

	return len(personIDs) > 0 && s.IsProxyFor(u.IDs, personIDs)
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
	if p.Creator != nil && intersects(p.Creator.IDs, u.IDs) {
		return true
	}
	for _, c := range p.Author {
		if c.Person != nil && intersects(c.Person.IDs, u.IDs) {
			return true
		}
	}
	for _, c := range p.Editor {
		if c.Person != nil && intersects(c.Person.IDs, u.IDs) {
			return true
		}
	}
	for _, c := range p.Supervisor {
		if c.Person != nil && intersects(c.Person.IDs, u.IDs) {
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
	if p.Creator != nil && intersects(p.Creator.IDs, u.IDs) {
		return true
	}
	for _, c := range p.Author {
		if c.Person != nil && intersects(c.Person.IDs, u.IDs) {
			return true
		}
	}
	for _, c := range p.Editor {
		if c.Person != nil && intersects(c.Person.IDs, u.IDs) {
			return true
		}
	}
	for _, c := range p.Supervisor {
		if c.Person != nil && intersects(c.Person.IDs, u.IDs) {
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
		if p.Creator != nil && intersects(p.Creator.IDs, u.IDs) {
			return true
		}
		if p.Creator != nil && s.IsProxyFor(u.IDs, p.Creator.IDs) {
			return true
		}
	}

	return false
}

func (s *Repo) isProxyForDataset(u *models.Person, d *models.Dataset) bool {
	var personIDs []string

	if d.Creator != nil {
		personIDs = append(personIDs, d.Creator.IDs...)
	}
	for _, c := range d.Author {
		if c.Person != nil {
			personIDs = append(personIDs, c.Person.IDs...)
		}
	}
	for _, c := range d.Contributor {
		if c.Person != nil {
			personIDs = append(personIDs, c.Person.IDs...)
		}
	}

	return len(personIDs) > 0 && s.IsProxyFor(u.IDs, personIDs)
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
	if d.Creator != nil && intersects(d.Creator.IDs, u.IDs) {
		return true
	}
	for _, c := range d.Author {
		if c.Person != nil && intersects(c.Person.IDs, u.IDs) {
			return true
		}
	}
	for _, c := range d.Contributor {
		if c.Person != nil && intersects(c.Person.IDs, u.IDs) {
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
	if d.Creator != nil && intersects(d.Creator.IDs, u.IDs) {
		return true
	}
	for _, c := range d.Author {
		if c.Person != nil && intersects(c.Person.IDs, u.IDs) {
			return true
		}
	}
	for _, c := range d.Contributor {
		if c.Person != nil && intersects(c.Person.IDs, u.IDs) {
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
		if d.Creator != nil && intersects(d.Creator.IDs, u.IDs) {
			return true
		}
		if d.Creator != nil && s.IsProxyFor(u.IDs, d.Creator.IDs) {
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

func intersects(s1 []string, s2 []string) bool {
	for _, e1 := range s1 {
		for _, e2 := range s2 {
			if e1 == e2 {
				return true
			}
		}
	}
	return false
}
