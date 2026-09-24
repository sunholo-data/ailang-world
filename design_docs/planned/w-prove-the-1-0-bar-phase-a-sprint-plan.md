# w-prove-the-1-0-bar — Phase A (row 92) sprint plan

- Sprint: `w-prove-1-0-phase-a` · queue row **92** · design doc
  [`w-prove-the-1-0-bar.md`](w-prove-the-1-0-bar.md) (read §2/§4/§7/§8/§9 and the Conflict
  Surface first) · Phase B (row 93) is **OUT** of this sprint.
- Base: `8dd867663d0b9db4d6f6f5e51bbdfb10acf939c6` · branch `sprint/w-prove-1-0-phase-a` ·
  worktree `/Users/voightkampff/dev/sunholo-data/.wt-world-iter182`
- Planned 2026-09-24 by mission-control iteration 182 (pi:openrouter/moonshotai/kimi-k3
  sprint-planner). Baseline gates measured green at base by this planner (see §8).
- **The artifact is the deliverable, not tooling** (doc §4.5, §9). New code is a fixture
  generator, the checked-in commit fixtures, and ONE small durability test. Total new Go
  ≈ 140 LOC; the whole sprint ≈ 270 LOC of code + the artifact.

## §1 What Phase A proves, and what counts as failure

Clause 5: on ≥3 REAL "why did X happen" questions from this mission's own recorded operation,
a provenance walk yields the **verified** answer in **≤5 minutes each**, where the pre-World
method was grep/log archaeology. The four anti-vacuity words (doc §2) are enforced by
construction in this plan:

| Word | How this plan enforces it |
|---|---|
| REAL | Questions are pre-harvested from charter rows 86/78/81 (+ one eligibility-gated candidate from the iter-181 record), not invented to fit the walk. The rejection log is prescribed below with named rejected candidates. |
| ≥3 | Three questions are mandatory; a fourth is taken iff its scripted eligibility check passes. |
| verified | Every question carries a verification source the walk did not produce: the surviving v0.30.0 binary, the live fleet script via `gh api`, this repo's git history, the in-git charter/log text. An answer with no independent confirmation is recorded **UNVERIFIED** and does not count. |
| pre-World method | Baselines are quoted from the mission record with file+section citations; where the record states no cost the cell says **UNMEASURED** — never imputed (doc §4.2, §8). |

**Honest failure is a result** (doc §4). The verdict closed set is:
`PASS` (verified AND ≤300 s) · `FAIL-TIME` (verified, >300 s) · `UNVERIFIED` (no independent
source confirmed) · `UNANSWERABLE` (the walk cannot produce the answer from the committed
record) · `NOT-RUN(TRANSPORT)` (the sandbox refused the listener or the pty mint — labelled,
never reported as pass/fail). A failed or unanswerable walk is clause-5 evidence and is
recorded, never substituted with an easier question.

## §2 The harvest (prescribed — M1 executes and documents it)

**Taken (mandatory 3):**

| # | Slug | Question | Ground truth (independent of the walk) | Record |
|---|---|---|---|---|
| Q1 | `iter171-index-row` | Why did iteration 171 have a full log entry and no index row? | The pinned **v0.30.0** binary lacks the `mission` command group: `ailang.v0.30.0 mission rotate-log world --keep 31` prints `Error: unknown command 'mission'`. Planner-verified 2026-09-24: `$HOME/.pinned-ailang/ailang.v0.30.0` survives on the rig and reproduces it. The index hole is visible in **git history** of `world-mission-index.md` (the row exists today because it was added manually after the iter-172 rule — `grep -c '^| 171 |'` is 1 NOW; that is the repair, not the incident). | charter row 86 + the `mission`-command bullet (`world-mission.md` ~lines 82–101); log `## 172` |
| Q2 | `iter154-unfenced-pi` | Why did three designer runs execute unfenced while the rulebook said they were sandboxed? | The fleet's `scripts/mission_pi_run.sh` invokes `pi --mode json --no-session --model "$MODEL" < "$DIRECTIVE"` with **no `-e` flag and no `PI_FENCE_ROOT` anywhere in the file** — planner-verified live 2026-09-24 via `gh api` (invocation at :165 today; row 78 recorded :155 at iter-154 — line drift, property unchanged); issue `sunholo-data/ailang#1043` is **OPEN**. | charter row 78; log `## 154` containment note |
| Q3 | `iter155-faithfulness-proof` | Why did a green `shasum -c` faithfulness proof pass over a destroyed commit split? | The prescribed manifest covers the **final tree**, identical whether the split is correct or collapsed; `git add design_docs host` staged disk state, MS1 swallowed MS3's files, MS3 committed nothing, and `shasum -c` returned OK on every file. Corroborated by two in-git records (row 81 + log `## 155`) and by the rebuilt row-59 split in git history. | charter row 81; log `## 155` |

