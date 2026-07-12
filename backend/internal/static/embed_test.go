package static

import (
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"testing/fstest"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/require"
)

func TestRegister_ServesFromEmbeddedFS(t *testing.T) {
	content := fstest.MapFS{
		"dist/index.html": &fstest.MapFile{Data: []byte("<html>embedded</html>")},
		"dist/app.js":     &fstest.MapFile{Data: []byte("console.log('embedded')")},
	}

	e := echo.New()
	Register(e, content, "")

	req := httptest.NewRequest(http.MethodGet, "/app.js", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "embedded")
}

func TestRegister_ServesRootAsIndex(t *testing.T) {
	content := fstest.MapFS{
		"dist/index.html": &fstest.MapFile{Data: []byte("<html>home</html>")},
	}

	e := echo.New()
	Register(e, content, "")

	req := httptest.NewRequest(http.MethodGet, "/", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.Equal(t, http.StatusOK, rec.Code)
	require.Contains(t, rec.Body.String(), "home")
}

func TestRegister_ServesDirectoryAsIndexFallback(t *testing.T) {
	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "index.html"), []byte("<html>dir</html>"), 0o644))
	require.NoError(t, os.Mkdir(filepath.Join(dir, "assets"), 0o755))

	e := echo.New()
	Register(e, nil, dir)

	req := httptest.NewRequest(http.MethodGet, "/assets/", nil)
	rec := httptest.NewRecorder()
	e.ServeHTTP(rec, req)
	require.True(t, rec.Code == http.StatusOK || rec.Code == http.StatusMovedPermanently)
}

func TestOpenAssetFS_PrefersEmbeddedWhenPresent(t *testing.T) {
	content := fstest.MapFS{
		"dist/index.html": &fstest.MapFile{Data: []byte("embedded")},
	}
	sub, err := openAssetFS(content, t.TempDir())
	require.NoError(t, err)
	f, err := sub.Open("index.html")
	require.NoError(t, err)
	defer f.Close()
	buf := make([]byte, 16)
	n, err := f.Read(buf)
	require.NoError(t, err)
	require.Equal(t, "embedded", string(buf[:n]))
}
