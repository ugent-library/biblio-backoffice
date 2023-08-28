package identifiers

var (
	BioStudies    = &BioStudiesType{}
	DOI           = &DOIType{}
	ENA           = &ENAType{}
	ENABioProject = &ENABioProjectType{}
	Ensembl       = &EnsemblType{}
	Handle        = &HandleType{}
	PubMed        = &PubMedType{}
	WebOfScience  = &WebOfScienceType{}
)

var types = map[string]Type{
	"BioStudies":    BioStudies,
	"DOI":           DOI,
	"ENA":           ENA,
	"ENABioProject": ENABioProject,
	"Ensembl":       Ensembl,
	"Handle":        Handle,
	"PubMed":        PubMed,
	"WebOfScience":  WebOfScience,
}

type Type interface {
	Validate(string) bool
	Normalize(string) (string, error)
	Resolve(string) string
}

func Resolve(name, id string) string {
	t, ok := types[name]
	if !ok {
		return ""
	}
	return t.Resolve(id)
}
