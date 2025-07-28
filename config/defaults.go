//  Copyright ©2017-2025  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License  https://opensource.org/license/bsd-3-clause/
//
// clog's config package
//

package config

import (
	"bytes"
	"fmt"
	"log/slog"
	"os"
	"runtime"
)

var rootConfigFilename = "core.clog.yaml"


// root keys  using viper.setDefaults survive a config file merge
// child keys using viper.setDefaults are lost

func (cfg *Config) setDefaults(cfgPathOverride *string) {
	//tell viper that we're reading YAML otherwise it silently fails.
	cfg.SetConfigType("yaml")

	//find the embedded config file
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
		msg := fmt.Sprintf("config.setDefaults() failed with clog's embedded file system: %s", err.Error())
		panic(msg)
	}
	//parse & load the config
	cfgReader := bytes.NewReader(rootConfig)
	if err = cfg.ReadConfig(cfgReader); err != nil {
		msg := fmt.Sprintf("config.setDefaults() failed reading clog's embedded file system: %s", err.Error())
		panic(msg)
	}

	//overlay various other configs with configCLI being the highest priority
	slog.Info(cfg.GetString("clog.jumbo.font"))
	searchPaths = cfg.GetStringSlice("clog.clogrc.search-paths")
	if cfgPathOverride != nil && len(*cfgPathOverride) > 2 {
		searchPaths = append(searchPaths, *cfgPathOverride)
	}



	//store the startup defaults
	cfg.SetDefault("isInteractive", false)
}

func init() {
	// log the order of the init files in case there are problems
	_, file, _, _ := runtime.Caller(0)
	slog.Debug("init " + file)
}
