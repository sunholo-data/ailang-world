You are working on Daneel (a virtual employee written in AILANG) through AILANG World's MCP tools ONLY. Your worktree root IS Daneel's `tools/` directory: paths are module files like `daneel_memory.ail`. You have no shell; use the world tools (ailang-read, ailang-write, ailang-edit, ailang-check, ailang-run, builtins-search, examples-search, ailang-cli).

TASK (Daneel issue #201): `std/string` already has native `length` and `find`, and several modules reimplement both — a string length helper (often `slen`, a `foldChars` over the whole string) and a substring search (`indexOf`/`scan`, a character walk). Swap in the native primitives. It is not a sed: `length` collides with `std/list.length`, which the same modules import, so each module needs an aliased import (e.g. `import std/string (length as strLength, find)` — first confirm the exact native names, signatures and return conventions with `builtins-search` and `ailang-cli` (`op: "iface", module: "std/string"`); in particular how `find` reports "not found" vs the local helper, so behaviour stays identical).

Candidate modules (no registry-package imports, so they can be checked here): daneel_brief.ail, daneel_candidates.ail, daneel_context.ail, daneel_evidence.ail, daneel_github.ail, daneel_memory.ail, daneel_reports.ail, daneel_sources.ail, daneel_tenant.ail, daneel_sync.ail. DO NOT edit daneel_docs.ail, daneel_intake.ail or daneel_share.ail (they import registry packages this environment cannot resolve) — just report whether they also contain reimplementations.

For EACH candidate module, one at a time:
1. `ailang-check` it first and note `passed` and the verify count — this is the baseline.
2. Read it; decide if it really reimplements string length and/or substring search. If not, skip it and say why.
3. Make the minimal edit: aliased import(s) + replace the helper's uses (remove the dead helper only if nothing else uses it). Preserve exact behaviour, including edge cases (empty strings, not-found results, offsets).
4. `ailang-check` again: it must pass and the verify count must not drop below the baseline. If a module has inline tests, run them via `ailang-cli` with `op: "test"` and the file path. If anything fails, fix or revert that module to its original content before moving on.

Finish with a short report: per module — changed / skipped (why) / reverted (why), baseline vs final verify counts, and anything you were unsure about. Do not attempt anything outside these files.
