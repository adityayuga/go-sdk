package jwt

import (
	"context"
	"errors"

	"time"

	"github.com/adityayuga/go-sdk/debug"
	"github.com/adityayuga/go-sdk/log"
	jwt "github.com/golang-jwt/jwt/v5"
)

type (
	module struct {
		accessTokenSecret  string
		refreshTokenSecret string

		accessTokenExpiryTimeInMinutes int
		refreshTokenExpireTimeInDay    int
	}

	Claims struct {
		Data
		jwt.RegisteredClaims
	}

	Data struct {
		UserID   int64                  `json:"user_id,omitempty"`
		Email    string                 `json:"email,omitempty"`
		Phone    string                 `json:"phone,omitempty"`
		Status   int                    `json:"status,omitempty"`
		Metadata map[string]interface{} `json:"metadata,omitempty"`
	}

	Config struct {
		AccessTokenSecret              string
		RefreshTokenSecret             string
		AccessTokenExpiryTimeInMinutes int
		RefreshTokenExpireTimeInDay    int
	}

	JWTMethod interface {
		GenerateAccessToken(ctx context.Context, data Data) (token string, expiryTime time.Time, err error)
		GenerateRefreshToken(ctx context.Context, data Data) (token string, expiryTime time.Time, err error)

		ParseAccessToken(ctx context.Context, accessToken string) (isValid bool, data Data)
		ParseRefreshToken(ctx context.Context, refreshToken string) (isValid bool, data Data)
	}
)

var (
	defaultAccessTokenExpireTime  int = 60 // 60 minutes
	defaultRefreshTokenExpireTime int = 30 // 30 days
)

func New(ctx context.Context, cfg Config) (*module, error) {
	// validate secret, must length 64 chars to support HS512
	if len(cfg.AccessTokenSecret) != 64 {
		return nil, errors.New("secret token must 64 characters")
	}
	if len(cfg.RefreshTokenSecret) != 64 {
		return nil, errors.New("secret refresh token must 64 characters")
	}
	if cfg.RefreshTokenExpireTimeInDay == 0 {
		cfg.RefreshTokenExpireTimeInDay = defaultRefreshTokenExpireTime // default 30 days
	}
	if cfg.AccessTokenExpiryTimeInMinutes == 0 {
		cfg.AccessTokenExpiryTimeInMinutes = defaultAccessTokenExpireTime // default 60 minutes
	}
	return &module{
		accessTokenSecret:              cfg.AccessTokenSecret,
		refreshTokenSecret:             cfg.RefreshTokenSecret,
		accessTokenExpiryTimeInMinutes: cfg.AccessTokenExpiryTimeInMinutes,
		refreshTokenExpireTimeInDay:    cfg.RefreshTokenExpireTimeInDay,
	}, nil
}

// GenerateAccessToken return short lived access token and error
func (m module) GenerateAccessToken(ctx context.Context, data Data) (token string, expiryTime time.Time, err error) {
	expiryTime = time.Now().UTC().Add(time.Duration(m.accessTokenExpiryTimeInMinutes) * time.Minute)
	claims := &Claims{
		Data: data,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiryTime),
		},
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	token, err = jwtToken.SignedString([]byte(m.accessTokenSecret))
	return
}

// GenerateRefreshToken return long lived token and error
func (m module) GenerateRefreshToken(ctx context.Context, data Data) (token string, expiryTime time.Time, err error) {
	expiryTime = time.Now().UTC().Add(time.Duration(m.refreshTokenExpireTimeInDay) * 24 * time.Hour)
	claims := &Claims{
		Data: data,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiryTime),
		},
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS512, claims)
	token, err = jwtToken.SignedString([]byte(m.refreshTokenSecret))
	return
}

// ParseRefreshToken
func (m module) ParseRefreshToken(ctx context.Context, refreshToken string) (isValid bool, data Data) {
	claims := &Claims{}
	debugID, _ := debug.GetDebugIDFromContext(ctx)

	parsedToken, err := jwt.ParseWithClaims(refreshToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.refreshTokenSecret), nil
	})
	if err != nil {
		log.Errorf(ctx, "[%s][pkgJwt] error: %s", debugID, err)
		return
	}

	if !parsedToken.Valid {
		return
	}

	return true, claims.Data
}

// ParseAccessToken
func (m module) ParseAccessToken(ctx context.Context, accessToken string) (isValid bool, data Data) {
	claims := &Claims{}
	debugID, _ := debug.GetDebugIDFromContext(ctx)

	parsedToken, err := jwt.ParseWithClaims(accessToken, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(m.accessTokenSecret), nil
	})
	if err != nil {
		log.Errorf(ctx, "[%s][pkgJwt] error: %s", debugID, err)
		return
	}
	if !parsedToken.Valid {
		return
	}

	return true, claims.Data
}

// SetMetadata adds or updates a metadata key-value pair in the data
func (d *Data) SetMetadata(key string, value interface{}) {
	if d.Metadata == nil {
		d.Metadata = make(map[string]interface{})
	}
	d.Metadata[key] = value
}

// GetMetadata retrieves a metadata value by key
func (d *Data) GetMetadata(key string) (interface{}, bool) {
	if d.Metadata == nil {
		return nil, false
	}
	val, exists := d.Metadata[key]
	return val, exists
}

// SetAllMetadata replaces all metadata with the provided map
func (d *Data) SetAllMetadata(metadata map[string]interface{}) {
	d.Metadata = metadata
}
