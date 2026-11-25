# Row 269 — UI Profile & Preset Onboarding (`w-ui-profile-preset-onboarding`)

> **Status:** Planned (r0, filed attended 2026-11-25)  
> **Depends on:** Row 243 V4 (operator runner + passkey-confirmed actions), PR #648 (`profile list`/`activate` + AILANG detector), PR #646 (`preset init --se-tools --mode --model`)  
> **Goal:** Let the operator inspect, onboard, and drive all repositories on the rig (exec profiles, presets, and headless/interactive episodes) directly from the AILANG World web/phone UI without dropping to a terminal.

---

## 1. Motivation

As of `dev` `f8fb3ac1`:
- `ailang-worldd profile init|list|activate` and `ailang-worldd preset init|list|show|check` make onboarding any Go, Python, Node/Vitest, or AILANG repo a two-step terminal flow.
- In `host/verbs/lookup.go:38-39`, both `profile` and `preset` are classified as `Other` (terminal-only).
- Row 243 V4 (Sprints 3–12) is building the passkey-confirmed operator runner for `episode start`, `episode launch`, and `episode finish` from the World UI.
- To drive **all** rig projects (`docparse`, `ailang-parse`, `cphu-aipla-app`, `ndttl`, `daneel`, `twilightgame`, `ailang-world`, etc.) from the World UI, the operator also needs to:
  1. See all active/proposed exec profiles and configured presets in the World UI (and pick a preset when launching an episode from the browser/phone).
  2. Propose, review, and passkey-approve new exec profiles (`profile init` → review JSON → `profile activate`) and presets (`preset init`) from the UI via the V4 operator runner, with live single-profile probing so no full `service restart` is required.

---

## 2. Design Outline

### Phase 1 — Read-Only UI & Address-Bar Visibility (zero runner dependency)
1. **Daemon read endpoints (localhost / authenticated surface only):**
   - `GET /v1/surface/profiles`: lists active (`<project>.json`) and proposed (`<project>.proposed.json`) exec profiles from `--exec-profiles-dir` alongside `/v1/health`'s live probe status (`loaded` vs `excluded` + reason).
   - `GET /v1/surface/presets`: lists configured presets from `--presets-dir` (name, repo, base branch, agent, mode, profile, digest, checks summary — never secrets).
2. **Address-bar `Read` verbs (`host/verbs/lookup.go`):**
   - Promote `profile list` and `preset list` / `preset show <name>` to `Read` verbs (keeping mutating subcommands `init`/`activate` terminal-only until Phase 2).
3. **World UI "Repos & Presets" panel (`host/surfacepage/`):**
   - Lists every configured preset with its repo, agent (`agy`/`codex`/`claude`), mode (`headless`/`interactive`), and exec profile health badge (`probed` / `excluded` / `proposed`).
   - One-click prefill into the V4 episode start bar (`episode start --preset <name> ...`).

### Phase 2 — Passkey-Approved Onboarding via V4 Operator Runner (after Row 243 V4 S8)
1. **Propose from UI:** operator enters a repo path (or picks from a configured workspace root allowlist) and project name. The operator runner executes `ailang-worldd profile init --repo <dir> --project <name>` (non-destructive: writes only `<project>.proposed.json`).
2. **Review & Passkey Confirm:** the World UI displays the exact proposed profile JSON (`root`, `commands`, `read_roots`, `path`, `probe`, `digest`) on a row-201 confirmation card.
3. **Activate & Live Probe:** after enrolled passkey confirmation, the operator runner validates & activates `<project>.proposed.json` → `<project>.json`, runs the 7-arm startup probe for that single profile, and signals the daemon to register the newly probed profile live (no full daemon restart).
4. **Create Preset from UI:** operator chooses preset name, agent (`agy`/`claude`/`codex`), mode (`headless`/`interactive`), `--se-tools`, and profile; passkey confirmation writes `<presets-dir>/<name>.json`.

---

## 3. Security Invariants
- **Read-only daemon preserved:** `worldd serve` never writes to `--exec-profiles-dir` or `--presets-dir`; all mutations execute exclusively in the V4 operator runner after passkey verification.
- **Exact-bytes confirmation:** passkey approval binds to the SHA-256 digest of the proposed profile/preset JSON (same exact-byte law as `episode finish --confirm-on-surface`).
- **Fail-closed probing:** an activated profile is never usable by `workspace-exec` until it passes all 7 confinement arms (bounded by `procbound.MaxOutstanding/2` per PR #653).

---

## 4. Acceptance Criteria
- Phase 1: `profile list` and `preset list` work as `Read` verbs in the World UI; the UI displays all active/proposed profiles and presets with live probe status.
- Phase 2: onboarding a new repo from the World UI (propose → passkey approve → live probe → create preset → launch headless `agy` episode) succeeds end-to-end with zero terminal steps.
