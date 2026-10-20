---
name: ailang-world
description: Use when you are working inside an AILANG World episode, that is, when your tools are World's mcp__world__* tools (ailang-read, ailang-write, ailang-edit, ailang-delete, workspace-exec, ailang-check, ...) rather than a shell. How to work so the episode's record is complete and the operator can grade it.
---

# Working inside a World episode

Every act you take goes through World and lands in its log. The record is only complete,
and the episode only evidence-grade, if **every** act goes through World's tools.

- **World's tools, never a shell.** Read, write, edit, delete and run only through the
  `mcp__world__*` tools. Paths are relative to the episode worktree.
- **`workspace-exec` runs only the commands the episode's exec profile names** (`go-test`,
  `pytest-file`, ...). A command you need and do not have is a request to the operator; say
  so in your summary. Never work around it.
- **Read with `offset`/`limit`** on large files instead of reading them whole. Delete your
  scratch files with the delete tool before you finish.
- **An `ok:false` refusal is a fact to report**, not something to route around. Quote it and
  carry on with what you can do.
- **The exec cap is about 9 s per call.** Narrow `-run` (or `-k`) and one package or one test
  file per call. What you could not run, hand back and say so plainly.

## Finish with a summary the operator can grade

1. What changed, with `file:line` for each change.
2. The tests and checks you ran, and their results.
3. What you could not run or check, and why (a missing command, the cap, a refusal).
