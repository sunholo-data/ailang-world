# w-mcp-effects-unrecorded-is-untyped — MCP callers get an untyped error after an effect has run (row 136)

**Status:** PARKED `needs-human-review` (iteration 236, designer lane claude-opus-5-5): quorum round 2 BLOCKED; one design-scope objection (option D, MCP idempotency key) awaits Mark as `D-WORLD-66`. Design only; nothing implemented.
**Target:** next World patch. **Priority:** clause-6 (+5); charter position after row 140, before 141 → 93 (D-WORLD-61=A).
**Estimated:** ~0.5 d host Go + docs. **Dependencies:** none for M1–M2. The client-visible part (M3) waits on an upstream release.
**Measured base:** worktree `sprint/row136-mcp-typed-errors` at `3d965f6`; pinned upstream `github.com/sunholo-data/ailang v0.47.2` (`go.mod:6`); upstream checked at `v0.47.2` (module cache), `v0.52.1` (latest tag) and `origin/dev` `b7028a7cb` (2026-10-06).

## Problem

Over `/a2a/`, a coordinator failure becomes a typed JSON-RPC error. `dispatchError`
(`host/projection/projection.go:416–470`) checks `EffectsUnrecordedError` first (`:437`) and returns
`EffectsUnrecordedPrefix` plus the effect-record refs (`:474–485`). Its only caller is the A2A path
(`:372`). Over `/mcp/`, `mcpAdapter.Invoke` returns the raw error from `Dispatch` (`mcp.go:130–134`),
and the released `mcphttp` handler turns *any* host error into `-32603 "host callback failed"`. So an
MCP agent whose effect ran but whose commit failed (R14 head race, finish failure, commit error) cannot
learn the record refs. MCP task ids are minted fresh for each call (`mcp.go:126`, `mintMCPTask` `:137`),
so a retry cannot reconcile against the earlier call. It runs the effect again: the call is
at-least-once, and the client gets nothing that would let it check first.

Second defect: every MCP call is recorded as A2A. The invocation id is `"a2a:" + episode + ":" + task`
(`host/coordinator/coordinator.go:108`, used at `:250`). Every entry header and object says
`writtenBy`/`Provenance` `"coordinator:a2a"` (`host/coordinator/plan.go:17,83,135`). The operator
line for an MCP refusal reads `ailang-worldd: a2a refusal: mcp invoke a2a:…` (`projection.go:501`,
`mcp.go:132`). The website already documents this as a limitation
(`website/docs/agents/provenance.md:45`, `website/docs/reference/mcp-and-a2a.md:204–207`).

## Measured facts

