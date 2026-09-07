# w-gate1-unverified-non-head-commit — an instrument that names a zero-check commit HEAD has moved past

- **Status:** planned (design only; nothing implemented, no git write performed)
- **Date:** 2026-09-07 (revised 2026-09-07 — design-quorum revision pass; **second revision — controller SPLIT**)
- **Owning queue row:** 70, `w-gate-1-cannot-see-an-unverified-commit-that-is-not-head`
- **Base commit measured at:** `ab5af8c337963d86a34dae9743294cd3bb6c5713` (== `origin/dev`), worktree `.wt-world-iter170`, branch `mission/world-iter170`
- **Verify profile:** `ailang-code`; pinned binary `AILANG_BIN=/Users/voightkampff/.pinned-ailang/ailang` (v0.30.0, measured). Go host gated by `go vet ./...` / `go test -run '^$' ./<pkg>/` as the compile fence — **never** `go build ./...` (it does not compile `_test.go`).

---

## Why this is a split, not another revision

Round 2 blocked with BOTH reviewers present and BOTH rejecting, and both objections landed on the SAME
surface — the tip-attribution predicate — while every other part of the document (the range enumeration,
the per-commit check-set read, the loud-failure paths, the fixtures, the parser rule, the CI wiring)
drew no objection in round 2 at all. That is a scope signal, not a maturity signal: the document bundles
a cheap, uncontested instrument with a genuinely hard attribution question. The hard half is carved out
to a new queue row (see **"Carved out — queue row 88"** below). **This document is reduced to the
uncontested half.**

The reduced instrument reports every non-HEAD zero-check commit as a `ZERO-CHECK-CANDIDATE` and makes
**NO claim about push tips**. Reporting is not alarming: a candidate is a finding a controller reads,
not a gate failure. The `UNKNOWN`/LOUD-FAIL alternative is **not** applied — round 1's accepted objection
(gemini) blocked exactly that, and row 70's own scope note warns that an assertion which is not about
tips "manufactures alarms by construction".

---

## Motivation

Row 70's claim, verbatim: **a direct-to-dev commit with zero check-runs is invisible to every gate in
this loop, because Gate 1 reads only `origin/dev`'s HEAD.** Measured iter-150: `e308577` — the commit
that added a CI job and thereby armed that iteration's red — has `checks=0`. The origin skill has a
rule for a HEAD whose check count is a true zero (fire `workflow_dispatch`, do not record a health
verdict); it has **nothing** for a zero on a commit that HEAD has since moved past, and that is the
likelier shape here because World's dev takes direct fleet pushes between fires. The red surfaced two
commits later wearing a merge's clothes.

