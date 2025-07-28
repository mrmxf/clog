//  Copyright ©2017-2025  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License  https://opensource.org/license/bsd-3-clause/

// clog's cfg package

package cfg

import (
	"fmt"
	"log/slog"
	"os"
	"runtime"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/rawbytes"
)

var rootConfigFilename = "core.clog.yaml"
var usrConfigBasename = "clog"
var usrConfigType = "yaml"

func (kfg *Config) setDefaults(cfgPathOverride *string) {
	fs, configPaths, err := FindEmbedded(rootConfigFilename)
	if err != nil {
		slog.Error("embedded config (" + rootConfigFilename + ") " + err.Error())
		os.Exit(1)
	}
	//set the coreFs for other tasks to find the core files
	coreFs = *fs

	//read the first root config found
	rootConfig, err := fs.ReadFile(configPaths[0])
	if err != nil {
		msg := fmt.Sprintf("cfg.setDefaults() failed with clog's embedded file system: %s", err.Error())
		panic(msg)
	}
	//parse & load the config
	if err = kfg.Load(rawbytes.Provider(rootConfig), yaml.Parser()); err != nil {
		msg := fmt.Sprintf("cfg.setDefaults() failed reading clog's embedded file system: %s", err.Error())
		panic(msg)
	}

	//overlay various other configs with configCLI being the highest priority
	searchPaths = kfg.Strings("clog.clogrc.search-paths")
	if cfgPathOverride != nil && len(*cfgPathOverride) > 2 {
		searchPaths = append(searchPaths, *cfgPathOverride)
	}

	if homeFolder, err := os.UserHomeDir(); err == nil {
		kfg.Set("clog.env.HOME", homeFolder)
	}

	//store the startup defaults
	kfg.Set("isInteractive", false)
}

func init() {
	// log the order of the init files in case there are problems
	_, file, _, _ := runtime.Caller(0)
	slog.Debug("init " + file)
}
