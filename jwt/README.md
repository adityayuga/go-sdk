# JWT Package

A Go JWT (JSON Web Token) package providing secure token generation and validation with support for both access and refresh tokens.

## Features

- **Dual Token Support** - Generate and validate both short-lived access tokens and long-lived refresh tokens
- **HS512 Signing** - Uses HMAC SHA-512 for token signing, requiring 64-character secrets
- **Context-Aware** - All operations take context as parameter for debug ID and tracing
- **Type-Safe Claims** - Strongly-typed Data struct with UserID, Email, Phone, and Status fields
- **Configurable Expiry** - Customizable expiry times for both access and refresh tokens
- **Debug Integration** - Automatic debug ID logging for error tracking
- **Error Handling** - Comprehensive validation of token format and secret length

## Installation

```bash
go get github.com/adityayuga/go-sdk/jwt
```

## Configuration

### Required Setup

```go
import "github.com/adityayuga/go-sdk/jwt"

config := jwt.Config{
    AccessTokenSecret:              "your-64-character-access-token-secret-key-here!!",
    RefreshTokenSecret:             "your-64-character-refresh-token-secret-key-here!",
    AccessTokenExpiryTimeInMinutes: 60, // optional, defaults to 60
    RefreshTokenExpireTimeInDay:     30, // optional, defaults to 30
}

jwtModule, err := jwt.New(context.Background(), config)
if err != nil {
    panic(err)
}
```

## Secret Key Generation

JWT secrets must be exactly 64 characters for HS512 signing. Generate secure secrets using:

```bash
# Generate a random 64-character secret
openssl rand -base64 64 | head -c 64

# Example output:
# W8x9Q2pL4mN6tK1vJ3hYbZ5cD7fGaR9sE2wU4oI8jX0P5yT6nM
```

Or use Go to generate:

```go
import "crypto/rand"
import "encoding/base64"

func generateSecret() string {
    bytes := make([]byte, 48)
    rand.Read(bytes)
    return base64.URLEncoding.EncodeToString(bytes)
}
```

## Usage

### Initialize JWT Module

```go
import (
    "context"
    "github.com/adityayuga/go-sdk/jwt"
)

ctx := context.Background()

config := jwt.Config{
    AccessTokenSecret:              "your-64-character-access-token-secret",
    RefreshTokenSecret:             "your-64-character-refresh-token-secret",
    AccessTokenExpiryTimeInMinutes: 60,
    RefreshTokenExpireTimeInDay:     30,
}

jwtModule, err := jwt.New(ctx, config)
if err != nil {
    log.Error(ctx, "Failed to initialize JWT module: %v", err)
    return
}
```

### Generate Access Token

Short-lived token for API authentication (default: 60 minutes).

```go
data := jwt.Data{
    UserID: 123,
    Email:  "user@example.com",
    Phone:  "1234567890",
    Status: 1,
}

token, expiryTime, err := jwtModule.GenerateAccessToken(ctx, data)
if err != nil {
    log.Error(ctx, "Failed to generate access token: %v", err)
    return
}

log.Info(ctx, "Access token generated", "token", token, "expires_at", expiryTime)
```

### Generate Refresh Token

Long-lived token for obtaining new access tokens (default: 30 days).

```go
data := jwt.Data{
    UserID: 123,
    Email:  "user@example.com",
    Phone:  "1234567890",
    Status: 1,
}

token, expiryTime, err := jwtModule.GenerateRefreshToken(ctx, data)
if err != nil {
    log.Error(ctx, "Failed to generate refresh token: %v", err)
    return
}

log.Info(ctx, "Refresh token generated", "token", token, "expires_at", expiryTime)
```

### Validate Access Token

```go
isValid, data := jwtModule.ParseAccessToken(ctx, accessToken)
if !isValid {
    log.Warn(ctx, "Invalid or expired access token")
    http.Error(w, "Unauthorized", http.StatusUnauthorized)
    return
}

log.Info(ctx, "Token valid", "user_id", data.UserID, "email", data.Email)
```

### Validate Refresh Token

```go
isValid, data := jwtModule.ParseRefreshToken(ctx, refreshToken)
if !isValid {
    log.Warn(ctx, "Invalid or expired refresh token")
    http.Error(w, "Unauthorized", http.StatusUnauthorized)
    return
}

// Generate new access token using validated data
newAccessToken, _, err := jwtModule.GenerateAccessToken(ctx, data)
```

