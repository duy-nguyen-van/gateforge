package auth

import (
	"fmt"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"

	"github.com/golang-jwt/jwt/v5"
)

// Claims carries the subject (user ID) and active tenant.
type Claims struct {
	TenantID string `json:"tenant_id,omitempty"`
	jwt.RegisteredClaims
}

// SelectionClaims is a short-lived token for tenant picker (sub only, no tenant).
type SelectionClaims struct {
	jwt.RegisteredClaims
}

// ProvideTokenService wires JWT for fx. Fails fast if JWT_SECRET is too short.
func ProvideTokenService(cfg *config.Config) (*TokenService, error) {
	return NewTokenService(cfg.JWTSecret, cfg.AppName, cfg.JWTAccessTTL)
}

// TokenService issues and verifies access tokens (HS256).
type TokenService struct {
	secret []byte
	issuer string
	ttl    time.Duration // access token lifetime
}

// NewTokenService builds a token service. secret must be non-empty.
func NewTokenService(secret string, issuer string, ttl time.Duration) (*TokenService, error) {
	if len(secret) < 32 {
		return nil, fmt.Errorf("jwt secret must be at least 32 bytes")
	}
	if ttl <= 0 {
		return nil, fmt.Errorf("jwt ttl must be positive")
	}
	return &TokenService{secret: []byte(secret), issuer: issuer, ttl: ttl}, nil
}

// SignAccessToken returns a Bearer token and its expiry time.
func (s *TokenService) SignAccessToken(userID, tenantID string) (token string, expiresAt time.Time, err error) {
	now := time.Now().UTC()
	expiresAt = now.Add(s.ttl)
	claims := Claims{
		TenantID: tenantID,
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			Issuer:    s.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := t.SignedString(s.secret)
	if err != nil {
		return "", time.Time{}, err
	}
	return signed, expiresAt, nil
}

// SignSelectionToken returns a short-lived token for tenant selection after auth.
func (s *TokenService) SignSelectionToken(userID string) (token string, expiresIn int64, err error) {
	now := time.Now().UTC()
	expiresAt := now.Add(constants.SelectionTokenTTL)
	claims := SelectionClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			Issuer:    s.issuer,
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(expiresAt),
		},
	}
	t := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	signed, err := t.SignedString(s.secret)
	if err != nil {
		return "", 0, err
	}
	return signed, int64(constants.SelectionTokenTTL.Seconds()), nil
}

// ParseAccessToken validates the token and returns user ID and tenant ID.
func (s *TokenService) ParseAccessToken(tokenString string) (userID, tenantID string, err error) {
	t, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return s.secret, nil
	})
	if err != nil {
		return "", "", err
	}
	claims, ok := t.Claims.(*Claims)
	if !ok || !t.Valid {
		return "", "", fmt.Errorf("invalid token claims")
	}
	if claims.Subject == "" {
		return "", "", fmt.Errorf("missing subject")
	}
	return claims.Subject, claims.TenantID, nil
}

// ParseSelectionToken validates a tenant-selection token and returns user ID.
func (s *TokenService) ParseSelectionToken(tokenString string) (userID string, err error) {
	t, err := jwt.ParseWithClaims(tokenString, &SelectionClaims{}, func(t *jwt.Token) (interface{}, error) {
		if t.Method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("unexpected signing method")
		}
		return s.secret, nil
	})
	if err != nil {
		return "", err
	}
	claims, ok := t.Claims.(*SelectionClaims)
	if !ok || !t.Valid {
		return "", fmt.Errorf("invalid token claims")
	}
	if claims.Subject == "" {
		return "", fmt.Errorf("missing subject")
	}
	return claims.Subject, nil
}
