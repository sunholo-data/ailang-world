# w-standing-session-allowances (queue row 224): standing per-preset session allowances — r0

**Status:** r0, written 2026-11-24 against `HEAD` (`f074dd4`). Governed by ratified ruling
**D-WORLD-96** (`design_docs/world-mission.md:1015`; `design_docs/planned/w-all-agent-work-through-world.md`
§4). Roadmap: [`w-all-agent-work-through-world`](w-all-agent-work-through-world.md) §2 row 6b,
§3 step 4 (after row 201). Charter: `design_docs/world-mission.md:6314`.

Slice 1 (pure allowance domain, canonical digest, domain-separated signature verification, and
draw-admission law) is implemented in `host/approval/allowance.go` alongside this r0 design doc.
Store schema v11, authority wiring, CLI verbs, and `/surface` endpoints (Slices 2–4) wait on the
open rulings in §4 (R1–R6).

---

## §0 What is wrong today

Row 224's own words: **"EVERY SESSION NEEDS ITS OWN MINT."**

- Even with launch presets (row 217) and key/passkey approvals (row 201), every `episode start
  --preset <name>` still blocks on a human act: either `y/N` at `/dev/tty`
  (`cmd/ailang-worldd/session_new.go:409-469`) or a per-request key/passkey approval
  (`host/authority/mint.go:100-175`).
- Host adapters (`integrations/herdr/ailang-world/bin/world-agent-popup`, `hosts init` launchers)
  cannot launch a confined agent pane without prompting Mark every time.
- **Ratified bounds (D-WORLD-96, 2026-10-27):**
  - Per preset, at most **10 sessions a day**, with the preset's grants and a **12 h TTL**
    (43,200 s).
  - Restricted to the **repos listed in the preset**.
  - Approval lasts **30 days** (2,592,000 s), **revocable at any time** from the CLI or the app.
  - Exceeding any bound **refuses** and requires an explicit approval (tap or typed `y`).
  - The **mission loop (row 196) is excluded** unless ruled separately.
- **Exit criterion (row 224):** a week of launches with no per-session prompt, every session
  traceable to its allowance and approval.

---

## §1 Premise verification log

Every premise below was **measured** directly in this worktree at base `f074dd4` (2026-11-24).

