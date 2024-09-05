package filter

import (
	"net/http"

	"github.com/gatepoint/gatepoint/pkg/log"
)

func RecordFilter(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Info("record filter")
		handler.ServeHTTP(w, r)
	})
}
