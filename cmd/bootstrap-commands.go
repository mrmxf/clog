//  Copyright ©2017-2025    Mr MXF   info@mrmxf.com
//  BSD-3-Clause License    https://opensource.org/license/bsd-3-clause/
//
// package cmd contains the default commands in a form that can be individually
// loaded by a fork of clog.

package cmd

import (
	"log/slog"
	"runtime"
)

func init() {
	_, file, _, _ := runtime.Caller(0)
	slog.Debug("init " + file)
}
