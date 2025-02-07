//  Copyright ©2017-2025  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License           https://opensource.org/license/bsd-3-clause/

package version

import (
	"fmt"
	"log/slog"
	"os"
	"runtime"

	"github.com/mrmxf/clog/config"
	"github.com/mrmxf/clog/crayon"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "version",
	Short: "Print the version string",
	Long:  `use clog -v to display the short semantic version string.`,
	Run: func(parent *cobra.Command, args []string) {
		_, file, _, _ := runtime.Caller(0)
		slog.Debug(fmt.Sprintf("run command: %s", file))

		if len(args) == 0 {
			fmt.Printf("%s (%s) %s\n",
				crayon.ColorCapitals(config.Cfg().GetString("title"), nil, nil),
				config.Cfg().GetString("app"),
				config.Cfg().GetString("clog.version.long"))
			os.Exit(0)
		}
		if args[0] == "short" {
			fmt.Println(config.Cfg().GetString("ver"))
			os.Exit(0)
		}
		if args[0] == "note" {
			fmt.Println(config.Cfg().GetString("clog.version.note"))
			os.Exit(0)
		}
		slog.Error("unknown version argument (" + args[0] + ")")
		os.Exit(1)
	},
}

func init() {
	_, file, _, _ := runtime.Caller(0)
	slog.Debug("init " + file)
}
