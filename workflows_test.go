//  Copyright ©2017-2025  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License           https://opensource.org/license/bsd-3-clause/

package main

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"
)

// The workflows in this repo have two failure modes that cost a consumer's
// pipeline rather than ours, and neither shows up in a local build:
//
//  1. a SARIF category collision. Since 2025-07 GitHub rejects several SARIF
//     runs sharing one category, and upload-sarif takes ONE category per call.
//     The fix was a matrix with `category: ${{ matrix.file }}`; collapsing that
//     back to a single upload is a one-line edit that looks tidier and breaks
//     every scan upload.
//
//  2. a relative `uses:` inside a REUSABLE workflow. `uses: ./…` there resolves
//     against the CALLER's checkout, not this repo, so it silently looks for
//     our action inside the consumer's tree. It works perfectly in this repo's
//     own CI - which is exactly why nothing here would catch it.
//
// Both are cheap to assert and expensive to discover.

func workflow(t *testing.T, name string) string {
	t.Helper()
	b, err := os.ReadFile(filepath.Join(".github", "workflows", name))
	if err != nil {
		t.Fatalf("read %s: %v", name, err)
	}
	return string(b)
}

func TestScanUploadIsOneCategoryPerFile(t *testing.T) {
	src := workflow(t, "build-check.yaml")

	if !strings.Contains(src, "matrix:") || !strings.Contains(src, "file: ${{ fromJSON(") {
		t.Error("build-check.yaml: the scan upload is no longer a per-file matrix; " +
			"GitHub rejects multiple SARIF runs sharing one category")
	}
	if !strings.Contains(src, "category: ${{ matrix.file }}") {
		t.Error("build-check.yaml: upload-sarif must pass `category: ${{ matrix.file }}` - " +
			"one category per file is the only shape it accepts")
	}
	// fail-fast off, so one malformed sweep cannot hide the others.
	if !strings.Contains(src, "fail-fast: false") {
		t.Error("build-check.yaml: the scan matrix needs `fail-fast: false`, " +
			"or the first bad sweep hides every other finding")
	}
}

// reusableUses matches a `uses:` line pointing at a path inside this repo.
var reusableUses = regexp.MustCompile(`(?m)^\s*uses:\s*\./`)

// usesSetupClog matches an actual invocation of the download action.
var usesSetupClog = regexp.MustCompile(`(?m)^\s*uses:.*setup-clog.*$`)

func TestReusableWorkflowsUseAbsoluteActionRefs(t *testing.T) {
	// These are called BY OTHER REPOS, so `./…` would resolve in their tree.
	for _, name := range []string{"build-check.yaml", "deploy-probe.yaml"} {
		src := workflow(t, name)
		if !strings.Contains(src, "workflow_call:") {
			t.Fatalf("%s is no longer a reusable workflow - this test needs rethinking", name)
		}
		if loc := reusableUses.FindString(src); loc != "" {
			t.Errorf("%s: relative `uses:` (%q) inside a reusable workflow resolves against "+
				"the CALLER's checkout, not this repo - it must be mrmxf/clog/...@<ref>",
				name, strings.TrimSpace(loc))
		}
	}
}

func TestSelfWorkflowsBuildFromThisCheckout(t *testing.T) {
	// The chicken-and-egg guard: if this repo ever downloads a clog release to
	// build itself, it can no longer go green before it has published one.
	for _, name := range []string{"self-build.yaml", "self-release.yaml"} {
		src := workflow(t, name)
		if !strings.Contains(src, "self-build: true") {
			t.Errorf("%s: must set `self-build: true` - this repo IS clog and must "+
				"never bootstrap from setup-clog", name)
		}
		// A `uses:` of setup-clog, not the word: both files warn against it in
		// prose, and a comment saying "never do this" is not doing it.
		if loc := usesSetupClog.FindString(src); loc != "" {
			t.Errorf("%s: %q - that would make the repo that produces every "+
				"release depend on a release existing", name, strings.TrimSpace(loc))
		}
	}
}

func TestArtifactPathsAreCanonical(t *testing.T) {
	// `tmp/` is what these dirs replaced, and it rots back easily because every
	// old snippet and every old habit still says tmp.
	stale := regexp.MustCompile(`(?m)(^|\s)tmp/`)
	dir := filepath.Join(".github", "workflows")
	entries, err := os.ReadDir(dir)
	if err != nil {
		t.Fatalf("read %s: %v", dir, err)
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".yaml") {
			continue
		}
		src := workflow(t, e.Name())
		if loc := stale.FindString(src); loc != "" {
			t.Errorf("%s: writes to tmp/ - artifacts belong in _clog_build/ or _clog_deploy/", e.Name())
		}
	}
}

// topLevelPermissions matches a workflow-level `permissions:` key.
var topLevelPermissions = regexp.MustCompile(`(?m)^permissions:`)

// Both of these shipped, and together they meant no GitHub target could ever
// publish through deploy-probe - found only when pihuw's v0.4.14 tag tried:
//
//  1. A reusable workflow's own permissions block can only NARROW what the
//     caller grants. `contents: read` there cut pihuw's write back to read:
//     the gh-pages push got 403.
//  2. Actions does not put its token in the environment, and the GitHub
//     deployers read GH_TOKEN from there: "could not read Username".
func TestDeployProbeCanPublishToGitHub(t *testing.T) {
	src := workflow(t, "deploy-probe.yaml")
	if topLevelPermissions.MatchString(src) {
		t.Error("deploy-probe.yaml: a top-level permissions block narrows the caller's grant; " +
			"remove it and let the caller grant contents: write for GitHub targets")
	}
	// the run lines themselves, not the header comment that mentions them
	for _, step := range []string{`clog CI deploy --target "$CLOG_TARGET"`, "clog CI probe || true"} {
		i := strings.Index(src, step)
		if i < 0 {
			t.Fatalf("deploy-probe.yaml no longer runs %q - this test needs rethinking", step)
		}
		start := strings.LastIndex(src[:i], "- name:")
		if !strings.Contains(src[start:i], "GH_TOKEN: ${{ github.token }}") {
			t.Errorf("deploy-probe.yaml: the step running %q does not set GH_TOKEN; "+
				"the GitHub deployers cannot authenticate without it", step)
		}
	}
}
