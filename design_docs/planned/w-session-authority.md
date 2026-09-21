# w-session-authority — the inbound `Authorization: Bearer` credential→(episode, grants) boundary the repo has never had, minted and resolved locally-first and failing closed under clause 3

- Status: **planned** · Date: 2026-09-21 · Designer: pi:ollama/deepseek-v4-flash:0731-cloud (iter-176) · Base commit: b65a9ecf53110a50b79515e8f1afc216e2a99928 · Owning queue row: 39 · Scope class: product (clause-3 authority) · Verify profile: ailang-code

## §1 Problem

The repo has **no inbound credential → session resolution at all**. This is not an unbuilt nicety; it is a contradiction with the bar. Clause 3 (quoted verbatim in §2, ratified and not relaxable) requires *"explicit authority end-to-end: every effect goes through the broker with a capability + budget check"*. But nothing at HEAD authenticates an inbound HTTP request, so there is no way for an external actor to enter the broker's authority domain except by the operator's own hands. `D-WORLD-26` settled the **envelope** — the session credential rides the standard `Authorization: Bearer` HTTP header, with constraints (i) it is a SESSION credential and never an API key, and (ii) the resolver fails closed — but it did **not** supply the **contents**: who mints a credential, where the credential→(episode, grants) mapping lives, what expires/revokes it, and how the resolver fails. This design builds that contents.

The measured absence is unambiguous and first-party:

- **F1** — `git grep 'Bearer' -- host/` → **0**; the look-up family → **0**; control `git grep 'Session' -- host/` → **181** lines; fresh absent literal `ZxqSessionBlobV9` → **0**. There is no bearer path and no session-by-ID look-up. ([controller-measured 2026-09-21 at b65a9ec])
- **F2** — `Authorization` appears only once, in a prose comment in `host/broker/approve.go`. Every one of the 128 `Credential` matches is **outbound** (the registry credential World presents upstream). `Authenticate` is all evidence-envelope. Nothing authenticates an inbound request. ([controller-measured 2026-09-21 at b65a9ec])
- **F3** — the daemon's whole HTTP surface is `host/daemon/daemon.go:555-565` (`GET /v1/health`, `GET /v1/head`, `GET /v1/worlds/{ref}`, `GET /v1/objects/{ref}`, `GET /v1/log/{index}`, `GET /v1/log`, `GET /v1/registry/{name...}`, `POST /v1/commit`, `GET /workbench`). `handleCommit` decodes the JSON and calls `d.store.Commit(commit)` **directly** — no session, no grant check, no authority on any route. The only `NewSession` call outside tests is a benchmark. ([controller-measured 2026-09-21 at b65a9ec]; I re-read the handler at HEAD, lines 517–577, and `POST /v1/commit` is the only mutation and calls `d.store.Commit(commit)` with no session.)
- **F4** — `broker.NewSession(s *store.Store, episodeID string, grants []Capability, registry Registry)` takes grants as an **argument**; nothing decides what grants a credential carries, and sessions are not persisted (`git grep -n 'Session' -- host/store/` → 0). `broker.Session` is in-process by design (clause 3), so the mapping is exactly the pre-existing gap. ([controller-measured 2026-09-21 at b65a9ec]; I read `broker.go:87` and the `Session` struct at HEAD.)

**What row 40 consumes (F6, verbatim from `design_docs/planned/w-a2a-session-projection.md`)**: an opaque `Authorization: Bearer` credential resolves to **`(episodeID, grants, expiry)`**, with **typed denial** on absent / malformed / unknown / expired, each distinguishable and each fail-closed, **no API-key fallback**, **no alternate-header fallback**, **no unauthenticated degrade path**, **timing-safe credential comparison** (this design's hash-then-index lookup realizes the contract's intent — D3/D5/§6), and **bounded lookup**. Row 40 is BLOCKED on this row: it re-binds its Decision-3 responsibility 1 to "the interface row 39 must land," and the resolver's prose name in `D-WORLD-26`'s ruling (`protocol.SessionResolver.ResolveSession(ctx, *http.Request)`) is "the ratified ENVELOPE description, not a code symbol." This design fixes the actual Go interface and package.

**Why a hard 1.0 blocker**: without this boundary, every remote path to the broker is ambient authority — a caller can `POST /v1/commit` to the sole mutation with no capability check, which is precisely the "no ambient-authority path from an agent to the outside world" that clause 3 forbids. It is the groomed queue's #2 product pick (row 8 closed today), and it gates row 40. Nothing under `F11` (the `go-verify` harness defect) is touched here; this item's gates run locally (F10) and are green at base.

## §2 Design

Decisions D1–D7 follow. Two ratified constants thread through every one and are quoted verbatim:

> **Clause 3 (the bar, verbatim)**: "explicit authority end-to-end: every effect goes through the broker with a capability + budget check; effect results are recorded (replay input); capsules run with a physical isolation floor beneath the semantic checks. No ambient-authority path exists from an agent to the outside world."

> **D-WORLD-26 (Mark, attended, directive `A` 2026-08-25T19:06:41Z)**: the session credential rides the standard `Authorization: Bearer` HTTP header. Constraint (i): a `Bearer` value here is a SESSION credential and NEVER an API key (the static `serve-api` key was measured process-wide at iter-24, so it cannot represent a session). Constraint (ii): the resolver must fail closed on an absent, malformed or unknown credential rather than degrading to an unauthenticated surface. The alternate header `X-World-Session` is REJECTED and must not be read, even as a fallback.

> **Clause 2 (local-first)**: the daemon runs on one machine over SQLite, zero cloud dependencies in the core. The session-authority mapping must be local-first too.

### D1 MINTING — a new `ailang-worldd session mint` client verb, operator-facing, one human act

**Decision**: mint through a new operator-facing CLI verb, **`ailang-worldd session mint`**, not through the workbench. `cmd/ailang-worldd` is the established flag-based CLI for operator one-acts (`serve`, `object get`, `log range`, `commit` — `flag.NewFlagSet("ailang-worldd <verb>")` in `cli.go`, top-level `switch verb` dispatch in `main.go:117`), and the F7 precedent (`tools/attended/self_mod_publish.sh` "preflight / approve / dry-run / live, one human act each") is that authority-relevant one-acts are explicit CLI verbs. The workbench (`GET /workbench`) is a read-only renderer — minting from it would put the credential minting surface inside the same request path it authorizes, which is circular for a first version.

