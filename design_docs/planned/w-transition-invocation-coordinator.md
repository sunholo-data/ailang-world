# w-transition-invocation-coordinator — The propose → verify → commit coordinator behind `/a2a/` `tasks/send` (row 106)

**Status**: Planned — **REVISION 1** after quorum round 1 (BLOCKED 2/2; see "Quorum log"). Design + prototype (iteration 196, designer `claude:claude-opus-5-5`; the previous rotation entry `pi:ollama/kimi-k3` was quota-cut after research only, so nothing is inherited from it).
**Item**: queue row 106 of `design_docs/world-mission.md` (`w-transition-invocation-coordinator`), regroom position 2, now the head of the clause-6 critical path after row 107 landed (`b7cfbdb`).
**Clauses**: clause-6 (*"the transition registry is served over MCP (capability-filtered per session) and an A2A agent card is published; no new wire protocols"*) and clause-3 (*"every effect goes through the broker with a capability + budget check; effect results are recorded (replay input); capsules run with a physical isolation floor … No ambient-authority path exists from an agent to the outside world"*). It is also the World arm of row 93 (non-inferiority floor run). Row 108 (MCP dispatch, blocked upstream) will dispatch into the same coordinator.
**Estimate**: ~2d as the row says. Seven milestones of ≤ ~150 production code lines each, sized from the prototype's measured code lines (V27). **First landable slice = M1–M5**, which makes `/a2a/` answer a real result. M6 (the real-interpreter daemon end-to-end run plus the QUICKSTART) closes the slice if velocity matches row 107, which landed 7 milestones in one iteration. The fallback split point is after M4: the coordinator lands with tests, and `/a2a/` keeps refusing until M5.
**Verified against**: dev `535330d` (worktree `.design-wt-iter196`, branch `design/iter196-w-transition-invocation-coordinator`). Every claim of the form "the codebase does X" has a Verification Log row with its command and observed output. Pinned interpreter: `$HOME/.pinned-ailang/ailang` = `AILANG v0.41.0 / Commit: 24ee108` (V0).

---

## Problem (measured at `535330d`)

World has a **catalogue** but no **server-side path that invokes an entry in it**:

- Row 107 landed a production writer: `world-publish transitions` → `transitionreg.PublishSet`. Row 40 landed the reader: the session-scoped A2A card at `GET /.well-known/agent.json`. The only production importers of `host/transitionreg` are exactly those two plus the daemon's reader construction (V2).
- `POST /a2a/` answers an **authorized** `skill_id` with the constant `-32603 "transition invocation is not available in this daemon"` (`host/projection/projection.go:318`, V4). The same constant also covers the resolution-error and registry-error paths (`:270`, `:312`).
- **Nothing in production runs AILANG transition source.** `host/capsule` and `host/replay` have **zero** production importers. Their only importers are two `host/broker` test files (V3).
- There is **no** propose/coordinate function in `host` or `cmd` (V1, with a positive control).
- The daemon's only mutation, `POST /v1/commit`, takes an **already-computed** `store.Commit` from the client (`host/daemon/handlers.go:546`, V18). A client can therefore append any log entry it likes, including one whose `TransitionFn` no interpreter ever executed. No path exists where *the World executes a transition and then commits its own result*.
- `transitionreg.Bind` / `Bound.Check` exist (`host/transitionreg/bind.go:61,87`) but have **no production Binder**. The broker's structural gate TR.C forbids constructing a `broker.Session`, or calling any `Invoke` selector, outside `host/broker` (V5). `host/daemon` documents that exact constraint as the reason it never built a session (`handlers.go:554-563`, V18).

Three inherited facts shape the design and are **new at this HEAD**:

1. **The capsule cannot pass an argument and cannot be cancelled.** `capsule.Runner.Run(entry Entry)` takes no context: it builds its own `context.WithTimeout(context.Background(), 60s)` (`capsule.go:171,194`, V6). It invokes `run --quiet --caps "" --entry main host/capsule/main.ail` with no arguments. The pinned interpreter *does* support a hermetic argument channel: `ailang run --args-file <path>` decodes a JSON value into the entry's parameters (V7).
2. **The publish-time check and the capsule disagree on module paths.** The publisher checks a source as `entry.ail` under `AILANG_RELAX_MODULES=1` (`host/archive/check.go:71-83`, V8). The capsule stages it at `host/capsule/main.ail` with **no** relax flag (`capsule.go:32,199`). A source declaring `module transitions/echo` therefore **passes publication and fails invocation** with `MOD010` under a `MkdirTemp` root. That was measured on the pinned binary (V9). It is a concrete instance of row 107's Residual 8.
3. **A source that checks can still be un-invocable.** A module whose `main` has arity 0 passes `ailang check` (rc 0). Invoked with an argument it fails with `ARG_DECODE_MISMATCH … expected null for unit type` (V10). This is exactly the negative fixture that Residual 8(c) asks for.

## Inherited obligations (row 107's Residuals 6, 7, 8 and 1): what the first slice closes

| Residual | Decision in this row | Why |
|---|---|---|
| **8 — invocation contract** (owned here) | **CLOSED by M1+M4+M6.** (a) The source-object convention `world/transition-source/v1` is **adopted unchanged** (Residual 1). The coordinator reads the object by `Descriptor.TransitionFn` and re-hashes it. (b) M6 publishes a self-contained transition through `PublishSet` with the real pinned interpreter and invokes that exact stored source through the production capsule path over `/a2a/`. (c) **Decision: invocation reports a typed incompatibility; publication is NOT widened.** The zero-arity fixture (V10) yields `*coordinator.IncompatibleError`, mapped to `-32603` with its own constant message. | Detecting "has a `main(string) -> string`" at publish time would need either a second interpreter run at publish time, or a parse of the interpreter's `check` output that the pinned CLI does not expose as an interface. Invocation already executes the pinned interpreter, and the mismatch is observable there as the interpreter's own typed failure. A publish-side guard can be added later behind `EnsureSourceLoadable` without changing this row's wire (Residual R-106-4). |
| **Fact 2 above (MOD010 divergence)** | **CLOSED by M1.** `capsule.RunContext` stages under `AILANG_RELAX_MODULES=1`, the **same** shape as the publish check. This is not a new relaxation: the invocation capsule adopts the one already controller-accepted in row 107. | Without it, every published module not literally named `host/capsule/main` is un-invocable (V9). The check and the run must agree, or publication verifies the wrong thing. |
| **7 — no landed `.ail` is publishable** | **Closed for the *invocation proof* only.** A hermetically self-contained transition (`echoSrc`, an inline Go-string fixture in `host/coordinator/coordinator_test.go`) is stored under row 107's source convention and invoked through the production capsule path with the real pinned interpreter (prototype, V25). In M6 it is published through `PublishSet` and invoked over `/a2a/`. It stays an inline Go fixture, following the `host/capsule/capsule_test.go` `source(...)` precedent. `verify_ail.sh` sweeps only `design_docs/` and `world/` against an exact module manifest (`scripts/verify_ail.sh:161-169`, V28), so a `.ail` testdata file would be unswept or would force a manifest edit. **Not closed for production card content**: a *useful* first skill is a package, per S3 ("why is this not a package?"). It is owned by row 93, whose floor run needs real skills (R-106-1). | A demo skill in `world/` would grow the kernel without an S3 answer. A test fixture proves the contract without claiming product content. |
| **6 — world-library staging** | **NOT closed; stays hermetic.** The capsule stages only the entry source, exactly like the publish check. `world/transitions.ail` remains un-invocable by design (V11: `LDR001`). | Staging a world library makes every invocation depend on library bytes that the descriptor does not pin. Replay would then be non-hermetic, which violates ratified D1. Pinning a library needs a descriptor field, so a registry interface change, so a new `InterfaceHash`. That is a row of its own (R-106-2). |

## Verification Log

All commands were run in `/Users/voightkampff/dev/sunholo-data/.design-wt-iter196` at `535330d` with `PATH=/opt/homebrew/bin:$PATH`. Every negative result is paired with a positive control in the same row.

