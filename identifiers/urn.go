package identifiers

type URNType struct{}

func (i *URNType) Validate(id string) bool {
	return true
}

func (i *URNType) Normalize(id string) (string, error) {
	return id, nil
}

func (i *URNType) Resolve(id string) string {
	return "https://nbn-resolving.org/" + id
}
