package middlewares

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/request"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestRequestContext(t *testing.T) {
	t.Run("uses incoming correlation id", func(t *testing.T) {
		e := echo.New()
		e.Use(RequestContext("gateforge-iam"))
		e.GET("/ctx", func(c echo.Context) error {
			cid, ok := request.CorrelationIDFromContext(c.Request().Context())
			require.True(t, ok)
			require.Equal(t, "existing-corr-id", cid)

			lang, ok := request.LanguageCodeFromContext(c.Request().Context())
			require.True(t, ok)
			require.Equal(t, "fr-CA", lang)

			ts, ok := request.RequestTimestampFromContext(c.Request().Context())
			require.True(t, ok)
			require.Positive(t, ts)

			url, ok := request.RequestURLFromContext(c.Request().Context())
			require.True(t, ok)
			require.Contains(t, url, "/ctx")

			return c.NoContent(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/ctx", nil)
		req.Header.Set(CorrelationIDHeaderKey, "existing-corr-id")
		req.Header.Set(LanguageCodeHeaderKey, "fr-CA,en")
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		require.Equal(t, "existing-corr-id", rec.Header().Get(CorrelationIDHeaderKey))
	})

	t.Run("generates correlation id when missing", func(t *testing.T) {
		e := echo.New()
		e.Use(RequestContext("myservice"))
		e.GET("/ctx", func(c echo.Context) error {
			cid, ok := request.CorrelationIDFromContext(c.Request().Context())
			require.True(t, ok)
			require.Contains(t, cid, "myservice-")
			return c.NoContent(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/ctx", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		require.NotEmpty(t, rec.Header().Get(CorrelationIDHeaderKey))
	})

	t.Run("defaults language to en", func(t *testing.T) {
		e := echo.New()
		e.Use(RequestContext("svc"))
		e.GET("/ctx", func(c echo.Context) error {
			lang, ok := request.LanguageCodeFromContext(c.Request().Context())
			require.True(t, ok)
			require.Equal(t, "en", lang)
			return c.NoContent(http.StatusOK)
		})

		req := httptest.NewRequest(http.MethodGet, "/ctx", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
	})
}

func TestParseLanguageCode(t *testing.T) {
	tests := []struct {
		name   string
		header string
		want   string
	}{
		{"empty", "", ""},
		{"single", "en", "en"},
		{"with quality", "en-US,fr;q=0.8", "en-US"},
		{"whitespace", "  de  ,en", "de"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, parseLanguageCode(tt.header))
		})
	}
}

func TestFormatRandomSuffix(t *testing.T) {
	tests := []struct {
		name string
		n    int64
		want string
	}{
		{"zero", 0, "000"},
		{"single digit", 5, "005"},
		{"double digit", 42, "042"},
		{"triple digit", 999, "999"},
		{"negative", -7, "007"},
		{"overflow", 1234, "234"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, formatRandomSuffix(tt.n))
		})
	}
}

func TestHeaderKeys(t *testing.T) {
	require.Equal(t, "X-Correlation-Id", CorrelationIDHeaderKey)
	require.Equal(t, "Accept-Language", LanguageCodeHeaderKey)
}
