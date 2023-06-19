package identifiers

type BioProjectType struct{}

func (i *BioProjectType) Validate(id string) bool {
	return true
}

func (i *BioProjectType) Normalize(id string) (string, error) {
	return id, nil
}

func (i *BioProjectType) Resolve(id string) string {
	return "https://www.ebi.ac.uk/ena/data/view/" + id
}
