//  Copyright ©2017-2025  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License  https://opensource.org/license/bsd-3-clause/

package crayon

// CrayonCmd allows scripts to access the same CLI colors as clog
//
// The default case prints a string for bash & zsh & sh to have colors:
//
//	source <(clog Core Crayon)
//	printf "pretty $cE red error $cI yellow info $cW magenta warning $cX\n"

import (
	"fmt"
	"log/slog"
	"runtime"

	"github.com/mrmxf/clog/crayon"
	"github.com/spf13/cobra"
)

var hideComment bool
var hideShellScript bool
var showExample bool
var DarkMode bool

var Command = &cobra.Command{
	Use:   "Crayon",
	Short: "get bash string for default colors",
	Long:  `Usage: source <(clog Core Crayon)`,

	Run: func(cmd *cobra.Command, args []string) {
		comment := "# Command  | Error     | Info      | File      | Header      | Success   | Text      | Url       | Warning   | eXit     | AMD64         | Arm64       | Linux       | Mac         | Windows\n"
		bashstr := crayon.GetBashString(false)

		if !hideComment {
			fmt.Print(comment)
		}
		if !hideShellScript {
			fmt.Print(bashstr + "\n")
		}
		if showExample {
			fmt.Println(crayon.SampleColors())
		}
	},
}

// set Crayon Flags to provide output control for printing colors.
//
//	clog Core Crayon --help # display help for the flags
func init() {
	Command.PersistentFlags().BoolVarP(&hideShellScript, "comment", "C", false, "print the comment, no Shell script")
	Command.PersistentFlags().BoolVarP(&hideComment, "script", "S", false, "print the Shell script, no comment")
	Command.PersistentFlags().BoolVarP(&showExample, "example", "E", false, "show an example of colors")
	Command.PersistentFlags().BoolVarP(&DarkMode, "darkmode", "D", false, "all colors for darkmode")
	_, file, _, _ := runtime.Caller(0)
	slog.Debug("init " + file)
}
