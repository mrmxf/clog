// Copyright ©2017-2025 Mr MXF   info@mrmxf.com
// BSD-3-Clause License   https://opensource.org/license/bsd-3-clause/

package gommi

import (
	"embed"
	"io/fs"
	"log/slog"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/go-chi/chi/v5"
)

func setContentType(w http.ResponseWriter, r *http.Request) {
	ext := filepath.Ext(r.RequestURI)
	switch ext {
	case ".css":
		w.Header().Set("Content-Type", "text/css")
	case ".eot":
		w.Header().Set("Content-Type", "application/vnd.ms-fontobject")
	case ".html":
		w.Header().Set("Content-Type", "text/html")
	case ".ico":
		w.Header().Set("Content-Type", "image/x-icon")
	case ".js":
		w.Header().Set("Content-Type", "text/javascript")
	case ".json":
		w.Header().Set("Content-Type", "application/json")
	case ".otf":
		w.Header().Set("Content-Type", "font/otf")
	case ".png":
		w.Header().Set("Content-Type", "image/png")
	case ".svg":
		w.Header().Set("Content-Type", "image/svg+xml")
	case ".ttf":
		w.Header().Set("Content-Type", "font/ttf")
	case ".txt":
		w.Header().Set("Content-Type", "text/plain")
	case ".wasm":
		w.Header().Set("Content-Type", "application/wasm")
	case ".woff":
		w.Header().Set("Content-Type", "font/woff")
	case ".woff2":
		w.Header().Set("Content-Type", "font/woff2")
	case ".xml":
		w.Header().Set("Content-Type", "application/xml")
	default:
		w.Header().Set("Content-Type", "text/plain")
	}
}

// FileServerFs sets up an http.FileServerFs handler to serve
// static files from a http.FileSystem.
func FileServerFs(r chi.Router, eFs embed.FS, route string, eFsRootPath string) error {
	if strings.ContainsAny(route, "{}*") {
		slog.Error("FileServerFs route does not permit any URL parameters.")
		return
	}

	fSys, err := fs.Sub(eFs, eFsRootPath)
	if err != nil {
		slog.Error("FileServerFs cannot find embedded files")
		return
	}

	// check for trailing slash
	if route != "/" && route[len(route)-1] != '/' {
		r.Get(route, http.RedirectHandler(route+"/", http.StatusMovedPermanently).ServeHTTP)
		route += "/"
	}
	route += "*"

	r.Get(route,
		func(w http.ResponseWriter, r *http.Request) {
			setContentType(w, r)
			rCtx := chi.RouteContext(r.Context())
			pathPrefix := strings.TrimSuffix(rCtx.RoutePattern(), "/*")
			fs := http.StripPrefix(pathPrefix, http.FileServer(http.FS(fSys)))
			fs.ServeHTTP(w, r)
		})
}
