//  Copyright ©2017-2025  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License  https://opensource.org/license/bsd-3-clause/
//
// clog's config package - wrap Viper for multi-project init
//
// usage:
// Cfg := config.New( &[]embed.FS{CoreFs, AppFs, OtherFs} )
//
//  if cfg := config.Cfg(); cfg==nil{
//		os.ExplodeFrontPanel()
//    os.Exit(255)
//  }
//
//
// clogrc search order - see core/clog.yaml

package config

import (
	"embed"
	"fmt"
	"log/slog"
	"os"
	"runtime"

	"github.com/spf13/viper"
)

// export the Fs with the `core/` folder - initialised by calling program
var coreFs embed.FS

// Config type embeds the Viper struct and extends it
type Config struct {
	*viper.Viper
}

// the global config variable used by other packages
var cfg *Config = nil

// a cache of the Fs slice for future extension
var fsCache []embed.FS = []embed.FS{}

// the searchPaths from the config can be exported
var searchPaths []string

// create a new config object
//
// fsSlice is a slice of embed.FS objects that are searched search for configs
// cfgPathOverride is the name of the config file to load
//
//	if nil or empty use `clog.yaml`
func New(fsSlice *[]embed.FS, cfgPathOverride *string) *Config {
	//initialise viper with logger that can be uses throughout clog
	cfg = &Config{
		viper.New(),
	}

	//preserve fsCache for future use
	if len(*fsSlice) > 0 {
		fsCache = append(fsCache, *fsSlice...)
	} else {
		slog.Error("No embed.FS passed to config.New(), cannot bootstrap clog")
		os.Exit(1)
	}

	// populate a new config object, load in the embedded config and set the
	// initial search paths to find other configs to overlay
	cfg.setDefaults(cfgPathOverride)

	// Merge the config file with defaults - ignore errors
	cfg.mergeAllConfigs()

	//enable auto-import of `env` variables declared in config
	// e.g. AWS_ACCESS_KEY_ID becomes cfg.GetString("AWS_ACCESS_KEY_ID")
	cfg.AutomaticEnv()
	//iterate through the environment variables declared and bind them in cfg
	for envIdentifier, envVariableName := range cfg.GetStringMapString("clogrc.env") {
		if false {
			fmt.Println(envIdentifier, envVariableName)
		}
		err := cfg.BindEnv(envVariableName)
		if err != nil {
			slog.Warn("Failed to bind environment variable: " + envVariableName + " Error: " + err.Error())
		}
	}
	return cfg
}

// get the config object
func Cfg() *Config {

	return cfg
}

// get the coreFs
func CoreFs() embed.FS {
	return coreFs
}

// get the ordered cache of embedded file systems searched
func FsCache() []embed.FS {
	return fsCache
}

func SearchPaths() *[]string {
	return &searchPaths
}

func init() {
	// log the order of the init files in case there are problems
	_, file, _, _ := runtime.Caller(0)
	slog.Debug("init " + file)
}
