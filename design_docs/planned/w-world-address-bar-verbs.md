# w-world-address-bar-verbs: operator verbs in the World app's address bar, one parser with the CLI

- **Status**: Planned. Designed in iteration 394, **r1** **[r1]** (unattended designer). r0
  (committed `ad144327`) was **BLOCKED** at quorum round 1, 3/3 REJECT (oc-glm-5-3, oc-kimi-k3,
  gpt6-1-sol; none disputed the direction). r1 applies the three fixes the controller adopted and
  nothing else; every changed passage is marked **[r1]**; §15 records each objection and its answer:
  (1) **no silent fallback** — only a schema-valid `outcome:"address"` ever navigates; every route
  failure is an explicit refusal with the view unchanged; (2) the **complete shadow set** of 27 CLI
  words and `world get <ref>` as a refused CLI verb; (3) the CORS, Shell-character and
  `serve --help` verification gaps.
  Parent design for queue row **243**, with child phases V1–V5. **Only V1 is specified to build
  (§2–§9).** V2–V5 are scoped and sized at parent level (§0.3) and each needs its own child design
  r0 before it is planned.
- **Authority**: **D-WORLD-107 = A** (Mark, attended 2026-11-02: "yeah the address bar accepting
  verbs sounds good"): the address bar accepts operator verbs in the CLI's grammar, shows the CLI's
  prompt, and is confirmed by row 201's key approval, never a typed `y` from a page. **D-WORLD-109 =
  A** (Mark, attended 2026-11-06): the loop may take row 243 after row 201 through design → quorum →
  plan → build → judge; it builds and judges against fakes and test keys only and never enrols a
  real key, mints or publishes; Mark's no-terminal exit run stays attended (§10). **D-WORLD-94**: the
  goal is all of Mark's agent work through World. **D-WORLD-89** (row 201's R1–R5) binds every
  section: the agent never approves; publishing stays typed-only (R4); the typed `y` keeps working
  at a terminal; approving is a POST bound to a verified signature, never a cookie or a session.
  Row 201's settled siblings (`-a0a1`, `-a1b`, `-a2`, `-a3`, `-a4`) are **not reopened** (§0.2).
- **Item** (row 243, V1): the World app's address bar (`#addr` in the app shell) learns the CLI's
  operator verbs through **one Go grammar package, `host/verbs`**, used by the CLI and by a new
  read-only daemon route `GET /surface/verb.json`. In V1 the bar **runs** `episode report <episode>`
  (byte-identical text to the CLI), **navigates** `approve <request-id|request-ref>` to row 201's
  existing key-approval cards, and for every mutating or host-process verb shows the exact CLI line
  and that the app cannot run it yet. V1 changes no World state from a page, adds no write, no
  process spawn and no `.ail`.
- **Designer lane**: `claude:claude-opus-5-5`.
- **Estimated**: whole item **≈ 5 phases, ≈ 10,000 changed lines** (parent, low confidence, §0.3);
  **V1 = 1 sprint ≈ 1,750 changed lines** (`-- cmd host docs`). Calibration: row 201 A4's four
  sprints measured 1,985 / 1,983 / 1,493 / 1,947 and A3a's measured 1,995 / 1,861 / 393 changed
  lines (V20), against plans of ≈ 1,950 / 1,800 / 1,500 / – (A4 §0 header) — about **1.0–1.1×**.
  V1 at 1,750 × 1.1 = 1,925 stays under the ~2,000 cap; §9 names the cut if it does not.
- **Measurement base**: `48c92a521fcbf46a0f1120609de80d656326ba5b` (origin/dev, PR #575's merge),
  worktree `/Users/voightkampff/dev/sunholo-data/ailang-world/.wt-world-iter394-row243`, branch
  `design/world-iter394-row243`, measured 2026-11-16 (V1). No browser, no daemon and no Go test was
  run in this fire; every claim below is from reading the tree at that sha (§14).

---

## 0. Scope, the amendment posture, and the sprint boundaries

### 0.1 What the row asks, measured against the tree

The row's sentence is "the address bar accepts operator verbs in the CLI's own grammar (for
example `start world <task>`, `finish <episode>`, `report <episode>`, `why <entry>`), with one
parser shared by the CLI and the view". Four premises do not hold as written, and the design
follows the tree:

- **P1 — the examples are not the CLI's grammar.** The CLI's operator verbs are `episode start
  <episode> …`, `episode launch <episode> --preset <name>`, `episode finish <episode> --preset
  <name>`, `episode report <episode>`, `why <target>` and `approve <id|ref>` (V2, V3). There is no
  `start world` and no top-level `finish` or `report`. Worse, **`world` is already a CLI verb**
  (`world get <ref>`, V4) **and** the address bar's navigation keyword (`world NAME[@N]`, V5).
  V1 keeps the CLI's words exactly (an optional leading `ailang-worldd` is dropped so a pasted CLI
  line works) and keeps `world …`/`back` as navigation. Short aliases are Q1 (§12).
- **P2 — there is no Go address parser to share.** The bar's parser is JavaScript,
  `parseAddress` in `host/surfacepage/surface.js:81`, and the app shell calls it through
  `addressTarget` (`app.js:325`) from `go(raw)` (`app.js:1201`) (V5, V6). The CLI parses argv with
  a hand-rolled verb `switch` (`main.go:380-…`) and one `flag.FlagSet` per verb, mostly through
  `parseInterleaved` (`session.go:181`) (V3, V7). "One parser" can therefore only be a Go parser
  that the CLI calls in-process and the page reaches through a daemon route. That is V1's shape.
- **P3 — "before anything runs, the app shows the same prompt the CLI shows" is not true of the CLI
  itself for `finish`.** `episode finish` commits in the worktree and prints `committed at <sha>`
  (`episode_finish.go:379`) **before** `finishPrompt` asks `Open this pull request? [y/N]`
  (`:409`, `episode_finish_operator.go:154-171`) (V8). The prompt's `sha:` line is the commit just
  made. A page-initiated finish must therefore either pre-approve the intent or bind its approval
  to a side-effect-free preview; that choice is V3's (Q3), not V1's.
- **P4 — half of the no-terminal start already exists, but both ends are CLI processes.**
  `episode start --request` stores a mint request and prints the request block
  (`mint_request.go:101-153`, `renderRequest` `:156-`); the A3a card `GET /surface/approve/{id}`
  takes a passkey or SSHSIG approval; and `approve <id> --await` is the minter that waits for it
  (`approve_await.go:1-8`) (V9, V10). But `--request` refuses a launch preset (`episode.go:249-250`)
  and the minter is a CLI process: someone still runs two commands. V4 (the host runner) closes
  that; V1 does not.

### 0.2 Amendment posture

- **V1 amends nothing.** It adds a read-only route under the A3a Host allowlist (`approvalRoute`,
  R2), reuses row 201's cards by navigation only, and leaves every POST, law, judge, store and
  `/v1` route untouched (§8).
- **Flagged for V3, not decided here.** A4 §2.6 L1 makes finish-shaped scopes **reserved** to their
  typed owner (`episode finish`) for every A4 decision format (`host/broker/decision_law.go:80-81`,
  `:116`) (V11). A key-confirmed finish from the bar would be `episode finish` itself accepting a
  verified key instead of the typed `y`, not an A4 decision — so L1's code is not contradicted —
  but A4's stated reason ("these scopes have their own typed owners … that show the exact bytes and
  run their own fences") is. V3's design must carry an explicit amendment section and a ledger ask
  (provisionally: "may `episode finish`'s confirmation be a verified enrolled key over the finish
  scope, as D-WORLD-89 R1 allows for mints?", recommended **yes**). D-WORLD-89 R4 (publishing
  typed-only) is untouched by every phase: `world-publish` is a different binary and never in the
  verb table (§2.1).

### 0.3 The phases (parent) and the sprint boundaries

