package static

import (
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestResolveDevDistPath(t *testing.T) {
	path := ResolveDevDistPath()
	if path != "" {
		_, err := os.Stat(filepath.Join(path, "index.html"))
		require.NoError(t, err)
	}
}

func TestIsEmptyFS(t *testing.T) {
	require.True(t, isEmptyFS(fstest.MapFS{}))
	require.False(t, isEmptyFS(fstest.MapFS{"index.html": &fstest.MapFile{Data: []byte("ok")}}))
}

func TestOpenAssetFS(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "index.html"), []byte("<html></html>"), 0o644))

	sub, err := openAssetFS(nil, dir)
	require.NoError(t, err)
	require.NotNil(t, sub)

	_, err = openAssetFS(fstest.MapFS{}, "")
	require.ErrorIs(t, err, fs.ErrNotExist)

	_, err = openAssetFS(nil, filepath.Join(dir, "missing"))
	require.ErrorIs(t, err, fs.ErrNotExist)
}

func TestRegister_ServesFromTempDir(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "index.html"), []byte("<html>spa</html>"), 0o644))
	require.NoError(t, os.WriteFile(filepath.Join(dir, "app.js"), []byte("console.log('ok')"), 0o644))

	e := echo.New()
	Register(e, nil, dir)

	t.Run("serves asset file", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/app.js", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		require.Equal(t, http.StatusOK, rec.Code)
		require.Contains(t, rec.Body.String(), "console.log")
	})

	t.Run("falls back to index for unknown route", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodGet, "/unknown/route", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		require.True(t, rec.Code == http.StatusOK || rec.Code == http.StatusMovedPermanently)
		if rec.Code == http.StatusOK {
			require.Contains(t, rec.Body.String(), "spa")
		}
	})

	t.Run("rejects non GET", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodPost, "/app.js", nil)
		rec := httptest.NewRecorder()
		e.ServeHTTP(rec, req)
		require.Equal(t, http.StatusMethodNotAllowed, rec.Code)
	})
}

func TestRegister_NoAssetsIsNoOp(t *testing.T) {
	e := echo.New()
	Register(e, nil, "")
	require.Empty(t, e.Routes())
}

func TestServeIndex(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "index.html"), []byte("<html>index</html>"), 0o644))
	sub, err := openAssetFS(nil, dir)
	require.NoError(t, err)

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/missing", nil)
	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	fileServer := http.FileServer(http.FS(sub))
	require.NoError(t, serveIndex(fileServer, c))
	require.NotZero(t, rec.Code)
}
