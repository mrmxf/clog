//  Copyright ©2017-2025  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License  https://opensource.org/license/bsd-3-clause/
//
// package log adds a log command to the clog command line tool

package logcmd

import (
	"fmt"
	"os"
	"runtime"
	"strings"

	slog "github.com/mrmxf/clog/slogger"

	"github.com/spf13/cobra"
)

var build bool
var debug bool
var emergency bool
var error bool
var fatal bool
var info bool
var success bool
var trace bool
var warn bool
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
	clog Log -B "$errCount" "$isProduction" "Base-Message"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		logMsg := strings.Join(args, " ")
		// most serious flag wins
		logFlag := "none"

		if up {
			//up one line, start of line, del EOL
			fmt.Print("\x1b[A\x1b[G\x1b[K")
		}
		// if user has many falgs, then the top-most case statement wins
		switch {
		case emergency:
			slog.Emergency(logMsg)
			logFlag = "X"
		case fatal:
			slog.Fatal(logMsg)
			logFlag = "F"
		case error:
			slog.Error(logMsg)
			logFlag = "E"
		case warn:
			slog.Warn(logMsg)
			logFlag = "W"
		case success:
			slog.Success(logMsg)
			logFlag = "S"
		case info:
			slog.Info(logMsg)
			logFlag = "I"
		case trace:
			slog.Trace(logMsg)
			logFlag = "T"
		case debug:
			slog.Debug(logMsg)
			logFlag = "D"
		case build:
			// a special case.
			// -  err + dev  builds: Warn and carry on (i.e. exit code 0)
			// -  err + prod builds: Fail and exit (i.e. exit code 1)
			// -  ok  + either     : print out success
			var exitCode int
			_, err := fmt.Sscanf(args[0], "%d", &exitCode)
			if len(args) < 4 || err != nil {
				slog.Error("clog Log -B requires 4 arguments")
				slog.Error("   arg[0]           \"$?\" exit code of command to log for")
				slog.Error("   arg[1]      \"$doPROD\" empty string for tolerant dev mode otherwise any string for fragile PROD mode")
				slog.Error("   arg[2]   \"OK message\" string to be logged for $?=0")
				slog.Error("   arg[3]  \"Err Message\" string to be logged for $?=0")
				os.Exit(1)
			}
			if exitCode == 0 {
				slog.Success(args[2])
			} else if len(args[1]) > 0 {
				// fragile production mode
				slog.Error(args[3])
				os.Exit(1)
			} else {
				// fragile production mode
				slog.Warn(args[3])
			}
			logFlag = "B"
		}
		// level, levelFile := slogger.GetLogLevel()
		slog.Debug("Log (-%s) (%s)", logFlag, logMsg)
	},
}

func init() {
	_, file, _, _ := runtime.Caller(0)
	slog.Debug("init " + file)

	Command.PersistentFlags().BoolVarP(&info, "info", "I", false, "clog Log -I \"Info message\"")
	Command.PersistentFlags().BoolVarP(&success, "success", "S", false, "clog Log -S \"Success message\"")
	Command.PersistentFlags().BoolVarP(&warn, "warn", "W", false, "clog Log -W \"Warn message\"")
	Command.PersistentFlags().BoolVarP(&error, "error", "E", false, "clog Log -E \"Error message\"")
	Command.PersistentFlags().BoolVarP(&trace, "trace", "T", false, "clog Log -T \"Trace message\"")
	Command.PersistentFlags().BoolVarP(&debug, "debug", "D", false, "clog Log -D \"Debug message\"")
	Command.PersistentFlags().BoolVarP(&fatal, "fatal", "F", false, "clog Log -E \"Fatal message\"")
	Command.PersistentFlags().BoolVarP(&emergency, "emergency", "X", false, "clog Log -X \"Emergency message\"")
	Command.PersistentFlags().BoolVarP(&build, "build", "B", false, "clog Log -B \"$?\" \"$emptyForDev\" \"okMsg $?=0\" \"errMsg $?>0\"")
	Command.PersistentFlags().BoolVarP(&up, "up", "U", false, "clog Log -UI \"up (overprint) Info message\"")
}
