//  Copyright ©2017-2025  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License  https://opensource.org/license/bsd-3-clause/
//
// clog's config package

package config

import (
	"log/slog"
	"os"
	"path/filepath"
	"runtime"
	"strings"
)

// mergeAllConfigs will search for the default override config files once
// the embedded config file has been loaded. When THE FINAL search location
// has been searched, it will then see if SearchPathList has changed. If so,
// it will call itself again until there is no more change to SearchPathList.
// Note that the max depth will be 4
func (cfg *Config) mergeAllConfigs() {
	slog.Debug("Merging user defined configs", "SearchPathList", searchPaths)

	fName := cfg.GetString("clog.clogrc.base") + "." + cfg.GetString("clog.clogrc.format")
	for _, path := range searchPaths {
		fPath := strings.Replace(path, "~", os.Getenv("HOME"), 1)
		fPath = strings.Replace(fPath, "$HOME", os.Getenv("HOME"), 1)
		if !(strings.HasSuffix(fPath, ".yaml") || strings.HasSuffix(fPath, ".json")) {
			fPath = filepath.Join(fPath, fName)
		}
		fPathAbs, err := filepath.Abs(fPath)
		if err != nil {
			slog.Debug("Error getting absolute path", "path", fPath, "error", err)
			continue
		}
		ioReader, err := os.Open(fPathAbs)
		if err == nil {
			slog.Debug("Found config file", "path", fPathAbs)
			err := cfg.MergeConfig(ioReader)
			if err != nil {
				slog.Error("Error merging config file", "path", fPathAbs, "error", err)
			}
			ioReader.Close()
		} else {
			slog.Debug("Did not find config file", "path", fPathAbs)
		}
	}
}

func init() {
	// log the order of the init files in case there are problems
	_, file, _, _ := runtime.Caller(0)
	slog.Debug("init " + file)
}
