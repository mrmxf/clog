//  Copyright ©2017-2025  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License           https://opensource.org/license/bsd-3-clause/

// Package cfg - wrap Koanf for multi-project init (koanf-based version of config package)
//
// usage:
// Kfg := cfg.New( &[]embed.FS{CoreFs, AppFs, OtherFs} )
//
//  if kfg := cfg.Kfg(); kfg==nil{
//		os.ExplodeFrontPanel()
//    os.Exit(255)
//  }
//
//
// clogrc search order - see core/clog.yaml

package cfg

import (
	"embed"
	"fmt"
	"log/slog"
	"os"
	"runtime"

	"github.com/knadh/koanf/v2"
)

// export the Fs with the `core/` folder - initialised by calling program
var coreFs embed.FS

// Config type embeds the Koanf struct and extends it
type Config struct {
	*koanf.Koanf
}

// the global config variable used by other packages
var kfg *Config = nil

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
	//initialise koanf with logger that can be uses throughout clog
	kfg = &Config{
		koanf.New("."),
	}

	//preserve fsCache for future use
	if len(*fsSlice) > 0 {
		fsCache = append(fsCache, *fsSlice...)
	} else {
		slog.Error("No embed.FS passed to cfg.New(), cannot bootstrap clog")
		os.Exit(1)
	}

	// populate a new config object, load in the embedded config and set the
	// initial search paths to find other configs to overlay
	kfg.setDefaults(cfgPathOverride)

	// Merge the config file with defaults - ignore errors
	kfg.mergeAllConfigs()

	//iterate through the environment variables declared and bind them in kfg
	if envMap := kfg.StringMap("clog.env"); envMap != nil {
		for envIdentifier, envVariableName := range envMap {
			if false {
				fmt.Println(envIdentifier, envVariableName)
			}
			if envValue := os.Getenv(envVariableName); envValue != "" {
				kfg.Set(envVariableName, envValue)
			}
		}
	}
	return kfg
}

// get the config object
func Kfg() *Config {
	return kfg
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
