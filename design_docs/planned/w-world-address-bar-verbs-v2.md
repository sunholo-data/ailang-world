# w-world-address-bar-verbs-v2: `why <target>` in the World app's address bar

- **Status**: Planned. Designed in iteration 396, **r0** (unattended designer). Child of
  [`w-world-address-bar-verbs.md`](w-world-address-bar-verbs.md) (r2, reviewed), queue row **243**,
  parent phase **V2**. V1 is **built and landed** (PR #582, merge `fb4d8b89`, judge PASS 91) and is
  the code this design measures, not the V1 plan. This doc carries **five explicit amendments** to
  the parent (§0.3) and **splits V2 into V2a/V2b** (§0.4, §8): the measured scope is ≈ 1,850 changed
  lines, above the parent's ≈ 1,600.
- **Authority**: **D-WORLD-107 = A** (the bar accepts operator verbs in the CLI's grammar; the row
  names `why <entry>`), **D-WORLD-109 = A** (the loop builds and judges V2 against fakes and test
  stores; Mark's exit run stays attended), **D-WORLD-89** R1–R5 untouched: V2 adds no POST, no
  write, no key path. Row 201's settled siblings are not reopened; V2 touches none of their files.
- **Item**: typing `why <target>` in the app's address bar answers with **the bytes
  `ailang-worldd why <target>` prints** (stdout, or its stderr line, and its exit status), produced
  by **one implementation** of the walk that the CLI and the daemon both call. The daemon serves it
  in-process through V1's `GET /surface/verb.json` — no shell-out, no loopback HTTP, no new route.
  Read-only; no World state changes.
- **Designer lane**: `claude:claude-opus-5-5`.
- **Estimated**: **V2a ≈ 950** (the seam: `host/why`, CLI repointed; no user-visible change) and
  **V2b ≈ 900** (the bar: admission, the in-process getter, the route branch, the page, the V1
  judge notes) changed lines, `git diff --numstat -- cmd host docs` **with git's default rename
  detection** (§8; W27 shows a whole-file move with renames counts 265, not 2,000). Calibration: V1
  planned 1,750 and landed 1,909 (W28), ≈ 1.09×.
- **Measurement base**: `200ea5d5b7c2944e561feeff16a4b2b3042333d1` (origin/dev, PR #583's merge;
  `cmd/` and `host/` byte-identical to V1's merge `fb4d8b89`, W1), worktree
  `/Users/voightkampff/dev/sunholo-data/ailang-world/.wt-world-iter396-row243-v2`, branch
  `design/world-iter396-row243-v2`, measured 2026-11-16.
- **Instruments**: the worktree's tree, `grep`/`sed`/`wc`, and one scratch git repo
  (`/tmp/v2rename.2oWw`, not in this repo) for the rename-detection probe (W27). No Go test, no
  daemon and no browser was run in this fire; every runtime premise is a planner obligation (§6).

---

## 0. Scope, the relationship to the parent, and the split

### 0.1 What the parent's V2 row says, measured against the tree

The parent's §0.3 V2 row: "`why <target>` in the bar: the `why` walk (`why.go`, 1,050 lines, V12)
behind a reader interface, served in-process; state change none; 1 sprint ≈ 1,600". Five premises,
each measured:

- **P1 — "the walk (`why.go`, 1,050 lines)": the count holds, the noun does not.** `why.go` is
  1,050 lines (W2) but three things: (a) **the CLI's shared read surface** — `worldReader`, the
  wire types, `getJSON`, `scanEntries`, `record` (lines 44–253) — used by **8 other files** of
  `cmd/ailang-worldd` with **7 non-test construction sites** (`demo_proof`, `demo_seed` ×3,
  `episode_report`, `logtail`, `provenance`) and 21 uses in `demo_test.go` alone (W3, W4); (b) the
  CLI front — help text, `parseWhyTarget`, `runWhy` (255–423); (c) the walk proper — `locate`,
  `byWorld`, `scanFor`, the chain, `buildChain`, `printChain` (424–1050). Moving "why.go" into a
  host package drags every CLI verb that reads World; V2 moves the file but keeps the CLI's
  `worldReader` as a thin adapter so none of those callers or their tests are edited (§2.1).
- **P2 — "behind a reader interface": no interface exists on the CLI side; one already exists on
  the daemon side.** `worldReader` is a concrete struct over the CLI's `*client` (W3). The daemon,
  however, already reads every `/v1` GET route through one seam, `readStore` (`d.reads`), and its
  tests already wrap it to inject failures (`aguiFailReads`, W7). The walk reads exactly seven `/v1`
  GET patterns (W5) whose handlers are registered on the daemon's mux (W6). So the interface V2 adds
  is at the **transport**: a `Getter` (one GET path → status and body), which the CLI satisfies with
  its HTTP client and the daemon satisfies by dispatching the same seven handlers in-process (§2.4).
  The walk then sees byte-identical inputs on both sides by construction.
- **P3 — "the same content" cannot hold past the daemon's read deadline.** The CLI walks under a
  60 s budget (`whyBudget`, W10); the route reads under `d.readCtx`, 10 s (W10). A walk that needs
  more than 10 s cannot answer in the bar; V2 refuses it explicitly (§2.5) — never truncated, never
  "not found".
- **P4 — the CLI turns an exhausted deadline into "not found".** `runWhy` prints `not found in the
  last N entries` for `context.DeadlineExceeded` as well as for an exhausted scan (W11). The bar
  must not copy that: a deadline is a refusal with its own text (§2.5).
- **P5 — the CLI turns some read failures into broken links.** `checkedObject` returns a read
  error as a ✗ with the error text, `buildChain` does the same for the world read, and
  `effectRecord` swallows its error (W12): at a terminal a transient 500 mid-walk prints an
  integrity refusal (exit 3). The bar must not present a read failure as an integrity verdict
  (§2.5). The CLI's behaviour is unchanged by V2 (RV3).

### 0.2 What V2 changes from a page

Nothing in World. The new code path is one more read-only arm of `GET /surface/verb.json`; its
callees are host/verbs (pure), host/why (pure but for an injected `Getter`) and seven existing `/v1`
GET handlers dispatched in-process. No POST, store write, credentials store, process or `.ail`
(§3, VS2).

### 0.3 Amendments to the parent (explicit)

| # | parent text | amended to | why |
|---|---|---|---|
| **A1** | §2.1(b) table row `why` → class `read`, "**`terminal`** answer (V2)" | `why` → class `read`, **runs** (§2.3's admission, §2.5's outcomes) | this is V2 |
| **A2** | §2.2 `result` = `{text | json | stderr+exit:1}` | `result` also carries **`exit: 3`** with `text` or `json` (why's integrity refusal); every `result` answer's `text`/`json`/`stderr` is **capped at 1 MiB** (`MaxVerbAnswerBytes`), episode report included; over the cap → `refused` with the size and the CLI line | why has three exits (W9); no route output may be unbounded |
| **A3** | §2.3 `refused → say(reason)` | a `refused` answer **that carries `cli`** opens the verb panel (`$ <cli>`, the reason, *Copy the line*); a `refused` without `cli` (Split refusals, the `world` 4b shape) stays a `say` | RS-SAY (§9): V2 adds refusals whose whole point is the CLI line (deadline, read failure, `--result`), and a 2.6 s fade loses it |
| **A4** | §2.3 "what a second submit does while one is pending is left to the planner" | **every Enter owns the bar**: an empty submit supersedes a pending ask, and a superseded turn never navigates, including after `loadIndex` | V1 judge note 1 (§9) |
| **A5** | §2.4 trace row `why` → `terminal` (shadows brane `why`) | `why` (bare) → `refused`, "ailang-worldd why: a target is required", `cli` = `ailang-worldd why` — still a verb, still shadowing brane `why` (RK1 unchanged; `world why` navigates) | the CLI's own answer (W8) |

The parent's §0.3 phase table is amended in one cell (V2's status → this doc). Nothing in row 201's
siblings, the `/v1` table, the CSP or `approvalRoute` changes.

### 0.4 The split

- **V2a — the seam (no user-visible change).** `host/why` (the moved walk behind `Getter`, `Run`,
  `ParseTarget`); the CLI's `why` repointed to it; the CLI's `worldReader` an adapter; why's flag
  set in `host/verbs` used by the CLI. Gate: **every existing CLI test passes unedited** except the
  two unit-test files that move with the code (§7).
- **V2b — the bar.** `verbs.AdmitWhy`; the daemon's in-process getter; the route's `why` branch and
  the A2 cap; the page (A2–A4); the V1 judge-note tests; live Chrome XV2; QUICKSTART §19.
- Nothing flows backwards: V2a is behaviour-neutral for the CLI and leaves the bar's `why` at
  `terminal`; V2b only adds a caller. Q1 (§11) asks whether to plan them as one sprint.

**Out of V2:** any mutating verb (V3/V4); `why --result`/`why -`/`--expect-store` from the bar
(§2.3, refused with the CLI line); changing the CLI's own deadline/read-failure wording (P4/P5,
RV3); streaming or paging a walk; any `.ail`; the static page's bar.

## 1. What the code does today (measured at `200ea5d5`)

1. **`runWhy`** (`why.go:349`) takes a leading target first (`args[0]` not starting with `-`, or
   exactly `-`, `:351`), then a `FlagSet` with four flags — `--expect-store`, `--result`, `--scan`
   (default 500), `--json` (`:357-360`) — parsed by `parseVerbFlags`, which is a plain `fs.Parse`
   (not interleaved; `mcpclient.go:415`); a second positional is `unexpected argument %q` (`:366`);
   `--scan` outside 1..5000 is refused (`:372`) (W8). Exit codes: 0 every link ✓, 1 usage or not
   found (`exitUsage = 1`), 3 any ✗ (`exitIntegrity`) (W9).
2. **Targets** (`parseWhyTarget`, `:304`): `--result <file>` reads a file (`os.ReadFile`, `:309`);
   `-` reads stdin up to 8 MiB (`:319`); `head`; `a2a:`/`mcp:`/`rest:<id>`; `sha256:<ref>` (parsed
   by `hashref.Parse`); otherwise a non-negative integer index (W13).
3. **`--expect-store`** compares `/v1/health`'s `db_path` with a path through `sameStorePath`
   (`filepath.Abs`, which reads the process's working directory; `doctor.go:369`) (W14).
4. **The reads.** Seven `/v1` GET patterns — `/v1/health`, `/v1/head` (via `headState`,
   `mcpclient.go:260`), `/v1/worlds/{ref}`, `/v1/objects/{ref}?payload=true`, `/v1/log/{index}`,
   `/v1/log?from&limit`, `/v1/receipts/{id}` — seven handlers
   (W5), all registered bare on the daemon's mux (`daemon.go:1110-1119`, W6) and unauthenticated
   (`isProtected` is only `POST /v1/commit`, W31). The CLI's client caps a response at 8 MiB
   (`maxClientResponseBytes`, `cli.go:21`, `:81`, W15).
5. **The daemon's read seam** `readStore` (`daemon.go:557`; field `reads`, `:431`) is what those
   handlers read through; tests replace it with failing wrappers (`agui_test.go:645`, W7).
   `handleReceipt` answers 400 for a non-`rest:` id; a handler's own read timeout is **503**
   `Timeout` "read deadline (10s) exceeded" (`handlers.go:349-350`, W29).
6. **Deadlines.** CLI: `budgetContext(whyBudget)`, 60 s (`why.go:36`, `:380`). Daemon: `readCtx` =
   `context.WithTimeout(r.Context(), d.readDeadline)`, `readDeadline = 10 * time.Second`
   (`handlers.go:287-289`, `daemon.go:144`) (W10). On `DeadlineExceeded` during `locate`, the CLI
   prints `not found in the last N entries` (`why.go:399-400`, W11). Read errors inside the chain
   become ✗ links (`:781-782`, `:834`) or are swallowed (`:525`) (W12).
7. **Shared helpers.** `elide`/`elideAt` live in `mcpclient.go` (`:595`, `:651`) and are called 4
   times in `why.go` and 5 times elsewhere (W16). `headState` has 9 non-test call lines (W35).
   Identifiers declared in `why.go` and used by other cmd files: `wireHeader`, `wireEntry`,
   `wireObject`, `wireWorld`, `invRecord`, `healthInfo`, `whyChain`, `whyLink`, `entryWire`,
   `worldWire`, `sumJSON`, `logPageMax`, `whyScanDefault`, `whyScanMax`, `whyBudget`, `contains`
   (10 files) (W4). Test files reference the unexported `whyInputLine`, `whyKeyField*`,
   `protectedLine`, `resultProtected` (`why_input_test.go`, `why_protected_test.go`) (W30).
8. **V1's route** (`surface_verb.go`): an answer struct with `exit,omitempty`; one deferred
   `writeApprovalJSON(200, a)` (`:45`); `case v.Words == "episode report"` runs it (`:56`); `why`
   (class `read`) falls to `a.Outcome = "terminal"` (`:63`); `lookup.go:99` gives `why` the reason
   "why reads World, but its walk runs in the CLI. The app cannot run it yet." (W17). The page's
   golden `verb-route.golden.json` records `why` → `terminal` (W18). VS1
   (`TestSurfaceVerbIsReadOnly`) pins host/verbs' imports to a pure set and `surface_verb.go`'s
   imports and `d.*` calls to exact sets (W32).
9. **V1's page** (`app.js`): `VERB_NEEDS.result = ['cli']` (`:485`) — **a `result` without `cli` is
   refused in code**, but no test row exercises it (0 rows; control: the terminal-without-cli row,
   1) (W19). A `result` with `stderr` must carry `exit === 1`, and the panel body is `stderr +
   'exit 1'`; `exit: 3` is unreadable today (`:499-500`, W20). `go(raw)` returns `out()` for an
   empty line **before** `askVerb` (`:1267-1268`) and `goAddress` awaits `loadIndex()` before
   acting (`:1283`) (W21). `say` fades after 2,600 ms (`:995`); Escape removes the panel (`:1520`);
   Enter calls `go` then focuses `#root` (`:1508`) (W24).
10. **V1's tests.** The never rows of `TestLookup` match reason prefixes without the CLI line
    (`lookup_test.go:35-37`; control: 3 rows that do contain `at a terminal: ailang-worldd`) (W22).
    `TestSplit` uses `MaxLine`/`MaxWords` symbolically; the literal `32` appears 0 times in
    `split_test.go` (control: 2 in `split.go`) (W23).
11. **Neither `host/why` exists nor does the daemon spawn a process** (W25, W26, with controls).

## 2. Design

### 2.1 `host/why`: the walk behind a `Getter` (V2a)

**The move.** `git mv cmd/ailang-worldd/why.go host/why/why.go`, then edit in place (package
`why`; exported names below; the CLI front deleted). A **new** file
`cmd/ailang-worldd/why_cli.go` holds what stays in the CLI. Keeping one destination file preserves
git's rename pairing (W27: a move with type renames and a 150-line deletion numstats as
`56 209 {cmd/ailang-worldd => host/why}/why.go`); the planner splits nothing further.

**The seam.**

```go
// Getter reads one GET path of the daemon's read surface: its status and body. A 404 is an
// answer, not an error; err is a transport or bound failure.
type Getter func(ctx context.Context, path string) (status int, body []byte, err error)

type Reader struct {
	Ctx         context.Context
	Get         Getter
	HealthCache *Health
}
```

Every method of today's `worldReader` (`getJSON`, `entry`, `page`, `object`, `world`, `headIndex`,
`health`, `record`, `scanEntries`, `locate`, `byWorld`, `scanFor`, `effectRecord`,
`checkedObject`, `buildChain`) becomes an exported method of `why.Reader` with an unchanged body
except `r.c.do(r.ctx, GET, p, nil)` → `r.Get(r.Ctx, p)` and `headState(r.ctx, r.c)` → a local
`r.head()` with `headState`'s exact body (12 lines; `headState` itself stays in cmd for its 8 other
callers). `http.StatusOK`/`StatusNotFound` become local constants so `host/why` imports no
`net/http` (VS2).

