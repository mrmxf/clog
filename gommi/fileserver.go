// Copyright ©2017-2025 Mr MXF   info@mrmxf.com
// BSD-3-Clause License   https://opensource.org/license/bsd-3-clause/

package gommi

import (
	"embed"
	"errors"
	"io/fs"
	"log/slog"
	"net/http"
	"os"
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

// NewFileServer sets up an http.FileServerFs handler to serve
// static files from a http.FileSystem mounted on the os filesystem.
func NewFileServer(r chi.Router, prefix string, mountPath string) error {
	// ensure that we can mount the folder
	abs, err := filepath.Abs(mountPath)
	if err != nil {
		msg := "gommi.NewEmbedFileServer cannot find mountPath"
		slog.Error(msg, "mountPath", mountPath)
		return errors.New(msg)
	}
	// make a new fs at the mount point
	webFs := os.DirFS(abs)
	return fileServerFs(r, webFs, prefix, mountPath)
}

// NewEmbedFileServer sets up an http.FileServerFs handler to serve
// static files from a http.FileSystem mounted on and embed.FS
func NewEmbedFileServer(r chi.Router, embedFs embed.FS, prefix string, mountPath string) error {
	// make a new fs at the mount point
	fs, err := fs.Sub(embedFs, mountPath)
	if err != nil {
		msg := "gommi.NewEmbedFileServer cannot find mountPath"
		slog.Error(msg, "mountPath", mountPath)
		return errors.New(msg)
	}

	// ensure that we can mount the folder
	return fileServerFs(r, fs, prefix, mountPath)
}

// FileServerFs sets up an http.FileServerFs handler to serve
// static files from a http.FileSystem.
func fileServerFs(r chi.Router, webFs fs.FS, route string, mountPath string) error {
	if strings.ContainsAny(route, "{}*") {
		msg := "gommi.FileServerFs route does not permit any URL parameters"
		slog.Error(msg)
		return errors.New(msg)
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
			fs := http.StripPrefix(pathPrefix, http.FileServer(http.FS(webFs)))
			fs.ServeHTTP(w, r)
		})
}
