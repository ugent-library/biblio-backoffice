package views

import (
	"net/url"
	"testing"

	"github.com/a-h/templ"
	"github.com/stretchr/testify/require"
)

func TestPath(t *testing.T) {
	u := URL(parseURL("https://user:Pa$$w0rd@example.com:8081/test/path/?query=string#fragment"))

	u.Path("foo", "123", "bar/", "456", "/baz", "789")
	assertUrl(t, "https://user:Pa$$w0rd@example.com:8081/test/path/foo/123/bar/456/baz/789?query=string#fragment", u)
}

func TestAddQueryParam(t *testing.T) {
	u := URL(parseURL("https://user:Pa$$w0rd@example.com:8081/test/path/?query=string#fragment"))

	u.AddQueryParam("foo", "123")
	assertUrl(t, "https://user:Pa$$w0rd@example.com:8081/test/path/?foo=123&query=string#fragment", u) // query params are added in alphabetical order

	u.AddQueryParam("bar", "456")
	assertUrl(t, "https://user:Pa$$w0rd@example.com:8081/test/path/?bar=456&foo=123&query=string#fragment", u)

	u.AddQueryParam("baz", "789")
	assertUrl(t, "https://user:Pa$$w0rd@example.com:8081/test/path/?bar=456&baz=789&foo=123&query=string#fragment", u)

	u.AddQueryParam("bar", "123")
	assertUrl(t, "https://user:Pa$$w0rd@example.com:8081/test/path/?bar=456&bar=123&baz=789&foo=123&query=string#fragment", u)
}

func TestSetQueryParam(t *testing.T) {
	u := URL(parseURL("https://user:Pa$$w0rd@example.com:8081/test/path/?f=a&f=b"))

	u.SetQueryParam("f", "c")
	assertUrl(t, "https://user:Pa$$w0rd@example.com:8081/test/path/?f=c", u)

	u.SetQueryParam("key", "val")
	assertUrl(t, "https://user:Pa$$w0rd@example.com:8081/test/path/?f=c&key=val", u)
}

func TestURL(t *testing.T) {
	origURLStr := "https://user:Pa$$w0rd@example.com:8081/test/path/?f=a&f=b"
	origURL := parseURL(origURLStr)

	u := URL(origURL)
	u.SetQueryParam("f", "c")
	require.Equal(t, origURL.String(), origURLStr)
}

type queryType struct {
	Query   string `query:"q,omitempty"`
	OrderBy string `query:"order_by,omitempty"`
	Limit   int    `query:"limit,omitempty"`
}

func TestQueryWithStruct(t *testing.T) {
	u := URL(parseURL("https://user:Pa$$w0rd@example.com:8081/test/path/?query=string#fragment"))

	u.Query(queryType{
		Query:   "foo",
		OrderBy: "bar",
		Limit:   123,
	})
	assertUrl(t, "https://user:Pa$$w0rd@example.com:8081/test/path/?limit=123&order_by=bar&q=foo#fragment", u)

	u.Query(queryType{
		OrderBy: "baz",
	})
	assertUrl(t, "https://user:Pa$$w0rd@example.com:8081/test/path/?order_by=baz#fragment", u)

	u.Query(queryType{})
	assertUrl(t, "https://user:Pa$$w0rd@example.com:8081/test/path/#fragment", u)
}

func TestQueryWithString(t *testing.T) {
	u := URL(parseURL("https://user:Pa$$w0rd@example.com:8081/test/path/?query=string#fragment"))

	u.Query("foo=123&bar=456")
	assertUrl(t, "https://user:Pa$$w0rd@example.com:8081/test/path/?foo=123&bar=456#fragment", u)

	u.Query("foo=321&bar=654&baz=987")
	assertUrl(t, "https://user:Pa$$w0rd@example.com:8081/test/path/?foo=321&bar=654&baz=987#fragment", u)
}

func TestClearQuery(t *testing.T) {
	u := URL(parseURL("https://user:Pa$$w0rd@example.com:8081/test/path/?foo=123&query=string#fragment"))

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

func parseURL(u string) *url.URL {
	parsed, err := url.Parse(u)
	if err != nil {
		panic(err)
	}

	return parsed
}
