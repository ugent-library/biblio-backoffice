package identifiers

type BioStudiesType struct{}

func (i *BioStudiesType) Validate(id string) bool {
	return true
}

func (i *BioStudiesType) Normalize(id string) (string, error) {
	return id, nil
}

func (i *BioStudiesType) Resolve(id string) string {
	return "https://www.ebi.ac.uk/biostudies/studies/" + id
}
