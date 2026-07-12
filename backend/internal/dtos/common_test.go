package dtos

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/utils/i18n"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func newEchoContext(t *testing.T) echo.Context {
	t.Helper()
	i18n.Init()
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	return e.NewContext(req, rec)
}

func TestNewPageableRequest(t *testing.T) {
	pr := NewPageableRequest()
	require.Equal(t, constants.DefaultPage, pr.Page)
	require.Equal(t, constants.DefaultPageSize, pr.PageSize)
}

func TestNewUnpaginatedRequest(t *testing.T) {
	pr := NewUnpaginatedRequest()
	require.Equal(t, constants.DefaultPage, pr.Page)
	require.Equal(t, constants.NoLimit, pr.PageSize)
}

func TestPageableRequest_ShouldPaginate(t *testing.T) {
	require.True(t, (&PageableRequest{PageSize: 10}).ShouldPaginate())
	require.True(t, (&PageableRequest{PageSize: 1}).ShouldPaginate())
	require.False(t, (&PageableRequest{PageSize: 0}).ShouldPaginate())
	require.False(t, (&PageableRequest{PageSize: -1}).ShouldPaginate())
}

func TestPageableRequest_GetLimit(t *testing.T) {
	require.Equal(t, 20, (&PageableRequest{PageSize: 20}).GetLimit())
	require.Equal(t, constants.NoLimit, (&PageableRequest{PageSize: 0}).GetLimit())
	require.Equal(t, constants.NoLimit, (&PageableRequest{PageSize: -5}).GetLimit())
}

func TestPageableRequest_GetOffset(t *testing.T) {
	require.Equal(t, 20, (&PageableRequest{Page: 3, PageSize: 10}).GetOffset())
	require.Equal(t, 0, (&PageableRequest{Page: 1, PageSize: 10}).GetOffset())
	require.Equal(t, 0, (&PageableRequest{Page: 5, PageSize: 0}).GetOffset())
	require.Equal(t, 0, (&PageableRequest{Page: 5, PageSize: -1}).GetOffset())
}

func TestMeta_HttpCode(t *testing.T) {
	tests := []struct {
		name      string
		meta      Meta
		wantCode  int
	}{
		{
			name:     "explicit Code takes precedence",
			meta:     Meta{Code: http.StatusCreated, ErrorCode: "400_BAD"},
			wantCode: http.StatusCreated,
		},
		{
			name:     "short error code defaults to 500",
			meta:     Meta{ErrorCode: "AB"},
			wantCode: http.StatusInternalServerError,
		},
		{
			name:     "empty error code defaults to 500",
			meta:     Meta{},
			wantCode: http.StatusInternalServerError,
		},
		{
			name:     "200 prefix",
			meta:     Meta{ErrorCode: "200_OK"},
			wantCode: http.StatusOK,
		},
		{
			name:     "400 prefix",
			meta:     Meta{ErrorCode: "400_BAD_REQUEST"},
			wantCode: http.StatusBadRequest,
		},
		{
			name:     "401 prefix",
			meta:     Meta{ErrorCode: "401_UNAUTHORIZED"},
			wantCode: http.StatusUnauthorized,
		},
		{
			name:     "403 prefix",
			meta:     Meta{ErrorCode: "403_FORBIDDEN"},
			wantCode: http.StatusForbidden,
		},
		{
			name:     "404 prefix",
			meta:     Meta{ErrorCode: "404_NOT_FOUND"},
			wantCode: http.StatusNotFound,
		},
		{
			name:     "500 prefix",
			meta:     Meta{ErrorCode: "500_INTERNAL"},
			wantCode: http.StatusInternalServerError,
		},
		{
			name:     "unknown prefix defaults to 500",
			meta:     Meta{ErrorCode: "503_UNAVAILABLE"},
			wantCode: http.StatusInternalServerError,
		},
		{
			name:     "non-numeric error code defaults to 500",
			meta:     Meta{ErrorCode: "BAD_REQUEST"},
			wantCode: http.StatusInternalServerError,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.wantCode, tt.meta.HttpCode())
		})
	}
}

func TestGetMeta(t *testing.T) {
	c := newEchoContext(t)
	meta := GetMeta(c, constants.Unauthorized, http.StatusUnauthorized)

	require.Equal(t, constants.Unauthorized, meta.ErrorCode)
	require.Equal(t, "Unauthorized", meta.Message)
	require.Equal(t, http.StatusUnauthorized, meta.Code)
	require.Equal(t, http.StatusUnauthorized, meta.HttpCode())
}

func TestGetMetaPaging(t *testing.T) {
	c := newEchoContext(t)
	pageable := &Pageable{Page: 2, PageSize: 25, Total: 100}
	meta := GetMetaPaging(c, constants.Success, pageable, http.StatusOK)

	require.Equal(t, constants.Success, meta.ErrorCode)
	require.Equal(t, "Success", meta.Message)
	require.Equal(t, http.StatusOK, meta.Code)
	require.Equal(t, 2, meta.Page)
	require.Equal(t, 25, meta.PageSize)
	require.Equal(t, int64(100), meta.Total)
}

func TestBaseResponse_JSON(t *testing.T) {
	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	resp := &BaseResponse[string]{
		Meta: Meta{Code: http.StatusOK, ErrorCode: constants.Success, Message: "ok"},
		Data: "hello",
	}
	require.NoError(t, resp.JSON(c))
	require.Equal(t, http.StatusOK, rec.Code)

	var body map[string]any
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
	require.Equal(t, "hello", body["data"])
	meta, ok := body["meta"].(map[string]any)
	require.True(t, ok)
	require.Equal(t, constants.Success, meta["error_code"])
}
