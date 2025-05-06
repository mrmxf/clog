// Copyright ©2017-2025 Mr MXF   info@mrmxf.com
// BSD-3-Clause License   https://opensource.org/license/bsd-3-clause/

package gommi

import (
	"embed"
	"html/template"
)

// Define some data template data structs
type TDErr struct {
	Title template.HTML
	Error template.HTML
}

type TDLanding struct {
	Title template.HTML
	Ptr   string // this is the path to the root folder from the route
}

// This struct is for managing all the templates used for rendering. A template
// is a set of named templates loaded as a group. The xxxSrc array holds the
// source filename in the embedded file system. the xxx template pointer holds
// the parsed group of templates. When executing, use the name "page" to ref
// the generic outer template
type DashTemplates struct {
	err        *template.Template
	errSrc     []string
	landing    *template.Template
	landingSrc []string
}

var tpl = DashTemplates{
	errSrc: []string{
		"templates/page.html",
		"templates/main-err.html",
	},
	landingSrc: []string{
		"templates/page.html",
		"templates/nav.html",
		"templates/main-landing.html",
	},
}

// iterate over the embedded source files and init the templates
func initTemplates(fs embed.FS) {
	// make a new logs template by parsing the collection of named templates
	tpl.err = template.Must(template.ParseFS(fs, tpl.errSrc...))
	tpl.landing = template.Must(template.ParseFS(fs, tpl.landingSrc...))
}