## API Reference

### `New(ctx context.Context, cfg Config) (*module, error)`

Creates a new JWT module instance.

**Parameters:**
- `ctx` - Context for operations
- `cfg` - Configuration struct with token secrets and expiry times

**Returns:**
- `*module` - JWT module instance
- `error` - Configuration validation error (secret length validation)

**Errors:**
- "secret token must 64 characters" - Access token secret length invalid
- "secret refresh token must 64 characters" - Refresh token secret length invalid

**Example:**
```go
jwtModule, err := jwt.New(ctx, config)
if err != nil {
    panic(err) // Secrets must be exactly 64 characters
}
```

### `GenerateAccessToken(ctx context.Context, data Data) (token string, expiryTime time.Time, err error)`

Generates a short-lived access token.

**Parameters:**
- `ctx` - Context for tracing
- `data` - User data to embed in token

**Returns:**
- `token` - Signed JWT token string
- `expiryTime` - Token expiration timestamp (UTC)
- `err` - Error if token generation fails

**Default Expiry:** 60 minutes

**Example:**
```go
token, expiry, err := jwtModule.GenerateAccessToken(ctx, userData)
if err != nil {
    return err
}
```

### `GenerateRefreshToken(ctx context.Context, data Data) (token string, expiryTime time.Time, err error)`

Generates a long-lived refresh token.

**Parameters:**
- `ctx` - Context for tracing
- `data` - User data to embed in token

**Returns:**
- `token` - Signed JWT token string
- `expiryTime` - Token expiration timestamp (UTC)
- `err` - Error if token generation fails

**Default Expiry:** 30 days

**Example:**
```go
refreshToken, expiry, err := jwtModule.GenerateRefreshToken(ctx, userData)
```

### `ParseAccessToken(ctx context.Context, accessToken string) (isValid bool, data Data)`

Validates and extracts data from an access token.

**Parameters:**
- `ctx` - Context for logging and tracing
- `accessToken` - JWT token string to validate

**Returns:**
- `isValid` - Whether token is valid and not expired
- `data` - Extracted user data (empty if invalid)

**Validation Checks:**
- Token signature validity
- Token expiration status
- Token format correctness

**Example:**
```go
isValid, userData := jwtModule.ParseAccessToken(ctx, token)
if !isValid {
    // Token is invalid or expired
    return
}
// Use userData
```

### `ParseRefreshToken(ctx context.Context, refreshToken string) (isValid bool, data Data)`

Validates and extracts data from a refresh token.

**Parameters:**
- `ctx` - Context for logging and tracing
- `refreshToken` - JWT token string to validate

**Returns:**
- `isValid` - Whether token is valid and not expired
- `data` - Extracted user data (empty if invalid)

**Example:**
```go
isValid, userData := jwtModule.ParseRefreshToken(ctx, token)
```

## Using Custom Metadata

The JWT package provides flexible metadata storage for custom user data beyond the standard fields.

### Basic Custom Metadata

```go
// Generate token with custom metadata
data := jwt.Data{
    UserID: 123,
    Email:  "user@example.com",
    Status: 1,
}

// Add custom fields
data.SetMetadata("role", "admin")
data.SetMetadata("department", "engineering")
data.SetMetadata("permissions", []string{"read", "write", "delete"})

token, _, err := jwtModule.GenerateAccessToken(ctx, data)
```

### Access Metadata After Parsing

```go
isValid, userData := jwtModule.ParseAccessToken(ctx, token)
if !isValid {
    return
}

// Retrieve custom metadata
role, exists := userData.GetMetadata("role")
if exists {
    fmt.Println("User role:", role) // Output: User role: admin
}

// Type assertion for specific types
permissions, _ := userData.GetMetadata("permissions")
permList := permissions.([]string)
```

### Bulk Metadata Assignment

```go
data := jwt.Data{
    UserID: 123,
    Email:  "user@example.com",
}

// Set all metadata at once
customData := map[string]interface{}{
    "role":       "manager",
    "tenant_id":  456,
    "features":   []string{"billing", "analytics", "reporting"},
    "last_login": time.Now(),
}
data.SetAllMetadata(customData)

token, _, err := jwtModule.GenerateAccessToken(ctx, data)
```

