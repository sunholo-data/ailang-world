# w-workbench-live-and-polish — row 149

- Status: **PLANNED — iteration 244 pick-time design, revision r0.** Not yet quorum-reviewed.
  Supersedes the row-149 parts of the attended DRAFT
  [`w-world-live-surface.md`](w-world-live-surface.md) (D6, D7, §4 row 149, §6 149(a)–(e)).
  Rows 150–151 stay there.
- Base `b702d73`. Premises are in §2. Probes are banked under
  `~/.ailang/state/world-iter244-design/` (`probe/` is a scratch Go module `replace`d onto this
  worktree; `scratch/` is a `git archive HEAD` copy used for compile and test probes).
- Binding: [coding-standards.md](../coding-standards.md) S2, S3, S6, S7; `D-WORLD-64`;
  [HUMAN-SURFACE.md](../HUMAN-SURFACE.md) P1, P3, P5, §5. Consumes the row-148 interface in
  [`../implemented/w-worldd-agui-event-stream.md`](../implemented/w-worldd-agui-event-stream.md) §1, D1–D6.

## 1. Problem and clause mapping

`GET /workbench` is correct but static. It renders the ratified grammar (grades, checked links,
the closed query grammar), but a human has to reload it to see a commit. It looks like an
inspector, not a product (V1, V2). Row 148 landed a live transport, `POST /agui/`. Nothing in
the human surface consumes it.

**Row 149 makes the workbench live and presentable without moving any rule server-side
logic owns.** One hand-written, embedded, same-origin script follows `/agui/`. When an entry
arrives, it re-fetches the page's own URL and swaps named server-rendered regions. The
workbench also gains: the mockups' token system (light and dark), a newest-first live list,
a quiet-is-health footer that makes a positive assertion, a deterministic server-laid-out
SVG graph, and a decisions pane fed by the approvals chain the broker actually writes.

**Clause mapping.** Clause 5 (the human surface works). This is a **release row** under
`D-WORLD-64`: it closes no UNMET clause, and the loop takes it only under that ruling's scoped
rule-(d) exception. Clause 6 is untouched: it adds no protocol, only an AG-UI consumer.

## 2. Verification Log (first-party, at `b702d73`)

Every "the code does X" sentence below cites one of these rows. A negative search is paired
with a positive control in the same row.

