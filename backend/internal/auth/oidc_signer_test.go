package auth

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/json"
	"encoding/pem"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gateforge-iam/gateforge-iam/internal/config"
	"github.com/gateforge-iam/gateforge-iam/internal/testutil"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/require"
)

func generateTestRSAPEM(t *testing.T) (pkcs1PEM, pkcs8PEM []byte) {
	t.Helper()
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	pkcs1PEM = pem.EncodeToMemory(&pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: x509.MarshalPKCS1PrivateKey(key),
	})
	pkcs8Bytes, err := x509.MarshalPKCS8PrivateKey(key)
	require.NoError(t, err)
	pkcs8PEM = pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: pkcs8Bytes,
	})
	return pkcs1PEM, pkcs8PEM
}

func TestProvideOIDCSigner_DevAutoGenerateKey(t *testing.T) {
	cfg := testutil.TestConfig()
	cfg.OIDCRSAPrivateKeyPEM = ""
	cfg.OIDCRSAPrivateKeyPath = ""

	signer, err := ProvideOIDCSigner(cfg)
	require.NoError(t, err)
	require.NotNil(t, signer)
	require.NotNil(t, signer.privateKey)
}

func TestProvideOIDCSigner_ProductionWithoutKeyFails(t *testing.T) {
	cfg := testutil.TestConfig()
	cfg.AppEnv = config.EnvironmentProduction
	cfg.OIDCRSAPrivateKeyPEM = ""
	cfg.OIDCRSAPrivateKeyPath = ""

	_, err := ProvideOIDCSigner(cfg)
	require.Error(t, err)
	require.Contains(t, err.Error(), "required in production")
}

func TestLoadOrGenerateRSAKey_FromPEMEnv(t *testing.T) {
	_, pkcs8PEM := generateTestRSAPEM(t)
	cfg := testutil.TestConfig()
	cfg.OIDCRSAPrivateKeyPEM = string(pkcs8PEM)

	key, err := loadOrGenerateRSAKey(cfg)
	require.NoError(t, err)
	require.NotNil(t, key)
}

func TestLoadOrGenerateRSAKey_FromPEMFile(t *testing.T) {
	pkcs1PEM, _ := generateTestRSAPEM(t)
	path := filepath.Join(t.TempDir(), "oidc.pem")
	require.NoError(t, os.WriteFile(path, pkcs1PEM, 0o600))

	cfg := testutil.TestConfig()
	cfg.OIDCRSAPrivateKeyPath = path

	key, err := loadOrGenerateRSAKey(cfg)
	require.NoError(t, err)
	require.NotNil(t, key)
}

func TestParseRSAPrivateKey_PKCS1AndPKCS8(t *testing.T) {
	pkcs1PEM, pkcs8PEM := generateTestRSAPEM(t)

	key1, err := parseRSAPrivateKey(pkcs1PEM)
	require.NoError(t, err)
	require.NotNil(t, key1)

	key8, err := parseRSAPrivateKey(pkcs8PEM)
	require.NoError(t, err)
	require.NotNil(t, key8)
}

func TestParseRSAPrivateKey_Errors(t *testing.T) {
	_, err := parseRSAPrivateKey([]byte("not-pem"))
	require.Error(t, err)
	require.Contains(t, err.Error(), "no PEM block")

	ecPEM := pem.EncodeToMemory(&pem.Block{Type: "EC PRIVATE KEY", Bytes: []byte("bad")})
	_, err = parseRSAPrivateKey(ecPEM)
	require.Error(t, err)
	require.Contains(t, err.Error(), "unsupported PEM type")
}

