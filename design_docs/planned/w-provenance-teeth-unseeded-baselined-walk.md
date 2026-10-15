# w-provenance-teeth-unseeded-baselined-walk — clause-5 comparison arm: UNSEEDED question, live-timed pre-World baseline, read-only walk

**Status**: r4 — premise-evidence correction under D-WORLD-52=A (attended, 2026-10-02); fresh full author-excluding quorum REQUIRED. No walk, baseline or credit is claimed here.
**Target**: World 1.0 bar, clause 5 (`world-mission.md:1032`), queue row **114** (`:1933`)
**Priority**: P0 — critical path `… → 93 → 114 → 139` (D-WORLD-58, `:944`); row 134 LANDED `f3a90ee` (`:6105`); row 93 FINAL HOLDS 2026-10-15 (§16 V2)
**Estimated**: ~1 day per question, opportunistic (row text `:1933`)
**Designer-Lane (r4 author)**: `claude:claude-opus-5-5` (Anthropic). Prior authors: r3 OpenAI `gpt-6.1-sol`, r2 `pi:ollama/kimi-k3:cloud`. The quorum must exclude Anthropic.
**Planner-Lane**: codex-ok — protocol-only; no compiler, kernel, package or harness code.
**Base**: `origin/dev` `0a1e21c283e6df65760583df90f8540f81e8c6be` (HEAD == origin/dev, V1).
**Scope**: the protocol, not new code (`:1933`). Raw evidence: `design_docs/verification/world-iter270/designer/evidence.md`.

## r4 changes (r3 quorum BLOCKED 3/3, `verification/world-iter222/quorum-disposition.json`)

| r3 objection | r4 fix |
|---|---|
| P4 cites V5, which §16 lacked (V8 also absent); P4's row-text claim unverified (gemini, oc-glm) | §16 rebuilt V1–V11, contiguous, every row re-measured this session with command + rc + verbatim output in the evidence file. P4 → V4 (the `29e1336` fix and "test-budget leak" grep'd on `:1933`). No dangling ids (§16 self-check). |
| P8 (row-92 F-5 caveats) unmapped (oc-glm) | P8 → V3 (F-5 text at artifact `:181`) and V8 (caveat (b) re-measured live: registry-key refusal; plus a NEW caveat — a HEAD-built daemon refuses an old store at registry bootstrap). |
| Confinement facility and capture corpus UNKNOWN ⇒ every attempt structurally NOT-RUN (oc-kimi) | D-WORLD-52=A settles both. **B confinement** = the existing `ailang_only` lane in read-only (question-kind) mode, cited file:line and controlled live (V5). **W channel** = existing Claude Code 2.1.300 tool allowlist on a served COPY, controlled live (V6). **Corpus** = three row-93 FINAL stores, positively DECLARED with sha256 (V7), served read-only from copies (V8), and shown to hold an unseeded answer to one real question (V9). |

---

## §1 Problem

Clause 5 (`:1032`): on ≥3 REAL "why did X happen" questions, a provenance walk yields the verified answer in ≤5 minutes each, where the pre-World method was grep/log archaeology. Row 92 walked four questions in 10–20 s, but read diagnoses the sprint had seeded, and its baseline column says `UNMEASURED` (10 occurrences, V3). Row 114 (`:1933`) asks for the missing comparison arm: a live-timed pre-World baseline by an agent without World access, and a walk over a World state where the answer was NOT hand-recorded.

**Premises** (each verified this session; ids → §16):

