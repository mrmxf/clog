//  Copyright ©2018-2025    Mr MXF   info@mrmxf.com
//  BSD-3-Clause License    https://opensource.org/license/bsd-3-clause/
//
// package cmd contains the default commands in a form that can be individually
// loaded by a fork of clog.
//
// a standard cobra app will do this in the `init()` function. The external
// bootstrap function allows the use of config functions & embedded Fs files.

package cmd

import (
	"log/slog"
	"runtime"

	"github.com/mrmxf/clog/cmd/cat"
	"github.com/mrmxf/clog/cmd/inc"
	"github.com/mrmxf/clog/cmd/jumbo"
	"github.com/mrmxf/clog/cmd/list"
	"github.com/mrmxf/clog/cmd/version"
	"github.com/mrmxf/clog/config"
	"github.com/mrmxf/clog/semver"
	"github.com/mrmxf/clog/ux"
	"github.com/spf13/cobra"
	// "github.com/mrmxf/clog/git"
	// "github.com/mrmxf/clog/scripts"
	// "github.com/mrmxf/clog/ux"
)

// define a chainable Bootstrap Function type
type BootStrapFunc func (rootCmd *cobra.Command, chain BootStrapFunc) error

func BootStrap(rootCmd *cobra.Command) error {
	cfg := 	config.Cfg()
	historyFilename := cfg.GetString("clog.history_file")

	// find the embedded release history file in the embedded file systems
	// last one found wins - this is usually the project's embedded fs
	// warn on errors - this should not cause a panic
	eFs, paths, err := config.FindEmbedded(historyFilename)
	if err != nil {
		slog.Warn("cannot find embedded release history "+ err.Error())
	}else{
		err = semver.Initialise(*eFs, paths[len(paths)-1])
		if err != nil {
			slog.Warn("cannot initialize semantic version "+ err.Error())
		}else{
			//override default empty strings with real semver info
			cfg.Set("clog.version.long", semver.Info().Long)
			cfg.Set("clog.version.note", semver.Info().Note)
			cfg.Set("clog.version.short", semver.Info().Short)
			cfg.Set("clog.version.appname", semver.Info().AppName)
			cfg.Set("clog.version.apptitle", semver.Info().AppTitle)
			cfg.RegisterAlias("ver", "clog.version.short")
			cfg.RegisterAlias("app", "clog.version.appname")
			cfg.RegisterAlias("title", "clog.version.apptitle")
		}

}

	//prepend cobra usage strings with build information
	rootCmd.SetUsageTemplate(cfg.GetString("clog.version.long") + rootCmd.UsageTemplate())

	rootCmd.AddCommand(cat.Command)     		// script helper include command
	rootCmd.AddCommand(inc.Command)     		// script helper include command
	rootCmd.AddCommand(jumbo.Command)       // Jumbo text output
	rootCmd.AddCommand(list.Command)        // list embedded files text output
	rootCmd.AddCommand(version.Command)			// version reporting


	//parse the config (key=snippets) for all snippets and add them before the scripts
	//clCmd.AddCommandSnippets(rootCmd)
	// clCmd.AddSnippets(rootCmd, "snippets")

	//look for all shell scripts in the clogrc folder &&  ignore errors
	// clScripts.FindScriptsToAdd(rootCmd, "clogrc/*.sh")

	//add in top level clog commands
	// rootCmd.AddCommand(clCmd.CatCmd)     // Cat an embedded file
	// rootCmd.AddCommand(clCheck.CheckCmd) // Check the current project
	// rootCmd.AddCommand(clCmd.InitCmd)    // Create a new config file
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
	// clCore.BootStrap(rootCmd)
	// clGit.BootStrap(rootCmd)
	// clDocker.BootStrap(rootCmd)

	//add in subcommands for the various sub-packages
	// once all the subcommands are loaded, the menus can be built
	ux.BuildMenus(rootCmd)

	
	// Finally, Execute the cobra command parser
		return RootCommand.Execute()
}

func init() {
	_, file, _, _ := runtime.Caller(0)
	slog.Debug("init " + file)
}
