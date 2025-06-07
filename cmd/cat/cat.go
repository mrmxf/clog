// Copyright ©2017-2025  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License  https://opensource.org/license/bsd-3-clause/
//
// command:
//   $ clog Cat

package cat

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"

	"github.com/mrmxf/clog/config"
	"github.com/mrmxf/clog/core"
	"github.com/spf13/cobra"
)

// CatCmd performs a cat of a resource in the clog file system.
//
// Errors go to stderr and an exit value of 1 can be used in scripts.
// Paths to the internal files are the same as clog Update and clog Init. You
// can source a script with the following snippet:
//
// ```
//
//	eval "$(clog Cat tpl/help-golang.sh)"
//
// ```
//
// this has the advantage of auto-updating with every new version of clog. For
// a list of files:
//
// ````
//
//	clog Ls
//
// ```

// Command define the cobra settings for this command
var Command = &cobra.Command{
	Use:   "Cat",
	Short: "copy an internal file to stdout",
	Long:  `returns error status 1 if file not found.`,

	Run: func(cmd *cobra.Command, args []string) {
		clogFs := config.CoreFs()
		srcPath := args[0]

		src, err := clogFs.Open(core.Clean(srcPath))
		if err != nil {
			slog.Error(fmt.Sprintf("cannot open embedded file %s", srcPath), "err", err)
			os.Exit(126)
		}
		defer src.Close()

		dst := os.Stdout

		nBytes, err := io.Copy(dst, src)
		if err != nil {
			slog.Error(fmt.Sprintf("cannot read embedded file %s (%d bytes copied)", srcPath, nBytes), "err", err)
			os.Exit(126)
		}
	},
}

func init() {
	_, file, _, _ := runtime.Caller(0)
	slog.Debug("init " + file)
}
