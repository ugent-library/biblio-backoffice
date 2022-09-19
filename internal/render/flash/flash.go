package flash

import (
	"html/template"
)

type Flash struct {
	Type         string
	Application  string
	Level        string
	Title        string
	Body         template.HTML
	Dismissable  bool
	DismissAfter uint
}

func SimpleFlash() Flash {
	return Flash{
		Type:         "simple",
		Application:  "Biblio",
		Level:        "info",
		DismissAfter: 1000,
	}
}

func ComplexFlash() Flash {
	return Flash{
		Type:         "complex",
		Application:  "Biblio",
		Level:        "info",
		DismissAfter: 1000,
	}
}
