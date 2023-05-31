package identifiers

type ENA struct{}

func (i *ENA) Validate(id string) bool {
	return true
}

func (i *ENA) Normalize(id string) (string, error) {
	return id, nil
}

func (i *ENA) Resolve(id string) string {
	return "https://www.ebi.ac.uk/ena/browser/view/" + id
}
