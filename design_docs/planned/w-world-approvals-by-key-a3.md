# w-world-approvals-by-key-a3: the surface approves, the minter mints

- **Status**: Planned. Designed in iteration 383, **r0**. Queue row 201, `w-world-approvals-by-key`,
  parent milestone **A3**. Parent: [`w-world-approvals-by-key.md`](w-world-approvals-by-key.md)
  (attended; its rulings are not reopened). Built siblings, settled and not reopened:
  [`-a0a1`](w-world-approvals-by-key-a0a1.md) (A0 request object, digest, challenge, expiry,
  consume; A1a), [`-a1b`](w-world-approvals-by-key-a1b.md) (SSHSIG approvers, `approver
  enroll|list|revoke`, `approve`), [`-a2`](w-world-approvals-by-key-a2.md) (the WebAuthn law,
  credentials schema 6, `approve --webauthn-link|--webauthn-assertion`, the ceremony page).
  **A3 is split in two (§0): A3a, the web surface flow, which the loop builds and judges headless
  in three sprints; and A3b, the Tauri-native notification and Touch ID convenience, which has no
  host to land in yet (§1.9) and is designed here as constraints plus an attended exit run only.**
- **Authority**: **D-WORLD-89** (Mark, attended 2026-10-21, "all defaults"): R1 a verified
  signature from an enrolled key may approve a mint, as well as the typed `y`; R2 the formats are
  WebAuthn (UV required) and OpenSSH SSHSIG, and "native app authenticators sign only in those
  formats"; R3 enrolment and revocation stay typed-only; R4 publishing stays typed-only.
  **D-WORLD-109 = A** (Mark, attended 2026-11-06): "The loop builds and judges against fakes and
  test keys, and never enrols a real key, mints or publishes." Every A3 test uses temp stores,
  `approvaltest` keys and Chrome's CDP **virtual** authenticator in a temp profile. Mark's
  real-passkey, real-phone, notification and Touch ID runs are the attended exit run (§10). The
  loop never runs them. The parent §3 invariants bind every section below: the agent never
  approves; publishing stays typed-only; the typed `y` keeps working; read routes stay as open as
  row 199 leaves them; approving is a POST bound to a verified signature, never a cookie or a
  session.
