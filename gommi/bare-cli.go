// Copyright ©2022-2025 Mr MXF   info@mrmxf.com
// BSD-3-Clause License   https://opensource.org/license/bsd-3-clause/

// package gommi provides a basic fileserver for use in a container that
// might be serving a static hugo site (for example). To aid a CLI, some
// flags are provided to control mounting of the file server to a given
// prefix and internal / external mount point
package gommi

import "flag"

var mountPath string
var mountPathFlag string = "path"
var mountPathDefault string = "/var/www"
var mountPathHelp string = "a path to mount the root file system, default=\"" + mountPathDefault + "\" (empty string==embedFS)"

var urlPrefix string
var urlPrefixFlag string = "prefix"
var urlPrefixDefault string = "/"
var urlPrefixHelp string = "a url prefix for the file server, default=\"" + urlPrefixDefault + "\""

func init() {
	flag.StringVar(&mountPath, mountPathFlag, mountPathDefault, mountPathHelp)
	flag.StringVar(&urlPrefix, urlPrefixFlag, urlPrefixDefault, urlPrefixHelp)
	flag.Parse()
}
