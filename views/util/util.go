package util

import (
	"html"
	"log"
	"strings"

	"jaytaylor.com/html2text"
	"mvdan.cc/xurls/v2"
)

var reURL, _ = xurls.StrictMatchingScheme("https")

func Linkify(text string) string {
	text = html.EscapeString(text)

	matches := reURL.FindAllStringIndex(text, -1)

	b := strings.Builder{}
	pos := 0
	for _, match := range matches {
		before := text[pos:match[0]]
		if len(before) > 0 {
			b.WriteString(before)
		}

		link := text[match[0]:match[1]]
		b.WriteString(`<a href="`)
		b.WriteString(link)
		b.WriteString(`" target="_blank">`)
		b.WriteString(link)
		b.WriteString(`</a>`)
		pos = match[1]
	}

	after := text[pos:]
	if len(after) > 0 {
		b.WriteString(after)
	}

	return b.String()
}

func StripHTML(text string, textOnly bool) string {
	result, err := html2text.FromString(text, html2text.Options{
		TextOnly: textOnly,
	})
	if err != nil {
		log.Printf("Error while stripping HTML from '%s': %s\n", text, err.Error())

		return text
	}

	return result
}
