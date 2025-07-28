//  Copyright ©2017-2025  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License           https://opensource.org/license/bsd-3-clause/

// Package main is a sample build of clog.
//
// All the sample commands appear here to show how you might build your own
// variant of clog.

package main

import (
	"embed"
	"log/slog"
	"os"
	"runtime"

	"github.com/mrmxf/clog/cmd/aws"
	"github.com/mrmxf/clog/cmd/cat"
	"github.com/mrmxf/clog/cmd/check"
	"github.com/mrmxf/clog/cmd/crayon"
	"github.com/mrmxf/clog/cmd/inc"
	"github.com/mrmxf/clog/cmd/jumbo"
	"github.com/mrmxf/clog/cmd/list"
	"github.com/mrmxf/clog/cmd/logcmd"
	"github.com/mrmxf/clog/cmd/should"
	"github.com/mrmxf/clog/cmd/snippets"
	"github.com/mrmxf/clog/cmd/source"
	"github.com/mrmxf/clog/cmd/version"
	"github.com/mrmxf/clog/config"
	"github.com/mrmxf/clog/core"
	"github.com/mrmxf/clog/scripts"
	"github.com/mrmxf/clog/semver"
	loggerPkg "github.com/mrmxf/clog/slogger"
	"github.com/mrmxf/clog/ux"
	"github.com/mrmxf/clog/ux/ui"
	"github.com/spf13/cobra"
)

//go:embed releases.yaml
var AppFs embed.FS

var CoreFs embed.FS = core.CoreFs
var cfgPathOverride string

// CLI flag for changing the configuration file path
var ConfigFilePath string

// turn on debug logging
var ConfigDebug bool

// CLI flag for showing the long version string
var ShowVersion bool

// CLI flag for showing the short version string
var ShowVersionShort bool

// CLI flag for showing the note associated with this version
var ShowVersionNote bool

// CLI flag for changing the logging level
var LogLevel int

const shortHelp = "Command Line Of Go - try clog --help for info"
const longHelp = `
Command Line Of Go (clog)
=========================
Clog aggregates:
  - snippets: command lines in your porject's colg.yaml
	-  scripts: files matching "clogrc/*.sh" - see below
	- commands: embedded functions compiled into clog

Create clog.yaml for a project
==========================================
clog Init  # run it twice to get a copy of the core.clog.yaml

Scripts in "clogrc/" must have the following 3 lines to be found by clog
==========================================
#  clog> commandName
# short> short help text
# extra> scripts need these 3 lines to be found by clog

Adding Snippets & macros
==========================================
edit clogrc/clog.yaml  # after you've made one

Running clog
==========================================
interactively: clog
as a web ui:   clog Svc && open localhost:8765
as api:      	 curl -H "Authorization: OAuth <ACCESS_TOKEN>" http://localhost:8765/api/version/command
`

const clogEmbeddedReleasesFile = "releases.yaml"

