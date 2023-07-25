package flash

import (
	"html/template"
)

type Flash struct {
	Type         string
	Application  string
	Code         string
	Level        string
	Title        string
	Body         template.HTML
	Dismissable  bool
	DismissAfter uint
}

func SimpleFlash() *Flash {
	return &Flash{
		Type:         "simple",
		Application:  "Biblio",
		Level:        "info",
		Dismissable:  true,
		DismissAfter: 5000,
	}
}

func ComplexFlash() *Flash {
	return &Flash{
		Type:         "complex",
		Application:  "Biblio",
		Level:        "info",
		Dismissable:  true,
		DismissAfter: 5000,
	}
}

func (f *Flash) WithLevel(level string) *Flash {
	f.Level = level
	return f
}

func (f *Flash) WithTitle(title string) *Flash {
	f.Title = title
	return f
}

func (f *Flash) WithBody(body template.HTML) *Flash {
	f.Body = body
	return f
}

func (f *Flash) WithApplication(app string) *Flash {
	f.Application = app
	return f
}

func (f *Flash) WithCode(code string) *Flash {
	f.Code = code
	return f
}

func (f *Flash) IsDismissable(dismissable bool) *Flash {
	f.Dismissable = dismissable
	return f
}

func (f *Flash) DismissedAfter(milliseconds uint) *Flash {
	f.DismissAfter = milliseconds
	return f
}
