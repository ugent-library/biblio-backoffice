package helpers

import (
	"html/template"
)

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"timeElapsed": timeElapsed,
	}
}
