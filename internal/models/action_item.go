package models

import "net/url"

type ActionItem struct {
	Template string
	URL      *url.URL
	Label    string
}
