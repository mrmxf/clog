# Plan — built-in Cloudflare helpers, and `clog Install clog`

Status: **proposed**, 2026-09-28. Nothing here is built yet except the first
helper (`EnsurePagesProject`, below), which shipped inside the Pages deployer in
v1.2.1.

## Why

Every Cloudflare site carries the same shell in its `.clog.yaml`: bonfire's
~300-line `deploy`, www-mrmxf-com's `bc-cf-api`, `bc-cf-zone-id`,
`bc-cf-pages-subdomain`, `bc-cf-dns-state` and `cutover`. They are generic, so by
the ownership rule they belong in util. They also drift. Bonfire created a
missing Pages project; the clog deployer refused to, until v1.2.1 — so a site
moving from its own snippets to the shipped verb silently lost a behaviour it
relied on.

The second part is smaller and separate. There is no way to make the **running**
clog a given version. `get-clog.sh` installs one, but only by piping a script
into bash, and `clog Update` belongs to the private edition.

## Part 1 — `clog CF`: built-in Cloudflare helpers

### Grammar

A new shipped namespace, `CF`, following [grammar.md](grammar.md). The verb says
what you get back:

| command | class | does |
|---|---|---|
| `clog CF pages ensure [project]` | action | create the Pages project if missing (production branch from the target) |
| `clog CF worker ensure [name]` | action | create the Worker if missing — see open question 1 |
| `clog CF pages get subdomain [project]` | one value | the project's real `*.pages.dev` host (never assume `<project>.pages.dev`) |
| `clog CF zone get id <domain>` | one value | the zone that owns a name, longest suffix wins |
| `clog CF dns get state <domain> [target]` | one value | `nozone` \| `missing` \| `ok` \| `foreign …` (A/AAAA/CNAME only) |
| `clog CF pages is attached <domain>` | question | exit 0 iff the domain is attached to the target's project |
| `clog CF cutover [dev\|prod]` | action | move a domain onto its Pages project — today's site `cutover` snippet |
| `clog CF api <METHOD> <path> [body]` | a document | a raw v4 call, the escape hatch |

- **Defaults come from the target.** With no argument, a `pages` verb reads
  `ci.targets.<CLOG_TARGET or the only cloudflare-pages target>.<mode>`:
  `project`, `domain` and `production-branch`. A site writes no names twice.
- **Credentials** are `$CLOUDFLARE_API_TOKEN` and `$CLOUDFLARE_ACCOUNT_ID`, as
  wrangler reads them. On a laptop an action re-runs itself under `clog CI run`,
  the same way `clog deploy` does.
- **`--dry-run`** on every action. It prints the calls it would make and changes
  nothing.

### What may happen by default, and what never does

| act | default? | why |
|---|---|---|
| create a Pages project / Worker | **yes**, on deploy | a new project serves only `*.pages.dev`, so no live site can move |
| attach a custom domain | **no** — `CF cutover` only | attaching can take a domain from whatever serves it today |
| write or delete DNS | **no** — `CF cutover` only | same reason, and deleting records is the one act that takes a site down |

This is the line bonfire's snippet drew with its "foreign DNS" guard, now drawn
in one place. `CF cutover` keeps every safety the site snippet has:

- it shows the records it would delete, and asks you to type the domain;
- it refuses without a terminal unless `CLOG_CUTOVER_YES=1`;
- it attaches before it touches DNS;
- it leaves MX, TXT, SPF and DKIM alone.

It also gains one step the snippet lacks. When the domain is attached to a
**different** Pages project (the `mrmxf-staging` → `www-mrmxf-staging` move), it
detaches it there first, after the same confirmation.

### Where the code goes

- `util/ci/cloudflare.go` holds the API client. `cfCall`, `cfResponse` and
  `EnsurePagesProject` exist already.
- New functions: `PagesSubdomain`, `ZoneID`, `DNSState`, `AttachDomain`,
  `DetachDomain`, `EnsureWorker`, `Cutover`.
- One function per act, shared by the command **and** the deployers. The
  cloudflare-pages deployer already calls `EnsurePagesProject`. A future
  cloudflare-worker deployer calls `EnsureWorker`.
- `util/ci/cf-command.go` holds the cobra `CF` namespace, mounted in clog's
  `main.go` beside `CI` and `BC`. The strict unknown-word rule from v1.0.3
  applies to it.
- Tests use `httptest` against `cfAPIBase`, as `cloudflare_test.go` does now.
  No test reaches Cloudflare.

### Migration

