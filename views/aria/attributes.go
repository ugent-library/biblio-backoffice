package aria

import (
	"log"
	"strings"

	"github.com/a-h/templ"
	"jaytaylor.com/html2text"
)

func Attributes(helpText string, helpFieldID string) templ.Attributes {
	if helpText != "" {
		if strings.Contains(helpText, "<") {
			return templ.Attributes{
				"aria-description": stripHTML(helpText),
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

func stripHTML(text string) string {
	result, err := html2text.FromString(text, html2text.Options{})
	if err != nil {
		log.Printf("Error while stripping HTML from '%s': %s\n", text, err.Error())

		return text
	}

	return result
}
