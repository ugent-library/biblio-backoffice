package locale

import (
	"context"
)

type contextKey int

const localeKey = contextKey(0)

// DEPRECATED
func Get(c context.Context) *Locale {
	if v := c.Value(localeKey); v != nil {
		return v.(*Locale)
	}
	return nil
}

// DEPRECATED
func Set(c context.Context, l *Locale) context.Context {
	return context.WithValue(c, localeKey, l)
}
