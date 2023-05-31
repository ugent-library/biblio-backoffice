package identifiers

type BioStudies struct{}

func (i *BioStudies) Validate(id string) bool {
	return true
}

func (i *BioStudies) Normalize(id string) (string, error) {
	return id, nil
}

func (i *BioStudies) Resolve(id string) string {
	return "https://www.ebi.ac.uk/biostudies/studies/" + id
}
