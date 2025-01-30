//  Copyright ©2019-2024  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License  https://opensource.org/license/bsd-3-clause/
//
// Package slogger wraps the slog library for consistent logging in clog

package slogger

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"log/slog"
	"runtime"

	"os"
)

var theLogger *slog.Logger = nil
var lvl = new(slog.LevelVar)
var crayon = ttycrayon.Color()

// GetLogger returns the default logger. If the logger has not been initialised
// return the default pretty logger with a default info logging level.
// to see the init order of packages - use a default order of slog.LevelDebug
// note that there is no programmatic way of turning that off. slogger() is
// initialised before everything else.
func GetLogger() *slog.Logger {
	if theLogger == nil {
		lvl.Set(slog.LevelInfo) // use slog.LevelDebug or LevelInfo
		opts := PrettyHandlerOptions{SlogOpts: slog.HandlerOptions{Level: lvl}}
		theLogger = slog.New(NewPrettyHandler(os.Stderr, opts))
	}

	return theLogger
}

// take a sting from config and convert it to slog.Level
func stringToLevel(level string) slog.Level {
	var l slog.Level
	switch level {
	case "DEBUG":
		l = slog.LevelDebug
	case "INFO":
		l = slog.LevelInfo
	case "WARN":
		l = slog.LevelWarn
	case "ERROR":
		l = slog.LevelError
	default:
		l = slog.LevelInfo
	}
	return l
}

// SetLevel sets the log level for the default logger.
func SetLevel(newLevel string) {
	lvl.Set(stringToLevel(newLevel))
}

type PrettyHandlerOptions struct {
	SlogOpts slog.HandlerOptions
}

type PrettyHandler struct {
	slog.Handler
	l *log.Logger
}

func (h *PrettyHandler) Handle(ctx context.Context, r slog.Record) error {
	level := r.Level.String() + ":"

	switch r.Level {
	case slog.LevelDebug:
		level = crayon.Dbg(level)
	case slog.LevelInfo:
		level = crayon.Info(level)
	case slog.LevelWarn:
		level = crayon.Warning(level)
	case slog.LevelError:
		level = crayon.Error(level)
	}

	fields := make(map[string]interface{}, r.NumAttrs())
	r.Attrs(func(a slog.Attr) bool {
		fields[a.Key] = a.Value.Any()

		return true
	})

	b, err := json.MarshalIndent(fields, "", "  ")
	if err != nil {
		return err
	}

	timeStr := r.Time.Format("[15:05:05.000]")
	msg := crayon.Command(r.Message)

	h.l.Println(timeStr, level, msg, crayon.Text(string(b)))

	return nil
}

func NewPrettyHandler(
	out io.Writer,
	opts PrettyHandlerOptions,
) *PrettyHandler {
	h := &PrettyHandler{
		Handler: slog.NewJSONHandler(out, &opts.SlogOpts),
		l:       log.New(out, "", 0),
	}

	return h
}

// all init functions log at Debug level that they are running init
// this means that slogger.init() is always the first to run

func init() {
	_, file, _, _ := runtime.Caller(0)
	GetLogger().Debug("init " + file)
}
