package identifiers

type WebOfScienceType struct{}

func (i *WebOfScienceType) Validate(id string) bool {
	return true
}

func (i *WebOfScienceType) Normalize(id string) (string, error) {
	return id, nil
}

func (i *WebOfScienceType) Resolve(id string) string {
	return "https://www.webofscience.com/wos/woscc/full-record/" + id
}
