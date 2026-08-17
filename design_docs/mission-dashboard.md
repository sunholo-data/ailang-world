# AILANG World — mission dashboard

*Snapshot, overwritten every iteration. History: charter STATUS + `world-mission-log.md`.*
**As of** 2026-08-17 (iter-89) · **dev** `e73b10d` · **CI** green, SHA-addressed, `checks=2` both
`success`. One gate fix landed; nothing routed. `metered=$0.00` of $5.

## 🔧 LANDED — the local verify gate was 100% dead on a green tree

`verify_ail.sh` refused `[NOT_A_RELEASE]` and **all 17** `host/verifygate` arms failed, with the
binary, modules and contracts all fine. `~/.ailang/state` crossed **322 MB** (`observatory.db`
272 MB) past the binary's 200 MB threshold, so every `ailang` call now prefixes an
`Observatory: 269MB` warning **on stderr** — and World merged stderr into output it *parses*, at
three sites (`verify_ail.sh:70` `2>&1|head -1`; `run_bounded`'s `stderr=subprocess.STDOUT` feeding a
JSON parser; `ail_binary_gate_test.go:46` `CombinedOutput()`). Both `run_bounded` callers parse JSON,
so the merge was wrong for every caller.

Not ours: a pristine `origin/dev` worktree failed **17/17** identically, and the same tree is green
on CI. After the fix: `verify_ail.sh` rc=0 (**10** identities / **39** tests / **9/9**, matching the
recorded pins), `go test ./...` rc=0 **0 FAIL**. Non-vacuity: still refuses `[DEV_BUILD]`. Mutation:
re-merging stderr reproduces the byte-identical red; restore byte-identical.
**An instrument that captures more than it parses can be voided by a process it has nothing to do
with** — the tell is that every arm dies at once and names the environment, not the assertion.

## ⚠ FOUR open asks — the queue is IDLE until one is answered

The headless queue is fully blocked, re-verified this iteration by reading all **24** rows' LIVE
heads (not inherited, and not through `~~` spans — iter-87's trap).

- **D-WORLD-5** — item 5 `w-mcp-projection`. **A** = import upstream `serveapi` at v0.33.1 ·
  **B** = stay dependency-free, build the seam locally. Unblocks the most (items 6 and 7 sit behind
  item 5 *landing*, which its discharged prerequisites are not).
- **D-WORLD-17** — one word (evidence-boundary architecture, re-parked after quorum round 4).
- **D-WORLD-18** — one word (daemon read cancellation, scope A/B).
- **D-WORLD-DRIVER-1 (NEW)** — see below.

## The spine — the driver this loop runs is not in git

`dev.ailang.mission-world.plist:14` executes this repo's **working-tree**
`tools/launchd/mission-control.sh`, modified and uncommitted since 2026-08-15, **109 diff lines**
from origin — plus untracked `derive-planner-lane.sh`, `test_mission_routing.sh`, `testdata/`, and a
`verify_go.sh` step that runs them. Gate 1's `cmp` guards `SKILL.md` and came back clean while the
**driver** diverged. The same bundle carried the **Human Decision Ledger** — Mark's 2026-08-15
"authoritative current state" — so `origin/dev` held **0** ledgers against **1** in the working tree.

**Fixed for the half World owns**: ledger + `scripts/mission_decisions.sh` committed
(`--check` → valid, 5 rows). `tools/launchd/*` is **frozen core** per CLAUDE.md, so the controller
may not land it → `D-WORLD-DRIVER-1`. Whole bundle backed up to
`~/.ailang/state/world-driver-backup-2026-08-17/`, 6/6 sha256 OK.

## Second — the two-day silence was a REFUSAL, not a death

Three slots (08-16 20:28, 08-17 00:29, 04:29) logged `NO usable controller … Refusing.`: every
Anthropic pref quota-limited **and** codex returning `try again at Aug 20th, 2026 5:34 AM`. That is
the new driver's fail-closed path working — zero tokens beyond probes, loud log — the opposite of
Standing rule 7's silent `rc=0` orphan. **Codex is exhausted fleet-wide until 2026-08-20 05:34**, so
the designer rotation's `codex:` entry will probe-fail on the next new-doc iteration.

Third: `derive-planner-lane.sh` is recorded as **absent** (iter-80) and now **exists** — returns
`opus fail-closed:env-pin`. A recorded *absence* goes stale too, and nothing re-reads one.

## Loop

launchd, ~6h watchdog. Controller `claude:claude-opus-5`; **no designer/planner/executor/evaluator
spawned**. Rotation pointer untouched at `codex:gpt-5.6-sol`. Weekly sweep: **0 orphans of 1**
enumerated open issue. Bookkeeping thread rotated off `#53` (Monday-07:00-local rule).
