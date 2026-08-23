# Sprint plan — `WB.A`–`WB.K` (queue item 14, `w-workbench-read-only`, clause-2)

**Status**: IN SPRINT — **`WB.B` LANDED 2026-08-23 (iteration 112)**, 2 of 11 milestones; `WB.C` is next · **Design doc**: [`w-workbench-read-only.md`](w-workbench-read-only.md) (AUTHORITATIVE — the doc wins on every divergence below) · **Planner**: mission-control iteration 110, opus lane (`derive-planner-lane.sh` → `opus fail-closed:env-pin`, token used VERBATIM).

- **Plan (machine)**: `.ailang/state/sprints/w-workbench-read-only.plan.json` — `jq -e .` passes.
- **Design doc (AUTHORITATIVE)**: `design_docs/planned/w-workbench-read-only.md`, 894 lines.
- **Base**: `3e0c34cd5b7373cbc028838640c8301cba8c7fa0` (`dev` == `origin/dev`, clean).
- **Planner**: sprint-planner (opus, model-pinned, lane token `fail-closed:env-pin`), iteration 110.
- **Platform for every measurement below**: darwin/arm64, go1.25.6, AILANG v0.30.0 pinned at
  `/tmp/ailang-v0300/ailang`. **The windows and ubuntu CI legs are UNRUN LOCALLY.**
- **THIS TRACKED FILE IS THE EXECUTOR'S COPY.** `.ailang/` is gitignored (`.gitignore:3`,
  `**/.ailang/`), so the machine plan does **not** appear in a fresh `git worktree add` checkout —
  which is where the sprint runs. All 27 prior plans in this repo share that property, which is why
  the `-sprint-plan.md` companion exists and travels with the design doc.

---

## 0a. Controller's first-party reproduction of the planner's load-bearing findings

A finding handed up by a sub-agent is a **claim**, whatever its provenance. Every finding below was
re-run by the controller at `3e0c34c` on darwin/arm64 before it was recorded anywhere. Three of the
planner's four vacuity findings were reproduced; **all three are real**, and the mechanism behind
two of them is sharper than the planner stated.

| Finding | Controller's command | Observed | Verdict |
|---|---|---|---|
| **D1** — AC7's `git diff` is blind to the files this sprint adds | wrote a real `host/daemon/zz_ctrl_probe_iter110.go`, then ran the doc-form `git diff --name-only 93e1ba5 -- ':!design_docs'` | **0** matches, **0** total lines | **REPRODUCED** |
| D1 positive control — the file genuinely exists | `ls host/daemon/zz_ctrl_probe_iter110.go`; `git status --porcelain \| grep -c` | file listed; status count **1** | control fires — the instrument is blind, the file is not absent |
| D1 repair — planner's `+ git ls-files --others --exclude-standard` form | same probe | **1** | repair sees it |
| **D2** — AC2 green at base with one of two named tests existing | doc-form (no `-v`) | **rc=0** | **REPRODUCED** |
| D2 mechanism | same selector **with** `-v` | exactly **1** top-level `=== RUN`: `TestBareNetHTTPExemptionIsPerGroup` | the absent arm is invisible in the rc |
| **D2 negative control — the sharpening** | `go test ./host/boundary -run 'TestZzNoSuchTestIter110' -count=1 -v` | **rc=0**, **0** `=== RUN` lines | `go test -run` **returns rc=0 when it matches nothing**, so AC2-as-written cannot distinguish *both tests pass* from *neither test exists* |
| **D4** — AC8 likewise green at base | doc-form (which does carry `-v`) | **rc=0**, exactly **1** top-level `=== RUN`: `TestReadCtxCancelledAfterHandler` | **REPRODUCED**; the `-v` is present, so the criterion is machine-checkable on the enumeration count, which is what the plan does |
| AC5 detector | `grep -c 'mux.HandleFunc' host/daemon/daemon.go`; control `grep -c 'func '` | **8**; control **22** | matches the doc and iteration 109 |
| Base gates | `AILANG_BIN=/tmp/ailang-v0300/ailang GOTOOLCHAIN=go1.25.6 ./scripts/verify_ail.sh` / `./scripts/verify_go.sh`, rc captured to file, never through a pipe | **rc=0** / **rc=0** | pristine base is green, including `host/capsule` |
| Doc base currency | `git diff --name-only 93e1ba5..HEAD -- ':!design_docs'` | **0** files; control (no pathspec) **7**, all under `design_docs/` | the doc's §11 rows are fresh at HEAD |

