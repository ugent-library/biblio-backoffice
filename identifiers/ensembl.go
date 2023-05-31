package identifiers

type Ensembl struct{}

func (i *Ensembl) Validate(id string) bool {
	return true
}

func (i *Ensembl) Normalize(id string) (string, error) {
	return id, nil
}

func (i *Ensembl) Resolve(id string) string {
	return "https://www.ensembl.org/id/" + id
}