| # | Claim | Evidence |
|---|---|---|
| V1 | `Invoker.Invoke` returns `(InvocationResult{Value json.RawMessage}, error)`. The result has no error, status or `isError` field. | upstream `serveapi/protocol/interfaces.go:16–23` (v0.47.2). The file is byte-identical at v0.52.1 and origin/dev (V14) |
| V2 | `callTool` returns any non-nil Invoke error as a host error. On success it wraps `Value` as `{content:[{type:"text",text:Value}], structuredContent:Value}`. `callResultJSON` has no `isError`. | `mcphttp/methods.go:23–26, 83–100` (v0.47.2); v0.52.1 / origin/dev `:25, :99`. `callTool` and `callResultJSON` are byte-identical across all three revisions; the file differs only in `toolJSON`/`listTools` (V14) |
| V3 | `serveMessages` answers the **whole POST** (every message in a batch) with `WriteMCPEnvelope(w, id, CallbackMessage(hostErr))`, code fixed at `-32603`. | `mcphttp/handler.go:178–183`; `protocol/envelope.go:33–51`. `serveMessages` is byte-identical across all three revisions (V14) |
| V4 | `CallbackMessage` maps only `ErrCallbackCapacity`, `DeadlineExceeded` and `Canceled`. Every other error becomes `"host callback failed"`. No interface hook exists. | `protocol/envelope.go:21–32`. The file is byte-identical at v0.52.1 and origin/dev (V14) |
| V5 | Upstream freezes this envelope with a test. | `serveapi/embedded_mcp_test.go:283` `TestEmbeddedMCPFrozenCallbackEnvelopes` |
| V6 | `hostcall.Run` derives `callCtx, cancel := context.WithTimeout(ctx, runner.timeout)` (`serveapi/protocol/hostcall/runner.go:42`, v0.47.2 module cache). Its final select returns `return zero, callCtx.Err()` on `case <-callCtx.Done():` (`runner.go:64–65`) without waiting for the callback goroutine, which keeps running and still sends its result to the buffered `done` channel. World passes one constant to both knobs: `invokeDeadline = 20 * time.Second` (`host/daemon/daemon.go:99`), `InvokeWait: invokeDeadline` (`:673`) and `CallbackTimeout: invokeDeadline` (`:674`). `runner.go` is byte-identical at v0.47.2, v0.52.1 and origin/dev (V14). | as cited |
| V7 | World advertises each tool's `outputSchema` on MCP. | `host/projection/mcpname.go:84` |
| V8 | The MCP operator line already carries the refs. `mcp.go:132` logs `err.Error()` for every Dispatch error, and `EffectsUnrecordedError.Error()` includes `records [sha256:…]` (`host/coordinator/errors.go:60–67`). The line's label still says `a2a refusal` (`projection.go:501`), and **no test asserts any `mcp invoke` line** (grep of `host/**/*_test.go`: zero hits). | as cited |
| V9 | "R14" is the store's head-race `ConflictError` (`coordinator.go:378`). The existing injection is `r14Store` with `Commit → &store.ConflictError{}`. Only the A2A e2e test uses it today. | `host/daemon/setools_e2e_test.go:719–770` `TestSeToolsA2APostEffectFailureIsMapped` |
| V10 | Changing the id prefix is safe for replay. Neither reconcile nor Replay derives a prefix: `committed()` compares the stored `rec.InvocationID` with the receipt key (`coordinator.go:199`), and `Replay` takes the full id and reads stored objects (`host/coordinator/effectful.go:329–`). | as cited |
| V11 | Namespaces: `/v1/commit` accepts only `rest:` (`host/daemon/handlers.go:352–358`), and the effect journal reserves `effect:` (`host/store/journal.go:308, 325, 372`). The store has no `a2a:` check (grep `host/store`). `why` accepts targets that start with `a2a:` or `rest:` (`cmd/ailang-worldd/why.go:317, 327`). No `.ail` module constrains id prefixes or `writtenBy` values (grep `*.ail`). | as cited |
| V12 | Object `Provenance` is first-writer-wins: `INSERT OR IGNORE` (`host/store/store.go:565, 1160`), and the object hash covers the payload only. So the provenance of an input/output object shared by two calls is not a reliable surface tag. The record object always includes `invocationId`, so it is unique to one call. The entry header's `WrittenBy` is part of the entry hash (`plan.go:135–140`) and is unique to its entry. | as cited |
| V13 | No upstream issue asks for MCP `isError` or typed host errors. | `gh issue list -R sunholo-data/ailang --search "mcphttp"`, `"host callback failed"`, `"isError mcphttp"`, run 2026-10-06 |
| V14 | Cross-revision identity. Each value is the first 16 hex digits of `sha256`. `mod` is `shasum -a 256 < $(go env GOMODCACHE)/github.com/sunholo-data/ailang@v0.47.2/<path>`; the other columns are `git show <rev>:<path> \| shasum -a 256` in `~/dev/sunholo-data/ailang`; origin/dev is `b7028a7cb`. | table below |

**V14 evidence** (run 2026-10-06):

| path (`serveapi/protocol/…`) | mod v0.47.2 | v0.47.2 | v0.52.1 | origin/dev |
|---|---|---|---|---|
| `interfaces.go` | `7a0ede1b1df8e6de` | `7a0ede1b1df8e6de` | `7a0ede1b1df8e6de` | `7a0ede1b1df8e6de` |
| `envelope.go` | `015396f49c73a4ad` | `015396f49c73a4ad` | `015396f49c73a4ad` | `015396f49c73a4ad` |
| `hostcall/runner.go` | — | `5fbd63356c6d405b` | `5fbd63356c6d405b` | `5fbd63356c6d405b` |
| `mcphttp/methods.go` | `36f983ee1cbadcca` | `36f983ee1cbadcca` | `5373f733ca2ab094` | `5373f733ca2ab094` |
| `mcphttp/handler.go` | `f4e310bd59c27464` | `f4e310bd59c27464` | `773f279b764b7541` | `773f279b764b7541` |

