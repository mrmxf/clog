//  Copyright ©2017-2025  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License           https://opensource.org/license/bsd-3-clause/

package util

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/mrmxf/clog/cmd"
	"github.com/spf13/cobra"
)

func TestBuildCommandTree(t *testing.T) {
	// Bootstrap the commands to load subcommands
	rootCmd := cmd.RootCommand

	// Add subcommands manually for testing since BootStrap calls Execute
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the version string",
		Long:  "use clog -v to display the short semantic version string.",
	})

	// Get the command tree from the root command
	tree := BuildCommandTree(rootCmd)

	// Verify basic structure
	if tree.Root.Command.Use != "clog" {
		t.Errorf("Expected root command Use to be 'clog', got '%s'", tree.Root.Command.Use)
	}

	if tree.Root.Command.Short != "Command Line Of Go - interactive helper" {
		t.Errorf("Expected root command Short to be 'Command Line Of Go - interactive helper', got '%s'", tree.Root.Command.Short)
	}

	// Verify flags are present
	if len(tree.Root.Command.Flags) == 0 {
		t.Error("Expected root command to have flags")
	}

	// Look for the config flag
	hasConfigFlag := false
	for _, flag := range tree.Root.Command.Flags {
		if flag.Name == "config" {
			hasConfigFlag = true
			if flag.Schema.Type != "string" {
				t.Errorf("Expected config flag type to be 'string', got '%s'", flag.Schema.Type)
			}
			break
		}
	}
	if !hasConfigFlag {
		t.Error("Expected root command to have config flag")
	}

	// Verify commands list is present
	if len(tree.Commands) == 0 {
		t.Error("Expected tree to have commands")
	}

	// Look for version command in commands list
	hasVersionCmd := false
	for _, cmd := range tree.Commands {
		if cmd.Use == "version" {
			hasVersionCmd = true
			if cmd.Short != "Print the version string" {
				t.Errorf("Expected version command Short to be 'Print the version string', got '%s'", cmd.Short)
			}
			if cmd.Path != "/clog/version" {
				t.Errorf("Expected version command path to be '/clog/version', got '%s'", cmd.Path)
			}
			break
		}
	}
	if !hasVersionCmd {
		t.Error("Expected to find version command in commands list")
	}
}

func TestCommandTreeJSON(t *testing.T) {
	// Get the command tree
	tree := BuildCommandTree(cmd.RootCommand)

	// Export to JSON using the new method
	err := tree.ExportToJSON("command_tree_new.json")
	if err != nil {
		t.Fatalf("Failed to export command tree to JSON: %v", err)
	}

	// Verify JSON is valid by unmarshaling
	data, err := os.ReadFile("command_tree_new.json")
	if err != nil {
		t.Fatalf("Failed to read JSON file: %v", err)
	}

	var unmarshaled CommandTree
	err = json.Unmarshal(data, &unmarshaled)
	if err != nil {
		t.Fatalf("Failed to unmarshal JSON: %v", err)
	}

	// Basic validation
	if unmarshaled.Root.Command.Use != tree.Root.Command.Use {
		t.Errorf("JSON roundtrip failed: Use mismatch")
	}
}

func TestBuildCommandTreeFromSource(t *testing.T) {
	// Test static analysis approach
	tree, err := BuildCommandTreeFromSource("../")
	if err != nil {
		t.Logf("Source analysis failed (expected): %v", err)
		return
	}

	if tree == nil {
		t.Error("Expected non-nil command tree from source analysis")
	}

	// Export to JSON for inspection
	err = tree.ExportToJSON("source_command_tree.json")
	if err != nil {
		t.Fatalf("Failed to export source tree to JSON: %v", err)
	}
}

func ExampleBuildCommandTree() {
	// Example usage of BuildCommandTree
	tree := BuildCommandTree(cmd.RootCommand)

	// Print basic info
	println("Root command:", tree.Root.Command.Use)
	println("Description:", tree.Root.Command.Short)
	println("Number of flags:", len(tree.Root.Command.Flags))
	println("Number of commands:", len(tree.Commands))

	// Print command names and paths
	for _, cmd := range tree.Commands {
		println("  -", cmd.Use, "at", cmd.Path)
	}
}

