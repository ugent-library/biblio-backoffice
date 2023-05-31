package identifiers

import "regexp"

var reHandle = regexp.MustCompile(`^https?://hdl.handle.net/`)

type Handle struct{}

func (i *Handle) Validate(id string) bool {
	return true
}

func (i *Handle) Normalize(id string) (string, error) {
	return reHandle.ReplaceAllString(id, ""), nil
}

func (i *Handle) Resolve(id string) string {
	return "https://hdl.handle.net/" + id
}
