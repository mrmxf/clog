//  Copyright ©2017-2025  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License  https://opensource.org/license/bsd-3-clause/
//
// command:
//   $ clog Core Cp core/file/path local/file/path

package copy

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"

	"github.com/mrmxf/clog/config"
	"github.com/spf13/cobra"
)

var Command = &cobra.Command{
	Use:   "Copy",
	Short: "copy an embedded file to a file system",
	Long:  `returns error status 1 if file not found.`,

	Run: func(cmd *cobra.Command, args []string) {
		fs := config.CoreFs()
		srcPath := args[0]
		dstPath := args[1]

		src, err := fs.Open(srcPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(126)
		}
		defer src.Close()

		dst, err := os.Create(dstPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(126)
		}
		defer dst.Close()

		nBytes, err := io.Copy(dst, src)
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s (%d bytes copied)\n", err, nBytes)
			os.Exit(126)
		}
	},
}

func init() {
	_, file, _, _ := runtime.Caller(0)
	slog.Debug("init " + file)
}
