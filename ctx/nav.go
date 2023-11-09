package ctx

import (
	"net/http"
)

func SetNav(nav string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := Get(r)
			c.Nav = nav
			next.ServeHTTP(w, r)
		})
	}
}
