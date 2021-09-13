package helpers

import (
	"html/template"

	"github.com/rvflash/elapsed"
	"github.com/ugent-library/biblio-backend/internal/engine"
)

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"searchArgs":  engine.NewSearchArgs,
		"timeElapsed": elapsed.LocalTime,
	}
}
