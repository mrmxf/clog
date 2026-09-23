# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with
code in this repository.

## Project Overview

**`mrmxf/clog`** is clog: the CLI **and** the reusable workflows and actions
that build it, in one repository.

That pairing is the whole design, and it exists for one reason: **lazy, happy
users.** Before this repo, anyone who wanted clog had to be told three names —
the binary shipped from `mrmxf/clog-sample` (a name that says "example"), the
CI assets lived in `mrmxf/util/.github/` (a name that says "library"), and
`github.com/mrmxf/clog` was an archived husk. Now a user forks one thing and
gets their own clog and their own CI in the same tree, with nothing to wire up.

Its releases are what `setup-clog` downloads, so **every consumer repo's CI
bootstraps from a release cut here.** A broken release here breaks every
pipeline, not just this one.

## Build and Release Commands

### The verbs
- `clog build` — checks, lints, scans, then writes `_clog_build/artifacts/`
- `clog build --fast` — skips every gate; makes no claim about the branch
- `clog build prod` — **refuses a HEAD that is not a clean release tag**
- `clog CI deploy --target release` — checksums the assets and publishes
- Use `go run . <verb>` when iterating on this repo's own code

### Plain Go
- `go build -o clog .`, `go test ./...`, `go vet ./...`

### Cutting a release
1. Add the version to `releases.yaml` (newest first).
2. Set `.clog-version` to the same tag.
3. Commit, `git push`, then `git tag vX.Y.Z && git push origin vX.Y.Z`.
4. `self-release.yaml` does the rest on the tag push.

The tag must be at `origin` before the deploy runs: `gh release create` would
otherwise invent a release pointing at the wrong commit.

## Architecture

### Entry point
`main.go` embeds `releases.yaml` and layers config in an order that matters:

1. `util/embedfs` `konfig.yaml` — the base every clog app shares: the `bc-*`
   build workers, the verbs, the `check:` groups. Without it the binary has no
   build tasks.
2. this repo's `.clog.yaml`.
3. the **working directory's** `.clog.yaml`, via a deferred `kfg.AutoMerge()`.

Step 3 is why `PreventAutoMerge: true` is set and `AutoMerge()` is called by
hand afterwards. When this binary runs in *another* repo's CI, that repo's
`.clog.yaml` must win. Let Konfigure do the merge and it fires first, and this
repo's config silently overrides the config of the repo being built.

### Two tag namespaces, never merged

| tag | names | moves? |
|---|---|---|
| `vX.Y.Z` | a binary release | never |
| `workflows-v1` | the CI assets | deliberately, by hand |
| `workflows-v1.N.M` | a pinnable CI asset version | never |

**Rule:** `workflows-v1` moves only by deliberate action, only to a commit that
passed `test-actions.yaml`, and **never in the same commit as a Go change**. One
repo now has two lifecycles sharing a commit graph; this is what keeps them
apart. `workflows-*` also cannot match a `v*` glob, so the release trigger
cannot fire on a CI asset tag.

### `uses:` resolution — the trap that breaks consumers

Inside a **reusable workflow**, `uses: ./…` resolves against the **caller's**
checkout, not this repo. So:

- `build-check.yaml` → `clog-prepare` must stay **absolute**
  (`mrmxf/clog/.github/actions/clog-prepare@<ref>`), or every external consumer
  breaks.
- `self-build.yaml` → `build-check.yaml` may be **relative**: a caller resolves
  against its own repo, and this repo *is* the caller.
- `clog-prepare` → `setup-clog` may be relative: the composite-action → action
  hop is the one place it is allowed.

Because of the first rule, a PR editing `clog-prepare` would otherwise be tested
against the *ref* rather than against the change. That is what
`test-actions.yaml` is for — a plain workflow, which *can* use `./…` from the
working tree. Do not delete it as redundant.

## Gotchas

- **Do not add build logic to Go here.** This is a thin wiring `main.go`.
  Behaviour goes in `.clog.yaml` or upstream in `mrmxf/util`.
- **`../go.work` hides broken pins.** `/home/bruce/gr/clogs/go.work` puts util's
  local tree into the workspace, so a local build uses util's *source*, not the
  versions in `go.mod`. Before claiming a change builds, check it the way CI
  does: `GOWORK=off go build ./...`. `ci.yaml` asserts this in CI and
  `pre-build` refuses a CI build that a workspace is masking.
- **Do not upgrade `github.com/mrmxf/clog` in other repos.** That module path
  served `v0.9.5` with `/config`, `/ux/ui`, `/gommi` and `/slogger`, and is
  genuinely imported by `clog-mrmxf`, `utbd` and every `www-podserver`. v1.0.0
  is `package main` and has none of them, so `go get -u` resolves v1.0.0 and
  fails loudly. The proxy is immutable, so `v0.9.5` is served forever — pin it.
- A change needed in util must be tagged **per module** (`buildinfo/v0.13.0`,
  not `v0.13.0`) and pushed before this repo can use it — plan work in that
  order. There are no `replace` directives.
- `golangci-lint` and `trivy` must be on `PATH` for `clog build` to pass. clog
  refuses to install them behind your back — it names the tool and prints the
  install line (`clog Install <tool>`).
- `releases.yaml` is history, not policy. Versions come from git tags via
  `clog BC gen buildinfo`.
- `_clog_build/` and `_clog_deploy/` are build output and are gitignored.
