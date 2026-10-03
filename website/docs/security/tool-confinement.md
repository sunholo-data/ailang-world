---
title: Tool confinement
sidebar_position: 2
description: The measured confinement facts for the eight coding tools — worktree root, path refusals, .git protection, the case-folded deny list, environment, cache, flag injection and the run outcome shapes — with the design's V-row evidence.
---

# Tool confinement

Every fact on this page was measured first-party on the AILANG **v0.51.0** tool binary and
re-proved on **v0.52.1**, the only release the daemon now accepts for `--tool-ailang-bin`
(Verification Log §13, V69–V80) and is cited by its Verification
Log row (`V…`) in `design_docs/planned/w-software-engineering-domain.md`. Most are also
regression tests that run against the real binary (the AC4.1 matrix).

A different tool binary changes these facts. The daemon refuses any other release, and the
policy is to stay on the last binary that passed the matrix and never widen (R-SE-8).

## The sandbox root is the episode worktree

| Fact | Evidence |
|---|---|
| The rendered policy's `fs_sandbox` is `<workspace-root>/<episode>`; each handler runs with that directory as its working directory. | §4.3, V39 |
| The episode ID must match `^[a-z0-9][a-z0-9-]{0,63}$`, and `<root>/<episode>` must resolve to itself (no symlink redirect onto another directory or a sibling episode). Otherwise the episode gets an empty handler registry and every effect is refused before it runs. | §4.3, AC4.4 |
| The policy file, the tool cache, the store and its archive must be outside the workspace root, or startup refuses: an agent that can edit its own policy has no policy. | §4.3, AC4.5 |
| The policy file sits outside the sandbox and is passed by absolute path; the supervised worker re-reads it and admits the run. | V33, V39 |

## Paths that are refused

| Attempt | Result | Evidence |
|---|---|---|
| `..` out of the worktree (op) | `refused … escapes sandbox` | V16, V36 |
| Absolute path such as `/etc/hosts` (op) | `refused … escapes sandbox` | V16, V36 |
| Symlink inside the worktree pointing outside (op) | `refused … escapes sandbox` | V36 |
| A **running program** reading `../ep2/secret.txt` | rc 1, `execution failed: readFile: path "…" escapes sandbox` | V34, V39 |
| A path that is empty, absolute, contains a `..` segment, or starts with `-` | Refused by the tool's own plan before any effect (a contracted path predicate in each module) | §4.5, L-ARGS |

The path predicate in the transitions is for a clear refusal message. The confinement itself is
the policy layer's.

## `.git` is read-only

| Attempt | Result | Evidence |
|---|---|---|
| Op `write` to `.git/config` and to `.GIT/config` | `refused: .git/ is read-only to the lane's tools` (both spellings) | V36 |
| A running program writing `.git/config`, `.GIT/config`, `.git/hooks/pre-commit` | rc 1, `E_FS_PROTECTED: … is under .git, which is read-only in restricted mode`; files byte-unchanged | V52 |
| In a real `git worktree add` worktree, where `.git` is a pointer **file** (`gitdir: …`): program or op writes to `.git` / `.GIT` | Refused; the pointer file stays byte-unchanged, so the gitdir redirect cannot be rewritten | V54 |

`.git` is case-folded upstream in both layers; the `.GIT` fix is in v0.50.1 and later (V19,
V32).

## The deny list, case-folded

The rendered policy makes these read-only to the agent:

- directories `.github/`, `.pi/`, `.claude/`, `.ailang/`
- files `.gitmodules`, `.gitattributes`

| Fact | Evidence |
|---|---|
| Writes to `.github/**` and `.claude/**` are refused for ops (`matches fs_deny_write ".github/**"`) and for programs (rc 1 `E_FS_PROTECTED`). | V36, V52 |
| `.gitmodules` and `.gitattributes` were writable in a worktree until they were added to the list. | V54 |
| `.ailang/**` is denied so an agent cannot poison the compile cache later checks and runs read; `run` and `check` still work. | V57 |
| On a case-insensitive volume (macOS APFS), v0.51.0 matched the deny list **case-sensitively**: `.CLAUDE/settings.json`, `.GITHUB/…`, `.AILANG/…`, `.Pi/…`, `.GITMODULES`, `.GitAttributes` and `.gitmoduleſ` (U+017F) all wrote through to the protected lowercase files, by `write`, `edit` and running programs. | V68 |
| **Fix, in the policy only:** each directory is rendered as every Unicode simple-fold variant of its name plus `/**` (196 literals), and each file as one `path.Match` class over its fold orbit, for example `.[gG][iI][tT][mM][oO][dD][uU][lL][eE][sſS]`. 198 patterns, 3102 bytes. All case-variant rows are refused, with the lowercase targets unchanged. | V68, R-SE-15 |
| On v0.52.1 the deny matcher folds case upstream (one matcher, `fileguard.Protection`): with only the six lowercase entries, all 39 case-variant rows are refused. World still renders the folded list (redundant, harmless). policy-tool also refuses `fmt --write` on a deny-listed path, and `test` no longer writes into the tested file's directory. | V72, V73, V75 |
| The only non-ASCII code points APFS folds onto ASCII letters are `ſ` (U+017F, onto `s`) and the Kelvin sign (U+212A, onto `k`); no code point is ignored. None of the four directory names contains `s` or `k`. | V68 |

## Environment and cache

