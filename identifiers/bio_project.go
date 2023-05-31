package identifiers

type BioProject struct{}

func (i *BioProject) Validate(id string) bool {
	return true
}

func (i *BioProject) Normalize(id string) (string, error) {
	return id, nil
}

func (i *BioProject) Resolve(id string) string {
	return "https://www.ebi.ac.uk/ena/data/view/" + id
}