| # | Claim | Command → observed |
|---|---|---|
| V0 | Pinned interpreter identity | `$HOME/.pinned-ailang/ailang --version` → `AILANG v0.41.0` / `Commit: 24ee108` |
| V1 | No propose/coordinate function exists in production Go | `grep -rniE "func .*(propose\|coordinat)" host cmd --include='*.go' \| grep -v _test.go \| wc -l` → **0**. **Control** (same regex can hit): `printf 'func proposeNext() {}\nfunc (c *Coordinator) x() {}\n' \| grep -ciE "func .*(propose\|coordinat)"` → **2** |
| V2 | Production importers of `host/transitionreg` | `grep -rln '"github.com/sunholo-data/ailang-world/host/transitionreg"' host cmd --include='*.go' \| grep -v _test.go` → `host/daemon/daemon.go`, `host/projection/projection.go`, `cmd/world-publish/transitions.go` (3; iter-187's "0" is superseded by rows 40 and 107) |
| V3 | `host/capsule` and `host/replay` have zero production importers | `grep -rln '"github.com/sunholo-data/ailang-world/host/\(capsule\|replay\)"' host cmd --include='*.go'` → only `host/broker/registry_publish_test.go`, `host/broker/episode_test.go`. Both are `_test.go`, which is the **positive control** that the instrument sees importers. Production: **0** |
| V4 | `/a2a/` refuses an authorized skill with the constant `-32603` | `grep -n "notAvailableMessage)" host/projection/projection.go` → `:270` (resolver error), `:312` (registry error), `:318` (authorized skill) |
| V5 | TR.C: `Invoke` selectors and `broker.Session`/`NewSession`/`NewReplaySession` are forbidden outside `host/broker`; the inside exemption count is exactly 3 | `sed -n 84,105p host/broker/invoke_boundary_test.go` (the selector detectors); `grep -n "wantCount = 3" host/broker/invoke_boundary_test.go` → `:274`; baseline `go test -run TestRegistryDispatchBindingBoundary ./host/broker/` → `ok` |
| V6 | `capsule.Run` takes no context and passes no argument | `grep -n "func (r \*Runner) Run\|exec.CommandContext\|context.Background" host/capsule/capsule.go` → `:171 func (r *Runner) Run(entry Entry)`, `:194 context.WithTimeout(context.Background(), r.execTimeout)`, `:196 exec.CommandContext(ctx, execPath,` (args: `run --quiet --caps <r.caps> --entry main host/capsule/main.ail`, `:197`) |
| V7 | The pinned interpreter has a hermetic argument channel (`--args-file`) that decodes a JSON string into `main(input: string)` | `ailang run --help` → `-args-file string  Path to a file containing JSON arguments`. Probe with `module host/capsule/main` / `export func main(input: string) -> string { input }` in a `mktemp -d` root, args file `"{\"a\":1}"`, `env -i AILANG_FS_SANDBOX=$R ailang run --quiet --caps "" --entry main --args-file args.json host/capsule/main.ail` → stdout `{"a":1}`, rc 0. **Control**: with no argument → `ARG_DECODE_MISMATCH: expected string, got <nil>`, rc 1 |
| V8 | The publish check stages `entry.ail` under `AILANG_RELAX_MODULES=1` | `grep -n "main.ail\|MkdirTemp\|RELAX\|\"check\"" host/archive/check.go` → `:66 MkdirTemp`, `:71 const checkFile = "entry.ail"`, `:77 … "check", checkFile`, `:83 … "AILANG_RELAX_MODULES=1"` |
| V9 | A module-named source passes publication but fails the capsule's strict run | Same root as V7 with `module transitions/echo`: `env -i AILANG_FS_SANDBOX=$R ailang run … host/capsule/main.ail` → `Error MOD010: module 'transitions/echo' doesn't match file path 'host/capsule/main'. Fix: use --relax-modules flag or set AILANG_RELAX_MODULES=1`, rc 1. **Control**: the same file with `module host/capsule/main` → `{"a":1}`, rc 0 |
| V10 | Residual 8(c) negative fixture: checks but cannot be invoked | `module host/capsule/main` / `export func main() -> string { "no-input" }`: `ailang check host/capsule/main.ail` → `✓ No errors found!`, rc 0. `ailang run … --args-file args.json` (args `"hi"`) → `ARG_DECODE_MISMATCH: expected (), got hi / expected null for unit type`, rc 1 |
| V11 | `world/transitions.ail` is refused standalone (Residual 7 premise re-measured) | `cp world/transitions.ail $T/host/capsule/main.ail; AILANG_RELAX_MODULES=1 ailang check host/capsule/main.ail` → `LDR001: module not found: world/logepoch` |
| V12 | Output shape: `--quiet` prints the returned string raw, with a trailing newline | V7 probe piped to `od -c` → `{ " a " : 1 } \n` |
| V13 | The upstream wire types used here exist in the pinned dependency | `grep -n "sunholo-data/ailang " go.mod` → `github.com/sunholo-data/ailang v0.33.2`. In that module, `serveapi/protocol/a2a_wire.go` has `A2ARequest`, `A2ATaskSendParams{ID, Message, Metadata}`, `A2AMessage{Role, Parts}`, `A2AContent{Type, Text, Data map[string]any}`, `A2AError`, `A2AResult`. `interfaces.go` has `Invoker{Invoke(ctx, Session, Invocation)}` |
| V14 | Daemon timeouts the invoke bound must nest inside | `grep -n "readDeadline = \|writeTimeout *= \|capsuleExecTimeout *= " host/daemon/daemon.go host/capsule/capsule.go` → `daemon.go:94 writeTimeout = 30 * time.Second`, `daemon.go:138 readDeadline = 10 * time.Second`, `capsule.go:30 capsuleExecTimeout = 60 * time.Second` |
| V15 | The projection is mounted over the daemon's own handles; MaxWait = readDeadline | `sed -n 520,545p host/daemon/daemon.go` → `projection.New(projection.Config{Resolver: d.resolver, Reader: transitionreg.NewReader(d.store), Heads: d.reads, … MaxWait: readDeadline})`; `grep -n 'POST /a2a/' host/daemon/daemon.go` → `:651 mux.HandleFunc("POST /a2a/", d.projection.A2A)` |
| V16 | The daemon archives the interpreter only when `cfg.AilangBin` is set | `sed -n 491,505p host/daemon/daemon.go` → `if cfg.AilangBin != "" { a := archive.New(cfg.DBPath); ref, err := a.Archive(cfg.AilangBin) …}`. Otherwise `release = unpinnedRelease` |
| V17 | The store's journal binds a commit to a durable intent; commit is CAS on the observed head | `sed -n 152,160p host/store/store.go` → `type Commit struct { InvocationID string; ObservedHead …; Objects …; NextWorld …; Entry … }`; `store.go:931-949`: `bindCommitIntentTx` then `ConflictError` on stale head; `store.go:1013` appends `'outcome'` in the same tx; `journal.go:411 AppendIntent`, `:324` refuses the `effect:` namespace |
| V18 | `POST /v1/commit` takes a client-computed commit and never builds a broker session (TR.C) | `sed -n 546,600p host/daemon/handlers.go` → `decodeCommit(request)` then `d.store.Commit(commit)`; comment `:554-563` "the broker's structural gate TR.C … forbids Session construction outside host/broker" |
| V19 | `Session.Bind` validates and copies the manifest; `BoundInvoker.Request` refuses undeclared triples | `sed -n 36,72p host/broker/confined.go` |
| V20 | `transitionreg.Bind`'s typed refusals | `sed -n 34,84p host/transitionreg/bind.go` → `TransitionAbsentError`, `AccessDeniedError{Label}`, `ProposalMismatchError{Field}`, and an untyped zero-snapshot error |
| V21 | Kernel commit laws the Go plan must satisfy | `sed -n 45,95p world/transitions.ail` → `verify` accepts iff `proposalMatchesWorld` (input world = current `stateRoot`); `applyRevision` `ensures result.revision == w.revision + 1 && stateRoot == outputWorld && logHead == nextLogHead` |
| V22 | Clause 7's package pipeline is attended-only (so it cannot be reused for server-side invocation) | `grep -n "func requireAttendedOperator\|func refuseAutomationEnvironment" cmd/world-publish/fences.go` → `:226`, `:182` (controlling-TTY probe + typed phrase + CI-env refusal, row 107 F9) |
| V23 | Baseline green for the packages this row touches | `AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./host/capsule/ ./host/transitionreg/ ./host/projection/ ./host/store/` → 4 × `ok` (30.5s / 16.0s / 0.9s / 11.0s) |

| V30 | (revision 1, objection 2) The durable steps take no context. The pool is a single connection. `busy_timeout` bounds lock waits only | `grep -n "func (s \*Store) AppendIntent\|func (s \*Store) Commit(" host/store/*.go` → `journal.go:411 AppendIntent(id string, intent JournalIntent)`, `store.go:874 Commit(c Commit) error` (no `context.Context` parameter on either, nor on `journal.go:813 GetReceipt(id string)`); `grep -n "s.db.Begin()"` → `store.go:909`, `journal.go:420` (`Begin`, not `BeginTx(ctx, …)`); `grep -n SetMaxOpenConns host/store/store.go` → `:305 db.SetMaxOpenConns(1)`; `sed -n 180,187p host/store/writer_lock.go` → `const busyTimeoutMillis = 2000`, commented "a measured margin, not a bound". **Control** (the instrument can see a ctx-taking store call): `grep -n "func (s \*Store) GetObject" host/store/store.go` → `:475 GetObject(ctx context.Context, ref hashref.HashRef)` |
| V31 | A commit-bound `Commit` writes its `committed` outcome in the same transaction as the world | `sed -n 992,1000p host/store/store.go` → `// Step 6: append the receipt in this same transaction.` … `Status: "committed", ResultRef: c.NextWorld.Ref`; `:1019` `tx.Commit()` is the single commit |
| V32 | (revision 1, objection 1) The skill ID is `Descriptor.ID`. It is the card's `skills[].id` and the key `Bind` looks up | `grep -n "^\tID " host/transitionreg/transitionreg.go` → `:22 ID string` (first field of `Descriptor`); `grep -n '"id": d.ID' host/projection/projection.go` → `:240 "id": d.ID, "name": d.Title, …`; `grep -n "snap.Lookup(id)" host/transitionreg/bind.go` → `:65 d, ok := snap.Lookup(id)`; `grep -n "SkillID: d.ID" host/coordinator/plan.go` → `:79` |
| V33 | Broker recovery reports an intent without an outcome and never resolves it | `sed -n 16,18p host/broker/recover.go` → "Recovery only reports this ambiguity; it never calls an effect handler or writes a journal outcome." |
| V34 | Pre-existing macOS flake in `host/capsule` (not caused by this row) | `TestOrdinaryOverflowCarriesNoESRCH` failed once in the revision-1 `go test ./host/coordinator/ ./host/capsule/ ./host/broker/` run: `capsule: overflow kill: operation not permitted` (EPERM from `kill(-pgid)`). Stress, with both binaries run concurrently under the same load: prototype `capsule.go` **3/320** failures, pristine `HEAD` `capsule.go` (built from `git show HEAD:…`, then restored byte-identical, `cmp` OK) **3/320**. The two agree, so the fault is pre-existing: macOS returns EPERM, not ESRCH, for a zombie-only group. Out of scope; filed as an observation for the controller (row 24's ESRCH handling does not cover EPERM) |

Prototype-era rows (V24+) are listed under **Prototype** below.

## Design

### Where it lives — and "why is this not a package?" (S3)

New host package **`host/coordinator`**. It is not an AILANG package because the coordinator *is* the host effect boundary (S2). It opens a capsule subprocess, reads and writes the store, and constructs a broker binder: three effects that `world/` must never perform. The **laws** it enforces are the existing kernel laws from `world/transitions.ail` (V21). The Go planner transcribes them and names each contract in a test. It invents none. It is also not `host/projection`, because MCP row 108 must dispatch into the same code (P4: "Protocol invocation enters the same broker/session/transaction path as every other agent action"). And it is not `host/broker`, because broker would then import `transitionreg`, `capsule` and `store` planning. `transitionreg` already imports `broker`, so that would be a cycle.

**Clause 7 reuse check.** The existing propose → verify → commit pipeline is `world-publish packet/approve/publish`. It is attended-only by construction: a TTY probe, a typed phrase, and a CI-env refusal (V22). Reusing it for agent-initiated invocation would either require a human per call or remove its fences. So this row reuses the pipeline's **parts** and not its CLI:

- `transitionreg.Bind` / `Check` (authorization + pin agreement)
- the broker's `Session.Bind` confinement
- the store's journal (`AppendIntent` + `Commit{InvocationID}`), the same durable intent/outcome pairing
- the capsule floor

### The one new broker seam (M2): `broker.OpenBinder`

```go
// host/broker/binder.go
// SessionBinder is the ONLY production handle on a session outside
// host/broker: it can Bind a descriptor manifest and nothing else. The raw
// Session (and its Invoke) stays unexported to every other package (TR.C).
type SessionBinder struct{ s *Session }

func OpenBinder(s *store.Store, episodeID string, grants []Capability) *SessionBinder
func (b *SessionBinder) Bind(m Manifest) (*BoundInvoker, error) // delegates to s.Bind
```

`OpenBinder` calls the unexported `newSession` with a nil handler registry. It never selects `Invoke` or `NewSession`, so TR.C's inside-exemption count stays **3**, and outside `host/broker` no finding appears (`*broker.SessionBinder` is not a `Session` selector). A consumer outside the broker can reach dispatch only through `BoundInvoker.Request`, which refuses undeclared triples (V19). The seam therefore *narrows* reachable authority relative to a raw session. The registry is nil because slice-1 transitions cannot request effects (see refusal R8).

### The capsule change (M1): `RunContext` with an argument

```go
// host/capsule/capsule.go
type Entry struct {
	Interpreter hashref.HashRef
	Source      []byte
	Args        []byte // NEW: one JSON value written to <root>/args.json; nil = legacy no-arg run
}
// RunContext is Run bounded ALSO by ctx: the child's context is
// context.WithTimeout(ctx, r.execTimeout), so the earlier of the caller's
// deadline and the exec allowance kills the process group.
func (r *Runner) RunContext(ctx context.Context, entry Entry) (Result, error)
func (r *Runner) Run(entry Entry) (Result, error) { return r.RunContext(context.Background(), entry) }
```

When `Args != nil`, the runner writes `args.json` inside the capsule root and appends `--args-file args.json`. The file is read by the interpreter host process, not by the sandboxed program, and never passes through a shell. The environment becomes `{"AILANG_FS_SANDBOX=" + root, "AILANG_RELAX_MODULES=1"}`: the publish check's own shape (V8/V9). Legacy `Run` behaviour is unchanged apart from the relax flag. That flag is a **strict widening** (a source that ran before still runs), and the existing capsule tests pin it (AC-CAP-LEGACY).

### Calling convention (the invocation contract; answers Residual 8)

A published transition is invocable iff its canonical source (`world/transition-source/v1` object at `TransitionFn`) declares **`export func main(input: string) -> string`**:

- **input**: the coordinator passes the *canonical JSON text* of the caller's single A2A `data` part, as one JSON **string** argument. The argument file is `json.Marshal(string(canonicalInput))`, so the program receives JSON text and decodes it itself. The host never interprets the transition's schema.
- **output**: the returned string, printed raw by `--quiet`. It must be a JSON **object** after stripping exactly one trailing `\n` (V12). A2A `data` parts are objects (`A2AContent.Data map[string]any`, V13).
- **effects**: none. The capsule runs `--caps ""`, so a descriptor with non-empty `DeclaredEffects` is refused *before* execution (R8). Wiring an effect bridge from the capsule child to the broker is R-106-3.
- **schemas**: `InputSchema`/`OutputSchema` are carried and unvalidated in slice 1 (R-106-5, no JSON-Schema validator is admitted as a dependency). The contract enforced is *JSON-object in, JSON-object out*. This is stated as such on the card description in M5 and not overclaimed.

### Coordinator API (M3, M4)

```go
package coordinator

// Store is the narrow store surface the coordinator needs (prod: *store.Store).
type Store interface {
	GetObject(ctx context.Context, ref hashref.HashRef) (store.Object, bool, error)
	SelectedHead(ctx context.Context) (hashref.HashRef, bool, error)
	GetWorld(ctx context.Context, ref hashref.HashRef) (store.World, bool, error)
	AppendIntent(id string, intent store.JournalIntent) (int64, hashref.HashRef, error)
	Commit(c store.Commit) error
}
// Runner is the capsule seam (prod: *capsule.Runner).
type Runner interface {
	RunContext(ctx context.Context, e capsule.Entry) (capsule.Result, error)
}
// BinderFor opens a session-scoped binder (prod: broker.OpenBinder over the store).
type BinderFor func(episodeID string, caps []broker.Capability) transitionreg.Binder

type Config struct {
	Store     Store
	Runner    Runner
	Binder    BinderFor
	Now       func() int64 // logical time for journal payloads; never read inside the pure core
	MaxInput  int         // canonical input cap, bytes (> 0; default refused)
	MaxOutput int         // output cap, bytes (> 0)
}
func New(cfg Config) (*Coordinator, error) // every field required; zero caps are construction errors

// Call is one authorized invocation request. Registry+Caps is the ONE
// snapshot pair the caller (projection / future MCP) already captured with
// transitionreg.NewRequest: dispatch never re-reads the registry (one-snapshot rule).
type Call struct {
	Request   transitionreg.Request
	EpisodeID string              // authority.SessionBinding.EpisodeID
	Grants    []broker.Capability // authority.SessionBinding.Caps (for the binder)
	SkillID   string
	TaskID    string              // A2A params.id — the idempotency key
	Input     map[string]any      // the single data part
	PinnedFn  *hashref.HashRef    // optional caller proposal pin (metadata.transition_fn)
}
type Result struct {
	Output       map[string]any
	OutputBytes  []byte // exact committed output bytes (replay compares these)
	InvocationID string
	WorldRef     hashref.HashRef
	EntryIndex   int64
	RecordRef    hashref.HashRef
}
func (c *Coordinator) Dispatch(ctx context.Context, call Call) (Result, error)
```

The method is **`Dispatch`**, not `Invoke`. A projection-side call `c.Invoke(...)` would be an `Invoke` selector outside `host/broker`, which TR.C reds (V5). `protocol.Invoker` is **not** implemented in slice 1: World serves its own `/a2a/` handler and nothing in-repo consumes an upstream `Invoker`. Row 108 adds a thin adapter when it lands (R-106-6), and a method *declaration* named `Invoke` is not a selector, so TR.C permits it.

### Order of operations — ONE bounded context

The **deadline source** is the `/a2a/` handler: `ctx, cancel := context.WithTimeout(r.Context(), h.invokeWait)`. Here `invokeWait` is a new required `projection.Config.InvokeWait`. The daemon sets it to `invokeDeadline = 20 * time.Second`, which is below `writeTimeout = 30s` (V14), leaving a 10 s margin for the durable steps and the response write in the common case. `projection.New` refuses `InvokeWait <= 0`, and the daemon's construction refuses `InvokeWait >= writeTimeout`. That same `ctx` flows unreplaced through resolution → snapshot → every coordinator step that accepts a context → `capsule.RunContext` (whose child context is `WithTimeout(ctx, 60s)`, so the earlier deadline wins).

**What the deadline bounds, stated exactly as strong as the code (revision 1, objection 2).** The invocation deadline bounds steps 1–9 and the pre-boundary check at the start of step 10. **It does not bound the durable steps.** `GetReceipt` (step 1b), `AppendIntent` and `Commit` (step 10) take no context, so no invocation deadline reaches them. They are bounded only by what the store itself bounds; see "Bounded waits — what is and is not bounded" below. **Nothing in this design guarantees that a JSON-RPC answer is written before `writeTimeout`.** What it guarantees instead is that (i) a commit that landed is never reported as a failure, and (ii) a client that got no answer can learn the outcome by resending the same task id, without anything re-executing (reconciliation, below).

| Step | What | Composes | Refusal (typed) |
|---|---|---|---|
| 0 | (projection) resolve session; capture ONE `transitionreg.Request` | `authority.ResolveContext`, `transitionreg.NewRequest` (existing) | existing denials |
| 1 | validate call: `TaskID` non-empty and ≤ 128 bytes of `[A-Za-z0-9._-]`; canonical input = `json.Marshal(Input)` ≤ `MaxInput` | pure | **R1** `*InvalidCallError{Field}` |
| 1b | **reconcile**: `GetReceipt(id)` (no ctx). **Resolved** → read the committed record/output/world back and return that result with `Reconciled: true`, with **no** execution and **no** commit. **Intent without outcome** → refuse. **Absent** → continue | store journal | **R15** `*NotCommittedError` |
| 2 | **propose**: `bound := transitionreg.Bind(call.Request.Registry, SkillID, call.Request.Caps, cfg.Binder(EpisodeID, grants))` | `transitionreg.Bind` → `broker.SessionBinder.Bind` | **R2** `*transitionreg.TransitionAbsentError`, **R3** `*transitionreg.AccessDeniedError`, **R4** bind error (wrapped) |
| 3 | **verify (pins)**: build `Proposal` from the bound descriptor, with `TransitionFn` replaced by `*PinnedFn` when the caller pinned one; `bound.Check(p)` | `Bound.Check` | **R5** `*transitionreg.ProposalMismatchError` |
| 4 | effect floor: `len(d.DeclaredEffects) != 0` → refuse (capsule is `--caps ""`) | pure | **R8** `*EffectsUnsupportedError` |
| 5 | load source `GetObject(ctx, d.TransitionFn)`; re-hash payload == `TransitionFn` | store | **R6** `*SourceError{Absent\|Corrupt}` |
| 6 | observe world: `SelectedHead(ctx)` + `GetWorld(ctx, head)` | store | **R7** `*WorldAbsentError` (no selected head). A head that names no world row is store damage, returned as an untyped wrapped error, i.e. `notAvailableMessage` |
| 7 | **execute**: `RunContext(ctx, Entry{d.Interpreter, src, Args: json.Marshal(string(canonicalInput))})` | capsule | **R9** `*IncompatibleError` (stderr carries the pinned interpreter's `ARG_DECODE_MISMATCH` or `entrypoint 'main' not found`); **R10** `*ExecutionError` (any other exec failure, output limit, capsule timeout); **R11** `ctx.Err()` (deadline/cancel, returned wrapped so `errors.Is(err, context.DeadlineExceeded)` holds) |
| 8 | output: strip one trailing `\n`, `len ≤ MaxOutput`, must decode to a JSON object | pure | **R12** `*OutputError` |
| 9 | **plan commit** (pure core, M3): input/output/record objects, next world, entry, intent | pure `planInvocation` | — (pure, total on validated inputs) |
| 10 | **commit boundary**: `if err := ctx.Err(); err != nil` → return R11 (**no durable mutation**). Then `AppendIntent(id, intent)`, then `Commit(store.Commit{InvocationID: id, …})`. Neither takes ctx. **After `Commit` returns nil the result is success, whatever the ctx state** | store journal | `AppendIntent` error → returned as-is, and `Commit` was never called, so nothing landed (a typed `*store.DuplicateInvocationError` is defensive only: step 1b answers a seen id first). **R14** `*store.ConflictError`, returned bare: the store compares the head before any write in the same transaction and rolls back (V17), so it definitively did not land. **R16** `*UnconfirmedError` for any other `Commit` error: the outcome is not confirmed to this caller |

The **invocation ID** is `"a2a:" + EpisodeID + ":" + TaskID`. The `effect:` namespace the journal refuses is never produced (V17). After step 10 begins, the store's own transaction makes the commit all-or-nothing. A commit-bound `Commit` writes its `committed` outcome in the same transaction as the world, log and head rows (V31), so **resolved receipt ⇔ commit landed**. An intent without an outcome is the landed journal's indeterminate state, which broker recovery only *reports* and never resolves (V33). The coordinator follows the same rule: it never writes an outcome itself, never re-executes, and never re-commits.

### Bounded waits — what is and is not bounded (revision 1, objection 2)

| Wait | Bounded by | Measured |
|---|---|---|
| Session resolution, registry/capability snapshot, source/head/world reads | the invocation `ctx` (row 23's plumbing: store reads inherit the caller's context) | existing, row 39/40 |
| Capsule execution | `min(ctx deadline, 60 s exec allowance)`, with process-group kill | AC-CAP-CTX, MUT-CTX-BG (V26) |
| Waiting for the **single pooled connection** in `GetReceipt` / `AppendIntent` / `Commit` | **nothing**. `s.db.Begin()` without a context (`store.go:909`, `journal.go:420`) on a `SetMaxOpenConns(1)` pool (`store.go:305`) waits behind every in-flight in-process store call. Context-carrying reads end at their own deadlines; context-free writers (`/v1/commit`, other invocations' commits) do not | **not bounded** (V30) |
| SQLite lock waits inside those transactions | `busy_timeout(2000)` per connection (`writer_lock.go:187`), which bounds SQLite *lock* waits only. The daemon's own comment calls the observed ~2.05 s "a measured margin, not a bound" | V30 |
| Transaction body and fsync | nothing in the store | **not bounded** |

So the durable tail is unbounded in the worst case, and the 10 s margin between `InvokeWait` and `writeTimeout` is a margin, **not a bound**. The fix, a context-accepting `Commit`/`AppendIntent` that is cancellable before the durable commit (all-or-nothing, with uncertain outcomes reconciled and never blindly retried), is ratified policy `D-WORLD-37` = A and is owned by **row 23's policy tranche**. This row does **not** change `store.Commit` (R-106-11).

### What the client sees when the durable tail overruns (revision 1, objection 2(b))

- **Commit lands, response written in time**: `A2AResult`.
- **Commit lands after the write deadline.** `net/http`'s `WriteTimeout` does not stop the handler. The handler finishes, the commit stays durable, and the response write fails, so the client sees a **transport error** (a closed connection or read timeout), not a JSON-RPC error. This is `net/http`'s documented behaviour, not re-measured here; M6 turns it into `TestA2AResendAfterWriteTimeout` on a real loopback server. The client is therefore **never sent a "failed" body for a commit that landed**. It is sent nothing. Resending the same task id returns the committed result from the journal (step 1b, `Reconciled: true`). `TestDispatchRefuses/R13_resend_reconciles_committed` pins that, and asserts the transition did not re-run.
- **The deadline fires *during* `Commit` and the commit lands**: `Dispatch` returns success, and ctx is not re-checked after `Commit` (`TestCommitBoundary/deadline_during_commit_reports_success`, killed by MUT-POSTCOMMIT-CTX).
- **`Commit` returns an untyped error**: `-32603 invocation outcome is not confirmed; resend the same task id`. The wire says *unconfirmed*, never *failed*. `R16_unconfirmed_then_reconciled` makes `Commit` land and then report an error. The first answer is `UnconfirmedError`, and the resend reconciles to the landed commit.
- **The resend finds an intent with no outcome**: `-32603 invocation was not committed; send a new task id` (R15). This holds because the store is single-connection and single-process-writer (writer lock, row 107 F8): the receipt read cannot run while that id's commit is in flight in this process, and no other process can hold the writer. An intent with no outcome therefore did not commit. **Stated limit:** this relies on the single-connection pool, and it becomes a residual if row 23 or any later row widens the pool (R-106-11).
- **Honest gap**: a client that never resends gets no answer after an overrun. The journal holds the truth, and there is no push channel (the card says `pushNotifications: false`).

### The pure core (M3): `planInvocation` — the kernel laws, transcribed

```go
type plan struct {
	Objects []store.Object      // input, output, record — all world/*/v1 semantic IDs
	Next    store.World
	Entry   store.LogEntry
	Intent  store.JournalIntent
}
func planInvocation(w store.World, id, episodeID string, d transitionreg.Descriptor,
	input, output []byte, logicalTime int64) plan
