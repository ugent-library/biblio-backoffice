package locale

import (
	"net/http"
)

// Detect is a middleware that sets the locale based on the Accept-Language header
// if not already set.
func Detect(l *Localizer) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := r.Context()
			if Get(c) == nil {
				loc := l.GetLocale(r.Header.Get("Accept-Language"))
				r = r.WithContext(Set(c, loc))
			}
			next.ServeHTTP(w, r)
		})
	}
}
