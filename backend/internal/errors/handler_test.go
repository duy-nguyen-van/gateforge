package errors

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/testutil"
	"github.com/gateforge-iam/gateforge-iam/internal/utils/i18n"

	"github.com/getsentry/sentry-go"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
	"gorm.io/gorm"
)

func TestMain(m *testing.M) {
	testutil.InitLogger()
	i18n.Init()
	os.Exit(m.Run())
}

func TestNewErrorHandler(t *testing.T) {
	require.NotNil(t, NewErrorHandler())
}

func TestErrorHandler_processError(t *testing.T) {
	h := NewErrorHandler()

	t.Run("nil error", func(t *testing.T) {
		require.Nil(t, h.processError(nil))
	})

	t.Run("existing AppError", func(t *testing.T) {
		appErr := ValidationError("invalid", nil)
		require.Equal(t, appErr, h.processError(appErr))
	})

	t.Run("echo HTTP errors", func(t *testing.T) {
		tests := []struct {
			name       string
			httpErr    *echo.HTTPError
			code       string
			errType    ErrorType
			httpStatus int
		}{
			{
				name:       "bad request",
				httpErr:    echo.NewHTTPError(http.StatusBadRequest, "bad input"),
				code:       constants.ValidationError,
				errType:    ErrorTypeValidation,
				httpStatus: http.StatusBadRequest,
			},
			{
				name:       "unauthorized",
				httpErr:    echo.NewHTTPError(http.StatusUnauthorized, "login required"),
				code:       constants.Unauthorized,
				errType:    ErrorTypeUnauthorized,
				httpStatus: http.StatusUnauthorized,
			},
			{
				name:       "forbidden",
				httpErr:    echo.NewHTTPError(http.StatusForbidden, "no access"),
				code:       constants.Forbidden,
				errType:    ErrorTypeForbidden,
				httpStatus: http.StatusForbidden,
			},
			{
				name:       "not found",
				httpErr:    echo.NewHTTPError(http.StatusNotFound, "missing"),
				code:       constants.NotFound,
				errType:    ErrorTypeNotFound,
				httpStatus: http.StatusNotFound,
			},
			{
				name:       "default status",
				httpErr:    echo.NewHTTPError(http.StatusTeapot, "teapot"),
				code:       constants.InternalError,
				errType:    ErrorTypeInternal,
				httpStatus: http.StatusTeapot,
			},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := h.processError(tt.httpErr)
				require.Equal(t, tt.code, got.Code)
				require.Equal(t, tt.errType, got.Type)
				require.Equal(t, tt.httpStatus, got.HTTPStatus)
			})
		}
	})

	t.Run("gorm errors", func(t *testing.T) {
		tests := []struct {
			name    string
			err     error
			errType ErrorType
			code    string
		}{
			{name: "record not found", err: gorm.ErrRecordNotFound, errType: ErrorTypeNotFound, code: constants.NotFound},
			{name: "invalid transaction", err: gorm.ErrInvalidTransaction, errType: ErrorTypeDatabase, code: constants.DatabaseError},
			{name: "not implemented", err: gorm.ErrNotImplemented, errType: ErrorTypeDatabase, code: constants.DatabaseError},
			{name: "missing where", err: gorm.ErrMissingWhereClause, errType: ErrorTypeDatabase, code: constants.DatabaseError},
			{name: "unsupported driver", err: gorm.ErrUnsupportedDriver, errType: ErrorTypeDatabase, code: constants.DatabaseError},
			{name: "registered", err: gorm.ErrRegistered, errType: ErrorTypeDatabase, code: constants.DatabaseError},
			{name: "invalid field", err: gorm.ErrInvalidField, errType: ErrorTypeDatabase, code: constants.DatabaseError},
			{name: "empty slice", err: gorm.ErrEmptySlice, errType: ErrorTypeDatabase, code: constants.DatabaseError},
			{name: "dry run unsupported", err: gorm.ErrDryRunModeUnsupported, errType: ErrorTypeDatabase, code: constants.DatabaseError},
			{name: "invalid db", err: gorm.ErrInvalidDB, errType: ErrorTypeDatabase, code: constants.DatabaseError},
			{name: "invalid value", err: gorm.ErrInvalidValue, errType: ErrorTypeDatabase, code: constants.DatabaseError},
			{name: "invalid value length", err: gorm.ErrInvalidValueOfLength, errType: ErrorTypeDatabase, code: constants.DatabaseError},
			{name: "preload not allowed", err: gorm.ErrPreloadNotAllowed, errType: ErrorTypeDatabase, code: constants.DatabaseError},
			{name: "generic database error", err: errors.New("database connection reset"), errType: ErrorTypeDatabase, code: constants.DatabaseError},
		}

		for _, tt := range tests {
			t.Run(tt.name, func(t *testing.T) {
				got := h.processError(tt.err)
				require.Equal(t, tt.code, got.Code)
				require.Equal(t, tt.errType, got.Type)
			})
		}
	})

	t.Run("context errors", func(t *testing.T) {
		canceled := h.processError(context.Canceled)
		require.Equal(t, ErrorTypeTimeout, canceled.Type)
		require.Equal(t, http.StatusRequestTimeout, canceled.HTTPStatus)

		deadline := h.processError(context.DeadlineExceeded)
		require.Equal(t, ErrorTypeTimeout, deadline.Type)
		require.Equal(t, http.StatusRequestTimeout, deadline.HTTPStatus)
	})

	t.Run("default internal error", func(t *testing.T) {
		got := h.processError(errors.New("unexpected"))
		require.Equal(t, constants.InternalError, got.Code)
		require.Equal(t, ErrorTypeInternal, got.Type)
		require.Equal(t, http.StatusInternalServerError, got.HTTPStatus)
	})
}

