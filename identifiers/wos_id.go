package identifiers

type WebOfScienceIDType struct{}

func (i *WebOfScienceIDType) Validate(id string) bool {
	return true
}

func (i *WebOfScienceIDType) Normalize(id string) (string, error) {
	return id, nil
}

func (i *WebOfScienceIDType) Resolve(id string) string {
	return "https://www.webofscience.com/wos/woscc/full-record/" + id
}
