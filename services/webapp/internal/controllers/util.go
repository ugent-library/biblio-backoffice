package controllers

import (
	"net/url"

	"github.com/go-playground/form/v4"
)

var (
	formDecoder = form.NewDecoder()
	formEncoder = form.NewEncoder()
)

func init() {
	formDecoder.SetTagName("form")
	formEncoder.SetTagName("form")
}

func DecodeForm(v interface{}, vals url.Values) error {
	return formDecoder.Decode(v, vals)
}

func EncodeForm(v interface{}) (url.Values, error) {
	return formEncoder.Encode(v)
}
