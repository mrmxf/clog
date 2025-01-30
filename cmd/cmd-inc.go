//  Copyright ©2019-2025  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License  https://opensource.org/license/bsd-3-clause/
//
// command:   $ clog Inc

package cmd

import (
	"embed"
	"runtime"

	"github.com/mrmxf/clog/slogger"
	"github.com/spf13/cobra"
)

// IncFs is the file system for the embedded shell file - override at runtime to use your own
//go:embed inc.sh
var IncFs embed.FS

var IncCmd = incCommand
var incCommand = &cobra.Command{
	Use:   "Inc",
	Short: "cat embedded " + incPath + " to stdout",
	Long:  `returns error status 1 if file not found.`,

	Run: func(cmd *cobra.Command, args []string) {
		catCommand.Run(cmd, []string{incPath})
	},
}

func init() {
	_, file, _, _ := runtime.Caller(0)
	slogger.GetLogger().Debug("init " + file)

}
