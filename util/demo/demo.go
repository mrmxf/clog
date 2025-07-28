//  Copyright ©2017-2025  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License           https://opensource.org/license/bsd-3-clause/

package main

import (
	"fmt"
	"os"
)

// DemoFunc represents a demo function signature
type DemoFunc func()

// DemoCommand represents a demo command configuration
type DemoCommand struct {
	Arg      string
	Function DemoFunc
	Help     string
}

// demoCommands defines all available demo commands
var demoCommands = []DemoCommand{
	{
		Arg:      "openapi",
		Function: OpenAPIDemo,
		Help:     "Generate OpenAPI 3.0 specification and command tree (machine-readable)",
	},
	{
		Arg:      "tree",
		Function: CommandTreeDemo,
		Help:     "Display command tree structure in human-readable format",
	},
}

// printHelp displays usage information
func printHelp() {
	fmt.Println("Usage: demo [COMMAND] [FLAGS]")
	fmt.Println()
	fmt.Println("Commands:")
	for _, cmd := range demoCommands {
		fmt.Printf("  %-10s %s\n", cmd.Arg, cmd.Help)
	}
	fmt.Println()
	fmt.Println("If no command is provided, 'openapi' is used by default.")
	fmt.Println()
	fmt.Println("Flags:")
	fmt.Println("  -H, --help     Show this help message")
}

// findCommand finds a demo command by argument
func findCommand(arg string) *DemoCommand {
	for i := range demoCommands {
		if demoCommands[i].Arg == arg {
			return &demoCommands[i]
		}
	}
	return nil
}

// parseArgs parses command line arguments using only Go native functionality
func parseArgs(args []string) (command string, showHelp bool) {
	// Default command
	command = "openapi"
	showHelp = false

	for i, arg := range args {
		switch arg {
		case "-H", "--help":
			showHelp = true
			return
		default:
			// First non-flag argument is the command
			if i == 0 {
				command = arg
			}
		}
	}

	return
}

func main() {
	// Get command line arguments (excluding program name)
	args := os.Args[1:]

	// Parse arguments
	command, showHelp := parseArgs(args)

	// Show help if requested
	if showHelp {
		printHelp()
		return
	}

	// Find and execute the command
	cmd := findCommand(command)
	if cmd == nil {
		fmt.Fprintf(os.Stderr, "Error: Unknown command '%s'\n", command)
		fmt.Fprintf(os.Stderr, "Use 'demo --help' to see available commands.\n")
		os.Exit(1)
	}

	// Execute the demo function
	cmd.Function()
}
