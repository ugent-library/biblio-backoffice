package views

import (
	"net/url"

	"github.com/a-h/templ"
	"github.com/ugent-library/bind"
)

type URLBuilder struct {
	url *url.URL
}

func URL(base *url.URL) *URLBuilder {
	return &URLBuilder{base}
}

func (builder *URLBuilder) WithPath(path ...string) *URLBuilder {
	builder.url = builder.url.JoinPath(path...)

	return builder
}

func (builder *URLBuilder) AddQuery(key string, value string) *URLBuilder {
	query := builder.url.Query()
	query.Add(key, value)

	builder.url.RawQuery = query.Encode()

	return builder
}

func (builder *URLBuilder) WithQuery(query any) *URLBuilder {
	vals, err := bind.EncodeQuery(query)
	if err != nil {
		return builder
	}

	builder.url.RawQuery = vals.Encode()
	return builder
}

func (builder *URLBuilder) ClearQuery() *URLBuilder {
	builder.url.RawQuery = ""

	return builder
}

func (builder *URLBuilder) URL() *url.URL {
	return builder.url
}

func (builder *URLBuilder) String() string {
	return builder.url.String()
}

func (builder *URLBuilder) SafeURL() templ.SafeURL {
	return templ.URL(builder.String())
}
