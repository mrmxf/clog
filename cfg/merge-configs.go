//  Copyright ©2017-2025  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License  https://opensource.org/license/bsd-3-clause/

// clog's cfg package

package cfg

import (
	"log/slog"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

// ExpandEnvVars will replace all elements like `$ENV_VAR` and `$THINGY` into
// the values taken from the environment.
// expandedStr will have all the valid environment variables expanded
// allValid will be false if any of the variables are zero length
func ExpandEnvVars(input string) (expandedStr string, allValid bool) {
	// Regex to match $VARIABLE_NAME where VARIABLE_NAME is valid env var format
	re := regexp.MustCompile(`\$([A-Z_][A-Z0-9_]*)`)

	allValid = true

	expandedStr = re.ReplaceAllStringFunc(input, func(match string) string {
		// Extract variable name (remove the $ prefix)
		varName := match[1:]

		// Get the environment variable value
		value := os.Getenv(varName)

		// If any replacement is empty, set allValid to false
		if value == "" {
			allValid = false
		}

		return value
	})

	return expandedStr, allValid
}

// mergeAllConfigs will search for the default override config files once
// the embedded config file has been loaded. When THE FINAL search location
// has been searched, it will then see if SearchPathList has changed. If so,
// it will call itself again until there is no more change to SearchPathList.
// Note that the max depth will be 4
func (kfg *Config) mergeAllConfigs() {
	slog.Debug("Merging user defined configs", "SearchPathList", searchPaths)

	for _, rawPath := range searchPaths {
		relPath := strings.Replace(rawPath, "~", "$HOME", 1)
		relPath, _ = ExpandEnvVars(relPath)
		path, err := filepath.Abs(relPath)
		if err != nil {
			slog.Debug("Error getting absolute path", "path", relPath, "error", err)
			continue
		}

		// Check if file exists
		if _, err := os.Stat(path); err == nil {
			slog.Debug("Found config file", "path", path)
			err := kfg.Load(file.Provider(path), yaml.Parser())
			if err != nil {
				slog.Error("Error merging config file", "path", path, "error", err)
			}
		} else {
			slog.Debug("Did not find config file", "path", path)
		}
	}
}

func init() {
	// log the order of the init files in case there are problems
	_, file, _, _ := runtime.Caller(0)
	slog.Debug("init " + file)
}
