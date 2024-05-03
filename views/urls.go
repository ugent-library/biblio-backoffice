package views

import (
	"net/url"

	"github.com/a-h/templ"
	"github.com/samber/lo"
	"github.com/ugent-library/bind"
)

type URLBuilder struct {
	url *url.URL
}

func URL(base *url.URL) *URLBuilder {
	return &URLBuilder{base}
}

func URLFromString(base string) *URLBuilder {
	return &URLBuilder{
		url: lo.Must(url.Parse(base)),
	}
}

func (builder *URLBuilder) Path(path ...string) *URLBuilder {
	builder.url = builder.url.JoinPath(path...)

	return builder
}

func (builder *URLBuilder) AddQueryParam(key string, value string) *URLBuilder {
	query := builder.url.Query()
	query.Add(key, value)

	builder.url.RawQuery = query.Encode()

	return builder
}

func (builder *URLBuilder) Query(query interface{}) *URLBuilder {
	if str, ok := query.(string); ok {
		builder.url.RawQuery = str
	} else {
		vals, err := bind.EncodeQuery(query)
		if err != nil {
			return builder
		}

		builder.url.RawQuery = vals.Encode()
	}

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
