# w-workbench-live-and-polish — row 149

- Status: **IMPLEMENTED 2026-10-08 (iteration 245, PR #234 → `b8b6a3f`) — was PLANNED — iteration 244/245 pick-time design, revision r2 (narrow-refinement
  carve-out after round 2)** (rounds 1 and 2 BLOCKED 3/3 present; no round 3; §12).
  Supersedes the row-149 parts of the attended DRAFT
  [`w-world-live-surface.md`](../planned/w-world-live-surface.md) (D6, D7, §4 row 149, §6 149(a)–(e)).
  Rows 150–151 stay there.
- Base `b702d73`. Premises are in §2. Probes are banked under
  `~/.ailang/state/world-iter244-design/` (`probe/` is a scratch Go module `replace`d onto this
  worktree; `scratch/` is a `git archive HEAD` copy used for compile and test probes). The r1
  probes (V29) are banked under `~/.ailang/state/world-iter245/design/` (`scratch/` is a fresh
  `git archive b702d73`; `v29_probe_test.go` sha256 `112e2e12…ba97fd`; `v29-probe.log`). The
  V29 probe and its log are committed as evidence under
  [`../verification/world-row149-design/`](../verification/world-row149-design/). The r2 probes
  (V31, V32, V33) are banked in the same `world-iter245/design/` directory.