Exact surface and flags:

```
ailang-worldd session mint --episode <id> \
                           [--grant EFFECT=SCOPE:BUDGET]... \
                           [--ttl <seconds>] \
                           [--out <path>]          # default: print to stdout once
```

- `--episode` (required, non-empty) — the episode the session binds to.
- `--grant` (repeatable, at least one) — `EFFECT=SCOPE:BUDGET`, each mapping to one `broker.Capability{Effect, Scope, Budget}`. `ExpiresAt` is derived at mint from TTL, not a flag — the operator chooses a TTL, not an absolute instant, so a long-lived mint can be expressed without the operator owning a clock (the rig's `date` is BSD, CI's is GNU; we never parse dates across that seam, F-batch §7).
- `--ttl` (default 3600) — lifetime in whole seconds; the minted credential expires `ttl` seconds after mint.
- `--out <path>` (optional) — write the credential **once** to `<path>` at mode 0600 (F8 precedent: mode-0600 file **outside the tree**); when omitted, print to stdout (also once). With `--out`, the CLI still prints a **one-line confirmation** (episode, grant count, expiry epoch, and the path) but **never the credential**.

`session mint` is a client verb that speaks **no** new HTTP endpoint itself — the F8 precedent is deliberately transport-agnostic operator-held-secret handling. It constructs the mapping **directly against the store DB path** (flag `--db`, default matching `serve --db`). Rationale for DB-path (not daemon-HTTP) minting: the daemon is not guaranteed to be running when an operator provisions a session, and letting the daemon both mint and hold the store write path adds a mint endpoint that, until D6 wires it, nothing authenticates. Minting against the store is local-first (clause 2) and keeps the operator's tty fence structural. Where that is insufficient (F13 residual, see §9), a daemon-owned mint endpoint is explicitly deferred, never silently unbuilt.

Mint is **tty-fenced** exactly as `cmd/world-publish/tty.go` fences `publish` (the "one candidate this loop is STRUCTURALLY unable to fake": a controlling terminal). `session mint` refuses when it has no controlling terminal (`os.Open("/dev/tty")` fails with `ENXIO`), because a mint that can be driven by a pasted/redirected input collapses "one human act" into "whoever runs the script gets a minted session credential." The human confirms by answering a Y/N prompt read from `/dev/tty` (the F7 one-act-precedent, honestly extended from the self-mod fences).

### D2 MAPPING STORE — a new SQLite table `session_credentials` in `store.Store`, persisted (survives restart)

**Decision**: the credential→(episodeID, grants, expiry) mapping lives in a new `store.Store` SQLite table, **persisted**, not memory-only and not a sidecar file.

- **Not memory-only**: a memory-only mapping dies with the daemon process, making every minted credential single-lifetime (F4's in-process `broker.Session` is a *consequence* of clause 3, not a reason to drop the mapping on the floor). A session credential minted by an operator outlives the daemon's next restart; single-lifetime credentials are a denial-of-service on the operator's own workflow. Persisting is the only locally-first choice that lets a session outlive a daemon restart.
- **Not a sidecar file**: `store.Store` already owns SQLite and the schema-version discipline (`host/store/store.go`, `go:embed schema.sql`, `PRAGMA user_version`, `freshInitTx` / `enforceSchemaVersion`); a table is the least-new-surface addition. There is no session/credential table today (`git grep -n 'Session' -- host/store/` → 0, F4).
- **Restart boundary resolved**: because the mapping is durable, a credential outlives restarts and `broker.Session` (in-process) is CONSTRUCTED from the mapping at resolve time (D5) and then lives only for the request/session that needs it. Restart does not invalidate a live credential; only expiry or revocation does.

Schema (added to `schema.sql`; version bump discussed below):

```sql
CREATE TABLE IF NOT EXISTS session_credentials (
    credential_id TEXT PRIMARY KEY,
    episode_id TEXT NOT NULL CHECK (episode_id <> ''),
    grants_json TEXT NOT NULL,
    expires_at INTEGER NOT NULL,
    created_at INTEGER NOT NULL
);
``` (round-2 quorum fix applied verbatim: the never-used `revoked` column is removed — D4 revokes by `DELETE`, so nothing ever reads a revoked flag; the secondary index on `credential_id` is removed — the TEXT PRIMARY KEY already carries the unique index a point lookup needs.)

The PK on `credential_id` already makes resolve an **indexed point lookup by exactly one column** (D5 bounded lookup): one `WHERE credential_id = ?` on the PK, never a scan. `grants_json` carries a copy of the grants the caller will hand `broker.NewSession`; it is the mapping's content, exactly the `grants []Capability` that `NewSession` takes as an argument (F4 — nothing decides a credential's grants today; this table is that missing decision).

**Schema version**: `currentSchemaVersion` is `2` (`store/store.go:42`), with `frozenFutureSchemaVersion = 3` in the version test. Adding a table is **not** an in-place mutation of an existing table, so the additive `CREATE TABLE IF NOT EXISTS` is safe — but the store's discipline is a **version-gated schema** (`enforceSchemaVersion` refuses a *future* or *legacy* version rather than silently altering). The design **bumps `currentSchemaVersion` `2 → 3`** and updates the frozen expectation constants and the migration test row so an un-migrated v2 DB is refused loudly (the `LegacySchemaVersionError` path already exists and is tested) rather than silently gaining a table mid-flight. This matches the established migration precedent (`1856bfb` "schema 1→2" landed `approval_claims` by version-gated DDL; `schema.sql` is re-applied via `db.Exec(schemaSQL)` on open). A fresh store initializes at v3; a v2 store is refused with the existing `legacy` error and a doc-pointer message, giving the operator a defined upgrade path rather than a silent table-add. **§6** records who else reads this schema surface.

### D3 CREDENTIAL FORMAT — opaque 256-bit token, hex, stored as its SHA-256, printed once

**Decision**:

