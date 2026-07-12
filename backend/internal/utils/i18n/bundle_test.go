package i18n

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/request"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestLanguageCandidates(t *testing.T) {
	require.Equal(t, []string{"en"}, languageCandidates(""))
	require.Equal(t, []string{"en"}, languageCandidates("  "))
	require.Equal(t, []string{"fr", "en"}, languageCandidates("fr"))
	require.Equal(t, []string{"en-US", "en", "en"}, languageCandidates("en-US"))
	require.Equal(t, []string{"en"}, languageCandidates("en"))
}

func TestLocalize_KnownAndUnknownKeys(t *testing.T) {
	Init()
	require.Equal(t, "Unauthorized", localize("en", constants.Unauthorized, nil))
	require.Equal(t, "UNKNOWN_CODE", localize("en", "UNKNOWN_CODE", nil))
	require.Equal(t, "Unauthorized", localize("fr", constants.Unauthorized, nil))
}

func TestT_NilContextUsesDefaultLanguage(t *testing.T) {
	Init()
	msg := T(nil, constants.Unauthorized, nil)
	require.Equal(t, "Unauthorized", msg)
}

func TestT_NilRequestUsesDefaultLanguage(t *testing.T) {
	Init()
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.SetRequest(nil)

	msg := T(c, constants.Unauthorized, nil)
	require.Equal(t, "Unauthorized", msg)
}

func TestTFromContext_NilContext(t *testing.T) {
	Init()
	msg := TFromContext(nil, constants.Unauthorized, nil)
	require.Equal(t, "Unauthorized", msg)
}

func TestLanguageFromContext(t *testing.T) {
	require.Equal(t, "en", languageFromContext(nil))
	ctx := request.NewLanguageCodeContext(context.Background(), "fr")
	require.Equal(t, "fr", languageFromContext(ctx))
	ctx = request.NewLanguageCodeContext(context.Background(), "")
	require.Equal(t, "en", languageFromContext(ctx))
}

func TestLocalize_NilMessagesAndEmptyValue(t *testing.T) {
	Init()
	saved := messages
	t.Cleanup(func() { messages = saved })

	messages = nil
	require.Equal(t, "KEY", localize("en", "KEY", nil))

	messages = map[string]map[string]string{"en": {"EMPTY": ""}}
	require.Equal(t, "EMPTY", localize("en", "EMPTY", nil))
}
