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
    uses: mrmxf/clog/.github/workflows/build-check.yaml@workflows-v1
    permissions: {contents: read, id-token: write, security-events: write}

  deploy:
    needs: [build]
    uses: mrmxf/clog/.github/workflows/deploy-probe.yaml@workflows-v1
    permissions: {contents: read, id-token: write}
```

Pin to `workflows-v1` (moving) or to a `workflows-v1.N.M` (fixed). Those tags
never collide with the `v*` release tags — different namespace, and `workflows-*`
cannot match a `v*` glob.

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
its own repo. The composite-action → action hop inside `clog-prepare` may also
be relative.

That rule is why `test-actions.yaml` exists: a PR editing `clog-prepare` would
otherwise be tested against the pinned ref instead of against the change. It is
a gate, not a convenience.

**2. `workflows-v1` moves by hand, never alongside a Go change.** One repo, two
lifecycles. Move the tag only to a commit that passed `test-actions.yaml`.

## Scan uploads

Since 2025-07 GitHub rejects several SARIF runs sharing one category, and
`upload-sarif` takes one category per call. `build-check.yaml` therefore lists
the SARIF files and publishes them from a matrix, one category each, with
`fail-fast: false` so one malformed sweep cannot hide the rest.
`workflows_test.go` asserts that shape — collapsing it back into a single upload
looks tidier and breaks every scan upload.