**Considered additionally (the rejection log — a harvest with no rejections is evidence the
bar was set after the fact, doc §4.1):**

| # | Candidate | Decision | Reason |
|---|---|---|---|
| C4 | Why did PR #141's remote CI go red on `BenchmarkRESTCommit` when every local gate was green? (iter-181) | **TAKE iff the scripted eligibility check passes** (M1, below) | Real, recorded (STATUS 2026-09-24 / iter-181), independent ground truth in-git: `a036062` is an ancestor of base and its diff adds `req.Header.Set("Authorization", auth)` to `bench_test.go` — the bench predated the session gate and 401'd. Planner-verified 2026-09-24. |
| C5 | Why did iteration 180 refuse every push/PR (prepush-gate unpassable, row 95)? | **REJECT** | Ground truth lives in a rig-global pi extension (`~/.pi/agent/extensions/prepush-gate.ts`) whose defective bytes have since been **replaced** by the attended fix (ailang `754fc7a43`); the original bytes no longer exist anywhere the executor may read, and the fix commit is in a different repo. The only surviving record is the charter text itself — the walk would be citing its own source (§2 circularity risk). |
| C6 | Why did a mutation-drill restore destroy the executor's uncommitted work (row 82)? | **REJECT** | The destroyed work was uncommitted by definition; no independent artifact of it survives anywhere. The row text is the only record, so verification would reduce to "the walk repeated the row" — exactly the circular shape §2's `verified` forbids. |

M1 records this table in the artifact verbatim, plus any further candidate the executor spots
in the record, each with take/reject and the reason. **The executor may not add a taken
question beyond C4 without a recorded eligibility proof** (real operation · independent ground
truth · baseline method recoverable).

**C4 eligibility check (run in M1, all must hold):**

```bash
git merge-base --is-ancestor a036062 HEAD && echo ELIGIBLE-1
git show a036062 -- host/daemon/bench_test.go | grep -c 'req.Header.Set("Authorization", auth)'   # expect >= 1
grep -c "bench predates the sprint" design_docs/world-mission.md                                  # expect >= 1
```

If any fails, C4 is recorded REJECTED with the failing check as the reason, and the sprint
proceeds with Q1–Q3.

## §3 Baselines (prescribed — M1 fills the artifact's baseline table)

For each taken question the artifact records: the pre-World method **as enumerated in the
record** (file + section citation, step count as recorded), and the measured cost **quoted
from the record** or the literal string `UNMEASURED (record states method, not cost)`.

Where to look (planner pre-read; the executor re-reads and quotes):

| Q | Method citation | Cost in the record |
|---|---|---|
| Q1 | `world-mission.md` row-86 bullet + the v0.30.0 `mission`-command bullet (~lines 82–101): grep the index for the row, control-grep a neighbour, run the prescribed regenerator, positive-control another command group | No duration stated → UNMEASURED expected |
| Q2 | `world-mission.md` row 78; log `## 154` containment note: read the recipe's two-flag invocation block, read the runner's actual invocation, compare | No duration stated → UNMEASURED expected |
| Q3 | `world-mission.md` row 81; log `## 155`: trust three green `git commit`s, get a green `shasum -c`, then catch it only by reading `git show --stat` per commit; rebuild; hit the zsh word-splitting trap on retry | No duration stated → UNMEASURED expected |
| C4 | `world-mission.md` STATUS 2026-09-24 (iter-181): local gates all green, remote CI red, read the CI failure, reproduce, judge round 2 | No duration stated → UNMEASURED expected |

A duration, tool-count, or dollar figure **attributed in the record to the diagnosis itself**
may be quoted (with its citation). An iteration's whole-fire metered cost is NOT the diagnosis
cost — never impute (doc §4.2, §8 "Not verified, and named as such").

## §4 The walk mechanics (measured against this tree at base)

Measured facts the plan relies on (planner, 2026-09-24, base `8dd8676`):

