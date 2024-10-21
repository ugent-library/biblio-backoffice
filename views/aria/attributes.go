package aria

import (
	"strings"

	"github.com/a-h/templ"
	"github.com/ugent-library/biblio-backoffice/views/util"
)

func Attributes(helpText string, helpFieldID string) templ.Attributes {
	if helpText != "" {
		if strings.Contains(helpText, "<") {
			return templ.Attributes{
				"aria-description": util.StripHTML(helpText, false),
				"aria-details":     helpFieldID,
			}
		} else {
			return templ.Attributes{
				"aria-describedby": helpFieldID,
			}
		}
	} else {
		return templ.Attributes{}
	}
}
