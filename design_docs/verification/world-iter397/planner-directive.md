MISSION-ROLE: planner

You are the sprint-planner for Ailang World mission iteration 397 (unattended run; no human is present).
Write the sprint plan for **row 243 `w-world-address-bar-verbs-v2`, sprint V2a — "the seam" (M1–M3)**. No re-design:
`design_docs/planned/w-world-address-bar-verbs-v2.md` r1 is reviewed (quorum r1 BLOCKED 3/3 → r1 revision → quorum r2
PROCEED 3/3, no absentee; independent sonnet audit PASS 84, 0 blocking; iteration 396) and is NOT reopened. Build V2a
as designed. **V2a is behaviour-neutral for the CLI and leaves the bar's `why` at `terminal`** (design §0.4): no route,
page, `/v1`, `.ail` or user-visible change. V2b (the bar: `AdmitWhy`, the daemon getter, the route's `why` branch, the
page) is a LATER sprint and is out of scope — do not build any of it.

## Scope (design §0.4, §2.1, §2.2, §3, §6 "V2a — the seam", §7, §8 rows M1–M3 — read them verbatim)
- **M1** — `git mv cmd/ailang-worldd/why.go host/why/why.go` + the host/why edits: the exports, `Getter`, `Run`,
  `ParseTarget`, `Elide` (§2.1). Keep host/why's diff from the original MINIMAL (no reflow, no comment rewrites) so git
  pairs the rename (§8: rename pairing needs > 50 % similarity; W27 measured `56 209` in a rough rehearsal, the audit
  `75 228`). **Size check:** after M1's commit `git diff --numstat <base> HEAD -- cmd host docs` must show
  `{cmd/ailang-worldd => host/why}/why.go`; if pairing is lost, stop and re-cut per §8 before going further.
- **M2** — `cmd/ailang-worldd/why_cli.go` (the CLI front + compatibility layer), `mcpclient.go` (`elide` delegation and
  the **[r1]** `headState` → `why.Head` delegation: cmd's `headState` body becomes a single `return why.Head(…)`), the
  two unit-test files that move with the code.
- **M3** — `WhyFlags` & friends in `host/verbs` + the CLI's use of them; tests VW1–VW5, VA1, VS2a.
- **Gate (design §0.4): every existing CLI test passes UNEDITED**, except the test files that move with the code (§7,
  as re-cut by audit input 2 below). The design's §6 "Must-live, unedited" list is binding.

## The seven audit/quorum inputs (design §13.2) — resolve each BY MEASUREMENT and record how in a §1.3/§1.4 row
1. `sum` (`why.go:774`) and `whyTarget{kind:…}` (`why.go:297`) are used by `demo_seed.go` (`:658,730,801,829,998,
   1064`) — add them to the compatibility layer, or list `demo_seed.go` as edited. Your call, recorded with the reason.
2. The two "moving" test files are not package-line-only: `why_protected_test.go:52,59,74` uses cmd's rig and
   `TestWhyProtectedLine` mixes pure subtests with an `e2e-exec` one (so it SPLITS); `why_input_test.go:31,58` call
   `elide(p)` (→ `why.Elide`). Re-cut §7 and VW1 accordingly and say exactly which lines of which pre-existing test
   files change, and why each change is a move rather than an edit of an assertion.
3. `c.get` collides with `func (c *client) get(path string) (int, string, error)` at `cli.go:94` (not a `Getter`). Use
   a closure or a new name; list `cli.go` if it is edited.
4–6 are V2b inputs (VS2's `d.*` set, the stale `lookup.go:99` reason, the cap's JSON-escape note) — **out of scope for
   V2a; do not act on them** beyond noting in Appendix C that they carry to V2b.
7. §2.1's call-site count (7 + `why.go:148`); `host/why` reaches `host/store` transitively via `host/broker`/
   `host/coordinator` — VS2a pins DIRECT imports only, so say so; size from the real patch, not the rehearsal.

## Where you work
- Sprint worktree (already created, branch `sprint/world-iter397-row243-v2a`, base = origin/dev
  `18d30a17a058e9adad71c42b6078b538374378d7`): `/Users/voightkampff/dev/sunholo-data/.wt-world-iter397-row243-v2a`
- The design was measured at `200ea5d5`. VERIFIED BY CONTROLLER: `git diff --stat 200ea5d5 18d30a17 -- cmd host docs`
  is EMPTY (only design docs changed since). Still re-measure every design line you rely on, live, at `18d30a17`.
- Pinned ailang: `$HOME/.pinned-ailang/ailang` = AILANG v0.41.0 (the Go suite's `AILANG_BIN`). V2a changes no `.ail`,
  so `./scripts/verify_ail.sh` must give the base's result.
- Rig load varies (it has been ≈ 200). Re-run a red in an untouched package once, alone, before believing it.
- Read the newest landed plan for the FORMAT and reproduce it: `design_docs/planned/w-world-address-bar-verbs-v1-sprint-plan.md`
  (iteration 395; §1.2 premise table, §1.3 decisions, §1.4 false premises, §2 patch protocol + awk extraction,
  §6 gate script, §7 executor notes, Appendix P fences, expected numstat per sub-milestone, Appendix C) and its JSON
  `design_docs/planned/sprint_w-world-address-bar-verbs-v1.json`.
- Scratch (compile your stages, run mutants) OUTSIDE the repo: `$HOME/.ailang/state/world-iter397/plan-scratch/`
  (exists; never under /tmp — location-sensitive tests go red there).
- Write exactly two files in the worktree, and DO NOT git add/commit them (the controller commits):
  - `design_docs/planned/w-world-address-bar-verbs-v2a-sprint-plan.md`
  - `design_docs/planned/sprint_w-world-address-bar-verbs-v2a.json`
- Leave every other file in the worktree untouched (the executor applies your patches there). No git writes.

## The move, the patches and the snapshot contract (new this sprint — read carefully)
- The executor's only writer is `patch -p1` (macOS `/usr/bin/patch`; check its version and whether it honours git
  rename headers — do NOT rely on them). Express the move as a whole-file DELETION diff of `cmd/ailang-worldd/why.go`
  plus a whole-file CREATION diff of `host/why/why.go` (and the same for each moved test file), so `patch -p1` alone
  reproduces the tree. Git detects the rename at diff time, so the commit still pairs if the content is similar enough.
  Prove it: in your round trip, commit base and patched tree in a scratch git repo and show `git diff --numstat`
  pairs `{cmd/ailang-worldd => host/why}/why.go`.
- Executor snapshots: cumulative, `$HOME/.ailang/state/world-iter397/snap/<M>/` holds every file added or changed
  since base (repo-relative paths), plus a `DELETED` file listing, one per line, every repo path removed since base,
  plus `MANIFEST.sha256` (exclude the two plan files). The controller's commit script copies the snapshot over the
  worktree and `git rm`s each `DELETED` path. Specify this in §7 and give the exact snapshot commands.

## Size rule
- **Cap: ≤ 1,300 changed lines** for V2a, measured `git diff --numstat <base> HEAD -- cmd host docs` with git's
  default rename detection (state this exact filter). Design estimate ≈ 970 (≈ 1,060 calibrated). A rename-paired file
  counts only its changed lines. If over, take the design's cut (§8): M3's tests before M3's CLI repointing (the
  repointing is behaviour-neutral and rides V2b). Prefer a smaller, fully proven sprint over a full, half-proven one.
  **Plan V2a only** (design Q1 default: two sprints); do not fold V2b in.