- `POST /v1/commit` is **session-gated** since `a036062` (row 39 landed iter-181): the
  middleware resolves `Authorization: Bearer <64-hex>`; all GET routes are open. Mint is
  **tty-fenced** (`cmd/ailang-worldd/session.go`: `os.OpenFile("/dev/tty")`, y/N confirm) and
  speaks to `--db` directly, never HTTP.
- Commit wire schema (`host/daemon/handlers.go`, decoder is `DisallowUnknownFields`):
  `observedHead` · `objects[]{hash,interfaceHash,semanticId,provenance,payload}` (payload
  base64; `hash` MUST be `sha256:<hex>` of the payload bytes — the store content-verifies) ·
  `nextWorld{ref,revision,stateRoot,logHead}` ·
  `entry{header{entryIndex,semanticsEpoch,transitionFn,interpreter,prevEntryHash,writtenBy},entryHash,transitionRef}`.
  Genesis: `observedHead: ""` is the ONE lenient field; `prevEntryHash` must be a REAL hash
  even at genesis (a zero prevEntryHash is writable but **not readable** —
  `TestGenesisRefLenienceIsExactlyOneField`).
- `serve --ailang-bin <path>` archives the interpreter by hashing the **raw file bytes**, so
  `health`'s `interpreter_ref` == `"sha256:" + shasum -a 256 <file>` (planner-verified against
  `host/archive/archive.go`; pinned binary sha256 at plan time:
  `1a67b0146858450182f48082299956ea5b08cdf30131979b188b339fcedb5b9f`).
- **Finding F-1 (walk entry point):** `GET /v1/objects/{ref}` accepts a hashref ONLY, and
  `POST /v1/commit` cannot register names in `epoch_registry_heads` (its only writers are
  `registry.Bootstrap`, broker approve, transitionreg). So the walk locates an incident by
  scanning `log range` and fetching each `transitionRef` object's metadata until the
  semanticId matches. Fine at 3–4 entries; it is a genuine provenance-query gap → the
  artifact's findings section proposes a queue row (`w-object-lookup-by-semantic-id`). A
  finding, NOT scope (doc §9).

### Fixture content (built by M2's generator from this table)

One commit per question, chained: Q1 genesis (entryIndex 0) → Q2 (1) → Q3 (2) → [C4 (3)].
Each commit's objects: one **incident object** (`semanticId: world/mission/incident/<slug>`,
`provenance: w-prove-1-0-phase-a`) whose payload is JSON
`{"question":…,"answer":…,"diagnosed_at":…,"sources":[<evidence object hashes>]}` plus one
**evidence object** per source (`semanticId: world/mission/evidence/<slug>/<n>`) whose payload
is JSON `{"kind":…,"ref":…,"check":…,"excerpt":…}`. `transitionRef` = the incident object's
hash; `header.interpreter` = sha256 of `$HOME/.pinned-ailang/ailang` bytes (the replay pin);
`writtenBy: w-prove-1-0-phase-a`; `semanticsEpoch: 1`.

Answers and evidence excerpts are the §2 table's ground truths, quoted from the record. The
generator carries them in a FACTS table; the executor does not paraphrase.

## §5 Milestones

| MS | Name | New code | Files | Gate |
|---|---|---|---|---|
| M1 | Harvest + baseline + artifact skeleton | 0 LOC | `design_docs/verification/w-1-0-value-demonstration.md` | G1 |
| M2 | Commit fixtures (+ generator) | ~130 LOC python | `design_docs/verification/w-1-0-value-demonstration/gen_fixtures.py`, `…/commits/q1-iter171-index-row.json`, `q2-iter154-unfenced-pi.json`, `q3-iter155-faithfulness-proof.json`, conditional `q4-iter181-ci-bench-401.json` | G1 + G2 |
| M3 | Durability test | ~140 LOC Go | `host/daemon/value_demonstration_test.go` | G1 + G3 |
| M4 | Live timed walk + artifact completion | 0 LOC | artifact tables filled + `…/walk-transcript.md` | G1 + G4 |

