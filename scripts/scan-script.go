//  Copyright ©2017-2025  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License  https://opensource.org/license/bsd-3-clause/

// Package cmd implements commands for the cobra CLI library

package scripts

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"regexp"
)

type scriptMeta struct {
	usage     string
	shortHelp string
	longHelp  string
	// option        []string
	defaultOption string
}

// process and options String
func processOptions(meta scriptMeta, metaMap map[string]string) {

}

// scan a script file and create a map of the data found
// First 3 lines should be of the format:
//
//	#  clog> command [ opt1 | opt | opt3 ]
//	# short> a short help explanation
//	# extra> long help printed when the help command is issued
func scriptMetadata(filePath string) (map[string]string, error) {
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Print(err)
		return nil, err
	}
	scanner := bufio.NewScanner(file)
	scanner.Split(bufio.ScanLines)

	metaMap := make(map[string]string)
	for scanner.Scan() {
		// extract keyword & text into matching strings
		// hash whitespace (word[2-12 chars]) whitespace > (the rest)
		// e.g. # short> This is short help
		r, _ := regexp.Compile(`#\s+([a-z]{2,12})\s*>\s+(.*)`)

		textLine := scanner.Text()
		match := r.FindStringSubmatch(textLine)
		if len(match) == 3 {
			// e.g. meta[short] = "This is short help"
			metaMap[match[1]] = match[2]
		}
	}

	file.Close()

	// put the data into the struct
	silent := true
	_, gotUsage := metaMap["usage"]
	if !gotUsage && !silent {
		fmt.Println(c.E("Error"), "script "+c.F(filePath)+c.E(" missing usage. Please prepend comments to script:"))
		fmt.Println(c.E("     "), "# ", c.W("usage"), ">", c.C("command [flags] [args]"))
		fmt.Println(c.E("     "), c.I("and optionally"))
		fmt.Println(c.E("     "), "# ", c.W("short"), "> some short help text for the menu")
		fmt.Println(c.E("     "), "#  ", c.W("long"), "> some long help text for the", c.C("--help"), "option")
		fmt.Println(c.E("     "), "#  ", c.W("opts"), "> option1 | option2 | [defaultOption] | option4")
		fmt.Println(c.E("     "), c.E("Ignoring script"))
		return nil, errors.New(filePath + " is missing Clog metadata at the start - see clog --help")
	}
	meta := scriptMeta{
		usage:         metaMap["usage"],
		shortHelp:     metaMap["short"],
		longHelp:      metaMap["long"],
		defaultOption: "",
	}
	_, gotOpts := metaMap["opts"]
	if gotOpts {
		processOptions(meta, metaMap)
	}
	return metaMap, nil
}
