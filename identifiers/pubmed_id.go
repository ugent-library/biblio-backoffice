package identifiers

type PubMedIDType struct{}

func (i *PubMedIDType) Validate(id string) bool {
	return true
}

func (i *PubMedIDType) Normalize(id string) (string, error) {
	return id, nil
}

func (i *PubMedIDType) Resolve(id string) string {
	return "https://www.ncbi.nlm.nih.gov/pubmed/" + id
}
