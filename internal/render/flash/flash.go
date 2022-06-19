package flash

import (
	"html/template"
	"time"
)

type Flash struct {
	Type         string
	Header       template.HTML
	Body         template.HTML
	Dismissable  bool
	DismissAfter time.Duration
}
