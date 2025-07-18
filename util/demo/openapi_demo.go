//  Copyright ©2017-2025  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License           https://opensource.org/license/bsd-3-clause/

package main

import (
	"fmt"
	"os"

	"github.com/mrmxf/clog/cmd"
	"github.com/mrmxf/clog/cmd/cat"
	"github.com/mrmxf/clog/cmd/list"
	"github.com/mrmxf/clog/cmd/version"
	"github.com/mrmxf/clog/util"
)

func OpenAPIDemo() {
	fmt.Println("Clog OpenAPI 3.0 Command Tree Analysis")
	fmt.Println("=======================================")

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
	fmt.Printf("Root Path: %s\n", tree.Root.Command.Path)
	fmt.Printf("Total Commands: %d\n", len(tree.Commands))

	// Print all commands with their API properties
	fmt.Printf("\nCommands and their API endpoints:\n")
	for _, cmd := range tree.Commands {
		fmt.Printf("  Command: %s\n", cmd.Use)
		fmt.Printf("    Path: %s\n", cmd.Path)
		fmt.Printf("    Short: %s\n", cmd.Short)
		fmt.Printf("    Flags: %d\n", len(cmd.Flags))
		fmt.Printf("    Args: %d\n", len(cmd.Args))
		fmt.Printf("    Has Args: %t\n", cmd.HasArgs)
		
		// Show available endpoints
		if cmd.Get != nil {
			fmt.Printf("    GET endpoint: %s\n", cmd.Path)
		}
		if cmd.Post != nil {
			fmt.Printf("    POST endpoint: %s (for args)\n", cmd.Path)
		}
		
		// Show some flags as query parameters
		if len(cmd.Flags) > 0 {
			fmt.Printf("    Query Parameters:\n")
			for i, flag := range cmd.Flags {
				if i >= 3 { // Show only first 3 flags
					fmt.Printf("      ... and %d more\n", len(cmd.Flags)-3)
					break
				}
				fmt.Printf("      %s (%s): %s\n", flag.Name, flag.Schema.Type, flag.Description)
			}
		}
		
		// Show args
		if len(cmd.Args) > 0 {
			fmt.Printf("    Arguments:\n")
			for _, arg := range cmd.Args {
				required := ""
				if arg.Required {
					required = " (required)"
				}
				fmt.Printf("      %s (%s)%s: %s\n", arg.Name, arg.Type, required, arg.Description)
			}
		}
		fmt.Println()
	}

	// Generate OpenAPI 3.0 specification
	fmt.Printf("Generating OpenAPI 3.0 specification...\n")
	spec := tree.GenerateOpenAPISpec("Clog API", "Command Line Of Go REST API", "1.0.0")
	
	// Export specifications
	err := spec.ExportToJSON("clog_openapi.json")
	if err != nil {
		fmt.Printf("Error exporting OpenAPI JSON: %v\n", err)
	} else {
		fmt.Printf("OpenAPI JSON exported to: clog_openapi.json\n")
	}
	
	err = spec.ExportToYAML("clog_openapi.yaml")
	if err != nil {
		fmt.Printf("Error exporting OpenAPI YAML: %v\n", err)
	} else {
		fmt.Printf("OpenAPI YAML exported to: clog_openapi.yaml\n")
	}

	// Export command tree in both formats
	err = tree.ExportToJSON("clog_command_tree.json")
	if err != nil {
		fmt.Printf("Error exporting command tree JSON: %v\n", err)
	} else {
		fmt.Printf("Command tree JSON exported to: clog_command_tree.json\n")
	}
	
	err = tree.ExportToYAML("clog_command_tree.yaml")
	if err != nil {
		fmt.Printf("Error exporting command tree YAML: %v\n", err)
	} else {
		fmt.Printf("Command tree YAML exported to: clog_command_tree.yaml\n")
	}

	// Create templates directory
	err = os.MkdirAll("templates", 0755)
	if err != nil {
		fmt.Printf("Error creating templates directory: %v\n", err)
		return
	}

	// Generate API templates for commands with args
	fmt.Printf("Generating API templates for commands with arguments...\n")
	err = tree.GenerateCommandTemplates("templates")
	if err != nil {
		fmt.Printf("Error generating command templates: %v\n", err)
	} else {
		fmt.Printf("API templates generated in: templates/\n")
	}

	// Summary
	fmt.Printf("\nSummary:\n")
	fmt.Printf("========\n")
	fmt.Printf("- OpenAPI 3.0 specification: clog_openapi.json, clog_openapi.yaml\n")
	fmt.Printf("- Command tree structure: clog_command_tree.json, clog_command_tree.yaml\n")
	fmt.Printf("- API templates: templates/*-api-template.json\n")
	fmt.Printf("- Total API endpoints: %d\n", countEndpoints(tree))
	fmt.Printf("- Commands with POST endpoints: %d\n", countCommandsWithArgs(tree))
}

func countEndpoints(tree *util.CommandTree) int {
	count := 0
	for _, cmd := range tree.Commands {
		if cmd.Get != nil {
			count++
		}
		if cmd.Post != nil {
			count++
		}
	}
	return count
}

func countCommandsWithArgs(tree *util.CommandTree) int {
	count := 0
	for _, cmd := range tree.Commands {
		if cmd.HasArgs {
			count++
		}
	}
	return count
}