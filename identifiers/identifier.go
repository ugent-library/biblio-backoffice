package identifiers

var (
	BioProject     = &BioProjectType{}
	BioStudies     = &BioStudiesType{}
	DOI            = &DOIType{}
	ENA            = &ENAType{}
	Ensembl        = &EnsemblType{}
	Handle         = &HandleType{}
	URN            = &URNType{}
	WebOfScienceID = &WebOfScienceIDType{}
)

var types = map[string]Type{
	"BioProject":     BioProject,
	"BioStudies":     BioStudies,
	"DOI":            DOI,
	"ENA":            ENA,
	"Ensembl":        Ensembl,
	"Handle":         Handle,
	"URN":            URN,
	"WebOfScienceID": WebOfScienceID,
}

type Type interface {
	Validate(string) bool
	Normalize(string) (string, error)
	Resolve(string) string
}

func GetType(name string) Type {
	return types[name]
}
