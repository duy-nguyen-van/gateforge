package middlewares

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gateforge-iam/gateforge-iam/internal/auth"
	"github.com/gateforge-iam/gateforge-iam/internal/constants"
	"github.com/gateforge-iam/gateforge-iam/internal/dtos"
	"github.com/gateforge-iam/gateforge-iam/internal/models"
	"github.com/gateforge-iam/gateforge-iam/internal/repositories"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

type stubUserRepo struct {
	user *models.User
	err  error
}

func (s *stubUserRepo) CreateWithPasswordHash(_ context.Context, _ *models.User, _ string) error {
	return nil
}
func (s *stubUserRepo) CreateUserOnly(_ context.Context, _ *models.User) error { return nil }
func (s *stubUserRepo) GetOneByID(_ context.Context, _ string) (*models.User, error) {
	return s.user, s.err
}
func (s *stubUserRepo) GetByEmailLower(_ context.Context, _ string) (*models.User, error) {
	return nil, nil
}
func (s *stubUserRepo) Count(_ context.Context) (int64, error)               { return 0, nil }
func (s *stubUserRepo) CountPlatformAdmins(_ context.Context) (int64, error) { return 0, nil }
func (s *stubUserRepo) SetPlatformAdmin(_ context.Context, _ string, _ bool) error {
	return nil
}
func (s *stubUserRepo) UpdateStatus(_ context.Context, _ string, _ constants.UserStatus) error {
	return nil
}
func (s *stubUserRepo) UpdateProfile(_ context.Context, _ string, _ repositories.UserProfilePatch) (*models.User, error) {
	return nil, nil
}
func (s *stubUserRepo) List(_ context.Context, _, _ string, _ *dtos.PageableRequest) (*dtos.DataResponse[models.User], error) {
	return nil, nil
}

func TestPlatformAdminAuth(t *testing.T) {
	okHandler := func(c echo.Context) error {
		return c.String(http.StatusOK, "ok")
	}

	t.Run("missing user id returns 401", func(t *testing.T) {
		e := echo.New()
		e.Use(PlatformAdminAuth(&stubUserRepo{}))
		e.GET("/admin", okHandler)

		req := httptest.NewRequest(http.MethodGet, "/admin", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		require.Equal(t, http.StatusUnauthorized, rec.Code)
		var body map[string]string
		require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &body))
		require.Equal(t, "User not authenticated", body["error"])
	})

	t.Run("repo error returns 403", func(t *testing.T) {
		e := echo.New()
		e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				c.Set(auth.EchoContextUserIDKey, "user-1")
				return next(c)
			}
		})
		e.Use(PlatformAdminAuth(&stubUserRepo{err: errors.New("db down")}))
		e.GET("/admin", okHandler)

		req := httptest.NewRequest(http.MethodGet, "/admin", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		require.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("non-admin user returns 403", func(t *testing.T) {
		e := echo.New()
		e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				c.Set(auth.EchoContextUserIDKey, "user-1")
				return next(c)
			}
		})
		e.Use(PlatformAdminAuth(&stubUserRepo{user: &models.User{IsPlatformAdmin: false}}))
		e.GET("/admin", okHandler)

		req := httptest.NewRequest(http.MethodGet, "/admin", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		require.Equal(t, http.StatusForbidden, rec.Code)
	})

	t.Run("platform admin passes through", func(t *testing.T) {
		e := echo.New()
		e.Use(func(next echo.HandlerFunc) echo.HandlerFunc {
			return func(c echo.Context) error {
				c.Set(auth.EchoContextUserIDKey, "admin-1")
				return next(c)
			}
		})
		e.Use(PlatformAdminAuth(&stubUserRepo{user: &models.User{IsPlatformAdmin: true}}))
		e.GET("/admin", okHandler)

		req := httptest.NewRequest(http.MethodGet, "/admin", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)

		require.Equal(t, http.StatusOK, rec.Code)
		require.Equal(t, "ok", rec.Body.String())
	})
}
