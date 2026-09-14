# w-check-no-personal-email — a real, non-vacuous privacy gate over the loop-written tracked surface, wired into CI's go-verify job directly

- Status: **planned** · Date: 2026-09-14 · Designer: pi:ollama/deepseek-v4-flash:0731-cloud (iter-175; rotation entry codex:gpt-6-astra was ration-blocked this fire) · Base commit: a295291bea069d793960e08f055e501c6f9661f4 · Owning queue row: 74 · Scope class: HARNESS — enforcement of the attended-ledger privacy contract · Verify profile: ailang-code

## §1 Problem

Row 74's premise is that the shared mission-control skill's ATTENDED LEDGER EDITS contract (Mark, attended 2026-09-02) declares `make check-no-personal-email` fails the build if a personal address reaches a tracked file — while that target existed **nowhere** in `sunholo-data/ailang-world`. The history is concrete and public:

- At iter-152, this PUBLIC repo carried a personal address in **7** doc locations. It was redacted in that iteration by substitution (7 → 0), but the redaction was a one-time manual sweep with no durable enforcement behind it.
- The enforcement contract was prose-only inside `~/.claude/` (through a symlink to the V1 checkout), invisible to this repo's CI, and unverifiable here.

**[controller-measured 2026-09-14 at a295291]** F1: `git grep -n 'check-no-personal-email' -- . ':!design_docs'` → **0 hits (rc=1)**. Positive control in the same tree: `git grep -c 'verify_ail.sh' -- .github` → **3**. So there is no such gate anywhere in World's tracked tree outside charter prose.

**[controller-measured 2026-09-14 at a295291]** F4: running the fleet's scan logic against World's tracked tree (34 in-scope files) yields **EXACTLY ONE hit**: `scripts/mission_answer.sh` line 69, `ATT_EMAIL="${MISSION_ATTENDED_EMAIL:-<personal address>}"` — the functional default that row 74 deliberately left alone. Everything else in scope is clean; the iter-152 redaction held.

**Why "passes today" is not a gate**: the single in-scope hit is a live personal address sitting in a tracked, loop-written script today, and no CI step would notice it. The contract is enforceable or it is decoration. "It happens to be clean except one address" is not enforcement — the instrument and its wiring do not exist. This item makes the contract real: an address-pattern check over the loop-written tracked surface, with the functional defaults handled, wired into CI, and proven non-vacuous with a mutation arm that plants an address and asserts the gate REDs.

## §2 Design

### D1 — The instrument: `scripts/check_no_personal_email.sh`

**Contract.**

- **Inputs**: `git ls-files | grep -E "$SCOPE_RE"` — the loop-written tracked surface, filtered by World's scope regex (see D2). Files are skipped when they match the binary-extension probe (the same extensions the fleet's precedent skips).
- **Address regex**: `PAT='[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}'`.
- **Exclusion list (verbatim, anchored)**:
  ```
  grep -vE 'users\.noreply\.github\.com|noreply@|@example\.(com|org|net)$|\.(invalid|test|localhost)$|gserviceaccount\.com|@sentry\.io'
  ```
  Note the `$` anchor added to `@example\.(com|org|net)$` — see D5 arm H.
- **Output lines (verbatim)**:
  - clean: `✓ check-no-personal-email: no personal addresses in the loop-written surface` → exit **0**.
  - a hit: `✗ personal email in loop-written file: <f>  ->  <addr>` (one per hit) followed by a three-line explanation, exit **1**.
  - instrument failure: the anti-vacuity floor message (below), exit **2**.
- **Exit codes**: `0` clean; `1` ≥1 address hit; **`2` instrument failure**.
- **Anti-vacuity floor (third exit code)**: if `git ls-files | grep -E "$SCOPE_RE"` yields **zero** in-scope files, print exactly:
  ```
  ✗ check-no-personal-email: zero in-scope files scanned — the gate must never print ✓ over an empty surface
  ```
  and exit **2**. A gate that scans nothing must not print `✓`. This is the floor that makes the gate non-vacuous on the surface-count side; arm F plus AC-M1-9 pin it.

