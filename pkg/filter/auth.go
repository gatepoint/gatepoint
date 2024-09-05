package filter

import (
	"net/http"

	"github.com/gatepoint/gatepoint/pkg/config"
)

func AuthFilter(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if config.EnableDebug() {
			handler.ServeHTTP(w, r)
			return
		}

		// do auth check
	})
}
