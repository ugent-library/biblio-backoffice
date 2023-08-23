package identifiers

var (
	BioProject     = &BioProjectType{}
	BioStudies     = &BioStudiesType{}
	DOI            = &DOIType{}
	ENA            = &ENAType{}
	Ensembl        = &EnsemblType{}
	Handle         = &HandleType{}
	PubMedID       = &PubMedIDType{}
	WebOfScienceID = &WebOfScienceIDType{}
)

var types = map[string]Type{
	"ENA BioProject": BioProject,
	"BioStudies":     BioStudies,
	"DOI":            DOI,
	"ENA":            ENA,
	"Ensembl":        Ensembl,
	"Handle":         Handle,
	"PubMedID":       PubMedID,
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
