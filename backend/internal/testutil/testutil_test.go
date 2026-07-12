package testutil

import (
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestInitLogger(t *testing.T) {
	InitLogger()
	require.NotNil(t, TestConfig())
}

func TestTestConfig(t *testing.T) {
	cfg := TestConfig()
	require.NotEmpty(t, cfg.JWTSecret)
	require.NotEmpty(t, cfg.DefaultTenantID)
}

func TestNewEchoContext(t *testing.T) {
	c, rec := NewEchoContext(http.MethodPost, "/api/v1/test", `{"ok":true}`)
	require.NotNil(t, c)
	require.Equal(t, http.MethodPost, c.Request().Method)
	require.Equal(t, http.StatusOK, rec.Code) // no handler invoked yet
}
