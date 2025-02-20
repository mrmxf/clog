//  Copyright ©2019-2024  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License         https://opensource.org/license/bsd-3-clause/
//
// command:
//   $ slog.Check  # check the repo tags for production consistency

package checklegacy

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/mrmxf/clog/scripts"
	"github.com/spf13/cast"
)

type dependencyKeyProps struct {
	key    string
	value  string
	exists bool
	errMsg string
	err    error
}

// check that a dependency key exists
func setDependencyKeyProps(depMap map[string]string, key string, depId string) dependencyKeyProps {
	props := dependencyKeyProps{}
	props.key = key
	props.value, props.exists = depMap[key]
	if !props.exists {
		props.errMsg = fmt.Sprintf("%*v: missing key '%v'", colWidth-2, depId, props.key)
		props.err = errors.New(props.errMsg)
	}
	return props
}

// load the dependencies from the config and run the snippets
func checkDependencies(configCheckMap map[string]interface{}, colWidth int) error {
	depList, depKeyExists := configCheckMap["dependencies"]
	if !depKeyExists {
		// nothing to report
		return nil
	}

	// iterate through the report list and build a map
	checks := cast.ToSlice(depList)
	slog.Info(fmt.Sprintf("%s:  performing %v dependency checks", strings.Repeat("-", colWidth-1), len(checks)))

	var depError error = nil

	for d, check := range checks {
		depId := fmt.Sprintf("d%02d", d+1)

		dependency := cast.ToStringMapString(check)

		name := setDependencyKeyProps(dependency, "name", depId)
		if !name.exists {
			slog.Error(name.errMsg)
			depError = name.err
			continue
		}

		snippet := setDependencyKeyProps(dependency, "snippet", depId)
		if !snippet.exists {
			slog.Error(snippet.errMsg)
			depError = snippet.err
			continue
		}

		helpMsg := setDependencyKeyProps(dependency, "help-msg", depId)
		if !helpMsg.exists {
			slog.Error(helpMsg.errMsg)
			depError = helpMsg.err
			continue
		}

		warnMsg := setDependencyKeyProps(dependency, "warn-msg", depId)
		errorMsg := setDependencyKeyProps(dependency, "error-msg", depId)

		value, status, err := scripts.CaptureShellSnippet(snippet.value, nil)
		// reset the cache of returned values
		statusStr := fmt.Sprintf("%v", status)
		errStr := ""
		if err != nil {
			errStr = err.Error()
		}
		env := map[string]string{"VALUE": value, "STATUS": statusStr, "ERR": errStr}

		_, wIfExists := dependency["warn-if"]
		if wIfExists {
			if !warnMsg.exists {
				slog.Error(warnMsg.errMsg)
				depError = warnMsg.err
				continue
			}
			depMap := cast.ToStringMap(check)
			// we only warn if the status is 0, so we use Warn as the oslog.er
			tmpError := handleIf(depMap["warn-if"], colWidth, depId+".warn-if", env, logINFO, logWARN)
			if tmpError != nil {
				//ignore the error for the purpose of failing Checkm, but print help info
				lMsg := fmt.Sprintf("%*v: %v", colWidth, depId+".warn-if.msg", warnMsg.value)
				vMsg := fmt.Sprintf("%*v: %v", colWidth, depId+".warn-if.value", value)
				sMsg := fmt.Sprintf("%*v: %v", colWidth, depId+".warn-if.status", status)
				hMsg := fmt.Sprintf("%*v: %v", colWidth, depId+".warn-if.help", helpMsg.value)
				slog.Warn(lMsg)
				slog.Warn(vMsg)
				slog.Warn(sMsg)
				slog.Warn(hMsg)
			}
		}

		_, eIfExists := dependency["error-if"]
		if eIfExists {
			if !errorMsg.exists {
				slog.Error(errorMsg.errMsg)
				depError = errorMsg.err
				continue
			}
			depMap := cast.ToStringMap(check)
			// we only error if the status is 0, so we use Error as the oslog.er
			tmpError := handleIf(depMap["error-if"], colWidth, depId+".error-if", env, logERROR, logNONE)
			if tmpError != nil {
				depError = tmpError
				lMsg := fmt.Sprintf("%*v: %v", colWidth-2, depId+".error-if.msg", errorMsg.value)
				vMsg := fmt.Sprintf("%*v: %v", colWidth-2, depId+".error-if.value", value)
				sMsg := fmt.Sprintf("%*v: %v", colWidth, depId+".warn-if.status", status)
				hMsg := fmt.Sprintf("%*v: %v", colWidth-2, depId+".error-if.help", helpMsg.value)
				slog.Error(lMsg)
				slog.Error(vMsg)
				slog.Warn(sMsg)
				slog.Error(hMsg)
			}
		}
	}
	if depError != nil {
		depError = errors.New("dependency checks failed")
	}
	return depError
}