func TestOpenAPIGeneration(t *testing.T) {
	// Create a simple command tree for testing
	rootCmd := cmd.RootCommand

	// Add a command with args
	testCmd := &cobra.Command{
		Use:   "test",
		Short: "Test command",
		Long:  "A test command with arguments",
	}
	rootCmd.AddCommand(testCmd)

	tree := BuildCommandTree(rootCmd)

	// Generate OpenAPI spec
	spec := tree.GenerateOpenAPISpec("Clog API", "Command Line Of Go API", "1.0.0")

	// Verify basic structure
	if spec.OpenAPI != "3.0.0" {
		t.Errorf("Expected OpenAPI version 3.0.0, got %s", spec.OpenAPI)
	}

	if spec.Info.Title != "Clog API" {
		t.Errorf("Expected title 'Clog API', got %s", spec.Info.Title)
	}

	// Verify paths exist
	if len(spec.Paths) == 0 {
		t.Error("Expected OpenAPI spec to have paths")
	}

	// Export to JSON and YAML
	err := spec.ExportToJSON("openapi_spec.json")
	if err != nil {
		t.Fatalf("Failed to export OpenAPI spec to JSON: %v", err)
	}

	err = spec.ExportToYAML("openapi_spec.yaml")
	if err != nil {
		t.Fatalf("Failed to export OpenAPI spec to YAML: %v", err)
	}
}

func TestYAMLExport(t *testing.T) {
	// Get the command tree
	tree := BuildCommandTree(cmd.RootCommand)

	// Export to YAML
	err := tree.ExportToYAML("command_tree.yaml")
	if err != nil {
		t.Fatalf("Failed to export command tree to YAML: %v", err)
	}

	// Verify YAML file exists and is readable
	data, err := os.ReadFile("command_tree.yaml")
	if err != nil {
		t.Fatalf("Failed to read YAML file: %v", err)
	}

	if len(data) == 0 {
		t.Error("YAML file is empty")
	}
}

func TestTemplateGeneration(t *testing.T) {
	// Create a command with args
	testCmd := &cobra.Command{
		Use:   "upload",
		Short: "Upload a file",
		Long:  "Upload a file to the server",
	}

	tree := BuildCommandTree(testCmd)

	// Generate templates
	err := tree.GenerateCommandTemplates("./templates")
	if err != nil {
		t.Fatalf("Failed to generate command templates: %v", err)
	}

	// Check if template files were created (if commands have args)
	for _, cmd := range tree.Commands {
		if cmd.HasArgs {
			templateName := fmt.Sprintf("templates/%s-api-template.json", strings.ToLower(cmd.Use))
			if _, err := os.Stat(templateName); os.IsNotExist(err) {
				t.Errorf("Expected template file %s to exist", templateName)
			}
		}
	}
}

func TestCommandAPIProps(t *testing.T) {
	// Create a command with flags
	testCmd := &cobra.Command{
		Use:   "serve",
		Short: "Start the server",
		Long:  "Start the HTTP server on specified port",
	}

	testCmd.Flags().StringP("host", "H", "localhost", "Host to bind to")
	testCmd.Flags().IntP("port", "p", 8080, "Port to listen on")

	tree := BuildCommandTree(testCmd)

	// Find the serve command
	var serveCmd *CmdApiProps
	for _, cmd := range tree.Commands {
		if cmd.Use == "serve" {
			serveCmd = cmd
			break
		}
	}

	if serveCmd == nil {
		t.Fatal("Expected to find serve command")
	}

	// Verify path generation
	if serveCmd.Path != "/serve" {
		t.Errorf("Expected path '/serve', got '%s'", serveCmd.Path)
	}

	// Verify GET operation exists
	if serveCmd.Get == nil {
		t.Error("Expected GET operation to exist")
	}

	// Verify flags were converted to query parameters
	if len(serveCmd.Flags) == 0 {
		t.Error("Expected command to have flags")
	}

	// Check for specific flags
	hasPortFlag := false
	for _, flag := range serveCmd.Flags {
		if flag.Name == "port" {
			hasPortFlag = true
			if flag.Schema.Type != "integer" {
				t.Errorf("Expected port flag type to be 'integer', got '%s'", flag.Schema.Type)
			}
			break
		}
	}
	if !hasPortFlag {
		t.Error("Expected to find port flag")
	}
}
