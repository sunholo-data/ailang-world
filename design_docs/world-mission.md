# Ailang World Mission — a semantic operating environment whose transaction language is AILANG

<!--
  Charter for the WORLD mission (from mission-charter-TEMPLATE.md, M-MISSION-PORTABILITY M3).
  DRAFTED AT BOOTSTRAP 2026-07-23 — written to be CHALLENGED, not to be final.
  Iteration 0 = ratify this charter (bar + queue + guardrails) through the design quorum
  with Mark, attended. The kill switch stays SET until ratification.
  Thesis document: DESIGN.md (v0.3) in this directory — this charter is its operational
  distillation; where they conflict, the ratified charter governs the loop and the
  conflict gets fixed in DESIGN.md.
-->

**Type**: Long-running mission (peer of the V1 mission in `sunholo-data/ailang`); advanced by a
scheduled outer loop on the always-on rig. **First mission outside the AILANG repo.**
**North star**: ship **Ailang World 1.0** — a local-first semantic operating environment
([DESIGN.md](DESIGN.md) §1: an immutable, typed, content-addressed world graph whose transactions
are AILANG programs; propose → verify → commit; capabilities, budgets, evidence, replay) — with
the software-engineering domain live — value proven by clause 5's provenance teeth + the R1
standing evidence, floored by clause 4's resident-agent non-inferiority check (motoko remains
the aspirational first *native* agent as an optional arm).
**Traces to**: [DESIGN.md](DESIGN.md) (thesis) and the AILANG program's operating model
([PROGRAM.md](https://github.com/sunholo-data/ailang/blob/dev/design_docs/PROGRAM.md) — frozen
minimal core, extension-routed evolution); language gaps found here route BACK to
`sunholo-data/ailang` as issues/backlog items, never worked around locally.
**Skill**: the SAME unforked `mission-control` skill every mission uses, reached via user-level
symlinks (`~/.claude/skills/*`) — this repo carries NO skill copies (a copy is a fork; forks stop
learning).
**Scheduling**: launchd `dev.ailang.mission-world` — **ARMED by Mark 2026-07-23** ahead of
ratification; sprint-routing stays BLOCKED until the charter is ratified (iter-0 honored this:
quorum-only, no sprint). Off switch: `~/.ailang/state/mission-world.disabled`. StartInterval
staggered vs the V1 loop (shared rig quota). Billing guard: subscription-or-nothing.
**Log**: [world-mission-log.md](world-mission-log.md) — append-only, one entry per iteration.
**Human-facing reporting**: GitHub issue #1 — every iteration posts its report there as a comment
(Mark follows by email); driver crashes post there too.

## Repo Profile (M-MISSION-PORTABILITY M2 — the per-mission values mission-control reads)

- **Repo slug**: `sunholo-data/ailang-world` (driver: `MISSION_REPO`)
- **Mission doc**: `design_docs/world-mission.md` (driver: `MISSION_DOC`)
- **Mission name / state namespace**: `world` (driver: `MISSION_NAME`; fully namespaced
  `~/.ailang/state/mission-world-*` paths — no collision with the V1 loop)
- **Bookkeeping issue**: `#1`, rotates weekly; live number in
  `~/.ailang/state/mission-world-gh-issue` (seeded by the v1 agent after ratification)
- **CI workflows Gate 3b / Gate 1 poll**: `CI` (job: "ailang-code verify gate")
- **Verify profile**: `ailang-code` — the shipped `ailang` binary IS the gate:
  `ailang check` (types), `ailang test` (tests, once test suites exist),
  `ailang ai-check` (unified check+verify: types + Z3 in one JSON — do not reinvent a split
  gate). Repo-wide sweep: `./scripts/verify_ail.sh` (what CI runs; fails loudly on zero modules).
  Go host code (daemon/store, from M1 onward) adds `go build ./... && go test ./...` to the gate
  when it first lands — extend the CI workflow in the same PR that introduces Go code.
- **Rig PATH (process fix, iter-18 — 3 instances in ONE iteration)**: the agent tool shell's `PATH`
  does **not** include `/opt/homebrew/bin`, so `gh`, `go`, `node` and `codex` all fail with
  `command not found` / `env: node: No such file or directory` (`rc=127`). Every Bash call that
  needs them must `export PATH=/opt/homebrew/bin:$PATH` first, and **every directive handed to a
  spawned sub-agent must say so** — otherwise a planner/executor reads it as a broken toolchain or
  a spent codex quota rather than a PATH gap (iter-18 lost a codex probe to exactly that
  misreading). The pinned AILANG binary lives OUTSIDE the repo at `/tmp/ailang-v0300/ailang`
  (v0.30.0, commit `e37b370`, clean) — always pass it explicitly as `AILANG_BIN`.
- **Codex model pins are NOT covered by the driver's probe (process fix, iter-19).** The driver's
  codex pre-flight (`tools/launchd/mission-control.sh:248`) runs `codex exec --skip-git-repo-check
  'reply with exactly: ok'` — **without `--model`** — so it exercises codex's DEFAULT model and
  reports the lane healthy even when the pinned model cannot run. Iter-19 hit exactly this:
  `MISSION_EXECUTOR_MODEL=codex:gpt-5.6-sol` failed with `400 … 'gpt-5.6-sol' requires a newer
  version of Codex` on codex-cli **0.137.0**, while a probe without `--model` returned rc=0 on
  default **`gpt-5.5`**. So: **always run the skill's own probe WITH `--model` before trusting the
  lane**, and read a model-availability 400 as a PIN problem, never as spent quota (that
  misdiagnosis is the iter-18 scar repeating in a new costume). The driver is FROZEN CORE — do not
  patch it here; route the fix upstream on both channels. Also guard the call site with
  `env -u OPENAI_API_KEY` so an ambient metered key cannot silently bill while `auth.json` is on
  ChatGPT subscription auth (`auth_mode=chatgpt`) — the `claude-sub` discipline, applied to OpenAI.
- **ENUMERATE the evaluator's non-blocking carry-forwards in the log entry (process fix, iter-19).**
  Iteration 18 recorded "three non-blocking carry-forwards → A2" as a bare COUNT; the items
  themselves were never written down anywhere durable and are unrecoverable — A2 could not fold
  what it could not read. Every Gate-4 entry now lists them **numbered, one line each, with the
  checkpoint they target** (iter-19: CF-1…CF-9), and the evaluator directive must ask for that list
  explicitly.
- **A fresh worktree NEVER contains the sprint plan JSON (process fix, iter-20).** `.gitignore`
  line 3 is `**/.ailang/`, so `.ailang/state/sprints/<item>.plan.json` — the executor's OWN plan —
  is structurally absent from every `git worktree add`. Iter-20's codex executor reported it as a
  mid-run blocker and fell back to the design doc (correct: the doc governs, so no harm — but the
  plan's `non_vacuity_requirements` would have been silently lost had the directive not quoted
  them). So every worktree-based executor directive must EITHER `cp` the plan into the worktree
  first (it stays gitignored, so it never pollutes the diff) OR quote the binding requirements
  inline. Do not assume a sub-agent can read mission state that lives outside git.
- **The weekly-rotation rule misfires on a thread created just before Monday 07:00 (process fix,
  iter-20 — 2nd instance).** Issue #9 was created `2026-07-27T05:51Z`, ~1h BEFORE that Monday's
  07:00 boundary, so the literal rule ("past the most recent Monday 07:00 AND the current issue was
  created before that boundary") reads ROTATE — which would open a second thread for the very week
  #9 already titles. Iter-19 and iter-20 both hit this and both judged NOT-DUE. **The intent test
  governs**: a thread whose TITLE names the current week IS this week's thread regardless of the
  clock; the `>80 comments` limb is the real size bound. Rotate on the boundary only when the
  current thread belongs to a PREVIOUS week.
- **`/tmp` scratch paths in the cross-provider recipe are SHARED between missions — namespace them
  (process fix, iter-21).** The mission-control skill's Gate-3 codex recipe names literal paths
  (`/tmp/codex_directive.txt`, `/tmp/codex_run.sh`, `/tmp/codex_last.txt`), and BOTH loops run on
  this one rig. Iter-21 went to write `/tmp/codex_run.sh` and found it holding the **V1 loop's**
  directive (a `design-subsumption` worktree), i.e. the two missions silently share one scratch
  buffer: a mid-flight collision could hand one mission's executor the other's directive, or
  truncate a running capture. The blast radius is real even though these are regenerated per use.
  So every `/tmp` artifact this loop writes gets a mission suffix — `codex_directive_worldb.txt`,
  `codex_run_worldb.sh`, `codex_last_worldb.txt`, `codex_out_worldb.log`. Treat the skill's literal
  paths as TEMPLATES, never as filenames. (The skill itself lives in the V1 checkout and World
  never edits it — this is proposed upstream, applied locally as a Repo Profile rule.)
- **A driver probe failing `rc=127` is a PATH gap — NEVER a spent quota and never an unusable model
  pin (process fix, iter-21; 4th instance of the PATH class, 3rd misdiagnosis costume).** Iter-21's
  fire exported `executor=opus` because the driver's codex pre-flight returned
  `rc=127 exec: codex: not found`: `tools/launchd/mission-control.sh:44` exports
  `PATH="$HOME/go/bin:$HOME/.local/bin:$PATH"` and omits `/opt/homebrew/bin`, so `claude` resolves
  and **`codex` cannot**. The skill's own probe WITH `--model` and PATH fixed returned rc=0 on
  codex-cli 0.145.0. Rule: **before honouring ANY driver-exported model fallback, re-probe** (the
  iter-19 rule, now load-bearing in the opposite direction — it catches false NEGATIVES too), and
  classify the rc: 127 = environment defect, fix or route it and honour the ratified pin; a 400 =
  pin problem; a `QUOTA_SIG` match = quota. The driver is FROZEN CORE → routed as
  `sunholo-data/ailang#493` (upstream `dev` has the identical line, so the V1 loop is demoting to
  opus every fire too), never patched locally.
- **STATUS rotation is a HAND edit, never a regex script (iter-18 scar).** A "move the 4th stamp to
  the archive" script that bounded the stamp with the next `---` deleted **293 lines** of this
  charter (bar, Conflict Surface, guardrails, routing policy, queue, Premise Verification Log) —
  the stamps are not `---`-delimited. Caught by `git diff --stat` before commit and restored with
  `git checkout --`. Always `git diff --stat` the charter before committing it.

---

## STATUS (rotation rule)

Newest **3** STATUS stamps live here; older ones move to `world-mission-status-archive.md`.
At Gate 4, after adding your stamp, move the now-4th stamp to the TOP of the archive file.

## STATUS 2026-07-27 (iteration 22) — **`w-worldd-m2` **ITEM COMPLETE**: M2.C LANDED (PR #14 → squash `73d3486`, dev CI green, both jobs); doc → `implemented/` with EVERY acceptance + design-freeze box checked. The clause-2 local daemon ships.** Mid-sprint EXECUTE ("Plan exists" lane): no new doc, no quorum, no planner. M2.C delivered the CLI client verbs over the now-complete REST surface (742 insertions, 9 files, ONE executor round, ZERO blocking evaluator findings): **`cmd/ailang-worldd/cli.go`** — all seven frozen routes reachable 1:1 as verbs (`health`/`head`/`world get`/`object get [--payload]`/`log get`/`log range --from/--limit`/`registry get`/`commit --file`) through **ONE `http.Client` and ONE transport path**, every call deriving `context.WithTimeout` from an **injectable** timeout field so no client call can hang and the deadline test need not wait 30 s, with response AND commit-file reads both bounded and non-2xx decoding the shared `daemon.APIError` envelope including the 409's machine-readable heads; **`cli_test.go`** — the milestone's centrepiece, an end-to-end episode against a **REAL SUBPROCESS** (`go build` the binary, spawn `serve --db <temp> --bind 127.0.0.1:0`, read the **announced** address from stdout — which is exactly why A2 announces it — then drive every verb through the CLI client functions, never raw HTTP, running genesis → commit → read plus a 409 conflict re-plan), plus a bounded-deadline test against a listener that accepts and never responds; and the two folded carry-forwards — **CF-B-1** (`handleHead`'s 404/500 move to the JSON envelope while the **success path stays canonical `text/plain` `algo:digest`**, the frozen contract, and the stale "arrives with M2.B" comment is corrected) and **CF-B-4** (405-from-the-real-mux and non-zero-offset log paging, both previously implemented and unasserted). **THE EXECUTOR'S THIRD CONSECUTIVE REFUSAL TO FABRICATE.** The codex `gpt-5.6-sol` sandbox denies loopback `bind(2)`, so it could not run ANY socket test or benchmark. For the third milestone running it authored them correctly anyway, quoted every sandbox denial verbatim, and **explicitly declined to invent the numbers** — writing into `BASELINE.md` that the controller must re-measure rather than relabelling M2.B's values as an M2.C measurement. That refusal is the only reason the baseline is trustworthy: a fabricated p95 here is undetectable after the fact and would poison every future sprint that diffs against the file. The controller ran all socket-dependent gates and re-measured **all six rows in ONE 200x invocation** outside the sandbox. **ROUTING — the ratified codex pin was HONOURED, and the driver was again wrong about it.** The driver exported `MISSION_EXECUTOR_MODEL=opus` because its codex pre-flight still fails `rc=127` on the `ailang#493` PATH gap (unfixed upstream; `tools/launchd/mission-control.sh:44` omits `/opt/homebrew/bin`). Per the Repo Profile's iter-19 rule the controller re-probed **WITH `--model`** and a fixed PATH → **rc=0, `ok`, codex-cli 0.145.0, `auth_mode=chatgpt`** — so the charter's ratified `codex:gpt-5.6-sol` pin ran, rather than a provably spurious opus fallback. **Controller's own independent evidence** (never laundering a sub-agent claim): both gates green on the pinned v0.30.0 with `host/replay` **RUNNING 13.6–14.0 s not SKIP**, `verify_ail.sh` at EXACTLY 4/4 identities / 9 modules / 14 tests, bench smoke green on all six names, `go test ./cmd/... -v` with **ZERO skips**, gofmt/vet clean, scope clean **by diff not by claim** (`host/store/**` incl. `schema.sql`, `host/{hashref,canon,archive,registry,replay}`, `world/**`, `scripts/**`, `.github/**`, `go.mod`, `go.sum` all byte-unchanged; no new deps, no new store methods); and **FIVE MUTATIONS, each RED then reverted GREEN** — CF-B-1 reverted (RED at two independent tests); the D7 client deadline swapped for `context.WithCancel` (**the deadline test HANGS until `go test`'s own 90 s timeout kills it** — the strongest possible proof it observes a real timeout rather than asserting the source calls `WithTimeout`); `--payload` parsed but never applied (RED, proving the e2e asserts response **content**, not exit codes); a GET route registered without its method prefix (RED at 405, proving the mux is what is under test); and the log-range offset never applied (RED, `indices begin at 0, want 37`). **A SIXTH MUTATION WAS REFUTED AS BEHAVIOUR-EQUIVALENT — recorded, not buried.** Collapsing the registry client's per-segment escaping into a whole-string `url.PathEscape` was expected to break the multi-segment `{name...}` route but stayed GREEN; re-checking the mutation instead of believing its result showed `PathEscape` turns `/` into `%2F` and Go's mux unescapes the wildcard to the **identical** `PathValue` (standalone probe: both URL forms → `200 name="world/epoch-registry/v1"`). The test was never at fault; the per-segment loop is defensive, not load-bearing. **Two further mutation attempts were DISCARDED before scoring because they failed to COMPILE** (`declared and not used`) — a build break proves nothing about a gate's strength, so each was reformulated into a compiling, behaviour-changing form first. **sprint-evaluator (sonnet; generator≠judge is a CROSS-PROVIDER split — codex/OpenAI vs Anthropic) PASS 89→96/100, ZERO blocking**, and it earned its keep: it **independently reproduced the mux/escaping refutation with its own probe** rather than accepting the controller's account, ran three of its own mutations, and re-ran every gate first-party. Its four non-blocking findings are enumerated CF-C-1…CF-C-4. **A measurement honesty correction also lands**: M2.B recorded the loopback transport tax as "~0.03 ms", but the same difference measured 0.03 ms and 0.10 ms across two runs because the **floor** moved (store-commit p95 0.5421 → 0.4717 ms) while the REST row barely did — at sub-millisecond magnitudes one run cannot resolve that to two decimals, so `BASELINE.md` now states only the durable claim (**well under 0.1 ms; commit cost is dominated by the kernel's fsync**) and records that rows are re-measured every milestone rather than carried forward. Gate 3b: PR run `30297461991` + dev-merge run `30297536500` both **completed/success**, both jobs each, verified by direct per-workflow query against `origin/dev` = `73d34862f` — not from the poll alone. Worktree removed, branch deleted. `metered=$0.00` (codex probe + executor on the ChatGPT subscription lane with `env -u OPENAI_API_KEY`; evaluator on a subscription Agent-tool pin; $5 ceiling untouched). Preflight clean: armed, billing CLEAN, `sunholo-voight-kampff`, dev==origin/dev, CI green, no `[nightly-eval]` issues, inbox triaged to zero (mission-v1's iter-107 report carried no request for World), no new `@MarkEdmondson1234` comment on #9 (8 comments) or predecessor #1, no rotation due. **Next: the queue head is now `w-effect-broker-m3`** (clause-3 effect broker — FS/Git/Model/Human.Approve handlers, capability + budget checks, the first physical isolation floor), a NEW-DOC item requiring design-doc-creator on the rotation designer + a pick-time quorum.

## STATUS 2026-07-27 (iteration 21) — **`w-worldd-m2` **MILESTONE M2.B LANDED**: the full REST v1 surface — PR #13 → squash `b412699`, dev CI green (both jobs).** Mid-sprint EXECUTE ("Plan exists" lane): no new doc, no quorum, no planner. Landed the remaining five frozen routes plus `POST /v1/commit` (1,218 insertions, 7 files): object reads with payload **opt-in behind `?payload=true`** so the default response is bounded; `GET /v1/log?from&limit` as a **bounded loop over the existing `GetLogEntry`** with **zero new store methods**, `limit` through the Z3-proven `clampLimit`; the multi-segment `{name...}` registry wildcard (registry IDs contain slashes); one shared JSON error envelope whose classes/statuses mirror the sketch's `httpStatus` exactly; and the commit body wrapped in `http.MaxBytesReader` at the Z3-proven `maxCommitBytes` → 413. **ROUTING FINDING — the codex lane was silently OFF for BOTH missions.** The driver exported `executor=opus` because its codex pre-flight failed **rc=127 `exec: codex: not found`**: `tools/launchd/mission-control.sh:44` exports `PATH="$HOME/go/bin:$HOME/.local/bin:$PATH"` and **omits `/opt/homebrew/bin`, where codex lives** — `claude` resolves, codex cannot. **rc=127 is a PATH gap, not a spent quota and not an unusable pin** (the iter-18/19 scar in a third costume). Per the Repo Profile's own iter-19 rule the controller re-probed WITH `--model` and PATH fixed → **rc=0, `ok`, codex-cli 0.145.0, `auth_mode=chatgpt`** — so the ratified `codex:gpt-5.6-sol` pin was HONOURED rather than a provably spurious fallback. **Pre-existing, and MASKED not caused by #486**: the old gate fell back only on `QUOTA_SIG`, and `command not found` does not match it, so a 127 false-greened the lane. Upstream `dev`'s copy has the identical line 44 + probe (fetched via `gh api`), so **the V1 loop demotes to opus every fire too** — the exact opposite of the codex flip's quota-relief intent. FROZEN CORE → no local patch; routed both channels as **`sunholo-data/ailang#493`** + `msg_20260727_183949_4951d6bc`. **A PRIOR FIRE WAS LOST**: the 16:12 run was killed by the driver's stall watchdog at 17:02 (rc=143) leaving **no artifacts** (verified: no worktree, no branch, no commit, no #9 comment, no log entry) — recorded so the numbering gap is not later read as a missing entry. **THE BLOCKING FINDING — an API that could not express a commit its own kernel accepts.** `POST /v1/commit` answered a genesis commit (zero observed head) with **400 `observedHead: hashref: empty hashref text`**, so the acceptance check "a genesis+commit episode driven **entirely** over REST" was unreachable. **Three parties found it independently**: the codex executor flagged the tension in its own final message and declined to invent a zero-ref encoding; the controller proved it with a first-party throwaway probe (embedded **accepts**, REST **400**) rather than forwarding the executor's claim; and the **sprint-evaluator BLOCKED round 1 (82/100)** having derived it separately, adding the cause nobody else named — the helper was called **`seedRESTGenesis` while calling `PutWorld`/`SelectHead` directly**, a misleading name that hid the gap from the executor's own review. **The fix, and its deliberate narrowness**: `parseGenesisRef` accepts `""` as the zero `HashRef` for **`observedHead` ONLY** — not a second encoding, since `HashRef.String()` **already** renders the zero value as `""`, so this round-trips the existing one. It **deliberately does NOT extend to `prevEntryHash`**: `store.Commit` **WRITES** a zero there that `store.GetLogEntry` **CANNOT READ BACK**, so accepting it would let a client append a log entry no reader can ever load — a worse defect than the one being closed. **The first draft of the fix DID include `prevEntryHash`; the new equivalence test failed against it, and that failure is what surfaced the M1 store asymmetry** (carried as CF-B-2, a kernel-side decision). **Controller's own independent evidence**: both gates green on the pinned v0.30.0 with `host/replay` **RUNNING 13.2s not SKIP**, `verify_ail.sh` at EXACTLY 4/4 identities / 9 modules / 14 tests, bench smoke rc=0, gofmt/vet clean, scope clean **by diff not by claim** (`host/store`, `world/`, `go.mod`, `go.sum`, `tools/`, `design_docs/`, `.github/` all byte-unchanged; no new deps, allowlist green); and **FIVE MUTATIONS, each RED then reverted GREEN** — clamp ceiling 500→1000 (RED "returned 510 items, want 500"); `MaxBytesReader` removed (RED at **400 not 413**, precisely the plan's predicted null case); the 409's machine-readable heads dropped **while keeping the prose message** (RED — proving a prose-only 409 is a dead end for a re-planning client); `observedHead` reverted to strict (RED, re-introduces the bug); the lenience over-widened to `prevEntryHash` (**RED at status 200** — it succeeds in writing the unreadable entry). **Evaluator round 2: PASS 89/100, ZERO blocking** — it re-derived the `prevEntryHash` asymmetry with its own probe rather than accepting the controller's account, cross-checked the N+1 arithmetic, and **caught a defect the CONTROLLER introduced in the fix**: the `parseGenesisRef` doc comment cited `TestGenesisRefLenienceIsExactlyTwoFields`, a test that does not exist after the fix narrowed from two fields to one. Fixed before merge — a comment naming a non-existent gate is how a future reader concludes there is no gate. **`bench/BASELINE.md` now covers the FULL surface with NO PENDING rows** (all six re-measured in ONE 200x invocation so no row comes from a different run than its neighbours): REST commit **0.5763 ms p95** (61×, a ~0.03 ms transport tax over the embedded floor — essentially all commit cost is the kernel's fsync); log range **1.248 ms** at limit=100 and **4.915 ms** at limit=500 — **3.94× the time for 5× the rows**, i.e. the deliberate N+1 is linear and 24× inside budget, so **no range-read store method is justified**, and the file now records what evidence would overturn that. Gate 3b: PR run `30287931351` + dev-merge run `30288051202` both **completed/success**, both jobs each; the three new benchmarks were confirmed to actually **RUN** on linux/amd64, not merely be wired; worktree removed, branch deleted. `metered=$0.00` (codex probe + executor on the ChatGPT subscription lane with `env -u OPENAI_API_KEY` — the key **was** set, so the strip was load-bearing; both evaluator rounds on subscription Agent-tool pins; $5 ceiling untouched). Preflight clean: armed, billing CLEAN, `sunholo-voight-kampff`, dev==origin/dev (`f61aafb`→now `b412699`), CI green, inbox empty, no `[nightly-eval]` issues, no new `@MarkEdmondson1234` comment (newest EQUALS the watermark `2026-07-27T08:55:11Z`, already actioned in iter-18; predecessor #1 re-checked), **no rotation due — #9 was created 05:51Z, AFTER the 05:00Z Monday boundary, which is iter-20's intent test applied for the first time**. **ONE process fix, no skill edit**: `/tmp` scratch paths for the cross-provider recipe (`codex_directive.txt`, `codex_run.sh`, `codex_last.txt`) are SHARED with the V1 loop on this rig and were found holding V1's content — now mission-namespaced (`*_worldb.*`). No routing-policy change. **Next: milestone M2.C** — CLI client verbs over the now-complete REST surface + a real-subprocess end-to-end episode + close-out, folding CF-B-1 and CF-B-4. **M2.C LANDS the item.** CF-A3-4 is closed EARLY (no PENDING baseline rows remain), so M2.C's baseline work is a re-measure-and-diff, not a fill-in.

## STATUS 2026-07-27 (iteration 20) — **`w-worldd-m2` M2.A checkpoint **A3 LANDED → MILESTONE M2.A COMPLETE**: PR #12 → squash `9579fe1`, dev CI green (both jobs).** Mid-sprint EXECUTE ("Plan exists" lane): no new doc, no quorum, no planner. Delivered Decision 6's three required parts + carry-forward CF-7 (300 insertions / 2 deletions, 5 files, ONE round, zero blocking findings): **`host/daemon/bench_test.go`** — `BenchmarkStoreCommit` (embedded `store.Commit`, the kernel floor, fresh **temp-file** store per run so fsync reality is in the number and A1's writer lock is exercised; each iteration commits a distinct world chained on the previous head, so compare-and-append does real work every time), `BenchmarkHeadRead` + `BenchmarkHealth` (real loopback round-trips against a running daemon on an ephemeral port, keep-alive warmed OUTSIDE the measured region), each collecting PER-ITERATION wall-clock samples reported as p50/p95 via `b.ReportMetric`, with the percentile helper correct at N==0/1 (the `-benchtime 1x` path CI takes); **`scripts/bench_worldd.sh --smoke`** asserting a **HARDCODED MANIFEST of benchmark NAMES**, not a line count, because `go test -bench` exits 0 on no-match — the V27/B1 vacuous-pass class; **`bench/BASELINE.md`**, the committed day-1 budget; and one new CI `go-verify` step. **THE CODEX EXECUTOR PIN WORKED — first successful real run since the iter-19 incident**: the skill's own probe WITH `--model` returned rc=0 on **codex-cli 0.145.0** (Mark's attended fix), `auth_mode=chatgpt`, invoked with `env -u OPENAI_API_KEY` — and the key **was** present in the tool shell, so the strip was load-bearing, not ceremonial. Executor ran one round, rc=0 inside its 30-min cap, **pin HONOURED, no fallback**. **The executor's headline behaviour was HONESTY UNDER A DEGRADED ENVIRONMENT**: the codex `workspace-write` sandbox denies loopback `bind(2)`, so the two HTTP benchmarks could not run inside it at all — it measured the store-commit row, recorded the others as **UNAVAILABLE with the sandbox error quoted**, wrote in `BASELINE.md` that the artifact must be refreshed outside the sandbox before A3 is accepted, and **explicitly declined to invent values**. Fabricated p95s here would have poisoned the day-1 budget permanently, since every later sprint diffs against this file. The controller measured the HTTP rows on the dev rig outside the sandbox and completed the table, recording the split provenance verbatim. **Controller's own independent evidence** (never laundering a sub-agent claim): both gates green on the pinned v0.30.0 with `host/replay` **RUNNING 12.5s not SKIP** and `verify_ail.sh` at EXACTLY 4/4 identities / 9 modules / 14 tests; `gofmt`/`go vet` clean; scope clean **by diff not by claim** (`daemon.go`, `go.mod`, `go.sum`, `host/store/**` and every other host package byte-unchanged; no new deps, allowlist still green); and **THREE MUTATIONS, each RED then reverted GREEN** — (1) renaming `BenchmarkHealth`→`BenchmarkHealthX` still emitted three benchmark lines and still exited **0** from `go test`, and the gate caught it **BY NAME** (a line-count check would have passed clean); (2) making `releaseFromVersion` return a different release per call turned the CF-7 test RED with a divergent-head error while the **same** mutation is GREEN against the pre-CF-7 body, proving CF-7 closed a genuinely open gap; (3) injecting 5 ms into the measured region moved head-read p50 **0.11 ms → 7.02 ms**, proving the percentiles track real wall clock. **Controller self-correction worth keeping**: mutation 2's FIRST form appended its counter *before* the newline split, so it was a silent no-op and its green result briefly read as "CF-7 is not load-bearing" — re-checking the mutation rather than believing its result turned a would-be false finding into the strongest proof in the set. **Measured day-1 baseline** (M4 Max, darwin/arm64, go1.26.4, one `-benchtime 200x` invocation): store commit p95 **0.6093 ms** (≤25), head read p95 **0.08596 ms** (≤5), health p95 **0.06288 ms** (≤2) — all inside budget with 32–58× headroom; REST-commit + both log-range rows explicitly **PENDING M2.B**. **sprint-evaluator (sonnet; generator≠judge is a CROSS-PROVIDER split here — codex/OpenAI vs Anthropic) PASS 89/100, ZERO blocking**; it independently confirmed all three mutations, the N==0/1/2 percentile math, and that the manifest is genuinely hardcoded rather than derived from the output it checks. **The judge refuted two controller claims and was right both times** (recorded — a judge that only ratifies is worthless): "`host/replay` does not `t.Skip`" holds only in `verify_go.sh`'s context, which exports `AILANG_BIN` (bare, it skips all 10); and the `go test | tee` exit-code propagation works only because the script is bash with `set -o pipefail`. **The controller in turn refuted the judge's one deduction** (−4, "hashing inside the measured window"): every `SumSHA256` is at lines 57–77, `start := time.Now()` is line **80** — the p50/p95 window contains only `s.Commit`; the judge visibly contradicted itself mid-sentence, its valid kernel (in-loop hashing inflates `ns/op`, Go's mean) is already disclosed in `BASELINE.md`, so it is recorded as **refuted-as-stated, not adopted**, with no code change. Gate 3b: PR run `30275747631` + dev-merge run `30276704692` both **completed/success**, both jobs each; **the new CI step was confirmed to actually RUN** (go-verify log shows the smoke gate on `linux/amd64 / AMD EPYC 7763` emitting real p50/p95, so the N==1 path works on a different OS+arch than the dev rig); worktree removed. `metered=$0.00` (codex probe + full executor run on the ChatGPT subscription lane; evaluator on a subscription Agent-tool pin; $5 ceiling untouched). Preflight clean: armed, billing CLEAN, `sunholo-voight-kampff`, dev==origin/dev (`bfbd94e`→now `9579fe1`), CI green, no `[nightly-eval]` issues, no new `@MarkEdmondson1234` comment on #9 (6 comments) or predecessor #1 (watermark unchanged at `2026-07-27T08:55:11Z`), inbox = 1 eval-suite FYI + this loop's OWN iter-19 outbound (neither outranking). **Two process fixes, no skill edit** (no gap reached the ≥2-instance bar): (1) a fresh worktree **never** contains the sprint plan JSON — `.gitignore`'s `**/.ailang/` makes `.ailang/state/sprints/*.plan.json` structurally absent from every `git worktree add`, and the codex executor reported it as a mid-run blocker; every worktree executor directive must copy it in or quote the binding requirements inline; (2) the weekly-rotation rule misfires on a thread created just before Monday 07:00 (#9 was created 05:51Z, ~1h early, so the literal rule reads ROTATE and would open a second thread for the week #9 already titles) — the intent test is now recorded. No routing-policy change (one clean codex run is not the charter's ≥3 rows). **Next: milestone M2.B** — the remaining five REST routes (`POST /v1/commit` with the 8 MiB cap and 409/400/404/413 mapping; `/v1/log` with the deliberate N+1 at limit=100 and the clamp max 500; object/head reads), folding CF-A3-2, CF-A3-3 and CF-5, extending both the bench manifest and `BASELINE.md`. Then M2.C lands the item.

## CURRENT GOAL

1. ~~Iteration 0: ratify the bar~~ — **DONE 2026-07-23 attended** (Mark: clause-4 fixed at
   −2pp / ≤25%; bar + Conflict Surface + guardrails + queue ratified as drafted; re-quorum run
   attended — see STATUS).
2. **NOW**: work the queue through the inner loop (design-doc → sprint-plan → execute →
   evaluate), one sprint-sized item per iteration, recording routing evidence every time.
   **D1 RATIFIED 2026-07-24 (Mark, attended: A+B-metadata) — queue UNBLOCKED.**
   Next item: `[NEXT] w-world-library-m1` (clause-1) — M1 may now freeze the log format per the
   settled decision doc.

## The bar — what "Ailang World 1.0" must meet (RATIFIED by Mark 2026-07-23, attended — see STATUS + issue #1 record)

<!-- Distilled from DESIGN.md §14–§17 milestones + §12.4 value thesis + §13.5 priorities.
     Deliberately fewer clauses than milestones: clauses are CHECKABLE END-STATES, not work items. -->

- **Clause 1 — deterministic kernel**: an immutable, content-addressed world store with an
  append-only transition log and **proven deterministic replay** — replaying any recorded episode
  reconstructs state bit-for-bit from pure transitions + recorded effect results. The log-epoch
  semantics decision (DESIGN.md open question 11: interpreter version pinning + content-addressed
  transition functions) is made and enforced BEFORE the log format freezes.
- **Clause 2 — local-first daemon**: `ailang-worldd` runs on one machine over SQLite with CLI +
  REST, **zero cloud dependencies in the core**. Cloud transports exist only as effect-handler
  extensions (or not at all in 1.0).
- **Clause 3 — explicit authority end-to-end**: every effect goes through the broker with a
  capability + budget check; effect results are recorded (replay input); capsules run with a
  physical isolation floor beneath the semantic checks. No ambient-authority path exists from an
  agent to the outside world.
- **Clause 4 — the M4 RESIDENT-AGENT NON-INFERIORITY FLOOR (do-no-harm kill-switch; thresholds
  ratified 2026-07-23; agent definition + framing AMENDED by Mark same day, attended)**: this
  clause is the FLOOR, deliberately NOT the value proof — agents are World's *residents*, not
  its destination, and a substrate that taxes its residents dies before its systemic value can
  accrue. World's VALUE is measured by clause 5's provenance teeth and the R1 standing evidence
  (below). The floor: **TWO reference agents from different providers — Claude Code (agent mode)
  and codex CLI** — each in paired arms: *shell* (native tools) vs *World* (MCP transition tools
  only), same model, same benchmark set (standard tier), N≥3 runs per arm. **The floor holds
  only if BOTH reference agents show: World pass-rate ≥ shell pass-rate − 2 percentage points
  AND median wall-clock overhead ≤ +25%** (non-inferiority, not superiority — by design).
  **Stability precondition (gate validity, not gate outcome)**: a reference agent is
  gate-eligible only if its SHELL arm completes N≥3 runs with zero harness-fault failures
  (api_error / resource_limit / harness classes) and pass-rate range ≤5pp; an ineligible agent
  is fixed or substituted — if both become ineligible the gate PAUSES for instrumentation (a
  measurement problem, explicitly NOT a value-park). If the gate fails on eligible agents after
  honest tuning, **World parks**. **motoko is an OPTIONAL third arm** — the local-first/R5
  narrative, run only if it meets the same precondition, informative and NEVER gate-blocking
  (amendment rationale: motoko's current instability would make the baseline metrologically
  unsound — Mark 2026-07-23).
- **Clause 5 — the human surface works, with PROVENANCE TEETH (the 1.0 value demonstration;
  strengthened by Mark 2026-07-23)**: one real goal expressed → proposals → approval inbox with
  evidence bundles → commit (SCENARIOS.md scenarios 1–3, live) — AND, measurably: **on ≥3 REAL
  "why did X happen" questions (arising from actual operation, not synthetic), a provenance walk
  yields the verified answer in ≤5 minutes each**, where the pre-World method was a grep/log
  archaeology session. This — capability the shell cannot express at any pass-rate — is what
  World is FOR; clause 4 only proves it isn't paid for with resident-agent regression.
  **Standing value evidence beyond the 1.0 bar: R1** — this mission itself migrating onto World,
  measured on incident classes eliminated + human attention per workstream against the 2026-07
  markdown-and-launchd baseline (tonight's operation IS the control arm).
- **Clause 6 — protocol-native boundary**: the transition registry is served over MCP
  (capability-filtered per session) and an A2A agent card is published; **no new wire protocols**
  (DESIGN.md §3.7).
- **Clause 7 — controlled self-modification proven once**: at least one World behavior change
  ships as an extension package through World's own propose → verify → commit pipeline, with the
  §14 boundaries enforced (mission-loop machinery, AILANG compiler, live daemon EXCLUDED from
  self-mod scope).

## Conflict Surface (RATIFIED 2026-07-23 with the bar; drafted post-iter-0 to answer quorum objection #2)

What this mission touches or overlaps, and the drawn boundaries:

- **`ailang serve-api` (ailang repo)** — the overlap is the PROTOCOL BOUNDARY only.
  **VERIFIED (code inspection 2026-07-23, `internal/apiserver/` at ailang HEAD)**: the package
  advertises NO persistence or scheduling API — zero sqlite/sql.Open/bolt/badger/scheduler hits
  package-wide; request-scoped serving of module exports over REST/MCP/A2A (`mcp.go`, `a2a.go`,
  `handler.go`, `routes_dispatch.go`). Precision note (quorum round-5): the load-bearing form of
  this claim is "no advertised state machinery to build a world kernel on" — verified from the
  package surface; exhaustive absence-of-state through transitive calls is neither claimed nor
  needed, because the daemon justification runs the other way: worldd needs a store + log +
  broker + scheduler that apiserver does not offer, while apiserver's protocol serving IS reused
  (path (a), demonstrated live). A new daemon is justified by state, not by protocol.
  **ALSO VERIFIED**: `internal/apiserver` is a Go `internal/` package → NOT importable from
  another repo. "Reuse" is therefore one of three concrete paths, in preference order:
  (a) **primary**: World exposes its transition registry AS `.ail` modules served by
  `ailang serve-api --mcp/--a2a` — **VERIFIED for static module exports (live test
  2026-07-23)**: served `design_docs/sketches/` over `--mcp-http`, and `tools/list` returned
  `plan`/`verify`/`commit` as MCP tools with schemas + effect signatures. **NOT claimed**
  (explicitly w-mcp-projection's ACCEPTANCE CRITERIA, not premises): dynamic worldd-backed
  registry, per-session capability filtering, propose→verify→commit enforcement at the
  boundary — if those fail there, fall to (b)/(c); (b) fallback: worldd runs `ailang
  serve-api` as a sidecar process; (c) only on evidence: upstream issue to export a public
  serving package. Grafting a world
  store onto serve-api itself stays ruled out — it would move OS concerns into the frozen
  language repo, the wrong direction under PROGRAM.md.
- **The `ailang` binary / compiler (frozen core)** — worldd consumes the RELEASED binary
  (`check`/`test`/`ai-check`/serve-api) and never links compiler internals; language gaps route
  upstream as issues (guardrail above).
- **Coordinator / Collaboration Hub (ailang repo)** — pattern overlap only (approval queue,
  messaging). Patterns port; schemas do not. World's human surfaces are projections from the
  world store (SCENARIOS.md scenario 1), not extensions of the Hub; no shared database.
- **The mission loop itself** — the loop that builds this repo stays OUTSIDE World's
  self-modification scope (DESIGN.md §14). When M4 puts motoko on World, the mission loop is
  unchanged; only motoko's execution substrate changes.
- **OPEN for ratification**: whether worldd lives in THIS repo as its own Go module (DESIGN.md
  §15 assumption, coordinator-recommended: keeps the language repo frozen, lets World iterate at
  its own cadence) vs an `ailang world` subcommand upstream. Revisit only on concrete
  binary-distribution pain.

## Guardrails (mission-specific; the skill's Standing Rules always apply on top)

- **Local-first is inviolable**: no cloud dependency may enter the World core. Cloud = effect
  handlers behind the broker, post-1.0 unless ratified otherwise.
- **Frozen-kernel discipline from day 1**: default every improvement to an extension; kernel
  changes require explicit human ratification (quorum evidence attached).
- **Never touch the shared loop machinery from this mission**: `tools/launchd/*` (shared driver +
  plist) and the symlinked skills are owned by the shared infrastructure; improvements route via
  Gate-5 retro in the shared skill so ALL missions benefit. NEVER copy skills into this repo.
- **Never touch the V1 mission's state or checkout** (`~/.ailang/state/mission-v1*`, legacy
  `mission-control.*` paths, the `sunholo-data/ailang` working tree).
- **Language gaps route upstream — two channels, always both**: anything AILANG can't express
  cleanly (or any core/binary change World needs) becomes (1) a GitHub issue on
  `sunholo-data/ailang` — the durable, triageable record with repro + version per that repo's
  conventions — AND (2) an agent message to the v1 loop's inbox so it is seen without waiting for
  a triage sweep: `ailang messages send mission-control "<summary + issue link>" --title "..."
  --from "world-coordinator"` (local send on this rig; the v1 session's start hook surfaces
  unread messages). The V1 mission owns routing it (extension vs core-floor per PROGRAM.md);
  World NEVER works around the compiler locally, no vendored forks.
- **CODING STANDARDS are binding** (ratified attended 2026-07-24): [coding-standards.md](coding-standards.md) — Z3 contracts on the pure core (S1), effects at the boundary (S2), **slim kernel / package-first** (S3: every kernel addition answers "why is this not a package?"), compiler-checked docs (S4), the AILANG fluency protocol (S5), honest non-vacuous gates (S6). Evaluators score against it; changes to it are ratification-class.
- **Compiler-checked docs discipline**: every `.ail` snippet that enters `design_docs/` ships as a
  checkable file under `design_docs/sketches/` (or a successor source tree) and passes the CI
  gate. A doc claiming "this compiles" without a checked artifact is a defect.
- **Kernel performance budget from day 1** (DESIGN.md §13.5): the M4 overhead gate gets HARDER as
  models get faster — fork cost and verify latency are design constraints on M1, not later
  optimizations.
- **Advisory quorum rounds are BOUNDED** (ratified by Mark 2026-07-23, from the 5-round
  iteration-0 ledger): on ratification-class docs (charter, bar changes), at most **2 quorum
  rounds per revision cycle**. The quorum's job there is objection-SURFACING for the human
  authority gate — it holds no veto, a synthesized PASS is not required, and rounds must never
  be re-run in search of one. When objections stop being load-bearing, close on the human gate
  and record the ledger. (Design-doc quorums at pick keep their normal semantics — BLOCKED
  parks the item.)
- **A MUTATION'S RESULT IS ONLY EVIDENCE IF THE MUTATION ITSELF IS VALID** (process fix, iter-22;
  **2nd instance** — iter-20 was the 1st). Mutation testing is how this mission proves a gate is
  non-vacuous, so a bad mutation silently produces a false finding in EITHER direction. Two
  distinct failure modes have now been observed:
  1. **The mutation does not compile.** Two iter-22 attempts failed `declared and not used`. A
     build break proves nothing about a test's strength — it is not a RED. Reformulate into a
     compiling, behaviour-changing form (e.g. keep the variable used: `from*0 + offset`,
     `*payload && len(path) < 0`) and re-run before scoring.
  2. **The mutation compiles but changes no behaviour**, so its GREEN reads as "the gate is
     vacuous" when the gate is fine. Iter-20's counter was appended before a newline split (a
     silent no-op whose green briefly read as "CF-7 is not load-bearing"); iter-22 collapsed the
     registry client's per-segment `url.PathEscape` into a whole-string one, which is
     behaviour-EQUIVALENT because `PathEscape` encodes `/` as `%2F` and Go's mux unescapes the
     `{name...}` wildcard to the identical `PathValue`.
  **The rule**: before recording ANY mutation result, confirm the mutation (a) COMPILED and
  (b) actually changed observable behaviour. An unexpected GREEN is a claim about your mutation
  first and about the gate second — re-check the mutation before concluding the gate is weak.
  Record refuted mutations in the log's Ruled out ledger rather than dropping them; both times,
  re-checking turned a would-be false finding into the strongest evidence in the set.
- **Kill switch stays until ratification**; only Mark (or the v1 agent on his instruction) arms
  the loop.

## Routing policy

Uses the **shared** per-role model routing from `mission-control` (controller /
designer-rotation / planner / executor / evaluator, generator≠judge enforced). Overrides for THIS
mission in `~/.config/ailang/mission-world.env`:

- **Executor default**: NON-Anthropic lane — `codex:gpt-5.6-sol` (`MISSION_EXECUTOR_MODEL`),
  per the shared quota plan. **iter-19 routing incident RESOLVED by Mark 2026-07-27 (attended,
  option c): codex-cli upgraded 0.137.0 → 0.145.0 (npm -g) and the pin LIVE-PROBED working**
  (`env -u OPENAI_API_KEY codex exec --model gpt-5.6-sol` → PIN-OK, ~21 tokens, ChatGPT-
  subscription OAuth — Mark verified server-codex OAuth2 same day, so the lane works headless
  in mission loops AND evals). The opus fallback stays the documented degrade path only.
  Upstream driver-probe gap (probe omits `--model`, false-greens an unusable pin) remains
  filed as ailang#486 for the v1 lane.
- **Evaluator**: must differ in provider from the executor (generator≠judge) — sonnet or
  gemini when the executor is codex.
- Otherwise: inherits the shared defaults.

## Queue (top = next; tags: [NEXT] [IN-SPRINT] [PARKED] [LANDED] [RULED OUT])

<!-- Every open item carries a clause tag. Estimates are honest guesses at bootstrap;
     iteration 0 re-scores. NEW-DOC items start with design-doc-creator. -->

**[LANDED 2026-07-24 (iter-13), CI green on dev `d0009c8`] w-m1-ailang-hardening** · clause-1 ·
The M1 AILANG surface now USES AILANG's distinguishing features. Doc
`design_docs/implemented/w-m1-ailang-hardening.md` (Fable r1 designer; **codex/gpt-5.6-sol iter-13
revision**; controller carve-out 2nd revision). Shipped via **PR #5 → squash `d0009c8`,
sprint-evaluator PASS 97/100** (sonnet, generator≠judge). **4 Z3-proven contracts**
(`transitions/applyRevision`, `contracts/isValidNextWorld`, `logepoch/sameRef`+`servesEntry`) +
**14 named inline tests** + a **hardcoded, bounded, non-vacuous required-check-manifest gate**
(`scripts/verify_ail.sh`). **iter-12→13 descope (empirically forced):** a contract on any
`Proposal`-taking predicate Z3-errors `unknown sort 'Proposal'` (Proposal.evidence is an ADT) and
`ai-check` exits 0 SILENTLY → the 3 Proposal predicates (`proposalMatchesWorld`,
`verificationMatchesProposal`, `commitAllowed`) are documented-limitation + tests-only (6 tests,
shared-predicate bodies UNCHANGED → anti-drift preserved); achievable proven set = **4, not 7**.
Upstream `sunholo-data/ailang#477`. **Two new toolchain findings landed as V-rows:** V26 (`ai-check
-timeout` is per-function Z3 only, not process-wide; `ailang test` has none → the gate wraps both
legs in a hardcoded-deadline process-group-killing `run_bounded`, Standing Rule 6), and **V27
(`ai-check` shells to an external `z3` and SKIPS SILENTLY without it — a bare `ubuntu-latest` runner
has no z3, so every contract vanished from `verify.results[]` and PR #5's first CI went red; fixed
by installing Z3 4.16.0 in the `ailang-verify` CI job, sha256-pinned)**. Retrofit of M1 milestone-1
(`world/{contracts,transitions,logepoch}.ail`; `types.ail` byte-identical); root cause was
discoverability (`.mcp.json` + upstream #476). Effects/package-extensions correctly out of scope
(frozen core). See memory `ailang-feature-discoverability-gap`. · queued iter-10 (Mark-directed).

1. [LANDED 2026-07-24] **w-log-epoch-decision** · clause-1 · ALL THREE decisions settled: D2
   (content-addressed transition fns) + D3 (SHA-256 + tagged HashRef) via quorum; **D1 RATIFIED
   by Mark attended: A+B-metadata hybrid** — authoritative pin = exact-binary content hash;
   release tag+commit + semanticsEpoch as compatibility metadata in every header (M8 promotion =
   policy change, never a format migration); hermeticity → replay-doubling in M1 acceptance.
   Doc: `design_docs/planned/w-log-epoch-decision.md` (SETTLED — ready for M1).
2. [LANDED 2026-07-27 (iter-16), all 6 milestones CI green on dev; doc → `implemented/`]
   **w-world-library-m1** · clause-1 · Sprint AUTHORIZED (Mark, attended, iter-8-prior). Sprint plan
   (`.ailang/state/sprints/w-world-library-m1.plan.json`): 6 CI-green milestones. **M1 (pure AILANG
   semantic library — `world/{logepoch,types,contracts,transitions}.ail`, 285 LOC) LANDED** via
   PR #2 → squash `9d61d663e`, evaluator PASS 93/100. **M2 (Go host bootstrap — `go.mod` +
   `host/hashref` + `host/canon`, 635 LOC / 45 tests, Decisions 2+3) + CI `go-verify` job LANDED**
   via PR #3 → squash `d5b155c`, evaluator PASS 97/100 (generator≠judge). **M3 (SQLite store +
   atomic append-only log, Decision 4 — `host/store/{schema.sql,store.go,store_test.go}`, 1004 LOC;
   pinned pure-Go `modernc.org/sqlite v1.54.0`, single-transaction compare-and-append with
   `ConflictError`, pair-keyed verify cache) LANDED** via PR #4 → squash `a901c30`, evaluator
   PASS 88/100 (generator≠judge; −5 non-blocking `store_heads` schema-split, carried to M4).
   **M4 (interpreter artifact archive + epoch-1 registry bootstrap, Decisions 5+6 —
   `host/archive/archive.go` single-stream hash-while-copy + atomic-rename + sidecar manifest +
   `ReplayError`; `host/registry/registry.go` `world/epoch-registry/v1` idempotent bootstrap;
   `store_heads` M3 carry-forward folded into `schema.sql`) LANDED** via PR #6 → squash `8133573`,
   evaluator PASS 88/100 (generator≠judge; 2 non-blocking carry-forwards to M5/M6: a `KindExecFailure`
   test + a sidecar-absent idempotence-recovery edge). **M5 (replay engine + replay-doubling +
   fixture episode, Decision 7 — `host/replay/replay.go` authoritative replay resolving the exe from
   each entry's `Interpreter` HashRef, delegating per-transition execution to the archived binary as a
   bounded subprocess [never reimplements the interpreter, §14]; `replay_test.go` 12 tests: bit-for-bit,
   replay-doubling A==B==recorded, epoch-candidate-cannot-redirect, pair-member→cache-miss; committed
   goldens; folded the 2 M4 carry-forwards) LANDED** via PR #7 → squash `ef06937`, evaluator PASS 73/100
   (generator≠judge). **B1 CI false-green fixed in-PR** (`8ba5fe9`): go-verify now downloads the pinned
   **v0.30.0** binary (by TAG, sha256-verified, version-asserted) + exports `AILANG_BIN`, so the replay
   tests actually RUN in CI instead of silently `t.Skip`-ing — verified in the go-verify log. **M6 (CI Go
   gate finalize + `scripts/verify_go.sh` + final green sweep + carry-forwards NB2/NB5) LANDED** via PR #8
   → squash `a07ac96`, evaluator PASS 96/100 (generator≠judge, ZERO blocking): `scripts/verify_go.sh` = the
   durable local Go gate with an anti-false-green guard (fails loudly if `AILANG_BIN` unset/≠v0.30.0 — closes
   the V27/B1 silent-skip class locally too); CI `go-verify` now calls it; NB2 resolved by a genuine
   end-to-end replay through a byte-distinct second interpreter (`TestInterpreterMemberChangeDrivesReal-
   ReplayEndToEnd`). **ITEM COMPLETE**: all 16 acceptance boxes checked, doc → `design_docs/implemented/
   w-world-library-m1.md`. The semantic world library ships: World/Proposal/Transition/Evidence in AILANG
   (ai-check green), Go host (SQLite store + content-addressed objects + append-only log + archive +
   epoch registry + replay), replay proven bit-for-bit — the embedded-host-only M1 floor before
   daemon/broker/isolation (clauses 2–4, queued next).
   The item: the semantic world library: World/Proposal/Transition/
   Evidence types in AILANG (ai-check green in CI), Go host for the SQLite store +
   content-addressed objects + append-only log; replay of a recorded episode proven bit-for-bit ·
   ~2–3d. Design doc `design_docs/planned/w-world-library-m1.md` written (codex/gpt-5.6-sol
   designer, rotation) and quorum-**direction-accepted**; net-axiom +12, all 9 epoch-doc M1
   implications mapped, verified against the **pinned released `AILANG v0.30.0`** (checksum
   `sha256:ac3174e0…`, not the rig's `-dirty` build). Two gemini-3-1-pro quorum rounds each
   surfaced ONE narrow, non-direction defect (r1 dirty-binary → FIXED; r2 `ailang replay`
   Conflict-Surface overlap → clarified in-doc, §14-forced). ~~SPRINT PARKED~~ → **AUTHORIZED
   (see tag above); next fire builds M1.**
3. [**LANDED 2026-07-27 (iter-22) — ITEM COMPLETE**, all milestones dev-CI-green (both jobs), doc → `design_docs/implemented/w-worldd-m2.md` with EVERY acceptance + design-freeze box checked] **w-worldd-m2** · clause-2 ·
   `ailang-worldd` local daemon: SQLite, REST API, CLI; zero cloud deps; kernel perf budget
   measured and recorded from the first commit. **Shipped across five merges**: **A1** `b0deedb`
   (PR #10) the RATIFIED single-writer kernel change — `store.Open` is the fail-closed writer path
   (non-waiting `flock(LOCK_EX|LOCK_NB)`, structured `WriterAlreadyActive`), additive
   `store.OpenReadOnly`, `:memory:` exempt, `schema.sql` byte-unchanged, enforcement proven
   CROSS-PROCESS (a `LOCK_EX`→`LOCK_SH` mutation reds the suite, so an in-process mutex cannot
   pass); **A2** `39b2115` (PR #11) the daemon shell — loopback guard mirroring the Z3-proven
   `isLoopbackHost` with **no override flag**, the D7 bound-constant block (all four
   `http.Server` timeouts + `maxCommitBytes` pinned to the Z3-proven `withinCommitBytes`),
   fail-closed lifecycle with stdout address announce and bounded drain, `/v1/health` + `/v1/head`,
   and **`TestDaemonDependencyAllowlist` — zero-cloud ENFORCED over the real build graph, not
   asserted**; **A3** `9579fe1` (PR #12) Decision 6's perf machinery — p50/p95 harness on a
   fresh temp-FILE store, `bench/BASELINE.md`, and a bench-smoke gate asserting a **hardcoded
   benchmark-NAME manifest** (a dropped benchmark still exits 0 from `go test` — the V27/B1
   vacuous-pass class, closed for benchmarks); **M2.B** `b412699` (PR #13) the full REST v1
   surface — payload opt-in, `GET /v1/log` as a bounded loop with **zero new store methods**, the
   multi-segment registry wildcard, one shared JSON error envelope, body capped → 413, and the
   genesis-over-REST gap found independently by three parties; **M2.C** `73d3486` (PR #14) the
   CLI client verbs 1:1 over the frozen table (ONE transport path, injectable D7 deadline,
   bounded response + commit-file reads), a **real-subprocess** end-to-end episode driving every
   verb through the CLI (genesis → commit → read, plus a 409 re-plan), a bounded-deadline test
   that genuinely observes a timeout, and carry-forwards CF-B-1 + CF-B-4 folded. Design doc
   authored by the Fable designer (rotation) and quorum-run over two rounds; the single-writer
   fork was **ratified by Mark attended (arm A)**. Executor `codex:gpt-5.6-sol` (ratified pin) on
   three of five merges; independent judge `sonnet` throughout (generator≠judge, cross-provider).
   **Open carry-forward OUT of the item**: **CF-B-2** — `store.Commit` writes a zero
   `PrevEntryHash` that `store.GetLogEntry` cannot read back; the daemon refuses it at the
   boundary but the embedded API can still trigger it. Kernel-side decision → a store-hardening
   queue item, NOT M2. Judge's M2.C non-blocking findings: CF-C-1 (`--limit 0` indistinguishable
   from unset), CF-C-2 (registry escaping untested for encoded characters), CF-C-3 (CF-B-2 needs
   a tracking item + repro fixture), CF-C-4 (405 asserted on 2 of 7 GET routes).
4. [**NEXT** — queue head as of iter-22; **NEW-DOC**: route to design-doc-creator on the rotation
   designer (next entry after `~/.ailang/state/mission-world-designer-rotation`, currently
   `codex:gpt-5.6-sol`), then the pick-time quorum, then sprint-planner. Before spawning the
   designer, run `grep -ri "w-effect-broker-m3" design_docs/` — a NEW-DOC tag is a claim, not a
   fact] **w-effect-broker-m3** · clause-3 · effect broker with FS / Git / Model (`std/ai`) /
   Human.Approve handlers; effect-result recording; capability + budget checks; first physical
   isolation floor · ~2–3d. Now unblocked: it builds on the landed `w-worldd-m2` daemon, whose
   single-writer authority is exactly the property that makes broker-mediated effects meaningful
   (an embedded writer bypassing capability/budget checks is the ambient-authority pattern
   clause 3 exists to end — Mark's ratification rationale, iter-18).
5. **w-mcp-projection** · clause-6 · project the transition registry over MCP + publish the A2A
   agent card (reuse `ailang serve-api --mcp/--a2a` machinery — do not reinvent) · ~1d
6. [PARKED until 2–5 land] **w-agent-floor-m4** · clause-4 · dual-reference NON-INFERIORITY
   floor: Claude Code + codex, shell arm vs World-MCP arm, paired N≥3, stability precondition
   checked first; motoko as optional third arm if eligible; report honestly; park World if the
   floor fails on eligible agents · ~3d (was w-motoko-m4 → w-value-gate-m4 → renamed with the
   2026-07-23 floor reframe; clause-5 provenance teeth carry the value burden)
6b. [DOC-DRAFTED 2026-07-27, attended — quorum + ratify at pick; GATES items 7 and all M6 work]
   **w-human-surface** · clause-5 · the founding UX design for the human↔AI surface:
   [HUMAN-SURFACE.md](HUMAN-SURFACE.md) — seven interaction principles (decision packets ·
   grounded prose · trust-gradient rendering PROVEN/TESTED/ATTESTED/CLAIMED · time-as-navigation
   · attention choreography · five zoom levels · speculation-as-gesture) + sign-the-type input
   grammar + anti-patterns (grade laundering is the cardinal sin). Mark named this the vision's
   key surface ("an ideal AI language in an AI state machine OS with a NOVEL human interaction
   surface"); no human-facing sprint routes before it is ratified · ratification session ~0.5d
7. [PARKED until 4 AND 6b land] **w-approval-inbox** · clause-5 · the approval inbox + provenance
   walk, first as CLI/generated projection (SCENARIOS.md scenario 1/3), **built to
   HUMAN-SURFACE.md** · ~2d
8. [PARKED until 4 lands] **w-self-mod-vertical** · clause-7 · one World extension shipped through
   World's own pipeline end-to-end · ~1–2d
9. [BACKLOG — infra, not critical-path] **w-verify-binary-lockfile** · clause-1-infra · a durable
   pin of the released `ailang` binary (version + sha256) for the `ailang-code` verify gate, so
   local controller verification stops depending on the rig's `-dirty` build (iter-6 pinned
   released `v0.30.0` ad-hoc — `sha256:ac3174e0…`; CI already sha256-verifies a released linux
   binary). Small; the mechanism may generalize to the SHARED `ailang-code` profile → confirm
   repo-local lockfile vs shared-skill fix with a human before implementing (do not hand-edit CI
   headless) · ~0.5d · surfaced iter-6, queued iter-7

## Premise Verification Log (quorum objection #1 — every load-bearing claim, with evidence)

| Premise | Verified | Evidence |
|---|---|---|
| `serve-api --mcp / --mcp-http / --a2a` exist (clause-6 machinery) | bootstrap + iter-0 | `ailang serve-api` help, v0.30.0 |
| `ai-check` = unified check+verify, always-JSON, **no `--json` flag** | bootstrap flag-test + iter-0 | error repro'd live; citations fixed here (`4c6080a`) and upstream (`ailang@aabb3a58c`, v1-ack on issue #1) |
| CI `CI` / job `ailang-code verify gate` green and NON-vacuous | every push since `3c8791d` | job log shows `checked 2 module(s)`, released linux binary, sha256-verified |
| DESIGN.md type sketches compile | CI-gated | `design_docs/sketches/*.ail` in every run |
| Mission state fully namespaced (no v1 collision) | dry-runs 2026-07-23 | `/tmp/ailang-mission-world.log`: distinct `mission-world.pid`, worldtest profile proof |
| Upstream routing channel works end-to-end | 2026-07-23 | defect report → v1 verified+fixed+ack'd < 1h (issue #1) |
| Registry + `ailang publish` cascade exist (clause-7 lane) | v1-mission operational history | live on multivac; World's local-first cascade mode is DESIGN.md §13/M-scope work, not assumed |
| `internal/apiserver` (serve-api impl) is stateless — no persistence, no scheduler | code inspection 2026-07-23 | package-wide grep: zero sqlite/sql.Open/bolt/badger/scheduler hits; request-scoped handlers only |
| `internal/apiserver` NOT importable cross-repo (Go `internal/`) | Go module rules + path | Conflict Surface reuse paths (a)/(b)/(c) sized accordingly; path (a) needs no upstream change |
| serve-api projects `.ail` exports as MCP tools (path (a), static case) | LIVE TEST 2026-07-23 | `ailang serve-api --mcp-http --port 8199 sketches` → `tools/list` returned `plan`/`verify`/`commit` with JSON schemas + effect rows in descriptions; server killed after, port freed. Dynamic/capability-filtered projection deliberately NOT claimed — w-mcp-projection acceptance criteria |
| `verify_ail.sh` fails loudly at N=0 (gate cannot pass vacuously) | LIVE TEST 2026-07-23 | script run against an empty `design_docs/` scratch tree → "✗ no .ail modules found — the gate would be vacuous; failing loudly", **exit code 1** |
| Driver sources `~/.config/ailang/mission-world.env` + respects role overrides | code + LIVE 2026-07-23 | `tools/launchd/mission-control.sh:46-47` sources `mission-${MISSION_PROFILE}.env`; `:238` `MISSION_EXECUTOR_MODEL` respected (default opus); dry-run log 19:45:39 echoes env-only values (repo-slug/doc/workdir) + resolved roles — sourcing proven end-to-end |
| Charter RATIFIED (authorization state, kill switch, sprint routing) | attended session 2026-07-23 | Mark's decisions recorded: clause-4 = −2pp/≤25% paired N≥3; bar+Conflict Surface+guardrails+queue as drafted; ratification comment on issue #1; STATUS stamp above is the in-doc record |
| Kill-switch + gh-issue state paths are world-namespaced in the LIVE driver (safety property) | LIVE 2026-07-23 | driver log 19:45:24/19:45:39/20:02:47: "kill switch present (`…/mission-world.disabled`) — skip" — the driver printed and HONORED the world-namespaced path 3× while the v1 loop ran unaffected; iter-0 posted to issue #1 = `mission-world-gh-issue` read correctly |
| `ailang messages` channel works end-to-end (guardrail's delivery leg) | live round-trip 2026-07-23 | sent `msg_…_2c6964d3` (defect report) + `msg_…_acc5edcc` (channel test); v1 agent RECEIVED and acted — upstream ack on issue #1 at 18:27Z citing the report, fix `ailang@aabb3a58c` |
| v1 session-start hook reads the message inbox | config + observed 2026-07-23 | ailang repo `.claude/settings.json` SessionStart → `scripts/hooks/session_start.sh` ("checks the user inbox … using the ailang messages CLI"); displayed 5 unread at this session's start |

---
**Document created**: 2026-07-23 (bootstrap, attended). **RATIFIED 2026-07-23** (iteration 0,
attended: Mark + World coordinator) — record on issue #1; advisory-quorum ledger in
`.ailang/state/mission-quorum/`. Sprint routing is authorized from the next loop fire.
