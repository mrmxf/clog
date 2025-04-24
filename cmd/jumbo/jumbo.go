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

// listFonts enumerates available font names
var listFonts bool

// display [text] in every available font. If [text] is not specificed then
// config setting clog.jumbo.sample is used
var showFonts bool

// fontOverride forces a font, otherwise the default from config setting
// clog.jumbo.font is used
var fontOverride string

// commentStyle controls the comment settings of the output word
// see: https://gist.github.com/dk949/88b2652284234f723decaeb84db2576c
// available comment styles with clog help Jumbo
var commentStyle string

// Command define the cobra settings for this command
var Command = &cobra.Command{
	Use:   "Jumbo [text] --font=thingy",
	Short: "Create large comment text for scripts to stdout",
	Long:  `Fonts can be chosen in [clogrc/clog.config.yaml]`,
	Run: func(cmd *cobra.Command, args []string) {
		jString := strings.Join(args, " ")
		if len(jString) == 0 {
			jString = config.Cfg().GetString("clog.jumbo.sample")
		}

		if listFonts {
			fmt.Println("Available fonts:")
			for _, font := range jumboFontList {
				fmt.Println("  " + font)
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

		prefix := ""
		suffix := ""
		switch commentStyle {
		case "go":
			prefix = "// "
			suffix = ""
		case "html":
			prefix = "<!-- "
			suffix = "-->"
		case "hugo":
			prefix = "{{/* "
			suffix = " */}}"
		case "js":
			prefix = "/* "
			suffix = " */"
		case "sh":
			prefix = "# "
			suffix = ""
		case "sql":
			prefix = "-- "
			suffix = ""
		case "tex":
			prefix = "% "
			suffix = ""
		case "none":
			prefix = ""
			suffix = ""
		default:
			prefix = "#"
			suffix = ""
		}

		myFigure := figure.NewFigure(strings.Join(args, " "), font, false)
		stringSlice := myFigure.Slicify()

		maxLen := 0
		for s := range stringSlice {
			if len(stringSlice[s]) > maxLen {
				maxLen = len(stringSlice[s])
			}
		}
		// print out the text with appropriate comments - padRight each line
		for s := range stringSlice {
			row := fmt.Sprintf("%s%*s%s", prefix, -maxLen, stringSlice[s], suffix)
			fmt.Println(row)
		}
	},
}

func init() {
	_, file, _, _ := runtime.Caller(0)
	slog.Debug("init " + file)

	Command.PersistentFlags().BoolVarP(&listFonts, "list", "L", false, "clog Jumbo -L         # list available fonts")
	Command.PersistentFlags().BoolVarP(&showFonts, "show", "S", false, "clog Jumbo -S text    # print text in all fonts")
	Command.PersistentFlags().StringVarP(&fontOverride, "font", "F", "", "clog Jumbo -F rounded \"text in rounded font\"")
	// https://gist.github.com/dk949/88b2652284234f723decaeb84db2576c
	Command.PersistentFlags().StringVarP(&commentStyle, "commentStyle", "C", "sh", "clog Jumbo -C none    # go|html|hugo|js|sh|sql|tex|none")
}
