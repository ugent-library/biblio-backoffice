package identifiers

import "regexp"

var reHandle = regexp.MustCompile(`^https?://hdl.handle.net/`)

type HandleType struct{}

func (i *HandleType) Validate(id string) bool {
	return true
}

func (i *HandleType) Normalize(id string) (string, error) {
	return reHandle.ReplaceAllString(id, ""), nil
}

func (i *HandleType) Resolve(id string) string {
	return "https://hdl.handle.net/" + id
}
