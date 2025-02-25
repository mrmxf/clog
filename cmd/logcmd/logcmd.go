//  Copyright ©2017-2025  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License  https://opensource.org/license/bsd-3-clause/
//
// package log adds a log command to the clog command line tool

package logcmd

import (
	"log/slog"
	"runtime"
	"strings"

	"github.com/spf13/cobra"
)

var info bool
var error bool
var warn bool
var success bool
var debug bool

// Command define the cobra settings for this command
var Command = &cobra.Command{
	Use:   "Log",
	Short: "log a message to the configured logger",
	Example: `
	clog Log -I "info message"
	clog Log -S "success message"
	clog Log -E "error message"
	clog Log -W "warning message"
	clog Log -D "debug message"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		logMsg := strings.Join(args, " ")
		// most serious flag wins
		logFlag := "none"
		if error {
			slog.Error(logMsg)
			logFlag = "E"
		}
		if warn {
			slog.Warn(logMsg)
			logFlag = "W"
		}
		if info {
			slog.Info(logMsg)
			logFlag = "I"
		}
		if debug {
			slog.Debug(logMsg)
			logFlag = "D"
		}
		// level, levelFile := slogger.GetLogLevel()
		slog.Debug("Log (-%s) (%s)", logFlag, logMsg)
	},
}

func init() {
	_, file, _, _ := runtime.Caller(0)
	slog.Debug("init " + file)

	Command.PersistentFlags().BoolVarP(&info, "info", "I", false, "clog Log -I \"Info message\"")
	Command.PersistentFlags().BoolVarP(&error, "error", "E", false, "clog Log -E \"Error message\"")
	Command.PersistentFlags().BoolVarP(&warn, "warn", "W", false, "clog Log -W \"Warn message\"")
	Command.PersistentFlags().BoolVarP(&debug, "debug", "D", false, "clog Log -D \"Debug message\"")
}
