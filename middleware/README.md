# HTTP Logger Middleware

An HTTP request/response logging middleware that integrates with the Go SDK logger.

## Features

- **Automatic Logging** - Logs all HTTP requests with method, URI, status code, duration, and remote address
- **Optional Request Body Logging** - Capture and log request payloads
- **Optional Response Body Logging** - Capture and log response payloads
- **Debug ID Integration** - Automatically includes debug ID from context in all logs
- **Multiple Handler Support** - Works with standard `http.Handler` and any compatible router (Chi, Gorilla Mux, etc.)
- **JSON Formatted Logs** - All logs are in JSON format for easy parsing and analysis
- **Setter Functions** - Simple API to enable/disable body logging globally

## Usage

### Basic Usage (No Body Logging)

```go
import "github.com/adityayuga/go-sdk/middleware"

// Use logger middleware with default settings (no body logging)
http.Handle("/api/", middleware.Logger(http.HandlerFunc(handler)))

// Or with a router
router := chi.NewRouter()
router.Use(middleware.Logger)
```

### Enable Request Body Logging

```go
import "github.com/adityayuga/go-sdk/middleware"

// Enable request body logging globally
middleware.EnableRequestBodyLogging()

// All Logger middleware after this will log request bodies
router.Use(middleware.Logger)
```

### Enable Response Body Logging

```go
import "github.com/adityayuga/go-sdk/middleware"

// Enable response body logging globally
middleware.EnableResponseBodyLogging()

// All Logger middleware after this will log response bodies
router.Use(middleware.Logger)
```

### Enable Both Request and Response Body Logging

```go
import "github.com/adityayuga/go-sdk/middleware"

// Enable both globally
middleware.EnableRequestBodyLogging()
middleware.EnableResponseBodyLogging()

// All Logger middleware after this will log both request and response bodies
router.Use(middleware.Logger)
```

### Disable Body Logging

```go
import "github.com/adityayuga/go-sdk/middleware"

// Disable request body logging
middleware.DisableRequestBodyLogging()

// Disable response body logging
middleware.DisableResponseBodyLogging()
```

## Log Output Example

### Basic Log (No Body Logging)
```json
{
  "level": "info",
  "ts": "2026-02-14T00:44:50.080+0700",
  "msg": "HTTP Request",
  "method": "GET",
  "uri": "/api/users/123",
  "proto": "HTTP/1.1",
  "status": 200,
  "remote_addr": "192.0.2.1:1234",
  "duration_ms": 45,
  "debug_id": "6323890720260214004450"
}
```

### With Request Body
```json
{
  "level": "info",
  "ts": "2026-02-14T00:44:50.080+0700",
  "msg": "HTTP Request",
  "method": "POST",
  "uri": "/api/users",
  "status": 201,
  "duration_ms": 12,
  "request_body": "{\"name\":\"John\",\"age\":30}",
  "debug_id": "0795567320260214004450"
}
```

### With Response Body
```json
{
  "level": "info",
  "ts": "2026-02-14T00:44:50.080+0700",
  "msg": "HTTP Request",
  "method": "GET",
  "uri": "/api/users/123",
  "status": 200,
  "duration_ms": 8,
  "response_body": "{\"id\":123,\"name\":\"John\",\"email\":\"john@example.com\"}",
  "debug_id": "1975032020260214004450"
}
```

### With Both Request and Response Bodies
```json
{
  "level": "info",
  "ts": "2026-02-14T00:44:50.080+0700",
  "msg": "HTTP Request",
  "method": "POST",
  "uri": "/api/users",
  "status": 201,
  "duration_ms": 15,
  "request_body": "{\"name\":\"Jane\",\"email\":\"jane@example.com\"}",
  "response_body": "{\"id\":456,\"name\":\"Jane\",\"email\":\"jane@example.com\"}",
  "debug_id": "5862261720260214004450"
}
```

## API Reference

### `Logger(next http.Handler) http.Handler`
Middleware that logs HTTP requests. Uses global configuration set by enable/disable functions.

### `EnableRequestBodyLogging()`
Enables request body logging globally for all Logger middleware instances.

### `DisableRequestBodyLogging()`
Disables request body logging globally for all Logger middleware instances.

