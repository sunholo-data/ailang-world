# w-world-address-bar-verbs-v3: confirm `episode finish` with an enrolled key on the surface

- **Status**: Planned **r0**. Designed in iteration 399, 2026-11-17. Child of
  [`w-world-address-bar-verbs.md`](w-world-address-bar-verbs.md) (r2, reviewed), queue row **243**,
  parent phase **V3**. V1 (PR #582) and V2 (V2a PR #585, V2b PR #587, merge `c251132f`) are landed
  and are the code this design measures. This doc carries **eight explicit amendments** (§0.3), one
  ledger ask, **D-WORLD-118** (§10), and splits V3 into **five sprints** (§8).
- **Authority**: **D-WORLD-107 = A** (the bar "shows the CLI's prompt, and is confirmed by row 201's
  key approval, never a typed y from a page"). **D-WORLD-109 = A** (the loop builds and judges
  against fakes and test keys only; never enrols a real key, mints or publishes). **D-WORLD-89**
  R1–R5 (R1 names "mints and `Human.Approve`" only, W33 — hence D-WORLD-118; R3 enrolment typed;
  R4 publishing typed). **D-WORLD-115 = A** (A4's derived-challenge amendment; not reopened).
  Row 201's settled siblings `-a0a1`, `-a1b`, `-a2`, `-a3`, `-a4` are reused, not reopened, except
  the two stated-reason amendments in §0.3 (AM6, AM7).
- **Item**: `ailang-worldd episode finish <episode> --preset <name> --confirm-on-surface` replaces
  the pull-request prompt's typed `y` with a verified enrolled key (passkey or SSHSIG) given on a
  World-app card, `/surface/confirm/<id>`, while the CLI waits at its terminal. The key signs the
  **real post-commit prompt bytes**, bound with the canonical argv, the finish scope, the title and
  the full body. The typed `y` is unchanged without the flag. `episode start` is **not** in V3 (P5).
- **Designer lane**: `claude:claude-opus-5-5`.
- **Estimated**: **five sprints, ≈ 4,600 raw / ≈ 6,000 calibrated changed lines** (`git diff
  --numstat -- cmd host docs`), each sprint planned ≤ 1,000 raw so that it lands ≤ 1,400 at the
  measured 1.35× (V2a planned ≈ 970, landed 1,299; V2b ≈ 930 → 1,255; W31). The parent's ≈ 4,500
  holds as a raw number; the calibration and the 1,400 cap make it five sprints, not 2–3 (§8).
- **Measurement base**: `994ba12358dc50246db2862d78f66d64e03b29f4` (origin/dev, PR #588's merge;
  `cmd/` and `host/` byte-identical to V2b's merge `c251132f`, W1), worktree
  `/Users/voightkampff/dev/sunholo-data/.wt-world-iter399-row243-v3`, branch
  `design/world-iter399-row243-v3`, measured 2026-11-17.
- **Instruments**: the worktree's tree, `grep`/`sed`/`git`. No Go test, daemon, browser or key was
  run in this fire; every runtime premise is a planner obligation (§6).

---

## 0. Scope, the relationship to the parent, and the split

### 0.1 The parent's V3 row, measured against the tree

The parent's §0.3 V3 row: "the verb-approval law: a `world-verb/v1` subject over the canonical argv
and the exact prompt bytes, A4's judges (SSHSIG/passkey), credentials schema 9, a card; `episode
finish` and `episode start` gain a 'confirm on the surface' mode that replaces the typed `y` with a
verified key (typed `y` unchanged at a terminal); 2–3 sprints ≈ 4,500". Seven premises, measured:

- **P1 — "the exact prompt bytes": they exist only after a commit, and no preview can reproduce
  them.** `episode finish` commits (`episode_finish.go:348`) and prints `committed at` (`:380`)
  before `finishPrompt` (`:410`); a decline says "commit kept at <sha>" (`:413`) (W3). The commit is
  made with no pinned date (zero `GIT_*_DATE`/`--date` in `host/worktree` and the CLI; control:
  `commit.go:48` is the commit call) so its sha depends on the wall clock (W7); the body embeds that
  sha (`add.Commit = sha` before the render, `:352-354`; `Commit: %s`, `episode_markdown.go:78`)
  (W8); and `--dry-run` is not side-effect-free: it records each check's outcome in World
  (`:290`) and makes, then resets, a real commit (`:355-358`) (W9). Q3 is decided accordingly
  (§2.1): **commit-then-approve**.
- **P2 — "replaces the typed `y`": finish has more than one.** Each unsandboxed check asks its own
  `Run this check? [y/N]` (`:266-273`) besides the pull-request prompt (W4, W43). V3 replaces only
  the pull-request `y` and refuses surface mode, before any check runs, when the pinned start has
  an unsandboxed check (§2.7, RV3).
- **P3 — "A4's judges": they are bound to A4's 91-byte decision message.** `JudgeDecisionSSH` and
  `JudgeDecisionAssertion` take `DecisionJudgeInput{Verb, Packet, Request, Expiry, Now}` and verify
  `world-decision/v1` (W14). The shape V3 needs — a **stored** request with a digest, a random
  challenge and an expiry, recorded by a surface POST, waited for by a CLI process — is A3a's
  (`JudgeAssertion`/`JudgeSSHSignature` over a `MintRequestRow`, W15; `approve --await`, W22). V3
  reuses the parts both share — `approval.Admit` (generic over bytes, W47), `checkApprover`,
  `assertionSteps(…, h, …)` and W15 (W16) — under a new message, and takes one step from A4: the
  daemon writes the world approval itself, under `humanMu` (W20).
- **P4 — "credentials schema 9": holds.** The store is at 8 (`credentials.go:44`, ladder to `< 8`,
  `:375`; no `verb_requests` anywhere, control `key_decisions` found) (W17). Neither existing
  request table can hold another id shape: both `mint_requests` and `surface_approvals` CHECK
  `request_id GLOB 'mr-[0-9a-f]*'` (W18). V3 adds exactly one table, so schema **9**.
- **P5 — "`episode start` gains a confirm-on-surface mode": already covered for mints, and the gap
  is V4's.** `episode start --request` stores a mint request that the A3a card approves with a key
  and `approve --await` mints (parent §0.1 P4) — D-WORLD-89 R1 already authorises keys for mints
  (W33). What it refuses is a **launch preset** (`episode.go:249`, "a request carries no launch
  callbacks or start records", `:241`; `mint_request.go:219`) (W24, W25): the missing piece is the
  minter writing a preset's start record and launching the agent — V4's host runner, not a new
  law. V3 drops `episode start` (AM1).
- **P6 — "a `world-verb/v1` subject over the canonical argv and the exact prompt bytes": incomplete.**
  The prompt shows only the body's first six lines and its length (`episode_finish_operator.go:160-171`,
  W12), while the approval the daemon records binds the title and body by SHA-256 inside the
  finish scope (`FormatFinishApprovalScope`, `approve.go:698-701`, W52). V3's subject therefore
  carries the scope, the title, the **full** body, the committer and the protected paths as well,
  and the daemon refuses a request whose prompt is not exactly the terminal's rendering of those
  fields (§2.2).
- **P7 — "2–3 sprints ≈ 4,500": the raw total holds; the count does not.** §8 measures ≈ 4,600 raw,
  five sprints under the 1,400 cap at 1.35× (W31).

### 0.2 What V3 changes from a page

After sprint 4 (and only if D-WORLD-118 is A or B), one POST, `POST /surface/confirm/{id}`, records
a verified key's confirmation of one stored verb request and writes one finish approval to the
world store, decided by that key. Nothing else from a page changes World: the bar still runs no
mutating verb (it answers `terminal` for `episode finish`, §2.8), no page ever submits or is asked
for a `y`, and the CLI process — not the daemon, not the page — commits, pushes and opens the PR.

### 0.3 Amendments to attended or ratified text (explicit)

| # | doc and text | amended to | why |
|---|---|---|---|
| **AM1** | parent §0.3 V3 row: "`episode finish` and `episode start` gain a 'confirm on the surface' mode"; "A4's judges"; "2–3 sprints ≈ 4,500" | **`episode finish` only**; A3a-shaped judges over a new message (`world-verb/v1`); **five sprints ≈ 4,600 raw / ≈ 6,000 calibrated** | P5 (start is covered for mints; its launch-preset gap is V4's); P3; P7 |
| **AM2** | parent §12 Q3 default: "one approval over the exact prompt bytes from a preview, with the run refusing if its prompt differs" | **commit-then-approve**: the key approves the real prompt after the commit, where the `y` sits today | P1: a preview's sha and body never equal the real ones (W7–W9); §2.1 |
| **AM3** | parent RK3: "V3 must choose intent-approval or a preview-bound approval" | the third option, commit-then-approve, chosen; "commit kept at <sha>" on decline, expiry or supersession; a re-run resumes | §2.1 |
| **AM4** | parent §10 step 4: "the app shows `finishPrompt`'s lines … he approves with his key" | holds, on `/surface/confirm/<id>`; **plus**: the preset's pinned checks must all be sandboxed, or surface mode refuses before any check runs | P2, RV3 |
| **AM5** | parent §2.1(b) `episode finish` row ("`terminal` answer (V3)") and the `approve` row (`mr-`/`sha256:` only) | `episode finish` stays `terminal` until V4, its reason naming `--confirm-on-surface` (sprint 5, only under A/B); `approve vr-<16 hex>` navigates to `/surface/confirm/<id>` | §2.8 |
| **AM6** | A4 §2.6 L1's **stated reason**: "these scopes have their own typed owners (`world-publish approve`, `episode finish`) that show the exact bytes and run their own fences; A4's card is not their card" | "…their own owners — `world-publish approve` (typed only, R4) and `episode finish` (its typed `y`, or, under D-WORLD-118, its own confirmation card `/surface/confirm/<id>`) — that show the exact bytes and run their own fences; A4's card is not their card." **L1's code and `reservedReason`'s text are unchanged** (`decision_law.go:80`, `:116-131`, W13): finish approvals are still made by `episode finish` | parent §0.2 flagged it; D-WORLD-118 rules it |
| **AM7** | A4 §2.6: "L2 for finish is not added. The finish record is reachable only behind the typed operator proof, and its approve phase creates its own request" | a finish approval may also be created by `POST /surface/confirm/{id}`, decided by a key (`decidedBy` `ssh:`/`webauthn:`); the record phase is unchanged (still the operator proof); L2 for finish stays absent, now **by design** (a key-decided finish approval is legitimate under A/B). Under B the POST refuses a non-draft scope | §2.6 |
| **AM8** | A4 §2.7 "`TestDaemonCredentialWriterSites` is edited to pin exactly these two" (now three, W21) | **four**: + `handleSurfaceConfirm` | §2.6, §7 |

A2 is untouched (its `assertionSteps` and W15 are called unchanged, W16). A3a is reused by shape
(the request/record/wait pattern and its guards), with no text changed. The ledger row D-WORLD-89
is not edited; D-WORLD-118 is a new row (§10).

### 0.4 Where V3 stops and V4 starts

- **V3**: the CLI runs **at a terminal** (finish's F1 — "no controlling terminal" refusal before any
  IO, `episode_finish.go:85-88`, W2 — is **kept**); with `--confirm-on-surface` it stores a request,
  prints its card's links, waits up to 10 minutes, and opens the PR after the daemon has recorded a
  key's confirmation. The bar answers `terminal` for `episode finish`.
- **V4**: who runs that CLI process with no terminal (relaxing F1 for surface mode only, parent Q4);
  launching it from the bar; `episode start` with a launch preset (the minter's start record and
  launch, P5); handing interactive launches to herdr; per-check confirmation of unsandboxed checks
  (RV3), if wanted.

## 1. What the code does today (measured at `994ba123`)

1. **The finish prompt.** After checks, the commit and `writeFinishBody`, finish resolves the
   repository and label (`seam.Prepare`), builds `broker.FinishApprovalScope{…, ExpiresAt:
   time.Now().Unix() + 900}` (`episode_finish.go:391-395`, W10) and, unless it is resuming a
   reusable approval, calls `finishPrompt(tty, scope, title, rendered, protected,
   pinned.Committer)` and reads one line: anything but `y` → "pull request declined … commit kept
   at <sha>" (`:405-414`) (W3). `finishPrompt` writes `committer:`, `repository:`, `branch:`,
   `sha:`, `base:`, `title:`, `body (N bytes), first lines:` (≤ 6 lines), `label:`, `draft|ready`,
   `protected paths changed:` and `Open this pull request? [y/N] ` (`episode_finish_operator.go:154-171`)
   (W12). `displayCommitter` turns `""` into `(repository git config)` (`committer.go:3`, 3 callers)
   (W45).
2. **The typed `y` is recorded by the daemon, behind the operator proof.** On `y`, the CLI writes a
   random proof's hash into `operator-proofs/` beside the credentials store
   (`writeOperatorProof`, `episode_finish_operator.go:30-73`, called at `:76`) and POSTs
   `{phase:"approve", scope}` to `/v1/episodes/{ep}/finish` (`episode_finish.go:415`); the daemon
   consumes the proof under `humanMu` and calls `broker.RecordOperatorApproval(ctx, d.store, scope,
   "operator:episode-finish", "operator:terminal", now)` (`host/daemon/episode_finish.go:158-178`)
   (W5). The record phase writes a proof with **no** `y` (`episode_finish.go:443`) (W5).
   `RecordOperatorApproval` takes `decidedBy` as a parameter and stores a `Human.Approve` request
   then its decision (`host/broker/approve.go:722-733`) (W6); the decision wire carries
   `decidedBy` (`:79`) (W37). The daemon's approve phase bounds the scope to `now < expires ≤
   now + 3600` (`host/daemon/episode_finish.go:169`) and the record phase re-checks the decision's
   scope against the record (`:186-195`) (W10, W38). Before dispatch the CLI re-checks HEAD,
   remote, title and body against the scope (F6, `episode_finish.go:423-430`) (W38).
3. **Resume.** A re-run on a World-provenance HEAD is mode `resume` (no checks re-run; the saved body
   is reused) or `resume-dispatch` (`finishMode`, `episode_finish_resume.go:19-42`); only
   `resume-dispatch` may reuse a recorded approval (`finishReusableApproval`, `:118-129`) (W11).
   F6 also requires the daemon to serve the preset's credentials store (`episode_finish.go:114-115`)
   (W36). The PR body is capped at `reportMarkdownMax = 60000` bytes (`episode_markdown.go:11`) (W34).
4. **Row 201's machinery.** The mint message is 90 bytes, `world-approval/v1` ‖ 0x00 ‖ digest ‖
   challenge ‖ u64be(expiry) (`request.go:26`, `:211-222`); a request lives 10 minutes
   (`RequestLifetime`, `:29`) (W15). `approval.Admit` is a pure verdict over `{Now, CreatedAt,
   ExpiresAt, Consumed, StoredDigest, BytesDigest}` (`admit.go:37-58`) (W47). The laws share
   `checkApprover` (P2–P5, `signed.go:75`) and `assertionSteps(response, a, h, hErr)` (W1–W14,
   `passkey_assert.go:60`); W15 is `:46` (W16). Consumes reuse `useApprover`
   (`approvers.go:245`) and `usePasskey` (`passkeys.go:225`) through `store.ApproverUse`
   (`mint_requests.go:300`) (W48). A4's decision message is 91 bytes under `world-decision@sunholo`
   (`decision.go:15-21`) (W14). L1 reserves finish-shaped scopes from A4's law (`decision_law.go:80`,
   `:121`) (W13); L2 refuses `ssh:`/`webauthn:` decisions at publish (`approve.go:579`) (W37).
5. **The A3a surface.** `approvalRoute` answers 421 for a `Host` outside the origin set before any
   handler (`surface_approve.go:76-89`); the approve POST checks R1 id, R3 Origin, R4
   `Sec-Fetch-Site`, R5 JSON, R6 size (`MaxApproveBodyBytes = 20 << 10`, `approve_model.go:27`),
   then admission, digest, the judge (R11), `Binds` (R12), a non-waiting `OpenCredentials` (R13)
   and the record (R14) (`surface_approve.go:203-300`) (W19, W49). Its body grammar is
   `ParseApproveBody` (`approve_model.go:197`: `v`, `digest`, exactly one of `webauthn`/`sshsig`)
   (W49). A4's POST adds R0 (any `Authorization`/`Cookie` → 403, `surface_decide.go:159`), takes
   `humanMu` (`:186`) and opens the writer in `keyArm` (`surface_decide_key.go:66`) (W20).
   `TestDaemonCredentialWriterSites` pins exactly `ensureCredentialsStore`, `handleSurfaceApprove`,
   `keyArm` (`surface_approve_post_guard_test.go:168`) (W21). The `/surface` approval routes are
   `daemon.go:1139-1148` (W30); the CSP has `connect-src 'self'` (`surface_page.go:19`) (W29).
   `commitBudget` is 3 s (`daemon.go:150`) (W35).
6. **The minter's wait.** `approve --await` polls every second (`awaitPollInterval`,
   `approve_await.go:24`) through `env.after` (`:130`), each poll one read-only open under its own
   3 s root (`pollSurfaceApproval`, `:70`), until `row.ExpiresAt` (`:105`) (W22);
   `TestProductionContextRoots` pins each such root (`context_roots_test.go:19`, `:30`) (W23).
7. **The bar.** `Lookup` gives `episode start|finish` class `Mutate`, reason "… changes World. The
   app cannot run it yet." (`lookup.go:90-92`); `ApproveTarget` knows `mr-` and `sha256:` only
   (`:120-131`) (W26); the page's navigate allowlist is the same two prefixes (`app.js:496`) (W44);
   *Needs you* is `needsYouHTML` (`app.js:424-448`) over `requests.json`'s `pending`, `decided`,
   `decisions` (`approve_model.go:67-72`) (W50). `host/verbs` may import only a pure set
   (`surface_verb_test.go:151`) (W46). The CLI's `approve` dispatches a `sha256:` ref to the decide
   mode and everything else to the mint path (`approve.go:228`) (W40).

## 2. Design

### 2.1 Q3: commit-then-approve

| option | what the key approves | what happens on a mismatch | verdict |
|---|---|---|---|
| **(a) intent** | the argv and a pre-commit summary, before checks | the PR's real body (check results, sha) is decided after the approval: nothing to mismatch, and nothing the human saw binds what opens | **rejected**: approves bytes that do not exist yet |
| **(b) preview-bound** | a side-effect-free preview's prompt; the real run refuses unless its prompt is byte-equal | **every run mismatches**: the commit sha has an unpinned date (W7); the body embeds the sha (W8) and the check-outcome entries the preview would have had to record (W9); `--dry-run` records outcomes and commits (W9), so "side-effect-free" does not exist | **rejected**: fails closed always, so it is an approval nobody can give |
| **(c) commit-then-approve** | the real prompt after the commit, exactly where the `y` sits (`episode_finish.go:410`) | none by construction: the request is built from the same `scope`, title, body, protected paths and committer that `finishPrompt` renders | **chosen** |

**Decline, expiry, supersession** each leave the commit, as a typed `n` does today: "… commit kept
at <sha>" plus `steps(sha)` (`episode_finish.go:157-159`). A re-run is mode `resume` (W11): no
check re-runs, the saved body is reused, so the new request's prompt bytes equal the old one's;
only the scope's `expires` and the request's id, challenge and expiry are new.

### 2.2 The subject: `world/verb-request/v1`, and the message `world-verb/v1`

**The request** (NEW, `host/approval/verb.go`, stdlib-only like the rest of the package, W15):

```go
const VerbRequestSchema   = "world/verb-request/v1"
const VerbRequestLifetime = 10 * time.Minute          // as RequestLifetime (request.go:29)
const MaxVerbRequestBytes = 128 << 10                 // body ≤ 60000 (W34) + prompt, scope, title, argv

type VerbRequest struct {        // canonical JSON in this field order (CanonicalVerb)
	Schema    string   `json:"schema"`    // VerbRequestSchema
	Verb      string   `json:"verb"`      // "episode finish" (the only admissible value)
	Argv      []string `json:"argv"`      // the canonical argv, §2.7
	Episode   string   `json:"episode"`
	Prompt    string   `json:"prompt"`    // broker.FinishPromptText(...): the terminal's lines, without "Open this pull request? [y/N] "
	Scope     string   `json:"scope"`     // broker.FormatFinishApprovalScope(scope), the approval the daemon records
	Title     string   `json:"title"`
	Body      string   `json:"body"`      // the full PR body (the file the dispatch sends)
	Committer string   `json:"committer"` // the pinned committer, raw ("" allowed)
	Protected []string `json:"protected"`
}
func (r VerbRequest) Validate() error            // NEW: schema; verb; episode ^[a-z0-9][a-z0-9-]{0,63}$; argv[0:4] == ailang-worldd episode finish <episode>, ≤ 64 words, each valid UTF-8 with no byte < 0x20 or 0x7f; prompt, scope, title non-empty; body ≤ 60000; every string valid UTF-8; canonical size ≤ MaxVerbRequestBytes
func CanonicalVerb(r VerbRequest) ([]byte, error) // NEW: Canonical's encoder settings (request.go:165-175)
func DecodeVerb(b []byte) (VerbRequest, error)    // NEW: Decode's strictness (request.go:178-194): unknown fields, trailing data, non-canonical → refused
```

**The message** (NEW, same file):

```
VerbMessage(digest, challenge, expiry) =
    "world-verb/v1"   (13 bytes, VerbPrefix)
  ‖ 0x00
  ‖ digest            (32: SHA-256 of CanonicalVerb(r), the stored bytes)
  ‖ challenge         (32: random, stored with the request, single use)
  ‖ u64be(expiry)     (8: the request's expires_at)
  = 86 bytes (VerbMessageLen)
VerbMessageHash = SHA-256(VerbMessage)        -- the WebAuthn challenge
VerbNamespace   = "world-verb@sunholo"        -- the SSHSIG namespace
```

- **Domain separation.** The three prefixes share `world-` and differ at byte 6 (`a`pproval,
  `d`ecision, `v`erb), so no message of one kind is a prefix of, or equal to, one of another; the
  three SSHSIG namespaces differ and are inside the signed blob (`ErrNamespace`); the passkey
  challenges are SHA-256 of different messages (VL4).
- **Why a stored random challenge here** (and not A4's derived one, D-WORLD-115): the request is a
  credentials-store row the CLI writes, exactly as a mint request is (A0); a stored challenge and
  `consumed_at` give freshness and single use with A0's own machinery. A4 derived its challenge
  because its subject (a committed packet) has no such row.
- **No verb byte.** There is one admissible action, confirm; there is no deny (RV5). A later deny
  would be a new prefix version, not a reinterpretation of these bytes.

**The prompt is a pure function of the other fields** (NEW in `host/broker`, beside
`FinishApprovalScope`, `approve.go:689-701`, W52):

```go
func DisplayCommitter(v string) string    // "" → "(repository git config)"; cmd's displayCommitter (committer.go:3) delegates
func FinishPromptText(s FinishApprovalScope, title, body string, protected []string, committer string) string
    // finishPrompt's body (episode_finish_operator.go:155-171) minus the final question
```

`finishPrompt` becomes `fmt.Fprint(tty, broker.FinishPromptText(scope, title, body, protected,
committer)+"Open this pull request? [y/N] ")`: byte-identical at the terminal (VT1; the existing
`episode_finish_m2/m3` prompt assertions stay unedited, W43).

**How the card shows what is signed** (§2.5): the canonical argv as one quoted line; the prompt
verbatim in a `<pre>` (the same lines the terminal prints); then "Open this pull request?" with
the key controls; then *What your key signs*: the digest, the challenge, the expiry, the scope
text, the title, the full body (in a `<details>`), the committer and the protected paths, every
field of the signed bytes rendered by one table, fail-closed.

### 2.3 The law and the judges (sprint 1)

**`host/approval/verb_signed.go`** (NEW), a copy of `signed.go`/`passkey_assert.go`'s shape with
the verb message:

```go
type VerbSignatureInput struct { Admit AdmitInput; Challenge [32]byte; Signature []byte; Approver Approver }
type VerbAssertionInput struct { Admit AdmitInput; Challenge [32]byte; Response []byte;  Approver Approver }
func ConfirmBySignature(in VerbSignatureInput) (*SignedVerb, error) // P1 Admit (admit.go:44) → P2–P5 checkApprover → VerifySSHSIG(·, VerbNamespace, VerbMessage(...))
func ConfirmByAssertion(in VerbAssertionInput) (*SignedVerb, error) // P1 → P2–P5, P3w → assertionSteps(h = VerbMessageHash) → W15 (as passkey_assert.go:46)
type SignedVerb struct { /* unexported: ok, digest, challenge, expiry, name, format, key, credentialID, signCount */ }
func (s *SignedVerb) Binds(digest, challenge [32]byte, expiry int64) error
func (s *SignedVerb) ApprovedBy() string                          // format ":" name; "" unless verified
func (s *SignedVerb) Use() (name, format string, key []byte, ok bool)
func (s *SignedVerb) PasskeyUse() (credentialID []byte, signCount uint32, ok bool)
```

W15 runs **in full** in the law (the mint path's rule): the counter moves inside the consume
(§2.4), and a crashed consume is not retried with the same artifact (RV2), so A4's deferred-W15
two-path judge is not needed.

**`host/authority/judge_verb.go`** (NEW): `JudgeVerbSSH(ctx, ro, row store.VerbRequestRow, armored,
now)` and `JudgeVerbAssertion(ctx, ro, row, response, now)`, `judge.go`'s steps and `*JudgeError`
step names (parse, store, unknown-key, stored-key, law) so texts are shared; a
`judgeVerbAdmit(row, now)` mirrors `judgeAdmit` (`judge.go:49`).

**`broker.FinishConfirmRefusal(f FinishConfirmFacts) string`** (NEW, pure, `host/broker`), the
daemon's step C10 (§2.6), "" when admissible:

| id | refused when | reason |
|---|---|---|
| F1 | `ParseFinishApprovalScope(f.Scope)` fails | "the request's scope is not a finish approval" |
| F2 | not `f.RequestExpires < scope.ExpiresAt ≤ f.Now + 3600` | "the request outlives its approval scope, or the scope is over an hour ahead" (the approve phase's own bound, W10) |
| F3 | `sha256(f.Title) ≠ scope.TitleSHA256` | "the title is not the one the scope approves" |
| F4 | `sha256(f.Body) ≠ scope.BodySHA256` | "the body is not the one the scope approves" |
| F5 | `f.Prompt ≠ FinishPromptText(scope, f.Title, f.Body, f.Protected, f.Committer)` | "the prompt is not the terminal's rendering of what is approved" |
| F6 (B only) | `!scope.Draft` | "D-WORLD-118 = B: a ready pull request is confirmed with a typed y" |

### 2.4 The store: credentials schema 9 `verb_requests` (sprint 2)

```sql
-- Verb requests (version 9, row 243 V3). One stored request per `episode finish
-- --confirm-on-surface`; inert until a verified key confirms it on the surface. The authority is
-- the world store's finish approval (decision_ref), never this row.
CREATE TABLE IF NOT EXISTS verb_requests (
    request_id   TEXT PRIMARY KEY CHECK (length(request_id) = 19 AND request_id GLOB 'vr-[0-9a-f]*'),
    request      BLOB NOT NULL CHECK (length(request) BETWEEN 1 AND 131072),
    digest       TEXT NOT NULL CHECK (length(digest) = 64),
    challenge    TEXT NOT NULL CHECK (length(challenge) = 64),
    episode_id   TEXT NOT NULL CHECK (episode_id <> ''),
    created_at   INTEGER NOT NULL,
    expires_at   INTEGER NOT NULL CHECK (expires_at > created_at),
    consumed_at  INTEGER,
    outcome      TEXT CHECK (outcome IS NULL OR outcome IN ('approved', 'superseded')),
    approved_by  TEXT CHECK (approved_by IS NULL OR approved_by <> ''),
    artifact     BLOB CHECK (artifact IS NULL OR length(artifact) BETWEEN 1 AND 16384),
    decision_ref TEXT CHECK (decision_ref IS NULL OR (length(decision_ref) = 71 AND decision_ref GLOB 'sha256:*')),
    CHECK ((consumed_at IS NULL) = (outcome IS NULL)),
    CHECK ((outcome IS 'approved') = (approved_by IS NOT NULL AND artifact IS NOT NULL)),
    CHECK (decision_ref IS NULL OR outcome IS 'approved')
);
```

`credentialsSchemaVersion = 9` (`credentials.go:44`); `upgradeCredentials` gains `if version < 9`
(after `:375`); the canonical-DDL test compares fresh and upgraded column for column, the A2/A3
pattern. `host/store/verb_requests.go` (NEW):

- `VerbRequestIDPattern = ^vr-[0-9a-f]{16}$`; `MaxPendingVerbRequests = 16`; sentinels
  `ErrVerbRequestUnknown|Consumed|Superseded|Expired|ClockBehind|DigestChanged`,
  `ErrTooManyPendingVerbRequests`, `ErrVerbDecisionLanded` (NEW; the mint sentinels' texts name
  mint requests, `mint_requests.go:49-54`, W39).
- **`PutVerbRequest(ctx, row VerbRequestRow, now) error`**, one transaction: S1 deadline, below v9 →
  `*OlderCredentialsStoreError{Feature: "verb requests", Need: 9}`; S2 grammar (id, 64-hex digest
  and challenge, `digest == sha256(request)`, sizes); S3 **supersede**: `UPDATE … SET consumed_at =
  now, outcome = 'superseded' WHERE episode_id = ? AND consumed_at IS NULL`; S4 prune unconsumed
  rows expired a day ago; S5 cap → `ErrTooManyPendingVerbRequests`; S6 `INSERT`.
- `VerbRequest(ctx, id)`, `PendingVerbRequests(ctx, now, limit)` (read-only, deadline-required).
- **`ConsumeVerbRequest(ctx, id, digest string, use ApproverUse, approvedBy string, artifact []byte,
  now) error`**, one transaction: re-read the row; `Admit`-order refusals (consumed → `Consumed` or
  `Superseded` by `outcome`; clock behind; expired; digest); `useApprover` (approver still active,
  bound by key); `usePasskey` when `use.Passkey != nil` (the counter moves here, exactly once);
  `UPDATE … consumed_at, outcome = 'approved', approved_by, artifact`.
- **`LandVerbDecision(ctx, id, decisionRef) error`**: `UPDATE … SET decision_ref = ? WHERE request_id
  = ? AND outcome = 'approved' AND decision_ref IS NULL`; zero rows → `ErrVerbDecisionLanded`.

The CLI writes (Put) with a short non-waiting writer; the daemon writes (Consume, Land) with a
short non-waiting writer per POST; both read through `OpenCredentialsReadOnly`.

### 2.5 The read routes and the card (sprints 2–3; read-only)

All behind `d.approvalRoute` (R2's 421, `surface_approve.go:76`), `no-store`, `nosniff`:

- **`GET /surface/confirmations/{id}`** (NEW, sprint 2; `{id}` ≠ `VerbRequestIDPattern` → 400;
  unknown → 404; store below v9 → 503 naming the `serve` restart):
  ```json
  {"v":1,"id":"vr-…","request":"<base64url canonical bytes>","digest":"<64 hex>","challenge":"<64 hex>",
   "created":…,"expires":…,"now":…,"state":"pending|approved|superseded|expired",
   "approvedBy":"","decision":"",
   "passkeys":[{"origin":"…","rp":"…","allow":["…"]}],"ssh":{"namespace":"world-verb@sunholo"}}
  ```
  Bytes and the raw challenge, **never a precomputed hash**, and **never the artifact**; `passkeys`
  grouped as `handleSurfaceRequestCard` groups them (`surface_approve.go:146`).
- **`GET /surface/requests.json`** (`surface_approve.go:92`) gains `"confirmations":[…]` from
  `PendingVerbRequests(now, 16)`: `{id, episode, repository, head, sha, draft, expires}` decoded from
  the stored bytes; a row that does not decode, or whose bytes' SHA-256 is not its digest, is listed
  with `"refused":"<reason>"` and no summary (A3's never-hidden rule); below v9 →
  `"confirmationsUnavailable":"<reason>"`, other lists unchanged.
- **`GET /surface/confirm/{id}` + `GET /surface/confirm.js`** (NEW, sprint 3; embedded
  `host/surfacepage/confirm/{index.html,confirm.js}` as `approve.go:13-18` embeds A3a's; unchanged
  `surfaceCSP`). `confirm.js`: (1) fetch the data (10 s `AbortController`, A3 §2.7's rule); (2)
  `SHA-256(bytes) == digest` or fail closed; (3) decode the bytes by one table `CONFIRM_FIELDS`
  (`schema`, `verb`, `argv`, `episode`, `prompt`, `scope`, `title`, `body`, `committer`,
  `protected`, with types): an unlisted key, a missing key or a wrong type → "a field the page cannot
  show; confirm at the terminal"; (4) render §2.2's layout with `textContent` only; (5) build the
  86-byte message and its hash **in the page**; (6) sprint 3 offers **Save the message** (a `blob:`
  `world-verb-<id>.msg`) and the line `ssh-keygen -Y sign -n world-verb@sunholo -f <key>
  world-verb-<id>.msg`; a state other than `pending` shows that state and no controls.
- **Sprint 4** adds the submit half (§2.6): **Approve with a passkey** (`navigator.credentials.get`,
  `challenge = VerbMessageHash`, UV required, `allowCredentials` from `passkeys`, A2's 240 s abort)
  and the SSHSIG paste box; the answer's `approvedBy`/`decision` or its refusal class and reason;
  then re-read. `confirm.js` repeats A3's encoding helpers (≈ 40 lines), as `decide.js` did.
- **app.js** *Needs you* (sprint 3) lists each confirmation under the decisions: "Confirm · episode
  finish <ep> · <repository> <head>@<sha12> · draft|ready · expires hh:mm" linking to
  `/surface/confirm/<id>`; a refused row shows its reason and no link; all text through `esc()`.

### 2.6 `POST /surface/confirm/{id}`: the confirmation (sprint 4; only under A or B)

Body: **`surfacepage.ParseApproveBody` unchanged** (`approve_model.go:197`, W49): `{"v":1,"digest":
"<64 hex>"}` plus exactly one of `webauthn` / `sshsig`, strict. There is **no typed arm**: a body
with `typed`, `y` or any other key is refused by the strict decoder (C6).

| id | step (in order) | refusal |
|---|---|---|
| C0 | any `Authorization` or `Cookie` header (A4's R0, `surface_decide.go:159`) | 403 `Forbidden` |
| C1 | `{id}` matches `VerbRequestIDPattern` | 400 |
| R2 | `Host` ∈ O (`approvalRoute`) | 421 |
| C3–C5 | `Origin` ∈ O; `Sec-Fetch-Site` empty or `same-origin`; `application/json` (`allowedApprovalOrigin`, `sameOriginFetch`, `jsonMediaType`, `surface_approve_guard.go:41-66`) | 403 / 403 / 415 |
| C6 | ≤ 20 KiB, `ParseApproveBody` | 413 / 400 |
| — | `d.humanMu.Lock()`; ctx = `commitBudget` (as `surface_decide.go:186-188`) | |
| C7 | the row exists (read-only handle) | 404 |
| C8 | `approval.Admit` on the row (consumed → "already confirmed by <approved_by>" or "superseded by a newer finish of <episode>"; clock; expired; stored digest = SHA-256(bytes)) | 409 |
| C9 | body `digest` == the row's ("the request changed since the card was shown; reload it") | 409 |
| C10 | `DecodeVerb` + `Validate` + `broker.FinishConfirmRefusal` (§2.3 F1–F5; F6 under B) | 409 `NotConfirmable` |
| C11 | the law: `JudgeVerbSSH` / `JudgeVerbAssertion` (unknown key, revoked, bad signature, namespace, origin, RP hash, UV, counter) | 422 `NotApproved` |
| C12 | `signed.Binds(digest, challenge, expires_at)` | 422 |
| C13 | `store.OpenCredentials` non-waiting: busy → 503 `Busy` + `Retry-After: 1`; below v9 → 503 `Restart` | 503 |
| C14 | `ConsumeVerbRequest` (single use; approver re-checked; counter moved) | 409 race / 422 approver changed or counter |
| C15 | `broker.RecordOperatorApproval(ctx, d.store, r.Scope, "operator:episode-finish", signed.ApprovedBy(), now)` — **unchanged** (W6) | 500 (row stays consumed, no decision: RV2) |
| C16 | `LandVerbDecision(id, decision)`; close the writer | 500 (decision exists, row unlinked: RV2) |
| ok | 200 `surfacepage.Confirmed(id, approvedBy, decision)` = `{"v":1,"confirmed":true,"request":id,"approvedBy":…,"decision":"sha256:…","next":"the waiting episode finish opens the pull request"}` | |

The writer is held across one credentials transaction, one world write and one more credentials
transaction — A4's own bound (A4 §2.7). `TestDaemonCredentialWriterSites` gains
`handleSurfaceConfirm` (AM8). Placement: `host/daemon/surface_confirm.go` (NEW).

### 2.7 The CLI: `episode finish … --confirm-on-surface` (sprint 5; only under A or B)

**Flag** (NEW, `runEpisodeFinishWith`'s set, `episode_finish.go:60-71`): `--confirm-on-surface`
(bool). Help text gains: "confirm the pull request with an enrolled key on the World app's card
instead of a typed y; waits up to 10 minutes; the terminal is still required".

**Where it acts** (all NEW lines, in this order):
1. after F1 (`:85-89`): `--confirm-on-surface` with `--dry-run` → `fail("--confirm-on-surface cannot
   be combined with --dry-run", "a dry run opens no pull request", "drop one of them")`;
2. in the fresh branch, after `validateFinishChecks` (`:251`) and **before the check loop** (`:263`):
   any `check.Unsandboxed` → `fail("check <name> is unsandboxed: it needs its own typed y", "--confirm-on-surface
   replaces only the pull request's y", "run finish without --confirm-on-surface")` — no check has
   run, no outcome is recorded, no commit exists;
3. at the prompt (`:409-419`): the `else` branch becomes `if *onSurface { approval.Approval, code =
   finishConfirmOnSurface(…) } else { today's typed y }`. The reusable-approval branch (`:405-408`)
   is unchanged: a resume-dispatch with a reusable approval creates no request.

**`finishConfirmOnSurface`** (NEW, `cmd/ailang-worldd/episode_finish_surface.go`):
1. build `VerbRequest{Argv: canonicalFinishArgv(fs, allowed, ep), Episode: ep, Prompt:
   broker.FinishPromptText(scope, actualTitle, rendered, protected, pinned.Committer), Scope:
   scopeText, Title: actualTitle, Body: rendered, Committer: pinned.Committer, Protected:
   protected}`; `Validate`, `CanonicalVerb`; over `MaxVerbRequestBytes` → refuse "confirm at the
   terminal". **`canonicalFinishArgv`** (NEW, pure): `ailang-worldd episode finish <ep>`, then each
   flag the user set, in `fs.Visit`'s lexical order, `--name` for a true bool and `--name value`
   otherwise, `--allow-failing` once per value in given order — so interleaved spellings of one
   command give one argv, and defaults never appear;
2. id `vr-` + 8 bytes, challenge 32 bytes, both from `env.random` (refuse if nil, as
   `createMintRequest`, `mint_request.go:102-120`); `expires = env.now() + 600`;
3. one non-waiting writer (`store.OpenCredentials(credentials)`, the store F6 proved the daemon serves,
   W36), `PutVerbRequest`, close — under its own 3 s root (`approveReadBudget`, `mint_request.go:24`);
4. print the prompt to the tty (`FinishPromptText` bytes: the terminal shows what the card shows),
   and to stdout:
   ```
   confirm request vr-0123456789abcdef on the surface (expires 2026-11-17T10:10:00Z, message hash 0123456789abcdef):
     open http://localhost:7644/surface/confirm/vr-0123456789abcdef
     open https://studio.tail….ts.net/surface/confirm/vr-0123456789abcdef
   the typed y is not read in this mode; Ctrl-C stops waiting and keeps the commit at <sha>
   ```
   (one `open` per distinct active passkey origin plus the preset daemon's loopback, as
   `writeAwaitLinks`, `approve_await.go:29-58`);
5. **wait** (`awaitVerbConfirmation`, NEW): every `awaitPollInterval` through `env.after`, one
   read-only open under its own 3 s root reading the row; the **wait ends at `expires_at +
   verbLandGrace`** (`verbLandGrace = 10 s`, above the POST's 3 s `commitBudget`, W35). The tty is
   **never read** in this mode;
6. on `outcome = approved` with a `decision_ref`: read the decision and its request through
   `GET /v1/objects/{ref}?payload=true` (`daemon.go:1113`; the client's existing object read,
   `host/why/why.go:142`) and require `decision = approve`, `decidedBy == approved_by`, request
   `effect = Human.Approve`, `requester = operator:episode-finish`, `scope == scopeText`; print
   `confirmed by <approved_by> (decision <ref>)`; return the ref. The flow continues at `finishStep("approve")`
   (`:418`), then the unchanged F6 re-check (`:423-430`), dispatch and record.

**Who waits, and every way the wait ends** (the full table is §2.9 P1): approved → continue;
expired → "no confirmation of request <id> by <t>"; superseded → "superseded by a newer finish of
<ep>"; approved without a landed decision past the grace → "confirmed by <who> but its decision
did not land"; a decision that does not approve this scope → "(F6)"; a store read error →
"cannot read request <id>". Every refusal ends with "commit kept at <sha>" and `steps(sha)`, exits
1, and dispatches nothing. **Ctrl-C** ends the process; the request stays pending until its expiry
or the next run's supersession (RV1). **Daemon restart mid-wait**: the CLI polls the credentials
store, not the daemon, so its wait is unaffected; the card is unreachable while the daemon is down
and the wait may end at expiry; a POST cut by the restart either committed C14–C16 or left the row
pending (C14 is one transaction), or left it consumed without a decision (RV2).
`TestProductionContextRoots` gains the put and the poll roots (W23).

### 2.8 The bar (sprint 5)

- `episode finish …` stays `terminal` (the app runs no mutating verb before V4); under A/B the reason
  becomes "episode finish changes World. The app cannot run it yet; at a terminal,
  --confirm-on-surface lets you approve its pull request here with your key." (`lookup.go:90-92`
  split; `episode start` unchanged). The golden `verb-route.golden.json` row is regenerated.
- `approve vr-<16 hex>` → `navigate` `/surface/confirm/<id>` (`ApproveTarget`, `lookup.go:124-131`);
  the page's allowlist gains `/surface/confirm/vr-` (`app.js:496`).
- The CLI's `approve vr-…` → `exitUsage`, "request <id> is a finish confirmation: confirm it on its
  card, /surface/confirm/<id>; the waiting episode finish opens the pull request" (a third arm
  before `approve.go:228`), so the bar's grammar and the CLI's agree on the id shape.
- **Never**: the bar runs `episode finish`, POSTs a confirmation, or sends or accepts a `y`.

### 2.9 Input and failure classes, traced

**P1 — the CLI (`--confirm-on-surface`)**

| input / failure | branch | the user sees | test |
|---|---|---|---|
| no tty | F1 (`:85`), unchanged | "refusing: no controlling terminal" | existing, unedited |
| with `--dry-run` | step 1 | combination refusal, exit 1, nothing read | CF1 |
| an unsandboxed pinned check (fresh) | step 2 | refusal naming the check; no outcome recorded | CF2 |
| mode `finished` | `:181-187`, unchanged | "already finished" | existing |
| resume-dispatch + reusable approval | `:405-408`, unchanged | "resuming under approval …"; no request | CF12 |
| request over 128 KiB | §2.7 1 | "confirm at the terminal"; commit kept | CF9 |
| `env.random` nil / failing | §2.7 2 | internal error; commit kept | CF9 |
| writer busy | §2.7 3 | "the credentials store is busy"; commit kept | CF8 |
| store below v9 | §2.7 3 | the `OlderCredentialsStoreError` text naming the restart; commit kept | CF8 |
| 16 pending | §2.7 3 | `ErrTooManyPendingVerbRequests`; commit kept | CF8 |
| card confirms | §2.7 5–6 | "confirmed by webauthn:<name>"; then the PR URL | CF3, XC3 |
| decision does not approve this scope | §2.7 6 | "(F6)"; nothing dispatched | CF7 |
| approved, no decision after grace | §2.7 5 | "did not land"; commit kept | CF6 |
| superseded by a newer run | §2.7 5 | "superseded"; commit kept | CF5 |
| expiry (+ grace) | §2.7 5 | "no confirmation … by <t>"; commit kept | CF4 |
| a poll read fails | §2.7 5 | "cannot read request <id>"; commit kept | CF4 |
| Ctrl-C | process ends | request pending until expiry or supersession | RV1 |
| re-run after any refusal | `resume`, new request | same prompt bytes, new id; old one superseded | CF10 |
| the typed path (no flag) | unchanged | today's prompt and `y` | all existing finish tests, unedited |

**P5 — `POST /surface/confirm/{id}`**: each row of §2.6's table is one branch; each has its own
test (CX0–CX16, §6). Two cross-cutting rows: a replayed artifact on the same request → C8 409
(consumed); the same artifact on another request → C11/C12 422 (another digest and challenge).

**P4 — the card**

| input / failure | page | test |
|---|---|---|
| data fetch: network error, non-2xx, timeout (10 s) | a refusal; no controls | CP1 |
| `SHA-256(bytes) ≠ digest` | fail closed | CP1 |
| unlisted / missing / mistyped field | "confirm at the terminal" | CP2 |
| hostile title, body, prompt, argv | rendered as text | CP3 |
| state approved / superseded / expired | the state; no controls | CP5 |
| passkey ceremony aborted / 240 s | "not approved"; re-read (sprint 4) | CP9 |
| POST refusal (any C step) | class and reason; re-read | CP9 |

**P6 — the bar and the CLI's `approve`**: `episode finish w1 --preset p` → `terminal` (CB1);
`approve vr-0123456789abcdef` → `navigate` (CB2); `approve vr-0123` → refused, the `ApproveTarget`
reason (CB2); a `navigate` answer with `/surface/confirm/x` → "refused: unexpected target" (CB3);
CLI `approve vr-…` → the card pointer (CB4). **P2/P3** (read routes) are §2.5's rows, CR1–CR4.

## 3. Pure core, effects at the boundary (coding standards S1–S3)

- **S1.** No `world/` function changes; S1's Z3 clause binds `world/` AILANG and nothing there is
  touched. The new laws are pure Go: `host/approval` (stdlib-only, no I/O, no clock; `now` and
  the random bytes are inputs) and `broker.FinishConfirmRefusal`/`FinishPromptText` (no I/O),
  table- and mutation-tested (§6).
- **S2.** Effects stay at the host boundary: the CLI's randomness, store writes and polls; the
  daemon's credentials transactions and world write, under `humanMu`.
- **S3 — "why is this not a package?"** It is the host's authority boundary: verifying WebAuthn and
  SSHSIG against the credentials store and writing an approval into the world store are host
  effects that row 201 placed in `host/approval`, `host/authority` and `host/store` (A0–A4); an
  AILANG package can hold neither the credentials store nor the key judges. No `.ail` changes, no
  kernel change, no registry package; the `/v1` table is unchanged.

## 4. Threats (V3)

| threat | answer |
|---|---|
| T1 forged signature (a key never enrolled) | C11 `unknown-key`, 422; nothing written (CX11) |
| T2 replayed artifact, same request | C8 / C14 `Consumed`, 409; `consumed_at` is durable (CX8, CX14) |
| T3 replayed artifact, another request / episode / argv | the digest covers argv, episode, scope, title, body; the challenge is per request → C11/C12 422 (VL5–VL7, CX12) |
| T4 expired request | C8 (`expires_at`, also signed in the message); the CLI's wait ends (CF4) |
| T5 a mint or decision signature used to confirm, or a verb signature to mint/decide | prefix byte 6 and three namespaces (VL4, both directions, `approvaltest` keys) |
| T6 a page submits `y` | there is no typed arm: `ParseApproveBody` is strict with exactly one key format (CX6); the CLI never reads the tty in this mode (CF3 asserts zero reads); the bar never runs finish (CB1) |
| T7 CSRF, DNS rebinding, foreign Origin | R2 421 (`approvalRoute`), C3 Origin ∈ O, C4 `Sec-Fetch-Site`, C5 JSON (preflight; no CORS header is ever sent), C0 `Authorization`/`Cookie` 403 (CX0–CX5) |
| T8 a signature over bytes the human was not shown | the page renders every field of the signed bytes through one table and fails closed (CP2); it checks the digest itself and builds the message itself (CP1, CP4); the daemon refuses a prompt that is not the terminal's rendering of the approved fields (C10 F5) |
| T9 TOCTOU between confirmation and dispatch | the CLI re-reads the world decision and requires `scope == scopeText` (CF7); then the existing F6 re-check of HEAD, remote, title and body (`:423-430`); the record route re-checks the scope (`host/daemon/episode_finish.go:186-195`) |
| T10 daemon restart | the request and the consume are SQLite rows; no in-memory state carries single use; a cut POST leaves either nothing or a consumed row (RV2) |
| T11 revocation or a cloned passkey between judge and consume | `useApprover` and `usePasskey` inside C14's transaction (CS6) |
| T12 a key-confirmed finish approval used as a publish stamp | the finish scope never parses as a publish scope, and L2 refuses `ssh:`/`webauthn:` at publish (`approve.go:579`) |
| T13 an agent writes a verb request | it needs the credentials store's write, the same fence as the operator proof (A4 RK4); the human still sees the request on the card and must sign it; no session reaches the POST (C0) |
| T14 request flooding | 16 pending per store, 128 KiB per request, one pending per episode (supersede) (CS3–CS4) |
| T15 hostile text (title, body, argv) | `textContent` / `esc()` only (CP3, CP7) |

## 5. By-construction argument

1. **No `y` from a page**: the only new POST has no typed arm (C6), and the CLI does not read its tty
   in surface mode.
2. **What opens is what was approved**: the key signs `SHA-256(CanonicalVerb(r))`, whose `scope`
   is the exact string the daemon records (C15 passes `r.Scope`), and the CLI dispatches only if
   that recorded scope equals its own `scopeText` (CF7) and F6 holds.
3. **What was approved is what was shown**: the page shows every field of `r`; the prompt is a pure
   function of the other fields (C10 F5) and is the terminal's own bytes (VT1).
4. **Single use**: `consumed_at` in one transaction with the approver and counter checks (C14), and
   `RecordOperatorApproval` creates a fresh request object per call (W6).
5. **Typed `y` unchanged**: without the flag no new branch is reachable (§2.7 steps 1–3 are guarded
   by `*onSurface`); every existing finish test passes unedited.

## 6. Test plan

Each row names the mutant it kills (S6). Keys: `host/approval/approvaltest` and the A2 recorded
authenticator fixtures; Chrome: the A3/A4 driver and its skip rule (`approve_page_chrome_test.go:21`,
`approve_decide_chrome_test.go:25`, W41); GitHub: the `finishDispatchFactory` fake
(`episode_finish_dispatch.go:33`). Never a real enrolment, mint or publish (D-WORLD-109).

### Sprint 1 — the law
- **VB8** (V2b note a) `TestAdmitWhy` gains a row whose argument needs quoting (`--result "a b"`):
  the reason ends with `at a terminal: ` + `Quote(...)`'s form. Kills: `Quote` → `strings.Join` at
  `flags.go:137` (the controller's survivor, W27).
- **VP12** (V2b note b) the placeholder suffix (`app.js:1015`), the keys row's `why T`
  (`app.js:1303`) and `index.html:67`'s placeholder each asserted (node for the two `app.js`
  strings, Go for `index.html`). Kills: each string edited (W28).
- **VL1** `CanonicalVerb`/`DecodeVerb` round trip; refused: unknown field, trailing data, reordered
  keys, whitespace. Kills: `DisallowUnknownFields` removed; the canonical re-encode check removed.
- **VL2** `Validate`: one case per rule of §2.2. Kills: each rule deleted (one mutant per rule).
- **VL3** `VerbMessage` is 86 bytes in §2.2's layout; negative expiry refused. Kills: two fields
  swapped; a wrong length.
- **VL4** domain separation, both directions: a mint SSHSIG and a decision SSHSIG fail
  `ConfirmBySignature`; a verb SSHSIG fails `ApproveBySignature` and `DecideBySignature`; a mint
  assertion fails `ConfirmByAssertion` at W6 and back. Kills: `VerbNamespace = Namespace`;
  `VerbPrefix = MessagePrefix`.
- **VL5** `ConfirmBySignature`: P1 each verdict, P2–P5 (revoked, wrong format, clock behind), bad
  signature, wrong digest/challenge/expiry. Kills: P1 or `checkApprover` removed.
- **VL6** `ConfirmByAssertion`: W1–W14 via the recorded fixtures, W15 regression refused. Kills: W15
  removed; the challenge computed with `MessageHash`.
- **VL7** `SignedVerb.Binds` per field; the zero value refuses; `ApprovedBy` is "" unverified.
  Kills: a field comparison removed.
- **VJ1** `JudgeVerbSSH`/`JudgeVerbAssertion`: each `JudgeError.Step`. Kills: steps reordered.
- **VT1** for three fixtures (draft/ready, create-label, protected paths, empty committer)
  `FinishPromptText(...) + "Open this pull request? [y/N] "` equals the bytes `finishPrompt`
  writes; the existing `episode_finish_m2/m3` prompt tests stay unedited. Kills: the moved text
  edited; `DisplayCommitter`'s `""` rule dropped.
- **VT2** `FinishConfirmRefusal` F1–F5 (and F6 with `DraftOnly`). Kills: each step deleted.

### Sprint 2 — the store and the read routes
- **CS1** schema 9 DDL: fresh equals upgraded from 8, column for column; `user_version = 9`. Kills:
  the `< 9` upgrade line removed.
- **CS2** `PutVerbRequest` below v9 → `OlderCredentialsStoreError{Need: 9}`. Kills: the check removed.
- **CS3** supersede: two puts for one episode → the first `superseded`; another episode's row
  untouched. Kills: the supersede removed; its `WHERE episode_id` dropped.
- **CS4** grammar (id, hex, `digest == sha256(request)`, sizes), the 16 cap, the day-old prune.
  Kills: each deleted.
- **CS5** the table's CHECKs refuse a raw `approved` row without `approved_by`, and a
  `decision_ref` on a `superseded` row. Kills: a CHECK dropped.
- **CS6** `ConsumeVerbRequest`: consumed, superseded, clock behind, expired, digest changed, approver
  revoked since the judge, passkey counter not advancing → each refuses and writes nothing;
  success sets four columns. Kills: each check removed; `usePasskey` skipped.
- **CS7** `LandVerbDecision` lands once; a second land, or a land on a non-approved row →
  `ErrVerbDecisionLanded`. Kills: the `decision_ref IS NULL` guard removed.
- **CR1** `GET /surface/confirmations/{id}`: 400, 404, 421 (foreign Host), 503 (v8), the 200 shape
  with no hash and no artifact. Kills: registered on the bare mux; the artifact or a message hash
  served.
- **CR2** `requests.json` `confirmations`: pending listed; an undecodable row listed with `refused`;
  consumed and superseded omitted; v8 → `confirmationsUnavailable` with the other lists intact.
  Kills: the refused row hidden.

### Sprint 3 — the card, read-only
- **CP1** digest mismatch, fetch failure, non-2xx, timeout → refusal, no controls. Kills: the
  digest check removed; the abort removed.
- **CP2** `CONFIRM_FIELDS` equals `VerbRequest`'s json tags (a Go test reads both, as A4's PD1);
  unlisted, missing and mistyped fields refuse. Kills: the table loosened.
- **CP3** hostile values render as text. Kills: `innerHTML`.
- **CP4** the page's message hash equals Go's `VerbMessageHash` for a golden written by a Go test.
  Kills: a field order changed in JS.
- **CP5** non-pending states show no controls. Kills: the state check removed.
- **CP6** *Save the message* is the 86 bytes; the signing line names `world-verb@sunholo`. Kills:
  the namespace changed.
- **CP7** *Needs you* line and link; a refused row has no link; `esc()`. Kills: the link on a refused row.
- **XC1** live Chrome: a seeded request's card shows the prompt equal to `FinishPromptText` of the
  fixture.

### Sprint 4 — the POST (waits for D-WORLD-118)
- **CX0–CX16** one test per §2.6 step (C0 `Authorization`, C0 `Cookie`, C1, R2, C3 missing/`null`/foreign
  Origin, C4, C5, C6 oversize/unknown key/`typed`/both formats, C7, C8 consumed/superseded/expired/
  clock/digest, C9, C10 F1–F5, C11 unknown/revoked/namespace/W6/W15, C12 via a judge seam, C13 busy
  and v8, C14 a concurrent consume lost, C15 a failing world-store wrapper → 500 and the row consumed
  with no decision, C16 a failing land → 500). Each kills its step's deletion.
- **CX17** success: 200; the world decision has `decidedBy` = `ssh:<test>` / `webauthn:<test>`,
  requester `operator:episode-finish`, scope = the request's; the row's `decision_ref` set; a second
  POST of the same artifact → 409. Kills: C15 passing `"operator:terminal"`; the land skipped.
- **CX18** `TestDaemonCredentialWriterSites` pins four sites; A4's `POST /surface/decide/{ref}` still
  refuses a finish-shaped request (D5), unedited. Kills: an extra writer site.
- **CX19** `humanMu`: a confirmation and a typed finish approve phase on one daemon serialize (the
  `surface_decide_race_test.go` shape). Kills: the lock removed.
- **CP8–CP10** `confirm.js` submit: passkey `challenge` = the hash, `allowCredentials` from groups;
  the SSHSIG paste posts the body grammar; answers and refusals rendered, then re-read. Kills: the
  challenge from the mint hash; a refusal shown as success.
- **XC2** live Chrome: the virtual authenticator confirms a seeded request; the world decision exists.

### Sprint 5 — the CLI and the bar (waits for D-WORLD-118)
- **CF1–CF14** §2.9 P1's rows, on the finish rig with a real loopback daemon; the test plays the card
  by POSTing an `approvaltest` SSHSIG. CF3 also asserts the tty reader is never read. CF13: the
  stored `prompt` equals what the typed path prints before its question. CF14: interleaved flag
  spellings give one canonical argv; defaults absent. Kills: the CLI reading `y`; the decision
  check removed (CF7); the unsandboxed refusal moved after the loop (CF2 sees an outcome recorded);
  an unbounded wait (CF4 under a fake clock and a test deadline).
- **CB1–CB4** §2.9 P6. Kills: the allowlist not extended / extended to any `/surface/confirm/`.
- **CQ1** QUICKSTART §20 is executed by the `quickstart_verbs_test.go` pattern (S7).
- **XC3** live Chrome end to end: the real CLI `episode finish … --confirm-on-surface` (fake gh) waits;
  the page confirms with the virtual authenticator; the CLI prints the fake PR URL.

## 7. Inventories and conflict surface

- **New**: `host/approval/verb.go`, `verb_signed.go` (+ tests); `host/authority/judge_verb.go`
  (+ test); `host/store/verb_requests.go` (+ test); `host/daemon/surface_confirm.go`,
  `surface_confirm_read.go` (+ tests); `host/surfacepage/confirm.go`, `confirm_model.go`,
  `confirm/{index.html,confirm.js}` (+ tests); `cmd/ailang-worldd/episode_finish_surface.go` (+ tests,
  XC3).
- **Edited**: `host/broker/approve.go` (`DisplayCommitter`, `FinishPromptText`,
  `FinishConfirmRefusal`; or a new `finish_prompt.go`), `host/store/credentials.go` and
  `credentials_schema.sql` (v9), `host/daemon/daemon.go` (four registrations),
  `surface_approve.go` (`requests.json`), `host/surfacepage/approve_model.go` (`Confirmations`),
  `app/app.js` (*Needs you*, allowlist), `cmd/ailang-worldd/episode_finish.go` (flag and three hooks),
  `episode_finish_operator.go` (`finishPrompt` delegates), `committer.go` (delegates), `approve.go`
  (`vr-` arm), `host/verbs/lookup.go` (reason, `ApproveTarget`), `docs/QUICKSTART.md` (§20).
- **Pre-existing tests edited (each with its reason)**: `surface_approve_post_guard_test.go` (fourth
  writer site, AM8); `host/store/context_roots_test.go` (two roots); the canonical-DDL test (v9);
  `lookup_test.go` and `verb-route.golden.json` (the finish reason, `vr-`); `TestAdmitWhy` and the
  page tests (sprint 1 pins); `TestSurfaceAppPure`'s export list if a pure helper is exported.
  Nothing else; in particular every `episode_finish*_test.go` stays unedited.
- **Not touched**: A4's law and routes (L1, `reservedReason`), `approve.js`, `decide.js`,
  `/surface/passkey`, the mint law, `world-publish`, `/v1`, `approvalRoute`, the CSP, any `.ail`,
  `tools/launchd/*`.

## 8. Sprints and sizes

Raw estimates; each sprint ≤ 1,000 raw so that × 1.35 (W31) ≤ 1,400.

| sprint | milestones | est. raw | depends on D-WORLD-118 |
|---|---|---|---|
| **S1 the law** | M0 V2b pins (≈ 45); M1 `host/approval` verb request, message, laws (≈ 240 + 320 tests); M2 judges (≈ 70 + 110); M3 `FinishPromptText`, `DisplayCommitter`, `FinishConfirmRefusal`, cmd delegations (≈ 75 + 90) | ≈ 950 | **no** |
| **S2 the store and reads** | M4 schema 9 + `verb_requests.go` (≈ 270 + 330); M5 `GET /surface/confirmations/{id}`, `requests.json` (≈ 150 + 240) | ≈ 990 | **no** |
| **S3 the card, read-only** | M6 page + `confirm.js` render/verify/message (≈ 310 + 220); M7 *Needs you* (≈ 25 + 60); M8 XC1 (≈ 110) | ≈ 725 | **no** |
| **S4 the POST** | M9 `surface_confirm.go` + `Confirmed` + registration (≈ 215 + 470); M10 submit half of `confirm.js` (≈ 90 + 90); M11 XC2 (≈ 120) | ≈ 985 | **yes** (A or B; B +≈ 15) |
| **S5 the CLI and the bar** | M12 flag, refusals, `episode_finish_surface.go`, hooks, roots (≈ 260 + 400); M13 bar, `approve vr-` (≈ 30 + 60); M14 QUICKSTART §20 + CQ1 (≈ 90), XC3 (≈ 150) | ≈ 990 | **yes** (A or B; B +≈ 15) |

**Total ≈ 4,640 raw, ≈ 6,000–6,300 calibrated.** Against the parent's ≈ 4,500 (2–3 sprints): the
raw sum is within 3 %; the difference is the measured over-run of the last two V-sprints (1.34×,
1.35×, W31) and the 1,400 cap. **Cuts** if a sprint measures over: S3 → XC1 moves to S4; S4 → XC2
moves to S5; S5 → XC3 moves to a follow-up sprint (QUICKSTART §20 stays: S7). Under **C**, S4–S5
are not built and S1–S3 are removed in one cleanup (≈ 2,700 lines of unused code, §10).

**Gates** (every sprint): `./scripts/verify_ail.sh` byte-for-byte the base's result (no `.ail`
changes); `go build ./... && go vet ./... && go test ./...`, with `go test -run '^$' ./...` as the
test-file compile fence; node tests non-skippable in CI (`WORLD_EXEC_NODE`).

## 9. The V2b judge notes, disposed

Both are taken in **sprint 1, M0** (tests only, no production change), because they are independent
of D-WORLD-118, small (≈ 45 lines), and sprint 1 is the first sprint that cannot be blocked by the
ledger: (a) `AdmitWhy`'s `Quote` in the refusal's CLI line — today every refusal row's arguments
need no quoting (`why_admit_test.go:83-88`), which is why a plain join survives (W27): VB8; (b) the
three page strings — no test or golden contains them (W28): VP12.

## 10. The ledger ask: D-WORLD-118 (row text for the controller to file)

> | D-WORLD-118 | OPEN | **Row 243 V3: may `episode finish`'s pull-request confirmation be a verified enrolled key, as D-WORLD-89 R1 allows for mints and `Human.Approve`?** Today `ailang-worldd episode finish` commits, then asks `Open this pull request? [y/N]` at its terminal, and the typed `y` is the only way a finish opens a PR (`cmd/ailang-worldd/episode_finish.go:410-414`); row 201 A4 §2.6 L1 keeps finish-shaped requests off the decision card. The finish approval is already stored as a `Human.Approve` request (`host/broker/approve.go:722-724`), so R1's words could be read to reach it, but A4 read them not to, and the loop asks rather than infers. The V3 design ([`w-world-address-bar-verbs-v3.md`](planned/w-world-address-bar-verbs-v3.md) r0, §2, §10) adds `episode finish … --confirm-on-surface`: after the commit, the CLI stores a `world/verb-request/v1` (the exact prompt lines, the canonical command, the finish scope, the title and the full PR body), prints its card's link and waits up to 10 minutes at the terminal; you approve on `/surface/confirm/<id>` with your passkey or SSH key; the daemon verifies the signature and records the finish approval with `decidedBy` = your key (`webauthn:<name>` / `ssh:<name>`); the CLI then opens the PR. The typed `y` keeps working without the flag; publishing stays typed-only (R4); enrolment stays typed (R3); the app's address bar still runs no mutating verb until V4. **A (recommended):** yes, for every finish, draft or ready — a key already mints sessions (R1), a broader authority than opening a pull request, and the key binds the exact bytes the terminal shows. **B:** yes for draft pull requests only; `--ready` stays typed-only (≈ +30 lines). **C:** no; finish stays typed-only, V3 stops, its sprints 1–3 are removed, and row 243's no-terminal exit run (parent §10 step 4) needs a terminal pane for finish (V4 hands it to herdr). **Default while unanswered:** nothing that confirms a finish from the surface is built or enabled, and the typed `y` stays the only confirmation: the loop builds V3 sprints 1–3 (the law, the store, a read-only card that can show a request but not confirm it — identical under A and B, removed under C) and does not plan sprints 4–5 (the POST that records a key's confirmation, the CLI flag). | Filed iteration 399 (unattended loop). Evidence: the V3 design's §0.1, §0.3 (AM6/AM7 amend A4 §2.6's stated reasons), §2.1, §2.6–§2.7 and its Verification Log. |

