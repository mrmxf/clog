//  Copyright ©2018-2025  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License           https://opensource.org/license/bsd-3-clause/
//
// pretry printing and other functions for semver

package semver

import (
	_ "embed"
	"unicode"
)

// iterate through a string and highlight it for display on a TTY.
//
// capital letters at the start of words use the capitals `c` highlighter,
// everything else uses the body `b` highlighter.
func highlightTitleCase(str string) string {
	var pen = Color()

	c := pen.Success
	b := pen.Info
	res := ""
	skipped := ""

	for _, ch := range str {
		if unicode.IsUpper(ch) && unicode.IsLetter(ch) {
			if len(skipped) > 0 {
				res += b(skipped)
			}
			res += c(string(ch))
			skipped = ""
		} else {
			skipped = skipped + string(ch)
		}
	}
	if len(skipped) > 0 {
		res += b(skipped)
	}
	return res
}