| # | Premise | Command | Output (abridged, verbatim values) |
|---|---|---|---|
| V1 | The workbench CSP forbids all script and all connections | `sed -n 20p host/daemon/workbench.go` | `const workbenchCSP = "default-src 'none'; style-src 'unsafe-inline'; base-uri 'none'; form-action 'self'; frame-ancestors 'none'"`. `setWorkbenchHeaders` (`:55-61`) also sets `nosniff`, `no-store` and `no-referrer`. |
| V2 | The page has no script and no live region | `grep -rn '<script\|script-src\|connect-src' host cmd` | Only hits are hostile-input strings in `host/workbench/render_test.go` (`:160`, `:187`, `:314`, `:318`, `:811`, `:812`). This is the positive control: the pattern matches test text. No production template contains a script. |
| V3 | The CSP literal is pinned in **two** test places, not one | `grep -rn "default-src" host cmd \| grep -v '^host/daemon/workbench.go'`; `grep -rn assertWorkbenchSecurityHeaders host` | `workbench_test.go:51` (`TestWorkbenchSecurityHeaders`, comment `:46-50`: "LITERAL, never the production `workbenchCSP` symbol") and `:386` (`assertWorkbenchSecurityHeaders`). The helper is called at `workbench_test.go:197,217` and `read_deadline_test.go:384,442`. **So "every existing workbench test unchanged" (row text, draft 149(b)) cannot hold literally: the widening must edit exactly these two literal lines.** |
| V4 | Routes, and the frozen-table test counts `/v1` only | `grep -n 'HandleFunc("GET /workbench"\|HandleFunc("POST /agui/"' host/daemon/daemon.go`; `sed -n 875,876p host/daemon/daemon.go`; `sed -n 34p;51p host/daemon/route_table_test.go` | `:875 mux := http.NewServeMux()`, `:876 GET /v1/health`, `:886 GET /workbench`, `:894 POST /agui/`. The route test regex is `` mux\.HandleFunc\("(GET\|POST) (/v1/[^\"]+)" `` and requires `len(mux) != 10`. A non-`/v1` route such as `GET /workbench/live.js` is outside the frozen set. |
| V5 | The row-148 contract this row consumes | `sed -n 17-25p host/daemon/agui.go`; `sed -n 10-16p host/agui/agui.go`; `docs/QUICKSTART.md:1012-1059` | `aguiRunBudget = writeTimeout - readDeadline - aguiWriteMargin` (18 s), `aguiTick 250ms`, `aguiPage 100`, **`aguiCap 16`**, `aguiMaxBody 64<<10`. `StateSchema = "world/agui-state/v1"`. **`customNames = []string{"world.entry.committed"}`** (`:16`) is the only CUSTOM. Cursor comes from `state.lastIndex` or `Last-Event-ID`. An unknown cursor gets 404 before any byte. A run ends with `RUN_FINISHED.result.lastIndex`. |
| V6 | A browser consumes the POST SSE stream via fetch + ReadableStream under the proposed CSP; inline script and cross-origin connect stay blocked | probe `probe/main.go` + `probe/live.js` (sha256 `01e8890e…`), Chrome headless (`/Applications/Google Chrome.app`), page CSP = the exact §3 D2 literal; server = `host/agui` encoder, 1.5 s runs, entry 2 committed at 3.3 s; `probe/js-on.log` | `0.690s REPORT external-script-ran cursor=1` · `REPORT cross-origin-fetch-blocked: TypeError` · `agui POST ct="application/json" body={"threadId":"workbench","runId":"wb-0","messages":[],"state":{"schema":"world/agui-state/v1","lastIndex":1}}` · `3.302s committed entry 2` · `3.323s agui emitted entry 2` · **`3.730s REPORT run 1 status=200 entries=1 firstEntryAtMs=1111 cursor=2 timeline=entries=3`**. The entry reached the script 1.11 s into a run, before the response ended (incremental delivery). The re-fetch and swap updated the timeline. Chrome console: `Executing inline script violates … 'script-src 'self''. … The action has been blocked.` and `Connecting to 'http://127.0.0.1:47812/x' violates … "connect-src 'self'"`. The cross-origin server logged **no** request. |
| V7 | JS off: the same page makes no stream request | same probe, Chrome `--blink-settings=scriptEnabled=false`; `probe/js-off.log` | `0.645s PAGE GET /workbench` (positive control: the page loaded), then only `3.302s committed entry 2`. **Zero** `agui POST` and zero `REPORT` lines. |
| V8 | `--dump-dom --virtual-time-budget` cannot observe a streamed swap | same probe, `--virtual-time-budget=8000 --dump-dom` under `perl -e 'alarm 25'` | `rc=142` (alarm). `dump.log` shows run 0 POSTed and the server committed entry 2, but no `REPORT run 0` ever arrived. A DOM-dump oracle is unusable, so the drill's oracle is the server-side request trace (§6). |
| V9 | `script-src 'self'` admits any same-origin GET body as a script; JSON objects fail to parse; `text/plain` executes | probe gadget routes (`probe/gadget.log`, `chrome-gadget.log`): `<script src="/json-gadget">` (writeJSON-style `{"a":"fetch(…)"}`, `application/json`, no nosniff) and `<script src="/text-gadget">` (valid JS, `text/plain`, no nosniff) | JSON: console `SyntaxError: Unexpected token ':'`, and no report. **text/plain: `REPORT TEXT-PLAIN-AS-SCRIPT-RAN`, so it executed.** |
| V10 | Which worldd GET bodies could be script sources | `grep -rn 'nosniff' host --include='*.go' \| grep -v _test`; `sed -n 129-133p;151-152p host/daemon/handlers.go`; `grep -n 'writeJSON(w, http.StatusOK,' host/daemon/*.go \| grep -v _test`; `sed -n 914-916p host/daemon/daemon.go`; `sed -n 272p host/projection/projection.go`; `sed -n 949-950p host/daemon/daemon.go` | `nosniff` is set only at `agui.go:162` and `workbench.go:59`. `writeJSON` sets `application/json`. Its 8 success callers (`handlers.go:438,470,511,573,601,660,691,789`) and `writeAPIError` (`APIError{…}`) all pass struct values. Health encodes `HealthResponse` (`daemon.go:916`), and the agent card is `map[string]any` (`projection.go:272`). So every JSON body starts with `{`. **`GET /v1/head` writes `text/plain` `ref.String()` (`sha256:<hex>`) without nosniff.** Its bytes are a daemon-minted hash and never caller-controlled. See §8 for why this bounds the widening. |
| V11 | Existing tests constrain new markup | `sed -n 948p host/daemon/workbench_test.go`; `sed -n 370p`; `sed -n 210-221p host/workbench/render_test.go`; `grep -n 'Count(' host/daemon/workbench*_test.go` | (a) `provenanceWalkSection` ends at `"</section>\n</main>"`. (b) `strings.Count(body, "<h3>entry ") != WorkbenchPageLimit` (also `:710`, `:729`). (c) Every `href` must be `/workbench` or `/workbench?…`, and the page may not contain `http://`, `https://`, `href="//`, `data:` or `javascript:`. (d) Body-wide counts of `committedBy: <a`, `transitionRef: <a`, `interpreter: <a` and `stateRoot: <a`. (e) Regions are located by the exact tags `<nav aria-label="world browser">`, `<section aria-label="timeline">` and `<article aria-label="selected entry">`. |
| V12 | Those constraints, measured | `scratch/`: apply the D2 CSP, edit the two V3 literals, add `<script src="/workbench/live.js" defer></script>`, a `section[aria-label="live"]` with an inline `<svg>` (no `xmlns`) before the provenance walk, and a `<footer>` after `</main>` (`scratch-green.diff`, 55 lines); `AILANG_BIN=~/.pinned-ailang/ailang go test ./host/workbench/ ./host/daemon/ -run 'Workbench\|Render\|ReadDeadline'` | Baseline: `ok` both, 42 `--- PASS` lines in daemon. **With the edits: `ok host/workbench`, `ok host/daemon`.** Mutant m1 (`<svg xmlns="http://www.w3.org/2000/svg"`): `--- FAIL: TestRenderEmitsOnlyLocalLinks`. Mutant m4 (`<h3>entry 7</h3>` in the live section): `--- FAIL: TestWorkbenchTimelineBound`, `--- FAIL: TestWorkbenchTimelinePaging`. Mutant m2 (live section moved after the walk): still `ok` with this content, so the placement rule in D4 is a precaution, not a measured breaker. |
| V13 | The field census only covers listed types | `sed -n 405-409p host/workbench/render_test.go`; scratch: add `Live LiveView` to `Page` | The list holds the 9 current view types. Unrendered `Page.Live` gives `--- FAIL … unrendered fields: [Page.Live]`. With `{{.Live.Cursor}}` rendered, **the unrendered `LiveView.Unrendered` field is NOT caught (`ok`)**, because `LiveView` is not in `workbenchViewTypes`. New view types must be appended there. |
| V14 | Grade markup is pinned byte-exact | `sed -n 124p host/workbench/render_test.go`; `sed -n 1003p host/daemon/workbench_test.go` | `` `<p><span>PROVEN</span></p>` `` (also the `TESTED`+verdict forms `:88`, `:107`). The daemon object page must contain none of `<span>PROVEN</span>` … `<span>CLAIMED</span>` (it always renders `GRADE UNAVAILABLE — `+`objectGradeUnavailableReason`, `workbench.go:42`). |
| V15 | No decision-packet producer exists | `grep -rln 'decision-packet\|DecisionPacket' host cmd world packages \| grep -v _test`; control `grep -rln 'epoch-registry/v1' host cmd \| grep -v _test` | Only `host/broker/testdata/metadata_world_core_0.1.0.json`, `world/types.ail`, `packages/world-core/world/types.ail` (type `:159`). Control: 7 production files. **No host path writes a `DecisionPacket`.** Row 148 emits no decision event (V5). |
| V16 | The decision facts the host does write: the approvals chain | `sed -n 16-22p;80-84p;100-120p;225-250p host/broker/approve.go`; `sed -n 133-139p host/broker/publish_op.go` | `Human.Approve` `PutObject`s a `world/approval-request/v1` object `{effect,scope,cost,requester,now}`. `DecideApproval` writes `world/approval-decision/v1` `{requestRef,decision,decidedBy,now}`. Each one prepends a `world/approvals/v1` head object `{previousHead,requestRef,decisionRef}` and moves registry head `world/approvals/v1`. `walkApprovalHead` is unexported and linear ("O(all approval-head objects)"). The producer is the publish path (`publish_op.go:135-138`). **These are registry and object writes, not log entries, so `/agui/` never streams them.** |
| V17 | There is no newest-first or "latest index" log read | `grep -rn 'entry_index DESC\|MAX(entry_index)' host --include='*.go'`; control `grep -rn 'ORDER BY entry_index' host --include='*.go' \| grep -v _test` | Only `reference_index_test.go:55` (an index DDL string in a test). Control: 6 ascending reads (`scan.go:67`, `log_after.go:23`, `object_commits.go:20`, `object_references.go:70-72`). `World` carries `LogHead` (a hash) and no index (`workbench.go:394-400`). **The workbench cannot learn the latest index without a new read.** |
| V18 | Adding a `Store` method to `readStore` breaks no test fake | scratch: add `LogEntriesLatest` to `*store.Store` and to `readStore` (`daemon.go:484-495`), then `go vet ./host/daemon/ ./host/store/ ./cmd/...`; control: add `NoSuchMethod()` | `vet rc=0`. Control: `agui_test.go:29:19: … *store.Store does not implement readStore (missing method NoSuchMethod)`. Every fake embeds `readStore` (`grep -rn 'func (.*) LogEntriesAfter'`: 3 test overrides, all on embedding types). |
| V19 | `host/workbench` is pure | `go list -deps ./host/workbench \| grep -E '^net/http$\|host/store\|^html/template$'` | Only `html/template` (the positive control). Neither `net/http` nor `host/store`. |
| V20 | `go:embed` precedent; deps | `grep -rn 'go:embed' host cmd --include='*.go' \| grep -v _test`; `sed -n 5-8p go.mod` | `host/store/store.go:41://go:embed schema.sql`. Direct requires: `ailang v0.47.2` and `modernc.org/sqlite v1.54.0`. This design adds none. |
| V21 | Mockup tokens | `sed -n 8-12p design_docs/mockups/approval-inbox.html`; `sed -n 8-12p design_docs/mockups/grounded-prose.html` | Both use `:root` plus `@media (prefers-color-scheme: dark)`. approval-inbox has 10 tokens (`--bg --card --ink --mut --line --acc --accbg --ok --warn --warnbg`). grounded-prose has the same plus `--okbg` (11). Light `--bg:#faf9f5`, dark `--bg:#1f1e1b`. |
| V22 | P5 demands a positive quiet state | `sed -n 102-108p design_docs/HUMAN-SURFACE.md` | "the quiet state MUST be a positive assertion with a timestamp — last verified heartbeat, next expected, what was checked — never the absence of items." |
| V23 | The workbench has no usage doc (S7 gap) | `grep -n -i workbench docs/QUICKSTART.md`; control `grep -c 'Watch the world live' docs/QUICKSTART.md` | One hit, `:320` "**Later:** the approval-inbox workbench (item 7)". Control: `1`. Also, numbering collides: `### 10. Build and test non-AILANG projects` (`:640`) and `## 10. Watch the world live` (`:1012`). |
| V24 | The closed grammar needs no change | `sed -n 45-53p;76-89p host/daemon/workbench.go`; `sed -n 226-227p host/daemon/workbench_test.go` | 7 accepted keys. The truth table enumerates `{"entry","from","object","payload","world"}`. V6 shows a same-URL re-fetch is enough for the refresh, so this design adds no key and no fragment endpoint. |
| V25 | The error page is script-free | `sed -n 63-65p host/daemon/workbench.go` | `workbenchErrorTemplate` is a fixed `<h1>{{.Class}}</h1><p>{{.Message}}</p>` page. Only `handleWorkbench`'s success path calls `workbench.Render`. |
| V26 | The workbench has no evidence read | `grep -n 'evidence\|Evidence' host/daemon/workbench.go host/workbench/render.go` | One hit, `workbench.go:42`, the *reason* string "does not resolve subject-bound evidence into an object grade". The draft D7's "transition → object → **evidence**" graph layer has no source. |
| V27 | Test tooling on the rigs | `node --version`; `sed -n 193p;216p .github/workflows/ci.yml`; Chrome path | Local `v26.0.0`. CI go-verify job: "node is the runner's own" and `echo "WORLD_EXEC_NODE=$(command -v node)" >> "$GITHUB_ENV"`. Chrome on CI is **UNVERIFIED** (not measured), so no CI gate depends on it. Locally, Chrome headless hangs inside this Bash sandbox (`about:blank --dump-dom` > 120 s) and runs unsandboxed (V6–V9). |
| V28 | Test context guard | `sed -n 10-15p host/daemon/m7b_deadline_test.go` | `boundedTestContext` = `context.WithTimeout(…, 30*time.Second)`, a stuck guard. |

**Draft claims false or stale at `b702d73`** (all in [`w-world-live-surface.md`](w-world-live-surface.md)):

1. 149(b) "with JS off, every existing workbench test passes unchanged". **False as worded.**
   Two literal CSP lines must change (V3). Refined in AC0.
2. §2.3 "`TestWorkbenchSecurityHeaders` pins this constant". **Incomplete.** A second literal in
   `assertWorkbenchSecurityHeaders` guards 4 more call sites (V3).
3. D7 / §4 "decision-packet stream". **No source.** No host path commits a `DecisionPacket`, and
   row 148 streams no decision event (V15, V5). Revised to the approvals chain (D7).
4. D7 "transition → object → **evidence** neighbourhood". **No evidence read exists** in the
   workbench (V26). Revised to entry → object (D6).
5. §4 "the fragment endpoints in the closed grammar". **Unneeded.** A same-URL re-fetch is
   enough (V6, V24), so the grammar stays byte-unchanged.
6. D6 "served … at a hashed path". **Revised** to a fixed path with `no-store` (D2). A hash only
   busts caches the route never fills.
7. §9 "CSP widening … limited to `'self'`". **Weaker than it reads.** `'self'` executes any
   same-origin `text/plain` body (V9). The bound is the content of the daemon's GET bodies
   (V10, §8), not the keyword.
8. §2.1 "routes at `daemon.go:800-821`". **Stale.** `mux :=` is now `:875` and the first route `:876`,
   with `/workbench` at `:886` (V4).
9. §6 149(c) "headless Chrome check". **Feasible only as a local drill.** A DOM dump cannot
   observe the swap (V8), and CI Chrome is unmeasured (V27). The CI weight moves to binary-free
   and node-driven tests (§6).

## 3. Decisions, with alternatives

### D1 — Live = follow `/agui/` with fetch, then re-fetch the page's own URL and swap named regions

The script is `host/workbench/live.js` (embedded, about 150 lines, no framework):

1. Read the cursor from `section[aria-label="live"]`'s `data-live-cursor`. This is the latest
   committed index at render time, or `-1` (D3).
2. `fetch("/agui/", {method:"POST"})` with body
   `{"threadId":"workbench","runId":"wb-<n>","messages":[],"state":{"schema":"world/agui-state/v1","lastIndex":<cursor>}}`
   (V5 and V6: accepted unchanged). Read `response.body.getReader()`. Split on `"\n\n"`, keep the
   unterminated tail, and parse `data: ` lines. This is row 148's V18 algorithm.
3. Track the cursor from `STATE_DELTA` (`replace` ops only; any other op throws) and from
   `RUN_FINISHED.result.lastIndex`. The cursor advances only on a **complete** frame.
4. On the first `CUSTOM world.entry.committed` of a burst, schedule a refresh 300 ms out
   (debounce; one refresh in flight). The refresh does `fetch(location.href)` and then
   `DOMParser`. For each selector in one fixed list, it calls
   `old.replaceWith(document.importNode(fresh, true))`. The list is
   `nav[aria-label="world browser"]`, `section[aria-label="timeline"]`,
   `section[aria-label="live"]`, `section[aria-label="world graph"]`,
   `section[aria-label="decisions"]` and `footer[aria-label="live status"]`. Newly arrived
   live rows get a CSS class `fresh` (animated in the stylesheet).
5. On `RUN_FINISHED`, re-POST at once with the new cursor. Never assume one run drains the
   log: a run is bounded (V5) and the script simply follows `lastIndex`.

Every pixel still comes from `html/template`. Grades, checked links and the closed grammar
stay server-side, and their tests keep guarding them. The script never builds markup from
stream data. It moves nodes that our own escaped template produced.

*Alternatives (rejected):* (a) **Fragment endpoints** (draft §4) would add grammar keys or
routes and duplicate render paths for no gain (V24). (b) **Client-side rendering from the
`CUSTOM` value** re-implements the edge and grade rules in JS, outside every test. (c) **A
`GET` SSE route for `EventSource`** is a second wire dialect next to AG-UI (§3.7), and row 148
D3 chose POST only. (d) **`<meta http-equiv=refresh>`** polls the whole page with no event,
which breaks P5 (it interrupts and redraws when nothing changed). (e) **htmx/Datastar** are
vendored behaviour outside the gates. The draft's own threshold (around 300 lines) is not
reached.

### D2 — One embedded script at a fixed path; the CSP widens by exactly two directives

- `//go:embed live.js` in `host/workbench` (`var LiveScript []byte`; package stays pure, V19).
- Route `GET /workbench/live.js` → `d.handleWorkbenchScript`. It sets
  `Content-Type: text/javascript; charset=utf-8`, `X-Content-Type-Options: nosniff`,
  `Cache-Control: no-store`, `Content-Security-Policy: default-src 'none'`, and writes
  `LiveScript`. Other methods get a mux 405. It is not protected:
  `isProtected(GET /workbench/live.js) == isProtected(GET /workbench)`, pinned by a test.
- The page head gains exactly `<script src="/workbench/live.js" defer></script>`. There is no
  other `<script` anywhere: no inline script, no `on*=` handler, no `eval`, no CDN. The error
  page stays script-free (V25).
- **New CSP literal** (the order measured in V6):
  `default-src 'none'; style-src 'unsafe-inline'; script-src 'self'; connect-src 'self'; base-uri 'none'; form-action 'self'; frame-ancestors 'none'`.
  The two V3 literal lines are updated to this string. They stay literals, never the symbol.

*Alternatives:* a hashed path (draft D6) costs a computed route to bust a cache the route never
fills. A nonce cannot be used: a per-request nonce makes the page non-deterministic and needs a
script-bearing template. A `'sha256-…'` hash-source with SRI would be narrower than `'self'`,
but `D-WORLD-64` and the row fix the widening as exactly `script-src 'self'; connect-src 'self'`.
It is recorded as a hardening candidate (§9).

### D3 — The cursor and the live list come from one new newest-first store read

Add `Store.LogEntriesLatest(ctx, limit int) ([]LogEntry, error)`:
`SELECT … FROM log_entries ORDER BY entry_index DESC LIMIT ?`, `1 ≤ limit ≤ MaxLogEntryPage`,
deadline required, quarantine refused. It reuses `scanLogAfterRow` and joins `readStore`
(V18). `handleWorkbench` calls it once with `limit = 10`. The result gives `Live.Recent`
(newest first) and `Live.Cursor = Recent[0].EntryIndex`, or `-1` on an empty log. It gives
`Live.Head`, the entry hash, or empty. Every render pays one extra indexed statement.

- **Why a store read:** V17. Nothing yields the latest index. A dense `GetLogEntry` scan
  stalls at gaps (row 148 V10/V26). Paging `LogEntriesAfter` to the end is O(log length).
- **Why the cursor must be an existing index:** `/agui/` returns 404 for an unknown cursor
  (V5). The max existing index is always valid. A commit landing between render and the first
  POST is delivered by the stream, never lost.
- **Why a separate live list:** the timeline is the oldest-first page `from=0` (100 rows,
  V11b). A commit at index ≥ 100 never appears there. The live list shows the tail.

*Alternative:* resume from `-1` and let the script skip the backlog. That streams the whole
log on every page load and holds a slot for up to 18 s per 100-entry page. Rejected.

### D4 — Polish within measured markup constraints

- **Tokens:** the union of the two mockups' sets (11 tokens, V21) goes on `:root`, with the
  dark values under `@media (prefers-color-scheme: dark)`. Cards, spacing and type scale follow
  the mockups. It stays one inline `<style>` (allowed by the unchanged
  `style-src 'unsafe-inline'`; no `<link href>`, V11c).
- **New regions:** `section[aria-label="live"]` (`aria-live="polite"`), `section[aria-label="world graph"]`
  and `section[aria-label="decisions"]` go inside `<main>` **before** the provenance walk. The
  walk stays the last section, so `"</section>\n</main>"` still closes it (V11a; precaution,
  V12 m2). `footer[aria-label="live status"]` goes after `</main>`.
- **Forbidden in new markup** (each measured red, or pinned by an existing count): `<h3>entry `
  (V12 m4); `xmlns="http://…"` or any absolute URL (V12 m1); the strings `committedBy: <a`,
  `transitionRef: <a`, `interpreter: <a`, `stateRoot: <a` (V11d); any `href` not built by
  `workbenchHref`.
- **Grades:** the grade markup is byte-unchanged (V14). Per-grade colour pills are a non-goal:
  the daemon never renders an available object grade (V14), and P3 forbids colour alone.
- **Census:** every new view type is appended to `workbenchViewTypes` (V13).

### D5 — Quiet is health: the footer is a positive assertion

The server renders `footer[aria-label="live status"]` with `role="status"`. It reads "Log at
entry N (head `<hash>`) when this page was rendered", or "The log is empty", plus a
`<span data-live-status>Live updates: off. Reload to refresh.</span>`. There is no server
wall clock, so the render stays deterministic. With JS on, the script rewrites only that span
after each `RUN_FINISHED`: "Live: checked HH:MM:SS, no new entries after entry N; next check
within 18 s". After a burst it reads "K new entries since HH:MM:SS". On failure it reads
"paused: <reason>, retrying in Ns" (D8). That is P5's "last verified heartbeat, next expected,
what was checked" (V22). A dead stream is never rendered as calm.

### D6 — The world graph: a deterministic bipartite SVG of entries → objects, from checked edges only

`workbench.LayoutGraph(entries []GraphEntry) GraphView` is pure. Each `GraphEntry` is
`{Index, Href, Edges []EdgeView}`, where the edges are the entry's already-checked
`transitionFn`/`interpreter`/`transitionRef` edges.

- **Where:** the home page uses the newest 5 entries of `Live.Recent`, each with
  `entryEdges` (≤ 15 extra `GetObject` reads). The selected-entry page uses `[Selected]` with
  the edges it already has. The object page renders "graph: select an entry" (a non-goal for
  now, §10).
- **Layout:** column 0 holds entries, newest first. Column 1 holds distinct object targets in
  first-appearance order (entry order, then relation order). Fixed geometry: box 280×28,
  row pitch 36, margin 16, column x 16 and 344, `viewBox="0 0 640 H"` with
  `H = 32 + 36·max(rows)`. Edges are straight lines from the entry's right-middle to the
  object's left-middle, labelled by relation in `<title>`. There is no map iteration, no float
  and no clock.
- **Links:** an available target is `<a href="{{workbenchHref .Href}}">`. An unstored target is
  a dashed box labelled `UNAVAILABLE` with **no** anchor (the checked-link rule). The `<svg>`
  carries no `xmlns` (V12 m1). Nodes carry no grade (V26, P3).
- **Golden:** `host/workbench/testdata/graph_home.golden` and `graph_selected.golden` are
  **generated** by `go test ./host/workbench -run TestGraphGolden -update-workbench-golden`
  (flag pattern: `host/broker/exec_settings_test.go:20`) and checked in as produced.

*Alternative:* a client-side force layout (d3-force) is non-deterministic, adds vendored JS
and breaks row 150's byte-identical captures (draft D7, re-affirmed).

### D7 — The decisions pane renders the approvals chain, the only decision facts the host writes

Add `broker.RecentApprovals(ctx, r ApprovalReader, maxHeads int) (ApprovalPage, error)`, where
`ApprovalReader = interface{ GetRegistryHead; GetObject }` (`readStore` satisfies it). It walks
at most `maxHeads = 40` head objects newest-first from `world/approvals/v1` (V16). A pure fold
`FoldApprovals(heads)` groups by `requestRef`: the newest decision head wins, and a request
with no decision is `pending`. It returns ≤ 20 summaries
`{RequestRef, Effect, Scope, Requester, Cost, Status, DecidedBy}` and a `Truncated` flag.
The pane shows each request as a checked `?object=` link, the effect/scope/requester, and the
status text. It never offers an action: `/workbench` stays GET-only, and a decision is a
transition through the authenticated path (HUMAN-SURFACE §5).

- No registry head means the pane says "No approval requests recorded", a positive empty state.
- A store I/O error goes through `writeWorkbenchStoreError` (5xx), like every workbench read.
- A malformed chain (wrong semantic id, undecodable payload) goes to
  `ErrApprovalChainMalformed`. The pane renders `UNAVAILABLE: approvals chain is malformed at
  <ref>` and the rest of the page still renders.
- The pane is **not** streamed: approvals are not log entries (V16). It refreshes whenever a
  log entry triggers the D1 swap, and shows "as of entry N". Streaming it needs a log-derivable
  source (§9).

*Alternatives:* render `world/decision-packet/v1` objects. There are none (V15), so that would
be an empty pane presented as a feature. A fixture-fed pane would be fabrication on the live
page. Both rejected.

### D8 — Resource posture: one stream slot per visible tab, never silent

`/agui/` admits 16 concurrent runs globally (V5). The script holds at most one, and only
while `document.visibilityState === "visible"`. On `visibilitychange` to hidden it aborts the
fetch (`AbortController`) and resumes on show. On 503 `StreamLimit` it waits `Retry-After`
seconds, then backs off exponentially (×2, cap 30 s). On a network error, a non-200, or
`RUN_ERROR`, it does the same. Every wait is shown in the D5 span. The page itself never
depends on the script.

### D9 — No `.ail` in this row (S1/S2/S3 answer)

All of row 149 is host rendering: an HTTP route, response headers, an embedded asset, HTML/SVG
templates, and a browser script. **Why not a package (S3):** an AILANG package cannot own a
route, set a CSP header, embed a file, or run in a browser. The semantics it shows (entries,
objects, approvals) already live in the kernel log and the broker. **Why no Z3 contract
(S1):** the only pure invariants are the layout geometry (in bounds, disjoint boxes) and the
approval fold (newest decision wins). Both are closed-form, and Go property tests check them
directly over the real inputs. A transcribed `.ail` twin would be a predicate nobody's code
runs, which S1 itself calls decoration. No AILANG syntax is claimed in this doc, so S5 is not
triggered.

## 4. Design sketch (files)

| Layer | File | Kind |
|---|---|---|
| Pure | `host/workbench/render.go` (tokens, regions, `LiveView`, `GraphView`, `DecisionsView`), `host/workbench/graph.go` (`LayoutGraph`), `host/workbench/live.js` + `embed.go` | Go/HTML/JS, no I/O |
| Pure | `host/broker/approvals_view.go` (`FoldApprovals`) | Go, no I/O |
| Effect | `host/store/log_latest.go` (`LogEntriesLatest`), `host/broker/approvals_view.go` (`RecentApprovals` walk), `host/daemon/workbench.go` (populate views, script route), `host/daemon/daemon.go` (one `HandleFunc`, the `readStore` method, route doc comment) | Go |
| Docs | `docs/QUICKSTART.md` new section "Open the live workbench" | S7 |
| CI | `.github/workflows/ci.yml`: one verbose PASS-loop step for the node-driven script test (M3) | gate |

## 5. Milestones (tests first; each ≤ ~250 production LOC; JS counts)

Executor packing: M1+M2 (one run), M3 (one run), M4 (one run), M5+M6 (one run).

**AC0 (every milestone; the refined 149(b)).** In `git diff -U0 b702d73 -- host/daemon/*_test.go host/workbench/*_test.go`,
the only changed lines inside **pre-existing** test functions and helpers are the two V3 CSP
literals and appended entries in `workbenchViewTypes`. New tests live in new functions or new
files. The evaluator checks this by hunk listing.

**M1 — Tokens, live/footer regions, view types (pure, `host/workbench`; ~120 LOC).**
- AC1.1 `TestTokensBothThemes`: parse the `<style>` block. `:root` and the dark `:root` declare
  the same 11 token names (V21). Null control: the test fails if it finds 0 tokens in either.
- AC1.2 `TestRenderLiveRegion`: `Live{Cursor:7, Recent:[7,6]}` renders `data-live-cursor="7"`,
  rows newest first, and no `<h3>entry `. `Cursor:-1` renders `data-live-cursor="-1"` and "The
  log is empty".
- AC1.3 `TestRenderQuietFooter`: both states render a positive sentence (index plus head, or
  empty log) and `<span data-live-status>Live updates: off`. The footer is outside `<main>`.
- AC1.4 `TestProvenanceWalkStaysLast`: the rendered page contains exactly one
  `"</section>\n</main>"`, and it closes the provenance walk.
- AC1.5 `workbenchViewTypes` gains `LiveView`. `TestWorkbenchViewFieldsAllRender` passes.
- AC1.6 All pre-existing `host/workbench` tests pass, with AC0.

**M2 — `LogEntriesLatest` and daemon wiring (effectful; ~90 LOC).**
- AC2.1 `TestLogEntriesLatestNewestFirst` (entries 0, 1, 5): limit 10 gives [5,1,0], limit 1
  gives [5], an empty store gives []. Each row deep-equals `GetLogEntry`.
- AC2.2 Limit 0 or 501 gives `InvalidLimitError`. No deadline gives `ErrNoDeadline`. A
  quarantined store is refused.
- AC2.3 `TestWorkbenchLiveCursorCrossesGap`: a daemon with commits at 0, 1, 5 (REST, as row 148
  AC3.14) renders `data-live-cursor="5"`. An empty store renders `-1`.
- AC2.4 `TestWorkbenchLiveCursorIsResumable`: POST `/agui/` with the rendered cursor as
  `state.lastIndex` gets 200, not 404. It uses `aguiBudget` 300 ms (stimulus) and a 10 s client
  context (stuck guard).
- AC2.5 `TestWorkbenchLiveStoreError`: a seam failing `LogEntriesLatest` gives the existing
  constant 500 HTML, with no store text.

**M3 — Script, route, CSP (effectful; ~200 LOC incl. ~150 JS).**
- AC3.1 The two V3 literals become the D2 literal. `TestWorkbenchSecurityHeaders` and
  `assertWorkbenchSecurityHeaders` callers pass.
- AC3.2 `TestWorkbenchLiveScriptTag`: the success page contains exactly one `<script` and it is
  `<script src="/workbench/live.js" defer></script>`. No `on[a-z]+=` attribute appears.
- AC3.3 `TestWorkbenchErrorPageInert`: 400/404/500 workbench pages contain no `<script`.
- AC3.4 `TestWorkbenchLiveScriptRoute`: GET returns 200, the D2 headers (literal values), and a
  body byte-equal to `workbench.LiveScript` with `len > 0`. POST gets 405.
  `isProtected` is equal for `GET /workbench/live.js` and `GET /workbench`.
- AC3.5 `TestLiveScriptContract` (binary-free): from the embedded script, extract the
  `STATE_SCHEMA` literal and assert it equals `agui.StateSchema`. Extract the one-line
  `REGIONS` JSON array and assert every selector's `aria-label` occurs exactly once in a
  rendered daemon page. Extract the POST path and assert it equals `/agui/`. Null control: the
  test fails if any extraction finds nothing.
- AC3.6 `TestLiveScriptPure` (node): runs `node` (from `WORLD_EXEC_NODE` or `PATH`) on the
  script plus a harness (the script exports `splitFrames`/`applyDelta`/`nextBackoff` only when
  `module` exists, and starts only when `document` exists). It feeds
  `host/agui/testdata/stream_fixture.golden` split at **every** byte offset into two chunks.
  The event list equals the whole-buffer parse, and the cursor never exceeds the last fully
  delivered entry. `applyDelta` with an `add` op throws. Backoff is 1, 2, 4 … capped at 30.
  Locally the test SKIPs without node. In CI a new verbose step runs it and fails unless
  `--- PASS: TestLiveScriptPure` appears (the ci.yml PASS-loop pattern).

**M4 — The SVG graph (pure layout + small wiring; ~200 LOC).**
- AC4.1 `TestGraphGolden`: the two fixtures render byte-equal to their goldens, which are
  generated by the `-update-workbench-golden` flag. **Anti-vacuity:** it fails if a golden is
  0 bytes, lacks `<svg`, or has a `<rect` count ≠ the fixture's node count.
- AC4.2 `TestGraphDeterministic`: two layouts and renders are byte-equal, and so is a layout of
  a deep copy.
- AC4.3 `TestGraphNodesDisjointAndInBounds`: for fixtures of 1–5 entries with shared and
  distinct targets, every box lies inside the viewBox and no two boxes intersect.
- AC4.4 `TestGraphUnavailableNodeHasNoLink`: an unstored target renders `UNAVAILABLE`, and the
  graph's `<a` count equals the count of available nodes.
- AC4.5 `TestGraphCarriesNoGrade`: the graph contains none of `PROVEN|TESTED|ATTESTED|CLAIMED`.
- AC4.6 `TestWorkbenchGraphWiring` (daemon): the home page graph shows the newest ≤ 5 entries.
  The selected-entry page graph shows exactly that entry's 3 edges. Pre-existing
  `TestRenderEmitsOnlyLocalLinks` and the body-wide count tests pass (V11).

**M5 — Decisions pane (effectful walk + pure fold; ~180 LOC).**
- AC5.1 `TestFoldApprovals`: request A, request B, decide A deny, request C, decide C approve
  gives the newest-first order [C approve, A deny, B pending]. A request decided twice shows the
  newest decision.
- AC5.2 `TestRecentApprovalsBounded`: a 200-head chain does ≤ 40 `GetObject` reads (counting
  fake) and sets `Truncated`.
- AC5.3 `TestRecentApprovalsMalformed`: a wrong-semantic-id head gives
  `ErrApprovalChainMalformed`. The workbench renders the UNAVAILABLE line with status 200.
- AC5.4 `TestWorkbenchDecisionsPane`: a real `HumanHandler` request plus `DecideApproval` on a
  daemon store renders a checked request link and "deny by <decidedBy>". With no head it renders
  "No approval requests recorded".
- AC5.5 `TestWorkbenchDecisionsReadOnly`: the pane contains no `<form`, no `<button` and no
  `method=`.

**M6 — S7 usage surface and the Chrome drill (effectful; docs + one env-gated test).**
- AC6.1 QUICKSTART gains "Open the live workbench": open `http://127.0.0.1:7644/workbench`,
  what the live list, footer, graph and decisions pane mean, the JS-off behaviour, the 16-slot
  posture (D8), and a commit appearing without reload. It is executed verbatim, with the
  transcript in the PR. It uses the next free section number (V23 collision noted, not fixed
  here). The `Handler` route doc comment names `GET /workbench/live.js`.
- AC6.2 `TestWorkbenchLiveChromeDrill` (runs only with `WORLD_CHROME=<path>`; the executor runs
  it and pastes the `-v` transcript; it is **not** a CI gate, V27). Setup:
  `httptest.NewServer` over `d.Handler()`, wrapped by a recorder that tees `POST /agui/` bodies
  and counts `GET /workbench`; `d.aguiBudget = 1s`; Chrome `--headless=new
  --user-data-dir=t.TempDir()`. Oracle: a POST with `lastIndex == k`, then (after the
  test's `d.store.Commit` of entry k+1) a `GET /workbench` re-fetch, **and** a later POST with
  `lastIndex == k+1`. The JS-off arm (`--blink-settings=scriptEnabled=false`) sees the page GET
  and zero POSTs in a 3 s window. This arm is an instrument control, not load-bearing.

Gate (every milestone): `go vet ./... && go build ./... && go test ./... && go test -race ./host/daemon ./host/workbench ./host/store ./host/broker`
with `AILANG_BIN=~/.pinned-ailang/ailang`.

## 6. Test timing (CI runs `-race` on small loaded linux runners)

| Test | Wall-clock bound | Kind |
|---|---|---|
| M1, M4, AC2.1–2.3, AC2.5, M3 (except 3.6), M5 | none (pure or synchronous `httptest.ResponseRecorder`); store ctx = `boundedTestContext` 30 s | stuck guard (V28) |
| AC2.4 cursor resumable | `aguiBudget` 300 ms; client ctx 10 s | 300 ms = **stimulus** (ends the run); 10 s = stuck guard. It asserts status and frames, never elapsed time |
| AC3.6 node | `exec.CommandContext` 30 s | stuck guard |
| AC6.2 Chrome drill | `aguiBudget` 1 s; commit after the first POST is observed; each wait ≤ 20 s; JS-off window 3 s | 1 s and the commit = **stimulus**; 20 s = stuck guard; the 3 s window is a labelled non-load-bearing control |

No test asserts that one bounded run delivers the whole log. AC2.4 and AC6.2 follow
`RUN_FINISHED.lastIndex` and the script's re-POST.

## 7. Failure modes

| Failure | Behaviour |
|---|---|
| JS disabled or blocked | Today's page, unchanged and complete. The footer says "Live updates: off. Reload to refresh" (V7, AC0) |
| `/agui/` at its cap (503) | Back off per `Retry-After`, ×2 to 30 s. The footer shows "paused: stream limit" (D8) |
| Daemon restarts or the network drops | The fetch rejects. Back off and resume from the same cursor (still valid: the log is append-only) |
| `RUN_ERROR` / non-200 | Same as above. The status span names the class |
| Connection cut mid-frame | The unterminated tail is discarded. The cursor advanced only on complete frames (AC3.6), and the deltas are idempotent (row 148 AC1.5) |
| Re-fetch returns 5xx | Regions are kept and the status span says "refresh failed (503 Timeout)". The next entry retries |
| Hidden tab | The stream is aborted and the slot released. It resumes on show |
| Log index gap | The cursor and stream cross it (AC2.3, row 148 AC3.14) |
| Malformed approvals chain | Pane UNAVAILABLE with a named reason. The page is 200 (AC5.3) |
| Graph target not stored | A dashed UNAVAILABLE node with no link (AC4.4) |

## 8. Conflict surface

- **The CSP widening loosens a pinned security header.** Bound: (1) the literal is pinned in
  both V3 places, so any further token (`'unsafe-inline'`, a host, `*`) is red (M-a, M-b).
  (2) Inline execution stays blocked, and so does cross-origin connect (V6, measured).
  (3) `'self'` can only load a same-origin GET body. Those bodies are the HTML page (nosniff,
  refused as script), `writeJSON` objects (a `SyntaxError` as script, V9), `GET /v1/head`'s
  daemon-minted `sha256:<hex>` (`text/plain`, executable in principle (V9) but never
  caller-controlled, V10), the agent card JSON, and `live.js`. (4) Getting a `<script>` tag
  into the page at all needs a template-escaping failure, which
  `TestRenderEscapesAllObjectText` and `TestWorkbenchLiveScriptTag` (exactly one `<script`)
  guard. Residual: a future GET route that echoes caller bytes as `text/plain` without nosniff
  would become a script gadget. That is filed as a finding (§9).
- **Edited:** `host/daemon/workbench.go`, `host/daemon/daemon.go` (one `HandleFunc`, the
  `readStore` method, the doc comment), `host/workbench/render.go`, `host/broker/` (new file),
  `host/store/` (new file), `docs/QUICKSTART.md`, `.github/workflows/ci.yml` (one step), and
  the two V3 literals.
- **Must stay unchanged** (no hunk; the evaluator checks `git diff --stat`): the ten `/v1`
  patterns and handlers, `isProtected`'s body, `acceptedWorkbenchKeys` and
  `supportedWorkbenchQuery`, `host/agui/*`, `host/daemon/agui.go`, every D7 constant,
  `workbenchErrorTemplate`, `tools/launchd/*`. No open critical-path row names these files at
  `b702d73`. Re-check at merge.

## 9. Findings to file (not fixed here)

1. `GET /v1/head` (and any future `text/plain` route) lacks `nosniff`, so under `script-src
   'self'` it is executable (V9, V10). It is harmless today and is a hardening candidate
   (nosniff on every daemon response, or a hash-source CSP).
2. Approval requests and decisions are not log entries, so no AG-UI event can carry them (V16).
   A `world.approval.*` CUSTOM name needs a log-derivable source first.
3. QUICKSTART section numbering collides at "10." (V23).
4. `w-world-live-surface.md` should point here for row 149 (done in this change).

## 10. Non-goals

- Rows 150–151: demo seed, captures, A2UI views, `world.view.a2ui`.
- Any change to `/agui/`, its CUSTOM table, or `host/agui`. Streaming approvals or decision
  packets (no source, V15/V16).
- Per-grade colour pills (V14); an object-page graph; graph evidence layers (V26).
- New query keys or fragment endpoints (V24). Actions of any kind from the workbench.
- A CI Chrome gate (V27). A JS test framework, `package.json`, or any npm dependency: node runs
  one plain file.

## 11. Mutation table (S6)

| Mutant | Test that must go red |
|---|---|
| M-a CSP gains `'unsafe-inline'` in `script-src` | `TestWorkbenchSecurityHeaders` (literal) |
| M-b CSP drops `connect-src 'self'` | `TestWorkbenchSecurityHeaders`, `assertWorkbenchSecurityHeaders` callers |
| M-c an inline `<script>` added to the page | `TestWorkbenchLiveScriptTag` (count) |
| M-d script tag added to the error template | `TestWorkbenchErrorPageInert` |
| M-e script route drops `nosniff` or serves `text/plain` | `TestWorkbenchLiveScriptRoute` |
| M-f script route accepts POST | `TestWorkbenchLiveScriptRoute` (405 arm) |
| M-g `/workbench/live.js` added to `isProtected` | `TestWorkbenchLiveScriptRoute` (posture arm) |
| M-h cursor via dense `GetLogEntry` scan | `TestWorkbenchLiveCursorCrossesGap` |
| M-i `LogEntriesLatest` ascending | `TestLogEntriesLatestNewestFirst` |
| M-j cursor rendered as latest+1 | `TestWorkbenchLiveCursorIsResumable` (404) |
| M-k live rows rendered as `<h3>entry N</h3>` | `TestWorkbenchTimelineBound` (measured, V12 m4) |
| M-l `<svg xmlns="http://www.w3.org/2000/svg">` | `TestRenderEmitsOnlyLocalLinks` (measured, V12 m1) |
| M-m new section placed after the provenance walk | `TestProvenanceWalkStaysLast` |
| M-n dark theme missing a token | `TestTokensBothThemes` |
| M-o footer empty-log branch renders nothing | `TestRenderQuietFooter` |
| M-p new view type not appended to the census | `TestWorkbenchViewFieldsAllRender` (after AC1.5 lists it; V13) |
| M-q graph links an unstored target | `TestGraphUnavailableNodeHasNoLink` |
| M-r graph row pitch < box height | `TestGraphNodesDisjointAndInBounds` |
| M-s object order from map iteration | `TestGraphDeterministic`, `TestGraphGolden` |
| M-t golden file truncated to 0 bytes | `TestGraphGolden` (anti-vacuity arm) |
| M-u graph renders a grade label | `TestGraphCarriesNoGrade` |
| M-v script splitter parses the unterminated tail | `TestLiveScriptPure` (byte-offset arm) |
| M-w script accepts a non-`replace` delta | `TestLiveScriptPure` (`add` arm) |
| M-x script's schema string drifts from `agui.StateSchema` | `TestLiveScriptContract` |
| M-y template renames `aria-label="timeline"` | `TestLiveScriptContract` (selector arm) |
| M-z fold: oldest decision wins | `TestFoldApprovals` |
| M-aa approvals walk unbounded | `TestRecentApprovalsBounded` |
| M-ab malformed chain → 500 page | `TestRecentApprovalsMalformed` |
| M-ac decisions pane gains a form/button | `TestWorkbenchDecisionsReadOnly` |

## 12. Quorum verification log

Not yet run. The pick-time quorum records here.