`methods.go` and `handler.go` differ between v0.47.2 and v0.52.1. `git diff v0.47.2 v0.52.1` has
exactly four hunks:
- `methods.go`: `toolJSON` gains `Title` and `Annotations`, and `listTools` copies them.
- `handler.go`: `Config` gains `Gate *protocol.BearerGate`, and `ServeHTTP` calls a new
  `admitGatedCalls`, an OAuth2 bearer gate that runs before dispatch.

Neither change touches the error path. Extracting the functions with
`awk '/^func \(h \*handler\) callTool/,/^}/'` (and the same for `serveMessages` and
`type callResultJSON`) and hashing each gives identical values at v0.47.2 and v0.52.1:
`callTool` `422c9a368c954ec4`, `serveMessages` `bb348556a44c1974`, `callResultJSON`
`c72048e7bc092348`. `git diff --stat v0.52.1 origin/dev -- serveapi/protocol` is empty. These hunks are the only reason line numbers
move at v0.52.1: by +2 in `methods.go` (the Invoke call moves `:84 → :86`) and by +30 in
`handler.go` (the host-error envelope moves `:182 → :208`).

## Options

**(A) Return the error as a successful tool result.** For example,
`structuredContent: {"error":{"kind":"EffectsUnrecorded",…}}`. **Rejected.**
1. Without `isError` (V2), every MCP client reads the call as a success, including agents that
   simply continue. The defect we want to remove is a client that misreads the outcome. This
   replaces a misread retry with a misread success.
2. The payload would break the `outputSchema` that World itself advertises (V7). MCP requires
   structured results to conform to the declared output schema, so a validating client rejects or
   mangles the result.
3. It is a local workaround of the library. The `Invoker` contract makes a non-nil error the failure
   channel (V1–V3). Sending failures through the success channel works around a library gap in
   World instead of routing it upstream, which the charter's frozen-core rule forbids (CLAUDE.md
   "Hard rules").

**(A′) Rewrite the response in middleware.** World would wrap `mcphttp`, buffer the frozen
envelope and replace its message. **Rejected.** World's own seam comment says the released handler
"owns … all wire bytes and the response envelope" (`mcp.go:27–28`). Upstream freezes those bytes (V5).
A batch's envelope has no per-message id to rewrite (V3). This would be a vendored fork in all but name.

**(B) Upstream ask only, with no host change now.** **Insufficient.** The client-visible fix
needs it, but it leaves the second defect (surface tags) and the untested operator line (V8) unfixed.
Both are host-only and correct today.

**(C) Split: host-side now, client-visible typed error gated on upstream. RECOMMENDED.**
- M1: make the operator log honest and tested. MCP refusals are labelled `mcp refusal`, and the R14
  e2e test proves the line carries effect-record refs that resolve in the store. Today this line is
  the only place an MCP operator can recover them.
- M2: tag the surface. MCP invocations get an `mcp:` id and `writtenBy: coordinator:mcp`.
- **M1 and M2 are separate patches.** Following the repo convention that each milestone is CI-green
  and mergeable alone, M1 lands first and its AC1.1 asserts the id the code mints *then*
  (`a2a:ep1:`). M2 changes that one AC1.1 assertion to `mcp:ep1:` in the same commit that changes
  the prefix. AC2.1b records the edit.
- M3: file the upstream ask and mark the client-visible typed error as gated on it. When it ships,
  a follow-up row bumps the pin and changes the M1 e2e assertion from the frozen envelope to the
  typed error.

**Pin bump: out of scope.** No release, including v0.52.1 and origin/dev, carries `isError` or a typed-error
hook (V1–V4). A bump would buy nothing and would trigger the S8 floor-raise inventory.

