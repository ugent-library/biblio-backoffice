package controllers

import (
	"net/url"

	"github.com/go-playground/form/v4"
)

var (
	formDecoder  = form.NewDecoder()
	queryDecoder = form.NewDecoder()
)

func init() {
	formDecoder.SetTagName("form")
	formDecoder.SetMode(form.ModeExplicit)
	queryDecoder.SetTagName("query")
	queryDecoder.SetMode(form.ModeExplicit)
}

func DecodeQuery(v interface{}, vals url.Values) error {
	return queryDecoder.Decode(v, vals)
}

func DecodeForm(v interface{}, vals url.Values) error {
	return formDecoder.Decode(v, vals)
}
