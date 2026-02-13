package middleware

import (
	"context"
	"net/http"
	"strings"
	"sync"
)

var (
	additionalHeadersToContext = []string{}
	headersMu                  sync.RWMutex
)

// AddHeadersToContext allows adding custom headers to be injected into the context
func AddHeadersToContext(headers []string) {
	headersMu.Lock()
	defer headersMu.Unlock()
	additionalHeadersToContext = append(additionalHeadersToContext, headers...)
}

func InjectHeadersToContext(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		listCommonHeaders := map[string]bool{
			"Authorization": true,
			"Host":          true,
		}

		listPrefixHeaders := []string{
			"x-",
		}

		// Get additional headers with read lock
		headersMu.RLock()
		additionalHeaders := make([]string, len(additionalHeadersToContext))
		copy(additionalHeaders, additionalHeadersToContext)
		headersMu.RUnlock()

		// Add additional headers to common headers map
		for _, header := range additionalHeaders {
			listCommonHeaders[header] = true
		}

		ctx := r.Context()
		for headerKey, headerValues := range r.Header {
			// Combine multiple header values with comma separator
			headerValue := strings.Join(headerValues, ",")

			// Check if header is in common headers
			if _, exists := listCommonHeaders[headerKey]; exists {
				ctx = context.WithValue(ctx, strings.ToLower(headerKey), headerValue)
				continue
			}

			// Check if header starts with any of the prefixes
			headerKey = strings.ToLower(headerKey)
			for _, prefix := range listPrefixHeaders {
				if strings.HasPrefix(headerKey, strings.ToLower(prefix)) {
					ctx = context.WithValue(ctx, strings.ToLower(headerKey), headerValue)
					break
				}
			}
		}

		r = r.WithContext(ctx)
		next.ServeHTTP(w, r)
	})
}
