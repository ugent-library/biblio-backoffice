package plato

import "github.com/ugent-library/biblio-backoffice/recordsources"

func init() {
	recordsources.Register("plato", New)
}

func New(conn string) (recordsources.Source, error) {
	return &platoSource{}, nil
}

type platoSource struct {
}
