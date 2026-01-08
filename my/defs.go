//  Copyright ©2017-2025  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License           https://opensource.org/license/bsd-3-clause/
// This file is part of clog.

// package my has the configuration structs and definitions for this version of clog.

package my

import (
	"os"

	"github.com/mrmxf/clog-mrmxf/embedfilesystem"
	"github.com/mrmxf/clog-mrmxf/kfg"
	"github.com/mrmxf/clog-mrmxf/semver"
)

// Import our custom middleware - update this path to match your project structure

// control variables returned by functions in the App struct
var DebugFlag bool
var DryrunFlag bool

// AppConfig represents the top-level configuration structure for the clog application.
// It maps to the "clog:" section in the YAML configuration file.
// The `yaml` tags tell the koanf library how to map YAML keys to struct fields.
type AppConfig struct {
	Env          EnvConfig   `yaml:"env"`           // env variable names used by subcommands
	Jumbo        JumboConfig `yaml:"jumbo"`         // settings for the jumbo text display
	Log          LogConfig   `yaml:"log"`           // logging behavior
	ReleasesPath string      `yaml:"releases-path"` // path of releases.yaml file  for version tracking
	StashPath    string      `yaml:"stash-path"`    // path of stash file for bash/zsh commands
	Title        string      `yaml:"title"`         // title of the app used for logging
	// properties not loaded from YAML
	Releases []kfg.AppRelease   // Releases information from the file pointed to by ReleasesPath
	SemVer   semver.VersionInfo // Semantic Version information via semver library
}

// EnvConfig defines the names of environment variables used by various external integrations.
// This allows users to customize which environment variables clog looks for.
type EnvGroup map[string]string
type EnvConfig map[string]EnvGroup

// JumboConfig contains settings for the jumbo text display feature.
// This feature can display large ASCII art text for various purposes.
type JumboConfig struct {
	Font   string `yaml:"font"`   // Font specifies which font style to use for jumbo text (e.g., "small")
	Sample string `yaml:"sample"` // Sample is example text used for testing or demonstration purposes
}

// LogConfig controls how clog outputs log messages.
// This affects both the verbosity and formatting of log output.
type LogConfig struct {
	Level string `yaml:"level"` // Level determines which log messages are shown (debug, info, warn, error)
	Style string `yaml:"style"` // Style controls the formatting of log output (plain, pretty, json)
}

// AppKey is the default configuration key for the main application configuration section.
// It is read from the embedded konfig.yaml with at `kfg.app-key`
var AppKey = "clog"

var appStruct = AppConfig{}

var App = &appStruct

func (app *AppConfig) IsDebug() bool {
	return DebugFlag
}

func (app *AppConfig) IsDryRun() bool {
	return DryrunFlag
}

var Use = "clog"
var ShortHelp = "Command Line Of Go - interactive helper"
var LongHelp = `
Command Line Of Go (clog)
=========================
Clog aggregates:
  - snippets: command lines in your project's clog.yaml
	-  scripts: files matching "clogrc/*.sh" - see below
	- commands: embedded functions compiled into clog

Create clog.yaml for a project
==========================================
clog Init  # run it twice to get a copy of the konfig.yaml

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

// initialise the boot options for my App
var BootOptions = kfg.KonfigureOpt{
	// boot configuration Fs, file
	AppFs:    embedfilesystem.CoreFs,
	FilePath: "konfig.yaml",
	// allow config to see the args before Cobra initialises
	AppArgs: &os.Args,
	// struct definitions to receive the App configuration
	AutoAppStruct:          App,
	AutoAppKey:             "clog",
	AutoAppAnnotationLabel: "yaml",
	AutoReleaseSlice:       &(App.Releases),
}
