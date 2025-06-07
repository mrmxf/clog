//  Copyright ©2017-2025    Mr MXF   info@mrmxf.com
//  BSD-3-Clause License    https://opensource.org/license/bsd-3-clause/
//
// package cmd contains the default commands in a form that can be individually
// loaded by a fork of clog.
//
// a standard cobra app will do this in the `init()` function. The external
// bootstrap function allows the use of config functions & embedded Fs files.

package cmd

import (
	"embed"
	"log/slog"
	"runtime"

	"github.com/mrmxf/clog/cmd/cat"
	"github.com/mrmxf/clog/cmd/check"
	"github.com/mrmxf/clog/cmd/copy"
	"github.com/mrmxf/clog/cmd/crayon"
	"github.com/mrmxf/clog/cmd/inc"
	initialise "github.com/mrmxf/clog/cmd/init"
	"github.com/mrmxf/clog/cmd/jumbo"
	"github.com/mrmxf/clog/cmd/list"
	"github.com/mrmxf/clog/cmd/logcmd"
	"github.com/mrmxf/clog/cmd/should"
	"github.com/mrmxf/clog/cmd/snippets"
	"github.com/mrmxf/clog/cmd/source"
	"github.com/mrmxf/clog/cmd/version"
	"github.com/mrmxf/clog/config"
	"github.com/mrmxf/clog/scripts"
	"github.com/mrmxf/clog/semver"
	"github.com/mrmxf/clog/ux"
	"github.com/spf13/cobra"
)

func BootStrap(bootCmd *cobra.Command) error {
	cfg := config.Cfg()
	historyFilename := cfg.GetString("clog.releases-path")

	// find the embedded release history file in the embedded file systems
	// last one found wins - this is usually the project's embedded fs
	eFs, paths, err := config.FindEmbedded(historyFilename)
	if err != nil {
		// warn on errors - this should not cause a panic
		slog.Warn("cannot find embedded release history " + err.Error())
	} else {
		initialiseConfigFromSemver(eFs, paths)
	}

	// prepend cobra usage strings with build information
	bootCmd.SetUsageTemplate(cfg.GetString("clog.version.long") + bootCmd.UsageTemplate())

	// load all the public builtin commands first
	bootCmd.AddCommand(cat.Command)        // script helper include command
	bootCmd.AddCommand(check.Command)      // copy an embedded file to a destination
	bootCmd.AddCommand(copy.Command)       // copy an embedded file to a destination
	bootCmd.AddCommand(crayon.Command)     // colored terminal commands
	bootCmd.AddCommand(inc.Command)        // script helper include command
	bootCmd.AddCommand(initialise.Command) // create a clogrc
	bootCmd.AddCommand(jumbo.Command)      // Jumbo text output
	bootCmd.AddCommand(list.Command)       // list embedded files text output
	bootCmd.AddCommand(logcmd.Command)     // list embedded files text output
	bootCmd.AddCommand(should.Command)     // logic helper for bash scripts
	bootCmd.AddCommand(source.Command)     // source a script or snippet
	bootCmd.AddCommand(version.Command)    // version reporting

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

	// load shell scripts so that they override snippets if there's a clash
	scripts.FindScripts(bootCmd, "clogrc/*.sh")

	// build the UX menus in case we're running interactively
	ux.BuildMenus(bootCmd)

	// Finally, Execute the cobra command parser on the configured hierarchy
	// the return value of the command is returned to the shell
	return bootCmd.Execute()
}

func initialiseConfigFromSemver(eFs *embed.FS, paths []string) {
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
}

func init() {
	_, file, _, _ := runtime.Caller(0)
	slog.Debug("init " + file)
}