**THE ITEM:** at Gate 1, after the HEAD check-set read, enumerate `dev@{previous fire}..origin/dev` and
report every non-HEAD commit whose check set is empty as a `ZERO-CHECK-CANDIDATE` — cheap, one API read
per commit (the check-set read), and it is the only instrument that would have named `e308577` at the
time. The instrument does **NOT** claim a candidate was a push tip: a zero-check commit may be a
separately-pushed tip or an intra-push parent of a later merge, and local git topology cannot distinguish
them (sol's objection, quoted verbatim in the carved-out row below). A candidate is a finding a
controller reads, not a gate failure.

This document designs the World-ownable deliverable: an **instrument that lives in this repo** — a
script under `scripts/` with tests — that a controller can RUN at Gate 1. It does not edit the shared
skill (fleet-owned) and does not edit `tools/launchd/*` (frozen core).

---

## Premises (hard constraints)

1. **P1 — `tools/launchd/*` is FROZEN CORE.** World may never modify it (ratified `D-WORLD-DRIVER-1`).
   Any fix requiring a driver edit is a FLEET hand-off, not a World milestone.
2. **P2 — the shared `mission-control` SKILL.md is NOT in this repo.** It lives in the V1 checkout
   (`/Users/voightkampff/dev/sunholo-data/ailang/.claude/skills/mission-control/`). A change to Gate 1's
   prose is a PROPOSAL to the fleet, never a World edit.
3. **P3 — the World-ownable deliverable is an instrument in this repo** (a script under `scripts/`,
   with tests) that a controller can RUN. The fleet half is the Gate-1 prose change that calls it.
4. **P4 — the lower bound `dev@{previous fire}` is not a git ref this loop can rely on.** The repo
   records a per-fire base in `~/.ailang/state/mission-world-base` (a TSV, format measured below). The
   fleet script `tools/launchd/mission-base.sh` (V1 checkout) reads/writes it with
   `record`/`snap`/`last <label>`/`drift <label>`. The instrument uses `mission-base.sh last gate1` as
   the lower bound.
5. **P5 — the events API is NOT a dependency, and neither is push-tip attribution.** The events API was
   measured unreliable (see "Ruled out" below: it missed `e308577`'s push, is missing the last ~3h of
   pushes, is pagination-limited, and is not chronologically ordered). The instrument does **NOT** claim
   to know push tips: it reports every non-HEAD zero-check commit as a `ZERO-CHECK-CANDIDATE` (Decision
   3), which is topology-independent and has zero network dependency beyond the single check-set read per
   commit. The push-boundary question is carved out to queue row 88.
6. **P6 — `go build ./...` is NOT a compile fence.** It does not compile `_test.go`. Use
   `go vet ./...` or `go test -run '^$' ./<pkg>/`. No acceptance criterion is designed around `go build`.

---

## Verification Log

Each row is a claim about what the codebase or the rig currently does, with its COMMAND and OBSERVED
OUTPUT. Empty/negative results carry a known-positive control in the same row. Rows V16–V17 are the
controller's first-party measurements of the reduced instrument's value, cited `VERIFIED BY CONTROLLER`.

| # | Claim | COMMAND | OBSERVED OUTPUT |
|---|---|---|---|
| V1 | Base SHA == origin/dev | `git rev-parse HEAD; git rev-parse origin/dev` | both `ab5af8c337963d86a34dae9743294cd3bb6c5713` |
| V2 | `e308577` has `checks=0` (controller's measurement, re-run) | `gh api repos/sunholo-data/ailang-world/commits/e308577/check-runs --jq .total_count` | `0` |
| V3 | The check-runs endpoint answers; `0` is a true zero, not instrument failure | `for sha in ab5af8c33 cc2b5b3fd f69873e62 eff748911 9befb67ad 3305e2ed0 9d72cfe82 37a78abbb; do gh api .../commits/$sha/check-runs --jq .total_count; done` | `2` for every one of the 8 newest commits |
| V4 | No instrument in this repo performs this enumeration today | `grep -rn "check-runs" scripts/ host/ Makefile` | rc=2, zero hits. **Positive control:** `grep -rln "check-runs" .` → hits in `design_docs/` (e.g. `world-mission.md`, `world-mission-log.md`) |
| V5 | `actions/runs?head_sha=` answers for every commit in the current range (needs full 40-char SHA) | `for sha in ab5af8c337… cc2b5b3fd… f69873e62…; do gh api .../actions/runs?head_sha=$sha --jq .total_count; done` | `1` for each of the three commits in `eff748911..ab5af8c` |
| V6 | `e308577` has no run | `gh api .../actions/runs?head_sha=e30857731fd37f2d41e72158de6bbe21108adfd0 --jq .total_count` | `0` |
| V7 | `e308577` is in the range of a MERGE push, not a fast-forward push | `git log --oneline --format='%h %p %s' d75b9c3e9..68403ea` | `68403ea e308577 725ad5a Merge…`, `e308577 d75b9c3 fix(mission)…`, `725ad5a d75b9c3 row 56…` — so `e308577` is a parent of the merge `68403ea` |
| V8 | `68403ea` is a merge commit (2 parents); `725ad5a` is not (1 parent) | `git rev-list --parents -n1 68403ea; git rev-list --parents -n1 725ad5a` | `68403ea e308577 725ad5a`; `725ad5a d75b9c3` |
| V9 | The mission-base file format is `<epoch>\t<iso8601>\tbase-<label>\t<n>\t<sha>` | `head ~/.ailang/state/mission-world-base` | `1788820068\t2026-09-07T22:27:48Z\tbase-gate1\t1\tab5af8c337963d86a34dae9743294cd3bb6c5713` (and 11 more rows) |
| V10 | `mission-base.sh last gate1` returns the previous fire's Gate-1 base | `bash /Users/voightkampff/dev/sunholo-data/ailang/tools/launchd/mission-base.sh last gate1` | `ab5af8c337963d86a34dae9743294cd3bb6c5713` (== origin/dev, so the range is empty at base) |
| V11 | The local reflog is empty (not a usable tip source) | `git reflog -15` | only `ab5af8c HEAD@{0}: reset: moving to HEAD`. **Positive control:** `git log --oneline -5` → 5 commits |
| V12 | The pinned binary is v0.30.0 | `/Users/voightkampff/.pinned-ailang/ailang version` | `AILANG v0.30.0` |
| V13 | `gh` is authenticated for the API reads | `gh auth status` | `Logged in to github.com account sunholo-voight-kampff` |
| V14 | The check-runs fixture for a verified commit has `total_count:2` | `gh api .../commits/ab5af8c337…/check-runs` | `{"total_count":2,"check_runs":[{…"name":"ailang-code verify gate"…},{…"name":"go host build + test gate"…}]}` |
| V15 | The check-runs fixture for a zero-check commit has `total_count:0` | `gh api .../commits/e30857731…/check-runs` | `{"total_count":0,"check_runs":[]}` |
| V16 | Sweep of the last 40 commits of `origin/dev` (39 non-HEAD commits classified): the reduced instrument's two classes give `VERIFIED = 35`, `ZERO-CHECK = 4`; under the now-REMOVED topology predicate `UNVERIFIED-TIP = 0` — the predicate produced ZERO false alarms over 39 real commits | controller's hand sweep of `origin/dev` check sets | `VERIFIED = 35`, `ZERO-CHECK (no checks) = 4`, `UNVERIFIED-TIP = 0` — **VERIFIED BY CONTROLLER** |
| V17 | Positive control, range `165b9fd9d..68403eaf5`: the reduced instrument SURFACES `e308577` (checks=0) and marks the checked commits VERIFIED | `gh api .../commits/<sha>/check-runs --jq .total_count` for `e30857731`, `d75b9c3e9`, `aa610dc71`, `725ad5a29` | `e30857731 checks=0` (SURFACED — row 70's whole point); `d75b9c3e9 checks=2` VERIFIED; `aa610dc71 checks=2` VERIFIED; `725ad5a29 checks=2` VERIFIED (it is a merge parent too, so a check-set-first ordering is what keeps it out of the candidate set) — **VERIFIED BY CONTROLLER** |

### Ruled out: the events API (kept as a record of WHY it was rejected)

The events API was the original candidate tip source. It is **not** a dependency of this design, and the
push-boundary question it was meant to answer is carved out to queue row 88. Its measurements are
preserved here so the rejection is auditable and not re-litigated:

| # | Claim | COMMAND | OBSERVED OUTPUT |
|---|---|---|---|
| (V5) | The events API returns PushEvent `before`/`head` (push tips) | `gh api repos/sunholo-data/ailang-world/events?per_page=100 --jq '.[] \| select(.type=="PushEvent") \| .payload'` | `{"before":"8b600fd9c…","head":"ae9a61575…","ref":"refs/heads/dev",…}` etc. |
| (V6) | The events API is NOT chronologically ordered | `gh api .../events?per_page=100 --jq '.[] \| "\(.created_at) \(.type)"'` | `2026-09-05T20:29:33Z PushEvent` precedes `2026-09-07T19:39:14Z PushEvent` in the array; the most recent PushEvent head is `ae9a61575` (2026-09-05) while `origin/dev` is `ab5af8c` (2026-09-07) |
| (V7) | The events API is missing the last ~3h of pushes | search events pages 1–3 for `ab5af8c`/`cc2b5b3`/`f69873e` as PushEvent heads | only `f69873e62b077be435e08a962b43b4c77f54484d` found; `ab5af8c` and `cc2b5b3` NOT found. **Positive control:** `f69873e` IS found |
| (V8) | The events API missed `e308577`'s push (permanent, not a lag — 5 days elapsed) | `gh api .../events?per_page=100 --jq '[.[] \| select(.type=="PushEvent") \| .payload.head] \| map(select(startswith("e308577")))'` | `[]` (not found). **Positive control:** `f69873e` found in the same call |
| (V9) | The events API is pagination-limited | `gh api .../events?per_page=100&page=4` | HTTP 422 `"pagination is limited for this resource"` |
| (V21) | The events fixture contains 25 PushEvents | `gh api .../events?per_page=100 --jq '[.[] \| select(.type=="PushEvent")] \| length'` | `25` |

**Why rejected:** a zero-check commit absent from the events API would fall through to `UNCLASSIFIED`
and trigger a loud-fail, so every normal multi-commit push would deterministically crash Gate 1 and park
the loop for hours until the upstream API caught up (V7, V8). The reduced instrument does not need a tip
source at all: it reports `ZERO-CHECK-CANDIDATE`s from the check-set read alone, which cannot lag or lose
events.

---

## Decisions

### Decision 1 — The instrument is a bash script under `scripts/`, runnable by a controller

`scripts/gate1_range_check.sh`. It orchestrates `git` (range enumeration) and `gh` (the check-set read).
This matches the repo's existing gate scripts (`verify_ail.sh`, `verify_go.sh`) and is directly runnable
at Gate 1. The classification logic is factored into pure, testable functions so the test script can
drive them against checked-in fixtures and a throwaway repo.

**Why bash, not Go:** the instrument is a Gate-1 orchestration script (git + API reads + exit codes), the
same shape as the repo's other gate scripts. The pure classification is small and testable in bash via
the existing `test_*.sh` pattern (`test_mission_answer.sh`). A Go host binary would add a build step for
no benefit here; the mission's Go host is for the kernel/broker, not for a one-shot gate probe.

### Decision 2 — The lower bound is `mission-base.sh last gate1`

The range is `$(mission-base.sh last gate1)..origin/dev`. `last gate1` returns the last SHA recorded
under `base-gate1` in `~/.ailang/state/mission-world-base` (V9, V10) — the previous fire's Gate-1 base,
which is exactly what the row means by `dev@{previous fire}`. The instrument reads the file directly
(fall back to `mission-base.sh last gate1` if the fleet script is reachable) and extracts the last
`base-gate1` row's 5th field. If no `base-gate1` row exists, the instrument LOUD-FAILs (no lower bound).

**Why `base-gate1`, not `base-gate4`:** the row's item runs at Gate 1, after the HEAD check-set read. The
previous fire's Gate-1 base is the correct lower bound for "commits since the last time Gate 1 ran".
`base-gate4` is a later stamp in the same fire and would narrow the range incorrectly.

### Decision 3 — The classification is two classes, topology-independent

Per-commit classification, driven **only** by the check-set read (one `gh api` read per commit):

| Class | Condition | Gate effect |
|---|---|---|
| `VERIFIED` | check set `total_count > 0` | OK |
| `ZERO-CHECK-CANDIDATE` | check set `total_count == 0` | exit 1 — **a controller must look** |

**Two classes. No third class. No claim about push tips.** A `ZERO-CHECK-CANDIDATE` is a finding a
controller reads, not a gate failure: it may be a separately-pushed tip or an intra-push parent of a
later merge, and the instrument does not (and cannot, from local topology) distinguish them. The
push-boundary question is carved out to queue row 88.

**Ordering (MUST-ANSWER 1) — the check-set read is the ONLY classification input.** There is no topology
predicate to order against. A commit with a check set (`total_count > 0`) is `VERIFIED` regardless of
anything else. This is what keeps a checked commit like `725ad5a` (a merge parent, V17) out of the
candidate set: it has a check set, so it is `VERIFIED`.

### Decision 4 — The verdict, with loud-failure paths

**Head exclusion (MUST-ANSWER 2) — topology-safe, and NOT `head^`.** `origin/dev`'s HEAD is already
covered by Gate 1's existing HEAD check-set read and by the shared skill's true-zero rule, so the
instrument must not double-report it. The enumeration rule is the reviewers' (round-3 quorum,
`gpt5-6-sol` and `gemini-3-1-pro` independently, applied verbatim under the narrow-refinement
carve-out):

> “Enumerate `git rev-list "$base..$head"`, then filter out only the commit whose SHA exactly equals
> the resolved `$head`; do not use `$head^`.” — `gpt5-6-sol`

> “Do not use Git parent shorthand (`^`) to exclude the head. Enumerate the full range
> `base..origin/dev`, and explicitly exclude the head inside the per-commit loop by comparing SHAs
> (e.g., `if [ "$commit" = "$head_sha" ]; then continue; fi`).” — `gemini-3-1-pro`

**Why `head^` is wrong, MEASURED BY THE CONTROLLER at `ab5af8c` (rule 3f — the objection was
reproduced, not forwarded).** When `$head` is a merge, `$head^` resolves to its FIRST parent only, so
every commit reachable exclusively through the other parent is silently dropped from the enumeration.
On the real range `165b9fd9d..68403eaf5` (whose head `68403ea` is a merge of `e308577` and `725ad5a`):

| Expression | Commits enumerated |
|---|---|
| `git rev-list $base..$head` | `68403eaf5 e30857731 725ad5a29 d75b9c3e9 aa610dc71` (5) |
| `git rev-list $base..$head^` — **the defect** | `e30857731 d75b9c3e9 aa610dc71` (3) — **`725ad5a29` silently dropped** |
| full range minus the exact head SHA — **the fix** | `e30857731 725ad5a29 d75b9c3e9 aa610dc71` (4) |

**One honest correction to the reviewers' framing, recorded rather than smoothed over:** both
reviewers wrote that `head^` “can silently miss the zero-check commit the instrument exists to
surface.” The general claim is correct and the defect is real, but in THIS measured instance
`e308577` is the merge's *first* parent, so it survives `head^`; the commit actually dropped is
`725ad5a29`. Had the merge parents been ordered the other way, `e308577` itself would have been
dropped. The fix is required either way — the enumeration must not depend on parent ordering.

The excluded commit is exactly the resolved `$head` SHA; every other commit in `base..head`,
through every parent lineage, is classified.

**Every `gh api` call is bounded with an explicit timeout, mapped to exit 2.** A hung or slow API call
must not leave the instrument waiting indefinitely; a timeout is a LOUD-FAIL, not a silent green.

Gate verdict and exit codes:

| Exit | Verdict | When |
|---|---|---|
| `0` | GREEN | every non-HEAD commit is `VERIFIED` |
| `1` | CANDIDATE | one or more `ZERO-CHECK-CANDIDATE`s found — **"a controller must look"**, NOT "a commit is unverified" |
| `2` | LOUD-FAIL | malformed/missing `total_count`, unreachable API, `gh` timeout, or missing lower bound |
| `3` | EMPTY | the range `last gate1..origin/dev` is empty (no new commits since the previous fire) |

The instrument **never** silently classifies a zero-check commit as anything other than a
`ZERO-CHECK-CANDIDATE`. The `EMPTY` case is explicit and observable (a distinct exit code + message),
never a silent green. The unreachable-API, `gh`-timeout, and zero-cardinality-enumeration paths still
fail loudly and visibly.

**Parser rule (sol's round-1 rule, applied exactly):** valid zero cardinality
(`{"total_count":0,"check_runs":[]}`) is **ACCEPTED** — it is a true zero, not a parse failure. Only a
missing, non-numeric, negative, or inconsistent `total_count` LOUD-FAILs (exit 2).

---

## Milestones

Total ~0.15 days. Two independently committable, CI-green milestones. Neither edits frozen core.

- **M1 — the enumeration + check-set read + classification + loud-fail paths + fixtures (~0.1 days).**
  `scripts/gate1_range_check.sh` resolves the range `last gate1..origin/dev`, excludes the head, reads
  each remaining commit's check set (one `gh api` read per commit, each with an explicit timeout mapped
  to exit 2), and classifies `VERIFIED` vs `ZERO-CHECK-CANDIDATE`. Verdict: GREEN (exit 0) if all
  non-HEAD commits are `VERIFIED`, else exit 1 (candidates found). Handles the `EMPTY` (exit 3),
  unreachable-API (exit 2), `gh`-timeout (exit 2), and malformed-`total_count` (exit 2) loud-failure
  paths. Ships fixtures + `scripts/test_gate1_range_check.sh`. **This is the row's "one API read per
  commit" core and is useful on its own: it surfaces zero-check commits that Gate 1 currently cannot
  see.**
- **M2 — CI wiring + the report format (~0.05 days).** Finalizes the per-commit report line (each line
  names the commit and its class; the exit-1 message states plainly that a candidate means "a controller
  must look", not "a commit is unverified"). Wires `scripts/test_gate1_range_check.sh` into
  `.github/workflows/ci.yml` (a new step in the `go-verify` job, or a new job) so the instrument's tests
  are CI-verified.

---

## Files to Create/Modify

- `scripts/gate1_range_check.sh` — **WORLD** — the instrument (M1 core, M2 report format).
- `scripts/test_gate1_range_check.sh` — **WORLD** — the test script (M1, extended in M2).
- `scripts/testdata/gate1_check_runs_verified.json` — **WORLD** — fixture, generated by
  `gh api repos/sunholo-data/ailang-world/commits/ab5af8c337963d86a34dae9743294cd3bb6c5713/check-runs`
  (V14). Checked in as produced.
- `scripts/testdata/gate1_check_runs_zero.json` — **WORLD** — fixture, generated by
  `gh api repos/sunholo-data/ailang-world/commits/e30857731fd37f2d41e72158de6bbe21108adfd0/check-runs`
  (V15). Checked in as produced.
- `scripts/testdata/gate1_check_runs_malformed.json` — **WORLD** — fixture, generated by
  `printf '{"check_runs":[]}' > scripts/testdata/gate1_check_runs_malformed.json`. Checked in as produced.
- `.github/workflows/ci.yml` — **WORLD** — add a step running `scripts/test_gate1_range_check.sh` (M2).
  This is the World's own CI workflow, not frozen core.
- `design_docs/planned/w-gate1-unverified-non-head-commit.md` — **WORLD** — this document.

**Fixture provenance (S7 / the parser rule):** every fixture is generated by the named command and
checked in byte-for-byte as produced — never hand-typed, never transcribed from this doc. The two real
check-runs fixtures come from genuinely different real invocations (a verified commit, V14; a zero-check
commit, V15). The malformed fixture is synthetic (a negative parser test) and is produced by the named
`printf` command, not transcribed from a real response. The parser must extract exactly one non-negative
integer `total_count`; it must accept `{"total_count":0,"check_runs":[]}` as a valid zero-check response
and LOUD-FAIL only when `total_count` is missing, non-numeric, negative, or inconsistent with the
`check_runs` array when that array is present (AC-M1-7, AC-M1-8).

---

## Acceptance Criteria

Each AC is a command with its expected output, **red at base** (base `ab5af8c337963d86a34dae9743294cd3bb6c5713`).
ACs marked M1 are green after M1; ACs marked M2 are green after M2. The test script supports `--only
<arm>` so each AC targets one behavior.

**M1 — enumeration + check-set read + classification + loud-fail paths:**

- **AC-M1-1** — the instrument exists and is executable.
  `test -x scripts/gate1_range_check.sh && echo EXISTS` → `EXISTS`. (Red at base: file absent.)

- **AC-M1-2** — the instrument resolves a range and reports per-commit check counts, GREEN. The
  `check-count` arm builds a throwaway linear repo and asserts each commit is reported `:2` with a GREEN
  verdict:
  `scripts/test_gate1_range_check.sh --only check-count` → `ok check-count` and exit 0. (Red at base: no
  instrument.)

- **AC-M1-3** — a zero-check commit is reported `ZERO-CHECK-CANDIDATE` → exit 1. The `candidate` arm
  builds a throwaway repo with one zero-check commit and asserts it is named `ZERO-CHECK-CANDIDATE` with
  exit 1:
  `scripts/test_gate1_range_check.sh --only candidate` → `ok candidate` and exit 0. (Red at base: no
  instrument.)

- **AC-M1-4** — the instrument LOUD-FAILs on an empty range (distinct exit 3, explicit message).
  `scripts/gate1_range_check.sh --base ab5af8c337963d86a34dae9743294cd3bb6c5713 --head ab5af8c337963d86a34dae9743294cd3bb6c5713; echo "rc=$?"` →
  `rc=3` and a message `no new commits since <base>; nothing to check`. (Red at base: no instrument.)

- **AC-M1-5** — the instrument LOUD-FAILs on an unreachable API (distinct exit 2).
  `scripts/gate1_range_check.sh --base eff748911… --head ab5af8c… --gh-bin /nonexistent/gh; echo "rc=$?"` →
  `rc=2` and a message naming the unreachable API. (Red at base: no instrument.)

- **AC-M1-6** — the instrument LOUD-FAILs on a `gh` timeout (distinct exit 2). Every `gh api` call is
  bounded with an explicit timeout; a timeout maps to exit 2, not a silent green.
  `scripts/gate1_range_check.sh --base eff748911… --head ab5af8c… --gh-timeout 0.001; echo "rc=$?"` →
  `rc=2` and a message naming the timed-out API. (Red at base: no instrument.)

- **AC-M1-7** — the parser LOUD-FAILs on a malformed check-runs response (missing `total_count`).
  `printf '{"check_runs":[]}' > /tmp/malformed.json; scripts/gate1_range_check.sh --check-runs-file /tmp/malformed.json; echo "rc=$?"` →
  `rc=2` and a message `missing or invalid total_count`. (Red at base: no instrument.)

- **AC-M1-8** — the parser ACCEPTS a valid zero-check response and proceeds to zero-check classification.
  `scripts/gate1_range_check.sh --check-runs-file scripts/testdata/gate1_check_runs_zero.json; echo "rc=$?"` →
  `rc=0` and `total_count=0` (valid zero cardinality is NOT a parse failure). (Red at base: no
  instrument.)

- **AC-M1-9** — the range's head is excluded (not double-reported). The `head-excluded` arm builds a
  throwaway repo with a **LINEAR history of 2+ commits** (a, then the head) and asserts the head is NOT
  reported as `ZERO-CHECK-CANDIDATE`, that the non-head commit IS classified, and that the verdict is
  GREEN: `scripts/test_gate1_range_check.sh --only head-excluded` → `ok head-excluded` and exit 0. (Red
  at base: no instrument.)
  **⚠ THE 2+-COMMIT RANGE IS LOAD-BEARING, NOT INCIDENTAL** (corrected at iteration 170 by the
  independent evaluator, and reproduced first-party by the controller before the fix). This arm was
  first written with a range of exactly ONE commit — base+1 == head. That is degenerate: under
  AC-M1-10's mutation `head^` resolves to `base` itself, so the range collapses to empty and the arm
  exits **3** ("no new commits"), i.e. it reds for a reason that has nothing to do with head exclusion.
  Measured both ways: with the one-commit range, AC-M1-10's mutation redded *both* arms
  (`merge-head-coverage` rc=1 **and** `head-excluded` rc=3), so the independence claimed below was
  **false as shipped**; with the widened range it reds `merge-head-coverage` (rc=1) while
  `head-excluded` stays GREEN (rc=0). The general claim was always true for 2+-commit linear ranges —
  the 3-commit `check-count` arm stayed green under the same mutation throughout — so the defect was in
  this arm's construction, not in the design. The third assertion (`classifies the non-head commit`) is
  what makes the collapse visible rather than silent.

- **AC-M1-10** — a MERGE head does not hide its second parent's lineage (round-3 quorum, `gpt5-6-sol`'s
  acceptance arm applied verbatim). The `merge-head-coverage` arm “creates a merge head with unique
  commits on both parents, assigns zero checks to a second-parent commit, and asserts that both parent
  histories are examined, the exact merge head is excluded, and the second-parent zero-check commit is
  reported as `ZERO-CHECK-CANDIDATE` with exit 1”:
  `scripts/test_gate1_range_check.sh --only merge-head-coverage` → `ok merge-head-coverage` and exit 0.
  (Red at base: no instrument.) **This arm is the sole killer of the `head^` defect: AC-M1-9's history is
  linear, so it cannot see the coverage hole — which is exactly the gap the round-3 quorum named.**

**M2 — CI wiring + the report format:**

- **AC-M2-1** — the instrument's test script is wired into CI and passes.
  `grep -n 'test_gate1_range_check' .github/workflows/ci.yml` → a step invoking it. (Red at base: no such
  step.)

- **AC-M2-2** — the report format names each commit's class and the exit-1 message states plainly that a
  candidate means "a controller must look", NOT "a commit is unverified". The `report` arm builds a
  throwaway repo with one zero-check commit and asserts the per-commit line names `ZERO-CHECK-CANDIDATE`
  and the exit-1 message reads `a controller must look: N zero-check candidate(s)`:
  `scripts/test_gate1_range_check.sh --only report` → `ok report` and exit 0. (Red at base: no
  instrument.)

---

## Non-Vacuity — Named RED Mutation for Every Gate

Each load-bearing assertion must have a named mutation that makes it red, and the single assertion that
fires.

- **AC-M1-2** — mutation: delete the check-set read (always report `0`). The per-commit line reads `:0`
  instead of `:2` → the `:2` assertion reds.
- **AC-M1-3** — mutation: classify every zero-check commit as `VERIFIED` (drop the candidate branch).
  `rc=0` instead of `rc=1` → the exit-1 assertion reds. **Catches the "candidate never surfaces" vacuity.**
- **AC-M1-4** — mutation: drop the empty-range branch (fall through to GREEN). `rc=0` instead of `rc=3` →
  the `rc=3` assertion reds. **Catches the zero-cardinality vacuity.**
- **AC-M1-5** — mutation: swallow the `gh` failure (exit 0 on unreachable API). `rc=0` instead of `rc=2` →
  the `rc=2` assertion reds. **Catches the unreachable-API vacuity.**
- **AC-M1-6** — mutation: drop the explicit timeout (no `--timeout` on the `gh api` call). A hung `gh`
  call is not mapped to exit 2 → the timeout assertion reds. **Catches the unbounded-API-call vacuity.**
- **AC-M1-7** — mutation: let the parser return 0 records silently (no loud failure). `rc=0` instead of
  `rc=2` → the `rc=2` assertion reds. **Catches the parser-loud-failure requirement.**
- **AC-M1-8** — mutation: treat valid zero cardinality as a parse failure. The zero fixture is rejected →
  the `total_count=0` assertion reds. **Catches the "empty array is a parse failure" self-contradiction.**
- **AC-M1-9** — mutation: include the range's head in the classification loop. The head is reported
  `ZERO-CHECK-CANDIDATE` → false exit 1 → the `ok head-excluded` assertion reds. **Catches the
  double-report / false alarm on the head.**
- **AC-M1-10** — mutation: replace the SHA-equality head filter with `git rev-list "$base..$head^"`.
  The second-parent zero-check commit vanishes from the enumeration, so it is never reported and the
  run exits 0 instead of 1 → the `ok merge-head-coverage` assertion reds. **Catches the merge-head
  coverage hole. AC-M1-9 stays GREEN under this same mutation — MEASURED at iteration 170, and only
  after AC-M1-9's arm was widened to a 2+-commit linear range; see the correction recorded under
  AC-M1-9 above. That green is what makes AC-M1-10 non-redundant rather than a second arm on the same
  property, and it is an assertion about this suite that must be re-measured, never assumed.**
- **AC-M2-1** — mutation: remove the CI step. `grep` finds nothing → the step assertion reds.
- **AC-M2-2** — mutation: the exit-1 message claims "commit is unverified" instead of "a controller must
  look". The message assertion reds. **Catches the "instrument adjudicates" overclaim.**

---

## Conflict Surface

- **`scripts/`** — the new `gate1_range_check.sh` and `test_gate1_range_check.sh` are new files. Existing
  scripts (`verify_ail.sh`, `verify_go.sh`, `bench_worldd.sh`, `mission_answer.sh`, `mission_decisions.sh`,
  `build_world_package.sh`, `verify_world_package.sh`) are untouched. No name collision.
- **`scripts/testdata/`** — the new check-runs fixtures are new files. Existing fixtures
  (`ailang_release_observed.txt`, `ailang_version_shim.sh`, `upstream_release_tags.txt`) are untouched.
  The events fixture is NOT created (the events API is ruled out).
- **`.github/workflows/ci.yml`** — M2 adds a step. This is the World's own CI workflow (not frozen core).
  The step must not break the existing `ailang-verify` / `go-verify` jobs; it is additive. The `go-verify`
  job already installs the pinned v0.30.0 binary and Z3; the new step runs a bash test script and needs no
  new tooling.
- **`~/.ailang/state/mission-world-base`** — the instrument READS this file (via `mission-base.sh last
  gate1` or directly). It never writes it. The fleet script owns writes. No write conflict.
- **The shared `mission-control` skill** — NOT touched by this design (P2). The instrument is a standalone
  script; the fleet's Gate-1 prose change (Fleet hand-off) is a separate proposal.
- **`tools/launchd/*`** — NOT touched (P1). The instrument does not call or modify the driver.

---

## Fleet hand-off

- **FLEET's half (a PROPOSAL, never a World edit):** a Gate-1 prose change in the shared
  `mission-control` skill (V1 checkout, `resources/gate-1-observe.md`) to call
  `scripts/gate1_range_check.sh` after the HEAD check-set read, and to route its verdicts (exit 1 = a
  controller must look; LOUD-FAIL parks for investigation; EMPTY is a no-op). This is a proposal to the
  fleet because the skill is not in this repo (P2). World does not edit it.
- **WORLD's half (complete and useful on its own):** the instrument `scripts/gate1_range_check.sh` +
  fixtures + tests. A controller can run it manually at Gate 1 today, without any fleet change. It is not
  gated on the fleet: the fleet proposal only wires it into the skill's prose; the instrument's value
  (naming a zero-check non-head commit) is available immediately.

---

## Carved out — queue row 88: `w-gate1-push-boundary-attribution`

The hard half of this row's original scope — attributing a zero-check commit to a push boundary — is
carved out to a new queue row. This document is reduced to the uncontested half (reporting
`ZERO-CHECK-CANDIDATE`s). The carved-out row inherits the following.

**The question:** how to distinguish a separately-pushed zero-check tip from an intra-push commit that
shares a push with a later merge. Local git topology cannot answer it. sol's round-2 objection, quoted
verbatim, is the problem statement:

> **Strongest objection:** The core topology predicate cannot determine push tips. A zero-check commit
> whose immediate child is a merge may be either a separately pushed tip or an intra-push parent included
> in the same push as that merge; both have identical local Git topology. V16 proves only that `e308577`
> is a merge parent, not that it was a push tip. The design therefore labels ambiguous commits
> `UNVERIFIED-TIP` while explicitly admitting this can manufacture false REDs, contradicting its claimed
> classification semantics.
>
> **Catch:** Verify the actual push boundary for `e308577` and provide evidence that the proposed
> predicate distinguishes separately pushed tips from same-push merge parents. Determinism of local
> topology does not establish correctness, and the rejected Events API does not supply the missing
> provenance.
>
> **Proposed fix:** Replace `EXPLAINED`/`UNVERIFIED-TIP` attribution with an honest topology-independent
> result such as `ZERO-CHECK-CANDIDATE`, reporting every non-HEAD commit with `total_count == 0` without
> claiming it was a push tip. If Gate 1 requires definitive tip classification, return
> `UNKNOWN`/LOUD-FAIL unless an authoritative, persisted push-boundary source is available. Add fixtures
> for two histories with identical Git topology but different push boundaries and require the design not
> to claim it can distinguish them. Also bound every `gh api` call with an explicit timeout and map timeout
> to exit 2.

**What was ruled out and why:** the GitHub events API (measured unreliable — not chronologically ordered,
~3h lag, permanent loss of `e308577`'s push, pagination-limited; see the "Ruled out" table above) and
local topology alone (sol's objection above).

**gemini-3-1-pro's round-2 fix, verbatim, which the new row inherits:**
> 1. In Decision 3 and V16-V21, replace `git rev-list --children --all` with
>    `git rev-list --children $base..origin/dev` to strictly bound the topology to the evaluated range.
> 2. Update Decision 3 to explicitly define multi-child handling: 'If a commit has multiple children in
>    the range, it is an UNVERIFIED-TIP if ANY child in the range is a merge'.

**Corroborating measurement (VERIFIED BY CONTROLLER):** while running the topology sweep by hand the
controller's own ad-hoc probe passed two SHAs to `git rev-list --parents -n1` and got
`fatal: ambiguous argument` — i.e. the multi-child case gemini names occurred immediately, in the first
hand-rolled implementation of the predicate. That is first-party evidence that gemini's clause 2 is
required, not hypothetical.

**What an answer would need:** an authoritative, persisted push-boundary source. Note that recording the
push boundary locally at Gate 1 time (each fire stamping the observed `origin/dev` tip) would accumulate
exactly that record going forward, and cannot recover history.

**Measurements handed off (VERIFIED BY CONTROLLER, taken at `ab5af8c`):**
```
Sweep of the last 40 commits of origin/dev (39 non-HEAD commits classified):
  VERIFIED = 35, ZERO-CHECK (no checks) = 4, and under the now-REMOVED topology predicate
  UNVERIFIED-TIP = 0  -> the predicate produced ZERO false alarms over 39 real commits.

Positive control, range 165b9fd9d..68403eaf5:
  e30857731 checks=0 -> the instrument SURFACES it   (this is row 70's whole point)
  d75b9c3e9 checks=2 -> VERIFIED
  aa610dc71 checks=2 -> VERIFIED
  725ad5a29 checks=2 -> VERIFIED  (it is a merge parent too, so a check-set-first ordering is
                                   what keeps it out of the candidate set)
```

**Explicitly:** the reduced instrument does NOT depend on this row. It is useful today.

---

## Risks

- **The instrument reports candidates, not verdicts on push tips.** A controller must read the report.
  This is by design: exit 1 means "a controller must look", not "a commit is unverified". The instrument
  surfaces; it does not adjudicate.
- **The range is empty at base (V10)** — `last gate1 == origin/dev`. The instrument's `EMPTY` path (exit
  3) handles this explicitly; it is not a silent green.
- **The head is excluded (MUST-ANSWER 2).** If Gate 1's existing HEAD read is ever removed, the instrument
  would silently stop covering the head. This is a fleet-prose concern (the HEAD read is in the shared
  skill), not a World-instrument defect; the instrument documents the exclusion so the coupling is visible.

---

## Deferred Scope

- **Push-boundary attribution** — carved out to queue row 88 (`w-gate1-push-boundary-attribution`). The
  reduced instrument does not depend on it.
- **A Go host implementation** of the classification logic. The bash instrument is sufficient and matches
  the repo's gate-script pattern; a Go port adds a build step for no benefit now.
- **Wiring the instrument into the fleet skill's Gate-1 prose** — that is the fleet's half (Fleet hand-off),
  deferred to the fleet's cadence.
- **A `workflow_dispatch` re-fire for a `ZERO-CHECK-CANDIDATE`** — the skill's existing HEAD-zero rule
  fires a dispatch; extending it to non-head candidates is a fleet decision, not a World instrument
  change.

---

## Quorum verification log

| Round | Reviewers present | Verdict | Surface each objection landed on | Disposition |
|---|---|---|---|---|
| 1 | `gemini-3-1-pro` (reject). ABSENT: `gpt6-astra`, `oc-glm-5-2` — both `unknown-model` on the pinned v0.30.0 binary (a CAPABILITY limit of the pinned binary, not budget). The OpenAI seat was RESTORED as `gpt5-6-sol` via `design-review --max-cost-usd 0.30` ($0.0481) and returned **reject**. | blocked | gemini: the events-API dependency (tip source). sol: AC-M2-4's parser self-contradiction. | Two surfaces → REVISE. Both fixes applied verbatim. |
| 2 | `gpt5-6-sol` (reject), `gemini-3-1-pro` (reject). No absentees. | blocked | **Both on ONE surface: the tip-attribution predicate.** Everything else drew none. | **CONTROLLER SPLIT** (Gate-2 localisation rule): the topology predicate REMOVED and carved out to queue row 88 with sol's objection and gemini's fix verbatim; sol's `ZERO-CHECK-CANDIDATE` fix applied verbatim to the retained half. |
| 3 | `gpt5-6-sol` (reject), `gemini-3-1-pro` (reject). No absentees. | blocked | Both on ONE narrow defect, and both proposed the SAME fix: `base..head^` drops a merge head's second-parent lineage. | **NARROW-REFINEMENT CARVE-OUT**: both objections carry a concrete reviewer-authored `proposed_fix` and neither disputes the design DIRECTION (completeness). Both fixes applied VERBATIM by the controller, and the defect was REPRODUCED first-party (rule 3f) before applying — see the Head-exclusion table above, including the correction that the commit actually dropped is `725ad5a29`, not `e308577`. New AC-M1-10 + its named mutation added. Routed to sprint-planner. |

**Metered quorum spend:** r1 $0.020758 (gemini) + restore $0.048105 (sol) + r2 $0.080378 + r3 $0.083044 = **$0.232285**.
