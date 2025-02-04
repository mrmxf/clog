//  Copyright ©2019-2024  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License  https://opensource.org/license/bsd-3-clause/

// Package scripts adds support for local bash scripts

package scripts

import (
	"path/filepath"

	"github.com/mrmxf/clog/crayon"
	"github.com/spf13/cobra"
)

var c = crayon.Color()
var scriptsMap = map[string]map[string]string{}

// Add scripts from clogrc folder
func FindScriptsToAdd(rootCmd *cobra.Command, folderGlob string) {

	//look for all shell scripts in the clogrc folder
	scripts, _ := filepath.Glob(folderGlob)

	//add each script found
	for _, script := range scripts {
		AddScript(rootCmd, script)
	}
}