### Complex Metadata Structure

```go
// Nested metadata
data := jwt.Data{
    UserID: 123,
    Email:  "user@example.com",
    Metadata: map[string]interface{}{
        "organization": map[string]interface{}{
            "id":   "org-123",
            "name": "Acme Corp",
        },
        "roles": []string{"admin", "developer"},
        "settings": map[string]interface{}{
            "notifications": true,
            "theme":         "dark",
        },
    },
}

token, _, err := jwtModule.GenerateAccessToken(ctx, data)

// Later, retrieve and use
isValid, userData := jwtModule.ParseAccessToken(ctx, token)
if isValid {
    orgData, _ := userData.GetMetadata("organization")
    org := orgData.(map[string]interface{})
    fmt.Println("Organization:", org["name"]) // Output: Organization: Acme Corp
}
```

## Data Structures

### `Data` Struct

Represents user data embedded in JWT tokens. Supports both standard fields and flexible custom metadata.

```go
type Data struct {
    UserID   int64                  `json:"user_id,omitempty"`
    Email    string                 `json:"email,omitempty"`
    Phone    string                 `json:"phone,omitempty"`
    Status   int                    `json:"status,omitempty"`
    Metadata map[string]interface{} `json:"metadata,omitempty"`
}
```

**Fields:**
- `UserID` - User identifier
- `Email` - User email address
- `Phone` - User phone number
- `Status` - User status code (e.g., 1=active, 0=inactive)
- `Metadata` - Custom key-value data for flexible storage

**Standard Usage:**
```go
data := jwt.Data{
    UserID: 123,
    Email:  "user@example.com",
    Phone:  "1234567890",
    Status: 1,
}
```

**With Custom Metadata:**
```go
data := jwt.Data{
    UserID: 123,
    Email:  "user@example.com",
    Phone:  "1234567890",
    Status: 1,
    Metadata: map[string]interface{}{
        "role":     "admin",
        "tenant":   "acme-corp",
        "features": []string{"billing", "analytics"},
    },
}
```

### Data Helper Methods

#### `SetMetadata(key string, value interface{})`
Adds or updates a metadata key-value pair.

```go
data := jwt.Data{UserID: 123, Email: "user@example.com"}
data.SetMetadata("role", "admin")
data.SetMetadata("tenant_id", 456)
data.SetMetadata("permissions", []string{"read", "write"})
```

#### `GetMetadata(key string) (interface{}, bool)`
Retrieves a metadata value by key.

```go
data := jwt.Data{
    UserID:   123,
    Metadata: map[string]interface{}{"role": "admin"},
}

role, exists := data.GetMetadata("role")
if exists {
    fmt.Println(role) // Output: admin
}
```

#### `SetAllMetadata(metadata map[string]interface{})`
Replaces all metadata with the provided map.

```go
data := jwt.Data{UserID: 123}
metadata := map[string]interface{}{
    "role":   "user",
    "tenant": "company-xyz",
}
data.SetAllMetadata(metadata)
```

### `Config` Struct

Configuration for JWT module initialization.

```go
type Config struct {
    AccessTokenSecret              string // Must be 64 characters
    RefreshTokenSecret             string // Must be 64 characters
    AccessTokenExpiryTimeInMinutes int    // Optional, defaults to 60
    RefreshTokenExpireTimeInDay    int    // Optional, defaults to 30
}
```

### `Claims` Struct

JWT claims structure (used internally).

```go
type Claims struct {
    Data
    jwt.RegisteredClaims
}
```

## Example: Complete Auth Flow

