package identifiers

type URN struct{}

func (i *URN) Validate(id string) bool {
	return true
}

func (i *URN) Normalize(id string) (string, error) {
	return id, nil
}

func (i *URN) Resolve(id string) string {
	return "https://nbn-resolving.org/" + id
}
