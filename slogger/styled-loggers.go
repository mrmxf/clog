// Copyright ©2018-2025 Mr MXF   info@mrmxf.com
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

func UsePrettyLogger(level slog.Level) {
	Logger = slog.New(
		console.NewHandler(os.Stderr,
			&console.HandlerOptions{Level: level}))
	slog.SetDefault(Logger)
}


func UsePlainLogger(level slog.Level) {
	Logger = slog.New(
		console.NewHandler(os.Stderr,
			&console.HandlerOptions{Level: level, NoColor: true}))
	slog.SetDefault(Logger)
}

// JobLogger is a no-color version of the PrettyLogger that is created
// to append to a job log folder. If the file cannot be opened for appending
// an error is returned
func NewJobLogger(path string, level slog.Level) (*slog.Logger, *os.File, error) {
	fileHandle, err := os.OpenFile(path, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	writer := bufio.NewWriter(fileHandle)

	newLogger := slog.New(
		console.NewHandler(writer,
			&console.HandlerOptions{Level: level, NoColor: true}))

	return newLogger, fileHandle, err
}

func UseJSONLogger(level slog.Level) {
	Logger = slog.New(slog.NewJSONHandler(os.Stderr,
		&slog.HandlerOptions{Level: level}))
	slog.SetDefault(Logger)
}

func init() {
	// trace init order for sanity
	_, file, _, _ := runtime.Caller(0)
	slog.Debug("init " + file)
}
