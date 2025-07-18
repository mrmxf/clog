//  Copyright ©2017-2025  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License           https://opensource.org/license/bsd-3-clause/

//go:build example

package util

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/mrmxf/clog/cmd"
)

// ExampleUsage demonstrates how to use the CommandTree function
func ExampleUsage() {
	// Bootstrap the commands to get the full command tree
	err := cmd.BootStrap(cmd.RootCommand)
	if err != nil {
		fmt.Printf("Failed to bootstrap commands: %v\n", err)
		return
	}

	// Generate the command tree
	tree := CommandTree(cmd.RootCommand)

	// Print summary
	fmt.Printf("Command Tree Analysis\n")
	fmt.Printf("====================\n")
	fmt.Printf("Root Command: %s\n", tree.Use)
	fmt.Printf("Description: %s\n", tree.Short)
	fmt.Printf("Flags: %d\n", len(tree.Flags))
	fmt.Printf("Subcommands: %d\n", len(tree.Subcommands))

	// Print all flags
	fmt.Printf("\nFlags:\n")
	for _, flag := range tree.Flags {
		shorthand := ""
		if flag.Shorthand != "" {
			shorthand = fmt.Sprintf(" (-%s)", flag.Shorthand)
		}
		fmt.Printf("  --%s%s [%s]: %s\n", flag.Name, shorthand, flag.DataType, flag.Usage)
	}

	// Print subcommands recursively
	fmt.Printf("\nSubcommands:\n")
	printSubcommands(tree.Subcommands, 0)

	// Export as JSON
	jsonData, err := json.MarshalIndent(tree, "", "  ")
	if err != nil {
		fmt.Printf("Failed to marshal to JSON: %v\n", err)
		return
	}

	err = os.WriteFile("full_command_tree.json", jsonData, 0644)
	if err != nil {
		fmt.Printf("Failed to write JSON file: %v\n", err)
		return
	}

	fmt.Printf("\nFull command tree exported to: full_command_tree.json\n")
}

// printSubcommands recursively prints subcommands with indentation
func printSubcommands(commands []*CommandInfo, depth int) {
	indent := ""
	for i := 0; i < depth; i++ {
		indent += "  "
	}

	for _, cmd := range commands {
		fmt.Printf("%s- %s: %s\n", indent, cmd.Use, cmd.Short)
		
		// Print flags for this subcommand
		if len(cmd.Flags) > 0 {
			fmt.Printf("%s  Flags:\n", indent)
			for _, flag := range cmd.Flags {
				shorthand := ""
				if flag.Shorthand != "" {
					shorthand = fmt.Sprintf(" (-%s)", flag.Shorthand)
				}
				fmt.Printf("%s    --%s%s [%s]: %s\n", indent, flag.Name, shorthand, flag.DataType, flag.Usage)
			}
		}

		// Print arguments for this subcommand
		if len(cmd.Args) > 0 {
			fmt.Printf("%s  Arguments:\n", indent)
			for _, arg := range cmd.Args {
				required := ""
				if arg.Required {
					required = " (required)"
				}
				fmt.Printf("%s    %s%s: %s\n", indent, arg.Name, required, arg.Description)
			}
		}

		// Recursively print subcommands
		if len(cmd.Subcommands) > 0 {
			printSubcommands(cmd.Subcommands, depth+1)
		}
	}
}