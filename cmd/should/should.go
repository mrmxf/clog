//  Copyright ©2017-2025  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License  https://opensource.org/license/bsd-3-clause/
//
// package log adds a log command to the clog command line tool

package should

import (
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"

	"github.com/mrmxf/clog/slogger"
	"github.com/spf13/cobra"
)

var debug bool

// Command define the cobra settings for this command
var Command = &cobra.Command{
	SilenceErrors: true,
	SilenceUsage:  true,
	Use:           "Should",
	Short:         "check env variable contains a space delimited word",
	Example: `
	export MAKE="data hugo exe ko"
	clog Should MAKE "hugo"    # $?==0 if $MAKE contains the word "hugo"
	clog Should NUKE "folders" # $?==0 if $NUKE contains the word "folders"
	`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) < 2 {
			cmd.Help()
			os.Exit(1)
		}

		if debug {
			slogger.UsePrettyLogger(slog.LevelDebug)
		}
		needle := args[1]
		haystackEnv := args[0]
		haystack, envExists := os.LookupEnv(haystackEnv)
		if envExists {
			words := strings.Split(haystack, " ")
			found := false
			for i := range words {
				if words[i] == needle {
					found = true
				}
			}
			if found {
				dbg := fmt.Sprintf("clog Should $%s \"%s\" - ✅ found in(%s)", haystackEnv, needle, haystack)
				slog.Debug(dbg)
				os.Exit(0)
			} else {
				dbg := fmt.Sprintf("clog Should $%s \"%s\" ❌ missing from(%s)", haystackEnv, needle, haystack)
				slog.Debug(dbg)
				os.Exit(1)
			}
		} else {
			dbg := fmt.Sprintf("clog Should $%s \"%s\" - ❌ missing env %s", haystackEnv, needle, haystackEnv)
			slog.Debug(dbg)
			os.Exit(2)
		}
	},
}

func init() {
	_, file, _, _ := runtime.Caller(0)
	slog.Debug("init " + file)

	Command.PersistentFlags().BoolVarP(&debug, "debug", "D", false, "clog Should -D $MAKE \"thingy\"")
}