**Why host code, not a package (S3):** this is transport error rendering, invocation-id minting and
operator logging at the host boundary (S2). The projection and coordinator are host Go by
construction. No `world/` kernel or `.ail` change is needed (V11).

## Design

### M1 — MCP operator line: correct label, tested refs (~0.15 d)

- `logRefusal(method, id, err)` (`projection.go:500`) becomes `logRefusal(surface, method, id, err)`
  and writes `ailang-worldd: <surface> refusal: <method> <id>: %q\n`. Every A2A call site passes
  `"a2a"`, so A2A bytes do not change. The four MCP call sites (`mcp.go:76, 81, 107, 132`) pass
  `"mcp"`. Their method labels become the JSON-RPC method names: `tools/list` for `mcp tools`, and
  `tools/call` for `mcp invoke`.
- No wire change, and no change to the logged cause content (`err.Error()`, `%q`-escaped, one physical line). Only the surface label and the method label change.
  The row-127 privacy boundary is unchanged: causes stay on the operator-only sink.
- `dispatchError` stays as it is. It is already a pure `error → (code, message)` function with a
  table test (`projection_test.go:~800–848`). M3's follow-up reuses it unchanged on the MCP path. M1
  adds no inert adapter type for an upstream interface that does not exist yet.

### M2 — Surface tag in the invocation id and the entry header (~0.25 d)

- `coordinator.Call` gains `Surface Surface`, with `type Surface string` and constants
  `SurfaceA2A = "a2a"` and `SurfaceMCP = "mcp"`. `validateCall` (R1) refuses any other value **and the
  empty value**. A forgotten tag then fails loudly as `InvalidCallError` instead of falling back to
  `a2a`. The strict refusal stays.
- **Call-literal sweep (part of M2, same commit).** `grep -rnE '\bCall\{' host cmd --include='*.go'`
  finds **5** `Call{…}` composite literals (the sixth hit, `effectful_test.go:640`, is a `[]Call{…}`
  slice whose elements come from the helper). Each one gets an explicit `Surface`:
  - production: `host/projection/projection.go:363` (A2A → `SurfaceA2A`) and
    `host/projection/mcp.go:130` (→ `SurfaceMCP`);
  - test helpers and literals: `host/coordinator/coordinator_test.go:155` (the `rig.call` helper,
    which feeds every coordinator/effectful test through it), `host/daemon/workspace_test.go:524`,
    and `host/daemon/commit_budget_test.go:154`.

  The last one matters. That test expects `context.DeadlineExceeded` from a held receipt read. R1
  `validateCall` (`coordinator.go:122–137`) runs before the receipt lookup, so without a `Surface` it
  would get `InvalidCallError` and turn red. That is loud, not silent, but the sweep must add
  `Surface: coordinator.SurfaceA2A` there, not change its expected error.
- `InvocationID(surface, episodeID, taskID)` returns `string(surface)+":"+episodeID+":"+taskID`.
  Callers: `coordinator.go:250`, `projection.go:374`, `mcp.go:132`, tests.
- `planInvocation`, `planEffectInvocation` and `planCommit` take the surface. Header `WrittenBy` and
  object `Provenance` become `"coordinator:" + surface`. Object provenance is advisory
  (first-writer-wins, V12). The surface tags of record are the id and the entry header.
- `why` (`why.go:317, 327`, usage `:276`, `main.go:29, 92`) accepts `mcp:<id>` the same way as
  `a2a:<id>`, via the record scan. The unused constant `coordinatorAuth` (`why.go:39`, zero references)
  is deleted.
- **Compatibility.** Existing stores keep their `a2a:` ids and `coordinator:a2a` headers, including
  those written by past MCP calls. Reconcile, Replay, `why a2a:…` and the integrity scan read stored
  strings and never derive a prefix (V10), so old stores replay unchanged. A2A clients send their own
  task ids, and their invocation ids stay `a2a:<ep>:<task>`, so resend and reconcile (R15/R17) behave
  as before. MCP task ids are random for each call, so no MCP idempotency key exists to break. Bonus:
  an A2A client can no longer pick a task id that collides with an MCP call in the same episode.
