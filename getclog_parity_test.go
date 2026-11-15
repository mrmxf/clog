//  Copyright ©2017-2025  Mr MXF   info@mrmxf.com
//  BSD-3-Clause License           https://opensource.org/license/bsd-3-clause/

package main

import (
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"runtime"
	"strings"
	"testing"
)

// get-clog.sh bootstraps every CI job; `clog Install clog` (alias `clog
// upgrade`, util/install/self.go) reinstalls clog on a laptop. They are two
// implementations of one algorithm, and if they disagree a laptop quietly runs
// a different clog from CI - toolchain-induced drift that no build would catch.
//
// This runs both against the same inputs, get-clog.sh with CLOG_PLAN=1 and
// clog with --dry-run, and requires the same version, repo, asset, URL and
// install directory. CLOG_INSTALL_DIR is always set: without it the two
// differ on purpose (the script writes ~/.local/bin, clog replaces itself).
// None of the cases touches the network.

var planLine = regexp.MustCompile(`(?m)^(version|repo|asset|url|dir)=.*$`)

func TestGetClogAndClogUpgradeResolveTheSamePlan(t *testing.T) {
	if testing.Short() {
		t.Skip("builds clog")
	}
	if runtime.GOOS == "windows" {
		t.Skip("get-clog.sh is a bash script")
	}
	bash, err := exec.LookPath("bash")
	if err != nil {
		t.Skip("bash not on PATH")
	}
	script, _ := filepath.Abs(filepath.Join(".github", "actions", "setup-clog", "get-clog.sh"))

	// A plain build carries no release repo, so clog's default is mrmxf/clog -
	// the same default get-clog.sh has.
	bin := filepath.Join(t.TempDir(), "clog")
	if out, err := exec.Command("go", "build", "-o", bin, ".").CombinedOutput(); err != nil {
		t.Fatalf("go build: %v\n%s", err, out)
	}

	cases := []struct {
		name    string
		env     map[string]string
		files   map[string]string
		wantErr bool
	}{
		{name: "pin file", files: map[string]string{".clog-version": "v1.3.5\n"}},
		{name: "pin file with comments", files: map[string]string{".clog-version": "# why\n\n  # more\n  v1.3.4  \nv9.9.9\n"}},
		{name: "CLOG_VERSION beats the file", env: map[string]string{"CLOG_VERSION": "v1.2.0"}, files: map[string]string{".clog-version": "v1.3.5\n"}},
		{name: "CLOG_VERSION_FILE", env: map[string]string{"CLOG_VERSION_FILE": "ci/pin"}, files: map[string]string{"ci/pin": "v1.0.9\n"}},
		{name: "CLOG_REPO", env: map[string]string{"CLOG_REPO": "someone/clog"}, files: map[string]string{".clog-version": "v1.3.5\n"}},
		{name: "no pin", wantErr: true},
		{name: "empty pin", files: map[string]string{".clog-version": "# none yet\n"}, wantErr: true},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			dir := t.TempDir()
			for name, body := range c.files {
				p := filepath.Join(dir, name)
				_ = os.MkdirAll(filepath.Dir(p), 0o755)
				if err := os.WriteFile(p, []byte(body), 0o644); err != nil {
					t.Fatal(err)
				}
			}
			env := []string{"PATH=" + os.Getenv("PATH"), "HOME=" + dir, "CLOG_INSTALL_DIR=" + filepath.Join(dir, "bin")}
			for k, v := range c.env {
				env = append(env, k+"="+v)
			}

			run := func(name string, args ...string) (string, error) {
				cmd := exec.Command(name, args...)
				cmd.Dir, cmd.Env = dir, env
				out, err := cmd.CombinedOutput()
				return strings.Join(planLine.FindAllString(string(out), -1), "\n"), err
			}
			sh, shErr := run(bash, "-c", `CLOG_PLAN=1 "$0"`, script)
			gp, goErr := run(bin, "Install", "clog", "--dry-run")
			up, upErr := run(bin, "upgrade", "--dry-run")

			if c.wantErr {
				if shErr == nil || goErr == nil || upErr == nil {
					t.Fatalf("all three must refuse: get-clog.sh %v, Install clog %v, upgrade %v", shErr, goErr, upErr)
				}
				return
			}
			if shErr != nil || goErr != nil || upErr != nil {
				t.Fatalf("get-clog.sh %v, Install clog %v, upgrade %v\nscript:\n%s\nclog:\n%s", shErr, goErr, upErr, sh, gp)
			}
			if sh == "" || sh != gp || gp != up {
				t.Errorf("plans differ\nget-clog.sh:\n%s\nclog Install clog:\n%s\nclog upgrade:\n%s", sh, gp, up)
			}
		})
	}
}