```go
import (
    "context"
    "net/http"
    "github.com/adityayuga/go-sdk/jwt"
    "github.com/adityayuga/go-sdk/log"
)

var jwtModule jwt.JWTMethod

func init() {
    ctx := context.Background()
    
    config := jwt.Config{
        AccessTokenSecret:              "your-64-char-access-secret-key-here",
        RefreshTokenSecret:             "your-64-char-refresh-secret-key-here",
        AccessTokenExpiryTimeInMinutes: 15,  // 15 minutes
        RefreshTokenExpireTimeInDay:     7,   // 7 days
    }
    
    var err error
    jwtModule, err = jwt.New(ctx, config)
    if err != nil {
        panic(err)
    }
}

// Login handler
func handleLogin(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    // Validate user credentials (pseudo-code)
    userData := jwt.Data{
        UserID: 123,
        Email:  "user@example.com",
        Phone:  "1234567890",
        Status: 1,
    }
    
    // Generate tokens
    accessToken, _, err := jwtModule.GenerateAccessToken(ctx, userData)
    if err != nil {
        log.Errorf(ctx, "Failed to generate access token: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    
    refreshToken, _, err := jwtModule.GenerateRefreshToken(ctx, userData)
    if err != nil {
        log.Errorf(ctx, "Failed to generate refresh token: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    
    // Return tokens to client
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "access_token":  accessToken,
        "refresh_token": refreshToken,
    })
}

// Protected route handler
func handleProtectedRoute(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    // Extract token from Authorization header
    authHeader := r.Header.Get("Authorization")
    token := strings.TrimPrefix(authHeader, "Bearer ")
    
    // Validate token
    isValid, userData := jwtModule.ParseAccessToken(ctx, token)
    if !isValid {
        log.Warnf(ctx, "Invalid token")
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    
    // Use userData
    log.Infof(ctx, "User %d accessed endpoint", userData.UserID)
}

// Refresh token handler
func handleRefresh(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    // Get refresh token from request
    var req struct {
        RefreshToken string `json:"refresh_token"`
    }
    json.NewDecoder(r.Body).Decode(&req)
    
    // Validate refresh token
    isValid, userData := jwtModule.ParseRefreshToken(ctx, req.RefreshToken)
    if !isValid {
        http.Error(w, "Unauthorized", http.StatusUnauthorized)
        return
    }
    
    // Generate new access token
    newAccessToken, _, err := jwtModule.GenerateAccessToken(ctx, userData)
    if err != nil {
        log.Errorf(ctx, "Failed to generate new access token: %v", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(map[string]string{
        "access_token": newAccessToken,
    })
}
```

## Important Notes

1. **Secret Key Protection**
   - Store secrets in environment variables or secure key management system
   - Never commit secrets to version control
   - Use different secrets for access and refresh tokens
   - Rotate secrets periodically

2. **Token Expiry**
   - Access tokens: Short-lived (default 60 minutes) for security
   - Refresh tokens: Long-lived (default 30 days) for user convenience
   - Adjust based on your security requirements

3. **Signature Algorithm**
   - Uses HS512 (HMAC SHA-512) for signing
   - Requires exactly 64-character secrets
   - Symmetric algorithm (same secret for signing and validation)

4. **Context Usage**
   - Always provide a valid context for tracing and debug ID extraction
   - Context enables automatic error logging with debug information

5. **Error Logging**
   - Parse failures are automatically logged with debug ID
   - Useful for monitoring and troubleshooting token validation issues

6. **Token Validation**
   - Always validate token signature and expiration
   - Invalid tokens return `false` for `isValid` without throwing errors
   - Errors are logged for debugging purposes

7. **JWT Standard Compliance**
   - Follows RFC 7519 JWT standard
   - Uses standard `github.com/golang-jwt/jwt/v5` library
   - Compatible with JWT debuggers and validators

## Best Practices

1. **Use HTTPS** - Always transmit tokens over HTTPS to prevent interception
2. **Token Storage** - Store refresh tokens securely (HttpOnly cookies preferred)
3. **Token Rotation** - Implement token rotation for enhanced security
4. **Scope/Permissions** - Consider adding permissions/scopes to Data struct
5. **Blacklisting** - Implement token blacklist for logout functionality
6. **Monitoring** - Monitor failed token validation attempts for security
7. **Clock Skew** - Account for minor time differences between servers (JWT library handles this)
8. **Custom Claims** - Extend Data struct for additional user information as needed

## Troubleshooting

### "secret token must 64 characters"
- Access token secret must be exactly 64 characters
- Use `openssl rand -base64 64 | head -c 64` to generate

### "secret refresh token must 64 characters"
- Refresh token secret must be exactly 64 characters
- Generate using the same method as access token secret

### Token parsing returns false (invalid)
- Token may be expired
- Token signature may be invalid (wrong secret used)
- Token format may be corrupted
- Check logs for detailed error information with debug ID

### Token expires too quickly
- Check `AccessTokenExpiryTimeInMinutes` configuration
- Ensure system time is synchronized (NTP)
- Verify token generation and validation use same UTC timezone