| phase | what lands | state change from a page? | size (est.) | status |
|---|---|---|---|---|
| **V1** | `host/verbs` (tokenizer, verb table, the shared `ParseInterleaved`, `episode report`'s flag set); `GET /surface/verb.json`; the bar runs `episode report`, navigates `approve`, explains the rest | **none** | 1 sprint ≈ 1,750 | **this doc** |
| V2 | `why <target>` in the bar: the `why` walk (`why.go`, 1,050 lines, V12) behind a reader interface, served in-process | none | 1 sprint ≈ 1,600 | child design r0 |
| V3 | the verb-approval law: a `world-verb/v1` subject over the canonical argv and the exact prompt bytes, A4's judges (SSHSIG/passkey), credentials schema 9, a card; `episode finish` and `episode start` gain a "confirm on the surface" mode that replaces the typed `y` with a verified key (typed `y` unchanged at a terminal) | only after a verified key | 2–3 sprints ≈ 4,500 | child design r0 + ledger ask (§0.2) |
| V4 | the host runner: who runs an approved verb's CLI process with no terminal, and the hand-off of an interactive `episode launch` to herdr (`integrations/herdr/ailang-world`, V13) or row 206's agent tab | only after a verified key | 1–2 sprints ≈ 2,000 | child design r0 (Q4) |
| V5 | readiness for Mark's exit run: QUICKSTART section, the attended script (§10) | none | < 300 | with V4 |

Nothing flows backwards: V1 lands the grammar and the read path before anything can mutate; V3
lands the law before V4 lands anything that can execute a mutating verb without a terminal (the
row-201 A1b/A2/A3 order: law first, then the surface, A4 §0).

**Out of V1:** `why` in the bar (V2); any mutating verb running from the bar (V3/V4); the static
page's address bar (`surface.js:331`, a page that may be served without a daemon; unchanged);
aliases (Q1); `--follow`, `--markdown`, `--interval`, `--addr` and `--help` from the bar (§2.1); any
change to `approve`'s CLI, the cards, or row 201's routes.

## 1. What the code does today (measured at `48c92a52`)

1. **The CLI's dispatch is a hand-rolled switch on the first word** after the global flags
   (`rest := globals.Args()`, `main.go:365`; `switch verb := rest[0]`, `main.go:380`), with
   `episode` dispatching `start|launch|finish|report` in `runEpisode` (`episode.go:165-187`) (V3).
   Thirty-four non-test `flag.NewFlagSet` calls exist in `cmd/ailang-worldd` (V7). `ownHelpVerbs`
   (`main.go:698-704`) lists the verbs that print their own help; it is not a full verb list (V14).
2. **`parseInterleaved(fs, args)`** (`session.go:181-193`) parses flags and positionals in any
   order by re-running `fs.Parse` on the remainder; it is used by `approve`, `approver`, `episode
   start`, `episode finish`, `episode launch`, `episode report`, `preset`, `session` and `session
   new` (11 matching lines including the definition) (V7).
3. **`episode report`** builds its own `FlagSet` with `--markdown`, `--json`, `--follow`,
   `--interval` (default 2s) and `--addr` (`episode_report.go:49-56`), refuses `--markdown` with
   `--json` or `--follow`, then GETs `/v1/episodes/<episode>/report` (`fetchReport`, `:105`) and
   prints `printReport` (`:225-254`) with `execCacheLine`, `readsLine`, `printRefusals`,
   `sortedCounts` (`:256-315`); a 404 prints `ailang-worldd episode report: episode <e> has no
   committed calls` and exits 1 (`reportOnce`, `:127-133`) (V15). Its pins are
   `TestEpisodeReportUsage`, `…TablePrintsEveryClass`, `…TablePrintsTheReadsLine`,
   `…JSONAndTableLoopback`, `…FollowReprintsOnChangeLoopback`, `…PrintsExecCache`,
   `TestEpisodeGlobalAddrRefusedForStartOnly`, `TestGlobalBareAddrIsNormalisedForEveryClientVerb`
   (`cmd/ailang-worldd/episode_report_test.go`) (V16).
4. **The report route** `GET /v1/episodes/{episode}/report` (`daemon.go:1120`) is handled by
   `handleEpisodeReport` (`episode_report.go:707-…`): it bounds the episode to
   `maxReportEpisodeBytes`, builds under `d.readCtx`, adds `Start`, `Finish`, `ExecCache`, answers
   404 `NotFound` when all are empty and otherwise `writeJSON` (`json.NewEncoder(w).Encode`,
   `handlers.go:130-134`) (V15). `/v1` reads are unauthenticated GETs (the `why` header,
   `why.go:11`, "The reads are unauthenticated GETs") (V12).
5. **The address bar.** The app shell's input is `<input id="addr" … placeholder="row-189@6 ·
   shared@0ffbb9c · back">` (`app/index.html:67`). `go(raw)` (`app.js:1201-1209`) calls
   `addressTarget`, which treats empty input as home and prefixes `world ` to anything not starting
   with `world`/`back` before `parseAddress` (`app.js:325-333`); `parseAddress` accepts exactly
   `world`, `back`, `world NAME` and `world NAME@AT` and otherwise answers `unknown address: try
   "world", "world NAME", "world NAME@N", "world NAME@HASH" or "back"` (`surface.js:11`, `:81-98`)
   (V5, V6). The pure part runs under node in `TestSurfaceAppPure` (`app_test.go:122-135`), which
   CI makes non-skippable with `WORLD_EXEC_NODE` (`ci.yml:250`, `:289-290`) (V17).
6. **Names a verb would shadow.** Episode names match `^[a-z0-9][a-z0-9-]{0,63}$`
   (`doctor.go:461`), so `episode`, `why` and `approve` are legal brane/episode names (V18). Today a
   bare `why` in the bar means `world why`.
7. **The surface's write guards (A3a).** `approvalRoute` computes the origin set `O` (the
   loopback `http://localhost:<port>` of the connection, every `--surface-origin`, every active
   passkey's origin) and answers **421** for a `Host` outside `O` before the handler (R2,
   `surface_approve.go:58-90`). The approve POST then checks `Origin ∈ O` (R3), `Sec-Fetch-Site`
   empty or `same-origin` (R4), `Content-Type: application/json` (R5), a size cap (R6), and judges a
   signature before its one non-waiting `OpenCredentials` (R11–R13, `:203-292`)
   (`surface_approve_guard.go:53-64`) (V19). `surfaceCSP` has `connect-src 'self'`
   (`surface_page.go:19`) (V21). The approval routes are pinned by name in `approvalRoutes`
   (`surface_approve_test.go:20`) (V22).
8. **The cards V1 navigates to exist.** `GET /surface/approve/{id}` and `GET /surface/decide/{ref}`
   are registered behind `approvalRoute` (`daemon.go:1141`, `:1145`); a mint request id is
   `^mr-[0-9a-f]{16}$` (`store.MintRequestIDPattern`, `mint_requests.go:46`) and a decision ref is
   `sha256:` and 64 hex (`approve.go:99`) (V10, V23).
9. **The daemon spawns no process.** No non-test file in `host/daemon` calls `exec.Command`, while
   `host/agentrun`, `host/sshsign` and `host/broker` do (positive control) (V24). V1 keeps it so.
10. **No `host/verbs` exists** (V25), and `TestDaemonDependencyAllowlist` allows any package of this
    repository (`daemon_test.go:899-…`, its classifier case `host/store` → module path, allowed as
    repo-internal) (V26), so `host/daemon` may import a new repo package.
11. **[r1] No CORS header is sent by the daemon.** The only `Access-Control`/`cors` matches in
    non-test `host/daemon` are three comments saying so (`surface_approve.go:31`, `:201`;
    `surface_approve_guard.go:60`); the same tree's `Header().Set(` calls set only
    `Cache-Control`, `Content-Security-Policy`, `Content-Type`, `Referrer-Policy`, `Retry-After` and
    `X-Content-Type-Options` (positive control: the instrument finds header writes) (V31). `/v1`
    GETs carry no auth wrapper: `GET /v1/episodes/{episode}/report` is registered as the bare
    `d.handleEpisodeReport` (`daemon.go:1120`, V15), and the CLI reads it with no credential in
    `TestEpisodeReportJSONAndTableLoopback` (V16).
12. **[r1] The page and the route ship in one binary.** `app.js`, `index.html` and `surface.js` are
    `//go:embed`ded in `host/surfacepage` (`app.go:15`, `:29`; `script.go:9`) and served by the
    daemon itself (`GET /surface/app.js` → `writeSurfaceAsset(…, surfacepage.AppScript)`,
    `surface_page.go:161-162`) with `Cache-Control: no-cache` (`:190`) (V32). A page that calls
    `/surface/verb.json` therefore comes from a daemon that serves it; a 404 can arise only from a
    tab left open across a downgrade or a proxy misroute. V1 needs **no legacy-daemon path**; a 404
    is one more explicit refusal (§2.3).
13. **[r1] The CLI's switch has exactly 28 arms** besides `default`: `approve approver call commit
    demo doctor episode head health help hosts log object outcome preset profile provenance
    registry serve service session setup surface-index tools version whoami why world` (V33; the
    controller's count, re-derived). `world` dispatches to `requireGet(args, "world", …)`, which
    takes exactly `get <ref>` (`cli.go:157-169`) (V34).
14. **[r1] `serve --help` returns before `serve` binds anything.** `run()` calls
    `genericVerbHelp(rest)` before the switch (`main.go:374-378`); for `serve` (not in
    `ownHelpVerbs`) with `--help` it returns the usage lines `usageLinesFor("serve")` finds
    (`main.go:708-740`, `:750-765`; usage line `main.go:86`), and `run` returns `exitOK`. The
    existing `TestEveryVerbHasHelp` already runs `run(["serve","--help"])` and `run(["help",
    "serve"])` expecting exit 0 (`help_test.go:9-33`) (V35). VB4 nevertheless calls no `run()`
    (§6).

## 2. Design (V1)

### 2.1 `host/verbs` — the grammar (pure Go, no I/O)

One package, imported by `cmd/ailang-worldd` and `host/daemon`. Four parts:

**(a) `Split(line string) ([]string, error)`** — the bar's tokenizer, the only thing the CLI does
not need (the shell splits its argv). Rules, each a named refusal:
- `TooLong`: more than 1,024 bytes; `TooMany`: more than 32 words; `NotText`: invalid UTF-8 or any
  byte < 0x20 or 0x7f (so no newline, no tab smuggling a second command);
- words are separated by spaces; `'…'` is literal; `"…"` is literal except `\"` and `\\`;
  `Unclosed`: an unclosed quote;