**G1 (every milestone, with `AILANG_BIN=$HOME/.pinned-ailang/ailang` exported):**
`go build ./...` rc=0 · `go vet ./...` rc=0 · `go test ./... -count=1` → **20 ok / 0 FAIL**
(measured at base) · `bash ./scripts/verify_ail.sh` → `✓ verify gate PASSED: 11 required
identities verified, 40 named tests pass` and `✓ world package gate PASSED: 9/9 steps`
(unchanged — no `.ail` is touched).
**G2 (M2):** `python3 design_docs/verification/w-1-0-value-demonstration/gen_fixtures.py --check`
rc=0 (recomputes every object hash == sha256(payload), asserts exact schema keys, contiguous
entryIndex chain, prevEntryHash chaining, interpreter == sha256 of the pinned binary).
**G3 (M3):** `go test ./host/daemon/ -run 'TestValueDemo' -count=1 -v` → both named tests PASS.
**G4 (M4):** artifact-completeness greps (below) all rc=0.

### M1 — Harvest + baseline + artifact skeleton (≈0.2d)

Executor directive sketch:

1. Create `design_docs/verification/w-1-0-value-demonstration.md` with EXACTLY these sections:
   `# 1.0 value demonstration (w-prove-1-0-phase-a)`, `## Harvest log` (the §2 taken table +
   C4 eligibility check output pasted verbatim + the C5/C6 rejection table), `## Baselines`
   (the §3 table filled by re-reading the record — quote, cite, or `UNMEASURED`), `## Walk
   protocol` (§6 of this plan, copied), `## Timing log` (empty prescribed table),
   `## Verification` (empty prescribed table), `## Comparison` (empty prescribed table),
   `## Findings` (F-1 and F-2 seeded verbatim from §4/§6), `## Walk transcript` pointer.
2. Run the C4 eligibility check (§2); paste its output; record TAKE or REJECT with the reason.
3. Run G1. Snapshot: copy the artifact into `.snap/M1/` (full content).

Acceptance criteria:

- **AC-M1-1** Harvest log lists ≥3 taken questions with ground-truth citations AND ≥2 rejected
  candidates with reasons.
  `grep -c "REJECT" design_docs/verification/w-1-0-value-demonstration.md` ≥ 2 and
  `grep -c "world/mission/incident/" design_docs/verification/w-1-0-value-demonstration.md` ≥ 3.
- **AC-M1-2** Every baseline row carries a record citation (`world-mission.md` /
  `world-mission-log.md` + section) and its cost cell is either a quoted figure with citation
  or the literal `UNMEASURED`. A scripted grep finds zero empty cost cells:
  `grep -c "UNMEASURED" <artifact>` ≥ 3 is EXPECTED (doc §8 anticipates this) — any row that is
  neither a quoted cost nor UNMEASURED fails review.
- **AC-M1-3** C4 decision recorded with its check output verbatim; if TAKE, a fourth incident
  slug `iter181-ci-bench-401` appears in the harvest log.
- G1 green. Snapshot `.snap/M1/` written.

### M2 — Commit fixtures + generator (≈0.25d)

Executor directive sketch:

1. Write `design_docs/verification/w-1-0-value-demonstration/gen_fixtures.py` (~130 LOC,
   stdlib only: `json, hashlib, base64, pathlib, argparse`). It carries the FACTS table
   (slugs, questions, answers, evidence items from §2/§4), computes all hashes, chains the
   commits, and emits `commits/q<k>-<slug>.json` (Q1 genesis: `observedHead: ""`,
   `prevEntryHash: sha256("w-1-0-value-demo/genesis-prev")`; each successor: `observedHead` =
   prior `nextWorld.ref`, `prevEntryHash` = prior `entryHash`). `header.interpreter` =
   `"sha256:" + sha256(file bytes of $HOME/.pinned-ailang/ailang)` — computed by reading the
   binary, never hard-coded. `--check` mode re-reads the emitted files and asserts: exact
   top-level/object/header key sets (the decoder is `DisallowUnknownFields`), every object
   `hash` == `sha256:<hex>` of the base64-decoded payload, entryIndex 0..N contiguous,
   prevEntryHash chaining, interpreter matches the pinned binary's current sha256.
2. Run it: `python3 …/gen_fixtures.py --out …/commits` then `--check` (G2).
3. Eyeball one fixture against the wire schema in `host/daemon/handlers.go:72-107`.
4. Run G1. Snapshot: generator + fixtures + (unchanged) artifact into `.snap/M2/` —
   cumulative, includes M1's file.

Acceptance criteria:

- **AC-M2-1** `python3 …/gen_fixtures.py --check` rc=0 (G2).
- **AC-M2-2** 3 fixtures exist, 4 iff C4 was taken:
  `ls design_docs/verification/w-1-0-value-demonstration/commits/q*.json | wc -l` ∈ {3,4} and
  matches the harvest log's taken count.
- **AC-M2-3** Every fixture's `semanticId`s carry the `world/mission/` prefix and its
  `provenance` is `w-prove-1-0-phase-a`: `grep -L 'w-prove-1-0-phase-a' …/commits/q*.json`
  prints nothing.
- G1 green. Snapshot `.snap/M2/` written (cumulative).

### M3 — Durability test (≈0.2d)

One new file, `host/daemon/value_demonstration_test.go` (~140 LOC), same package so it reuses
`newHandlerDaemon`, `authHeader`, and the existing httptest patterns. It reads the checked-in
fixtures from `../../design_docs/verification/w-1-0-value-demonstration/commits/` (test
working directory is the package dir), POSTs them in order through `httptest.NewServer(d.Handler())`
with the session header, then performs the walk IN-PROCESS (GET `/v1/log?from=0`, scan
transitionRefs, GET objects) and asserts:

- **TestValueDemo_AnswersMatchRecord** — for each taken question, the incident object located
  by semanticId `world/mission/incident/<slug>` exists and its payload's `answer` field
  byte-equals the expectation table entry (the table in the test file mirrors the FACTS
  answers). **Kills mutations MU-1 and MU-3.**
- **TestValueDemo_EvidenceChainComplete** — every hash in every incident payload's `sources[]`
  resolves via `GET /v1/objects/<hash>?payload=true` (200, payload present) and each evidence
  payload's `excerpt` byte-equals the expectation table. **Kills mutation MU-2.**

Executor directive sketch:

1. Write the test file; the expectation table is keyed by slug and holds `answer` +
   `[]excerpt` exactly as in the generator's FACTS (copy, do not paraphrase).
2. `go test ./host/daemon/ -run 'TestValueDemo' -count=1 -v` (G3) → 2 PASS, non-vacuous
   (each test fails its own expectation-table entry if the entry is emptied — self-check by
   inspection, asserted by MU-1 in the battery).
3. Run G1 (full suite still **20 ok / 0 FAIL** — the test lands in the existing `host/daemon`
   package, no new package). Snapshot: cumulative `.snap/M3/` (M1+M2 files + the test).
4. **STOP.** The pre-registered mutation battery (§7) is the controller's next act.

Acceptance criteria:

- **AC-M3-1** G3: both named tests PASS against the checked-in fixtures replayed through the
  real daemon stack (httptest + session-gated POST). Killers: MU-1, MU-3.
- **AC-M3-2** The evidence chain is asserted end-to-end: every `sources[]` hash resolves and
  its excerpt matches. Killer: MU-2.
- **AC-M3-3** G1 green with 20 ok / 0 FAIL; `verify_ail.sh` byte-identical outcome (no `.ail`).
- Snapshot `.snap/M3/` written (cumulative).

### M4 — Live timed walk + artifact completion (≈0.35d)

The walk runs against a scratch DB inside the worktree (`.walk-scratch/`, untracked, never
committed, never snapshotted). All client verbs use `--addr 127.0.0.1:7654` (non-default port,
avoids colliding with any rig daemon on the 7644 default).

**Setup (once; untimed, but its wall-clock is recorded in the transcript as context):**

```bash
cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter182
export AILANG_BIN=$HOME/.pinned-ailang/ailang
mkdir -p .walk-scratch
go build -o .walk-scratch/ailang-worldd ./cmd/ailang-worldd
.walk-scratch/ailang-worldd --addr 127.0.0.1:7654 serve --db .walk-scratch/walk.db \
  --bind 127.0.0.1:7654 --ailang-bin "$AILANG_BIN" > .walk-scratch/serve.log 2>&1 &
for i in $(seq 1 50); do .walk-scratch/ailang-worldd --addr 127.0.0.1:7654 health >/dev/null 2>&1 && break; sleep 0.2; done
.walk-scratch/ailang-worldd --addr 127.0.0.1:7654 health
# replay-pin pre-check: health's interpreter_ref must equal the value baked into the fixtures:
test "$(.walk-scratch/ailang-worldd --addr 127.0.0.1:7654 health | python3 -c 'import json,sys; print(json.load(sys.stdin)["interpreter_ref"])')" \
     = "sha256:$(shasum -a 256 "$AILANG_BIN" | awk '{print $1}')" && echo REPLAY-PIN-OK
# tty-fenced mint through a pty wrapper (the fence is intended behavior — row 39 D1):
printf 'y\n' | python3 -c "import pty,sys; pty.spawn(sys.argv[1:])" \
  .walk-scratch/ailang-worldd session mint --db .walk-scratch/walk.db \
  --episode w-1-0-value-demo --grant fs.read=/tmp:1 --ttl 3600 --out .walk-scratch/token
test "$(wc -c < .walk-scratch/token | tr -d ' ')" = "64" && echo TOKEN-OK
# commit the fixtures in order:
for f in design_docs/verification/w-1-0-value-demonstration/commits/q*.json; do
  .walk-scratch/ailang-worldd --addr 127.0.0.1:7654 commit --file "$f" --session "$(cat .walk-scratch/token)"
done
```