- Sub-milestones of ≤ ~250–400 changed lines each (the rename-paired move may be its own sub-milestone), each
  independently green and safe at commit (your call, measured).

## Retro directives (iterations 374–396; each was paid for)
1. **Ship compiled patches**: every line of code, test and doc as a unified diff in Appendix P
   (```` ````diff <M> ```` fences), generated LAST from compiled, vetted, gofmt-clean, green scratch stages; then
   extract them FROM THE PLAN FILE with §2's awk, apply to a fresh `git archive 18d30a17` export, and `diff -r`
   against your work tree — byte-identical or regenerate. The controller re-does this round trip independently.
   **Ship the patches BEFORE running the full mutant set, and checkpoint the plan file** (write a first complete
   version, then append the mutant ledger).
2. **Expected numstat after each sub-milestone's commit** (measured after commit, `-- cmd host docs`, renames on).
3. **Mutants anchored to the DIFF**: one mutant per new or changed production branch, enumerated from the diff, not
   from the test plan — in particular: `Run`'s exit/stream mapping (0 / 1 / 3, stdout vs stderr), the `Getter` 404 vs
   error classes, `ParseTarget`'s forms, `Head`, `Elide`, every exported constant (VW4), `WhyFlags`' registrations and
   defaults, the CLI adapter's error lines, and every compatibility-layer shim. Run each against the named test row AND
   the package; record killed/survived; name equivalents with reasons. A mutant that does not compile is not a mutant
   (pair with a killed control). **Name the test that killed each mutant; a mutant killed only by an
   inventory/static/AST test is not behavioural proof — re-run with that test excluded and record both.** Write the
   ledger to a file in scratch; every re-run appends its own line. **The plan's mutant table ids must be the ledger's.**
4. **Behaviour neutrality is the deliverable**: for a seam sprint the proof is that the CLI's observable behaviour is
   byte-identical. For every existing `why` test and every caller of the moved identifiers (W3, W4, W35, W39), trace
   input → output and say which existing test drives it end to end. Name any path through the moved code that NO
   existing test drives (that is where a refactor can change behaviour silently) and either add a characterisation test
   pinned to the BASE's output (generated by running the base binary/test rig, never hand-typed) or name it a residual.
5. **Repo-wide inventories**: run every design §7 inventory and any test that enumerates packages, imports, exported
   symbols, files, CLI verbs or context roots, by name with `-v`, on your final scratch tree (incl.
   `TestDaemonDependencyAllowlist`, `TestWorldBoundaryDependencyAllowlist`, `TestProductionContextRoots`,
   `TestProductionGoSurface`, `TestEveryVerbHasHelp`, VS1 in `host/daemon/surface_verb_test.go`). Search for any such
   test before assuming none — a new package `host/why` will meet every package allowlist. If the diff adds a production
   context root, `go test ./host/store/ -run '^TestProductionContextRoots$'` is a named gate (charter Repo Profile).