func TestStringHelpers(t *testing.T) {
	require.True(t, isGORMError(errors.New("database connection reset")))
	require.True(t, isGORMError(errors.New("unique constraint violated")))
	require.True(t, containsIgnoreCase("database", "database"))
	require.True(t, containsIgnoreCase("database", "data"))
	require.True(t, containsIgnoreCase("connection", "tion"))
	require.True(t, containsSubstring("hello world", "lo wo"))
	require.True(t, containsSubstring("anything", ""))
	require.False(t, containsSubstring("short", "longer"))
	require.False(t, isGORMError(errors.New("plain failure")))
}

func TestErrorHandler_HandleError(t *testing.T) {
	h := NewErrorHandler()
	c, rec := testutil.NewEchoContext(http.MethodGet, "/api/v1/users", "")

	err := h.HandleError(c, ValidationError("invalid user", nil))
	require.NoError(t, err)
	require.Equal(t, http.StatusBadRequest, rec.Code)

	var resp dtos.BaseResponse[map[string]interface{}]
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
	require.Equal(t, constants.ValidationError, resp.Meta.ErrorCode)
	require.Equal(t, http.StatusBadRequest, resp.Meta.Code)
}

func TestErrorHandler_HandleError_WithSentryHub(t *testing.T) {
	h := NewErrorHandler()
	c, rec := testutil.NewEchoContext(http.MethodGet, "/api/v1/users", "")
	hub := sentry.CurrentHub().Clone()
	req := c.Request().WithContext(sentry.SetHubOnContext(c.Request().Context(), hub))
	c.SetRequest(req)

	appErr := InternalError("server blew up", errors.New("cause")).
		WithOperation("save").
		WithResource("user").
		WithContext("field", "email")
	require.NoError(t, h.HandleError(c, appErr))
	require.Equal(t, http.StatusInternalServerError, rec.Code)
}

func TestErrorHandler_SuccessResponse(t *testing.T) {
	h := NewErrorHandler()

	t.Run("without pagination", func(t *testing.T) {
		c, rec := testutil.NewEchoContext(http.MethodGet, "/api/v1/users", "")
		err := h.SuccessResponse(c, "ok", map[string]string{"id": "1"}, nil)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rec.Code)

		var resp dtos.BaseResponse[map[string]string]
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		require.Equal(t, constants.Success, resp.Meta.ErrorCode)
		require.Equal(t, "ok", resp.Meta.Message)
		require.Equal(t, "1", resp.Data["id"])
	})

	t.Run("with pagination", func(t *testing.T) {
		c, rec := testutil.NewEchoContext(http.MethodGet, "/api/v1/users", "")
		page := &dtos.Pageable{Page: 2, PageSize: 10, Total: 25}
		err := h.SuccessResponse(c, "listed", []string{"a"}, page)
		require.NoError(t, err)
		require.Equal(t, http.StatusOK, rec.Code)

		var resp dtos.BaseResponse[[]string]
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		require.Equal(t, 2, resp.Meta.Page)
		require.Equal(t, 10, resp.Meta.PageSize)
		require.Equal(t, int64(25), resp.Meta.Total)
	})
}

