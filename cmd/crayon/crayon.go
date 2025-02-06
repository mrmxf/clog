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
var hideBashStr bool
var showSample bool

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
		if !hideBashStr {
			fmt.Print(bashstr + "\n")
		}
		if showSample {
			fmt.Println(crayon.SampleColors())
		}
	},
}



// set Crayon Flags to provide output control for printing colors.
//
//	clog Core Crayon --help # display help for the flags
func init() {
	Command.PersistentFlags().BoolVarP(&hideBashStr, "hideBashStr", "b", false, "exclude / hide the bash string default=show")
	Command.PersistentFlags().BoolVarP(&hideComment, "hideComment", "a", false, "hide the annotation of the bash string default=show")
	Command.PersistentFlags().BoolVarP(&showSample,  "showSample", "s", false, "show a sample of colors default=hide")
	_, file, _, _ := runtime.Caller(0)
	slog.Debug("init " + file)
}
