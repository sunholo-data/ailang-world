# Sprint plan — floor PreregTasks selects preregs by kind (iteration 269)

**Design:** `design_docs/planned/w-floor-prereg-test-kind-filter.md` (`7b771aa`) · **Milestones:** 1 · **Scope:** test-only
**Quorum:** PROCEED at 1 of 3 external reviewers — gemini-4-argon PASS; oc-kimi-k3, gpt6-1-sol, oc-glm-5-3 absent (429: ollama weekly limit / openai no credits; individually re-run, still absent)

**Quorum fold-in Q1 (gemini-4-argon, non-blocking):** the design's skip guard `kind != PREREG_KIND`
restates the branch condition and is tautological. Replace it: every skipped file must ALSO lack a
top-level `live_config` key (a real FINAL prereg whose `kind` got corrupted still carries
`live_config`, so it is caught), and `len(selected) >= 4`.

## M1 — kind filter in `PreregTasks.test_reads_the_draft_this_harness_writes`

Replace lines 196-204 of `scripts/floor/test_final.py` with exactly:

```python
    def test_reads_the_draft_this_harness_writes(self):
        drafts = glob.glob(os.path.join(run.REPO, 'design_docs', 'verification', 'world-floor-*', 'preregistration*.json'))
        self.assertTrue(drafts)
        selected, skipped = [], []
        for p in drafts:
            with open(p) as f:
                pre = json.load(f)
            if pre.get('kind') != run.PREREG_KIND:
                # Quorum fold-in Q1: a skipped file must not look like a FINAL prereg.
                self.assertNotIn('live_config', pre, p)
                skipped.append(p)
                continue
            selected.append(p)
            tasks = report.prereg_tasks(pre)
            self.assertEqual(len(tasks), 23, p)
            self.assertEqual(sorted(tasks), sorted(pre['live_config']['corpus']['tasks']), p)
        self.assertGreaterEqual(len(selected), 4, (selected, skipped))
```

Do not touch `test_refuses_an_empty_prereg`, `report.py`, `run.py`, any `.json` under
`design_docs/verification/`, or `.github/workflows/ci.yml`. No skip/xfail/expectedFailure. Do not
assert the smoke path by name. Commit: `fix(floor): PreregTasks selects preregs by kind (iter269)`.

## Gate commands (run verbatim, from the worktree root unless stated)

1. Focused: `cd scripts/floor && python3 -m unittest test_final -v` → `OK`.
2. Full floor suite: `python3 -m unittest discover -s scripts/floor -p 'test_*.py'` → `OK`.
   Record `echo "FLOOR_REQUIRE_TOOL=${FLOOR_REQUIRE_TOOL-unset}"`; unset locally is acceptable
   (CI sets `1` with a v0.52.1 tool binary).
3. No Go/AIL code is touched, so `go build`/`verify_ail.sh` are not required locally; CI runs them on the PR.

## Acceptance criteria (executable)

Let `F=design_docs/verification/world-floor-final-2026-10-09/preregistration.json`.

- **AC1** Gates 1 and 2 print `OK`; evidence states the `FLOOR_REQUIRE_TOOL` value.
- **AC2 (test bites + Q1 bites)**
  - (a) `jq 'del(.live_config)' $F > /tmp/x269 && cp /tmp/x269 $F`; run gate 1 → must FAIL/ERROR
    (`SystemExit ... live_config.corpus.tasks` on `$F`). `git checkout -- $F`.
  - (b) `jq '.kind = "corrupted"' $F > /tmp/x269k && cp /tmp/x269k $F`; run gate 1 → must FAIL
    with `AssertionError: 'live_config' unexpectedly found in {...}` naming `$F` (proves Q1).
    `git checkout -- $F`.
  - After both: `git status --porcelain` shows nothing beyond the intended test/plan changes
    (`$F` absent). Bank all three outputs.
- **AC3 (revert control)** `cp scripts/floor/test_final.py /tmp/new269.py &&
  git show origin/dev:scripts/floor/test_final.py > scripts/floor/test_final.py`; run gate 1 →
  original `SystemExit: the preregistration names no corpus tasks (live_config.corpus.tasks)`
  ERROR reproduces. `cp /tmp/new269.py scripts/floor/test_final.py`; gate 1 → `OK`;
  `git diff --quiet HEAD -- scripts/floor/test_final.py` (rc 0) after the M1 commit.
- **AC4** PR CI: "go host build + test gate" and "ailang-code verify gate" both SUCCESS.
- **AC5** `git diff --stat origin/dev...HEAD` lists exactly `scripts/floor/test_final.py`,
  the design doc and this plan.

## Risk

Low. Selected count today = 4 (V3); the smoke prereg has no `live_config` (V3), so Q1 passes on
current evidence. A future non-FINAL prereg carrying `live_config` will fail loudly — intended.
