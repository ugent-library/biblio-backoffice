package helpers

import (
	"fmt"
	"html/template"
	"net/url"
	"time"

	"github.com/rvflash/elapsed"
	"github.com/ugent-library/biblio-backoffice/identifiers"
	"github.com/ugent-library/biblio-backoffice/localize"
	"github.com/ugent-library/biblio-backoffice/models"
	"github.com/ugent-library/friendly"
)

func FuncMap() template.FuncMap {
	return template.FuncMap{
		"searchArgs":        models.NewSearchArgs,
		"timeElapsed":       elapsed.LocalTime,
		"formatRange":       FormatRange,
		"formatBool":        FormatBool,
		"formatBytes":       friendly.Bytes,
		"formatTime":        FormatTime,
		"languageName":      localize.LanguageName,
		"resolveIdentifier": identifiers.Resolve,
		"pathEscape":        url.PathEscape,
	}
}

func FormatRange(start, end string) string {
	var v string
	if len(start) > 0 && len(end) > 0 && start == end {
		v = start
	} else if len(start) > 0 && len(end) > 0 {
		v = fmt.Sprintf("%s - %s", start, end)
	} else if len(start) > 0 {
		v = fmt.Sprintf("%s -", start)
	} else if len(end) > 0 {
		v = fmt.Sprintf("- %s", end)
	}

	return v
}

func FormatBool(b bool, t, f string) string {
	if b {
		return t
	}
	return f
}

func FormatTime(t time.Time, loc *time.Location, fmt string) string {
	return t.In(loc).Format(fmt)
}
