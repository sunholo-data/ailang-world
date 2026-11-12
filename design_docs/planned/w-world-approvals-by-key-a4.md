# w-world-approvals-by-key-a4: decide a `Human.Approve` decision packet with an enrolled key

- **Status**: Planned. Designed in iteration 388, **r2** (r2 by the controller's narrow-refinement
  carve-out: quorum round 2 BLOCKED 2/2 on one surface, the K0 retry order; both reviewers' fixes applied
  in their own text, marked **[r2]**; §15). r0 was **BLOCKED** at quorum round 1
  (2/2 present; §15). r1 applies both reviewers' fixes and nothing else; changed text is marked
  **[r1]**: (1) D4 now requires a **complete** request-wide packet discovery (§2.4, §2.5 D4a/D4b,
  RD9b–RD9d); (2) the `key_decisions` DDL is tightened verbatim and a crashed, unlanded key decision
  is retried **idempotently with its exact artifact**, which also closes the counter-0 replay
  (§2.7, KD5–KD8, RD13). Queue row 201, `w-world-approvals-by-key`,
  parent milestone **A4**. Parent: [`w-world-approvals-by-key.md`](w-world-approvals-by-key.md)
  (attended; its rulings are not reopened). Built siblings, settled and not reopened:
  [`-a0a1`](w-world-approvals-by-key-a0a1.md) (the mint request object, digest, challenge, expiry,
  consume), [`-a1b`](w-world-approvals-by-key-a1b.md) (SSHSIG approvers, `approver
  enroll|list|revoke`, `approve`), [`-a2`](w-world-approvals-by-key-a2.md) (the WebAuthn law),
  [`-a3`](w-world-approvals-by-key-a3.md) r2 (the surface flow A3a: `surface_approvals`, the shared
  judge, the five `/surface` approval routes, *Needs you*, `approve --await`). The row-172 producer
  is [`w-decision-packet-production-projection.md`](w-decision-packet-production-projection.md).
- **Authority**: **D-WORLD-89** (Mark, attended 2026-10-21, "all defaults"): **R1** a verified
  signature from an enrolled key may approve "for mints and `Human.Approve`" decisions, as well as
  the typed `y`; **R2** the formats are WebAuthn (UV required) and OpenSSH SSHSIG; **R3** enrolment
  and revocation stay typed-only; **R4** publishing stays typed-only. **D-WORLD-73 = A** (row 176):
  the daemon binds `Human.Approve` for every episode, serialized, and "never decides a request".
  **D-WORLD-109 = A**: the loop builds and judges against fakes and test keys and never enrols a
  real key, mints or publishes; §10 is Mark's attended run, never the loop's. The parent §3
  invariants bind every section: the agent never approves; publishing stays typed-only; the typed
  `y` keeps working; read routes stay as open as row 199 leaves them; approving is a POST bound to
  a verified signature or a typed proof, never a cookie or a session.
- **Item** (A4): an operator decides, **approve or deny**, a pending `Human.Approve` request that an
  agent's invocation raised in a daemon episode, and whose decision packet the coordinator
  committed (row 172/176). Three ways, one law:
  - an **enrolled SSH key** (`ailang-worldd approve <request-ref> [--deny] --sign-with <key>` or
    `--signature <file>`), or an **enrolled passkey** (the surface card, or `--webauthn-assertion`);
  - the **typed `y`** (`ailang-worldd approve <request-ref> [--deny]` at a terminal), the root.
  The daemon verifies, then lands one immutable `world/approval-decision/v1` through the existing
  approvals head, with `decidedBy` derived from what verified (`ssh:<name>`, `webauthn:<name>`,
  `operator:terminal`). It changes no `.ail` code.
- **Designer lane**: `claude:claude-opus-5-5`.
- **Estimated**: **3 sprints, about 5,250 changed lines** (§9): sprint 1 ≈ 1,950 **[r1]** (+≈ 250 for the complete discovery, DS1, DL4a/b, RD9b–RD9d), sprint 2 ≈ 1,800 **[r1]** (+≈ 200 for K0, KD5–KD8, RD13),
  sprint 3 ≈ 1,500. Calibration (V31): A3a's four sprints measured 1,543, 1,995, 1,861 and 393
  changed lines outside `design_docs/` against plans of 1,650, 1,900 and 1,950 (≈ 1.05×).