- `/v1/commit` still accepts only `rest:` (V11), so `mcp:` cannot be forged over REST. M2 adds a row
  to the existing `commit_budget_test.go:356` id table.
- Docs (S7), updated together with M2: `docs/QUICKSTART.md:39` gets the refusal-line format for
  both surfaces. In `website/docs/`: `agents/provenance.md:45,50,91` (delete the "cannot tell"
  sentence), `concepts/sessions-and-episodes.md:16`, `guides/provenance-walks.md:79`,
  `reference/cli.md:37,254,284–286`, `reference/http-api.md`, `getting-started/coding-tools.md`,
  and `website/static/AGENTS.md`. Sweep with `grep -rn 'coordinator:a2a\|a2a:<' website docs`.
  Existing frozen verification logs under `design_docs/verification/` are evidence and stay unchanged.

### M3 — Upstream ask; client-visible typed error gated (~0.1 d now)

- File an issue on `sunholo-data/ailang` (area:serveapi) with V1–V5 as the repro. The ask: an
  `Invoker` error that implements an optional interface, e.g.
  `interface{ JSONRPCError() (code int, message string) }`, is answered as a **per-message**
  JSON-RPC error carrying that code and message. Other messages in a batch keep their results (V3).
  `CallbackMessage` keeps its frozen fallback for every other error. Upstream may instead choose
  `InvocationResult.IsError` → `isError:true`. Either works, provided the error path is not a
  schema-bound `structuredContent` (V7). World's A2A precedent favours a typed JSON-RPC error,
  because it carries the same information over both surfaces.
- Send `ailang messages send mission-control` with the issue link. The body is positional
  (memory: the `send` flag shape).
- Update the warning at `website/docs/reference/mcp-and-a2a.md:204–207`. It should name the operator
  line (M1) as today's recovery path and link the upstream issue.
