package urlviews

import (
	"net/url"

	"github.com/ugent-library/bind"
)

func Query(u *url.URL, v any) *url.URL {
	vals, err := bind.EncodeQuery(v)
	if err != nil {
		return u
	}
	newU := *u
	newU.RawQuery = vals.Encode()
	return &newU
}

func QueryClear(u *url.URL) *url.URL {
	newU := *u
	newU.RawQuery = ""
	return &newU
}

func QuerySet(u *url.URL, k string, v string) *url.URL {
	newU := *u
	q := u.Query()
	q.Set(k, v)
	newU.RawQuery = q.Encode()
	return &newU
}
