package locale

import (
	"net/http"
)

// DEPRECATED
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