### `EnableResponseBodyLogging()`
Enables response body logging globally for all Logger middleware instances.

### `DisableResponseBodyLogging()`
Disables response body logging globally for all Logger middleware instances.

---

# Inject Headers to Context Middleware

Middleware that extracts HTTP headers and injects them into the request context for use by downstream handlers.

## Features

- **Automatic Header Injection** - Extracts headers and adds them to context
- **Common Headers** - Pre-configured for Authorization and Host headers
- **Custom Headers** - Dynamically add custom headers to be injected
- **Prefix-Based Headers** - All headers starting with `x-` are automatically injected
- **Thread-Safe** - Safe to register custom headers from multiple goroutines

## Usage

### Basic Usage (Common Headers Only)

```go
import "github.com/adityayuga/go-sdk/middleware"

// Use middleware - automatically injects Authorization, Host, and x-* headers
router := chi.NewRouter()
router.Use(middleware.InjectHeadersToContext)
```

### Add Custom Headers

```go
import "github.com/adityayuga/go-sdk/middleware"

// Register custom headers to be injected into context during middleware setup
middleware.AddHeadersToContext([]string{"X-API-Key", "X-Request-ID", "X-Client-Version"})

router := chi.NewRouter()
router.Use(middleware.InjectHeadersToContext)
```

### Access Headers from Context

```go
func handleRequest(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    // Get Authorization header
    auth := ctx.Value("authorization").(string)
    
    // Get custom header
    apiKey := ctx.Value("x-api-key").(string)
    
    // Get any x-* header
    requestID := ctx.Value("x-request-id").(string)
}
```

## API Reference

### `InjectHeadersToContext(next http.Handler) http.Handler`
Middleware that extracts HTTP headers and injects them into the request context. Headers are stored with lowercase keys.

**Automatically Injects:**
- `authorization` - Authorization header
- `host` - Host header
- All `x-*` headers (case-insensitive, stored as lowercase)

**Custom Headers:** Added via `AddHeadersToContext()`

### `AddHeadersToContext(headers []string)`
Registers custom headers to be injected into context by the middleware.

**Parameters:**
- `headers` - Array of header names to register (e.g., `[]string{"X-API-Key", "X-Custom-Header"}`)

**Thread-Safety:** Safe to call from multiple goroutines

**Example:**
```go
middleware.AddHeadersToContext([]string{"X-API-Key", "X-Tenant-ID"})
```

## Important Notes

1. **Header Key Format** - Headers are stored in context with lowercase keys for consistent access
2. **Multiple Values** - If a header appears multiple times, values are joined with commas
3. **Thread-Safe Configuration** - Use mutex protection for concurrent-safe header registration
4. **Prefix Matching** - All headers starting with `x-` are automatically included
5. **Custom Headers Override** - Custom headers take precedence in the common headers map

## Example Setup

```go
import (
    "github.com/go-chi/chi/v5"
    "github.com/adityayuga/go-sdk/middleware"
)

// Setup during application initialization
middleware.AddHeadersToContext([]string{
    "X-API-Key",
    "X-Client-ID",
    "X-Request-Source",
})

router := chi.NewRouter()
router.Use(middleware.InjectHeadersToContext)
router.Use(middleware.Logger)

router.Get("/api/data", func(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    apiKey := ctx.Value("x-api-key")
    clientID := ctx.Value("x-client-id")
    // Use extracted headers...
})
```

---

# Panic Recovery Middleware

Middleware that catches and recovers from panics in HTTP handlers, preventing server crashes.

## Features

- **Panic Recovery** - Gracefully handles panics in request handlers
- **Automatic Logging** - Logs panic errors with request path and error details
- **Error Response** - Returns 500 Internal Server Error when panic occurs
- **Context-Aware** - Uses context for debug ID and error logging
- **Standard HTTP Handler** - Works with any standard `http.Handler` compatible router

## Usage

### Basic Usage

```go
import "github.com/adityayuga/go-sdk/middleware"

router := chi.NewRouter()
router.Use(middleware.Recover)
```

### With Logger Middleware

```go
import (
    "github.com/go-chi/chi/v5"
    "github.com/adityayuga/go-sdk/middleware"
    "github.com/adityayuga/go-sdk/log"
)

// Initialize logger
log.Init("info")

router := chi.NewRouter()
// Recovery middleware should typically be early in the middleware chain
router.Use(middleware.Recover)
router.Use(middleware.Logger)
```

