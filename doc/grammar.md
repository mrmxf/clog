# The clog grammar (v1.0.0)

`clog <Namespace> [<sub>] <verb> [<noun>] [args]` — **the verb is always the
first word after the namespace path.**

v1.0.0 is a promise, so the names had to be worth promising. Every rename below
is **permanent and loud**: the old name still resolves, but only to an error
that names its replacement. There are no aliases and no silent compatibility,
because a mass rename is only safe when a missed call site fails where it runs
instead of quietly doing the wrong thing.

## Rule 1 — case says who owns the verb

**Capitalised = shipped. lowercase = yours.**

| you type | you get |
|---|---|
| `clog Build` | the shipped implementation — always available, never shadowed |
| `clog build` | your `.clog.yaml` override if you wrote one, else the shipped one |

- Every verb **responds to either case** when only one exists, so nothing forces
  you to hold shift.
- When both exist, **exact case wins**. That is the point: `build` is yours,
  `Build` is the reference you are testing against.
- A snippet that overrides a built-in **is always written lowercase**. The
  capitalised name is reserved so it stays reachable.

The fallback is a custom resolver (`snips.ResolveCase`), deliberately **not**
cobra's `EnableCaseInsensitive`: that global flag would make `build` and `Build`
collide and resolve by registration order, destroying the pair. The required
order is exact match → case-insensitive match against shipped names → unknown
command.

**Carve-out:** worker snippets (`bc-golang`, `bc-hugo`, `verb-args`, …) stay
lowercase-hyphenated. The `bc-` prefix already says "machinery, not for you",
they are not verbs anyone types, and `Bc-golang` reads badly.

## Rule 2 — the verb says what you get back

| class | verb | contract | caller writes |
|---|---|---|---|
| one value | `get` | prints **one** scalar; unset → nothing, exit 0 | `X=$(clog CI get ci.artifact)` |
| many values | `list` | prints **0..n** bare lines, no headers | `for t in $(clog CI target list)` |
| a document | `show` | prints JSON / `KEY=value` / a table | `eval "$(clog CI mode show)"` |
| a question | `is` `has` `should` `require` | **exit code only**, message to stderr | `if clog CI should deploy; then` |
| an action | `build` `deploy` `run` `install` | mutates the world | `clog build prod` |

Why `list` and not `print`: every printer prints, so `print` carries no
information. The question a caller needs answered is **cardinality**.

Why not `errifbad`: it names the *failure mode*, not the question, so it reads
as an implementation detail. Predicates read as questions.

## Retired — `CI`

| retired | canonical |
|---|---|
| `ci resolve` | `CI show event` |
| `ci policy` | `CI show policy` |
| `ci scan` | `CI show scan` |
| `ci mode` | `CI mode show` |
| `ci targets` | `CI target list` |
| `ci stack` (bare) | `CI stack list` |
| `ci stack get tools\|chk\|make\|names` | `CI stack list <which>` |
| `ci scan --list-targets` | `CI scan list targets` |
| `ci env` | retired outright (already deprecated) |
| `ci get` · `mode get` · `target get` · `stack get watch\|type` | unchanged — scalars |
| `ci should` · `require` · `run` · `deploy` | unchanged |

## Retired — `BC`

| retired | canonical |
|---|---|
| `BC releases version\|date\|flow\|note\|build` | `BC releases get <field>` |
| `BC releases yaml` | `BC releases get path` |
| `BC is <field> <value>` | `BC releases is <field> <value>` |
| `BC git branch\|suffix` | `BC git get branch\|suffix` |
| `BC git tag head\|origin\|prod\|ref` | `BC git tag get <which>` |
| `BC git hash head\|origin\|prod\|ref` | `BC git hash get <which>` |
| `BC git tree clean\|ahead\|behind\|unstaged` | `BC git tree is <state>` — **polarity flips** |
| `BC stash has error` | name unchanged — **polarity flips** |
| `BC stash get error` | name unchanged — **now prints to stdout** |
| `BC stashLog` | `BC stash log` |
| `BC genBuildinfo` | `BC gen buildinfo` |
| `BC semver <needs> <have>` | `BC semver satisfies <needs> <have>` |
| `BC linkerpath` | `BC get linkerpath` |
| `BC flow` · `git tag tidy` · `git checkout production` | unchanged — actions |

## The four behavioural bugs this fixed

1. **Inverted predicates.** `BC git tree ahead` exited **0 when NOT ahead**; same
   for `behind` and `unstaged`. `BC stash has` exited **1 when it HAD** an error.
2. **`not found` as a value.** `BC git hash prod|ref` printed it to **stdout**
   then exited 1, so `H=$(clog BC git hash prod)` yielded the string
   `"not found"`.
3. **`BC stash get error` wrote to stderr**, breaking the `get` contract:
   `X=$(…)` was silently empty.
4. **Help text named commands that failed.**

The polarity flips are the one change **not** protected by a loud error — the
name survives and only the exit code changes, so a missed caller goes green when
it should go red. Every call site was flipped in one wave, and
`sub-git-tree_test.go` asserts exit 0 **iff** the tree is ahead.

## What the retirement looks like

```console
$ clog CI policy
Error: `clog CI policy` was retired in v1.0.0
  why: a bare noun says nothing about whether it prints a value, a list or a document
  use: clog CI show policy
```

Permanent, by design. These are the "I told you that would bite you" reminder,
and they are what made renaming ~96 call sites in the embedded konfig — plus
every consumer `.clog.yaml` — a safe operation rather than a hopeful one.

A static test (`retired_config_test.go`) greps the shipped config for any
retired name and fails the build, which turns the runtime error into a
compile-time one for anything this repo ships.
