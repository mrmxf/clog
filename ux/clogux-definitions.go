// Package to provide types for clog's ux (User eXperience)
//
// Exposes a command line interface (cli) and a web user interface (wui)

package ux

import (
	"github.com/spf13/cobra"
)

const AliasVocabulary = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@"

type uiMode int

// create integer constants for the uiMode
const (
	CLI uiMode = iota
	WEB
)

// wrapper from cobra commands providing a single letter short code for
// each unique child of a parent
type CmdChild struct {
	Name    string
	Short   string
	Alias   string
	Command *cobra.Command
}

// stats for each of the menu pages to help with formatting
type CmdChildStats struct {
	MaxNameLen  int
	MaxShortLen int
	AliasList   string
}

// MenuForm struct represents a level in a set of hierarchical menus
type MenuForm struct {
	Parent   *MenuForm
	Cmd      *cobra.Command
	Name     string
	Short    string
	Key      string
	Children []*MenuForm
	Id       int
}

// Option struct represents a simple set of choices.
// each choice is represented by a key and a caption
type Option struct {
	Key     string
	Caption string
}

// ActionForm struct represents a simple set of options to be performed
// after displaying a set of strings to the user
type ActionFormData struct {
	Title       string
	Description []string
	Options     []Option
}
