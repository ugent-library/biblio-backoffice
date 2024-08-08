package flash

type Flash struct {
	Type         string
	Application  string
	Level        string
	Title        string
	Body         string
	Dismissible  bool
	DismissAfter uint
}

func SimpleFlash() *Flash {
	return &Flash{
		Type:         "simple",
		Application:  "Biblio",
		Level:        "info",
		Dismissible:  true,
		DismissAfter: 5000,
	}
}

func ComplexFlash() *Flash {
	return &Flash{
		Type:         "complex",
		Application:  "Biblio",
		Level:        "info",
		Dismissible:  true,
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

func (f *Flash) WithBody(body string) *Flash {
	f.Body = body
	return f
}

func (f *Flash) IsDismissible(dismissible bool) *Flash {
	f.Dismissible = dismissible
	return f
}

func (f *Flash) DismissedAfter(milliseconds uint) *Flash {
	f.DismissAfter = milliseconds
	return f
}
