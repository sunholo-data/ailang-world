# Sprint plan — w-a2a-session-projection (iteration 187)

**Design doc**: `design_docs/planned/w-a2a-session-projection.md` (quorum-closed iter-187 r2
under the narrow-refinement carve-out). **Base commit**: `4ff25f9` (sprint branch, detached
planner worktree). **Planner**: planner-role, iter-187 continuation run (the first run built the
prototype and was cut off by a provider quota error before the drill; this run re-measured
everything first-party).

**PROTOTYPE-FIRST.** The whole design was implemented in the disposable planner worktree
BEFORE this plan was written. Every LOC figure, every gate reading, and every KILLED/SURVIVED
verdict below was **measured on that prototype**. The prototype does not ship; the executor
rebuilds from this plan. The full prototype (all 11 files + `prototype.diff` + `SHA256SUMS`)
is preserved at
`/Users/voightkampff/.ailang/state/mission-world-iter187-evidence/prototype/` with relative
paths intact; per-file sha256 pins are in §8 so the executor can reuse the exact code.

---

## 1. Baselines (measured this run)

Pristine tree = `git archive HEAD` (4ff25f9) extracted to `/tmp/iter187-pristine`, gates run
there; prototype gates run in the planner worktree. `AILANG_BIN=$HOME/.pinned-ailang/ailang`
(v0.41.0, commit 24ee108) for the ailang-dependent gates.

