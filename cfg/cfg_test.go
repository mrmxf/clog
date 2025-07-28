package cfg_test

import (
	"embed"
	"testing"

	"github.com/mrmxf/clog/cfg"
	"github.com/mrmxf/clog/core"
)

//go:embed test_releases.yaml
var TestFs embed.FS

func TestNew(t *testing.T) {
	// Test creating a new cfg with embedded file systems
	CoreFs := core.CoreFs
	configEmbedFsList := &[]embed.FS{CoreFs, TestFs}

	kfg := cfg.New(configEmbedFsList, nil)
	if kfg == nil {
		t.Fatal("Expected cfg.New() to return a valid Config, got nil")
	}
}

func TestKfg(t *testing.T) {
	// Initialize first
	CoreFs := core.CoreFs
	configEmbedFsList := &[]embed.FS{CoreFs, TestFs}
	cfg.New(configEmbedFsList, nil)

	// Test getting the global config
	kfg := cfg.Kfg()
	if kfg == nil {
		t.Fatal("Expected cfg.Kfg() to return a valid Config, got nil")
	}
}

func TestCoreFs(t *testing.T) {
	// Initialize first
	CoreFs := core.CoreFs
	configEmbedFsList := &[]embed.FS{CoreFs, TestFs}
	cfg.New(configEmbedFsList, nil)

	// Test getting the core filesystem
	coreFs := cfg.CoreFs()

	// Try to read a known file from core
	_, err := coreFs.ReadFile("core.clog.yaml")
	if err != nil {
		t.Fatalf("Expected to read core.clog.yaml from CoreFs, got error: %v", err)
	}
}

func TestFsCache(t *testing.T) {
	// Initialize first
	CoreFs := core.CoreFs
	configEmbedFsList := &[]embed.FS{CoreFs, TestFs}
	cfg.New(configEmbedFsList, nil)

	// Test getting the filesystem cache
	fsCache := cfg.FsCache()
	if len(fsCache) < 2 {
		t.Fatalf("Expected fsCache to have at least 2 filesystems, got %d", len(fsCache))
	}
}

func TestSearchPaths(t *testing.T) {
	// Initialize first
	CoreFs := core.CoreFs
	configEmbedFsList := &[]embed.FS{CoreFs, TestFs}
	cfg.New(configEmbedFsList, nil)

	// Test getting search paths
	paths := cfg.SearchPaths()
	if paths == nil {
		t.Fatal("Expected cfg.SearchPaths() to return a valid slice pointer, got nil")
	}
}

func TestConfigValues(t *testing.T) {
	// Initialize first
	CoreFs := core.CoreFs
	configEmbedFsList := &[]embed.FS{CoreFs, TestFs}
	kfg := cfg.New(configEmbedFsList, nil)

	// Test reading some known config values
	logLevel := kfg.String("clog.log.level")
	if logLevel == "" {
		t.Error("Expected to read clog.log.level from config, got empty string")
	}

	logStyle := kfg.String("clog.log.style")
	if logStyle == "" {
		t.Error("Expected to read clog.log.style from config, got empty string")
	}
}
