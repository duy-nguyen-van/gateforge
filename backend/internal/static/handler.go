package static

import (
	"io/fs"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"strings"

	"github.com/labstack/echo/v4"
)

// ResolveDevDistPath returns a frontend dist directory on disk for local debug builds
// (go run / delve without -tags embedfrontend). Checks paths relative to cwd.
func ResolveDevDistPath() string {
	candidates := []string{
		"../../internal/static/dist", // backend/cmd/server
		"internal/static/dist",       // backend/
	}
	for _, candidate := range candidates {
		indexPath := filepath.Join(candidate, "index.html")
		if st, err := os.Stat(indexPath); err == nil && !st.IsDir() {
			return candidate
		}
	}
	return ""
}

// Register mounts the embedded SPA and static assets. Must be called after API/OIDC routes.
// When embed FS is empty, diskPath is used (for local debug without -tags embedfrontend).
func Register(r *echo.Echo, content fs.FS, diskPath string) {
	sub, _ := openAssetFS(content, diskPath)
	if sub == nil {
		return
	}

	fileServer := http.FileServer(http.FS(sub))

	r.GET("/*", func(c echo.Context) error {
		method := c.Request().Method
		if method != http.MethodGet && method != http.MethodHead {
			return echo.ErrMethodNotAllowed
		}

		reqPath := c.Request().URL.Path
		if reqPath == "/" {
			reqPath = "index.html"
		} else {
			reqPath = strings.TrimPrefix(path.Clean(reqPath), "/")
		}

		f, err := sub.Open(reqPath)
		if err != nil {
			return serveIndex(fileServer, c)
		}
		defer f.Close()

		stat, err := f.Stat()
		if err != nil || stat.IsDir() {
			return serveIndex(fileServer, c)
		}

		fileServer.ServeHTTP(c.Response(), c.Request())
		return nil
	})
}

func openAssetFS(content fs.FS, diskPath string) (fs.FS, error) {
	if content != nil {
		if sub, err := fs.Sub(content, "dist"); err == nil && !isEmptyFS(sub) {
			return sub, nil
		}
	}

	if diskPath == "" {
		return nil, fs.ErrNotExist
	}

	indexPath := filepath.Join(diskPath, "index.html")
	if st, err := os.Stat(indexPath); err != nil || st.IsDir() {
		return nil, fs.ErrNotExist
	}

	return os.DirFS(diskPath), nil
}

func isEmptyFS(fsys fs.FS) bool {
	entries, err := fs.ReadDir(fsys, ".")
	return err != nil || len(entries) == 0
}

func serveIndex(fileServer http.Handler, c echo.Context) error {
	r := c.Request().Clone(c.Request().Context())
	r.URL.Path = "/index.html"
	fileServer.ServeHTTP(c.Response(), r)
	return nil
}
