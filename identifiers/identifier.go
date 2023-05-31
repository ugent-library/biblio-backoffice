package identifiers

var types = map[string]Type{
	"BioProject": &BioProject{},
	"BioStudies": &BioStudies{},
	"DOI":        &DOI{},
	"ENA":        &ENA{},
	"Ensembl":    &Ensembl{},
	"Handle":     &Handle{},
	"URN":        &URN{},
}

type Type interface {
	Validate(string) bool
	Normalize(string) (string, error)
	Resolve(string) string
}

func GetType(name string) Type {
	return types[name]
}
