# 1.0 value demonstration (w-prove-1-0-phase-a)

Sprint `w-prove-1-0-phase-a` · queue row **92** (clause 5) · design doc
[`w-prove-the-1-0-bar.md`](../planned/w-prove-the-1-0-bar.md) · Phase B (row 93) is OUT of
scope. This artifact is the deliverable (doc §4.5, §9). It records the harvest, baselines,
walk protocol, and — after M4 runs — the timing, verification, and comparison tables. M1–M3
seed everything the walk needs; M4 is a separate run by the controller after the mutation battery.

## Harvest log

Questions are pre-selected from the mission's own recorded operation (charter + mission log),
NOT invented to fit the walk. Every taken question arose from real operation, has independent
ground truth the walk did not produce, and its original method/baseline is recoverable from the
record. Candidates considered and rejected are recorded below with reasons — a harvest with no
rejections would be evidence the bar was set after the fact (doc §4.1).

### Taken (mandatory 3 + C4, eligibility-gated)

| # | Slug (semanticId) | Question | Ground truth (independent of the walk) | Record |
|---|---|---|---|---|
| Q1 | `iter171-index-row` (`world/mission/incident/iter171-index-row`) | Why did iteration 171 have a full log entry and no index row? | The pinned **v0.30.0** binary lacks the `mission` command group: `ailang.v0.30.0 mission rotate-log world --keep 31` prints `Error: unknown command 'mission'`. Planner-verified 2026-09-24: `$HOME/.pinned-ailang/ailang.v0.30.0` survives on the rig and reproduces it. The index hole is visible in **git history** of `world-mission-index.md` (the row exists today because it was added manually after the iter-172 rule — `grep -c '^| 171 |'` is 1 NOW; that is the repair, not the incident). | charter row 86 + the `mission`-command bullet (`world-mission.md` ~lines 82–101); log `## 172` |
| Q2 | `iter154-unfenced-pi` (`world/mission/incident/iter154-unfenced-pi`) | Why did three designer runs execute unfenced while the rulebook said they were sandboxed? | The fleet's `scripts/mission_pi_run.sh` invokes `pi --mode json --no-session --model "$MODEL" < "$DIRECTIVE"` with **no `-e` flag and no `PI_FENCE_ROOT` anywhere in the file** — planner-verified live 2026-09-24 via `gh api` (invocation at :165 today; row 78 recorded :155 at iter-154 — line drift, property unchanged); issue `sunholo-data/ailang#1043` is **OPEN**. | charter row 78; log `## 154` containment note |
| Q3 | `iter155-faithfulness-proof` (`world/mission/incident/iter155-faithfulness-proof`) | Why did a green `shasum -c` faithfulness proof pass over a destroyed commit split? | The prescribed manifest covers the **final tree**, identical whether the split is correct or collapsed; `git add design_docs host` staged disk state, MS1 swallowed MS3's files, MS3 committed nothing, and `shasum -c` returned OK on every file. Corroborated by two in-git records (row 81 + log `## 155`) and by the rebuilt row-59 split in git history. | charter row 81; log `## 155` |
| Q4/C4 | `iter181-ci-bench-401` (`world/mission/incident/iter181-ci-bench-401`) | Why did PR #141's remote CI go red on `BenchmarkRESTCommit` when every local gate was green? | `a036062` is an ancestor of base and its diff adds `req.Header.Set("Authorization", auth)` to `bench_test.go` — the bench predated the session gate and 401'd. Ground truth in-git (STATUS 2026-09-24 / iter-181). **TAKEN** — all three C4 eligibility-check legs passed (see below). | STATUS 2026-09-24 (iter-181), `world-mission.md`; git `a036062` |

### C4 eligibility check (run 2026-09-24 in M1, all legs must hold)

```bash
$ git merge-base --is-ancestor a036062 HEAD && echo ELIGIBLE-1
ELIGIBLE-1
$ git show a036062 -- host/daemon/bench_test.go | grep -c 'req.Header.Set("Authorization", auth)'
1
$ grep -c "bench predates the sprint" design_docs/world-mission.md
1
```

All three legs pass → **C4 TAKEN**. The fourth incident slug `iter181-ci-bench-401` is
therefore taken and appears in the taken table above.

### Considered and rejected (a harvest with no rejections is evidence the bar was set after the fact)

| # | Candidate | Decision | Reason |
|---|---|---|---|
| C5 | Why did iteration 180 refuse every push/PR (prepush-gate unpassable, row 95)? | **REJECT** | Ground truth lives in a rig-global pi extension (`~/.pi/agent/extensions/prepush-gate.ts`) whose defective bytes have since been **replaced** by the attended fix (ailang `754fc7a43`); the original bytes no longer exist anywhere the executor may read, and the fix commit is in a different repo. The only surviving record is the charter text itself — the walk would be citing its own source (§2 circularity risk). |
| C6 | Why did a mutation-drill restore destroy the executor's uncommitted work (row 82)? | **REJECT** | The destroyed work was uncommitted by definition; no independent artifact of it survives anywhere. The row text is the only record, so verification would reduce to "the walk repeated the row" — exactly the circular shape §2's `verified` forbids. |

No further candidate was spotted in the record during M1; had one appeared, it would be
recorded here with take/reject and the reason.

## Baselines

For each taken question the pre-World method is as enumerated in the record (file + section
citation), and the measured cost is quoted from the record **or** the literal string
`UNMEASURED (record states method, not cost)` — never imputed (doc §4.2, §8). An iteration's
whole-fire metered cost is NOT the diagnosis cost and is never substituted.