| # | premise | measured source | status |
|---|---|---|---|
| P1 | D-WORLD-96 bounds: per preset ≤ 10 sessions/day, preset's grants, 12 h TTL, preset's repos, 30-day approval validity, revocable any time from app or CLI, over-bounds asks for a tap, mission loop (row 196) excluded. | `design_docs/world-mission.md:1015, 6314`; `design_docs/planned/w-all-agent-work-through-world.md:142-148` | measured |
| P2 | **Presets currently name a single repo, not a list.** `launchPreset` (`schema: "world/preset/v1"`) has `Repo string` (`json:"repo"`), `Profile`, `Agent`, `Model`, `TTLS`, `Grants`, `Daemon`, `Launch`, `Checks`. `loadLaunchPreset` computes `resolvedPreset.SHA256 = serviceHash(data)` over the exact preset file bytes and resolves `ProfileDigest` from the exec profile file. | `cmd/ailang-worldd/preset.go:47-68, 327-455` | measured |
| P3 | `runPresetEpisodeStart` rejects `--grant` overrides (F7), allows `--agent`, `--model`, `--ttl`, `--branch` overrides, checks daemon health (`checkPresetDaemon`), constructs `store.EpisodeStart` (`schema: "world/episode-start/v1"`), and calls `provisionEpisode`. | `cmd/ailang-worldd/episode_preset.go:53-320` | measured |
| P4 | **Credentials store is at schema version 10** (`credentialsSchemaVersion = 10`, `CredentialsBindingWriteVersion = 9`). Tables at v10: `approvers`, `episode_starts`, `key_decisions`, `mint_requests`, `operator_action_steps`, `operator_actions`, `operator_runner_state`, `session_credentials`, `surface_approvals`, `verb_requests`, `webauthn_credentials`. Only `serve` upgrades (`UpgradeCredentials`); CLI verbs never upgrade and refuse older schemas with `OlderCredentialsStoreError`. | `host/store/credentials.go:44-52, 81-84, 217-224, 344-370`; `host/store/credentials_schema.sql:1-133`; `host/store/operator_schema.sql:1-62`; `host/store/mint_requests.go:63-74` | measured |
| P5 | All session mints go through `host/authority.mintRecorded`, which clamps grant expiry to session expiry, hashes the 32-byte raw token with SHA-256 (`credential_id`), and writes `session_credentials` plus its companion record (`episode_starts` or `mint_requests` consume) in **one SQLite transaction** under `credentials.db.writer.lock`. | `host/authority/mint.go:108-198`; `host/store/credentials.go:195-210` | measured |
| P6 | `host/approval` is pure Go (standard library only, no I/O, no clock, no randomness) and enforces strict domain separation across approval kinds using distinct message prefixes and SSHSIG namespaces: mint (`world-approval/v1`, `world-approval@sunholo`), verb (`world-verb/v1`, `world-verb@sunholo`), operator action (`world-operator-action/v1`, `world-operator-action@sunholo`). `host/authority/judge.go` provides one shared judge for CLI and daemon. | `host/approval/request.go:1-28`; `host/approval/sshsig.go:19, 262, 353`; `host/approval/verb.go:32-35`; `host/approval/operator_action.go:24-27`; `host/authority/judge.go:1-100` | measured |
| P7 | `provisionEpisode` opens `/dev/tty` unless `r.consume != nil && r.consume.Signed != nil`. `approver enroll\|revoke` uses a fixed 120 s `/dev/tty` prompt (`confirmApproverAct`) with no store lock held while waiting. | `cmd/ailang-worldd/session_new.go:405-469`; `cmd/ailang-worldd/approver.go:25-45` | measured |
| P8 | The daemon registers `/surface/*` approval routes behind `d.approvalRoute(...)`, reads `credentials.db` through `OpenCredentialsReadOnly`, and takes the credentials writer lock only for short per-request writes. | `host/daemon/daemon.go:1149-1162` | measured |
| P9 | `EpisodeStart` (`host/store/episode_starts.go:18-35`) stores immutable start metadata (`preset`, `repo`, `base`, `branch`, `agent`, `model`, `overrides`, `task`, `checks`, `ttl_s`, `resumes`, `launches`), but has no `allowance` provenance field yet. | `host/store/episode_starts.go:18-35` | measured |
| P10 | Row 196 (mission loop executors through World) is still GATED and has no World caller yet. Exclusion must be enforced both in the allowance payload and at draw admission. | `design_docs/planned/w-all-agent-work-through-world.md:67, 147` | measured |

---

## §2 Design

### §2.1 Credentials store schema v11 (`standing_allowances` + `allowance_draws`)

Schema version bumps from `10` to `11` (`credentialsSchemaVersion = 11`, upgraded only by `serve`
via `UpgradeCredentials`). Two new tables are added to `credentials.db`:

```sql
-- Standing session allowances (version 11, row 224, D-WORLD-96).
-- One row per approved standing allowance for a launch preset. Revoking sets
-- revoked_at (never deletes), so every drawn session keeps its audit chain.
CREATE TABLE IF NOT EXISTS standing_allowances (
    allowance_id    TEXT PRIMARY KEY CHECK (length(allowance_id) = 19 AND allowance_id GLOB 'sa-[0-9a-f]*'),
    preset_name     TEXT NOT NULL CHECK (length(preset_name) BETWEEN 1 AND 51),
    preset_sha256   TEXT NOT NULL CHECK (length(preset_sha256) = 64),
    profile_project TEXT CHECK (profile_project IS NULL OR profile_project <> ''),
    profile_digest  TEXT CHECK ((profile_digest IS NULL) = (profile_project IS NULL) AND (profile_digest IS NULL OR profile_digest <> '')),
    payload         BLOB NOT NULL CHECK (length(payload) BETWEEN 1 AND 131072),
    digest          TEXT NOT NULL CHECK (length(digest) = 64),
    challenge       TEXT NOT NULL CHECK (length(challenge) = 64),
    max_per_day     INTEGER NOT NULL CHECK (max_per_day BETWEEN 1 AND 10),
    ttl_s           INTEGER NOT NULL CHECK (ttl_s BETWEEN 1 AND 43200),
    approved_at     INTEGER NOT NULL CHECK (approved_at > 0),
    expires_at      INTEGER NOT NULL CHECK (expires_at > approved_at AND expires_at - approved_at <= 2592000),
    approved_by     TEXT NOT NULL CHECK (approved_by <> ''),
    artifact        BLOB CHECK (artifact IS NULL OR length(artifact) BETWEEN 1 AND 16384),
    revoked_at      INTEGER CHECK (revoked_at IS NULL OR revoked_at >= approved_at),
    revoked_by      TEXT CHECK ((revoked_at IS NULL) = (revoked_by IS NULL) AND (revoked_by IS NULL OR revoked_by <> '')),
    revoked_reason  TEXT CHECK (revoked_reason IS NULL OR revoked_at IS NOT NULL),
    CHECK ((approved_by = 'tty') = (artifact IS NULL))
);

-- At most one unrevoked allowance per preset name at a time (see R4).
CREATE UNIQUE INDEX IF NOT EXISTS standing_allowances_active_preset
    ON standing_allowances(preset_name)
    WHERE revoked_at IS NULL;

-- Draw ledger (version 11, row 224). One row per session drawn from an allowance.
-- Inserted in the same transaction as session_credentials and episode_starts.
CREATE TABLE IF NOT EXISTS allowance_draws (
    draw_id       TEXT PRIMARY KEY CHECK (length(draw_id) = 19 AND draw_id GLOB 'ad-[0-9a-f]*'),
    allowance_id  TEXT NOT NULL REFERENCES standing_allowances(allowance_id),
    day_key       TEXT NOT NULL CHECK (length(day_key) = 10 AND day_key GLOB '[0-9][0-9][0-9][0-9]-[0-9][0-9]-[0-9][0-9]'),
    episode_id    TEXT NOT NULL CHECK (episode_id <> ''),
    credential_id TEXT NOT NULL UNIQUE CHECK (length(credential_id) = 64),
    repo          TEXT NOT NULL CHECK (repo <> ''),
    drawn_at      INTEGER NOT NULL CHECK (drawn_at > 0),
    ttl_s         INTEGER NOT NULL CHECK (ttl_s BETWEEN 1 AND 43200)
);

CREATE INDEX IF NOT EXISTS allowance_draws_daily
    ON allowance_draws(allowance_id, day_key);
```

**Traceability in `episode_starts`:** `store.EpisodeStart` gains an optional field:

```go
type StartAllowance struct {
    AllowanceID string `json:"allowance_id"`
    DrawID      string `json:"draw_id"`
    ApprovedBy  string `json:"approved_by"`
    DayKey      string `json:"day_key"`
    DrawNumber  int    `json:"draw_number"` // 1..max_per_day
}
```

Because `episode_starts.start_json` is immutable and survives `session revoke`, every session drawn
from an allowance remains permanently traceable to `(allowance_id, draw_id, approved_by)` even after
the session expires or is revoked.

---

### §2.2 Cryptographic & attended approval binding (`host/approval/allowance.go`)

A standing allowance is a canonical JSON document (`schema: "world/standing-allowance/v1"`) with
strictly closed fields:

