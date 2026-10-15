# Floor PreregTasks test selects preregistrations by kind

**Status:** Planned (iteration 269, dev-red repair)
**Planner-Lane:** opus
**Scope:** test-only; one file (`scripts/floor/test_final.py`)

## Problem

Dev CI job "go host build + test gate" is RED on origin/dev `f29f7e6`, `4e8299c` and `9b6f6a8`
(runs 38738979590, 38739309525, 38747016289). Parent `ba7e309` was green. The only failure is
step "Floor harness suite (row 93 M1+M3)" (`python3 -m unittest discover -s scripts/floor -p
'test_*.py' -v`, `FLOOR_REQUIRE_TOOL=1`): `test_final.PreregTasks.test_reads_the_draft_this_harness_writes`
ERRORs with `SystemExit: the preregistration names no corpus tasks (live_config.corpus.tasks)`,
raised by `scripts/floor/report.py` `prereg_tasks()`. That step is the job's last real step, so
nothing else is hidden behind it.

## Cause

The test selects files **by path alone**:
`design_docs/verification/world-floor-*/preregistration*.json`. Attended PR #269 (`f29f7e6`)
added `world-floor-smoke-2026-10-15/preregistration.json`, a SMOKE preregistration: it has no
`kind` and no `live_config` (its corpus lives under a top-level `corpus{commit,source,vs_76a5aef}`).
The glob matched it and `prereg_tasks` correctly refused it.

The other four matches (`world-floor-final-2026-10-09`, `-10-14`, `-10-15` and
`world-floor-m4-2026-10-06/preregistration.draft.json`) all carry
`kind == "row93-floor-preregistration"` and 23 tasks under `live_config.corpus.tasks`.

The test's stated subject is "the draft this harness writes". The harness defines that itself:
`scripts/floor/run.py:84` `PREREG_KIND = 'row93-floor-preregistration'`; `prereg_draft` writes it
(line 327) and the FINAL runner refuses anything else (lines 362-363:
`if pre.get('kind') != PREREG_KIND: raise RunRefused`). So the path glob was an approximation of
the subject; `kind` is the subject.

## Design (prescriptive)

Change only `PreregTasks.test_reads_the_draft_this_harness_writes`:

1. Keep the glob as the candidate set.
2. Load each candidate. If `pre.get('kind') == run.PREREG_KIND`, it is **selected**: assert
   `len(report.prereg_tasks(pre)) == 23` and `sorted(tasks) == sorted(pre['live_config']['corpus']['tasks'])`,
   with the path as the assertion message (unchanged checks).
3. Otherwise record it in `skipped`.
4. Non-vacuity controls:
   - `assertTrue(selected)` and `assertGreaterEqual(len(selected), 4, selected)` — the four known
     FINAL/M4 preregs. A filter that silently matches nothing (typo in the kind, moved directory)
     fails.
   - For every skipped path, assert `pre.get('kind') != run.PREREG_KIND` (restating the skip
     reason as an assertion so no other branch can skip). The kind is the ONLY skip reason: a
     file that claims `row93-floor-preregistration` and lacks `live_config` is selected and
     FAILS via `prereg_tasks`' SystemExit — it is never skipped.
   - Do NOT assert that the smoke path specifically is in `skipped` (that would couple the test
     to an evidence directory name).

`test_refuses_an_empty_prereg` is unchanged. No `skip`, no `xfail`, no `expectedFailure`.

**Do not modify:** `scripts/floor/report.py`, `scripts/floor/run.py`, any `.json` under
`design_docs/verification/`, `.github/workflows/ci.yml`.

## Rejected alternatives

- **(a) Narrow the glob to `world-floor-final-*`.** Silently drops the M4 draft
  (`world-floor-m4-2026-10-06/preregistration.draft.json`), which is literally the draft
  `prereg_draft` wrote, and couples the test's meaning to directory naming conventions that
  attended work chooses freely (the smoke dir proved that naming is not a contract).
- **(b) Make `report.prereg_tasks` accept the smoke shape.** Changes the contract of the FINAL
  verdict tool while a paid attended FINAL is executing from it. A verdict tool that accepts more
  shapes is a weaker refusal, and the smoke prereg is not an input `report.py` is ever given.
  Out of scope and risky.
- **(c) Delete or rename the smoke prereg.** It is attended, digest-sealed evidence; the loop
  never edits evidence to make a test pass.

## Acceptance criteria

- **AC1** `cd scripts/floor && python3 -m unittest discover -s . -p 'test_*.py'` passes locally
  (equivalently from repo root with `-s scripts/floor`). Record whether `FLOOR_REQUIRE_TOOL` was
  set; unset locally is acceptable (CI sets it to `1`) — state which in the evidence.
- **AC2 (test still bites)** Throwaway local edit: remove the `live_config` key from
  `design_docs/verification/world-floor-final-2026-10-09/preregistration.json` (keep `kind`),
  run `python3 -m unittest test_final.PreregTasks` → must FAIL/ERROR on that path. Then
  `git checkout -- design_docs/verification/world-floor-final-2026-10-09/preregistration.json`
  and show `git status --porcelain` is empty except the test change. Bank the output.
- **AC3 (control)** With only the test change reverted (`git stash` / checkout of the test file),
  the original `SystemExit ... live_config.corpus.tasks` ERROR reproduces; re-apply and it passes.
- **AC4** PR CI: "go host build + test gate" and "ailang-code verify gate" both SUCCESS.
- **AC5** `git diff --stat origin/dev...HEAD` touches exactly `scripts/floor/test_final.py`
  plus design/plan docs.

## Files

- `scripts/floor/test_final.py` — `PreregTasks.test_reads_the_draft_this_harness_writes` only.

## Verification Log

| # | Fact | Command | Result |
|---|------|---------|--------|
| V1 | Three dev runs red | `gh run view <id> --json headSha,conclusion` for 38738979590, 38739309525, 38747016289 | f29f7e6 / 4e8299c / 9b6f6a8 all `failure` |
| V2 | Glob matches 5 files | `ls design_docs/verification/world-floor-*/preregistration*.json` | final-10-09, final-10-14, final-10-15, m4-10-06 draft, smoke-10-15 |
| V3 | Kinds and task counts | `python3 -c` loop printing `kind`, `len(live_config.corpus.tasks)` | four `row93-floor-preregistration` 23; smoke `None` 0 |
| V4 | Harness kind definition | `grep -n PREREG_KIND scripts/floor/run.py` | 84 def; 327 written; 362-363 FINAL refuses other kinds |
| V5 | Local repro | `cd scripts/floor && python3 -m unittest test_final.PreregTasks` | `FAILED (errors=1)` |
| V6 | CI step + env | `grep -n "Floor harness suite" -A3 .github/workflows/ci.yml` | line 384, `FLOOR_REQUIRE_TOOL: '1'` |

All measured in worktree `iter269-devred` at `9b6f6a8` on 2026-10-15.
