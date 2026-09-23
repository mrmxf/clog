//  Copyright ©2017-2025  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License           https://opensource.org/license/bsd-3-clause/

// Package main is clog: one vocabulary for building and deploying, on whatever
// runs your CI.
//
// This repo carries the CLI **and** the workflows and actions that build it, so
// forking one repository gives you your own clog and your own CI in the same
// tree. That is the whole point — lazy, happy users. It wires together
// github.com/mrmxf/util/* modules only, and imports nothing private.
package main

import (
	"embed"
	"log/slog"
	"os"
	"runtime"

	"github.com/mrmxf/util/bc"
	semver "github.com/mrmxf/util/buildinfo"
	"github.com/mrmxf/util/check"
	"github.com/mrmxf/util/ci"
	"github.com/mrmxf/util/cmdlog"
	"github.com/mrmxf/util/crayon"
	"github.com/mrmxf/util/embedfs"
	"github.com/mrmxf/util/install"
	"github.com/mrmxf/util/kfg"
	"github.com/mrmxf/util/slogger"
	"github.com/mrmxf/util/snippets"
	"github.com/mrmxf/util/snips"
	"github.com/mrmxf/util/source"
	"github.com/spf13/cobra"
)

//go:embed releases.yaml
var ReleasesFs embed.FS

var debugFlag bool

var rootCmd = &cobra.Command{
	Use:   "clog",
	Short: "Build and deploy with one vocabulary, on any CI",
	Long: `clog gives a repo one set of verbs — build, deploy, check, CI — that run
identically on a laptop, on GitHub Actions and on GitLab CI. What a repo is, and
what its CI may do, are declared once in .clog.yaml; the workflows in this
repository read that rather than repeating it.

Case says who owns a verb: ` + "`clog Build`" + ` is always the shipped implementation,
` + "`clog build`" + ` is your override if you wrote one. Where only one exists, either
spelling works.

Run without arguments to see this help. Use --debug for verbose logging.`,
	RunE: func(cmd *cobra.Command, args []string) error {
		return cmd.Help()
	},
}

func main() {
	defer func() { _ = slogger.CloseLogger() }()
	// Nothing useful to do with a logger-setup error: there is no logger yet
	// to report it with, and the slog default still works.
	_, _ = slogger.UsePrettyWithDbgTmpLogger(slog.LevelInfo)

	// Config is layered, and the order matters.
	//
	//   1. util's embedded konfig — the base every clog app shares. It carries
	//      the bc-* build workers (bc-hugo, bc-golang, bc-image), the
	//      git:/project:/install: snippet trees and the check: blocks. Without
	//      it this binary has no build tasks to run.
	//   2. this repo's own .clog.yaml.
	//   3. the working directory's .clog.yaml, via a deferred AutoMerge.
	//
	// Step 3 is the whole point when another repo runs this binary in CI: that
	// repo's .clog.yaml is where its ci.policy/modes/targets live. AutoMerge is
	// deferred rather than left to Konfigure so it lands AFTER step 2 - run by
	// Konfigure it would fire first, and this repo's own config would then
	// override the config of the repo being built.
	bootOpts := kfg.KonfigureOpt{
		AppFs:            embedfs.CoreFs,
		FilePath:         "konfig.yaml",
		PreventAutoMerge: true,
		PreventAutoApp:   true,
	}
	if err := kfg.Konfigure(&bootOpts); err != nil {
		slog.Warn("base config load failed", "err", err)
	}
	if err := kfg.AutoMerge(); err != nil {
		slog.Debug("working-directory config merge failed", "err", err)
	}

	// set SemVer so commands that format version strings have data
	info := semver.Info()

	// `clog --version` is part of the contract, not a convenience: the
	// clog-prepare action runs it to prove the binary it just built works, and
	// setup-clog reports it as the action's `version` output. Cobra supplies
	// the flag as soon as Version is non-empty. `clog buildinfo` stays the
	// long form.
	rootCmd.Version = info.Long
	if rootCmd.Version == "" {
		rootCmd.Version = info.Short
	}

	if err := bootStrap(rootCmd); err != nil {
		slog.Error("bootstrap failed", "err", err)
		os.Exit(1)
	}

	// Case is meaningful in clog: a Capitalised command is shipped, a
	// lowercase one is yours, and where both exist they are two commands on
	// purpose - `clog build` runs your override, `clog Build` runs the default
	// it overrides. ResolveCase supplies the other half of that bargain: when
	// only one case exists there is no ambiguity to preserve, so either
	// spelling should just work. Exact matches are left alone, which is what
	// keeps the pair reachable.
	rootCmd.SetArgs(snips.ResolveCase(rootCmd, os.Args[1:]))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

// bootStrap wires all util commands into the cobra tree.
func bootStrap(root *cobra.Command) error {
	// version / buildinfo
	root.AddCommand(semver.Command)

	// check: try/ok/catch/finally blocks from config
	root.AddCommand(check.Command)

	// declarative tool installer
	root.AddCommand(install.Command)

	// ── build tasks ──────────────────────────────────────────────────────────
	// The embedded bc-* workers call all of these, so a binary without them can
	// load the build config and then fail part-way through running it:
	//   clog BC      flow, gen buildinfo, git, stash log, releases, semver
	//   clog Log     worker logging
	//   clog Crayon  terminal colours the workers eval
	//   clog CI      show, mode, target, stack, run, require, deploy, probe
	//   clog Source  `clog Source project config`
	root.AddCommand(bc.Command)
	root.AddCommand(ci.Command)
	root.AddCommand(cmdlog.Command)
	root.AddCommand(crayon.Command)
	root.AddCommand(source.Command)

	// snippet tree from .clog.yaml
	snippetsKey := "snippets"
	rawSnippetsMap, _ := kfg.Unmarshal[snips.RawSnippets](snippetsKey, "")
	if rawSnippetsMap != nil {
		opts := snippets.Command{
			Use:     "Snippets",
			Key:     snippetsKey,
			Verbose: false,
			Plain:   false,
			Raw:     *rawSnippetsMap,
		}
		root.AddCommand(snippets.Bootstrap(root, opts))
	}

	return nil
}

func init() {
	_, file, _, _ := runtime.Caller(0)
	slog.Debug("init " + file)

	root := rootCmd
	root.PersistentFlags().BoolVar(&debugFlag, "debug", false, "enable debug logging")
	root.PersistentPreRun = func(cmd *cobra.Command, args []string) {
		if debugFlag {
			_, _ = slogger.UsePrettyWithDbgTmpLogger(slog.LevelDebug)
			slog.Debug("debug logging enabled")
		}
	}
}