If `serve` cannot bind, or both the `pty.spawn` wrapper and the `script -q /dev/null …`
alternate cannot satisfy the tty fence: verdict **NOT-RUN(TRANSPORT)** for every live-walk
row, labelled in the artifact (never pass/fail); the M3 durability test stands as the
regression guard; the transport block is reported as a finding. (Pre-registered PF-1/PF-2.)

**Per question (Q1 shown; Q2/Q3/[Q4] identical in shape with their own slugs and verification
blocks). The timer starts at the question — before the first query:**

```bash
T0=$(date +%s); date -u +%Y-%m-%dT%H:%M:%SZ          # start clock: "why did iteration 171 have a full log entry and no index row?"
# 1. locate: scan the log, read object metadata until the incident semanticId matches
.walk-scratch/ailang-worldd --addr 127.0.0.1:7654 log range --from 0
.walk-scratch/ailang-worldd --addr 127.0.0.1:7654 object get <transitionRef>     # repeat per entry until semanticId == world/mission/incident/iter171-index-row
# 2. read the incident payload (base64) -> question, answer, sources[]
.walk-scratch/ailang-worldd --addr 127.0.0.1:7654 object get <hash> --payload | python3 -c 'import json,sys,base64; print(base64.b64decode(json.load(sys.stdin)["payload"]).decode())'
# 3. follow every source hash the same way -> evidence payloads
# 4. assemble the answer from the incident object's "answer" field
# 5. VERIFY against ground truth OUTSIDE World:
"$HOME/.pinned-ailang/ailang.v0.30.0" mission rotate-log world --keep 31 2>&1 | grep -F "unknown command 'mission'"   # expect a match
git log --oneline -S'| 171 |' -- design_docs/world-mission-index.md | tail -3   # the row's first appearance is the manual repair, not iter-171's Gate 4
T1=$(date +%s); echo "duration_s=$((T1-T0))"          # stop clock
```

Verification blocks per question:

| Q | Independent verification (recorded verbatim into the transcript) |
|---|---|
| Q1 | `ailang.v0.30.0 mission rotate-log …` grep match (above) AND `git log -S'| 171 \|'` showing the row first appears in the manual-repair commit. Do NOT check against the current pin: v0.41.0 HAS the command (planner-measured) — that is the pin bump, not a refutation. |
| Q2 | `gh api repos/sunholo-data/ailang/contents/scripts/mission_pi_run.sh --jq .content \| base64 -d \| grep -n 'pi --mode\|PI_FENCE_ROOT\| -e '` → the invocation matches, ZERO `-e`/`PI_FENCE_ROOT` hits; `gh issue view 1043 --repo sunholo-data/ailang --json state --jq .state` → `OPEN`. If the network/`gh` is unavailable: fall back to the in-git row-78 text as primary and record the live check `UNAVAILABLE` (PF-3); the answer still counts iff the in-git record independently confirms it. |
| Q3 | `grep -n "manifest is over the FINAL TREE" design_docs/world-mission.md` (row 81) AND `grep -n "commit reconstruction was wrong and I rebuilt it" design_docs/world-mission-log.md` (log 155) AND `git show --stat` on the landed row-59 commits (2026-09-05) naming each milestone's own files — the rebuilt split corroborates the incident narrative. |
| Q4 | `git show a036062 -- host/daemon/bench_test.go \| grep -c 'req.Header.Set("Authorization", auth)'` ≥ 1 AND `grep -n "bench predates the sprint" design_docs/world-mission.md`. |

