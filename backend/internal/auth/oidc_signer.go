package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"math/big"
	"os"
	"strings"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/config"

	"github.com/golang-jwt/jwt/v5"
)

// OIDCSigner issues and verifies OIDC access and ID tokens (RS256, kid in header).
type OIDCSigner struct {
	privateKey    *rsa.PrivateKey
	keyID         string
	previous      *rsa.PrivateKey
	previousKeyID string
	issuer        string
	accessTTL     time.Duration
	idTTL         time.Duration
	jwksJSON      []byte
}

// ProvideOIDCSigner loads RSA key from config or generates one in non-production when unset.
func ProvideOIDCSigner(cfg *config.Config) (*OIDCSigner, error) {
	issuer := cfg.AppBaseURL
	if issuer == "" {
		issuer = "http://localhost:3000"
	}

	key, err := loadOrGenerateRSAKey(cfg)
	if err != nil {
		return nil, err
	}
	var previous *rsa.PrivateKey
	if strings.TrimSpace(cfg.OIDCPreviousPrivateKeyPEM) != "" {
		previous, err = parseRSAPrivateKey([]byte(cfg.OIDCPreviousPrivateKeyPEM))
		if err != nil {
			return nil, fmt.Errorf("previous OIDC RSA key: %w", err)
		}
	}

	accessTTL := cfg.OIDCAccessTTL
	if accessTTL <= 0 {
		accessTTL = cfg.JWTAccessTTL
	}
	idTTL := cfg.OIDCIDTokenTTL
	if idTTL <= 0 {
		idTTL = accessTTL
	}

	signer := &OIDCSigner{
		privateKey:    key,
		keyID:         cfg.OIDCKeyID,
		previous:      previous,
		previousKeyID: cfg.OIDCPreviousKeyID,
		issuer:        issuer,
		accessTTL:     accessTTL,
		idTTL:         idTTL,
	}
	raw, err := signer.marshalJWKS()
	if err != nil {
		return nil, err
	}
	signer.jwksJSON = raw
	return signer, nil
}

func loadOrGenerateRSAKey(cfg *config.Config) (*rsa.PrivateKey, error) {
	if cfg.OIDCRSAPrivateKeyPath != "" {
		pemBytes, err := os.ReadFile(cfg.OIDCRSAPrivateKeyPath)
		if err != nil {
			return nil, fmt.Errorf("read OIDC RSA key file: %w", err)
		}
		return parseRSAPrivateKey(pemBytes)
	}
	if cfg.OIDCRSAPrivateKeyPEM != "" {
		return parseRSAPrivateKey([]byte(cfg.OIDCRSAPrivateKeyPEM))
	}
	if cfg.AppEnv == config.EnvironmentProduction {
		return nil, fmt.Errorf("OIDC_RSA_PRIVATE_KEY_FILE or OIDC_RSA_PRIVATE_KEY_PEM is required in production")
	}
	k, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return nil, fmt.Errorf("generate OIDC RSA key: %w", err)
	}
	return k, nil
}

func parseRSAPrivateKey(pemBytes []byte) (*rsa.PrivateKey, error) {
	block, _ := pem.Decode(normalizePEM(pemBytes))
	if block == nil {
		return nil, fmt.Errorf("OIDC RSA PEM: no PEM block")
	}
	switch block.Type {
	case "RSA PRIVATE KEY":
		k, err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("parse PKCS1 RSA private key: %w", err)
		}
		return k, nil
	case "PRIVATE KEY":
		k, err := x509.ParsePKCS8PrivateKey(block.Bytes)
		if err != nil {
			return nil, fmt.Errorf("parse PKCS8 private key: %w", err)
		}
		rk, ok := k.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("OIDC private key is not RSA")
		}
		return rk, nil
	default:
		return nil, fmt.Errorf("unsupported PEM type %q (expected RSA PRIVATE KEY or PRIVATE KEY)", block.Type)
	}
}

// normalizePEM turns a one-line env value into PEM. Double-quoted .env values
// often keep the two-character sequence \n instead of a real line break.
func normalizePEM(pemBytes []byte) []byte {
	s := strings.TrimSpace(string(pemBytes))
	s = strings.Trim(s, `"'`)
	if strings.Contains(s, `\n`) {
		s = strings.ReplaceAll(s, `\n`, "\n")
	}
	return []byte(s)
}

// OIDCAccessClaims is the access token JWT payload (resource server / userinfo).
type OIDCAccessClaims struct {
	jwt.RegisteredClaims
	Scope    string `json:"scope,omitempty"`
	ClientID string `json:"client_id,omitempty"`
}

// OIDCIDClaims is the OpenID Connect ID token payload.
type OIDCIDClaims struct {
	jwt.RegisteredClaims
	Nonce         string `json:"nonce,omitempty"`
	AtHash        string `json:"at_hash,omitempty"`
	Email         string `json:"email,omitempty"`
	EmailVerified bool   `json:"email_verified,omitempty"`
	Name          string `json:"name,omitempty"`
	GivenName     string `json:"given_name,omitempty"`
	FamilyName    string `json:"family_name,omitempty"`
}

