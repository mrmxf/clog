//  Copyright ©2019-2024  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License         https://opensource.org/license/bsd-3-clause/
//
// command:
//   $ clog Check        # run checks from the `default:` key in the config
//   $ clog Check  thing # run checks from the `thing:` key in the config

package checklegacy

import (
	"errors"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"

	"github.com/mrmxf/clog/config"
	"github.com/spf13/cobra"
)

// make a type alias for the log function like log.Error
// type logFuncType = func(msg string, args ...any)

// define the default column width for readability
var colWidth = 30

// define the stores for each section of the check
var vals = map[string]string{}
var errs = map[string]error{}

// define the enum log levels
type logAsType int

const (
	logNONE  logAsType = iota
	logINFO  logAsType = iota
	logWARN  logAsType = iota
	logERROR logAsType = iota
)

func logAs(msg string, level logAsType) {
	switch level {
	case logINFO:
		slog.Info(msg)
	case logWARN:
		slog.Warn(msg)
	case logERROR:
		slog.Error(msg)
	}
}

// define the enum test conditions

// define the conditions for matching strings / action keys
const (
	testEQ = iota
	testNE
	testGT
	testGE
	testLT
	testLE
)

var Command = &cobra.Command{
	Use:   "CheckLegacy",
	Short: "run checks defined in the check:thing section of config",
	Long:  `returns error = number of fatal issues found.`,

	RunE: func(cmd *cobra.Command, args []string) error {
		// return success if there is nothing to check
		if config.Cfg().Get("check") == nil {
			return nil
		}
		defaultKeys := []string{"default", "default-core"}
		var keyToCheck string
		for _, keyToCheck = range defaultKeys {
			if config.Cfg().Get("check."+keyToCheck) != nil {
				break
			}
		}
		if len(keyToCheck) == 0 {
			slog.Error("no default or default-core key found in clog.config.yaml")
			os.Exit(1)
		}
		checkFailed := errors.New("clog Check failed")
		if len(args) > 0 {
			keyToCheck = args[0]
			cfgPath := "check." + keyToCheck
			exists := config.Cfg().Get(cfgPath)
			if exists == nil {
				slog.Error("no key found for " + cfgPath + " in clog.config.yaml")
				os.Exit(1)
			}
			checkFailed = errors.New("clog Check " + keyToCheck + " failed")
		}

		configCheckMap := config.Cfg().GetStringMap("check." + keyToCheck)
		objErr := []error{}

		// report that we did something
		slog.Info(fmt.Sprintf("%s:  starting clog Check %v", strings.Repeat("=", colWidth-1), keyToCheck))

		//iterate over Git report section (if it exists)
		err := checkReport(configCheckMap, colWidth)
		if err != nil {
			objErr = append(objErr, err)
		}

		//iterate over the dependencies in the clog.config
		err = checkDependencies(configCheckMap, colWidth)
		if err != nil {
			objErr = append(objErr, err)
		}

		// check the error objects for failure. The overall check will fail when
		// there are errors from any of the modules
		if len(objErr) > 0 {
			objErr = append(objErr, checkFailed)
			for _, e := range objErr {
				slog.Error(e.Error())
			}
			slog.Info(fmt.Sprintf("%s:  completed clog Check %v", strings.Repeat("=", colWidth-1), keyToCheck))
			os.Exit(1)
		}
		slog.Info(fmt.Sprintf("%s:  completed clog Check %v", strings.Repeat("=", colWidth-1), keyToCheck))
		return nil
	},
}

func init() {
	_, file, _, _ := runtime.Caller(0)
	slog.Debug("init " + file)
}
