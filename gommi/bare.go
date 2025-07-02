// Copyright ©2017-2025 Mr MXF   info@mrmxf.com
// BSD-3-Clause License   https://opensource.org/license/bsd-3-clause/

// gommi  - Golang Mimimal Modular InterWeb
//
// simple static web server tools that will serve a static site from inside
// a container.
//
// 	    r, _ := gommi.Bare()
//    	r.NewEmbedFileServer(eFS, "/", "embedWWW/")
//     	http.ListenAndServe("0.0.0.0:8080", r)
//
// Usually pages are served from:
//  - templates in the embedded file system for basic server side rendering
//  - the embedded file system for non-templates
//  - the containing pod's `/var/www` if it exists
//  - 404 otherwise

package gommi

import (
	"io/fs"
	"log/slog"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	slogchi "github.com/samber/slog-chi"
)

type ChiMux struct {
	*chi.Mux       // the chi Mux that we're extending
	webFs    fs.FS // the file system that will be used to serve pages
}

// the top level mux
var mux *ChiMux

// Bare is a a bare mux with no routes - just a slog logger & recoverer
func Bare(opts ...Options) (*ChiMux, error) {
	processOptions(opts)

	slog.SetDefault(opt.Logger)
	mux = &ChiMux{chi.NewRouter(), nil}

	// Middleware
	mux.Use(slogchi.New(opt.Logger))
	mux.Use(middleware.Recoverer)

	return mux, nil
}

// get the slog logger in use in gommi.Bare
func GetLogger() *slog.Logger {
	return opt.Logger
}

// get the slog logger in use in gommi.Bare
func GetMux() *ChiMux {
	return mux
}
