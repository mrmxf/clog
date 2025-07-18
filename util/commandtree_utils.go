//  Copyright ©2017-2025  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License           https://opensource.org/license/bsd-3-clause/

package util

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"gopkg.in/yaml.v3"
)

// Helper functions

// buildCommandPath builds the API path for a command
func buildCommandPath(use, parentPath string) string {
	commandName := strings.Split(use, " ")[0] // Get just the command name, ignore args
	if parentPath == "" {
		return "/" + commandName
	}
	return parentPath + "/" + commandName
}

// generateOperationId generates an OpenAPI operation ID
func generateOperationId(method, commandUse string) string {
	commandName := strings.Split(commandUse, " ")[0]
	return method + strings.Title(commandName)
}

// mapFlagTypeToOpenAPI maps cobra flag types to OpenAPI types and formats
func mapFlagTypeToOpenAPI(flagType string) (string, string) {
	switch flagType {
	case "bool":
		return "boolean", ""
	case "int", "int8", "int16", "int32":
		return "integer", "int32"
	case "int64":
		return "integer", "int64"
	case "uint", "uint8", "uint16", "uint32":
		return "integer", "int32"
	case "uint64":
		return "integer", "int64"
	case "float32":
		return "number", "float"
	case "float64":
		return "number", "double"
	case "string":
		return "string", ""
	case "stringArray":
		return "array", ""
	default:
		return "string", ""
	}
}

// parseDefaultValue parses default values based on type
func parseDefaultValue(defaultVal, dataType string) interface{} {
	switch dataType {
	case "boolean":
		return defaultVal == "true"
	case "integer":
		if defaultVal == "0" {
			return 0
		}
		return defaultVal
	case "number":
		if defaultVal == "0" {
			return 0.0
		}
		return defaultVal
	default:
		return defaultVal
	}
}

// buildArgsSchema builds OpenAPI schema for command arguments
func buildArgsSchema(args []ArgProperty) Schema {
	properties := make(map[string]Schema)
	requiredFields := []string{}
	
	for _, arg := range args {
		properties[arg.Name] = Schema{
			Type: arg.Type,
		}
		if arg.Required {
			requiredFields = append(requiredFields, arg.Name)
		}
	}
	
	schema := Schema{
		Type: "object",
	}
	
	// Note: OpenAPI schema properties and required fields would need to be added
	// to the Schema struct to fully support this
	return schema
}

// extractAllCommands extracts all commands from a tree into a flat list
func extractAllCommands(cmd *cobra.Command, parentPath string) []*CmdApiProps {
	var commands []*CmdApiProps
	
	// Add current command
	cmdApi := buildCmdApiProps(cmd, parentPath)
	commands = append(commands, cmdApi)
	
	// Add all subcommands recursively
	for _, subCmd := range cmd.Commands() {
		if !subCmd.Hidden {
			subCommands := extractAllCommands(subCmd, cmdApi.Path)
			commands = append(commands, subCommands...)
		}
	}
	
	return commands
}

// Export functions

// ExportToJSON exports a command tree to JSON format
func (ct *CommandTree) ExportToJSON(filename string) error {
	data, err := json.MarshalIndent(ct, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

// ExportToYAML exports a command tree to YAML format
func (ct *CommandTree) ExportToYAML(filename string) error {
	data, err := yaml.Marshal(ct)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

// ExportCommandToJSON exports a single command to JSON format
func (cmd *CmdApiProps) ExportToJSON(filename string) error {
	data, err := json.MarshalIndent(cmd, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

// ExportCommandToYAML exports a single command to YAML format
func (cmd *CmdApiProps) ExportToYAML(filename string) error {
	data, err := yaml.Marshal(cmd)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

// GenerateOpenAPISpec generates a complete OpenAPI 3.0 specification
func (ct *CommandTree) GenerateOpenAPISpec(title, description, version string) *OpenAPISpec {
	paths := make(map[string]map[string]Operation)
	
	for _, cmd := range ct.Commands {
		pathMethods := make(map[string]Operation)
		
		if cmd.Get != nil {
			pathMethods["get"] = *cmd.Get
		}
		
		if cmd.Post != nil {
			pathMethods["post"] = *cmd.Post
		}
		
		if len(pathMethods) > 0 {
			paths[cmd.Path] = pathMethods
		}
	}
	
	return &OpenAPISpec{
		OpenAPI: "3.0.0",
		Info: OpenAPIInfo{
			Title:       title,
			Description: description,
			Version:     version,
		},
		Paths: paths,
	}
}

// ExportOpenAPISpec exports OpenAPI specification to JSON or YAML
func (spec *OpenAPISpec) ExportToJSON(filename string) error {
	data, err := json.MarshalIndent(spec, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

func (spec *OpenAPISpec) ExportToYAML(filename string) error {
	data, err := yaml.Marshal(spec)
	if err != nil {
		return err
	}
	return os.WriteFile(filename, data, 0644)
}

// GenerateCommandTemplates generates API template files for commands with args
func (ct *CommandTree) GenerateCommandTemplates(outputDir string) error {
	for _, cmd := range ct.Commands {
		if cmd.HasArgs {
			templateName := fmt.Sprintf("%s-api-template.json", strings.ToLower(cmd.Use))
			templatePath := filepath.Join(outputDir, templateName)
			
			template := map[string]interface{}{
				"command": cmd.Use,
				"path":    cmd.Path,
				"method":  "POST",
				"args":    cmd.Args,
				"example": generateExamplePayload(cmd.Args),
			}
			
			data, err := json.MarshalIndent(template, "", "  ")
			if err != nil {
				return err
			}
			
			err = os.WriteFile(templatePath, data, 0644)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

// generateExamplePayload generates an example payload for command arguments
func generateExamplePayload(args []ArgProperty) map[string]interface{} {
	example := make(map[string]interface{})
	
	for _, arg := range args {
		switch arg.Type {
		case "string":
			example[arg.Name] = "example_" + arg.Name
		case "integer":
			example[arg.Name] = 1
		case "boolean":
			example[arg.Name] = false
		case "number":
			example[arg.Name] = 1.0
		default:
			example[arg.Name] = "example_value"
		}
	}
	
	return example
}