**Fill the artifact** with the prescribed tables (formats are contracts — the completeness
greps pattern-match them):

```
## Timing log
| question | start_utc | stop_utc | duration_s | verdict |
## Verification
| question | answer (from walk) | verification source (command, verbatim) | confirmed? |
## Comparison
| question | baseline method (citation) | baseline cost | walk time (s) | verification source | verdict |
```

Append the raw command output to `design_docs/verification/w-1-0-value-demonstration/walk-transcript.md`
per question (commands + stdout, unedited). Stop the daemon (`kill %1` or SIGTERM) and report
the drain line.

Acceptance criteria:

- **AC-M4-1** Every taken question has a Timing-log row with a verdict from the closed set
  {PASS, FAIL-TIME, UNVERIFIED, UNANSWERABLE, NOT-RUN(TRANSPORT)} and a numeric `duration_s`
  (or `-` iff NOT-RUN(TRANSPORT)). `grep -cE '\| (PASS|FAIL-TIME|UNVERIFIED|UNANSWERABLE|NOT-RUN)'`
  over the artifact ≥ taken-count, and zero `TBD`/empty cells remain:
  `grep -c "TBD" <artifact>` == 0.
- **AC-M4-2** Every PASS row has `duration_s ≤ 300` and its Verification row names a source
  OUTSIDE World whose command is pasted verbatim in the transcript.
- **AC-M4-3** The Comparison table has the six prescribed columns and ≥3 rows; the harvest
  log (with rejections) is present in the same artifact (doc §4.5: the artifact IS the
  deliverable).
- **AC-M4-4** G1 green; `verify_ail.sh` outcome unchanged; the Findings section carries F-1
  (semantic-ID query gap → proposed queue row) and F-2 (registry names not settable via
  commit) plus any transport findings (PF-1/PF-2/PF-3 if fired).
- Snapshot `.snap/M4/` written (cumulative; `.walk-scratch/` is EXCLUDED — untracked scratch).

## §6 Findings pre-seeded for the artifact

- **F-1 — no semantic-ID query route.** `GET /v1/objects/{ref}` is hashref-only; a walk must
  linear-scan `log range` + per-entry `object get` metadata to locate an incident by
  semanticId. At mission scale (thousands of entries) this breaks the ≤5-minute bar's spirit.
  Proposed queue row: `w-object-lookup-by-semantic-id` (read route `GET /v1/objects/by-semantic-id/{name...}`
  or a log index over semanticIds). **Finding, not scope** (doc §9).
- **F-2 — registry names are not commit-settable.** `epoch_registry_heads` is written only by
  `registry.Bootstrap`/broker/transitionreg, so a walker cannot mint a well-known registry
  name for an incident class via `POST /v1/commit`. Compounds F-1; same proposed row covers it.

## §7 Pre-registered mutation battery (controller-run, post-M3 — NOT executor work)

Discipline: `cp -p <file> <file>.bak` → `sed -i ''` with an asserted occurrence count →
sha256 landed-proof (pre ≠ post) → compile/vet assert BEFORE any test result is read → the
scoped observable read → the named test required as the SOLE failure → restore by `cp` from
the `.bak`, asserted by `shasum -a 256` equality (and `test -x` where a mode bit matters) →
NEVER `git checkout --`.

| ID | Target | Mutation | Sole killer | Guards |
|---|---|---|---|---|
| MU-1 | `host/daemon/value_demonstration_test.go` | Replace Q1's expected `answer` string in the expectation table with a wrong value (sed count asserted = 1) | `TestValueDemo_AnswersMatchRecord` | AC-M3-1 (proves the assertion compares content, not presence) |
| MU-2 | same | Collapse the evidence-fetch loop to `_ = incident.Sources` (skip resolution) | `TestValueDemo_EvidenceChainComplete` | AC-M3-2 |
| MU-3 | `…/commits/q1-iter171-index-row.json` | Flip one base64 character in the incident object's `payload` (sed count = 1) | `TestValueDemo_AnswersMatchRecord` (replay POST 400s on the store's content-verify, or the answer mismatches) | AC-M3-1; proves the fixtures are load-bearing, not decorative |
| MU-4 | same test file | GREEN CONTROL — reword one comment, no behavior change | none — ALL tests stay green | non-vacuity of the battery |

## §8 Baseline gate measurement (planner, 2026-09-24, this worktree at `8dd8676`)

