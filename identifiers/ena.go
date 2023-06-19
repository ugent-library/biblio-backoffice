package identifiers

type ENAType struct{}

func (i *ENAType) Validate(id string) bool {
	return true
}

func (i *ENAType) Normalize(id string) (string, error) {
	return id, nil
}

func (i *ENAType) Resolve(id string) string {
	return "https://www.ebi.ac.uk/ena/browser/view/" + id
}
