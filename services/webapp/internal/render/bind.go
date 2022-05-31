package render

import (
	"net/url"

	"github.com/gorilla/schema"
)

var decoder = schema.NewDecoder()

func init() {
	decoder.SetAliasTag("form")
	decoder.IgnoreUnknownKeys(true)
}

func Bind(v interface{}, vals url.Values) error {
	return decoder.Decode(v, vals)
}