| Fact | Evidence |
|---|---|
| Both handler branches get a **minimal** environment, never the daemon's: `HOME=<db-dir>/cache/<episode>`, `PATH=/usr/bin:/bin`, `LANG=C`, `LC_ALL=C`, `AILANG_CACHE_DIR=<db-dir>/cache/<episode>`, and `AILANG_EXAMPLES` only when a corpus is configured. The registry credential cannot reach a tool. | `childEnv` in `host/broker/handlers_ailang.go` |
| `policy-tool` honours `AILANG_CACHE_DIR`; `run --policy` ignores it and writes its compile cache to `<worktree>/.ailang/`, which the deny list makes read-only to the agent. | V57, R-SE-4 |
| Because `HOME` is the per-episode cache, `examples-search` never finds the operator's `~/.ailang/examples` by itself; it uses only the corpus `serve --examples-dir` names, which must be outside the workspace root. Without one it refuses and never runs the op, which also stops the binary's CWD fallback from searching the agent's own worktree. | V65 |

## No shell, a fixed op set

| Fact | Evidence |
|---|---|
| There is no shell tool. Each effect name has a fixed op allowlist: `Workspace.Read` → `read`; `Workspace.Write` → `write`, `edit`; `Ailang.Check` → `ai_check`; `Ailang.Discover` → `builtins_list`, `examples_search`; `Ailang.CLI` → the policy summary's `cli` list; `Ailang.Run` → the run path only. | §4.3, V40 |
| `ailang-cli`'s allowlist is the policy summary's 20 `cli` ops: `agent_prompt`, `ai_check`, `axioms`, `builtins_list`, `builtins_show`, `check`, `devtools_prompt`, `docs_search`, `examples_list`, `examples_search`, `examples_show`, `examples_tags`, `fmt`, `iface`, `pkg_docs`, `policy_check`, `prompt`, `test`, `tree`, `version`. Not read, write, edit or run. | V40 |
| Per-op flag admission is `policy-tool`'s own: for example `fmt --check --write`, `test --allow-skips --json --no-color --package`, `check --json --quiet --strict-syntax`, `docs_search --json --limit`, `builtins_list --json`; most ops admit none. | V56, V67 |
| `policy-tool` silently ignores unknown top-level request keys, so each transition refuses any argument key outside its schema (L-ARGS). | V56 |
| An op with `write`/`edit` on `Workspace.Read` is refused by the handler, with no file change. | AC4.2 |

## Flag injection

| Fact | Evidence |
|---|---|
| `ailang run --policy P -policy.ail` parses the path as a flag (`flag provided but not defined`). | V58 |
| Closed twice: the plan's path predicate refuses a leading `-`, and the run handler always passes the path after `--` (`ailang run --policy P [--args-json J] -- <path>`). | §4.3, AC4.3 |
| `policy-tool` refuses a hyphen-leading `path` (`looks like a flag`), and sandbox-checks `--package`. | V67 |

## `ailang-cli` never writes

| Fact | Evidence |
|---|---|
| `fmt --write` rewrote deny-listed files (`.claude/x.ail`, `.ailang/x.ail`) through `policy-tool`, because the formatter does not consult `fs_deny_write`. Only the exact key `write` maps to `--write`, and its value is ignored (`"false"` writes). `FLAGS`/`Flags` alias `flags`, and duplicate `flags` objects merge. | V67 |
| **Fix, in both layers:** the `ailang-cli` plan refuses `fmt` with a `write` key in any `flags` entry (a Z3-verified predicate, zero-effect refusal), and the handler refuses the same plus any case-folded alias of `op` or `flags`, forwarding only its own re-encoding. `fmt` returns the formatted text; save it with `ailang-write`. | V67, R-SE-13 |
| Remaining upstream gap: `test` writes a transient `_namedtest_body_<random>.ail` into the tested file's directory, including a deny-listed one, then removes it. | V67, R-SE-13 |

## `ailang-run` outcome shapes

The run handler composes its result `{admitted, exit_code, decision, limit, stdout, stderr}`
from exactly four observed shapes, and records anything else as a handler failure. It never
infers success from the exit code.

| Shape | What the binary does | Evidence |
|---|---|---|
| (a) admitted | rc 0 or the program's rc; program output on stdout; one `policy:` JSON line on stderr | V33, V39 |
| (b) admitted, refused at runtime | rc 1; `execution failed: … escapes sandbox` plus the `policy:` line | V34 |
| (c) refused before execution | non-zero rc (2 for a policy violation or type error, 1 for an unreadable file); the decision JSON on **stdout**, stderr empty; `admitted:false`, with `error_kind` such as `policy_violation` or `read_failed` | V49, V63 |
| (d) supervisor limit | rc 3; partial stdout; the `policy:` line plus a `policy-result:` line with `reason` (for example `timeout`) and `stage` | V61 |

| Fact | Evidence |
|---|---|
| A refused inner effect can still exit 0: a program printed `E_NET_IP_BLOCKED … 127.0.0.1` and returned rc 0. | V35 |
| A program declaring `Net` is refused statically under the IO/FS policy (`missing_from_policy: ["Net"]`). | V49 |
| `security_mode` defaults to `restricted` if omitted; the rendered policy states it anyway, and a test pins the exact key set. | V62 |
| The policy's `timeout_ms` is 8000, below the handler's 10 s cap, so AILANG's supervisor kills first and reports shape (d). | §4.3, V61 |

## Not covered by the broker

The **inner** file operations of a program run by `ailang-run` are confined and budgeted by the
policy (`[budgets] FS = 1000`, `timeout_ms`), but they are not brokered, journaled or replayable
individually. That is an owner-ratified trade-off (D-SE-3 = A), recorded as residual R-SE-11.
