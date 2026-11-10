# w-world-approvals-by-key-a2: passkey approvers (WebAuthn), verified with the standard library

- **Status**: Planned. Designed in iteration 379, **r2** (controller narrow-refinement carve-out applying quorum r2's verbatim fixes; see §15). Queue row 201, `w-world-approvals-by-key`.
  **[r1]** r0 (commit `72e20eb4`) was **BLOCKED 3 of 3** at quorum round 1. None of the three
  disputed the direction, and the controller measured each premise as real. r1 applies all three
  fixes (§15), and the changed text is marked **[r1]**:
  - the enrolment code is now a 128-bit fingerprint, under an attacker who reads the genuine file
    before replacing it (§2.8, §2.12, §4, §5);
  - the page's `blob:` save is measured working under the exact `surfaceCSP` and asserted by XB1
    (§2.14, V31);
  - the parent's and the ledger's ruling text are quoted (V32, V33), and the credential-ID floor is
    dropped to the spec's (§2.5, §2.10, V34).
  **Three sprints: sprint 1 builds the pure verification law (M0–M3); sprint 2 builds the store and
  the approval binding (M4–M5); sprint 3 builds the CLI paths and the ceremony page (M6–M7). Each
  sprint has a named fallback cut (§9).**
  This doc covers parent milestone **A2**. Parent:
  [`w-world-approvals-by-key.md`](w-world-approvals-by-key.md). Built siblings:
  [`w-world-approvals-by-key-a0a1.md`](w-world-approvals-by-key-a0a1.md) (A0, A1a; PRs #533, #535)
  and [`w-world-approvals-by-key-a1b.md`](w-world-approvals-by-key-a1b.md) (A1b; PRs #537, #539,
  #541, #543). Their direction and rulings are settled and not reopened here.
- **Authority**: **D-WORLD-89** (Mark, attended 2026-10-21, "all defaults"). R1: a verified
  signature from an enrolled key may approve a mint, as well as the typed `y`. R2: the formats are
  **WebAuthn (passkeys and FIDO2, user verification REQUIRED)** and OpenSSH SSHSIG. R3: enrolment and
  revocation stay typed-only at the daemon's terminal. R4: publishing stays typed-only.
  **D-WORLD-109 = A** (Mark, attended 2026-11-06): the loop builds and judges row 201 against
  **fakes, temp stores and test keys only**. It never enrols a real key, never opens a real
  credentials store, never mints in one and never publishes. Parent §2.1 fixes the verifier to the
  Go standard library (`crypto/ecdsa` P-256, `crypto/ed25519`, `crypto/sha256`), with an in-repo
  bounded CBOR subset, **no attestation verification (`fmt: none`)** and no new Go dependency.
  Parent §2.2 fixes the challenge (the message hash), the relying-party ID (the surface's host
  name: a tailnet `*.ts.net` name or `localhost`), origin binding, the sign counter, one-time
  consumption and expiry. Parent §2.3 fixes enrolment: it completes in a browser, with a one-time
  enrolment code the operator types at the terminal, and nothing enrols without a typed `y`.
  **Mark's real-passkey check on the tailnet is the attended exit run (§10). The loop never runs it.**
  **[r1]** The ruling text is quoted from the parent and the ledger in V32 and V33. Every clause
  above was found verbatim, except one. "Never opens a real credentials store" is A1b's
  restatement of D-WORLD-109's "never enrols a real key, mints or publishes", kept as A1b kept it.
  **One DEPARTED-FROM-PARENT point, in shape only:** parent §2.3 sketches `approver enroll …
  --format webauthn|ssh [--pubkey <file>]`. A2's webauthn form adds `--origin` and
  `--registration` (§2.12), because the browser's result reaches the CLI as a file until A3's
  POST exists. The rule itself is unchanged: an enrolment code shown by the surface, typed at the
  terminal, then a typed `y`.
- **Item**: four things.
  - The **pure WebAuthn law** in `host/approval`: a bounded CBOR subset, ES256 and EdDSA COSE keys,
    authenticator data, `clientDataJSON`, the registration law and the assertion law.
  - **Credentials schema 6**: a `webauthn_credentials` table beside `approvers`, and the consume
    advancing the sign counter in the mint's own transaction.
  - The CLI: `approver enroll --format webauthn` (enrolment code + typed `y`), and `approve <id>
    --webauthn-link` / `--webauthn-assertion <file>`, recording `approved_by = webauthn:<name>`.
  - **One static ceremony page**, `GET /surface/passkey` and its script. It calls
    `navigator.credentials.create`/`get` and lets the operator save the result. It has no POST and no
    `/v1` route, and it reads no store.
  It changes no `.ail` code.
- **Designer lane**: `claude:claude-opus-5-5`.
- **Estimated**: 8 milestones in three sprints, about **5,200 changed lines** in all (§9): sprint 1
  ≈ 1,800, sprint 2 ≈ 1,600, sprint 3 ≈ 1,800. **Calibration (V30):** A1b was estimated at about
  3,900 lines and landed 6,571 over four sprints (1.7×). These estimates are sized from A1b's
  *actual* per-area sizes, not from its estimates.
- **Measurement base**: `98374909` (= origin/dev). Worktree
  `/Users/voightkampff/dev/sunholo-data/.wt-world-iter379-row201-a2`, branch
  `sprint/world-iter379-row201-a2`. sha256 prefixes (V1): `approval/approver.go` `1c69b80407e05229`,
  `approval/signed.go` `483a713329be77a8`, `approval/sshsig.go` `a40850200782b3e4`,
  `approval/request.go` `cf9897ec3f10881e`, `store/approvers.go` `e927ecdb5c691df6`,
  `store/mint_requests.go` `c130d108a3c9c5be`, `store/credentials.go` `be329ce3e2a2772c`,
  `authority/mint.go` `3bd26c47a976e740`, `cmd/…/approve.go` `7d0647d65a4be07e`,
  `cmd/…/approve_sig.go` `446c0f6df1c49ad3`, `cmd/…/approver.go` `1ab536f84fece934`,
  `host/sshsign/sshsign.go` `5973a5e645aded81`.
- **Instruments**: `go1.26.6 darwin/arm64`; `Google Chrome 157.0.8103.50`, driven headless through
  the DevTools protocol with a **virtual authenticator** (a test key, never a real one); `node
  v26.13.0`; pinned `AILANG v0.41.0` at `$HOME/.pinned-ailang/ailang`. Throwaway probes lived in
  `/tmp/w379-designer/` and are not in the repo (V20–V24). The W3C WebAuthn texts were read at
  `https://www.w3.org/TR/webauthn-3/` (L3) and `https://www.w3.org/TR/webauthn-2/` (L2) for
  section and step numbers (V25). The designer's shell has no controlling terminal (V29).

---

## 0. Scope, and the sprint boundaries

**In A2:**
- **The verification law (pure, sprint 1).** A bounded CBOR subset (§2.3), COSE keys for ES256
  and EdDSA (§2.4), authenticator data (§2.5), `clientDataJSON` (§2.6), the response envelope
  (§2.7), the registration law with its enrolment code (§2.8), and the assertion law (§2.9). The
  test authenticator `approvaltest.Passkey` and fixtures recorded from a real Chrome with a virtual
  authenticator are also in sprint 1, because every later test needs them (§2.15).
- **The store and the binding (sprint 2).** Credentials schema 6 adds `webauthn_credentials`
  (§2.10). The consume binds the passkey's credential ID and advances its sign counter in the
  mint's transaction (C13–C18). `authority` derives `approved_by = webauthn:<name>` from the
  verified approval (§2.11).
- **The operator paths (sprint 3).** `approver enroll --format webauthn --origin <o>
  --registration <file>` (§2.12); `approver list`/`revoke` show passkeys; `approve <id>
  --webauthn-link` and `--webauthn-assertion <file>` (§2.13); and the static ceremony page
  `GET /surface/passkey` (§2.14), which is what lets Mark's attended exit run use a real passkey.
- **Carried forward from the A1b judge (iteration 378), as test-only items in sprint 1 M0**:
  **K8-ERRCODE**, where `approver revoke`'s post-write listing failures return `exitFatal` unpinned
  (`approver.go:449` and `:455`; both mutants to `exitOK` survive `TestApprover*`, V3); and
  **sshsign G1**, where `host/sshsign` has no unit test for a failed signer run (the mutant survives
  `./host/sshsign/`, V4).

**Out of A2 (A3 and later):**
- the surface's *needs you* card and request card;
- an Approve button that POSTs an assertion to the daemon (A3: "a POST bound to a verified
  signature", parent §3);
- the native notification and the Tauri Touch ID convenience;
- the graphical decision list;
- `Human.Approve` packets (A4).
In A2 the browser hands its result to the operator **as a file**, and the CLI consumes it. So
A2 adds no write route.

**Sprint boundaries.**
- **Sprint 1** is pure: no store, no CLI, no route. It can be judged on its own.
- **Sprint 2** depends on sprint 1's types. It changes how a mint is consumed, but no operator
  path reaches the new branch yet.
- **Sprint 3** depends on both. It makes enrolment and approval reachable.
- Nothing flows backwards. This follows A1b §0's reasoning: the law is reviewed before anything
  can write.

## 1. What the code does today (measured at `98374909`)

1. **Exit codes** are unchanged from A1b §1.1: a refusal before the write is `exitUsage` (1); a
   failure at or after the write is `exitFatal` (2) (A1b V2; re-read `approve.go:117-121`, V5).
2. **The message a key signs** is `approval.SignedMessage(digest, challenge, expiry)`, 90 bytes:
   `"world-approval/v1" ‖ 0x00 ‖ digest ‖ challenge ‖ expiry (u64 BE)` (`request.go:209-222`).
   **`approval.MessageHash` already exists** (`request.go:224-231`). It is SHA-256 of that message,
   pinned by `fixtureMessageHashHex = "29181d22…fef5"` (`request_test.go:35`, `:295-300`), and
   **no production code calls it yet** (V6).
3. **The enrolment law refuses WebAuthn, naming A2.** `FormatWebAuthn = "webauthn"`
   (`approver.go:18`); `NewApprover` E2 returns `ErrFormatNotYet` (`:53-55`). It is pinned by
   `approver_test.go:150-151` and the CLI arm `approver_enroll_test.go:230-232` (stderr contains
   `row 201 A2`) (V7). `Approver` is `{Name, Format, Key SSHPublicKey, EnrolledAt, LastUsedAt,
   RevokedAt}` (`:28-35`).
4. **The signature law** (`signed.go`):
   - `SignedApproval{ok, digest, challenge, expiry, name, format, key}` (`:16-24`).
   - `ApproveBySignature` runs P1–P6 (`:55-77`). P3 refuses any format but `ssh` (`:62`).
   - `Binds` (B1–B4), `ApprovedBy` (B5) and `Use` (B6) (`:80-112`).
   - RC1 reflects over it (`signed_test.go:177`) (V8).
5. **The store** (V9):
   - `credentialsSchemaVersion = 5` (`credentials.go:42`).
   - The table set per version is `credentialsTables` (`:71-82`). Version 5 holds `approvers`,
     `episode_starts`, `mint_requests` and `session_credentials`.
   - `upgradeCredentials` ends with `if version < 5 { approversSQL }` (`:349-360`).
   - A bound, started or resumed mint refuses on any store below the binary's version (`:437`).
   - `approvers` (`credentials_schema.sql:63-70`) has `public_key BLOB NOT NULL UNIQUE`, and
     `format` is checked in code only.
   - `EnrolApprover` refuses any format but `ssh` (`approvers.go:89`). It has `Approver`,
     `ApproverByKey`, `ListApprovers` and `RevokeApprover`.
   - `useApprover` (C8) and `classifyApproverUse` (C9–C12) are at `:245-278`.
   - The credentials DB sets no `PRAGMA foreign_keys`. Only the world store does (`store.go:397`),
     so a `REFERENCES` clause would not be enforced (V10).
6. **The consume** (`mint_requests.go:284-360`):
   - `RequestConsume{ID, Digest, ApprovedBy, Challenge, Approver *ApproverUse, Clock, Check}`.
   - `ApproverUse{Name, Format, PublicKey}`.
   - C1–C5 run before the transaction, in this order: completeness, approver completeness,
     `tty`⇔nil, `format:name`, schema 5.
   - Inside it: the guarded request `UPDATE` (C6), then `useApprover` (C8), then the session insert.
   - **Two existing test arms set `Format = "webauthn"` on an ssh approver**
     (`consume_binding_test.go:138`, expecting `ErrInconsistentApproval` from C4; `:289`, expecting
     `ErrApproverKeyChanged` from C11). §2.10 orders the new checks so the first stays green
     unedited, and names the one edit the second needs (V11).
7. **Authority** (`mint.go:79-103`):
   - `requestConsume` refuses a caller-set `ApprovedBy`/`Approver` (A1) and a bad digest (A2).
   - It sets `tty` when `Signed == nil` (A5). Otherwise it calls `Binds` (A4), then copies
     `Use()` into `ApproverUse` and `ApprovedBy()` into `ApprovedBy` (A3).
   - `RequestMint.Signed` is `*approval.SignedApproval` (`:67`).
8. **`approve`** (`approve.go:117-307`):
   - It has three exclusive modes, `--signature`, `--sign-with` and `--write-message`. F1 is the
     `modes > 1` check (`:147-156`).
   - It admits the row, decodes the challenge (F4, `:225-228`), then judges the signature before
     any network call (`:235-250`).
   - Steps 2–5 re-derive the request. Step 6 sets `RequestMint.Signed` (`:298-299`), and
     `provisionEpisode` skips the terminal when `Signed != nil` (A1b §2.8).
   - `approveSignature` reads its file through `readRegularFileCapped` (open non-blocking, `fstat`
     the descriptor, `LimitReader`) (`approver.go:84-114`) (V5).
9. **`approver`** (`approver.go`):
   - `confirmApproverAct` opens `/dev/tty`, writes the prompt and reads one line through the
     error-aware `readTermLineErrWithin`, bounded by `approverConfirmBudget = 120 s`. It confirms
     only on a complete `y` line with no read error (`:27-77`).
   - `runApproverEnroll` runs N1–N14 (`:279-350`).
   - `approverKeyLabel` assumes an SSH blob (`:261-267`), so a passkey row would print as
     `invalid-key SHA256:…`.
   - `runApproverRevoke`'s K8 listing returns `exitFatal` on a failed read-only reopen (`:449`)
     and on a failed `DecidedMintRequests` (`:455`). Neither is pinned (V3).
10. **`host/sshsign.Sign`** has G1 (`cmd.Run` error), G2 (over 16 KiB) and G3 (empty)
    (`sshsign.go:53-61`). Its unit tests are `TestCappedBufferKeepsItsBound` and
    `TestSignHandsTheSignerAScrubbedEnvironment`. Neither runs a failing signer, and the G1 mutant
    survives them; only the real-tool `cmd` tests kill it (V4).
11. **The daemon** binds loopback only (`DefaultBindHost = "127.0.0.1"`, `daemon.go:163`; checked
    at `:713`).
    - It serves `/surface` and its assets under `surfaceCSP = "default-src 'none'; …; script-src
      'self'; connect-src 'self'; …; frame-ancestors 'none'"` (`surface_page.go:19`).
    - It sets no `Permissions-Policy` header, so the WebAuthn features keep their default `self`
      allowlist.
    - It does no `Host`-header check: the only `r.Host` read is `session_self.go:69` (V12).
    - `TestSurfaceAppRoutes` pins the `/surface/` asset set ("every other name under /surface/ is
      a 404", `surface_app_test.go:17-18`).
    - `TestFrozenV1RouteTableMatchesMux` pins 13 `/v1` routes (`route_table_test.go:51-53`). A2
      adds none (V13).
12. **The live-Chrome pattern exists.** `WORLD_CHROME`-gated tests launch the operator's Chrome
    headless (`workbench_live_chrome_test.go:68-73`). They drive pages by HTTP reports and **do
    not use the DevTools protocol**, so none can add a virtual authenticator. The AC10 census skips
    `_test.go` files (`registry_publish_test.go:1094`), so a test-only launcher adds no production
    exec site (V14).
13. **Go 1.26 randomness.** `ecdsa.GenerateKey` and `ecdsa.SignASN1` ignore their `io.Reader`
    since Go 1.26. Deterministic tests use `testing/cryptotest.SetGlobalRandom`, which "cannot be
    used in parallel tests". `ecdsa.ParseRawPrivateKey` builds a P-256 key from fixed bytes, and
    `ecdsa.ParseUncompressedPublicKey` refuses an off-curve point (V15, V23).
14. **Nothing of A2 exists** (V16). No production file parses CBOR, COSE or `clientDataJSON`, and
    none mentions `webauthn_credentials`, `sign_count` or `/surface/passkey`. Controls in the same
    command find `FormatWebAuthn` and `MessageHash`.
15. **A1b's seam note** (A1b §2.10, `w-world-approvals-by-key-a1b.md:829-833`) planned
    `sign_count` "by `ALTER TABLE ADD COLUMN` at schema 6". **This design departs from that
    (§2.10):** a separate table can carry `UNIQUE(credential_id)`, which `ADD COLUMN` cannot. It
    also keeps the existing canonical-DDL test's one-table-per-version pattern.

## 2. Design

### 2.1 What a passkey signs: the challenge is the message hash

**Rule.** The WebAuthn `challenge` of an approval assertion is
`approval.MessageHash(digest, challenge, expiry)`: the SHA-256 of the existing 90-byte
`SignedMessage`. In `clientDataJSON` it appears as the unpadded base64url of those 32 bytes (L3
§7.1 step 8 and L2 §7.2 step 12: "equals the base64url encoding of … challenge") (V25).

**Why this, and not something else:**
- **The parent fixes it.** Parent §2.2 says: "For WebAuthn, the message hash is the challenge."
- **It binds all three values.** The request row's 32-byte challenge alone would not cover the
  digest or the expiry. The message hash covers digest, challenge and expiry exactly as the SSHSIG
  message does, so `Binds` (B2–B4) means the same thing for both formats.
- **It is already pinned.** `MessageHash` and its golden exist (§1.2).
- **It was measured with a real browser.** Chrome 157 was given this exact fixture hash as
  `challenge`. Its `clientDataJSON` carried
  `"challenge":"KRgdIk2YL2WdG9RwMjjT-EieJpaP_ioivXhvo78V_vU"`, the base64url of
  `29181d22…fef5` (V21). The signature covers it: the assertion signature is over
  `authenticatorData ‖ SHA-256(clientDataJSON)` (L2 §7.2 step 20), and the measured ES256 and
  EdDSA signatures verified over exactly that with `crypto/ecdsa.VerifyASN1` and
  `ed25519.Verify` (V22).
- **The comparison is a string equality** with
  `base64.RawURLEncoding.EncodeToString(hash[:])`. That is the spec's own comparison, and it leaves
  no decoding ambiguity: padded or non-canonical forms simply do not match (V23).

**Enrolment's challenge** is 32 fresh bytes from `env.random`, held in the `approver enroll`
process and in the link it prints (§2.12). It is never stored. The enrolment is one process, so
there is nothing to replay across processes.

### 2.2 Origin and relying-party ID

`ParseCeremonyOrigin(s string) (rpID string, err error)` is pure (`host/approval/origin.go`). It
accepts exactly:

| origin | rpID |
|---|---|
| `https://<host>.ts.net` or `https://<host>.ts.net:<port>`, every label `[a-z0-9]([a-z0-9-]{0,61}[a-z0-9])?`, at most 253 characters | the host |
| `http://localhost`, `http://localhost:<port>`, `https://localhost`, `https://localhost:<port>` | `localhost` |

- **[r2]** A port is 1–65535 with no leading zero, and must not be the scheme's default (443 for https, 80 for http); the page compares against `location.origin`, which elides the default port, so an explicit default port could never enrol. **O3** covers a bad or default port → `ErrOrigin`. (Quorum r2, `oc-kimi-k3`, applied verbatim; measured V36.) OR1 gains refusal arms `https://x.ts.net:443`, `https://localhost:443` and `http://localhost:80` → `ErrOrigin`, each with guidance naming the default-port rule.
- Everything else is refused: uppercase, a path or a trailing `/`, userinfo, a query, any other
  domain, and **any IP literal**. L3 §5.1.3 step 7 says "Only the domain format of host is allowed
  here" (V25).
- **O1** grammar → `ErrOrigin`; **O2** a host that is neither `localhost` nor `*.ts.net` →
  `ErrOriginHost`; **O3** a bad or default port → `ErrOrigin` **[r2]**.
- **rpID = host, exactly.** The spec would also allow a registrable suffix. A2 is stricter: one
  origin enrols one rpID, and parent §4 ("the surface is the only origin enrolled") is literal.
- **Consequence.** The daemon serves on `127.0.0.1:7644`, and `http://127.0.0.1:7644` cannot be an
  origin. The operator opens `http://localhost:7644/surface/passkey`, or the tailnet name in front of
  it (§10). QUICKSTART §16 says so.
- **Both are stored.** The rpID and the full origin are stored per credential at enrolment (§2.10).
  An assertion must match both **exactly** (W7, W10).

### 2.3 The bounded CBOR subset (`host/approval/cbor.go`)

A2 decodes two CBOR documents: the `attestationObject` (RFC 8949 major types 2, 3 and 5), and the
COSE key and the optional extensions inside authenticator data (types 0, 1, 2, 3, 5; extension
values may also use 4 and the simple values `false`/`true`/`null`). The decoder is a recursive
reader over a byte slice:

```go
type cborLimits struct{ maxDepth, maxItemsPerContainer, maxTotalItems, maxStringBytes int }
var attestationLimits = cborLimits{maxDepth: 4, maxItemsPerContainer: 16, maxTotalItems: 64, maxStringBytes: 2048}

// decodeCBOR reads one item and returns it and the bytes after it.
func decodeCBOR(b []byte, lim cborLimits) (cborValue, []byte, error)
// decodeCBORExact reads one item and refuses trailing bytes.
func decodeCBORExact(b []byte, lim cborLimits) (cborValue, error)
```

Each refusal is its own branch, with a typed error:

| id | refusal | error |
|---|---|---|
| K1 | input ends inside a head or argument | `ErrCBORTruncated` |
| K2 | additional information 28–30 (reserved, RFC 8949 §3) | `ErrCBORMalformed` |
| K3 | additional information 31 (indefinite length, §3.2) | `ErrCBORIndefinite` |
| K4 | an argument not in its shortest form (preferred serialization, §4.2.1; CTAP2 canonical) | `ErrCBORNonCanonical` |
| K5 | a string or container length beyond the remaining bytes | `ErrCBORTruncated` |
| K6 | a byte or text string over `maxStringBytes` | `ErrCBORTooLarge` |
| K7 | a map or array over `maxItemsPerContainer` entries | `ErrCBORTooLarge` |
| K8 | nesting over `maxDepth` | `ErrCBORTooDeep` |
| K9 | major type 6 (any tag) | `ErrCBORUnsupported` |
| K10 | major type 7 other than 20/21/22 (`false`/`true`/`null`): floats, `undefined`, other simple values, `break` | `ErrCBORUnsupported` |
| K11 | a text string that is not valid UTF-8 | `ErrCBORMalformed` |
| K12 | a map key that is not an integer or a text string | `ErrCBORUnsupported` |
| K13 | a duplicate map key (RFC 8949 §5.6) | `ErrCBORDuplicateKey` |
| K14 | an integer outside `int64` (major 0 argument > 2⁶³−1, or major 1 argument ≥ 2⁶³) | `ErrCBORTooLarge` |
| K15 | more than `maxTotalItems` items in one document | `ErrCBORTooLarge` |
| K16 | `decodeCBORExact`: trailing bytes | `ErrCBORTrailing` |

- **Map key order is not checked.** Chrome's output is in CTAP2 canonical order, `fmt`, `attStmt`,
  `authData`, and COSE `1, 3, −1, −2, −3` (V21). But the outer `attestationObject` is assembled by
  the browser, and other browsers' order is UNVERIFIED. Order carries no security weight: what is
  stored is a **re-encoded** canonical COSE key (§2.4), and every trusted byte is checked through a
  hash or a signature. Duplicates *are* refused (K13), so a key cannot carry two meanings.
- **K4 is enforced.** CTAP2 requires canonical encoding of what the authenticator emits, and
  Chrome's measured bytes conform. A browser that emitted non-shortest heads would be refused; the
  attended exit run is where that would surface (residual R3).
- **Bounds.** A credential ID is at most 1023 bytes (§2.5), so `authData` stays well under 2,048.
  Chrome's measured `attestationObject` is 194 bytes (ES256) and 159 (EdDSA), and its COSE keys are
  77 and 42 bytes (V21).

### 2.4 COSE keys: ES256 and EdDSA only (`host/approval/cose.go`)

```go
type COSEKey struct {
	Alg int64  // -7 (ES256) or -8 (EdDSA)
	X   []byte // 32 bytes
	Y   []byte // 32 bytes for ES256; nil for EdDSA
}
func ParseCOSEKey(v cborValue) (COSEKey, error)
func (k COSEKey) Canonical() []byte                       // CTAP2 canonical re-encoding: what is stored
func (k COSEKey) verify(signed, sig []byte) error          // ES256: VerifyASN1 over SHA-256(signed); EdDSA: ed25519.Verify
```

| id | refusal | error |
|---|---|---|
| Q1 | not a map | `ErrCOSEKey` |
| Q2 | `kty` (1), `alg` (3) or `crv` (−1) missing or not an integer | `ErrCOSEKey` |
| Q3 | `alg` not −7 or −8 (for example −257 RS256 or −35 ES384) | `ErrCOSEAlg` |
| Q4 | `alg` −7 with `kty` ≠ 2 (EC2) or `crv` ≠ 1 (P-256) | `ErrCOSEKey` |
| Q5 | `alg` −8 with `kty` ≠ 1 (OKP) or `crv` ≠ 6 (Ed25519) | `ErrCOSEKey` |
| Q6 | `x` (−2), or for ES256 `y` (−3), not a 32-byte string | `ErrCOSEKey` |
| Q7 | ES256: `0x04 ‖ x ‖ y` refused by `ecdsa.ParseUncompressedPublicKey` (off the curve) | `ErrCOSEKey` |
| Q8 | any other label (EdDSA's `y` included) | `ErrCOSEKey` |

- **Which algorithms, and why.**
  - ES256 is the one algorithm every platform passkey provider offers, and it is what Chrome's
    authenticator made when offered `[-7]` (V21).
  - EdDSA (Ed25519) is accepted because the standard library verifies it, the SSH side already
    trusts Ed25519, and Chrome's authenticator makes it when offered `[-8]` (V21).
  - Both are exactly what parent §2.1 names (`crypto/ecdsa` P-256, `crypto/ed25519`).
  - **RS256 (−257) is refused in A2.** It would need `crypto/rsa`, which the parent did not list.
    Some Windows Hello TPMs are reported to offer only RS256 (UNVERIFIED in this session). That is
    residual R4 and non-blocking question Q1.
- **The labels and values are measured, not recalled.** Chrome's ES256 key is
  `a5 01 02 03 26 20 01 21 58 20 …`: kty 2, alg −7, crv 1, x, y. Its EdDSA key is
  `a4 01 01 03 27 20 06 21 58 20 …`: kty 1, alg −8, crv 6, x (V21, V23). The section citations,
  RFC 9052 §7 (key objects) and RFC 9053 §2.1/§2.2 (ECDSA/EdDSA) and §7.1/§7.2 (EC2/OKP
  parameters), are cited from memory and **not re-fetched** (UNVERIFIED as citations, not as values).
- **ES256 signatures** are ASN.1 DER. Chrome's was 72 bytes starting `0x30` (V22). A signature over
  72 bytes is refused before verification. An Ed25519 signature must be 64 bytes.
- **`Canonical()`** re-encodes in CTAP2 order. It is byte-equal to Chrome's key bytes for both
  algorithms (test CK3). The stored `approvers.public_key` of a passkey is `Canonical()`, so
  `UNIQUE(public_key)` cannot be dodged by re-encoding the same key.

### 2.5 Authenticator data (`host/approval/authdata.go`)

```go
type AuthData struct {
	RPIDHash      [32]byte
	UP, UV, BE, BS bool       // flags bits 0, 2, 3, 4 (L3 §6.1)
	SignCount     uint32
	CredentialID  []byte     // registration only
	AAGUID        [16]byte   // registration only
	Key           COSEKey    // registration only
}
func ParseAuthData(b []byte, registration bool) (AuthData, error)
```

The layout is L3 §6.1: `rpIdHash` (32) ‖ `flags` (1) ‖ `signCount` (4, big-endian), then, if AT,
`aaguid` (16) ‖ `credentialIdLength` (2) ‖ `credentialId` ‖ the COSE key, then, if ED, an
extensions map. The flag bits are from L3 §6.1, as read (V25): UP bit 0, UV bit 2, BE bit 3,
BS bit 4, AT bit 6, ED bit 7.

| id | refusal | error |
|---|---|---|
| D1 | under 37 bytes | `ErrAuthDataShort` |
| D2 | BS set while BE is not (L3 §7.1 step 17) | `ErrAuthDataFlags` |
| D3 | registration without AT, or an assertion with AT | `ErrAuthDataFlags` |
| D4 | registration: truncated inside `aaguid`, the length or the credential ID | `ErrAuthDataShort` |
| D5 **[r1]** | an empty credential ID, or one over 1023 bytes (L3 §7.1 step 25) | `ErrCredentialID` |
| D6 | the COSE key does not decode (K1–K15) or parse (Q1–Q8) | wrapped |
| D7 | ED set but what follows is not one CBOR map | `ErrAuthDataExtensions` |
| D8 | trailing bytes after the last expected part | `ErrAuthDataTrailing` |

- **The RFU bits (1 and 5) are ignored**, as the spec reserves them.
- **Extension contents are ignored.** They are covered by the signature (assertion) or the
  enrolment code (registration). A2 requests no extension, but some authenticators add one
  unasked (UNVERIFIED), and refusing it would refuse a real passkey.
- **[r1] No 16-byte floor.** r0 refused credential IDs under 16 bytes, citing L3's definition as
  UNVERIFIED. Read in r1 (V34), L3 §4 says a credential ID is "At most 1023 bytes long". It then
  says "Credential IDs are generated by authenticators in two forms". Only the first form is "At
  least 16 bytes that include at least 100 bits of entropy". The second form, an encrypted
  credential source, has no stated minimum, and "Relying Parties do not need to distinguish these
  two Credential ID forms". So 16 bytes is a rule for authenticators generating IDs, not a
  relying-party check. The only relying-party length step is L3 §7.1 step 25 (≤ 1023). D5,
  `EnrolPasskey`'s validation, the DDL `CHECK`, C14 and the test arms now use **1–1023**. Chrome's
  IDs are 32 bytes (V21).
- **Measured:** Chrome's registration `authData` had flags `0x45` (UP|UV|AT) and signCount 1, with
  AAGUID `01020304…0708` (the virtual authenticator's). Its assertions had flags `0x05` (UP|UV),
  were 37 bytes, and counted 2 and then 3 (V21).

### 2.6 `clientDataJSON` (`host/approval/clientdata.go`)

```go
type ClientData struct {
	Type, Challenge, Origin string
	CrossOrigin             bool
	HasTopOrigin            bool
}
func ParseClientData(b []byte) (ClientData, error)
```

`ParseClientData` walks `encoding/json` tokens. It does not `Unmarshal` into a struct, because
that keeps the last of two duplicate keys silently.

| id | refusal | error |
|---|---|---|
| J1 | empty, or over 2 KiB | `ErrClientData` |
| J2 | not valid UTF-8 | `ErrClientData` |
| J3 | not exactly one JSON object, or trailing data | `ErrClientData` |
| J4 | `type`, `challenge`, `origin`, `crossOrigin` or `topOrigin` appears twice | `ErrClientDataDuplicate` |
| J5 | `type`, `challenge` or `origin` missing or not a string | `ErrClientData` |
| J6 | `crossOrigin` present and not a boolean; `topOrigin` present and not a string | `ErrClientData` |
| J7 | an unknown member nested over depth 4 | `ErrClientData` |

- **Unknown members are allowed** (bounded by J1 and J7). L3 §5.8.1 lets clients add members, and
  the L3 §5.8.1.2 "limited verification algorithm" is a prefix match over the spec's serialization
  order. A2 parses instead of prefix-matching, so the result does not depend on member order. Chrome
  is reported to inject a random `other_keys_can_be_added_here` member sometimes (UNVERIFIED: absent
  in both measured runs). An arm in CD1 includes one.
- **Measured:** Chrome emitted
  `{"type":"webauthn.create","challenge":"…","origin":"http://localhost:56547","crossOrigin":false}`
  and `{"type":"webauthn.get",…}` with the same member order (V21).

### 2.7 The response files (`host/approval/passkey_json.go`)

The browser's output reaches the CLI as a JSON file. The page (§2.14) writes it, and it is also
exactly the subset of L3's `PublicKeyCredential.toJSON()` that Chrome 157 emits (V21):

```json
{"id":"<b64url>","rawId":"<b64url>","type":"public-key",
 "response":{"clientDataJSON":"<b64url>","attestationObject":"<b64url>"}}                  // registration
{"id":"<b64url>","rawId":"<b64url>","type":"public-key",
 "response":{"clientDataJSON":"<b64url>","authenticatorData":"<b64url>","signature":"<b64url>","userHandle":"<b64url>"|null}}  // assertion
```

- **Every field is decoded with `base64.RawURLEncoding.Strict()`.** A padded or non-canonical
  string is refused (measured, V23).
- **Unknown members are ignored.** `toJSON()` adds `publicKey`, `publicKeyAlgorithm`,
  `transports`, `authenticatorAttachment` and `clientExtensionResults` (V21). None of them is
  trusted: the key comes from `authData`, never from `publicKey`.
- **Duplicate envelope keys are not refused**, because no envelope value is trusted unchecked:
  - `rawId` must equal `authData`'s credential ID at registration (R15).
  - At assertion, `rawId` only selects which enrolled key must verify the signature (W2).
  - The three byte fields are covered by the enrolment code or the signature.
- **E1** the file is over 16 KiB (checked by the CLI's `readRegularFileCapped`); **E2** the JSON
  does not decode, or a required member is missing → `ErrPasskeyResponse`; **E3** `type` ≠
  `public-key`, or `id` ≠ `rawId` → `ErrPasskeyResponse`; **E4** a field is not strict base64url →
  `ErrPasskeyResponse`.

### 2.8 The registration law and the enrolment code (`host/approval/passkey.go`)

```go
const FormatWebAuthn = "webauthn" // existing (approver.go:18)

type Passkey struct {
	RPID, Origin   string
	CredentialID   []byte
	Key            COSEKey
	AAGUID         [16]byte
	BackupEligible bool   // stored; an assertion must agree (W13)
	BackupState    bool   // shown at enrolment only (synced or device-bound); not stored
	SignCount      uint32 // the registration's counter
}

// Approver gains one field (nil for ssh approvers):
//	Passkey *Passkey

type RegistrationInput struct {
	Name, Origin string
	Challenge    [32]byte // the enrol process's fresh challenge
	Response     []byte   // the registration file
	Now          int64
}

func NewPasskeyApprover(in RegistrationInput) (Approver, error)
func EnrolmentCode(clientDataJSON, attestationObject []byte) string // [r1] 26 characters, 128 bits
func GroupEnrolmentCode(code string) string                         // [r1] "XXXXX-XXXXX-XXXXX-XXXXX-XXXXXX"
func RenderPasskeyEnrolment(a Approver) string
```

`NewPasskeyApprover`'s checks run in this order. The first that fails is returned:

| id | check | spec | error |
|---|---|---|---|
| R1 | the name grammar | — | `ErrApproverName` |
| R2 | `ParseCeremonyOrigin(Origin)` → rpID (O1–O3) | §2.2 | `ErrOrigin*` |
| R3 | the envelope (E2–E4) | — | `ErrPasskeyResponse` |
| R4 | `ParseClientData` (J1–J7) | — | `ErrClientData*` |
| R5 | `type == "webauthn.create"` | L3 §7.1 step 7 | `ErrClientDataType` |
| R6 | `challenge == b64url(Challenge)` | step 8 | `ErrChallenge` |
| R7 | `origin == Origin` | step 9 | `ErrOriginMismatch` |
| R8 | `crossOrigin` is not `true`, and there is no `topOrigin` | steps 10–11 | `ErrCrossOrigin` |
| R9 | the `attestationObject` is a map of exactly `fmt` (text), `attStmt` (map) and `authData` (bytes) (K1–K16) | — | `ErrAttestationObject` |
| R10 | `fmt == "none"` and `attStmt` is the empty map | step 21; §8.7 | `ErrAttestationFormat` |
| R11 | `ParseAuthData(authData, true)` (D1–D8) | — | wrapped |
| R12 | `rpIdHash == SHA-256(rpID)` | step 14 | `ErrRPIDHash` |
| R13 | UP set | step 15 | `ErrUserPresence` |
| R14 | UV set | step 16 | `ErrUserVerification` |
| R15 | `rawId == authData.CredentialID` | — | `ErrCredentialIDMismatch` |

- **R16** (the algorithm is −7 or −8, step 20) is Q3, reached through R11.
- **Step 26** (the credential ID is not already registered) is the store's (`UNIQUE` and the
  `EnrolPasskey` pre-check, §2.10).
- **No attestation is verified** (parent §2.1). `fmt: none` is what the page requests, and what
  Chrome returned (`attStmt: {}`, V21).
- The step numbers above are from L3 §7.1 as read (V25).

**What `fmt: none` means for trust.** Without attestation, World cannot prove that the enrolled key
lives in hardware or that its UV flag is honest. A software authenticator can set UV. **The trust
comes from the enrolment ceremony.** The key is accepted only when the operator, at the daemon's
terminal, types the code his own browser showed for *that* registration and then types `y` (§2.12).
After that, an assertion's UV bit is the authenticator the operator chose vouching for user
verification, which is the parent's "a key that cannot be used without the person". Residual R1
states the rest.

**The enrolment code, [r1]: a 128-bit fingerprint, not a secret.**
```
EnrolmentCode = base32(RFC 4648 alphabet A–Z 2–7, uppercase, no padding) of the FIRST 16 BYTES of
                SHA-256("world-enrolment-code/v1" ‖ 0x00 ‖ SHA-256(clientDataJSON) ‖ SHA-256(attestationObject))
              = 26 characters (128 bits; the last character carries 3 bits)
```

- **[r1] The attacker model.** r0 claimed that an agent "would have to guess 40 bits blind". That
  is false, and quorum round 1 showed why. The code is computed only from `clientDataJSON` and
  `attestationObject`, and both are in the registration file. So **the registration response and
  its fingerprint are public to a file-substitution attacker.** Such an attacker can:
  1. read the genuine file the operator's browser saved;
  2. compute its code;
  3. grind registrations of its own. With `fmt: none` nothing is attested, so the attacker's
     key, credential ID and every unsigned field vary freely;
  4. look for one whose code equals the genuine code;
  5. replace the file before the CLI reads it.
  At 40 bits that is about 2⁴⁰ hashes of offline work.
- **[r1] What the code now needs, and has.** Enrolment requires a fingerprint with **at least 128
  bits of target-preimage resistance**, compared against the operator's browser through the
  terminal.
  - The attacker must find a registration `X ≠ G` with `code(X) = code(G)`, for a `G` that it did
    not choose. `G` carries the operator's authenticator's fresh key and credential ID. That is a
    second preimage on 128 bits of truncated SHA-256: about 2¹²⁸ work.
  - A birthday collision (about 2⁶⁴) does not help, because the attacker cannot choose the genuine
    registration's bytes.
  - The comparison is **single-attempt**. One code is typed per `approver enroll` process, and a
    mismatch refuses and exits (WN10). A substituted file therefore cannot be retried against the
    operator's typing.
- **[r1] What the comparison binds.** The value the operator transfers from his own browser to
  `/dev/tty` must equal the fingerprint of the exact bytes the CLI read. Those bytes are
  everything the registration law trusts. `rawId` is checked against `authData`'s credential ID
  (R15), so it is covered too.
  - A substituted file fails WN10, whether or not the attacker read the genuine one first.
  - The typed `y` follows, showing the verified key's details.
- **[r1] Display and typing (operator UX).**
  - The page shows the code in **five groups**, `XXXXX-XXXXX-XXXXX-XXXXX-XXXXXX` (5, 5, 5, 5 and
    6 characters). It is in a selectable, read-only monospace field. Measured on the V21
    registration: `YDOR4-X36PQ-NFRTX-BHWAV-4W2XG4` (V35).
  - The operator may type it, or select and paste it into the terminal: 30 characters, far below
    the terminal's line limit.
  - The CLI prompt shows the same grouping as a template.
  - Before comparing, the CLI **normalises**: it uppercases, then removes `-`, spaces and tabs.
    The result must be exactly 26 characters of `A–Z2–7`, or it is refused with WN10's message.
  - The base32 alphabet has no `0`, `1`, `8` or `9`, so `O`/`0` and `I`/`1` cannot be confused in
    what the operator types: a typed `0` or `1` is simply refused.
  - The comparison is `subtle.ConstantTimeCompare` over the normalised 26 bytes.
- The page computes the code with `crypto.subtle.digest`, which is available in the secure contexts
  `https:` and `http://localhost`. The CLI computes it from the file it read (WN10).
- The code has no branches. A golden test pins it against the committed Chrome fixture (PE4). The
  live browser test pins the page's computation and grouping against Go's (XB1).

`RenderPasskeyEnrolment` prints:
- the name, `webauthn`, the rpID and origin, and the algorithm (`ES256`/`EdDSA`);
- the credential ID as `cred <first 16 hex of SHA-256(credentialId)>`;
- the AAGUID in hex;
- `synced passkey` (BE and BS set), `backup-eligible` (BE only) or `device-bound`;
- the counter, and `approved_by=webauthn:<name>`;
- a fixed paragraph, golden-pinned:

> This passkey will approve mint requests with no terminal, each time after your device verifies
> you (Face ID, Touch ID, a PIN or a security-key touch). World cannot check which device holds
> it; the 26-character enrolment code you typed, which matched this exact registration, is what ties
> it to the browser you just used. Enrol only a passkey you created yourself just now.

`NewApprover` E2's message changes to "webauthn approvers enrol from a browser registration
(approver enroll --format webauthn --origin … --registration …)". `ErrFormatNotYet` is renamed
`ErrUseRegistration`, and the branch stays.

### 2.9 The assertion law (`host/approval/passkey_assert.go`)

```go
type AssertionInput struct {
	Admit     AdmitInput
	Challenge [32]byte // the request row's challenge
	Response  []byte   // the assertion file
	Approver  Approver // the enrolled row the assertion's credential ID names, with Passkey set
}
func ApproveByAssertion(in AssertionInput) (*SignedApproval, error)
```

The checks run in this order:

| id | check | spec | error |
|---|---|---|---|
| P1 | `Admit` | — | `*NotAdmittedError` |
| P2–P5 | name grammar; **format `webauthn` and `Passkey != nil`** (P3w); not revoked; `Now ≥ EnrolledAt` | — | as A1b |
| W1 | the envelope (E2–E4) | — | `ErrPasskeyResponse` |
| W2 | `rawId == Passkey.CredentialID` | L2 §7.2 steps 5–6, allowCredentials and ownership (step numbers UNVERIFIED) | `ErrCredentialIDMismatch` |
| W3 | a present, non-null `userHandle` equals `[]byte(Name)` (the user ID the page sends) | step 6 | `ErrUserHandle` |
| W4 | `ParseClientData` | — | — |
| W5 | `type == "webauthn.get"` | step 11 | `ErrClientDataType` |
| W6 | `challenge == b64url(MessageHash(StoredDigest, Challenge, ExpiresAt))` | step 12 | `ErrChallenge` |
| W7 | `origin == Passkey.Origin` | step 13 | `ErrOriginMismatch` |
| W8 | not `crossOrigin`; no `topOrigin` | L3 §7.2 (step number UNVERIFIED) | `ErrCrossOrigin` |
| W9 | `ParseAuthData(authenticatorData, false)` | — | wrapped |
| W10 | `rpIdHash == SHA-256(Passkey.RPID)` | step 15 | `ErrRPIDHash` |
| W11 | UP set | step 16 | `ErrUserPresence` |
| W12 | UV set (user verification is required, parent R2) | step 17 | `ErrUserVerification` |
| W13 | `BE == Passkey.BackupEligible` | L3 §7.2 (step number UNVERIFIED) | `ErrBackupEligibility` |
| W14 | `Key.verify(authenticatorData ‖ SHA-256(clientDataJSON), signature)` | step 20 | `ErrBadSignature` |
| W15 | if `new ≠ 0 || stored ≠ 0` then `new > stored` | step 21; L3 §6.1.1 | `ErrSignCountRegression` |

The §7.2 step numbers are L2's, as read (V25).

- **Order.** Admission and the approver's state come first, as in A1b §2.3, so a replay says
  "already used" and a revoked key says "revoked". The credential ID (W2) comes before any parsing
  of signed content. The counter (W15) comes last, so it is judged only on a genuine signature.
- **P2–P5 are shared.** They move into one unexported `checkApprover(admit, a, format)` called by
  both laws, so a mutant there is killed by both laws' tests.
- **The counter (W15).** "Both zero" passes and updates nothing. That is the spec's rule for
  authenticators without a counter (L2 §7.2 step 21; L3 §6.1.1, V25); synced passkeys are reported
  to always send 0 (UNVERIFIED). Otherwise the new value must be strictly greater. On failure the
  spec says "This is a signal that the authenticator may be cloned", and leaves the decision to the
  relying party. **World refuses, and never auto-revokes** (Q2). The refusal names
  `approver revoke <name>` as the fix if the operator did not cause it. The store re-checks it
  atomically (C18).
- **`SignedApproval` gains two unexported fields**, `credentialID []byte` and `signCount uint32`.
  It also gains `PasskeyUse() (credentialID []byte, signCount uint32, ok bool)`, which is `ok` only
  when the value is verified and its format is `webauthn` (**B7**: unverified → false; **B8**: ssh →
  false). `Binds`, `ApprovedBy` and `Use` are unchanged, so `ApprovedBy()` is `webauthn:<name>`
  and `Use()`'s key is the canonical COSE key. RC1 must classify both new fields.

### 2.10 Store: credentials schema 6 (`host/store/passkeys.go`)

```sql
-- Passkeys (version 6, row 201 A2). One row per webauthn approver (approvers.name, whose format is
-- 'webauthn' and whose public_key is the canonical COSE key). The credential ID is never enrolled
-- twice; the counter only moves forward, inside the consume that uses it. The credentials store
-- does not enforce foreign keys, so EnrolPasskey writes both rows in one transaction.
CREATE TABLE IF NOT EXISTS webauthn_credentials (
    name            TEXT PRIMARY KEY CHECK (length(name) BETWEEN 1 AND 48),
    credential_id   BLOB NOT NULL UNIQUE CHECK (length(credential_id) BETWEEN 1 AND 1023),
    rp_id           TEXT NOT NULL CHECK (rp_id <> ''),
    origin          TEXT NOT NULL CHECK (origin <> ''),
    alg             INTEGER NOT NULL CHECK (alg IN (-7, -8)),
    aaguid          BLOB NOT NULL CHECK (length(aaguid) = 16),
    backup_eligible INTEGER NOT NULL CHECK (backup_eligible IN (0, 1)),
    sign_count      INTEGER NOT NULL CHECK (sign_count BETWEEN 0 AND 4294967295)
);
```

- **The upgrade.** `credentialsSchemaVersion` goes 5 → 6. `credentialsTables(6)` is `["approvers",
  "episode_starts", "mint_requests", "session_credentials", "webauthn_credentials"]`, sorted by
  name. `upgradeCredentials` gains `if version < 6 { passkeysSQL }`, byte-identical (normalised) to
  the schema file under `TestCredentialsSchemaMatchesCanonicalDDL`.
- **Why a table and not columns.** SQLite's `ADD COLUMN` cannot add a `UNIQUE` column, and
  `credential_id` must be unique (L3 §7.1 step 26). A table also keeps every SSH row free of
  passkey-only nullable columns.
- **The cost.** As with A1b's R2, a v5 store refuses bound, started and resumed mints from this
  build until `serve` restarts once (`credentials.go:437`). `episode start` already forces that
  restart through its same-build check.

```go
const passkeysSchemaVersion = 6
type PasskeyRow struct {
	Name, RPID, Origin string
	CredentialID, AAGUID []byte
	Alg                  int64
	BackupEligible       bool
	SignCount            int64 // 0 .. 2^32-1
}
type PasskeyUse struct {
	CredentialID []byte
	SignCount    int64 // the assertion's counter
}
// ApproverUse gains:   Passkey *PasskeyUse // non-nil exactly when Format == "webauthn"
var (
	ErrPasskeyCredentialEnrolled = errors.New("store: that passkey credential is already enrolled (a credential is never enrolled twice, even after revoke)")
	ErrSignCountRegression       = errors.New("store: the passkey's sign counter did not advance (a cloned authenticator, or an older assertion)")
	ErrPasskeyChanged            = errors.New("store: the approver's passkey is not the one that signed")
)
```

| Method | What it does |
|---|---|
| `EnrolPasskey(ctx, a ApproverRow, p PasskeyRow)` | Requires a deadline. Below v6 → `*OlderCredentialsStoreError{Need: 6, Feature: "passkeys"}`. **Validation, each its own `ErrBadApprover` branch (S1–S9):** `a.Format != "webauthn"`; `a.Name != p.Name`; an empty name or `PublicKey`; `EnrolledAt <= 0`; a non-zero `LastUsedAt`/`RevokedAt`; a credential ID outside 1–1023 bytes **[r1]**; `Alg` not −7/−8; an `AAGUID` not 16 bytes; an empty `RPID`/`Origin`; `SignCount` outside 0..2³²−1. **In one transaction:** a taken name → `ErrApproverNameTaken`; an enrolled key → `ErrApproverKeyEnrolled`; an enrolled credential ID → `ErrPasskeyCredentialEnrolled` (each naming the holder); then both `INSERT`s. Driver error text is never parsed |
| `PasskeyByCredentialID(ctx, id)` | `(ApproverRow, PasskeyRow, ok, err)`, joined on `name`. Below v6 → `Older` |
| `Passkey(ctx, name)` | `(PasskeyRow, ok, err)`. Below v6 → `Older` |
| `ActivePasskeys(ctx)` | the active webauthn approvers with their passkey rows, ordered by name, capped at 64 (for `--webauthn-link`). Below v6 → `Older` |

`EnrolApprover` keeps refusing every format but `ssh` (`approvers.go:89`, unchanged), so a
passkey can only enter through `EnrolPasskey`.

**The consume** (`MintSessionForRequest`) gains three pre-transaction checks. **[r1] Where they
sit, exactly.** Today's order before `BeginTx` is C1 (`:306`), C2 (`:309`), C3 (`:312`), C4
(`:315`), the session-row check (`:318`), the mint-request version check (`:321`) and C5 (`:324`),
then `BeginTx` (`:327`) (V11). **C13, C14 and C15 are inserted in that order immediately after C5
and before `BeginTx`.** So all three run after C4 and after C5. The invariant this pins: the
existing arm `consume_binding_test.go:138` (`Format = "webauthn"`, `ApprovedBy "ssh:x"`, no
`Passkey`) is still refused by **C4** with `ErrInconsistentApproval`, **unedited**. Moving C13
before C4 turns that arm red (CS7's ordering mutant).
- **C13** `Approver != nil && (Approver.Format == "webauthn") != (Approver.Passkey != nil)` →
  `ErrIncompleteConsume`.
- **C14** `Passkey != nil` with a credential ID outside **[r1]** 1–1023 bytes, or a `SignCount` outside
  0..2³²−1 → `ErrIncompleteConsume`.
- **C15** `Passkey != nil && c.version < 6` → `Older{Need: 6, Feature: "passkeys"}`.
- Inside the transaction, after C8's `useApprover`, when `Passkey != nil`:
  `UPDATE webauthn_credentials SET sign_count = ? WHERE name = ? AND credential_id = ? AND
  (sign_count < ? OR (sign_count = 0 AND ? = 0))`, with the new count in the first, fourth and
  fifth places. If that does not change exactly one row, `classifyPasskeyUse` reads the row in the
  transaction:
  - **C16** absent → `ErrApproverUnknown`;
  - **C17** another credential ID → `ErrPasskeyChanged`;
  - **C18** the stored counter is not below the new one (and they are not both zero) →
    `ErrSignCountRegression`.
  
  Any of these rolls back the request update and the approver's `last_used_at`. The request stays
  pending and no session row is written.

**The one named edit.** The existing CS3 arm "other-format" (`consume_binding_test.go:289`) sets
`Format = "webauthn"` on the ssh approver `x` with no `Passkey`, and expects C11's
`ErrApproverKeyChanged`. Under C13 it would get `ErrIncompleteConsume`. The arm gains
`rc.Approver.Passkey = &PasskeyUse{CredentialID: 32 bytes, SignCount: 1}`, so it passes C13–C15
and still reaches C11. That is the arm's original intent: the stored format differs.

**The store's law over `approved_by`**, extended from A1b §2.5: every `webauthn:<name>` is written
in the same transaction as a one-row update of an active `approvers` row with that name, format
and COSE key, **and** a one-row forward move of that name's `sign_count` for that credential ID.

### 2.11 Authority (`host/authority/mint.go`)

`requestConsume` keeps A1–A5. After A3 it adds one branch:
- **A6** when `rm.Signed.PasskeyUse()` is `ok`, it sets `consume.Approver.Passkey =
  &store.PasskeyUse{CredentialID, SignCount}`.

There is deliberately no separate "webauthn without `PasskeyUse`" refusal here. `PasskeyUse()`
is `ok` exactly when the approval is verified and its format is `webauthn` (B7, B8), so such a
branch would be unreachable from outside `host/approval` and its mutant unkillable. If A6 were
skipped, the store's C13 would refuse the consume (`ErrIncompleteConsume`), and AU3 pins that.
Nothing else in `authority` changes. RC5 needs no new `RequestMint` field.

### 2.12 `approver enroll --format webauthn` (sprint 3, `cmd/ailang-worldd/approver_passkey.go`)

```
approver enroll --name <name> --format webauthn --origin <origin> --registration <file> [--db <path>] [--credentials-db <path>]
```

`runApproverEnroll` dispatches on `--format webauthn` to `runApproverEnrollPasskey`, before
N4 and `NewApprover`. The steps:

1. **Flags.**
   - **WN1** `--pubkey` with `webauthn`, or `--origin`/`--registration` with `ssh` → help,
     `exitUsage`.
   - **WN2** `--origin` or `--registration` missing → `exitUsage`.
   - N2 and N3 as today.
2. **WN3** `ParseCeremonyOrigin` fails → `exitUsage`, naming the accepted forms and that
   `127.0.0.1` cannot be one.
3. **WN4** the `--registration` path already exists (`Lstat`) → `exitUsage` ("move it away; the
   browser saves the new registration there").
4. **Pre-checks** on a read-only handle under a 3 s root, closed before any terminal: N6 below
   **v6** → the restart text; N7 the name is taken.
5. **WN5** `env.random` fails, or is nil → `exitFatal`. It fills the 32-byte challenge.
6. **The link.** It prints to stdout:
   - `<origin>/surface/passkey#<b64url(JSON)>`, where the JSON is
     `{"v":1,"op":"create","rp":…,"origin":…,"name":…,"challenge":…}`;
   - the instruction: open the link in a browser on this origin, create the passkey, save the file
     to `<file>`, then type the code the page shows.
7. **Fence, part 1: the code.** It goes through `confirmApproverAct`'s terminal handling, refactored
   so one terminal serves two bounded reads (§2.12a):
   - N9a/N9 as today;
   - **[r1]** prompt `Enrolment code shown by the browser (XXXXX-XXXXX-XXXXX-XXXXX-XXXXXX): `.
     The operator types or pastes it (§2.8);
   - **WN6** no line within `passkeyCeremonyBudget = 300 s` → `exitUsage`, `no code within 5m0s;
     nothing was enrolled`;
   - **WN7** EOF, a read error or a partial line → `exitUsage`.
8. **WN8** reading `--registration` through `readRegularFileCapped(…, 16 KiB)` fails →
   `exitUsage`.
9. **WN9** `NewPasskeyApprover` fails → `exitUsage`, naming the typed error.
10. **WN10 [r1]** the typed code, normalised (uppercase, without `-`, spaces or tabs), is not
    exactly 26 characters of `A–Z2–7`, or is not `subtle.ConstantTimeCompare`-equal to the
    26-character `EnrolmentCode(...)` of the file → `exitUsage`: "the code you typed is not the code
    of <file>; nothing was enrolled". There is one attempt per process: the process exits on a
    mismatch.
11. **WN11** a second read-only pre-check under a new 3 s root: the COSE key or the credential ID is
    already enrolled → `exitUsage`, naming the holder. This is the third `Background` root.
12. **Fence, part 2: the y/N**, on the same terminal, within `approverConfirmBudget` (120 s): prompt
    `RenderPasskeyEnrolment(a) + "Enrol approver <name>? [y/N] "`. N10a and N10 are as today.
13. **The write** under the writer lock, taken only now: `EnrolPasskey`. N11 and N13 are as today;
    N12 gains `ErrPasskeyCredentialEnrolled`.
14. **N14** prints `enrolled approver <name>: passkey <ES256|EdDSA> <rpID> cred <16 hex>; approvals
    by it record approved_by=webauthn:<name>`.

**Why 300 s for the code.** Between the link and the code the operator opens a browser, passes
Face ID or Touch ID (or uses a phone across devices), and saves a file. A1b's 120 s is for reading
a fingerprint. 300 s is a constant, not a flag (A1b Q7's reasoning). No store handle is open while
either prompt waits.

#### 2.12a One terminal, two bounded reads

`confirmApproverAct` (`approver.go:33-55`) is split, behaviour-preserving:
- `openApproverTerminal(env, stderr, prefix, nothing) (term, code)` holds N9a and N9.
- `askApproverLine(term, prompt, budget, after) (line string, ok, complete bool)` holds the
  error-aware bounded read. `complete` means `err == nil` and the line ends in `\n`.
- `confirmApproverAct` becomes `openApproverTerminal` + `askApproverLine` + the `isYes` test, with
  the **same messages and exit codes**. CA4, CA9 and CA10 stay green **unedited**, and they are the
  refactor's oracle.

#### 2.12b `approver list` and `revoke` show passkeys

`approverKeyLabel(row, passkey)` gains a branch:
- **L4** a `webauthn` row with its passkey prints `passkey <alg> <rpID> cred <16 hex>`;
- **L5** a `webauthn` row with no passkey row prints `invalid-passkey`.
`approver list` calls `Passkey(ctx, name)` for each `webauthn` row, on the same read-only handle
(at most 256 rows, under the existing 3 s root). Revoked passkeys are listed too, unlike
`ActivePasskeys`. On a v5 store it calls no passkey method and lists only `ssh` approvers, because
a webauthn row cannot exist below v6. The revoke law and the
K8 listing are unchanged: K8 already filters by `row.Format + ":" + name`, so it lists
`webauthn:<name>` approvals.

### 2.13 `approve --webauthn-link` and `--webauthn-assertion` (sprint 3)

```
approve <request-id> --webauthn-link [--db <path>] [--credentials-db <path>]
approve <request-id> --webauthn-assertion <file> [--db <path>] [--credentials-db <path>]
```

- **F1** is widened to five exclusive modes. F2 (`--list` with any mode) is unchanged.
- **`--webauthn-link`** runs after admission and F4, like `--write-message`.
  - **WF1** no active passkey → `exitUsage`, "no passkey approver is enrolled; see approver list".
  - Otherwise it prints, for each distinct origin among the active passkeys:
    `<origin>/surface/passkey#<b64url(JSON)>`, where the JSON is
    `{"v":1,"op":"get","rp":…,"origin":…,"request":id,"expires":expires_at,
    "hash16":first16hex(MessageHash),"challenge":b64url(MessageHash),"allow":[credential IDs at
    that origin]}`.
  - It also prints the next step: `ailang-worldd approve <id> --webauthn-assertion <file> --db …
    --credentials-db …`.
  - It opens no terminal and writes no row, and the request stays pending.
- **`--webauthn-assertion <file>`**, after admission and F4, before any network call:
  - **WF2** `readRegularFileCapped(…, 16 KiB)` fails → `exitUsage`.
  - **WF3** the envelope does not parse (E2–E4) → `exitUsage`.
  - On the same read-only handle as A1b's lookup: **WF4** below v6 → `exitUsage`, naming the
    restart; **WF5** no approver holds that credential ID → `exitUsage`, "the assertion's
    credential `cred <16 hex>` is not an enrolled passkey; see approver list".
  - **WF6** `ApproveByAssertion` fails → `exitUsage`, naming the reason.
  - Then steps 2–6 are unchanged. `Signed` is the passkey approval, so `provisionEpisode` skips the
    terminal and prints `approved by webauthn:<name>` (A1b §2.8 V2, unchanged).
- **No new context root.** Both paths use `runApprove`'s existing 3 s read root.

### 2.14 The ceremony page: what lands in A2, and what is A3

**In A2.** Two GET routes, registered next to the existing `/surface/*` assets:
`GET /surface/passkey` (HTML) and `GET /surface/passkey.js`.
- They serve embedded bytes (`host/surfacepage/passkey/{index.html,passkey.js}`, `go:embed`) with
  `setSurfaceHeaders` and `surfaceCSP` (V12).
- They read no store and no request.
- The page makes **no network request at all**: no `fetch`, no `XMLHttpRequest`, no form `action`.

The script:
1. **Reads `location.hash`** and decodes its strict JSON, `v == 1`. If `origin !==
   location.origin` or `rp !== location.hostname`, it shows which origin to open and stops.
2. **Shows what will happen**, with `textContent` only and never `innerHTML`.
   - create: "Create a passkey for approver <name> at <rp>".
   - get: "Approve request <id>, expires <local time>, message hash <hash16>". The operator
     compares this with what the CLI printed.
3. **[r2]** Each ceremony creates an AbortController and starts a fixed 240000 ms application timer. Pass controller.signal alongside publicKey to navigator.credentials.create/get. On expiry, abort the operation, mark the attempt expired, disable saving, and display an explicit timeout error. Ignore any late result; clear the timer in finally. Disable the ceremony button while an attempt is pending. (Quorum r2, `gpt6-1-sol`, applied verbatim: WebAuthn's `timeout` is a client-adjustable hint, not an enforced bound.) The calls below therefore carry `signal: controller.signal` beside `publicKey`.
   **On a button click**, calls the API. A user gesture is required by some browsers.
   - create: `navigator.credentials.create({publicKey: {challenge, rp: {id: rp, name: "AILANG
     World"}, user: {id: utf8(name), name, displayName: name}, pubKeyCredParams: [{type:
     "public-key", alg: -7}, {type: "public-key", alg: -8}], authenticatorSelection:
     {userVerification: "required", residentKey: "preferred"}, attestation: "none", timeout:
     240000}})`.
   - get: `navigator.credentials.get({publicKey: {challenge, rpId: rp, allowCredentials: allow,
     userVerification: "required", timeout: 240000}})`.
4. **Builds the §2.7 JSON itself.** It does not depend on `toJSON()`, whose support outside
   Chrome is UNVERIFIED.
5. **Shows the result** in a read-only `<textarea>`, and for create the grouped enrolment code
   (§2.8).
6. **[r1] Saves it through one measured mechanism.** The page sets a "Save" link up before the
   operator clicks:
   - `a.href = URL.createObjectURL(new Blob([json + "\n"], {type: "application/json"}))`;
   - `a.download = "world-registration-<name>.json"` or `"world-assertion-<id>.json"`.
   The operator's own click triggers the download.
   - **Measured (V31) in Chrome 157 headless, under a page served with exactly `surfaceCSP`** (the
     string pinned at `host/daemon/surface_fonts_test.go:64`), saving through a real
     DevTools-dispatched mouse click:
     - `Browser.downloadWillBegin` fired with the `blob:http://localhost:…` URL and the suggested
       name;
     - `downloadProgress` reached `completed`;
     - the file landed in the download directory, byte-equal to the textarea's JSON plus `\n`;
     - `securitypolicyviolation` and the console recorded **nothing**.
   - **The listener works.** A control run of the same page, which also loads one blocked image,
     recorded `img-src http://127.0.0.1:9/control.png` and Chrome's `default-src 'none'` console
     message, and still completed the download.
   - So `blob:` downloads are not blocked by this CSP in Chrome, and **`surfaceCSP` is not
     changed**.
   - The textarea stays only as a *display* of what was saved, so the operator can see it. It is
     not a fallback for a broken link.
   - **XB1 asserts the save path itself.** It sets `Browser.setDownloadBehavior` to a temp
     directory, clicks Save with a dispatched mouse event, and requires `completed`, a file
     byte-equal to the textarea, and zero CSP violations. Go then verifies **the saved file**, not
     the textarea.
   - Safari and Firefox remain UNVERIFIED (R3), and the attended run (§10) checks the save there.

**Why a page in A2 at all.** The parent puts "a real browser check on the tailnet with a passkey"
in A2. A passkey ceremony needs a page at the rpID's origin. File and IP origins cannot be one
(§2.2), and the daemon is what serves the surface's origin. The page is the minimal piece: static,
same-origin, no state, and no write.

**Read-route posture (parent §3).** The page adds nothing readable: it holds no data and reads
none. Approving stays bound to a verified signature, never a cookie or a session. In A2 the signed
artifact travels as a file to the CLI. In A3 it travels in a POST body, and the same
`ApproveByAssertion` judges it.

**A3 owns:**
- the *needs you* card and request card, fetched from the daemon rather than carried in a fragment;
- the Approve POST;
- moving this page into the surface app;
- the notification and Touch ID convenience;
- the graphical decision list.
A3 reuses A2's law, store and page script unchanged.

### 2.15 The test authenticator and the recorded fixtures

**`approvaltest.Passkey`** (`host/approval/approvaltest/passkey.go`, imported only by `_test.go`
files, as TS3 guards):
- `NewPasskey(alg int64, seed byte)`. EdDSA uses `ed25519.NewKeyFromSeed`; ES256 uses
  `ecdsa.ParseRawPrivateKey(P256, seed-derived 32 bytes)` (V15).
- `Register(origin, rpID, name string, challenge [32]byte, opts …) []byte` and
  `Assert(origin, rpID, name string, challenge [32]byte, counter uint32, opts …) []byte` produce
  §2.7 files in **Chrome's measured layout**: `clientDataJSON` member order, CTAP2 key order,
  flags `0x45` and `0x05`, a 32-byte credential ID, `fmt: none` with an empty `attStmt`.
- Knobs (options) produce each refusal arm: flags, type, origin, `crossOrigin`, `topOrigin`,
  duplicate members, extra COSE labels, non-shortest heads, tags, trailing bytes, a wrong `alg`,
  an off-curve point, a counter, the BE/BS bits, a `userHandle`, and an extensions map.
- EdDSA output is byte-deterministic. ES256 signatures are random (Go 1.26), so ES256 tests check
  verification, not bytes. Where bytes matter, a test uses `cryptotest.SetGlobalRandom` and does not
  run in parallel.

**Recorded fixtures** (`host/approval/testdata/webauthn/`):
- `chrome-es256.json` and `chrome-eddsa.json`. Each records the Chrome version, origin, rpID, the
  challenge hex, a `registration`, an `assertion` (counter 2) and an `assertion2` (counter 3).
- `SOURCE.txt` and `generate.sh`.
- They are produced by **`generate.sh`, which runs the M7 Node/DevTools driver against Chrome's
  virtual authenticator**: a test key, never a real authenticator (D-WORLD-109).
- The challenge for both the create and the get is the **existing fixture message hash**
  `29181d22…fef5`. So the recorded assertions approve the existing fixture request (`fixtureDigest`,
  `fixtureChallenge`, `fixtureExpiry`, `request_test.go`). This is exactly what the designer's probe
  did (V20–V22).
- The files are JSON text with **no NUL**, so `scripts/tracked_binary_gate.sh` passes.
  `.gitattributes` gains `host/approval/testdata/webauthn/* -text`. Their SHA-256s are pinned by
  `TestPasskeyFixturesAreTheGeneratedBytes` (XW0).
- **The designer's probe outputs are not committed.** The sprint regenerates them with the
  committed generator.

**How the tests show the fixtures are what a real browser emits:**
- **XW1, layout.** `approvaltest.Passkey`'s output has the same skeleton as Chrome's: CBOR major
  types, key order and lengths for the attestation object and COSE key; `authData` lengths (37, and
  55 + L + key); `clientDataJSON` member order. A test authenticator that drifts from real browsers
  goes red.
- **XW2, the real bytes through the production law.** Both Chrome fixtures enrol through
  `NewPasskeyApprover` (with the recorded origin and challenge). Their assertions approve the
  fixture request through `ApproveByAssertion`, and `assertion` after `assertion2` is refused by
  W15. The byte layout each test relies on is the spec's: L3 §6.1 (authData), §5.8.1
  (`CollectedClientData`) and §8.7 (`none`).
- **XB1, the page in a live Chrome** (M7, gated). The production page and the production law agree
  with a real browser end to end, including the enrolment code.

## 3. Pure core, effects at the boundary (coding standards S1–S3)

| Piece | Where | Pure? | Tested by |
|---|---|---|---|
| CBOR subset, COSE, authData, clientData, origin, envelope | `host/approval/{cbor,cose,authdata,clientdata,origin,passkey_json}.go` | yes | CB1–CB4, CK1–CK4, AD1–AD3, CD1–CD2, OR1, EN1 |
| `NewPasskeyApprover`, `EnrolmentCode`, `RenderPasskeyEnrolment` | `host/approval/passkey.go` | yes | PE1–PE5, RC4 |
| `ApproveByAssertion`, `SignedApproval.PasskeyUse` | `host/approval/passkey_assert.go`, `signed.go` | yes | PA1–PA4, RC1 |
| test authenticator | `host/approval/approvaltest/passkey.go` (non-production) | yes | TP1, XW1, TS3 |
| schema 6, `EnrolPasskey`, lookups, the counter in the consume | `host/store/passkeys.go`, `mint_requests.go` | effect | PS1–PS6, CS7–CS9, RC6 |
| `approved_by = webauthn:` | `host/authority/mint.go` | no I/O beyond the store call | AU3 |
| `/dev/tty`, the clock, the timer, randomness | `sessionEnv` (`openTerminal`, `now`, `after`, `random`): no new seam | seam | CW1–CW9 |
| the page | `host/surfacepage/passkey/*`, served by `host/daemon` | static asset | PG1, XB1 |

- **Why not `.ail` (S1).** As A1b §3 says, the pinned v0.41.0 `std/crypto` has no ECDSA and no
  Ed25519 (A1b V22). A2 adds a byte-level CBOR parser and two signature verifiers, none of which
  has an invariant Z3 could state without the cryptography. **A2 writes no `.ail`**, so the AILANG
  fluency protocol (S5) has nothing to apply to, and `verify_ail.sh` is expected exactly green.
- **Why not a package (S3).** The consumers are `cmd/ailang-worldd`, `host/authority` and (A3) the
  daemon, all host Go. `host/approval` already holds the SSHSIG half of the same law.
- **S2 kept.** `host/approval` stays standard-library-only with no I/O. The new imports are
  `crypto/subtle` (not on the forbidden list) and packages it already imports
  (`imports_test.go:12-16`; V17).

## 4. Threats (A2)

| Threat | What A2 does | Pinned by |
|---|---|---|
| **An agent enrols a passkey** | Enrolment needs the typed code *and* a typed `y` at `/dev/tty` (WN6–WN10, N9–N10). Agents have no terminal (parent §1). `EnrolPasskey` has one production caller, after `isYes` | CW3, CW4 |
| **[r1] An agent reads the genuine registration file, then substitutes its own** (its own software authenticator, the printed challenge; it can compute the genuine code) | The code is **public** to this attacker, and it is not relied on as a secret. It is a 128-bit fingerprint of the exact bytes the CLI reads. A substitute must match the genuine code: a 128-bit second preimage, about 2¹²⁸ work. The attacker cannot choose the genuine bytes, so birthday collisions do not help. One attempt per process, a constant-time compare, then the typed `y` (§2.8, WN10) | CW5 (the attacker knows and types the genuine code; its file is refused; 0 rows) |
| **Replay of an assertion** | A new challenge per request (the message hash covers that row's challenge, digest and expiry). The request is consumed once (A1b C6). A replay is refused by admission (P1) | PA2, CW7 |
| **An assertion for request R1 used for R2** | W6 compares against R2's own message hash. C6 binds the row's challenge | PA1 (W6 arm), CS7 |
| **Counter regression, or a cloned authenticator** | W15 before any network call, and C18 inside the mint transaction. Refused, request pending, no auto-revoke (Q2). Both-zero counters skip the check (spec), so residual R2 | PA3, CS8 |
| **rpID or origin mismatch, or a phishing page on another tailnet host** | The browser scopes the passkey to the rpID. W7 and W10 require the stored origin and rpID exactly. rpID is `localhost` or `*.ts.net` (O2) and never an IP | PA1 (W7, W10 arms), OR1 |
| **UV absent** (a passkey used without biometric or PIN) | R14 and W12 refuse. The page requests `required` | PE1, PA1 |
| **A wrong algorithm** (RS256, ES384) | Q3 refuses at enrolment, so no such key is ever stored | CK1 |
| **Oversize or malformed CBOR** (deep nesting, indefinite lengths, tags, floats, duplicate keys, a huge length) | K1–K16 with fixed bounds. The file is at most 16 KiB before any parse | CB1–CB3 |
| **Two encodings of one key, to enrol it twice** | The stored key is `Canonical()`. `UNIQUE(public_key)`, plus `UNIQUE(credential_id)` | CK3, PS2 |
| **Credential ID collision** (one ID at two approvers) | `UNIQUE(credential_id)` and the `EnrolPasskey` pre-check. An assertion is judged against the key of the approver its ID names, so a collision cannot transfer a signature | PS2 |
| **A revoked or replaced approver racing the consume** | C8 (A1b) and C16/C17 inside the mint transaction. The `Clock` seam commits a revoke inside the window, as A1b's CS3 does | CS9 |
| **An SSH approver's name colliding with a passkey approver's** | `approvers.name` is the primary key across formats, so one name is one key and one format forever. A consume claiming `webauthn:x` for ssh approver `x` fails C11 (format differs) | CS9 (the edited "other-format" arm) |
| **A same-origin page swapping the challenge** (an injected script asks the operator's passkey to sign another pending request's hash) | `surfaceCSP` has `script-src 'self'`, and the page uses no inline script or `innerHTML`. The page shows the request id, expiry and `hash16`, which the CLI also prints. The assertion approves only the request whose hash it signed, and only when someone runs `approve <that id>`. **Residual R5:** the A2 display is unauthenticated (from the fragment); A3 shows the request from the daemon | PG1, XB1 |
| **The page leaks data** | It reads no store and makes no request (PG1's labelled controls). The fragment never reaches a server | PG1 |
| **A process that can rewrite the daemon binary, the page or `credentials.db`** | Out of scope, unchanged. It can already insert rows | — |

## 5. By-construction argument

1. **No passkey enrols without a typed code and a typed `y`.**
   - `EnrolPasskey` has one production caller, `runApproverEnrollPasskey`. It runs only after
     `openApproverTerminal` succeeded, a complete code line matched the file's `EnrolmentCode`
     (WN10), and a complete `y` line was read (N10).
   - **[r1] No substituted registration enrols, even when its author read the genuine one.** The
     code is not secret. It is a 128-bit fingerprint of the bytes the CLI read, and the operator
     copies it from his own browser. A substitute enrols only if its fingerprint equals the
     genuine one, which is a 128-bit second preimage on SHA-256 (§2.8). There is one comparison
     per process.
   - `EnrolApprover` still refuses every format but `ssh`.
2. **`approved_by = webauthn:<name>` means `<name>`'s enrolled passkey, with UV, signed this exact
   request.**
   - The value comes only from a `SignedApproval` returned by `ApproveByAssertion`. That requires
     W6 (the message hash of the stored digest, the row's challenge and the row's expiry), W12 (UV)
     and W14 (the signature under the stored COSE key).
   - Authority re-checks `Binds` (A4).
   - The store writes the value only alongside a one-row update of that active approver's row (C8)
     and a one-row forward move of its counter for that credential ID (C16–C18).
3. **A counter never moves backwards.** The only `UPDATE` of `sign_count` is guarded by
   `sign_count < new OR both zero`, inside the consume transaction.
4. **The terminal is skipped only with a verified approval.** This is unchanged from A1b §5.2:
   `Signed` is set at one site, from either law's return.
5. **Nothing minted differs from the request.** Unchanged from a0a1 §5.2: the passkey path runs
   `approve`'s steps 2–6 exactly as the SSH path does.
6. **World never holds a private key.** The browser's authenticator signs, and World receives only
   public data.

## 6. Test plan

**Rules the planner must carry (iterations 373–378):**
- **One mutant per new production branch.** The list: K1–K16, Q1–Q8, D1–D8, J1–J7, O1–O3, E2–E4,
  R1–R15, P3w, W1–W15, B7–B8, S1–S9 (the `EnrolPasskey` validations), the three `EnrolPasskey`
  pre-checks, C13–C18, A6 (authority), WN1–WN11, WF1–WF6, the widened F1, L4–L5, and the M0
  carry-forwards. The planner copies this list into each sprint plan. The evaluator applies each
  mutant and records which assertion fires (S6). A branch the planner adds gets its own row.
- **Reflective completeness** for every struct a security check claims to cover. **RC1**
  `SignedApproval` (+`credentialID`, `signCount`). **RC3** `RequestConsume`, `ApproverUse` and the
  new element type `PasskeyUse`. **RC4** `Approver` (+`Passkey`) and the element type `Passkey`:
  each field is shown by `RenderPasskeyEnrolment` or excluded with a reason. **RC6** `PasskeyRow`
  against `PRAGMA table_info(webauthn_credentials)`, with a distinct non-zero value per field read
  back.
- **Exit codes** are asserted by constant. A refusal arm checks the typed error with
  `errors.Is`/`errors.As` **and** that nothing was written.
- **No `t.Parallel`** in A2 tests. `cryptotest.SetGlobalRandom` forbids it, and
  `TestNoParallelWriteForkPackagesOutsideVerifygate` refuses it in write+fork packages
  (`host/daemon`, `cmd/ailang-worldd`).

### Sprint 1: the verification law

**M0, the carry-forwards (test-only, about 90 lines)**

| Id | Test | Asserts | Mutant killed |
|---|---|---|---|
| K8a | `TestApproverRevokeListingReopenFailureIsFatal` (`approver_revoke_test.go`) | the `approverRevoke` seam runs the real revoke, then renames the credentials file away → `exitFatal`; stderr has `revoked, but cannot read the credentials store`; stdout has `revoked approver <name>` | `approver.go:449` `return exitFatal` → `return exitOK` (**survives today**, V3) |
| K8b | `TestApproverRevokeListingQueryFailureIsFatal` | a new test seam `approverDecided = (*store.Credentials).DecidedMintRequests` (the `approverRevoke` pattern, `approver.go:376`) returns an error → `exitFatal`; stderr has `revoked, but cannot list what it approved` | `:455` `exitFatal` → `exitOK` (**survives today**, V3) |
| SG1 | `TestSignRefusesAFailedSignerEvenWithOutput` (`host/sshsign/sshsign_test.go`) | a temp script prints `partial-signature` and exits 3 → `Sign` returns a non-nil error wrapping `*exec.ExitError` and **nil** bytes | `sshsign.go:53` G1 dropped (the script's output would come back as a signature; **survives `./host/sshsign/` today**, V4) |

**M1, CBOR and COSE (`cbor_test.go`, `cose_test.go`)**

| Id | Test | Asserts | Mutants killed |
|---|---|---|---|
| CB1 | `TestCBORRefusals` | one arm per K1–K16, each a hand-built byte string: `0x18` alone (K1), `0x1c` (K2), `0x5f…` (K3), `0x18 0x05` (K4), `0x58 0x20` + 3 bytes (K5), a 2,049-byte string (K6), a 17-entry map (K7), 5 nested arrays (K8), `0xc0 0x00` (K9), `0xf9 0x00 0x00` and `0xf7` (K10), text `0x61 0xff` (K11), a map keyed by bytes (K12), `{1:1,1:2}` (K13), `0x1b ff…ff` and `0x3b 80 00…` (K14), 65 items over nested maps (K15), `0x00 0x00` via `decodeCBORExact` (K16); each `errors.Is` its error | each K deleted; each bound `>` turned `>=` (the boundary arm below) |
| CB2 | `TestCBORBoundsAreInclusive` | exactly 2,048 bytes, 16 entries, depth 4 and 64 items are **accepted** | off-by-one in each bound |
| CB3 | `TestCBORAcceptsChromesBytes` | both fixtures' `attestationObject` and COSE keys decode with 0 trailing bytes | an over-strict subset (for example, refusing `false` in `crossOrigin`-free documents) |
| CB4 | `FuzzDecodeCBOR` (seeded with CB1/CB3) | never panics; a successful decode re-encodes no longer than its input; ran for 30 s by the sprint with `-fuzz` and logged | — (robustness; not a mutant oracle) |
| CK1 | `TestCOSEKeyRefusals` | Q1–Q8: an array; no `kty`; alg −257 and −35 (Q3 → `ErrCOSEAlg`); −7 with kty 1, −7 with crv 2; −8 with crv 1; a 31-byte x; an off-curve y (Q7); label 2 (`kid`) present; EdDSA with label −3 | each Q deleted; Q3 accepting −257 |
| CK2 | `TestCOSEKeyVerifies` | Chrome's ES256 and EdDSA keys verify Chrome's assertion signatures over `authData ‖ SHA-256(cdj)`; one flipped byte → false; a 73-byte ES256 signature and a 63-byte Ed25519 signature are refused before verification | hashing with SHA-512; verifying over `authData` alone |
| CK3 | `TestCOSECanonicalIsChromesBytes` | `Canonical()` of each parsed Chrome key is byte-equal to the key bytes inside Chrome's `authData`; a key decoded from a permuted-order map re-encodes to the same canonical bytes | storing the received bytes |
| CK4 | `TestCOSEKeyHasOnlyDeclaredFields` (reflective) | `COSEKey`'s fields are exactly `Alg`, `X`, `Y`, each set by a test arm | a field added that `Canonical()` ignores |

**M2, authData, clientData, origin and the envelope**

| Id | Test | Asserts | Mutants killed |
|---|---|---|---|
| AD1 | `TestAuthDataRefusals` | D1–D8 from `approvaltest.Passkey` knobs: 36 bytes; BS without BE; registration without AT; assertion with AT; a truncated credential length; **[r1]** credential IDs of 0 and 1,024 bytes refused, and of 1 and 1,023 bytes accepted; a bad COSE key; ED with an integer; 1 trailing byte | each D deleted; **[r1]** a 16-byte floor reintroduced (the 1-byte arm reds) |
| AD2 | `TestAuthDataReadsChrome` | Chrome registration: flags UP, UV, AT and not BE/BS, counter 1, AAGUID `0102…08`, a 32-byte credential ID; assertions: 37 bytes, UP and UV, counters 2 and 3 | the flag bit positions swapped (UV read from bit 1) |
| AD3 | `TestAuthDataIgnoresExtensionsContents` | ED with `{"credProtect": 2}` → parses, `Key` intact | D7 over-strict (refusing any extension) |
| CD1 | `TestClientDataRefusals` | J1–J7; plus acceptance of an unknown member `other_keys_can_be_added_here` and of reordered members | duplicate-key detection deleted (a duplicated `challenge` arm goes green) |
| CD2 | `TestClientDataReadsChrome` | Chrome's two `clientDataJSON`s → `type`, `challenge` (base64url of the fixture hash), `origin`, `crossOrigin=false` | — |
| OR1 | `TestCeremonyOrigin` | accepted: `http://localhost`, `http://localhost:7644`, `https://localhost:8443`, `https://voights-mac-studio.tail97eda0.ts.net`; refused: **[r2]** `https://x.ts.net:443`, `https://localhost:443`, `http://localhost:80` (each with guidance naming the default-port rule), `http://127.0.0.1:7644`, `https://[::1]`, `https://LOCALHOST`, `https://x.ts.net/`, `https://x.ts.net:0`, `https://x.ts.net:065535`, `https://evil.example`, `https://ts.net`, `https://a..ts.net`, `http://x.ts.net` (http only for localhost), `https://u@x.ts.net` | each O deleted; `http` allowed for `.ts.net` |
| EN1 | `TestPasskeyEnvelope` | E2–E4: bad JSON, no `response`, `type` `"x"`, `id ≠ rawId`, a padded field, a non-canonical last character; acceptance of Chrome's `toJSON()` output with its extra members | E3 or E4 deleted |

**M3, the laws, the code, `SignedApproval`, the test authenticator, fixtures and layout**

| Id | Test | Asserts | Mutants killed |
|---|---|---|---|
| PE1 | `TestNewPasskeyApproverRefusals` | one arm per R1–R15 (R16 via CK1), each with the typed error and the zero `Approver` | each R deleted; R14 (UV) deleted specifically, with a UV-clear arm |
| PE2 | `TestNewPasskeyApproverAcceptsChrome` | both Chrome registrations → `Format "webauthn"`, `Passkey` fields equal to the fixture's (rpID `localhost`, the recorded origin, the credential ID, alg, AAGUID, counter 1, BE false), `EnrolledAt == now` | a field not copied |
| PE3 | `TestRenderPasskeyEnrolmentGolden` | golden text for the Chrome ES256 registration: rp, origin, `ES256`, `cred` prefix, AAGUID, `device-bound`, `approved_by=webauthn:<name>`, the §2.8 paragraph verbatim; with BE and BS set (test authenticator) → `synced passkey` | a dropped field or paragraph |
| PE4 **[r1]** | `TestEnrolmentCodeGolden` | `EnrolmentCode` of the committed Chrome ES256 registration equals a pinned **26-character** value, computed once from the committed fixture and recorded in the plan. The designer's probe registration gave `YDOR4X36PQNFRTXBHWAV4W2XG4` (V35), which shows the format; the pin is the fixture's own. Further arms: length exactly 26, alphabet `A–Z2–7`; `GroupEnrolmentCode` gives 5-5-5-5-6 groups; a 1-byte change to either input changes it | hashing one input only; truncating to fewer than 16 bytes (the length arm); a different grouping |
| PE5 | `TestNewApproverPointsWebAuthnAtRegistration` (edits `approver_test.go:150-151`) | `NewApprover(…, "webauthn", …)` → `ErrUseRegistration` | E2 deleted (webauthn → "must be ssh") |
| PA1 | `TestApproveByAssertionOrder` | one arm per P1–P5, P3w, W1–W15, each from a test-authenticator knob on an otherwise valid assertion; P4 (revoked) with a good assertion says revoked; W2 before W14 (another credential's good assertion says credential mismatch) | each check deleted; reordered pairs |
| PA2 | `TestApproveByAssertionAcceptsChrome` | Chrome `assertion` (counter 2) approves the fixture request against the passkey from PE2 with stored counter 1 → `ApprovedBy() == "webauthn:<name>"`, `PasskeyUse()` → `(credential ID, 2, true)` | the counter not carried |
| PA3 | `TestSignCountLaw` | stored 1 / new 2 → ok; 2/2, 3/2, 5/0 → `ErrSignCountRegression`; 0/0 → ok; 0/7 → ok; Chrome `assertion` after `assertion2` → refused | `<=`/`<` swapped; the both-zero clause dropped (0/0 refused); the `||` turned `&&` (5/0 passes) |
| PA4 | `TestSignedApprovalPasskeyUse` | zero/nil → `ok false` (B7); a verified **ssh** approval → `ok false` (B8); a verified passkey approval → the credential ID copy and count | B7/B8 deleted; an uncloned slice |
| RC1 | `TestSignedApprovalFieldsAreAllBound` (edited) | classifies `credentialID` and `signCount` as `PasskeyUse` and perturbs each | a new field unclassified |
| RC4 | `TestEnrolmentPromptShowsEveryPersistedField` (edited) | `Approver.Passkey` and the element type `Passkey`: every field rendered or excluded with a reason (`BackupState`: "shown as synced/device-bound"; `Key`: "shown as the algorithm") | a new `Passkey` field the prompt never shows |
| TP1 | `TestTestPasskeyVerifies` | for both algorithms, `Register` → `NewPasskeyApprover` accepts; `Assert` → `ApproveByAssertion` accepts | — (it must be correct; XW1/XW2 tie it to Chrome) |
| XW0 | `TestPasskeyFixturesAreTheGeneratedBytes` | the SHA-256 of each fixture file equals its pin; `SOURCE.txt` names the Chrome version and `virtual authenticator` | a hand-edited fixture |
| XW1 | `TestTestPasskeyMatchesChromeLayout` | the skeleton comparison of §2.15 for both algorithms | a test authenticator emitting another key order, flag value or member order |
| XW2 | (PE2 + PA2 + PA3's last arm) | the real bytes through the production law | — |
| — | `TestApprovalImportsOnlyTheStandardLibrary`, `TestApprovalTestIsTestOnly` (existing) | still green with the new files | an `os` import; production importing `approvaltest` |

### Sprint 2: the store and the binding

**M4, schema 6 and the passkey store (`host/store/passkeys_test.go`, edited credentials tests)**

| Id | Test | Asserts | Mutants killed |
|---|---|---|---|
| PS1 | `TestCredentialsSchemaMatchesCanonicalDDL` (edited) | five tables; `webauthn_credentials` DDL equals the schema file and `passkeysSQL`; `user_version` 6 | drift between the two DDL copies |
| PS2 | `TestEnrolPasskeyRefusals` | v5 → `Older{Need 6}` with the text `passkeys need schema 6`; each S1–S9 arm → `ErrBadApprover` (**[r1]** credential ID 0 and 1,024 bytes refused, 1 byte accepted); a taken name, an enrolled COSE key, an enrolled credential ID (each also after the holder is **revoked**) → the typed error naming the holder; 0 rows in **both** tables each time | each validation deleted; a pre-check deleted (the `UNIQUE` backstop error is not `errors.Is` the typed error); one `INSERT` committed without the other |
| PS3 | `TestPasskeyLookups` | `PasskeyByCredentialID`, `Passkey` and `ActivePasskeys` (revoked excluded, ordered, capped at 64); misses; v5 → `Older` | the lookup by prefix; revoked rows included |
| PS4 | `TestCredentialsUpgradeToV6` | stores at v1, v3, v4 and v5 → v6 with every row unchanged, idempotently; a CLI `OpenCredentials` leaves a v5 store alone | the `version < 6` step missing |
| PS5 | `TestEnrolApproverStillSSHOnly` | `EnrolApprover` with format `webauthn` → `ErrBadApprover` (existing branch, `approvers.go:89`) | — (guards that passkeys enter only through `EnrolPasskey`) |
| RC6 | `TestPasskeyRowMatchesTheTable` | reflective `PasskeyRow` ↔ `PRAGMA table_info`; distinct values round-trip | a column dropped from a `SELECT` |
| edited | `TestCredentialsMethodsRequireADeadline`, `TestCredentialsUpgradeToV5`, version pins | +4 methods; "current" pins move to `credentialsSchemaVersion` (6) | — |

**M5, the consume and authority**

| Id | Test | Asserts | Mutants killed |
|---|---|---|---|
| CS7 | `TestPasskeyConsumeCompleteness` | C13 (webauthn with no `Passkey`; ssh with a `Passkey`), C14 (**[r1]** an empty and a 1,024-byte credential ID, a counter of 2³²), C15 (v5 store) → typed errors, Clock not read, request pending | C13, C14, C15 deleted; **[r1]** C13–C15 moved before C4 (the unedited `consume_binding_test.go:138` arm then gets `ErrIncompleteConsume` instead of `ErrInconsistentApproval` and reds) |
| CS8 | `TestPasskeyConsumeAdvancesTheCounter` | stored 1, new 2 → ok: `sign_count = 2`, `last_used_at = at`, `approved_by = webauthn:x`, 1 session; again with new 2 on a fresh request → `ErrSignCountRegression` (C18), request pending, counter still 2, `last_used_at` unchanged; 0/0 → ok, counter 0; 5/0 → refused | the counter `UPDATE` deleted (a regression mints); its `WHERE` widened to `<=`; the both-zero clause dropped |
| CS9 | `TestPasskeyConsumeChecksTheRowInTheTransaction` | the `Clock` seam (A1b CS3's, V36/V37 there) on a second connection: revokes `x` → `ErrApproverRevoked`; rewrites its credential ID → `ErrPasskeyChanged` (C17); deletes its passkey row → `ErrApproverUnknown` (C16); raises its counter to 9 → `ErrSignCountRegression`; each rolls back `last_used_at`, the request and the session | `RowsAffected` unchecked; the classify arms swapped |
| edited | CS3 "other-format" (`consume_binding_test.go:289`) | gains `Passkey: &PasskeyUse{…}` and still expects `ErrApproverKeyChanged` (§2.10) | — |
| unedited | CS2 "approver-format-not-ssh" (`:138`) | still `ErrInconsistentApproval` from C4 (C13 placed after C4) | C13 placed before C4 (this arm then reds) |
| RC3 | `TestRequestConsumeCompleteness` (edited) | `ApproverUse.Passkey` (optional, with a meaning: required exactly for webauthn) and the element type `PasskeyUse` (`CredentialID`, `SignCount`: required) | an unclassified field |
| AU3 | `TestMintForRequestDerivesPasskeyUse` (`host/authority`) | a verified passkey `Signed` (from `approvaltest.Passkey` through `ApproveByAssertion`) → the row's `approved_by = webauthn:<name>`, the counter advanced, 1 session; a zero `SignedApproval` → `ErrNotVerified` (A4), nothing written | A6 deleted (the consume then reaches C13 → `ErrIncompleteConsume`, and the arm expects success) |

### Sprint 3: the operator paths and the page

**M6, the CLI (`cmd/ailang-worldd/approver_passkey_test.go`, `approve_passkey_test.go`)**

The tests use the A1b rigs (`spyEnv`, `approveRig`, temp stores). Registrations and assertions come
from `approvaltest.Passkey` with the challenge the run printed: the test parses the printed link's
fragment. Keys are never real.

| Id | Test | Asserts | Mutants killed |
|---|---|---|---|
| CW1 | `TestPasskeyEnrollFlagRules` | WN1 (each pairing), WN2, WN3 (`http://127.0.0.1:7644` → names `localhost`), WN4 (the file exists) → `exitUsage`, no terminal, 0 rows | each WN deleted |
| CW2 | `TestPasskeyEnrollEnrolsWhatTheBrowserMade` | the spy terminal answers the right code then `y` → one approvers row (`webauthn`, canonical COSE key) and one passkey row; the link printed to stdout has the challenge the file carries; the y/N text contains `RenderPasskeyEnrolment`; N14's line | storing the received COSE bytes; the challenge not from `env.random` |
| CW3 | `TestPasskeyEnrollWithoutATerminalEnrolsNothing` | `noTerminalEnv` → `sessionTTYFix`, 0 rows | fence bypassed |
| CW4 | `TestPasskeyEnrollDeclinedOrUnanswered` | code timeout (WN6: the recorded budget is exactly 300 s, the message `no code within 5m0s`); code EOF/partial (WN7); right code then `n`, EOF, `y`+EOF, or y timeout (120 s) → `exitUsage`, terminal closed, 0 rows; **while each prompt waits the store opens writable** (no lock held) | a 120 s code budget; an unbounded read (the 5 s guard fires); the writer opened before the y/N |
| CW5 **[r1]** | `TestPasskeyEnrollRefusesASubstituteEvenWithTheGenuineCode` | **read-then-substitute.** The test plays the attacker. Authenticator A makes the genuine registration `G` for the printed challenge, and the test computes `code(G)`, as an attacker who read `G` can. Authenticator B makes `X` for the same challenge and writes it to `--registration`. The spy terminal types `code(G)`. Result: WN10, `exitUsage`, the y/N never shown, **0 rows**. Further arms: `G`'s file with `code(G)` with its **last** character changed → refused; with its first group changed → refused; a 25- and a 27-character code → refused; a code containing `0` → refused; `G`'s code typed in lowercase, with `-`, spaces and a tab → accepted; and the same pasted as one 26-character string → accepted | WN10 deleted (X enrols); a prefix compare (`code[:8]`, killed by the last-character arm); the length check dropped (the 25-character prefix arm); the compare made case-sensitive |
| CW6 | `TestPasskeyEnrollRaceAndRefusals` | WN5 (nil `random` → `exitFatal`); WN8 (a FIFO → refused by `fstat`); WN9 (a UV-clear file); WN11 (the credential enrolled by another handle while the code prompt waited) → `exitUsage`; N12 `ErrPasskeyCredentialEnrolled` from the seam → `exitUsage` | WN11 deleted (the arm reaches the y/N spy) |
| CW7 | `TestApproveByPasskey` | after enrolment: `--webauthn-link` prints one link per origin with the request's `hash16` and the allow list, opens no terminal, writes nothing (WF1 on an empty store); `--webauthn-assertion` with the test authenticator's assertion over that fragment's challenge → `exitOK` with **the spy terminal failing the test if opened**, `approved_by = webauthn:<name>`, counter advanced, and A1b's `approveEffectProblems` passes; the same file again → refused as already used | the terminal opened; WF1 deleted |
| CW8 | `TestApprovePasskeyRefusalsBeforeAnyEffect` | WF2 (a FIFO, 16 KiB+1), WF3, WF4 (v5 store, names the restart), WF5 (an unknown credential), WF6 (counter regression, revoked, another request's hash) → `exitUsage`, the `stubDaemon` sees 0 requests, request pending; F1 with any two of the five modes | each WF deleted; the lookup after the network call |
| CW9 | `TestApproverListAndRevokeShowPasskeys` | `list` prints `passkey ES256 localhost cred …` (L4) and `invalid-passkey` for a hand-made row (L5); `revoke` shows the label and then lists `webauthn:<name>` approvals (K8) | L4/L5 deleted |
| edited | `TestEveryVerbHasHelp` | usage lines `approver enroll … --format webauthn`, `approve <request-id> --webauthn-link`, `--webauthn-assertion <file>` | — |
| edited | `TestProductionContextRoots` | `cmd/ailang-worldd/approver_passkey.go\|runApproverEnrollPasskey\|Background: 3` | — |
| unedited | CA1–CA10, CR1–CR13 and every A1b CLI test | green, with `confirmApproverAct`'s refactor (§2.12a) | — |

**M7, the page, the live browser and QUICKSTART §16**

| Id | Test | Asserts | Mutants killed |
|---|---|---|---|
| PG1 | `TestPasskeyPageRoutes` (`host/daemon`) | `GET /surface/passkey` and `/surface/passkey.js` → 200, the embedded bytes, `surfaceCSP` and the surface headers; `POST` → 405; served with the world store closed (no store read). **Labelled instrument-health controls (not load-bearing):** the script contains no `fetch(`, `XMLHttpRequest`, `innerHTML` or `eval(` | the route serving another asset; a missing CSP |
| edited | `TestSurfaceAppRoutes` | the asset set gains the two names | — |
| XB1 | `TestPasskeyPageInChrome` (`host/daemon`, gated: `WORLD_CHROME` + `node` on `PATH`; `t.Skip("SKIP COUNTED: …")` unless `WORLD_REQUIRE_PASSKEY_BROWSER=1`, then `t.Fatal`) | an `httptest` daemon on `localhost`; `testdata/passkey_cdp.mjs` (the designer's probe, generalised) starts Chrome with a temp profile and `--remote-debugging-port=0`, adds a **virtual authenticator** (`ctap2`, `internal`, UV true), opens the create link the production code built, and clicks Create. **[r1]** It then sets `Browser.setDownloadBehavior` to a temp directory and clicks **Save** with a dispatched mouse event. It returns the downloaded file, the textarea text, the displayed grouped code, the download events and every `securitypolicyviolation`. Go requires: download state `completed`; the file byte-equal to the textarea plus `\n`; **zero** CSP violations (the page is served with the production `surfaceCSP`); `NewPasskeyApprover` accepting **the saved file**; and `EnrolmentCode`, grouped, equal to the page's code. Then the get link for the fixture request: `ApproveByAssertion` accepts the page's assertion. Chrome is killed and the profile removed | the page computing the code over one input, or truncating it; **[r1]** the Save link removed or blocked (no `completed` download); sending `userVerification: "preferred"` with a UV-false authenticator variant (a second subtest, which must be **refused** by R14); the page using its own challenge |
| XB2 **[r2]** | `TestPasskeyPageAbortsAPendingCeremony` (`host/daemon`, same gating as XB1) | Quorum r2 (`gpt6-1-sol`), applied verbatim: a browser test with an injected credentials operation that remains pending until aborted; it asserts the exact deadline, signal abortion, explicit error, and absence of any downloadable result after expiry. The page's 240000 ms timer is the production constant; the test drives it with CDP `Emulation.setVirtualTimePolicy` (or an injected clock seam the plan names) rather than waiting four minutes, and the injected `navigator.credentials.create`/`get` returns a promise that settles only when its `signal` aborts. Go requires: the abort fired at exactly 240000 ms of page time; `signal.aborted` true; the explicit timeout error shown; the Save link absent or disabled and no `downloadWillBegin`; a late resolution after expiry ignored; the ceremony button disabled while pending | no `signal` passed (the operation never aborts); the timer not started, or a different constant; saving left enabled after expiry; a late result accepted; the button left enabled while pending |
| — | `generate.sh` | re-records the §2.15 fixtures with the same driver | — |
| — | QUICKSTART §16 | "Approve with a passkey" *(attended — pending first verbatim run)*. **[r1]** It shows the 26-character code in its five groups, says the code may be typed or pasted, and says it is a fingerprint, not a secret: `serve` on the §0 scratch store, `http://localhost:7644`, enrol, request, link, assertion, decided list, revoke, a refused approval | — (S7) |

**Must stay green, unedited:** every a0a1 and A1b test, except the named edits:
- `approver_test.go:150-151` (PE5);
- `approver_enroll_test.go:230-232`, whose `--format webauthn --pubkey` arm now expects WN1's
  usage refusal, and whose `row 201 A2` substring changes to `--registration`;
- `consume_binding_test.go:289`;
- the schema-version pins and the reflective tests' new classifications.
Also `TestDaemonDependencyAllowlist` and the quickstart tests.

## 7. Repo-wide inventories (iteration 378's retro)

Each sprint runs these by name, beyond its own `-run`. "Change" says what A2 does to the inventory.

| Inventory test | Where | Change |
|---|---|---|
| `TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential` (AC10 exec-site census) | `host/broker/registry_publish_test.go:1143` | **No new production exec site.** XB1's Chrome and Node launches are in `_test.go` files, which the census skips (`:1094`) |
| `TestDaemonDependencyAllowlist` | `host/daemon/daemon_test.go:899` | no new module |
| `TestWorldBoundaryDependencyAllowlist`, `TestBareNetHTTPExemptionIsPerGroup`, `TestEveryCommandDirectoryIsAProtectedGroup` | `host/boundary/allowlist_world_test.go:758`, `:895`, `:1121` | no new transport import in `cmd/ailang-worldd`, `host/store` or `host/replay`; no new `cmd/` directory |
| `TestProductionContextRoots`, `TestProductionGoSurface` | `host/store/context_roots_test.go:116`, `:221` | **+1 root entry** (`approver_passkey.go\|runApproverEnrollPasskey\|Background: 3`); no Go file outside `host/` and `cmd/` |
| `TestNoNewDeadlineFreeStoreReads`, `TestReadGettersHonorContext` | `host/store/context_read_test.go:377`, `:52` | the four new store methods take a deadline context |
| `TestCredentialsMethodsRequireADeadline` | `host/store/credentials_test.go:260` | +4 methods |
| `TestApprovalImportsOnlyTheStandardLibrary`, `TestApprovalTestIsTestOnly` | `host/approval/imports_test.go:27`, `approvaltest/signer_test.go:86` | `crypto/subtle` added (allowed); `approvaltest` stays test-only |
| `TestFrozenV1RouteTableMatchesMux` | `host/daemon/route_table_test.go:15` | **no `/v1` route**; stays 13 |
| `TestSurfaceAppRoutes` | `host/daemon/surface_app_test.go:19` | +2 assets (edited) |
| `TestDaemonConstructorCallerInventory` | `host/daemon/startup_callers_test.go:136` | unchanged |
| `TestNoParallelWriteForkPackagesOutsideVerifygate` | `host/verifygate/forklocked_write_test.go:717` | no `t.Parallel` added |
| Reflective RC1, RC3, RC4, RC5, RC2/RC6 | `signed_test.go:177`, `consume_binding_test.go:197`, `approver_test.go:342`, `host/authority/mint_derive_test.go:241`, `approvers_test.go:388` | new fields classified; RC6 new |
| tracked-binary gate | `scripts/verify_go.sh` → `scripts/tracked_binary_gate.sh` | fixtures are JSON text, with no NUL |

**A2 adds:** no production exec site, **no new Go package** (`host/surfacepage` gains an embed
directory; `approvaltest` gains a file), no new module, one new standard-library import
(`crypto/subtle`, in `host/approval`), one context-root entry, two `/surface` GET routes, and no
`/v1` route.

## 8. Conflict surface

| Sprint | File | Change |
|---|---|---|
| 1 | `host/approval/{cbor,cose,authdata,clientdata,origin,passkey_json,passkey,passkey_assert}.go` (+ tests) | **new**, pure |
| 1 | `host/approval/approver.go` | `Approver.Passkey`; E2's error renamed; `checkApprover` shared with `signed.go` |
| 1 | `host/approval/signed.go` | `credentialID`, `signCount`, `PasskeyUse` |
| 1 | `host/approval/approvaltest/passkey.go` (+ test) | **new**, test support |
| 1 | `host/approval/testdata/webauthn/*`, `.gitattributes` | **new** fixtures; +1 line |
| 1 | `cmd/ailang-worldd/approver.go` (+ `approver_revoke_test.go`) | the `approverDecided` seam (M0 K8b only) |
| 1 | `host/sshsign/sshsign_test.go` | +SG1 |
| 2 | `host/store/credentials_schema.sql`, `credentials.go`, `passkeys.go` (new), `mint_requests.go` (`ApproverUse.Passkey`, defined at `:294`; C13–C18), `approvers.go` (`useApprover` unchanged) + tests | schema 6 |
| 2 | `host/authority/mint.go` (+ test) | A6 |
| 3 | `cmd/ailang-worldd/approver_passkey.go` (new), `approver.go` (the §2.12a split, labels), `approve.go`, `approve_sig.go` (or a new `approve_passkey.go`), `main.go`, `help_test.go` | the CLI |
| 3 | `host/surfacepage/passkey/{index.html,passkey.js}`, `host/surfacepage/passkey.go` (embed), `host/daemon/daemon.go` (+2 `HandleFunc`), `surface_page.go` | the page |
| 3 | `host/daemon/testdata/passkey_cdp.mjs`, `host/daemon/passkey_page_*_test.go`, `host/approval/testdata/webauthn/generate.sh` | the browser driver |
| 3 | `docs/QUICKSTART.md` | §16 |

**Not touched**: `world/*.ail`, `packages/`, `scripts/verify_ail.sh`, `tools/launchd/*`, any
`/v1` route or `design_docs/sketches/worlddapi.ail`, `cmd/world-publish/*` (R4), `host/sshsign/sshsign.go`
(test only), the typed-`y` path, `session *`, presets and `--resume`. No open PRs at measurement
time (V18).

## 9. Milestones, sprints and gates

Every milestone ends with `go build ./... && go vet ./<pkg>/` plus its own tests. `go vet` is the
compile fence for test files (MEMORY: go-build-is-not-a-compile-fence). Every sprint ends with:
- `./scripts/verify_ail.sh` (no `.ail` change, so exactly green);
- `AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh` (it includes the tracked-binary
  gate);
- `go build ./... && go test ./...`;
- **the §7 inventory tests by name**: `go test ./host/broker/ -run
  '^TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential$' -count=1`; `go test
  ./host/daemon/ -run '^(TestDaemonDependencyAllowlist|TestFrozenV1RouteTableMatchesMux|TestSurfaceAppRoutes|TestDaemonConstructorCallerInventory)$'
  -count=1`; `go test ./host/boundary/ -run
  '^(TestWorldBoundaryDependencyAllowlist|TestEveryCommandDirectoryIsAProtectedGroup)$' -count=1`;
  `go test ./host/store/ -run
  '^(TestProductionContextRoots|TestProductionGoSurface|TestNoNewDeadlineFreeStoreReads|TestCredentialsMethodsRequireADeadline)$'
  -count=1`; `go test ./host/verifygate/ -run '^TestNoParallelWriteForkPackagesOutsideVerifygate$'
  -count=1`.

**Sprint 1 = M0–M3, "the verification law", about 1,800 lines.**
- **M0: carry-forwards.** About 10 production lines (the seam) and 80 test lines. K8a, K8b, SG1.
  `go test ./cmd/ailang-worldd/ -run '^TestApproverRevokeListing' -count=1`;
  `go test ./host/sshsign/ -count=1`.
- **M1: CBOR and COSE.** About 270 production and 430 test lines. CB1–CB4, CK1–CK4.
  `go test ./host/approval/ -run 'CBOR|COSE' -count=1`.
- **M2: authData, clientData, origin and the envelope.** About 230 + 380 lines. AD1–AD3, CD1–CD2,
  OR1, EN1. `go test ./host/approval/ -run 'AuthData|ClientData|CeremonyOrigin|PasskeyEnvelope' -count=1`.
- **M3: the laws, the test authenticator and the fixtures.** About 260 + 220 production and
  test-support lines, about 430 test lines, and the fixtures (about 60 lines, plus the generator,
  which ships in M7). PE1–PE5, PA1–PA4, RC1, RC4, TP1, XW0–XW2.
  `go test ./host/approval/... -run 'Passkey|Assertion|SignCount|Enrolment|SignedApproval|NewApprover|Imports|ApprovalTest' -count=1`.
- **Sprint 1 fallback cut:** the fixtures need the M7 driver to regenerate. **If the plan's estimate
  exceeds 2,000 lines, or the driver is not yet committed**, sprint 1 commits the fixtures with
  `generate.sh` marked *regeneration lands in sprint 3*. XW0 still pins them, and `SOURCE.txt`
  names the iteration-379 recording. If it still exceeds 2,000, M3's PA* half (the assertion law)
  moves to the head of sprint 2.

**Sprint 2 = M4–M5, "the store and the binding", about 1,600 lines.**
- **M4: schema 6 and the passkey store.** About 230 + 520 lines. PS1–PS5, RC6.
  `go test ./host/store/ -run 'Credentials|Passkey|EnrolApprover|OlderCredentials' -count=1`.
- **M5: the consume and authority.** About 110 + 560 lines (the CS9 rig is the expensive part).
  CS7–CS9, the CS2/CS3 arms, RC3, AU3.
  `go test ./host/store/ ./host/authority/ -run 'Consume|Passkey|ApprovedBy|MintForRequest|RequestMint' -count=1`.
- **Sprint 2 fallback cut:** AU3 and authority's A6 move to the head of sprint 3. The store
  then accepts a passkey consume that no caller builds yet, which is safe.

**Sprint 3 = M6–M7, "the operator paths and the page", about 1,800 lines.**
- **M6: the CLI.** About 400 + 820 lines. CW1–CW9 and the edited help, roots and enrol arms.
  `go test ./cmd/ailang-worldd/ -run '^(TestPasskey|TestApproveByPasskey|TestApprovePasskey|TestApproverListAndRevokeShowPasskeys|TestApprover|TestApprove|TestEveryVerbHasHelp)' -count=1`.
- **M7: the page, the live browser and QUICKSTART.** About 230 page lines (HTML, JS and Go), about
  300 test and driver lines, and 50 documentation lines. PG1, XB1 and the generator.
  `WORLD_CHROME="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" WORLD_REQUIRE_PASSKEY_BROWSER=1 go test ./host/daemon/ -run '^TestPasskeyPageInChrome$' -count=1 -v`
  must print `--- PASS` for both subtests, **on the controller's machine** (CI does not run it:
  Chrome and Node on the runner are UNVERIFIED, and the existing `WORLD_CHROME` tests are not in CI
  either).
- **Sprint 3 fallback cut:** M7 becomes its own fourth sprint. A2's attended exit run needs M7, so
  it waits for it.

## 10. The attended exit run (Mark; NOT run by the loop, D-WORLD-109)

On a scratch store, with his own browser and passkey provider:
1. **Serve.** `ailang-worldd serve` on the scratch store. Then either use
   `http://localhost:7644/surface/passkey` on the daemon's machine, or put the tailnet name in front
   of it with `tailscale serve` (HTTPS, tailnet only, never Funnel). Opening the daemon to the
   tailnet exposes its other routes too; that is row 199's posture, not A2's.
2. **Enrol.** `approver enroll --name <his-name> --format webauthn --origin
   https://<host>.ts.net --registration ~/Downloads/world-registration-<his-name>.json`. Open the
   printed link in Safari or Chrome on the Mac. Create the passkey with Touch ID or iCloud Keychain,
   or with the phone by the browser's cross-device QR code, and save the file. Type or paste the page's
   26-character code, **[r1]** in five groups, and check the shown rp, algorithm and `synced passkey`/`device-bound`, then type `y`.
3. **Approve.** `episode start … --request`, then `approve <id> --webauthn-link`. Open it, pass
   Touch ID or Face ID, and save. `approve <id> --webauthn-assertion <file>` mints with no World
   prompt.
4. **Check.** `approve --list --decided` shows `webauthn:<his-name>`. Running the same assertion
   file again is refused as already used.
5. **Revoke.** `approver revoke <his-name>` (typed). A new request with the same passkey is refused
   as revoked.
6. **Record**, in the iteration log, **public values only**: the browser and provider, `fmt`, the
   algorithm, the flags (BE/BS), whether the counter stayed 0, and whether the CBOR was canonical.
   This answers R3 and Q3. The credential ID is a stable identifier of his passkey, so it is not
   committed unless Mark says so (Q4).

QUICKSTART §16 is executed verbatim in that run, and its *attended — pending* tag is removed. The
loop's judgement uses only M0–M7's tests.

## 11. Residuals and risks

- **R1 `fmt: none` means no hardware proof.** World cannot tell a hardware passkey from a software
  one, and UV is self-reported by whichever authenticator holds the key. The binding is the typed
  enrolment code and `y` (§2.8): a key enrolled from the operator's own ceremony. That is stronger
  than A1b's R6, since a passkey requires UV per use, but it rests on the enrolment.
- **R2 Counterless passkeys.** Synced passkeys are reported to send counter 0 (UNVERIFIED), so W15
  and C18 detect no clone for them. The spec's own rule.
- **R3 Browser variance.** Only Chrome 157 was measured. **[r1]** In Chrome the `blob:` save under
  the exact `surfaceCSP` is measured working (V31). Safari's and Firefox's CBOR canonicity,
  `toJSON()` support and `blob:` save are UNVERIFIED. K4's strictness and the save link are the
  likely failure points. The attended run (§10.6) is the check, and the fix
  would be a narrow change to K4 or the page.
- **R4 RS256 refused.** A provider that offers only RS256 (reported for some Windows Hello TPMs,
  UNVERIFIED) cannot enrol. `crypto/rsa` is standard library and would be a small, separate change
  (Q1).
- **R5 The A2 page's display is unauthenticated.** It shows the request id, expiry and `hash16`
  from the fragment, and the operator compares them with the CLI. A3 replaces this with the request
  fetched from the daemon.
- **R6 Schema 6 costs one `serve` restart** (`credentials.go:437`), as A1b's R2.
- **R7 Size.** A1b ran at 1.7× its estimate (V30). The fallback cuts in §9 are pre-named.
- **R8 The browser test is local-only.** XB1 runs on the controller's machine with `WORLD_CHROME`.
  CI covers everything else through the recorded fixtures.

## 12. Open questions

**None blocks the loop.** Each item below has a default the loop builds on, and Mark can revisit
any of them at the exit run.
- **Q1** Accept RS256 (−257) too? **Default: no in A2** (parent §2.1's library list); revisit if
  Mark's provider needs it.
- **Q2** A counter regression refuses the approval and suggests a revoke; it never auto-revokes.
  **Default: yes** (A1b Q4's reasoning: one typed act never becomes silent ones).
- **Q3** rpID is exactly `localhost` or a `*.ts.net` host, and never a custom domain. **Default:
  yes**, from parent §2.2's wording.
- **Q4** Commit Mark's exit-run registration and assertion (public data) as an attended fixture?
  **Default: no.** Record the public properties only.

## 13. What this is NOT doing

- No surface card, no Approve POST, no notification or Touch ID, no graphical decision list (A3);
  no `Human.Approve` (A4).
- No attestation verification, no RS256, ES384 or other algorithm, no extensions processing, no
  `topOrigin`/cross-origin use, no custom rpID domains.
- No `/v1` route, no `.ail`, no package, no new Go module, no production exec site.
- No change to the SSH path's behaviour, the typed `y`, `session *`, presets, `--resume` or
  publishing (R4 of the parent).
- No real key, no real store, no mint in one and no publish, by the loop or its tests
  (D-WORLD-109). Every passkey is a test authenticator's or Chrome's **virtual** authenticator's,
  in a temp profile.

## 14. Verification Log

Commands ran in the worktree at `98374909` unless stated otherwise. Empty results are paired with a
positive control in the same command.

| # | Claim | Command | Observed |
|---|---|---|---|
| V1 | Base and instruments | `git log --oneline -1`; `shasum -a 256 host/approval/{approver,signed,sshsig,request}.go host/store/{approvers,mint_requests,credentials}.go host/authority/mint.go cmd/ailang-worldd/{approve,approve_sig,approver}.go host/sshsign/sshsign.go \| awk '{print substr($1,1,16), $2}'`; `go version`; `$HOME/.pinned-ailang/ailang --version`; `"/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" --version`; `node --version` | `98374909 Merge pull request #545 …`; the twelve prefixes in the header; `go1.26.6 darwin/arm64`; `AILANG v0.41.0`; `Google Chrome 157.0.8103.50`; `v26.13.0` |
| V2 | A1b's shape and seams | `wc -l design_docs/planned/w-world-approvals-by-key*.md`; `sed -n 829,833p design_docs/planned/w-world-approvals-by-key-a1b.md` | a1b 1,304 lines, a0a1 1,012, parent 93; the seam note "`sign_count` by `ALTER TABLE ADD COLUMN` at schema 6, `approved_by = webauthn:<name>`, and the browser half of enrolment (a one-time code typed at this same y/N)" |
| V3 | **K8-ERRCODE: both K8 exits are unpinned (premise TRUE)** | `cp cmd/ailang-worldd/approver.go /tmp/w379-designer/approver.go.bak`; `go test ./cmd/ailang-worldd/ -run 'TestApprover' -count=1`; `sed -i '' '449s/return exitFatal/return exitOK/'`; the same `go test`; restore; then the same for line `455`; restore; `git status --short` | before: `ok … 0.817s`; `:449 return exitOK` → `ok … 0.837s` (**survives**); `:455 return exitOK` → `ok … 0.751s` (**survives**); `git status` empty after each restore. Lines `:446-450` are the read-only reopen; `:451-456` the `DecidedMintRequests` call |
| V4 | **sshsign G1: no unit test runs a failing signer (premise TRUE)** | `grep -n "func Test" host/sshsign/sshsign_test.go`; backup; `sed` the G1 line to `if err := cmd.Run(); false && err != nil {`; `go test ./host/sshsign/ -count=1`; `go test ./cmd/ailang-worldd/ -run 'SignWith' -count=1 -v`; restore; `git status --short` | tests: `:18 TestCappedBufferKeepsItsBound`, `:32 TestSignHandsTheSignerAScrubbedEnvironment`; mutant: `ok …/host/sshsign 0.465s` (**survives**); the cmd real-tool tests catch it (`--- FAIL: TestSignWithSSHKeygenRefusesBadOutput`, `--- FAIL: TestSignWithSSHKeygenEndToEnd`), which need `ssh-keygen`; restored, status empty |
| V5 | `approve`'s modes, flow and file reader | `cat -n cmd/ailang-worldd/approve.go`; `cat -n cmd/ailang-worldd/approve_sig.go`; `sed -n 81,114p cmd/ailang-worldd/approver.go` | `:117 func runApprove`; `:147-156` F1 `modes > 1` over `--signature`, `--sign-with`, `--write-message`; `:225-228` F4; `:235-250` signature judged before the network; `:298-299 req.consume = &authority.RequestMint{… Signed: signed, …}`; `approve_sig.go:30 approveSignature` → `readRegularFileCapped(path, approveMaxSignatureBytes)`; `approver.go:84 func readRegularFileCapped` (`O_RDONLY\|syscall.O_NONBLOCK`, `f.Stat()`, `IsRegular`, `readCapped` via `io.LimitReader(r, max+1)`) |
| V6 | `MessageHash` exists, is pinned, and has no production caller (with control) | `grep -rn "MessageHash" --include='*.go' host cmd` | `request.go:224-225` (definition); `request_test.go:35 const fixtureMessageHashHex = "29181d224d982f659d1bd4703238d3f8489e26968ffe2a22bd786fa3bf15fef5"`, `:295-300` the golden; no other match |
| V7 | The enrolment law refuses webauthn and is pinned | `cat -n host/approval/approver.go`; `grep -rn "ErrFormatNotYet\|row 201 A2\|\"webauthn\"" --include='*_test.go' host cmd` | `:18 FormatWebAuthn = "webauthn"`; `:39 ErrFormatNotYet`; `:53-55` E2; `:28-35 type Approver struct`; tests `approver_test.go:150-151`, `approver_enroll_test.go:230-232` (`want: []string{"row 201 A2"}`), `consume_binding_test.go:138`, `:289` |
| V8 | The signature law | `cat -n host/approval/signed.go`; `grep -n "func TestSignedApprovalFieldsAreAllBound" host/approval/signed_test.go` | `:16-24 SignedApproval{ok, digest, challenge, expiry, name, format, key}`; `:55 ApproveBySignature`, `:62 if in.Approver.Format != FormatSSH { // P3`; `:80 Binds`; `:99 ApprovedBy`; `:107 Use`; RC1 at `signed_test.go:177` |
| V9 | Schema 5, table sets, the upgrade, the bound-mint refusal, the approvers DDL | `grep -n "credentialsSchemaVersion\|func credentialsTables\|version < \|case version" host/store/credentials.go`; `sed -n '/^func credentialsTables/,/^}/p'`; `sed -n '/^func upgradeCredentials/,/^}/p'`; `sed -n 435,440p`; `grep -n approvers -A12 host/store/credentials_schema.sql`; `cat -n host/store/approvers.go` | `:42 const credentialsSchemaVersion = 5`; `case version >= 5: return []string{approversTable, "episode_starts", mintRequestsTable, credentialsTable}`; `if version < 5 { statements = append(statements, approversSQL) }`; `:437 if (row.Bound() \|\| start != nil \|\| resume != nil) && c.version < credentialsSchemaVersion {`; schema `:63-70` (`public_key BLOB NOT NULL UNIQUE`); `approvers.go:89 if row.Format != "ssh"`; `:245 useApprover`; `:262 classifyApproverUse` (C9–C12) |
| V10 | The credentials DB does not enable foreign keys (with control) | `grep -rn "foreign_keys" --include='*.go' host/store` | only `store.go:397` (the world store) and `durable_driver_fault_test.go:100`; none in `credentials*.go` |
| V11 | The consume's order, and the two existing webauthn arms | `cat -n host/store/mint_requests.go \| sed -n 279,360p`; `sed -n 118,160p host/store/consume_binding_test.go`; `sed -n 280,300p` same | `:306` C1, `:309` C2, `:312` C3, `:315` C4, `:321` mint-request version, `:324` C5, `:327 BeginTx`, `:332 at := rc.Clock()`, `:336` C6, `:348-352` C8, `:353` session insert; arm `:138 {"approver-format-not-ssh", … rc.Approver.Format = "webauthn" }` expects `ErrInconsistentApproval`; arm `:288-290 {"other-format", … rc.Approver.Format, rc.ApprovedBy = "webauthn", "webauthn:x" }, ErrApproverKeyChanged}` |
| V12 | Daemon bind, surface CSP and headers, no Permissions-Policy, no Host check | `grep -n "BindHost\|DefaultBindHost =" host/daemon/daemon.go`; `grep -n "surfaceCSP\s*=" host/daemon/*.go`; `grep -n "func setWorkbenchHeaders" -A7 host/daemon/workbench.go`; `grep -rn "Permissions-Policy" --include='*.go' host \| grep -v _test`; `grep -rn "r.Host\|allowedHost" --include='*.go' host/daemon \| grep -v _test`; `sed -n 1091,1124p host/daemon/daemon.go` | `:163 DefaultBindHost = "127.0.0.1"`, `:713 if !isLoopbackHost(cfg.BindHost)`; `surface_page.go:19 const surfaceCSP = "default-src 'none'; style-src 'unsafe-inline'; script-src 'self'; connect-src 'self'; font-src 'self'; base-uri 'none'; form-action 'self'; frame-ancestors 'none'"`; headers `Content-Type`, `Cache-Control: no-store`, CSP, `nosniff`, `Referrer-Policy: no-referrer`; Permissions-Policy: none; Host: only `session_self.go:69`; routes include `GET /surface`, `/surface/app.js`, … and the 13 `/v1` |
| V13 | Route pins | `sed -n 1,60p host/daemon/route_table_test.go`; `sed -n 15,20p host/daemon/surface_app_test.go`; `sed -n 266,280p host/daemon/handlers_test.go` | `:51-53 if len(mux) != 13 \|\| len(sketch) != 13 …` over `/v1/` only; "every other name under /surface/ is a 404"; `TestGETRoutesRejectOtherMethods` names only `/v1/health`, `/v1/head` |
| V14 | Live-Chrome pattern and the census's test-file exclusion | `grep -rn "remote-debugging\|devtools\|--headless" --include='*.go' host cmd scripts`; `sed -n 20,26p host/daemon/surface_live_chrome_test.go`; `grep -n "func enumerateSubprocessSites" -A30 host/broker/registry_publish_test.go \| grep _test.go`; `grep -n "chrome\|Chrome" .github/workflows/ci.yml` | `workbench_live_chrome_test.go:68 args := []string{"--headless=new", …}`, `:73 exec.CommandContext(ctx, chrome, args...)`; **no** `remote-debugging` anywhere (control: the `--headless` hit); `os.Getenv("WORLD_CHROME")` … `t.Skip("WORLD_CHROME unset; …")`; `:1094 … strings.HasSuffix(path, "_test.go") {` (skipped); `ci.yml`: no Chrome |
| V15 | Go 1.26 crypto APIs | `go doc crypto/ecdsa SignASN1`; `go doc crypto/ecdsa GenerateKey`; `go doc testing/cryptotest SetGlobalRandom`; `go doc crypto/ecdsa ParseRawPrivateKey`; `go doc crypto/ecdsa ParseUncompressedPublicKey`; `go doc crypto/ecdsa VerifyASN1` | "Since Go 1.26, a secure source of random bytes is always used, and the Reader is ignored unless GODEBUG=cryptocustomrand=1 … use testing/cryptotest.SetGlobalRandom"; `SetGlobalRandom` "cannot be used in parallel tests"; `ParseRawPrivateKey` "fixed-length big-endian integer … SEC 1 … 2.3.6"; `ParseUncompressedPublicKey` "returns an error if the point is not in uncompressed form, is not on the curve, or is the point at infinity"; `VerifyASN1(pub, hash, sig []byte) bool` |
| V16 | Nothing of A2 exists (with controls) | `grep -rln "webauthn_credentials\|sign_count\|/surface/passkey\|clientDataJSON\|ParseCOSE\|decodeCBOR" --include='*.go' host cmd \| grep -v _test`; control `grep -rn "FormatWebAuthn = \|func MessageHash" --include='*.go' host` | none; control `approver.go:18`, `request.go:225` |
| V17 | `host/approval`'s import guard | `sed -n 1,78p host/approval/imports_test.go` | forbidden: `os`, `os/exec`, `os/signal`, `net`, `net/http`, `syscall`, `io/ioutil`, `crypto/rand`, `math/rand`, `math/rand/v2`, `unsafe`, and any path with a dot in its first segment (exact-path match, so `crypto/subtle` and `net/url` are not forbidden) |
| V18 | No open PRs | `gh pr list --state open --json number,headRefName,title --limit 20` | `[]` |
| V19 | The inventory tests named in §7 exist at those lines | `grep -n "^func Test" host/broker/registry_publish_test.go host/boundary/allowlist_world_test.go host/daemon/startup_callers_test.go host/store/context_read_test.go host/store/context_roots_test.go`; `grep -rn "func TestDaemonDependencyAllowlist" host`; `grep -n "func TestNoParallelWriteForkPackagesOutsideVerifygate" host/verifygate/forklocked_write_test.go`; `grep -rn "func TestApprovalTestIsTestOnly\|…\|func TestCredentialsSchemaMatchesCanonicalDDL" --include='*_test.go' host cmd`; `sed -n 17,25p host/store/context_roots_test.go` | `registry_publish_test.go:1143`; `allowlist_world_test.go:758, :895, :1121`; `startup_callers_test.go:136`; `context_read_test.go:52, :377`; `context_roots_test.go:116, :221`; `daemon_test.go:899`; `forklocked_write_test.go:717`; `signer_test.go:86`, `imports_test.go:27`, `credentials_test.go:64, :260`, `consume_binding_test.go:197`, `approver_test.go:342`, `mint_derive_test.go:241`, `approvers_test.go:388`, `help_test.go:20`; roots `:20-24` (`approve.go\|runApprove: 1`, `approve_sig.go\|signForApproval: 1`, `approver.go\|runApproverList: 1`, `\|runApproverEnroll: 2`, `\|runApproverRevoke: 3`) |
| V20 | **A real Chrome emits WebAuthn create/get with a virtual authenticator, over the fixture message hash** (throwaway) | `/tmp/w379-designer/browser/probe.mjs` (sha `01a17c0d7bab9333`): Node starts a `127.0.0.1` HTTP server and Chrome `--headless=new --remote-debugging-port=0 --user-data-dir=<mktemp>`; DevTools `WebAuthn.enable`, `WebAuthn.addVirtualAuthenticator {protocol: ctap2, ctap2Version: ctap2_1, transport: internal, hasResidentKey, hasUserVerification, isUserVerified, automaticPresenceSimulation}`; navigates to `http://localhost:<port>/`; `Runtime.evaluate` runs `create` (`rp.id: "localhost"`, `pubKeyCredParams` `[-7]` then `[-8]`, `userVerification: "required"`, `attestation: "none"`) and `get` twice with challenge = `29181d22…fef5`; bounded by a 90 s `date +%s` poll; `pgrep -fl w379-designer` after | `chrome: Chrome/157.0.8103.50`, `enable: ok`, `authenticatorId …`; outputs `es256.json` (sha `c0370982dbca416a`), `eddsa.json` (sha `1544721fc45dd65a`); the probe's own profile removal raced Chrome once (`ENOTEMPTY`) and was removed by hand; `no probe process` |
| V21 | The recorded byte layout (throwaway Go decoder over V20's output) | `/tmp/w379-designer/verify/main.go` (sha `d7267aceda7b9fd7`), `GOFLAGS=-mod=mod GOPROXY=off go run . <file>` for both | ES256: `attestationObject 194 bytes, trailing=0`; `top-level key order: [fmt attStmt authData]  fmt: none  attStmt: map[]`; `authData 164 bytes; rpIdHash=49960de5…9763` = `sha256(localhost)`; `flags=0x45; signCount=1`; `aaguid=01020304050607080102030405060708 credIdLen=32 credId==rawId:true`; `COSE key 77 bytes`, `order: [1 3 -1 -2 -3]  kty: 2  alg: -7  crv: 1`; reg `clientDataJSON: {"type":"webauthn.create","challenge":"KRgdIk2YL2WdG9RwMjjT-EieJpaP_ioivXhvo78V_vU","origin":"http://localhost:56547","crossOrigin":false}`; assertion `authData 37 bytes flags=0x05 signCount=2; second get signCount=3`; get `{"type":"webauthn.get",…}`. EdDSA: `attestationObject 159 bytes`, `authData 129`, `COSE key 42 bytes`, `order: [1 3 -1 -2]  kty: 1  alg: -8  crv: 6`, same flags and counters. `toJSON()` present with `publicKey`, `publicKeyAlgorithm`, `transports`, `authenticatorAttachment`, `clientExtensionResults` |
| V22 | Chrome's assertion signatures verify with the standard library over `authData ‖ SHA-256(clientDataJSON)` | same program | `ES256: sig 72 bytes (DER, first byte 0x30); VerifyASN1=true; tampered=false`; `EdDSA: sig 64 bytes; Verify=true; tampered=false`; `challenge b64url of fixture message hash: KRgdIk2YL2WdG9RwMjjT-EieJpaP_ioivXhvo78V_vU` |
| V23 | COSE heads, the on-curve check, strict base64url (throwaway) | `/tmp/w379-designer/p2/main.go`, `go run .` | `es256 cose head: a5010203262001215820 len 77`; `eddsa cose head: a4010103272006215820 len 42`; `on-curve parse err: <nil>`; one flipped byte of y → `P256 point not on curve`; `RawURLEncoding.Strict()` on a padded string → `illegal base64 data at input byte 43`; a non-canonical last character → `illegal base64 data at input byte 42` under `Strict()`, but `<nil>` without it |
| V24 | Probes stayed out of the repo | `ls /tmp/w379-designer`; `git status --short` (worktree) | `approver.go.bak browser p2 sshsign.go.bak verify`; worktree: only this doc |
| V25 | WebAuthn section and step numbers (WebFetch of the W3C TRs) | `https://www.w3.org/TR/webauthn-3/` read in 100 k-character windows; `https://www.w3.org/TR/webauthn-2/` window 300 k–400 k | L3: §5.8.1 "Client Data Used in WebAuthn Signatures (dictionary CollectedClientData)", §5.8.1.2 "Limited Verification Algorithm", §6.1 flags "UP: Bit 0 … UV: Bit 2 … BE: Bit 3 … BS: Bit 4 … AT: Bit 6 … ED: Bit 7", §6.1.1 (counter mismatch only when "either is non-zero"; "less than or equal … cloned"), §7.1 steps 7 (type), 8 (challenge), 9 (origin), 10–11 (crossOrigin, topOrigin), 14 (rpIdHash), 15 (UP), 16 (UV), 17 ("If the BE bit … is not set, verify that the BS bit is not set"), 20 (alg ∈ pubKeyCredParams), 21 (fmt), 25 ("credentialId is ≤ 1023 bytes"), 26 ("not yet registered"), §8.7 None, §5.1.3 step 7 ("Only the domain format of host is allowed here"). L2 §7.2: step 11 (type `webauthn.get`), 12 (challenge), 13 (origin), 15 (rpIdHash), 16 (UP), 17 (UV), 20 ("verify that sig is a valid signature over the binary concatenation of authData and hash"), 21 (counter: "If authData.signCount is nonzero or storedSignCount is nonzero …"). **Not read (UNVERIFIED):** L3 §7.2's step numbers (BE consistency, topOrigin), §8.7's exact wording, Chrome's `other_keys_can_be_added_here` injection, synced passkeys' zero counter, Windows Hello RS256 |
| V26 | `approver`'s confirm helper, labels and the K8 code | `cat -n cmd/ailang-worldd/approver.go` | `:27 approverConfirmBudget = 120 * time.Second`; `:33-55 confirmApproverAct` (N9a, N9, N10a, N10 with `rerr != nil \|\| !strings.HasSuffix(line, "\n") \|\| !isYes(line)`); `:61 readTermLineErrWithin`; `:261-267 approverKeyLabel` (SSH-only); `:276 approverEnrol` seam; `:376 approverRevoke` seam; `:446-456` K8 |
| V27 | Authority's derivation | `cat -n host/authority/mint.go \| sed -n 55,111p` | `:62-69 RequestMint{Request, RequestCreatedAt, Challenge, ExpiresAt, Signed, Consume}`; `:84` A1; `:88` A2; `:92` A5; `:96` A4; `:99-101` A3 |
| V28 | `sessionEnv` seams (no new seam needed) | `grep -n "sessionEnv struct" -A16 cmd/ailang-worldd/session.go` | `openTerminal`, `now`, `random` ("cryptographically random bytes (crypto/rand in production)"), `after`, `sign` |
| V29 | The designer's shell has no controlling terminal | `( : < /dev/tty ) 2>&1 && echo openable \|\| echo "tty: not openable"` | `device not configured: /dev/tty`, `tty: not openable` |
| V30 | A1b's measured sprint sizes | `for m in b7d2e629 290084ad f0ce26f7 9429e798; do git diff --numstat $m^1 $m \| awk …; done` | code+tests+fixtures `1706`, `1910`, `1923`, `1032` (= 6,571; A1b §0 estimated "about 3,900") |
| V31 [r1] | **The page's `blob:` save works under the exact `surfaceCSP` in Chrome** (quorum r1, gemini-4-argon; throwaway) | the pin: `sed -n 64p host/daemon/surface_fonts_test.go`. The probe: `/tmp/w379-designer/blob/probe.mjs` (sha `efa451f25307bc89`). A Node server sends **exactly** that CSP (plus `nosniff`, `no-store`, `no-referrer`) on an HTML page and its `/p.js`. The script records `securitypolicyviolation` into `window.__v`, builds `new Blob([textarea + "\n"], {type: "application/json"})` and sets `a.href = URL.createObjectURL(…)` and `a.download = "world-registration-probe.json"`. Chrome runs `--headless=new` with a temp profile. DevTools: `Browser.setDownloadBehavior {behavior: allow, downloadPath: <tmp>, eventsEnabled: true}`, `Log.enable`, `Runtime.enable`, then a real `Input.dispatchMouseEvent` press and release on the link. The run is polled to completion under a 90 s `date +%s` bound. **Control:** `probe_control.mjs` (sha `fca80ee51cbd7e04`) is identical except that it also loads `http://127.0.0.1:9/control.png` | the pin: `const wantCSP = "default-src 'none'; style-src 'unsafe-inline'; script-src 'self'; connect-src 'self'; font-src 'self'; base-uri 'none'; form-action 'self'; frame-ancestors 'none'"`. Probe (`out.json`, sha `633b70d559b30bae`): `Chrome/157.0.8103.50`, `setDownloadBehavior: ok`, `downloadWillBegin url: blob:http://localhost:59537/39a9cc81-… suggestedFilename: world-registration-probe.json`, `downloadProgress … state: completed`, `cspViolations: []`, `consoleAndLog: []`, `filesInDownloadDir: [{f: world-registration-probe.json, content: "{\"probe\":true}\n"}]`. Control (`out_control.json`, sha `1244ce7676c68535`): `violations: ['img-src http://127.0.0.1:9/control.png']`, log `Loading the image 'http://127.0.0.1:9/control.png' violates the following Content Security Policy directive: "default-src 'none'" …`, the file still saved, `last state: completed`. So the instrument sees violations, and the `blob:` download raised none. Temp directories removed; `no probe process` |
| V32 [r1] | **The parent's ruling text, verbatim** (quorum r1, oc-kimi-k3) | `grep -n "standard library only\|fmt: none\|message hash is the \`challenge\`\|relying-party ID is\|one-time enrolment code\|No key can enrol without a typed\|POST bound to a verified signature\|only origin enrolled\|sign counters checked" design_docs/planned/w-world-approvals-by-key.md` | `:32 - **Verification uses the standard library only:** \`crypto/ecdsa\` (P-256), \`crypto/ed25519\` and \`crypto/sha256\`. … (a bounded CBOR subset, no attestation, \`fmt: none\`), … There are no new Go dependencies`; `:42 - **For WebAuthn,** the message hash is the \`challenge\`. The relying-party ID is the surface's host name: the tailnet \`*.ts.net\` name, or \`localhost\` on the daemon's machine.`; `:49 … WebAuthn enrolment completes in the browser: the surface shows a one-time enrolment code the operator types into the terminal prompt. **No key can enrol without a typed \`y\`** …` (the same line sketches `approver enroll --name "Mark's iPhone" --format webauthn\|ssh [--pubkey <file>]`, the one DEPARTED-FROM-PARENT shape point in the header); `:66 … Approving is a POST bound to a verified signature, never a cookie or a session.`; `:74 \| Replay or reuse of an approval \| One-time challenge, 10-minute expiry, consumed on use; WebAuthn sign counters checked. \|`; `:77 … The surface is the only origin enrolled.` Also `:27` (table): "WebAuthn assertion (passkeys, FIDO2), user verification required" |
| V33 [r1] | **The ledger's ruling text, verbatim** (quorum r1, oc-kimi-k3) | `grep -n "^\| D-WORLD-89 \|^\| D-WORLD-109 " design_docs/world-mission.md` (lines `1008`, `1028`); `sed -n 1008p … \| grep -o "ANSWERED[^\|]*"`; `sed -n 1028p … \| grep -oE "The loop builds and judges against fakes and test keys, and never enrols a real key, mints or publishes\.\|ANSWERED — A\. Mark, attended 2026-11-06[^(]{0,120}"` | D-WORLD-89: `Mark answered: "all defaults". R1 = YES: for mints and Human.Approve decisions, a verified signature from an enrolled key may approve, as well as the typed y; … R2 = YES: the formats are WebAuthn (passkeys or FIDO2, with user verification) and OpenSSH signatures (SSHSIG); … R3 = YES: enrolment and revocation stay typed-only at the daemon's terminal. R4 = YES: publishing stays typed-only. R5 = YES: …` `(Mark Edmondson, attended 2026-10-21 …)`. D-WORLD-109: `ANSWERED — A. Mark, attended 2026-11-06 in-session: "d-world-109 - a".` and `The loop builds and judges against fakes and test keys, and never enrols a real key, mints or publishes.` ("Never opens a real credentials store" is A1b's restatement, as the header says) |
| V34 [r1] | **The credential-ID length rule** (quorum r1, oc-kimi-k3) | WebFetch of `https://www.w3.org/TR/webauthn-3/`, window 0–100 k, asking for §4's "Credential ID" definition verbatim | "A probabilistically-unique byte sequence identifying a public key credential source and its authentication assertions." / "At most 1023 bytes long." / "Credential IDs are generated by authenticators in two forms:" / "1. At least 16 bytes that include at least 100 bits of entropy, or" / "2. The public key credential source, without its Credential ID or mutable items, encrypted so only its managing authenticator can decrypt it." / "Relying Parties do not need to distinguish these two Credential ID forms." So the 16-byte floor binds authenticators generating form 1 only. The relying-party step is §7.1 step 25 (≤ 1023, V25). D5, S-validation, the DDL `CHECK` and C14 → 1–1023 |
| V35 [r1] | The 128-bit enrolment code's format, on real browser bytes (throwaway) | `/tmp/w379-designer/code/main.go`, `GOFLAGS=-mod=mod GOPROXY=off go run . /tmp/w379-designer/browser/es256.json`: SHA-256 of `"world-enrolment-code/v1\x00" ‖ SHA-256(clientDataJSON) ‖ SHA-256(attestationObject)`, first 16 bytes, `base32.StdEncoding.WithPadding(NoPadding)`, grouped 5-5-5-5-6 | `len: 26 bits encoded: 128 (16 bytes)`; `grouped: YDOR4-X36PQ-NFRTX-BHWAV-4W2XG4`. This is the designer's probe registration, not the committed fixture, so PE4 pins the fixture's own value |
| V36 [r2] | Browser origin serialisation of an explicit default port (quorum r2, `oc-kimi-k3`) | controller, 2026-11-10: `node -e 'new URL(u).origin'` for each URL, and Chrome headless `--dump-dom` of a `data:` page writing `new URL(u).origin` (temp profile, removed) | both print `https://x.ts.net:443/ -> https://x.ts.net`, `https://localhost:443/ -> https://localhost`, `http://localhost:80/ -> http://localhost`, and the control `http://localhost:7644/ -> http://localhost:7644`. Chrome output: `https://x.ts.net | https://localhost | http://localhost | http://localhost:7644`. The default port is elided, so an origin stored with `:443`/`:80` can never equal `location.origin`: refuse it (O3). Note: this measures Chrome's URL origin serializer, which `location.origin` uses; the reviewer's literal `Runtime.evaluate` of `location.origin` on a served `https://<host>:443/` page was not run (it needs TLS on 443) |

## 15. Quorum verification log

- **Round 1** (r0, commit `72e20eb4`, author `claude:claude-opus-5-5`; seats gpt6-1-sol,
  gemini-4-argon, oc-kimi-k3): **BLOCKED 3 of 3**. None disputed the direction. The controller
  measured each premise first, and all three were real.
  - **gpt6-1-sol, on §2.8 (the enrolment code) and the §4 threat row.** The code is derived only
    from the registration file, so an attacker who reads the genuine file can compute it and grind
    a 40-bit collision: "guess 40 bits blind" was false. **Applied in r1** (the reviewer's route):
    - a 128-bit, 26-character fingerprint, stated as public, with the read-then-substitute attacker
      modelled (§2.8);
    - single-attempt comparison, then the typed `y` (WN10);
    - five-group display, paste allowed, normalisation (§2.8, §2.12 step 7);
    - CW5 rewritten to that attacker, plus prefix, length and alphabet arms;
    - PE4's 26-character golden, `GroupEnrolmentCode`;
    - §4, §5.1, §10 step 2, and QUICKSTART §16's text; V35.
  - **gemini-4-argon, on §2.14 step 5 (the `blob:` save) and R3.** The save was UNVERIFIED under
    `surfaceCSP`, with a fallback. **Applied in r1:** measured in Chrome 157 under the exact pinned
    CSP. The download completed, the file matched, and there were zero violations, with a
    violation-detecting control (V31). `surfaceCSP` is unchanged, the textarea is no longer called a
    fallback, and XB1 now asserts the save itself (`setDownloadBehavior`, a real click, `completed`,
    byte-equality, no violations, Go verifying the saved file). R3 narrowed to Safari and Firefox.
  - **oc-kimi-k3, on the header's Authority block, §2.5 D5, §2.10 (DDL, S-validation, C14) and
    §2.10's ordering sentence.** No Verification Log row read the parent's or the ledger's ruling
    text, and the 16-byte credential-ID floor was UNVERIFIED. **Applied in r1:**
    - V32 and V33 quote both. Every clause was found, and the one shape departure is named
      DEPARTED-FROM-PARENT in the header;
    - V34 quotes L3 §4. The floor is an authenticator generation rule, so it is dropped to 1–1023
      everywhere, with AD1, PS2 and CS7 arms at 0, 1 and 1,024 bytes;
    - §2.10 now states that C13–C15 sit immediately after C5 and before `BeginTx`. The unedited
      `consume_binding_test.go:138` arm pins that, with an ordering mutant in CS7.
- **Round 2** (r1, commit `967ed666`; same seats): **BLOCKED 2 of 3**; `gemini-4-argon` **PASS**. Neither reject disputed the direction, and each carried a concrete `proposed_fix`. The controller measured both premises first (rule 3f), and both were real.
  - **oc-kimi-k3, on §2.2 (origin grammar) against §2.14 step 1 (`location.origin`).** An explicit default port was accepted and stored, while browsers elide it from `location.origin`, so such an origin could never enrol. Measured V36 (Chrome 157 and Node's WHATWG URL both elide `:443`/`:80`).
  - **gpt6-1-sol, on §2.14 step 3 (the ceremony wait).** Only WebAuthn's `timeout` hint bounded the browser wait; no AbortSignal or application deadline existed (measured: no `AbortController`/`signal` anywhere in the r1 text).
  - Per round: r1 = enrolment-code secrecy + CSP save + authority/floor; r2 = origin default-port + ceremony abort. Different surfaces each round, no localisation, so no SPLIT signal (gate-2 rule (b)).
- **r2, controller narrow-refinement carve-out** (ratified; precedents iterations 372, 373, 375). Applied verbatim:
  - §2.2: the default-port rule and O3 wording, and the OR1 refusal arms for `https://x.ts.net:443`, `https://localhost:443`, `http://localhost:80` (kimi's primary option, refusal; the normalisation alternative was not taken);
  - §2.14 step 3: the AbortController paragraph, word for word, and XB2 (the reviewer's browser test);
  - V36 records the measurement. One named adaptation: V36 measures Chrome's URL origin serializer through `new URL(...).origin` rather than a `Runtime.evaluate` of `location.origin` on a served `:443` page, which would need TLS on port 443.
  - No other text changed; the planner should treat OR1's three new arms and XB2 as part of M2 and M7 respectively.
