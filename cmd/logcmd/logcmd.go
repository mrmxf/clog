//  Copyright ©2017-2025  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License  https://opensource.org/license/bsd-3-clause/
//
// package log adds a log command to the clog command line tool

package logcmd

import (
	"fmt"
	"runtime"
	"strings"

	slog "github.com/mrmxf/clog/slogger"

	"github.com/spf13/cobra"
)

var trace bool
var debug bool
var info bool
var success bool
var warn bool
var error bool
var fatal bool
var emergency bool
var up bool

// Command define the cobra settings for this command
var Command = &cobra.Command{
	Use:   "Log",
	Short: "log a message to the configured logger",
	Example: `
	clog Log -T  "trace message"
	clog Log -D  "debug message"
	clog Log -W  "warning message"
	clog Log -I  "info message"
	clog Log -S  "success message"
	clog Log -E  "error message"
	clog Log -F  "fatal message"
	clog Log -X  "emergency message"
	clog Log -UI "up one line (overprint) an info message"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		logMsg := strings.Join(args, " ")
		// most serious flag wins
		logFlag := "none"
		if up {
			//up one line, start of line, del EOL
			fmt.Print("\x1b[A\x1b[G\x1b[K")
		}

		if trace {
			slog.Trace(logMsg)
			logFlag = "T"
		}
		if debug {
			slog.Debug(logMsg)
			logFlag = "D"
		}
		if info {
			slog.Info(logMsg)
			logFlag = "I"
		}
		if success {
			slog.Success(logMsg)
			logFlag = "S"
		}
		if warn {
			slog.Warn(logMsg)
			logFlag = "W"
		}
		if error {
			slog.Error(logMsg)
			logFlag = "E"
		}
		if fatal {
			slog.Fatal(logMsg)
			logFlag = "F"
		}
		if emergency {
			slog.Emergency(logMsg)
			logFlag = "X"
		}
		// level, levelFile := slogger.GetLogLevel()
		slog.Debug("Log (-%s) (%s)", logFlag, logMsg)
	},
}

func init() {
	_, file, _, _ := runtime.Caller(0)
	slog.Debug("init " + file)

	Command.PersistentFlags().BoolVarP(&trace, "trace", "T", false, "clog Log -T \"Trace message\"")
	Command.PersistentFlags().BoolVarP(&debug, "debug", "D", false, "clog Log -D \"Debug message\"")
	Command.PersistentFlags().BoolVarP(&info, "info", "I", false, "clog Log -I \"Info message\"")
	Command.PersistentFlags().BoolVarP(&success, "success", "S", false, "clog Log -S \"Success message\"")
	Command.PersistentFlags().BoolVarP(&warn, "warn", "W", false, "clog Log -W \"Warn message\"")
	Command.PersistentFlags().BoolVarP(&error, "error", "E", false, "clog Log -E \"Error message\"")
	Command.PersistentFlags().BoolVarP(&fatal, "fatal", "F", false, "clog Log -E \"Fatal message\"")
	Command.PersistentFlags().BoolVarP(&emergency, "emergency", "X", false, "clog Log -X \"Emergency message\"")
	Command.PersistentFlags().BoolVarP(&up, "up", "U", false, "clog Log -UI \"up (overprint) Info message\"")
}
