package middleware

import (
	"bytes"
	"io"
	"net/http"
	"sync"
	"time"

	"github.com/adityayuga/go-sdk/log"
)

var (
	mu              sync.RWMutex
	logRequestBody  = false
	logResponseBody = false
)

// EnableRequestBodyLogging enables request body logging in the Logger middleware
func EnableRequestBodyLogging() {
	mu.Lock()
	defer mu.Unlock()
	logRequestBody = true
}

// DisableRequestBodyLogging disables request body logging in the Logger middleware
func DisableRequestBodyLogging() {
	mu.Lock()
	defer mu.Unlock()
	logRequestBody = false
}

// EnableResponseBodyLogging enables response body logging in the Logger middleware
func EnableResponseBodyLogging() {
	mu.Lock()
	defer mu.Unlock()
	logResponseBody = true
}

// DisableResponseBodyLogging disables response body logging in the Logger middleware
func DisableResponseBodyLogging() {
	mu.Lock()
	defer mu.Unlock()
	logResponseBody = false
}

// Logger is a middleware that logs HTTP requests (compatible with both standard and chi routers)
func Logger(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		ctx := r.Context()

		// Get current configuration
		mu.RLock()
		logReqBody := logRequestBody
		logRespBody := logResponseBody
		mu.RUnlock()

		// Capture request body if logging is enabled
		var requestBody []byte
		if logReqBody && r.Body != nil {
			var err error
			requestBody, err = io.ReadAll(r.Body)
			if err != nil {
				log.Errorf(ctx, "failed to read request body: %v", err)
				requestBody = []byte{}
			}
			// Restore the body for the next handler
			r.Body = io.NopCloser(bytes.NewBuffer(requestBody))
		}

		// Create a response writer wrapper to capture status code and body
		wrapped := &responseWriter{
			ResponseWriter: w,
			statusCode:     http.StatusOK,
			logBody:        logRespBody,
		}

		// Serve the request
		next.ServeHTTP(wrapped, r)

		// Build log fields
		fields := []interface{}{
			"method", r.Method,
			"uri", r.RequestURI,
			"proto", r.Proto,
			"status", wrapped.statusCode,
			"remote_addr", r.RemoteAddr,
			"duration_ms", time.Since(start).Milliseconds(),
		}

		// Add request body if logging is enabled
		if logReqBody && len(requestBody) > 0 {
			fields = append(fields, "request_body", string(requestBody))
		}

		// Add response body if logging is enabled
		if logRespBody && len(wrapped.body) > 0 {
			fields = append(fields, "response_body", string(wrapped.body))
		}

		log.Info(ctx, "HTTP Request", fields...)
	})
}

type responseWriter struct {
	http.ResponseWriter
	statusCode int
	body       []byte
	logBody    bool
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.logBody {
		rw.body = append(rw.body, b...)
	}
	return rw.ResponseWriter.Write(b)
}