func TestErrorHandler_ConvenienceResponses(t *testing.T) {
	h := NewErrorHandler()

	tests := []struct {
		name      string
		call      func(echo.Context) error
		status    int
		errorCode string
	}{
		{
			name: "validation",
			call: func(c echo.Context) error {
				return h.ValidationErrorResponse(c, "invalid", map[string]string{"email": "bad"})
			},
			status:    http.StatusBadRequest,
			errorCode: constants.ValidationError,
		},
		{
			name: "not found",
			call: func(c echo.Context) error {
				return h.NotFoundErrorResponse(c, "user")
			},
			status:    http.StatusNotFound,
			errorCode: constants.NotFound,
		},
		{
			name: "unauthorized",
			call: func(c echo.Context) error {
				return h.UnauthorizedErrorResponse(c, "login required")
			},
			status:    http.StatusUnauthorized,
			errorCode: constants.Unauthorized,
		},
		{
			name: "forbidden",
			call: func(c echo.Context) error {
				return h.ForbiddenErrorResponse(c, "denied")
			},
			status:    http.StatusForbidden,
			errorCode: constants.Forbidden,
		},
		{
			name: "internal",
			call: func(c echo.Context) error {
				return h.InternalErrorResponse(c, "boom", errors.New("cause"))
			},
			status:    http.StatusInternalServerError,
			errorCode: constants.InternalError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c, rec := testutil.NewEchoContext(http.MethodGet, "/api/v1/resource", "")
			require.NoError(t, tt.call(c))
			require.Equal(t, tt.status, rec.Code)

			var resp dtos.BaseResponse[map[string]interface{}]
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
			require.Equal(t, tt.errorCode, resp.Meta.ErrorCode)
		})
	}
}

func TestErrorHandler_errorResponse(t *testing.T) {
	h := NewErrorHandler()

	t.Run("includes context data", func(t *testing.T) {
		c, rec := testutil.NewEchoContext(http.MethodPost, "/api/v1/users", "")
		appErr := ValidationErrorWithDetails("invalid", nil, map[string]string{"email": "bad"})
		require.NoError(t, h.errorResponse(c, appErr))
		require.Equal(t, http.StatusBadRequest, rec.Code)

		var resp dtos.BaseResponse[map[string]interface{}]
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		require.Equal(t, "bad", resp.Data["email"])
	})

	t.Run("falls back to app error message", func(t *testing.T) {
		c, rec := testutil.NewEchoContext(http.MethodGet, "/api/v1/users", "")
		appErr := NewAppError("UNKNOWN_CODE", "custom message", ErrorTypeInternal, http.StatusInternalServerError)
		require.NoError(t, h.errorResponse(c, appErr))

		var resp dtos.BaseResponse[any]
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &resp))
		require.Equal(t, "custom message", resp.Meta.Message)
	})
}

func TestErrorHandler_logErrorBranches(t *testing.T) {
	h := NewErrorHandler()
	c, _ := testutil.NewEchoContext(http.MethodGet, "/api/v1/users", "")

	logTypes := []ErrorType{
		ErrorTypeValidation,
		ErrorTypeNotFound,
		ErrorTypeUnauthorized,
		ErrorTypeForbidden,
		ErrorTypeInternal,
		ErrorTypeDatabase,
		ErrorTypeExternal,
		ErrorTypeCache,
		ErrorTypeTimeout,
		ErrorTypeConflict,
	}

	for _, errType := range logTypes {
		t.Run(string(errType), func(t *testing.T) {
			appErr := NewAppError(constants.InternalError, "log me", errType, http.StatusInternalServerError).
				WithOperation("op").
				WithResource("resource").
				WithContext("key", "value")
			h.logError(c, appErr)
		})
	}
}
