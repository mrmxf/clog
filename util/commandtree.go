//  Copyright ©2017-2025  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License           https://opensource.org/license/bsd-3-clause/

package util

import (
	"go/ast"
	"go/parser"
	"go/token"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

// QueryParam represents a query parameter for OpenAPI
type QueryParam struct {
	Name        string `json:"name" yaml:"name"`
	In          string `json:"in" yaml:"in"`
	Description string `json:"description" yaml:"description"`
	Required    bool   `json:"required" yaml:"required"`
	Schema      Schema `json:"schema" yaml:"schema"`
}

// Schema represents an OpenAPI schema
type Schema struct {
	Type        string                 `json:"type" yaml:"type"`
	Format      string                 `json:"format,omitempty" yaml:"format,omitempty"`
	Default     interface{}            `json:"default,omitempty" yaml:"default,omitempty"`
	Enum        []string               `json:"enum,omitempty" yaml:"enum,omitempty"`
	Items       *Schema                `json:"items,omitempty" yaml:"items,omitempty"`
	Properties  map[string]Schema      `json:"properties,omitempty" yaml:"properties,omitempty"`
	Required    []string               `json:"required,omitempty" yaml:"required,omitempty"`
	Example     interface{}            `json:"example,omitempty" yaml:"example,omitempty"`
	Description string                 `json:"description,omitempty" yaml:"description,omitempty"`
}

// RequestBody represents OpenAPI request body
type RequestBody struct {
	Description string                 `json:"description" yaml:"description"`
	Required    bool                   `json:"required" yaml:"required"`
	Content     map[string]ContentType `json:"content" yaml:"content"`
}

// ContentType represents OpenAPI content type
type ContentType struct {
	Schema   Schema                 `json:"schema,omitempty" yaml:"schema,omitempty"`
	Example  interface{}            `json:"example,omitempty" yaml:"example,omitempty"`
	Examples map[string]ExampleObj  `json:"examples,omitempty" yaml:"examples,omitempty"`
	Encoding map[string]Encoding    `json:"encoding,omitempty" yaml:"encoding,omitempty"`
}

// ExampleObj represents OpenAPI example object
type ExampleObj struct {
	Summary       string      `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description   string      `json:"description,omitempty" yaml:"description,omitempty"`
	Value         interface{} `json:"value,omitempty" yaml:"value,omitempty"`
	ExternalValue string      `json:"externalValue,omitempty" yaml:"externalValue,omitempty"`
}

// Encoding represents OpenAPI encoding object
type Encoding struct {
	ContentType   string            `json:"contentType,omitempty" yaml:"contentType,omitempty"`
	Headers       map[string]Header `json:"headers,omitempty" yaml:"headers,omitempty"`
	Style         string            `json:"style,omitempty" yaml:"style,omitempty"`
	Explode       bool              `json:"explode,omitempty" yaml:"explode,omitempty"`
	AllowReserved bool              `json:"allowReserved,omitempty" yaml:"allowReserved,omitempty"`
}

// Response represents OpenAPI response
type Response struct {
	Description string                 `json:"description" yaml:"description"`
	Headers     map[string]Header      `json:"headers,omitempty" yaml:"headers,omitempty"`
	Content     map[string]ContentType `json:"content,omitempty" yaml:"content,omitempty"`
	Links       map[string]Link        `json:"links,omitempty" yaml:"links,omitempty"`
}

// Header represents OpenAPI header
type Header struct {
	Description string  `json:"description,omitempty" yaml:"description,omitempty"`
	Required    bool    `json:"required,omitempty" yaml:"required,omitempty"`
	Deprecated  bool    `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	Schema      Schema  `json:"schema,omitempty" yaml:"schema,omitempty"`
	Example     interface{} `json:"example,omitempty" yaml:"example,omitempty"`
}

// Link represents OpenAPI link
type Link struct {
	OperationRef string            `json:"operationRef,omitempty" yaml:"operationRef,omitempty"`
	OperationId  string            `json:"operationId,omitempty" yaml:"operationId,omitempty"`
	Parameters   map[string]string `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	RequestBody  interface{}       `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
	Description  string            `json:"description,omitempty" yaml:"description,omitempty"`
	Server       *OpenAPIServer    `json:"server,omitempty" yaml:"server,omitempty"`
}

// Operation represents OpenAPI operation
type Operation struct {
	Summary     string                 `json:"summary,omitempty" yaml:"summary,omitempty"`
	Description string                 `json:"description,omitempty" yaml:"description,omitempty"`
	OperationId string                 `json:"operationId,omitempty" yaml:"operationId,omitempty"`
	Tags        []string               `json:"tags,omitempty" yaml:"tags,omitempty"`
	Parameters  []QueryParam           `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	RequestBody *RequestBody           `json:"requestBody,omitempty" yaml:"requestBody,omitempty"`
	Responses   map[string]Response    `json:"responses" yaml:"responses"`
	Deprecated  bool                   `json:"deprecated,omitempty" yaml:"deprecated,omitempty"`
	Security    []map[string][]string  `json:"security,omitempty" yaml:"security,omitempty"`
	Servers     []OpenAPIServer        `json:"servers,omitempty" yaml:"servers,omitempty"`
}

// CmdApiProps represents the API properties of a single command
type CmdApiProps struct {
	Use         string             `json:"use" yaml:"use"`
	Short       string             `json:"short" yaml:"short"`
	Long        string             `json:"long" yaml:"long"`
	Aliases     []string           `json:"aliases,omitempty" yaml:"aliases,omitempty"`
	Path        string             `json:"path" yaml:"path"`
	Get         *Operation         `json:"get,omitempty" yaml:"get,omitempty"`
	Post        *Operation         `json:"post,omitempty" yaml:"post,omitempty"`
	HasArgs     bool               `json:"hasArgs" yaml:"hasArgs"`
	Flags       []QueryParam       `json:"flags,omitempty" yaml:"flags,omitempty"`
	Args        []ArgProperty      `json:"args,omitempty" yaml:"args,omitempty"`
}

// ArgProperty represents command argument properties
type ArgProperty struct {
	Name        string `json:"name" yaml:"name"`
	Description string `json:"description" yaml:"description"`
	Required    bool   `json:"required" yaml:"required"`
	Type        string `json:"type" yaml:"type"`
}

// CommandNode represents a node in the command tree
type CommandNode struct {
	Command  *CmdApiProps   `json:"command,omitempty" yaml:"command,omitempty"`
	Children []*CommandNode `json:"children,omitempty" yaml:"children,omitempty"`
}

// CommandTree represents the complete command tree
type CommandTree struct {
	Root     *CommandNode `json:"root" yaml:"root"`
	Commands []*CmdApiProps `json:"commands" yaml:"commands"`
}

// OpenAPISpec represents a complete OpenAPI 3.0 specification
type OpenAPISpec struct {
	OpenAPI    string                           `json:"openapi" yaml:"openapi"`
	Info       OpenAPIInfo                      `json:"info" yaml:"info"`
	Servers    []OpenAPIServer                  `json:"servers,omitempty" yaml:"servers,omitempty"`
	Paths      map[string]map[string]Operation  `json:"paths" yaml:"paths"`
	Components *OpenAPIComponents               `json:"components,omitempty" yaml:"components,omitempty"`
}

// OpenAPIInfo represents OpenAPI info section
type OpenAPIInfo struct {
	Title          string         `json:"title" yaml:"title"`
	Description    string         `json:"description" yaml:"description"`
	Version        string         `json:"version" yaml:"version"`
	Contact        *OpenAPIContact `json:"contact,omitempty" yaml:"contact,omitempty"`
	License        *OpenAPILicense `json:"license,omitempty" yaml:"license,omitempty"`
	TermsOfService string         `json:"termsOfService,omitempty" yaml:"termsOfService,omitempty"`
}

// OpenAPIServer represents OpenAPI server object
type OpenAPIServer struct {
	URL         string                    `json:"url" yaml:"url"`
	Description string                    `json:"description,omitempty" yaml:"description,omitempty"`
	Variables   map[string]ServerVariable `json:"variables,omitempty" yaml:"variables,omitempty"`
}

// ServerVariable represents OpenAPI server variable
type ServerVariable struct {
	Enum        []string `json:"enum,omitempty" yaml:"enum,omitempty"`
	Default     string   `json:"default" yaml:"default"`
	Description string   `json:"description,omitempty" yaml:"description,omitempty"`
}

// OpenAPIContact represents OpenAPI contact information
type OpenAPIContact struct {
	Name  string `json:"name,omitempty" yaml:"name,omitempty"`
	URL   string `json:"url,omitempty" yaml:"url,omitempty"`
	Email string `json:"email,omitempty" yaml:"email,omitempty"`
}

// OpenAPILicense represents OpenAPI license information
type OpenAPILicense struct {
	Name string `json:"name" yaml:"name"`
	URL  string `json:"url,omitempty" yaml:"url,omitempty"`
}

// OpenAPIComponents represents OpenAPI components section
type OpenAPIComponents struct {
	Schemas         map[string]Schema         `json:"schemas,omitempty" yaml:"schemas,omitempty"`
	Responses       map[string]Response       `json:"responses,omitempty" yaml:"responses,omitempty"`
	Parameters      map[string]QueryParam     `json:"parameters,omitempty" yaml:"parameters,omitempty"`
	RequestBodies   map[string]RequestBody    `json:"requestBodies,omitempty" yaml:"requestBodies,omitempty"`
	SecuritySchemes map[string]SecurityScheme `json:"securitySchemes,omitempty" yaml:"securitySchemes,omitempty"`
}

// SecurityScheme represents OpenAPI security scheme
type SecurityScheme struct {
	Type             string            `json:"type" yaml:"type"`
	Description      string            `json:"description,omitempty" yaml:"description,omitempty"`
	Name             string            `json:"name,omitempty" yaml:"name,omitempty"`
	In               string            `json:"in,omitempty" yaml:"in,omitempty"`
	Scheme           string            `json:"scheme,omitempty" yaml:"scheme,omitempty"`
	BearerFormat     string            `json:"bearerFormat,omitempty" yaml:"bearerFormat,omitempty"`
	Flows            *OAuthFlows       `json:"flows,omitempty" yaml:"flows,omitempty"`
	OpenIdConnectUrl string            `json:"openIdConnectUrl,omitempty" yaml:"openIdConnectUrl,omitempty"`
}

// OAuthFlows represents OpenAPI OAuth flows
type OAuthFlows struct {
	Implicit          *OAuthFlow `json:"implicit,omitempty" yaml:"implicit,omitempty"`
	Password          *OAuthFlow `json:"password,omitempty" yaml:"password,omitempty"`
	ClientCredentials *OAuthFlow `json:"clientCredentials,omitempty" yaml:"clientCredentials,omitempty"`
	AuthorizationCode *OAuthFlow `json:"authorizationCode,omitempty" yaml:"authorizationCode,omitempty"`
}

// OAuthFlow represents OpenAPI OAuth flow
type OAuthFlow struct {
	AuthorizationUrl string            `json:"authorizationUrl,omitempty" yaml:"authorizationUrl,omitempty"`
	TokenUrl         string            `json:"tokenUrl,omitempty" yaml:"tokenUrl,omitempty"`
	RefreshUrl       string            `json:"refreshUrl,omitempty" yaml:"refreshUrl,omitempty"`
	Scopes           map[string]string `json:"scopes" yaml:"scopes"`
}

// BuildCommandTree generates a hierarchical structure of all cobra commands
func BuildCommandTree(rootCmd *cobra.Command) *CommandTree {
	root := buildCommandNode(rootCmd, "")
	commands := extractAllCommands(rootCmd, "")
	
	return &CommandTree{
		Root:     root,
		Commands: commands,
	}
}

// buildCommandNode recursively builds CommandNode from a cobra.Command
func buildCommandNode(cmd *cobra.Command, parentPath string) *CommandNode {
	cmdApi := buildCmdApiProps(cmd, parentPath)
	node := &CommandNode{
		Command: cmdApi,
	}

	// Recursively process subcommands
	for _, subCmd := range cmd.Commands() {
		if !subCmd.Hidden {
			childNode := buildCommandNode(subCmd, cmdApi.Path)
			node.Children = append(node.Children, childNode)
		}
	}

	return node
}

// buildCmdApiProps builds CmdApiProps from a cobra.Command
func buildCmdApiProps(cmd *cobra.Command, parentPath string) *CmdApiProps {
	path := buildCommandPath(cmd.Use, parentPath)
	flags := extractQueryParams(cmd)
	args := extractArgProperties(cmd)
	hasArgs := len(args) > 0

	cmdApi := &CmdApiProps{
		Use:     cmd.Use,
		Short:   cmd.Short,
		Long:    cmd.Long,
		Aliases: cmd.Aliases,
		Path:    path,
		HasArgs: hasArgs,
		Flags:   flags,
		Args:    args,
	}

	// Create GET operation
	cmdApi.Get = &Operation{
		Summary:     cmd.Short,
		Description: cmd.Long,
		OperationId: generateOperationId("get", cmd.Use),
		Tags:        []string{"commands"},
		Parameters:  flags,
		Responses: map[string]Response{
			"200": {
				Description: "Command executed successfully",
				Content: map[string]ContentType{
					"application/json": {
						Schema: Schema{
							Type: "object",
							Properties: map[string]Schema{
								"output": {
									Type:        "string",
									Description: "Command output",
								},
								"exitCode": {
									Type:        "integer",
									Description: "Exit code",
								},
								"duration": {
									Type:        "number",
									Description: "Execution duration in milliseconds",
								},
							},
							Required: []string{"output", "exitCode"},
						},
					},
				},
			},
			"400": {
				Description: "Bad request",
				Content: map[string]ContentType{
					"application/json": {
						Schema: Schema{
							Type: "object",
							Properties: map[string]Schema{
								"error": {
									Type:        "string",
									Description: "Error message",
								},
								"code": {
									Type:        "integer",
									Description: "Error code",
								},
							},
							Required: []string{"error"},
						},
					},
				},
			},
			"500": {
				Description: "Internal server error",
				Content: map[string]ContentType{
					"application/json": {
						Schema: Schema{
							Type: "object",
							Properties: map[string]Schema{
								"error": {
									Type:        "string",
									Description: "Error message",
								},
								"code": {
									Type:        "integer",
									Description: "Error code",
								},
							},
							Required: []string{"error"},
						},
					},
				},
			},
		},
	}

	// Create POST operation if command has args
	if hasArgs {
		cmdApi.Post = &Operation{
			Summary:     cmd.Short + " (with arguments)",
			Description: cmd.Long,
			OperationId: generateOperationId("post", cmd.Use),
			Tags:        []string{"commands"},
			Parameters:  flags,
			RequestBody: &RequestBody{
				Description: "Command arguments",
				Required:    true,
				Content: map[string]ContentType{
					"application/json": {
						Schema: buildArgsSchema(args),
					},
				},
			},
			Responses: map[string]Response{
				"200": {
					Description: "Command executed successfully",
					Content: map[string]ContentType{
						"application/json": {
							Schema: Schema{
								Type: "object",
								Properties: map[string]Schema{
									"output": {
										Type:        "string",
										Description: "Command output",
									},
									"exitCode": {
										Type:        "integer",
										Description: "Exit code",
									},
									"duration": {
										Type:        "number",
										Description: "Execution duration in milliseconds",
									},
								},
								Required: []string{"output", "exitCode"},
							},
						},
					},
				},
				"400": {
					Description: "Bad request",
					Content: map[string]ContentType{
						"application/json": {
							Schema: Schema{
								Type: "object",
								Properties: map[string]Schema{
									"error": {
										Type:        "string",
										Description: "Error message",
									},
									"code": {
										Type:        "integer",
										Description: "Error code",
									},
								},
								Required: []string{"error"},
							},
						},
					},
				},
				"500": {
					Description: "Internal server error",
					Content: map[string]ContentType{
						"application/json": {
							Schema: Schema{
								Type: "object",
								Properties: map[string]Schema{
									"error": {
										Type:        "string",
										Description: "Error message",
									},
									"code": {
										Type:        "integer",
										Description: "Error code",
									},
								},
								Required: []string{"error"},
							},
						},
					},
				},
			},
		}
	}

	return cmdApi
}

// extractQueryParams extracts query parameters from cobra command flags
func extractQueryParams(cmd *cobra.Command) []QueryParam {
	var params []QueryParam

	// Process local flags
	cmd.LocalFlags().VisitAll(func(flag *pflag.Flag) {
		params = append(params, createQueryParam(flag))
	})

	// Process inherited flags
	cmd.InheritedFlags().VisitAll(func(flag *pflag.Flag) {
		params = append(params, createQueryParam(flag))
	})

	return params
}

// createQueryParam creates a QueryParam from a pflag.Flag
func createQueryParam(flag *pflag.Flag) QueryParam {
	paramType, format := mapFlagTypeToOpenAPI(flag.Value.Type())
	schema := Schema{
		Type:   paramType,
		Format: format,
	}
	
	if flag.DefValue != "" {
		schema.Default = parseDefaultValue(flag.DefValue, paramType)
	}
	
	// Handle array types
	if paramType == "array" {
		schema.Items = &Schema{
			Type: "string",
		}
	}

	return QueryParam{
		Name:        flag.Name,
		In:          "query",
		Description: flag.Usage,
		Required:    false, // pflag doesn't have a direct Required field
		Schema:      schema,
	}
}

// extractArgProperties attempts to extract argument information through static analysis
func extractArgProperties(cmd *cobra.Command) []ArgProperty {
	var args []ArgProperty

	// Since cobra.Args functions can't be compared directly, we use static analysis
	// based on command patterns and documentation
	args = analyzeArgsFromSource(cmd)

	return args
}

// analyzeArgsFromSource attempts to analyze arguments from source code
func analyzeArgsFromSource(cmd *cobra.Command) []ArgProperty {
	var args []ArgProperty

	// This is a simplified approach - in a real implementation,
	// you might want to analyze the Run function's source code
	// to understand how it uses the args parameter

	// For now, we'll make educated guesses based on command patterns
	use := strings.ToLower(cmd.Use)
	
	switch {
	case strings.Contains(use, "copy"):
		args = append(args, ArgProperty{
			Name:        "source",
			Description: "Source file or path",
			Required:    true,
			Type:        "string",
		})
		args = append(args, ArgProperty{
			Name:        "destination",
			Description: "Destination file or path",
			Required:    true,
			Type:        "string",
		})
	case strings.Contains(use, "version"):
		args = append(args, ArgProperty{
			Name:        "format",
			Description: "Version format (short, note, or full)",
			Required:    false,
			Type:        "string",
		})
	case strings.Contains(use, "cat"):
		args = append(args, ArgProperty{
			Name:        "file",
			Description: "File to display",
			Required:    true,
			Type:        "string",
		})
	case strings.Contains(use, "check"):
		args = append(args, ArgProperty{
			Name:        "path",
			Description: "Path to check",
			Required:    false,
			Type:        "string",
		})
	}

	return args
}

// BuildCommandTreeFromSource analyzes cobra commands from source files
func BuildCommandTreeFromSource(rootDir string) (*CommandTree, error) {
	fset := token.NewFileSet()
	
	// Parse all Go files in the cmd directory
	cmdDir := filepath.Join(rootDir, "cmd")
	packages, err := parser.ParseDir(fset, cmdDir, nil, parser.ParseComments)
	if err != nil {
		return nil, err
	}

	// Build command tree from AST
	rootCmd := &CmdApiProps{
		Use:   "root",
		Short: "Root command",
		Long:  "Root command discovered from source analysis",
		Path:  "/",
	}

	var commands []*CmdApiProps
	commands = append(commands, rootCmd)

	for _, pkg := range packages {
		for _, file := range pkg.Files {
			sourceCommands := extractCommandsFromAST(file)
			commands = append(commands, sourceCommands...)
		}
	}

	return &CommandTree{
		Root: &CommandNode{
			Command: rootCmd,
		},
		Commands: commands,
	}, nil
}

// extractCommandsFromAST extracts command information from Go AST
func extractCommandsFromAST(file *ast.File) []*CmdApiProps {
	var commands []*CmdApiProps

	ast.Inspect(file, func(n ast.Node) bool {
		// Look for cobra.Command struct literals
		if composite, ok := n.(*ast.CompositeLit); ok {
			if selector, ok := composite.Type.(*ast.SelectorExpr); ok {
				if ident, ok := selector.X.(*ast.Ident); ok {
					if ident.Name == "cobra" && selector.Sel.Name == "Command" {
						cmd := parseCommandFromComposite(composite)
						if cmd != nil {
							commands = append(commands, cmd)
						}
					}
				}
			}
		}
		return true
	})

	return commands
}

// parseCommandFromComposite parses a cobra.Command from a composite literal
func parseCommandFromComposite(composite *ast.CompositeLit) *CmdApiProps {
	cmd := &CmdApiProps{}

	for _, elt := range composite.Elts {
		if kv, ok := elt.(*ast.KeyValueExpr); ok {
			if ident, ok := kv.Key.(*ast.Ident); ok {
				switch ident.Name {
				case "Use":
					if lit, ok := kv.Value.(*ast.BasicLit); ok {
						cmd.Use = strings.Trim(lit.Value, "\"")
						cmd.Path = "/" + strings.ToLower(cmd.Use)
					}
				case "Short":
					if lit, ok := kv.Value.(*ast.BasicLit); ok {
						cmd.Short = strings.Trim(lit.Value, "\"")
					}
				case "Long":
					if lit, ok := kv.Value.(*ast.BasicLit); ok {
						cmd.Long = strings.Trim(lit.Value, "\"")
					}
				case "Aliases":
					if composite, ok := kv.Value.(*ast.CompositeLit); ok {
						for _, alias := range composite.Elts {
							if lit, ok := alias.(*ast.BasicLit); ok {
								cmd.Aliases = append(cmd.Aliases, strings.Trim(lit.Value, "\""))
							}
						}
					}
				}
			}
		}
	}

	return cmd
}