**Why D1 is the finding of this iteration, and why the doc's author could not have seen it.** AC7 is
*"only priced files changed"* — the criterion that stops scope creep. It is vacuous for exactly the
files this sprint creates, and the reason is a property of **this loop's own executor recipe**, not
of the design: the sandboxed cross-provider executor performs no git write operations, so every new
file is untracked for the whole sprint, and `git diff` does not see untracked files. The doc was
reviewed by two quorum rounds plus a restored third reviewer and none of them could have caught it,
because none of them knows how this loop runs its executor. *Guard the helper, miss the call site* —
the mission's own recurring shape, arriving this time from the tooling side of the boundary rather
than the design side.

**The negative control on D2 generalises past this AC.** `go test -run <selector>` exits **0** on an
empty match set. Any acceptance criterion in this repo of the form *"`go test -run 'TestA|TestB'`
passes"* is therefore green before either test is written, and stays green if a rename silently
orphans the selector. That is a repo-wide instrument property, not an AC7/AC2 accident, and it is
recorded as such rather than patched here.


## 0. Read this first, executor

You are a cross-provider sandboxed lane (`codex:gpt-5.6-sol`, fallback `opus`) that has never read
this repository's skills. **This plan is your gate.** Everything you need is in
`w-workbench-read-only.plan.json` and in this file. Read the design doc for context, but the
plan's `features[].tasks` are what you execute.

Four rules that override anything you might infer:

1. **You perform NO git write operations at all.** No `add`, `commit`, `checkout`, `switch`,
   `stash`, `reset`, `worktree`, `push`, `branch`, `rm`, `mv`, `clean`. In a linked worktree the
   `.git` entry is a *file* pointing at a gitdir **outside** `--sandbox workspace-write`, so git
   writes are either denied or operate on state you cannot see. Instead, after each milestone you
   snapshot the cumulative tree into `.snap/M<k>/` (see plan `executor_git_policy`). The
   controller creates the commits.
2. **`git checkout -- <file>` is FORBIDDEN, especially during the mutation drill.** In an
   uncommitted worktree it *deletes the milestone*, it does not restore a mutant. Restore is
   always `cp` from `.snap/backup/`.
3. **Scope is the doc's, exactly.** Plan `out_of_scope.list` is §10 of the doc, binding.
   **Discovery of a requirement for any item on it STOPS THE SPRINT and returns to design.
   It is never absorbed as "small glue."** Write `.snap/STOP.md` and halt.
4. **Every wait is bounded.** The blocking-store arms use a 2 s watchdog; the gates have their own
   ceilings. There is no unbounded poll anywhere in this plan, and you must not add one.

---

## 1. Milestones

Eleven milestones, `WB.A` … `WB.K`. Each is independently buildable and testable, so the history
bisects; each is sized for a single 30-minute executor run. The four drill milestones are
additionally resumable at per-mutant granularity because `MUTATIONS.md` is append-only — a run
that ends mid-drill resumes at the first mutant with no section.

