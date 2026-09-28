# claude-backlog.md — known issues

Not read every query. Read it before touching a listed file, or when picking up debt work.

Severity: **B** blocking / incorrect output · **M** moderate debt · **C** cosmetic · **L** long range

## Found 2026-09-28 — from the pihuw tidy

| # | Where | Issue |
|---|---|---|
| M-01 | `.github/workflows/build-check.yaml:146-160` | `kodata` is hard-coded in the artifact upload beside `_clog_build` and `_clog_deploy`. It is a legacy name from the ko container build, not a canonical clog dir, and it should come out. **Not yet safe:** as of 2026-09-28 four Hugo sites still set `publishDir: kodata` — `www-chiddingfoldbonfire`, `www-mrmxf-com`, `www-mrmxf-com-clog1` and `0_mrmxf/fohuw`. Each must move to `publishDir: _clog_build/public` (and its `.clog.yaml` target `dir:`/`ref:`) first, or its deploy silently publishes nothing. pihuw moved on 2026-09-28 and is the worked example. Removing the line is a CI-only change: it moves `workflows`, but a consumer pinned to `@vX.Y.Z` only sees it after a release. |
| L-01 | `.clog.yaml` `ci.artifact` + `build-check.yaml` | Long range: let a repo name extra dirs to carry from build to deploy, e.g. `ci.artifact: {name: pihuw-site, extra-dirs: [site/]}`, so clog stops hard-coding per-tool paths like `kodata`. Only worth it if a repo appears that genuinely cannot publish under `_clog_build/`. Until then the canonical-dirs-only rule is simpler and is the one to keep. |
