package identifiers

var (
	BioProject     = &BioProjectType{}
	BioStudies     = &BioStudiesType{}
	DOI            = &DOIType{}
	ENA            = &ENAType{}
	Ensembl        = &EnsemblType{}
	Handle         = &HandleType{}
	PubMedID       = &PubMedIDType{}
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
	"PubMedID":       PubMedID,
	"URN":            URN,
	"WebOfScienceID": WebOfScienceID,
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