**Run** — the CLI front's result half, moved out of `runWhy`:

```go
type Options struct { Scan int; JSON bool }
type Answer  struct { Stdout, Stderr string; Exit int }   // exactly what the CLI writes
func ParseTarget(arg string) (Target, error)              // head, ids, sha256, index; pure
func Run(r Reader, t Target, o Options) Answer            // locate → buildChain → print
```

`Run` contains `runWhy`'s lines from `r.locate` to the exit choice (`why.go:394-422`) unchanged in
text: the `not found in the last %d entries (%s)` line, `ailang-worldd why: %v`, `MarshalIndent` +
newline for `--json`, `printChain` otherwise, exit 3 when `!ch.OK`. `ParseTarget` is
`parseWhyTarget` minus its `--result` and `-` arms (the two effects), same error texts.

**What stays in cmd** (`why_cli.go`): `whyHelp`; `runWhy` (flag set from `verbs.WhyFlags`, §2.2;
the leading-target and trailing-argument rules; the `--scan` bound; `parseWhyTarget` as a wrapper
that performs the file and stdin reads and otherwise calls `why.ParseTarget`; `--expect-store` with
`serviceMessage` as today; then `ans := why.Run(r.why(), t, opts)` and two `Fprint`s); and the
**compatibility layer** so no other CLI file or test is edited:

