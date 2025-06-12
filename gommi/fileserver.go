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

var urlPrefix string

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
		w.Header().Set("Content-Type", "text/html")
	}
}

// NewFileServer sets up an http.FileServerFs handler to serve
// static files from a http.FileSystem mounted on the os filesystem.
// When building a container with ko.build that will be mounted
func (r *ChiMux) NewFileServer(prefix string, mountPath string) error {
	// ensure that we can mount the folder
	abs, err := filepath.Abs(mountPath)
	if err != nil {
		msg := "gommi.NewEmbedFileServer cannot find mountPath"
		slog.Error(msg, "mountPath", mountPath)
		if abortOnError {
			os.Exit(1)
		}
		return errors.New(msg)
	}
	// make a new fs at the mount point
	r.webFs = os.DirFS(abs)
	slog.Info("initialising os file server on prefix", "prefix", prefix, "mountPath", abs)
	return r.fileServerFs(prefix)
}

// NewEmbedFileServer sets up an http.FileServerFs handler to serve
// static files from a http.FileSystem mounted on and embed.FS
func (r *ChiMux) NewEmbedFileServer(embedFs embed.FS, prefix string, mountPath string) error {
	// make a new fs at the mount point
	fs, err := fs.Sub(embedFs, mountPath)
	if err != nil {
		msg := "gommi.NewEmbedFileServer cannot find mountPath"
		slog.Error(msg, "mountPath", mountPath)
		return errors.New(msg)
	}
	r.webFs = fs
	slog.Info("initialising embed file server on prefix", "prefix", prefix, "mountPath", mountPath)
	// ensure that we can mount the folder
	return r.fileServerFs(prefix)
}

// fileServerFs sets up an http.FileServerFs handler to serve
// static files from a http.FileSystem.
func (r *ChiMux) fileServerFs(prefix string) error {
	if strings.ContainsAny(prefix, "{}*") {
		msg := "gommi.fileServerFs route does not permit any URL parameters"
		slog.Error(msg)
		return errors.New(msg)
	}

	// chi prefix must start with slash
	if !strings.HasPrefix(prefix, "/") {
		prefix = "/" + prefix
	}
	// check for trailing slash on prefix
	switch {
	case prefix == "":
		prefix = "/"
	case strings.HasSuffix(prefix, "/"):
		break
	default:
		prefix += "/"
	}
	// setup a GET handler serving static files for `prefix/*`
	r.Get(prefix+"*", r.serveFile)
	return nil
}

func (r *ChiMux) serveFile(w http.ResponseWriter, req *http.Request) {
	setContentType(w, req)
	rCtx := chi.RouteContext(req.Context())
	pathPrefix := strings.TrimSuffix(rCtx.RoutePattern(), "/*")
	fs := http.StripPrefix(pathPrefix, http.FileServer(http.FS(r.webFs)))
	fs.ServeHTTP(w, req)
}
