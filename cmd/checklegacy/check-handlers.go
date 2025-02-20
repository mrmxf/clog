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
	"reflect"
	"regexp"
	"strings"

	"github.com/mrmxf/clog/scripts"
	"github.com/spf13/cast"
)

func checkCondition(condition string) (int, error) {
	switch {
	case condition == "==":
		return testEQ, nil
	case condition == "!=":
		return testNE, nil
	case condition == ">":
		return testGT, nil
	case condition == ">=":
		return testGE, nil
	case condition == "<":
		return testLT, nil
	case condition == "<=":
		return testLE, nil
	default:
		msg := fmt.Sprintf("unknown group test condition: %v", condition)
		return 0, errors.New(msg)
	}
}

// the report actions are slices of maps with a single key
// return that key (or an error len() is not 1)
func getActionKey(action map[string]interface{}, index int) (string, error) {
	if len(action) == 0 {
		msg := fmt.Sprintf("empty report action line %*v", colWidth-2, index)
		return "", errors.New(msg)
	}
	if len(action) > 1 {
		msg := fmt.Sprintf("too many keys (%v) in action line %*v", len(action), colWidth-2, index)
		return "", errors.New(msg)
	}
	key := ""
	for k := range action {
		key = k
	}
	return key, nil
}

// the report actions are slices of maps with a single key
// return that key (or an error len() is not 1)
func getSnippetActionKeyValue(action map[string]string, id string) (string, string, error) {
	if len(action) == 0 {
		msg := fmt.Sprintf("empty report action line %*v", colWidth-2, id)
		return "", "", errors.New(msg)
	}
	if len(action) > 1 {
		msg := fmt.Sprintf("too many keys (%v) in action line %*v", len(action), colWidth-2, id)
		return "", "", errors.New(msg)
	}
	key := ""
	for k := range action {
		key = k
	}
	return key, action[key], nil
}

// handle a snippet. okLog and errLog can be set to nil to suppress output
// otherwise check the return status and use the logger to report the value.
func handleSnippet(action map[string]string, colWidth int, id string, env map[string]string, okLog logAsType, errLog logAsType) error {
	actionKey, snippet, err := getSnippetActionKeyValue(action, id)
	if err != nil {
		return err
	}

	//if the key is a duplicate (e.g. a repeat in warn-if) then do a warning
	_, exists := vals[actionKey]
	if exists {
		slog.Warn(fmt.Sprintf("%*v: %v", colWidth-1, "duplicate Snippet key - overwriting", actionKey))
	}

	// if the key is empty then just log the value
	if len(snippet) == 0 {
		snippet = "echo \"$" + actionKey + "\""
	}

	value, status, err := scripts.CaptureShellSnippet(snippet, env)
	vals[actionKey] = value
	errs[actionKey] = err

	// we return an error:
	// - when errLog != logNONE and there is an error (a report error)
	// - when  okLog >= logWARN and there is no error (warn-if && error-if)
	var resErr error = nil
	logMsg := fmt.Sprintf("%*v: %v", colWidth-1, id+"."+actionKey, strings.TrimSpace(value))

	if err != nil || status > 0 {
		logAs(logMsg, errLog)
		if errLog != logNONE {
			resErr = errors.New(logMsg)
		}
	} else {
		logAs(logMsg, okLog)
		// during error-if we want to return an error if the condition was true
		if okLog >= logWARN {
			resErr = errors.New(logMsg)
		}
	}
	return resErr
}

// iterate through the array of snippets
// rawActions is created by the yaml reader as []interface{} (preserve yaml order)
// each element of the array is map[string]string
// each map has a single entry or an error will occur
func handleIf(rawActions interface{}, colWidth int, id string, env map[string]string, okLog logAsType, errLog logAsType) error {
	tRawActions := reflect.TypeOf(rawActions)
	if tRawActions != reflect.TypeOf([]interface{}{}) {
		logMsg := fmt.Sprintf("%*v: %v", colWidth-1, id, `Bad syntax, expected array: "- some-key: shellSnippetString"`)
		slog.Error(logMsg)
		return errors.New(logMsg)
	}

	var lastErr error = nil
	actionSlice := cast.ToSlice(rawActions)

	for k := range actionSlice {
		action := cast.ToStringMapString(actionSlice[k])
		lastErr = handleSnippet(action, colWidth, id, env, okLog, errLog)
	}
	return lastErr
}

// iterate through the array of group checks
func handleGroup(rawTests interface{}, colWidth int, id string, okLog logAsType, errLog logAsType) error {
	rawList := cast.ToStringSlice(rawTests)
	parseTest := regexp.MustCompile(`([^\s]+)\s+([^\s]+)\s+([^\s]+)`)
	parseString := regexp.MustCompile(`^"(.+)"$`)
	for testIndex, rawTest := range rawList {
		test := strings.TrimSpace(cast.ToString(rawTest))
		testId := fmt.Sprintf("%v.%v", id, testIndex)
		bits := parseTest.FindStringSubmatch(test)

		snippet1, exists := vals[bits[1]]
		if !exists {
			logMsg := fmt.Sprintf("%*v: %v", colWidth-1, testId, fmt.Sprintf("snippet %v not found", bits[1]))
			logAs(logMsg, errLog)
			vals[testId] = "error"
			errs[testId] = errors.New(logMsg)
			return errs[testId]
		}
		condition, err := checkCondition(bits[2])
		if err != nil {
			logMsg := fmt.Sprintf("%*v: %v", colWidth-1, testId, err.Error())
			slog.Error(logMsg)
			vals[testId] = "error"
			errs[testId] = errors.New(logMsg)
			return errs[testId]
		}

		testStrings := parseString.FindStringSubmatch(bits[3])
		var testString string
		if len(testStrings) > 0 {
			testString = testStrings[1]
		} else {
			testString, exists = vals[bits[3]]
			if !exists {
				logMsg := fmt.Sprintf("%*v: %v", colWidth-1, testId, fmt.Sprintf("snippet %v not found", bits[31]))
				logAs(logMsg, errLog)
				vals[testId] = "error"
				errs[testId] = errors.New(logMsg)
				return errs[testId]
			}
		}

		//check the values against the condition to issue a warning/error
		// if no warning/error to be issued then continue
		switch condition {
		case testEQ:
			if !(snippet1 == testString) {
				continue
			}
		case testNE:
			if !(snippet1 != testString) {
				continue
			}
		case testGT:
			if !(snippet1 > testString) {
				continue
			}
		case testGE:
			if !(snippet1 >= testString) {
				continue
			}
		case testLT:
			if !(snippet1 < testString) {
				continue
			}
		case testLE:
			if !(snippet1 <= testString) {
				continue
			}

		}
		logBody := fmt.Sprintf("failed test: %v(%v) %v %v(%v)", bits[1], snippet1, bits[2], bits[3], testString)
		logMsg := fmt.Sprintf("%*v: %v", colWidth-1, testId, logBody)
		logAs(logMsg, errLog)
		vals[testId] = "error"
		errs[testId] = errors.New(logMsg)
	}
	return nil
}