| Gate | Base (pristine 4ff25f9) | Prototype |
|---|---|---|
| `go vet ./...` | rc=0 | rc=0 |
| `AILANG_BIN=… go test ./... -count=1` | **20 ok / 0 FAIL** | **21 ok / 0 FAIL** (new: `host/projection`) |
| `AILANG_BIN=… bash scripts/verify_ail.sh` | rc=0 — "11 required identities verified, 40 named tests pass" | rc=0 — identical line |
| `gofmt -l host cmd` | 3 files, ALL pre-existing: `host/daemon/session_middleware_test.go`, `host/store/schema_version_test.go`, `host/store/store.go` | only the 2 `host/store` files remain; the prototype **fixes** `session_middleware_test.go` (comment realignment). No file the sprint touches is unformatted. |
| `bash scripts/bench_worldd.sh --smoke` | rc=0 | rc=0 |
| Dependency closure `go list -deps ./host/daemon/... ./cmd/ailang-worldd/... \| sort -u` | **250** (controller's base list `~/.ailang/state/iter187-deps-base.txt`) | **253** — see §5, deviation D3 |
| M1 deadline test (`TestResolveContext_DeadlineWhileSingleConnectionHeld`) | n/a | **0.10 s** PASS |

The `go.sum`/`go.mod` module movement, measured (§5/D3–D4): `+ailang v0.33.2`,
`go-isatty v0.0.20→v0.0.22`, `x/sys v0.46.0→v0.47.0`, **`x/sync v0.21.0→v0.22.0` (module-graph
only — the doc predicted only the first two bumps)**.

## 2. Milestones, files, LOC (all LOC measured on the prototype)

The design fixes the sequence M1 → M2 → M3. **Commit plan (bisectability decision, required by
the directive): TWO commits.** `go mod tidy` drops an unimported requirement and the doc forbids
a dead anchor file, so **M2 and M3 land as ONE commit** — the go.mod change rides with the first
real `protocol` import. Commit 1 = M1; commit 2 = M2+M3. Each commit builds and passes
`go test ./...` (verified on the prototype: M1-only state compiles — `ResolveContext` is additive;
the fused M2+M3 state is the full prototype, 21 ok / 0 FAIL).

### M1 — P6.A-CTX: bounded context-aware resolution in `host/authority` (~0.1d)

**Files**: `host/authority/resolver.go` (+41/−12), `host/authority/resolvecontext_test.go`
(NEW, 278 lines).

**Surface added** (exact signatures, as prototyped):

```go
// host/authority/resolver.go — added to the Resolver interface:
ResolveContext(ctx context.Context, header string, now int64) (ResolveOutcome, error)

// Resolve keeps today's behaviour byte-for-byte: it delegates and collapses a
// store error to DenialUnknown EXACTLY as before (the declared residual).
func (r *resolver) Resolve(header string, now int64) ResolveOutcome
func (r *resolver) ResolveContext(ctx context.Context, header string, now int64) (ResolveOutcome, error)
```

`ResolveContext` is the old `Resolve` body with exactly two changes: (1) the store call is
`r.st.ResolveSession(ctx, credentialID)` (was `context.Background()`); (2) a store error returns
`(ResolveOutcome{}, fmt.Errorf("authority: resolve: %w", err))` — never a `DenialUnknown`.
Header shape, token shape, hashing, the indexed PK lookup, expiry (`now > row.ExpiresAt`) and
grant decoding are byte-identical to the old policy. `Resolve` =
`ResolveContext(context.Background(), …)` with `err != nil → DenialUnknown`, so `/v1/commit`'s
middleware is untouched.

**Fake-Resolver scan (directive requirement): none needed.** `grep` shows NO fake `Resolver`
implementations in `host/daemon` tests — the session-middleware tests run the REAL resolver over
a real store (`newHandlerDaemon`). The interface widening compiles with zero stub changes
(measured: full build green).

**Named tests** (all PASS on the prototype):

| Test | Asserts |
|---|---|
| `TestResolveContext_PolicyEquivalence` | Over one fixture table (absent, 2 malformed shapes, unknown, expiry boundary BOTH sides, success): `Resolve` and `ResolveContext(background)` agree with each other AND with pinned expectations — kills policy drift the bare compare can't see. |
| `TestResolveContext_StoreErrorIsErrorNotDenial` | Failing store → non-nil error wrapping the store error + zero outcome from `ResolveContext`; same failure through `Resolve` still collapses to `DenialUnknown` (the residual); a `context.DeadlineExceeded` store error surfaces as an error. |
| `TestResolveContext_DeadlineWhileSingleConnectionHeld` | REAL one-connection `sql.DB` pool (`SetMaxOpenConns(1)`, the store.go:305 shape), the only connection held by an open transaction (self-rollback after 3 s so a ctx-dropping mutation REDs instead of hanging): `ResolveContext` with a 100 ms deadline returns `errors.Is(err, context.DeadlineExceeded)` **within the bound** (measured 0.10 s), zero outcome. |
| `TestResolve_NoGoroutineInResolvePath` | Source scan of every non-test authority file: no statement beginning `go ` outside comments; fails vacuous-green check (`checked == 0` → fatal). |

**LOC**: production +41/−12 (net +29, no new file); test 278. Design estimate said 40–60 LOC +
~10 test; measured is inside the production range, test larger (the deadline fixture is a real
pool, not a mock).

### M2 — P6.D: dependency admission (~0.15d) — COMMIT-FUSED WITH M3

**Files**: `go.mod` (+6/−3), `go.sum` (+8/−7), `host/daemon/daemon_test.go` (ONE allowlist line +
doc-comment rewrite + narrowness test ≈ 52 of its +274 lines).

- `go get github.com/sunholo-data/ailang@v0.33.2` (measured: succeeds under go 1.26.6; bumps as
  in §1).
- `allowedDepModules` gains EXACTLY ONE entry, the PACKAGE path
  `"github.com/sunholo-data/ailang/serveapi/protocol"` (never the module root); the doc comment
  is rewritten to say "module roots or package paths" and to record that entries match as path
  prefixes.
- **Named test**: `TestAilangProtocolAdmissionIsNarrow` — with the entry in place,
  `disallowedDeps(["…/internal/apiserver", "…/serveapi", "cloud.google.com/go/storage"])` must
  refuse ALL THREE (the facade leg doubles as MUT-FACADE-IMPORT's shape, the cloud leg as
  MUT-CLOUD-DEP's classifier proof); control leg: the wire package path AND a future subpackage
  (`…/serveapi/protocol/subpkg`) are admitted, so the gate can't pass by the entry going missing.
- Merge criterion: `TestDaemonDependencyAllowlist` green (measured), narrowness green, and the
  probe (delete the line → allowlist test REDs naming exactly one intruder) — re-verified this
  run via MUT-ALLOWLIST-ROOT/MUT-FACADE-IMPORT (§4).

### M3 — P6.B-A2A-CARD: `host/projection` + additive mount (~0.4d)

**Files** (LOC measured):

| File | LOC | Kind |
|---|---|---|
| `host/projection/projection.go` | 364 new | production |
| `host/projection/projection_test.go` | 1247 new | test (20 tests) |
| `host/daemon/daemon.go` | +48/−0 | production (mount + injection) |
| `host/daemon/middleware.go` | +20/−10 | production (extract `writeSessionDenial`, see §5/D6) |
| `host/broker/broker.go` | +12/−0 | production (`NewCapabilitySnapshot`, see §5/D1) |
| `host/daemon/daemon_test.go` | ≈ +222 of +274 | test (2 helpers + 4 mount/integration tests) |
| `host/daemon/session_middleware_test.go` | +4/−4 | test (gofmt realignment only) |

**Types/signatures added** (as prototyped):

```go
// host/projection/projection.go
type HeadReader interface { GetRegistryHead(ctx context.Context, name string) (hashref.HashRef, bool, error) }
type DenyWriter  func(w http.ResponseWriter, kind authority.DenialKind)      // daemon injects writeSessionDenial
type ErrorWriter func(w http.ResponseWriter, class, message string, status int) // daemon injects writeAPIError
type Config struct {
    Resolver authority.Resolver; Reader transitionreg.Reader; Heads HeadReader
    Deny DenyWriter; Fail ErrorWriter; Agent protocol.AgentInfo; MaxWait time.Duration
}
func New(cfg Config) (*Handler, error)            // rejects nil seams and MaxWait <= 0
func (h *Handler) AgentCard(w http.ResponseWriter, r *http.Request)  // GET /.well-known/agent.json
func (h *Handler) A2A(w http.ResponseWriter, r *http.Request)        // POST /a2a/
type bindingCaps struct{ caps []broker.Capability }                  // transitionreg.CapabilitySource over the binding's immutable grants (D1)
func (c bindingCaps) CapabilitySnapshot(now int64) broker.CapabilitySnapshot
// host/broker/broker.go
func NewCapabilitySnapshot(grants []Capability, now int64) CapabilitySnapshot  // Epoch 0; detached copy
// host/daemon/middleware.go
func writeSessionDenial(w http.ResponseWriter, k authority.DenialKind)  // THE ONE denial→HTTP mapping, extracted verbatim from Wrap (D6)
```

**Handler behaviour** (implements Decisions 1–6 / B1–B5): both routes derive ONE bounded ctx
(`context.WithTimeout(r.Context(), h.maxWait)`), resolve via `ResolveContext`, deny BEFORE any
registry touch. Card route: denial → injected `Deny` (F2 mapping, daemon envelope); resolve/snap
failure → 503/504 via injected `Fail` (504 iff `errors.Is(err, context.DeadlineExceeded)`);
absent head (B3 `GetRegistryHead` pre-check, `store.TransitionRegistryV1`) → zero-skills card at
200; skills = `transitionreg.NewRequest(ctx, Reader, bindingCaps{b.Caps}, now).Allowed()` —
verbatim IDs, registry bytewise order; card = `map[string]any` with the upstream keyset (F5).
`/a2a/`: every outcome a `protocol.A2AError` at HTTP 200 — -32001 absent/unknown/expired, -32600
malformed/undecodable/non-"2.0", -32601 method ≠ `tasks/send`, -32602 undecodable params or
unlisted/guessed/stale `skill_id` (`params.Metadata["skill_id"].(string)`, the upstream read
site), -32603 constant `"transition invocation is not available in this daemon"` for an
AUTHORIZED skill and for resolution/registry failure. 1 MiB body cap (`http.MaxBytesReader`).

**Mount** (`host/daemon/daemon.go`, additive): `d.projection` built in `New` over the ONE
resolver, `transitionreg.NewReader(d.store)`, the `d.reads` head seam, `Deny:
writeSessionDenial`, `Fail: writeAPIError`, `Agent: protocol.AgentInfo{Name: "ailang-worldd",
…, Version: Version}`, `MaxWait: readDeadline` (the D7 read deadline — a projection request IS
store reads under the transport); construction failure aborts at `StageConfig`. Routes:
`mux.HandleFunc("GET /.well-known/agent.json", d.projection.AgentCard)` and
`mux.HandleFunc("POST /a2a/", d.projection.A2A)` — NOT in `isProtected`.

**Named tests** (all PASS; `host/projection` 20 + `host/daemon` 4):

| Test | Asserts (AC) |
|---|---|
| `TestAgentCard_ExactSkillSetPerSession` | Two sessions, unequal caps → EXACT own ID sets, verbatim (dot/slash grammar CallerSurface would reject), registry order, no extras (AC2, F6). |
| `TestAgentCard_TracksSessionSnapshot` | Same handler, two sessions → different cards; refetch stable per snapshot (AC3). |
| `TestAgentCard_AmbientExportsAbsent` | `exit`/`writeBytes`/8×`std/io`/`submit_feedback` never projected; non-vacuous control (AC4). |
| `TestAgentCard_UpstreamKeySet` | Card = EXACT upstream keyset + exact skill-entry keyset (AC1/AC8, F5). |
| `TestAgentCard_DenialMatrix` | 4 kinds → 401/401/401/400 + envelope class via the INJECTED writer; messages constant (same bad credential twice → byte-identical body) (AC-CARD-DENIALS). |
| `TestAgentCard_AbsentHeadZeroSkills` | No head → 200, explicit `"skills":[]`, still session-gated (AC-ABSENT-HEAD). |
| `TestAgentCard_HeadReadFailuresAre5xx` | Head-check error → 503; snapshot failure after PRESENT head → 503 (AC-ABSENT-HEAD, B3 race 2nd half). |
| `TestAgentCard_HeadRaceAbsentThenSucceeds` | Check absent, head published before snapshot → NewRequest's result is USED (B3 race 1st half). |
| `TestA2A_DenialMatrix` | 4 kinds on `/a2a/` → HTTP 200 + A2AError, -32001/-32600, constant messages, no REST envelope (AC-DENIAL-JSONRPC). |
| `TestA2A_CodeMatrix` | Full code matrix: -32600 (undecodable, non-"2.0"), -32601 (method), -32602 (unlisted/guessed/stale/params), constant -32603 (authorized); no `result` member ever (AC5/AC8/AC-A2A-CODES). |
| `TestA2A_NeverWritesStore` | Full denial+admission battery leaves world head AND registry head byte-identical (AC5/AC8). |
| `TestProjection_DeniedNeverTouchesRegistry` | All 4 kinds × both routes: counting `GetRegistryHead` + counting snapshot reader see ZERO calls; valid session + parse failure also zero (AC-DENIAL-JSONRPC). |
| `TestProjection_OneSnapshotPerRequest` | One request = EXACTLY 1 head check + 1 snapshot read, both routes (AC7). |
| `TestAgentCard_HeadChangeWithoutRestart` | Publish rev 2 → next card changes without restart (AC12). |
| `TestProjection_CarrierFixed_OnlyBearerHeader` | X-World-Session-only (with a VALID token) fails closed absent on both routes; key-shaped 64-hex Bearer resolves to nothing; standard-carrier control (AC6, D-WORLD-26). |
| `TestProjection_BoundedWait` | Fault-injected blocking resolver and blocking head check: card → 504, `/a2a/` → -32603, each within the bound (100 ms + 1.5 s slack), never a hang; seams have a 2 s escape hatch so a ctx-dropping mutation REDs (Decision 6 / retained AC13). |
| `TestProjection_ConfigValidation` | Zero/negative MaxWait and each nil seam → construction error (Decision 6 startup half). |
| `TestProjection_WireOwnershipSource` | Source gate: imports protocol; NO `"jsonrpc"` literal or json-tagged wire field; NO CallerSurface/ValidateMCPName/X-World-Session in code (AC1/F6/AC6). |
| `TestProjection_NoDeadlineTampering` | Source gate: no ResponseController/SetWriteDeadline/SetReadDeadline (AC14). |
| `TestProjection_ZeroSkipSource` | Source gate: no `t.Skip*` in projection tests (zero-skip clause). |
| `TestProjectionRoutes_MountedAndNotProtected` (daemon) | Routes mounted; NOT `isProtected` (predicate probe + `/v1/commit` control); unauthenticated card = projection's OWN 401 SessionAbsent; unauthenticated `/a2a/` = 200 + -32001 (proving no middleware wrap); health route untouched (AC9/B5). |
| `TestCardDenial_ByteIdenticalToCommitMiddleware` (daemon) | Same bad credential (5 shapes incl. minted-expired) → card and `/v1/commit` denial bodies BYTE-IDENTICAL (B1). |
| `TestAgentCard_ViaDaemon_SessionScoped` (daemon) | Full mounted stack (daemon→resolver→store→transitionreg): exact per-session sets incl. slash-grammar ID; AgentInfo on the wire; authorized → constant -32603, guessed → -32602 (AC2/AC5 integration). |
| `TestAgentCard_ViaDaemon_AbsentHeadZeroSkills` (daemon) | Fresh daemon (epoch registry only, no transition head — the F8 production default) → 200 `"skills":[]` (AC-ABSENT-HEAD integration). |

**LOC totals (measured)**: production 364+48+20+12 = **444** (+ net middleware/broker/authority
M1 +29, go.mod 6 → **≈ 479**); test 1247+278+274+4 = **≈ 1803**. M3 production exceeds 250 LOC
but is ONE package with ONE mount seam; splitting it would break the fused M2+M3 commit, so it
is kept whole with the per-file breakdown above (design Files table already scopes it as one
milestone).

## 3. Gates (every milestone; readings from the prototype)

```
go vet ./...                                                   # rc=0
AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./... -count=1  # 21 ok / 0 FAIL (base 20/0)
AILANG_BIN=$HOME/.pinned-ailang/ailang bash scripts/verify_ail.sh  # rc=0, 11 identities / 40 named tests
gofmt -l host cmd                                              # empty for touched files (2 pre-existing host/store hits remain, base-confirmed)
bash scripts/bench_worldd.sh --smoke                           # rc=0
```
Plus the closure gate: `go list -deps ./host/daemon/... ./cmd/ailang-worldd/... | sort -u | wc -l`
= **253** with ADD = {`…/ailang/serveapi/protocol`, `…/host/projection`, `…/host/transitionreg`},
removed set EMPTY (§5/D3 — supersedes the doc's 250→251).

## 4. Mutation drill — every mutation in the doc, run against the prototype

Method per mutation: pristine copy saved to `/tmp/iter187-mut/`, exact `old` anchor verified
unique (`grep -c -F` = 1 unless noted), mutation applied, **only the named test** run, then the
file restored with `cp` (never git); restore verified by a final green build + test. Sole-killer
status = full-package run under the mutation. **Score: 30 KILLED / 0 SURVIVED** (14 sole-killer,
15 multi-killer, 1 toolchain-layer kill). Three mutations needed the directive's compile-fail
reshape — recorded with their new anchors.

### M1 (P6.A-CTX) — file `host/authority/resolver.go`

| ID | old anchor (count) | new text | Named test → result |
|---|---|---|---|
| MUT-RESOLVE-BACKGROUND | `row, ok, err := r.st.ResolveSession(ctx, credentialID)` (1) | `row, ok, err := r.st.ResolveSession(context.Background(), credentialID)` | `TestResolveContext_DeadlineWhileSingleConnectionHeld` → **KILLED** (3.01 s escape, "err = nil … want error wrapping DeadlineExceeded"), **sole killer** |
| MUT-CTX-ERR-AS-UNKNOWN | `return ResolveOutcome{}, fmt.Errorf("authority: resolve: %w", err)` (1) | `denied := DenialUnknown` + `return ResolveOutcome{Denied: &denied}, nil` | `TestResolveContext_StoreErrorIsErrorNotDenial` → **KILLED**, not sole (also Deadline test: error collapses to denial there too) |
| MUT-RESOLVE-POLICY-DRIFT | `if now > row.ExpiresAt {` (1) | `if now >= row.ExpiresAt {` | `TestResolveContext_PolicyEquivalence` → **KILLED** ("expiry boundary live" case), **sole killer** |

### M2 (P6.D) — files `host/daemon/daemon_test.go`, `host/projection/projection.go`

| ID | old anchor (count) | new text | Named test → result |
|---|---|---|---|
| MUT-ALLOWLIST-ROOT | `"github.com/sunholo-data/ailang/serveapi/protocol", // w-a2a-session-projection P6.D: ONE pinned A2A wire package (ailang v0.33.2), never the module root` (1) | `"github.com/sunholo-data/ailang", // MUT-ALLOWLIST-ROOT: module root widens the gate` | `TestAilangProtocolAdmissionIsNarrow` → **KILLED**, **sole killer** (`TestDaemonDependencyAllowlist` stays GREEN — measured) |
| MUT-FACADE-IMPORT | `"github.com/sunholo-data/ailang/serveapi/protocol"` (1, in `host/projection/projection.go` import block) | same line + `\n\n\t_ "github.com/sunholo-data/ailang/serveapi" // MUT-FACADE-IMPORT` | `TestDaemonDependencyAllowlist` → **KILLED**, **sole killer** — RED names intruders BY NAME: "29 package(s) outside the zero-cloud allowlist … `github.com/sunholo-data/ailang/serveapi` … `modelcontextprotocol/go-sdk/…`". **RESHAPED**: the literal mutation does not compile (facade needs `modelcontextprotocol/go-sdk/mcp`, absent from go.sum); reshape = `go get github.com/sunholo-data/ailang/serveapi@v0.33.2` to admit the facade's own deps, then the import compiles and the gate REDs. Anchor unchanged; go.mod/go.sum restored after. (Finding: the facade is ALSO stopped one layer earlier, at the go.sum boundary.) |

### M3 (P6.B) — files as noted

| ID | File | old anchor (count) → new | Named test → result |
|---|---|---|---|
| MUT-PROTO-OWNER | projection.go | `type bindingCaps struct{ caps []broker.Capability }` (1) → prepend parallel struct `type a2aWireMirror struct { JSONRPC string \`json:"jsonrpc"\` }` | `TestProjection_WireOwnershipSource` → **KILLED**, **sole** |
| MUT-CARDSURFACE | projection.go | `skills := make([]map[string]any, 0, len(allowed))` (1) → prepend a `protocol.ToolDescriptor` build + `protocol.CallerSurface(tds)` gate that 503s on error | `TestAgentCard_ExactSkillSetPerSession` → **KILLED** (real IDs rejected → 503), not sole (9 tests incl. `TestProjection_WireOwnershipSource` source half) |
| MUT-SESSION-UNION | projection.go | `return req.Allowed(), nil` (1) → `return req.Registry.List(), nil // MUT` | `TestAgentCard_ExactSkillSetPerSession` → **KILLED**, not sole (+ `TracksSessionSnapshot`) |
| MUT-CARD-GLOBAL | projection.go | same edit as MUT-SESSION-UNION (applied separately) | `TestAgentCard_TracksSessionSnapshot` → **KILLED**, not sole (+ `ExactSkillSetPerSession`) |
| MUT-UNFILTERED-PROJECTION | projection.go | `card := map[string]any{` (1) → prepend hardcoded `skills = append(skills, map[string]any{"id":"exit", …})` | `TestAgentCard_AmbientExportsAbsent` → **KILLED**, not sole (6 tests) |
| MUT-A2A-STALE | projection.go | `protocol.A2AError(w, req.ID, codeInvalidParams, msgNotAuthorized)` (1) → `protocol.A2AError(w, req.ID, codeInternal, notAvailableMessage)` | `TestA2A_CodeMatrix` → **KILLED**, **sole** |
| MUT-A2A-FAKE-SUCCESS | projection.go | the authorized branch anchored on `// An AUTHORIZED skill_id: the constant not-available refusal (B4).\n\t\t\tprotocol.A2AError(w, req.ID, codeInternal, notAvailableMessage)` (1) → `protocol.A2AResult(w, req.ID, map[string]any{"status": "accepted"})` | `TestA2A_CodeMatrix` → **KILLED**, **sole** |
| MUT-DEFAULT-CAPS | resolver.go | the `if !ok { denied := DenialUnknown; return …Denied…}` block (1) → return a `Success` binding with process-default caps | `TestProjection_CarrierFixed_OnlyBearerHeader` → **KILLED** (key-shaped Bearer → 200), not sole (4 projection + 4 authority tests) |
| MUT-ALT-HEADER | projection.go | `r.Header.Get("Authorization")` (**2 occurrences, both replaced BY DESIGN**) → inline func falling back to `"Bearer " + r.Header.Get("X-World-Session")` when empty | `TestProjection_CarrierFixed_OnlyBearerHeader` → **KILLED**, not sole (+ `TestProjection_WireOwnershipSource`, both denial matrices, `AbsentHeadZeroSkills`) |
| MUT-KEY-AS-SESSION | resolver.go | `credentialID := hashHexToken(token)` (1) → prepend `if token == "cccc…64" { return Success binding }` | `TestProjection_CarrierFixed_OnlyBearerHeader` → **KILLED** ("key-shaped Bearer card status=200, want 401"), **sole** |
| MUT-SPLIT-SNAPSHOT | projection.go | `_, hasHead, err := h.heads.GetRegistryHead(ctx, store.TransitionRegistryV1)` (1) → same + a second discarded `GetRegistryHead` call | `TestProjection_OneSnapshotPerRequest` → **KILLED** (heads.calls=2), not sole (+ `DeniedNeverTouchesRegistry`) |
| MUT-A2A-CODEFLIP | projection.go | swap the two codes: authorized branch `codeInternal`→`codeInvalidParams` (comment-anchored, 1) AND final `codeInvalidParams, msgNotAuthorized`→`codeInternal, msgNotAuthorized` (1) | `TestA2A_CodeMatrix` → **KILLED**, **sole** |
| MUT-SECOND-OPEN | daemon.go | `\tproj, err := projection.New(projection.Config{` (1) → prepend `if _, oerr := store.Open(cfg.DBPath); oerr != nil { return nil, d.abort(StageConfig, "projection opened a second store handle", oerr) }` | `TestProjectionRoutes_MountedAndNotProtected` → **KILLED** with the doc's exact observation: startup aborts "…another process already holds the writer lock…" (`WriterAlreadyActive`), not sole (startup-fatal: 52 daemon tests RED) |
| MUT-ISPROTECTED-EXPAND | daemon.go | `return r.Method == http.MethodPost && r.URL.Path == "/v1/commit"` (1) → `return (… "/v1/commit") \|\| r.URL.Path == "/.well-known/agent.json" \|\| r.URL.Path == "/a2a/"` | `TestProjectionRoutes_MountedAndNotProtected` → **KILLED**, **sole** |
| MUT-DEADLINE-RELAX | projection.go | `func (h *Handler) AgentCard(w http.ResponseWriter, r *http.Request) {` (1) → append `\n\t_ = http.NewResponseController(w).SetWriteDeadline(time.Time{})` | `TestProjection_NoDeadlineTampering` → **KILLED**, **sole** |
| MUT-CARD-DENIAL-TABLE | middleware.go | `writeAPIError(w, classFor(k), msgFor(k), http.StatusBadRequest)` (1) → `… http.StatusUnauthorized)` | `TestSessionMiddleware_MalformedHeader` (PRE-EXISTING) → **KILLED** ("status=401, want 400"), **sole**. **Deviation D2**: because B1 is implemented as ONE shared writer (D6), the mapping lives in the daemon, pinned by the middleware suite + byte-identity test — there is no projection-side table to mutate. |
| MUT-CARD-ENVELOPE | daemon.go | `Deny:     writeSessionDenial,` (1) → inline func writing plain `{"denied":"…"}` JSON at 401 | `TestCardDenial_ByteIdenticalToCommitMiddleware` → **KILLED** (bodies differ), not sole (+ `ProjectionRoutes_MountedAndNotProtected` class assertion) |
| MUT-ABSENT-HEAD | projection.go | `if !hasHead {\n\t\t\treturn nil, nil\n\t\t}` (1) → `if !hasHead {\n\t\t\treturn nil, fmt.Errorf("projection: absent head treated as a read error (MUT-ABSENT-HEAD): %w", rerr)\n\t\t}` | `TestAgentCard_AbsentHeadZeroSkills` → **KILLED**, **sole**; control `TestAgentCard_HeadReadFailuresAre5xx` stays GREEN (verified). **RESHAPED**: the first shape (delete the branch) did not compile (`hasHead` unused) — new anchor recorded here. |
| MUT-A2A-MESSAGE-INTERP | projection.go | authorized branch (comment anchor, 1) → message `fmt.Sprintf("transition %q is not available in this daemon", skillID)` | `TestA2A_CodeMatrix` → **KILLED**, **sole** |
| MUT-DENIAL-401 | projection.go | `protocol.A2AError(w, nil, codeSessionDenied, msgAbsent)` (1) → `http.Error(w, msgAbsent, http.StatusUnauthorized)` | `TestA2A_DenialMatrix` → **KILLED** (status assertion), not sole (+ `CarrierFixed` `/a2a/` leg) |
| MUT-DENIAL-SNAPSHOT-FIRST | projection.go | the `AgentCard` head block `…WithTimeout…\n\tdefer cancel()` (1) → append `\n\t_, _, _ = h.heads.GetRegistryHead(ctx, store.TransitionRegistryV1)` BEFORE `ResolveContext` | `TestProjection_DeniedNeverTouchesRegistry` → **KILLED** (counting observer ≥1 on denied requests), not sole (+ `OneSnapshotPerRequest`) |
| MUT-DROP-DEADLINE-PROJ | projection.go | `out, err := h.resolver.ResolveContext(ctx, r.Header.Get("Authorization"), time.Now().Unix())` (**2 occurrences, both replaced BY DESIGN**) → `… ResolveContext(context.Background(), …)` | `TestProjection_BoundedWait` → **KILLED** (2 s escape hatch, status/time assertions RED), **sole** |
| MUT-SKIP-SOCKET | projection_test.go | `func TestProjection_ConfigValidation(t *testing.T) {` (1) → append `\n\tt.Skip("mutation: skip on listen failure")` | `TestProjection_ZeroSkipSource` → **KILLED**, **sole** (re-shaped per the doc: invocation-half socket skip is deferred; the gate is the zero-skip source scan) |
| MUT-CLOUD-DEP | projection.go | import anchor (1) → add `_ "cloud.google.com/go/storage"` | **KILLED AT THE TOOLCHAIN LAYER**: does not compile — "missing go.sum entry for module providing package cloud.google.com/go/storage (imported by …/host/projection)" — the intruder is NAMED BY NAME before any test runs. No compilable reshape exists: `go get cloud.google.com/go/storage@v1.43.0` was measured to DESTABILIZE the graph (downgraded `google.golang.org/protobuf v1.36.12→v1.34.2` and dropped the ailang requirement). Classifier-level proof that the allowlist reports this exact path by name lives in `TestAilangProtocolAdmissionIsNarrow`'s cloud leg (green on the prototype). |
| MUT-STARTUP-CACHE | projection.go | `func (h *Handler) allowedDescriptors(…` (1) → prepend `var startupAllowedCache …` + early-return cache hit; `return req.Allowed(), nil` (1) → cache-set + return | `TestAgentCard_HeadChangeWithoutRestart` → **KILLED**, not sole (8 tests — a process-global cache poisons the whole package run, as expected) |

**Deferred (NOT drilled — moved verbatim to the design's "Deferred to invocation" section,
gated on F7/F8):** `MUT-COMMIT-BOUNDARY-LIE`, `MUT-LEAK-CONN`; `MUT-SKIP-SOCKET`'s original
socket-closure target is deferred with them (drilled here in its re-shaped zero-skip form).

## 5. Deviations from the design doc (each forced by a measurement)

- **D1 — `broker.NewCapabilitySnapshot` instead of `broker.NewSession` in the capability path.**
  The design's B2 says "construct a broker session from the binding (`broker.NewSession(store,
  b.EpisodeID, b.Caps, registry)` … the broker session's `CapabilitySnapshot(now)` is the
  `CapabilitySource`)". Measured impossible as written: `NewSession` takes a `*store.Store`
  (broker.go:99) — handing the projection a store handle violates P5 ("projection never opens a
  store"; the prototype makes it structural: the projection holds NO store), and TR.C's
  dispatch-binding gate forbids constructing a live `*Session` outside `host/broker`. Fix: one
  additive pure constructor `broker.NewCapabilitySnapshot(grants, now)` (+12 LOC, Epoch 0 — a
  binding has no ledger; `Allows` reads only grants and Now) plus a 2-method adapter
  `bindingCaps` in the projection. The capability POLICY still lives exactly in `broker.Allows`
  (pinned by `TestAgentCard_ExactSkillSetPerSession`).
- **D2 — the card-denial status mapping lives in the daemon, not the projection.** Consequence
  of B1's writer injection (D6): `MUT-CARD-DENIAL-TABLE`'s sole killer is the PRE-EXISTING
  `TestSessionMiddleware_MalformedHeader`, and `MUT-CARD-ENVELOPE` is killed by the byte-identity
  test — the doc's table assumed a projection-side mapping. The coverage is equivalent (the
  mapping exists exactly once and is pinned at that one place).
- **D3 — closure moves 250 → 253, not 251.** Measured ADD set vs the controller's base list:
  `github.com/sunholo-data/ailang/serveapi/protocol` (P6.D's +1, the doc's number) PLUS two
  same-module packages newly reachable from the daemon core: `host/projection` (new) and
  `host/transitionreg` (F7 measured ZERO production importers at base; `daemon.go` now imports it
  for `transitionreg.NewReader(d.store)`). Removed set EMPTY. No allowlist impact: both extra
  packages sit under the already-admitted `github.com/sunholo-data/ailang-world` root.
  `TestDaemonDependencyAllowlist` is GREEN at 253 (measured). AC16's "250 → 251" is the marginal
  effect of the pin alone; the plan's gate records the honest end-state 250 → 253.
- **D4 — a THIRD module bump the doc didn't predict:** `golang.org/x/sync v0.21.0 → v0.22.0`
  (module-graph only, pulled by ailang v0.33.2's go.mod; NO `x/sync` package enters the 253
  closure — verified against both dep lists). AC16's "ONLY other graph movement is the two
  measured bumps" becomes three; all three roots are already allowlisted or absent from the
  closure, so no allowlist change. Any FOURTH movement still fails the row.
- **D5 — reshaped mutations** (compile-fail is not a kill): `MUT-ABSENT-HEAD` (branch-inversion
  form, new anchor in §4), `MUT-FACADE-IMPORT` (needs the facade's own `go get` to compile),
  `MUT-CLOUD-DEP` (no compilable form exists — toolchain-layer kill, measured destabilization of
  the `go get` reshape).
- **D6 — `host/daemon/middleware.go` touched despite "/v1/commit middleware stays
  byte-identical".** The ONLY change is extracting Wrap's inline denial `switch` into
  `writeSessionDenial` — statement-for-statement identical logic (same cases, same order, same
  status mapping, same `classFor`/`msgFor`/`writeAPIError` calls) — so the daemon can inject THE
  SAME writer into the projection (`Deny: writeSessionDenial`), which is how B1's "card-route
  denial byte-identical to the middleware" is achieved with the envelope daemon-owned (AC1:
  projection never formats it). **Behaviour-preserving, proven by**: the entire pre-existing
  middleware suite (`TestSessionMiddleware_AbsentHeader/MalformedHeader/UnknownToken/
  ExpiredAndRevokedIsDistinct/NoAlternateHeaderFallback/SuccessReachesHandler`) passes UNCHANGED,
  and `TestCardDenial_ByteIdenticalToCommitMiddleware` proves wire byte-identity across 5
  credential shapes. `session_middleware_test.go`'s only change is a gofmt comment realignment
  (the file was already unformatted at base — baseline table, §1).
- **D7 — `host/broker/broker.go` touched**: purely additive `NewCapabilitySnapshot` (D1); no
  existing line changed; the full `host/broker` suite (53.6 s) passes unchanged.
- **D8 — fake-Resolver stubs unnecessary** (the directive anticipated them): none exist in
  `host/daemon` tests (grep-verified; middleware tests use the real resolver); M1 widens the
  interface with zero stub churn.
- **D9 — deadline fixture is a REAL one-connection pool**, not a store-internal hack: the design
  allowed "an open transaction or equivalent"; `heldConnDB` holds the pool's only connection with
  an open transaction and a 3 s self-rollback (so a ctx-dropping mutation REDs instead of
  hanging). Measured 0.10 s per run; MUT-RESOLVE-BACKGROUND REDs in 3.01 s (escape), never
  hangs.

## 6. AC → test map

| AC | Test(s) |
|---|---|
| AC1 wire reuse | `TestProjection_WireOwnershipSource`, `TestAgentCard_UpstreamKeySet`, `TestAgentCard_ExactSkillSetPerSession` (verbatim IDs), `TestCardDenial_ByteIdenticalToCommitMiddleware` (envelope stays daemon-owned) |
| AC2 exact surface | `TestAgentCard_ExactSkillSetPerSession`, `TestAgentCard_ViaDaemon_SessionScoped` |
| AC3 card tracks session | `TestAgentCard_TracksSessionSnapshot` |
| AC4 ambient absent | `TestAgentCard_AmbientExportsAbsent` |
| AC5 admission | `TestA2A_CodeMatrix`, `TestA2A_NeverWritesStore`, `TestAgentCard_ViaDaemon_SessionScoped` |
| AC6 carrier fixed (D-WORLD-26=A) | `TestProjection_CarrierFixed_OnlyBearerHeader`, `TestAgentCard_DenialMatrix`, `TestA2A_DenialMatrix` |
| AC7 one snapshot | `TestProjection_OneSnapshotPerRequest` |
| AC8 A2A conformance | `TestA2A_CodeMatrix`, `TestA2A_DenialMatrix`, `TestAgentCard_UpstreamKeySet` |
| AC9 landed behaviour | `TestProjectionRoutes_MountedAndNotProtected` + the unchanged `host/daemon` suite (20→21 ok, zero REST regressions) |
| AC11 dep floor | `TestDaemonDependencyAllowlist` |
| AC12 dynamic source | `TestAgentCard_HeadChangeWithoutRestart` |
| AC13 | MOVED → Deferred; retained half: `TestProjection_BoundedWait` |
| AC14 no deadline tampering | `TestProjection_NoDeadlineTampering` (+ same-run `TestBoundedWaitsAndBodyLimit` / D7 tests as the REST control) |
| AC16 pinned dep | `TestDaemonDependencyAllowlist` + `TestAilangProtocolAdmissionIsNarrow` + the measured closure row (§3, D3/D4) |
| AC-CARD-DENIALS | `TestAgentCard_DenialMatrix`, `TestCardDenial_ByteIdenticalToCommitMiddleware` |
| AC-ABSENT-HEAD | `TestAgentCard_AbsentHeadZeroSkills`, `TestAgentCard_HeadReadFailuresAre5xx`, `TestAgentCard_HeadRaceAbsentThenSucceeds`, `TestAgentCard_ViaDaemon_AbsentHeadZeroSkills` |
| AC-A2A-CODES | `TestA2A_CodeMatrix`, `TestA2A_DenialMatrix` |
| AC-DENIAL-JSONRPC | `TestA2A_DenialMatrix`, `TestAgentCard_DenialMatrix`, `TestProjection_DeniedNeverTouchesRegistry` |
| P6.A-CTX ACs | `TestResolveContext_PolicyEquivalence`, `TestResolveContext_StoreErrorIsErrorNotDenial`, `TestResolveContext_DeadlineWhileSingleConnectionHeld`, `TestResolve_NoGoroutineInResolvePath` |
| zero-skip clause | `TestProjection_ZeroSkipSource` |

## 7. Residuals and risks

1. **`/v1/commit` still resolves without a request deadline** (`Resolve` →
   `ResolveContext(context.Background())`, store error → `DenialUnknown`) — the design's declared
   residual; this sprint does not change it.
2. **No invocation success path exists** (F7): every authorized `/a2a/` call gets the constant
   -32603. The deferred half re-opens when F7 (coordinator) and F8 (production registry publish)
   land.
3. **Transition registry is never populated in production today** (F8): the shipped card is the
   zero-skills-at-200 shape until a publisher lands; both shapes are tested.
4. **Module-graph risk (AC16):** the pin moves three versions (isatty, x/sys, x/sync — D4). The
   executor must re-run the closure + bump check and fail on any FOURTH movement.
5. **`x/sync` bump was NOT predicted by the doc** — recorded as D4; if the quorum wants the doc
   amended, that is a one-line AC16 edit at ratification time, not a code change.
6. **Source-gate style**: three gates are line-oriented source scans (wire ownership, deadline
   tampering, zero-skip) — deliberate, matching the repo's existing `disallowedDeps` style; they
   skip full-line comments and have vacuity controls.
7. **`MUT-SECOND-OPEN` is startup-fatal by design** (52 tests RED) — the named killer is
   `TestProjectionRoutes_MountedAndNotProtected`; the breadth is expected for a mount-time store
   open.
8. **F6 grammar conflict travels to the MCP child**: card IDs are verbatim and never pass
   `CallerSurface`; the MCP dispatch child must solve naming itself (already recorded in both
   docs).

## 8. Prototype reference (sha256-pinned; executor may reuse verbatim)

Evidence root: `/Users/voightkampff/.ailang/state/mission-world-iter187-evidence/prototype/`
(relative paths preserved; `prototype.diff` = the full `git diff`, 661 lines; `SHA256SUMS`).

| File | sha256 |
|---|---|
| `go.mod` | `f4829312475397a1c83e73983e1637bfc1149630ca948463aa7afebc28709a56` |
| `go.sum` | `d83ebcdb0f73eff3f92bb686b291a0987562cac38f0643c3c6703e45e2c4051c` |
| `host/authority/resolver.go` | `052bc6502cd59467d02db4fc1ea62b44dac3d2a413ca902c88ca8bf7f498790d` |
| `host/authority/resolvecontext_test.go` | `dba4ad26dd8b9f25e657bbda9bae27f49adf7861b6681599d54a3b26256466e2` |
| `host/broker/broker.go` | `ce9682173b4c88db1c48b9b40bba6bba735e40cdeb9479ce2c3a9e5feaf26729` |
| `host/daemon/daemon.go` | `34bea896fda12e686c225d3c07ce9cc6ceca245b82695060bbe449f9f77a9627` |
| `host/daemon/daemon_test.go` | `ff149d9f909149cdf73cb54b365d162a9a9d5889a57f80956b587df1aef6acbc` |
| `host/daemon/middleware.go` | `892d136217c588253008e4313d1d7d8e00380880fab8239bd305ee4c7c86d54d` |
| `host/daemon/session_middleware_test.go` | `946b0cce9b482b9fe8d316127938f6007a57a1f2310977f4f2ceace52e4a9763` |
| `host/projection/projection.go` | `ac35404cba5fa8725f3d6fc5c4d78ff698af512eb6d7f8d6fa7d1104efb9f474` |
| `host/projection/projection_test.go` | `cf904d78e5d9264698e34dcce1e3bf0195958008522cf52a40b9526dabf409df` |