```

The **skill ID is `d.ID`**, the descriptor's own `ID string` field (`host/transitionreg/transitionreg.go:22`). It is the same field the A2A card emits as `skills[].id` (`host/projection/projection.go:240`, `"id": d.ID`) and the same key `transitionreg.Bind` looks up (`snap.Lookup(id)`, `bind.go:65`). So the planner needs no separate `skillID` parameter: the descriptor it receives *is* the one bound under that ID (V32). The record's `skillId` is written as `SkillID: d.ID` (`host/coordinator/plan.go:79`), and `TestPlanRecordSkillID` kills MUT-SKILLID.

- Objects: `world/invocation-input/v1` (canonical input bytes), `world/invocation-output/v1` (output bytes), `world/invocation-record/v1` (the canonical JSON record: invocation ID, episode, **skill ID = `Descriptor.ID`**, `TransitionFn`, `Interpreter`, `SemanticsEpoch`, input ref, output ref). Each object's `Hash = SHA-256(payload)` and `InterfaceHash = SHA-256(semanticID)`, the `replay.SourceObject` convention.
- **Law `applyRevision` (V21):** `Next.Revision == w.Revision + 1`, `Next.StateRoot == outputRef`, `Next.LogHead == Entry.EntryHash`.
- **Law `proposalMatchesWorld` (V21):** the commit's `ObservedHead == w.Ref`. `store.Commit`'s CAS turns a moved head into R14, so the verify contract is enforced *by the store at the commit boundary* and not only by the planner.
- Entry: `EntryIndex = w.Revision + 1`, `PrevEntryHash = w.LogHead`, `TransitionFn`/`Interpreter`/`SemanticsEpoch` from the descriptor, `WrittenBy = "coordinator:a2a"`, `TransitionRef = recordRef`. `EntryHash = SHA-256(canonical JSON of header + TransitionRef)`. `Next.Ref = SHA-256(canonical JSON of {revision, stateRoot, logHead})`.
- Intent: every field copied from `Next`/`Entry`, as `bindCommitIntentTx` compares them (V17).

Every law has a named test (AC-LAW-REV, AC-LAW-STATE, AC-LAW-LOG, AC-LAW-PREV, AC-LAW-OBSERVED) and a mutation (see Non-Vacuity). The planner is pure Go, not `.ail`. S1 applies to `world/`. A Z3-verified `.ail` twin of `applyRevision` already exists, and the Go transcription is bound to it by named law tests. The broker's decision law uses the same arrangement (S3 ledger note).

### What is recorded for replay (clause 3)

Each committed invocation stores three immutable objects: exact input bytes, exact output bytes, and the record. The log entry pins `(TransitionFn, Interpreter, SemanticsEpoch)`. That is everything needed to re-run: source object + archived interpreter + input → byte-compare output. Slice 1 proves it with **AC-REPLAY**: the test reads the committed record back through the store and re-executes through `capsule.RunContext`, and the output bytes must be equal. Teaching `host/replay.Engine` to pass `Args` is R-106-7. Replay's current `runPinnedTransition` passes none, and it stages a world library this row deliberately does not use. No effect results exist to record, because R8 refuses effectful transitions.

### Session authority

Only the session resolved by the daemon's ONE resolver (row 39, `authority.ResolveContext`) can reach `Dispatch`. The binding's `EpisodeID` names the journal ID and the broker binder. Its `Caps` feed **both** the snapshot filter (existing `Allowed`) and `Bind`'s `broker.Allows` re-check against the *same* captured capability snapshot. The coordinator has no credential path, no header reading, and no ambient caps. `Now` is logical time for journal payloads only, and capability expiry uses the snapshot's captured `now`.

### `/a2a/` wiring (M5): replacing the constant `-32603`

The admission sequence is unchanged up to the authorized-skill branch (`projection.go:314-320`). There the handler now:

1. requires `len(params.Message.Parts) == 1 && Parts[0].Type == "data" && Parts[0].Data != nil`, else `-32602 invalid params`;
2. reads optional `params.Metadata["transition_fn"]` (a string HashRef; unparsable → `-32602 invalid params`);
3. calls `h.coord.Dispatch(ctx, Call{Request: req, …})` with the **same** `transitionreg.Request` it admitted against. `allowedDescriptors` is refactored to return the `Request` so there is one snapshot;
4. on success writes `protocol.A2AResult(w, req.ID, task)`, where `task = {"id": TaskID, "status": {"state": "completed"}, "artifacts": [{"parts": [{"type": "data", "data": Output}]}], "metadata": {"invocation_id", "world_ref", "entry_index"}}`.

When the daemon has **no** pinned interpreter (`cfg.AilangBin == ""`, V16), it has no archive and so no runner. It wires `Coordinator: nil`, and the authorized branch keeps today's constant `-32603 notAvailableMessage`. That behaviour is preserved and tested (AC-UNPINNED).

**Error → JSON-RPC mapping (all on HTTP 200 via `protocol.A2AError`; every message is a constant):**

| Refusal | Code | Message (constant) |
|---|---|---|
| R1 invalid call, bad parts, bad `transition_fn` | -32602 | `invalid params` (existing `msgInvalidParams`) |
| R2 absent (raced out after admission), R3 access denied | -32602 | `not authorized` (existing `msgNotAuthorized`; never names the gap) |
| R5 proposal pin mismatch | -32602 | `proposal does not match the registered transition` |
| R13 resent task id, committed | — | **success**: `A2AResult` carrying the originally committed result (reconciled from the journal) |
| R15 resent task id, intent without outcome | -32603 | `invocation was not committed; send a new task id` |
| R16 `Commit` returned an untyped error | -32603 | `invocation outcome is not confirmed; resend the same task id` |
| `*store.DuplicateInvocationError` (defensive; unreachable after step 1b) | -32602 | `task id already used in this session` |
| R4 bind error, R6 source error, other store errors | -32603 | `transition invocation is not available in this daemon` (existing `notAvailableMessage`) |
| R7 no world | -32603 | `no world is selected; commit a genesis world first` |
| R8 effects declared | -32603 | `transitions that declare effects cannot be invoked in this daemon` |
| R9 incompatible | -32603 | `transition does not implement the invocation calling convention` |
| R10 execution failed | -32603 | `transition execution failed` |
| R11 deadline/cancel | -32603 | `invocation exceeded its deadline` |
| R12 output invalid | -32603 | `transition output is not a JSON object` |
| R14 conflict | -32603 | `world head moved during invocation; not committed; send a new task id` |

The mapping is a single `switch` over `errors.As` in `projection.dispatchError`, and each arm is a mutation unit (see Non-Vacuity). No stderr, path, skill name or store detail is ever interpolated.

**No new wire protocol**: A2A `tasks/send` request/response only, emitted through `protocol.A2AError` / `protocol.A2AResult` from the pinned `serveapi/protocol` (V13). No streaming and no push; the card keeps `streaming: false`.

## Milestones (each ≤ ~150 production code lines)

| M | Scope | Files | Est. code LOC |
|---|---|---|---|
| **M1** | `capsule.RunContext` + `Entry.Args` + `--args-file` + relax env; `Run` delegates | `host/capsule/capsule.go`, `capsule_test.go` | ~35 |
| **M2** | `broker.OpenBinder` / `SessionBinder` (TR.C count stays 3; add the NEG detector control `NEG-session-binder`) | `host/broker/binder.go`, `binder_test.go`, `invoke_boundary_test.go` (+1 control row) | ~20 |
| **M3** | coordinator pure core: `planInvocation`, record codec, refusal types + the five law tests | `host/coordinator/plan.go` (94), `errors.go` (28), law tests | 122 (measured) |
| **M4a** | seams + construction: `Store`/`Runner`/`BinderFor`, `Config`, `New`, `Call`/`Result`, `InvocationID`, `validateCall`, `classifyExec`, `parseOutput` | `host/coordinator/coordinator.go` (first half) | ~90 |
| **M4b** | `Dispatch` composition (steps 1–10), fake-runner tests for every refusal, real-interpreter tests (echo / zero-arity / replay) | `host/coordinator/coordinator.go` (Dispatch), `coordinator_test.go` | ~100 |
| **M4c** | (revision 1) step 1b reconciliation: `GetReceipt` lookup, `committed` read-back, `NotCommittedError`/`UnconfirmedError`, bare-conflict rule, no post-commit ctx check; tests R13/R15/R16 + `deadline_during_commit_reports_success` | `host/coordinator/coordinator.go`, `errors.go`, `coordinator_test.go` | ~55 (M4a+M4b+M4c = 231 + 40 errors.go, measured) |
| **M5** | `/a2a/` wiring: `InvokeWait`, `Coordinator` config, `allowedDescriptors` returns the `Request`, `dispatchError` mapping, `A2AResult`; daemon constructs the runner and coordinator iff an interpreter is archived; card description sentence updated | `host/projection/projection.go`, `projection_test.go`, `host/daemon/daemon.go` | ~120 |
| **M6** | Daemon end-to-end (real pinned interpreter, in-process `httptest` recorder, no sockets): genesis commit → `PublishSet` the echo transition → `/a2a/ tasks/send` → result + `GET /v1/log/{i}` shows the entry → AC-REPLAY; `docs/QUICKSTART.md` §7 "invoke a published transition" (S7), marked *attended — pending first verbatim run* | `host/daemon/invoke_e2e_test.go`, `docs/QUICKSTART.md` | ~30 prod-doc + test |

**First landable slice: M1–M5**, with M6 in the same iteration if velocity holds. **The prototype already implements M1–M4c** (see Prototype). **Split point if not**: M1–M4c land as "coordinator exists, tested, no production caller", and `/a2a/` stays constant until M5. That is honest, but the clause-6 value appears only at M5.

## Acceptance Criteria

All commands run with `export PATH=/opt/homebrew/bin:$PATH; export AILANG_BIN=$HOME/.pinned-ailang/ailang`.

| AC | Criterion | Command |
|---|---|---|
| AC-CAP-ARGS | `RunContext` with `Args` delivers the JSON string to `main(input)` and returns its stdout | `go test ./host/capsule/ -run 'TestRunContextArgs'` |
| AC-CAP-CTX | A cancelled `ctx` kills the capsule before its exec allowance | `go test ./host/capsule/ -run 'TestRunContextCancel'` |
| AC-CAP-RELAX | A source whose module header ≠ `host/capsule/main` runs (publish/run shapes agree; closes V9) | `go test ./host/capsule/ -run 'TestRunContextRelaxedModule'` |
| AC-CAP-LEGACY | All pre-existing capsule tests stay green | `go test ./host/capsule/` |
| AC-TRC | TR.C stays green with exemption count 3; the new NEG control holds | `go test ./host/broker/ -run TestRegistryDispatchBindingBoundary` |
| AC-BINDER | `OpenBinder(...).Bind` returns a `BoundInvoker` whose `Request` refuses an undeclared triple | `go test ./host/broker/ -run TestOpenBinder` |
| AC-LAW-REV / -STATE / -LOG / -PREV / -OBSERVED | Each kernel law named above holds on the plan | `go test ./host/coordinator/ -run 'TestPlanLaw'` |
| AC-HAPPY | echo transition: `Dispatch` returns the output object, commits exactly one entry at revision+1, and the head moves | `go test ./host/coordinator/ -run 'TestDispatchEchoRealInterpreter'` |
| AC-R1…R16 | Each refusal branch returns its typed error. Every pre-boundary refusal leaves **no durable mutation** (selected head unchanged, no receipt). R14 asserts a *bare* conflict (not `UnconfirmedError`). R15 and R13 assert the transition is **not re-executed** on a resend | `go test ./host/coordinator/ -run 'TestDispatchRefuses'` |
| AC-SKILLID | (revision 1) the record's `skillId` equals `Descriptor.ID` and is non-empty | `go test ./host/coordinator/ -run 'TestPlanRecordSkillID'` |
| AC-RECONCILE | (revision 1) a resent task id returns the originally committed result (`Reconciled: true`, same world/entry/record/output) without executing; a commit that landed but reported an error reconciles on resend | `go test ./host/coordinator/ -run 'TestDispatchRefuses/(R13_resend\|R16_unconfirmed)'` |
| AC-POSTCOMMIT | (revision 1) a deadline that expires during the durable steps does not turn a landed commit into an error | `go test ./host/coordinator/ -run 'TestCommitBoundary/deadline_during_commit'` |
| AC-RESEND-WIRE (M6) | (revision 1) on a real loopback server with a `writeTimeout` shorter than an injected `Commit` delay, the first request gets a transport error, the commit is durable, and a resend of the same task id gets `A2AResult` with the committed output | `go test ./host/daemon/ -run 'TestA2AResendAfterWriteTimeout'` |
| AC-BOUNDARY | Cancel immediately before step 10 → no intent, head unchanged; a completed Dispatch → exactly one intent + one outcome (receipt) | `go test ./host/coordinator/ -run 'TestCommitBoundary'` |
| AC-REPLAY | Re-executing the committed record's pins + input reproduces the committed output bytes | `go test ./host/coordinator/ -run 'TestReplayCommittedInvocation'` |
| AC-INCOMPAT | Residual 8(c): the zero-arity fixture (checks, V10) returns `*IncompatibleError` | `go test ./host/coordinator/ -run 'TestDispatchIncompatibleRealInterpreter'` |
| AC-A2A-MAP | Every refusal maps to its code + constant message; success emits `A2AResult` with the task shape | `go test ./host/projection/ -run 'TestA2ADispatch'` |
| AC-UNPINNED | With no coordinator configured, an authorized skill still gets the constant `-32603` | `go test ./host/projection/ -run 'TestA2AAuthorizedSkillNoCoordinator'` |
| AC-E2E | Publish-through-`PublishSet` then invoke over `/a2a/` against the real pinned interpreter; the log entry is readable over `GET /v1/log/{i}` | `go test ./host/daemon/ -run 'TestA2AInvokesPublishedTransition'` |
| AC-VERIFY | Whole gate | `./scripts/verify_ail.sh && go vet ./... && go test ./...` |

## Non-Vacuity — one named RED mutation per gate and per refusal branch

The prototype column reports the mutations executed in this worktree (see **Prototype**). "planned" means the target code is not prototyped yet, so the executor must run the mutation.

| Mutation | File: exact change (applied verbatim by the harness) | Killing test | Prototype |
|---|---|---|---|
| MUT-ARGS-DROP | `host/capsule/capsule.go`: `args = append(args, "--args-file", argsFile)` → `_ = argsFile` | `TestRunContextArgs` | **KILLED** |
| MUT-CTX-BG | `capsule.go`: `context.WithTimeout(parent, r.execTimeout)` → `context.WithTimeout(context.Background(), r.execTimeout)` | `TestRunContextCancel` (red after 29.3 s: the child ran to completion) | **KILLED** |
| MUT-RELAX-DROP | `capsule.go`: delete `, "AILANG_RELAX_MODULES=1"` from `cmd.Env` | `TestRunContextRelaxedModule` | **KILLED** |
| MUT-TRC-SCOPE | `host/broker/binder.go`: add `func (b *SessionBinder) Raw() *Session { return b.s }`, and `var _ = (*broker.Session)(nil)` in `host/coordinator/coordinator.go` | `TestRegistryDispatchBindingBoundary/outside_broker_is_clean` (TR.C scans the new package) | **KILLED** |
| MUT-REV | `host/coordinator/plan.go`: `Revision: w.Revision + 1,` → `+ 2,` | `TestPlanLawRevision` | **KILLED** |
| MUT-STATE | `plan.go`: `StateRoot: out.Hash` → `StateRoot: in.Hash` | `TestPlanLawStateRoot` | **KILLED** |
| MUT-LOGHEAD | `plan.go`: `LogHead: entryHash}` → `LogHead: rec.Hash}` | `TestPlanLawLogHead` | **KILLED** |
| MUT-PREV | `plan.go`: `PrevEntryHash: w.LogHead` → `PrevEntryHash: w.Ref` | `TestPlanLawPrevEntry` | **KILLED** |
| MUT-OBSERVED | `plan.go`: `ObservedHead: w.Ref, Objects` → `ObservedHead: hashref.HashRef{}, Objects` | `TestPlanLawObservedHead` | **KILLED** |
| MUT-R1-TASK | `coordinator.go`: `if !validTaskID(call.TaskID) {` → `if false && !validTaskID(call.TaskID) {` | `TestDispatchRefuses/R1_bad_task_id` | **KILLED** |
| MUT-R1-INPUT | `coordinator.go`: `if call.Input == nil {` → `if false {` | `TestDispatchRefuses/R1_nil_input` | **KILLED** |
| MUT-AMBIENT-CAPS (R3) | `coordinator.go`: pass `broker.NewCapabilitySnapshot([]broker.Capability{{Effect: "Admin", …}}, 0)` to `transitionreg.Bind` instead of `call.Request.Caps` (ambient authority) | `TestDispatchRefuses/R3_access_denied` | **KILLED** |
| MUT-R5 | `coordinator.go`: `p.TransitionFn = *call.PinnedFn` → `_ = call.PinnedFn` | `TestDispatchRefuses/R5_pin_mismatch` | **KILLED** |
| MUT-R6-ABSENT | `coordinator.go`: the source `if !ok {` guard → `if false {` | `TestDispatchRefuses/R6_absent_source` | **KILLED** |
| MUT-R6-HASH | `coordinator.go`: `\|\| sum != d.TransitionFn {` → `\|\| (false && sum != d.TransitionFn) {` | `TestDispatchRefuses/R6_corrupt_source` | **KILLED** (first form `\|\| false` did not compile, `sum` unused: re-expressed, not counted twice) |
| MUT-R7 | `coordinator.go`: the `SelectedHead` `if !ok { return …WorldAbsentError` → `if false {` | `TestDispatchRefuses/R7_no_world` | **KILLED** (against the final code, which separates "no head" from "head without a world row". See the Prototype notes) |
| MUT-R8 | `coordinator.go`: `if len(d.DeclaredEffects) != 0 {` → `if false {` | `TestDispatchRefuses/R8_effects_declared` | **KILLED** |
| MUT-R9-FAKE | `coordinator.go`: `return &IncompatibleError{Err: err}` → `return &ExecutionError{Err: err}` | `TestDispatchRefuses/R9_incompatible` | **KILLED** |
| MUT-R9 (real) | same change | `TestDispatchIncompatibleRealInterpreter` (pinned v0.41.0, zero-arity `main`) | **KILLED** |
| MUT-R10 | `coordinator.go`: `return Result{}, classifyExec(err)` → `_ = classifyExec(err)` | `TestDispatchRefuses/R10_exec_failed` | **KILLED** |
| MUT-R11-EXEC | `coordinator.go`: `if cerr := ctx.Err(); cerr != nil {` → `if cerr := ctx.Err(); false && cerr != nil {` | `TestDispatchRefuses/R11_cancelled_during_exec` | **KILLED** |
| MUT-R11-BOUNDARY | `coordinator.go`: `if err := ctx.Err(); err != nil { // R11` → `… false && err != nil …` | `TestCommitBoundary/cancel_before_commit` | **KILLED** |
| MUT-R12 | `coordinator.go`: `err := json.Unmarshal(out, &obj); err != nil \|\| obj == nil` → `err := json.Unmarshal(out, new(any)); err != nil` | `TestDispatchRefuses/R12_output_not_object` | **KILLED** |
| MUT-R13-JOURNAL | `plan.go`: `InvocationID: id, ObservedHead: w.Ref, Objects` → `InvocationID: "", …` (commit unbound from its intent) | `TestCommitBoundary/receipt` | **KILLED** |
| MUT-R14-REBASE | `coordinator.go`: before `planInvocation`, re-read `SelectedHead`/`GetWorld` and plan against the new world (silent re-base) | `TestDispatchRefuses/R14_head_moved` | **KILLED** |
| MUT-SKILLID | (rev 1) `plan.go`: `SkillID: d.ID,` → `SkillID: "",` | `TestPlanRecordSkillID` | **KILLED** |
| MUT-RECONCILE-SKIP | (rev 1) `coordinator.go`: `if seen {` → `if false && seen {` | `TestDispatchRefuses/R13_resend_reconciles_committed` | **KILLED** |
| MUT-R15 | (rev 1) `coordinator.go`: `if rc.State == store.ReceiptResolved {` → `if true {` | `TestDispatchRefuses/R15_resend_not_committed` | **KILLED** |
| MUT-R16 | (rev 1) `coordinator.go`: `return Result{}, &UnconfirmedError{InvocationID: id, Err: err} // R16` → `return Result{}, err // R16` | `TestDispatchRefuses/R16_unconfirmed_then_reconciled` | **KILLED** |
| MUT-R14-TYPED | (rev 1) `coordinator.go`: `if store.IsConflict(err) { // R14` → `if false { // R14` (conflict becomes "unconfirmed") | `TestDispatchRefuses/R14_head_moved` (asserts a bare conflict) | **KILLED** |
| MUT-POSTCOMMIT-CTX | (rev 1) `coordinator.go`: insert `if err := ctx.Err(); err != nil { return Result{}, err }` after a successful `Commit` | `TestCommitBoundary/deadline_during_commit_reports_success` | **KILLED** |
| MUT-ID-NS | `coordinator.go`: `return "a2a:" + episodeID` → `return "effect:" + episodeID` | `TestDispatchEchoRealInterpreter` (the journal refuses the namespace) | **KILLED** |
| MUT-R4-BIND | `Dispatch`: swallow a binder error (fake `BinderFor` whose `Bind` fails) | `TestDispatchRefuses/R4_bind_error` (to add) | planned. The prototype has no R4 test; a validated registry cannot produce a failing `Bind`, so it needs a fake binder |
| MUT-MAP-<Rn> (×12) | `projection.dispatchError`: map the arm to `codeInternal, notAvailableMessage` (for -32602 arms) or to `codeInvalidParams` (for -32603 arms) | `TestA2ADispatch/<Rn>` | planned (M5) |
| MUT-ONE-SNAPSHOT | `projection.A2A`: pass a freshly captured `NewRequest` to `Dispatch` instead of the admitted one | `TestA2ADispatch/registry_moved_after_admission` | planned (M5) |
| MUT-UNPINNED | `daemon.go`: construct a coordinator even without an archived interpreter | `TestA2AAuthorizedSkillNoCoordinator` (daemon variant) | planned (M5) |
| MUT-SKIP-SOCKET | M6: stop tracking `ConnState` closure in the deadline test | `TestA2AInvokeDeadlineClosesSocket` | planned (M6) |

**R2** (`TransitionAbsentError`) has no coordinator-side branch: the refusal is `transitionreg.Bind`'s own, and existing `host/transitionreg` tests cover it. `TestDispatchRefuses/R2_absent` pins that the coordinator propagates it typed.

## Conflict Surface

- `host/capsule/capsule.go`: `Run` now delegates to `RunContext`, and the env gains `AILANG_RELAX_MODULES=1`. This strictly widens which module headers run. Existing tests (`capsule_test.go`, `cleanup_test.go`) must stay byte-green. The cleanup tests exercise `collectOutput` via fakes, so they are unaffected.
- `host/broker`: a new exported `SessionBinder`. TR.C's detector cannot see it (not a `Session` selector), so review must: the file's doc comment is the contract, and MUT-TRC-SCOPE (executed, KILLED) proves TR.C still catches a raw-session leak from the new package. TR.C's inside-exemption count stays **3** (measured: the gate is green with `binder.go` present, V26).
- `host/store/context_roots_test.go` (`contextRootPins`, S8-class inventory): **did not fire**. `capsule.Run` keeps its single `context.Background()` root (now as the argument to `RunContext`), and `host/coordinator` has none. The full `host/store` suite is green with the prototype (V26).
- `host/projection/projection.go`: the authorized branch changes from a constant refusal to dispatch. **Every existing `/a2a/` test that asserts the constant for an authorized skill must now construct the handler without a coordinator** (AC-UNPINNED keeps that behaviour pinned). `allowedDescriptors` changes its return type. The card route consumes only `.Allowed()`, so card tests are unaffected.
- `host/daemon/daemon.go`: construction order. The runner needs the archive the daemon already builds (V16). The daemon's card `Description` string changes (a test may pin it: grep before editing). The two inventory gates that fired for row 107 (`contextRootPins` in `host/store`, S8-class) will fire again if `Dispatch` introduces a `context.Background` root. The design has none (the ctx always comes from the request), so a red here means a design violation, not an inventory update.
- `host/store`: **no change.** The journal and `Commit` are used as landed.
- `docs/QUICKSTART.md`: a new §, attended-verbatim pending (S7).

## Axiom Compliance (coding-standards.md)

- **S1**: no new `world/` function. The Go planner transcribes `applyRevision`/`proposalMatchesWorld` (already Z3-verified in `world/transitions.ail`) under named law tests with mutations.
- **S2**: all effects (subprocess, store) stay in `host/`, and the transition runs with `--caps ""`.
- **S3**: "why not a package?" is answered above (host boundary, shared by A2A and MCP). The demo transition is a test fixture, not kernel.
- **S4**: this doc contains no `.ail` code block. The only AILANG it states is the calling-convention signature, which was run on the pinned binary (V7, V25). Transition fixtures are inline Go strings executed by the pinned interpreter in tests (the capsule precedent). They are not `design_docs` snippets, and `verify_ail.sh` does not sweep `host/` (V28).
- **S5**: every syntax claim here (`--args-file`, `main(input: string) -> string`, MOD010/MOD014/ARG_DECODE_MISMATCH texts) was run on the pinned v0.41.0 (V7–V12).
- **S6**: every refusal branch has its own mutation. TR.C's count is asserted, not assumed. AC-UNPINNED guards the null configuration.
- **S7**: QUICKSTART §7 in M6. The card description states the calling convention.

## Residuals (named owners)

- **R-106-1**: first *useful* published skill (production card content). Owner: row 93, as packages.
- **R-106-2**: world-library staging (row 107 Residual 6). This needs a descriptor field pinning library bytes, so a registry interface change. Owner: a new row, raised only when a real skill needs `world/*`.
- **R-106-3**: effect bridge from the capsule child to `BoundInvoker.Request`, which would make R8 unnecessary. Owner: a new row; clause 3's "effect results are recorded" then applies to transition-initiated effects.
- **R-106-4**: optional publish-time invocability check (widen `EnsureSourceLoadable` with a probe run). Owner: executor's judgement or a later row. Slice 1 reports R9 at invocation.
- **R-106-5**: JSON-Schema validation of input/output against the descriptor schemas. This needs dependency admission. Owner: a later row.
- **R-106-6**: `protocol.Invoker` adapter for MCP. Owner: row 108.
- **R-106-7**: `host/replay.Engine` passes `Args` for invocation entries. Owner: the next replay row. AC-REPLAY covers slice 1.
- **R-106-8**: the transition sees only its input, not the current world state. Passing `{state, input}` needs a convention for state-object payloads, which today are arbitrary bytes. Owner: row 93's design, if its skills need state.
- **R-106-9**: ~~idempotent replay of the result for a retried identical task~~. **Resolved in revision 1** (step 1b, M4c): a resent task id is answered from the journal.
- **R-106-11** (revision 1, objection 2). **Unbounded durable tail.** `GetReceipt`, `AppendIntent` and `Commit` take no context (V30), so the invocation deadline does not bound them. Owner: **queue row 23's policy tranche**, ratified `D-WORLD-37` = A ("`Commit` accepts a ctx and may be cancelled before the durable commit, all-or-nothing; an uncertain outcome is reconciled, never blindly retried"). This row does not change `store.Commit`. **When it lands, this row adopts it as follows:** pass the invocation `ctx` unreplaced into the ctx-accepting `GetReceipt`/`AppendIntent`/`Commit`. Delete the separate pre-boundary `ctx.Err()` check if `Commit`'s own pre-durable cancellation subsumes it (keep MUT-R11-BOUNDARY's test either way). Map a pre-durable cancellation to R11 ("not committed"). Map an uncertain outcome to R16, reconciled by the step-1b resend. Keep "no ctx check after a successful `Commit`". It also owns the R15 inference's dependence on the single-connection pool: if the pool widens, R15 ("not committed") must become "indeterminate; resend later" until the in-flight commit is excluded.
- **R-106-10**: an invocation does not debit the session's budget. The binding's caps are an immutable per-session snapshot (row 39). Clause 3's budget check applies to *effects*, and slice 1 has none (R8). Owner: R-106-3.

## Estimate honesty

~2d for M1–M6, matching the row. The prototype below implements M1–M4 in real code; M5 and M6 are design only. The riskiest unknown was the capsule argument channel, and V7 plus the prototype retired it. The next risk is M5's test churn in `projection_test.go` (existing authorized-skill tests).

## Open decisions for the human

**None required.** The design uses no new wire, no new dependency, no new charter-level policy. The one policy-shaped choice is that invocation, not publication, reports incompatibility (Residual 8(c)). Row 107 assigned that choice to this row's designer, and it is reversible (R-106-4).

## Related Documents

- `design_docs/implemented/w-a2a-session-projection.md`: "Deferred to invocation" (P4 dispatch half, `protocol.Invoker`, AC13, socket closure), which this row discharges except for OS socket closure (Bounded wait below).
- `design_docs/implemented/w-transition-registry-production-publisher.md`: Residuals 1, 6, 7, 8.
- `design_docs/implemented/w-effect-journal.md`: the intent/outcome journal reused at step 10.
- `world/transitions.ail`: the kernel laws transcribed by `planInvocation`.
- `design_docs/planned/w-mcp-dispatch-projection.md`: row 108, the second consumer of `Dispatch`.

### AC13 carry-over (bounded waits, socket closure)

AC13's invocation half is discharged at the Go level by AC-CAP-CTX (a blocked capsule terminates at the bound), AC-BOUNDARY (both sides of the commit boundary) and AC-R11. The **OS-level socket-closure** leg (`MUT-SKIP-SOCKET`) needs a real listener. It runs in M6 as `TestA2AInvokeDeadlineClosesSocket` on a loopback `http.Server` with `ConnState` tracking, using a fixture transition that loops forever and an `InvokeWait` of 1s. Because loopback tests run for real outside the sandbox, this is an executor obligation and not a residual.

## Prototype (executed in this worktree, iteration 196)

**What is prototyped (real code, tests, mutations):** M1 (capsule `RunContext` + `Args` + relax env), M2 (`broker.OpenBinder`), M3 (pure planner + refusal types), M4a/M4b (`Coordinator.New`/`Dispatch`). **What is not:** M5 (`/a2a/` wiring, `dispatchError` mapping, daemon construction, `InvokeWait`), M6 (daemon E2E through `PublishSet` over `/a2a/`, socket-closure test, QUICKSTART), the `NEG-session-binder` TR.C control row, and the R4 test. In the prototype tests the registry revision is seeded directly (the `seedTransitionRegistry` shape), not through `PublishSet`. Publishing through `PublishSet` is M6's obligation (Residual 8(b)).

| File | Status | Code lines (`grep -vcE '^\s*(//\|$)'`) |
|---|---|---|
| `host/capsule/capsule.go` | modified (+27/−4 per `git diff --stat`) | n/a |
| `host/capsule/runcontext_test.go` | new: `TestRunContextArgs`, `TestRunContextCancel`, `TestRunContextRelaxedModule` | test |
| `host/broker/binder.go` | new | 7 |
| `host/broker/binder_test.go` | new: `TestOpenBinder` | test |
| `host/coordinator/errors.go` | new | 40 (rev 1; was 28) |
| `host/coordinator/plan.go` | new | 94 |
| `host/coordinator/coordinator.go` | new | 231 (rev 1; was 188), split M4a/M4b/M4c |
| `host/coordinator/coordinator_test.go` | new: 5 law tests + `TestPlanRecordSkillID`, 3 real-interpreter tests, 18 refusal subtests (rev 1: R13 → resend-reconciles, + R15, R16; R14 asserts a bare conflict), 3 boundary subtests (rev 1: + `deadline_during_commit_reports_success`), construction test | test |

**Verification rows for the prototype:**

| # | Claim | Command → observed |
|---|---|---|
| V24 | Vet clean | `go vet ./...` → no output, `VET-OK` |
| V25 | Coordinator tests pass, with the real-interpreter tests **executed, not skipped** | `AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./host/coordinator/ -v` → `--- PASS: TestDispatchEchoRealInterpreter (1.04s)`, `TestReplayCommittedInvocation (1.06s)`, `TestDispatchIncompatibleRealInterpreter (1.00s)`. All 16 `TestDispatchRefuses/*`, both `TestCommitBoundary/*` and the 5 `TestPlanLaw*` PASS. The echo transition's committed output is `{"echo":{"msg":"hi"}}` |
| V26 | No regression anywhere, loopback daemon/projection tests included (this lane is not sandboxed, so they ran for real) | `AILANG_BIN=… go test ./host/... ./cmd/...` → 24 × `ok`, 0 `FAIL`, `rc=0`, incl. `host/broker 77.4s` (TR.C), `host/capsule 44.1s`, `host/daemon 25.9s`, `host/projection 3.5s`, `host/store 11.3s`. `go list ./... \| grep -vE '/(host\|cmd)/'` → empty (host+cmd is the whole module) |
| V27 | Milestone sizes | `grep -vcE '^\s*(//\|$)'` → `coordinator.go:188`, `plan.go:94`, `errors.go:28`, `binder.go:7` |
| V28 | `verify_ail.sh` sweeps only `design_docs/` and `world/` | `grep -n "ROOTS=" -A3 scripts/verify_ail.sh` → `"design_docs\|."`, `".\|world"`; exact manifest `LEG1_MODULES` at `:169` |
| V29 | Mutation harness: every file restored byte-identical | `/tmp/iter196mut/mut.py` backs up with `shutil.copy`, applies each edit only if its anchor occurs exactly once, runs `go test -count=1 <pkg> -run <killer>`, restores from the backup, and asserts `filecmp.cmp(..., shallow=False)`. After the run, `git status --short` lists only this row's files (capsule.go modified; binder*, runcontext_test.go, host/coordinator/, this doc new) |

**Mutation results (revision 1): 32 of 32 executed mutations KILLED, 0 survived.** The whole harness was re-run after the revision-1 code change, including all 26 original mutations, whose anchors all still matched exactly once. Earlier, the original round was **26 of 26 KILLED, 0 survived** (the table in Non-Vacuity; 5 planned rows are M5/M6/R4 obligations). Two honest notes from the run. (1) `MUT-R6-HASH`'s first form did not compile (`sum` unused); it was re-expressed and then killed. (2) Before execution, a read of the first draft found that `MUT-R7` would survive there: a second `GetWorld` guard also returned `WorldAbsentError`. This was reasoned, not executed against that draft. The draft was changed so that a head without a world row is store damage (untyped) and not "no world". The typed branch is now single and load-bearing.

**Design observations the prototype surfaced (folded into the design above):**

- A parent-deadline expiry inside the capsule surfaces as `*capsule.TimeoutError{Limit: execTimeout}` (the capsule labels any `DeadlineExceeded` with its own allowance). The coordinator therefore checks **its own** `ctx.Err()` first on a runner error (R11 before R9/R10). MUT-R11-EXEC proves that order is load-bearing.
- Input canonicalisation is `json.Marshal(map[string]any)` (sorted keys). Numbers arrive as `float64` from the upstream `A2AContent.Data` decoding, so integers above 2^53 lose precision *before* World sees them. That is an upstream wire-type property, recorded here and not fixed.

**Revision-1 re-verification.** `go vet ./...` → clean. `AILANG_BIN=… go test -count=1 ./host/coordinator/ ./host/capsule/ ./host/broker/` → coordinator `ok`, broker `ok`; capsule hit the pre-existing EPERM flake once (V34). A standalone re-run of `go test -count=1 ./host/capsule/` → `ok` (19.0 s).

## Quorum log

### Round 1 (iteration 196): BLOCKED 2/2 (gemini-3-1-pro, gpt6-astra present; glm and kimi ABSENT on the Ollama weekly limit)

**Objection 1 (gemini-3-1-pro), verbatim:** "The pure core signature `planInvocation` in M3 is missing the `skillID` parameter. Step 9 explicitly requires the `skill ID` to be written into the `world/invocation-record/v1` canonical JSON record. Because `transitionreg.Descriptor` does not store the skill ID (it is the registry key), the pure function has no way to read it and cannot construct the record object as specified."

*Controller measurement:* premise **FALSE**. `transitionreg.Descriptor` has `ID string` (`host/transitionreg/transitionreg.go:22`), and the prototype's `plan.go:79` already writes `SkillID: d.ID`. The doc never *said* that the skill ID is `Descriptor.ID`, so a reader could not see it.

*Changed:* the pure-core section now states that the skill ID is `d.ID`, and that this is the same field the card emits as `skills[].id` (`projection.go:240`) and the key `Bind` looks up (`bind.go:65`). The record bullet says **skill ID = `Descriptor.ID`**. The signature block now matches the prototype's real parameter order. New Verification row **V32** (greps with file:line). New **AC-SKILLID**, new test `TestPlanRecordSkillID`, and new mutation **MUT-SKILLID** (`SkillID: d.ID,` → `SkillID: "",`): **KILLED**.

**Objection 2 (gpt6-astra), verbatim:** "The claimed end-to-end invocation deadline stops at step 10. AppendIntent and Commit take no context, and checking ctx.Err() before calling them does not bound either call. The document provides no verified upper bound for their lock acquisition, database waits, or transaction completion. Consequently, the claim that a 20-second invocation deadline inside a 30-second write timeout guarantees an answer is unsupported and violates the bounded-waits gate."

*Controller measurement:* premise **TRUE**. `func (s *Store) Commit(c Commit) error` (`store.go:874`) opens `s.db.Begin()` with no context. The pool is `SetMaxOpenConns(1)` (`store.go:305`), so an in-process caller can wait for the single connection with no bound. `busy_timeout` is 2000 ms (`writer_lock.go:187`), which bounds SQLite lock waits only. The ratified policy `D-WORLD-37` = A is owned by row 23's policy tranche (not landed), and `store.Commit` must not change in this row.

*Changed:* (a) The claim is weakened to match the code. The deadline paragraph now says the invocation deadline bounds steps 1–9 and the pre-boundary check. It does **not** bound `GetReceipt`/`AppendIntent`/`Commit`, and "nothing in this design guarantees that a JSON-RPC answer is written before `writeTimeout`". The prior phrase "so the JSON-RPC answer can always be written" is **removed**, and the 10 s gap is called a margin, not a bound. A new "Bounded waits" table lists what bounds each wait and what does not (**V30**, with a positive control; `AppendIntent` was checked the same way: `journal.go:411`/`:420`, no ctx, `Begin()`). (b) A new section, "What the client sees when the durable tail overruns", covers the write timeout firing mid-commit. The handler finishes and the commit is durable, so the client gets a transport error, **never a "failed" body**. A resend of the same task id is answered from the journal (resolved ⇒ the committed result, and nothing re-runs). An untyped `Commit` error is answered **"not confirmed; resend the same task id"**, never "failed". A conflict is answered "not committed" only because the store proves the rollback (V17). This required new code, implemented and tested in the prototype (M4c): step 1b reconciliation, `NotCommittedError` (R15), `UnconfirmedError` (R16), the bare-conflict rule, and no ctx check after a successful `Commit`. It carries 5 new mutations (MUT-RECONCILE-SKIP, MUT-R15, MUT-R16, MUT-R14-TYPED, MUT-POSTCOMMIT-CTX), **all KILLED**, and 4 new/changed tests. The honest remaining gap is stated: a client that never resends gets no answer after an overrun. (c) New residual **R-106-11** names row 23's policy tranche (`D-WORLD-37` = A) as the owner of a ctx-accepting `Commit`/`AppendIntent`/`GetReceipt`, and states what this row adopts when it lands. `store.Commit` is untouched. (d) The prototype's `Dispatch` doc comment now says the durable steps take no context and are bounded only by the store. `go vet ./...` is clean. The full mutation harness (32 mutations, including the 26 originals) was re-run: 32/32 KILLED.

**Also found during revision 1 (not an objection):** V34, the pre-existing macOS EPERM flake in `host/capsule`'s `TestOrdinaryOverflowCarriesNoESRCH`. It fails 3/320 on both the pristine and the prototype `capsule.go` under the same load. Reported, not fixed here.