## API Reference

### `Recover(next http.Handler) http.Handler`
Middleware that recovers from panics in HTTP handlers and logs the error.

**Behavior:**
- Catches any panic during request handling
- Logs error with context debug ID
- Returns HTTP 500 (Internal Server Error) to client
- Prevents application crash

**Log Output:**
```json
{
  "level": "error",
  "ts": "2026-02-14T00:44:50.080+0700",
  "msg": "panic recovered",
  "path": "/api/users",
  "error": "runtime error: index out of range",
  "debug_id": "6323890720260214004450"
}
```

## Important Notes

1. **Middleware Order** - Place `Recover` early in the middleware chain to catch panics from all downstream handlers
2. **Error Response** - When a panic is caught, a generic 500 error is returned to the client
3. **Logging** - Panics are logged using the logger package with full error context
4. **Debug ID** - The recovery middleware includes the debug ID from context in error logs
5. **Not a Replacement** - Should be used alongside proper error handling, not as a substitute

## Example: Comprehensive Middleware Setup

```go
import (
    "github.com/go-chi/chi/v5"
    "github.com/adityayuga/go-sdk/middleware"
    "github.com/adityayuga/go-sdk/log"
)

// Initialize logger
log.Init("info")

// Register custom headers
middleware.AddHeadersToContext([]string{"X-API-Key", "X-User-ID"})

router := chi.NewRouter()

// Middleware order matters - Recovery first to catch all panics
router.Use(middleware.Recover)

// Header injection
router.Use(middleware.InjectHeadersToContext)

// Enable body logging for debugging (optional)
middleware.EnableRequestBodyLogging()
middleware.EnableResponseBodyLogging()
router.Use(middleware.Logger)

// Define routes
router.Post("/api/users", createUser)
router.Get("/api/users/{id}", getUser)
```

## Important Notes

1. **Global Configuration**: The enable/disable functions set global configuration. All Logger middleware instances share the same settings.

2. **Thread-Safe**: All configuration changes are protected with mutex locks for concurrent safety.

3. **Performance**: Logging request/response bodies incurs a performance cost. Use only when necessary for debugging.

4. **Memory Usage**: Large request/response bodies are buffered in memory. Avoid enabling for endpoints that handle large file uploads/downloads.

5. **Security**: Be careful with sensitive data in body logging (passwords, tokens, credit cards). Consider filtering sensitive fields in production.

6. **Request Body Restoration**: The middleware automatically restores the request body after reading it, so downstream handlers can still access it.

7. **Debug ID**: Every request automatically receives a unique debug ID for tracing across the application.

## Example with Chi Router

```go
import (
    "github.com/go-chi/chi/v5"
    "github.com/adityayuga/go-sdk/middleware"
)

router := chi.NewRouter()

// Enable body logging if needed
// middleware.EnableRequestBodyLogging()
// middleware.EnableResponseBodyLogging()

// Use logger middleware globally
router.Use(middleware.Logger)

router.Post("/api/users", createUser)
router.Get("/api/users/{id}", getUser)
```

## Example: Enable Body Logging for Debugging

```go
import (
    "github.com/go-chi/chi/v5"
    "github.com/adityayuga/go-sdk/middleware"
)

router := chi.NewRouter()

// Enable both request and response body logging
middleware.EnableRequestBodyLogging()
middleware.EnableResponseBodyLogging()

// Use logger middleware - will now log request and response bodies
router.Use(middleware.Logger)
```

## Example: Conditional Body Logging

```go
import (
    "os"
    "github.com/go-chi/chi/v5"
    "github.com/adityayuga/go-sdk/middleware"
)

router := chi.NewRouter()

// Enable body logging only in development
if os.Getenv("ENV") == "development" {
    middleware.EnableRequestBodyLogging()
    middleware.EnableResponseBodyLogging()
}

router.Use(middleware.Logger)
```

## Testing

Run tests with:

```bash
go test -v ./middleware
```

The middleware includes comprehensive tests for:
- Basic logging without body content
- Request body logging
- Response body logging  
- Both request and response body logging
- Error responses (5xx)
- Multiple response writes (chunks)
