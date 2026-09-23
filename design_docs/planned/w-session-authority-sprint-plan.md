# w-session-authority — Sprint Plan (mission iteration 176)

- **Sprint id:** `w-session-authority` · **Queue row:** 39 · **Status:** planned
- **Base commit:** `cf6cbf029a720a06b95f94b8b9ac449b6cf754a7` (detached worktree HEAD, clean tree)
- **Branch:** `sprint/w-session-authority` · **Executor worktree:** `/Users/voightkampff/dev/sunholo-data/.wt-world-iter176` (sibling of the repo, never /tmp)
- **Planner:** mission-control iteration 176, pi:ollama/kimi-k3:cloud sprint-planner
- **Design doc (WHAT, single source of truth):** [`design_docs/planned/w-session-authority.md`](w-session-authority.md) — quorum round-2 fixes applied verbatim by the controller (middleware in `host/daemon`, no `revoked` column, no secondary index)
- **Machine-readable record:** `sprint_w-session-authority.json` (repo root)

This plan adds only the HOW: order, per-milestone gates, snapshot points, executor constraints, and the pre-registered post-execution mutation battery. No design decision the quorum settled is re-litigated.

---

## §0 Baseline measured at base commit (this planner, in this worktree, 2026-09-21)

With `AILANG_BIN=$HOME/.pinned-ailang/ailang` exported ("AILANG v0.41.0", commit 24ee108, built 2026-09-21T15:50:02Z):

| Gate | Command | Observed |
|---|---|---|
| vet | `go vet ./...` | **rc=0** |
| AILANG gate | `bash ./scripts/verify_ail.sh` | **rc=0** — `✓ verify gate PASSED: 11 required identities verified, 40 named tests pass` and `✓ world package gate PASSED: 9/9 steps performed non-zero work` |
| Go tests | `go test ./... -count=1` | **rc=0 — 19 ok, 0 FAIL** |

All three gates are **GREEN at base**, matching the controller's measurement at `b65a9ec` (design doc F10/V13). Nothing red at base locally.

**Base conditions (each load-bearing for the executor):**

- **(a)** `scripts/verify_go.sh` is **NOT a gate** on this rig (charter row 76). It never appears in an acceptance criterion.
- **(b)** Remote CI's `go-verify` job is **RED at base** due to a KNOWN harness defect being escalated separately (`.github/workflows/ci.yml` still carries the old v0.30.0 pin where the host now wants v0.41.0; the world-package leg pins v0.41.0 into `$HOME/.pinned-ailang`). **All sprint gates run LOCALLY** and local gates are green at base (table above).
- **(c)** A bare `go test ./...` **WITHOUT** `AILANG_BIN` fails — always export `AILANG_BIN=$HOME/.pinned-ailang/ailang`. This is a base condition, never a regression.
- **(d)** `tools/launchd/*` is frozen fleet core — never touch.

Also re-confirmed at base: no `host/authority` directory exists; `git grep Session -- host/store/` → 0 rows (V7/V15); `host/store/store.go:42` reads `currentSchemaVersion = 2`; `host/store/schema_version_test.go:17-19` reads `frozenFutureSchemaVersion = 3` / `expectedCurrentSchemaVersion = 2`; nine mux routes at `host/daemon/daemon.go:555-565` with `POST /v1/commit` the sole mutation; `switch verb` dispatch at `cmd/ailang-worldd/main.go:117`; tty-fence precedent at `cmd/world-publish/tty.go:26-49`.

---

## §1 Goal

Deliver the inbound `Authorization: Bearer` credential→(episode, grants) boundary the repo has never had (design doc §1): mint/resolve/revoke primitives in a new `host/authority` package over a persisted `session_credentials` table (schema v2→3, fail-closed `LegacySchemaVersionError` for v2 stores), operator-facing `ailang-worldd session mint|revoke` verbs behind a controlling-terminal fence, a daemon-side `SessionMiddleware` protecting `POST /v1/commit`, and `handleCommit` constructing its broker session from the resolved binding. Eight GET routes stay unauthenticated as **declared residual R1**. Fails closed on all four typed denials (absent/malformed/unknown/expired), per D-WORLD-26.

## §2 Approach