// this is the root commmand  that is executed when you type "clog"
var rootCommand = &cobra.Command{
	Use:   "clog",
	Short: shortHelp,
	Long:  longHelp,
	Run: func(cmd *cobra.Command, args []string) {

		// Show the version string (and exit) if flags are set
		if ShowVersion {
			version.Command.Run(cmd, []string{})
		}
		if ShowVersionShort {
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

// The main program pulls in all the clog configs and `embed` file systems. All
// config files will be loaded in the order they are found.

func main() {
	// Initialize configuration system
	// this happens in main (and not init) because all flags must be parsed so
	// that we can use an external file system override.
	configEmbedFsList := &[]embed.FS{CoreFs, AppFs}
	config.New(configEmbedFsList, &cfgPathOverride)

	// Execute the cobra command parser on the configured hierarchy
	// the return value of the command is returned to the shell
	err := rootCommand.Execute()
	if err != nil {
		loggerPkg.Error(err.Error())
		os.Exit(1)
	}
}

// if you want to log the `init()` order for this application then set the
// default log level to LevelDebug inside the slogger package
func init() {
	_, file, _, _ := runtime.Caller(0)
	slog.Debug("init " + file)

	loggerPkg.UsePrettyLogger(loggerPkg.LevelInfo)

	// find the embedded release history file in the embedded file systems
	// last one found wins - this is usually the project's embedded fs
	eFs, paths, err := config.FindEmbedded(clogEmbeddedReleasesFile)
	if err != nil {
		// warn on errors - this should not cause a panic
		slog.Warn("cannot find embedded release history " + err.Error())
	} else {
		updateConfigFromSemver(eFs, paths)
	}

	// Define persistent (global) flags and any flags for the root command
	rootCommand.PersistentFlags().StringVarP(&ConfigFilePath, "config", "c", "", "clog -c myClogfig.yaml   # clog Cat core.clog.yaml > myClogfig.yaml")
	rootCommand.PersistentFlags().BoolVarP(&ConfigDebug, "debug", "d", false, "clog --debug             # set debug level to --debug")
	rootCommand.PersistentFlags().BoolVar(&ShowVersion, "version", false, "clog --version           # shows the full version string")
	rootCommand.PersistentFlags().BoolVarP(&ShowVersionShort, "v", "v", false, "clog -v                  # shows just the semantic version")
	rootCommand.PersistentFlags().BoolVarP(&ShowVersionNote, "note", "n", false, "clog --note              # shows just the version note")
	rootCommand.PersistentFlags().IntVarP(&LogLevel, "loglevel", "l", 0, "clog --loglevel 1        # 0:OFF 1:DEBUG 2:INFO 3:WARN 4:ERROR")

	// add in the subcommands for the rootCommand
	// load all the public builtin commands first
	rootCommand.AddCommand(aws.Command)        // aws day-to-day management commands
	rootCommand.AddCommand(cat.Command)        // script helper include command
	rootCommand.AddCommand(check.Command)      // copy an embedded file to a destination
	rootCommand.AddCommand(copy.Command)       // copy an embedded file to a destination
	rootCommand.AddCommand(crayon.Command)     // colored terminal commands
	rootCommand.AddCommand(inc.Command)        // script helper include command
	rootCommand.AddCommand(initialise.Command) // create a clogrc
	rootCommand.AddCommand(jumbo.Command)      // Jumbo text output
	rootCommand.AddCommand(list.Command)       // list embedded files text output
	rootCommand.AddCommand(logcmd.Command)     // list embedded files text output
	rootCommand.AddCommand(should.Command)     // logic helper for bash scripts
	rootCommand.AddCommand(source.Command)     // source a script or snippet
	rootCommand.AddCommand(version.Command)    // version reporting


	// load shell scripts so that they override snippets if there's a clash
	scripts.FindScripts(rootCommand, "clogrc/*.sh")

	// build the UX menus in case we're running interactively
	ux.BuildMenus(rootCommand)
}

// runtimeInit cannot run during init() because all the flags have to be 
// parsed and the configs loaded once the embeddedInit is done.
func runtimeInit(){
	cfg := config.Cfg()

	// create a new snippets command from the clog.snippets cfg() branch
	branchKey := "snippets"
	opts := snippets.SnippetsCmdOpts{
		Use:     "Snippets",
		Key:     branchKey,
		Verbose: false,
		Plain:   false,
		Raw:     cfg.GetStringMap(branchKey),
	}
	snippetsTree := snippets.NewSnippetsCommand(bootCmd, opts)
	bootCmd.AddCommand(snippetsTree) // main snippets
}
func updateConfigFromSemver(eFs *embed.FS, paths []string) {
	err := semver.Initialise(*eFs, paths[len(paths)-1])
	if err != nil {
		slog.Warn("cannot initialize semantic version " + err.Error())
	} else {
		// override default empty strings with real semver info
		config.Cfg().Set("clog.version.long", semver.Info().Long)
		config.Cfg().Set("clog.version.note", semver.Info().Note)
		config.Cfg().Set("clog.version.short", semver.Info().Short)
		config.Cfg().Set("clog.version.appname", semver.Info().AppName)
		config.Cfg().Set("clog.version.apptitle", semver.Info().AppTitle)
		config.Cfg().RegisterAlias("ver", "clog.version.short")
		config.Cfg().RegisterAlias("app", "clog.version.appname")
		config.Cfg().RegisterAlias("title", "clog.version.apptitle")
	}
	// prepend cobra usage strings with build information
	rootCommand.SetUsageTemplate(config.Cfg().GetString("clog.version.long") + rootCommand.UsageTemplate())
}
