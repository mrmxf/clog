//  Copyright ©2017-2025  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License  https://opensource.org/license/bsd-3-clause/
//
// package init creates a basic clogrc folder and files if they are missing.
//
// It does 4 things:
// 1. creates a clogrc folder if it doesn't exist
// 2. copies core/sample/clog.config.yaml to clogrc/clog.config.yaml if missing
// 3. copies core/clog.core.config.yaml to clogrc/tmp-clog.core.config.yaml
// 4. copies all files in core/sample/init to clogrc if missing

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

		// 1, create the clogrc folder if it doesn't exist
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

		// 2. create clog.config.yaml if it doesn't exist
		dst := dstFolder + "/clog.config.yaml"
		_, err = os.Stat(dst)
		if err != nil {
			// no config exists - try a copy and ignore any error
			copy.Command.Run(cmd, []string{sample, dst})
		}

		// 3. create tmp-clog.core.config.yaml
		dst = "clogrc/tmp-clog.core.config.yaml"
		copy.Command.Run(cmd, []string{core, dst})
	},
}

func init() {
	_, file, _, _ := runtime.Caller(0)
	slog.Debug("init " + file)
}
