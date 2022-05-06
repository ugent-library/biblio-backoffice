package forms

import (
	"net/url"

	"github.com/go-playground/form/v4"
)

var formDecoder = form.NewDecoder()
var formEncoder = form.NewEncoder()

func Decode(v interface{}, vals url.Values) error {
	return formDecoder.Decode(v, vals)
}

func Encode(v interface{}) (url.Values, error) {
	return formEncoder.Encode(v)
}
