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
	copyBase := *base
	return &URLBuilder{url: &copyBase}
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

func (builder *URLBuilder) QueryAdd(pairs ...string) *URLBuilder {
	if len(pairs)%2 != 0 {
		panic("QueryAdd arguments should be an even sized-list of key value pairs")
	}

	query := builder.url.Query()

	for i := 0; i < len(pairs); i += 2 {
		query.Add(pairs[i], pairs[i+1])
	}

	builder.url.RawQuery = query.Encode()

	return builder
}

func (builder *URLBuilder) SetQueryParam(key string, value string) *URLBuilder {
	query := builder.url.Query()
	query.Set(key, value)

	builder.url.RawQuery = query.Encode()

	return builder
}

func (builder *URLBuilder) QuerySet(pairs ...string) *URLBuilder {
	if len(pairs)%2 != 0 {
		panic("QueryAdd arguments should be an even-sized list of key value pairs")
	}

	query := builder.url.Query()

	for i := 0; i < len(pairs); i += 2 {
		query.Set(pairs[i], pairs[i+1])
	}

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
