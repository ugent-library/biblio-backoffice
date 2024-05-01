package views

import (
	"net/url"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/require"
)

func TestWithPath(t *testing.T) {
	u := NewURLBuilder("https://user:Pa$$w0rd@example.com:8081/test/path/?query=string#fragment")

	u.WithPath("foo", "123", "bar/", "456", "/baz", "789")
	assertUrl(t, "https://user:Pa$$w0rd@example.com:8081/test/path/foo/123/bar/456/baz/789?query=string#fragment", u)
}

func TestAddQuery(t *testing.T) {
	u := NewURLBuilder("https://user:Pa$$w0rd@example.com:8081/test/path/?query=string#fragment")

	u.AddQuery("foo", "123")
	assertUrl(t, "https://user:Pa$$w0rd@example.com:8081/test/path/?foo=123&query=string#fragment", u) // query params are added in alphabetical order

	u.AddQuery("bar", "456")
	assertUrl(t, "https://user:Pa$$w0rd@example.com:8081/test/path/?bar=456&foo=123&query=string#fragment", u)

	u.AddQuery("baz", "789")
	assertUrl(t, "https://user:Pa$$w0rd@example.com:8081/test/path/?bar=456&baz=789&foo=123&query=string#fragment", u)

	u.AddQuery("bar", "123")
	assertUrl(t, "https://user:Pa$$w0rd@example.com:8081/test/path/?bar=456&bar=123&baz=789&foo=123&query=string#fragment", u)
}

type queryType struct {
	Query   string `query:"q,omitempty"`
	OrderBy string `query:"order_by,omitempty"`
	Limit   int    `query:"limit,omitempty"`
}

func TestWithQuery(t *testing.T) {
	u := NewURLBuilder("https://user:Pa$$w0rd@example.com:8081/test/path/?query=string#fragment")

	u.WithQuery(queryType{
		Query:   "foo",
		OrderBy: "bar",
		Limit:   123,
	})
	assertUrl(t, "https://user:Pa$$w0rd@example.com:8081/test/path/?limit=123&order_by=bar&q=foo#fragment", u)

	u.WithQuery(queryType{
		OrderBy: "baz",
	})
	assertUrl(t, "https://user:Pa$$w0rd@example.com:8081/test/path/?order_by=baz#fragment", u)

	u.WithQuery(queryType{})
	assertUrl(t, "https://user:Pa$$w0rd@example.com:8081/test/path/#fragment", u)
}

func TestClearQuery(t *testing.T) {
	u := NewURLBuilder("https://user:Pa$$w0rd@example.com:8081/test/path/?foo=123&query=string#fragment")

	u.ClearQuery()
	assertUrl(t, "https://user:Pa$$w0rd@example.com:8081/test/path/#fragment", u)
}

func assertUrl(t *testing.T, expected string, actual *URLBuilder) {
	// assert as string
	require.Equal(t, expected, actual.String())

	// assert as URL struct
	url, err := url.Parse(expected)
	if err != nil {
		t.Errorf("invalid expected url: %s", err)
	}
	require.Equal(t, url, actual.URL())

	// assert as SafeURL (implicit string)
	require.Equal(t, templ.SafeURL(expected), actual.SafeURL())
}
