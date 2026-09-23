# clog

One vocabulary for building and deploying, on whatever runs your CI.

`clog` gives a repository one set of verbs — `build`, `deploy`, `Check`, `CI` —
that behave the same on a laptop, on GitHub Actions and on GitLab CI. What a
repo *is*, and what its CI is *allowed to do*, are declared once in
`.clog.yaml`. The workflows read that file rather than repeating it, which is
the only arrangement in which CI and a laptop cannot drift apart.

**This repository carries both halves**: the CLI, and the workflows and actions
that build it. Fork it and you have your own clog and your own CI in one tree,
with nothing to wire together. That is the point.

## Install

```bash
CLOG_REPO=mrmxf/clog bash <(curl -fsSL \
  https://github.com/mrmxf/clog/releases/latest/download/get-clog.sh)
```

The installer verifies the binary's sha256 against the release's
`checksums.txt` and refuses to install one it cannot verify.

In a GitHub workflow:

```yaml
- uses: mrmxf/clog/.github/actions/setup-clog@workflows-v1
  with:
    clog-ref: v1.0.0        # optional; default reads .clog-version
```

## The verbs

```bash
clog build              # checks, lints, scans, then writes _clog_build/artifacts/
clog build --fast       # skips every gate; makes no claim about the branch
clog build prod         # refuses a HEAD that is not a clean release tag
clog CI deploy          # publish to each ci.targets entry for this run's mode
clog CI probe           # ask the LIVE targets what they look like from outside
clog CI show policy     # what this run may do, and why
```

### Case says who owns a verb

| you type | you get |
|---|---|
| `clog Build` | the shipped implementation — always available, never shadowed |
| `clog build` | your `.clog.yaml` override if you wrote one, else the shipped one |

Where only one of the pair exists, either spelling works, so nothing forces you
to hold shift. Where both exist, exact case wins — which is what lets you A/B
your override against the default it replaces instead of replacing it blind.

### The verb says what you get back

| class | verb | contract |
|---|---|---|
| one value | `get` | prints **one** scalar; unset → nothing, exit 0 |
| many values | `list` | prints **0..n** bare lines, no headers |
| a document | `show` | prints JSON / `KEY=value` / a table |
| a question | `is` `has` `should` `require` | **exit code only**, message to stderr |
| an action | `build` `deploy` `run` `install` | mutates the world |

Full inventory, including every name retired in v1.0.0 and what replaced it:
[doc/grammar.md](doc/grammar.md).

## Where things live

```
_clog_build/          what a build produces
├── artifacts/        binaries, rendered site — what a deploy publishes
├── check/            stash, lint, vet, test output
└── scan/             SARIF from the build-time sweeps

_clog_deploy/         what a deploy produces
├── receipts/         what was published, where, at what digest
└── probe/            findings against the LIVE target
```

Both are gitignored. Nothing else is guessed at: a debugger looks rather than
having to know.

## Configuration

`.clog.yaml` in the repo root, and nothing else. `clog CI --config-help` prints
the whole schema with examples.

## Licence

BSD-3-Clause. See [LICENSE](LICENSE).
