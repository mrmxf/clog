# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Commands

### Build
- **Build**: `go build -o clog clog-sample.go` - Builds the main clog executable from the sample
- **Build all targets**: Use the built-in workflow `clog bc-golang` which builds for multiple platforms (linux/amd64, darwin/amd64, linux/arm64, darwin/arm64) and outputs to `tmp/` directory

### Test
- **Run all tests**: `go test ./...` - Runs all tests across all packages
- **Run specific test**: `go test ./package-name` - Runs tests for a specific package
- **Test with verbose output**: `go test -v ./...`

### Development
- **Run clog locally**: `./clog` (after building) or `go run clog-sample.go`
- **Check configuration**: `clog check pre-build` - Runs pre-build validation checks
- **List available commands**: `clog --help`

## Architecture

**Clog** is a command-line tool and configuration management system built in Go. It follows a modular architecture:

### Core Components
- **`cmd/`**: Command implementations using Cobra framework
  - `cmd/root.go`: Root command and bootstrap logic
  - Individual command packages: `aws/`, `check/`, `init/`, `list/`, etc.
- **`config/`**: Configuration management using Viper, supports YAML configs with search paths and embedded filesystems
- **`core/`**: Core embedded files and templates, contains `core.clog.yaml` with default configuration
- **`scripts/`** & **`shell/`**: Script execution and shell integration
- **`slogger/`**: Structured logging with pretty-printing support
- **`semver/`**: Semantic versioning support with release tracking

### Configuration System
- Uses embedded filesystems (`embed.FS`) to bundle configuration files
- Searches multiple paths for `clog.yaml` files: `/var/clogrc/`, `$HOME/.config/clogrc/`, `./clogrc/`, etc.
- Supports configuration overlaying and environment variable binding
- Main config file: `core/core.clog.yaml` with extensive snippet definitions

### Command Architecture
- Built on Cobra CLI framework
- Commands can execute shell snippets defined in configuration
- Supports complex workflow commands like `bc-golang`, `bc-hugo`, `bc-ko` for building different project types
- Check commands validate environment and prerequisites

### Key Features
- **Snippets**: Reusable shell commands defined in YAML configuration
- **Multi-target builds**: Cross-platform Go binary compilation
- **Environment checks**: Validates tools and dependencies before builds
- **Logging**: Structured logging with colored output and multiple formats
- **Web server**: Optional `gommi/` package provides HTTP server capabilities
- **AWS integration**: DNS and cloud resource management commands

### Testing Strategy
- Uses Go's standard testing framework
- Test files follow `*_test.go` convention
- System tests in `testclog/` package
- Individual unit tests across packages: `slogger/`, `semver/`, `cmd/`, etc.

### Build Process
The project uses a sophisticated build system defined in YAML snippets that supports:
- Multi-architecture compilation
- Version injection via linker flags
- Release management through `releases.yaml`
- Docker/container builds with `ko`
- Hugo static site generation