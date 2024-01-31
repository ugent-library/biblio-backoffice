package identifiers

type EGAType struct{}

func (i *EGAType) Validate(id string) bool {
	return true
}

func (i *EGAType) Normalize(id string) (string, error) {
	return id, nil
}

func (i *EGAType) Resolve(id string) string {
	return "https://ega-archive.org/datasets/" + id
}
