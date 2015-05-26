package core

import (
	"net/http"

	"github.com/zenazn/goji/web"
)

// Makes sure templates are stored in the context
func (b *Blog) ApplyTemplates(c *web.C, h http.Handler) http.Handler {
	fn := func(w http.ResponseWriter, r *http.Request) {
		c.Env["Template"] = b.Template
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