| Milestone | Deliverable | Doc ACs closed | Mutations (claim / discharge) | h |
|---|---|---|---|---|
| ~~WB.A~~ **LANDED** `83f1973` | `host/workbench` package, view model, grade/verdict constructors | — | claims M17, M18, M19, M21 (still discharged by WB.H) | 0.75 |
| ~~WB.B~~ **LANDED** | `Render` + one parsed `html/template`, landmarks, escaping, local-only links, unavailable states | — | claims M14, M15, M16, M20 (still discharged by WB.H; **M20's pin was hollow as specified — see §7c**) | 0.75 |
| WB.C | ninth registration `GET /workbench`, `handleWorkbench` happy path, security headers, §3.5 comment | **AC5, AC6** | claims M1, M22, M23 | 1.0 |
| WB.D | closed query grammar + every refusal branch | — | claims M2–M9, M13, M31, M32 | 1.0 |
| WB.E | payload opt-in, 64 KiB cap, 100-entry timeline cap | — | claims M10, M11, M12 | 0.75 |
| WB.F | `TestWorkbenchReadDeadline` + `/workbench` in the cancelled-after-handler table | **AC8** | claims M29, M30 | 0.75 |
| WB.G | boundary gate: transport-free renderer + daemon positive arm, overlay-driven | **AC2** | claims M24, M25 | 1.0 |
| WB.H | mutation drill 1/4 | — | **discharges M14–M21** | 0.75 |
| WB.I | mutation drill 2/4 | — | **discharges M1–M9** | 0.75 |
| WB.J | mutation drill 3/4 | — | **discharges M10–M13, M22, M23, M29–M32** | 1.0 |
| WB.K | mutation drill 4/4 + full gates + final acceptance | **AC1, AC3, AC4, AC7** | **discharges M24–M28** | 1.25 |

Total 9.5 h + 2.5 h contingency = 12 h ≈ **1.5 days**, inside the doc's §9 price of 1.5–2 days.

**The `mutations` column for WB.A–WB.G is a CLAIM, not a result.** It is discharged only by
WB.H–WB.K, and only by the protocol in §3 below.

---

## 2. Baselines — every acceptance command, measured on the pristine tree

All rows **VERIFIED BY ME** at `3e0c34c`, run outside any sandbox, exit codes captured to a file
and never through a pipe (`${PIPESTATUS[0]}` is silently empty in zsh).

| Command | rc | Note |
|---|---|---|
| `verify_ail.sh` (both env vars) | **0** | 11 modules; 10 identities / 40 named tests; 9/9 package steps |
| `verify_go.sh` (both env vars) | **0** | `✓ go gate PASSED … AILANG v0.30.0`; `host/capsule` green on this run |
| **AC1 doc form** `go test ./host/workbench ./host/daemon -count=1` | **1** | `directory not found`; control `go test ./host/daemon` → rc=0 |
| **AC1 executor form** (`-v` + 12 named `=== RUN`) | **1** | gotest rc=1, runlines=0 |
| `go build ./host/workbench/...` (narrowest gate) | **1** | `lstat ./host/workbench/: no such file or directory` |
| **AC2 doc form** | **0** | ⚠️ **GREEN AT BASE**; with `-v`: exactly 1 `=== RUN` line |
| **AC2 executor form** (`-v` + 2 named `=== RUN`) | **1** | gotest rc=0, runlines=1 |
| **AC3** `verify_go.sh` | **0** | regression pin only — never feature evidence |
| **AC4** `verify_ail.sh` | **0** | regression pin only — never feature evidence |
| **AC5 detector** `grep -c 'mux.HandleFunc' host/daemon/daemon.go` | prints **8** | |
| **AC5 executor form** (criterion `==9` + composition) | **1** | |
| **AC6** compound (doc form, unchanged) | **1** | exits at the root `test -d host/workbench` |
| **AC7 doc form** `git diff --name-only 93e1ba5 -- ':!design_docs'` | **0**, empty | ⚠️ **BROKEN — see D1**; control (no pathspec) prints 7 paths, all `design_docs/` |
| **AC8 doc form** | **0** | ⚠️ **GREEN AT BASE**; exactly 1 `=== RUN` line |
| **AC8 executor form** (2 named `=== RUN` + blocking-store arm) | **1** | gotest rc=0, runlines=1 |
| `gofmt -l host/ cmd/` | **0**, empty | |

Base currency: `git diff --name-only 93e1ba5 -- ':!design_docs'` returns **0 files** with a
known-positive control of 7 in the same call — so the doc's declared measurement base `93e1ba5`
is still current for every non-doc claim and its §11 Verification Log is fresh at `3e0c34c`.

**Narrowest-gate rule.** Prefer `go build ./host/workbench/...` (base rc=1) and
`go test ./host/daemon -run '^TestWorkbench' -count=1` over `verify_go.sh` during implementation.
`verify_go.sh` runs a plain leg and a `-race -timeout 8m` leg and takes 6–9 minutes on this rig;
**slow is not hung** — do not kill it before 12 minutes. It is budgeted once, in WB.K.

---

## 3. Non-vacuity is the spine

The doc's §6 mutation table has **32 rows**, M1–M32, including M31/M32 for the closed-grammar
query-parameter refusal the restored reviewer's verbatim fix introduced. Plan `mutations[]`
carries all 32 with, per row: the exact one-line mutant, the file, the named test, the milestone
that CLAIMS it, and the milestone that DISCHARGES it.

**A "kills which mutation" column is a CLAIM.** It is discharged only by:

1. applying the one-line mutant **and proving it landed** with `grep -n` on the mutated literal —
   a mutation that silently did not apply reports a false SURVIVED;
2. asserting the tree **BUILDS** (`go build ./...` rc=0) **before reading any test result**. For
   the two `_test.go` mutants (M25, M28) `go build` does not compile the file, so the build check
   is `go vet ./host/boundary` rc=0 **and** `go test ./host/boundary -run '^$' -count=1` rc=0;
3. **enumerating the red set by RUNNING it, never by predicting it** — record the complete
   `--- FAIL:` set from the classification arm, not just the expected member;
4. restoring **byte-identically with `cp`** from `.snap/backup/`, verified with `shasum -a 256`,
   then re-running the pristine control — **never `git checkout --`**, which in an uncommitted
   sprint worktree deletes the milestone.

A mutant whose named test does not red is **recorded as SURVIVED with its full output**. Do not
repair the mutant, do not repair the test to make it red, do not omit the row.

---

## 4. The closed grammar is load-bearing

§2.4's route/parameter posture is the restored reviewer's own replacement text, applied verbatim:

> **The workbench query grammar is closed. The only accepted keys are `world`, `object`, `from`,
> `entry`, and `payload`. Any unknown key, duplicate scalar key, or unsupported parameter
> combination returns `400 Bad Request` with a constant HTML error message; no parameter is
> ignored and no precedence fallback is applied.**

This is the defect the restored reviewer found — the prior rule "unknown parameters are ignored
only if they do not select data" was a **silent fallback**, and `?paylod=1` would have rendered a
different view instead of refusing. Pin it with `TestWorkbenchRefusalBranches/unknown-parameter`
and `/duplicate-parameter` (mutations M31/M32). Nothing is ignored. There is no precedence
fallback. Do not add one "for usability".

---

## 5. Worktree convention

The sprint worktree **must be a SIBLING of the repo** —
`/Users/voightkampff/dev/sunholo-data/.wt-iter110-workbench` — and **never under `/tmp`**.
`host/verifygate`, `host/boundary` and `host/daemon/read_deadline_test.go` all derive `repoRoot`
from `runtime.Caller`; a `/tmp`-rooted checkout reds CWD-resolving tests **for the LOCATION**, and
CI never reproduces it. On this rig `/tmp` is additionally a symlink to `/private/tmp`, which
`find` declines as a traversal root.

**The controller creates the worktree** (`git worktree add` is a git write) and must also write
`.snap/PRISTINE_MANIFEST.txt` and `.snap/PRISTINE_SHA256.txt` before the executor starts — see
plan `preconditions_the_controller_must_satisfy_before_the_executor_starts`.

---

## 6. Doc-vs-plan divergences — **the DOC wins; quoted verbatim, not silently reconciled**

Full detail in plan `doc_defects`. Summary, worst first.

### D1 (high) — AC7's command cannot see the files this sprint adds

Doc, §7 AC7, verbatim:

> ```sh
> git diff --name-only 93e1ba5 -- ':!design_docs'
> ```
>
> Baseline: empty output, rc=0 at revision time (V14). … For implementation, pass is that every
> listed path is one of the files enumerated in §8. Fail trigger: any unpriced non-`design_docs`
> path appears.

`git diff` **does not see untracked files**. All four files §8.1 adds are untracked until someone
commits, and this executor makes no commits — so AC7 as written reports EMPTY, i.e. **PASSES**,
with every new file present, and would pass identically with forty unpriced new files present.

**Measured, VERIFIED BY ME**: created `host/daemon/zz_probe_untracked.go`; the doc's form matched
the probe **0** times, the repaired form (`git diff` ∪ `git ls-files --others --exclude-standard`)
matched it **1** time. Probe removed; `git status --porcelain` re-confirmed empty.

**Repair**: the executor runs a **git-free** manifest + sha256 form (plan
`acceptance_criteria.AC7.executor_command`) — exactly 4 additions, 0 removals, exactly 3 sha256
mismatches and they must be `daemon.go`, `read_deadline_test.go`, `allowlist_world_test.go`. The
controller re-runs the repaired git form after committing. The doc's *intent* is preserved
exactly; only the instrument is repaired.

### D2 (medium) — AC2's command cannot produce the enumeration its own Pass clause demands

Doc, §7 AC2, verbatim:

> ```sh
> go test ./host/boundary -run 'Test(BareNetHTTPExemptionIsPerGroup|WorkbenchPackageRemainsTransportFree)' -count=1
> ```
> … Pass requires test enumeration/log output proving both names ran.

No `-v`, so `go test` prints only `ok  <pkg>`. And `go test -run <non-matching>` is not an error
on go1.25.6, so **the command is rc=0 GREEN at base** with only one of the two named tests in
existence — **VERIFIED BY ME**. Repair: add `-v`, machine-check for exactly 2 matching `=== RUN`
lines. Repaired baseline rc=1 — VERIFIED BY ME.

### D3 (medium) — AC1 goes green as soon as the directory exists

Doc, §7 AC1: `go test ./host/workbench ./host/daemon -count=1` … "Pass: all new renderer and
handler tests pass." From milestone WB.A onward the exit code is 0 while ten of the twelve §5/§6
named tests are unwritten. Repair: `-v` plus a 12-name `=== RUN` enumeration check. Repaired
baseline rc=1 — VERIFIED BY ME.

### D4 (low) — AC8's command is likewise green at base

Doc already states the right Pass condition ("`=== RUN` enumeration showing both names ran … a
blocking-store arm answering `503` with class `Timeout` on `/workbench`") but leaves it as prose
over a command measured rc=0 at base with one `=== RUN` line. Repair: machine-check it.

### D5 (low) — §5 omits two tests §6 requires

§5 says "Every property intended to survive implementation is attached to a committed Go test",
but its 16-row table contains neither `TestWorkbenchTimelineBound` (M12) nor
`TestWorkbenchSecurityHeaders` (M22, M23); both appear only in §6. §6 is a strict superset, so
this is a gap in §5's enumeration, not a contradiction. **Resolved as the UNION** — the plan
writes both (WB.E and WB.C).

### D6 (medium implementation hazard) — §8.2's "route table" is a shared helper

Doc, §8.2, verbatim:

> | `host/daemon/read_deadline_test.go` | add `GET /workbench` to the cancelled-after-handler route table and a workbench blocking-store arm (AC8) |

The instruction is right; the hazard it does not name is that the table comes from the **shared**
helper `seedReadRoutes` (`read_deadline_test.go:53–70`), consumed by five tests, **four of which
`json.Unmarshal` the error body** via `assertErrorClass` (`handlers_test.go:127–135`):
`TestDaemonReadDeadline/real-store-expired-deadline`, `TestDaemonReadDisconnect`,
`TestTimeoutStatusMirrorsSketch`, `TestInternalErrorsAreSanitized`. `/workbench` answers
`text/html`, so editing the shared helper reds all four for the wrong reason. **The plan forbids
touching `seedReadRoutes` (INV8) and prescribes a local `routes = append(...)` inside
`TestReadCtxCancelledAfterHandler` only**, with the HTML 503 arms in a separate
`TestWorkbenchReadDeadline`.

### D7 (informational) — §7 assigns no AC to any milestone

The doc has no milestone structure, so "make the plan's milestone AC lists identical to §7's" is
satisfied as: every AC title quoted is §7's own heading text verbatim; the **union** of the eleven
milestone AC lists is exactly {AC1…AC8}, none invented, none dropped; and no AC appears twice.
Verified mechanically against the JSON.

### D8 (informational) — `payload` spelling differs between the two routes

§2.2 specifies `payload=0|1` for `/workbench`; the existing JSON route uses `payload=true`
(`handlers.go:393`). Not a defect — the doc wins for `/workbench`, and §10 forbids changing frozen
`/v1` semantics. Recorded so nobody "harmonises" it later.

---

## 7. Three under-specified points — flagged by the planner, **ALL THREE ADJUDICATED by the controller**

Plan `open_questions_for_controller`. In each case the plan states a reading so the executor is not
blocked, and marks it as the plan's reading in the task text. **A design-direction ruling belongs
to the controller.**

- **Q1 — what is the 503's body format?** §2.4 says the deadline expiry returns "`503`, the same
  class the JSON routes emit (V18)", while the same section mandates HTML error pages and §2.5
  sets `Content-Type: text/html; charset=utf-8`. **Plan reading**: an HTML error page whose body
  contains the literal token `Timeout`; `writeAPIError`/`writeReadTimeout` (JSON) are NOT called
  from the workbench path. If the controller rules otherwise, `TestWorkbenchReadDeadline` changes
  from a body-substring check to `assertErrorClass`.
- **Q2 — which parameter combinations are "supported"?** §2.4 refuses an "unsupported parameter
  combination" without enumerating the supported set; §2.2 enumerates four query states.
  **Plan reading**: the accepted key sets are exactly `{}`, `{world}`, `{from, entry}`,
  `{object}`, `{object, payload}` — notably **`?from=0` alone is refused**. A more permissive
  reading is defensible but is a design decision, and any "ignore the rest" variant would be the
  silent fallback the carve-out forbids.
- **Q3 — malformed `payload` value?** §2.4's malformed-value list names only `world`, `object`,
  `from`, `entry`. **Plan reading**: `?payload=true` → 400, because ignoring it is a silent
  fallback. Pinned by a **PLAN-ADDED** subtest `TestWorkbenchRefusalBranches/malformed-payload`,
  labelled as such rather than laundered into the doc's list.

### 7a. Controller's ruling on Q1–Q3 (iteration 110) — **all three plan readings UPHELD**

None of the three is a design-direction dispute, which is why none of them is parked as a human
ask. Each is a **completeness gap** the doc's own text — plus, for Q1, one measurement — closes
uniquely. Rule 3f applies: a reviewer's or a planner's open premise is a claim to be *measured*,
not forwarded.

- **Q1 — UPHELD. The 503 body is HTML carrying the literal class token `Timeout`; the workbench
  path does NOT call `writeAPIError`/`writeReadTimeout`.** Measured: `writeReadTimeout` is
  `writeAPIError(w, "Timeout", …, http.StatusServiceUnavailable)` and `writeAPIError`
  (`host/daemon/handlers.go:134-136`) is `writeJSON(...)`. So §2.4's *"the same class the JSON
  routes emit"* names the **class token**, and the envelope is a separate thing the existing code
  keeps separate. §2.4 already fixes the media type for the sibling branches (*"a small HTML error
  page"*, *"a constant HTML error message"*) and §2.5 fixes `Content-Type: text/html` for the
  route, so a JSON 503 would contradict two sections to satisfy a phrase that was never about the
  envelope. **Class ≠ envelope** — that is the whole ambiguity.
- **Q2 — UPHELD. The accepted key sets are exactly `{}`, `{world}`, `{from, entry}`, `{object}`,
  `{object, payload}`; `?from=0` alone is `400`.** This is a conjunction of two doc sections, not a
  controller invention: §2.2 pairs `from` with `entry`, while §6's **M10** (*"payload remains
  opt-in"*, arm `TestWorkbenchPayloadPreviewBound/default-off`) requires `{object}` **without**
  `payload` to be accepted — which is why both `{object}` and `{object, payload}` are in the set. A
  permissive reading of `from` alone would have to invent a default for the absent `entry`, and
  §2.4's verbatim carve-out text forbids exactly that: *"no parameter is ignored and no precedence
  fallback is applied."*
- **Q3 — UPHELD. A malformed `payload` value is `400` with the constant message; keep the
  PLAN-ADDED subtest.** §2.4's omission of `payload` from its malformed-value list is an
  **enumeration gap, not a permission**. Treating `payload=true` as `0` is a silent fallback, i.e.
  the precise defect (`?paylod=1` rendering a different view instead of refusing) that the restored
  `gpt5-6-sol` reviewer found at round 3 and that the carve-out text exists to close.

The plan JSON carries the same three rulings under
`open_questions_for_controller[].controller_resolution`, so the executor reads them without needing
this file. **The executor is not blocked on any human answer.**


---

## 7b. Precondition 3's manifest count is wrong for the tree it is checked in (iteration 111)

`preconditions_the_controller_must_satisfy_before_the_executor_starts[3]` states the pristine
manifest is *"VERIFIED BY ME to be 1079 lines at 3e0c34c in the main checkout"*. Run where the
precondition actually applies — the **sprint worktree** — the identical command returns **157**.

| reading | value |
|---|---|
| `find … \| sort \| wc -l` in the sprint worktree | **157** |
| the same command in the MAIN checkout | **1079** |
| control: `git ls-files \| grep -v '^design_docs/' \| wc -l` | **156** |

The two trees differ by **922 untracked files** — `tools/` alone is **844** in the main checkout
against **13** tracked, plus 66 under `packages/` and 25 under `world/`. A fresh worktree has only
tracked files, so **157 = tracked + 1** is the honest count there.

**The `+1` is a second defect.** The command excludes `-not -path './.git/*'`, which assumes `.git`
is a **directory**. In a linked worktree `.git` is a **FILE**, so it is not under `./.git/` and it
lands in the manifest — measured, the sole difference against `git ls-files` is the literal entry
`.git`. §5 of this plan explains at length that a linked worktree's `.git` is a file pointing
outside the sandbox; this command does not act on what that prose knows.

**AC7 is NOT broken by either.** The delta is computed pristine-vs-post **inside one tree**, so a
constant offset and a self-cancelling `.git` row both drop out. What breaks is the **cross-check**:
a controller comparing its manifest against the stated 1079 sees a 6.9× mismatch and can only
conclude its worktree is broken. Use the worktree's own count, or add
`-not -name '.git'` and re-derive.

---

## 7c. `M20`'s named killer did not kill it — the row is repaired (iteration 112)

WB.B task 14 says: *extend `TestGradeViewRequiresTestVerdict/fail` to also `Render` a `Page`
carrying that `GradeView` and assert the body contains BOTH `TESTED` and `FAIL`. **This is what
makes M20 killable.*** The executor implemented that instruction **verbatim and correctly**.
Measured, it does not make M20 killable.

| arm | mutant LANDED | BUILDS | package rc | red set |
|---|---|---|---|---|
| M20 as the plan specified the pin | 2 literal occurrences, 0 remaining `{{.Grade.Verdict}}` actions | rc=0 | **0** | **empty — SURVIVED** |
| M20 against the repaired pin | same | rc=0 | 1 | `TestGradeViewRequiresTestVerdict/fail` **alone** (`-skip` that arm → rc=0) |

**Why the specified pin is hollow.** The FAIL span is
`<span class="verdict-fail" aria-label="test verdict FAIL">✗ verdict: {{.Grade.Verdict}}</span>`,
and the branch is selected by `{{if eq .Grade.Verdict "FAIL"}}` — which reads the **data**, not the
rendered action. So under M20 the page still emits the literal `aria-label="test verdict FAIL"`,
`strings.Contains(rendered, "FAIL")` is satisfied by a string the mutation never touches, and the
assertion passes. This is the shared skill's rule 3i extension exactly: **the observable's value
set is larger than the mechanism's** — ask not only *which write does this read?* but *what else
writes this value?*

**The repair, which is the ROW and not the code.** The assertion now requires `verdict: FAIL` — a
string only the `{{.Grade.Verdict}}` action can produce, because the aria-label spells it without
a colon — and, separately, the accessible label, keeping §2.6's dual channel pinned in its own
right. Task 14's text is corrected in `plan.json` to specify those two observables instead of the
bare token.

**The tension with §3 rule 5, stated rather than papered over.** §3 says *"A mutant whose named
test does not red is recorded as SURVIVED with its full output. Do not repair the mutant, do not
repair the test to make it red, do not omit the row."* That rule exists to stop a **drill** being
laundered into a pass, and it is right. This is not a drill: it is a spot-check of a *claim the
milestone's own task text makes about itself*, the survival is published here and in the commit
message with its full measurement, and the fix changes **which value is observed** rather than
manufacturing a red. Left alone, WB.H's drill would record M20 SURVIVED and produce this identical
repair one milestone later. Read strictly, though, §3 rule 5 covers this case and forbids it —
so the honest reading is that §3 rule 5 needs the scope it was written with: **it governs the
discharge milestones WB.H–WB.K, not a spec defect found outside them.**

**What is NOT discharged.** WB.H still owns M14–M21. This section records a spot-check of four
claims (M14 sole killer `TestRenderEscapesAllObjectText`; M15 red set of two, named arm included
and the second member explained — it asserts the same local-href form; M16 sole killer
`TestRenderUnavailableProvenanceEdge`; M20 sole killer after the repair). Every mutant was
asserted LANDED by occurrence count and BUILDS rc=0 before any test result was read, and restored
by `cp` from a controller backup with `shasum -a 256 -c` byte-identity — never `git checkout --`.

---

## 8. Known base hazards — **not this sprint's fault, do not absorb**

- **`host/capsule` `TestF5WallClockTimeoutHasElapsedBound`** asserts an absolute 2 s wall-clock
  ceiling and has been observed red *only* inside `verify_go.sh`, where `host/broker` and
  `host/verifygate` run concurrently; 10/10 green in isolation. **It passed on my baseline run at
  `3e0c34c`.** If it reds: re-run `go test ./host/capsule -count=1` in isolation, record both
  results, report it as a load artifact. **That is charter queue row 32.**
- **`verify_go.sh`'s driver drift gate** over `tools/launchd/` and `scripts/mission_decisions.sh`.
  Those are **frozen core** — never touch them. A red there means "the fleet must commit", never
  "absorb it into your change".
- **Both gates require `AILANG_BIN` AND `GOTOOLCHAIN=go1.25.6`.** `verify_go.sh` fails loudly if
  `AILANG_BIN` is unset or ≠ v0.30.0. The doc's AC3/AC4 commands omit `GOTOOLCHAIN`; the executor
  forms in the plan add it.

---

## 9. Definition of done

See plan `definition_of_done`. In short: eleven snapshots with evidence; AC1/AC2/AC5/AC6/AC8 pass
in their **executor** forms with the measured base rc recorded beside the post rc (all five are
rc=1 at base); AC3/AC4 rc=0 with exact totals, cited **only** as regression pins; AC7's manifest
delta exactly 4 additions / 0 removals / 3 sha mismatches; **all 32 mutation rows discharged by
measurement**; INV1–INV10 hold; `gofmt -l host/ cmd/` empty; no §10 item built; no git write
command anywhere in the sprint's shell history.
