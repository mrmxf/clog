// Copyright ©2017-2025 Mr MXF   info@mrmxf.com
// BSD-3-Clause License   https://opensource.org/license/bsd-3-clause/

package slogger

// package slogger defines a consistent set of styled loggers for clog and
// other apps. It silently initializes to a Pretty Logger with LogInfo logging.
//
// If you want a different default logger then use a different UseXXXLogger
// in your main.init()

import (
	"log/slog"
	"runtime"
)

// the exported default logger
var Logger *slog.Logger

// these are exported with a get function for clog Log
var logLevel slog.Level
var logLevelFile slog.Level

type SlogStyle int

const (
	Plain SlogStyle = iota
	Pretty
	JSON
	JOB
)

// the default logging level for the default logger
var defaultLogLevel = slog.LevelInfo

// var defaultLogLevel = slog.LevelDebug  //use this for init tracing

// use this function to set a log level from a config file
func SetLogger(level slog.Level, style SlogStyle) {

}

// get the active log levels for a split console / cicd logging experience
func GetLogLevel() (logLevel slog.Level, logLevelFile slog.Level) {
	return logLevel, logLevelFile
}

func init() {
	// uncomment this line to see init order
	UsePrettyLogger(defaultLogLevel)

	// trace init order for sanity
	_, file, _, _ := runtime.Caller(0)
	slog.Debug("init " + file)
}
