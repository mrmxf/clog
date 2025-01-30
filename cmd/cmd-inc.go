//  Copyright ©2019-2025  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License  https://opensource.org/license/bsd-3-clause/
//
// command:   $ clog Inc

package cmd

import (
	"embed"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/mrmxf/clog/crayon"
	"github.com/mrmxf/clog/slogger"
	"github.com/spf13/cobra"
)

// IncFs is the file system for the embedded shell file - override at runtime to use your own
//go:embed inc.sh
var IncFs embed.FS

var DarkMode bool = false
var JustCrayon bool = false

var CmdInc = &cobra.Command{
	Use:   "Inc",
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
	slogger.GetLogger().Debug("init " + file)
	CmdInc.PersistentFlags().BoolVarP(&DarkMode, "darkmode", "d", false, "all colors for darkmode")
	CmdInc.PersistentFlags().BoolVarP(&JustCrayon, "justcrayon", "j", false, "return color helpers, not script helpers")
}
