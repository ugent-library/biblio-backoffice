package identifiers

import "github.com/caltechlibrary/doitools"

type DOIType struct{}

func (i *DOIType) Validate(id string) bool {
	return true
}

func (i *DOIType) Normalize(id string) (string, error) {
	return doitools.NormalizeDOI(id)
}

func (i *DOIType) Resolve(id string) string {
	return "https://doi.org/" + id
}
