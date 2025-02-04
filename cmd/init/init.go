//  Copyright ©2019-2024  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License  https://opensource.org/license/bsd-3-clause/
//
// command:
//   $ clog Init

package init

import (
	"log/slog"
	"os"
	"runtime"

	"github.com/mrmxf/clog/cmd/copy"
	"github.com/spf13/cobra"
)

// Command define the cobra settings for this command
var Command = &cobra.Command{
	Use:   "Init",
	Short: "create clogrc/clog.config.yaml if missing",
	Long:  `create missing clogrc/clog.config.yaml and clogrc/tmp-clog.core.config.yaml.`,

	Run: func(cmd *cobra.Command, args []string) {
		sample := "core/sample/clog.config.yaml"
		core := "core/clog.core.config.yaml"
		dstFolder := "clogrc"
		_, err := os.Stat(dstFolder)
		if err != nil {
			//make the folder if it doesn't exist
			slog.Warn("creating folder " + dstFolder)
			if err = os.MkdirAll(dstFolder, 0755); err != nil {
				slog.Error("failed to create folder " + dstFolder)
				os.Exit(1)
			}
		}
		dst := dstFolder + "/clog.config.yaml"
		_, err = os.Stat(dst)
		if err != nil {
			// no config exists - try a copy and ignore any error
			copy.Command.Run(cmd, []string{sample, dst})
		}
		dst = "clogrc/tmp-clog.core.config.yaml"
		copy.Command.Run(cmd, []string{core, dst})
	},
}

func init() {
	_, file, _, _ := runtime.Caller(0)
	slog.Debug("init " + file)
}
