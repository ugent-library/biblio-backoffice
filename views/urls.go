package views

import (
	"net/url"

	"github.com/ugent-library/bind"
)

func urlWithQuery(u *url.URL, v any) *url.URL {
	vals, err := bind.EncodeQuery(v)
	if err != nil {
		return u
	}
	newU := *u
	newU.RawQuery = vals.Encode()
	return &newU
}