- `schema`: `"world/standing-allowance/v1"`
- `preset_name`: preset basename (1–51 chars, episode grammar)
- `preset_sha256`: 64-hex SHA-256 of the exact preset file bytes (`resolvedPreset.SHA256`)
- `profile`: `{project, digest}` (both empty if preset has no exec profile; both set otherwise)
- `repos`: non-empty, sorted, deduplicated list of absolute clean repository paths
- `grants`: non-empty list of `{effect, scope, budget}` in preset order
- `ttl_s`: `1..43200` (at most 12 h per D-WORLD-96)
- `max_per_day`: `1..10` (at most 10/day per D-WORLD-96)
- `valid_days`: `1..30` (at most 30 days per D-WORLD-96)
- `world_db`, `credentials_db`, `workspace_root`: canonical absolute paths binding the store & root
- `exclude_mission_loop`: must be `true` (enforced by `CanonicalAllowance`)

**Domain separation (fourth approval domain):**

| domain | prefix | SSHSIG namespace | message length |
|---|---|---|---|
| mint request (row 201) | `world-approval/v1` | `world-approval@sunholo` | 89 bytes |
| verb confirm (row 243 V3) | `world-verb/v1` | `world-verb@sunholo` | 86 bytes |
| operator action (row 243 V4) | `world-operator-action/v1` | `world-operator-action@sunholo` | 97 bytes |
| **standing allowance (row 224)** | **`world-allowance/v1`** | **`world-allowance@sunholo`** | **91 bytes** |

`AllowanceMessage(digest, challenge, expiresAt)` is:
`"world-allowance/v1\n" || digest (32 B) || challenge (32 B) || big-endian int64(expiresAt) (8 B)`.

A signature over a mint request, verb, decision, or operator action can **never** verify as a
standing allowance, and vice versa.

**Tamper & drift protection:**
1. At draw time, the store recomputes `sha256(payload)` and compares it to `standing_allowances.digest`.
2. `episode start --preset <p> --from-allowance` re-loads the preset and exec profile from disk and
   passes their current `preset_sha256` and `profile_digest` to `AdmitAllowanceDraw`.
3. If either the preset file or the exec profile changed by even 1 byte since the allowance was
   approved, the draw is **refused** (forces explicit approval or a fresh `allowance grant`).

---

### §2.3 Atomic draw & daily-cap check

`authority.MintFromAllowance` runs inside a **single `BEGIN IMMEDIATE` transaction** under
`credentials.db.writer.lock`:

1. Read the active `standing_allowances` row for `preset_name` (or explicit `allowance_id`).
2. If `approved_by != "tty"`, check `approvers` to confirm the signing key is still enrolled and
   `revoked_at IS NULL` (see R6).
3. Compute `dayKey = approval.AllowanceDayKey(now)` (UTC `YYYY-MM-DD`, see R3) and query:
   `SELECT COUNT(*) FROM allowance_draws WHERE allowance_id = ? AND day_key = ?`.
4. Call pure `approval.AdmitAllowanceDraw(AllowanceDrawInput{...})`, checking in order:
   - **D1 (integrity):** `sha256(payload) == stored_digest` and canonical JSON valid
   - **D2 (mission loop):** `caller != "mission-loop"` and `allowance.ExcludeMissionLoop == true`
   - **D3 (revocation):** `revoked_at == nil` (and signing approver not revoked)
   - **D4 (clock & expiry):** `now >= approved_at`, `now < expires_at`, `expires_at - approved_at <= 30d`
   - **D5 (preset & profile binding):** `preset_name`, `preset_sha256`, `profile.project`, `profile.digest` match
   - **D6 (repo bound):** requested `repo` is in `allowance.Repos`
   - **D7 (grants & TTL bounds):** requested `grants` match `allowance.Grants` exactly; `0 < ttl_s <= allowance.TTLS <= 43200`
   - **D8 (daily cap):** `draws_today < allowance.MaxPerDay` (and `MaxPerDay <= 10`)
