package identifiers

type ENABioProjectType struct{}

func (i *ENABioProjectType) Validate(id string) bool {
	return true
}

func (i *ENABioProjectType) Normalize(id string) (string, error) {
	return id, nil
}

func (i *ENABioProjectType) Resolve(id string) string {
	return "https://www.ebi.ac.uk/ena/data/view/" + id
}
