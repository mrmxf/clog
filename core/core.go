//  Copyright ©2017-2025  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License           https://opensource.org/license/bsd-3-clause/

// Package care manages the embedded files of clog
// if you are forking core, then you probably want to override this package
// to include only the files that you want in your own fork.

package core

import "embed"

//go:embed all:sample
//go:embed all:sh
//go:embed core.clog.yaml

var CoreFs embed.FS