`set -u`, safe under bash 3.2 (macOS): no associative arrays, no `mapfile`, no `${var,,}`. `ROOT` derived from `dirname "$0"/..`; `cd "$ROOT"` before scanning.

### D2 — Scope regex is World's own

**[controller-measured 2026-09-14 at a295291]** F7: World's in-scope tracked files under the fleet's regex are `design_docs/world-mission.md`, `world-mission-log.md`, `world-mission-log-archive.md`, `world-mission-status-archive.md`, `world-mission-dashboard.md`, `world-mission-index.md` (6) + 28 under `scripts/` = **34**. World has **NO `.claude/skills/`** directory (the mission-control skill resolves through a symlink into the V1 checkout), so that regex alternative matches nothing here.

**Decision: drop** `.claude/skills/.*` and **keep** `scripts/.*` and the mission-doc alternative. World's scope regex is:
```
SCOPE_RE='^(design_docs/[^/]*mission[^/]*\.md|scripts/.*)$'
```
A regex alternative that can never match is a vacuous claim: it costs nothing to run but asserts coverage that does not exist, and it invites a future editor to believe World owns a skills tree it does not. Dropping it keeps the gate honest about the surface it actually guards. This is a World-local instrument modelled on the fleet's — the scope regex is **World's own by explicit decision**, not inherited by copy-paste blind.

### D3 — The functional default: noreply, not an allow-list

**[controller-measured 2026-09-14 at a295291]** F5: the fleet resolved the SAME default in ITS `scripts/mission_answer.sh` by switching it to the GitHub noreply identity — V1 origin/dev line 108 reads `ATT_EMAIL="${MISSION_ATTENDED_EMAIL:-3155884+MarkEdmondson1234@users.noreply.github.com}"` (commit `8369877d9`, 2026-09-04). The noreply form is attributable to the same GitHub account, is what GitHub itself stamps on web commits, and is **allowed** by the gate's exclusion list — so no allow-list mechanism is needed at all.

**[controller-measured 2026-09-14 at a295291]** F6: World's `scripts/mission_answer.sh` guards are `FLEET_PATTERN="${MISSION_FLEET_ACCOUNT:-sunholo-voight-kampff}"`, refusing when `$ATT_EMAIL` contains the fleet pattern or `$ATT_NAME` matches `*[Bb]ot*`. Its suite `scripts/test_mission_answer.sh` has arms 1, 2, 3, 4a–4d, 5a; **every arm sets `MISSION_ATTENDED_EMAIL` explicitly**, so no existing arm exercises the default.

**Decision: adopt the noreply default** with a one-line change at `scripts/mission_answer.sh:69`:
```
ATT_EMAIL="${MISSION_ATTENDED_EMAIL:-3155884+MarkEdmondson1234@users.noreply.github.com}"
```
The noreply default does not contain `sunholo-voight-kampff`, so arms 4a–4d are unaffected. The exclusion list already permits `users.noreply.github.com`, so the gate stays simple: **the row's "allow-list for functional defaults" is explicitly declared NOT built**, for exactly this reason. An allow-list would re-introduce the very escape (a private allow-listed domain) the gate exists to close; the noreply default removes the need for one.

**One new arm in `scripts/test_mission_answer.sh`**: with `MISSION_ATTENDED_EMAIL` UNSET, `--dry-run` on the fixture must succeed (rc 0) and the script's resolved identity must be the noreply address. Precise form: the arm runs the fixture under `MISSION_ATTENDED_EMAIL= env -u MISSION_ATTENDED_EMAIL bash scripts/mission_answer.sh … --dry-run` and greps a `--dry-run` line that names the identity. If the script does not print the identity on `--dry-run`, this milestone adds **one line** to the script's `note`/dry-run output that prints the resolved `ATT_EMAIL`. The arm asserts rc 0 AND the noreply address on that printed line.

### D4 — CI wiring

