package middleware

import (
	"net/http"

	"github.com/adityayuga/go-sdk/log"
)

// Recover is a middleware that recovers from panics (compatible with both standard and chi routers)
func Recover(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				ctx := r.Context()
				log.Errorf(ctx, "panic recovered, path: %s, error: %+v", r.URL.Path, err)
				http.Error(w, "Internal Server Error", http.StatusInternalServerError)
			}
		}()

		next.ServeHTTP(w, r)
	})
}
