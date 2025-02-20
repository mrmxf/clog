//  Copyright ©2019-2024  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License         https://opensource.org/license/bsd-3-clause/
//
// command:
//   $ clog Check  # check the repo tags for production consistency

package checklegacy

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/spf13/cast"
)

// checkReport runs snippets and reports on them
// see the annotated clog.core.config.yaml for examples.
//
// report is an array. Each element is a map[string]interface
// current supported actions are:
// ```yaml
// - NAME:                      # echo the env variable NAME
// - NAME: snippet              # run the shell snippet, report the value
// - warn-if:
//   - NAME1: snippet           # run snippet NAME1 report the value, WARN if error
//   - NAME2: snippet           # run snippet NAME2 report the value, WARN if error
//
// - warn-group-GNAME:          # runs a conditional group of tests entitled GNAME
//   - NAME1 cond NAME2         # WARN if test true [ cond: ==, !=, >, <, >+, <= ]
//   - NAME1 cond "string"      # WARN if test true [ cond: ==, !=, >, <, >+, <= ]
//
// - error-if                     # ERROR version of warn-if
// - error-group-GNAME            # ERROR version of warn-group-GNAME
// ```
func checkReport(configCheckMap map[string]interface{}, colWidth int) error {
	reportList, reportKeyExists := configCheckMap["report"]
	if !reportKeyExists {
		// nothing to report
		return nil
	}

	// iterate through the report list and build a map
	actions := cast.ToSlice(reportList)
	slog.Info(fmt.Sprintf("%s:  performing %v report actions", strings.Repeat("-", colWidth-1), len(actions)))

	// zero the cache of returned values
	vals = map[string]string{}
	errs = map[string]error{}
	var lastErr error = nil

	//iterate over all the report actions
	for index, action := range actions {
		actionMap := cast.ToStringMap(action)
		id := fmt.Sprintf("a%02d", index)

		key, err := getActionKey(actionMap, index)
		if err != nil {
			return err
		}

		switch {
		case key == "warn-if":
			// warnings are not errors so ignore them
			_ = handleIf(actionMap[key], colWidth, id+".warn", nil, logNONE, logWARN)
		case key == "error-if":
			lastErr = handleIf(actionMap[key], colWidth, id+".error", nil, logNONE, logERROR)
		case strings.HasPrefix(key, "warn-group-"):
			prefix, _ := strings.CutPrefix(key, "warn-group-")
			id = fmt.Sprintf("%v.%v", id, prefix)
			// warnings are not errors so ignore them
			_ = handleGroup(actionMap[key], colWidth, id, logNONE, logWARN)
		case strings.HasPrefix(key, "error-group-"):
			prefix, _ := strings.CutPrefix(key, "error-group-")
			id = fmt.Sprintf("%v.%v", id, prefix)
			lastErr = handleGroup(actionMap[key], colWidth, id, logNONE, logERROR)
		default:
			action := cast.ToStringMapString(actionMap)
			//ignore warnings and errors
			lastErr = handleSnippet(action, colWidth, id, nil, logINFO, logNONE)
		}
	}
	if lastErr == nil {
		return nil
	}
	return errors.New("report checks failed")
}