**[controller-measured 2026-09-14 at a295291]** F2: World has NO `Makefile` (`ls Makefile` → No such file or directory). CI is `.github/workflows/ci.yml` with two jobs: `ailang-verify` and `go-verify` ("go host build + test gate"). The `go-verify` job currently ends with these steps, each `timeout-minutes: 2`:
- `Gate-1 range instrument suite` → `bash scripts/test_gate1_range_check.sh`
- `Queue census instrument suite` → `bash scripts/test_queue_census.sh`
- `Queue census (live charter, controlled)`
- `Gate-0 self-notice instrument suite` → `bash scripts/test_gate0_self_notices.sh`

**Decision: wire into CI's `go-verify` job DIRECTLY as its own steps — NOT `scripts/verify_go.sh`.** Charter row 76 established that `scripts/verify_go.sh` is NOT usable as a gate on this rig (it short-circuits on driver drift). Two new steps appended after the last one, matching the sibling instrument-suite pattern:

```yaml
      - name: Refuse a personal email in the loop-written surface
        timeout-minutes: 2
        run: bash scripts/check_no_personal_email.sh
      - name: Personal-email gate self-test
        timeout-minutes: 2
        run: bash scripts/test_check_no_personal_email.sh
```

This is exactly the pattern the sibling suites (Gate-1, Queue census, Gate-0) already follow.

### D5 — The self-test: `scripts/test_check_no_personal_email.sh`

Five arms A–E modelled on the fleet's, in a throwaway `git init` repo, each fixture address **assembled at runtime by string concatenation** so the test file itself carries no address. Plus three World-only arms F, G, H:

- **A — real clean repo** (rc 0, `✓` line).
- **B — a personal address in a mission doc** → rc 1 (fixture assembled at runtime).
- **C — GitHub noreply allowed** (rc 0).
- **D — reserved-TLD placeholders allowed** (`someone@example.com` / `@example.org` / `@example.net` → rc 0).
- **E — governance/code files out of scope** (rc 0).
- **F — anti-vacuity floor**: empty scope → rc **2**, and the `✓` line must NOT print (the instrument-failure message does). This proves the D1 floor.
- **G — a hit in `scripts/`** (not only in a mission doc) is caught → rc 1. Pins D2 that `scripts/.*` is a live alternative.
- **H — the exclusion is exact and anchored**: an address at a domain that merely CONTAINS `example.com` as a prefix, e.g. `a@example.com.evil.net`, is a **HIT**, not allowed.

**On H's anchoring question** (the design tension called out in the row): the fleet regex `@example\.(com|org|net)` does **NOT** anchor at the end, so `a@example.com.evil.net` WOULD be excluded under fleet parity. **Decision: anchor with `$`** (`@example\.(com|org|net)$`) and add arm H to prove the anchored behavior. The reserved TLDs are meant as *final* TLD placeholders, not as a prefix licence that lets `example.com.<attacker-owned>` become an exempt surface — a privacy gate must not whitelist a substring. The single deviation from fleet parity is this one `$` anchor, and it is deliberate, documented here, and pinned by arm H. (If parity were preferred instead, the residual would be recorded and a follow-up row filed — this doc goes the anchored route.)

**D6 rule 3k (as the row directs):** every arm runs `bash scripts/check_no_personal_email.sh` as a process in the throwaway repo — exactly what CI runs — and asserts on stdout/stderr/exit code. The self-test never re-implements the scan; it drives the real binary. The `✓`/`✗`/floor lines are asserted verbatim, so they are part of the contract (D1).

### D6 — Output contract

The three line classes are the contract and are grepped verbatim by the self-test:
- clean → `✓ check-no-personal-email: no personal addresses in the loop-written surface`
- per-hit → `✗ personal email in loop-written file: <f>  ->  <addr>` then a three-line explanation
- floor → `✗ check-no-personal-email: zero in-scope files scanned — the gate must never print ✓ over an empty surface`

The per-hit line embeds the matched address, so the line itself carries no personal information beyond what is already in the tracked file being flagged — acceptable, since the whole point of the ✗ branch is to name the offender for the fixer.

## §3 Files to create / modify