| Q | Method citation | Cost in the record |
|---|---|---|
| Q1 | `world-mission.md` row-86 bullet + the v0.30.0 `mission`-command bullet (~lines 82–101): grep the index for the row, control-grep a neighbour, run the prescribed regenerator, positive-control another command group | `UNMEASURED (record states method, not cost)` |
| Q2 | `world-mission.md` row 78; log `## 154` containment note: read the recipe's two-flag invocation block, read the runner's actual invocation, compare | `UNMEASURED (record states method, not cost)` |
| Q3 | `world-mission.md` row 81; log `## 155`: trust three green `git commit`s, get a green `shasum -c`, then catch it only by reading `git show --stat` per commit; rebuild; hit the zsh word-splitting trap on retry | `UNMEASURED (record states method, not cost)` |
| C4 | `world-mission.md` STATUS 2026-09-24 (iter-181): local gates all green, remote CI red, read the CI failure, reproduce, judge round 2 | `UNMEASURED (record states method, not cost)` |

All four pre-World baselines are **unmeasured**: the mission record (rows 86/78/81, log `## 154`
and `## 155`, STATUS iter-181) states what went wrong and how it was found, but never a duration
for the diagnosis itself. Per doc §4.2 and §8 ("Not verified, and named as such"), they are
reported as unmeasured rather than estimated. The walk-time comparison therefore uses the
standard configured bar: **≤300 s per question** (plan §1 verdict closed set). M4 will record the
measured walk durations against this.

## Walk protocol

Copied from the sprint plan §4 (walk mechanics) and §1 (verdict closed set). The walk runs in M4
against a scratch DB (`.walk-scratch/`, untracked) with a non-default port; M1–M3 only build the
fixtures and the durable replay test.

- **Commit path (M4 setup):** `POST /v1/commit` is session-gated; the walker mints a session via
  the tty-fenced `session mint` command (through a pty wrapper — the fence is intended, row 39
  D1) and POSTs the fixtures in entry-index order. All GET routes are open.
- **Walk entry point (F-1):** `GET /v1/objects/{ref}` is hashref-only, and `POST /v1/commit`
  cannot register `epoch_registry_heads` names. So the walk locates an incident by scanning
  `log range --from 0`, fetching each `transitionRef` object's metadata until the semanticId
  matches `world/mission/incident/<slug>`. Fine at 3–4 entries; a genuine provenance-query gap
  carried as finding F-1.
- **Read the incident payload** (base64) → `question`, `answer`, `sources[]`.
- **Follow every source hash** the same way → evidence payloads.
- **Assemble the answer** from the incident object's `answer` field.
- **Verify outside World** against each question's independent ground truth (the deployed
  v0.30.0 binary, fleet script via `gh api`, this repo's git history, the in-git charter/log
  text). An answer produced by the walk alone is recorded **UNVERIFIED** and does not count
  (doc §2).
- **The timer starts at the question, before the first query** — never at the first successful
  query (doc §4.3).
- **Verdict closed set:** `PASS` (verified AND ≤300 s) · `FAIL-TIME` (verified, >300 s) ·
  `UNVERIFIED` (no independent source confirmed) · `UNANSWERABLE` (the walk cannot produce the
  answer from the committed record) · `NOT-RUN(TRANSPORT)` (the sandbox refused the listener or
  the pty mint — labelled, never reported as pass/fail). Honest failure is a result; a
  failed/unanswerable walk is recorded, never substituted with an easier question (doc §4).
- **Transport blocks** (pre-registered PF-1/PF-2/PF-3): if `serve` cannot bind, or neither the
  `pty.spawn` wrapper nor the `script -q /dev/null` alternate satisfies the tty fence, every
  live-walk row is `NOT-RUN(TRANSPORT)`; the M3 durability test stands as the regression guard,
  and the transport block is reported as a finding.
- **Replay-pin:** `serve --ailang-bin` archives the interpreter by hashing the raw file bytes;
  health's `interpreter_ref` must equal the value baked into the fixtures
  (`"sha256:" + shasum -a 256 "$HOME/.pinned-ailang/ailang"`). M4 asserts REPLAY-PIN-OK before
  the timed walk begins.

## Timing log

Populated by M4 (per-question start_utc / stop_utc / duration_s / verdict).

| question | start_utc | stop_utc | duration_s | verdict |

## Verification

Populated by M4 (per-question answer-from-walk / verification source command verbatim /
confirmed?).

| question | answer (from walk) | verification source (command, verbatim) | confirmed? |

## Comparison

Populated by M4 (six prescribed columns; the artifact IS the deliverable, doc §4.5).

| question | baseline method (citation) | baseline cost | walk time (s) | verification source | verdict |

## Findings

Seeded from the sprint plan §6 (verbatim). No code was written to work around these — they are
findings, not scope (doc §9).

- **F-1 — no semantic-ID query route.** `GET /v1/objects/{ref}` is hashref-only; a walk must
  linear-scan `log range` + per-entry `object get` metadata to locate an incident by
  semanticId. At mission scale (thousands of entries) this breaks the ≤5-minute bar's spirit.
  Proposed queue row: `w-object-lookup-by-semantic-id` (read route `GET /v1/objects/by-semantic-id/{name...}`
  or a log index over semanticIds). **Finding, not scope** (doc §9).
- **F-2 — registry names are not commit-settable.** `epoch_registry_heads` is written only by
  `registry.Bootstrap`/broker/transitionreg, so a walker cannot mint a well-known registry
  name for an incident class via `POST /v1/commit`. Compounds F-1; same proposed row covers it.

Transport findings (PF-1/PF-2/PF-3), if fired during M4's live walk, are recorded here as they
occur.

## Walk transcript pointer

The raw, unedited command output of the M4 live timed walk is appended per question to
`w-1-0-value-demonstration/walk-transcript.md` (sibling of this file, in the same verification
directory). It is written by M4 — the controller's live-walk run after the mutation battery.