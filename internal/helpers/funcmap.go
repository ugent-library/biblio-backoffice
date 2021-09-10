package helpers

import (
	"html/template"

	"github.com/ugent-library/biblio-backend/internal/engine"
)

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"searchArgs":  engine.NewSearchArgs,
		"timeElapsed": TimeElapsed,
	}
}
