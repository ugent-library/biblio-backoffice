package urls

import (
	"html"
	"html/template"
	"net/url"
	"strings"

	"github.com/go-playground/form/v4"
	"github.com/nics/ich"
	"mvdan.cc/xurls/v2"
)

var (
	queryEncoder = form.NewEncoder()
	reURL, _     = xurls.StrictMatchingScheme("https")
)

func init() {
	queryEncoder.SetTagName("query")
	queryEncoder.SetMode(form.ModeExplicit)
}

// TODO split into mux and query packages
func FuncMap(r *ich.Mux, scheme, host string) template.FuncMap {
	return template.FuncMap{
		"urlFor":     urlFor(r, scheme, host),
		"pathFor":    pathFor(r),
		"query":      query,
		"querySet":   querySet,
		"queryAdd":   queryAdd,
		"queryDel":   queryDel,
		"queryClear": queryClear,
		"linkify":    linkify,
	}
}

func urlFor(r *ich.Mux, scheme, host string) func(string, ...string) *url.URL {
	return func(name string, pairs ...string) *url.URL {
		u := r.PathTo(name, pairs...)
		u.Host = host
		u.Scheme = scheme
		return u
	}
}

func pathFor(r *ich.Mux) func(string, ...string) *url.URL {
	return r.PathTo
}

func query(v any, u *url.URL) (*url.URL, error) {
	vals, err := queryEncoder.Encode(v)
	if err != nil {
		return u, err
	}

	newU := *u
	newU.RawQuery = vals.Encode()

	return &newU, nil
}

func querySet(k, v string, u *url.URL) (*url.URL, error) {
	newU := *u
	q := u.Query()
	q.Set(k, v)
	newU.RawQuery = q.Encode()

	return &newU, nil
}

func queryAdd(k, v string, u *url.URL) (*url.URL, error) {
	newU := *u
	q := u.Query()
	q.Add(k, v)
	newU.RawQuery = q.Encode()

	return &newU, nil
}

func queryDel(k string, u *url.URL) (*url.URL, error) {
	newU := *u
	q := u.Query()
	q.Del(k)
	newU.RawQuery = q.Encode()

	return &newU, nil
}

func queryClear(u *url.URL) (*url.URL, error) {
	newU := *u
	newU.RawQuery = ""

	return &newU, nil
}

func linkify(text string) template.HTML {
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

	return template.HTML(b.String())
}