- **Measurement base**: `b0fbcd3eee7b5b29190e678e296635a3dd3a0d0c` (PR #561's merge), worktree
  `/Users/voightkampff/dev/sunholo-data/.wt-world-iter388-a4`, branch
  `design/world-iter388-row201-a4`. sha256 prefixes (V1): `host/broker/approve.go`
  `aeef311293e1c42d`, `host/broker/approvals_view.go` `198bb050061c92b9`,
  `host/coordinator/decisionpacket.go` `ff50de9b6af8bf92`, `decisionpacket_validate.go`
  `ac56423782a187aa`, `host/daemon/daemon.go` `0495878c6aedce4b`, `host/daemon/episode_finish.go`
  `4e320c537e689620`, `host/daemon/surface_approve.go` `6f3cafda604f4c23`,
  `host/authority/judge.go` `d3a6620f0aa00022`, `host/approval/request.go` `cf9897ec3f10881e`,
  `host/approval/sshsig.go` `a40850200782b3e4`, `cmd/ailang-worldd/approve.go` `6a38843161c7d75f`,
  `cmd/ailang-worldd/episode_finish_operator.go` `0adb363562679cd5`, `host/surfacepage/app/app.js`
  `d2495278f7b17590`.
- **Instruments**: `go1.26.6 darwin/arm64` (the worktree's toolchain; the scratch probe ran the
  cached `golang.org/toolchain@v0.0.1-go1.26.6.darwin-arm64` with `GOTOOLCHAIN=local`). One
  throwaway Go probe lived in `/Users/voightkampff/.ailang/state/world-iter388/designer-scratch/redecide/`
  (`main.go` sha `d9162ffe5e808aa8`) and is not in the repo (V9). No browser was run in this fire:
  every browser premise is a planner obligation (§1.12).

---

## 0. Scope, the amendment posture, and the sprint boundaries

**What A4 is, measured.** The parent's whole A4 text is one line, "`Human.Approve` decision packets
(row 172) approved the same way" (V2). Today a packet produced on the daemon's normal path shows
`approvals chain: pending` and **no operator command or route can decide it** (§1.1–§1.3, V3–V6).
So A4 is not a second way to do something that exists: it is the first operator decide path for
these packets. The typed `y` arm is built here too, because the parent's "typed `y` as root" has no
decide path to be the root of.

**In A4:**
- **Sprint 1, "the law and the typed root" (M1–M3).** The decision message and its pure law; the
  broker refuses a second decision of the same request; the reserved-scope refusal (publish and
  finish shapes, §2.6) and the publish consumer's key-decision refusal (§2.6 L2); `GET
  /surface/decisions/{ref}` (the subject); `POST /surface/decide/{ref}` with the **typed arm only**;
  `ailang-worldd approve <request-ref> [--deny]` typed. After sprint 1 an operator at a terminal can
  approve or deny a packet. No key can decide anything yet.
- **Sprint 2, "the keys" (M4–M5).** Credentials schema 8 `key_decisions` (the evidence and the
  passkey counter); the decision judge in `host/authority`; the POST's `sshsig` and `webauthn` arms;
  the CLI's `--sign-with`, `--signature`, `--write-message` and `--webauthn-assertion` on a request
  ref. All judged with `approvaltest` keys and recorded authenticator fixtures.
- **Sprint 3, "the surface" (M6–M7).** `requests.json` gains `decisions`; *Needs you* lists pending
  packets; the decision card `GET /surface/decide/{ref}` + `decide.js` (approve / deny with a
  passkey or a pasted SSHSIG); the live-Chrome tests on the A2/A3 driver; QUICKSTART §18
  *(attended — pending first verbatim run)*.

**Amendment posture (r0).** The parent is attended. A4 departs from its §2.2 text in one place, and
says so here rather than later:
- **(a) The decision message is not the mint message, and its challenge is derived, not stored.**
  Parent §2.2 says a `Human.Approve` packet "carries … a one-time challenge: 32 random bytes from the
  daemon; an expiry: 10 minutes by default", and that "the signed message is a fixed-prefix encoding
  of (`"world-approval/v1"`, request digest, challenge, expiry)"; for SSHSIG "the namespace is
  `world-approval@sunholo`" (V2). A4 instead signs `"world-decision/v1" ‖ 0x00 ‖ verb ‖ packet digest
  ‖ request digest ‖ expiry` (91 bytes) under SSHSIG namespace `world-decision@sunholo` (§2.3).
  - **Why a different message.** The mint message has no verb (V10): a mint has only "approve". A
    decision is approve **or deny**, so the verb must be signed, or a captured deny could be posted
    as an approve. A distinct prefix and namespace also make a mint signature unusable as a decision
    and the reverse (domain separation, §5.2).
  - **Why no stored random challenge.** The approval request object has no challenge or expiry and
    its wire is strict (`DisallowUnknownFields`, V4); the packet's v1 bytes are frozen ("Amendments
    are /v2 semantic IDs", V7). A stored challenge would need a new credentials table written by an
    issuing route. A4 gets the challenge's two jobs another way: **freshness** from a signer-chosen
    expiry that the daemon bounds to `now < expiry ≤ min(now + 600, packet deadline)`; **single use**
    from a new broker rule that a request is decided at most once, enforced under `humanMu` on the
    durable approvals chain (§2.4). The packet's object digest (unique per packet, minted by the
    coordinator, not chosen by the signer) fills the challenge's slot.
  - **Ledger ask:** **D-WORLD-115** (§12), default **A** (keep the derived challenge and the
    decision message); **B** = the parent's literal stored random challenge, priced at ≈ +550 lines
    (a schema-8 issuance column set, an issuing POST, its tests) in sprint 2. If Mark answers B, only
    §2.3's challenge slot and sprint 2 change; the verb, the prefix and the namespace stay, because
    deny needs them under either answer.
- **No other departure.** The reserved-scope refusal (§2.6) narrows nothing the parent allows: R4
  says publishing stays typed-only, and A4 sends publish approvals back to the typed command that
  owns them. The first-decision-final rule (§2.4) changes broker behaviour that no production caller
  reaches today (V3, V9). Neither is put to Mark; the quorum may disagree, and §15 records it.

**Out of A4:** binding `Human.PollApproval` in the daemon (an effect-surface widening that, like
row 176's binding, needs its own ruling; §2.9); escalation, timeout transitions, or any
`EscalateBounded`/`ExecuteIfGranted` policy (the packet's v1 policy is always `Cancel`, V7); a
decision made by an agent, ever; enrolment or revocation from the surface (R3); any change to
mints, the A0–A3a law, `surface_approvals`, `/surface/passkey`, `approve.js`, publishing, or
`episode finish`; the Tauri app (A3b, row 188 N3).

**Sprint boundaries.** Sprint 1 lands the law, the broker rules and the typed decide path; no key
can reach a decision. Sprint 2 adds keys behind the same law, CLI first. Sprint 3 adds the pages.
Nothing flows backwards: the law and the store are reviewed before anything can reach them (A1b §0,
A2 §0, A3 §0).

## 1. What the code does today (measured at `b0fbcd3e`)

1. **The operator decide entry point has no production caller.** `broker.DecideApproval`
   (`host/broker/approve.go:181`) wraps `decideApproval` (`:191`). In non-test Go the only calls of
   `decideApproval` are `RecordOperatorApproval` (`:696`) and `mintAttendedApproval`
   (`publish_op.go:153`); 5 `_test.go` files call `DecideApproval(` (V3).
2. **Both production deciders create the request they decide.** `RecordOperatorApproval` has one
   caller, `host/daemon/episode_finish.go:173` (`"operator:episode-finish"`, decided by
   `"operator:terminal"`), and puts its own request (`approve.go:686-697`). `MintAttendedApproval`
   has one caller, `cmd/world-publish/main.go:310`, which mints the request through `Human.Approve`
   and decides it in the same process (`publish_op.go:129-174`). `cmd/world-publish/tty.go:15-16`:
   "it is not a HEADLESS/ATTENDED fence — a headless loop can mint one by calling
   broker.DecideApproval. It is a SINGLE-USE fence." (V5)
3. **A packet on the daemon path stays pending.** The daemon's binder binds `Human.Approve` only
   (`daemon.go:954`, `reg[broker.EffectHumanApprove] = d.human`), serialized by `humanMu`
   (`:438-443`, `serializedHandler` `:959-968`). The coordinator mints a packet only for one allowed,
   successful `Human.Approve` (`effectful.go:340-352`). QUICKSTART documents "it never decides a
   request", "decisions stay operator-only" (`docs/QUICKSTART.md:1593`) and that the page "shows
   `approvals chain: pending` until an operator decides" (`:1603`). With §1.1–§1.2 there is no
   operator decide path for such a packet. **The controller's inference is CONFIRMED** (V3–V6).
4. **The approval request object.** `Human.Approve` stores `world/approval-request/v1` =
   `{effect, scope, cost, requester, now}` (`approve.go:66-72`, `:141-145`), where `effect` is always
   `"Human.Approve"`, `scope`/`cost`/`now` are the invocation's and `requester` is the agent's
   payload. It carries **no challenge and no expiry**. Every reader decodes it strictly
   (`decodeApprovalJSON`, `DisallowUnknownFields`, `:320-327`). Object hashes are `SumSHA256(payload)`
   (`host/broker/record.go:132-140`), so a page can recompute a ref from the payload bytes (V4).
5. **Identical requests collide.** Two `Human.Approve` executions with the same scope, cost,
   requester and `now` return **the same** `requestRef` (probe, V9). The coordinator's `now` is Unix
   seconds (`daemon.go:876`), so two episodes with the same grant scope and requester in the same
   second share one request.
6. **A request can be decided twice, and the newest decision wins.** `decideApproval` checks only
   that the request exists and is on the chain (`approve.go:198-212`); it appends a second decision.
   The probe decided one request `deny` and then `approve` with no error; `Human.PollApproval` then
   returned the `approve`, and `RecentApprovals` showed `status=approve by=ssh:mallory` (V9).
   `FoldApprovals` keeps the newest decision per request (`approvals_view.go:37-51`).
7. **`Human.PollApproval` does not bind the request to the poller's scope.** A poll under scope
   `some-other-scope` returned the decision of a `worktree` request (V9; the handler,
   `approve.go:154-172`, reads only `requestRef`). It is bound in no daemon binder (V6), so no agent
   reaches it today; it is a residual for whoever binds it (§11 RK6).
8. **The decision packet** (`world/decision-packet/v1`) is `{packetHash, proposalHash, requestRef,
   createdAt, deadlineAt, escalationsRemaining, policy}`, compact JSON in that frozen order;
   `deadlineAt` is the authorizing grant's expiry and `policy` is always `Cancel`
   (`decisionpacket.go:90-181`). `ValidatePacket` runs seven ordered checks (bytes, carried,
   proposal, request, anchored, lifecycle, pair-unique) and returns the chain's `ApprovalSummary`
   (`decisionpacket_validate.go:13-38`). The anchor reads the newest 40 heads and the duplicate scan
   the first 100 packets (`QUICKSTART.md:1603`, `packetScanLimit = 100`). Packets list by
   `ObjectsBySemanticID`, which orders by `hash_ref`, not time (`objects_by_semantic_id.go:20`) (V7,
   V8).
9. **Publish consumption reads any approve decision with a publish-shaped scope.**
   `validatePublishApproval` (`approve.go:517-650`) follows the payload's `approvalRef` to an
   `approve` decision, its request (`effect == "Human.Approve"`, `cost == PublishCost` = 1,
   `registry_publish.go:31`), then `ParsePublishApprovalScope` and the byte terms. **It never reads
   `decidedBy`.** `world-publish publish --live --approval-ref R` names R (`main.go:347-351`), and
   `--decided-by` is free text (`:172`). Both `approve` and `publish` sit behind
   `requireAttendedOperator` (`main.go:294`, `:399`). So a `Human.Approve` request with a
   publish-shaped scope that **any** decider approves is a valid publish stamp (V11).
10. **The finish record reads any approve decision with a finish-shaped scope.**
    `finishDecisionScope` (`episode_finish.go:243-280`) follows a decision ref to its request and
    `ParseFinishApprovalScope`; it reads neither `requester` nor `decidedBy` (V12). The route itself
    is behind the operator proof (§1.11).
11. **The typed operator fence that exists for a daemon-held write.** `POST
    /v1/episodes/{episode}/finish` refuses any `Authorization` header (`episode_finish.go:139-143`),
    takes `humanMu`, then `consumeOperatorProof`: a 0700 `operator-proofs/` directory beside the
    credentials store, an owned 0600 `<name>.json` holding `{sha256, phase, expires_at}` (≤ 300 s),
    the header `X-World-Operator-Proof` hashing to it, and an unlink "before every store mutation"
    (`:90-133`, `:158-163`). The CLI writes it after the typed `y` (`writeOperatorProof`,
    `cmd/ailang-worldd/episode_finish_operator.go:30-73`; POST `:75-96`) (V13).
12. **Browser premises: none measured in this fire.** A3 measured (A3 V17, V23–V26) that Chrome 157
    headless rebuilds A3's 90-byte message with `crypto.subtle` byte-equal to Go and sends
    `Origin`/`Sec-Fetch-Site: same-origin` on a same-origin JSON POST. A4's 91-byte message, its
    verb byte and its second digest are **not** measured: **planner obligation P1** (§9 sprint 3,
    before M6 is planned) is to re-run A3's V23 probe shape with A4's message and two golden
    subjects, and record it as a V-row in the sprint plan.
13. **The A3a surface.** `approvalRoute` applies the Host allowlist to the five approval routes
    (`surface_approve.go:76-87`); `handleSurfaceApprove` runs R1, R3–R14 then a non-waiting
    `store.OpenCredentials` "closed at once" (`:191-296`). That is the daemon's **only**
    `OpenCredentials(` call in non-test code (`:280`), pinned by `TestDaemonCredentialWriterSites`
    (`surface_approve_post_guard_test.go:132`). `RequestsJSON` is `{v, now, pending, decided,
    unavailable}` (`host/surfacepage/approve_model.go:64-70`). `app.js` (1,470 lines) renders *Needs
    you* from `/surface/requests.json` (`app.js:424-429`, `:921-926`) (V14, V15).
14. **The route tables.** `/v1` is pinned at **13** routes, identical to `worlddapi.ail`'s
    `routes()` (`route_table_test.go:51-55`; 13 `path: "/v1/` rows in the sketch); the thirteenth
    came with a ruling (D-WORLD-100 R6, row 218) (V16). `surfaceCSP` is unchanged since A3
    (`surface_page.go:19`) (V17).
15. **The keys' law and its reusable parts.** `SignedMessage` is `MessagePrefix ("world-approval/v1")
    ‖ 0x00 ‖ digest ‖ challenge ‖ u64be(expiry)`, 90 bytes (`request.go:26`, `:209-222`).
    `VerifySSHSIG(s, namespace, message, key)` takes the namespace as a parameter (`sshsig.go:262`);
    `VerifyApproval` fixes `Namespace = "world-approval@sunholo"` (`:19`, `:353-363`).
    `ApproveByAssertion` computes `MessageHash(StoredDigest, Challenge, ExpiresAt)` and compares it
    with the client data's challenge (W6), then origin (W7), cross-origin (W8), RP-ID hash (W10) and
    the counter (W15, `ad.SignCount <= p.SignCount` refused when either is non-zero)
    (`passkey_assert.go:33-94`). `checkApprover` (P2–P5) is shared (`signed.go:75-92`). The judges
    take a `store.MintRequestRow` (`judge.go:49-129`), so they cannot judge a decision unchanged (V10).
16. **Grants are exact.** `Decide` matches effect and scope by string equality
    (`host/broker/decide.go:110-111`); `--grant Human.Approve=worktree:1` is the documented grant
    (`QUICKSTART.md:1654`). Whatever scope string an operator grants, an agent can raise exactly that
    scope (V18).

## 2. Design

### 2.1 The shape: the daemon decides, because the daemon holds the world store

```
agent invocation ─► Human.Approve (serialized) ─► approval request + approvals head   [today, row 172/176]
                 └► coordinator commit carries world/decision-packet/v1                [today]
operator:
  CLI   approve <request-ref> [--deny] [--sign-with k | --signature f | --webauthn-assertion f]
  page  /surface → Needs you → /surface/decide/<request-ref>
        GET /surface/decisions/<request-ref>  ─► packet bytes, request bytes, plan ref, entry, now,
                                                 the window, the law's pre-check, passkey groups
        the client checks both SHA-256 refs, renders every field, builds the 91-byte message
        signs (passkey / ssh-keygen -Y sign -n world-decision@sunholo) or asks y at /dev/tty
        POST /surface/decide/<request-ref>   ─► guards → THE DECISION LAW (§2.5) under humanMu →
                                                 [key: key_decisions row, counter] →
                                                 broker.DecideApproval (first decision only)
agent: reads the decision on the approvals chain (/v1/registry, /v1/objects, the packet page)
```

**Why the daemon decides here, when in A3 the CLI mints.** A3 moved the mint to the CLI because a
mint is `git worktree add` plus the credentials writer plus a token file, all in `package main`
(A3 §2.1). A decision is one object and one head move in the **world** store, which the daemon
holds as the sole writer (`world-publish` refuses with `IsWriterAlreadyActive` when a writer is
active, `main.go:298-305`, V11). So the daemon writes it, as it already writes the finish approval
(§1.11). The CLI never opens the world store for A4.

**One route, three formats.** `POST /surface/decide/{ref}` takes exactly one of `sshsig`,
`webauthn` or `typed`. The typed arm needs the row-236 operator proof (§1.11); the key arms need a
verified signature. Putting all three behind one route and one law is the reuse the directive asks
for, and keeps `/v1` at 13 (no sketch edit, no new ruling; V16). A `/v1/approvals/{ref}/decide`
route was considered and rejected: it widens the frozen `/v1` table, which took a ruling last time.

### 2.2 What is signed: the subject

The subject of a decision is **the packet and its request, both by content**:
- the **packet object digest** `P` = the 32 bytes of the packet object's ref (`SumSHA256` of its
  payload). The packet binds `proposalHash` (the invocation's plan), `requestRef`, `createdAt`,
  `deadlineAt`, `escalationsRemaining` and `policy` (§1.8);
- the **request object digest** `Q` = the 32 bytes of `requestRef`. The request binds `effect`,
  `scope`, `cost`, `requester` and `now` (§1.4).

`P` already binds `Q` through `requestRef`; `Q` is signed as well so the message names the thing the
chain is keyed on, and so a check `packet.requestRef == Q` is a byte comparison, not a parse.

**"The exact request"** is therefore: the effect (always `Human.Approve`), the scope (the authority
being asked for), the cost, the requester, the logical time, the proposal plan's ref (whose content
is bound by its hash and shown via the existing object page), the packet's deadline and policy. The
evidence roster the workbench shows (`packetRosterFromValidated`) is derived from the plan and its
invocation record; it is bound through `proposalHash`, displayed by link, and **not** part of the
card's fail-closed field table (§2.8), because it is not in either signed object's bytes.

### 2.3 The decision message (`host/approval/decision.go`, pure)

```
DecisionMessage(verb, P, Q, expiry) =
    "world-decision/v1"  (17 bytes, DecisionPrefix)
  ‖ 0x00
  ‖ verb                 (1 byte: 0x01 approve, 0x02 deny; anything else refused)
  ‖ P                    (32 bytes, the packet object digest — the challenge's slot)
  ‖ Q                    (32 bytes, the request object digest)
  ‖ u64be(expiry)        (8 bytes, Unix seconds)
  = 91 bytes
DecisionMessageHash = SHA-256(DecisionMessage)        -- the WebAuthn challenge
SSHSIG namespace    = "world-decision@sunholo"         -- DecisionNamespace
```

- **Domain separation, measured against today's encoding (V10).** The mint message is 90 bytes
  starting `"world-approval/v1" 0x00`; the decision message is 91 bytes starting `"world-decision/v1"
  0x00`. The prefixes are the same length and differ at byte 6, so neither message is a prefix or
  a suffix-extension of the other, and their SHA-256 values differ. For SSHSIG the namespace is part
  of the signed blob, so a mint SSHSIG fails `VerifySSHSIG(…, DecisionNamespace, …)` on the namespace
  check (`ErrNamespace`) before any message comparison; for WebAuthn the challenge is the message
  hash. Tests DM3/DM4 (§6) prove both directions with `approvaltest` keys.
- **Why the verb is signed.** Without it, one signature would decide either way.
- **What changes in existing code: nothing.** `SignedMessage`, `MessageHash`, `Namespace`,
  `VerifyApproval` and `ApproveBy*` are untouched (the mint law is not re-opened). A4 adds
  `DecisionMessage`, `DecisionMessageHash`, `DecisionNamespace`, `VerifyDecisionSSHSIG`, and two
  laws `DecideBySignature` / `DecideByAssertion` that reuse `checkApprover` (P2–P5) and the
  assertion's W1–W15 steps through one unexported function parameterized by the expected challenge
  hash (a refactor that keeps `ApproveByAssertion`'s behaviour; every A2 test stays green unedited).
- **`SignedDecision`** mirrors `SignedApproval`: unexported fields, `ok` set only by the two laws,
  `Binds(verb, P, Q, expiry) error`, `DecidedBy() string` = `format + ":" + name`, `Use()` and
  `PasskeyUse()` for the evidence row (§2.7).

### 2.4 Single use, durable across restart

- **First decision is final (broker).** `decideApproval` gains one check after the request lookup:
  if `findApprovalDecision(requestRef)` finds a decision → `ErrApprovalAlreadyDecided` (a new
  sentinel), nothing written. The two production deciders decide a request they have just created
  (§1.2), so they are unaffected; the probe's deny-then-approve (V9) becomes a refusal.
- **Serialization.** `appendApprovalHead` reads then sets without compare-and-set
  (`daemon.go:438-440`). The decide handler takes `d.humanMu`, as `handleEpisodeFinish` does
  (`episode_finish.go:158`) and as `serializedHandler` does for `Human.Approve`, so a decision, a
  new request and a finish approval never interleave. Inside the lock the law re-reads the chain;
  "undecided" is decided under the same lock as the write.
- **Durable.** The decision object and the head are in the SQLite world store; after a restart the
  chain still holds the decision, so a replayed signature (same verb, same subject) is refused by
  `ErrApprovalAlreadyDecided`, whatever its expiry. No in-memory state carries single use.
- **Ambiguous subjects are refused.** Because identical requests collide (§1.5), a request ref may
  be carried by more than one packet (different proposals). The law refuses a request ref that is
  not carried by **exactly one** validated packet: "request <ref> is carried by N decision packets;
  deciding it would decide all of them; it cannot be decided here". (One decision on the chain
  would answer every invocation that shares the request; A4 refuses rather than guess.)
- **[r1] "Exactly one" needs a complete count (gpt6-1-sol's fix).** There is no request index for
  packets: they list only by `ObjectsBySemanticID`, ordered by `hash_ref`, and `ValidatePacket`'s
  own pair-unique scan stops at `packetScanLimit = 100` and still returns code 0 with
  `UniqueTruncated` set (`decisionpacket_validate.go:34`, `:190-215`, V33). So a second packet
  carrying the request can sort outside any bounded prefix. A4 therefore does **not** reuse that
  scan or the listing's view. Under `humanMu`, the decide handler runs its own **discovery**:
  keyset pages of `ObjectsBySemanticID(DecisionPacketV1, after, 500)` (the store's
  `MaxSemanticIDPage`, a unique total order on the primary key, V34), following the cursor, each
  packet read with `GetObject` and decoded, up to `decisionDiscoveryLimit = 2000` packets **plus one
  sentinel read**. It yields `Discovery{Complete, Scanned, Carrying, Undecodable}`:
  - `Complete` is true **only** when the cursor ran out (a page shorter than its limit) before the
    sentinel; reaching the sentinel means the set is not established;
  - `Carrying` counts every decoded packet whose `requestRef == Q`, whatever its validity;
  - `Undecodable` counts packets that did not decode (a packet that cannot be read cannot be shown
    not to carry `Q`).
  Packet discovery must establish a complete request-wide candidate set. If the bounded global scan
  reaches its overflow sentinel, or `Undecodable > 0`, the law refuses with 409 (D4a); it never
  treats one visible match as exactly one match. Then `Carrying` must be exactly 1 and that packet
  must validate with code 0 (D4b). **Listing truncation warnings do not authorize decisions**: the
  `requests.json` `decisions` view (§2.8) is display only, and the POST always re-discovers.
  The cost is at most 2,001 envelope rows and 2,000 payload reads per decide POST, all under the
  read deadline (RK3b names the cap).
- **The row-176 residual.** A `Human.Approve` that lost its commit leaves a request and a pending
  head with **no** committed packet. Such a request has zero validated packets, so the law refuses
  it ("no committed decision packet carries request <ref>"). *Needs you* never lists it (§2.8).
- **The 40-head bound.** `ValidatePacket`'s anchor reads the newest 40 heads (§1.8). A packet whose
  request is older than that is refused with its validator sentence (code 5, "not established"),
  never decided on a guess. RK3 (§11) names the cost.

### 2.5 The decision law (pure; `host/broker/decision_law.go`)

The facts are gathered by the daemon (IO); the law is a pure function of them, in this order, the
first failure winning. Each refusal has its own mutant (§6).

| id | refused when | reason (text names the fix) | HTTP |
|---|---|---|---|
| D1 | `{ref}` is not `sha256:` + 64 lowercase hex | not a request ref | 400 |
| D2 | the object is absent, not `world/approval-request/v1`, does not decode strictly, or its `effect` ≠ `Human.Approve` | no such approval request | 404 |
| D3 | it is not on the approvals chain (`findApprovalRequest`) | not an approval request of this world | 404 |
| D4a **[r1]** | the discovery is not complete: the scan reached its overflow sentinel (`!Complete`), or `Undecodable > 0` (§2.4) | "the decision packets carrying request <ref> cannot be counted completely (more than 2000 stored packets, or one that does not decode); it cannot be decided here" | 409 |
| D4b | `Carrying ≠ 1` (0: row-176 residual; ≥ 2: §2.4), or the one carrying packet's `ValidatePacket` code ≠ 0 (its sentence; a code-5 bound names "not established") | the count, or the packet sentence | 409 |
| D5 | **reserved scope** (§2.6): publish-shaped or finish-shaped | "publish approvals are minted by `world-publish approve` at a terminal (R4)" / "finish approvals are made by `ailang-worldd episode finish`" | 403 |
| D6 | already decided | "request <ref> was already decided <verb> by <decidedBy>" | 409 |
| D7 | `now ≥ packet.deadlineAt` | "the packet's deadline passed at <t>; its policy is Cancel" | 409 |
| D8 | the posted `packet` / `request` digests ≠ the stored `P` / `Q` | "the request changed since it was shown; reload it" | 409 |
| D9 | verb ∉ {approve, deny} | bad verb | 400 |
| D10 | expiry outside `now < expiry ≤ min(now + 600, deadlineAt)` | "the decision's expiry is stale or too far ahead (at most 10 minutes, and not after the packet's deadline)" | 409 |
| D11 | format arm (below) refuses | the arm's reason | 422 / 403 |

**Format arms (D11):**
- **sshsig / webauthn**: `authority.JudgeDecisionSSH` / `JudgeDecisionAssertion` (§2.7) over the
  91-byte message for `(verb, P, Q, expiry)`; then `SignedDecision.Binds(verb, P, Q, expiry)`. The
  approver is looked up by key on a read-only credentials handle; revoked, unknown, wrong format,
  clock behind, bad signature, wrong origin, RP hash, UV absent, counter not advancing → 422 with
  the law's reason, the same texts as the mint path where the step is shared.
- **typed**: `consumeOperatorProof("decision-" + hex(Q), "decide-" + verb, proof)` under `humanMu`,
  verbatim from the finish route (§1.11). The phase carries the verb, so a proof made for `deny`
  cannot decide `approve`. Any failure → 403 with the proof's reason.

**`decidedBy` is derived, never supplied.** The POST body has no `decidedBy` field (strict decode
refuses one). The handler passes `signed.DecidedBy()` (`ssh:<name>` / `webauthn:<name>`) or the
constant `operator:terminal` to `broker.DecideApproval`.

**Pure-core form.** `DecisionRefusal(f DecisionFacts) (code int, reason string)` where
`DecisionFacts` holds booleans and values the daemon read (ref grammar, request decoded, on chain,
**[r1]** the discovery result `Discovery{Complete bool, Scanned, Carrying, Undecodable int}` and
the carrying packet's validator code, scope class, decided-by, now, deadline, posted digests, verb,
expiry). The format arm runs after it in the daemon. The law has no IO and no clock. **[r1]** D4a is
evaluated from `Discovery` alone, before D4b reads `Carrying`, so no count from an incomplete scan
reaches D4b.

### 2.6 R4: publishing stays typed-only, by construction (two layers)

**What makes a decision a publish stamp today (§1.9).** An `approve` decision whose request has
`effect = Human.Approve`, `cost = 1` and a scope that `ParsePublishApprovalScope` accepts with
`effect=Registry.Publish`. An operator can grant `Human.Approve` on any exact scope string (§1.16),
so an agent can raise a publish-shaped request if such a grant exists; an A4 decision of it would
be a publish stamp. Two layers stop that.

- **L1, the A4 law refuses reserved scopes for every format, typed included (D5).**
  `broker.ReservedApprovalScope(scope) (owner string, reserved bool)`, pure:
  - **publish**: `ParsePublishApprovalScope(scope)` succeeds, **or** the scope starts with
    `registry:`, **or** contains `effect=Registry.Publish` (near-miss shapes fail closed: a future
    parser change cannot make a refused shape acceptable);
  - **finish**: `ParseFinishApprovalScope(scope)` succeeds, **or** the scope starts with
    `github-pr:`.
  Typed is refused as well because these scopes have their own typed owners (`world-publish
  approve`, `episode finish`) that show the exact bytes and run their own fences; A4's card is not
  their card. *Needs you* lists such a packet with `"reserved":"publish"` and the owning command,
  and no decide control.
- **L2, the publish consumer refuses key decisions.** `validatePublishApproval` gains one check
  after `decision.Decision == "approve"`: a `decidedBy` that starts with `ssh:` or `webauthn:` →
  `ErrPublishApprovalDenied` ("publishing stays typed-only (D-WORLD-89 R4): decision <ref> was made
  by <decidedBy>"). Every key decision carries exactly such a `decidedBy`, because the handler takes
  it from `SignedDecision.DecidedBy()` (§2.5), which is `format + ":" + name` as `ApprovedBy` is
  (`signed.go:113-119`), and the only formats are `"ssh"` and `"webauthn"` (`approver.go:16`,
  `:18`, V19). So even if
  L1 were deleted, a key-decided publish-shaped request fails at the consumer.
  - **Cost to existing callers:** `world-publish --decided-by` is free text (§1.9). An operator who
    types `--decided-by ssh:mark` would now be refused at publish. V19 measured that no test or
    doc uses such a value; the refusal names the reason. This is the only behaviour change on the
    publish path, and it can only refuse.
- **L2 for finish is not added.** The finish record is reachable only behind the typed operator
  proof (§1.11), and its approve phase creates its own request. L1 suffices; RK5 records it.
- **Tests:** R4a–R4e (§6), including a must-die control (L1 deleted → L2 still refuses) and a
  must-live control (a typed `world-publish approve` stamp still publishes in the existing
  `registry_publish_test.go` arms, unedited).

### 2.7 The keys (sprint 2): judge, evidence, counter

- **Judge** (`host/authority/judge_decision.go`): `JudgeDecisionSSH(ctx, ro, in DecisionJudgeInput,
  armored) (*approval.SignedDecision, error)` and `JudgeDecisionAssertion(…, response)`, the same
  steps and `*JudgeError` step names as `judge.go` (parse, store, unknown-key, stored-key, law), so
  CLI texts are shared. `DecisionJudgeInput{Verb, P, Q, Expiry, Now}`; no `MintRequestRow`.
- **Credentials schema 8, `key_decisions`** (`host/store/key_decisions.go`): one row per key
  decision.
  ```sql
  -- Key decisions of Human.Approve requests (version 8, row 201 A4). The evidence of what verified,
  -- kept beside the approvers; the authority is the world store's decision, never this row.
  CREATE TABLE key_decisions (
      request_ref   TEXT PRIMARY KEY CHECK (length(request_ref) = 71 AND request_ref GLOB 'sha256:*'),
      packet_ref    TEXT NOT NULL CHECK (length(packet_ref) = 71 AND packet_ref GLOB 'sha256:*'),
      verb          TEXT NOT NULL CHECK (verb IN ('approve', 'deny')),
      format        TEXT NOT NULL CHECK (format IN ('ssh', 'webauthn')),
      approver      TEXT NOT NULL CHECK (approver <> ''),
      artifact      BLOB NOT NULL CHECK (length(artifact) BETWEEN 1 AND 16384),
      expiry        INTEGER NOT NULL,
      decided_at    INTEGER NOT NULL,
      decision_ref  TEXT CHECK (decision_ref IS NULL OR (length(decision_ref) = 71 AND decision_ref GLOB 'sha256:*'))
                                  -- NULL until the world decision landed
  )
  ```
  **[r1]** Both `packet_ref` and `decision_ref` CHECKs are gemini-4-argon's text, verbatim (KD5, KD6).
  `RecordKeyDecision(ctx, row, signed, now)`, one transaction: K1 deadline, below v8 →
  `OlderCredentialsStoreError{Feature: "key decisions", Need: 8}`; K2 grammar; K3+K4 the mint
  consume's own in-transaction helpers, **reused unchanged**: `useApprover` (the approver still
  active, bound by key, `last_used_at = now`; `approvers.go:245`) and, for webauthn, `usePasskey`
  (`UPDATE webauthn_credentials SET sign_count = ? WHERE name = ? AND credential_id = ? AND
  (sign_count < ? OR (sign_count = 0 AND ? = 0))`, `passkeys.go:221-230`; C16–C18) — so a revoke
  since the read-only judge, or a counter that did not advance, refuses inside the transaction; K5
  `INSERT`, or replace a row whose
  `decision_ref IS NULL` (a verified decision that did not land). `LandKeyDecision(ctx, request,
  decisionRef)` sets `decision_ref`.
- **Order inside the handler (under `humanMu`).** law D1–D10 → judge (read-only) → writer open
  (non-waiting; busy → 503 `Busy`, `Retry-After: 1`) → `RecordKeyDecision` → `broker.DecideApproval`
  → `LandKeyDecision` → close. A crash after the decision and before the land leaves `decision_ref`
  NULL: the chain is the authority, so *Decided* still shows it, and `approve --list --decided`
  marks the evidence "not linked". No path decides twice (D6).
- **[r1] A crash, or a failed world write, between `RecordKeyDecision` and `DecideApproval`
  (gemini-4-argon's objection; the reviewer's second option, chosen).** The unlanded row (`decision_ref
  IS NULL`) has already moved the counter, so a retry with the same hardware-counter assertion would
  fail `usePasskey` (stored `N`, assertion `N`: `sign_count < ?` is false; V32). r1 makes the retry
  **idempotent on the exact artifact**:
  - **K0 (new, inside the handler, before K3/K4).** **[r2] Superseded by the two-path judge below:
    r1 put the judge's W15 counter check before K0, so for a hardware-counter passkey the retry
    died at W15 (stored `N`, assertion `N`) before K0 could recognise the artifact — K0 was
    unreachable on exactly the path it was added for (quorum round 2, both reviewers).**
  - **[r2] The two-path decision judge (gpt6-1-sol's fix, verbatim):** "Replace the retry paragraph
    with an explicit two-path decision judge: verify the signature, challenge binding, origin, RP
    hash, UV, and approver identity/status first, but defer counter acceptance. On the non-waiting
    credentials writer, transactionally reload the approver and passkey and compare the existing
    unlanded row against the exact artifact and all signed subject fields. Only an exact match to
    previously committed evidence may bypass counter advancement; otherwise enforce the existing W15
    and usePasskey rules. Recheck revocation transactionally even on the exact-match path. Leave the
    mint judge unchanged."
    Concretely (the controller's reading of that fix, audited VALID — §15): `JudgeDecisionAssertion` runs W1–W14 and P2–P5 on the read-only handle and does **not**
    run W15; it returns the verified assertion's counter. Inside `RecordKeyDecision`'s transaction:
    K0 reloads the approver (`useApprover`'s active/unrevoked check runs on both paths) and the
    passkey row, then — **exact match** (an unlanded row, `decision_ref IS NULL`, whose `artifact` is
    byte-equal to the POSTed one and whose `verb`, `packet_ref`, `request_ref`, `expiry`, `format`
    and `approver` equal the verified ones) → skip `usePasskey`, no counter move, go to
    `broker.DecideApproval` + `LandKeyDecision`; **no exact match** → W15 (the A2 law's counter rule, `passkey_assert.go:88`) and `usePasskey`, then K5. The A2
    verifier gains a defer-W15 form for this judge only; the mint path keeps calling the full W1–W15.
    A counter refused here, in the writer, answers **422** with the same `ErrSignCountRegression` text
    the read-only judge uses, and writes nothing (the transaction rolls back). SSHSIG has no counter: its exact-match path still
    re-checks the approver. gemini-4-argon's equivalent fix ("check for an existing unlanded
    `key_decisions` row … BEFORE running W15 in the judge … skips W15 … and K0 then skips
    `useApprover`/`usePasskey`") is the same ordering; the one difference — gemini also skips
    `useApprover` — is resolved by gpt6-1-sol's "Recheck revocation transactionally even on the
    exact-match path", which is the stricter of the two and is kept. The mint judge (`judge.go`) is
    unchanged.
    The law D1–D10 still runs first on every retry, and D6 still refuses if the decision already
    landed.
  - **Otherwise** (no row, or a different artifact) the record path runs as above: K3/K4 move the
    counter, and K5 replaces an unlanded row only through that path. So a **fresh** assertion
    (counter `N+1`) also lands, by replacement; a different artifact can never skip the counter.
  - **The counter-0 replay, measured.** `usePasskey`'s SQL is `… AND (sign_count < ? OR (sign_count
    = 0 AND ? = 0))` (`passkeys.go:226`), and the A2 law refuses only `(ad.SignCount != 0 ||
    p.SignCount != 0) && ad.SignCount <= p.SignCount` (`passkey_assert.go`, W15; V32). So for an
    authenticator that always reports 0 (stored 0, assertion 0), **the counter accepts the same
    assertion twice**. What refuses a second *use*: (1) once landed, D6 on the world chain refuses
    any decision of that request, before the judge, under `humanMu` (BR1, RD8); (2) while unlanded,
    the same artifact takes K0's idempotent path, which lands the **one** decision that artifact
    signs and writes no second row (the request is the primary key); (3) the artifact signs one
    `(verb, P, Q, expiry)`, so it cannot decide another request or the other verb (DM5, DM6). A
    counter-0 assertion therefore lands at most one decision, ever; the counter adds nothing for such
    authenticators, as on the mint path (A2 §2.9). KD8 and RD13 test it.
- **Writer sites.** This is the daemon's second `store.OpenCredentials(` (beside A3a's, V14);
  `TestDaemonCredentialWriterSites` is edited to pin exactly these two. The writer is never held
  across anything but the two credentials transactions and one world-store write.
- **R6 of A2 recurs:** schema 8 costs one `serve` restart, named by the existing
  `OlderCredentialsStoreError` text.

### 2.8 The routes

**`GET /surface/decisions/{ref}`** (sprint 1; `approvalRoute`, so the Host allowlist applies):
```json
{"v":1,"request":"sha256:…","requestBytes":"<base64url payload>",
 "packet":"sha256:…","packetBytes":"<base64url payload>","entry":42,
 "proposal":"sha256:…","now":1794414000,"window":{"min":1794414001,"max":1794414600},
 "refusal":"","reserved":"","decidedBy":"","decision":"",
 "passkeys":[{"origin":"https://studio.tail….ts.net","rp":"studio.tail….ts.net","allow":["…"]}],
 "ssh":{"namespace":"world-decision@sunholo"}}
```
- It sends **bytes, never a precomputed message or hash.** The client checks `SHA-256(requestBytes)
  == request`, `SHA-256(packetBytes) == packet`, the packet's `requestRef == request`, and builds the
  message itself, as A3's card does (A3 §2.7).
- `refusal` is the law's D2–D7 reason at `now` (display; the POST re-runs the law); `reserved` is
  L1's owner; `window.max` is `min(now + 600, deadlineAt)`.

**`POST /surface/decide/{ref}`** (typed arm sprint 1; key arms sprint 2): A3a's guards verbatim
(R2 Host via `approvalRoute`; R3 `Origin` ∈ `O`; R4 `Sec-Fetch-Site`; R5 JSON media type; R6 ≤ 20
KiB strict JSON) plus **R0: any `Authorization` or `Cookie` header → 403** (the finish route's
posture: the agent holds a session, so a session-bearing request is one an agent could make).
Body, strict, exactly one of `sshsig`, `webauthn`, `typed`:
```json
{"v":1,"verb":"approve","packet":"<64 hex P>","request":"<64 hex Q>","expiry":1794414540,"sshsig":"-----BEGIN SSH SIGNATURE-----…"}
{"v":1,"verb":"deny","packet":"…","request":"…","expiry":…,"typed":true}   + header X-World-Operator-Proof
```
Then, under `humanMu`, the law (§2.5) and the arm; success → 200 `{"v":1,"decided":"approve",
"request":"sha256:…","decision":"sha256:…","decidedBy":"ssh:mark"}`.
- **The CLI uses the same route.** It sends `Origin: http://127.0.0.1:<port>` (in `O` by
  construction) and no `Sec-Fetch-Site`. The browser guards exist against cross-site pages, not
  against local processes; a local process still needs a key or a typed proof.
- **`OPTIONS` and other methods** → 405 (mux patterns); no CORS header is ever written.
- **Placement.** `host/daemon/surface_decide.go` (handlers), `host/surfacepage/decide_model.go`
  (JSON shapes and the body grammar, pure).

**`GET /surface/requests.json`** (sprint 3) gains `"decisions":[…]`: for each pending request of
`RecentApprovals` (≤ 20 summaries of ≤ 40 heads), its validated packet (one scan of ≤ 101 packets
per read, `packetScanLimit`), as `{request, packet, effect, scope, cost, requester, deadline,
reserved, refusal}`. A pending request with zero packets is omitted (row-176 residual, §2.4); with
two, it is listed with its refusal. `"decisionsTruncated": true` when either bound was hit, and
*Needs you* says so. `pending`/`decided` (mint requests) are unchanged. **[r1]** This listing is
display only: a request it shows with one packet may still be refused by the POST's complete
discovery (D4a/D4b, §2.4), and the card shows that refusal.

**`GET /surface/decide/{ref}` + `GET /surface/decide.js`** (sprint 3): the decision card, static
embedded bytes under the unchanged `surfaceCSP`, `approvalRoute`. The script:
1. fetches `/surface/decisions/{ref}` (10 s AbortController deadline; A3 §2.7 [r2] rules);
2. checks both hashes and the packet's `requestRef`; **fails closed** on any mismatch;
3. renders **every field of both signed objects** by one table (`DECISION_FIELDS`): request
   `effect`, `scope`, `cost`, `requester`, `now`; packet `packetHash`, `proposalHash`, `requestRef`,
   `createdAt`, `deadlineAt`, `escalationsRemaining`, `policy`; plus entry and the proposal link
   (`/workbench?object=<proposal>`). A key path the table does not list, a listed path missing, or
   an unexpected JSON type → refuse ("a field the page cannot show; decide at a terminal"). PD1 ties
   the table to `approvalRequestWire` and `packetWire` json tags at build time;
4. offers **Approve** and **Deny** with a passkey (`navigator.credentials.get`, `challenge =
   DecisionMessageHash`, UV required, A2's 240 s abort rule), or **Save the message** (91-byte
   `blob:` download) + the signing line `ssh-keygen -Y sign -n world-decision@sunholo -f <key>
   world-decision-<verb>-<ref16>.msg` + a paste box. Expiry = `window.max` at render, re-fetched if
   less than 60 s remain;
5. POSTs and shows `decidedBy`, or the refusal's class and reason, then re-reads.
`approve.js`, `/surface/passkey` and `passkey.js` are not changed (PG2 pins stay). `decide.js`
repeats A3's encoding helpers (≈ 40 lines) rather than editing `approve.js`.

**`app.js`** (sprint 3): *Needs you* lists `decisions` under the mint requests, one line each:
"Decide · <scope> · requested by <requester> · deadline <hh:mm>" linking to `/surface/decide/<ref>`;
a reserved one reads "publish approval — use world-publish approve at a terminal" with no link; the
opt-in notification (A3 §2.9) fires for a new decision id too. All text through `esc()`.

### 2.9 What the agent observes, and why it cannot decide

- **Observation.** The decision is an ordinary object on the approvals chain: the open read routes
  `GET /v1/registry/world/approvals/v1` and `GET /v1/objects/{ref}` serve it (`handleRegistry` is
  generic, `handlers.go:671-693`, V20), and the packet's object page shows `approvals chain:
  approve by ssh:mark` (`workbench_packet.go:89-92`). `Human.PollApproval` is **not** bound in the
  daemon (§1.3, V6); binding it is an effect-surface widening outside A4 (§0), and whoever binds it
  must first fix §1.7 (RK6). A4 does not change what an agent's session can call.
- **Why no grant reaches a decide path.** A decision is made only by `POST /surface/decide/{ref}`
  (no `/v1`, no MCP, no AG-UI, no A2A route, no effect handler calls `DecideApproval`; XG1 is the
  inventory test). That route refuses any `Authorization`/`Cookie` (R0), so a session is useless
  there; it then needs (a) a signature by a key enrolled with a typed `y` (R3), which an agent does
  not hold, or (b) an operator proof file in the 0700 directory beside the credentials store plus
  its secret, which the CLI writes only after reading `y` at `/dev/tty` (§1.11). (b) is row 236's
  fence, inherited unchanged: it holds exactly as far as an agent cannot write that directory (RK4).

### 2.10 The CLI: `ailang-worldd approve <request-ref> …`

`approve` already dispatches on `mr-` ids (§1.13, `approve.go:197-201`). A4 adds a second id
shape: `sha256:` + 64 hex selects the **decision mode**, in `cmd/ailang-worldd/approve_decide.go`.
- `approve <request-ref> [--deny]` — typed: GET the subject, check both hashes, print every field,
  the verb in capitals and the deadline, ask `y` at `/dev/tty` (`openTerminal`, as `session.go`),
  write the proof, POST `typed`. No terminal → refused, naming the key flags.
- `approve <request-ref> [--deny] --sign-with <key>` / `--signature <file>` / `--write-message
  <file>` / `--webauthn-link` / `--webauthn-assertion <file>` — the A1b/A2 flag names, the decision
  message instead of the mint message; `--write-message` writes the 91 bytes with the expiry it
  chose (`window.max`) and prints it, and `--signature` must be given `--expiry <t>` with it.
- `approve --list` adds a "decisions" section from `requests.json`'s `decisions` (sprint 3; sprint 1
  prints it from `GET /surface/decisions` by ref only).
- `--await` with a request ref → `exitUsage` ("decisions are made by the daemon at once; there is
  nothing to await").
- `--deny` with an `mr-` id → `exitUsage` (mints have no deny; A0's request simply expires).
- `TestEveryVerbHasHelp` gains the decision usage lines; QUICKSTART §18 (sprint 3) documents the
  full path, including building the message and signing it (S7).

## 3. Pure core, effects at the boundary (coding standards S1–S3)

| Piece | Where | Pure? | Tested by |
|---|---|---|---|
| the decision message, its hash, the namespace, the SSHSIG/assertion laws, `SignedDecision` | `host/approval/decision.go` | yes (stdlib only, A2 §3 S2) | DM1–DM8 |
| the decision law `DecisionRefusal`, `ReservedApprovalScope` | `host/broker/decision_law.go` | yes | DL1–DL12, DL4a, DL4b, R4a |
| **[r1]** the request-wide packet discovery `DiscoverRequestPackets` | `host/coordinator/decisionpacket_discover.go` | store reads only (`PacketReader`) | DS1, RD9b–RD9d |
| first-decision-final, publish L2 | `host/broker/approve.go` | store via the existing interfaces | BR1–BR3, R4b–R4d |
| the decision judge | `host/authority/judge_decision.go` | store reads only | JD1–JD3 |
| schema 8 `key_decisions`, K0's unlanded-row read | `host/store/key_decisions.go` | effect | KD1–KD8 |
| body grammar, subject JSON, `decisions` listing fold | `host/surfacepage/decide_model.go` | yes | DMD1–DMD3 |
| the routes | `host/daemon/surface_decide.go` | effect (HTTP, stores) | RD1–RD12 |
| the CLI | `cmd/ailang-worldd/approve_decide.go` | effect | CD1–CD9 |
| the card, *Needs you* lines | `host/surfacepage/decide/*`, `app/app.js` | static assets | PD1–PD2, XD1–XD4 |

- **Why no `.ail` (S1).** The law is a refusal ladder over booleans and integers that mirrors no
  sketch law; the cryptography is not encodable in pinned AILANG (A2 §3). Its invariants are
  carried by the mutation rows in §6. `verify_ail.sh` is expected exactly green, with no `.ail`
  changed (no `/v1` route, so `worlddapi.ail` is untouched, V16).
- **Why no package (S3).** Every consumer is host Go (`host/daemon`, `cmd/ailang-worldd`,
  `host/broker`) or the daemon's embedded assets; the decide path is a host authority boundary.
- **No new Go dependency.** Only the standard library and in-repo packages;
  `TestDaemonDependencyAllowlist` is a sprint gate (§9).

## 4. Threats (A4)

| # | threat | answer |
|---|---|---|
| T1 | an agent POSTs a decision | R0 refuses its session; it has no enrolled key (R3) and cannot make an operator proof (§2.9) |
| T2 | a mint signature replayed as a decision, or the reverse | different prefix, length and SSHSIG namespace (§2.3); DM3/DM4 |
| T3 | a deny replayed as an approve | the verb is signed (DM5); the typed proof's phase carries the verb (RD6) |
| T4 | the same decision replayed after it landed, or after a restart | D6 on the durable chain under `humanMu` (BR1, RD8) |
| T5 | a signature for packet A used on packet B | `P` and `Q` are in the message; D8 and `Binds` (DM6, RD7) |
| T6 | a key-approved publish stamp | L1 (D5) and L2 (§2.6); R4a–R4e |
| T7 | one decision silently answering two invocations | D4b refuses a request carried by ≠ 1 packet; **[r1]** D4a refuses when the count is not complete (DL4a, DL4b, RD9, RD9b–RD9d) |
| T15 **[r1]** | a crash or failed world write after the counter moved; a counter-0 assertion replayed | K0's exact-artifact idempotent retry; D6 after landing; the request is the row's key (KD7, KD8, KD9, RD13, RD14) **[r2]**: W15 deferred to the writer so the exact-match retry is reachable |
| T8 | deciding a request whose invocation lost its commit | D4 (zero packets) (DL4) |
| T9 | a stale signature held back and posted late | D10's 10-minute window and the packet deadline (DL9, DL10) |
| T10 | a cloned security key | the counter advances in SQL inside `RecordKeyDecision` (KD2) |
| T11 | a revoked key between the judge and the write | K3 re-checks in SQL (KD3) |
| T12 | a cross-site page posting | A3a's R2–R5 unchanged (RD2) |
| T13 | the card shows one thing and the key signs another | the page builds the message from bytes it hashed and rendered in full, fail-closed (PD1, XD2) |
| T14 | a request ref decided in one world store replayed into another with byte-identical objects | the signature then decides exactly the bytes the operator saw; accepted, named here |

## 5. By-construction argument

1. **Only a key holder or a typed operator decides.** `DecideApproval` gains exactly one production
   caller, the decide handler (XG1 pins the census at three: `RecordOperatorApproval`,
   `mintAttendedApproval`, `surface_decide.go`). The handler reaches it only after the law and one
   arm, and the `decidedBy` it passes comes from `SignedDecision.DecidedBy()` or a constant.
2. **Domain separation.** A key signature verifies only over the message whose prefix, verb, `P`, `Q`
   and expiry it was made for; the mint law builds a different prefix and length under a different
   namespace (§2.3). Neither law accepts the other's artifact (DM3/DM4).
3. **Single use.** A request on the chain has at most one decision object reachable from the head,
   because `decideApproval` refuses a second one and every writer of a decision runs under
   `humanMu` or in a process that holds the world-store writer alone (`world-publish`). The chain is
   on disk.
4. **R4.** For a key decision to stamp a publish it would need (L1) the A4 law to accept a
   publish-shaped scope, which D5 refuses for every format, and (L2) the consumer to accept a
   `decidedBy` beginning `ssh:` or `webauthn:`, which it refuses. Every key decision carries such a
   `decidedBy` (point 1). Deleting either layer leaves the other (R4c, R4d).
5. **The typed `y` stays the root.** Without any enrolled key, sprint 1's typed arm decides every
   decidable packet; enrolment and revocation stay at the terminal (R3).

## 6. Test plan

Each row names the mutation it kills (S6). Keys are `approvaltest` keys or recorded authenticator
fixtures (A2); the browser rows use the A2/A3 driver with Chrome's virtual authenticator.

### Sprint 1: the law and the typed root

| Id | Test (package) | Asserts | Mutant killed |
|---|---|---|---|
| DM1 | `TestDecisionMessageLayout` (`host/approval`) | 91 bytes; prefix, `0x00`, verb at 18, `P` at 19, `Q` at 51, BE expiry at 83; golden hex for two subjects | field order swapped; LE expiry |
| DM2 | `TestDecisionMessageRefusesBadInput` | verb 0/3 and negative expiry refused | verb not checked |
| DL1–DL10 | `TestDecisionRefusal<Dn>` (`host/broker`), one test per D1–D10 | each fact false alone → its code and reason; all true → 0 | **each Dn branch deleted** (unit = the branch) |
| DL4a **[r1]** | `TestDecisionRefusalD4a` | `Discovery{Complete: false}` with `Carrying: 1` → D4a (409), never 0; `Undecodable: 1` with `Complete: true, Carrying: 1` → D4a | **D4a deleted** (one visible match treated as exactly one); D4a reading `Carrying` instead of `Complete` |
| DL4b **[r1]** | `TestDecisionRefusalD4b` | complete with `Carrying` 0 and 2 → D4b; complete, 1, validator code 5 → D4b with the "not established" sentence | D4b `Carrying >= 1`; the validator code ignored |
| DS1 **[r1]** | `TestDiscoverRequestPacketsCompleteness` (`host/coordinator`) | fake `PacketReader` with page limit 500: 0, 499, 500, 1,999 and exactly 2,000 packets → `Complete: true` (the sentinel read returns none); 2,001 → `Complete: false`; the cursor follows `after` across pages; one undecodable payload → `Undecodable: 1` | off-by-one at the sentinel (`>=` for `>`); the cursor not advanced (the first page re-read); undecodable payloads skipped silently |
| DL11 | `TestDecisionRefusalOrder` | two facts false → the lower Dn wins | order swapped |
| DL12 | `TestDecisionRefusalMustLive` | a well-formed approve and a well-formed deny → 0 (control) | a law that refuses everything |
| R4a | `TestReservedApprovalScope` | `PublishApprovalScopeFor(…)` → publish; `registry:x`, `…effect=Registry.Publish…` near-misses → publish; `FormatFinishApprovalScope(…)`, `github-pr:x` → finish; `worktree`, `world` → not reserved (control) | prefix arm deleted; parse arm deleted |
| BR1 | `TestDecideApprovalRefusesASecondDecision` | decide deny then approve → `ErrApprovalAlreadyDecided`, head unchanged, poll still deny (the V9 probe as a test) | the new check deleted |
| BR2 | `TestDecideApprovalFirstDecisionMustLive` | an undecided request decides (control); `RecordOperatorApproval` and `mintAttendedApproval` tests unedited | refusing every decide |
| R4b | `TestPublishApprovalRefusesAKeyDecision` | a publish-shaped request decided `ssh:mark` / `webauthn:mark` → `ErrPublishApprovalDenied` with the R4 text | L2 deleted |
| R4c | `TestPublishApprovalTypedStampMustLive` | the same request decided `operator:terminal` / a free-text name → accepted (control); every `registry_publish_test.go` arm unedited | L2 matching every `decidedBy` |
| RD1 | `TestDecisionRoutesRegistered` (`host/daemon`) | `GET /surface/decisions/{ref}`, `POST /surface/decide/{ref}` exist; `/v1` still 13; `TestSurfaceAppRoutes` edited | a route missing |
| RD2 | `TestDecidePOSTGuards` | R0 `Authorization`, R0 `Cookie`, R2 foreign Host, R3 no/`null`/other Origin, R4 `same-site`, R5 `text/plain`, R6 unknown field / two formats / `decidedBy` field → refused, **0 decisions** | each guard deleted |
| RD3 | `TestDecisionSubjectJSON` | bytes decode to exactly the stored payloads; no message or hash field; `refusal`/`reserved` set by the law | a precomputed hash sent |
| RD4 | `TestDecidePOSTLawArms` | a daemon store built per D2–D10 → status, class and **0 new decisions** | the handler skipping the law |
| RD5 | `TestDecideTypedLands` | a test proof (written as `writeOperatorProof` does) → one decision `operator:terminal`, the proof file unlinked before the decision | proof not consumed |
| RD6 | `TestDecideTypedProofBindsTheVerb` | a `decide-deny` proof with `verb: approve` → 403, 0 decisions | phase not carrying the verb |
| RD7 | `TestDecidePOSTStaleDigests` | `P` or `Q` from another packet → 409 | D8 dropped |
| RD8 | `TestDecideReplayAfterRestart` | decide, close and reopen the daemon on the same store, POST again → 409 D6 | in-memory single use |
| RD9 | `TestDecideRefusesASharedRequest` | two invocations producing one request ref and two packets → 409 D4b | D4b counting only ≥ 1 |
| RD9b **[r1]** | `TestDecideDiscoveryFindsAHiddenSecondPacket` (`host/daemon`) | a store with **150 packets**, two carrying the target request, the second chosen so its `hash_ref` sorts **after the first 101** (the test asserts its position in `ObjectsBySemanticID` order before posting); a valid typed POST for the first → **409 D4b**, **0 new decisions, head unchanged**; sprint 2 re-runs it with a valid SSHSIG and also asserts **0 `key_decisions` rows** | discovery reusing `ValidatePacket`'s 100-packet prefix or `UniqueTruncated`-blind code 0 (gpt6-1-sol's scenario) |
| RD9c **[r1]** | `TestDecideDiscoveryRefusesAnOverflow` | **2,001 packets**, exactly one carrying the request → **409 D4a**, 0 decisions (sprint 2 adds the SSHSIG arm and 0 `key_decisions` rows) | the overflow sentinel ignored |
| RD9d **[r1]** | `TestDecideDiscoveryFindsTheOnlyPacketBeyondThePrefix` (must-live) | 150 packets, exactly one carrying the request, sorted after position 101 → the typed POST lands one decision | a discovery that refuses everything past the prefix |
| R4d | `TestDecideRefusesReservedScopesTyped` | a publish-shaped and a finish-shaped request, typed → 403 with the owner; **must-die control: with L1 compiled out, R4b's L2 still refuses the resulting key decision** (sprint 2 extends it with a key) | L1 bypassed for typed |
| CD1 | `TestApproveDecideTyped` (`cmd/ailang-worldd`) | a fake terminal answering `y` → POST `typed` with the verb; `n` → no POST; every field printed | the prompt skipped |
| CD2 | `TestApproveDecideNoTerminal` | no terminal → `exitUsage` naming `--sign-with` | a silent typed default |
| CD3 | `TestApproveDecideFlagRules` | `--deny` with `mr-`; `--await` with a request ref; two modes | F1 not widened |
| CD4 | `TestApproveDecideRefusesAMismatchedSubject` | a daemon answering bytes whose hash ≠ ref → refused before any prompt | the client trusting the daemon |
| XG1 | `TestDecideApprovalCallerCensus` | non-test callers of `decideApproval`/`DecideApproval` are exactly the three named in §5.1 | a fourth caller (an effect handler, MCP) |

### Sprint 2: the keys

| Id | Test | Asserts | Mutant killed |
|---|---|---|---|
| DM3 | `TestMintSignatureIsNotADecision` (`host/approval`) | an `approvaltest` SSHSIG and an assertion over the **mint** message for the same digests → refused by both decision laws | shared namespace; shared prefix |
| DM4 | `TestDecisionSignatureIsNotAMint` | a decision SSHSIG/assertion → refused by `ApproveBySignature`/`ApproveByAssertion` | — |
| DM5 | `TestDecisionBindsTheVerb` | an approve signature posted as deny → refused | verb not in the message |
| DM6 | `TestSignedDecisionBinds` | `Binds` on each field | B-checks deleted |
| DM7 | `TestDecideByAssertionSteps` | W6–W15 each refused for the decision challenge; the A2 matrix still green unedited | the refactor changing `ApproveByAssertion` |
| DM8 | `TestDecideBySignatureApprover` | P2–P5 each | `checkApprover` skipped |
| JD1–JD3 | `TestJudgeDecision*` (`host/authority`) | the step names and texts agree with `judge.go`'s for unknown key, stored key, parse, law | a judge skipping the lookup |
| KD1 | `TestCredentialsSchema8KeyDecisions` (`host/store`) | v7 → v8 upgrade; table set; DDL equals `credentials_schema.sql`; v7 refuses with the restart text | upgrade step deleted |
| KD2 | `TestRecordKeyDecisionCounter` | a non-advancing counter → `ErrSignCountRegression`, 0 rows; advancing → updated | the SQL counter check deleted |
| KD3 | `TestRecordKeyDecisionRevokedSince` | revoke between judge and record → refused | K3 deleted |
| KD4 | `TestRecordKeyDecisionReplaceUnlanded` | an unlanded row is replaced; a landed one is not | K5 replacing landed rows |
| KD5 **[r1]** | `TestKeyDecisionsPacketRefCheck` | raw `INSERT`s with `packet_ref` of length 71 not starting `sha256:` (e.g. `sha512:` + 64 hex) and of length 70 → constraint error; a `sha256:` + 64 hex value → accepted (control) | the `GLOB 'sha256:*'` clause dropped from `packet_ref` |
| KD6 **[r1]** | `TestKeyDecisionsDecisionRefCheck` | `decision_ref` NULL → accepted; `sha256:` + 64 hex → accepted; length 71 without `sha256:`, length 70, empty string → constraint error | the `decision_ref` CHECK dropped, or `IS NULL OR` dropped (NULL refused) |
| KD7 **[r1] [r2]** | `TestRecordKeyDecisionIdempotentRetry` | an unlanded row from a hardware-counter assertion (counter 7, stored now 7); the same artifact retried → K0 match → **no `usePasskey` call, counter still 7**, `DecideApproval` + `LandKeyDecision` run once; a **second, independently made** valid assertion that also carries counter 7 → `ErrSignCountRegression`, 0 decisions | K0 deleted (the retry dies on the counter: the objection's failure); K0 matching on request ref alone, without the byte-equal artifact (lets a different artifact skip the counter) |
| KD8 **[r1]** | `TestDecidePOSTCounterZeroAssertionLandsOnce` (`host/daemon`) | a counter-0 passkey (stored 0, assertion 0): first POST lands; the identical POST again → **409 D6**, still one decision, one row; a crash simulated between record and decide, then the identical POST twice → one decision | D6 deleted (the counter-0 assertion lands twice); K0 inserting a second row instead of landing the first |
| KD9 **[r2]** | `TestRecordKeyDecisionExactMatchRechecksTheApprover` (`host/store`) | an unlanded row from counter 7; the approver revoked **after** the read-only judge and **before** `RecordKeyDecision` (store level: the judge is not in this test, so nothing else can refuse) → the exact-match call refuses with the revoked-approver error, counter still 7, row still unlanded, 0 decisions | the exact-match path skipping `useApprover` (gemini-4-argon's variant: a revoked approver's retry lands) |
| RD10 | `TestDecidePOSTSSHSIG` / `TestDecidePOSTPasskey` | → 200, one decision `ssh:<name>` / `webauthn:<name>`, one landed `key_decisions` row | recording before the law |
| RD11 | `TestDecidePOSTKeyDecidedByIsDerived` | a body with a forged `decidedBy` → 400; the landed `decidedBy` equals the verified approver | caller-supplied `decidedBy` |
| RD13 **[r1] [r2]** | `TestDecidePOSTRetryAfterAFailedWorldWrite` (`host/daemon`) | a fault-injected `DecideApproval` failure (the store's write seam) after `RecordKeyDecision`, for (a) an SSHSIG, (b) a counter-N passkey, (c) a counter-0 passkey → 500, one unlanded row, 0 decisions; the **identical** body retried → 200, one decision, the row landed, the counter unchanged by the retry; a fresh assertion at `N+1` instead → 200 by replacement | the retry path running `usePasskey` again; the unlanded row left unlanded |
| RD14 **[r2]** | `TestDecidePOSTRetryAfterAFailedWorldWriteAcrossRestart` (`host/daemon`) | gpt6-1-sol's test, verbatim: "commits stored counter 7, injects a world-write failure, restarts the daemon, and retries the identical assertion successfully without moving the counter; a different assertion carrying counter 7 must still fail"; and gemini-4-argon's: "the full HTTP/judge path succeeds on a counter-7 retry when `webauthn_credentials.sign_count` is already 7 in SQLite". | the read-only judge running W15 before the exact-match check (the round-2 defect: the retry dies at W15); exact match on the artifact alone without the signed subject fields |
| RD12 | `TestDecidePOSTBusyWriter` | a held credentials writer → 503 `Busy`, **0 decisions** | deciding before the evidence |
| R4e | `TestDecideRefusesReservedScopesByKey` | publish- and finish-shaped requests with a valid key → 403; with L1 compiled out → the decision lands but `validatePublishApproval` refuses it (must-die) | either layer alone deleted |
| CD5–CD9 | `TestApproveDecide{SignWith,Signature,WriteMessage,WebAuthnLink,WebAuthnAssertion}` | the 91-byte message on disk; `ssh-keygen -Y sign -n world-decision@sunholo` (an `approvaltest` key file); the POST; refusals print the judge's text | the mint message written |
| edited | `TestDaemonCredentialWriterSites` (two sites), `TestCredentialsMethodsRequireADeadline`, `TestProductionContextRoots`, `TestEveryVerbHasHelp` | — | — |

### Sprint 3: the surface (live Chrome gated as XB1, as A3)

| Id | Test | Asserts | Mutant killed |
|---|---|---|---|
| DMD1–DMD3 | `TestDecideModel*` (`host/surfacepage`) | the `decisions` fold: pending only, zero-packet requests omitted, two-packet listed refused, reserved marked, truncation flags | omitting the reserved marker |
| PD1 | `TestDecidePageFieldTable` | `DECISION_FIELDS` equals the json tags of `approvalRequestWire` and `packetWire` | a field added without the page |
| PD2 | `TestDecidePageAssets` | `GET /surface/decide/{ref}` and `decide.js` under the unchanged `surfaceCSP`; A2/A3 page hashes unchanged (PG2) | CSP widened |
| XD1 | `TestDecidePageApproveWithPasskey` | virtual authenticator, enrolled test passkey, Approve → one decision `webauthn:<name>`; Deny on a second packet → `deny` | the page signing the mint challenge |
| XD2 | `TestDecidePageFailsClosed` | an injected unknown key, a hash mismatch, a wrong `requestRef` → refused, no `get` call | rendering a subset |
| XD3 | `TestNeedsYouListsDecisions` | the line, the link, the reserved text without a link | — |
| XD4 | `TestDecidePageSSHSIGPaste` | the saved 91 bytes equal Go's `DecisionMessage`; a pasted `approvaltest` signature decides | the blob not the message |

**Must stay green, unedited:** every `TestApprove*`, `TestPasskey*`, `TestApprovePage*`,
`TestSurfaceApprove*`, `registry_publish_test.go`, `episode_finish*_test.go`,
`decisionpacket*_test.go`, `TestFrozenV1RouteTableMatchesMux`, `TestDaemonDependencyAllowlist`.

## 7. Repo-wide inventories

- `TestFrozenV1RouteTableMatchesMux`: unchanged, 13 (no `/v1` route).
- `TestSurfaceAppRoutes`: +2 (sprint 1), +2 (sprint 3).
- `TestDaemonCredentialWriterSites`: 1 → 2 (sprint 2).
- `TestCredentialsMethodsRequireADeadline`: +`RecordKeyDecision`, `LandKeyDecision`, `KeyDecision`.
- `TestProductionContextRoots`: +1 per new CLI root (`approve_decide.go`).
- `TestDaemonDependencyAllowlist`: unchanged.
- **XG1 (new)**: the `decideApproval` caller census (§5.1).

## 8. Conflict surface

| file | who else edits it | risk |
|---|---|---|
| `host/broker/approve.go` | row 172/236 (landed) | two small additions (BR, L2); no wire change |
| `host/daemon/daemon.go` | every daemon row | +2/+4 `HandleFunc` lines only |
| `host/surfacepage/app/app.js` | row 188 style work (1,470 lines) | one new block in `needsYouHTML`; sprint 3 only |
| `cmd/ailang-worldd/approve.go` | A3a (landed) | the id-shape dispatch; the decision mode lives in a new file |
| `host/store/credentials.go` | any credentials row | schema 8 step, A2/A3 pattern |
| `docs/QUICKSTART.md` | every row | new §18 |

## 9. Milestones, sprints and gates

Every milestone ends with `go build ./... && go vet ./<pkg>/` (the compile fence for test files)
and its own tests. Every sprint ends with `./scripts/verify_ail.sh`;
`AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_go.sh`; `go build ./... && go test ./...`;
`go test ./host/daemon/ -run '^(TestFrozenV1RouteTableMatchesMux|TestSurfaceAppRoutes|TestDaemonDependencyAllowlist|TestDaemonCredentialWriterSites)$' -count=1 -v`
(each must print `--- PASS`).

**Sprint 1 = M1–M3, "the law and the typed root", about 1,950 lines [r1].**
- **M1 the message, the law and the discovery** (≈ 280 + 700 **[r1]**): DM1, DM2, DL1–DL12, DL4a, DL4b, DS1, R4a.
  `go test ./host/approval/ ./host/broker/ ./host/coordinator/ -run '^(TestDecisionMessage|TestDecisionRefusal|TestReservedApprovalScope|TestDiscoverRequestPackets)' -count=1 -v`
- **M2 the broker rules** (≈ 40 + 260): BR1, BR2, R4b, R4c, XG1.
  `go test ./host/broker/ -run '^(TestDecideApproval|TestPublishApproval|TestDecideApprovalCallerCensus)' -count=1 -v`;
  `go test ./host/broker/ -count=1` (the publish arms unedited).
- **M3 the subject, the typed POST, the typed CLI** (≈ 260 + 590 **[r1]**): RD1–RD9, RD9b–RD9d, R4d, CD1–CD4.
  `go test ./host/daemon/ -run '^(TestDecisionRoutesRegistered|TestDecidePOST|TestDecisionSubjectJSON|TestDecideTyped|TestDecideReplay|TestDecideRefuses|TestDecideDiscovery)' -count=1 -v`;
  `go test ./cmd/ailang-worldd/ -run '^(TestApproveDecide|TestEveryVerbHasHelp)' -count=1 -v`
- **Fallback cut:** CD1–CD4 (the typed CLI) move to sprint 2's head (**[r1]** this brings sprint 1 back to ≈ 1,550); the discovery, D4a and RD9b–RD9d are never cut; the route and its tests
  stay, and §10's typed step waits for them. M1 and M2 are never cut.

**Sprint 2 = M4–M5, "the keys", about 1,800 lines [r1].**
- **M4 the laws, the judge, schema 8** (≈ 350 + 700 **[r1]**): DM3–DM8, JD1–JD3, KD1–KD7.
  `go test ./host/approval/ ./host/authority/ ./host/store/ -run 'Decision|Decide|JudgeDecision|KeyDecision|Credentials' -count=1 -v`;
  the A2 matrix `go test ./host/approval/ -count=1` unedited.
- **M5 the key POST arms and the key CLI** (≈ 230 + 520 **[r1]**, + ≈ 60 **[r2]**): RD10–RD14, KD8, R4e, CD5–CD9, edited inventories.
  `go test ./host/daemon/ ./cmd/ailang-worldd/ -run '^(TestDecidePOST|TestDecideRefusesReservedScopesByKey|TestApproveDecide|TestDaemonCredentialWriterSites)' -count=1 -v`
- **Fallback cut:** the webauthn arm (RD10 passkey half, CD8–CD9, DM7) moves to sprint 3, which
  needs it for XD1 anyway; SSHSIG ships alone. R4e is never cut. **[r1]** K0, KD7, KD8 and RD13 move with the webauthn arm only in their passkey halves; RD13 (a) (SSHSIG) and the DDL CHECKs (KD5, KD6) are never cut.

**Sprint 3 = M6–M7, "the surface", about 1,500 lines.** **Planner obligation P1** (§1.12) runs
first: measure A4's 91-byte message with `crypto.subtle` in Chrome against Go's goldens, as A3 V23.
- **M6 the listing and the card** (≈ 120 + 300 page + 60 Go + 400 tests): DMD1–DMD3, PD1, PD2.
- **M7 *Needs you*, the live browser, QUICKSTART §18** (≈ 80 `app.js` + 450 tests + 90 docs): XD1–XD4.
  `WORLD_CHROME="/Applications/Google Chrome.app/Contents/MacOS/Google Chrome" WORLD_REQUIRE_PASSKEY_BROWSER=1 go test ./cmd/ailang-worldd/ -run '^(TestDecidePage|TestNeedsYouListsDecisions|TestApprovePage|TestPasskeyPage)' -count=1 -v`
  must print `--- PASS` for every subtest on the controller's machine (CI has no Chrome step, A3 V13).
- **Fallback cut:** the card page (M6 page half, XD1, XD2, XD4) becomes a fourth sprint; the
  `decisions` listing, the *Needs you* line (pointing at `ailang-worldd approve <ref>`) and QUICKSTART
  stay. PD1 is never separated from the page.

**Priced whole-A4 fallback.** If only sprint 1 lands, the typed `y` decides every decidable
packet and no key can (R1 is then unused for decisions, not violated). If sprints 1–2 land, keys
decide from any terminal; the surface lists nothing new.

## 10. The attended exit run (Mark; NOT run by the loop, D-WORLD-109)

QUICKSTART §18 *(attended — pending first verbatim run)*. Mark, on the studio, with a real episode:
1. Start an episode whose grants include `--grant Human.Approve=worktree:1`; run a skill whose plan
   requests one `Human.Approve`; confirm the packet page reads `approvals chain: pending`.
2. **Typed root:** `ailang-worldd approve <request-ref> --deny` at a terminal; read every printed
   field; type `y`; confirm `deny by operator:terminal` on the packet page.
3. **SSH key:** a second packet; `approve <ref> --sign-with <enrolled key>` (Secretive or a YubiKey
   `sk` key); confirm the holder's prompt and `approve by ssh:<name>`.
4. **Passkey:** a third packet; open *Needs you* on the phone over the tailnet, open the card, read
   it, Approve with Face ID; confirm `approve by webauthn:<name>`; record which `Host` arrived (A3
   §10.2's open question).
5. **Replays and R4:** re-POST step 3's body (saved by `--write-message`/`--signature`) → refused
   "already decided"; confirm a publish-shaped request (a test grant on a throwaway store) shows
   "use world-publish approve" and no decide control.
6. Record the run in the iteration log and §15; any refusal text that confused Mark becomes a row.

## 11. Residuals and risks

- **RK1** The derived challenge (§0 (a)) gives no server-side "the daemon asked for this" fact;
  freshness is the 10-minute signer-chosen window. D-WORLD-115 B restores the parent's literal.
- **RK2** Identical requests share a ref (§1.5); A4 refuses them (D4) rather than fixing the
  request identity, which would change the frozen request wire. The fix belongs to a row that owns
  `world/approval-request/v2`.
- **RK3** The 40-head anchor bound (§2.4): a busy world can make an old packet undecidable here; it
  then expires under `Cancel`. An indexed approvals surface (the comment at `approve.go:261-264`)
  removes it.
- **RK3b [r1]** The complete-discovery cap (§2.4): with more than 2,000 stored packets in one world, D4a refuses **every** A4 decision; the packets then expire under `Cancel`, and the typed owners are unaffected. A request-indexed packet lookup with a verified completeness contract (the reviewer's alternative) lifts it; it needs a store index and is not built here.
- **RK4** The typed arm inherits row 236's operator-proof fence: it is as strong as the agent's
  inability to write the directory beside the credentials store.
- **RK5** No L2 on the finish record (§2.6).
- **RK6** `Human.PollApproval` ignores scope (§1.7); harmless while unbound in the daemon, a defect
  the day it is bound.
- **RK7** Schema 8 costs one `serve` restart (A2 R6).
- **RK8** Browser premises for the 91-byte message are unmeasured until P1 (§1.12).

## 12. Open questions

- **D-WORLD-115 (proposed; one word, A/B).** "Row 201 A4: a `Human.Approve` decision signs
  `world-decision/v1 ‖ verb ‖ packet digest ‖ request digest ‖ expiry` under namespace
  `world-decision@sunholo`, with no stored random challenge (single use from the approvals chain,
  freshness from a ≤ 10-minute expiry). **A** (default): keep it. **B**: add the parent's stored
  32-byte daemon challenge (≈ +550 lines in sprint 2; the verb, prefix and namespace stay)." Until
  answered, sprints 1–2 build A; nothing in sprint 1 depends on the answer.
- **Q2 (planner, not Mark).** Whether `approve --list` should fetch decisions from the daemon in
  sprint 1 or wait for sprint 3's `requests.json`; the design allows either.

## 13. What this is NOT doing

- No escalation, timeout transition, or policy other than `Cancel`.
- No binding of `Human.PollApproval`; no change to what a session can call.
- No change to the mint path, the A0–A3a law, `surface_approvals`, `approve.js`, or the passkey page.
- No `/v1` route, no `.ail` change, no new Go dependency.
- No decision of publish- or finish-shaped requests by any A4 format.
- No Tauri code (A3b).

## 14. Verification Log

Commands ran in the worktree at `b0fbcd3e` unless stated. Empty or negative results are paired with
a positive control in the same command.

| # | Claim | Command | Observed |
|---|---|---|---|
| V1 | Base, branch, toolchain, file hashes | `git rev-parse HEAD`; `git branch --show-current`; `go version`; `shasum -a 256 host/broker/approve.go host/broker/approvals_view.go host/coordinator/decisionpacket.go host/coordinator/decisionpacket_validate.go host/daemon/daemon.go host/daemon/episode_finish.go host/daemon/surface_approve.go host/authority/judge.go host/approval/request.go host/approval/sshsig.go cmd/ailang-worldd/approve.go cmd/ailang-worldd/episode_finish_operator.go host/surfacepage/app/app.js \| awk '{print substr($1,1,16), $2}'` | `b0fbcd3eee7b5b29190e678e296635a3dd3a0d0c`; `design/world-iter388-row201-a4`; `go1.26.6 darwin/arm64`; the prefixes in the header |
| V2 | The parent's A4 line, §2.2 and §6 | `cat design_docs/planned/w-world-approvals-by-key.md` (93 lines, read in full) | "**A4:** `Human.Approve` decision packets (row 172) approved the same way."; §2.2 "A **mint request**, and in S2 any `Human.Approve` decision packet … carries: … a **one-time challenge**: 32 random bytes from the daemon; an **expiry**: 10 minutes by default. The signed message is a fixed-prefix encoding of `("world-approval/v1", request digest, challenge, expiry)`"; "the namespace is `world-approval@sunholo`"; R1 "for **mints and `Human.Approve`**"; R4 "**publishing stays typed-only**" |
| V3 | `DecideApproval` has no production caller (control: test callers) | `grep -rn "DecideApproval(\|decideApproval(" --include="*.go" . \| grep -v _test.go`; `grep -rln "DecideApproval(" --include="*_test.go" . \| wc -l` | `host/broker/approve.go:181` (def), `:188`, `:191` (def), `:696`; `host/broker/publish_op.go:153`; control `5` |
| V4 | The request wire, strict decode, object hash | `cat -n host/broker/approve.go \| sed -n 1,340p`; `grep -n "func brokerObject" -A8 host/broker/*.go` | `:66-72 approvalRequestWire {effect, scope, cost, requester, now}`; `:141-145` built from `req`; `:320-327 dec.DisallowUnknownFields()`; `record.go:132-140 Hash: hashref.SumSHA256(payload)` |
| V5 | The two production deciders and `tty.go` | `grep -rn "RecordOperatorApproval(\|MintAttendedApproval(" --include="*.go" . \| grep -v _test.go`; `grep -n "DecidedBy\|DecideApproval\|headless" cmd/world-publish/main.go cmd/world-publish/tty.go`; `sed -n 120,200p host/broker/publish_op.go` | `cmd/world-publish/main.go:310 broker.MintAttendedApproval(activeCtx, db, plan)`; `host/daemon/episode_finish.go:173 … "operator:episode-finish", "operator:terminal", now)`; defs `approve.go:686`, `publish_op.go:125`; `tty.go:15-16` "a headless loop can mint one by calling broker.DecideApproval. It is a SINGLE-USE fence." |
| V6 | The daemon binds `Human.Approve` and not `Human.PollApproval` (control: Approve found by the same grep) | `grep -rn "EffectHumanPollApproval\|EffectHumanApprove" --include='*.go' host/daemon \| grep -v _test.go`; `sed -n 940,968p host/daemon/daemon.go` | only `host/daemon/daemon.go:954: reg[broker.EffectHumanApprove] = d.human // row 172 D4 (MD2 target)`; `serializedHandler` `:959-968` locks `mu`; PollApproval: no match |
| V7 | The packet wire, its producer and freeze | `cat -n host/coordinator/decisionpacket.go` | `:86-88` "The byte form is FROZEN … Amendments are /v2 semantic IDs"; `:98-106 packetWire`; `:174-181 mintPacket` `DeadlineAt: grantExpiresAt`, `Policy: "Cancel"`; `effectful.go:340-352` producer for one OK `Human.Approve` |
| V8 | The validator's result and bounds; listing order | `sed -n 1,60p host/coordinator/decisionpacket_validate.go`; `grep -rn "ORDER BY hash" host/store/*.go`; `sed -n 1603p docs/QUICKSTART.md` | `PacketResult{Code, Entry, Packet, TransitionFn, Approval, AnchorTruncated, MalformedRef, UniqueTruncated}`; `packetScanLimit = 100`; `objects_by_semantic_id.go:20 … ORDER BY hash_ref LIMIT ?`; QUICKSTART "shows `approvals chain: pending` until an operator decides … The anchor check reads the newest 40 approval heads" |
| V9 | **Identical requests collide; a request is decided twice and the newest wins; poll ignores scope** (throwaway probe) | `/Users/voightkampff/.ailang/state/world-iter388/designer-scratch/redecide/main.go` (sha `d9162ffe5e808aa8`; `go.mod` `replace github.com/sunholo-data/ailang-world => <this worktree>`): temp `store.Open`, `NewHumanHandler`, two `Execute(Human.Approve, scope worktree, cost 1, now 100, requester agent)`, `DecideApproval(deny, operator:terminal, 101)`, `DecideApproval(approve, ssh:mallory, 102)`, `Execute(Human.PollApproval, scope some-other-scope)`, `RecentApprovals(40)`; run `GOTOOLCHAIN=local GOFLAGS=-mod=mod GOPROXY=off ~/go/pkg/mod/golang.org/toolchain@v0.0.1-go1.26.6.darwin-arm64/bin/go run .` | `execute1 …"requestRef":"sha256:db66eb2f…3ff9"}`; `execute2` the same; `same=true`; `decide deny: sha256:c48483a5… err=<nil>`; `decide approve: sha256:15336edf… err=<nil>`; `poll(other scope) {"status":"decided","decision":{…"decision":"approve","decidedBy":"ssh:mallory","now":102}} err=<nil>`; `view sha256:db66eb2f4cd08 status=approve by=ssh:mallory` |
| V10 | The mint message, namespace, verifier and laws | `sed -n 200,232p host/approval/request.go`; `grep -n "MessagePrefix *=\|Namespace *=" host/approval/*.go \| grep -v _test`; `sed -n 340,380p host/approval/sshsig.go`; `cat -n host/approval/signed.go \| sed -n 1,140p`; `sed -n 28,95p host/approval/passkey_assert.go`; `cat -n host/authority/judge.go` | `request.go:26 MessagePrefix = "world-approval/v1"`; `SignedMessage` 90 bytes, no verb; `sshsig.go:19 Namespace = "world-approval@sunholo"`; `:262 VerifySSHSIG(s, namespace, message, key)`; `:353 VerifyApproval` fixes `Namespace`; `passkey_assert.go` `MessageHash(in.Admit.StoredDigest, in.Challenge, in.Admit.ExpiresAt)` W6, W7, W8, W10, W15 `ad.SignCount <= p.SignCount`; `signed.go:75 checkApprover` P2–P5, `:113 ApprovedBy` = `format + ":" + name`; `judge.go:65`, `:102` take `store.MintRequestRow` |
| V11 | **Publish consumption ignores `decidedBy`; both world-publish steps are attended; the writer is exclusive** | `cat -n host/broker/approve.go \| sed -n 517,650p`; `grep -n "PublishCost *=" host/broker/*.go`; `sed -n 280,360p cmd/world-publish/main.go`; `grep -n "decidedBy\|approval-ref\|requireAttendedOperator" cmd/world-publish/main.go` | `:544 decision.Decision != "approve"`, `:578 request.Effect != EffectHumanApprove`, `:583 request.Cost != PublishCost`, `:589 ParsePublishApprovalScope`; no `DecidedBy` read in `:517-650`; `registry_publish.go:31 const PublishCost = 1`; `main.go:172 "decided-by" … "who is granting the approval"`; `:294`, `:399 requireAttendedOperator`; `:298-305 store.Open` + `IsWriterAlreadyActive`; `:347-351 requireApprovalRef(opts.approvalRef)` |
| V12 | The finish record ignores requester and `decidedBy` | `grep -n "func (d \*Daemon) finishDecisionScope" -A40 host/daemon/*.go` | `episode_finish.go:243-280`: decodes `{requestRef, decision}`, `decision.Decision != "approve"`, `ParseFinishApprovalScope(request.Scope)`, expiry bound; no `Requester`/`DecidedBy` |
| V13 | The operator-proof fence | `cat -n host/daemon/episode_finish.go \| sed -n 60,200p`; `cat -n cmd/ailang-worldd/episode_finish_operator.go \| sed -n 1,110p`; `grep -rn "consumeOperatorProof\|X-World-Operator-Proof\|operator-proofs\|humanMu" --include='*.go' . \| grep -v _test.go` | `:90-133 consumeOperatorProof` (0700 dir, owned 0600 file, `p.ExpiresAt > now+300` refused, `ConstantTimeCompare`, `os.Remove` before mutation); `:139-143` Authorization refused; `:158-160 humanMu.Lock(); consumeOperatorProof`; CLI `:30-73 writeOperatorProof`, `:84-85` header; `daemon.go:438-440` "humanMu serializes Human.Approve across episodes … without compare-and-set" |
| V14 | A3a's POST and the daemon's one credentials writer site | `sed -n 180,296p host/daemon/surface_approve.go`; `grep -rn "OpenCredentials(" --include='*.go' host/daemon \| grep -v _test.go`; `grep -rn "func TestDaemonCredentialWriterSites" host/daemon/*_test.go` | R1, R3–R14 as cited; only `surface_approve.go:280 wr, err := store.OpenCredentials(d.credentialsPath) // R13: non-waiting, closed at once`; test at `surface_approve_post_guard_test.go:132` |
| V15 | `RequestsJSON`, *Needs you* in `app.js` | `grep -n "type RequestsJSON" -A8 host/surfacepage/approve_model.go`; `grep -n "requests.json\|Needs you" host/surfacepage/app/app.js`; `wc -l host/surfacepage/app/app.js` | `:64-70 {V, Now, Pending, Decided, Unavailable}`; `app.js:424`, `:429 Needs you ·`, `:921-926 getJSON('/surface/requests.json')`; `1470` |
| V16 | `/v1` pinned at 13, tied to the sketch; the 13th came with a ruling | `sed -n 15,60p host/daemon/route_table_test.go`; `grep -c 'path: "/v1/' design_docs/sketches/worlddapi.ail`; `grep -n 'HandleFunc' host/daemon/daemon.go` | "Thirteen since row 218 M2 added POST /v1/episodes/{episode}/finish (the operator route, D-WORLD-100 R6)"; `len(mux) != 13 \|\| len(sketch) != 13`; `13`; routes `daemon.go:1108-1149` incl. `:1137-1141` the five A3a approval routes |
| V17 | `surfaceCSP` | `grep -n "^const surfaceCSP" host/daemon/*.go` | `surface_page.go:19 … "default-src 'none'; style-src 'unsafe-inline'; script-src 'self'; connect-src 'self'; font-src 'self'; base-uri 'none'; form-action 'self'; frame-ancestors 'none'"` |
| V18 | Grants are exact-string; the documented grant | `sed -n 1,130p host/broker/decide.go`; `grep -n "Human.Approve" docs/QUICKSTART.md` | `decide.go:110 effectNameMatches … c.Effect == effect`; `:111 scopeMatches … c.Scope == scope`; `QUICKSTART.md:1654 --grant Human.Approve=worktree:1` |
| V19 | No test or doc passes a key-shaped `decidedBy` (L2's only behaviour change); the format names | `grep -rn 'decided-by[ =]\+"\?\(ssh\|webauthn\):' --include='*.go' --include='*.md' --include='*.sh' . \| grep -v design_docs/planned/w-world-approvals-by-key-a4.md \| wc -l` (control: `grep -rn 'decided-by' --include='*.go' cmd/world-publish \| wc -l`); `grep -rn 'DecidedBy: *"\(ssh\|webauthn\):' --include='*_test.go' host/broker cmd/world-publish \| wc -l` (control: `'DecidedBy: *"'`); `grep -rn 'ecideApproval(.*"\(ssh\|webauthn\):' --include='*_test.go' . \| wc -l` (control: `'ecideApproval('`); `grep -n 'FormatSSH *=\|FormatWebAuthn *=' host/approval/*.go` | `0` (control `2`); `0` (control `6`); `0` (control `13`); `approver.go:16 FormatSSH = "ssh"`, `:18 FormatWebAuthn = "webauthn"` |
| V20 | The registry read route is generic | `grep -n "func (d \*Daemon) handleRegistry" -A25 host/daemon/handlers.go` | `handlers.go:671-693`: any `{name...}` → `GetRegistryHead(ctx, name)` → `{name, head}` |
| V21 | Credentials schema 7; the mint id pattern | `grep -n "credentialsSchemaVersion = \|case version >= 7\|if version < 7" host/store/credentials.go`; `grep -n "MintRequestIDPattern\b.*=" host/store/mint_requests.go`; `sed -n 100,108p host/store/mint_requests.go` | `:43 const credentialsSchemaVersion = 7`; `:74 case version >= 7:`; `:369 if version < 7 {`; `mint_requests.go:46 var MintRequestIDPattern = regexp.MustCompile(`^mr-[0-9a-f]{16}$`)`; `:105` is a **use** (`if !MintRequestIDPattern.MatchString(row.RequestID)`), not the definition |
| V22 | Charter rows and ledger | `grep -n '^17[26]\. \|^201\. ' design_docs/world-mission.md`; `grep -n "\| D-WORLD-73 \|" …`, `-89`, `-109`, `-114`; `grep -o 'D-WORLD-1[0-9][0-9]' design_docs/world-mission.md \| sort -t- -k3 -n \| uniq \| tail -3` | `:6240` row 172 LANDED `d36b148`; `:6247` row 176 LANDED `8511bad` "`d.binder` binds `Human.Approve` … `serializedHandler`"; `:6284` row 201 IN PROGRESS, A3a sprint 4 LANDED `4134ce89`; `:992` D-WORLD-73 RESOLVED; `:1008` D-WORLD-89 R1 "for mints and Human.Approve decisions"; `:1028` D-WORLD-109 "The loop builds and judges against fakes and test keys, and never enrols a real key, mints or publishes."; `:1033` D-WORLD-114 OPEN; highest id `D-WORLD-114` |
| V23 | The approve CLI dispatches on `mr-` ids and its modes | `sed -n 130,200p cmd/ailang-worldd/approve.go` | F1 over `--signature`, `--sign-with`, `--write-message`, `--webauthn-link`, `--webauthn-assertion`, `--await`; `:197-201 if !store.MintRequestIDPattern.MatchString(id) { return refuse("%q is not a mint request id …` |
| V24 | The packet page shows the chain's status | `cat -n host/daemon/workbench_packet.go` | `:89-92 status := res.Approval.Status; if res.Approval.DecidedBy != "" { status += " by " + res.Approval.DecidedBy }` |
| V25 | The coordinator's clock is Unix seconds | `grep -n "Now:\|\.Now()" host/daemon/daemon.go` | `:876 Binder: d.binder, Now: func() int64 { return time.Now().Unix() }`; `:829-830 finishNow`, `approveNow` wall clock |
| V26 | `approvalRoute` applies the Host allowlist | `grep -n "func (d \*Daemon) approvalRoute" -A12 host/daemon/surface_approve.go` | `:76-87`: `approvalOrigins`, `allowedApprovalHost(r.Host, o)` → 421 `MisdirectedRequest` |
| V27 | The A3a guard predicates are reusable | `cat -n host/daemon/surface_approve_guard.go \| grep -n "func "` | `surfaceOrigins :21`, `checkSurfaceOrigin :32`, `allowedApprovalHost :41`, `allowedApprovalOrigin :53`, `sameOriginFetch :57`, `jsonMediaType :61` |
| V28 | The consume's in-transaction counter and approver helpers (reused by K3+K4) | `grep -n "sign_count\|SignCount" host/store/*.go \| grep -v _test`; `sed -n 215,230p host/store/passkeys.go`; `grep -n "func useApprover\|func usePasskey" host/store/*.go` | `passkeys.go:35 sign_count … CHECK (sign_count BETWEEN 0 AND 4294967295)`; `:57 ErrSignCountRegression`; `:221-230 usePasskey` "the consume's counter binding … only forward (both zero: unchanged)", `UPDATE webauthn_credentials SET sign_count = ? WHERE name = ? AND credential_id = ? AND (sign_count < ? OR (sign_count = 0 AND ? = 0))`; `approvers.go:245 func useApprover(ctx, tx, a ApproverUse, at int64)` |
| V29 | The world-demo plan requests one `Human.Approve` with a requester | `grep -n "requester" packages/world-demo/approve.ail` | `:109 kv("payload", jo([kv("requester", js("demo"))]))`; `:116` plan `{"effect":"Human.Approve","scope":"worktree","cost":1,"payload":{"requester":"demo"}}` |
| V30 | Probes stayed out of the repo | `git status --short` (worktree) after writing this doc | `?? design_docs/planned/w-world-approvals-by-key-a4.md` only (the controller re-checks) |
| V31 | A3a's measured sprint sizes (calibration) | `git log --merges --oneline -40 \| grep -i "a3a\|row201-a3"`; `for m in $(git log --merges --format=%h -40 --grep='a3a'); do git diff --numstat $m^1 $m -- . ':!design_docs' \| awk '{a+=$1+$2} END {print a}'; done` | `b44b5f7d` (s1) `1543`; `d1027477` (s2) `1995`; `70823879` (s3) `1861`; `4134ce89` (s4) `393` |
| **V32 [r1]** | **The counter-0 rule: the same counter-0 assertion passes the counter twice** | `cat -n host/store/passkeys.go \| sed -n 221,250p`; `grep -n "W15" host/approval/passkey_assert.go` | `passkeys.go:226 UPDATE webauthn_credentials SET sign_count = ? WHERE name = ? AND credential_id = ? AND (sign_count < ? OR (sign_count = 0 AND ? = 0))`; `:231 if n … == 1 { return err }` (accepted); `:244-246` C18 `ErrSignCountRegression` otherwise; `passkey_assert.go:88 if (ad.SignCount != 0 \|\| p.SignCount != 0) && ad.SignCount <= p.SignCount { // W15`. Stored 0 / assertion 0 → W15 false and the UPDATE matches → accepted, every time; stored N / assertion N → W15 refuses and the UPDATE matches no row |
| **V33 [r1]** | **`ValidatePacket`'s uniqueness scan is a 100-packet hash-ordered prefix and still returns code 0 when truncated** | `grep -n "packetScanLimit\|UniqueTruncated" host/coordinator/*.go host/daemon/*.go \| grep -v _test`; `cat -n host/coordinator/decisionpacket_validate.go \| sed -n 185,216p` | `decisionpacket_validate.go:34 const packetScanLimit = 100`; `:190 ObjectsBySemanticID(ctx, DecisionPacketV1, "", packetScanLimit+1)`; `:194 res.UniqueTruncated = len(others) > packetScanLimit`; `:195 for _, o := range others[:min(len(others), packetScanLimit)]`; `:214 f[6] = true` → code 0 with `UniqueTruncated` set; its only reader `workbench_packet.go:78` appends a display sentence |
| **V34 [r1]** | **The discovery's page contract: keyset by `hash_ref`, unique total order, payload nil, page ≤ 500** | `cat -n host/store/objects_by_semantic_id.go \| sed -n 13,45p` | `:13 const MaxSemanticIDPage = 500`; `:15-18` "a keyset page … it never depends on OFFSET or SQLite rowids: hash_ref is the PRIMARY KEY, so it is a unique total order"; `:20 … WHERE semantic_id = ? AND hash_ref > ? ORDER BY hash_ref LIMIT ?`; `:24-25` "Payload is always nil: callers fetch the one payload they want with GetObject"; `:34` a limit outside `1..500` → `InvalidLimitError` |

## 15. Quorum verification log

- **Round 1 (2026-11-12T12:44:28Z): BLOCKED 2/2 present (gemini-4-argon, gpt6-1-sol); oc-glm-5-3 and oc-kimi-k3 absent (ollama 429 weekly limit); objections and the r1 changes:** artifact `/Users/voightkampff/.ailang/state/mission-quorum-world/w-world-approvals-by-key-a4-2026-11-12T12-44-28Z.json`. Neither disputed the direction.
  - **(gpt6-1-sol) D4's "exactly one packet" rested on a 101-packet hash-ordered scan with no completeness branch** → lands on §2.4 (the complete request-wide discovery, `Discovery{Complete, Scanned, Carrying, Undecodable}`, overflow sentinel, "listing truncation warnings do not authorize decisions"), §2.5 (D4 split into D4a completeness and D4b count, `DecisionFacts` carries the discovery), §2.8 (the listing is display only), §3, §4 T7, §6 (DL4a, DL4b, DS1, RD9b with 150 packets and the second outside the first 101, RD9c with 2,001, RD9d must-live), §9 gates and sizes, §11 RK3b, V33, V34. The reviewer's first option (refuse on an incomplete scan) is applied; the request-indexed alternative is named as RK3b's lift, not built.
  - **(gemini-4-argon) a crash between `RecordKeyDecision` and `DecideApproval` strands the counter, a counter-0 assertion could be replayed, and two refs lacked the `sha256:*` CHECK** → lands on §2.7 (both CHECKs verbatim; the reviewer's second option chosen: K0, an idempotent retry on the byte-equal artifact of an unlanded row that skips `useApprover`/`usePasskey` and goes to `DecideApproval` + `LandKeyDecision`; the counter-0 analysis), §4 T15, §6 (KD5, KD6, KD7, KD8, RD13), §9 gates and sizes, V32. Measured (V32): the counter-0 replay **is** accepted twice by the counter alone; what refuses a second use is D6 after landing and K0's single row while unlanded, so one counter-0 assertion lands at most one decision.
- **Round 2 (2026-11-12T12:53:54Z): BLOCKED 2/2 present (gemini-4-argon, gpt6-1-sol); oc-glm-5-3 and oc-kimi-k3 absent (ollama 429 weekly limit, unreachable).** Artifact `/Users/voightkampff/.ailang/state/mission-quorum-world/w-world-approvals-by-key-a4-2026-11-12T12-53-54Z.json`. Both objections land on ONE surface, §2.7's retry order (r1's K0 sat after the judge's W15, so a counter-N retry could never reach it); round 1's discovery fix drew no objection. Neither disputed the direction; both carried a concrete `proposed_fix`.
  - **r2, by the controller's narrow-refinement carve-out (gate 2), not a designer run:** §2.7's K0 bullet replaced by gpt6-1-sol's two-path judge in its own words (W15 deferred to the writer transaction; exact match bypasses only the counter; revocation re-checked on both paths; mint judge unchanged); gemini-4-argon's equivalent ordering cited, its `useApprover` skip not taken because gpt6-1-sol's stricter revocation re-check is kept; KD7/RD13 marked; RD14 added with both reviewers' test text verbatim. After the independent audit (evaluator sonnet, CARVE-OUT VALID 84/100, zero blocking): its N1 (an end-to-end revoked-approver arm cannot kill its mutant, because the read-only judge refuses a revoked approver first) moved that arm to the store-level KD9; N2/N3 the writer's late counter refusal is a 422 with the judge's text, and the A2 verifier gains a defer-W15 form; N5 RD14/KD9 added to M5 and T15; N6 the controller's elaboration is labelled. N4 (fold K0 into the Judge and handler-order bullets) is left to the planner: the two-path bullet sits directly under them. No other section changed. An independent evaluator audits this carve-out before any plan.
- **Amendment confirmation:** (a) the decision message and the derived challenge (§0) — confirmed /
  reverted; filed as D-WORLD-115 on (date).
- **Anchors found false by the designer:** `MintRequestIDPattern` is defined at
  `host/store/mint_requests.go:46`; `:105` is a use (V21). Credentials schema is 7 at the base (V21),
  as the directive's unverified list implied.