- type aliases `wireHeader = why.Header`, `wireEntry = why.Entry`, `wireObject = why.Object`,
  `wireWorld = why.World`, `invRecord = why.Record`, `healthInfo = why.Health`, `whyChain =
  why.Chain`, `whyLink = why.Link`, `entryWire = why.EntryWire`, `worldWire = why.WorldWire`
  (all exported-field structs; positional literals in `demo_test.go` keep compiling because field
  order is unchanged);
- constants `logPageMax`, `whyScanDefault`, `whyScanMax` = `why.…`; `whyBudget` stays a cmd
  constant (60 s, the CLI's own budget); `func sumJSON(v any) string { return why.SumJSON(v) }`;
  `contains` re-declared in cmd (it is a generic helper used by 10 files, W4);
- `worldReader` keeps its three fields and gains `func (r worldReader) why() why.Reader`; each of
  its 15 methods becomes a one-line delegation (same names, same signatures);
- `elide`/`elideAt` **move** to host/why (`why.Elide`, `why.ElideAt`) and `mcpclient.go`'s become a
  one-line delegation and a constant alias.

`why_input_test.go` and `why_protected_test.go` move with the code (`git mv`, package line only:
the names they use stay unexported inside host/why). `why_test.go` (8 tests, the CLI end to end)
stays in cmd **unedited** and is V2a's main pin (W30).

### 2.2 One parse for `why`'s arguments (V2a, used by V2b)

`host/verbs` gains the CLI's own parse, as V1 did for `episode report`:

- `WhyFlags(fs *flag.FlagSet) *WhyArgs` registers exactly `expect-store`, `result`, `scan`
  (default `WhyScanDefault = 500`), `json` with today's usage strings;
- `WhyLeadingTarget(args) (target string, rest []string)` — `runWhy`'s `:351` rule;
- `WhyTrailingTarget(fs, target) (string, error)` — the `:363-369` rule, error text `unexpected
  argument %q`;
- `WhyScanMax = 5000`, and `WhyScanOK(n) bool`.

The CLI's `runWhy` calls these between its unchanged `parseVerbFlags` and `parseWhyTarget`. host/verbs
stays pure (no new import outside VS1's set). A test pins `verbs.WhyScanDefault == why.ScanDefault`
and `verbs.WhyScanMax == why.ScanMax` (VA1).

### 2.3 The bar's admission: every argument form (V2b)

`verbs.AdmitWhy(args) (target string, scan int, asJSON bool, reason string)` runs §2.2's parse on a
fresh `FlagSet` (output discarded), in the CLI's order, then the bar profile. Then the route calls
`why.ParseTarget(target)`. Every form:

| form (after `why`) | CLI does | bar answers | exact reason |
|---|---|---|---|
| *(nothing)* | usage, exit 1 | `refused` + `cli` | `ailang-worldd why: a target is required` |
| `<index>`, `head` | walks | **runs** | – |
| `sha256:<64 hex>` | walks (dispatch on semanticId, or world ref) | **runs** | – |
| `sha256:<malformed>` | usage | `refused` + `cli` | `ailang-worldd why: "<arg>" is not a content address: <hashref error>` |
| `a2a:<id>`, `mcp:<id>` | scans back | **runs** | – |
| `rest:<id>` | reads the receipt | **runs** (a 400 from `handleReceipt` is the CLI's own `result` exit 1, §2.5) | – |
| anything else (`x`, `-3`) | usage | `refused` + `cli` | `ailang-worldd why: "<arg>" is not an entry index, head, sha256:<ref>, a2a:/mcp:/rest:<id>, or -` |
| `--json` | JSON chain | **runs** (`json`) | – |
| `--scan N`, 1 ≤ N ≤ 5000 | bounded scan | **runs** | – |
| `--scan N` out of range | usage | `refused` + `cli` | `ailang-worldd why: --scan must be 1..5000` |
| a second positional | usage | `refused` + `cli` | `ailang-worldd why: unexpected argument "<arg>"` |
| unknown flag, bad value | usage | `refused` + `cli` | `ailang-worldd why: <flag package error>` |
| `-h`, `--help` | prints help | `refused` + `cli` | `see ailang-worldd why --help at a terminal` |
| `-` | reads stdin | `refused` + `cli` | `why - is not run from the app: the app has no stdin; at a terminal: <cli>` |
| `--result <file>` | reads a file | `refused` + `cli` | `--result is not run from the app: the app reads no file on the daemon's host; at a terminal: <cli>` |
| `--expect-store <path>` | compares stores | `refused` + `cli` | `--expect-store is not run from the app: the app reads its own daemon's store; at a terminal: <cli>` |

Order: leading target → parse (help, error) → trailing argument → scan bound → bar profile
(`--result`, `--expect-store`, `-`, in that order) → `ParseTarget`. Every input is already bounded by
V1's `Split` (1,024 bytes, 32 words, W32); no further length rule is needed (a `sha256:` target is
71 bytes, an id is opaque and bounded by the line).

### 2.4 The daemon's in-process `Getter` (V2b)

`host/daemon/surface_verb_why.go`:

- `whyReadRoutes` — a fixed list of the seven patterns of W6 that the walk uses (`GET /v1/health`,
  `/v1/head`, `/v1/worlds/{ref}`, `/v1/objects/{ref}`, `/v1/log/{index}`, `/v1/log`,
  `/v1/receipts/{id}`);
- `d.whyMux()` — built once (`sync.Once`), a private `http.ServeMux` registering exactly those
  patterns to **the same method values** `Handler()` registers (`d.handleHealth`, …);
- `d.whyGetter(ctx, tally)` returns a `why.Getter` that, per call:
  1. builds `GET <path>` with `http.NewRequestWithContext(ctx, …)`; a build error → `err`;
  2. `h, pattern := mux.Handler(req)`; **`pattern == ""` → `err` ("unrouted in-process read"),
     never a 404** — an unregistered path must not read as "absent";
  3. serves into a `boundedRecorder` capped at `why.MaxResponseBytes` (= 8 MiB, the CLI's cap,
     pinned equal by VW4); overflow → `err` with the CLI's wording `response from <path> exceeds N
     bytes`;
  4. tallies: status 503 → `tally.timeout`; any other status ≥ 500 → `tally.failed`; any `err` →
     `tally.failed` (and `d.writeInternalErrorLog`); returns `(status, body, nil)` otherwise.

The seven handlers run their own `d.readCtx(req)` derived from the walk's context, so every read is
bounded by the walk's deadline as well as its own. `trackActivityRequests` and the session
middleware are not traversed: the outer `/surface/verb.json` request is the activity, and none of
the seven routes is protected (W31).

### 2.5 The route's `why` branch and its outcome precedence (V2b)

`handleSurfaceVerb` gains `case v.Words == "why": d.verbWhy(r, args, &a)` before the
`Read|Mutate|Host → terminal` case. `verbWhy`:

1. `AdmitWhy` → a reason → `refused` (with `cli` from `Lookup`, already set). `ParseTarget` error
   → `refused`, reason `ailang-worldd why: <err>` (the CLI's first stderr line, without the help).
2. `ctx, cancel := d.readCtx(r)`; `ans := why.Run(why.Reader{Ctx: ctx, Get: d.whyGetter(ctx,
   &tally)}, t, why.Options{Scan, JSON})`.
3. Precedence, first match wins:

| # | condition | outcome | reason / members |
|---|---|---|---|
| a | `r.Context()` cancelled (the page aborted: superseded or its 15 s budget) | `refused` (written to a closed connection; nobody reads it) | – |
| b | `ctx.Err() == DeadlineExceeded` **or** `tally.timeout` | `refused` + `cli` | `read deadline (10s) exceeded: the walk did not finish in the app; at a terminal it has 60 s: <cli>` |
| c | `tally.failed` | `refused` + `cli` | `the daemon could not read World during the walk, so no chain is shown; at a terminal: <cli>` (no raw error; the handlers and the getter logged it) |
| d | `ans` has both or neither of `Stdout`/`Stderr`, or `Exit ∉ {0,1,3}`, or `Stderr` with `Exit ≠ 1` | `refused` + `cli` | `the walk answered something the app cannot show; at a terminal: <cli>` |
| e | `len(Stdout)+len(Stderr) > MaxVerbAnswerBytes` (1 MiB) | `refused` + `cli` | `the answer is <n> bytes, over the app's 1 MiB cap; at a terminal: <cli>` |
| f | otherwise | `result` + `cli` | `Exit 0/3` → `text` (or `json` with `--json`) = `Stdout`, `exit` = 3 when 3; `Exit 1` → `stderr` = `Stderr`, `exit: 1` |

(b) precedes (c) because a deadline surfaces as 503s and as context errors inside the walk; both
precede (f) because P4/P5 would otherwise print a deadline as "not found" or a read failure as ✗.
The A2 cap (e) is applied **once, in the deferred writer**, to every `result` answer (episode
report's too), so no `result` of any verb leaves the route unbounded.

### 2.6 The page (V2b)

1. **`verbOutcome` (A2).** `result` keeps `cli` required; exactly one of `text`/`json`/`stderr`;
   `exit` ∈ {absent, 1, 3}; `stderr` ⇔ `exit === 1`; `exit === 3` requires `text` or `json`. Body =
   the shown member + (`exit` ? `'exit ' + exit` : ''). (Every CLI output ends in a newline, so the
   status sits on its own line, as V1's `stderr + 'exit 1'` already does.)
2. **Refusals with a CLI line (A3).** `refused` with a non-empty `cli` → `{panel: {head: '$ ' + cli,
   body: reason, copy: cli}}`; without `cli` → `{say: reason}` as today.
3. **Every Enter owns the bar (A4).** `go(raw)` becomes `goWith(raw, fx)` (exported to the node
   harness; `go` passes the real effects):
   ```js
   let barTurn = 0;
   async function goWith(raw, fx) {
     const turn = ++barTurn, current = () => turn === barTurn;
     if (blank(raw)) { dropVerb(); return fx.out(); }          // supersedes a pending ask
     const act = await fx.ask(raw);
     if (!current()) return;                                     // a newer Enter owns the bar
     return verbAct(act, raw, { ...fx.act, address: r => fx.address(r, current) });
   }
   ```
   `dropVerb()` marks the pending `verbJob` superseded and aborts it (the body of `askVerb`'s
   first line, extracted). `goAddress(raw, current)` re-checks `current()` after `await
   loadIndex()` and does nothing if a newer Enter came. `askVerb`'s own supersession and its
   `finally { if (verbJob === job) verbJob = null; }` are unchanged.
4. **Discoverability.** The placeholder gains `· why <target>`; the `?` keys row reads `in the bar:
   a place, episode report E, why T, approve ID`.

### 2.7 Input and failure classes, traced to their outcome

Entry: `GET /surface/verb.json` → `approvalRoute` (R2 Host check) → V1 steps 1–3 → `Lookup` (`why`
row, class `read`) → §2.5. Page: `goWith` → `askVerb` (15 s) → `verbOutcome` → `verbAct`.

| input / condition | route | page shows | navigates? |
|---|---|---|---|
| `why 3` (entry exists, chain ✓) | `result`, `text`, exit 0 | panel `$ ailang-worldd why 3` + the CLI's stdout | no |
| `why head --json` | `result`, `json` | panel + the CLI's JSON | no |
| `why <index>` with a ✗ link | `result`, `text`, `exit:3` | panel + stdout + `exit 3` | no |
| `why 999` (no such entry) | `result`, `stderr` `ailang-worldd why: log entry 999 not found`, exit 1 | panel + stderr + `exit 1` | no |
| `why a2a:nobody:x` (scan exhausted) | `result`, `stderr` `… not found in the last 500 entries (invocationId)`, exit 1 | panel | no |
| `why head` with no selected head | `result`, `stderr` `ailang-worldd why: no world head has been selected yet`, exit 1 | panel | no |
| `why sha256:<unknown 64 hex>` | object 404 → world 404 → `result`, stderr `… is neither an object nor a world in this store`, exit 1 | panel | no |
| `why rest:x` with a pending receipt | `result`, stderr `receipt rest:x is pending: nothing committed`, exit 1 | panel | no |
| `why`, `why x`, `why sha256:zz`, `why --scan 0 head`, `why 1 2`, `why --bogus` | `refused` + `cli` (§2.3) | refusal panel (A3) | no |
| `why -`, `why --result f`, `why --expect-store p head`, `why --help` | `refused` + `cli` (§2.3) | refusal panel | no |
| a read 500s mid-walk (`d.reads` fails) | (c) `refused` + `cli` | refusal panel; never ✗, never "not found" | no |
| a read times out (503) or the walk passes 10 s | (b) `refused` + `cli` | refusal panel with the 60 s note | no |
| an unrouted path (a host/why change reads a new route) | getter `err` → (c) | refusal panel | no |
| a sub-response over 8 MiB | getter `err` → (c) | refusal panel | no |
| the answer over 1 MiB (huge `--json` worktree delta) | (e) `refused` + size | refusal panel | no |
| the daemon accepts but never answers / the body stalls | – | at 15 s: `say('the verb route did not answer in time')` (V1) | no |
| HTTP 404/400/421/500, network error, unreadable body | – | V1's explicit `say`, view unchanged | no |
| a newer Enter while `why` is pending | (a) | the older answer is dropped; the newer one answers | no |
| an **empty** Enter while any ask is pending | (a) | home (`out()`); the older answer is dropped | no — **fixed** (A4) |
| an Enter during `goAddress`'s `await loadIndex()` | – | the older `openState` is skipped | no — **fixed** (A4) |
| `world why` | `address` (V1 step 3) | brane `why` | yes (only path that navigates) |

No `why` answer, success or failure, reaches `addressTarget`: the only navigation is a 200,
schema-valid `outcome:"address"` (V1, unchanged), and `Lookup` gives `why` no address class.

## 3. Pure core, effects at the boundary (coding standards S1–S3)

- **S1.** No `world/` or `.ail` change. `host/verbs` stays pure (VS1's import set). `host/why` does
  no I/O itself: every read goes through the injected `Getter`; its direct imports exclude `os`,
  `os/exec`, `net`, `net/http`, `io/fs`, `database/sql` and `host/store` (VS2). The file read and
  stdin read of `--result`/`-` stay in cmd.
- **S2.** The only effect V2 adds is seven existing read handlers dispatched in-process under
  `d.readCtx`. The CLI's effects (HTTP, file, stdin, `filepath.Abs`) are unchanged and stay in cmd.
- **S3.** No kernel change, no registry package, no new route, `/v1` table unchanged.

## 4. Threats (V2)

| threat | answer |
|---|---|
| T1. A page changes World through `why` | `verbWhy` reaches no write: VS2 pins `surface_verb_why.go`'s `d.*` calls to `{readCtx, readDeadline, whyMux, writeInternalErrorLog}` and its imports; the getter's mux registers only GET read handlers (VR6 asserts the list) |
| T2. File read on the daemon host (`--result`) / stdin | refused before any read (§2.3, VA2) |
| T3. SSRF / another daemon | there is no `--addr` on `why`'s flag set; global flags are V1 step 2; the getter never opens a socket |
| T4. A path outside the seven routes (e.g. a future host/why change that reads `/v1/session/self`) | unrouted → error, never 404, never another handler (VR5) |
| T5. Hostile bytes in a chain (inputs are agent-written) | `textContent` only (V1's VP3), a hostile-input chain row in the page test (VP9) |
| T6. Denial of service | the line is ≤ 1,024 bytes; `--scan` ≤ 5,000; the walk is bounded by 10 s and each read by 10 s; each sub-response ≤ 8 MiB; the answer ≤ 1 MiB. A 500-entry log page measured 15.65 ms p50 (W34), so a 5,000-entry scan is ~10 pages plus 5,000 object reads, expected well under the deadline; past it, (b) |
| T7. A read failure shown as an integrity verdict | (c) precedes (f) (VR3) |
| T8. Confusion: brane `why` | unchanged from V1 (RK1): `world why` navigates (VP4 still passes) |

## 5. By-construction argument

1. **Same walk.** The CLI and the route both call `why.Run` (one function); the CLI's `runWhy`
   writes `Answer.Stdout`/`Stderr` verbatim and returns `Answer.Exit`.
2. **Same inputs.** The walk's bytes come from the same seven handler methods, read through the same
   `d.reads`; the CLI gets them over HTTP, the bar in-process (VR4 compares the two transports'
   bytes per route).
3. **Same parse.** `WhyFlags`, `WhyLeadingTarget`, `WhyTrailingTarget`, `WhyScanOK` and
   `ParseTarget` are called by both; the bar differs only by the named profile of §2.3.
4. **Never silent.** Every non-`f` row of §2.5 is a `refused` with the CLI line; no `why` answer is
   `address`; the page navigates only on a schema-valid `address` (V1) and only for the current turn
   (A4).

## 6. Test plan

Each row names the mutant it kills. "Real CLI" = `run([]string{…})` / `whyRun` in-process against a
loopback daemon (`newWalkRig`, `why_test.go:52`); expected bytes are **generated by that run**,
never typed.

### V2a — the seam

- **Must-live, unedited:** all of `why_test.go` (8 tests), `demo_test.go`, `demo_proof_test.go`,
  `outcome_test.go`, `mcpclient_test.go`, `help_test.go`, `service_pointers_test.go`,
  `cliwalk_e2e_test.go`, every `episode_report*_test.go`, `logtail`/`provenance` tests.
- **VW1** moved unit tests (`host/why/why_input_test.go`, `why_protected_test.go`) pass with only
  the package line changed. Kills: a moved helper edited in passing.
- **VW2** `Run` exit and stream contract over a **replay `Getter`** whose bodies are recorded from
  the walk rig's real routes (recorded in the test, not typed): exit 0 → stdout only; ✗ (one body
  tampered) → stdout, exit 3; absent index → stderr `ailang-worldd why: log entry 9 not found\n`,
  exit 1; scan exhausted → `not found in the last 1 entries (invocationId)`. Kills: `Run`
  returning exit 1 for ✗; a stream swapped.
- **VW3** `Getter` classes: a getter returning 404 for `/v1/objects/…` → "does not resolve"; one
  returning `err` → the CLI's `ailang-worldd why: <err>` line. Kills: a 404 treated as an error;
  an error treated as absent.
- **VW4** constants: `why.MaxResponseBytes == maxClientResponseBytes`, `why.ScanDefault == 500`,
  `why.ScanMax == 5000`, `why.ElideAt == 160` (literals, so a constant changed on both sides fails).
  Kills: a drifted cap.
- **VA1** `WhyFlags` registers exactly `expect-store, result, scan, json` (`VisitAll`, defaults,
  usage strings equal to the base's); `verbs.WhyScan* == why.Scan*`. Kills: a flag renamed on one
  side.
- **VS2a** `host/why` direct imports exclude the §3 set; `host/verbs` still passes VS1's pure set.
  Kills: an `os`/`net/http` import into the walk.
- **Size check (planner):** after M1, `git diff --numstat` must show the move as
  `{cmd/ailang-worldd => host/why}/why.go`; if rename pairing is lost, stop and re-cut (§8).

### V2b — the bar

**Admission (`host/verbs`)**
- **VA2** one row per line of §2.3's table, asserting `reason` exactly. Kills: any refusal deleted;
  `--result` admitted; the profile order changed so `--help` wins over `--result` differently from
  the table.
- **VA3** parity with the CLI's parse, **generated**: for each usage row of §2.3, run the real CLI
  (`whyRun`) and assert the bar's reason equals the CLI's **first** stderr line. Kills: the bar's
  parse diverging from the CLI's (e.g. interleaving flags that the CLI's `fs.Parse` would not).

**Route (`host/daemon`)**
- **VR1** `why` outcomes end to end through `approvalRoute` on a seeded walkable store: §2.7's
  rows 1–9. Kills: the `why` case falling to `terminal`; `json` and `text` swapped.
- **VR2** deadline (b): `d.readDeadline = expiredReadDeadline` (the existing helper,
  `episode_report_test.go:331`) and a `d.reads` wrapper that blocks past the deadline → the (b)
  reason, never `result`, never "not found". Kills: (b) deleted (the CLI's "not found" leaks).
- **VR3** read failure (c): `d.reads` wrapped to fail the Nth `GetObject` (the `aguiFailReads`
  pattern) for N in {first, mid-chain, the effect record's} → the (c) reason; the response body
  contains no store error text. Kills: (c) deleted (a ✗ chain shown); (c) after (f).
- **VR4** transport parity: for every path the walk rig's walk requests (recorded by a proxy
  around `whyRun`), the in-process getter's `(status, body)` equals `Handler()`'s over HTTP, byte
  for byte (health compared on `interpreter_ref`/`interpreter_version` only: its `active` block is
  per-request). Kills: a private mux bound to a different handler.
- **VR5** unrouted and oversize: a getter call for `/v1/session/self` and `/v1/nowhere` → `err`
  (not 404); a stub handler writing 8 MiB + 1 → `err`. Kills: unrouted read as absent; the cap
  removed.
- **VR6** `whyReadRoutes` equals the set of patterns the recorded walk reaches (VR4's proxy log,
  generated) and is a subset of `Handler()`'s GET patterns. Kills: a route added to the private
  mux beyond the walk's needs.
- **VR7** the A2 cap: a `result` over 1 MiB (a replay `Getter` returning a large worktree delta
  under `--json`) → (e); exactly 1 MiB → `result`. Applies to `episode report` too (one row).
  Kills: the cap removed or applied only to `why`.
- **VS2** `surface_verb_why.go`'s imports and `d.*` calls are exact sets (§4 T1); VS1's
  `surface_verb.go` call set gains only `verbWhy`. Kills: a write or spawn added.

**The bar's output is the CLI's (`cmd/ailang-worldd`, generated)**
- **VE1** on `newWalkRig`: for `3`, `head`, `head --json`, every `sha256:` kind of
  `TestWhyTargetFormsResolveToOneEntry`, an `a2a:` and an `mcp:` id, `999`, and an exhausted
  `--scan 1`, run the real CLI and the route against the same daemon; assert `text`/`json` ==
  stdout, `stderr` == stderr, `exit` == the exit code. Kills: the renderer forked; the JSON
  re-marshalled.
- **VE2** exit 3 parity: a store fixture whose committed record disagrees with its entry header
  (planner obligation: a `makeCLICommit`-built entry the commit route accepts; if the store refuses
  every such fixture, run VE2 in `host/daemon` with a tampering `d.reads` wrapper and the expected
  bytes from `why.Run` over an HTTP `Getter` against the **same** wrapped daemon's `Handler()`).
  Kills: exit 3 dropped by the route.
- **VE3** the page golden `host/surfacepage/testdata/verb-why.golden.json` is written by VE1's run
  under `WORLD_UPDATE_VERB_GOLDEN=1` and compared byte for byte otherwise; the planner first proves
  the rig deterministic (two fresh rigs, equal answers) or records which members vary and keeps
  them out of the golden.

**Page (`host/surfacepage`, node, non-skippable in CI)**
- **VP7** (A4, judge note 1) with `goWith` and fakes: Enter `row-189@6` (ask pending) → Enter `''`
  → `out` once; release the first answer as a schema-valid `address` → `fx.address` **not** called.
  Then: Enter `row-189@6` → `address` answer → during the deferred `loadIndex`, Enter `log` →
  `openState` not called. Kills: `dropVerb` removed; the `current()` check after `ask` or after
  `loadIndex` removed.
- **VP8** (X-JS-FINALLY-NULL) three asks A, B, C: A pending; B supersedes A; A's `finally` runs; C
  must abort B (B's signal aborted, B's answer `{superseded:true}`). Kills: the `verbJob === job`
  guard removed.
- **VP9** `why` answers from the generated golden: `text` → panel body = stdout; `json` → body =
  JSON; stderr exit 1 → body + `exit 1`; a golden body with `exit:3` → body + `exit 3`; a hostile
  input line in a chain stays text. Kills: exit 3 unreadable; exit appended twice.
- **VP10** (judge note 3) `result` shape refusals: no `cli`; `exit:3` with `stderr`; `exit:2`;
  `text` and `json` both present → `UNREADABLE`. Kills: `result: []` in `VERB_NEEDS`; the exit rule
  loosened.
- **VP11** (A3, RS-SAY) a `refused` with `cli` → panel with *Copy the line*; without `cli`
  (`TooLong`, `world a b c`) → `say`. V1 rows that asserted `say` for `log`, `ailang-worldd world get
  world-1`, `episode` move to the panel expectation. Kills: the panel branch deleted.
- **V1 page rows re-pointed:** rows that used `why` as the *terminal* example (`app_verb_test.go:106`,
  `:138`, `:161`) use `episode launch w1 --preset p`; VP2's failure rows keep `why` as an input.

**V1 judge notes, grammar side (`host/verbs`)**
- **VB6** (judge note 2) for every word in `NeverWords` and `OtherWords`, bare and with one
  argument, `reason` **ends with** `"; at a terminal: " + cli`. Kills: the CLI line dropped from a
  never reason.
- **VB7** (`MaxWords`/`MaxLine`) `MaxWords == 32` and `MaxLine == 1024` as literals. Kills: a
  bound changed with its test.

**Live Chrome (on the A3 driver, its skip rule):** **XV2** type `why head` in `#addr`, Enter, read
the panel; equal to the CLI's stdout on the same daemon; then **Escape** removes `#verb` (judge
note 5, Escape half).

**Docs:** QUICKSTART §19 gains `why <target>` with one verbatim example and the 10 s note.

## 7. Inventories and conflict surface

- **Moved:** `cmd/ailang-worldd/why.go` → `host/why/why.go`; `why_input_test.go`,
  `why_protected_test.go` → `host/why/` (package line only).
- **New:** `cmd/ailang-worldd/why_cli.go`, `host/why/run_test.go`,
  `host/daemon/surface_verb_why.go` + `_test.go`, `cmd/ailang-worldd/verb_why_e2e_test.go`,
  `host/surfacepage/testdata/verb-why.golden.json`, the XV2 test.
- **Edited:** `mcpclient.go` (`elide` delegation), `host/verbs/flags.go`/`lookup.go` (+ tests),
  `surface_verb.go` (+ VS1's call set, the `why` golden row), `app.js`, `index.html`,
  `app_verb_test.go` (re-pointed rows, §6), `docs/QUICKSTART.md`.
- **Pre-existing tests edited (each with its reason):** `surface_verb_test.go` VS1 call set
  (+`verbWhy`); `surface_verb_outcomes_test.go` (`why` → refused row; golden regenerated);
  `lookup_test.go` (`why` row); `app_verb_test.go` (A3/A5 re-pointing). Nothing else.
- **Not touched:** row 201's files, `approvalRoute`, `/v1` handlers, `readStore`, CSP, any `.ail`,
  `tools/launchd/*`.

## 8. Milestones and sizes

| sprint | milestone | content | est. changed lines |
|---|---|---|---|
| V2a | M1 | `git mv` + host/why edits (exports, `Getter`, `Run`, `ParseTarget`, `Elide`) | ≈ 410 (rename-paired) |
| V2a | M2 | `why_cli.go` (front + compatibility layer), `mcpclient.go`, tests moved | ≈ 260 |
| V2a | M3 | `WhyFlags` & friends + CLI use; VW1–VW4, VA1, VS2a | ≈ 280 |
| V2b | M4 | `AdmitWhy`, the `why` row, VA2–VA3, VB6–VB7 | ≈ 170 |
| V2b | M5 | the getter, `verbWhy`, the A2 cap; VR1–VR7, VS2 | ≈ 400 |
| V2b | M6 | VE1–VE3, the page (A2–A4), VP7–VP11, XV2, QUICKSTART | ≈ 330 |

V2a ≈ 950, V2b ≈ 900, together ≈ 1,850 (× 1.09 ≈ 2,015). **The cut** if V2a over-runs: M3's tests
before M3's CLI repointing (the repointing is behaviour-neutral and can ride V2b). **If rename
pairing is lost** (git pairs a rename only above 50 % similarity), M1 counts ≈ 1,700 on its own; the
planner then keeps host/why's diff from the original minimal (no reflow, no comment rewrites).

## 9. The V1 judge notes, disposed

| note (W36) | measured at HEAD | disposition |
|---|---|---|
| 1. an empty submit does not supersede a pending ask (and `await loadIndex()` before `openState`) | true (W21) | **in V2b**: A4 + VP7 |
| 2. a never-verb's refusal can lose its CLI line untested | true (W22) | **in V2b**: VB6 |
| 3. a `result` without `cli` is accepted by the page | **false as code** (`VERB_NEEDS.result = ['cli']`), **true as test** (no row; W19) | **in V2b**: VP10 pins it |
| 4a. third submit's supersession untested (`X-JS-FINALLY-NULL`) | true (W19's file has two-ask cases only) | **in V2b**: VP8 |
| 4b. `MaxWords` not pinned apart from its constant | true (W23) | **in V2b**: VB7 (with `MaxLine`) |
| 4c. Escape closes the panel; Enter refocus | untested (W24) | Escape **in V2b** (XV2); Enter refocus **out**: V2 does not change focus handling and the refocus is a DOM-only effect with no node seam — residual RS-FOCUS |
| 4d. RS-SAY: refusal text in a 2.6 s `say` | true (W24) | **in V2b** for refusals with `cli` (A3, VP11); `say` stays for refusals without a CLI line |
| (also logged) 6. the moved printer's `start:`/`uncommitted intents:` lines; 7. the 256-byte episode boundary | – | **out**: report-specific, pre-existing, untouched by V2 |

## 10. Residuals and risks

- **RV1** the bar answers only walks that finish in 10 s; a deep `--scan` on a slow store refuses
  with the CLI line (§2.5 b). No paging or streaming.
- **RV2** the 1 MiB answer cap can refuse a legitimate `--json` chain with a large worktree delta;
  the refusal names the size and the CLI line.
- **RV3** the CLI keeps P4/P5 (deadline as "not found", read failure as ✗). V2 does not change the
  CLI's output; a follow-up row may adopt §2.5's (b)/(c) at the terminal.
- **RV4** the compatibility layer in cmd (aliases and 15 delegations) is the price of not editing 8
  CLI files and their tests; a later cleanup may repoint them.
- **RV5** rename pairing is a size-accounting assumption (W27), not a behaviour; §8 states the
  fallback.
- **RS-FOCUS** Enter's refocus stays untested (§9).

## 11. Open questions (each with the default the loop applies)

- **Q1. One sprint or two?** The measured scope is ≈ 1,850 (≈ 2,015 calibrated), over the parent's
  1,600 and at the ~2,000 cap. Default: **two sprints, V2a then V2b**, as §8. A planner whose
  compiled patches measure ≤ 1,900 for both may ask the controller to run them as one.
- **Q2. Should the bar's `--scan` be capped below 5,000 to fit 10 s?** Default: **no** — the same
  bound as the CLI; past the deadline the answer is an explicit refusal (RV1).
- **Q3. Does A3 (refusals with `cli` open the panel) need Mark?** Default: **no** — it is within
  D-WORLD-107 = A (the bar "shows the exact CLI line"), changes no World state, and only moves text
  from a fading header into the panel V1 already uses.
- **Q4. Ledger ask?** Default: **no** — read-only, within D-WORLD-107/109.

## 12. What this is NOT doing

- Not running `why --result`, `why -` or `--expect-store` from the bar; not reading a file on the
  daemon host.
- Not adding a route, a POST, a store write, a process, a key path, or any `.ail`.
- Not changing the CLI's output, exit codes or wording (RV3), nor `headState`, nor `/v1`.
- Not moving the CLI's other World readers (`provenance`, `log tail`, `demo`) into host/why.
- Not designing V3–V5.

## 13. Verification Log

All commands ran in `/Users/voightkampff/dev/sunholo-data/ailang-world/.wt-world-iter396-row243-v2`
at `200ea5d5` on 2026-11-16 (zsh; output trimmed), except W27 (scratch repo, stated).

| # | claim | command | observed |
|---|---|---|---|
| W1 | base sha; cmd/host unchanged since V1's merge (negative, with control) | `git rev-parse HEAD`; `git diff --stat fb4d8b89 HEAD -- cmd host \| tail -1`; control `… -- design_docs \| tail -1` | `200ea5d5b7c2944e561feeff16a4b2b3042333d1`; *(empty)*; control `19 files changed, 1121 insertions(+), 138 deletions(-)` |
| W2 | why.go is 1,050 lines in three sections | `wc -l cmd/ailang-worldd/why.go`; `grep -n '^// [a-z]\|^func ' …/why.go` | `1050`; `44: // the read surface`, `255: // why`, `349: func runWhy`, `425: func (r worldReader) locate`, `587: // the chain` |
| W3 | `worldReader` is shared by other CLI files | `grep -lw worldReader cmd/ailang-worldd/*.go \| grep -v '/why' \| wc -l`; `grep -rn 'worldReader{' cmd/ailang-worldd/*.go \| grep -v _test \| grep -v '/why' \| wc -l`; `grep -c worldReader …/demo_test.go` | `8`; `7` (`demo_proof.go:211`, `demo_seed.go:191`, `:588`, `:603`, `episode_report.go:135`, `logtail.go:79`, `provenance.go:58`); `21` |
| W4 | why.go identifiers used elsewhere | `for id in wireHeader wireEntry … contains; do grep -lw "$id" cmd/ailang-worldd/*.go \| grep -v '/why'; done` | `wireEntry`: demo_proof_test, demo_seed, demo_test, logtail, provenance; `invRecord`: demo_test, provenance; `healthInfo`: cliwalk_e2e_test, demo_test, marketing_chrome_test; `whyChain`: demo_seed, demo_test; `whyBudget`: episode_markdown_test, episode_report, provenance; `contains`: 10 files; `errScanExhausted`, `frozenRigPort`, `printChain`: none |
| W5 | the paths the walk reads | `grep -on '"/v1/[a-z]*' cmd/ailang-worldd/why.go \| sort -t: -k2 -u`; `grep -n '"/v1/head"' …/mcpclient.go` | `173:"/v1/health`, `122:"/v1/log`, `136:"/v1/objects`, `454:"/v1/receipts`, `142:"/v1/worlds`; `mcpclient.go:260: c.do(ctx, http.MethodGet, "/v1/head", nil)` |
| W6 | those routes' handlers on the mux | `grep -n '"GET /v1/health"\|"GET /v1/head"\|… "GET /v1/receipts/{id}"' host/daemon/daemon.go` | `1110 handleHealth`, `1111 handleHead`, `1112 handleWorld`, `1113 handleObject`, `1115 handleLogEntry`, `1116 handleLogRange`, `1119 handleReceipt` |
| W7 | the daemon's read seam and its failure wrapper | `sed -n 543,567p host/daemon/daemon.go`; `grep -n 'reads readStore' …`; `grep -n 'aguiFailReads{' host/daemon/agui_test.go` | `// readStore is the daemon's request-read surface … let a test wrap the real store`; `557: type readStore interface {`; `431: reads readStore`; `645: d.reads = aguiFailReads{d.reads, tc.err}`, `802: …` |
| W8 | runWhy's parse | `grep -n 'fs.String("expect-store"\|fs.String("result"\|fs.Int("scan"\|fs.Bool("json"\|args\[0\] == "-"\|parseVerbFlags(fs\|unexpected argument\|--scan must be' …/why.go`; `grep -n 'if err := fs.Parse(args)' …/mcpclient.go`; `grep -n 'a target is required' …/why.go` | `351: if len(args) > 0 && (!strings.HasPrefix(args[0], "-") \|\| args[0] == "-")`; `357`–`360` the four flags (`scan` default `whyScanDefault`); `361: parseVerbFlags(fs, args, whyHelp, …)`; `366: "ailang-worldd why: unexpected argument %q\n"`; `372: "--scan must be 1..%d"`; `mcpclient.go:415: if err := fs.Parse(args)`; `317: return whyTarget{}, errors.New("a target is required")` |
| W9 | exit codes | `grep -n "exitOK    = 0\|exitUsage = 1\|exitIntegrity = 3" cmd/ailang-worldd/main.go …/mcpclient.go` | `main.go:74: exitOK = 0`, `:75: exitUsage = 1`, `mcpclient.go:41: const exitIntegrity = 3` |
| W10 | CLI 60 s vs daemon 10 s | `grep -n "whyBudget      =\|budgetContext(whyBudget)" …/why.go`; `grep -n "readDeadline = " host/daemon/daemon.go`; `grep -n "func (d \*Daemon) readCtx" -A2 host/daemon/handlers.go` | `36: whyBudget = 60 * time.Second`; `380: ctx, cancel := budgetContext(whyBudget)`; `144: readDeadline = 10 * time.Second`; `287: … return context.WithTimeout(r.Context(), d.readDeadline)` |
| W11 | deadline printed as "not found" | `grep -n 'errors.Is(err, context.DeadlineExceeded)\|not found in the last %d entries' …/why.go` | `399: if errors.Is(err, errScanExhausted) \|\| errors.Is(err, context.DeadlineExceeded) {`; `400: "ailang-worldd why: not found in the last %d entries (%s)\n"` |
| W12 | read errors become ✗ or are swallowed | `grep -n 'case err != nil:\|return obj, false, err.Error()\|why = werr.Error()\|if err != nil \|\| !ok \|\| obj.SemanticID != broker.EffectRecordV1' …/why.go` | `781: case err != nil:` `782: return obj, false, err.Error()` (checkedObject); `834: why = werr.Error()` (world link); `525: if err != nil \|\| !ok \|\| obj.SemanticID != … { return …, false }` (effectRecord) |
| W13 | `--result` and `-` are a file and a stdin read | `grep -n 'os.ReadFile(resultFile)\|io.ReadAll(io.LimitReader(stdin' …/why.go` | `309: b, err := os.ReadFile(resultFile)`; `319: b, err := io.ReadAll(io.LimitReader(stdin, maxClientResponseBytes+1))` |
| W14 | `--expect-store` uses the process's cwd | `grep -n 'func sameStorePath' -A4 cmd/ailang-worldd/doctor.go` | `369: func sameStorePath(a, b string) bool { ca, errA := filepath.Abs(a) …` |
| W15 | the CLI's response cap | `grep -n "maxClientResponseBytes = \|exceeds %d bytes" cmd/ailang-worldd/cli.go` | `21: maxClientResponseBytes = 8388608`; `81: "response from %s%s exceeds %d bytes"` |
| W16 | `elide` lives in mcpclient.go and is shared | `grep -n "^func elide\|^const elideAt" …/mcpclient.go`; `grep -n "elide(" cmd/ailang-worldd/*.go \| grep -v _test \| grep -v /why.go \| wc -l`; `grep -c "elide(" …/why.go` | `595: const elideAt = 160`, `651: func elide(s string) string`; `5`; `4` |
| W17 | V1's route gives `why` a terminal answer | `grep -n 'a.Outcome = "terminal"\|case v.Words == "episode report"\|defer func() { writeApprovalJSON' host/daemon/surface_verb.go`; `grep -n 'case first == "why"' host/verbs/lookup.go` | `45: defer func() { writeApprovalJSON(w, http.StatusOK, a) }()`; `56: case v.Words == "episode report":`; `63: a.Outcome = "terminal"`; `lookup.go:99: case first == "why":` |
| W18 | the page golden records `why` → terminal | `grep -o '"why": "[^}]*outcome[^,]*' host/surfacepage/testdata/verb-route.golden.json` | `"why": "{…\"verb\":\"why\",\"class\":\"read\",\"outcome\":\"terminal\"` |
| W19 | `result` needs `cli` in code; no test row for it (negative, with control) | `grep -n "const VERB_NEEDS" …/app.js`; `grep -c "result.*cli: ''\|B('episode report w1'), cli" host/surfacepage/app_verb_test.go`; control `grep -c "B('why'), cli: ''" …` | `485: const VERB_NEEDS = { address: [], result: ['cli'], … }`; `0`; control `1` |
| W20 | `stderr` requires exit 1; exit 3 unreadable | `grep -n "b.stderr !== undefined && b.exit !== 1\|'exit 1'" …/app.js` | `499: if (shown.length !== 1 \|\| (b.stderr !== undefined && b.exit !== 1)) return UNREADABLE;`; `500: … body: b.stderr !== undefined ? b.stderr + 'exit 1' : shown[0] …` |
| W21 | an empty Enter returns before superseding; loadIndex awaited before acting | `grep -n "return out();\|const act = await askVerb\|await loadIndex()" …/app.js` | `1267: if (String(raw == null ? '' : raw).trim() === '') return out();`; `1268: const act = await askVerb(raw, …)`; `1283: r = addressTarget(raw, await loadIndex(), parse)`; `1285: if (r.home) return out();` |
| W22 | never rows don't pin the CLI line (with control) | `sed -n 35,37p host/verbs/lookup_test.go`; control `grep -c 'at a terminal: ailang-worldd' …/lookup_test.go` | `"approver is never run from the app: enrolling"`, `"minting and revoking a session"`, `"host administration"`; control `3` |
| W23 | `MaxWords` is pinned only symbolically (negative, with control) | `grep -n "MaxLine *=\|MaxWords *=" host/verbs/split.go`; `grep -n MaxWords host/verbs/split_test.go`; `grep -c "32" …/split_test.go`; control `grep -c "32" …/split.go` | `17: MaxLine = 1024`, `18: MaxWords = 32`; `29: strings.Repeat("w ", MaxWords)`, `41: … MaxWords+1 …`; `0`; control `2` |
| W24 | say fades at 2.6 s; Escape removes the panel; Enter refocuses | `grep -n "2600" …/app.js`; `grep -n "Escape" …/app.js`; `sed -n 1506,1509p …/app.js` | `995: function say(t) { … setTimeout(() => { m.style.opacity = '0'; }, 2600); }`; `1520: else if (k === 'Escape') { if ($('verb')) $('verb').remove(); …`; `1508: if (e.key === 'Enter') { go(e.target.value); $('root').focus({ preventScroll: true }); }` |
| W25 | no `host/why` yet (negative, with control) | `ls host/why`; control `ls -d host/verbs host/daemon` | `ls: host/why: No such file or directory`; control `host/daemon host/verbs` |
| W26 | the daemon spawns nothing (negative, with control) | `grep -rln 'exec.Command' host/daemon --include='*.go' \| grep -v _test \| wc -l`; control `grep -rln 'exec.Command' host --include='*.go' \| grep -v _test \| head -2` | `0`; control `host/agentrun/agentrun.go`, `host/sshsign/sshsign.go` |
| W27 | git pairs a whole-file move with renames (scratch repo `/tmp/v2rename.2oWw`) | `git init`; copy why.go, commit; `git mv cmd/ailang-worldd/why.go host/why/why.go`; `sed` package line + `worldReader→Reader`, `wireEntry→Entry`, `invRecord→Record`, `whyChain→Chain`, `whyLink→Link`, delete lines 258–410; new `cmd/ailang-worldd/why_cli.go` with those 153 lines; `git diff --cached --numstat`; `git config --get diff.renames`; `git --version` | `153 0 cmd/ailang-worldd/why_cli.go`; `56 209 {cmd/ailang-worldd => host/why}/why.go`; *(unset: git's default, renames on)*; `git version 2.54.0` |
| W28 | V1's measured size | `git diff --numstat fb4d8b89^1 fb4d8b89 -- cmd host docs \| awk '{a+=$1;d+=$2} END{print a, d, a+d}'` | `1759 150 1909` (32 files; the moved printer counted `105 0` + `6 99`) |
| W29 | receipt 400; handler timeout is 503 | `grep -n 'restInvocationID(id)' host/daemon/handlers.go`; `grep -n 'func writeReadTimeout' -A2 …/handlers.go` | `420: if !restInvocationID(id) {` → `BadRequest … 400`; `349: writeAPIError(w, "Timeout", "read deadline (%s) exceeded", http.StatusServiceUnavailable)` |
| W30 | why's tests; which reference moved unexported names | `grep -n "^func Test" cmd/ailang-worldd/why*_test.go`; `grep -ohE '\b(whyInputLine\|protectedLine\|resultProtected\|whyKeyField[A-Za-z]*\|runWhy\|whyChain\|whyLink)\b' cmd/ailang-worldd/*_test.go \| sort \| uniq -c` | `why_test.go`: `TestWhyTargetFormsResolveToOneEntry`, `…ResolvesMCPInvocationID`, `…BrokenLinkIsIntegrityRefusal`, `…ScanBound`, `…ShowsTheWorktreeDelta`, `…ShowsTheFailureLine`, `TestLogTailFollowPrintsEachCommitOnce`, `TestProvenanceTrailer`; `why_protected_test.go` 2, `why_input_test.go` 4; `5 whyInputLine`, `2 whyKeyFieldCount`, `1 whyKeyFieldMax`, `1 protectedLine`, `3 resultProtected`, `5 runWhy`, `6 whyChain`, `5 whyLink` |
| W31 | `/v1` GETs are unprotected | `grep -n "func (d \*Daemon) isProtected" -A2 host/daemon/daemon.go` | `1169: return r.Method == http.MethodPost && r.URL.Path == "/v1/commit"` |
| W32 | VS1 pins exact import and call sets; V1's input bounds | `sed -n 130,180p host/daemon/surface_verb_test.go`; `grep -n "MaxLine *=\|MaxWords *=" host/verbs/split.go` | `pure := map[string]bool{"errors", "flag", "fmt", "io", "regexp", "strings", "time", "unicode/utf8"}`; `surface_verb.go` imports `{bytes, encoding/json, fmt, net/http, strings, host/verbs}`; `want := {"readCtx", "episodeReport", "writeInternalErrorLog", "verbEpisodeReport", "readDeadline"}`; `MaxLine = 1024`, `MaxWords = 32` |
| W33 | placeholder and keys row | `grep -n 'id="addr"' host/surfacepage/app/index.html`; `grep -o "\['Enter', 'in the bar[^]]*\]" …/app.js` | `67: … placeholder="row-189@6 · shared@0ffbb9c · back · episode report &lt;episode&gt;"`; `['Enter', 'in the bar: a place, episode report E, approve ID', 'the same']` |
| W34 | a 500-entry log page is ~15 ms | `grep -n "LogRange" bench/BASELINE.md \| head -2` | `268: BenchmarkLogRange/limit_100-16 … 3.064 p50_ms`; `269: BenchmarkLogRange/limit_500-16 … 15.65 p50_ms 17.55 p95_ms` |
| W35 | `headState` stays in cmd for other callers | `grep -n "headState(" cmd/ailang-worldd/*.go \| grep -v _test \| wc -l` | `9` (in demo_seed, demo_proof, episode_report, mcpclient, why) |
| W36 | the V1 judge notes handed to V2 | `sed -n 1859p design_docs/world-mission-log.md` | `Non-blocking, logged here and handed to V2 …: (1) a never-verb's refusal text can lose its CLI line unchecked; (2) an empty submit does not supersede a pending ask … plus X-JS-FINALLY-NULL …; (3) MaxWords not pinned apart from its constant; (4) a result answer without cli is accepted by the page; (5) Escape-closes-panel and Enter refocus untested; (6) … start:/uncommitted intents: … untested; (7) the exact 256-byte episode boundary untested; (8) … RS-SAY` |
| W37 | V1's page tests use `why` as the terminal example | `grep -n "'why'" host/surfacepage/app_verb_test.go \| cut -c1-90` | `106: for (const q of ['why', 'episode finish w1 --preset p'])` (terminal); `129:` VP2 inputs; `138: { ...B('why'), cli: '' }`; `145`, `149`; `161: … app.askVerb('why', rig(golden('why')))` |
| W38 | the walk rig and its daemon (for VE1) | `grep -n "func newWalkRig" -A9 cmd/ailang-worldd/why_test.go` | `52: func newWalkRig(t *testing.T) *walkRig { … daemon.New(…) … srv := httptest.NewServer(d.Handler())` |
