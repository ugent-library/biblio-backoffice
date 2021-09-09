package presenters

import "github.com/ugent-library/biblio-backend/internal/fields"

type Presenter interface {
	Process() map[string]fields.FieldSet
}