- **Format**: `32` random bytes from `crypto/rand` (`rand.Read`), rendered as **64 lowercase hex chars**. No structure, no embedded episode/grants/timestamp — opaque by design so the token reveals nothing and no field can be guessed from its prefix (F6 requires opaque).
- **Stored form**: the mapping stores **`sha256(token_hex)`** (the hex of the token's SHA-256, 64 chars) in `credential_id`. The raw token is **never persisted**; only the hex you hand over the wire and its hash are ever materialized. This is the D-WORLD-26-faithful choice: a leaked `session_credentials` table rows (e.g. a DB backup) does **not** reveal usable credentials, and the operator still prints the raw token exactly once at mint for out-of-band handoff.
- **Timing safety — hash-then-index-lookup, not a Go-level compare** (F9): the timing-safe mechanism is that the presented token is **SHA-256-hashed to a fixed-length value before any comparison**, and equality is the SQLite **indexed point lookup** (`WHERE credential_id = ?` on the PK). No Go-level byte-for-byte compare happens. This is secure because the attacker cannot control the **stored** hashes — they are hashes of 256-bit random tokens the attacker does not know — so a B-tree prefix-comparison timing difference reveals nothing about token material: leaking a hash prefix does not reveal the preimage token. We state plainly: a **Go-level `crypto/subtle.ConstantTimeCompare` after an index hit is vacuous** — the index already performed the exact equality to find the row, so a subsequent Go compare would merely compare the query parameter against the row that matched it (a value against itself) — and this design does **NOT** claim one. `crypto/subtle` is therefore kept **OUT of the resolve path**, for the reason just named: an indexed lookup on a credential hash is already timing-safe against token-prefix leakage, so an added constant-time compare would contribute nothing but a vacuous claim. F9 is respected as observed (nothing in `host/`/`cmd/` uses `subtle` today) and the design keeps it that way.
- **What prints vs what persists**: at mint the raw hex prints **exactly once** (stdout or `--out`, D1). The store row holds `credential_id = sha256(token_hex)` + grants + expiry — never the raw token. Nothing logs the raw token; any error/result the mint path emits includes only the `credential_id` (the hash), never the secret (consistent with F9's `redactedMarker` spirit and F8's "contents never read" handling of operator-held secrets).

### D4 EXPIRY / REVOCATION — TTL at mint, `expires_at` enforced at resolve, delete-revokes, live sessions immutable

**Decision**:

- **Expiry**: TTL is set at mint (D1 `--ttl`); the store row records `expires_at = mintUnix + ttl`; **expiry is enforced at resolve**, not at mint. A resolve whose `now > expires_at` returns the **distinct** `Expired` denial (D5), HTTP 401; the row is untouched (it stays visible so diagnostics can distinguish "expired" from "unknown"). The broker does not read a wall clock (F4: `Session` snapshots take `now` from the caller) — the resolver supplies `now` (RFC 3339 clock read once at the `session` boundary, not per grant).
- **Revocation**: revoke is a CLI verb **`ailang-worldd session revoke <id>`** (id = the `credential_id` hash, printed at mint), which **deletes the mapping row** (`DELETE FROM session_credentials WHERE credential_id = ?`). A deleted row means the next resolve of that credential is an **unknown** credential — the design explicitly states the equivalence: **delete-from-mapping = the credential becomes "unknown", not "revoked"**, because there is nothing left to distinguish it from a token that never existed. Revocation is therefore total and immediate (no check at the middleware beyond the row's absence). The revocation `DELETE` is applied in a transaction against the store's single-connection discipline (`store.Store` `SetMaxOpenConns(1)`); it never races a concurrent resolve on the same process.
- **mid-flight LIVE `broker.Session`**: when a credential is revoked or expires, an already-constructed, live in-process `broker.Session` **is not torn down**. Sessions are immutable grants: `broker.Session` carries its `grants` snapshot in-process, and `CapabilitySnapshot` is explicitly "an immutable copy of a session capability ledger at **one instant**" (`broker.go:60-63`). Expiry/revocation **gate NEW resolutions** (D5); a session already handed out finishes its grant budget (`Budget` debit logic unchanged). This is a **declared residual** (§9, residual R2): a long-lived in-process session holding a revoked grant can still spend until it returns or its budget is exhausted. The alternative — tearing down live sessions — is rejected because there is no handle-identification path from the store table to every in-process `Session` (sessions are not persisted, F4), and killing them would break the broker's single-`epoch` serialization contract without a compensating benefit in this sprint.

### D5 RESOLVER INTERFACE + TYPED DENIALS — `authority.Resolver` in a new `host/authority` package

**Decision**: a new package **`host/authority`** owns mint, resolve, and denial. The interface is named **`authority.Resolver`** (the ruling's `protocol.SessionResolver` prose name is explicitly "an envelope description, not a code symbol" — F6) and lives at the `host` boundary, not in `protocol/` (that namespace belongs to the A2A projection layer, row 40). Exact types:

```go
package authority

// ResolveOutcome is the typed result of credential resolution. Exactly one of
// its four denial kinds (or Success) is set per Resolve call. Timing-safe
// hash-then-index lookup (D3) and bounded lookup are inside Resolve, not the
// caller's job; there is NO crypto/subtle Go-level compare in this path.
type ResolveOutcome struct {
    Success   *SessionBinding // set iff the credential is known and unexpired
    Denied    *DenialKind     // set iff not success; each kind is distinct
}

type DenialKind uint8
const (
    DenialAbsent    DenialKind = iota // no Authorization: Bearer header at all
    DenialMalformed                   // header present but not "Bearer <token>" or token malformed
    DenialUnknown                     // hash present-shaped but no mapping row (or row revoked/deleted)
    DenialExpired                     // mapping row exists but expires_at < now
)

// SessionBinding is the credential→(episode, grants, expiry) mapping this sprint
// hands to row 40 (F6). ExpiredAt is the derived `expires_at` so the consumer
// can make its own liveness call without re-hashing.
type SessionBinding struct {
    EpisodeID string
    Caps      []broker.Capability
    ExpiresAt int64
    CreatedAt int64
    /* CredentialID (the stored hash) is deliberately NOT exposed to row 40 —
       no consumer needs the secret's hash as a capability signal. */
}

type Resolver interface {
    // Resolve maps an opaque Authorization: Bearer credential to a binding.
    // header is the ENTIRE Authorization header value, or "" if absent. It
    // returns a typed outcome; it NEVER returns a nil-outcome-on-non-error.
    Resolve(header string, now int64) ResolveOutcome
}
```

- **Fail-closed taxonomy and HTTP mapping**: absent → 401 with a distinct `apiErr.Class`; malformed → 400; unknown → 401; expired → 401. Each is an **independent, distinguishable error** (F6 demands each distinguishable), delivered through `daemon.writeAPIError`. Absent is `401 Unauthorized` (nothing to name); malformed is `400 Bad Request` (the header is present but not a valid `Bearer`); unknown and expired are both `401` but **distinct classes** (`SessionUnknown` vs `SessionExpired`) so a consumer can tell a bad token from an expired one — the row 40 consumer does this today to render "please mint again" vs "your session is stale."
- **Absent vs malformed**: absent is the header string being empty; malformed is a non-empty header that fails `strings.Fields(header)` == `["Bearer", token]` shape or has an empty/whitespace token or a token whose hex length ≠ 64 / non-hex chars. A credential that parses into hex but has no mapping row is **unknown**, never malformed.
- **Timing safety — equality is the indexed lookup, no Go-level compare** (F9): the resolver performs **no** `crypto/subtle.ConstantTimeCompare`, consistent with D3. The presented token is hashed to a fixed 64-hex value and equality is the SQLite PK lookup; the unknown-vs-known tokens traverse the **identical** code path — one and the same single query shape. A constant-time compare after the hit would be vacuous, so it is not claimed. `crypto/subtle` stays OUT of the resolve path by design. The guard that this stays true is AC-M3-1 (source-scan, §4) and its killer mutation M1 (§5).
- **Bounded lookup**: the lookup key **IS the credential hash**; it is re-selected as `credential_id` in the result so the row's identity stays available to the resolver's diagnostics (harmless; it is the key, not a re-typed comparator). exactly one indexed query per resolve — `SELECT credential_id, episode_id, grants_json, expires_at, created_at FROM session_credentials WHERE credential_id = ?` with `? = sha256(presentedToken)`. There is no scan, no `LIKE`, no fallback table. If the hash is empty/absent, `DenialAbsent` **without** touching the store (absent header is refused before any query — bounded by construction).
- **No fallback of any kind**: the resolver reads **only** the `Authorization` header; it never reads `X-World-Session` (REJECTED by D-WORLD-26) and never treats the static `serve-api` API key as a session (constraint (i)). The mutation M6 (fall back to `X-World-Session`) makes a named test red.
- **`Resolve` returns a typed outcome, never `(nil, nil)`**: the interface makes the "unauthenticated degrade path" structurally impossible — there is no code path where a missing/unknown credential yields a *usable* session. Every non-success branch is one of the four denials.

### D6 WIRING — a middleware wrapper on `d.Handler()`'s mux for the session-bearing routes; read routes defer as a named residual

**Decision**: this sprint delivers authority **primitives** (mint/resolve/deny) + enforcement on the one route that exists behind a session boundary today, and records full-`/v1/*` read-route enforcement as a **named residual**, NOT silently skipped.

The daemon's `Handler()` (F3) returns a `*http.ServeMux` with the nine routes. The recommended wiring is a **daemon-side middleware** — `host/daemon/middleware.go` (round-2 quorum fix applied verbatim: the middleware lives in the daemon package so it can call the package-local `writeAPIError` directly; `host/authority` must NOT import `host/daemon`, which the original placement would have cycled with) — that wraps the mux: `d.Handler()` becomes `authSessionMiddleware(d.resolver, protected, mux)`, where `protected` is the set of routes that require a session this sprint. Concretely this sprint enforces **`POST /v1/commit`** only, for the reason in D7. The eight GET routes (`/v1/health`, `/v1/head`, `/v1/worlds`, `/v1/objects`, `/v1/log/…`, `/v1/log`, `/v1/registry/…`, `/workbench`) stay unauthenticated **this sprint** and are recorded as **residual R3 (§9)**: full `/v1/*` read enforcement is a follow-up/queue-row, because row 40's projection is the session-scoped *write* consumer and putting auth on every read would change byte-visible behavior for every existing REST reader (P6 in row 40 freezes the REST contract) without a compensating benefit in this sprint. The middleware, however, is route-agnostic by construction (it takes a `protected` set), so flipping on the read routes later is an `if`-set change, not a rewrite.

The middleware logic:

```go
func (m *SessionMiddleware) Wrap(protected func(r *http.Request) bool, mux http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        if !protected(r) {
            mux.ServeHTTP(w, r)             // unprotected route passes through (residual, §9)
            return
        }
        header := r.Header.Get("Authorization")
        out := m.resolver.Resolve(header, time.Now().Unix())
        switch {
        case out.Denied != nil:
            switch *out.Denied {
            case authority.DenialAbsent, authority.DenialUnknown, authority.DenialExpired:
                writeAPIError(w, classFor(*out.Denial), msgFor(*out.Denial), http.StatusUnauthorized) // in-package (daemon), no host/authority→host/daemon edge
            case authority.DenialMalformed:
                writeAPIError(w, "InvalidSession", msgFor(*out.Denied), http.StatusBadRequest)
            }
            return
        default: // Success
            ctx := context.WithValue(r.Context(), sessionCtxKey{}, out.Success)
            mux.ServeHTTP(w, r.WithContext(ctx))
        }
    })
}
```

`handleCommit` then reads `binding, ok := authority.FromContext(ctx)` (a tiny typed accessor; context keys are unexported and typed, not a global string bag) and binds the episode/grants before building the session (D7). The X-World-Session header is **never read** anywhere in this file (D-WORLD-26; §6 records the rejection).

### D7 handleCommit — today it writes with no authority; this sprint binds it to the resolved session

**Decision**: **this sprint changes `handleCommit` to require a session.** Because `POST /v1/commit` is the daemon's sole mutation (F3: "the surface's only mutation"), leaving it ambient would leave clause 3's "no ambient-authority path" unbuilt at the one place it matters. Under the D6 middleware, a caller must present a valid, unexpired, known session or the request is denied before `decodeCommit` runs. On success, `handleCommit` constructs the broker session from the resolved binding — `broker.NewSession(d.store, binding.EpisodeID, binding.Caps, registry)` — rather than calling `d.store.Commit` with zero authority. **The store `Commit` semantics are unchanged**; the session is the capability/budget gate between the HTTP body and the store write. This is honest and consistent with clause 3: authority is now a real hop between the wire and the kernel, exactly where the bar requires it.

Whether the commit should additionally be *re-recorded as an effect outcome through `broker.Session`'s effect path* (replay input per clause 3) is carried as a refinement **residual R4** for the row-40 consumer, which owns the A2A propose → verify → commit sequence; this sprint's job is the boundary, not re-plumbing every mutation into the effect-ledger. The residual is stated plainly so it is not mistaken for done.

**Worked example** (four denial kinds then a success; operator's tty is a controlling terminal):

```
$ ailang-worldd session mint --db ./world.db --episode ep-42 --grant fs.read=/tmp/b:1 \
     --ttl 3600
Confirm mint for episode ep-42 (1 grant, expiry +3600s) with a session credential? [y/N] y
credential_id=9f2c…<64-hex-hash>
<raw-64-hex-token-printed-ONCE>

$ curl -i -X POST http://127.0.0.1:7644/v1/commit -d '…'                      # no header
HTTP/1.1 401          class=SessionAbsent
$ curl -i -X POST http://127.0.0.1:7644/v1/commit -H 'Authorization: Key abc' -d '…'
HTTP/1.1 400          class=InvalidSession          # malformed
$ curl -i -X POST http://127.0.0.1:7644/v1/commit -H 'Authorization: Bearer deadbeef…' -d '…'
HTTP/1.1 401          class=SessionUnknown           # unknown token
$ curl -i -X POST http://127.0.0.1:7644/v1/commit \
     -H "Authorization: Bearer $(mint --ttl 1 …)" ; sleep 2 ; curl …            # expired
HTTP/1.1 401          class=SessionExpired
# happy path:
$ curl -i -X POST http://127.0.0.1:7644/v1/commit \
     -H "Authorization: Bearer $(the-once-printed-hex)" -d '{"…":…}'
HTTP/1.1 200          class=none (commit accepted under session ep-42)
```

`X-World-Session` is never read: a request carrying ONLY `X-World-Session: …` and no `Authorization` is `SessionAbsent` (the middleware never looks at it).

## §3 Files to create / modify

| Path | Action | Est. LOC | Notes |
|---|---|---|---|
| `host/authority/resolver.go` | **new** | ~80 | `Resolver` interface, `DenialKind`, `ResolveOutcome`, `SessionBinding`, `FromContext`/typed context keys, timing-safe hash-then-index + bounded hash lookup; no `crypto/subtle`. |
| `host/authority/mint.go` | **new** | ~70 | Token generation (`crypto/rand`), hashing, `INSERT INTO session_credentials`, one-line print / `--out` 0600 write. |
| `host/authority/revoke.go` | **new** | ~30 | `DELETE FROM session_credentials WHERE credential_id = ?` in a transaction. |
| `host/daemon/middleware.go` | **new** | ~50 | `SessionMiddleware.Wrap(protected, mux)` — daemon-side (round-2 quorum fix), calls in-package `writeAPIError`; route-agnostic protected set; typed denial → HTTP; no `X-World-Session` read. |
| `host/authority/<resolver,mint>_test.go` | **new** | ~150 | Named tests for ACs + mutation matrix (§4, §5). |
| `host/store/schema.sql` | **modify** | +9 | add `session_credentials` table + `idx_session_credential_lookup` (D2 DDL). |
| `host/store/store.go` | **modify** | +~35 | `currentSchemaVersion` 2→3; `MintSession`, `ResolveSession`, `RevokeSession` store methods (single indexed query / DELETE). |
| `host/store/schema_version_test.go` + `journal_test.go` fixtures | **modify** | +~20 | bump `expectedCurrentSchemaVersion`/`frozenFutureSchemaVersion`, add v3 migration row; keep fixtures non-vacuous. |
| `cmd/ailang-worldd/main.go` | **modify** | +~10 | add `case "session"` dispatch. |
| `cmd/ailang-worldd/session.go` | **new** | ~90 | `session mint` / `session revoke` verbs + flags + tty fence (port of `cmd/world-publish/tty.go` pattern). |
| `host/daemon/daemon.go` + `host/daemon/handlers.go` | **modify** | +~40 | wire `authority` middleware in `Handler()`; route `POST /v1/commit` protected; `handleCommit` builds session from resolved binding. |
| `design_docs/planned/w-session-authority.md` | this doc | — | — |

Nothing under `tools/launchd/*` (fleet-frozen core, F12) and **no `.ail` files** (F12). No `.ail` gate is hit.

## §4 Acceptance criteria

All run **LOCALLY** at the sprint tree with `AILANG_BIN=$HOME/.pinned-ailang/ailang` (F10, v0.41.0). `scripts/verify_go.sh` is **NOT** a gate (charter row 76), so none of the criteria below require it. Each AC is a COMMAND + EXPECTED OUTPUT.

- **AC-M1-1 mint prints exactly once.** `ailang-worldd session mint --db <tmp>/world.db --episode e1 --grant fs.read=/tmp/b:1 --ttl 3600` under a PTTY (test drives `os.Open("/dev/tty")` via the fencing seam) → prints the raw 64-hex token **exactly once** and stores one row. Count raw-token occurrences in the output: **exactly 1**.
- **AC-M1-2 distinct credentials per mint.** Mint twice for the same episode; the two printed tokens differ (and both rows have distinct `credential_id`). `test "$a" != "$b"` → **0 (distinct)**.
- **AC-M1-3 resolve returns episode+grants+expiry.** `Resolve("Bearer <token>", now<expires)` → `Success` with `EpisodeID==e1`, `Caps` matching the minted grant, `ExpiresAt==mint+3600`. Resolve **never** returns the raw token.
- **AC-M2-1 absent header.** `Resolve("", now)` → `DenialAbsent`; over HTTP, `POST /v1/commit` with no `Authorization` → **401** `SessionAbsent`.
- **AC-M2-2 malformed header.** `Resolve("Key abc", …)` and `Resolve("Bearer", …)` and `Resolve("Bearer z" /*1 char*/ , …)` → `DenialMalformed`; HTTP → **400** `InvalidSession`. A non-64-hex token is malformed → **400**, never 401.
- **AC-M2-3 unknown token.** `Resolve("Bearer deadbeef…64hex-unknown", …)` → `DenialUnknown`; HTTP → **401** `SessionUnknown`. The static `serve-api` key value as `Bearer` → **401** (constraint (i) — never a session).
- **AC-M2-4 expired.** Mint with `--ttl 1`, `sleep 2`, resolve → `DenialExpired`; HTTP → **401** `SessionExpired`, distinct from `SessionUnknown`.
- **AC-M2-5 revocation.** `session revoke <credential_id-hash>` then resolve the same token → `DenialUnknown` (delete-from-mapping = unknown, D4).
- **AC-M2-6 no alternate-header fallback.** Request with ONLY `X-World-Session: <a valid minted token>` and no `Authorization` → **401** `SessionAbsent` (must NOT resolve). Under mutation M6 it resolves — so a green control proves the absence is really enforced.
- **AC-M3-1 no direct secret comparison.** A **source-scan test** over `host/authority/*.go` **FAILS** if any direct equality (`==`, `bytes.Equal`, `strings.EqualFold`, `strings.Contains`) is applied to token or credential-hash material — pinning that no future edit introduces a direct secret comparison into the resolver. Behavioral note: unknown-vs-known tokens traverse the identical code path (the same single `SELECT … WHERE credential_id = ?` query shape, D3/D5), so there is no branch an online adversary can time-difference. Killers: mutation M1 (§5).
- **AC-M3-2 bounded lookup.** Assert the resolve path issues **exactly one** SQL query: wrap the store/DB seam (test injects a counting `*sql.DB`-like observer or uses the store's existing injected-`objectStore` seam `daemon.go:267`) and assert query count == 1 with the single indexed `SELECT … WHERE credential_id = ?`.
- **AC-M4-1** `go test ./... -count=1` → **19 ok or better** (19 packages stay green / 0 FAIL; the new `host/authority` package adds its own test package that must also pass). `go vet ./...` → **rc=0**.
- **AC-M4-2** `bash ./scripts/verify_ail.sh` with `AILANG_BIN=$HOME/.pinned-ailang/ailang` → **"✓ verify gate PASSED: 11 required identities verified, 40 named tests pass"** and **"✓ world package gate PASSED: 9/9 steps"** (unchanged — this item adds no `.ail`).
- **AC-M4-3** store schema bump: a v2 store DB opened by the new binary is **refused** (`LegacySchemaVersionError`); a fresh DB initializes at v3 and the migrate row/constraints in `schema_version_test.go` are updated and green.

## §5 Test plan / mutation matrix

Each mutation: `cp <file> <file>.bak && sed -i '…' <file>; <run target test>`; restore `cp <file>.bak <file>`; restore **asserted** by `shasum -a 256 <file>` == the pre-mutation digest AND `test -x <file>` (the `.bak` is `cp -p` so the mode bit survives; `test -x` guards a mode that matters, e.g. if a test file were made executable); **never** `git checkout --` (the tree may be dirty; a checkout would silently discard unrelated work).

| # | Mutation | Which named test goes RED (sole killer) | Mechanism |
|---|---|---|---|
| M1 | Introduce a direct `==` comparison of token material into the resolver (e.g. replace the hash-then-index with `if presentedHash == storedHash`) | `TestAuthority_NoDirectTokenCompare` | the AC-M3-1 source-scan test over `host/authority/*.go` reds on any direct equality (`==`/`bytes.Equal`/`strings.EqualFold`/`strings.Contains`) of token or credential-hash material (sole killer). Restore discipline as written above. |
| M2 | Collapse `DenialUnknown` into the absent branch (unknown token → `SessionAbsent`/401) | `TestResolve_TypedDenialTaxonomy` | the taxonomy test asserts unknown and absent are distinct kinds (AC-M2-1 vs AC-M2-3). |
| M3 | Remove the `now < expires_at` check at resolve (ignore expiry) | `TestResolve_ExpiryEnforced` | the expiry test mints `--ttl 1`, advances `now`, asserts `DenialExpired` (AC-M2-4). |
| M4 | Mint prints/logs the token twice (duplicate print) | `TestMint_PrintsOnce` | asserts raw token appears exactly once in mint output (AC-M1-1). |
| M5 | Revocation leaves the row (the `DELETE` is omitted — revoke becomes a no-op against the store) | `TestRevoke_DeletesRow` | asserts the row is absent post-revoke and resolve of a revoked token is `Unknown` (AC-M2-5). |
| M6 | Fall back to reading `X-World-Session` when `Authorization` is absent | `TestNoAlternateHeader` | the carrier test sends only `X-World-Session` and demands `DenialAbsent` (AC-M2-6). |
| M7 (green control) | Reword a comment in the resolver (no behavior change) | **all tests stay green** | proves the mutant-detection matrix is not vacuously green. |

## §6 Conflict surface

- **Daemon Handler/mux (F3)**: the nine routes in `Handler()` (listed §1) are the only HTTP surface; wiring the middleware wraps this mux and touches no route's handler signature except `handleCommit`'s internal session binding. `handleCommit` is currently the only route that reads a body and the only mutation; adding the session hop changes its authority behavior without renaming/removing any route (P6 in row 40 freezes the REST contract — this respects it for the eight GETs, which keep byte-compatible responses).
- **`broker.NewSession` readers**: the only non-test callers are `cmd/ailang-worldd` commit path and, once wired, `handleCommit`; the benchmark `host/daemon/bench_test.go:474` passes a literal capability grant and does **not** go through the resolver — the new session layer must not change its behavior (it benchmarks the broker directly, not HTTP). the new `host/authority` package imports `host/broker` for `Capability` and `NewSession`; the imports are strictly one-way (`host/authority` → `host/broker`; `host/daemon` → `host/authority` and `host/broker`; `host/authority` does NOT import `host/daemon` — the round-2 fix moved the middleware into the daemon package precisely so the denial→HTTP mapping uses the in-package `writeAPIError`), so no import cycle is introduced (authority is above broker, daemon above authority).
- **Store schema (D2)**: `store/store.go` owns `go:embed schema.sql`, `currentSchemaVersion=2`, `enforceSchemaVersion`, `freshInitTx`, and the legacy/future/invalid refusal paths (all tested in `schema_version_test.go`). The migration bump touches exactly those constants and the additive table DDL. `schema_version_test.go:19` asserts `currentSchemaVersion == expectedCurrentSchemaVersion` — both are bumped together. Nobody else reads or writes a session/credential table today (`git grep -n 'Session' -- host/store/` → 0, F4).
- **Row 40 (F6)**: this interface satisfies row 40's dependency contract **verbatim**: opaque credential → `(episodeID, grants, expiry)` = `SessionBinding{EpisodeID, Caps, ExpiresAt}`; typed denial on all four kinds, each distinct and fail-closed = `DenialKind`; no API-key fallback & no alternate-header fallback = only-header-read + `X-World-Session` never read; no degrade path = `Resolve` returns a typed outcome, never usable-on-error; **the contract's "constant-time credential comparison"** — its INTENT (no timing side-channel on credential material) is satisfied by the **hash-then-index mechanism** (D3/D5): the presented token is hashed to a fixed length before any comparison, equality is the SQLite indexed PK lookup, and leaked B-tree prefix-timing reveals no token material because the stored values are SHA-256 hashes of unknown 256-bit random tokens. The replacement guard for the phrase is the **source-scan test of AC-M3-1** (FAILS on any direct equality of token/credential-hash material); `crypto/subtle` is deliberately NOT used (a post-hit compare would be vacuous). bounded lookup = single indexed PK query. Row 40's `ResolveSession(ctx, r)` can be an adapter over `authority.Resolver.Resolve(r.Header.Get("Authorization"), now)`.
- **Mint is operator-facing (D1)**: `cmd/ailang-worldd` dispatch (`main.go:117 switch`) gains `case "session"`; the tty fence mirrors `cmd/world-publish/tty.go` (controlling-terminal requirement) and the one-act precedent of `tools/attended/self_mod_publish.sh`. `session mint` takes `--db`, never `--addr`, so it does not depend on a running daemon and does not cross the HTTP surface it would otherwise authorize.
- **X-World-Session rejection**: the middleware and resolver never read that header — not as fallback, not as preference. A request with only `X-World-Session` is `DenialAbsent` (AC-M2-6), guarded by M6.

## §7 Verification Log

| # | Claim | Provenance / Command | Observed |
|---|---|---|---|
| V1 | No bearer token handling inbound | `git grep 'Bearer' -- host/` → **0** | [controller-measured 2026-09-21 at b65a9ec] (F1) |
| V2 | No session-by-ID resolve | `git grep -E 'func .*(GetSession\|LookupSession\|ResolveSession\|SessionByID\|FindSession)' -- host/ cmd/` → **0**; control `git grep 'Session' -- host/` → **181** | [controller-measured 2026-09-21 at b65a9ec] (F1) |
| V3 | Absence instrument is real (fresh literal) | `git grep -E 'ZxqSessionBlobV9' -- host/` → **0** | [controller-measured 2026-09-21 at b65a9ec] (F1) |
| V4 | No inbound `Authorization` code path | `git grep 'Authorization' -- host/` → **1** (prose comment in `approve.go`); `Credential` matches all outbound; `Authenticate` evidence-only | [controller-measured 2026-09-21 at b65a9ec] (F2); **my re-read** of `host/broker/approve.go` shows only the comment path |
| V5 | Daemon HTTP surface + `handleCommit` writes with no authority | `host/daemon/daemon.go:555-565` (nine routes); `handlers.go:517-577` decodes then `d.store.Commit(commit)` | [controller-measured 2026-09-21 at b65a9ec] (F3); **I ran** `sed -n '540,580p' host/daemon/daemon.go` and `sed -n '517,580p' host/daemon/handlers.go` — confirms routes and direct Commit |
| V6 | `NewSession` grants are an argument; only non-test caller is none | `broker.go:87`; `git grep -n 'NewSession' -- host/daemon/ cmd/` → only `bench_test.go:474` | [controller-measured 2026-09-21 at b65a9ec] (F4); **I read** `broker.go:40-120` (`Session` struct, `CapabilitySnapshot`) and `bench_test.go:465-485` |
| V7 | Sessions not persisted | `git grep -n 'Session' -- host/store/` → **0** | [controller-measured 2026-09-21 at b65a9ec] (F4); **I ran** at HEAD, 0 lines |
| V8 | No constant-time compare in repo; design does NOT adopt it | `git grep -n 'subtle\|ConstantTime' -- host/ cmd/` → **0** (F9); the resolve path keeps `crypto/subtle` OUT (hash-then-index, D3/D5) — a post-hit Go compare would be vacuous | [controller-measured 2026-09-21 at b65a9ec] (F9); [quorum round-1, V16] |
| V9 | CLI verb precedent + dispatch | `main.go:78-159` `switch verb` (`serve|health|head|world|object|log|registry|commit|help`); `cli.go` `flag.NewFlagSet` | **I ran** `sed -n '74,165p' cmd/ailang-worldd/main.go` and `grep -n` of `cli.go` verbs |
| V10 | Store schema-version discipline | `store.go:42 currentSchemaVersion=2`; `schema_version_test.go:17-19` `frozenFutureSchemaVersion=3`/`expectedCurrentSchemaVersion=2`; `enforceSchemaVersion` legacy/future/invalid refusal; migration precedent `1856bfb` (1→2, `approval_claims`) | **I ran** `sed -n '300,370p' store.go`, `sed -n '1,30p' schema_version_test.go`, `git log --oneline -3 -- host/store/schema.sql` |
| V11 | TTY/one-act precedent | `cmd/world-publish/tty.go:26-49` ("controlling terminal … STRUCTURALLY unable to fake"); `tools/attended/self_mod_publish.sh` (one human act each, `/dev/tty` fence) | **I ran** `grep -in` on `cmd/world-publish/tty.go` and read the script header |
| V12 | Row-40 contract consumed | `w-a2a-session-projection.md:150-156` (opaque→(episodeID,grants,expiry); typed denial; no fallback; no degrade; timing-safe comparison; bounded) — contract intent met by hash-then-index (D3/D5), not `crypto/subtle` | **I ran** `sed -n '140,215p'` and grep of the row-40 doc; [quorum round-1, V16] |
| V13 | Gates green at base | `go vet ./...` rc=0; `verify_ail.sh` PASSED (11 identities, 40 tests) + world package 9/9; `go test ./... -count=1` **19 ok, 0 FAIL** — on v0.41.0 pin | [controller-measured 2026-09-21 at b65a9ec] (F10) — **I did NOT re-run** (trusting first-party F10; the doc's ACs define the post-change target) |
| V14 | `go-verify` harness defect is out of scope | CI's `go-verify` still pins v0.30.0 while host/verifygate wants v0.41.0; NOT designed here | [controller-measured 2026-09-21 at b65a9ec] (F11) |
| V15 | No `host/authority` package yet | `ls host/` → no `authority` entry (`ls host/authority` → not found) | **I ran** `ls host/` at HEAD, 12 dirs + no `authority` |
| V16 | Quorum round-1 block on constant-time framing | round 1 (2026-09-21T22:00:09Z): `gemini-3-1-pro` REJECTED the Go-level `crypto/subtle.ConstantTimeCompare`-after-index-hit claim; `gpt5-6-sol` absent auth | [quorum round-1 synthesis 2026-09-21 22:00:09Z] |
| V17 | Controller first-party confirmation of the objection | controller ran the same reasoning against the D5 text: SELECT lacked `credential_id`, and the index hit already establishes equality (a Go compare would be vacuous by construction) | [controller-measured 2026-09-21] |
| V18 | Quorum round-2 block closed under the narrow-refinement carve-out (controller, verbatim fixes) | round 2 (2026-09-21T22:15:05Z): `gemini-3-1-pro` REJECTED on (a) the `authority`→`daemon` import cycle its middleware placement creates, (b) the never-used `revoked` column, (c) the redundant secondary index on the PK; `gpt5-6-sol` absent (auth) again. All three fixes are the reviewer's own `proposed_fix`, applied VERBATIM by the controller: middleware moved to `host/daemon/middleware.go` calling in-package `writeAPIError`; schema reduced to the reviewer's exact DDL; `CREATE INDEX` removed. No design direction disputed. | [quorum round-2 artifact w-session-authority-2026-09-21T22-15-05Z; carve-out per gate-2 (ratified V1 iter-95; used on this mission iter-175)] |

## §8 Milestones

- **M1 (≤ 1d)**: `host/authority` package + `store.Store` schema bump (D2, D3, D5) + the resolve/mint/revoke tests. Store `currentSchemaVersion 2→3`, add `session_credentials` table, `MintSession`/`ResolveSession`/`RevokeSession` methods, timing-safe hash-then-index + bounded-lookup + taxonomy/expiry/revoke tests. Exit: AC-M1-3, AC-M2-1..6, AC-M3-1..2 green on `host/authority`.
- **M2 (≤ 1d)**: `cmd/ailang-worldd session mint|revoke` verbs + tty fence (D1, D4) + daemon middleware wiring + `handleCommit` session binding (D6, D7) + full AC battery + update `schema_version_test`. Exit: all of §4 green, `go test ./... -count=1` 19-ok-or-better, `go vet` rc=0, `verify_ail.sh` PASSED.
- **Controller milestone**: file the `w-a2a-session-projection.md` unblock per §6 (row-40 adapter) and a queue/bookkeeping follow-up if `M2`'s residual list (§9) spawns rows. The planner refines; keep to these 2–3.

## §9 Risks and declared residuals

- **R1 — full `/v1/*` read-route enforcement deferred (D6).** The eight GET routes (health/head/worlds/objects/log/registry/workbench) stay unauthenticated **this sprint**. This is a real authority gap on READ paths, so it is a named residual, recorded for a follow-up queue row, NOT silently skipped. The middleware's route-agnostic `protected` set makes it a config flip later.
- **R2 — live in-process sessions survive revoke/expiry (D4).** A session already constructed into a live `broker.Session` can finish spending its grant budget even if its credential is revoked/expired meanwhile (no handle path from the table to every in-process session; sessions are not persisted, F4). Expiry/revocation gate NEW resolutions only. Plainly declared; the "tear down live sessions" alternative is rejected this sprint (§2 D4).
- **R3 — commit re-record into the effect-ledger (D7).** `handleCommit` is session-gated but does NOT (this sprint) re-record the commit as an effect *outcome* through `broker.Session`'s replay path. Clause 3's "effect results are recorded (replay input)" refinement is owned by row 40's propose → verify → commit sequence; recorded as residual, not claimed done.
- **R13 — daemon-owned mint endpoint (D1) not built.** `session mint` writes to the store DB directly (local-first). If a remote/telnet operator needs to mint against a live daemon, a daemon mint endpoint is NOT in this sprint; it is named as a follow-up. (PENDING if ever needed.)
- **Trust model — credential theft = session theft.** The design defends secrets at rest (table stores only `sha256`), at transit (Bearer over TLS is the operator's deployment concern; the daemon does not force TLS), and at print (once, 0600). It does **not** defend against an attacker who already holds the raw token; that is inherent to bearer credentials and out of scope. Timing-safe hash-then-index lookup + bounded lookup + no fallback reduce online-timing/amplification leakage but are not full mitigation against a compromised channel. Everything unmeasured is labelled PENDING in §7 (none are load-bearing beyond the F-rows' first-party values).

## §10 Fleet note

Nothing in the shared fleet harness (`sunholo-data/ailang`) is touched by this design — no `.ail` changes, no core edits, no skills, no `tools/launchd/*` (fleet-frozen, F12). There is no proposal to the fleet: `host/authority` is World-local product code (`go` host + CLI + store), and the resolver is a World boundary. The only fleet-adjacent surface is the operator-facing verb's tty-fence reuse of the controlling-terminal pattern already demonstrated by `cmd/world-publish/tty.go`, which is a World-owned port, not an upstream change. Charter queue row 40's unblock is the sole cross-row dependency, and it is satisfied by this design's interface (F6) rather than by any fleet change.
