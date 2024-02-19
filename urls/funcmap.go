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

var queryEncoder = form.NewEncoder()

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
		"renderMsg":  renderMsg,
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

func renderMsg(oldVal string) template.HTML {
	oldVal = html.EscapeString(oldVal)

	oldVal = strings.ReplaceAll(oldVal, "\r\n", "<br>")

	re, _ := xurls.StrictMatchingScheme("https")
	urlIndexPairs := re.FindAllStringIndex(oldVal, -1)

	newValBuilder := strings.Builder{}
	startPos := 0
	for _, pair := range urlIndexPairs {
		// NON URL
		prefix := oldVal[startPos:pair[0]]
		if len(prefix) > 0 {
			newValBuilder.WriteString(prefix)
		}

		//URL
		postfix := oldVal[pair[0]:pair[1]]
		newValBuilder.WriteString("<a href=\"")
		newValBuilder.WriteString(postfix)
		newValBuilder.WriteString("\" target=\"_blank\">")
		newValBuilder.WriteString(postfix)
		newValBuilder.WriteString("</a>")
		startPos = pair[1]
	}
	prefix := oldVal[startPos:]
	if len(prefix) > 0 {
		newValBuilder.WriteString(prefix)
	}

	return template.HTML(newValBuilder.String())
}