- `go build ./...` rc=0 · `go vet ./...` rc=0
- `go test ./... -count=1` (AILANG_BIN exported) → **20 ok, 0 FAIL**
- `bash ./scripts/verify_ail.sh` → `✓ verify gate PASSED: 11 required identities verified, 40
  named tests pass` · `✓ world package gate PASSED: 9/9 steps performed non-zero work` ·
  compiler pinned by exact bytes: **AILANG v0.41.0** on Darwin/arm64
- A bare `go test ./...` WITHOUT `AILANG_BIN` exported fails — base condition, never a
  regression. `scripts/verify_go.sh` is NOT a gate on this rig (charter row 76).
- Ground truths spot-verified by this planner: `ailang.v0.30.0` reproduces
  `Error: unknown command 'mission'`; the fleet `mission_pi_run.sh` still wires no `-e` flags
  (#1043 OPEN); `a036062` is an ancestor of base and carries the bench Authorization fix.

## §9 Doc defects found while planning (recorded; doc NOT edited)

- **DEF-1 (doc §3, candidate 1):** "Checkable against the binary and the index" is
  time-sensitive. The CURRENT pin (v0.41.0) HAS the `mission` group (planner-measured: it
  errors `failed to read mission registry missions`, not `unknown command`), so a naive live
  check against `$HOME/.pinned-ailang/ailang` would CONTRADICT the ground truth; and the index
  no longer shows the hole (`| 171 |` present — the manual repair). Resolution: the plan's Q1
  verification uses the surviving `$HOME/.pinned-ailang/ailang.v0.30.0` (planner-measured,
  reproduces) and `git log -S` over the index, never the current pin or the current index.
- **DEF-2 (row 78 / doc §3, candidate 2):** line drift — the fleet script's pi invocation sits
  at `:165` at plan time; row 78 recorded `:155` (iter-154). Property unchanged; issue #1043
  OPEN. Resolution: verification asserts the PROPERTY (no `-e`, no `PI_FENCE_ROOT`), never the
  line number.
- **DEF-3 (doc §4, mechanics gap):** the doc never says how a walker FINDS the incident
  object. Measured at plan time: `object get` is hashref-only and commits cannot register
  registry names, so the walk entry point is a `log range` scan. Resolution: §4/§6 encode the
  scan and seed finding F-1 (+proposed queue row). Not scope.
- **DEF-4 (doc predates row 39's land):** Phase A's commit path is now session-gated with a
  tty-fenced mint (`a036062`, iter-181); the doc's Phase A steps mention no credential setup.
  Resolution: M4 carries the pty-mint step and the PF-1/PF-2 pre-registrations.

## §10 Executor constraints (binding)

1. Work IN `/Users/voightkampff/dev/sunholo-data/.wt-world-iter182` on
   `sprint/w-prove-1-0-phase-a`. **NO git write operations** (no add/commit/stash/push/
   checkout/branch); the controller commits per milestone from `.snap/M<k>/`. Read-only
   `git log/show/grep` is REQUIRED for the verification steps and is not a mutation.
2. After EACH milestone: copy every created-or-modified file into `.snap/M<k>/` — cumulative
   full content (M4's snapshot includes M1–M3 files). `.snap/` and `.walk-scratch/` are
   untracked working artifacts, never committed; `.walk-scratch/` is EXCLUDED from snapshots.
3. Gates per milestone with `AILANG_BIN=$HOME/.pinned-ailang/ailang` exported: G1 always;
   G2 at M2; G3 at M3; G4 at M4. A gate red at BASE is a finding to REPORT, never fix.
4. Never touch `tools/launchd/*` (frozen fleet core), `~/.ailang/state/mission-v1*`, or any
   `.ail` file (verify_ail's outcome must stay byte-identical).
5. The verdict closed set is exhaustive; honest failure (FAIL-TIME / UNVERIFIED /
   UNANSWERABLE / NOT-RUN(TRANSPORT)) is a RESULT — record it, never substitute an easier
   question (doc §4).
6. Where the walk needs a feature World lacks, record a FINDING (proposed queue row) — do not
   build provenance tooling (doc §9).
7. Executor STOPS after `.snap/M3/` for the mutation battery, resumes for M4 only after the
   controller confirms the battery ran; final STOP after `.snap/M4/`. The mutation battery and
   all commits are controller acts.