6. **Fixtures from the real tool**: VW2's replay `Getter` bodies must be RECORDED from the walk rig's real routes in
   the test (design §6), never typed. Any characterisation output likewise.
7. The §6 gate script for the controller (one bash shell, `tee` every leg to `$HOME/.ailang/state/world-iter397/gates/`,
   never only through `tail`), incl. `go test -run '^$' ./...` as the test-file compile fence, `-race` on touched Go
   packages, the inventories, `./scripts/verify_ail.sh`, and
   `VERIFY_GO_LEGS=plain AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh` after commit. Measure each gate
   at the base first. **Never start a gate under `nohup`.**
8. Executor notes as the V1 plan's §7: no git writes; `patch -p1` is the only writer; runs of ~2–4 sub-milestones;
   cumulative snapshots per the contract above with MANIFEST.sha256; stop conditions; numstat checks; exact acceptance
   commands per sub-milestone; **an explicit "never `nohup`, never background a test leg; pass a 600000 ms timeout to
   every long Bash call" line**.
9. Acceptance criteria (§9-style), numbered, that the judge can score against — including the behaviour-neutrality
   criterion and the exact list of pre-existing test files changed (with the reason for each).
10. **Every premise gets a §1.2 row**: any sentence saying "the code currently does X" cites the command and its
    observed output; an empty result needs a known-positive control in the same call. Where the design and HEAD
    disagree, measure, route minimally, record it in §1.4, and list every file you touch outside the design's §7.
    **For every type, field or function your patches name, cite where it exists at HEAD or mark it new** (396 retro).
11. Close with Appendix C: anything moved to V2b (with size), inputs 4–6, and the V2b handoff.

## Hard constraints
- D-WORLD-109: test stores, temp dirs and test keys only. No real key, no `~/.ssh`, no `~/.ailang/world`, no real
  `/dev/tty`, no real browser profile, no mint or publish, no real network listener beyond `httptest`/loopback in tests.
- `host/why` keeps effects out of its core: no `os`, `net/http`, file or stdin reads in the walk (design §3, VS2a);
  the CLI's `--result` file read, stdin read and `--expect-store` cwd logic stay at the CLI boundary. No new `/v1`
  route, no `surface_verb.go` change, no page change, no `.ail`.
- Do not touch: `tools/launchd/*`, any skill, `world/*.ail`, `packages/`, `scripts/` (except reading),
  `cmd/world-publish/*`, `~/.ailang/state/mission-v1*`, row 201's files, `host/surfacepage/*`, `host/daemon/surface_verb*.go`.
- Export `AILANG_BIN=$HOME/.pinned-ailang/ailang` in every shell that runs Go tests.
- Your Bash tool's shell is zsh: run loops/globs/`$VAR` work under `bash -c` or `bash <<'EOF'`. macOS has no
  `timeout` command. Pass `timeout: 600000` on any Bash call that may run long; the default is 120 s.
- Known load-sensitive / flaky tests (re-run once alone, do not chase): `TestSeedRestartIgnoresASurvivingOldCopy`
  (host/broker), `TestHostsLaunchersExecHost`, host/broker exec-probe timeouts incl.
  `TestExecProtectedCheckTimeout/reserve`, `TestPublishApprovalRefusalSetWithALandedPositiveControl`, local `-race` of
  host/broker, host/daemon, host/verifygate, `TestMCPPostBudgetProductionConstants/batch`, `TestBackoff`
  (host/surfacefollow), `TestDemoSeedPreflightRefusals`, `TestWorktreeSnapshotRefusesAConfiguredFilter`,
  `TestMissionConfigDefaultFlowIntegration` (host/verifygate), `scripts/marketingcheck` supervisor tests.
- **Time budget: aim to finish within ~90 minutes; hard stop at 130.** The iteration's 6 h wall clock started
  21:39Z, and the executor, gates, judge, CI and record still follow you. If you run short, take the design's §8 cut
  fully proven rather than everything half-proven.

## Waiting discipline (binding on you)
If you background anything (a long test run, a mutant loop), poll its ARTIFACT (an output/rc file it writes on
completion) inside a bounded `date +%s` loop in the same turn. A Monitor is an event stream, not a wait; nothing will
re-invoke you when it fires. Never end your turn announcing a wait, a pause, or "once it completes" — end only when
both files are written and the round trip has passed.

## Your final message
Report: the plan path, scope (what is in, anything moved to V2b and why, with measured numbers), the sub-milestone list
with expected numstat after each (and the rename-pairing proof), mutant counts (total / killed / named equivalents) and
the ledger path, the round-trip result, base-gate results, every premise found false (§1.4), how you resolved audit
inputs 1, 2, 3 and 7, the exact pre-existing test files changed, and anything the controller must decide.
