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

var defaultConfigFilename = "core.clog.yaml"

// root keys  using viper.setDefaults survive a config file merge
// child keys using viper.setDefaults are lost

func (cfg *Config) setDefaults(configCLI *string) {
	cfg.SetConfigName("clog.config")
	cfg.SetConfigType("yaml")

	fs, configPaths, err := FindEmbedded(defaultConfigFilename)
	if err != nil {
		slog.Error("embedded config ()" + defaultConfigFilename + ") " + err.Error())
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
	if err = cfg.ReadConfig(bytes.NewBuffer(rootConfig)); err != nil {
		msg := fmt.Sprintf("config.setDefaults() failed reading clog's embedded file system: %s", err.Error())
		panic(msg)
	}

	//overlay various other configs with configCLI being the highest priority
	searchPaths = cfg.GetStringSlice("clog.clogrc.search-order")
	if configCLI != nil && len(*configCLI) > 2 {
		searchPaths = append(searchPaths, *configCLI)
	}

	if homeFolder, err := os.UserHomeDir(); err == nil {
		cfg.Set("clog.homeFolder", homeFolder)
	}

	//store the startup defaults
	cfg.SetDefault("isInteractive", false)
	cfg.SetDefault("clog.clogrc.config_base", ".clog")
	cfg.SetDefault("clog.clogrc.config_format", "yaml")

	//override the defaults in viper for finding other configs
	if fn := cfg.GetString("clog.clogrc.config_base"); len(fn) > 0 {
		cfg.SetConfigName(fn)
	}
	if ff := cfg.GetString("clog.clogrc.config_format"); len(ff) > 0 {
		cfg.SetConfigType(ff)
	}
}

func init() {
	// log the order of the init files in case there are problems
	_, file, _, _ := runtime.Caller(0)
	slog.Debug("init " + file)
}