| File | Action | Delta | Notes |
|---|---|---|---|
| `scripts/check_no_personal_email.sh` | create | +~70 | new executable instrument (D1), `SCOPE_RE` from D2, anchored exclusion, exit-2 floor |
| `scripts/test_check_no_personal_email.sh` | create | +~90 | new self-test, arms A–H (D5), runtime-concatenated fixtures |
| `scripts/mission_answer.sh` | modify | 1 line | line 69 default → noreply form (D3); +1 dry-run identity line if absent |
| `scripts/test_mission_answer.sh` | modify | +1 arm | new default-identity arm (D3) |
| `.github/workflows/ci.yml` | modify | +2 steps | appended to `go-verify` (D4) |

Nothing under `tools/launchd/` (frozen fleet core). No `.ail` files.

**Frozen-core note (explicit, per row 74's own requirement):** this is a World-local instrument modelled on the fleet's `f0d44915d`, NOT a vendored fork of frozen core. `scripts/` is World-owned; only `tools/launchd/*` is frozen fleet core. The scope regex is World's own (D2); World cannot and does not edit the shared skill.

## §4 Acceptance criteria

Each AC is a command with expected output. M denotes a milestone executor step; C denotes a controller gate step.

- **AC-M1-1** `test -x scripts/check_no_personal_email.sh` → true. (Mode bit: a previous iteration lost one to `mv` from /tmp; the executable bit is asserted explicitly.)
- **AC-M1-2** `bash scripts/check_no_personal_email.sh; echo rc=$?` → `✓ check-no-personal-email: no personal addresses in the loop-written surface`, `rc=0`, on the sprint tree (after D3's fix).
- **AC-M1-3** the SAME command at BASE a295291bea069d793960e08f055e501c6f9661f4 → `rc=1`, the ✗ line naming `scripts/mission_answer.sh` (proves the instrument sees the very hit row 74 named, BEFORE D3's fix).
- **AC-M1-4** `bash scripts/test_check_no_personal_email.sh` → `8 passed, 0 failed` (arms A–H).
- **AC-M1-5** `bash scripts/test_mission_answer.sh` → all arms pass, including the new default-identity arm (arm asserting rc 0 and the noreply identity under `env -u MISSION_ATTENDED_EMAIL`).
- **AC-M1-6** `/usr/bin/grep -c 'check_no_personal_email' .github/workflows/ci.yml` → `2` (the two new `run:` steps).
- **AC-M1-7** `git diff --stat a295291 -- tools/launchd/` → empty (no frozen-core touch).
- **AC-M1-8** `go vet ./...` → rc 0; `go test ./... -count=1` with `AILANG_BIN=$HOME/.pinned-ailang/ailang` → rc 0; `bash ./scripts/verify_ail.sh` → rc 0 (11 identities / 40 named tests); profile `ailang-code`. (No `.ail` files touched by this item; `go build` is not a compile fence for `_test.go`.)
- **AC-M1-9** `git ls-files | grep -cE '^(design_docs/[^/]*mission[^/]*\.md|scripts/.*)$'` → `34` — the in-scope count, the D1 floor's positive control (the same surface the gate scans must number 34, so the floor is not silently vacuous).
- **AC-M1-10** after D3: in-scope address scan over the tree → 0 hits, AND `grep -c 'users\.noreply\.github\.com' scripts/mission_answer.sh` → `1` (the literal noreply form appears exactly once, on line 69).
- **C** queue row 74 tag → LANDED; repo profile note records the gate's location.

## §5 Test plan / mutation matrix

Each mutation is applied by `cp` backup + `sed`, then restored by `cp` back; restore is asserted by `shasum -a 256` (content) AND `test -x` (mode) — NEVER `git checkout --` (charter row 82). `git checkout --` is forbidden because the worktree is loop-managed and a checkout is a controller-only action.

| # | Mutation | Which arm/criterion goes RED | Sole killer? |
|---|---|---|---|
| M1 | on a hit, `exit 1` → `exit 0` | arm B (expects rc 1) | yes |
| M2 | delete the exclusion `grep -vE` | arm C (noreply) or D (reserved TLD) red | yes |
| M3 | scope regex narrowed to mission docs only | arm G (a `scripts/` hit) red | yes |
| M4 | remove the empty-scope floor | arm F (expects rc 2) red | yes |
| M5 | hits counter never incremented | arm B (expects ✗ + rc 1) red | yes |
| M6 | `$` anchor removed from the exclusion (reverts to fleet parity) | arm H (expects `a@example.com.evil.net` to be a HIT) red | yes |
| M7 | `scripts/mission_answer.sh` default reverted to a raw address | AC-M1-2 red: the LIVE gate catches it. This is the ONE mutation the live gate kills, and it is the row's whole point. | yes |
| M8 | GREEN control: reword a comment inside the script | every arm stays green | n/a — control |

M7 is the mutation that closes the loop: it proves the gate is not merely self-test-consistent but actually guards the live surface it is wired to scan — the single hit row 74 named dies under the real gate.

## §6 Conflict surface

- **What else reads `scripts/mission_answer.sh`'s identity**: the attended-ruling channel itself — Gate 0.6(b)/(c) of the shared mission-control skill relies on the resolved `ATT_EMAIL`/`ATT_NAME` to bind an attended ruling to a human. D3's change only alters the DEFAULT; any real attended run that sets `MISSION_ATTENDED_EMAIL` is untouched, so the attended-ruling channel's behaviour is preserved except when no override is provided (where the noreply form is still attributable to the same GitHub account).
- **What else writes to `ci.yml`**: sibling instrument-suite rows already append steps to `go-verify`. Adding two steps is the same append-only pattern; the `go-verify` job's other steps and their `timeout-minutes: 2` are untouched. Nothing renames or reorders existing steps.
- **The fleet's copy as precedent**: no shared file is edited. `f0d44915d` is cited and mirrored in design intent but is a separate repo tree. World cannot edit the shared skill (resolved through the V1 checkout / symlink); this gate makes the fleet's contract locally enforceable — which is the point.

## §7 Verification Log

Provenance key: `[controller-measured 2026-09-14 at a295291]` = measured first-party by the controller at this base; "this session" = measured by the designer.

| # | Provenance | Command | Result |
|---|---|---|---|
| V1 | `[controller-measured 2026-09-14 at a295291]` | `git grep -n 'check-no-personal-email' -- . ':!design_docs'` | 0 hits (rc=1) |
| V2 | `[controller-measured 2026-09-14 at a295291]` | `git grep -c 'verify_ail.sh' -- .github` | 3 |
| V3 | `[controller-measured 2026-09-14 at a295291]` | `ls Makefile` | No such file or directory (F2) |
| V4 | `[controller-measured 2026-09-14 at a295291]` | fleet repo inspection of `scripts/check_no_personal_email.sh` + `scripts/test_check_no_personal_email.sh` @ `f0d44915d` (2026-09-02) | both present; wired in ITS ci.yml lines 262/265 as `make …`; regex, exclusion, five arms A–E as quoted (F3) |
| V5 | `[controller-measured 2026-09-14 at a295291]` | V1's scan logic over World's 34 in-scope files | exactly one hit: `scripts/mission_answer.sh` line 69 functional default (F4) |
| V6 | `[controller-measured 2026-09-14 at a295291]` | V1 origin/dev line 108 of `scripts/mission_answer.sh` | `ATT_EMAIL="${MISSION_ATTENDED_EMAIL:-3155884+MarkEdmondson1234@users.noreply.github.com}"` @ `8369877d9` (F5) |
| V7 | `[controller-measured 2026-09-14 at a295291]` | World's `scripts/mission_answer.sh` guards + `scripts/test_mission_answer.sh` arms | guards as quoted; arms 1,2,3,4a–4d,5a all set `MISSION_ATTENDED_EMAIL`; none touch the default (F6) |
| V8 | `[controller-measured 2026-09-14 at a295291]` | `git ls-files` under V1 scope regex | 6 mission docs + 28 scripts = 34; NO `.claude/skills/` in World (F7) |
| V9 | `[controller-measured 2026-09-14 at a295291]` | `bash ./scripts/verify_ail.sh`; `go vet ./...`; `go test ./... -count=1` (`AILANG_BIN=$HOME/.pinned-ailang/ailang` v0.30.0) | rc 0 across; 11 identities / 40 named tests; `scripts/verify_go.sh` NOT a gate (row 76) (F8, F9) |
| V10 | this session | `git rev-parse HEAD` | `a295291bea069d793960e08f055e501c6f9661f4` |
| V11 | this session | `ls scripts/check_no_personal_email.sh scripts/mission_answer.sh scripts/test_mission_answer.sh` | `check_no_personal_email.sh` → No such file (gate absent); `mission_answer.sh` and `test_mission_answer.sh` present |
| V12 | this session | `sed -n '60,75p' scripts/mission_answer.sh` | line 69 `ATT_EMAIL="${MISSION_ATTENDED_EMAIL:-<personal address>}"`; guards and `FLEET_PATTERN` as quoted |

No designer-run measurement contradicts an F-row; this session only confirmed existence/location (V10–V12), no new arithmetic.

## §8 Milestones

- **M1 — executor**: implement the five files (new instrument, self-test, one-line default change, +1 test arm, +2 CI steps) and run all AC-M1-* commands to green.
- **M2 — controller**: queue row 74 tag → LANDED; Repo Profile note records that the gate exists and where (`scripts/check_no_personal_email.sh`, wired in `go-verify`); charter mention that the fleet has its own at `f0d44915d`.

Estimate ~**0.3d** (the instrument and self-test are near-copies of a known-good precedent; the only bespoke work is the World scope regex, the exit-2 floor, arms F/G/H, the anchor, and the one-line default change).

## §9 Risks and declared residuals

- **Scope is deliberately narrow**: design docs under `planned/`/`implemented/` are NOT scanned. These are sprint artifacts (ephemeral working docs, this one included), not durable loop-written surfaces, and their transient content is reviewed at ratification. Widening the surface is a follow-up row, not this one. The gate is deliberately scoped to mission docs + `scripts/` (the 34-file durable loop-written surface).
- **Tracked-files-only**: the gate reads `git ls-files` — an untracked draft is invisible until `git add`. That is correct: a draft not staged for landing is not yet "in the loop-written surface," and CI runs on tracked state only. Flagged so a future editor does not assume draft protection.
- **Regex false positives on `user@host`-style shell strings in `scripts/`**: possible in theory, measured none today (V5 — the only in-scope hit is the functional default, itself fixed by D3). Any future shell string shaped like an address can be handled by a comment + a deliberate scope decision, not by weakening the exclusion.
- **The one deviation from fleet parity** (the `$` anchor, D5/H) could surprise a future maintainer expecting byte-identical behaviour. It is documented in the script header and here, and pinned by arm H.
- **Empty-surface vacuosity**: the floor (exit 2) plus AC-M1-9's count pin guard against a scope typo silently scanning nothing.
- Residual: the fleet's `make` wrapper (`make check-no-personal-email`) has no World equivalent because World has no Makefile; CI invokes the script directly (D4). This is intentional, not a loss of coverage.

## §10 Fleet note

World's `scripts/check_no_personal_email.sh` / `scripts/test_check_no_personal_email.sh` mirror the fleet's `f0d44915d` precedent (same address regex, same exclusion set, same per-hit and clean output lines, same runtime-concatenated fixture technique in the self-test). The deliberate differences are: (1) a World-owned scope regex limited to mission docs + `scripts/` (no `.claude/skills/` — World has no such tree); (2) a third exit code `2` for the empty-scope anti-vacuity floor; (3) extra self-test arms F/G/H, including an anchored `$` on `@example\.(com|org|net)` that is a single, deliberate deviation from the fleet's unanchored exclusion; and (4) no Makefile — CI runs the scripts directly in `go-verify`. World's `scripts/mission_answer.sh` adopts the same noreply default the fleet adopted at `8369877d9`, so there is no divergence between the two identity resolutions. Nothing here requires the fleet to change; the gate only makes the fleet's own contract locally enforceable inside World. The controller sends a courtesy note at Gate 5 informing the fleet of the mirror and the one anchoring deviation, in case it wants to adopt the `$` anchor upstream.