- **Item** (A3a): four things.
  - **A recorded approval** (credentials schema 7, `surface_approvals`): the daemon stores a
    POSTed signed artifact *only after* the A1b/A2 law verified it, and the store refuses it for a
    used, expired or changed request.
  - **The minter**, `approve <request-id> --await`: no terminal; it waits for a recorded approval,
    **re-judges it with the same law**, then runs `approve`'s steps 2–6 unchanged. The mint stays
    in the CLI process, where it is today (§1.3).
  - **Four routes under `/surface`** (none under `/v1`): `GET /surface/requests.json` (*needs
    you* and the decision list), `GET /surface/requests/{id}` (the request card's data),
    `GET /surface/approve/{id}` + `GET /surface/approve.js` (the card page), and
    `POST /surface/approve/{id}` (the Approve). The new routes check `Host`; the POST also checks
    `Origin`, `Sec-Fetch-Site` and `Content-Type`.
  - **The surface app's *Needs you*** lists the pending requests and links each to its card; a
    **decision list** shows who approved what; an opt-in **browser notification** fires on a new
    request.
  It changes no `.ail` code.
- **Designer lane**: `claude:claude-opus-5-5`.
- **Estimated**: A3a in **3 sprints, about 5,100 changed lines** (§9): sprint 1 ≈ 1,650, sprint 2
  ≈ 1,700, sprint 3 ≈ 1,750. **Calibration (V20):** A2's four sprints measured 1,969, 1,843, 1,686
  and 2,768 changed lines outside `design_docs/`; its sprint 4 was planned at 1,054 lines (2.6×).
  The estimates here are sized from those actuals, and each sprint has a pre-named fallback cut.
- **Measurement base**: `0140bd7f` (PR #552's head: A2 sprint 4 built, judged PASS 91, not yet
  merged because of a CI runner outage). Worktree
  `/Users/voightkampff/dev/sunholo-data/.wt-world-iter383-a3`, branch
  `design/world-iter383-row201-a3`. sha256 prefixes (V1): `host/daemon/daemon.go`
  `9a9357770a4b0360`, `host/daemon/surface_page.go` `f769ee06b047f113`, `host/surfacepage/model.go`
  `4ace112b84f959b4`, `host/surfacepage/app/app.js` `8ba97c83d7e78a69`,
  `host/surfacepage/passkey/passkey.js` `0064c6c90099a356`, `cmd/ailang-worldd/approve.go`
  `b28c5c4b47cb5b6a`, `approve_passkey.go` `bd95db229e9498cf`, `approve_sig.go`
  `446c0f6df1c49ad3`, `host/store/mint_requests.go` `b7113f7d5df7e885`, `host/store/credentials.go`
  `0e32e12f8ef7f150`, `cmd/ailang-worldd/session_new.go` `587f39f99ace013d`,
  `cmd/ailang-worldd/testdata/passkey_cdp.mjs` `2a336e5e810f8ed6`.
- **Instruments**: `go1.26.6 darwin/arm64`; Google Chrome at
  `/Applications/Google Chrome.app/Contents/MacOS/Google Chrome`, headless, temp profiles, a
  virtual authenticator only; `node v26.13.0`. One throwaway probe lived in `/tmp/w383-designer/`
  and is not in the repo (V17). The designer's shell has no controlling terminal (A2 V29).

---

## 0. Scope, and the sprint boundaries

**The cut, and why (measured).** The parent's A3 sentence has two halves: "the surface flow:
*needs you*, the request card and Approve" and "plus a native notification and the optional Touch
ID convenience in the Tauri app".
- The first half is web code the daemon serves, and the A2 harness already drives such a page in a
  real headless Chrome with a virtual authenticator (`TestPasskeyPageInChrome` passes at the base,
  V15). The loop can build and judge it: **A3a**.
- The second half has **no Tauri app to land in.** The repo's only Tauri code is the design fixture
  `design_docs/mockups/surface-app-v3/`, which "is not product code, and no test or build reads
  it" (V12). Row 188's N3, "the Tauri shell in the repo (`desktop/`)", is not built: `desktop/`
  does not exist, the N3 plan says "The World sandbox has no `cargo`", and CI has no Rust, Tauri or
  Chrome step (V12, V13). Under D-WORLD-109 = A the loop also has no GUI to judge a native
  notification or a Touch ID prompt with. So **A3b** is designed here as constraints (§2.11) and an
  attended exit run (§10.3), and is built when row 188 N3 exists.
- **What A3a still delivers of the second half, honestly:** a *browser* notification (the
  Notifications API, opt-in, gradable headless with a CDP permission grant and an injected
  recorder, §2.9), and Touch ID **through the passkey path**: on the Mac, Safari's and Chrome's
  WebAuthn platform authenticator is the Touch ID prompt, in the WebAuthn format R2 already
  admits. The parent's "optional native convenience … signs in one of the two formats above, so
  the daemon learns nothing new" (parent §2.1) is therefore met in the browser by A3a. A3b's
  remaining value is the native window and a native notification.

**In A3a:**
- **Sprint 1, "the record and the minter" (M1–M3).** Credentials schema 7 adds
  `surface_approvals` (§2.2). The judge that `approve --webauthn-assertion` and `--signature` run
  moves into `host/authority` so the daemon and the CLI call the same code (§2.5). `approve <id>
  --await` is the minter (§2.6). Nothing writes the new table yet except tests.
- **Sprint 2, "the routes" (M4–M5).** `GET /surface/requests.json`, `GET /surface/requests/{id}`,
  `POST /surface/approve/{id}`, with the Host, Origin, `Sec-Fetch-Site` and Content-Type guards
  (§2.3, §2.4), all judged with `httptest` and `approvaltest` keys. An in-process end-to-end test
  POSTs a test passkey's assertion and an `approvaltest` SSHSIG, and `approve --await` mints them.
- **Sprint 3, "the page and the app" (M6–M8).** The card page `GET /surface/approve/{id}` and
  `approve.js` (§2.7); *Needs you*, the decision list and the opt-in notification in `app.js`
  (§2.8, §2.9); the live-Chrome tests XA1–XA4 on the A2 driver, extended (§6); QUICKSTART §17
  *(attended — pending first verbatim run)*.

**Out of A3a:** the Tauri app and anything native (A3b); `Human.Approve` decision packets (A4);
a request *created* from the surface (row 243's no-terminal episode, and parent §3's "an agent may
*propose* a mint request in S5"); enrolment or revocation from the surface (R3: typed-only); any
change to `/surface/passkey`, `passkey.js`, the A2 law, the typed `y`, `session *`, presets,
`--resume` or publishing.

**Sprint boundaries.** Sprint 1 adds a table, a refactor that keeps every existing CLI text, and a
CLI mode; no route can write. Sprint 2 adds the routes that write the table; no page posts yet.
Sprint 3 adds the pages. Nothing flows backwards, which follows A1b §0's and A2 §0's rule: the
law and the store are reviewed before anything can reach them.

## 1. What the code does today (measured at `0140bd7f`)

1. **Requests and their lists.** `episode start --request` stores a mint request with "no
   terminal, no worktree" (`cmd/ailang-worldd/episode.go:117`, `:211`). The store holds at most
   `MaxPendingMintRequests = 16` pending rows (`host/store/mint_requests.go:39`).
   `PendingMintRequests(ctx, now, limit)` (`:171`) and `DecidedMintRequests(ctx, approvedBy,
   limit)` (`:209`, cap 64 at `:205`) list them; `approve --list` and `--list --decided` print
   them (`approve.go:205-222`, `:337-352`) (V5).
2. **The canonical request is JSON.** `approval.Canonical` is `json.Encoder` with
   `SetEscapeHTML(false)`, trailing newline trimmed (`host/approval/request.go:165-176`). Its
   fields include the episode, worktree, repo, git binary, branch, every grant, `ttl_s`, the
   profile, `credentials_db`, `session_file`, `mcp_config`, `agent`, `model`, `task` (a path) and
   the daemon ref (`:33-55`). The 90-byte signed message is `"world-approval/v1" ‖ 0x00 ‖ digest
   ‖ challenge ‖ expiry (u64 BE)` (`:211-222`); `MessageHash` is its SHA-256 (`:225`) (V6).
   **A browser can recompute both** from the stored bytes, the challenge and the expiry.
3. **The mint lives in the CLI process.** `runApprove` (`approve.go:129-333`) admits the row,
   judges a signature or assertion "before any network call" (`:255-276`), checks the stores
   (`:278-287`), re-derives the request from the world now (`:289-317`: `prepareMint`, the MCP
   file, `fetchDaemonHealth`, `buildMintRequest`, byte-equal), then calls `provisionEpisode`
   (`:319-328`). `provisionEpisode` (`session_new.go:404`) opens the terminal only when no
   verified approval is set (`:408-418`), takes the credentials writer (`:486`), **runs `git
   worktree add`** (`:493-508`), mints (`:510-526`) and writes the token to the request's
   `session_file` (`:531-543`), rolling back on failure (`:545-570`) (V7).
4. **The daemon never keeps the credentials writer.** `ensureCredentialsStore` opens it writable
   only to create or upgrade at `serve` start, "with a short writable open that is closed again
   at once, so the daemon never keeps the credentials writer lock" (`host/daemon/credentials.go:49-80`).
   Every request-time read is `store.OpenCredentialsReadOnly` (`episode_start.go:17`,
   `binding.go:84`, `surface_page.go:199`, `surface_follow_episodes.go:68`), and
   `episode_start.go:11-12` says "No session route has a writer for attended start history" (V8).
   The writer lock is non-waiting (`host/store/credentials.go:194`), and
   `ensureCredentialsStore` already treats a busy writer as `store.IsWriterAlreadyActive`
   (`host/daemon/credentials.go:68`) (V8).
5. **The judges.** `approveAssertion` (`approve_passkey.go:75-115`): read the file capped at
   `approval.MaxPasskeyResponseBytes = 16 KiB` (`host/approval/passkey_json.go:15`), parse,
   `PasskeyByCredentialID` on a read-only handle, `ParseStoredCOSEKey`, then
   `approval.ApproveByAssertion` (`:102`). `approveArmored` (`approve_sig.go:42-78`):
   `ParseSSHSIG`, `ApproverByKey`, then `approval.ApproveBySignature`. Both live in `package main`,
   so the daemon cannot call them today (V9).
6. **The routes.** The mux has 13 `/v1` routes plus `/workbench`, `/surface*`, the agent card,
   `/a2a/`, `/mcp/` and `/agui/` (`host/daemon/daemon.go:1092-1127`).
   `TestFrozenV1RouteTableMatchesMux` pins the 13 (`route_table_test.go`).
   `/surface/passkey` and `/surface/passkey.js` are at `:1115-1116`. The only session-protected
   request is `POST /v1/commit` (`isProtected`, `:1139-1141`); every other route is open on the
   loopback bind (`DefaultBindHost = "127.0.0.1"`, `DefaultBindPort = 7644`, `:163-165`).
   `TestSurfaceAppRoutes` says "every other name under /surface/ is a 404"
   (`surface_app_test.go:15-18`) (V3, V10).
7. **The surface's headers.** `surfaceCSP = "default-src 'none'; style-src 'unsafe-inline';
   script-src 'self'; connect-src 'self'; font-src 'self'; base-uri 'none'; form-action 'self';
   frame-ancestors 'none'"` (`surface_page.go:19`); `setSurfaceHeaders` adds it to
   `setWorkbenchHeaders`' `no-store`, `nosniff` and `no-referrer` (`workbench.go:57-63`). The
   daemon sends **no CORS header** and reads no `Origin` or `Sec-Fetch-*` header in production
   code, and its only `r.Host` read is `session_self.go:69` (V10).
8. ***Needs you* has no source.** `surfacepage.Model.NeedsYou []string` says "empty in S1: nothing
   the daemon reads records a wait on you" (`host/surfacepage/model.go:56-57`); both builders set
   it to `[]string{}` (`model.go:281`, `follow_model.go:201`). `app.js`'s `nowCard` renders
   "Needs you · N" and "Nothing needs you" from it (`app/app.js:583`, `:605-606`). `app.js` is
   1,413 lines and had 5 commits since 2026-10-25 (row 188's style work) (V11). The app uses no
   `Notification` (V14).
9. **No Tauri app exists.** Only `design_docs/mockups/surface-app-v3/` (a 10-line `main.rs`,
   `tauri = "2"`, `"csp": null`), which "no test or build reads" (V12). `desktop/` does not exist;
   N3 is planned as "operator, then Claude" because the sandbox has no `cargo`
   (`w-world-os-surface-sprint-plan.md:138`). `ci.yml` (526 lines) has 0 matches for `cargo`,
   `tauri` or `WORLD_CHROME` (V13).
10. **The A2 browser harness**, which A3 reuses: `cmd/ailang-worldd/testdata/passkey_cdp.mjs`
    (172 lines) drives Chrome over CDP with commands `tab`, `open`, `ceremony`, `uv` and `advance`
    (`:120`), page injections `abort`, `late` and `eddsa`, a CSP-violation watcher, and a virtual
    authenticator (`ctap2`, `internal`). The Go side is `passkey_page_driver_test.go`
    (`startPasskeyBrowser`, gated on `WORLD_CHROME` and `node`, `SKIP COUNTED` unless
    `WORLD_REQUIRE_PASSKEY_BROWSER=1`, `:63-72`). Its tests are `TestPasskeyPageInChrome`,
    `TestPasskeyPageAbortsAPendingCeremony` and `TestPasskeyPageRecordsTheFixtures`. **They live in
    `cmd/ailang-worldd/`, not `host/daemon/` as A2 §6 planned** (V15).
11. **The A2 page** reads only its fragment, makes no request, and builds its `get` from the
    fragment's server-computed `challenge` (`passkey.js:62-64`, `:115-131`). It shows `hash16`
    for the operator to compare with the CLI: A2's residual R5, "the A2 display is
    unauthenticated" (V16).
12. **Credentials schema is 6** (`credentials.go:43`); `credentialsTables(6)` is `approvers`,
    `episode_starts`, `mint_requests`, `session_credentials`, `webauthn_credentials` (`:72-85`);
    the upgrade chain ends `if version < 6` (`:364`) (V4).
13. **Browser behaviour (V17, a throwaway probe).** Chrome 157 headless, two Node servers on
    `localhost:47811` and `localhost:47812`:
    - a same-origin `fetch` POST with `Content-Type: application/json` arrived with
      `Origin: http://localhost:47811` and `Sec-Fetch-Site: same-origin`;
    - the same POST from the other port's page was **preflighted** (`OPTIONS`, `Origin:
      http://localhost:47812`, `Sec-Fetch-Site: same-site`), the server answered without CORS
      headers, **no POST arrived**, and the page saw `TypeError: Failed to fetch`.
    So two `localhost` ports are *same-site* but not *same-origin*: a guard must require
    `same-origin`, not `same-site`.

## 2. Design

### 2.1 The shape: the surface approves, the minter mints

```
episode start --request          ─► mint_requests row (pending)               [today, A0]
approve <id> --await (no tty)    ─► polls, read-only                           [A3a, the minter]
browser: /surface → Needs you → /surface/approve/<id>
   GET /surface/requests/<id>    ─► the stored bytes, digest, challenge, expiry, passkey origins
   page: SHA-256(bytes) == digest; challenge := MessageHash(digest, challenge, expiry) computed HERE
   navigator.credentials.get     ─► Face ID / Touch ID / security key (the operator's own)
   POST /surface/approve/<id>    ─► daemon: guards, then THE SAME LAW (ApproveByAssertion /
                                    ApproveBySignature); only a verified artifact is recorded
                                 ─► surface_approvals row
approve <id> --await             ─► reads the row, RE-JUDGES it with the same law, steps 2–6,
                                    provisionEpisode with Signed set (no terminal), the mint
```

**Why the daemon does not mint (considered and rejected).** The parent's §2.4 says "The daemon
verifies the signature … It then mints". Measured, minting is `provisionEpisode`: `git worktree
add`, the credentials writer for the whole mint, a token file at `session_file` and the MCP
config, with rollback (§1.3), all in `package main`. Doing that in the daemon would add a
production exec site (the AC10 census), move `provisionEpisode` out of `package main`, keep the
credentials writer across a git subprocess, and reverse §1.4's "the daemon never keeps the
credentials writer". The parent's *verify* half is kept: the daemon verifies with the same law
before it records anything. The *mint* half runs where it runs today, in the minter. This is a
**DEPARTED-FROM-PARENT point in mechanism only**: the operator's act is still one passkey or
signature gesture on the surface, and nothing mints without that verified signature.
**An in-memory relay** (the daemon holds the artifact; the minter fetches it over a GET) was also
rejected: it is lost on restart, and it serves signed artifacts on an open read route.

**Who runs the minter.** Whoever made the request: a launcher script, the controller, a herdr
pane, or the operator. It needs no terminal (`provisionEpisode`'s V1 skips `/dev/tty` when
`Signed != nil`, §1.3). This is the same trust model as A1b/A2 today: anyone may run `approve <id>
--signature <file>`; only the signature decides, and the mint equals the request.

### 2.2 The record: credentials schema 7 (`host/store/surface_approvals.go`)

```sql
-- Approvals POSTed to the surface (version 7, row 201 A3). One row per request: the artifact the
-- daemon verified, kept for the minter, which re-judges it. Never a credential, never trusted.
CREATE TABLE surface_approvals (
    request_id  TEXT PRIMARY KEY CHECK (length(request_id) = 19 AND request_id GLOB 'mr-[0-9a-f]*'),
    format      TEXT NOT NULL CHECK (format IN ('webauthn', 'ssh')),
    approver    TEXT NOT NULL CHECK (approver <> ''),          -- the verified "<format>:<name>"
    artifact    BLOB NOT NULL CHECK (length(artifact) BETWEEN 1 AND 16384),
    digest      TEXT NOT NULL CHECK (length(digest) = 64),     -- the request's digest when recorded
    received_at INTEGER NOT NULL
)
```

- `credentialsSchemaVersion = 7`; `credentialsTables(7)` adds `surface_approvals`; the upgrade
  gains `if version < 7 { surfaceApprovalsSQL }`, the A2 pattern (one table per version, the
  canonical-DDL test compares `credentials_schema.sql` and the upgrade column for column).
- **`RecordSurfaceApproval(ctx, row SurfaceApprovalRow, now int64) error`**, one transaction:
  - **S1** `requireDeadline`; below v7 → `*OlderCredentialsStoreError{Feature: "surface approvals",
    Need: 7}`;
  - **S2** grammar: the id pattern, `format` ∈ {`webauthn`, `ssh`}, `approver` is
    `format + ":" + name`, `1 ≤ len(artifact) ≤ 16384`, `digest` 64 lowercase hex → `ErrBadSurfaceApproval`;
  - **S3** the request exists → else `ErrRequestUnknown`;
  - **S4** it is unconsumed, `created_at ≤ now < expires_at`, and its `digest` equals the row's →
    else `ErrRequestConsumed` / `ErrRequestExpired` / `ErrRequestClockBehind` /
    `ErrRequestDigestChanged` (the existing sentinels, `mint_requests.go:50-54`);
  - **S5** `INSERT … ON CONFLICT(request_id) DO UPDATE` (a later verified approval replaces an
    earlier one while the request is pending, so an artifact the minter later refuses, e.g. after a
    counter moved, can be replaced by a fresh one, §2.6 M4).
- **`SurfaceApproval(ctx, id) (SurfaceApprovalRow, bool, error)`**: read-only, deadline-required.
- **Pruning.** `PutMintRequest`'s existing prune (unconsumed requests expired a day ago,
  `mint_requests.go:95-97`) also deletes `surface_approvals` rows whose request is gone, in the same
  transaction. A consumed request's row stays, as the record of which artifact approved it.
- **Writer.** The daemon opens the store with `store.OpenCredentials` (non-waiting) per POST, runs
  `RecordSurfaceApproval`, and closes it at once (§2.4 R13). A busy writer is a 503, never a wait.
  This is the second short writable open in the daemon, beside `ensureCredentialsStore`'s; §1.4's
  "never *keeps* the writer" holds, and §7 adds an inventory test that pins exactly these two
  production sites.
- **R6 of A2 recurs:** schema 7 costs one `serve` restart, named by the existing
  `OlderCredentialsStoreError` text.

### 2.3 The read routes: *needs you*, the card's data

**Host allowlist (all five new routes, §2.3–§2.4, §2.7).** `allowedApprovalHost(r)`: the `Host`
header must be `localhost:<p>` or `127.0.0.1:<p>`, where `<p>` is the daemon's bound port, or the
exact host (and port, if any) of an **active passkey's stored origin** (`ActivePasskeys`, read
once per request on the read-only handle). Anything else → **421 Misdirected Request**, inert
JSON `{"v":1,"error":"MisdirectedRequest",…}`. This defeats DNS rebinding (an attacker's name
resolving to 127.0.0.1 sends its own `Host`). It is applied to **the new routes only**: the
existing routes stay exactly as row 199 leaves them (parent §3), and RT7 pins that
`/surface/model.json` still answers a foreign `Host`.
- Behind `tailscale serve`, whether the daemon sees the `*.ts.net` name or `127.0.0.1:7644` is
  **UNVERIFIED**. Both pass: the first because the passkey enrolled at that origin, the second
  because it is loopback. The attended run records which (§10.2).

**`GET /surface/requests.json`** (application/json, `no-store`, `nosniff`):
```json
{"v":1,"now":1794414000,
 "pending":[{"id":"mr-0123456789abcdef","episode":"w201x","agent":"claude","model":"…",
             "grants":3,"ttl_s":43200,"expires":1794414600,"hash16":"29181d224d982f65",
             "recorded":"webauthn:mark"}],
 "decided":[{"id":"mr-…","episode":"w201w","approvedBy":"webauthn:mark","at":1794410000,
             "credential16":"…"}],
 "unavailable":""}
```
- `pending` is `PendingMintRequests(now, 16)`; each summary is decoded from the stored bytes
  (`approval.Decode`). A row that does not decode, or whose bytes' SHA-256 is not its digest, is
  listed with `"refused":"<reason>"` and no summary (never hidden: *needs you* must not drop a
  request silently).
- `recorded` is the verified approver of a `surface_approvals` row, or `null`: "approved, waiting
  for the minter".
- `decided` is `DecidedMintRequests("", 20)`: the decision list (parent §2.4 step 4).
- **What it does not carry:** the request bytes, the challenge, any artifact, paths. `hash16` is
  `approve --webauthn-link`'s and the CLI's (§1.11), for comparison.
- A store below v4, absent, or unreadable → 200 with `"unavailable":"<reason>"` and empty lists, as
  `/surface/model.json` reports `sessionsUnavailable`.

**`GET /surface/requests/{id}`** (application/json): `{id}` must match `MintRequestIDPattern`
(else 400); unknown → 404.
```json
{"v":1,"id":"mr-…","request":"<base64url of the stored canonical bytes>",
 "digest":"<64 hex>","challenge":"<64 hex, the row's nonce>","created":…,"expires":…,
 "consumed":false,"approvedBy":"",
 "passkeys":[{"origin":"https://studio.tail….ts.net","rp":"studio.tail….ts.net",
              "allow":["<base64url credential id>"]}],
 "ssh":{"namespace":"world-approval@sunholo"}}
```
- It sends **the stored bytes and the raw nonce, never a precomputed message hash.** The page
  computes the digest and the message hash itself (§2.7), so what it displays is what it signs.
- `passkeys` groups active passkeys by origin, as `writePasskeyLinks` does (`approve_passkey.go:43-54`).
  Credential IDs are public WebAuthn identifiers, sent in every `allowCredentials`; they appear here
  on an open route (threat T8, §4).

### 2.4 The Approve: `POST /surface/approve/{id}`

**Request.** Body ≤ 20 KiB (`http.MaxBytesReader`), strict JSON (`DisallowUnknownFields`, exactly
one value), exactly one of `webauthn` and `sshsig`:
```json
{"v":1,"digest":"<64 hex the page displayed>","webauthn":{ …the A2 §2.7 assertion envelope… }}
{"v":1,"digest":"<64 hex>","sshsig":"-----BEGIN SSH SIGNATURE-----\n…\n-----END SSH SIGNATURE-----\n"}
```

**Checks, in order** (each refusal writes nothing and is inert JSON `{"v":1,"error":<class>,
"reason":<text>}` with `no-store`, `nosniff`):

| id | check | status, class |
|---|---|---|
| R1 | `{id}` matches `MintRequestIDPattern` | 400 `BadRequest` |
| R2 | `Host` allowed (§2.3) | 421 `MisdirectedRequest` |
| R3 | `Origin` present and equal to `http://localhost:<p>` or an active passkey's origin (for the loopback Host) — the origin whose host is the request's `Host` | 403 `CrossOrigin` |
| R4 | `Sec-Fetch-Site`, when present, is exactly `same-origin` (V17: two localhost ports are `same-site`) | 403 `CrossOrigin` |
| R5 | `Content-Type` media type is `application/json` (a `text/plain` cross-origin POST needs no preflight; this makes every browser cross-origin POST a preflighted one, which the daemon never answers with CORS, V17) | 415 `UnsupportedMediaType` |
| R6 | body size, JSON grammar, `v == 1`, exactly one format, `digest` 64 hex | 413 / 400 |
| R7 | the request row exists (read-only handle) | 404 `NotFound` |
| R8 | `approval.Admit`: unconsumed, not expired, clock not behind, stored digest = SHA-256(bytes) | 409 `Conflict` with the CLI's own refusal text (`approveAdmission`) |
| R9 | body `digest` == the row's digest ("the request changed since the card was shown") | 409 `Conflict` |
| R10 | the row's challenge is 64 hex (F4) | 409 `Conflict` |
| R11 | **the law**: `authority.JudgeAssertion` or `authority.JudgeSSHSignature` (§2.5) | 422 `NotApproved`, the law's reason |
| R12 | the verified approval's `Binds(digest, challenge, expiry)` | 422 `NotApproved` |
| R13 | `OpenCredentials` (non-waiting); busy → 503 `Busy`, `Retry-After: 1`; below v7 → 503 `Restart` naming the `serve` restart | 503 |
| R14 | `RecordSurfaceApproval` (S1–S5); a race lost to a consume or expiry → 409 | 409 |

**Success** → 200 `{"v":1,"recorded":true,"request":"mr-…","approvedBy":"webauthn:mark",
"next":"the minter (approve mr-… --await) mints it; it then appears under Decided"}`.

- **No cookie, no session, no Authorization.** The handler reads no `Cookie` and no
  `Authorization` header; `isProtected` is unchanged (`POST /v1/commit` only). A valid session
  Bearer does not approve anything (RT11).
- **`OPTIONS` and every other method** on the path → 405 by the mux's method patterns (no CORS
  response is ever written).
- **It never mints.** The handler has no path to `provisionEpisode`, `authority.MintForRequest` or
  the session table; the only write is `RecordSurfaceApproval`.
- **Placement.** `host/daemon/surface_approve.go` (new): the handlers and guards; the JSON shapes in
  `host/surfacepage/approve_model.go` (pure, tested without a daemon).

### 2.5 One judge for the daemon and the CLI (`host/authority/judge.go`)

```go
// JudgeAssertion: ParseAssertionResponse, PasskeyByCredentialID on ro, ParseStoredCOSEKey, then
// approval.ApproveByAssertion against the admitted row at now. JudgeSSHSignature: ParseSSHSIG,
// ApproverByKey, then approval.ApproveBySignature.
func JudgeAssertion(ctx context.Context, ro *store.Credentials, row store.MintRequestRow, response []byte, now int64) (*approval.SignedApproval, error)
func JudgeSSHSignature(ctx context.Context, ro *store.Credentials, row store.MintRequestRow, armored []byte, now int64) (*approval.SignedApproval, error)
```
- Typed errors (`*JudgeError{Step: "parse"|"unknown-key"|"stored-key"|"law", Err}`), so
  `approveAssertion` and `approveArmored` become thin wrappers that print **byte-identical**
  refusal texts: every WF2–WF6 and F11–F14 test stays green unedited (§6 "must stay green").
- `host/authority` already imports `host/approval` and `host/store` (`mint.go`), and `host/daemon`
  already imports `host/authority` (V9), so no new package edge and no new module.
- **Why not `host/approval`:** it is standard-library-only with no I/O (A2 §3, S2), and the judge
  reads the store.

### 2.6 The minter: `approve <request-id> --await [--await-for <duration>]`

- **F1** widens to six exclusive modes (`approve.go:161-170`); F2 (`--list` with a mode) is
  unchanged. `--await-for` without `--await` → `exitUsage`.
- After admission and F4 (as every mode), with no terminal:
  - **M1** each poll is one read-only open under a 3 s context root (`approveReadBudget`), reading
    the request row and `SurfaceApproval(id)`; polls are 1 s apart through `env.after` (the fake
    timer in tests);
  - **M2** the request consumed elsewhere → `exitUsage` with the existing "already used … by …"
    text; expired → the existing `expiredRefusal`;
  - **M3** a recorded row → `authority.JudgeAssertion`/`JudgeSSHSignature` on **its artifact**,
    against the row read in the same poll. **The row is re-judged, never trusted**: its `approver`
    and `digest` columns are display only;
  - **M4** a refused artifact → print `refused the surface's approval: <reason>; waiting for
    another` once per distinct `received_at`, and keep waiting;
  - **M5** a verified one → steps 2–6 exactly as today, with `Signed` set, so `provisionEpisode`
    opens no terminal and prints `approved by <format>:<name>`;
  - **M6** the wait ends at `min(expires_at, now + --await-for)`: `exitUsage`, "no approval for
    request <id> by <time>; it stays pending until <expiry>".
- It prints, at start: the request id, `hash16`, and `open <surface>/surface/approve/<id>`
  for each distinct active passkey origin plus `http://localhost:<port>` (no fragment, nothing
  secret).
- `TestProductionContextRoots` gains one entry (`approve.go|runApproveAwait|Background: 3`).

### 2.7 The card page: `GET /surface/approve/{id}` and `GET /surface/approve.js`

- Static embedded bytes (`host/surfacepage/approve/{index.html,approve.js}`, `go:embed`),
  `setSurfaceHeaders` and the **unchanged** `surfaceCSP`. `{id}` grammar as R1 (400 otherwise);
  the same bytes for every id. Host allowlist as §2.3.
- **The script** (text via `textContent` only; no `innerHTML`, no `eval`, no inline script):
  1. reads `{id}` from `location.pathname`; `fetch('/surface/requests/' + id)` (same origin,
     `connect-src 'self'`);
  2. **binds what it shows to what it signs:** `bytes = unb64u(request)`;
     `SHA-256(bytes) == digest` or **refuse** ("the daemon's request does not match its digest;
     nothing to approve"); parse `bytes` as JSON and **render the card from the parsed bytes**:
     episode, agent and model, every grant as a pill (scope and budget), TTL, profile, worktree,
     branch, session file; the expiry in local time; `hash16`;
  3. **computes the challenge itself**: `msg = "world-approval/v1" ‖ 0x00 ‖ digest ‖ nonce ‖
     u64be(expires)` (90 bytes), `challenge = SHA-256(msg)` with `crypto.subtle`; it never takes a
     hash from the server;
  4. **Approve with passkey** (a button, a user gesture): `navigator.credentials.get({publicKey:
     {challenge, rpId: location.hostname, allowCredentials: <the group whose origin ==
     location.origin>, userVerification: "required", timeout: 240000}, signal})` with A2's
     240000 ms AbortController bound and late-result rule (A2 §2.14 [r2]). No group for this
     origin → the button is hidden and the page says which origin to open;
  5. builds the A2 §2.7 assertion JSON itself (as `passkey.js:49-51`) and `fetch` POSTs `{v,
     digest, webauthn}` with `Content-Type: application/json`; shows the response's
     `approvedBy` and `next`, or the refusal's class and reason, then re-reads
     `/surface/requests/{id}`;
  6. **SSHSIG**: a "Save the message" link (`blob:` of the 90 bytes, `download =
     "world-approval-<id>.msg"`, the save mechanism A2 measured under this CSP, A2 V31), the
     signing line `ssh-keygen -Y sign -n world-approval@sunholo -f <key> world-approval-<id>.msg`,
     and a paste box (or a file input read with `FileReader`) whose text is POSTed as `{v, digest,
     sshsig}`.
- **The A2 page is not changed.** `/surface/passkey` and `passkey.js` stay byte-identical (PG2 pins
  their sha256): they remain enrolment's page (R3 needs the typed code) and the CLI link path.
  `approve.js` repeats about 40 lines of A2's encoding helpers rather than editing `passkey.js`.
  XA1 asserts both pages' `get` options are the same shape against the same virtual authenticator.
- **Why a page under `/surface` rather than inside `app.js`**: `app.js` is 1,413 lines under
  active style work (§1.8); a separate document keeps the code that touches a passkey small,
  auditable and driven by the A2 driver's id-based `pageState`. The app links to it (§2.8). This is
  how A3 reads A2 §2.14's "moving this page into the surface app": the approval ceremony moves from
  a CLI-printed fragment link to a page the surface app opens from *Needs you*, fed by the daemon.

### 2.8 *Needs you* and the decision list in the surface app (`app.js`)

- `nowCard`'s "Needs you" section (`app.js:605-606`) is fed by `GET /surface/requests.json`,
  fetched on the app's existing refresh tick (no new timer). Each pending request is one line:
  episode, agent, grant count, "expires in N min", a link `/surface/approve/<id>`, and
  "approved — waiting for the minter" when `recorded` is set; a `refused` row shows its reason.
  `N` jumps to the first one (the key already exists, `surface.js:257`).
- **Decided**: under it, the last few decided requests: "You · Minted · w201w · webauthn:mark ·
  14:02" (parent §2.4 step 4), and a "show all" toggling the full list of 20.
- `model.json`'s `needsYou` stays `[]` and is not repurposed (`model.go:56-57`), so no model test
  changes; A4 may use it for decision packets.
- All text through the app's existing `esc()`.

### 2.9 The opt-in browser notification

- A "Notify me" control in the Needs-you section calls `Notification.requestPermission()` (a user
  gesture). Never prompted unasked.
- On a refresh whose `pending` holds an id not seen before in this tab, with permission
  `granted`: `new Notification("World needs you", {body: "Approve mint request for <episode> (<n>
  grants), expires <hh:mm>", tag: <id>})`; `onclick` focuses the tab and opens
  `/surface/approve/<id>`. One per id per tab.
- A notification carries no authority: it is a link to the card, which shows the request from
  the daemon, and approving still needs the passkey (T9, §4).
- Secure context: `http://localhost` and the `https://*.ts.net` origin both qualify.

### 2.10 QUICKSTART §17 "Approve from the surface" *(attended — pending first verbatim run)*

`serve` on the §0 scratch store; enrol a passkey (§16); `episode start … --request`; `approve <id>
--await` in a background shell; open `/surface`, Needs you, the card, Approve; the minter mints;
the decision list shows the approver; a second Approve of the same request is refused as already
used.

### 2.11 A3b: the Tauri app's notification and Touch ID (designed as constraints; not built in A3a)

Built when row 188 N3 lands `desktop/` (a Tauri v2 shell loading `/surface` from the daemon,
`w-world-os-surface-sprint-plan.md:138`). Binding constraints for that work:
- **C1** the window loads the daemon's `/surface` (no bundled copy), so *Needs you*, the card and
  the POST are A3a's, unchanged.
- **C2** the native notification is a mirror of §2.9 (same trigger, same text, a click opens the
  card). It is never an approval.
- **C3** Touch ID signs **only** in WebAuthn (the webview's platform authenticator) or SSHSIG (an
  SSH agent such as Secretive), per R2, and reaches the daemon through `POST /surface/approve/{id}`.
  No new format, no daemon change, no key held by World.
- **C4** if the webview offers no WebAuthn, the app opens the card in the default browser instead
  of adding a native signer (simplest path; whether WKWebView in Tauri exposes WebAuthn is
  UNVERIFIED, Q1).
- **C5** the app sets a CSP no weaker than `surfaceCSP` (the fixture's `"csp": null` is a mockup
  setting, V12).

## 3. Pure core, effects at the boundary (coding standards S1–S3)

| Piece | Where | Pure? | Tested by |
|---|---|---|---|
| the request summaries, card JSON, POST body grammar, refusal classes | `host/surfacepage/approve_model.go` | yes | AM1–AM3 |
| the Host/Origin/`Sec-Fetch-Site`/Content-Type predicates | `host/daemon/surface_approve_guard.go`, pure functions of the header values and the allowed set | yes | RG1 |
| the law | `host/approval` (A1b, A2), **unchanged** | yes | existing |
| the judge (lookups + law) | `host/authority/judge.go` | store reads only | JG1–JG2 |
| schema 7, `RecordSurfaceApproval`, `SurfaceApproval`, prune | `host/store/surface_approvals.go` | effect | SA1–SA4 |
| the routes | `host/daemon/surface_approve.go` | effect (HTTP, store) | RT1–RT11 |
| the minter | `cmd/ailang-worldd/approve_await.go`, seams `env.now`, `env.after` (no new seam) | effect | AW1–AW7, E2E1 |
| the card page, the app's Needs you, the notification | `host/surfacepage/approve/*`, `app/app.js` | static assets | PG2, XA1–XA4 |

- **Why no `.ail` (S1).** A3 is HTTP guards, a table and browser code over A1b/A2's law; none of
  it has an invariant Z3 could state that is not already the law's (A2 §3: pinned AILANG has no
  ECDSA/Ed25519). `verify_ail.sh` is expected exactly green.
- **Why no package (S3).** Every consumer is host Go (`cmd/ailang-worldd`, `host/daemon`) or the
  daemon's own embedded assets.
- **S2 kept.** The law stays I/O-free; effects sit in the store, the daemon handler and the CLI.

## 4. Threats (A3)

| # | Threat | What A3 does | Pinned by |
|---|---|---|---|
| T1 | **A script or agent POSTs to the Approve route** (agents reach the loopback daemon, row 199) | A POST records only an artifact the law verified: a valid WebAuthn assertion or SSHSIG by an enrolled key over this request's message. The agent has no key and no enrolment path (R3). Garbage → 422, nothing written. A **hand-inserted row** still mints nothing: the minter re-judges (M3) | RT4, AW7 |
| T2 | **CSRF / a cross-origin POST from another tailnet or localhost page** | R3 `Origin` must be the surface's own origin; R4 `Sec-Fetch-Site: same-origin` (V17: another localhost port is `same-site`); R5 JSON only, so a browser preflights, and the daemon never answers CORS (V17: no POST arrived). Even a POST that got through needs a signature that origin's page cannot make: the passkey's rpID binds it to the enrolled origin (A2 W7, W10) | RT8, XA2 |
| T3 | **Replay of a POST body** | The artifact is bound to this request's digest, nonce and expiry (W6/B1–B4). While pending, a replay replaces the row with the same artifact (harmless). After the consume, R8 refuses ("already used"), and the store's C6 consume is once-only. A passkey's counter (W15/C18) also refuses a stale assertion | RT9, AW2 |
| T4 | **The request swapped between the card and the signing** (a modified card response, a tampered proxy) | The page renders the card **from the bytes it hashes**, refuses when SHA-256(bytes) ≠ digest, and computes the challenge from that digest, the nonce and the expiry itself. Bytes swapped consistently with a forged digest → the page signs the forged digest → W6 refuses at the daemon (422) and at the minter. The body's `digest` (R9) names the mismatch | XA2 (both arms), RT4 (R9) |
| T5 | **A request changed after the card was signed** (the world moved) | Unchanged: the minter re-derives the request and requires byte equality (`approve.go:309-317`) | existing CR tests |
| T6 | **The *needs you* list as an information leak on the open read routes** | The new routes are Host-checked (T7). The list carries summaries only (episode, agent, model, grant count, TTL, expiry, `hash16`); the card carries the request bytes (paths, grants) and passkey origins and credential IDs. Any loopback process (an agent, row 199) can read them, as it can read `/surface/model.json`'s sessions today. They carry **no token, no secret and no authority**: the nonce is not secret (it is signed over, not kept), and a credential ID is sent to every page that asks for an assertion. **Residual R1**: closed by row 199's read auth, not here (parent §3) | RT2, RT3 (no `artifact`, no token field) |
| T7 | **DNS rebinding / a forged `Host`** | The five new routes refuse any `Host` that is not the daemon's loopback name and port or an enrolled passkey's origin host (421). Existing routes unchanged (row 199) | RT7 |
| T8 | **Credential IDs on an open route** (tracking a passkey) | Public WebAuthn identifiers, scoped to World's rpID. Shown only for active passkeys. A2's Q4 (do not *commit* Mark's) is unaffected | RT3 |
| T9 | **Notification spoofing** (another site, or an agent, raises a "World needs you" notification leading to a look-alike page) | A notification carries no authority. A look-alike page on another origin cannot use World's passkey (rpID). A genuine card shows the request from the daemon, derived from the hashed bytes. **Residual R4**: an operator could approve a *genuine but unwanted* request after a spoofed nudge; the card's full request is the check, as at the typed `y` | XA4 (only opted-in, one per id), XA1 |
| T10 | **A cookie or session approves** | The handler reads no `Cookie` and no `Authorization`; `isProtected` unchanged | RT11 |
| T11 | **The minter trusts the daemon's verdict** | It does not: M3 re-judges the artifact against the row it just read. The `approver` column is display only | AW7 |
| T12 | **The daemon's writer blocks a mint, or the reverse** | Non-waiting writer, opened per POST and closed at once; busy → 503 `Retry-After: 1`, the page retries 3 times then says so | RT10 |
| T13 | **An injected script on the surface origin swaps the challenge** | `surfaceCSP` `script-src 'self'`, no inline script, no `innerHTML` in `approve.js` (PG2's labelled controls); the page shows the request it signs. A2's residual R5 (the unauthenticated fragment display) is **closed** for the surface path: the card comes from the daemon and is bound to its digest | PG2, XA2 |
| T14 | **Oversize or malformed bodies** | 20 KiB `MaxBytesReader`; strict JSON; the A2 envelope's own 16 KiB bound and parser | RT4 |
| T15 | **A process that can rewrite the daemon binary, the pages or `credentials.db`** | Out of scope, unchanged (A2 §4) | — |

## 5. By-construction argument

1. **Nothing mints without a verified signature.** The only mint paths are `provisionEpisode`'s.
   With no typed `y`, it needs `Signed != nil` (`session_new.go:408`), set only from
   `JudgeAssertion`/`JudgeSSHSignature`, i.e. from `ApproveByAssertion`/`ApproveBySignature`. The
   minter calls them on the artifact itself (M3). The POST handler has no mint call at all.
2. **What the operator saw is what was signed.** The page renders from bytes `b` with
   SHA-256(`b`) = `d`, and signs `MessageHash(d, nonce, expiry)`. W6 accepts only the stored
   row's `MessageHash`, so a signature over a displayed request that differs from the stored one
   is refused.
3. **Approving is a POST bound to a verified signature, never a cookie or a session** (parent
   §3). The POST's only write follows R11–R12. No cookie is set or read anywhere in A3 (PG2's
   control greps `Set-Cookie`/`document.cookie`).
4. **Read routes stay as open as row 199 leaves them.** No existing route's handler, middleware or
   predicate changes; the new routes add a Host check that only narrows.
5. **The agent never approves; publishing stays typed-only; the typed `y` keeps working.** No key,
   no enrolment (R3, untouched); `cmd/world-publish` untouched; `approve <id>` with no mode still
   asks `y` at `/dev/tty` (unchanged).
6. **World never holds a private key.** The browser's authenticator or the operator's SSH key
   signs.

## 6. Test plan

Every browser test uses the A2 harness (`passkey_page_driver_test.go`, `passkey_cdp.mjs`),
extended, never replaced. Keys are `approvaltest` keys or Chrome's virtual authenticator.

### Sprint 1: the record and the minter

| Id | Test | Asserts | Mutants killed |
|---|---|---|---|
| SA1 | `TestCredentialsSchema7SurfaceApprovals` (`host/store`) | v6 → v7 upgrade; `credentialsTables(7)`; the DDL equals `credentials_schema.sql`'s; v6 stores refuse `RecordSurfaceApproval` with the restart text | the upgrade step deleted; the table set unchanged |
| SA2 | `TestRecordSurfaceApprovalGuards` | S2 each arm; S3 unknown; S4 consumed / expired / clock-behind / digest-changed → the named sentinel and **0 rows**; S5 replace while pending | each S check deleted; `INSERT` without the request join |
| SA3 | `TestSurfaceApprovalPrunedWithItsRequest` | the prune removes orphans; a consumed request's row stays | prune not extended |
| SA4 | inventory edits | `TestCredentialsMethodsRequireADeadline` +2; `TestNoNewDeadlineFreeStoreReads` | a read without a deadline |
| JG1 | `TestJudgeAssertionIsTheCLIJudge` (`host/authority`) | over the A2 refusal matrix (`approvaltest.Passkey` knobs), `JudgeAssertion`'s `Step` and the CLI wrapper's text agree; a verified one `Binds` | a wrapper that skips `ParseStoredCOSEKey`; a judge that skips the lookup |
| JG2 | `TestJudgeSSHSignatureIsTheCLIJudge` | the same for SSHSIG (F11–F14) | same |
| AW1 | `TestApproveAwaitMintsARecordedPasskeyApproval` (`cmd/ailang-worldd`) | a row recorded mid-wait → one session, `approved_by = webauthn:<name>`, the token in `session_file`, `openTerminal` never called | the wait not polling; `Signed` unset (a terminal opens) |
| AW2 | `TestApproveAwaitRejudgesAndWaits` | a recorded artifact whose counter regressed → the refusal printed once, still waiting; a replacement → mints | refusal ends the wait; refusal printed every poll |
| AW3 | `TestApproveAwaitStopsAtExpiry` | fake clock: ends exactly at `expires_at` (and at `--await-for`), `exitUsage`, 0 sessions | an unbounded wait; `>` vs `≥` |
| AW4 | `TestApproveAwaitStopsWhenUsedElsewhere` | a consume by another process mid-wait → "already used", `exitUsage` | the consumed check deleted |
| AW5 | `TestApproveAwaitFlagRules` | F1 six modes; `--await` with each other mode; `--await-for` alone; `--list --await` | F1 not widened |
| AW6 | `TestApproveAwaitSSHSIG` | a recorded `approvaltest` SSHSIG mints `ssh:<name>` | the format switch wrong |
| **AW7** | `TestApproveAwaitNeverTrustsTheRow` | a row **inserted directly** (no daemon) with (a) an unenrolled key's valid signature, (b) a passkey assertion over another request, (c) `approver = "webauthn:mark"` over a garbage artifact → each refused, **0 sessions** | **the minter trusting the row** (minting from the `approver` column) |
| edited | `TestEveryVerbHasHelp`, `TestProductionContextRoots` | the `--await` usage line; +1 root | — |

### Sprint 2: the routes

| Id | Test | Asserts | Mutants killed |
|---|---|---|---|
| AM1–AM3 | `TestApproveModel*` (`host/surfacepage`) | summaries from bytes; a non-decoding row → `refused`; the POST body grammar (both/neither format, unknown field, trailing value, bad digest) | each grammar arm |
| RG1 | `TestApprovalGuards` (pure) | Host, Origin, `Sec-Fetch-Site`, media type tables, incl. `same-site`, `null`, `localhost:<other port>`, upper-case hosts, `application/json; charset=utf-8` accepted, `text/plain` refused | `same-site` accepted; prefix-matching hosts |
| RT1 | `TestApprovalRoutesRegistered` (`host/daemon`) | the five routes; `/v1` still 13; `TestSurfaceAppRoutes` edited (+`approve.js`, `approve/{id}`, `requests.json`, `requests/{id}`) | a route missing |
| RT2 | `TestRequestsJSON` | pending/decided shapes and caps; `recorded`; `unavailable` on a missing store; **no `request`, `challenge`, `artifact` or token key** | leaking the artifact |
| RT3 | `TestRequestCardJSON` | `request` decodes to exactly the stored bytes; digest/challenge/expiry are the row's; **no message-hash field**; passkey groups by origin | a precomputed hash sent |
| RT4 | `TestApprovePOSTRefusalArms` | R1–R14 each → its status/class, **0 `surface_approvals` rows**, 0 sessions | each R deleted |
| RT5 | `TestApprovePOSTRecordsAVerifiedPasskey` | an `approvaltest.Passkey` assertion → 200, one row, `approver = webauthn:<name>`, 0 sessions | recording before the law |
| RT6 | `TestApprovePOSTRecordsAVerifiedSSHSIG` | the same for SSHSIG | — |
| RT7 | `TestApprovalRoutesRefuseForeignHost` | `Host: rebind.example:7644` → 421 on all five; loopback and the enrolled `*.ts.net` host pass; **control:** `/surface/model.json` with the foreign Host still 200 (row 199 posture) | the check missing on one route; applied to an existing route |
| RT8 | `TestApprovePOSTCrossOrigin` | no `Origin`, `Origin: null`, another port, `Sec-Fetch-Site: same-site`/`cross-site`, `text/plain` → 403/415; `OPTIONS` → 405 with no `Access-Control-*` header | — |
| RT9 | `TestApprovePOSTReplay` | the same body twice while pending → 200, one row; after the minter consumed → 409 "already used" | — |
| RT10 | `TestApprovePOSTWriterBusy` | a held writer → 503, `Retry-After: 1`, 0 rows; the daemon holds no writer after any response | a waiting open; a kept writer |
| RT11 | `TestApprovePOSTNoSessionNoCookie` | a valid session Bearer and a cookie with no signature → refused (R6/R11); `isProtected` unchanged | — |
| RT12 | `TestDaemonCredentialWriterSites` (inventory, new) | exactly two production `store.OpenCredentials` call sites in `host/daemon`: `ensureCredentialsStore` and the POST | a third writer |
| E2E1 | `TestSurfaceApproveThenAwaitMints` (`cmd/ailang-worldd`) | in process: a daemon `Handler()` on temp stores; `episode start --request`; POST a test passkey's assertion; `approve --await` mints; `requests.json` lists it under `decided` | the end-to-end seam |

### Sprint 3: the page and the app (live Chrome, gated as XB1)

| Id | Test | Asserts | Mutants killed |
|---|---|---|---|
| PG2 | `TestApprovePageRoutes` (`host/daemon`) | the embedded bytes, `surfaceCSP` unchanged, 400 for a bad id; `passkey.js`/`index.html` sha256 unchanged. Labelled controls (not load-bearing): `approve.js` has no `innerHTML`, `eval(`, `document.cookie` | — |
| XA1 | `TestApprovePageInChrome` (`cmd/ailang-worldd`) | virtual authenticator; enrol through the A2 page + CLI (XB1's helpers); `--request`; open `/surface/approve/<id>`; the card's grants equal the request's; click Approve; 200 recorded; `approve --await` mints; zero CSP violations, zero console errors; the A2 page's `get` for the same request carries the same `challenge` bytes | the page using another challenge |
| **XA2** | `TestApprovePageRefusesASwappedRequest` | CDP `Fetch.enable` rewrites the card response: (a) bytes changed, digest kept → the page refuses, **no `credentials.get` call** (injected recorder); (b) bytes and digest both forged → the page signs, the POST is **422**, 0 rows | **`approve.js` using a server hash** or skipping the digest check |
| XA3 | `TestNeedsYouListsPendingRequests` | `/surface` shows "Needs you · 1", the link opens the card; after the mint, Decided shows `webauthn:<name>` | the list not wired |
| XA4 | `TestNeedsYouNotifiesOnlyWhenAskedOnce` | with no permission: zero notifications; after `Browser.grantPermissions` + the opt-in click, a new request → exactly one notification (tag = id, body names the episode); the next refresh → none | a notification per refresh; notifying without opt-in |
| edited | the driver | `passkey_cdp.mjs` gains commands `fetchRewrite`, `grantNotifications`, `click`, and the `notify`/`getRecorder` injections; `pageState` stays for A2 tests | — |

**Must stay green, unedited:** every a0a1, A1b and A2 test (the CLI refusal texts are kept by
§2.5), `TestDaemonDependencyAllowlist`, the quickstart tests, and every existing route test,
except the named edits (`TestSurfaceAppRoutes`, the schema pins, help and context roots).

## 7. Repo-wide inventories

| Inventory test | Change |
|---|---|
| `TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential` | **no new production exec site** (the minter reuses `provisionEpisode`'s existing git call; Chrome/Node only in `_test.go`) |
| `TestDaemonDependencyAllowlist` | no new module |
| `TestWorldBoundaryDependencyAllowlist`, `TestEveryCommandDirectoryIsAProtectedGroup` | no new transport import in `cmd/ailang-worldd` or `host/store`; no new `cmd/` directory |
| `TestProductionContextRoots`, `TestProductionGoSurface` | +1 root (`runApproveAwait`) |
| `TestNoNewDeadlineFreeStoreReads`, `TestCredentialsMethodsRequireADeadline` | +2 store methods, deadline-required |
| `TestFrozenV1RouteTableMatchesMux` | **no `/v1` route**; stays 13 |
| `TestSurfaceAppRoutes` | +4 names (edited) |
| `TestDaemonCredentialWriterSites` | **new** (RT12) |
| `TestApprovalImportsOnlyTheStandardLibrary` | unchanged: `host/approval` not touched |
| `TestNoParallelWriteForkPackagesOutsideVerifygate` | no `t.Parallel` added |
| tracked-binary gate | no binaries (pages are text) |

**A3a adds:** no exec site, **no new Go package** (`host/surfacepage` gains an embed directory;
`host/authority` a file), no module, one schema version, five `/surface` routes, no `/v1` route.

## 8. Conflict surface

| Sprint | File | Change |
|---|---|---|
| 1 | `host/store/credentials_schema.sql`, `credentials.go`, `surface_approvals.go` (new), `mint_requests.go` (prune) + tests | schema 7 |
| 1 | `host/authority/judge.go` (new) + test | the shared judge |
| 1 | `cmd/ailang-worldd/approve_passkey.go`, `approve_sig.go` (wrappers), `approve.go` (F1, the mode), `approve_await.go` (new), help | the minter |
| 2 | `host/daemon/daemon.go` (+5 `HandleFunc`), `surface_approve.go`, `surface_approve_guard.go` (new), `host/surfacepage/approve_model.go` (new) + tests | the routes |
| 3 | `host/surfacepage/approve/{index.html,approve.js}`, `approve.go` (embed) (new); `host/surfacepage/app/app.js` (`nowCard`'s Needs-you block, the refresh tick, the notification) | the pages |
| 3 | `cmd/ailang-worldd/testdata/passkey_cdp.mjs`, `passkey_page_driver_test.go` (extended), `approve_page_*_test.go` (new) | the browser tests |
| 3 | `docs/QUICKSTART.md` | §17 |

**Live conflict risk:** `app.js` (row 188's style work, 5 commits since 2026-10-25, V11): A3 edits
only the Needs-you block and adds one fetch to the existing refresh. `daemon.go`'s route block
(every surface row adds lines there). PR #552 (A2 sprint 4) must merge first: A3 builds on
`passkey_cdp.mjs` and the driver it adds.

**Not touched:** `host/approval/*` (the law), `host/surfacepage/passkey/*`, `world/*.ail`,
`packages/`, `tools/launchd/*`, any `/v1` route, `cmd/world-publish/*` (R4), the typed-`y` path,
`session *`, presets, `--resume`, `design_docs/mockups/*`.

## 9. Milestones, sprints and gates

Every milestone ends with `go build ./... && go vet ./<pkg>/` (the compile fence for test files)
plus its own tests. Every sprint ends with `./scripts/verify_ail.sh`;
`AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh`; `go build ./... && go test ./...`;
and the §7 inventories by name, as A2 §9 lists them, plus `go test ./host/daemon/ -run
'^TestDaemonCredentialWriterSites$' -count=1` from sprint 2.

**Sprint 1 = M1–M3, "the record and the minter", about 1,650 lines.**
- **M1 schema 7** (≈ 180 + 420): SA1–SA4.
  `go test ./host/store/ -run 'Credentials|SurfaceApproval|OlderCredentials' -count=1`.
- **M2 the shared judge** (≈ 140 production incl. the wrapper rewrite + 300 test): JG1, JG2; every
  `TestApprove*` and `TestPasskey*` CLI test unedited.
  `go test ./host/authority/ ./cmd/ailang-worldd/ -run 'Judge|TestApprove|TestPasskey' -count=1`.
- **M3 `approve --await`** (≈ 170 + 440): AW1–AW7, help, roots.
  `go test ./cmd/ailang-worldd/ -run '^(TestApproveAwait|TestEveryVerbHasHelp)' -count=1`;
  `go test ./host/store/ -run '^TestProductionContextRoots$' -count=1`.
- **Fallback cut:** if the plan exceeds 2,000 lines, M3 moves to the head of sprint 2 (the table
  then has no writer and no reader outside tests, which is safe).

**Sprint 2 = M4–M5, "the routes", about 1,700 lines.**
- **M4 the read routes and the guards** (≈ 260 + 520): AM1–AM3, RG1, RT1–RT3, RT7.
- **M5 the POST** (≈ 200 + 720): RT4–RT6, RT8–RT12, E2E1.
  `go test ./host/surfacepage/ ./host/daemon/ -run 'ApproveModel|ApprovalGuards|Approval|RequestsJSON|RequestCard|ApprovePOST|SurfaceAppRoutes|FrozenV1|DaemonCredentialWriter' -count=1`;
  `go test ./cmd/ailang-worldd/ -run '^TestSurfaceApproveThenAwaitMints$' -count=1`.
- **Fallback cut:** the SSHSIG arm of the POST (R11's `sshsig`, RT6, AW6's route half) moves to
  sprint 3; the passkey POST ships alone. The guards (R2–R5) are never cut.

**Sprint 3 = M6–M8, "the page and the app", about 1,750 lines.**
- **M6 the card page** (≈ 300 page + 60 Go): PG2.
- **M7 Needs you, Decided, the notification** (≈ 160 `app.js`).
- **M8 the live browser and QUICKSTART** (≈ 220 driver + 700 tests + 60 docs): XA1–XA4.
  `WORLD_CHROME="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" WORLD_REQUIRE_PASSKEY_BROWSER=1 go test ./cmd/ailang-worldd/ -run '^(TestApprovePage|TestNeedsYou|TestPasskeyPage)' -count=1 -v`
  must print `--- PASS` for every subtest **on the controller's machine** (CI has no Chrome step,
  V13; R3).
- **Fallback cut:** M7's notification and Decided list (XA4, half of XA3) become a fourth sprint;
  the card, Approve and the Needs-you link (XA1–XA2, XA3's list) stay. A3a's exit run waits for
  the fourth sprint only for its notification step.

**A3b: no sprint in this design.** It becomes a child design when `desktop/` exists (row 188 N3).

## 10. The attended exit run (Mark; NOT run by the loop, D-WORLD-109)

On a scratch store, with his own devices and passkey (enrolled in A2's run, §10 there):
1. **Serve and expose.** `ailang-worldd serve` on the scratch store; `tailscale serve` HTTPS in
   front of `127.0.0.1:7644` (tailnet only, never Funnel).
2. **Record the Host.** Open `/surface` from the phone. Record which `Host` the daemon saw (the
   `*.ts.net` name or `127.0.0.1:7644`; §2.3 UNVERIFIED). Both are designed to pass.
3. **Approve from the phone.** In a shell (or a herdr pane), `episode start … --request` then
   `approve <id> --await`. On the phone: Needs you · 1 → the card → check grants, TTL, worktree →
   Approve → Face ID. The minter mints with no World prompt. Decided shows `webauthn:<name>`.
4. **On the Mac.** The same with Safari or Chrome and Touch ID (the platform authenticator); turn
   on "Notify me" and see one notification for a new request.
5. **SSHSIG from the card.** Save the message, `ssh-keygen -Y sign … `, paste; the minter mints
   `ssh:<name>`.
6. **Refusals.** Approve the same request again (refused: already used); let one expire (the
   minter exits naming the expiry).
7. **A3b (only if `desktop/` exists by then).** Open the Tauri window; check the notification
   mirror and whether the webview runs WebAuthn (Q1).
8. **Record** public values only: browsers, the Host seen, notification behaviour per browser, and
   any refusal text that confused.

QUICKSTART §17 is executed verbatim in that run, and its *attended — pending* tag is removed.

## 11. Residuals and risks

- **R1 Open read routes.** Pending requests' summaries and cards are readable by any loopback
  process (and, behind `tailscale serve`, the tailnet), like the existing surface routes. No
  secret or authority is exposed (T6); row 199 closes it.
- **R2 The minter must be running.** A Face ID approval with no minter waiting is recorded and
  shown as "waiting for the minter" until the request expires (default 10 min). The UX cost is
  real; the alternative (the daemon mints) was rejected in §2.1 with measurements.
- **R3 Browser tests are local-only** (as A2's R8): CI has no Chrome. CI covers sprint 1–2 fully.
- **R4 A genuine but unwanted request after a spoofed nudge** (T9). The card is the check.
- **R5 Safari and Firefox** behaviour for `fetch` POST headers, `crypto.subtle` and
  Notifications is UNVERIFIED (Chrome 157 measured, V17); §10.4 checks Safari.
- **R6 Schema 7 costs one `serve` restart.**
- **R7 Size.** A2 ran 1.0–2.6× its sprint estimates (V20). The cuts in §9 are pre-named.
- **R8 A3b's host is not scheduled.** Row 188 N3 is "operator, then Claude"; until it lands,
  the "native" half of the parent's A3 sentence is met only by the browser (§0).
- **R9 tailscale `Host` behaviour** UNVERIFIED (§2.3); designed to pass either way.

## 12. Open questions

**None blocks the loop.** Each has a default the loop builds on.
- **Q1** A3b: is a native signer in the Tauri app wanted at all, given that Touch ID already works
  through the browser's passkey (§0)? **Default: no native signer; the app mirrors notifications
  and opens the card (C4), when N3 exists.**
- **Q2** Should `episode start --request` gain `--await` (store, then wait, in one command)?
  **Default: no in A3** (two commands; one new mode only). A small follow-up if the exit run wants
  it.
- **Q3** The decision list shows 20 decided requests. **Default: 20** (the store's cap is 64).

## 13. What this is NOT doing

- No Tauri, Rust, native notification or native Touch ID code (A3b). No `Human.Approve` (A4).
- No daemon-side mint, no new exec site, no kept credentials writer.
- No change to the WebAuthn or SSHSIG law, to `/surface/passkey`/`passkey.js`, or to enrolment
  and revocation (typed-only, R3).
- No request created from the surface; no agent proposal path (S5).
- No change to any existing route's openness, no `/v1` route, no cookie, no session, no CORS.
- No `.ail`, no package, no module.
- No real key, no real store, no mint in one, no publish, by the loop or its tests (D-WORLD-109).

## 14. Verification Log

Commands ran in the worktree at `0140bd7f` unless stated otherwise. Empty results are paired with
a positive control in the same or an adjacent command.

| # | Claim | Command | Observed |
|---|---|---|---|
| V1 | Base and file hashes | `git rev-parse HEAD`; `git branch --show-current`; `shasum -a 256 host/daemon/daemon.go host/daemon/surface_page.go host/surfacepage/model.go host/surfacepage/app/app.js host/surfacepage/passkey/passkey.js cmd/ailang-worldd/approve.go cmd/ailang-worldd/approve_passkey.go cmd/ailang-worldd/approve_sig.go host/store/mint_requests.go host/store/credentials.go cmd/ailang-worldd/session_new.go cmd/ailang-worldd/testdata/passkey_cdp.mjs \| awk '{print substr($1,1,16), $2}'`; `go version` | `0140bd7fb26be0fa9a93f367b9dde7cd6456512d`; `design/world-iter383-row201-a3`; the prefixes in the header; `go1.26.6 darwin/arm64` |
| V2 | The parent's A3 sentence and §3 invariants | `sed -n 60,93p design_docs/planned/w-world-approvals-by-key.md` (read in full, 93 lines) | `:84 - **A3:** the surface flow: *needs you*, the request card and Approve, plus a native notification and the optional Touch ID convenience in the Tauri app.`; `:63-66` the four invariants; `:31` "Optional native convenience … signs in one of the two formats above" |
| V3 | The route block and `isProtected` | `grep -n "HandleFunc\|mux.Handle" host/daemon/*.go \| grep -v _test`; `grep -n "func (d \*Daemon) isProtected" -A3 host/daemon/daemon.go` | `daemon.go:1092-1104` 13 `/v1`; `:1107-1119` surface (incl. `:1115 GET /surface/passkey`, `:1116 GET /surface/passkey.js`); `:1124-1127` agent card, a2a, mcp, agui; `:1139-1141 return r.Method == http.MethodPost && r.URL.Path == "/v1/commit"` |
| V4 | Credentials schema 6 and its table set | `grep -n "credentialsSchemaVersion = \|case version >= \|if version < " host/store/credentials.go`; `sed -n 70,90p host/store/credentials.go` | `:43 const credentialsSchemaVersion = 6`; `:74 case version >= 6: return []string{approversTable, "episode_starts", mintRequestsTable, credentialsTable, passkeysTable}`; `:364 if version < 6 {` |
| V5 | Request listing | `sed -n 1,95p host/store/mint_requests.go`; `sed -n 165,240p host/store/mint_requests.go`; `cat -n cmd/ailang-worldd/approve.go \| sed -n 100,355p` | `:39 MaxPendingMintRequests = 16`; `:171 PendingMintRequests`; `:205 MaxDecidedMintRequests = 64`; `:209 DecidedMintRequests`; `approve.go:205-222` list, `:337-352` decided |
| V6 | Canonical JSON, the signed message | `grep -n "^func \|^type MintRequest struct" host/approval/request.go`; `sed -n '/^type MintRequest struct/,/^}/p'`; `sed -n '/^func Canonical/,/^}/p'`; `sed -n '/^func SignedMessage/,/^}/p'` | `:33 type MintRequest struct` (fields incl. `session_file`, `task`, `daemon`); `:165 Canonical` = `json.NewEncoder`, `SetEscapeHTML(false)`, `TrimSuffix "\n"`; `:211 SignedMessage` 90 bytes; `:225 MessageHash` |
| V7 | **The mint is in the CLI process (premise of §2.1)** | `cat -n cmd/ailang-worldd/approve.go \| sed -n 100,355p`; `grep -n "^func provisionEpisode" -A80 cmd/ailang-worldd/*.go`; `sed -n 484,600p cmd/ailang-worldd/session_new.go` | `approve.go:255-276` judged "before any network call"; `:301 fetchDaemonHealth`; `:326 provisionEpisode(req, …)`; `session_new.go:404 func provisionEpisode`; `:410 if !signed { // V1` opens the terminal; `:486 openCredentialsForSession`; `:498/:500 worktree.AddExisting/worktree.Add`; `:519 authority.MintForRequest`; `:537 writeTokenFile(r.outPath, rawToken)`; `:565 worktree.Remove` rollback |
| V8 | **The daemon never keeps the credentials writer** | `grep -rn "OpenCredentials\|openCredentials" --include='*.go' host/daemon cmd/ailang-worldd \| grep -v _test`; `cat -n host/daemon/credentials.go`; `sed -n 1,26p host/daemon/episode_start.go`; `grep -n "writer lock\|NON-WAITING\|func OpenCredentials" host/store/credentials.go` | in `host/daemon` only `OpenCredentialsReadOnly` (`episode_start.go:17`, `credentials.go:60`, `surface_follow_episodes.go:68`, `binding.go:84`, `surface_page.go:199`) plus `store.UpgradeCredentials` at `credentials.go:66`; comment `:52-53` "so the daemon never keeps the credentials writer lock"; `episode_start.go:11-12` "No session route has a writer for attended start history"; `store/credentials.go:194` "NON-WAITING writer lock first" |
| V9 | The judges, and the import edges | `cat -n cmd/ailang-worldd/approve_passkey.go`; `sed -n 42,80p cmd/ailang-worldd/approve_sig.go`; `go list -f '{{join .Imports "\n"}}' ./host/daemon \| grep sunholo`; `grep -rln 'host/approval"' host/daemon host/authority` | `approve_passkey.go:75 approveAssertion`, `:102 approval.ApproveByAssertion`; `approve_sig.go:42 approveArmored`, `ApproveBySignature`; `host/daemon` imports `host/authority`, `host/store`, `host/surfacepage` (and 16 others); `host/authority/mint.go` imports `host/approval`; `host/daemon` does not |
| V10 | Surface CSP and headers; no CORS, Origin or `Sec-Fetch` handling; bind | `cat -n host/daemon/surface_page.go \| head -40`; `grep -n "func setWorkbenchHeaders" -A7 host/daemon/workbench.go`; `grep -rn "Access-Control\|Sec-Fetch\|\"Origin\"" --include='*.go' host/daemon \| grep -v _test`; `grep -n "DefaultBindHost =\|r.Host\b" host/daemon/*.go \| grep -v _test`; `grep -n "DefaultBindPort\|BindPort int" host/daemon/daemon.go` | `surface_page.go:19 const surfaceCSP = "default-src 'none'; style-src 'unsafe-inline'; script-src 'self'; connect-src 'self'; font-src 'self'; base-uri 'none'; form-action 'self'; frame-ancestors 'none'"`; `workbench.go:57-63`; the CORS/Origin grep: **no match**; `daemon.go:163 DefaultBindHost = "127.0.0.1"`, `:165 DefaultBindPort = 7644`, `:231 BindPort int`; `session_self.go:69 addr := "http://" + r.Host` (the only `r.Host`) |
| V11 | *Needs you* today; `app.js` size and churn | `grep -rn "needs you\|needsYou\|needs-you\|needs_you" --include='*.go' --include='*.js' --include='*.html' host cmd \| grep -v _test`; `sed -n 50,60p host/surfacepage/model.go`; `sed -n 578,612p host/surfacepage/app/app.js`; `wc -l host/surfacepage/app/app.js`; `git log --oneline --since=2026-10-25 -- host/surfacepage/app/app.js \| wc -l` | `model.go:56-57` "NeedsYou is empty in S1: nothing the daemon reads records a wait on you." `NeedsYou []string`; `model.go:281`, `follow_model.go:201` set `[]string{}`; `app.js:583 needs = ne.needsYou \|\| []`, `:605-606 Needs you · … Nothing needs you`; `surface.js:257` key `N`; `1413`; `5` |
| V12 | **No Tauri app in the repo** | `grep -rli tauri . \| grep -v '^./.git/'`; `ls -d desktop`; `ls -la design_docs/mockups/surface-app-v3`; `cat …/README.md`; `cat …/Cargo.toml …/tauri.conf.json`; `wc -l …/src/main.rs`; `grep -n "no test or build reads it" …/README.md` | Tauri only under `design_docs/` (mockup, plans, archives); `ls: desktop: No such file or directory`; `main.rs` 10 lines; `tauri = { version = "2" }`; `"security": { "csp": null }`; README:3 "It is not product code, and no test or build reads it." |
| V13 | N3 unbuilt; CI has no Rust/Tauri/Chrome | `grep -n "cargo\|tauri" design_docs/planned/w-world-os-surface-sprint-plan.md`; `grep -n "^188\. " design_docs/world-mission.md`; `grep -c "" .github/workflows/ci.yml`; `grep -ciE "cargo\|tauri\|WORLD_CHROME" .github/workflows/ci.yml` | plan `:138` N3 "(`desktop/`) … The World sandbox has no `cargo`, so the operator builds it"; `526`; `0` |
| V14 | The app uses no Notifications API | `grep -c Notification host/surfacepage/app/app.js`; `grep -n Notification host/surfacepage/app/app.js host/surfacepage/surface.js` | `0`; no match |
| V15 | **The A2 harness exists, where, and passes at the base** | `ls cmd/ailang-worldd/testdata/passkey_cdp.mjs host/daemon/testdata`; `grep -rln "passkey_cdp\|addVirtualAuthenticator" host cmd`; `grep -rn "func TestPasskeyPage" cmd/ailang-worldd/*_test.go`; `WORLD_CHROME="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" WORLD_REQUIRE_PASSKEY_BROWSER=1 go test ./cmd/ailang-worldd/ -run '^TestPasskeyPageInChrome$' -count=1 -timeout 240s -v` | driver at `cmd/ailang-worldd/testdata/passkey_cdp.mjs`; `host/daemon/testdata` holds no `.mjs`; tests `passkey_page_chrome_test.go:25`, `passkey_page_abort_test.go:18`, `passkey_fixture_record_test.go:23`; `--- PASS: TestPasskeyPageInChrome (4.13s)` with subtests `enrol-then-approve`, `user-verification-required`, `links-for-another-page`; `ok … 4.522s` |
| V16 | The A2 page signs a fragment-supplied challenge | `cat -n host/surfacepage/passkey/passkey.js` | `:62-64 navigator.credentials.get({ publicKey: { challenge, rpId: f.rp, allowCredentials: f.allow…`; `:54 const challenge = unb64u(f.challenge)`; `:117` the fragment parse; `:124` the origin check |
| V17 | **Chrome's POST headers and the preflight** (throwaway probe; quorum may re-run it) | `/tmp/w383-designer/hdr/probe.mjs` (two Node servers on `127.0.0.1:47811/47812`; Chrome `--headless=new --virtual-time-budget=3000 --dump-dom` on a page that `fetch`-POSTs `{"v":1}` with `Content-Type: application/json` to its own origin, and on a page at the other port posting cross-origin), run with `WORLD_CHROME=… node /tmp/w383-designer/hdr/probe.mjs` | `same-origin page result: ok posted`; `cross-origin page result: err TypeError: Failed to fetch`; server log: `{"port":47811,"method":"POST","origin":"http://localhost:47811","sfs":"same-origin","ct":"application/json"}`, then `{"port":47811,"method":"OPTIONS","origin":"http://localhost:47812","sfs":"same-site","ct":null}`; **no cross-origin POST arrived** |
| V18 | Route, asset and dependency pins green at the base | `go test ./host/daemon/ -run '^(TestFrozenV1RouteTableMatchesMux\|TestSurfaceAppRoutes\|TestPasskeyPageRoutes\|TestDaemonDependencyAllowlist)$' -count=1 -timeout 270s -v` | `--- PASS` ×4; `ok …/host/daemon 0.549s` |
| V19 | Decision-ledger rulings | `grep -n "\| D-WORLD-88 \|\\\|\| D-WORLD-89 \|\\\|\| D-WORLD-109 \|" design_docs/world-mission.md` | `:1007` D-WORLD-88 R1 = Tauri v2 ("Native desktop shell now"); `:1008` D-WORLD-89 "all defaults", R1–R5 = YES; `:1028` D-WORLD-109 "The loop builds and judges against fakes and test keys, and never enrols a real key, mints or publishes." ANSWERED A, 2026-11-06 |
| V20 | A2's measured sprint sizes (calibration) | `for m in d5d90e49 3e55451f 0ce117c5; do git diff --numstat $m^1 $m -- . ':!design_docs' \| awk '{a+=$1+$2} END {print a}'; done`; `git diff --numstat ab4b2402 0140bd7f -- . ':!design_docs' \| awk …`; `git log --oneline ab4b2402..0140bd7f` | `1969`, `1843`, `1686`; sprint 4 `2768`; its plan commit `d688d22d` says "7 compiled patches, 1,054 lines" |
| V21 | The daemon's port is configurable and may be ephemeral (the Host allowlist must use the bound port) | `grep -n "BindPort\|Port \+int\|BindHost \+string" host/daemon/daemon.go` | `:164-165 DefaultBindPort = 7644`; `:228 BindHost string`; `:229-231` "BindPort is the TCP port; 0 asks the kernel for an ephemeral port, whose resolved value is announced on stdout" |
| V22 | Probe stayed out of the repo | `git status --short` (worktree) after writing this doc | only this doc (the controller re-checks) |

## 15. Quorum verification log

*(Empty: the controller fills it.)*