- `Shell`: an unquoted `$`, backtick, `;`, `|`, `&`, `<`, `>`, `(`, `)` — **nine characters
  [r1]** (the r0 test plan's "ten" was wrong) — the bar runs no shell, and
  refusing these keeps "what the bar shows" equal to "what a shell would run" for the line it
  prints.
- **`Quote(argv) string`** is its inverse for display (`ailang-worldd episode report w1`), with
  the property `Split(Quote(a)) == a` for every admissible `a` (tested, §6).

**(b) `Lookup(argv []string) (Verb, []string)`** — the verb table. Order matters and is the whole
classification (traced in §2.4):
1. a leading word `ailang-worldd` is dropped once;
2. a first word starting with `-` → `class=refused`, reason "global flags are not accepted: the app
   talks to its own daemon" (the CLI's `--addr`, `--db` etc. never reach the bar);
3. **[r1]** the bare navigation shapes only: `world`, `world NAME`, `world NAME@AT` (one or two
   words) or `back` → `class=address` (navigation keeps both words; P1);
3b. **[r1]** `world get <ref>` (first word `world`, second `get`, three or more words) →
   `class=other`, refused with the CLI line ("`world get` is a CLI command the app does not run; at a
   terminal: `ailang-worldd world get <ref>`") (it is a CLI verb, V34). Step 3 now captures only
   bare navigation shapes, so a pasted CLI command is never reinterpreted as brane navigation. 3b
   shadows nothing: `parseAddress` accepts only one- or two-word `world` input (`surface.js:84-86`,
   V5), so a three-word `world …` is an address error today. Any other `world` input of three or
   more words falls to step 5 and gets `parseAddress`'s own error, as today; `world get` (two words)
   is step 3, navigation to a brane named `get`, as today;
4. a table verb, longest match first:

| words | class | V1 in the bar | why this class holds |
|---|---|---|---|
| `episode report` | **read** | runs (§2.2) | GETs one route, writes nothing (§1.3) |
| `approve` | **navigate** | `approve mr-<16 hex>` → `/surface/approve/<id>`; `approve sha256:<64 hex>` → `/surface/decide/<ref>`; any flag or other shape refused | the decision is the card's key flow (row 201 A3a/A4), never the bar's; the bar only opens it |
| `why` | **read** | `terminal` answer (V2) | read-only, but its walk lives in the CLI (§1, V12) |
| `episode start` | **mutate** | `terminal` answer (V3/V4) | mints a session or stores a mint request (§0.1 P4) |
| `episode finish` | **mutate** | `terminal` answer (V3) | commits, records an approval, opens a PR (§0.1 P3) |
| `episode launch` | **host** | `terminal` answer (V4) | starts an agent process that needs a terminal (`episode_launch_interactive.go:18-21`, V27) |
| `approver`, `session`, `service`, `serve`, `setup` | **never** | refused, naming the reason | enrolment/revocation are typed-only (D-WORLD-89 R3); mint/revoke and host administration stay at a terminal |
| **[r1]** `commit`, `outcome`, `call`, `demo` | **other** (writes World) | refused: "`<verb>` is a CLI command the app does not run; at a terminal: `<line>`" | they write World with a session credential (`commit` = `POST /v1/commit`; `outcome record`; `call` = MCP `tools/call` on `/mcp/`; `demo seed --db --session`, main.go:109, :164, V33); none has a key-approval path, so none may ever run from a page before V3's law, and none is an operator verb |
| **[r1]** `doctor`, `head`, `health`, `help`, `hosts`, `log`, `object`, `preset`, `profile`, `provenance`, `registry`, `surface-index`, `tools`, `version`, `whoami` | **other** (reads or local tooling) | refused, same text | not an operator verb (D-WORLD-107's list is start, launch, finish, approve, report; the row adds why) |

5. anything else → `class=address` (the page's existing `parseAddress` path).

**[r1] The classification is exhaustive and its sets are fixed by the switch.** The 28 arms (§1.13,
V33) split exactly as: navigation token `world` (steps 3/3b); operator `{why, episode, approve}`
(3); never `{approver, session, service, serve, setup}` (5); other `{commit, outcome, call, demo,
doctor, head, health, help, hosts, log, object, preset, profile, provenance, registry,
surface-index, tools, version, whoami}` (19), plus the `world get` shape of step 3b. 1 + 3 + 5 + 19
= 28. VB4 (§6) asserts the go/ast-derived arm set minus the navigation token `world` (and never
counting `back` or the dropped `ailang-worldd`, which are not arms) equals operator ∪ never ∪
other, each arm in exactly one set, so a CLI verb added or removed without a class fails the
build. These 27 words are the **shadow set** (RK1).

**(c) `ParseInterleaved(fs *flag.FlagSet, args []string) ([]string, error)`** — moved verbatim from
`session.go:181-193`; the CLI's ten call sites call `verbs.ParseInterleaved` (§1.2). Behaviour is
identical by construction (same body), pinned by every CLI test that reaches those verbs, unedited.

**(d) The per-verb flag sets.** `EpisodeReportFlags(fs *flag.FlagSet) *EpisodeReportArgs` registers
exactly today's five flags with today's names, defaults and usage strings (§1.3). The CLI's
`runEpisodeReport` calls it instead of declaring its own; the bar's admission calls the same
function on a fresh `FlagSet`, runs `ParseInterleaved`, then applies the **bar profile**: allow
`--json`; refuse `--follow` and `--interval` ("the app shows one report; it re-reads on Enter"),
`--markdown` ("the PR body is rendered by `episode finish`"), `--addr` ("the app reads its own
daemon"), `-h`/`--help` ("see `ailang-worldd episode report --help`"); require exactly one
positional. So the CLI and the bar run **the same flag definitions and the same parse**; they differ
only in a named admission profile applied after it. V3 migrates `episode start`/`finish` the same
way when the bar learns to run them; V1 migrates only the verb it runs.

### 2.2 `GET /surface/verb.json?q=<line>` (read-only)

Registered behind `d.approvalRoute` (so R2's Host allowlist and 421 apply exactly as on the
approval routes). It has no POST, no `OpenCredentials`, no store write, no `exec`.

1. **[r1]** `q` absent → 400 `BadRequest` (the page never sends that). A `q` over 1,024 bytes is an
   operator-input refusal, not an HTTP error: **200** `{outcome:"refused", reason:"TooLong"}`
   (checked on the raw query value before `Split`, so no oversized input is tokenized).
2. `Split` → a refusal answers **200** `{outcome:"refused", reason}` (an operator typo is not an HTTP
   error).
3. `Lookup` → by class:
   - `address` → `{outcome:"address"}`; the page runs its existing path unchanged.
   - `read` + `episode report` → admission (§2.1d) → the in-process report: `handleEpisodeReport`'s
     body is extracted unchanged into `d.episodeReport(ctx, ep) (EpisodeReport, found bool, err)`,
     which both the `/v1` handler and this route call. Found → `{outcome:"result", cli, text}` where
     `text` is `daemon.WriteEpisodeReportText(rep)` — `printReport` and its four helpers **moved**
     from `cmd/ailang-worldd/episode_report.go` into `host/daemon`, and the CLI's `printReport`
     becomes a one-line call to it. With `--json`, `json` is the `json.NewEncoder(...).Encode` bytes
     the `/v1` route writes. Not found → `{outcome:"result", exit:1, stderr:"ailang-worldd episode
     report: episode <e> has no committed calls"}`, the CLI's line. A timeout or store error → 200
     `{outcome:"refused", reason}` with the `/v1` route's text, never a raw store error (the
     `internalErrorMessage` rule, `handlers.go:136`).
   - `navigate` → `{outcome:"navigate", href}`, `href` built server-side from an id matching
     `MintRequestIDPattern` or a `sha256:`+64-hex ref; the card itself checks existence.
   - `read` (`why`), `mutate`, `host` → `{outcome:"terminal", class, cli, reason}`; `cli` is
     `Quote(["ailang-worldd", argv…])`.
   - `never`, `other` → `{outcome:"refused", class, reason, cli}`.
4. Every answer carries `{v:1, schema:"world/surface-verb/v1", q, argv, verb, class}`.

The response is a GET with no state change, so it needs no Origin/CSRF arm: a cross-origin page
cannot read it (no CORS headers are sent, §1.11 **[r1]**, V31; same-origin policy), and the only
data it returns is the episode report that `GET /v1/episodes/{episode}/report` already serves
unauthenticated (§1.4, §1.11).

### 2.3 The page (`app.js`)

`go(raw)` becomes: empty → home (unchanged, no fetch); else `fetch('/surface/verb.json?q=' +
encodeURIComponent(raw))`:
- **[r1] No silent fallback.** Only a successful, schema-validated `outcome: address` invokes the
  existing `addressTarget` path. "Schema-validated" means: HTTP 200, a JSON object, `schema ===
  "world/surface-verb/v1"`, `v === 1`, and `outcome` one of the five known values with that
  outcome's required fields. Network errors, non-2xx responses (400, 404, 421, 5xx) and invalid
  responses display an explicit refusal — `say('the verb route did not answer: HTTP <status>')`,
  `say('the verb route is unreachable: <error>')` or `say('the verb route answered something this
  page cannot read')` — and **leave the current view unchanged** (no navigation, no panel, the
  input kept). A 404 is not treated as "an older daemon": the page and the route ship embedded in
  one daemon binary (§1.12, V32), so V1 offers **no** "Navigate as an address" action; if one is
  ever wanted, it must be a separate user-selected action that shows the address interpretation
  before proceeding and is never invoked automatically (gpt6-1-sol's fix, recorded for later). The
  typed `world NAME` form needs the route too; while the route fails, nothing navigates from the
  bar, which is the honest state.
- `address` (schema-valid) → the existing `addressTarget` path, unchanged;
- `result` → a verb panel in the stage: a header line `$ <cli>` and the `text` (or `json`, or the
  `stderr` line and `exit 1`) set with `textContent` — never `innerHTML`;
- `navigate` → `location.assign(href)` only if `href` starts with `/surface/approve/mr-` or
  `/surface/decide/sha256:`; otherwise `say('refused: unexpected target')`;
- `terminal` → the panel: "`<verb>` changes World (or: needs a host process). The app cannot run it
  yet. At a terminal: `<cli>`" with a copy button (copying is a page-local act);
- `refused` → `say(reason)`.

The placeholder gains `· episode report <episode>`. The `?` keys popover gains one row. The pure
part (`verbOutcome(body) → action`) is exported to the node harness like `addressTarget`.

### 2.4 Input-class traces (no earlier step pre-empts a new branch)

Entry points: the mux pattern `GET /surface/verb.json` → `approvalRoute` (its only pre-handler
step is R2, which refuses on `Host` alone, before reading `q`) → the handler's steps 1–4. On the
page: `go(raw)` → empty check → fetch → outcome switch.

| input | page | R2 | step 1 | Split | Lookup step | outcome |
|---|---|---|---|---|---|---|
| `` (empty) | home, no fetch | – | – | – | – | unchanged |
| `row-189@6` | fetch | pass | pass | `[row-189@6]` | 5 | `address` → `parseAddress('world row-189@6')`, unchanged |
| `world why` | fetch | pass | pass | `[world why]` | 3 (before 4) | `address` → brane `why` |
| `back` | fetch | pass | pass | `[back]` | 3 | `address` |
| `why` | fetch | pass | pass | `[why]` | 4 | `terminal` (shadows brane `why`; Q2) |
| `episode report w1` | fetch | pass | pass | 3 words | 4 | `result` |
| `ailang-worldd episode report w1 --json` | fetch | pass | pass | 5 words | 1 then 4 | `result` with `json` |
| `episode report w1 --follow` | fetch | pass | pass | | 4 → admission | `refused` |
| `--addr x episode report w1` | fetch | pass | pass | | 2 (before 4) | `refused` |
| `episode report w1; rm -rf ~` | fetch | pass | pass | `Shell` | – | `refused` |
| `approve mr-0123456789abcdef` | fetch | pass | pass | | 4 | `navigate` |
| `approve mr-0123456789abcdef --deny` | fetch | pass | pass | | 4 | `refused` |
| `episode finish w1 --preset p` | fetch | pass | pass | | 4 | `terminal` |
| `approver enroll --name x` | fetch | pass | pass | | 4 | `refused` (R3) |
| `log tail` | fetch | pass | pass | | 4 (other) | `refused` — **[r1]** a deliberate change: today `parseAddress('world log tail')` answers "unknown address: …"; V1 answers the CLI refusal and line instead (both are errors; the text and the path change) |
| **[r1]** `log` (bare) | fetch | pass | pass | `[log]` | 4 (other) | `refused` with `cli` — **shadow-set change**: today it navigates to a brane named `log` (`world log`); RK1 |
| **[r1]** `world session` | fetch | pass | pass | `[world session]` | 3 (before 4) | `address` → brane `session`, unchanged (the shadowed name stays reachable) |
| **[r1]** `ailang-worldd world get world-1` | fetch | pass | pass | 4 words | 1, then 3b | `refused`, `cli` = `ailang-worldd world get world-1`; not address |
| **[r1]** `world get` | fetch | pass | pass | 2 words | 3 | `address` → brane `get`, unchanged |
| **[r1]** > 1,024 bytes | fetch | pass | **200** `refused`, `TooLong` | – | – | `say(reason)`; view unchanged |
| **[r1]** `why`, with a foreign `Host` | fetch | **421** | – | – | – | `say('the verb route did not answer: HTTP 421')`; view unchanged; **not** `world why` |
| **[r1]** `why`, page bug drops `q` | fetch | pass | **400** | – | – | `say('… HTTP 400')`; view unchanged |
| **[r1]** `why`, route answers 404 (tab kept across a downgrade, or a proxy misroute) | fetch | – | – | – | – | `say('… HTTP 404')`; view unchanged; **not** `world why` |
| **[r1]** `why`, route answers 500 | fetch | pass | (handler fault) | – | – | `say('… HTTP 500')`; view unchanged |
| **[r1]** `why`, network error (daemon down) | fetch throws | – | – | – | – | `say('the verb route is unreachable: …')`; view unchanged |
| **[r1]** `why`, 200 with a body missing `schema` or with an unknown `outcome` | fetch | pass | pass | | | `say('… cannot read')`; view unchanged |

Step 3 precedes step 4, and step 3b catches `world get <ref>`, so `world` is never treated as the
CLI's `world get` and a pasted `world get` line is never navigation **[r1]**. Step 2 precedes step
4, so a global flag never reaches a verb's admission. **[r1]** No row reaches `addressTarget`
except through a 200, schema-valid `outcome:"address"`.

## 3. Pure core, effects at the boundary (coding standards S1–S3)

- **S1.** No `world/` function changes. `host/verbs` is pure Go (no I/O, no clock, no randomness);
  `Split`/`Quote`/`Lookup`/admission are table-tested and property-tested (§6). S1's Z3 clause
  applies to `world/` AILANG; nothing there is touched.
- **S2.** The only effect V1 adds is a read inside the daemon's existing report builder, under
  `d.readCtx`, behind an existing route's code path.
- **S3.** No `.ail` changes; no kernel change; no registry package. The `/v1` table stays at 13
  routes (the new route is under `/surface`, which the sketch's `routes()` does not mirror; §8).

## 4. Threats (V1)

| threat | answer in V1 |
|---|---|
| T1. A script or an agent in a browser submits `episode finish …` to change World | V1 has no mutating path: the handler has no store write, no `OpenCredentials`, no `exec` (VS1 pins it statically; `TestDaemonCredentialWriterSites` is unedited and still passes). A mutating verb only ever answers `terminal`. |
| T2. SSRF via `--addr` (the CLI's report fetches any URL) | the bar reads the report in-process; `--addr` and global flags are refused (§2.1 step 2, §2.1d). |
| T3. File reads via a path flag (`why --result`, `episode start --task`) | not admissible in V1: `why` and `start` answer `terminal` before any flag parse. |
| T4. Hostile text in a report (tool names, episode names, refusal reasons from agents) | rendered with `textContent`; a hostile-value node test (VP3). |
| T5. DNS rebinding or a foreign `Host` | R2 via `approvalRoute` (421), as the approval routes. |
| T6. Cross-origin read of the answer | no CORS headers (§1.11, V31 **[r1]**); and the data is already open at `/v1` (§1.4). |
| T7. Open redirect through `navigate` | `href` is built server-side from two strict patterns; the page re-checks the prefix. |
| T8. Confusion: a brane named after any of the 27 shadow-set words (RK1) **[r1]** | the verb wins on the first word; `world NAME` always navigates (Q2, VP4, VP6). A route failure never turns a verb into navigation (§2.3, VP2) **[r1]**. |
| T9. A copied `cli` line runs something else in a shell | `Split` refuses shell metacharacters and `Quote` round-trips (VB2). |
| T10. Denial of service | `q` ≤ 1,024 bytes; one report build bounded by `d.readCtx`, the same as the open `/v1` route. |

## 5. By-construction argument

1. **No page input changes World in V1**, because the only new server code path is a GET handler
   whose callees are `Split`, `Lookup`, the admission, `d.episodeReport` (the `/v1` builder,
   read-only) and pure renderers. VS1 asserts the handler file and `host/verbs` import neither
   `host/store` write APIs nor `os/exec`, and that the daemon still has exactly its pinned
   credential writer sites.
2. **One parser**: the CLI and the route call the same `verbs.ParseInterleaved` on a `FlagSet` built
   by the same `verbs.EpisodeReportFlags`; the bar's difference is a named admission applied after
   the parse (§2.1d). VB5 mutates the flag set in one place and both sides change.
3. **Same output**: the CLI's `printReport` *is* `daemon.WriteEpisodeReportText` after the move;
   VR3 compares the route's `text` with the real CLI's stdout against one loopback daemon and store.
4. **[r1] Navigation never regresses for inputs outside the enumerated shadow set (RK1)**: every
   other input reaches the unchanged `addressTarget` path through a 200, schema-valid
   `outcome:"address"`; every route failure is an explicit refusal that leaves the view unchanged
   and never selects navigation (§2.3, VP2); and the shadow set itself is pinned exhaustively
   (VB3/VB4), never silently extended. A shadowed name stays reachable as `world NAME` (VP6).

## 6. Test plan (V1, one sprint)

Each row names the mutant it kills. "Real CLI" means `run([]string{…})` in-process against a
loopback daemon, the existing `TestEpisodeReportJSONAndTableLoopback` pattern — fixtures are what
the real CLI and the real route produce, never hand-typed.

**`host/verbs` (unit)**
- **VB1** `Split`: one case per refusal (`TooLong` at 1,025 bytes and not at 1,024; `TooMany` at 33
  and not 32; `NotText` for `\n`, `\t`, 0x7f, invalid UTF-8; `Unclosed`; `Shell` for each of the nine
  characters **[r1]** unquoted, and accepted inside quotes). Kills: each refusal deleted; off-by-one bounds.
- **VB2** `Quote` round-trip property over generated argv (seeded `math/rand`, 2,000 cases including
  quotes, spaces, backslashes) and every argv in VB3. Kills: a `Quote` that drops escaping.
- **VB3** `Lookup` table: every row of §2.1(b)'s table and every row of §2.4's trace, **[r1]**
  including bare `log` → `other`/refused (was brane navigation), `world session` → address,
  `ailang-worldd world get world-1` → `other`/refused with `cli` = `ailang-worldd world get
  world-1`, `world get` → address, and one bare word from each of the never and other sets.
  Kills: steps 2/3 swapped with 4; `world` classified as the CLI verb; **[r1]** step 3b deleted
  (a pasted `world get` line becomes navigation) or widened to two words (brane `get` lost); a
  missing table row; longest-match broken (`episode` alone → refused "needs a subcommand", not
  `episode report`).
- **VB4 [r1]** parity with the CLI's switch, **without invoking `run()`**: the test parses
  `cmd/ailang-worldd/main.go` with `go/ast`, finds `run`'s `switch verb := rest[0]` and collects
  its non-default `case` string literals (28 at the base, V33); it asserts that set minus the
  navigation token `world` equals operator ∪ never ∪ other from `host/verbs`, each arm in exactly
  one set (`back` and the dropped `ailang-worldd` are excluded: they are not arms). Help check, also
  without `run()`: every word in the three sets has a help path — `ownHelpVerbs[w]`, or
  `genericVerbHelp([]string{w, "--help"})` returns ok, or `w == "demo"`/`"help"` (their own
  branches, `main.go:370-372`, `:710-716`). Control: a fabricated arm list with one extra word
  (`nonsuch`) fails the equality, and one with `log` removed fails it too. Kills: a CLI verb added
  without a class; a table word the CLI does not have; a word in two sets; `world` or `back`
  counted as an arm. (`serve --help`'s safety is already proven by `TestEveryVerbHasHelp`, V35;
  VB4 does not depend on it.)
- **VB5** one flag set: `EpisodeReportFlags` registers exactly `markdown, json, follow, interval,
  addr` with today's defaults (`interval` = 2s) — checked by `VisitAll`; the bar admission's five
  refusals and one allow each asserted. Kills: an admission that allows `--addr`/`--follow`; a
  changed default.

**CLI behaviour unchanged (must-live, unedited):** every test in
`cmd/ailang-worldd/episode_report_test.go` (V16), `session*_test.go`, `approve*_test.go`,
`episode_*_test.go`, `preset*_test.go`, `help_test.go` and `cli_test.go` passes **without edits**
(the planner lists any edit as a conflict-surface item with its reason).

**Route (`host/daemon`)**
- **VR1** registration: `GET /surface/verb.json` behind `approvalRoute` — a foreign `Host` gets 421
  (copying `TestApprovalRoutesRefuseForeignHost`'s shape). Kills: registered on the bare mux.
- **VR2** each outcome class from §2.4's trace, end to end through `httptest`. Kills: any arm's
  branch swapped or deleted (one mutant per arm, recorded in the plan).
- **VR3** same output: seed a store with the existing report fixtures' episode, run the real CLI
  `episode report <e>` and `--json` against the loopback daemon, and compare byte-for-byte with the
  route's `text` and `json`; and the not-found case's stderr line and exit 1. Kills: the renderer
  copied instead of moved and then edited; the JSON re-marshalled with indentation.
- **VR4 [r1]** 400 for a missing `q`; **200 `{outcome:"refused", reason:"TooLong"}`** for a
  1,025-byte `q` and a normal answer at 1,024; a store error answers `refused` with no raw error
  text. Kills: the cap removed; the cap answered as 400 (an operator refusal made an HTTP error);
  the error text leaked.
- **VS1** static: `host/verbs` imports no I/O package (`os`, `os/exec`, `net`, `io/fs`);
  `surface_verb.go` does not call `OpenCredentials`, `exec.Command` or any `store` method outside the
  read set; `TestDaemonCredentialWriterSites` unchanged. Kills: a write or spawn added to the route.

**Page (`host/surfacepage`, node, non-skippable in CI via `WORLD_EXEC_NODE`)**
- **VP1** `verbOutcome` maps each `outcome` to its action. Kills: `terminal` treated as `result`.
- **VP2 [r1]** no silent fallback. For the inputs `why`, `log`, `row-189@6` and `world
  shared@0ffbb9c`, each of: HTTP 404, 400, 421, 500, a thrown fetch (network error), a 200 that is
  not JSON, a 200 without `schema`, and a 200 with an unknown `outcome` → the page calls `say` with
  the explicit message of §2.3, and **neither `addressTarget`/`openState`/`out` nor
  `location.assign` is called** (spies), and the input value is unchanged. Must-live: a 200
  schema-valid `outcome:"address"` for `row-189@6`, `world shared@0ffbb9c` and `back` reaches
  `addressTarget` with today's answers. A `TooLong` refusal shows its reason. **Named mutant:
  "fallback on any failure"** (any failure branch calling `addressTarget`), plus "fallback on 404
  only" and "accept a 200 without schema validation".
- **VP3** hostile values: a report `text` and a `reason` containing `<img src=x onerror=…>` render
  as text (`textContent`), not markup. Kills: `innerHTML`.
- **VP4** precedence: `world why` → address; `why` → verb route. Kills: the page short-circuiting on
  its own word list.
- **VP5** navigate prefix check: `href` `https://evil/…` and `/surface/model.json` are refused.
  Kills: the prefix check deleted.
- **VP6 [r1]** the shadow set at page level, with the route's real answers recorded by VR2 (not
  hand-typed): bare `log` → the refusal panel with `cli` (was navigation), `world session` →
  `addressTarget` → brane `session`, `ailang-worldd world get world-1` → refusal with its `cli`
  echo, never navigation. Kills: the page treating a refused `world …` answer as an address.

**Live Chrome (one, on the A3 driver):** **XV1** type `episode report <e>` in `#addr` on a served
daemon, press Enter, read the panel's text; equal to VR3's CLI text. Planner obligation: re-use the
A3a/A4 live-Chrome driver and skip rule exactly as `approve_page_chrome_test.go` does.

## 7. Repo-wide inventories

- **Callers of `parseInterleaved`** (to be repointed): 11 matching lines in 9 files incl. the
  definition (V7). The executor re-greps at plan time and lists them in the plan.
- **Callers of `printReport`/`execCacheLine`/`readsLine`/`printRefusals`/`sortedCounts`**: the
  executor greps `cmd/` before the move; any other caller (e.g. `episode_markdown.go`) keeps
  compiling through the moved, exported names.
- **Every place that lists `/surface` routes**: `daemon.go:1124-1146` (registrations) and
  `approvalRoutes` (`surface_approve_test.go:20`) (V22); V1 adds one entry to a new
  `verbRoutes` list rather than editing `approvalRoutes`.
- **Docs**: QUICKSTART gains a short §19 "Operator verbs in the address bar (V1: report;
  approve opens the card)" after §18 (`docs/QUICKSTART.md:1861`, V28), marked attended-pending
  like §17/§18.

## 8. Conflict surface

| item | change | why it is safe / who must agree |
|---|---|---|
| `cmd/ailang-worldd/session.go` | `parseInterleaved` moves to `host/verbs` | identical body; every caller's tests unedited |
| `cmd/ailang-worldd/episode_report.go` | flag set from `verbs.EpisodeReportFlags`; printer moved to `host/daemon` | V16's tests unedited; VR3 |
| `host/daemon/episode_report.go` | handler body extracted into `d.episodeReport` | the `/v1` route's tests unedited; `/v1` stays 13 routes (`route_table_test.go`) |
| `host/daemon/daemon.go` | one registration under `approvalRoute` | no `/v1` change; `worlddapi.ail` untouched |
| `host/surfacepage/app/app.js`, `index.html` | `go(raw)` asks the route first; placeholder, keys row | **[r1]** no fallback: only a schema-valid `outcome:"address"` navigates; every failure is an explicit refusal with the view unchanged (VP2). The 27-word shadow set is a deliberate, enumerated behaviour change (RK1). The static page (`surface.js`'s own bar) is untouched |
| row 201 | none (navigation to existing GET cards only) | no amendment (§0.2) |
| row 206 / row 222 | none in V1 | V4 consumes them |

## 9. Milestones, sprints and gates

**V1 = one sprint** (≈ 1,750 changed lines `-- cmd host docs`; cap 2,000):

| milestone | content | est. lines |
|---|---|---|
| M1 | `host/verbs`: `Split`, `Quote`, `Lookup` + table, `ParseInterleaved` (move), `EpisodeReportFlags` + bar admission; VB1–VB5 | ≈ 700 (≈ 300 code, ≈ 400 tests) |
| M2 | CLI repoint (10 call sites), `episode report` on the shared flags; printer moved to `host/daemon` | ≈ 250 (a move counts both sides) |
| M3 | `GET /surface/verb.json`, `d.episodeReport` extraction; VR1–VR4, VS1 | ≈ 500 |
| M4 | `app.js`/`index.html`; VP1–VP5; XV1; QUICKSTART §19 | ≈ 300 |

**Gates:** `./scripts/verify_ail.sh` green (no `.ail` changed, so it must be byte-for-byte the base's
result); `go build ./... && go vet ./... && go test ./...` green, with `go test -run '^$' ./...` as
the compile fence for test files (memory: `go build` does not compile `_test.go`); node tests run
(CI's `WORLD_EXEC_NODE`). **Fallback cut**, if M1–M3 measure over 1,500: M4's live Chrome XV1 and
the `approve` navigate arm move to V2's first sprint (the route answers `terminal` for `approve`
meanwhile); the plan records the cut.

**V2–V5**: one child design each (r0 → quorum → plan), in order; V3 carries the ledger ask (§0.2).

## 10. The attended exit run (Mark; NOT run by the loop, D-WORLD-109)

The row's exit criterion needs V1–V4 and Mark's own enrolled key; the loop never runs it, never
enrols a real key, never mints and never publishes. After V4 lands, Mark, at the World app only:
1. types `episode start <episode> --preset <name> --task <file>` in the address bar; the app shows
   the CLI's request block (`renderRequest`'s lines: digest, expiry, episode, worktree, grants,
   session, profile, agent); he approves with his passkey or SSH key on the card;
2. the agent launches (headless, or an interactive pane handed to herdr / row 206's tab);
3. types `episode report <episode>` (V1) and reads it;
4. types `episode finish <episode> --preset <name>`; the app shows `finishPrompt`'s lines
   (committer, repository, branch, sha, base, title, body preview, label, draft/ready, protected
   paths); he approves with his key; the PR URL is shown;
5. opens no terminal at any step.
**V1's own attended check** (optional, Mark): QUICKSTART §19 verbatim — `episode report` in the bar
equals the terminal's output, and `approve mr-…` opens the card.

## 11. Residuals and risks

- **RK1 (V1) [r1]** the verb grammar shadows bare-name navigation for every first word the CLI
  switch accepts. The complete set (V33, §2.1(b)): `why`, `episode`, `approve` (operator);
  `approver`, `session`, `service`, `serve`, `setup` (never); `commit`, `outcome`, `call`, `demo`,
  `doctor`, `head`, `health`, `help`, `hosts`, `log`, `object`, `preset`, `profile`, `provenance`,
  `registry`, `surface-index`, `tools`, `version`, `whoami` (other) — 27 words. A brane so named is
  reachable only as `world NAME` (Q2). This is a deliberate, enumerated behaviour change from the
  base; VB4's AST-derived arm list is the single source these sets must equal, so an unclassified
  new CLI verb (and any removed one) fails the build. Multi-word inputs that start with these words
  were address errors before and are CLI refusals now (§2.4, `log tail`).
- **RK2 (V1)** the route rebuilds a report on every Enter; bounded like `/v1`, no cache.
- **RK3 (V3)** P3: a finish's prompt exists only after a commit; V3 must choose intent-approval or a
  preview-bound approval (Q3), and either way a declined finish leaves a local commit, as the
  terminal does today ("commit kept at <sha>", `episode_finish.go:412`).
- **RK4 (V3)** the A4 §2.6 L1 reason text (§0.2) needs a ruling before V3 builds.
- **RK5 (V4)** the runner must run as Mark with `gh` auth (finish's `seam.Prepare` needs GitHub,
  `episode_finish.go:385-388`); a launchd daemon's environment may lack it. V4 measures.
- **RK6** the sizes for V2–V5 are parent estimates without child designs; V3 is the uncertain one.

## 12. Open questions (each with the default the loop applies)

- **Q1. Aliases (`report <e>` for `episode report <e>`, as the row's examples)?** Default: **no** —
  the CLI's words only, plus an optional leading `ailang-worldd`; one grammar, not two. Revisit with
  Mark at the exit run.
- **Q2. Verb vs a brane of the same name?** Default: **the verb wins on the first word**; `world
  NAME` always navigates (VP4).
- **Q3. (V3) One approval bound to a side-effect-free preview, or an intent approval then the
  CLI's own prompt?** Default for V3's design to evaluate first: **one approval over the exact prompt
  bytes from a preview**, with the run refusing if its prompt differs from the approved bytes.
- **Q4. (V4) Who runs an approved verb with no terminal?** Default: **not the daemon** (it spawns
  nothing today, §1.9); V4's design compares a user-level runner (a `service`-installed companion)
  with the herdr plugin (`integrations/herdr/ailang-world`), and hands interactive launches to herdr.
- **Q5. Does V1 need a ledger ask?** Default: **no** — it is read-only and within D-WORLD-107 = A.

## 13. What this is NOT doing

- Not running any mutating verb from a page in V1; not adding any POST, store write, process spawn
  or key path.
- Not changing row 201's laws, routes, cards, `approve`'s CLI or the typed `y` anywhere.
- Not putting `world-publish`, `approver enroll|revoke`, `session mint|revoke` or `service` in the
  bar, ever (R3, R4).
- Not changing the static page's bar, the `/v1` table, `worlddapi.ail` or any `.ail`.
- Not designing V2–V5 beyond their scope, size and open questions.

## 14. Verification Log

All commands ran in `/Users/voightkampff/dev/sunholo-data/ailang-world/.wt-world-iter394-row243`
at `48c92a52` on 2026-11-16 (zsh; output trimmed).

| # | claim | command | observed |
|---|---|---|---|
| V1 | base sha | `git log -1 --oneline`; `git rev-parse HEAD` | `48c92a52 Merge pull request #575 …`; `48c92a521fcbf46a0f1120609de80d656326ba5b` |
| V2 | the CLI's episode verbs and usage lines | `sed -n 165,187p cmd/ailang-worldd/episode.go`; `grep -n 'episode start\|…launch\|…finish\|…report' cmd/ailang-worldd/main.go` | `case "start"`, `"launch"`, `"finish"`, `"report"`; `main.go:141: ailang-worldd episode finish <episode> --preset <name> …`, `:144: … episode launch <episode> --preset <name> …`, `:140: … episode report <episode> [--json \| --markdown] [--follow] …` |
| V3 | top-level dispatch is a switch after global flags | `grep -n 'rest :=' cmd/ailang-worldd/main.go`; `sed -n 370,505p …/main.go` | `365: rest := globals.Args()`; `switch verb := rest[0]; verb {` with arms `version … why … approve … episode … help`, `default: unknown command %q` |
| V4 | `world` is a CLI verb | `grep -n 'usage: …world <\|  ailang-worldd .* world ' …`; `grep -n '"world":' main.go` | `main.go:101: ailang-worldd [--addr <url>] world get <ref>`; `402: case "world":` |
| V5 | the bar's parser is JS `parseAddress` | `grep -n 'parseAddress' -A50 host/surfacepage/surface.js` | `81: function parseAddress(text, index, base)`; `85: … words[0] === "world" …`; `11: ADDRESS_HELP = 'try "world", "world NAME", …'` |
| V6 | the app shell's `go(raw)` and `addressTarget` | `sed -n 300,420p …/app.js`; `sed -n 1190,1215p …/app.js`; `grep -n 'id="addr' …/index.html` | `325: function addressTarget(raw, index, parse)`, `329: parse(/^(world\|back)(\s\|$)/i.test(t) ? t : 'world ' + t, …)`; `1201: async function go(raw)`; `index.html:67: <input id="addr" … placeholder="row-189@6 · shared@0ffbb9c · back">` |
| V7 | `parseInterleaved` and FlagSet counts | `grep -n 'func parseInterleaved' -A25 …`; `grep -c 'parseInterleaved(' cmd/ailang-worldd/*.go \| grep -v ':0'`; `grep -rn 'flag.NewFlagSet' cmd/ailang-worldd/*.go \| grep -v _test \| wc -l` | `session.go:181`; `approve.go:1 approver.go:1 episode.go:1 episode_finish.go:1 episode_launch.go:1 episode_report.go:1 preset.go:1 session.go:2 session_new.go:2`; `34` |
| V8 | finish commits before its prompt | `grep -n 'Fprint\|confirm\|y/N' cmd/ailang-worldd/episode_finish.go`; `sed -n 360,420p …`; `grep -n 'func finishPrompt' -A30 …` | `379: "committed at %s\nbody: %s\n"`; `409: finishPrompt(tty, scope, …)`; `412: "opening a PR requires typed y at the terminal; commit kept at "+sha`; `episode_finish_operator.go:170: "repository: %s\nbranch: %s\nsha: %s\nbase: %s\ntitle: %s\n… Open this pull request? [y/N] "` |
| V9 | `--request` exists and refuses a launch preset | `sed -n 225,330p cmd/ailang-worldd/episode.go`; `grep -n 'func createMintRequest' -A80 …` | `211: request := fs.Bool("request", … "store a mint request instead of minting (no terminal; approve it with ailang-worldd approve <id>)")`; `249: case *preset != "" && *preset != "se-tools": return fail("--request cannot be combined with a launch preset …")`; `mint_request.go:149: fmt.Fprint(stdout, renderRequest(…))`; `:151: "nothing is minted until a person approves it at a terminal"` |
| V10 | the minter `approve --await` and the card routes | `sed -n 1,80p cmd/ailang-worldd/approve_await.go`; `grep -n '"GET /surface/approve/{id}"\|"GET /surface/decide/{ref}"' -B1 -A1 host/daemon/*.go` | `approve_await.go:3: … approve <request-id> --await …, the minter. It opens no terminal`; `daemon.go:1141: mux.HandleFunc("GET /surface/approve/{id}", d.approvalRoute(…))`; `:1145: "GET /surface/decide/{ref}"` |
| V11 | A4 L1 reserves finish scopes | `sed -n 401,437p design_docs/planned/w-world-approvals-by-key-a4.md`; `grep -rn 'reserved' host/broker/decision_law.go` | `**finish**: ParseFinishApprovalScope(scope) succeeds, or the scope starts with github-pr:`; `decision_law.go:80: if owner, reserved := ReservedApprovalScope(f.Scope); reserved {`; `:116: func ReservedApprovalScope(…)` |
| V12 | `why` is a 1,050-line client walk over unauthenticated GETs | `wc -l cmd/ailang-worldd/why.go`; `sed -n 1,60p …/why.go`; `grep -n 'func runWhy' -A40 …` | `1050 cmd/ailang-worldd/why.go`; `// The reads are unauthenticated GETs (V15)`; `349: func runWhy(addr string, args []string, stdin io.Reader, …)`, `--result`, `--scan`, `--json`, `parseWhyTarget(target, *resultFile, stdin)` |
| V13 | herdr integration lives in `integrations/`, not in Go | `grep -rln 'herdr' cmd host --include='*.go'`; `find integrations/herdr -type f` | only `cmd/ailang-worldd/hosts_awareness.go`, `hosts_hooks_test.go` (agent-state hooks, `:67`, `:158`); `integrations/herdr/ailang-world/{herdr-plugin.toml,bin/world-open,bin/world-panel,…}` |
| V14 | `ownHelpVerbs` is not a full verb list | `sed -n 693,706p cmd/ailang-worldd/main.go` | `"version", "service", "tools", "call", "why", "provenance", "setup", "doctor", "session", "approve", "approver", "episode", "profile", "surface-index", "hosts", "whoami"` (no `log`, `world`, `head`, `serve`) |
| V15 | report CLI and route bodies | `sed -n 1,104p cmd/ailang-worldd/episode_report.go`; `sed -n 127,157p …`; `sed -n 225,256p …`; `sed -n 707,760p host/daemon/episode_report.go`; `grep -n 'func writeJSON' -A8 host/daemon/*.go` | flags `markdown json follow interval(2s) addr`; `"refusing: --markdown cannot combine with --json or --follow"`; `"ailang-worldd episode report: episode %s has no committed calls\n"`; `func printReport(w io.Writer, r daemon.EpisodeReport)`; `handleEpisodeReport … buildEpisodeReport(ctx, d.intents, d.reads, episode) … writeAPIError(w, "NotFound", …)`; `handlers.go:133: _ = json.NewEncoder(w).Encode(value)` |
| V16 | report's CLI pins | `grep -n '^func Test' cmd/ailang-worldd/episode_report_test.go` | `TestClientBaseURLAcceptsABareHostPort`, `TestGlobalBareAddrIsNormalisedForEveryClientVerb`, `TestEpisodeGlobalAddrRefusedForStartOnly`, `TestEpisodeReportUsage`, `…TablePrintsEveryClass`, `…TablePrintsTheReadsLine`, `…JSONAndTableLoopback`, `…FollowReprintsOnChangeLoopback`, `…PrintsExecCache` |
| V17 | node tests are non-vacuous in CI | `sed -n 120,135p host/surfacepage/app_test.go`; `grep -rn 'WORLD_EXEC_NODE' .github/workflows/` | `TestSurfaceAppPure … t.Skip("node unavailable locally")` only without `WORLD_EXEC_NODE`; `ci.yml:250: echo "WORLD_EXEC_NODE=$(command -v node)"`, `:289: : "${WORLD_EXEC_NODE:?WORLD_EXEC_NODE is required}"` |
| V18 | episode-name pattern | `grep -rn 'workspaceEpisodePattern *=' cmd host` | `doctor.go:461: regexp.MustCompile(`^[a-z0-9][a-z0-9-]{0,63}$`)` |
| V19 | A3a guards R2–R5 | `grep -n 'func \|// R[0-9]' host/daemon/surface_approve.go`; `sed -n 56,92p …`; `grep -n 'func allowedApprovalOrigin\|func sameOriginFetch\|func jsonMediaType' -A12 host/daemon/*.go` | `approvalOrigins … surfaceOrigins(port, d.cfg.SurfaceOrigins, enrolled)`; `if !allowedApprovalHost(r.Host, o) { … StatusMisdirectedRequest`; `R3 Origin`, `R4 Sec-Fetch-Site`, `R5 jsonMediaType`, `R13 OpenCredentials … closed at once`; `sameOriginFetch(site) = site == "" \|\| site == "same-origin"` |
| V20 | comparable sprint sizes | `for m in 338ff821 9ea7b45f e6761873 842b9331 d1027477 70823879 4134ce89 b44b5f7d; do git diff --shortstat $m^1 $m -- cmd host docs; done` | A4 s1 `1982+/3-`; A4 s2 `1929+/54-`; A4 s3 `1423+/70-`; A4 s4 `1935+/12-`; A3a s2 `1991+/4-`; A3a s3 `1828+/33-`; A3a s4 `357+/36-`; A3a s1 `1467+/76-` |
| V21 | CSP | `grep -n 'surfaceCSP *=' -A3 host/daemon/*.go` | `surface_page.go:19: … connect-src 'self'; … form-action 'self'; frame-ancestors 'none'` |
| V22 | approval routes pinned by name | `grep -n 'GET /surface/requests.json' host/daemon/surface_approve_test.go` | `20: var approvalRoutes = []string{"GET /surface/approve.js", "GET /surface/approve/{id}", "GET /surface/requests.json", "GET /surface/requests/{id}", "POST /surface/approve/{id}"}` |
| V23 | id and ref shapes | `grep -rn 'MintRequestIDPattern *=' host/store/*.go`; `grep -n 'sha256:' cmd/ailang-worldd/approve.go` | `mint_requests.go:46: ^mr-[0-9a-f]{16}$`; `approve.go:99: approve <request-ref> (sha256: and 64 hex, row 201 A4)` |
| V24 | the daemon spawns nothing (negative, with control) | `grep -rln 'exec.Command' host/daemon --include='*.go' \| grep -v _test \| wc -l`; `grep -rln 'exec.Command' host --include='*.go' \| grep -v _test \| head -3` | `0`; control: `host/agentrun/agentrun.go`, `host/sshsign/sshsign.go`, `host/broker/handlers.go` |
| V25 | no `host/verbs` yet (negative, with control) | `ls host/verbs`; `ls host` | `ls: host/verbs: No such file or directory`; control lists `approval authority … surfacepage … worktree` |
| V26 | repo packages pass the dependency gate | `sed -n 895,940p host/daemon/daemon_test.go` | `TestDaemonDependencyAllowlist`: "the daemon core's full transitive build graph must consist only of the Go standard library, this repository, and the explicitly enumerated local-only modules" |
| V27 | interactive launch needs a controlling terminal | `sed -n 1,60p cmd/ailang-worldd/episode_launch_interactive.go` | `// Tests supply a private read-write file; production requires a controlling tty.`; `os.OpenFile("/dev/tty", os.O_RDWR, 0)` |
| V28 | QUICKSTART's last approval sections | `grep -n '^## 1[78]' docs/QUICKSTART.md` | `1841: ## 17. Approve from the surface *(attended — pending first verbatim run)*`; `1861: ## 18. Decide an agent's approval request *(attended — …)*` |
| V29 | the row and rulings quoted | `grep -n '^243\. ' design_docs/world-mission.md`; `sed -n 1026p …`; `sed -n 1028p …` | row 243 at `6330`; D-WORLD-107 `ANSWERED — A … never a typed y from a page`; D-WORLD-109 `ANSWERED — A … never enrols a real key, mints or publishes` |
| V30 | the route list (no verb route yet) | `grep -rn 'HandleFunc\|mux.Handle(' host/daemon/*.go \| grep -v _test \| grep -o '"[A-Z]* */[^"]*"' \| sort \| uniq` | 41 patterns incl. `"GET /surface/requests.json"`, `"POST /surface/decide/{ref}"`, `"GET /v1/episodes/{episode}/report"`; none matches `verb` |

| V31 **[r1]** | no CORS header in the daemon (negative, with control) | `grep -rniE "Access-Control\|cors" host/daemon --include='*.go' \| grep -v _test`; control: `grep -rn 'Header().Set(' host/daemon --include='*.go' \| grep -v _test \| grep -o 'Set("[A-Za-z-]*"' \| sort \| uniq -c` | 3 lines, all comments: `surface_approve.go:31: … no-store, nosniff, and never a CORS header.`; `:201: … answers no CORS, and never mints …`; `surface_approve_guard.go:60: … the daemon answers no preflight with CORS.`; control: `10 Cache-Control`, `5 Content-Security-Policy`, `15 Content-Type`, `1 Referrer-Policy`, `2 Retry-After`, `10 X-Content-Type-Options` (no `Access-Control-*`) |
| V32 **[r1]** | the page is embedded in and served by the daemon, `no-cache` | `grep -rn 'go:embed' host/surfacepage/*.go`; `grep -n 'AppScript\|handleSurfaceAppScript' host/daemon/*.go \| grep -v _test`; `grep -n 'func writeSurfaceAsset' -A10 host/daemon/*.go` | `app.go:15: //go:embed app/app.js`, `app.go:29: //go:embed app/index.html`, `script.go:9: //go:embed surface.js`; `daemon.go:1129: mux.HandleFunc("GET /surface/app.js", d.handleSurfaceAppScript)`; `surface_page.go:162: writeSurfaceAsset(w, "text/javascript; charset=utf-8", surfacepage.AppScript)`; `:190: w.Header().Set("Cache-Control", "no-cache")` |
| V33 **[r1]** | the switch's 28 arms (re-derived, not copied); what `commit`/`outcome`/`call`/`demo` do | `git log -1 --oneline` (`ad144327`, the r0 commit; `cmd/` unchanged since `48c92a52`); `sed -n 375,510p cmd/ailang-worldd/main.go \| grep -oE 'case "[a-z-]+"' \| sort \| tr '\n' ' '`; `… \| wc -l`; `grep -n '"POST /v1/commit"\|"POST /mcp/"' host/daemon/daemon.go`; `grep -n 'outcome' cmd/ailang-worldd/main.go`; `sed -n 12p cmd/ailang-worldd/demo.go` | `approve approver call commit demo doctor episode head health help hosts log object outcome preset profile provenance registry serve service session setup surface-index tools version whoami why world`; `28`; `1118: POST /v1/commit`, `1153: POST /mcp/`; `109: … outcome record --file <outcome.json\|-> [--session …]`, `164: Session credential (tools, call, commit, outcome)`; `usage: ailang-worldd [--addr URL] demo seed --db PATH --session FILE` |
| V34 **[r1]** | `world` takes exactly `get <ref>`; `parseAddress` takes one or two words | `grep -n 'func requireGet' -A8 cmd/ailang-worldd/cli.go`; `grep -n 'func runWorld' -A4 …`; `sed -n 84,86p host/surfacepage/surface.js` | `158: if len(args) != 2 \|\| args[0] != "get" {`; `166: return requireGet(args, "world", …)`; `84: if (words.length === 1 && (words[0] === "world" \|\| words[0] === "back"))`; `85: const m = words.length === 2 && words[0] === "world" ? …`; `86: if (!m) return {error: "unknown address: " + ADDRESS_HELP};` |
| V35 **[r1]** | `serve --help` returns from `genericVerbHelp` before the switch; already exercised | `grep -n 'func genericVerbHelp' -A30 cmd/ailang-worldd/main.go`; `sed -n 340,378p …`; `grep -n 'func usageLinesFor' -A15 …`; `grep -n '^  ailang-worldd serve' …`; `sed -n 8,35p cmd/ailang-worldd/help_test.go` | `374: if text, ok := genericVerbHelp(rest); ok { fmt.Fprint(stdout, text); return exitOK }` (before `switch`); `718: if ownHelpVerbs[verb] … return "", false` (serve not in it); `730: lines := usageLinesFor(verb)`; `735: if verb == "serve" { … "serve flags:" …`; `86:   ailang-worldd serve --db <path> …`; `help_test.go:11: {"serve"}, {"health"}, …`; `TestEveryVerbHasHelp … code := run(argv, …); if code != exitOK` for `serve --help` and `help serve` |

## 15. Quorum verification log

### Round 1 (quorum 2026-11-16T09:53:06Z on r0 `ad144327`): BLOCKED 3/3 — answered in r1

Artifact: `/Users/voightkampff/.ailang/state/mission-quorum-world/w-world-address-bar-verbs-2026-11-16T09-53-06Z.json`.
Reviewers oc-glm-5-3, oc-kimi-k3, gpt6-1-sol; all REJECT; none disputed the direction. The
controller adopted the fixes below; r1 applies them and nothing else.

1. **Silent fallback (gpt6-1-sol; oc-glm-5-3 on the same surface).** Objection: "§2.3 explicitly
   violates the no-silent-fallback axiom: a failed verb request is reinterpreted as navigation …
   `why` … becomes `world why` after a network error or an older daemon's 404" (gpt6-1-sol); "the
   route's 400 being swallowed … a 500 … would degrade to the address path" (glm). Fix adopted
   (gpt6-1-sol's text, verbatim): "Only a successful, schema-validated `outcome: address` invokes
   the existing addressTarget path. Network errors, non-2xx responses, and invalid responses display
   an explicit refusal and leave the current view unchanged. If legacy-daemon compatibility is
   needed, offer a separate user-selected 'Navigate as an address' action … never invoke it
   automatically." **r1:** §2.3 (the rule, the schema check, the messages); §2.4 (trace rows for
   404, 400, 421, 500, a network error and an invalid 200, each an explicit message with the view
   unchanged); §5.4; §8; VP2 (all failure classes × `why`, `log`, `row-189@6`, `world
   shared@0ffbb9c`; spies prove no navigation; named mutant "fallback on any failure"). glm (a):
   over-1,024-byte `q` → 200 `{outcome:"refused", reason:"TooLong"}`, 400 only for a missing `q`
   (§2.2 step 1, VR4). **The 404 question:** the page is `//go:embed`ded in and served by the same
   daemon with `Cache-Control: no-cache` (§1.12, V32), so V1 offers no legacy "Navigate as an
   address" action; a 404 is one more explicit refusal.
2. **Incomplete shadow set; `world get` (oc-kimi-k3).** Objection: "'Navigation never regresses' …
   is falsified by the doc's own verb table … a brane named `log`, `session`, … loses bare-name
   navigation — yet RK1 and T8 enumerate only `why`, `episode`, `approve` … `ailang-worldd world get
   <ref>` … is silently re-interpreted as navigation to brane `get`." Fix adopted (kimi's terms):
   RK1 replaced with the full enumerated set; §5.4 now reads "Navigation never regresses for inputs
   outside the enumerated shadow set (RK1) …"; step 3b `world get <ref>` → class other, refused with
   the CLI line, bare `world NAME`/`world NAME@AT` stay navigation; VB4 asserts the go/ast arm set
   equals operator ∪ never ∪ other, excluding the navigation tokens and the dropped binary name;
   VB3/§2.4/VP6 rows for bare `log`, `world session`, `ailang-worldd world get world-1`. **r1:**
   §1.13 (V33: the 28 arms re-derived), §2.1(b) steps 3/3b and two new table rows classifying
   every arm (`commit`, `outcome`, `call`, `demo` as other-writes-World with the reason; the 15
   read/tooling words as other), the exhaustiveness paragraph (1 + 3 + 5 + 19 = 28), §2.4, T8, RK1,
   §5.4, §8, VB3, VB4, VP6. The controller's arm list was re-derived independently (V33) and
   matched.
3. **Verification gaps (oc-glm-5-3 catch).** (i) CORS: §1.11 and V31 (three comment-only matches,
   positive control on `Header().Set(`); the `/v1` no-auth note cites the bare registration
   (`daemon.go:1120`) and `TestEpisodeReportJSONAndTableLoopback`. (ii) Shell characters: the list
   has **nine**; §2.1(a) and VB1 now both say nine. (iii) `serve --help`: traced by reading
   (`genericVerbHelp` returns before the switch; §1.14, V35) and already exercised by
   `TestEveryVerbHasHelp`; per glm's fix (f) VB4 nevertheless invokes no `run()` — it uses the
   go/ast arm list plus `ownHelpVerbs`/`genericVerbHelp` checks.

*(Round 2: for the controller.)*