**Which sprints depend on the answer:** S4 and S5 wait for A or B (B adds the F6 draft check to S4
and a `--ready` refusal to S5, ≈ 15 lines each). S1–S3 do not: they are the same under A and B,
change nothing from a page, and are removed under C.

## 11. Residuals and risks

- **RV1 Ctrl-C leaves a pending request.** It can still be confirmed until its expiry; the daemon
  then records an approval no CLI uses. It is inert: `finishReusableApproval` reuses only a recorded
  finish (`episode_finish_resume.go:119`, W11), the scope expires ≤ 15 minutes after creation, and the
  next run supersedes the request. A SIGINT handler that marks it cancelled is not designed.
- **RV2 a POST cut after C14.** The request is spent (consumed, `decision_ref` NULL, or a decision
  with an unlinked row). The CLI refuses after the 10 s grace with "commit kept"; a re-run resumes
  and makes a new request. A4's idempotent same-artifact retry (K0) is deliberately not copied: a
  new request costs one re-run and one gesture.
- **RV3 unsandboxed checks refuse surface mode** (P2). Whether Mark's real presets carry any is
  UNVERIFIED (no presets exist on this machine, W42). If they do, the §10 exit run needs V4 to design
  per-check confirmation, or the checks to be sandboxed.
- **RV4 F1 is kept.** V3's finish still needs a terminal; V4 decides who runs it without one.
- **RV5 no deny on the card.** Decline is expiry or Ctrl-C. An unsigned deny would let any local page
  cancel a finish; a signed one needs a verb byte (A4's) and a new message version.
- **RV6 the request bytes on a GET.** The full PR body (≤ 60,000 bytes) and local paths in `argv`
  (`--presets-dir` if given) are served behind the Host allowlist, as A3a serves a mint request's
  paths; the body is the text that goes public with the PR.
- **RV7** schema 9 costs one `serve` restart (A2's R6 recurs); a CLI against a v8 store refuses with
  the restart named (CF8).
- **RV8** the finish record names the approval by ref only; who confirmed is in the decision
  (`decidedBy`), not in `episode report`. V3 does not change the report.
- **RV9** sizes: five sprints rest on the 1.35× of two sprints; §8 names the cuts.

## 12. Open questions (each with the default the loop applies)

- **Q1. Relax F1 in V3 for surface mode?** Default **no** (V4; RV4).
- **Q2. A deny on the card?** Default **no** (RV5).
- **Q3. Per-check key confirmation of unsandboxed checks?** Default **no**: refuse (RV3); V4 or a
  later row.
- **Q4. Ledger ask?** **Yes**: D-WORLD-118 (§10).
- **Q5. Five sprints, or merge S2/S3?** Default **five**; a planner whose compiled patches measure
  ≤ 1,400 together may ask the controller to merge S2 and S3.

## 13. What this is NOT doing

- Not running any verb from the bar; not adding a `y` anywhere on a page; not relaxing finish's F1.
- Not changing `episode start`, `approve --await`, the mint law, A4's law, L1's code, `world-publish`,
  `/v1`, `approvalRoute`, the CSP, or any `.ail`.
- Not adding a deny, a SIGINT cancel, or per-check confirmation.
- Not designing V4 (the runner, launch presets, herdr) or V5.

## 14. Verification Log

All commands ran in `/Users/voightkampff/dev/sunholo-data/.wt-world-iter399-row243-v3` at `994ba123`
on 2026-11-17 (zsh; output trimmed); the batch is `/tmp/w399-designer/wlog.sh`. `F` =
`cmd/ailang-worldd`.

| # | claim | command | observed |
|---|---|---|---|
| W1 | base; cmd/host unchanged since V2b (negative, with control) | `git rev-parse HEAD`; `git diff --stat c251132f HEAD -- cmd host \| tail -1`; control `… -- design_docs \| tail -1` | `994ba12358dc50246db2862d78f66d64e03b29f4`; *(empty)*; control `20 files changed, 974 insertions(+), 143 deletions(-)` |
| W2 | F1 refuses without a tty first | `sed -n 85,88p $F/episode_finish.go` | `tty, err := env.openTerminal()` … `"refusing: no controlling terminal", "finish is an attended act"` |
| W3 | commit, then the prompt; decline keeps the commit | `grep -n 'CommitWithCommitter(ctx\|"committed at %s\|finishPrompt(tty\|commit kept at "+sha' $F/episode_finish.go` | `348: sha, err = worktree.CommitWithCommitter(…)`; `377`; `380: "committed at %s\nbody: %s\n"`; `410: finishPrompt(tty, scope, …)`; `413: "pull request declined", "…; commit kept at "+sha` |
| W4 | unsandboxed checks ask their own `y` | `grep -n 'if check.Unsandboxed\|Run this check? \[y/N\]\|unsandboxed code requires' $F/episode_finish.go` | `266`, `270: "Run this check? [y/N] "`, `273: "unsandboxed code requires its own typed y"` |
| W5 | typed approve phase and the proof | `grep -n 'Phase: "approve", Scope: scopeText\|Phase: "record", Record' $F/episode_finish.go`; `grep -n 'RecordOperatorApproval(ctx, d.store' host/daemon/episode_finish.go`; `grep -n 'proof, err := writeOperatorProof' $F/episode_finish_operator.go` | `415: …{Phase: "approve", Scope: scopeText}…`; `443: …{Phase: "record", Record: &finish}…`; `173: broker.RecordOperatorApproval(ctx, d.store, input.Scope, "operator:episode-finish", "operator:terminal", now)`; `76` |
| W6 | `RecordOperatorApproval` takes `decidedBy`; a `Human.Approve` request | `sed -n 722,725p host/broker/approve.go`; `sed -n 732p …` | `func RecordOperatorApproval(ctx, s, scope, requester, decidedBy string, now int64)`; `Effect: EffectHumanApprove, Scope: scope, Cost: 1, Requester: requester, Now: now`; `decideApproval(ctx, s, obj.Hash, "approve", decidedBy, now)` |
| W7 | no pinned commit date (negative, with control) | `grep -rn 'GIT_COMMITTER_DATE\|GIT_AUTHOR_DATE\|--date' host/worktree $F --include='*.go' \| grep -v _test \| wc -l`; control `grep -n '"commit", "--allow-empty"' host/worktree/commit.go` | `0`; control `48: args = append(args, "commit", "--allow-empty", "--cleanup=verbatim", "--author", …)` |
| W8 | the body embeds the sha | `grep -n 'add.Commit = sha\|rendered, renderErr = renderReportMarkdown' $F/episode_finish.go`; `grep -n 'Commit: %s' $F/episode_markdown.go` | `352: add.Commit = sha`; `354`; `78: "\nCommit: %s\n\nD-WORLD-60: …", add.Commit` |
| W9 | dry-run records outcomes and commits | `grep -n 'entry, err := record(ctx, c, token, taskOutcomeObject\|if \*dry \|\| renderErr != nil\|"reset", "-q", "--mixed", previous' $F/episode_finish.go` | `290: entry, err := record(ctx, c, token, taskOutcomeObject(canonical), …)`; `355: if *dry \|\| renderErr != nil {`; `358: query("reset", "-q", "--mixed", previous)` |
| W10 | scope +900 s; route bound 3600 s | `grep -n 'ExpiresAt: time.Now().Unix() + 900' $F/episode_finish.go`; `grep -n 'scope.ExpiresAt > now+3600' host/daemon/episode_finish.go` | `394`; `169` (approve phase), `276` (record phase) |
| W11 | resume modes; reuse only in resume-dispatch | `grep -n 'return "resume-dispatch"\|return "resume", trailer\|if mode != "resume-dispatch"' $F/episode_finish_resume.go` | `38`, `41`, `119` |
| W12 | the prompt's text | `grep -n 'func finishPrompt\|Open this pull request? \[y/N\]' $F/episode_finish_operator.go` | `154: func finishPrompt(tty io.Writer, scope broker.FinishApprovalScope, title, body string, protected []string, committer ...string)`; `170: "repository: %s\nbranch: %s\nsha: %s\nbase: %s\ntitle: %s\nbody (%d bytes), first lines:\n%s\nlabel: %s\n%s\nprotected paths changed: %s\nOpen this pull request? [y/N] "` |
| W13 | L1 reserves finish scopes | `grep -n 'ReservedApprovalScope(f.Scope); reserved\|func ReservedApprovalScope\|HasPrefix(scope, "github-pr:")\|finish approvals are made by' host/broker/decision_law.go` | `80`, `116`, `121`, `131: "finish approvals are made by `ailang-worldd episode finish`; this request cannot be decided here"` |
| W14 | A4's message and judge input | `grep -n 'DecisionPrefix = \|DecisionNamespace = \|DecisionMessageLen = ' host/approval/decision.go`; `grep -n 'Packet, Request \[32\]byte' host/authority/judge_decision.go` | `15: "world-decision/v1"`, `18: "world-decision@sunholo"`, `21: … + 1 + 1 + 32 + 32 + 8`; `19: Packet, Request [32]byte` |
| W15 | A3a's judges; the mint message; lifetime | `grep -n '^func JudgeAssertion\|^func JudgeSSHSignature' host/authority/judge.go`; `grep -n 'MessagePrefix = \|RequestLifetime = \|make(\[\]byte, 0, 90)' host/approval/request.go` | `65: … row store.MintRequestRow, response []byte, now int64)`, `102: … row store.MintRequestRow, armored …`; `26: "world-approval/v1"`, `29: 10 * time.Minute`, `215` |
| W16 | shared law parts | `grep -n '^func assertionSteps\|^func checkApprover\|// W15' host/approval/passkey_assert.go host/approval/signed.go` | `passkey_assert.go:46: … ad.SignCount <= p.SignCount { // W15`; `:60: func assertionSteps(response []byte, a Approver, h [32]byte, hErr error)`; `signed.go:75: func checkApprover(…)` |
| W17 | schema 8; no v9 (negative, with control) | `grep -n 'const credentialsSchemaVersion' host/store/credentials.go`; `grep -n 'if version < ' … \| tail -1`; `grep -c 'version < 9\|verb_requests' host/store/credentials.go host/store/credentials_schema.sql`; control `grep -c key_decisions host/store/credentials_schema.sql` | `44: … = 8`; `375: if version < 8 {`; `0`, `0`; control `1` |
| W18 | both request tables CHECK `mr-` | `grep -n "request_id  *TEXT PRIMARY KEY" host/store/credentials_schema.sql` | `45` and `92`: `… request_id GLOB 'mr-[0-9a-f]*'` |
| W19 | A3a POST order | `grep -n '// R10\|// R11\|// R12\|// R13\|// R14' host/daemon/surface_approve.go` | `257 R10`, `267 R11: the law, before anything is recorded`, `282 R12`, `286 R13: store.OpenCredentials … non-waiting, closed at once`, `292 R14` |
| W20 | A4 POST: R0, `humanMu`, law, writer | `sed -n 158,200p host/daemon/surface_decide.go`; `grep -n 'store.OpenCredentials(d.credentialsPath)' host/daemon/surface_decide_key.go` | `159: if len(r.Header.Values("Authorization")) != 0 \|\| len(r.Header.Values("Cookie")) != 0 { // R0`; `186: d.humanMu.Lock()`; `198: broker.DecisionRefusal(f) … // 1. the law`; `66: store.OpenCredentials(…) // non-waiting` |
| W21 | writer sites pinned | `grep -n 'if want := map\[string\]string{"ensureCredentialsStore"' host/daemon/surface_approve_post_guard_test.go` | `168: {"ensureCredentialsStore": "UpgradeCredentials", "handleSurfaceApprove": "OpenCredentials", "keyArm": "OpenCredentials"}` |
| W22 | the minter's wait | `grep -n 'const awaitPollInterval\|deadline := row.ExpiresAt\|<-env.after(awaitPollInterval)\|func pollSurfaceApproval' $F/approve_await.go` | `24: … = time.Second`; `70`; `105: deadline := row.ExpiresAt`; `130: <-env.after(awaitPollInterval)` |
| W23 | context roots pinned | `grep -n 'mint_request.go\|createMintRequest\|approve_await.go\|pollSurfaceApproval' host/store/context_roots_test.go` | `19: "cmd/ailang-worldd/mint_request.go\|createMintRequest\|Background": 1`; `30: "…/approve_await.go\|pollSurfaceApproval\|Background": 1` |
| W24 | `--request` refuses a launch preset | `grep -n 'a request carries no launch callbacks\|case \*preset != "" && \*preset != "se-tools":' $F/episode.go` | `241: why := "a request carries no launch callbacks or start records (row 201 §2.5)"`; `249` |
| W25 | a request has no start record | `grep -n '"start":  ' $F/mint_request.go` | `219: "start": "nil: set only by episode_preset.go; --request refuses launch presets"` |
| W26 | the bar's finish and approve rows | `grep -n 'case "start", "finish":\|The app cannot run it yet."}\|mintRequestID = \|decisionRef   = \|return "/surface/approve/"\|return "/surface/decide/"' host/verbs/lookup.go` | `90`, `92: … "changes World. The app cannot run it yet."`, `95`, `120: ^mr-[0-9a-f]{16}$`, `121: ^sha256:[0-9a-f]{64}$`, `127`, `130` |
| W27 | V2b note (a) | `grep -n 'cli := Quote(append' host/verbs/flags.go`; `sed -n 1776p design_docs/world-mission-log.md`; `grep -n 'at +' host/verbs/why_admit_test.go` | `137: cli := Quote(append([]string{Binary, "why"}, args...))`; `(1) AdmitWhy's Quote … unpinned — a plain join survives the whole host/verbs package; reproduced by the controller (mutant host/verbs/flags.go:137 …)`; rows `83`–`88`: args `-`, `- --scan 3`, `--result f`, `--expect-store p head`, … (none needs quoting) |
| W28 | V2b note (b): the strings, unpinned (negative, with control) | `grep -n` each string in `index.html`/`app.js`; `grep -rlE "why &lt;target&gt;\|why <target>'\|why T, approve ID" host/surfacepage/*_test.go host/surfacepage/testdata \| wc -l`; control `grep -rl "did not answer in time" host/surfacepage/*_test.go` | `index.html:67 … · why &lt;target&gt;">`; `app.js:1015 … + ' · episode report <episode> · why <target>'`; `app.js:1303 ['Enter', 'in the bar: a place, episode report E, why T, approve ID', 'the same']`; `0`; control `host/surfacepage/app_verb_test.go` |
| W29 | CSP | `grep -n "connect-src 'self'" host/daemon/surface_page.go` | `19: const surfaceCSP = "default-src 'none'; … connect-src 'self'; …"` |
| W30 | the `/surface` approval routes | `sed -n 1139,1148p host/daemon/daemon.go \| grep -o '"[A-Z]* /surface/[^"]*"'` | `GET /surface/requests.json`, `GET /surface/requests/{id}`, `POST /surface/approve/{id}`, `GET /surface/approve/{id}`, `GET /surface/approve.js`, `GET /surface/decisions/{ref}`, `POST /surface/decide/{ref}`, `GET /surface/decide/{ref}`, `GET /surface/decide.js`, `GET /surface/verb.json` |
| W31 | calibration | `for m in fb4d8b89 2d530cbe c251132f b44b5f7d d1027477 70823879; do git diff --numstat $m^1 $m -- cmd host docs \| awk '{a+=$1;d+=$2} END{print a,d,a+d}'; done` | V1 `1759 150 1909` (plan 1,750); V2a `924 375 1299` (plan ≈ 970); V2b `1186 69 1255` (plan ≈ 930); A3a s1 `1543`, s2 `1995`, s3 `1861` |
| W32 | D-WORLD-118 is free; nothing OPEN (negative, with control) | `grep -c 'D-WORLD-118' design_docs/world-mission.md`; control `grep -c 'D-WORLD-117' …`; `grep -c '\| OPEN \|' …`; control `grep -c '\| RESOLVED \|' …` | `0`; control `3`; `0`; control `104` |
| W33 | R1's words | `sed -n 1008p design_docs/world-mission.md` | `R1 = YES: for mints and Human.Approve decisions, a verified signature from an enrolled key may approve, as well as the typed y` |
| W34 | body cap | `grep -n 'const reportMarkdownMax' $F/episode_markdown.go` | `11: const reportMarkdownMax = 60000` |
| W35 | commit budget | `grep -n 'commitBudget = ' host/daemon/daemon.go` | `150: commitBudget = 3 * time.Second` |
| W36 | F6: the daemon serves the preset's credentials | `sed -n 114,115p $F/episode_finish.go` | `if same, _ := store.SameDatabase(h.CredentialsDB, credentials); !same { … "the operator proof must be beside the served credentials store"` |
| W37 | `decidedBy` on the wire; L2 | `grep -n 'DecidedBy  string' host/broker/approve.go`; `grep -n 'HasPrefix(decision.DecidedBy, "ssh:")' …` | `79: DecidedBy string \`json:"decidedBy"\``; `579` |
| W38 | the post-approval re-checks | `grep -n 'record scope differs from approval\|record does not match an approved finish scope' host/daemon/episode_finish.go`; `grep -n 'the commit or body changed since approval (F6)' $F/episode_finish.go` | `188`, `192`; `429` |
| W39 | the mint cap and its text | `grep -n 'const MaxPendingMintRequests\|ErrTooManyPendingRequests = ' host/store/mint_requests.go` | `39: … = 16`; `49: "store: 16 mint requests are already pending; …"` |
| W40 | `approve`'s id dispatch | `grep -n 'if decisionRefPattern.MatchString(id) {' $F/approve.go` | `228` |
| W41 | the live-Chrome drivers | `grep -n '^func Test' $F/approve_page_chrome_test.go $F/approve_decide_chrome_test.go $F/verb_bar_chrome_test.go` | `TestApprovePageInChrome` `:21`; `TestApproveDecideWebAuthnLinkInChrome` `:25`; `TestVerbBarEpisodeReportInChrome` `:22` |
| W42 | presets on this machine (UNVERIFIED for Mark's) | `ls ~/.ailang/world/presets/` | `No such file or directory` |
| W43 | finish's two typed prompts; the prompt's pins | `grep -n '\[y/N\]' $F/episode_finish.go $F/episode_finish_operator.go`; `grep -rln 'Open this pull request' cmd host --include='*_test.go'` | `episode_finish.go:270`, `episode_finish_operator.go:170`; `episode_finish_m3_test.go` (and `committer:` in `episode_finish_m2_test.go:310`) |
| W44 | the page's navigate allowlist | `grep -n "b.href.startsWith('/surface/approve/mr-')" host/surfacepage/app/app.js` | `496: … startsWith('/surface/approve/mr-') \|\| … startsWith('/surface/decide/sha256:') ? { assign: b.href } : { say: 'refused: unexpected target' }` |
| W45 | `displayCommitter` | `grep -n 'func displayCommitter' $F/committer.go`; `grep -rn 'displayCommitter(' $F/*.go \| grep -v _test \| wc -l` | `3`; `3` |
| W46 | host/verbs' pure import pin | `grep -n 'pure := map' host/daemon/surface_verb_test.go` | `151: {"errors", "flag", "fmt", "io", "regexp", "strings", "time", "unicode/utf8"}` |
| W47 | `Admit` is generic over bytes | `sed -n 37,58p host/approval/admit.go` | `type AdmitInput struct { Now, CreatedAt, ExpiresAt int64; Consumed bool; StoredDigest, BytesDigest [32]byte }`; `func Admit(in AdmitInput) Verdict` consumed → clock → expired → digest |
| W48 | consume helpers | `grep -n 'func useApprover\|func usePasskey\|type ApproverUse\|type PasskeyUse' host/store/*.go` | `approvers.go:245`, `passkeys.go:225`, `mint_requests.go:300`, `passkeys.go:216` |
| W49 | A3a page embeds and body grammar | `grep -n 'go:embed' host/surfacepage/approve.go`; `grep -n '^func ParseApproveBody\|MaxApproveBodyBytes = ' host/surfacepage/approve_model.go` | `13: approve/index.html`, `18: approve/approve.js`; `197`; `27: 20 << 10` |
| W50 | *Needs you* and its model | `grep -n 'Needs you\|function needsYouHTML' host/surfacepage/app/app.js`; `grep -n 'Decisions\|Pending \|Decided ' host/surfacepage/approve_model.go` | `424`, `429`, `448: /surface/decide/sha256:…`; `67 Pending`, `68 Decided`, `71 Decisions`, `72 DecisionsTruncated` |
| W51 | the finish route refuses sessions; consumes the proof under `humanMu` | `sed -n 135,163p host/daemon/episode_finish.go` | `137: "… next: run ailang-worldd episode finish at a terminal"`; `158: d.humanMu.Lock()`; `160: d.consumeOperatorProof(episode, input.Phase, …)` |
| W52 | the finish scope's frozen format | `sed -n 689,701p host/broker/approve.go` | `type FinishApprovalScope struct { Repository, Remote, Head, SHA, Base, TitleSHA256, BodySHA256 string; Draft, CreateLabel bool; ExpiresAt int64 }`; `"github-pr:%s#remote=%s&head=%s&sha=%s&base=%s&title-sha256=%s&body-sha256=%s&label=ailang-world&draft=%t&create-label=%t&expires=%d"` |
