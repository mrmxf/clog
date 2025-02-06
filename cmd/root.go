//  Copyright ©2017-2025    Mr MXF   info@mrmxf.com
//  BSD-3-Clause License    https://opensource.org/license/bsd-3-clause/
//
// command:
//   clog
//
// This is the default implementation of a rootCmd along with the default
// help text and command names. This file also defines shared package vars

package cmd

import (
	"log/slog"
	"runtime"

	"github.com/mrmxf/clog/cmd/version"
	"github.com/mrmxf/clog/ux/ui"
	"github.com/spf13/cobra"
)

// CLI flag for changing the configuration file path
var ConfigFilePath string

// CLI flag for showing the long version string
var ShowVersion bool

// CLI flag for showing the short version string
var ShowVersionShort bool

// CLI flag for showing the note associated with this version
var ShowVersionNote bool

// CLI flag for changing the logging level
var LogLevel int

var RootCommand = &cobra.Command{
	Use:   "clog",
	Short: "Command Line Of Go - interactive helper",
	Long: `
Command Line Of Go (clog)
=========================
Clog aggregates:
  - snippets: command lines in your porject's clog.config.yaml
	-  scripts: files matching "clogrc/*.sh" - see below
	- commands: embedded functions compiled into clog

Create clog.config.yaml for a project
==========================================
clog Init  # run it twice to get a copy of the clog.core.config.yaml

Scripts in "clogrc/" must have the following 3 lines to be found by clog
==========================================
#  clog> commandName
# short> short help text
# extra> scripts need these 3 lines to be found by clog

Adding Snippets & macros
==========================================
edit clogrc/clog.config.yaml  # after you've made one

Running clog
==========================================
interactively: clog
as a web ui:   clog Svc && open localhost:8765
as api:      	 curl -H "Authorization: OAuth <ACCESS_TOKEN>" http://localhost:8765/api/version/command
`,
	Run: func(cmd *cobra.Command, args []string) {

	// Show the version string (and exit) if flags are set
	if ShowVersion {
		version.Command.Run(cmd, []string{})
	}
	if ShowVersionShort{
		version.Command.Run(cmd, []string{"short"})
	}
	if ShowVersionNote {
		version.Command.Run(cmd, []string{"note"})
	}
	
		//we get here if no subcommands appeared in the command line
		//or the home key was pressed on the server
		_, err := ui.HomeMenu(cmd)
		if err != nil {
			panic(err)
		}
	},
}

func init() {
	_, file, _, _ := runtime.Caller(0)
	slog.Debug("init " + file)

	// Define persistent (global) flags and any flags for the root command
	RootCommand.PersistentFlags().StringVarP(&ConfigFilePath, "config", "c", "", "clog -c myClogfig.yaml   # clog Core Cat clogrc/core/clog.clConfig.yaml > myClogfig.yaml")
	RootCommand.PersistentFlags().BoolVar(&ShowVersion, "version", false, "clog --version           # shows the full version string")
	RootCommand.PersistentFlags().BoolVarP(&ShowVersionShort, "v", "v", false, "clog -v                  # shows just the semantic version")
	RootCommand.PersistentFlags().BoolVarP(&ShowVersionNote, "note", "n", false, "clog --note              # shows just the version note")
	RootCommand.PersistentFlags().IntVarP(&LogLevel, "loglevel", "L", 0, "clog --loglevel 1        # 0:OFF 1:DEBUG 2:INFO 3:WARN 4:ERROR")
}
