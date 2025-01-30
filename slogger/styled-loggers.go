// Copyright ©2022-2024 Mr MXF   info@mrmxf.com
// BSD-3-Clause License   https://opensource.org/license/bsd-3-clause/

package slogger

// package log defines the logger for the app

import (
	"bufio"
	"log/slog"
	"os"
	"runtime"

	"github.com/phsym/console-slog"
)

var Logger *slog.Logger

func UsePrettyDebugLogger() {
	Logger = slog.New(
		console.NewHandler(os.Stderr,
			&console.HandlerOptions{Level: slog.LevelDebug}))
	slog.SetDefault(Logger)
}

func UsePrettyInfoLogger() {
	Logger = slog.New(
		console.NewHandler(os.Stderr,
			&console.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(Logger)
}

// JobLogger is a no-color version of the PrettyInfoLogger that is created
// to append to the job log folder
func JobLogger(path string) (*slog.Logger, *os.File) {
	fileHandle, _ := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	writer := bufio.NewWriter(fileHandle)

	newLogger := slog.New(
		console.NewHandler(writer,
			&console.HandlerOptions{Level: slog.LevelInfo, NoColor: true}))

	return newLogger, fileHandle
}

func UseJSONInfoLogger() {
	Logger = slog.New(slog.NewJSONHandler(os.Stderr,
		&slog.HandlerOptions{Level: slog.LevelInfo}))
	slog.SetDefault(Logger)
}

func UseProductionJSONErrorLogger() {
	Logger = slog.New(slog.NewJSONHandler(os.Stderr,
		&slog.HandlerOptions{Level: slog.LevelError}))
	slog.SetDefault(Logger)
}

func init() {
	// uncomment this line to see init order
	// UsePrettyInfoLogger()

	// trace init order for sanity
	_, file, _, _ := runtime.Caller(0)
	slog.Debug("init " + file)
}
