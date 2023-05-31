package identifiers

import "github.com/caltechlibrary/doitools"

type DOI struct{}

func (i *DOI) Validate(id string) bool {
	return true
}

func (i *DOI) Normalize(id string) (string, error) {
	return doitools.NormalizeDOI(id)
}

func (i *DOI) Resolve(id string) string {
	return "https://doi.org/" + id
}
