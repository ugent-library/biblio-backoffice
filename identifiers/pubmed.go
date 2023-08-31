package identifiers

type PubMedType struct{}

func (i *PubMedType) Validate(id string) bool {
	return true
}

func (i *PubMedType) Normalize(id string) (string, error) {
	return id, nil
}

func (i *PubMedType) Resolve(id string) string {
	return "https://www.ncbi.nlm.nih.gov/pubmed/" + id
}
