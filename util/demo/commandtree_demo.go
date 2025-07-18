//  Copyright ©2017-2025  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License           https://opensource.org/license/bsd-3-clause/

package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mrmxf/clog/cmd"
	"github.com/mrmxf/clog/cmd/version"
	"github.com/mrmxf/clog/cmd/list"
	"github.com/mrmxf/clog/cmd/cat"
	"github.com/mrmxf/clog/util"
)

func CommandTreeDemo() {
	fmt.Println("Clog Command Tree Analysis")
	fmt.Println("==========================")

	// Create a copy of the root command to avoid Execute() being called
	rootCmd := cmd.RootCommand
	
	// Add some example subcommands to demonstrate
	rootCmd.AddCommand(version.Command)
	rootCmd.AddCommand(list.Command)
	rootCmd.AddCommand(cat.Command)
	
	// Generate the command tree
	tree := util.BuildCommandTree(rootCmd)

	// Print summary
	fmt.Printf("Root Command: %s\n", tree.Root.Command.Use)
	fmt.Printf("Description: %s\n", tree.Root.Command.Short)
	fmt.Printf("Flags: %d\n", len(tree.Root.Command.Flags))
	fmt.Printf("Commands: %d\n", len(tree.Commands))

	// Print all flags
	fmt.Printf("\nRoot Command Flags:\n")
	for _, flag := range tree.Root.Command.Flags {
		shorthand := ""
		if flag.Name != "" {
			shorthand = fmt.Sprintf(" (-%s)", flag.Name)
		}
		fmt.Printf("  --%s%s [%s]: %s\n", flag.Name, shorthand, flag.Schema.Type, flag.Description)
	}

	// Print commands
	fmt.Printf("\nCommands:\n")
	for _, cmd := range tree.Commands {
		fmt.Printf("  - %s: %s (Path: %s)\n", cmd.Use, cmd.Short, cmd.Path)
	}

	// Export as JSON
	jsonData, err := json.MarshalIndent(tree, "", "  ")
	if err != nil {
		fmt.Printf("Failed to marshal to JSON: %v\n", err)
		return
	}

	err = os.WriteFile("command_tree_demo.json", jsonData, 0644)
	if err != nil {
		fmt.Printf("Failed to write JSON file: %v\n", err)
		return
	}

	fmt.Printf("\nFull command tree exported to: command_tree_demo.json\n")
}

// CommandTreeDemo is exported for use by the demo orchestrator