5. Insert `session_credentials`, `episode_starts` (with `StartAllowance`), and `allowance_draws`
   atomically, then `COMMIT`.

Because `credentials.db` uses a single writer lock (`credentials.db.writer.lock`) and SQLite
`BEGIN IMMEDIATE`, two concurrent `episode start --from-allowance` calls serialize cleanly; the
11th draw in a day is guaranteed to observe `draws_today == 10` and refuse.

**Post-mint failure rollback:** If writing `<sessions>/<ep>.session` or `<ep>.mcp.json` fails
after commit (matching `session_new.go:542-555`), the session credential is revoked via
`authority.Revoke`. The `allowance_draws` row is kept for audit (or marked failed if ruled); keeping
it ensures a failing caller cannot loop to probe mints without hitting the daily cap.

---

### §2.4 Revocation (CLI + App/Surface)

- **CLI:** `ailang-worldd allowance revoke <allowance-id|preset-name> [--reason TEXT] [--revoke-sessions]`
  - Sets `revoked_at = now`, `revoked_by = "tty"`, `revoked_reason = reason` in one transaction.
  - Rows are never deleted (same invariant as `approvers`).
  - Default leaves already-minted live sessions intact (they expire naturally); `--revoke-sessions`
    also deletes any live `session_credentials` rows linked through `allowance_draws`.
- **App / Surface (`/surface`):**
  - `GET /surface/allowances.json` (read-only): lists active (and optionally recent expired/revoked)
    allowances with preset name, `approved_by`, `expires_at`, and `draws_today / max_per_day`.
  - `POST /surface/allowances/{id}/revoke`: sets `revoked_at = now`, `revoked_by = "surface"`
    (see R5 for auth model). Immediately blocks the next `--from-allowance` draw.

---

### §2.5 CLI verbs

1. **`ailang-worldd allowance grant --preset <name> [--max-per-day 10] [--ttl 43200] [--days 30] [--replace] [--signature <file> | --webauthn-assertion <file>]`**
   - Loads and validates the preset (`loadLaunchPreset`), checks daemon health, builds canonical
     `StandingAllowance`, verifies the key signature (or attended confirmation per R2), and inserts
     into `standing_allowances`.
   - Refuses `max-per-day > 10`, `ttl > 43200`, or `days > 30` (hard ceilings from D-WORLD-96).
2. **`ailang-worldd allowance list [--all] [--json]`**
   - Read-only via `OpenCredentialsReadOnly`. Displays active allowances with `draws_today / max_per_day`,
     `expires_at`, `preset_sha256` prefix, and `approved_by`. `--all` includes revoked and expired.
3. **`ailang-worldd allowance revoke <allowance-id|preset-name> [--reason TEXT] [--revoke-sessions]`**
   - Marks the allowance revoked.
4. **`ailang-worldd episode start [episode] --preset <name> --from-allowance [--launch] [--headless] [--task FILE] ...`**
   - Skips `/dev/tty` prompt and calls `authority.MintFromAllowance`.
   - Rejects `--grant` (already rejected by F7), rejects `--ttl` exceeding the allowance's `ttl_s`
     (or 12 h), rejects any `--repo` outside `allowance.Repos`.
   - Every refusal prints a three-part message (row 227: what failed, why, and the exact next
     command — e.g. `ailang-worldd episode start --preset <name> ...` without `--from-allowance` for
     an explicit approval, or `ailang-worldd allowance grant --preset <name>` if the preset changed
     or expired).

---

### §2.6 Mission loop exclusion (row 196)

D-WORLD-96 explicitly excludes the mission loop (row 196). Enforcement is layered:
1. `StandingAllowance.ExcludeMissionLoop` is required to be `true` in `CanonicalAllowance`.
2. `AdmitAllowanceDraw` unconditionally refuses `Caller == "mission-loop"`.
3. `episode start --from-allowance` refuses if mission-loop environment markers (e.g.
   `WORLD_MISSION_LOOP` or `MISSION_LOOP_RUN_ID`) are set, printing a refusal pointing to row 196's
   separate delegated-minting ruling.