| # | Premise | V |
|---|---|---|
| P1 | Row 92's walks were seeded retrieval; all its baselines are `UNMEASURED`, so no speedup can be claimed from them | V3 |
| P2 | D-WORLD-52=A (attended): B runs in the existing `ailang_only` lane, read-only, one root (repo + logs), no shell, no World; corpus = real operation recorded by World after row 134 (M6 smoke, row 93's runs, repo work); no invented questions | V2 |
| P3 | Row 134 LANDED; row 93's FINAL HOLDS (2026-10-15); D-WORLD-58 puts 114 next after 93 | V2 |
| P4 | Row 114's own text carries the C1 diagnosis (fix `29e1336`, "test-budget leak"), so any root containing the charter is contaminated for C1 | V4 |
| P5 | The lane's read-only mode is a tool LIST, not a policy: a question-kind task on pi gets `{Read, AilangRead, AilangCheck}` ∩ profile; `ailang_only` has no Bash; the AilangRead op is root-anchored and refuses `..`, absolute and symlink escapes; a directory read lists entries; pi has no grep tool | V5 |
| P6 | Claude Code 2.1.300 `--tools Bash --allowedTools "Bash(<prefix> verb *)" --permission-mode dontAsk --setting-sources ""` admits only the listed client verbs; `commit`, `cat`, `ls` and a `&&`-chained command are denied, and the denial is recorded mechanically in `permission_denials` | V6 |
| P7 | Five World stores from real row-93 operation survive on the rig, none open by any process; copies hash-identical to originals | V7 |
| P8 | Row-92 F-5 caveats hold: `--addr` needs a scheme; an exported `AILANG_REGISTRY_API_KEY` refuses every verb (re-measured); single-writer lock. NEW: a daemon built from HEAD refuses an old store at registry bootstrap — the rig's own `bin/ailang-worldd` serves it | V3, V8 |
| P9 | A served copy answers `head`, `provenance` and `why <index>` (all 7 links verified), and serving it leaves the archival copy byte-identical | V8 |
| P10 | The 2026-10-14 store holds, as operation-produced effect records, each codex World episode's `Workspace.Write` of `app_config.json` alongside a solution program that only reads it — and nobody hand-recorded a diagnosis in the store | V9 |
| P11 | The 2026-10-14 diagnosis entered the repo in `c21f303` (the verdict commit); the rig's checkout is the earlier `fd321b6`; the rig's logs carry no diagnosis phrase | V10 |
| P12 | World records carry no agent token counts (the 2026-10-09 overhead mechanism) | V11 |

## §2 Anti-vacuity words

| Word | Forbids here |
|---|---|
| **UNSEEDED** | Any object written for the walk; W reads only a DECLARED cohort produced by operation, frozen before the question existed |
| **BASELINED** | A ≤300 s verdict without that question's own live-timed B duration beside it |
| **CONFINED** | An isolation claim without the PR controls (allow succeeds, deny is recorded) run before the question is revealed |

## §3 Decisions

| Decision | Chosen by | Status |
|---|---|---|
| D1 Active harvest: real operation's next "why" question is taken; no invented questions | human (`:1933`; D-52) | frozen |
| D2 Three fresh roles: B (baseline), W (walker), V (verifier, one per arm); none is the diagnosing controller or this designer | designer | proposal |
| D3a B confinement = `ailang_only` read-only lane, one root | **human** (D-52, V2) | frozen |
| D3b W channel = Claude Code tool allowlist (client GET verbs only) + served COPY + head-invariant audit | designer (measured, V6/V8) | proposal |
| D4 B's root is a `.git`-free export at a pre-diagnosis commit plus allowlisted raw logs | designer | proposal |
| D5 Clause-5 ≤300 s compares ONLY the World envelope (issue → verified); baseline wall recorded, no threshold | designer | proposal |
| D6 Zero code, zero harness edits | human (`:1933`) | frozen |

**Freeze.** D2/D3b/D4/D5 are freeze items; D-52A says execution requires an explicit freeze. The controller records the freeze after this quorum passes; whether Mark ratifies it personally is a controller/attended question (UNKNOWN, §17).

## §4 Declared corpus (E3) — the capture cohorts

| Cohort | Store (original, never opened rw) | sha256 (archival copy == original, V7) | Producer | Completeness scope |
|---|---|---|---|---|
| K1 | `~/.ailang/state/floor-final-2026-10-14/store/world.db` | `bf1d4db5…a5d6d27` | row-93 FINAL 2026-10-14, 6 episodes `fl-{cc,cx}-r{1,2,3}` | every MCP/effect invocation of those episodes, entries 1–495 |
| K2 | `~/.ailang/state/floor-final-2026-10-15/store/world.db` | `33b63a55…6fd366` | row-93 FINAL 2026-10-15 | same, 495 entries |
| K3 | `~/.ailang/state/floor-final/store/world.db` | `e36c7fe2…aa3c70ce` | row-93 FINAL 2026-10-09 | same, 489 entries |

`floor-smoke-2026-10-15` and `floor-tune` exist (V7) but are not declared: no question from them has been harvested. A new cohort (e.g. Daneel `.ail` work through World, D-58) is declared the same way before its question is revealed.

**What the cohorts record:** World-side effects (tool inputs, effect plans, requests, outcomes) per episode. **What they do not record:** agent token usage (V11), grader verdicts, or anything outside World's tool calls. A question whose answer lives only there is `UNANSWERABLE(evidence)` for the walk — the finding the row allows (`:1933`), and it names the capture gap.

## §5 Protocol

**P0 Candidate.** At Gate 2 the controller writes one harvest-log line (question, event ref, ground-truth source, cohort) and opens the attempt's exposure log. Normal diagnosis recording is never paused; isolation comes from the roots, not from what the record says.

**P1 Eligibility.**

| Leg | Test |
|---|---|
| E1 REAL | Arose from operation; event ref exists; row-92 harvest rules apply (taken + rejected lists) |
| E2 GROUND TRUTH | A check neither arm produced (grader output, re-execution, CI log, git) — the same for both verifiers |
| E3 CORPUS | Cohort declared per §4 before reveal; none ⇒ walk `NOT-RUN(NO-CORPUS)`, baseline still runs |
| E4 SOURCE ISOLATION | B's root passes the D4 export check (P3); W's environment holds no repo, docs or memory |
| E5 QUESTION FROZEN | Exact bytes + sha256, identical for B and W |
| E6 ROLES FRESH | B/W/V spawned fresh by the controller; neither the controller nor this designer is B or W for a question whose diagnosis they have seen; each role's initial context is captured and token-scanned |

**PR — confinement controls, before reveal, in the exact role environment.**
- **B** (D3a; V5): pi with `ailang pi tool-profile` flags for the question-kind read set (`--no-builtin-tools --tools ailang_read` [+ `ailang_check`]), `AILANG_AGENT_POLICY` = a restricted policy with `fs_sandbox` = B's export root, on the AILANG release pinned at run time (V5 ran a dev build). Controls: read of a root file ok; directory listing ok; `..`, absolute-path and symlink escapes refused; a request for any tool outside the set refused by pi. Writes into the disposable export are harmless and not used.
- **W** (D3b; V6): fresh `claude -p` (pinned client version), cwd = an empty scratch dir holding only a copy of the rig's `ailang-worldd`, `--tools Bash`, `--allowedTools` = `Bash(env -u AILANG_REGISTRY_API_KEY ./worldd-client --addr http://127.0.0.1:<port> <verb> *)` for verbs `head`, `why`, `object get`, `object find`, `log get`, `log range`, `provenance` only; `--permission-mode dontAsk --setting-sources "" --strict-mcp-config --output-format stream-json`. Controls: `head` succeeds; `commit`, `call`, `cat <db>`, `ls <repo>`, a `&&`-chained command and a `curl -X POST` each appear in `permission_denials`. No session token is placed anywhere W can read.
- Any failed or missing control ⇒ `NOT-RUN(ISOLATION)` before reveal, preserved.

**P2 Freeze.** Before either arm: repo HEAD; question sha256; cohort archival sha256; served-copy `head` ref and `log range` count (the invariants); B root manifest (paths + sha256) and its token scan; spawn commands; `max_attempt_wall_s` = 3600. The served copy is prepared and started BEFORE reveal: `env -u AILANG_REGISTRY_API_KEY <rig>/bin/ailang-worldd serve --db <copy> --bind 127.0.0.1:<port> --ailang-bin ~/.pinned-ailang/ailang` (V8), never on an original.

**P3 B's root (D4).** One directory: `git archive <pre-diagnosis commit>` (no `.git`) plus an allowlist of raw operation logs from the rig (phase logs, agent transcripts, fixtures, raw run rows). Excluded whatever its ancestry: World stores, `private/`, verdict prose (`README.md`, `summary.md`), later charter/log text, memory files. The controller token-scans the finished root for the question's diagnosis phrases with a positive control (the same scan on the diagnosis commit must hit). Zero hits plus a manifest review is required; a hit ⇒ fix the export or leave the candidate `CONTAMINATED-PENDING`.

**P4 Baseline FIRST.** Question issued → B answers from its root, transcript timestamped → V_B verifies against E2 → sealed. Nothing from B reaches W.
**P5 Walk SECOND.** Question issued → W answers through the allowlisted verbs only, every call logged → V_W verifies against the same E2 → sealed.
**P6 Accounting.** Archival sha256 unchanged; served-copy `head` and entry count equal the P2 invariants (any change ⇒ `INVALID(seeding)`); `permission_denials` and the tool-call log quoted; question sha256 equal across arms.

### Timer rules

- Answer duration per arm: question issued → answer declared. Serve start and environment build are before reveal and recorded separately.
- **World envelope** = W question issued → V_W verified result, continuous wall including handoff and verifier setup. **Only this is compared with 300 s.** Baseline envelope (issued → V_B verified) is recorded without a threshold.
- Timestamps come from the controller's clock and the stream-json events; declared vs transcript durations agree within 5 s; no retrospective timing.
- Over `max_attempt_wall_s` ⇒ `INVALID(clock)`.

## §6 Verdicts and integrity audit

Per arm, closed set: `PASS` (verified; for W also envelope ≤300 s) · `FAIL-TIME` · `UNVERIFIED` · `UNANSWERABLE(evidence)` (declared cohort queried, answer absent) · `NOT-RUN(ISOLATION|NO-CORPUS|TRANSPORT)` · `INVALID(clock|contamination|seeding|pin)`. `CONTAMINATED-PENDING` is a candidate state, not a verdict. Replays and row-92's walks are never qualifying.

V attests per question: (1) archival hash pre == post; (2) served `head` and entry count unchanged; (3) PR controls quoted for both arms before reveal; (4) exposure log and initial-context scans clean; (5) no W tool call outside the allowlist succeeded.

## §7 Artifact (`design_docs/verification/w-provenance-teeth-unseeded-baselined-walk/`)

`index.md` (harvest log with rejections, accounting, comparison table, findings) and per question `Q<n>-<slug>/`: `pins`, `prereq-transcript`, `baseline-transcript`, `walk-transcript` (+ tool-call log, `permission_denials`), `verify-b`, `verify-w`, `envelope` (`t_issue_b`, `t_verified_b`, `t_issue_w`, `t_verified_w`, `world_wall_s`, `baseline_wall_s`), `exposure-log`, `verdict`. Transcripts stay on the rig if large; the artifact carries their sha256.

## §8 Candidates and completion

Row-complete at **≥3 qualifying successes**: E1-real; PR passed; identical question bytes; audit clean; same-question live baseline; W answer verified by an independent V against the same E2; World envelope ≤300 s. Partials never count; every attempt is preserved.

| # | Candidate | Cohort | Status |
|---|---|---|---|
| Q1 | "Why did codex's World arm fail `config_file_parser` 0/3 in the 2026-10-14 FINAL, while its shell arm passed 2/3?" | K1 | **ELIGIBLE pending PR.** Answer reachable from operation-produced effect records (V9). Pre-diagnosis commit `fd321b6` (V10); diagnosis text lives in `c21f303` README, the 2026-10-15 README, the 10-15 preregistration and later charter/log — all outside B's root. E2: re-run a World-arm solution from K1 in a clean dir without `app_config.json` and compare with the grader rows. B's root includes the rig transcripts, which hold the same raw writes (V10) — fair: they are the pre-World logs. |
| Q2 | "Why was codex's World overhead +26.6% in the 2026-10-09 FINAL?" | K3 | **CANDIDATE.** Diagnosis is in D-WORLD-70 (`:956`) and the 10-09 README; its pre-diagnosis commit is not yet measured (PENDING at harvest). Expected walk outcome: `UNANSWERABLE(evidence)` for the token mechanism (V11) — a capture-gap finding, not a success. |
| Q3+ | Next real question from K2, a smoke/tune store once declared, or World-recorded repo work (D-58) | — | none harvested |
| C1/C2/C3 | iter-206/207/210 CI and suite incidents | none | **NOT-RUN(NO-CORPUS) by construction** — they predate row 134 and no World store recorded them. Kept in the rejection log; C1 is also contaminated by `:1933` (V4). |

Honest expectation: Q1 can qualify; Q2 probably yields the gap finding; a third success needs a further real question. The row stays open until then — no invented questions.

## §9 Acceptance criteria

| AC | Criterion |
|---|---|
| AC1 | Harvest log with taken and rejected candidates (row-92 rules) |
| AC2 | Question sha256 identical for B and W |
| AC3 | B root: no `.git`, no World store, no verdict prose; manifest + token scan 0 hits with a positive control |
| AC4 | PR controls for both arms quoted with timestamps BEFORE reveal (B: allow/escape refusals; W: `head` ok + six denials in `permission_denials`) |
| AC5 | Archival sha256 pre == post; served `head` and entry count unchanged |
| AC6 | ≤300 s claims cite only the World envelope; baseline envelope recorded; ≤5 s cross-check |
| AC7 | Both answers verified against the same E2 by fresh verifiers; B sealed before W starts |
| AC8 | Verdicts from the §6 set; failures preserved |
| AC9 | ≥3 qualifying successes; partials never count |
| AC10 | No speedup ratio without that question's live baseline; every ratio states the D-52 baseline conditions (read/list only, no grep — V5) |
| AC11 | Zero code and harness diffs (D6) |
| AC12 | Stores used only as copies; no daemon on an original; writer locks untouched |

## §10 Adversarial controls

| MU | Mutation | Expected kill |
|---|---|---|
| MU-1 | W commits an object to the served copy | allowlist denial (AC4); else `head` changes ⇒ `INVALID(seeding)` |
| MU-2 | Diagnosis file left in B's root | AC3 scan ⇒ `INVALID(contamination)` |
| MU-3 | W chains a shell command onto an allowed verb | `permission_denials` (V6) |
| MU-4 | W inherits project memory or CLAUDE.md | empty cwd + `--setting-sources ""`; initial-context scan ⇒ `INVALID(contamination)` |
| MU-5 | Question bytes differ between arms | sha256 ⇒ `INVALID(contamination)` |
| MU-6 | Clock starts at first successful query | stream-json timestamps ⇒ `INVALID(clock)` |
| MU-7 | Daemon started on an original store | AC12; original sha256 re-checked ⇒ `INVALID(pin)` |
| MU-8 | An easier question swapped in after a failure | immutable candidate ids; rejection logged |

## §11 Conflict surface and non-goals

No parser, typechecker, kernel or package surface. Read-only use of the daemon's GET verbs, the rig's binaries and stores (copies). Non-goals: no capture recorder, no confinement code, no filter or proxy, no AILANG-language claims, no change to the floor harness, no synthetic dry-run counted as evidence. A capture gap (e.g. token usage, V11) may file ONE optional scoped row.

## §12 Risks

| Risk | Mitigation |
|---|---|
| Baseline without grep is slower than a real grep session, inflating any World advantage | AC10: ratios state the D-52 conditions; the bias is named in `index.md`; the bar is ≤300 s for W, not a ratio |
| A fresh W cannot finish in 300 s through `why` one entry at a time | `provenance --episode` narrows the range first; `FAIL-TIME` is recorded honestly; loosening is attended-only |
| A transcript in B's root carries a diagnosis written during operation | P3 token scan + manifest review; `CONTAMINATED-PENDING` if it cannot be removed without destroying the evidence |
| Rig stores deleted before execution | archival copies under `~/.ailang/state/world-iter270/` (sha256 in §4) |

## §13 Axiom compliance

A1 +1 (pins and transcripts make verdicts re-derivable) · A2 +1 (frozen archival copies) · A3 +1 (operation-produced effect records only) · A4 +1 (authority narrowed by tested allowlists, not credential absence) · A5 +1 (bounded envelope) · A7 +1 (closed verdicts, fixed schema) · A9 +1 (durations recorded, none imputed) · A11 +1 (NOT-RUN/INVALID preserved) · A12 +1 (verification outside the answering arm) · others 0. **Net +9.**

## §14 References

Row 114 `:1933`; clause 5 `:1032`; D-WORLD-52 `:938`; D-WORLD-58 `:944`; row 134 `:6105`; row-92 artifact `verification/w-1-0-value-demonstration.md` (Harvest log `:9`, F-5 `:181`); r3 quorum `verification/world-iter222/`; SE-domain lane measurements `planned/w-software-engineering-domain.md` V33–V36, V52, V54; FINAL records `verification/world-floor-final-2026-10-{09,14,15}/`.

## §16 Verification Log — r4 (2026-10-15, iteration 270)

Commands, exit codes and verbatim outputs: `design_docs/verification/world-iter270/designer/evidence.md`, one section per row. Nothing below is copied from r3.

| ID | Measurement |
|---|---|
| V1 | `git rev-parse HEAD origin/dev` → both `0a1e21c2…`. |
| V2 | `grep -n` on the charter: row 114 `:1933`; clause 5 `:1032`; D-52 `:938` ANSWERED A with the lane and corpus text quoted; D-58 `:944` order `… 93 → 114 …`; row 134 `:6105` LANDED `f3a90ee`. 2026-10-15 README `:15` "Overall: HOLDS". |
| V3 | Row-92 artifact: `grep -c UNMEASURED` → 10; `grep -n F-5` → `:181` (caveats a–c quoted); `## Harvest log` `:9`, rejections `:40`. |
| V4 | `sed -n 1933p` charter → contains `` 29e1336` ``; `grep -c 'test-budget leak'` → 1. |
| V5 | ailang @ `514ca3bc3`: `toolpolicy.go:40,53-54,71`, `provider_executor.go:141-152,223-231`, `policytool.go:144-154`, `fs_ops.go:41-42`. `ailang pi tool-profile AilangRead` rc 0 → `--no-builtin-tools --tools ailang_read`; `… Grep` → error "pi has no tool for". `policy-tool` live: in-root read ok, dir list ok, `..`/absolute/symlink refused "escapes sandbox", `run` refused "unknown op", write ok (so read-only = tool list). Binary v0.54.10-dirty — re-run on the pinned release at PR. |
| V6 | `claude` 2.1.300, haiku, allowlist of client verbs: `head` ran (`sha256:1710ddd7…`); `commit`, `cat ../f1014.db`, `ls <repo>/design_docs` denied; `--output-format json` → `permission_denials` = the `why 120 && cat /etc/hosts` command. Without `env -u` in the prefix the allowed verb exits 2 on the registry-key refusal. |
| V7 | `pgrep`/`lsof` on the 10-14 store: none. Five stores + writer locks listed with sizes. `cp -p` + `shasum -a 256`: original == copy for K1/K2/K3. K1 table counts: log_entries 495, objects 4446, journal 1976, commit_objects 1977; semantic-id histogram and episode prefixes `fl-{cc,cx}-r{1,2,3}`. K2 495, K3 489 entries. |
| V8 | HEAD-built `ailang-worldd --help` lists `why`, `provenance`, `object`, `log`. Serve on a copy: with the key → refusal; HEAD binary → "registry: existing head … diverges from bootstrap revision"; rig binary + `--ailang-bin` pinned v0.41.0 → `health` ok, `head` `sha256:1710ddd7…`, `provenance` per-episode ranges (`fl-cx-r1 entries=81-164`), `why 120` "all 7 link(s) verified". After stop: K1 archival sha256 unchanged. |
| V9 | `sqlite3 -readonly` on the K1 copy: payloads matching `app_config` 25 (control `fizzbuzz` 26); 3 `Workspace.Write` requests of `app_config.json`, each with a same-second solution write that imports only reads; their effect-intents are `effect:fl-cx-r1:18`, `fl-cx-r2:16`, `fl-cx-r3:17`, records `allowed:true`. Designer feasibility probe only. |
| V10 | `git log` on the 10-14 dir: `fd321b6` → `4f15e17` → `c21f303` (diagnosis + verdict); `merge-base --is-ancestor 4f15e17 c21f303` rc 0; diagnosis-phrase count 0 at `4f15e17` charter, 1 in `c21f303` README (control). Rig repo HEAD `fd321b6`. Rig grep for diagnosis phrases (excluding repo): no files; control `config_file_parser` → `phase1.log`, `phase2.log`; transcripts hold raw `app_config.json` writes. |
| V11 | K3 copy: payloads with `input_tokens`/`inputTokens` 0; control `fizzbuzz` 28. |

**Self-check:** premises cite V1–V5 and V7–V11; D3b and PR cite V6; no other ids are used.

## §17 Still UNKNOWN

- A live pi run of B under the read-only flags (only the flags and the policy-tool layer were measured, on a dev build).
- Whether D2/D3b/D4/D5 freeze on quorum or need Mark personally (D-52A says "explicit freeze").
- Q2's pre-diagnosis commit and E2 source; whether any third real question will arise.
- A fresh W's real envelope. The designer has seen the Q1 answer (V9) and is barred from B, W and V for Q1.

**Verdict on readiness:** executable for Q1 once frozen. The confinement facilities are named and controlled live, the cohort is declared with hashes, and an unseeded answer exists in it. Row completion still depends on real questions the protocol cannot create.
