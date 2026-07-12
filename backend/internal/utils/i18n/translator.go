package i18n

import (
	"context"

	"github.com/gateforge-iam/gateforge-iam/internal/request"

	"github.com/labstack/echo/v4"
)

const defaultLanguage = "en"

// T resolves an error code to a localized message for the request language.
func T(c echo.Context, errorCode string, param map[string]interface{}) string {
	lang := defaultLanguage
	if c != nil && c.Request() != nil {
		lang = languageFromContext(c.Request().Context())
	}
	return localize(lang, errorCode, param)
}

// TFromContext resolves an error code using the language stored on context.Context.
func TFromContext(ctx context.Context, errorCode string, param map[string]interface{}) string {
	return localize(languageFromContext(ctx), errorCode, param)
}

func languageFromContext(ctx context.Context) string {
	if ctx == nil {
		return defaultLanguage
	}
	if lang, ok := request.LanguageCodeFromContext(ctx); ok && lang != "" {
		return lang
	}
	return defaultLanguage
}
