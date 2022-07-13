package bind

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/go-playground/form/v4"
)

type Flag int

const (
	Vacuum Flag = iota
)

var (
	PathValuesFunc func(r *http.Request) url.Values

	pathDecoder  = form.NewDecoder()
	formDecoder  = form.NewDecoder()
	queryDecoder = form.NewDecoder()
	pathEncoder  = form.NewEncoder()
	formEncoder  = form.NewEncoder()
	queryEncoder = form.NewEncoder()
)

func init() {
	pathDecoder.SetTagName("path")
	pathDecoder.SetMode(form.ModeExplicit)
	formDecoder.SetTagName("form")
	formDecoder.SetMode(form.ModeExplicit)
	queryDecoder.SetTagName("query")
	queryDecoder.SetMode(form.ModeExplicit)
	pathEncoder.SetTagName("path")
	pathEncoder.SetMode(form.ModeExplicit)
	formEncoder.SetTagName("form")
	formEncoder.SetMode(form.ModeExplicit)
	queryEncoder.SetTagName("query")
	queryEncoder.SetMode(form.ModeExplicit)
}

func PathValues(r *http.Request) url.Values {
	if PathValuesFunc != nil {
		return PathValuesFunc(r)
	}
	return nil
}

func VacuumValues(vals url.Values) url.Values {
	for k, v := range vals {
		vacuumed := vacuumStringSlice(v)
		if len(vacuumed) > 0 {
			vals[k] = vacuumed
		} else {
			delete(vals, k)
		}
	}
	return vals
}

func RequestPath(r *http.Request, v any, flags ...Flag) error {
	return Path(PathValues(r), v, flags...)
}

func Path(vals url.Values, v any, flags ...Flag) error {
	if hasFlag(flags, Vacuum) {
		vals = VacuumValues(vals)
	}
	return pathDecoder.Decode(v, vals)
}

func RequestQuery(r *http.Request, v any, flags ...Flag) error {
	return Query(r.URL.Query(), v, flags...)
}

func Query(vals url.Values, v any, flags ...Flag) error {
	if hasFlag(flags, Vacuum) {
		vals = VacuumValues(vals)
	}
	return queryDecoder.Decode(v, vals)
}

func RequestForm(r *http.Request, v any, flags ...Flag) error {
	r.ParseForm()
	return Form(r.Form, v, flags...)
}

func Form(vals url.Values, v any, flags ...Flag) error {
	if hasFlag(flags, Vacuum) {
		vals = VacuumValues(vals)
	}
	return formDecoder.Decode(v, vals)
}

func Request(r *http.Request, v any, flags ...Flag) error {
	if err := RequestPath(r, v, flags...); err != nil {
		return err
	}

	if err := Query(r.URL.Query(), v, flags...); err != nil {
		return err
	}

	r.ParseForm()

	return Form(r.Form, v, flags...)
}

// include encoding helpers as a convenience

func EncodePath(v any) (url.Values, error) {
	return pathEncoder.Encode(v)
}

func EncodeForm(v any) (url.Values, error) {
	return formEncoder.Encode(v)
}

func EncodeQuery(v any) (url.Values, error) {
	return queryEncoder.Encode(v)
}

// helpers

func vacuumStringSlice(vals []string) []string {
	var tmp []string
	for _, v := range vals {
		v = strings.TrimSpace(v)
		if v != "" {
			tmp = append(tmp, v)
		}
	}
	return tmp
}

func hasFlag(flags []Flag, flag Flag) bool {
	for _, f := range flags {
		if f == flag {
			return true
		}
	}
	return false
}
