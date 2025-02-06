//  Copyright ©2017-2025  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License  https://opensource.org/license/bsd-3-clause/
//
// package snips provide handling functions to enable snippets

package snips

import (
	"fmt"
	"log/slog"
	"os"
	"runtime"

	"github.com/mrmxf/clog/scripts"
	"github.com/spf13/cobra"
)

// create helper types
type Snippet string
type SnippetGroup map[Snippet]any
type RawSnippets map[string]any

// a parsed snippets struct
type ParsedSnippets struct {
	Snippets  SnippetGroup
	ParentCmd *cobra.Command // the parent command of the snippet group
}

type ListSnippetsData struct {
	Title   string
	Key     string
	Parsed  *ParsedSnippets
	Verbose bool
	Plain   bool
}

// add snippets to the main list of root commands, found with given key
func ParseSnippets(rootCmd *cobra.Command, raw RawSnippets) (ParsedSnippets, error) {
	pSnips := ParsedSnippets{
		Snippets:  SnippetGroup{},
		ParentCmd: rootCmd,
	}

	recurseRawMap(pSnips.ParentCmd, pSnips.Snippets, 0, raw)

	return pSnips, nil
}

func kmdPath(c *cobra.Command, kmd string) string{
	return fmt.Sprintf("%s %s",c.CommandPath(), kmd)
}

func recurseRawMap(parentCmd *cobra.Command, group SnippetGroup, depth int, raw RawSnippets) {
	for kmd, snip := range raw {
		switch skript := snip.(type) {

		case int:
			// an int is a valid ZSH / BASH command !!
			slog.Debug(fmt.Sprintf("%d.leaf - %s %T", depth, kmd, skript))
			cmd := &cobra.Command{
				Use:   kmd,
				Short: kmdPath(parentCmd, kmd),
				Run: func(cmd *cobra.Command, args []string) {
					strInt := fmt.Sprintf("%d", skript)
					slog.Debug(fmt.Sprintf("snippet: %s\n$ %s\n", kmd, strInt))
					exitStatus := scripts.StreamShellSnippet(strInt, nil)
					os.Exit(exitStatus)
				},
			}
			parentCmd.AddCommand(cmd)
			group[Snippet(kmd)] = skript

		case string:
			slog.Debug(fmt.Sprintf("%d.leaf - %s %T", depth, kmd, skript))
			cmd := &cobra.Command{
				Use:   kmd,
				Short: "snippet " + kmdPath(parentCmd, kmd),
				Run: func(cmd *cobra.Command, args []string) {
					slog.Debug(fmt.Sprintf("snippet: %s\n$ %s\n", kmd, skript))
					exitStatus := scripts.StreamShellSnippet(skript, nil)
					os.Exit(exitStatus)
				},
			}
			parentCmd.AddCommand(cmd)
			group[Snippet(kmd)] = skript

		case map[string]interface{}:
			slog.Debug(fmt.Sprintf("%d.node - %s %T", depth, kmd, snip))
			cmd := &cobra.Command{
				Use:   kmd,
				Short: kmdPath(parentCmd, kmd),
				Run: func(cmd *cobra.Command, args []string) {
            cmd.Help()
					os.Exit(1)
				},
			}
			// add this command stub to the tree and descend
			parentCmd.AddCommand(cmd)
			newGroup := SnippetGroup{}
			group[Snippet(kmd)] = &newGroup
			recurseRawMap(cmd, newGroup, depth+1, snip.(map[string]interface{}))
		default:
			slog.Debug(fmt.Sprintf("%d.WARNING ignoring unexpected snippet (%s) of type %s", depth, kmd, snip))
		}

	}
}

func init() {
	// log the order of the init files in case there are problems
	_, file, _, _ := runtime.Caller(0)
	slog.Debug("init " + file)
}