// SignAccessTokenOIDC returns an RS256 JWT for the OAuth2 access token.
func (s *OIDCSigner) SignAccessTokenOIDC(userID, audience, scope, clientID string) (string, time.Time, error) {
	now := time.Now().UTC()
	exp := now.Add(s.accessTTL)
	claims := OIDCAccessClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			Issuer:    s.issuer,
			Audience:  jwt.ClaimStrings{audience},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
		Scope:    scope,
		ClientID: clientID,
	}
	t := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	t.Header["kid"] = s.keyID
	signed, err := t.SignedString(s.privateKey)
	if err != nil {
		return "", time.Time{}, err
	}
	return signed, exp, nil
}

// SignIDToken returns an RS256 OpenID Connect ID token.
func (s *OIDCSigner) SignIDToken(userID, audience, nonce, accessToken string, u *OIDCUserClaims) (string, error) {
	now := time.Now().UTC()
	exp := now.Add(s.idTTL)
	claims := OIDCIDClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Subject:   userID,
			Issuer:    s.issuer,
			Audience:  jwt.ClaimStrings{audience},
			IssuedAt:  jwt.NewNumericDate(now),
			ExpiresAt: jwt.NewNumericDate(exp),
		},
		Nonce: nonce,
	}
	if accessToken != "" {
		claims.AtHash = AccessTokenHash(accessToken)
	}
	if u != nil {
		claims.Email = u.Email
		claims.EmailVerified = u.EmailVerified
		claims.Name = u.Name
		claims.GivenName = u.GivenName
		claims.FamilyName = u.FamilyName
	}
	t := jwt.NewWithClaims(jwt.SigningMethodRS256, claims)
	t.Header["kid"] = s.keyID
	return t.SignedString(s.privateKey)
}

// OIDCUserClaims holds optional profile claims for the ID token.
type OIDCUserClaims struct {
	Email         string
	EmailVerified bool
	Name          string
	GivenName     string
	FamilyName    string
}

// ParseAccessTokenOIDC validates an RS256 access token and returns claims.
func (s *OIDCSigner) ParseAccessTokenOIDC(tokenString string) (*OIDCAccessClaims, error) {
	parser := jwt.NewParser(
		jwt.WithValidMethods([]string{"RS256"}),
		jwt.WithIssuer(s.issuer),
	)
	t, err := parser.ParseWithClaims(tokenString, &OIDCAccessClaims{}, s.verificationKey)
	if err != nil {
		return nil, err
	}
	claims, ok := t.Claims.(*OIDCAccessClaims)
	if !ok || !t.Valid {
		return nil, fmt.Errorf("invalid token claims")
	}
	if claims.Subject == "" {
		return nil, fmt.Errorf("missing subject")
	}
	return claims, nil
}

// JWKSResponse is the JSON body for /.well-known/jwks.json.
type JWKSResponse struct {
	Keys []JWK `json:"keys"`
}

// JWK is a minimal RSA public key for JWKS.
type JWK struct {
	Kty string `json:"kty"`
	Kid string `json:"kid"`
	Use string `json:"use"`
	Alg string `json:"alg"`
	N   string `json:"n"`
	E   string `json:"e"`
}

func (s *OIDCSigner) verificationKey(t *jwt.Token) (interface{}, error) {
	if t.Method != jwt.SigningMethodRS256 {
		return nil, fmt.Errorf("unexpected signing method")
	}
	kid, _ := t.Header["kid"].(string)
	if s.previous != nil && kid != "" && kid == s.previousKeyID {
		return &s.previous.PublicKey, nil
	}
	if kid == "" || kid == s.keyID {
		return &s.privateKey.PublicKey, nil
	}
	return nil, fmt.Errorf("unknown signing key")
}

// MarshalJWKS returns the precomputed JWKS JSON for the active and previous public keys.
func (s *OIDCSigner) MarshalJWKS() ([]byte, error) {
	if len(s.jwksJSON) > 0 {
		out := make([]byte, len(s.jwksJSON))
		copy(out, s.jwksJSON)
		return out, nil
	}
	return s.marshalJWKS()
}

func (s *OIDCSigner) marshalJWKS() ([]byte, error) {
	keys := []JWK{publicJWK(s.privateKey, s.keyID)}
	if s.previous != nil && s.previousKeyID != "" {
		keys = append(keys, publicJWK(s.previous, s.previousKeyID))
	}
	return json.Marshal(JWKSResponse{Keys: keys})
}

func publicJWK(key *rsa.PrivateKey, kid string) JWK {
	pub := key.PublicKey
	return JWK{
		Kty: "RSA",
		Kid: kid,
		Use: "sig",
		Alg: "RS256",
		N:   base64.RawURLEncoding.EncodeToString(pub.N.Bytes()),
		E:   base64.RawURLEncoding.EncodeToString(big.NewInt(int64(pub.E)).Bytes()),
	}
}

// AccessTokenHash computes the at_hash claim (first half of SHA-256 of access_token), base64url.
func AccessTokenHash(accessToken string) string {
	sum := sha256.Sum256([]byte(accessToken))
	half := sum[:len(sum)/2]
	return base64.RawURLEncoding.EncodeToString(half)
}