- Binding: [coding-standards.md](../coding-standards.md) S2, S3, S6, S7; `D-WORLD-64`;
  [HUMAN-SURFACE.md](../HUMAN-SURFACE.md) P1, P3, P5, §5. Consumes the row-148 interface in
  [`../implemented/w-worldd-agui-event-stream.md`](w-worldd-agui-event-stream.md) §1, D1–D6.

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
| V5 | The row-148 contract this row consumes | `sed -n 17-25p host/daemon/agui.go`; `sed -n 10-16p host/agui/agui.go`; `docs/QUICKSTART.md:1012-1059` | `aguiRunBudget = writeTimeout - readDeadline - aguiWriteMargin` (18 s), `aguiTick 250ms`, `aguiPage 100`, **`aguiCap 16`**, `aguiMaxBody 64<<10`. `StateSchema = "world/agui-state/v1"`. **`customNames = []string{"world.entry.committed"}`** (`:16`) is the only CUSTOM. Cursor comes from `state.lastIndex` or `Last-Event-ID`. An unknown cursor gets 404 before any byte; `-1` is genesis, never 404 (V29). A run ends with `RUN_FINISHED.result.lastIndex`. |
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
| V29 | `state.lastIndex = -1` is genesis on an empty **and** a non-empty store (r1, quorum objection 2) | Code: `sed -n 51-80p;139-152p host/daemon/agui.go`. Runtime: `v29_probe_test.go` (committed as `design_docs/verification/world-row149-design/v29_probe_test.go.txt`, log `v29-probe.log` beside it) copied into `world-iter245/design/scratch/host/daemon/`, `AILANG_BIN=~/.pinned-ailang/ailang go test ./host/daemon/ -run TestV29 -v -count=1` (`v29-probe.log`). Body `{"threadId":"workbench","runId":"wb-0","messages":[],"state":{"schema":"world/agui-state/v1","lastIndex":-1}}` | Code: `aguiCursor` starts `after := int64(-1)` and refuses only `after < -1` ("invalid state.lastIndex"). `handleAGUI` calls `GetLogEntry` (404 "unknown cursor") only `if after >= 0`. **(a) Empty store:** `status=200 content-type="text/event-stream"`; lines `data: {"type":"RUN_STARTED","threadId":"workbench","runId":"wb-0","protocolVersion":"1.0"}`, `data: {"type":"STATE_SNAPSHOT","snapshot":{"schema":"world/agui-state/v1","lastIndex":-1,"logHead":null}}`; entry 0, committed after the run's first poll, arrives as `data: {"type":"CUSTOM","name":"world.entry.committed","value":{"header":{"entryIndex":0,…` (`V29a DELIVERED entry 0`). **(b) Store 0..2:** `status=200`, first frame `data: {"type":"RUN_STARTED","threadId":"workbench","runId":"wb-0","protocolVersion":"1.0"}`, CUSTOM indexes `[0 1 2]`. **Controls (the instrument sees rejection):** `lastIndex=3` on 0..2 → `404 {"error":{"class":"NotFound","message":"unknown cursor"}}`; `lastIndex=-2` → `400 … "invalid state.lastIndex"`; `lastIndex=0` on the empty store → `404 … "unknown cursor"` (the M-ad premise). Contract text: row 148 D4 (`../implemented/w-worldd-agui-event-stream.md:159-166`): "`after = state.lastIndex`, which must be an integer ≥ −1" and "A non-negative `after` must name an **existing** entry"; `:118-119` "`logHead`: <hash of entry n, or null when n = -1>". |
| V30 | `isProtected` is an exact single-route match; every GET is unprotected by construction (r1, quorum objection 1) | `sed -n 903,905p host/daemon/daemon.go`; positive control `sed -n 1082,1083p host/daemon/daemon_test.go`; `grep -rn isProtected host --include='*.go'` | `func (d *Daemon) isProtected(r *http.Request) bool {` / `return r.Method == http.MethodPost && r.URL.Path == "/v1/commit"` / `}`. Exact `==` on method and path, no prefix. Control: `daemon_test.go:1082` asserts `isProtected(POST /v1/commit)` is true ("the predicate probe is vacuous" otherwise), so the predicate does return true for its one route. Other call sites: `daemon.go:895` (`Wrap(d.isProtected, mux)`), `agui_test.go:454` (row 148's posture pin), `daemon_test.go:1077`. **So `GET /workbench/live.js` and `GET /workbench` are both `false`, and the body needs no change.** |

| V31 | `readStore` satisfies D7's `ApprovalReader` without an import cycle (r2, objection B2) | `sed -n 484,495p host/daemon/daemon.go`; scratch (`world-iter245/design/scratch`): add `host/broker/v31_approval_reader.go` declaring `ApprovalReader interface{ GetRegistryHead(ctx, name string) (hashref.HashRef, bool, error); GetObject(ctx, hashref.HashRef) (store.Object, bool, error) }` (copy `v31_approval_reader.go.txt`) and `host/daemon/v31_assert.go` = `var _ broker.ApprovalReader = (readStore)(nil)`, then `go vet ./host/broker/ ./host/daemon/`; control: add `NoSuchMethod()` to the interface; `go list -deps ./host/broker \| grep -c 'host/daemon$'` | `readStore` holds `GetObject(ctx context.Context, ref hashref.HashRef) (store.Object, bool, error)` and `GetRegistryHead(ctx context.Context, name string) (hashref.HashRef, bool, error)`. **vet `rc=0`, empty output.** Control: `vet: host/daemon/v31_assert.go:5:31: … readStore does not implement broker.ApprovalReader (missing method NoSuchMethod)`. Deps count `0`: `host/broker` does not import `host/daemon`; the assertion lives in `host/daemon`, which already imports `host/broker` (`daemon.go:56`). |
| V31b | `D-WORLD-64` grants the scoped rule-(d) exception and contains **no** CSP language; the CSP literal is queue row 149's text (r2, objection C1) | `grep -n '\| D-WORLD-64 \|' design_docs/world-mission.md`; `grep -n '^149\. ' design_docs/world-mission.md`; `sed -n 950p design_docs/world-mission.md \| grep -o -i 'csp\|script-src\|connect-src' \| wc -l`; control `grep -c "script-src" design_docs/world-mission.md` | Ledger row `:950`, verbatim exception: "(4) **Placement:** they move no UNMET clause, so they never displace a routable critical-path row (the D-WORLD-58/61 order as amended by later attended rulings). **Exception to rule (d), scoped to these four rows only:** when no critical-path row is routable, the loop MAY pick the lowest-numbered open row of 148–151 instead of a no-pick block report, so blocked fires produce release work. Every other rule-(d) restriction stands." Check (a): the scope covers release rows that close no UNMET clause, **yes** ("they move no UNMET clause", "scoped to these four rows only"). Check (b): CSP language in the row, **none** (`wc -l` → `0`). Control: the whole charter has exactly `1` line with `script-src`, and it is queue row 149. Positive control: queue row `:6122` contains "CSP widened by exactly `script-src 'self'; connect-src 'self'`, pinned by a literal-string test; no inline script, no CDN, no build chain". |
| V32 | Region census per page variant at `b702d73` (r2, objection B1) | `v32_regions_probe_test.go` (sha256 `14efc3b6…dd65a`) in the scratch `host/daemon/`, `AILANG_BIN=~/.pinned-ailang/ailang go test ./host/daemon/ -run TestV32 -v -count=1` (`v32-probe.log`); counts `aria-label="<label>"` occurrences | `home-empty-store /workbench status=200 world browser=1 timeline=1 inspector=1 selected entry=0 provenance walk=1 live=0 world graph=0 decisions=0 live status=0`; `home` identical; `entry /workbench?from=0&entry=0 status=200 world browser=1 timeline=1 inspector=1 selected entry=1 provenance walk=1 live=0 …`; `object /workbench?object=sha256:d9b46… status=200 world browser=1 timeline=1 …` (same). Control: `/workbench?entry=0` (no `from`) is `status=400` with every count 0 — the error page has no regions (V25), which is why a non-200 re-fetch is never parsed. **So the two existing swapped regions occur exactly once on every 200 variant; `live`, `world graph` and `decisions` do not exist yet, and this row must render them unconditionally on every variant (D1 missing-node contract, AC3.5 table).** |
| V33 | Hash-source vs `'self'` for an **external** script, in real browsers (r2, objection C1 re-decision) | `world-iter245/design/csp-hash-probe/main.go` (sha256 `da56df5a…a1c502`): four pages, CSP `default-src 'none'; <variant>; connect-src 'self'`, one `<script src="/s.js" defer>` whose body POSTs `/report?m=RAN-<variant>`; Chrome `--headless=new` 154.0.8037.98 and Firefox `--headless --screenshot` 157.0.1, unsandboxed; `server.log`, `server-ff.log` | Script hash `sha256-2B0qxfomDLK6lImzhAYYm/2gRYZrU3BaG6xPZKEWlKk=`. **Chrome:** `REPORT RAN-self`, `REPORT RAN-hash-sri`; `hash-nosri` and `hash-wrongcsp`: page fetched, **no** script fetch, no report; console "Loading the script 'http://127.0.0.1:47901/s.js' violates the following Content Security Policy directive: "script-src 'sha256-2B0q…'"". **Firefox:** `REPORT RAN-self`, `REPORT RAN-hash-sri`; `hash-nosri` and `hash-wrongcsp` page-only. Safari: **UNMEASURED** (installed, but not headless-drivable from this rig cheaply). Chrome rc=142 is the `perl alarm 40` (V8's dump-dom hang); the reports landed before it. |
| V34 | `MaxLogEntryPage` | `sed -n 11,12p host/store/log_after.go` | `// MaxLogEntryPage bounds each keyset log read.` / `const MaxLogEntryPage = 500` |
| V35 | The edge helper the graph reuses checks exactly three entry relations | `sed -n 155-176p host/daemon/workbench.go`; `sed -n 152p host/workbench/render.go`; `grep -n 'stateRoot: <a' host/daemon/workbench_test.go` | `func (d *Daemon) entryEdges(ctx context.Context, entry store.LogEntry) ([]workbench.EdgeView, error)` at `:159`, relations `transitionFn`, `interpreter`, `transitionRef`, each via `checkedEdge` (`:182`, one `d.reads.GetObject`). `stateRoot` is rendered from the world nav (`render.go:152` `{{template "edge" .World.StateRoot}}`) and pinned there (`workbench_test.go:453`, `:791`). |
| V36 | Row 148's 503 sets `Retry-After: 1` | `grep -n 'Retry-After' host/daemon/*.go \| grep -v _test` | One hit: `host/daemon/agui.go:122: w.Header().Set("Retry-After", "1")`. |

**Draft claims false or stale at `b702d73`** (all in [`w-world-live-surface.md`](../planned/w-world-live-surface.md)):

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
   committed index at render time, or `-1` on an empty log (D3). `-1` is genesis on the
   `/agui/` side, never an unknown cursor (V29).
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
   `section[aria-label="live"]`, `section[aria-label="world graph"]` and
   `section[aria-label="decisions"]`. **The footer is never swapped.** The script owns
   `<span data-live-status>`; on every applied swap it re-renders the current status
   (live/paused/failed) from script state, so a refresh can never revert the footer to the
   server-rendered off text. The server-rendered span text is the no-JS state only. (The
   log-position sentence that must refresh lives in the swapped live section, D5.) **Missing
   node contract (r2):** every one of the five regions is rendered unconditionally on every
   200 workbench variant (home, `?from=&entry=`, `?object=`; V32), so a one-sided miss means
   the layout changed under the page. If a selector is absent on either side, `applySwap`
   skips that region and `renderStatus` reports "paused: page layout changed, reload
   required" — never silently, and it never calls `replaceWith` on null. The stream stops
   (no re-POST) until reload. A non-200 re-fetch is never parsed (§7). Newly arrived live
   rows get a CSS class `fresh` (animated in the stylesheet).
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
  This holds with `isProtected`'s body unchanged: it is an exact match on `POST /v1/commit`
  only, so every GET is `false` by construction (V30).
- The page head gains exactly `<script src="/workbench/live.js" defer></script>`. There is no
  other `<script` anywhere: no inline script, no `on*=` handler, no `eval`, no CDN. The error
  page stays script-free (V25).
- **New CSP literal** (the order measured in V6):
  `default-src 'none'; style-src 'unsafe-inline'; script-src 'self'; connect-src 'self'; base-uri 'none'; form-action 'self'; frame-ancestors 'none'`.
  The two V3 literal lines are updated to this string. They stay literals, never the symbol.

*Where the literal comes from (r2 correction):* the widening "exactly `script-src 'self';
connect-src 'self'`, pinned by a literal-string test" is **queue row 149's text** (and the
attended draft's D6), not `D-WORLD-64`. The ruling routes rows 148–151 and grants the scoped
rule-(d) exception; it contains no CSP language (V31b). r0/r1 wrongly attributed the pin to the
ruling.

*Re-decision: `'self'` versus a `'sha256-…'` hash-source (designer's decision, r2).* Because the
ruling does not pin the keyword, the choice is made here on the merits:

| Criterion | `script-src 'self'` | `script-src 'sha256-<live.js>'` + `integrity=` |
|---|---|---|
| What executes | Any same-origin GET body served as script; today that adds `GET /v1/head`'s daemon-minted `text/plain` hash (V9, V10) | Only the exact `live.js` bytes (V33 `hash-wrongcsp` control: a mismatched hash blocks) |
| Browser support for an **external** script | CSP2; measured Chrome 154 and Firefox 157 run it (V6, V33 `self`) | CSP3 ties hash-sources on external scripts to a matching `integrity` attribute. Measured: Chrome 154 and Firefox 157 run it **with** `integrity` and block it **without** (V33). Safari **UNMEASURED**. A browser that does not support it fails closed to the JS-off page |
| Determinism | Fixed literal | Also deterministic: the hash is a function of the embedded bytes, computable in Go at init (`crypto/sha256` over `LiveScript`), so no build chain |
| Pin and churn | One literal in the two V3 places, stable across script edits | The literal changes on every `live.js` edit (two test literals plus the template's `integrity`), or the test derives it from the symbol, which V3's literal discipline forbids |
| Scope | Matches row 149's text | Narrower than row 149's text; changing an executed row's acceptance text is a queue-row change, not a design choice |

**Decision: keep `'self'`.** Grounds: (1) row 149's text specifies it, and the r2 carve-out
forbids widening or changing direction; (2) the extra surface it admits is bounded today (no
same-origin GET returns caller-controlled bytes as non-`nosniff` script, V10, §8); (3) Safari
support for the hash form is unmeasured; (4) a hash pin turns every script edit into a security
-header edit. The hash-source form is now **measured feasible** in Chrome and Firefox and is
filed as the hardening candidate in §9 item 1, for Mark to take as a row-text change if wanted.

*Other alternatives:* a hashed path (draft D6) costs a computed route to bust a cache the route
never fills. A nonce cannot be used: a per-request nonce makes the page non-deterministic and
needs a script-bearing template.

### D3 — The cursor and the live list come from one new newest-first store read

Add `Store.LogEntriesLatest(ctx, limit int) ([]LogEntry, error)`:
`SELECT … FROM log_entries ORDER BY entry_index DESC LIMIT ?`, `1 ≤ limit ≤ MaxLogEntryPage`
(= 500, `host/store/log_after.go:12`, V34),
deadline required, quarantine refused. It reuses `scanLogAfterRow` and joins `readStore`
(V18). `handleWorkbench` calls it once with `limit = 10`. The result gives `Live.Recent`
(newest first) and `Live.Cursor = Recent[0].EntryIndex`, or `-1` on an empty log. It gives
`Live.Head`, the entry hash, or empty. Every render pays one extra indexed statement.

- **Why a store read:** V17. Nothing yields the latest index. A dense `GetLogEntry` scan
  stalls at gaps (row 148 V10/V26). Paging `LogEntriesAfter` to the end is O(log length).
- **Why the cursor must be an existing index or `-1`:** `/agui/` returns 404 for an unknown
  non-negative cursor (V5, V29 control `lastIndex=3`). The max existing index is always valid.
  On an empty log the cursor is `-1`, which row 148 D4 defines as genesis ("must be an integer
  ≥ −1"; only a non-negative cursor must name an existing entry). Measured: `-1` on an empty
  store gets 200 and the first committed entry 0 is delivered on that stream (V29a); `0` on
  an empty store would be 404 (V29 control). A commit landing between render and the first
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

### D5 — Quiet is health: a positive log assertion plus a script-owned status footer

The log-position sentence lives in the **swapped** `section[aria-label="live"]`, so every
refresh updates it: "Log at entry N (head `<hash>`) when this page was rendered", or "The log
is empty". The server renders `footer[aria-label="live status"]` with `role="status"`
containing only `<span data-live-status>Live updates: off. Reload to refresh.</span>`. That
text is the no-JS state only. The footer is never swapped (D1 step 4). There is no server
wall clock, so the render stays deterministic. With JS on, the script owns the span. It
renders it from script state at start, after each `RUN_FINISHED`, on every wait, and after
every applied swap: "Live: checked HH:MM:SS, no new entries after entry N; next check within
18 s". After a burst it reads "K new entries since HH:MM:SS". On failure it reads "paused:
<reason>, retrying in Ns" (D8). That is P5's "last verified heartbeat, next expected, what was
checked" (V22). A dead stream is never rendered as calm, and a live one is never rendered as
off.

| Event | Status span |
|---|---|
| JS off or blocked | Server text "Live updates: off. Reload to refresh." (unchanged forever) |
| refresh swap during a run | The status span is re-asserted from script state; the footer never reverts to off while the stream is live |
| Wait (503, network, `RUN_ERROR`, hidden tab) | "paused: <reason>, retrying in Ns" — survives any swap, because the footer is not swapped and the span is re-asserted |

### D6 — The world graph: a deterministic bipartite SVG of entries → objects, from checked edges only

`workbench.LayoutGraph(entries []GraphEntry) GraphView` is pure. Each `GraphEntry` is
`{Index, Href, Edges []EdgeView}`, where the edges are the entry's already-checked
`transitionFn`/`interpreter`/`transitionRef` edges.

- **Where:** the home page uses the newest 5 entries of `Live.Recent`, each with the existing
  `d.entryEdges` (`host/daemon/workbench.go:159-176`, one `checkedEdge` → `GetObject` per
  relation, V35) (≤ 15 extra `GetObject` reads). The selected-entry page uses `[Selected]` with
  the edges it already has. The object page renders "graph: select an entry" (a non-goal for
  now, §10), so the region still exists on that variant (D1 missing-node contract).
- **Why only these three relations:** they are exactly the entry's outbound edges that
  `entryEdges` checks (V35). `stateRoot` is a **world** edge, already shown once in the world
  nav (`render.go:152`), and `committedBy` is an **inbound** object walk (the provenance
  section). Neither is an entry→object edge. The graph also never renders any relation in the
  `relation: <a` partial form (labels go in SVG `<title>`), so the body-wide V11d counts of
  `committedBy: <a` and `stateRoot: <a` are untouched.
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
The daemon maps each summary into a `workbench.DecisionRow` inside `workbench.DecisionsView`
(so `host/workbench` stays pure, V19). `readStore` already satisfies `ApprovalReader` (V31).
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
seconds (row 148 sends `Retry-After: 1`, `host/daemon/agui.go:122`, V36), then backs off
exponentially (×2, cap 30 s). Retry-After absent or unparsable → use the first backoff step
(1 s). On a network error, a non-200, or
`RUN_ERROR`, it does the same. Every wait is shown in the D5 span, and a refresh swap never
erases it (the footer is not swapped; the span is re-asserted after each swap, D1 step 4). The page itself never
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
  rows newest first, no `<h3>entry `, and inside the live section "Log at entry 7 (head …)".
  `Cursor:-1` renders `data-live-cursor="-1"` and "The log is empty" inside the live section.
- AC1.3 `TestRenderQuietFooter`: both states render a footer that is outside `<main>` and holds
  `<span data-live-status>Live updates: off`. The log-position sentence is **not** in the footer
  (it is in the swapped live section, D5).
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
  context (stuck guard). **Empty-store arm (r1):** on an empty store, the cursor read from the
  rendered `data-live-cursor` (`-1`) is POSTed over an `httptest.NewServer`; the response is 200;
  after the run's first poll (the `aguiPollReads` pattern of row 148's
  `TestAGUISeesDirectStoreCommit`) the test commits entry 0 and asserts a `CUSTOM` with
  `"entryIndex":0` arrives on that same stream (V29a measured exactly this). Stuck guard 5 s per
  wait.
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
  `isProtected` is equal (both `false`) for `GET /workbench/live.js` and `GET /workbench`
  (V30).
- AC3.5 `TestLiveScriptContract` (binary-free): from the embedded script, extract the
  `STATE_SCHEMA` literal and assert it equals `agui.StateSchema`. Extract the one-line
  `REGIONS` JSON array; assert it has exactly 5 selectors and none names `footer` or
  `live status` (D1 step 4). **Table-driven per variant (r2):** for each row of
  {home on an empty store, home on a seeded store, `?from=0&entry=0`, `?object=<stored ref>`}
  (the V32 variants), the daemon page is 200 and every selector's `aria-label` occurs exactly
  once. Extract the POST path and assert it equals `/agui/`. Null control: the test fails if
  any extraction finds nothing or the variant table is empty.
- AC3.6 `TestLiveScriptPure` (node): runs `node` (from `WORLD_EXEC_NODE` or `PATH`) on the
  script plus a harness (the script exports `splitFrames`/`applyDelta`/`nextBackoff`/`applySwap`/
  `renderStatus` only when `module` exists, and starts only when `document` exists). It feeds
  `host/agui/testdata/stream_fixture.golden` split at **every** byte offset into two chunks.
  The event list equals the whole-buffer parse, and the cursor never exceeds the last fully
  delivered entry. `applyDelta` with an `add` op throws. Backoff is 1, 2, 4 … capped at 30. A
  missing or unparsable `Retry-After` gives the first step, 1 (D8).
  **Footer arm (r1):** the harness builds a minimal stub document (plain objects with
  `querySelector`/`replaceWith`/`textContent`; no DOM library), rewrites the status span to the
  script's live status, then calls `applySwap` with a fresh document whose footer holds the
  server-rendered `Live updates: off. Reload to refresh.` span, and asserts the span still
  carries `renderStatus(state)` and that the footer stub's `replaceWith` was never called.
  A second sub-arm makes the re-assert itself load-bearing: the span holds the server "off"
  text, script state is `paused`, and after `applySwap` the span must read
  `renderStatus(paused)`. **Missing-region sub-arm (r2):** a stub fresh document that omits
  `section[aria-label="decisions"]` (and, separately, an old document that omits it):
  `applySwap` does not throw, swaps the other four regions, does not touch the missing one,
  and the span reads "paused: page layout changed, reload required".
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
  `TestRenderEmitsOnlyLocalLinks` and the body-wide count tests pass (V11). **Census (r2):**
  `GraphView` (and every struct type it holds, e.g. its node and edge row types) is appended
  to `workbenchViewTypes`, and `TestWorkbenchViewFieldsAllRender` continues to pass (V13).

**M5 — Decisions pane (effectful walk + pure fold; ~180 LOC).**
- AC5.1 `TestFoldApprovals`: request A, request B, decide A deny, request C, decide C approve
  gives the newest-first order [C approve, A deny, B pending]. A request decided twice shows the
  newest decision.
- AC5.2 `TestRecentApprovalsBounded`: a 200-head chain does ≤ 40 head-object GetObject reads
  and ≤ 2×kept-summaries additional body reads (≤ 80 total on the counting fake) and sets
  Truncated. (Checked against D7: the walk reads ≤ 40 head objects; each of the ≤ 20 kept
  summaries needs its request body for effect/scope/requester/cost and, when decided, its
  decision body for `DecidedBy`, so ≤ 40 + 2×20 = 80. The registry head is one
  `GetRegistryHead`, not a `GetObject`.)
- AC5.3 `TestRecentApprovalsMalformed`: a wrong-semantic-id head gives
  `ErrApprovalChainMalformed`. The workbench renders the UNAVAILABLE line with status 200.
- AC5.4 `TestWorkbenchDecisionsPane`: a real `HumanHandler` request plus `DecideApproval` on a
  daemon store renders a checked request link and "deny by <decidedBy>". With no head it renders
  "No approval requests recorded". **Census (r2):** `DecisionsView` (and `DecisionRow`) is
  appended to `workbenchViewTypes`, and `TestWorkbenchViewFieldsAllRender` continues to pass
  (V13).
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
| AC2.4 cursor resumable | `aguiBudget` 300 ms; client ctx 10 s. Empty-store arm: `aguiBudget` 5 s, commit after the first poll, each wait ≤ 5 s | 300 ms and the post-poll commit = **stimulus**; 10 s and 5 s = stuck guards. It asserts status and frames, never elapsed time |
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
| A region selector is missing on either side of a swap | `applySwap` skips it, never throws; the span reads "paused: page layout changed, reload required"; the stream stops until reload (D1, AC3.6 missing-region sub-arm) |
| 503 without a parsable `Retry-After` | First backoff step, 1 s (D8) |
| Refresh swap during a run | The footer is not swapped; the status span is re-asserted from script state, so it never reverts to the server "off" text while the stream is live, and a paused/retrying text survives (D1 step 4, D5; AC3.5, AC3.6 footer arm) |
| Fresh world (empty log) | Cursor `-1` (genesis): 200, and entry 0 is delivered when it commits (V29a; AC2.4 empty-store arm) |
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
  patterns and handlers, `isProtected`'s body (exact match on `POST /v1/commit`, so the new GET
  is unprotected without an edit, V30), `acceptedWorkbenchKeys` and
  `supportedWorkbenchQuery`, `host/agui/*`, `host/daemon/agui.go`, every D7 constant,
  `workbenchErrorTemplate`, `tools/launchd/*`. No open critical-path row names these files at
  `b702d73`. Re-check at merge.

## 9. Findings to file (not fixed here)

1. `GET /v1/head` (and any future `text/plain` route) lacks `nosniff`, so under `script-src
   'self'` it is executable (V9, V10). It is harmless today and is a hardening candidate
   (nosniff on every daemon response, or a hash-source CSP; the latter is measured working on
   an external script with `integrity` in Chrome 154 and Firefox 157, Safari unmeasured, V33;
   adopting it changes row 149's text, so it is Mark's call).
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
| M-o live section's empty-log branch renders nothing | `TestRenderLiveRegion` (`Cursor:-1` arm) |
| M-p any of the three new view types (`LiveView` AC1.5, `GraphView` AC4.6, `DecisionsView` AC5.4, or their row types) not appended to the census | `TestWorkbenchViewFieldsAllRender` (V13: an unlisted type's unrendered field is not caught, so each AC's append is what arms it; the mutant is: append the type, leave one field unrendered) |
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
| M-aa approvals walk unbounded (> 40 head reads) or body reads not bounded by kept summaries (> 80 total) | `TestRecentApprovalsBounded` |
| M-ab malformed chain → 500 page | `TestRecentApprovalsMalformed` |
| M-ac decisions pane gains a form/button | `TestWorkbenchDecisionsReadOnly` |
| M-ad empty-log cursor rendered as `0` instead of `-1` | `TestWorkbenchLiveCursorIsResumable` (empty-store arm: status 404 ≠ 200; premise measured, V29 control) |
| M-ae `footer[aria-label="live status"]` re-added to `REGIONS` | `TestLiveScriptContract` (footer-exclusion arm) |
| M-af `applySwap` skips the status re-assert | `TestLiveScriptPure` (footer arm, second sub-arm: span still reads the server "off" text) |
| M-ag `applySwap` throws on, or silently skips, a one-sided missing region | `TestLiveScriptPure` (missing-region sub-arm: throw, or span not "paused: page layout changed…") |
| M-ah a new region rendered only on some variants (e.g. graph omitted on `?object=`) | `TestLiveScriptContract` (per-variant table row) |
| M-ai `Retry-After` missing → wait 0 or NaN | `TestLiveScriptPure` (backoff arm: `nextBackoff(undefined)` and `nextBackoff("x")` = 1) |

## 12. Quorum verification log

**Round 1: BLOCKED.** 3 of 4 seats present; `gpt6-1-sol` unreachable, so N−1. Each objection's
premise was measured first-party before revising (controller code reads, then this designer's
own commands in V29/V30). No objection disputed direction.

| Reviewer | Verdict | Objection (one line) | Premise measured | Surface | r1 answer |
|---|---|---|---|---|---|
| gemini-3-1-pro | BLOCK | D2/§8 claim `isProtected` stays unchanged and is equal for the two GETs, but §2 never shows exact vs prefix matching | REAL gap in evidence, claim TRUE: exact `==` on `POST /v1/commit` (V30) | §2, D2, §8, AC3.4 | V30 added with the source lines and a positive control; D2 and §8 cite it; the posture-pin test (AC3.4, M-g) is kept |
| oc-glm-5-3 (and oc-kimi-k3, second half) | BLOCK | `lastIndex = -1` on an empty log was never probed; if refused, a fresh world's live surface is dead on arrival; no empty-store AC | REAL gap in evidence, claim TRUE: 200 on both stores, entry 0 delivered (V29) | §2, V5, D1, D3, AC2.4, §6, §7, §11 | V29 probe on an empty store and a 0..2 store, with 404/400 controls; row 148 D4's "≥ −1" quoted; AC2.4 empty-store arm; mutant M-ad |
| oc-kimi-k3 | BLOCK | Swapping the footer reverts a working stream's status span to the server "off" text (and wipes D8's paused text) on every refresh | REAL (by construction of r0's D1 list) | D1, D5, D8, AC1.2, AC1.3, AC3.5, AC3.6, §7, §11 | Footer removed from the swap list; script owns and re-asserts `<span data-live-status>` after every swap; log-position sentence moved into the swapped live section; D5 status table row "refresh swap during a run"; AC3.6 footer arm; AC3.5 exclusion arm; mutants M-ae, M-af |

**Round 2: BLOCKED.** 3 of 4 seats present (`gpt6-1-sol` absent again). All three present seats
rejected r1. Every objection carried a concrete fix and none disputed direction, so the
controller took the **narrow-refinement carve-out**: r2 applies the fixes and goes to the
planner with **no round 3**. Premises were measured first (controller, then V31–V36).

| Reviewer | Verdict on r1 | Objection (one line) | Surface | r2 answer |
|---|---|---|---|---|
| gemini-3-1-pro | REJECT | `GraphView`/`DecisionsView` are never appended to `workbenchViewTypes` (V13 census rule) | AC4.6, AC5.4, §11 M-p | Both ACs now require the append (plus their row types) and that `TestWorkbenchViewFieldsAllRender` still passes; M-p covers all three additions |
| oc-glm-5-3 | REJECT | (1) a REGIONS selector missing on either side makes `replaceWith(null)` throw and freezes the status as live | §2 V32, D1 step 4, AC3.5, AC3.6, §7, §11 | V32 census of all page variants; regions rendered unconditionally; missing-node contract (skip + "paused: page layout changed, reload required", never silent); AC3.5 table-driven per variant; AC3.6 missing-region sub-arm; M-ag, M-ah |
| oc-glm-5-3 | REJECT | (2) `readStore` satisfying `ApprovalReader` is asserted, not measured | §2 V31, D7 | V31: vet-green assertion in `host/daemon` with a `NoSuchMethod` red control; no import cycle |
| oc-glm-5-3 | REJECT | (3) `MaxLogEntryPage`, the `entryEdges` helper, the edge exclusions, and `Retry-After` are uncited; absent `Retry-After` undefined | §2 V34–V36, D3, D6, D8, AC3.6, §7, §11 | Cited (`log_after.go:12`, `workbench.go:159`, `agui.go:122`); D6 states why `stateRoot`/`committedBy` are excluded against V11d; D8: absent or unparsable `Retry-After` → 1 s; M-ai |
| oc-kimi-k3 | REJECT | (1) the CSP pin is attributed to `D-WORLD-64`, which contains no CSP language; re-decide `'self'` vs hash-source | §2 V31b, V33, D2, §9 | V31b quotes the exception verbatim (scope check yes, CSP language none); attribution corrected to row 149's text; D2 re-decision table on the merits with Chrome/Firefox measured and Safari UNMEASURED; **decision: keep `'self'`**, hash-source filed as measured-feasible hardening |
| oc-kimi-k3 | REJECT | (2) §12 lacks a round-2 record | §12 | This table and the carve-out note |
| oc-kimi-k3 | REJECT | (3) AC5.2's "≤ 40 GetObject reads" ignores the body reads per summary | AC5.2, §11 M-aa | Amended verbatim (≤ 40 head reads, ≤ 2×kept-summaries body reads, ≤ 80 total); checked against D7's walk (40 + 2×20 = 80), so the reviewer's numbers stand; M-aa updated |
