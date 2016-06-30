// Package acl provides http.Handlers to prevent access to pages for
// authenticated users and for non-authenticated users.
package acl

import (
	"net/http"

	"github.com/blue-jay/blueprint/lib/session"
)

// DisallowAuth does not allow authenticated users to access the page.
func DisallowAuth(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess := session.Instance(r)

		// If user is authenticated, don't allow them to access the page
		if sess.Values["id"] != nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		h.ServeHTTP(w, r)
	})
}

// DisallowAnon does not allow anonymous users to access the page.
func DisallowAnon(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sess := session.Instance(r)

		// If user is not authenticated, don't allow them to access the page
		if sess.Values["id"] == nil {
			http.Redirect(w, r, "/", http.StatusFound)
			return
		}

		h.ServeHTTP(w, r)
	})
}
