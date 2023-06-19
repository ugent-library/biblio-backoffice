package identifiers

type EnsemblType struct{}

func (i *EnsemblType) Validate(id string) bool {
	return true
}

func (i *EnsemblType) Normalize(id string) (string, error) {
	return id, nil
}

func (i *EnsemblType) Resolve(id string) string {
	return "https://www.ensembl.org/id/" + id
}