Two executor milestones, in fixed order (M1 → M2 is structural: M2's CLI and middleware call store methods and resolver types that only M1 creates — see `split_verdict`). The controller builds commits after each milestone; the executor never runs a git write. After **each** milestone the executor snapshots every created-or-modified file into `.snap/M<k>/` (cumulative full content, untracked, never committed). The §5 mutation battery (M1–M7) is **pre-registered here and runs AFTER the executor finishes**, driven by the controller/evaluator against the landed tree with cp/sed restore discipline.

Gates run per milestone (all local, all with `AILANG_BIN` exported):

1. `go vet ./...` → rc=0
2. `go test ./host/authority/... -count=1` → ok (M1 onward)
3. `go test ./... -count=1` → **20 ok / 0 FAIL** (19 base packages + the new `host/authority` test package; "19 ok or better" per AC-M4-1)
4. `bash ./scripts/verify_ail.sh` → unchanged `11 identities / 40 tests` and `9/9 steps` (this sprint adds no `.ail`)

Any verdict about a process actually **serving a socket** (live daemon + curl batteries) is `UNINFORMATIVE UNDER SANDBOX` — labelled as such, never reported pass/fail; the informative assertions are the in-process `httptest` middleware/handler tests. A gate red **at base** is a finding to REPORT, never an obstacle and never silently fixed.

## §3 Milestone breakdown (files + per-file LOC, from design doc §3)

### M1 — `M1_AUTHORITY_PACKAGE_STORE_V3` (≤1d; plan estimate 0.6d; ~394 LOC)

`host/authority` package + `store.Store` schema bump (D2, D3, D5) + resolve/mint/revoke tests.

| Path | Action | Est. LOC | Notes |
|---|---|---|---|
| `host/authority/resolver.go` | new | ~80 | `Resolver` interface, `DenialKind` (Absent/Malformed/Unknown/Expired), `ResolveOutcome`, `SessionBinding{EpisodeID, Caps, ExpiresAt, CreatedAt}` (no CredentialID exposed), `FromContext` typed-context accessor; hash-then-index resolve: `sha256(presentedToken)` → single indexed `SELECT … WHERE credential_id = ?`; **no `crypto/subtle`** (post-hit compare would be vacuous, D3/D5) |
| `host/authority/mint.go` | new | ~70 | `crypto/rand` 32 bytes → 64 lowercase hex; store persists `sha256(token_hex)` only; `INSERT INTO session_credentials`; print raw token **once** (stdout) or `--out` write at mode 0600 with a one-line confirmation that never carries the credential |
| `host/authority/revoke.go` | new | ~30 | `DELETE FROM session_credentials WHERE credential_id = ?` in a transaction (delete-from-mapping = unknown, D4) |
| `host/authority/resolver_test.go`, `host/authority/mint_test.go` | new | ~150 | Named tests for ACs + the §5 matrix's killers: `TestMint_PrintsOnce`, `TestResolve_TypedDenialTaxonomy`, `TestResolve_ExpiryEnforced`, `TestRevoke_DeletesRow`, `TestNoAlternateHeader`, `TestAuthority_NoDirectTokenCompare` (source-scan, AC-M3-1), bounded-lookup query-count==1 test (AC-M3-2) |
| `host/store/schema.sql` | modify | +9 | Append the D2 DDL **exactly as the round-2 quorum fixed it**: `CREATE TABLE IF NOT EXISTS session_credentials (credential_id TEXT PRIMARY KEY, episode_id TEXT NOT NULL CHECK (episode_id <> ''), grants_json TEXT NOT NULL, expires_at INTEGER NOT NULL, created_at INTEGER NOT NULL);` — **no `revoked` column, NO secondary `CREATE INDEX`** (see §7 defect DEF-1) |
| `host/store/store.go` | modify | +~35 | `currentSchemaVersion` **2→3**; `MintSession`, `ResolveSession`, `RevokeSession` (single indexed query / single DELETE) |
| `host/store/schema_version_test.go` + `journal_test.go` fixtures | modify | +~20 | `expectedCurrentSchemaVersion` 2→3, `frozenFutureSchemaVersion` 3→4, add v3 migration row; keep fixtures non-vacuous. **Travels with M1** — deferring it reds `host/store` at `schema_version_test.go:19` (see §7 defect DEF-3) |

**M1 exit:** AC-M1-3, AC-M2-1..6 (resolver level), AC-M3-1, AC-M3-2, AC-M4-3 green on `host/authority`/`host/store`; gates 1–4 above all green (20 ok); `.snap/M1/` written.

### M2 — `M2_CLI_MIDDLEWARE_COMMIT_BINDING` (≤1d; plan estimate 0.6d; ~190 LOC)

CLI verbs + tty fence (D1, D4) + daemon middleware wiring (round-2 placement) + `handleCommit` session binding (D6, D7) + full AC battery.

| Path | Action | Est. LOC | Notes |
|---|---|---|---|
| `cmd/ailang-worldd/session.go` | new | ~90 | `session mint --episode <id> [--grant EFFECT=SCOPE:BUDGET]... [--ttl 3600] [--out <path>]` and `session revoke <credential_id-hash>`; `--db` flag default matching `serve --db`; **tty fence** ported from `cmd/world-publish/tty.go` (controlling terminal via `os.Open("/dev/tty")`, ENXIO refuse, Y/N prompt read from `/dev/tty`); mint speaks directly to the store DB path, **no daemon HTTP endpoint** (residual R13) |
| `cmd/ailang-worldd/main.go` | modify | +~10 | `case "session"` dispatch in the `switch verb` at :117 |
| `host/daemon/middleware.go` | new | ~50 | `SessionMiddleware.Wrap(protected, mux)` in the **daemon package** (round-2 quorum fix: calls in-package `writeAPIError`; `host/authority` must NOT import `host/daemon`); route-agnostic `protected` set; typed denial → HTTP: absent/unknown/expired → 401 (`SessionAbsent`/`SessionUnknown`/`SessionExpired` distinct classes), malformed → 400 `InvalidSession`; success → `context.WithValue` + typed `authority.FromContext` accessor; **never reads `X-World-Session`** |
| `host/daemon/daemon.go` + `host/daemon/handlers.go` | modify | +~40 | `Handler()` returns `authSessionMiddleware(d.resolver, protected, mux)`; `protected` = `POST /v1/commit` only this sprint (eight GETs unauthenticated, residual R1); `handleCommit` reads the binding via `authority.FromContext(ctx)` and builds `broker.NewSession(d.store, binding.EpisodeID, binding.Caps, registry)` instead of a zero-authority `d.store.Commit` |

**M2 exit:** all of design-doc §4 green (AC-M1-1..AC-M4-3); gates 1–4 all green (20 ok / 0 FAIL; vet rc=0; verify_ail unchanged); `TestMint_PrintsOnce`/taxonomy/expiry/revoke/no-alternate-header/source-scan/bounded-lookup all green including their HTTP halves; `.snap/M2/` written (cumulative — includes M1's files).

### Controller bookkeeping milestone (NOT an executor feature row)

File the `w-a2a-session-projection.md` unblock per design-doc §6 (row-40 adapter over `authority.Resolver.Resolve`), and spawn queue rows for residuals R1/R2/R3/R13 if the controller so decides. Controller/human acts, not executor work.

## §4 Day-by-day (cap: 2 days)

| Day | Work | Exit state |
|---|---|---|
| **Day 1** (M1) | 1. `schema.sql` DDL + `store.go` v3 bump + three `Mint/Resolve/Revoke` store methods. 2. Move the three schema-version constants/rows together (`currentSchemaVersion` 2→3, `expectedCurrentSchemaVersion` 2→3, `frozenFutureSchemaVersion` 3→4, v3 migration row, journal fixtures). 3. `host/authority` resolver/mint/revoke. 4. Named tests incl. source-scan + bounded-lookup. 5. Gates 1–4. 6. `.snap/M1/`. | AC-M1-3, AC-M2-1..6 (resolver), AC-M3-1..2, AC-M4-3 green; 20 ok / 0 FAIL |
| **Day 2** (M2) | 1. `cmd/ailang-worldd/session.go` verbs + tty fence; `main.go` dispatch. 2. `host/daemon/middleware.go`; wire `Handler()`; protect `POST /v1/commit`. 3. `handleCommit` session binding. 4. Full AC battery incl. HTTP halves (`httptest`, in-process; live-socket runs labelled `UNINFORMATIVE UNDER SANDBOX`). 5. Gates 1–4. 6. `.snap/M2/`. Then the executor STOPS; the pre-registered mutation battery (§6) is the controller's next act. | All §4 ACs green; 20 ok / 0 FAIL; vet rc=0; verify_ail unchanged |

## §5 Acceptance criteria (copied verbatim from design doc §4; milestone mapping)

All run **LOCALLY** at the sprint tree with `AILANG_BIN=$HOME/.pinned-ailang/ailang` (F10, v0.41.0). `scripts/verify_go.sh` is **NOT** a gate (charter row 76), so none of the criteria below require it. Each AC is a COMMAND + EXPECTED OUTPUT.

- **AC-M1-1 mint prints exactly once.** `ailang-worldd session mint --db <tmp>/world.db --episode e1 --grant fs.read=/tmp/b:1 --ttl 3600` under a PTTY (test drives `os.Open("/dev/tty")` via the fencing seam) → prints the raw 64-hex token **exactly once** and stores one row. Count raw-token occurrences in the output: **exactly 1**. *(M2; mint primitive lands in M1)*
- **AC-M1-2 distinct credentials per mint.** Mint twice for the same episode; the two printed tokens differ (and both rows have distinct `credential_id`). `test "$a" != "$b"` → **0 (distinct)**. *(M2; primitive in M1)*
- **AC-M1-3 resolve returns episode+grants+expiry.** `Resolve("Bearer <token>", now<expires)` → `Success` with `EpisodeID==e1`, `Caps` matching the minted grant, `ExpiresAt==mint+3600`. Resolve **never** returns the raw token. *(M1)*
- **AC-M2-1 absent header.** `Resolve("", now)` → `DenialAbsent`; over HTTP, `POST /v1/commit` with no `Authorization` → **401** `SessionAbsent`. *(resolver half M1; HTTP half M2)*
- **AC-M2-2 malformed header.** `Resolve("Key abc", …)` and `Resolve("Bearer", …)` and `Resolve("Bearer z" /*1 char*/ , …)` → `DenialMalformed`; HTTP → **400** `InvalidSession`. A non-64-hex token is malformed → **400**, never 401. *(M1 + M2 halves)*
- **AC-M2-3 unknown token.** `Resolve("Bearer deadbeef…64hex-unknown", …)` → `DenialUnknown`; HTTP → **401** `SessionUnknown`. The static `serve-api` key value as `Bearer` → **401** (constraint (i) — never a session). *(M1 + M2 halves)*
- **AC-M2-4 expired.** Mint with `--ttl 1`, `sleep 2`, resolve → `DenialExpired`; HTTP → **401** `SessionExpired`, distinct from `SessionUnknown`. *(M1 + M2 halves)*
- **AC-M2-5 revocation.** `session revoke <credential_id-hash>` then resolve the same token → `DenialUnknown` (delete-from-mapping = unknown, D4). *(resolver/store half M1; CLI half M2)*
- **AC-M2-6 no alternate-header fallback.** Request with ONLY `X-World-Session: <a valid minted token>` and no `Authorization` → **401** `SessionAbsent` (must NOT resolve). Under mutation M6 it resolves — so a green control proves the absence is really enforced. *(M1 + M2 halves)*
- **AC-M3-1 no direct secret comparison.** A **source-scan test** over `host/authority/*.go` **FAILS** if any direct equality (`==`, `bytes.Equal`, `strings.EqualFold`, `strings.Contains`) is applied to token or credential-hash material — pinning that no future edit introduces a direct secret comparison into the resolver. Behavioral note: unknown-vs-known tokens traverse the identical code path (the same single `SELECT … WHERE credential_id = ?` query shape, D3/D5), so there is no branch an online adversary can time-difference. Killers: mutation M1 (§5). *(M1)*
- **AC-M3-2 bounded lookup.** Assert the resolve path issues **exactly one** SQL query: wrap the store/DB seam (test injects a counting `*sql.DB`-like observer or uses the store's existing injected-`objectStore` seam `daemon.go:267`) and assert query count == 1 with the single indexed `SELECT … WHERE credential_id = ?`. *(M1)*
- **AC-M4-1** `go test ./... -count=1` → **19 ok or better** (19 packages stay green / 0 FAIL; the new `host/authority` package adds its own test package that must also pass). `go vet ./...` → **rc=0**. *(M2; re-proved at both milestone exits — 20 ok expected)*
- **AC-M4-2** `bash ./scripts/verify_ail.sh` with `AILANG_BIN=$HOME/.pinned-ailang/ailang` → **"✓ verify gate PASSED: 11 required identities verified, 40 named tests pass"** and **"✓ world package gate PASSED: 9/9 steps"** (unchanged — this item adds no `.ail`). *(M2; re-proved at both milestone exits)*
- **AC-M4-3** store schema bump: a v2 store DB opened by the new binary is **refused** (`LegacySchemaVersionError`); a fresh DB initializes at v3 and the migrate row/constraints in `schema_version_test.go` are updated and green. *(M1; re-proved at M2 exit)*

**No AC is amended** (`acceptance_criteria_amendments: []` in the JSON) — every criterion is satisfiable as written; the §7 defects are file-table/cross-reference/sequencing issues, not AC contradictions.

## §6 Pre-registered mutation battery (runs AFTER the executor finishes — controller/evaluator act)

Protocol per arm (verbatim from design doc §5): `cp <file> <file>.bak && sed -i '…' <file>; <run target test>`; restore `cp <file>.bak <file>`; restore **asserted** by `shasum -a 256 <file>` == the pre-mutation digest AND `test -x <file>` where a mode bit matters (the `.bak` is `cp -p` so the mode bit survives); **never** `git checkout --` (the tree may be dirty; a checkout would silently discard unrelated work). Each arm: assert the mutation LANDED (pre/post sha256 differ), assert the mutant still compiles/vets before reading any test result, read the scoped observable (never the exit code alone), require the named test to be the SOLE failure, restore, re-verify green control before the next arm.

| # | Mutation | Sole killer (named test) | AC guarded |
|---|---|---|---|
| M1 | Introduce a direct `==` comparison of token material into the resolver (e.g. replace hash-then-index with `if presentedHash == storedHash`) | `TestAuthority_NoDirectTokenCompare` | AC-M3-1 |
| M2 | Collapse `DenialUnknown` into the absent branch (unknown token → `SessionAbsent`/401) | `TestResolve_TypedDenialTaxonomy` | AC-M2-1 vs AC-M2-3 |
| M3 | Remove the `now < expires_at` check at resolve (ignore expiry) | `TestResolve_ExpiryEnforced` | AC-M2-4 |
| M4 | Mint prints/logs the token twice (duplicate print) | `TestMint_PrintsOnce` | AC-M1-1 |
| M5 | Revocation leaves the row (the `DELETE` is omitted — revoke becomes a no-op against the store) | `TestRevoke_DeletesRow` | AC-M2-5 |
| M6 | Fall back to reading `X-World-Session` when `Authorization` is absent | `TestSessionMiddleware_NoAlternateHeaderFallback` (host/daemon) | AC-M2-6 — **corrected evaluator round-2:** the mutation is only expressible at the middleware, so the originally-named `TestNoAlternateHeader` (resolver-level) stays green under it; reproduced first-party by the controller. |
| M7 (green control) | Reword a comment in the resolver (no behavior change) | **all tests stay green** | non-vacuity of the whole battery |
| M8 (added round-2, judge finding D-1) | Neuter `handleCommit`'s `authority.FromContext` presence check (`!ok` → `!ok && false` — the handler wired without the middleware no longer fails closed) | `TestHandleCommit_DirectCallWithoutBindingFailsClosed` (host/daemon, added round-2) | defense-in-depth gate for direct handler wiring; **proved round-2 by the controller: mutant compiles, fail-count exactly 1 (sole killer), restore byte-identical, package green after restore** |

Also pre-registered: **baseline gate outcomes** (§0) as the base state; **schema-version expectations** — `currentSchemaVersion` 2→3, `expectedCurrentSchemaVersion` 2→3, `frozenFutureSchemaVersion` 3→4, one new v3 migration row, fresh store initializes at v3, v2 store refused via `LegacySchemaVersionError` (AC-M4-3); the six **load-bearing named tests** listed above.

## §7 Doc defects found (recorded, NOT silently fixed)

- **DEF-1 — §3 file table contradicts the round-2 quorum fix on the index.** The `host/store/schema.sql` row says "add `session_credentials` table + `idx_session_credential_lookup` (D2 DDL)", but §2 D2 (round-2 quorum fix applied verbatim) and §7 V18 removed the secondary `CREATE INDEX` ("the TEXT PRIMARY KEY already carries the unique index a point lookup needs"). **Plan instruction:** the executor follows D2/V18 — table DDL only, no secondary index. Doc left as-is.
- **DEF-2 — residual cross-reference numbers in §2 don't match §9.** D6 cites "residual R3 (§9)" for deferred read-route enforcement and D7 cites "residual R4" for effect-ledger re-recording; §9 numbers these **R1** and **R3** respectively, and no R4 exists. This plan uses §9 numbering (R1 read routes, R2 live sessions, R3 effect-ledger, R13 daemon mint endpoint).
- **DEF-3 — §8 sequencing would red M1.** §8 lists "update `schema_version_test`" in M2 while M1 owns the `currentSchemaVersion` 2→3 bump; `schema_version_test.go:19` asserts `currentSchemaVersion == expectedCurrentSchemaVersion`, so deferring the test bump leaves M1's `host/store` red. **Plan instruction:** the test-constant bump + v3 migration row + journal fixtures travel with M1 (§3's file table is milestone-neutral; no AC text changes).
- **DEF-4 — trivial pseudo-code field typo (D6).** The middleware sketch calls `classFor(*out.Denial)`; `ResolveOutcome` (D5) has no `Denial` field — the field is `Denied`. Pseudo-code only; the executor writes compiling code.

## §8 Risks (design doc §9 residuals + the schema-bump risk)

- **R1 — full `/v1/*` read-route enforcement deferred (D6).** Eight GET routes stay unauthenticated this sprint; a real authority gap on read paths, named for a follow-up queue row. The middleware's route-agnostic `protected` set makes it a config flip later.
- **R2 — live in-process sessions survive revoke/expiry (D4).** An already-constructed `broker.Session` finishes its grant budget; expiry/revocation gate NEW resolutions only. Declared, not hidden.
- **R3 — commit effect-ledger re-record deferred (D7).** `handleCommit` is session-gated but the commit is not re-recorded through `broker.Session`'s replay path; owned by row 40's propose → verify → commit sequence.
- **R13 — no daemon-owned mint endpoint (D1).** Mint writes the store DB directly (local-first); a daemon mint endpoint is an explicit follow-up if ever needed.
- **Schema-version-bump risk — v2 stores are refused.** An existing v2 operator DB opened by the new binary is refused with `LegacySchemaVersionError` until re-initialized/migrated — a deliberately loud operator-facing behavior change, covered by **AC-M4-3** and the doc-pointer message on the refusal path.
- **Trust model (informational):** credential theft = session theft is inherent to bearer credentials; TLS is the operator's deployment concern (§9).

Process risks the handoff guards: a red-at-base gate reported as a sprint regression (§0 records the green base; CI go-verify red is KNOWN, base condition (b)); socket-serving verdicts reported as real passes (must be labelled `UNINFORMATIVE UNDER SANDBOX`); a mutation "kill" that never landed (sha256 landed-proof + compile assert before any red is scored).

## §9 Handoff to the sprint-executor (binding constraints)

1. Run **IN** `/Users/voightkampff/dev/sunholo-data/.wt-world-iter176` on `sprint/w-session-authority` (created by the controller from `cf6cbf0`).
2. **NO git write operations at all** — no `git add/commit/stash/checkout/branch`. The controller builds commits.
3. After **EACH** milestone, snapshot every created-or-modified file into `.snap/M<k>/` — cumulative full content (`.snap/M2/` includes M1's files). `.snap/` is an untracked working artifact, never committed.
4. Gates per milestone, all with `AILANG_BIN=$HOME/.pinned-ailang/ailang` exported: `go vet ./...`; `go test ./host/authority/... -count=1` (M1 onward); `go test ./... -count=1` (both milestones — expect **20 ok / 0 FAIL**); `bash ./scripts/verify_ail.sh` (unchanged 11/40 and 9/9).
5. Anything serving a socket is `UNINFORMATIVE UNDER SANDBOX` — label it, never report it as pass/fail. The informative HTTP assertions are in-process (`httptest`).
6. A gate red **at base** is a finding to REPORT, not an obstacle — never fix it, never plan around it silently.
7. Never touch `tools/launchd/*` (frozen fleet core). No `.ail` changes — the verify gate must stay byte-identical in outcome (AC-M4-2).
8. Restore discipline for any probing: `cp -p` backup + `shasum -a 256` assertion, never `git checkout --`.
9. Executor STOPS after M2's snapshot. The §6 mutation battery and the controller bookkeeping milestone are not executor work.

## §10 Velocity and pricing

- **LOC estimate:** 584 (doc headline "~+570"): M1 ≈ 394 (80+70+30+150+9+35+20), M2 ≈ 190 (50+10+90+40).
- **Estimates:** M1 0.6d, M2 0.6d → plan total **1.2d** against the doc's ≤2d (each milestone ≤1d). Within the ≤2-day cap.
- **Agree with the doc's sizing.** The one line the plan calls tight: the store schema bump is +~64 LOC but touches three coupled constants plus a migration row plus refusal-path expectations (`schema_version_test.go:17-19`, measured) — a slip there reds `host/store` in a way LOC doesn't express; M1's 0.6d carries it.
- **Precedent for honesty:** base gates measured twice now (controller at `b65a9ec`, this planner at `cf6cbf0`) — 19 ok / 0 FAIL, vet rc=0, verify_ail PASSED. The "20 ok" post-M1 target is arithmetic from a measured 19, not a hope.