1. **util `ci/v0.26.0`** — the `CF` namespace, the helpers and `cutover`.
2. **clog v1.3.0** — mounts `CF`. Its `retired_config_test.go` gains the site
   snippet names, so a consumer still calling `clog bc-cf-api` from a
   hand-written script gets pointed at `clog CF api`.
3. **www-mrmxf-com** — delete `bc-cf-*` and `cutover` from `.clog.yaml`, and
   change `check.deploy-ready` to call `clog CF …`. The site then has no
   Cloudflare shell at all.
4. **bonfire** — drop its own `build`/`deploy` for the shipped verbs plus a
   `pages` target. That is the same move www-mrmxf-com made in v1.1.
5. **cf_tools Workers** (both sites) — once `CF worker ensure` and a
   `cloudflare-worker` target kind exist, `cf_tools/deploy.sh` becomes a target
   block. That is a later plan; this one only makes the helper exist.

## Part 2 — `clog Install clog [vX.Y.Z|latest|pinned]`

Make the running clog a specific version by reinstalling it in place.

```console
$ clog Install clog              # latest release of mrmxf/clog
$ clog Install clog v1.2.0       # exactly that one - downgrades allowed
$ clog Install clog pinned       # whatever this repo's .clog-version says
```

### Behaviour

1. **Resolve the version.**
   - `latest` (the default) follows `releases/latest`, as `get-clog.sh` does.
   - `pinned` reads the first non-comment line of `./.clog-version`.
   - Anything else must be a `vX.Y.Z` tag.
2. **Download** `clog-<cpu>-<os>` and `checksums.txt` from that release. The
   asset names come from the same convention `get-clog.sh` uses (one Go function
   that the script's convention is tested against).
3. **Verify** the binary's sha256 against `checksums.txt`, and refuse on any
   mismatch or missing entry. `CLOG_SKIP_VERIFY` is not honoured here, because a
   self-replacement is the worst place to skip it.
4. **Replace** `os.Executable()` atomically:
   - write to a temp file in the **same directory**, `chmod 0755`, then `rename`
     over the old binary. A running Linux or macOS process keeps its open inode,
     so this is safe mid-run.
   - If the directory is not writable, say so and print the `sudo` or
     `CLOG_INSTALL_DIR` alternative. It never escalates by itself.
5. **Report** `v1.1.1+basic → v1.2.0+basic` and run the new binary's `version`
   to prove it starts.
6. **No-op** when the running version already equals the target. `--force`
   reinstalls anyway.

### Guards

- **Edition.** A `+mrmxf` (clog-mrmxf) binary must not be replaced by the public
  `+basic` one, because the private commands would silently disappear. Refuse
  unless `--repo` names the matching repo, or `--force` is given.
- **A dev build** (`+dev` version, `go run`) refuses, because replacing
  `go-build/…/exe/clog` is meaningless. Say to run `go install` instead.
- **Windows:** no Windows release is built, so the command isn't offered there.

### Where the code goes

- `util/install`, as a special case of `runInstall` for the tool name `clog` (a
  new strategy `self`, recipe `recipes/clog.yaml`). `clog Install list` then
  shows it, and `clog Install have clog` works for free.
- The repo default (`mrmxf/clog`) lives in one Go constant, which `get-clog.sh`'s
  default is tested against.
- Ships with util `install/v0.18.0` → clog v1.3.0, beside Part 1.

### Tests

- A fake release server (httptest) serves assets and `checksums.txt`.
- **Good checksum:** the binary is replaced and the old one is gone.
- **Bad checksum:** refused, and the original binary is byte-identical afterwards.
- **Read-only install directory:** refused with the sudo hint, nothing changed.
- **Edition mismatch:** refused without `--force`.
- **`pinned`:** reads `.clog-version` with its comment header.

## Open questions

1. **What does "ensure a Worker" create?**
   - Option (a): an empty placeholder script, so routes and secrets can be set
     before the first real deploy.
   - Option (b): nothing new — `wrangler deploy` of the Worker's own directory,
     done only when the Worker is missing.

   (b) never deploys code nobody built, but it makes `ensure` a deploy. Proposal:
   (b), named `CF worker ensure --from <dir>`.
2. Should `CF cutover` also retire the container-registry target's hook rows in
   utbd? Proposal: no. It prints the reminder, because utbd is private and clog
   is public.
3. `clog Install clog` on a CI runner: allowed (it's just a tool install), or
   refused in favour of `.clog-version`? Proposal: allowed, with a warning that
   the pin is the source of truth.