---

## §3 Vertical Slice Plan

| Slice | Scope | Package(s) | Status |
|---|---|---|---|
| **S1** | Pure allowance domain, canonical JSON + digest, domain-separated SSHSIG signing/verification, and full D-WORLD-96 draw-admission law + unit tests | `host/approval/allowance.go`, `host/approval/allowance_test.go` | **Implemented (this PR)** |
| **S2** | Credentials schema v11 (`standing_allowances`, `allowance_draws`, v10→v11 upgrade), `EpisodeStart.Allowance`, `authority.MintFromAllowance` / `GrantAllowance` / `RevokeAllowance` | `host/store/`, `host/authority/` | Pending R1–R6 |
| **S3** | CLI verbs (`allowance grant\|list\|revoke`, `episode start --preset <p> --from-allowance`), three-part refusals, herdr `world-agent-popup` integration | `cmd/ailang-worldd/`, `integrations/herdr/` | Pending S2 |
| **S4** | `/surface/allowances.json`, `POST /surface/allowances/{id}/revoke`, and allowance status/revoke card in the World app | `host/daemon/` | Pending S2 |

---

## §4 Open Rulings

| # | Question | Options | Recommended Default |
|---|---|---|---|
| **R1** | **Multi-repo presets vs single-repo presets.** `launchPreset` (`cmd/ailang-worldd/preset.go:50`) currently has `repo: string`, while D-WORLD-96 says "the repos listed in the preset". | (A) Keep `StandingAllowance.Repos []string` in S1/S2 (1 element for current presets) and add optional `repos: []string` to `world/preset/v1` in S3; (B) Restrict presets and allowances to 1 repo per preset permanently. | **(A)** — future-proof without breaking existing presets; 1-repo presets simply populate a 1-element `Repos` slice. |
| **R2** | **Approval mechanism for `allowance grant`.** D-WORLD-96 says "Mark approves a preset once with his key". | (A) Require an enrolled key (SSH or WebAuthn passkey) — matches D-WORLD-96 literally, since row 201 is built; (B) Allow attended `/dev/tty` `y` as well (recorded as `approved_by = "tty"`). | **(A)** — require an enrolled key (SSH or WebAuthn); allow `--tty` only if explicitly requested for bootstrap environments without enrolled keys. |
| **R3** | **Daily cap boundary (10 sessions/day).** | (A) UTC calendar day (`YYYY-MM-DD` UTC) — deterministic across DST and timezones; (B) Operator local calendar day (`time.Local`); (C) Rolling 24-hour window (`drawn_at > now - 86400`). | **(A) UTC calendar day** — pure, deterministic, no DST edge cases; refusal message prints next reset in local time. |
| **R4** | **Granting an allowance when one is already active for the same preset.** | (A) Refuse unless `--replace` is passed (atomically revokes prior allowance with `revoked_reason = "superseded"`); (B) Always auto-supersede; (C) Allow multiple concurrent allowances per preset. | **(A)** — prevents accidental silent replacement while keeping preset updates one command (`--replace`). |
| **R5** | **Authentication for `/surface` revocation (`POST /surface/allowances/{id}/revoke`).** | (A) Loopback + CSRF/origin check without requiring a key tap (`revoked_by = "surface"`) — zero-friction emergency stop; (B) Require a passkey/SSH signature like `POST /surface/approve/{id}`. | **(A)** — revocation is strictly authority-reducing (emergency stop); zero friction is safer than gating stop behind a key tap. |
| **R6** | **Approver key revocation cascading to allowances.** | (A) Revoking an enrolled approver key (`approvers.revoked_at IS NOT NULL`) immediately blocks draws from allowances signed by that key; (B) Allowances remain valid until separately revoked or expired. | **(A)** — if a key is compromised and revoked, standing allowances signed by it must stop immediately. |
