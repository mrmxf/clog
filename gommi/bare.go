// Copyright ©2017-2025 Mr MXF   info@mrmxf.com
// BSD-3-Clause License   https://opensource.org/license/bsd-3-clause/

// gommi  - Golang Mimimal Modular InterWeb
//
// simple static web server tools that will serve a static site from inside
// a contaienr.
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
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/mrmxf/clog/slogger"
	slogchi "github.com/samber/slog-chi"
)

type ChiMux struct {
	*chi.Mux       // the chi Mux that we're extending
	webFs    fs.FS // the file system that will be used to serve pages
}

// the default logger
var logger *slog.Logger

// the top level mux
var mux *ChiMux

// the port to use
var Port = 8080

var abortOnError = true

// Bare is a a bare mux with no routes - just a slog logger & recoverer
func Bare(continueOnError ...bool) (*ChiMux, error) {
	slogger.UsePrettyLogger(slog.LevelInfo)
	if len(continueOnError) > 0 {
		abortOnError = !continueOnError[0]
	}
	mux = &ChiMux{chi.NewRouter(), nil}
	// Create a slog logger, which:
	//   - Logs to stdout.
	logger = slog.New(slog.NewTextHandler(os.Stdout, nil))

	// Middleware
	mux.Use(slogchi.New(logger))
	mux.Use(middleware.Recoverer)

	return mux, nil
}

// get the slog logger in use in gommi.Bare
func GetLogger() *slog.Logger {
	return logger
}

// get the slog logger in use in gommi.Bare
func GetMux() *ChiMux {
	return mux
}
