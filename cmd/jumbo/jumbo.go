//  Copyright ©2017-2025  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License  https://opensource.org/license/bsd-3-clause/
//
// command:
//   $ clog Jumbo "Make this text huge" -font=default

package jumbo

import (
	"fmt"
	"log/slog"
	"runtime"
	"slices"
	"strings"

	"github.com/common-nighthawk/go-figure"
	"github.com/mrmxf/clog/config"
	"github.com/spf13/cobra"
)

const defaultFont = "standard"

// various flags for the command
var listFonts bool
var showFonts bool
var bareString bool
var fontOverride string

// Command define the cobra settings for this command
var Command = &cobra.Command{
	Use:   "Jumbo \"Some Text\" --font=thingy",
	Short: "Create large comment text for scripts",
	Long:  `Fonts can be chosen in [clogrc/clog.config.yaml]`,
	Run: func(cmd *cobra.Command, args []string) {
		jString := strings.Join(args, " ")
		if len(jString) == 0 {
			jString = config.Cfg().GetString("clog.jumbo.sample")
		}

		if listFonts {
			fmt.Println("Available fonts:")
			for _, font := range jumboFontList {
				cmd.Println("  " + font)
			}
			return
		}

		if showFonts {
			fmt.Println("Showing available fonts:")
			for _, font := range jumboFontList {
				fmt.Println("  Font: " + font)
				fig := figure.NewFigure(jString+" ("+font+")", font, false)
				fmt.Println(fig.String())
			}
			return
		}

		font := config.Cfg().GetString("clog.jumbo.font")
		if len(fontOverride) > 0 {
			font = fontOverride
		}

		// check that font is in list
		if !slices.Contains(jumboFontList, font) {
			// bad juju - use a standard font and log an error
			slog.Warn("font " + font + " not found, using " + defaultFont)
			font = defaultFont
		}

		myFigure := figure.NewFigure(strings.Join(args, " "), font, false)
		stringSlice := myFigure.Slicify()

		for s := range stringSlice {
			if bareString {
				fmt.Println(stringSlice[s])
			} else {
				fmt.Println("# " + stringSlice[s])
			}
		}
	},
}

func init() {
	_, file, _, _ := runtime.Caller(0)
	slog.Debug("init " + file)

	Command.PersistentFlags().BoolVarP(&listFonts, "list", "L", false, "clog Jumbo -l       # list available fonts")
	Command.PersistentFlags().BoolVarP(&showFonts, "show", "S", false, "clog Jumbo -s       # show available fonts")
	Command.PersistentFlags().BoolVarP(&bareString, "bare", "B", false, "clog Jumbo -b \"Text without leading #\"")
	Command.PersistentFlags().StringVarP(&fontOverride, "font", "F", "", "clog Jumbo -f rounded \"text in rounded font\"")
}