func TestOIDCSigner_SignAndParseAccessTokenOIDC(t *testing.T) {
	signer := newTestOIDCSigner(t)

	token, exp, err := signer.SignAccessTokenOIDC("user-1", "client-aud", "openid profile", "client-1")
	require.NoError(t, err)
	require.NotEmpty(t, token)
	require.True(t, exp.After(time.Now()))

	claims, err := signer.ParseAccessTokenOIDC(token)
	require.NoError(t, err)
	require.Equal(t, "user-1", claims.Subject)
	require.Equal(t, "openid profile", claims.Scope)
	require.Equal(t, "client-1", claims.ClientID)
	require.Contains(t, claims.Audience, "client-aud")
}

func TestOIDCSigner_ParseAccessTokenOIDC_Errors(t *testing.T) {
	signer := newTestOIDCSigner(t)

	_, err := signer.ParseAccessTokenOIDC("not-a-jwt")
	require.Error(t, err)

	other, err := ProvideOIDCSigner(testutil.TestConfig())
	require.NoError(t, err)
	token, _, err := other.SignAccessTokenOIDC("user-1", "aud", "scope", "client")
	require.NoError(t, err)
	_, err = signer.ParseAccessTokenOIDC(token)
	require.Error(t, err)
}

func TestOIDCSigner_SignIDToken(t *testing.T) {
	signer := newTestOIDCSigner(t)

	accessToken, _, err := signer.SignAccessTokenOIDC("user-1", "client-aud", "openid", "client-1")
	require.NoError(t, err)

	idToken, err := signer.SignIDToken("user-1", "client-aud", "nonce-123", accessToken, &OIDCUserClaims{
		Email:         "user@example.com",
		EmailVerified: true,
		Name:          "Jane Doe",
		GivenName:     "Jane",
		FamilyName:    "Doe",
	})
	require.NoError(t, err)
	require.NotEmpty(t, idToken)

	parser := jwt.NewParser(jwt.WithValidMethods([]string{"RS256"}), jwt.WithIssuer(signer.issuer))
	parsed, err := parser.ParseWithClaims(idToken, &OIDCIDClaims{}, func(t *jwt.Token) (any, error) {
		return &signer.privateKey.PublicKey, nil
	})
	require.NoError(t, err)

	claims, ok := parsed.Claims.(*OIDCIDClaims)
	require.True(t, ok)
	require.Equal(t, "user-1", claims.Subject)
	require.Equal(t, "nonce-123", claims.Nonce)
	require.Equal(t, AccessTokenHash(accessToken), claims.AtHash)
	require.Equal(t, "user@example.com", claims.Email)
	require.True(t, claims.EmailVerified)
	require.Equal(t, "Jane Doe", claims.Name)
}

func TestOIDCSigner_MarshalJWKS(t *testing.T) {
	signer := newTestOIDCSigner(t)
	cfg := testutil.TestConfig()
	signer.keyID = cfg.OIDCKeyID

	raw, err := signer.MarshalJWKS()
	require.NoError(t, err)

	var resp JWKSResponse
	require.NoError(t, json.Unmarshal(raw, &resp))
	require.Len(t, resp.Keys, 1)
	require.Equal(t, "RSA", resp.Keys[0].Kty)
	require.Equal(t, cfg.OIDCKeyID, resp.Keys[0].Kid)
	require.Equal(t, "sig", resp.Keys[0].Use)
	require.Equal(t, "RS256", resp.Keys[0].Alg)
	require.NotEmpty(t, resp.Keys[0].N)
	require.NotEmpty(t, resp.Keys[0].E)
}

func TestAccessTokenHash(t *testing.T) {
	hash := AccessTokenHash("my-access-token")
	require.NotEmpty(t, hash)
	require.Equal(t, hash, AccessTokenHash("my-access-token"))
	require.NotEqual(t, hash, AccessTokenHash("other-token"))
}

func newTestOIDCSigner(t *testing.T) *OIDCSigner {
	t.Helper()
	cfg := testutil.TestConfig()
	cfg.OIDCAccessTTL = time.Hour
	cfg.OIDCIDTokenTTL = time.Hour
	signer, err := ProvideOIDCSigner(cfg)
	require.NoError(t, err)
	return signer
}
