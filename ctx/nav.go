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

func SetSubNav(subNav string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := Get(r)
			c.SubNav = subNav
			next.ServeHTTP(w, r)
		})
	}
}

func SetBreadcrumbs(breadcrumbs ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := Get(r)
			c.Breadcrumbs = breadcrumbs
			next.ServeHTTP(w, r)
		})
	}
}

func AddBreadcrumb(breadcrumb string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			c := Get(r)
			c.Breadcrumbs = append(c.Breadcrumbs, breadcrumb)
			next.ServeHTTP(w, r)
		})
	}
}
