//  Copyright ©2019-2024  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License  https://opensource.org/license/bsd-3-clause/
//
// command:
//   $ clog <shell snippet from config>

package snippets

import (
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"sort"

	"github.com/mrmxf/clog/config"
	"github.com/mrmxf/clog/scripts"
	"github.com/spf13/cobra"
)

// process snippets

const shellSnippetsKey = "snippets.Sh"
const snippetsKey = "snippets"
const runSnippetString = "Sh"

var verboseListing bool

func reportSnippets(cmd *cobra.Command, title string, key string, cmdPath string, snippets map[string]string) {
	fmt.Println(">>>" + title + " in config key `" + key + "`")
	if len(snippets) == 0 {
		slog.Warn("No " + title + " found in config key `" + key + "`")
		return
	}

	keys := make([]string, 0, len(snippets))
	for k := range snippets {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, k := range keys {
		// check if the key has children
		snippet := snippets[k]
		if len(snippet) == 0 {
			childKey := key + "." + k
			children := config.Cfg().GetStringMapString(childKey)
			if len(children) > 0 {
				newCmdPath := "  " + cmdPath + k + " "
				reportSnippets(cmd, title, childKey, newCmdPath, children)
			} else {
				slog.Error("no commands for snippet (" + childKey + ") in clog.config.yaml")
			}
		} else {
			if verboseListing {
				fmt.Println("  " + cmdPath + k + ":   " + snippet)
			} else {
				fmt.Println("  " + cmdPath + k)
			}
		}
	}

}

func OldreportSnippets(name string, key string, cmdStr string, snippets map[string]string) {
	fmt.Println(name + " in config key `" + key + "`")
	if len(snippets) == 0 {
		slog.Warn("No " + name + " found in config key `" + key + "`")
	} else {
		keys := make([]string, 0, len(snippets))
		for k := range snippets {
			keys = append(keys, k)
		}
		sort.Strings(keys)

		for _, k := range keys {
			if verboseListing {
				fmt.Println("  " + cmdStr + k + ":   " + snippets[k])
			} else {
				fmt.Println("  " + cmdStr + k)
			}
		}
	}
}

var OldListSnippetsCmd = listSnippetsCommand
var OldlistSnippetsCommand = &cobra.Command{
	Use:   "Snippets",
	Short: "list snippets found in the config keys " + shellSnippetsKey + ", " + snippetsKey,
	Long:  `local config adds & overwrites the core snippets`,

	Run: func(cmd *cobra.Command, args []string) {
		doCmdStr := cmd.Parent().CommandPath() + " "
		snippets := config.Cfg().GetStringMapString(snippetsKey)
		OldreportSnippets("Command Snippets", snippetsKey, doCmdStr, snippets)

		doCmdStr = cmd.Parent().CommandPath() + " " + runSnippetString + " "
		snippets = config.Cfg().GetStringMapString(shellSnippetsKey)
		OldreportSnippets("Shell Snippets", shellSnippetsKey, doCmdStr, snippets)

		if !verboseListing {
			fmt.Println("\nuse clog Snippets --show to show full shell snippet stings")
		}
	},
}

// var ShCmd = shellSnippetCommand
// var shellSnippetCommand = &cobra.Command{
// 	Use:   runSnippetString,
// 	Short: "run a shell snippet found in the config key `" + shellSnippetsKey + "`",
// 	Long:  `local config adds & overwrites the core snippets`,

// 	Run: func(cmd *cobra.Command, args []string) {
// 		log := slogger.GetLogger()
// 		snippet := config.Cfg().GetString(shellSnippetsKey + "." + args[0])

//			if len(snippet) == 0 {
//				slog.Error("snippet (" + args[0] + ") not found in clog.config.yaml")
//				os.Exit(1)
//			}
//			slog.Debug("running shell snippet: "+args[0], "bash", snippet)
//			scripts.ShellSnippet(snippet)
//		},
//	}
var ListSnippetsCmd = listSnippetsCommand
var listSnippetsCommand = &cobra.Command{
	Use:   "Snippets",
	Short: "list snippets found in the config key " + snippetsKey,
	Long:  `local config adds & overwrites the core snippets`,

	Run: func(cmd *cobra.Command, args []string) {
		root := cmd.Root()
		rootCmdPath := root.CommandPath() + " "
		snippets := config.Cfg().GetStringMapString(snippetsKey)
		reportSnippets(root, "Snippets", snippetsKey, rootCmdPath, snippets)

		if !verboseListing {
			fmt.Println("\nclog Snippets --show   # show full shell snippet strings")
		}
	},
}

// add command snippets to the main list of root commands
func AddCommandSnippets(rootCmd *cobra.Command) {
	snippets := config.Cfg().GetStringMapString(snippetsKey)
	for k, snippet := range snippets {
		cmd := new(cobra.Command)
		cmd.Use = k
		cmd.Short = "command snippet " + k
		if len(snippet) == 0 {
			slog.Error("no commands for snippet (" + k + ") in clog.config.yaml")
			os.Exit(1)
		}
		cmd.Run = func(cmd *cobra.Command, args []string) {
			slog.Debug("running command snippet: "+k, "bash", snippet)
			result := scripts.StreamShellSnippet(snippet, nil)
			//return the status of the command
			os.Exit(result)
		}
		rootCmd.AddCommand(cmd)
	}
}

// add snippets to the main list of root commands, found with given key
func AddSnippets(rootCmd *cobra.Command, key string) {
	snippets := config.Cfg().GetStringMapString(key)
	for k, snippet := range snippets {
		cmd := new(cobra.Command)
		cmd.Use = k
		cmd.Short = "snippet " + k
		if len(snippet) == 0 {
			// check if the key has children
			newKey := key + "." + k
			children := config.Cfg().GetStringMapString(newKey)
			if len(children) > 0 {
				AddSnippets(cmd, newKey)
			} else {
				slog.Error("no commands for snippet (" + newKey + ") in clog.config.yaml")
				os.Exit(1)
			}
		}
		cmd.RunE = func(cmd *cobra.Command, args []string) error {
			slog.Debug("running command snippet: "+k, "bash", snippet)
			exitStatus := scripts.StreamShellSnippet(snippet, nil)
			// if exitStatus != 0 {
			// 	// prevent Cobra from displaying command help on error
			// 	cmd.SilenceErrors = true
			// 	cmd.SilenceUsage = true
			// 	return fmt.Errorf("%d", exitStatus)
			// }
			// return nil
			os.Exit(exitStatus)
			return nil
		}
		rootCmd.AddCommand(cmd)
	}
}

func init() {
	// log the order of the init files in case there are problems
	_, file, _, _ := runtime.Caller(0)
	slog.Debug("init " + file)

	ListSnippetsCmd.PersistentFlags().BoolVarP(&verboseListing, "show", "s", false, "clog Snippets -s       # show script strings")
}
