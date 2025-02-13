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

	"github.com/mrmxf/clog/cmd"
	"github.com/mrmxf/clog/config"
	"github.com/mrmxf/clog/slogger"
)

//go:embed releases.yaml
var AppFs embed.FS

//go:embed core
var CoreFs embed.FS

// The main program pulls in all the clog configs and `embed` file systems. All
// config files will be loaded in the order they are found.

func main() {
	// create a list of embedded file systems to init the system
	configEmbedFsList := &[]embed.FS{CoreFs, AppFs}

	//Overlay all configs by searching the `embed` File Systems
	config.New(configEmbedFsList, nil)

	// BootStrap clog with the RootCommand
	// load commands, snippets, scripts & update config as needed
	err := cmd.BootStrap(cmd.RootCommand)
	if err != nil {
		slog.Error(err.Error())
		os.Exit(1)
	}
}

// if you want to log the `init()` order for this application then set the
// default log level to LevelDebug inside the slogger package
func init() {
	_, file, _, _ := runtime.Caller(0)
	slog.Debug("init " + file)

	slogger.UsePrettyLogger(slog.LevelError)
	// slogger.UsePrettyLogger(slog.LevelDebug)
}