- Charter: add a follow-up row, gated on the upstream release: bump the pin (S8 inventory), have
  `mcpAdapter.Invoke` return an error carrying `dispatchError(err)`, and flip the M1 e2e wire
  assertion. This row does not do that work.

  The row is **externally blocked with an active predicate**, which the loop runs at Gate 1 like any
  externally blocked row. The tripwire in R4 only fires on a pin bump, so it is not enough. The row
  text carries this command verbatim, with `<N>` set to the issue number filed above:

  ```bash
  gh issue view <N> -R sunholo-data/ailang --json state -q .state
  tag=$(gh release list -R sunholo-data/ailang --limit 1 --exclude-drafts --exclude-pre-releases --json tagName -q '.[0].tagName')
  gh api "repos/sunholo-data/ailang/contents/serveapi/protocol/mcphttp/methods.go?ref=$tag" -q .content | base64 -d | grep -cE 'isError|IsError'
  gh api "repos/sunholo-data/ailang/contents/serveapi/protocol/envelope.go?ref=$tag" -q .content | base64 -d | awk '/^func CallbackMessage/,/^}/' | grep -cE 'errors\.As'
  ```

  **UNBLOCKED** iff the issue state is `CLOSED` **and** at least one count is ≥ 1 (an `isError` field
  in `mcphttp`, or a typed-error hook in `CallbackMessage`). Otherwise the row stays BLOCKED and the
  loop skips it. The predicate was dry-run on 2026-10-06: `tag=v0.52.1`, counts `0` and `0` (BLOCKED,
  as expected). The `gh issue view` form returned `CLOSED` for a known issue (#885). A hook under
  another name will not match, so a count of 0 against a CLOSED issue is a prompt to read the release
  notes, not proof of absence.
- **Known limit, recorded for the follow-up (V6, verified).** `Invoke` runs under `Run`'s `callCtx`
  (20 s, `runner.go:42`). Its own `WithTimeout(ctx, invokeWait)` (`mcp.go:90`) adds nothing, because
  it is the same 20 s.

  The handler budget reserves `HandlerHeadroom = FinishPhaseBudget + PostEffectBudget` = 6 s
  (`host/coordinator/effectful.go:37–41`, `handlerBudget` `:112–118`), so normally an
  `EffectsUnrecordedError` returns inside the deadline. But the finish and post-effect phases run on
  `detached` contexts (`effectful.go:122–123`) that ignore the caller's deadline. If they push Dispatch
  past 20 s, `Run` takes its `case <-callCtx.Done(): return zero, callCtx.Err()` branch
  (`runner.go:64–65`), and the wire says `host callback timed out` (V4). The late result goes into the
  buffered `done` channel and is discarded.

  Even an upstream fix cannot carry the typed error in that window. The M1 operator line is still
  written, because `mcp.go:131–132` runs inside the callback goroutine after `Dispatch` returns.

## Acceptance criteria

Each criterion is load-bearing and names one test and one mutation (S6). Mutations are
implementation-stage; the evaluator runs them.

| AC | Test (new unless noted) | Asserts | Mutant it kills |
|---|---|---|---|
| AC1.1 | `host/daemon/setools_e2e_test.go` `TestSeToolsMCPPostEffectFailureIsLogged` (r14Store rig, `ErrorLog` captured, `/mcp/` `tools/call ailang-read`) | the wire is exactly `-32603 "host callback failed"` (pins today's upstream, V3–V4). The log has exactly one line starting `ailang-worldd: mcp refusal: tools/call a2a:ep1:` at M1 (`mcp:ep1:` from M2 on, AC2.1b). Its quoted cause contains a `sha256:` ref that `GetObject` resolves to `broker.EffectRecordV1` with `Allowed` and `Effect == Workspace.Read`. `entryCount` is unchanged | **MUT-MCP-LOG-NOCAUSE**: `mcp.go` logs a constant instead of `err` → the ref-resolves assertion fires. **MUT-MCP-LOG-LABEL**: MCP passes `"a2a"` → the prefix assertion fires |
| AC1.2 | existing `a2a_wire_test.go:208–349` `assertRefusalLog` cases, unchanged | A2A operator bytes are unchanged | **MUT-A2A-LABEL**: A2A passes `"mcp"` → the existing `ailang-worldd: a2a refusal:` check at `a2a_wire_test.go:496` fires |
| AC1.3 | `projection_test.go` dispatchError table (existing, `EffectsUnrecorded` rows) | the typed message and refs are unchanged | **MUT-UNRECORDED-ORDER**: move the `unrecorded` case below `conflict` → the `EffectsUnrecorded` row fires |
| AC2.1 | `mcp_conformance_test.go:127,137` (edited to `mcp:ep-a:`) | each MCP receipt id is `InvocationID(SurfaceMCP,"ep-a",task)` and starts `mcp:ep-a:` | **MUT-MCP-SURFACE**: `mcp.go` sets `SurfaceA2A` → the prefix assertion fires |
| AC2.1b | AC1.1's prefix assertion, edited in the M2 commit from `a2a:ep1:` to `mcp:ep1:` | the MCP refusal line names the `mcp:` id end to end over `/mcp/` under R14 | **MUT-MCP-SURFACE** (as AC2.1) → AC1.1's prefix assertion fires |
| AC2.2 | `coordinator_test.go` `TestDispatchRefusesUntaggedSurface` | `Call{Surface:""}` and `Call{Surface:"rest"}` → `InvalidCallError`. Store unchanged | **MUT-SURFACE-DEFAULT**: `validateCall` defaults empty to `a2a` → the empty case fires |
| AC2.3 | `setools_e2e_test.go` `TestSeToolsSurfaceTaggedInLog` (one MCP call and one A2A call on one rig) | the MCP entry header has `writtenBy == "coordinator:mcp"` and its record `invocationId` starts `mcp:ep1:`. The A2A entry has `coordinator:a2a` and `a2a:ep1:<task>`, and its A2A `metadata.invocation_id` equals the record's id | **MUT-WRITTENBY-CONST**: `planCommit` ignores the surface → the MCP header assertion fires |
| AC2.4 | `why_test.go` new row `mcp:ep1:x` (synth fixture with `coordinator:mcp`) | `why mcp:…` resolves the entry. `why bogus:…` is still refused | **MUT-WHY-PREFIX**: drop `mcp:` from `parseTarget` → the new row fires |
| AC2.5 | `commit_budget_test.go:356` table (+ `"mcp:ep:task"`) | `/v1/commit` refuses an `mcp:` id | **MUT-REST-ACCEPTS-MCP**: `restInvocationID` also accepts `mcp:` → the new row fires |
| AC2.6 | compatibility: the existing `why_test.go:112–122` synth fixtures (`coordinator:a2a`, `a2a:` ids), `effectful_test.go:509–622` Replay and `coordinator_test.go:990` reconcile, **byte-unchanged** | an old-shape store replays, reconciles and walks | **MUT-A2A-ID-FORMAT**: `InvocationID` renders A2A as anything but `a2a:<ep>:<task>` (e.g. surface after episode) → `coordinator_test.go:818` (re-signed `InvocationID(SurfaceA2A,"ep","task") == "a2a:ep:task"`) and the `:990` reconcile message case fire |
| AC3.1 | the upstream issue URL and the `ailang messages` id are recorded in the iteration log; `mcp-and-a2a.md` warning updated | the gate exists and is addressable | n/a: evidence, not a test (labelled instrument) |

Gate: `./scripts/verify_ail.sh` (with `AILANG_BIN`, though no `.ail` changes) and `go vet ./... && go test
./...` with `AILANG_BIN` set. CI runs `-race`, so add `-race` to the executor's gate list.

## Non-goals

- Any MCP wire change in this row. The client still sees `host callback failed` until upstream ships.
- Returning invocation ids or entry indexes in MCP results (`mcp.go:135`). That is QUICKSTART
  item 4 of row 137, a separate row.
- Re-tagging historical entries, or any store migration.
- Changing the A2A wire, the refusal semantics, deadlines, or the logged cause content (`err.Error()`, `%q`, one physical line).

## Risks

- **R1: docs drift.** `a2a:` / `coordinator:a2a` appear in about 8 website/docs files (V11 grep).
  Mitigation: the M2 sweep grep is an executor gate step. The row-137 binding tests catch QUICKSTART
  drift.
- **R2: downstream consumers matching `coordinator:a2a`.** None in repo code other than the dead
  constant (V11). Mixed old and new stores are expected and documented.
- **R3: the upstream ask is declined or slow.** The MCP refs then stay recoverable only by the
  operator (M1). The follow-up row stays gated. No local workaround (A/A′ rejected above).
- **R4: the frozen envelope assertion in AC1.1 breaks on a future pin bump.** This is intended. It is
  the tripwire that tells the follow-up row the upstream fix has arrived.

## Quorum verification log

Round-1 artifact: `~/.ailang/state/mission-quorum-world/w-mcp-effects-unrecorded-is-untyped-2026-10-06T07-29-36Z.json`.
The synthesis said PROCEED at N-1. A solo re-run restored the kimi seat, which the synthesis had marked
"invalid", and its REJECT made a revision owed. Round 2 below is the single allowed revision. It has
not been re-reviewed.

| Round | Reviewer / lane | Verdict | Findings | Disposition |
|---|---|---|---|---|
| 1 | gemini-3-1-pro | PASS | A strict empty-`Surface` refusal breaks every existing `Call{…}` literal | APPLIED r2: M2 "Call-literal sweep" counts 5 literals by file:line, keeps the strict refusal, and flags `commit_budget_test.go:154` |
| 1 | oc-glm-5-3 | PASS | The upstream-gate tripwire is passive (fires only on a pin bump) | APPLIED r2: M3 follow-up row carries an active Gate-1 predicate (exact `gh` commands, UNBLOCKED rule, dry-run 2026-10-06 = BLOCKED) |
| 1 | oc-kimi-k3 (restored by solo re-run) | REJECT | (a) V6 lacked file:line; (b) cross-revision identity in V1/V2/V4 unevidenced; (c) AC1.1 asserted `mcp:` before M2 existed; (d) the "what the operator line logs" non-goal was ambiguous | APPLIED r2: (a) V6 cites `runner.go:42, 64–65` and `daemon.go:99, 673, 674`, and the M3 known limit is re-derived with the `effectful.go` budgets; (b) V14 sha256 table. **Correction:** `methods.go` and `handler.go` are NOT byte-identical at v0.52.1, but the only hunks are `toolJSON` Title/Annotations and the OAuth bearer `Gate`; `callTool`, `serveMessages` and `callResultJSON` hash identical; (c) M1 and M2 are split patches, AC1.1 asserts `a2a:ep1:` at M1, and AC2.1b moves it to `mcp:ep1:`; (d) reworded to "logged cause content" |
| 1 | gpt6-1-sol | ABSENT | OpenAI 429 "no credits remaining"; cannot be restored | none; round is N-1 |
| 1 | controller | PASS | none | none |
| 2 | (single allowed revision; designer claude-opus-5-5) | REVISED | items 1–6 above | awaiting controller disposition |
| 2 | gemini-3-1-pro | REJECT | The M3 Gate-1 predicate's `grep -c` exits 1 on zero matches, so under `set -e` it aborts instead of evaluating to BLOCKED | OPEN. Reviewer fix: append `\|\| true` to both `grep -cE` commands |
| 2 | oc-glm-5-3 | REJECT | M1 changes the operator line (`a2a refusal: mcp invoke` → `mcp refusal: tools/call`), but every doc edit and the consumer sweep are scoped to M2; catch: V14 has `—` for the module-cache `hostcall/runner.go` hash that V6 cites | OPEN. Reviewer fix: move the doc step and a `grep -rnE "a2a refusal\|mcp invoke\|mcp tools" website docs` sweep (recorded as V15) into M1; fill the runner.go cell |
| 2 | oc-kimi-k3 (quorum seat) | REJECT | **Design scope:** no option makes an MCP retry safe host-side. A caller-supplied idempotency key minted as `mcp:<ep>:<key>` would reach the existing `committed()` dedup (`coordinator.go:199`) and Replay | OPEN, **needs a human ruling (D-WORLD-66)**. Reviewer fix: add option (D) and either adopt it with an AC, or reject it with file:line evidence that `committed()`/Replay cannot be reached from `mcp.go`'s Dispatch path |
| 2 | oc-kimi-k3 (solo re-run, $0.40 cap) | REJECT | V14's line-shift arithmetic contradicts itself (`:182 → :208` is +26, not +30; V3 anchors at `:178–183`) | OPEN. Reviewer fix: re-run the `grep -n "CallbackMessage(hostErr)"` extraction at both revisions and state the measured shift; specify `-` as the id at `mcp.go:76, :81` |
| 2 | gpt6-1-sol | ABSENT | OpenAI 429 "no credits remaining" | none; round is N-1 |
| 2 | controller | PASS (in-session) | — | — |

**Round-2 outcome (iteration 236): BLOCKED, PARKED `needs-human-review`.** Round-2 artifact:
`~/.ailang/state/mission-quorum-world/w-mcp-effects-unrecorded-is-untyped-2026-10-06T07-37-35Z.json`.
Gemini's, GLM's and solo-kimi's objections carry verbatim fixes and do not dispute direction. The
quorum-kimi objection asks for a new option, and the reviewer's fix offers adopt-or-reject. The
reject branch's stated condition (`committed()`/Replay unreachable from `mcp.go`) is FALSE: they are
reachable if a key exists. So the narrow-refinement carve-out does not apply, and the choice is a
design-direction call for Mark. Facts the controller measured for that call at `3d965f6`:
(1) the released `callTool` decodes only `name` and `arguments` (`mcphttp/methods.go:73–76`,
v0.47.2), so a `_meta` idempotency key never reaches the host; (2) 9 of the 10 `additionalProperties`
keys in `packages/se-tools/transitions.json` are `false`, so a key carried in `arguments` changes
every advertised tool input schema; (3) a key derived from a hash of tool+arguments would collapse
two legitimate identical calls (e.g. two reads of a file that changed in between) into one.
