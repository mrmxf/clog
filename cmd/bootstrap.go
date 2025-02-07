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
	"github.com/mrmxf/clog/cmd/copy"
	"github.com/mrmxf/clog/cmd/crayon"
	"github.com/mrmxf/clog/cmd/inc"
	initialise "github.com/mrmxf/clog/cmd/init"
	"github.com/mrmxf/clog/cmd/jumbo"
	"github.com/mrmxf/clog/cmd/list"
	"github.com/mrmxf/clog/cmd/snippets"
	"github.com/mrmxf/clog/cmd/source"
	"github.com/mrmxf/clog/cmd/version"
	"github.com/mrmxf/clog/config"
	"github.com/mrmxf/clog/scripts"
	"github.com/mrmxf/clog/semver"
	"github.com/mrmxf/clog/ux"
	"github.com/spf13/cobra"
	// "github.com/mrmxf/clog/git"
	// "github.com/mrmxf/clog/scripts"
	// "github.com/mrmxf/clog/ux"
)

func BootStrap(rootCmd *cobra.Command) error {
	cfg := config.Cfg()
	historyFilename := cfg.GetString("clog.history-file")

	// find the embedded release history file in the embedded file systems
	// last one found wins - this is usually the project's embedded fs
	eFs, paths, err := config.FindEmbedded(historyFilename)
	if err != nil {
		// warn on errors - this should not cause a panic
		slog.Warn("cannot find embedded release history " + err.Error())
	} else {
		initialiseConfigFromSemver(eFs, paths)
	}

	//prepend cobra usage strings with build information
	rootCmd.SetUsageTemplate(cfg.GetString("clog.version.long") + rootCmd.UsageTemplate())

	// load all the builtin commands first
	rootCmd.AddCommand(cat.Command)        // script helper include command
	rootCmd.AddCommand(copy.Command)       // copy an embedded file to a destination
	rootCmd.AddCommand(crayon.Command)     // colored terminal commands
	rootCmd.AddCommand(inc.Command)        // script helper include command
	rootCmd.AddCommand(initialise.Command) // create a clogrc
	rootCmd.AddCommand(jumbo.Command)      // Jumbo text output
	rootCmd.AddCommand(list.Command)       // list embedded files text output
	rootCmd.AddCommand(source.Command)     // source a script or snippet
	rootCmd.AddCommand(version.Command)    // version reporting

	// create a new snippets command from the clog.snippets cfg() branch
	branchKey := "snippets"
	opts := snippets.SnippetsCmdOpts{
		Use:     "Snippets",
		Key:     branchKey,
		Verbose: false,
		Plain:   false,
		Raw:     cfg.GetStringMap(branchKey),
	}
	snippetsTree := snippets.NewSnippetsCommand(rootCmd, opts)
	rootCmd.AddCommand(snippetsTree) // main snippets

	// load shell scripts so that they override snippets if there's a clash
	scripts.FindScripts(rootCmd, "clogrc/*.sh")

	//add in top level clog commands
	// rootCmd.AddCommand(clCheck.CheckCmd) // Check the current project
	// rootCmd.AddCommand(clCmd.LintCmd)    // Lint the current project with megalinter
	// rootCmd.AddCommand(clCmd.ShCmd)           // run a snippet
	// rootCmd.AddCommand(clCmd.ListSnippetsCmd) // list snippets

	// add in the top level menus for any child commands
	// rootCmd.AddCommand(clCi.CiCmd)         // Core commands
	// rootCmd.AddCommand(clCore.CoreCmd)     // Core commands
	// rootCmd.AddCommand(clGit.GitCmd)       // Git tools
	// rootCmd.AddCommand(clDocker.DockerCmd) // Docker tools

	// now bootstrap child commands
	// clCi.BootStrap(rootCmd)
	// clGit.BootStrap(rootCmd)
	// clDocker.BootStrap(rootCmd)

	//add in subcommands for the various sub-packages
	// once all the subcommands are loaded, the menus can be built
	ux.BuildMenus(rootCmd)

	// Finally, Execute the cobra command parser
	return RootCommand.Execute()
}

func initialiseConfigFromSemver(eFs *embed.FS, paths []string) {
	err := semver.Initialise(*eFs, paths[len(paths)-1])
	if err != nil {
		slog.Warn("cannot initialize semantic version " + err.Error())
	} else {
		//override default empty strings with real semver info
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
