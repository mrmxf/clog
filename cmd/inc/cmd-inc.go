//  Copyright ©2017-2025  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License  https://opensource.org/license/bsd-3-clause/
//
// command:   $ clog Inc

package inc

import (
	"embed"
	"fmt"
	"io"
	"log/slog"
	"os"
	"runtime"

	"github.com/mrmxf/clog/crayon"
	"github.com/spf13/cobra"
)

var CommandString = "Inc"

// IncFs is the file system for the embedded shell file - override at runtime to use your own
//
//go:embed inc.sh
var IncFs embed.FS

var DarkMode bool = false
var JustCrayon bool = false

// Command define the cobra settings for this command
var Command = &cobra.Command{
	Use:   CommandString,
	Short: "send embedded helper script to stdout",
	Long:  `returns error status 126 if embedded file not found.`,

	Run: func(cmd *cobra.Command, args []string) {

		src, err := IncFs.Open("inc.sh")
		if err != nil {
			fmt.Fprintf(os.Stderr, "%s\n", err)
			os.Exit(126)
		}
		defer src.Close()

		// start with the color helpers
		fmt.Println(crayon.GetBashString(DarkMode))

		dst := os.Stdout
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
	Command.PersistentFlags().BoolVarP(&DarkMode, "darkmode", "D", false, "all colors for darkmode")
}
