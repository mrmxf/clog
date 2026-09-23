# Workflows and actions

Everything here builds clog, and clog is built by it. The repo carries both so
that forking it gives you a working CLI and a working pipeline in one tree.

## What a consumer uses

| asset | what it does |
|---|---|
| `build-check.yaml` | reusable — **builds and gates**, everything into `_clog_build/` |
| `deploy-probe.yaml` | reusable — **publishes and probes**, everything into `_clog_deploy/` |
| `actions/clog-prepare` | checkout → Go → clog → `ci.policy` into `$GITHUB_ENV` |
| `actions/setup-clog` | install a pinned, checksum-verified clog release |

```yaml
jobs:
  build:
    uses: mrmxf/clog/.github/workflows/build-check.yaml@workflows
    permissions: {contents: read, id-token: write, security-events: write}

  deploy:
    needs: [build]
    uses: mrmxf/clog/.github/workflows/deploy-probe.yaml@workflows
    permissions: {contents: write, id-token: write}   # write for a GitHub target; read otherwise
```

Pin to `workflows` (moving: the newest CI that passed `test-actions.yaml`) or
to a release tag `vX.Y.Z` (fixed). A tag names a commit of the whole repo, not
of a folder, so a release tag pins the workflows exactly as well as the binary;
there is no separate CI version number to track. The moving tag is not called
`v1` because `self-release.yaml` fires on `v*`.

## What only this repo uses

| workflow | why |
|---|---|
| `self-build.yaml` | this repo's CI, calling `./build-check.yaml` with `self-build: true` |
| `self-release.yaml` | `v*` tag → binaries + `checksums.txt` + `get-clog.sh` |
| `test-actions.yaml` | PR gate exercising `./.github/actions/*` **from the working tree** |
| `ci.yaml` | `GOWORK=off` Go build/vet/test, and actionlint |
| `dump-context.yaml`, `test-setup-clog.yaml` | diagnostics, run on demand |

## Two rules worth knowing before you edit

**1. `uses:` resolves differently in a reusable workflow.** Inside
`build-check.yaml` and `deploy-probe.yaml`, `uses: ./…` resolves against the
**caller's** checkout — so the reference to `clog-prepare` is absolute on
purpose. A caller (`self-build.yaml`) may use `./…`, because it resolves against
its own repo. The same is true inside a composite action, so `clog-prepare`'s
reference to `setup-clog` is absolute too; `ci.yaml` enforces that.

That rule is why `test-actions.yaml` exists: a PR editing `clog-prepare` would
otherwise be tested against the pinned ref instead of against the change. It is
a gate, not a convenience.

**2. `workflows` moves by hand, never alongside a Go change.** One repo, two
lifecycles. Move the tag only to a commit that passed `test-actions.yaml`.

## Scan uploads

Since 2025-07 GitHub rejects several SARIF runs sharing one category, and
`upload-sarif` takes one category per call. `build-check.yaml` therefore lists
the SARIF files and publishes them from a matrix, one category each, with
`fail-fast: false` so one malformed sweep cannot hide the rest.
`workflows_test.go` asserts that shape — collapsing it back into a single upload
looks tidier and breaks every scan upload.
