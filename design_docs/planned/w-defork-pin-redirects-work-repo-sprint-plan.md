# Sprint plan: w-defork-pin-redirects-work-repo

**Iteration:** World 166  
**Design:** `design_docs/planned/w-defork-pin-redirects-work-repo.md`  
**Planning base:** `0195b139444739e5d4df9fa4a146818182db9fcf`  
**Executor branch:** `mission/world-iter166-defork-pin`  
**Executor worktree:** `/Users/voightkampff/dev/sunholo-data/.wt-world-iter166`  
**Status:** planned; implementation and external mutation have not started

## Outcome and scope

This is a deliberately small sprint. Its operational outcome is one rig-local line:

```text
AILANG_DRIVER_PIN=0
```

in `/Users/voightkampff/.config/ailang/mission-world.env`, followed by proof from a new controller
fire that World again starts in `/Users/voightkampff/dev/sunholo-data/ailang-world` rather than the
pin worktree of `sunholo-data/ailang`.

The executor cannot make or certify that change. It runs under `codex exec --sandbox
workspace-write`; the env file, launchd, the live controller, and the real sibling clones are
outside its writable worktree. Any executor statement that it changed the env file, exercised
launchd, or certified those outside paths is **UNINFORMATIVE UNDER SANDBOX** and must not be used
as acceptance evidence.

The executor's honest repo-local deliverable is the **charter queue row plus bookkeeping
artifacts** candidate:

1. refresh existing queue row 68 in `design_docs/world-mission.md`, whose claim that the pin was
   “inert today” is now stale; and
2. create `design_docs/planned/w-defork-pin-redirects-work-repo-fleet-issue.md`, the exact issue
   body the controller will file in `sunholo-data/ailang` for the M2 fleet fix.

This is preferable to a `scripts/` check because a repository script cannot honestly certify the
rig-local file or launchd state from the executor sandbox. It is preferable to a new verification
harness because AC-F1/F2/F3/F11 are already short, reproducible commands and the fleet implementation
does not exist in World. A harness here would duplicate the design while still being unable to
certify the external paths. The executor share is therefore small, and the estimate stays small.

World must not edit, temporarily mutate, vendor, or sync any path under `tools/launchd/`. M2 is a
fleet-owned implementation. M3 is later World work and is not started by this sprint.

## Ownership split

| Work | Owner | Writable target | Evidence authority |
|---|---|---|---|
| Refresh queue row 68 | executor | `design_docs/world-mission.md` | executor, repo-local static checks |
| Create the fleet issue body | executor | `design_docs/planned/w-defork-pin-redirects-work-repo-fleet-issue.md` | executor, repo-local static checks; controller reproduces external-path facts before filing |
| Add `AILANG_DRIVER_PIN=0` | controller only | `/Users/voightkampff/.config/ailang/mission-world.env` | controller outside sandbox |
| Observe a new live fire | controller only | launchd/log/live controller session | controller outside sandbox |
| File the fleet issue | controller only | GitHub repository `sunholo-data/ailang` | controller after reviewing the body |
| Implement or test M2 in `tools/launchd/*` | fleet only | fleet repository | future fleet sprint |
| Re-enable the pin for M3 | controller only, later | rig-local env file | only after the hand-back predicate is satisfied |

## Milestone M1A — repo-local hand-off artifacts

**Owner:** executor  
**Estimate:** 0.15 day / 1.2 hours  
**Files:** exactly two: one modified, one created

- Modify `design_docs/world-mission.md` at existing queue row 68. Do not create a duplicate row.
  Preserve its history, but make the current state explicit: the formerly latent condition fired
  after de-fork commit `e92594c`; World uses the fleet driver; the pin moved the controller to an
  `ailang` worktree; M1 is the controller-owned opt-out; M2 is the fleet-owned guard; and World
  never edits `tools/launchd/*` under `D-WORLD-DRIVER-1`.
- Create `design_docs/planned/w-defork-pin-redirects-work-repo-fleet-issue.md`. It is a ready-to-file
  GitHub issue body, not an implementation. It must contain:
  - the observed wrong-repository consequence and the three real origin URLs;
  - the `git-common-dir` counterexample for the separate `ailang-motoko` clone;
  - the requested M2 semantics: authoritative-when-set
    `AILANG_DRIVER_MISSION_IS_DE_FORKED`, normalized-origin inference by default, loud STALE on an
    indeterminate origin, and `_set_pin_workdir` mutating the current shell;
  - the round-2 binding: the production call is the plain
    `_set_pin_workdir "$wt" "$src" || return 1`, never a command substitution;
  - the exact AC-F1 through AC-F12 commands and expected results from the design doc, with AC-F5,
    AC-F10, and AC-F12 called out as the wrong-discriminator, loud-indeterminate, and current-shell
    regression killers; and
  - ownership: the fleet implements and lands the code, while World supplies evidence and does not
    edit frozen core.

### M1A acceptance

All commands run from `/Users/voightkampff/dev/sunholo-data/.wt-world-iter166`.

**AC-E1 — exactly one updated row 68, and its stale “inert today” assertion is gone.**

```bash
export PATH=/opt/homebrew/bin:$PATH
row_count=$(grep -c '^68\. .*w-driver-pin-named-world-points-at-the-wrong-repo' design_docs/world-mission.md || true)
active_count=$(grep -c 'THE LATENT CONDITION FIRED' design_docs/world-mission.md || true)
stale_count=$(grep -c 'It is inert \*\*today\*\* only because' design_docs/world-mission.md || true)
printf 'row=%s active=%s stale=%s\n' "$row_count" "$active_count" "$stale_count"
test "$row_count" -eq 1
test "$active_count" -eq 1
test "$stale_count" -eq 0
```

Expected: `row=1 active=1 stale=0`. The row and active-marker counts are the known-positive controls
for the zero stale-text claim. Named RED mutation: retain the old “inert today” sentence; the last
assertion fails.

**AC-E2 — the fleet issue body is complete enough to file without re-deriving the fix.**

```bash
export PATH=/opt/homebrew/bin:$PATH
f=design_docs/planned/w-defork-pin-redirects-work-repo-fleet-issue.md
test -s "$f"
for needle in \
  'AILANG_DRIVER_MISSION_IS_DE_FORKED' \
  '_set_pin_workdir "$wt" "$src" || return 1' \
  'AC-F1' 'AC-F2' 'AC-F3' 'AC-F4' 'AC-F5' 'AC-F6' \
  'AC-F7' 'AC-F8' 'AC-F9' 'AC-F10' 'AC-F11' 'AC-F12' \
  'D-WORLD-DRIVER-1' 'UNINFORMATIVE UNDER SANDBOX'; do
  grep -Fq "$needle" "$f" || { printf 'missing issue-body needle: %s\n' "$needle" >&2; exit 1; }
done
printf 'issue-body-contract=complete\n'
```

Expected: `issue-body-contract=complete`. Named RED mutation: delete any one required fleet AC label
or the plain-call form; the loop identifies the missing contract element and exits nonzero.

**AC-E3 — scope is exactly the two hand-off artifacts and frozen core is untouched.**

```bash
export PATH=/opt/homebrew/bin:$PATH
base=0195b139444739e5d4df9fa4a146818182db9fcf
git diff --quiet "$base" -- tools/launchd
test ! -e tools/launchd/lib/pin-root.sh
if git diff --quiet "$base" -- design_docs/world-mission.md; then
  echo 'world queue row was not updated' >&2
  exit 1
fi
test -s design_docs/planned/w-defork-pin-redirects-work-repo-fleet-issue.md
changed=$(git status --short -- design_docs/world-mission.md design_docs/planned/w-defork-pin-redirects-work-repo-fleet-issue.md | wc -l | tr -d ' ')
printf 'handoff-paths=%s frozen-core-diff=0\n' "$changed"
test "$changed" -eq 2
```

Expected: `handoff-paths=2 frozen-core-diff=0`. The changed charter and nonempty issue body are the
known-positive controls paired with the zero frozen-core diff. Named RED mutation: edit any tracked
`tools/launchd/*` file; `git diff --quiet` fails.

**AC-E4 — no forbidden placeholders in either artifact.**

```bash
export PATH=/opt/homebrew/bin:$PATH
files='design_docs/world-mission.md design_docs/planned/w-defork-pin-redirects-work-repo-fleet-issue.md'
positive=$(rg -l 'w-defork-pin-redirects-work-repo|w-driver-pin-named-world-points-at-the-wrong-repo' $files | wc -l | tr -d ' ')
bad=$(rg -n 'MILESTONE[_-]ID|T[O]DO|auto-parse[[:space:]]failed' $files || true)
printf 'positive-files=%s forbidden-hits=%s\n' "$positive" "$([ -n "$bad" ] && printf 1 || printf 0)"
test "$positive" -eq 2
test -z "$bad"
```

Expected: `positive-files=2 forbidden-hits=0`. The two positive file matches prove the searched file
set is live. Named RED mutation: insert any forbidden marker; `test -z` fails.

**AC-E5 — static content only; no broad product gate is attributed to this documentation change.**

```bash
export PATH=/opt/homebrew/bin:$PATH
jq -e . design_docs/planned/sprint_w-defork-pin-redirects-work-repo.json >/dev/null
test -s design_docs/planned/w-defork-pin-redirects-work-repo-sprint-plan.md
printf 'planning-artifacts=valid\n'
```

Expected: `planning-artifacts=valid`. No socket, launchd, external-path, Go, or AIL verdict is claimed
by this milestone.

## Milestone M1B — controller-only rig mitigation and proof

**Owner:** controller, outside the executor sandbox  
**Estimate:** 0.10 day / 0.8 hour active work, plus the wait for the next scheduled fire  
**Repository files changed:** none

The executor must stop after M1A and hand the uncommitted worktree to the controller. The controller
performs every command below outside `codex exec --sandbox workspace-write`.

### M1B.0 — capture the live-log baseline before editing

```bash
export PATH=/opt/homebrew/bin:$PATH
before=$(grep -cF 'pin disabled via AILANG_DRIVER_PIN=0' /tmp/ailang-mission-world.log 2>/dev/null || true)
printf '%s\n' "${before:-0}" > /tmp/w-defork-pin-disabled-count.before
printf 'pin-disabled-before=%s\n' "${before:-0}"
```

This makes the later log assertion a delta, not a pass on an old line.

### M1B.1 — exact idempotent edit with pre-edit backup

Backup path:
`/Users/voightkampff/.config/ailang/mission-world.env.pre-w-defork-pin-redirects-work-repo-iter166`

```bash
export PATH=/opt/homebrew/bin:$PATH
env_file=/Users/voightkampff/.config/ailang/mission-world.env
backup=/Users/voightkampff/.config/ailang/mission-world.env.pre-w-defork-pin-redirects-work-repo-iter166
test -f "$env_file" || { echo "missing $env_file" >&2; exit 1; }
if [ ! -e "$backup" ]; then cp -p "$env_file" "$backup"; fi
exact=$(awk '$0 == "AILANG_DRIVER_PIN=0" {n++} END {print n+0}' "$env_file")
conflicts=$(awk '$0 ~ /^[[:space:]]*AILANG_DRIVER_PIN=/ && $0 != "AILANG_DRIVER_PIN=0" {n++} END {print n+0}' "$env_file")
test "$conflicts" -eq 0 || { echo 'conflicting active AILANG_DRIVER_PIN assignment; refusing' >&2; exit 1; }
case "$exact" in
  0) printf '\nAILANG_DRIVER_PIN=0\n' >> "$env_file" ;;
  1) : ;;
  *) echo 'duplicate AILANG_DRIVER_PIN=0 assignments; refusing' >&2; exit 1 ;;
esac
```

Running the block twice does not add the line twice. It refuses rather than guessing if an active,
conflicting assignment or a pre-existing duplicate is present. The backup is never overwritten by a
rerun.

Exact revert command:

```bash
export PATH=/opt/homebrew/bin:$PATH
cp -p /Users/voightkampff/.config/ailang/mission-world.env.pre-w-defork-pin-redirects-work-repo-iter166 /Users/voightkampff/.config/ailang/mission-world.env
```

Do not use the revert until rollback is actually required. M3 later removes or comments the opt-out
only after the fleet hand-back predicate below is met.

### M1B acceptance

**AC-C1 — the active opt-out occurs exactly once and the backup exists.**

```bash
export PATH=/opt/homebrew/bin:$PATH
env_file=/Users/voightkampff/.config/ailang/mission-world.env
backup=/Users/voightkampff/.config/ailang/mission-world.env.pre-w-defork-pin-redirects-work-repo-iter166
count=$(awk '$0 == "AILANG_DRIVER_PIN=0" {n++} END {print n+0}' "$env_file")
printf 'pin-optout=%s backup=%s\n' "$count" "$([ -f "$backup" ] && printf present || printf absent)"
test "$count" -eq 1
test -f "$backup"
```

Expected: `pin-optout=1 backup=present`. Base measurement was `pin-optout=0`, with the existing
`MISSION_WORKDIR=` line found once as the same-file positive control. Named RED mutation: delete or
comment the opt-out; `count` becomes 0.

**AC-C2 — re-running the edit is byte-idempotent.**

```bash
export PATH=/opt/homebrew/bin:$PATH
env_file=/Users/voightkampff/.config/ailang/mission-world.env
backup=/Users/voightkampff/.config/ailang/mission-world.env.pre-w-defork-pin-redirects-work-repo-iter166
before=$(shasum -a 256 "$env_file" | awk '{print $1}')
test -f "$env_file" || { echo "missing $env_file" >&2; exit 1; }
if [ ! -e "$backup" ]; then cp -p "$env_file" "$backup"; fi
exact=$(awk '$0 == "AILANG_DRIVER_PIN=0" {n++} END {print n+0}' "$env_file")
conflicts=$(awk '$0 ~ /^[[:space:]]*AILANG_DRIVER_PIN=/ && $0 != "AILANG_DRIVER_PIN=0" {n++} END {print n+0}' "$env_file")
test "$conflicts" -eq 0 || { echo 'conflicting active AILANG_DRIVER_PIN assignment; refusing' >&2; exit 1; }
case "$exact" in
  0) printf '\nAILANG_DRIVER_PIN=0\n' >> "$env_file" ;;
  1) : ;;
  *) echo 'duplicate AILANG_DRIVER_PIN=0 assignments; refusing' >&2; exit 1 ;;
esac
after=$(shasum -a 256 "$env_file" | awk '{print $1}')
printf 'before=%s after=%s\n' "$before" "$after"
test "$before" = "$after"
```

Expected: identical hashes. Named RED mutation: replace the guarded append with an unconditional
append; the second hash differs and AC-C1 also reads 2.

**AC-C3 — sourcing the real env file yields the opt-out and preserves the plist work directory.**

```bash
export PATH=/opt/homebrew/bin:$PATH
MISSION_WORKDIR=/Users/voightkampff/dev/sunholo-data/ailang-world /bin/bash -c '. /Users/voightkampff/.config/ailang/mission-world.env; printf "PIN=%s WD=%s\n" "${AILANG_DRIVER_PIN:-<unset>}" "$MISSION_WORKDIR"' | grep -Fx 'PIN=0 WD=/Users/voightkampff/dev/sunholo-data/ailang-world'
```

Expected: the exact line shown. Named RED mutation: delete the opt-out; output contains
`PIN=<unset>` and the exact grep fails.

**AC-C4 — the sourced opt-out reaches the real fleet helper's disabled branch without re-exec.**

```bash
export PATH=/opt/homebrew/bin:$PATH
env -u AILANG_DRIVER_PINNED -u AILANG_DRIVER_SRC -u AILANG_DRIVER_DRIFT -u AILANG_DRIVER_REF -u AILANG_DRIVER_PIN_GATE_REFRESHED MISSION_WORKDIR=/Users/voightkampff/dev/sunholo-data/ailang-world /bin/bash -c '. /Users/voightkampff/.config/ailang/mission-world.env; [ "${AILANG_DRIVER_PIN:-}" = 0 ] || { echo "env opt-out absent" >&2; exit 1; }; . /Users/voightkampff/dev/sunholo-data/ailang/tools/launchd/lib/pin-root.sh; pin_root_to_committed_ref; printf "PIN_STATUS=%s WD=%s\n" "$PIN_STATUS" "$MISSION_WORKDIR"' | grep -Fx 'PIN_STATUS=disabled WD=/Users/voightkampff/dev/sunholo-data/ailang-world'
```

Expected: the exact line shown. The precondition exits before calling the helper when the env-file
mutation is absent, so this criterion does not accidentally initiate a pin/re-exec. Named RED
mutation: delete the env line; it exits with `env opt-out absent`.

**AC-C5 — a new fire, not historical log residue, reports the opt-out.**

After the next scheduled World fire:

```bash
export PATH=/opt/homebrew/bin:$PATH
before=$(cat /tmp/w-defork-pin-disabled-count.before)
after=$(grep -cF 'pin disabled via AILANG_DRIVER_PIN=0' /tmp/ailang-mission-world.log 2>/dev/null || true)
printf 'pin-disabled-before=%s after=%s\n' "$before" "${after:-0}"
test "${after:-0}" -gt "$before"
```

Expected: `after` is strictly greater than `before`. Named RED mutation: remove the opt-out before
the new fire; no new disabled marker is written and the comparison fails. This replaces the design
doc's `>= 1` check, which could pass on a historical marker.

**AC-C6 — the newly fired controller is in World's repository and can see its charter.**

Run inside the newly fired controller session before it creates a sprint worktree:

```bash
export PATH=/opt/homebrew/bin:$PATH
actual=$(pwd -P)
origin=$(git config --get remote.origin.url)
test "$actual" = /Users/voightkampff/dev/sunholo-data/ailang-world
test "$origin" = https://github.com/sunholo-data/ailang-world.git
test -f design_docs/world-mission.md
printf 'cwd=%s origin=%s charter=present\n' "$actual" "$origin"
```

Expected: `cwd=/Users/voightkampff/dev/sunholo-data/ailang-world`, the World origin, and
`charter=present`. Named RED mutation: remove the opt-out; the current pin behavior starts in
`/Users/voightkampff/.ailang-driver-pin/world`, whose measured origin is `sunholo-data/ailang` and
which lacks `design_docs/world-mission.md`.

Only AC-C1 through AC-C6 can accept M1B. Executor-side results for these commands are
**UNINFORMATIVE UNDER SANDBOX**.

## Controller filing step after M1A review

The controller, not the executor, reproduces the external-path discriminator inputs outside the
sandbox and then files the prepared body:

```bash
export PATH=/opt/homebrew/bin:$PATH
for repo in /Users/voightkampff/dev/sunholo-data/ailang /Users/voightkampff/dev/sunholo-data/ailang-motoko /Users/voightkampff/dev/sunholo-data/ailang-world; do
  printf 'repo=%s origin=%s common=%s\n' "$repo" "$(git -C "$repo" config --get remote.origin.url)" "$(git -C "$repo" rev-parse --path-format=absolute --git-common-dir)"
done
gh issue create --repo sunholo-data/ailang --title 'pin-root: preserve a de-forked mission work repository across driver pin re-exec' --body-file design_docs/planned/w-defork-pin-redirects-work-repo-fleet-issue.md
```

Expected measurements before filing:

- `ailang` and `ailang-motoko` both have origin
  `https://github.com/sunholo-data/ailang.git`, while their absolute common dirs differ;
- `ailang-world` has origin `https://github.com/sunholo-data/ailang-world.git` and its own common
  dir.

If those inputs change, stop and update the issue evidence rather than filing stale facts. Filing
the issue is an external side effect reserved to the controller.

## Gate 3b posture

Gate 3b is expected to report two inherited red checks on this documentation-only branch:

- `go host build + test gate` invokes missing `tools/launchd/test_mission_routing.sh` through
  `scripts/verify_go.sh`;
- `launchd drivers (bash 3.2)` invokes deleted driver tests and cannot read the deleted
  `tools/launchd/mission-control.sh`.

The latest measured `dev` run at planning time was GitHub Actions run `34033464096`, head
`9166de05c42d0c3430d0ebb87c99db378b72420c`: `ailang-code verify gate` succeeded and those two
jobs failed. This predates the branch and is parked by the controller-supplied `D-WORLD-34`
instruction. The executor and controller must not fix, mask, vendor, or reintroduce deleted
`tools/launchd/*` files. Gate 3b records the two failures as inherited; it does not call the sprint
green and does not reinterpret them as failures caused by the hand-off artifacts.

## Hand-back predicate for M3 — later, not this sprint

M3 may re-enable the pin only after all of the following are true:

1. a fleet-authored commit implementing the design's M2 guard has landed in
   `sunholo-data/ailang`;
2. fleet AC-F4 through AC-F12 pass outside the executor sandbox, including AC-F5 against the
   separate motoko clone, AC-F10's caller-visible `PIN_STATUS=STALE`, and AC-F12's no-command-
   substitution call site;
3. with the opt-out removed, a fresh World fire reports `PIN_STATUS=pinned` while AC-C6 still
   reports the World cwd, World origin, and present charter.

Until all three hold, `AILANG_DRIVER_PIN=0` remains active. M3 is not estimated or assigned here.

## Planner Findings

1. **P1 — the planning base differs from the design's measured code base only by the design doc.**
   `export PATH=/opt/homebrew/bin:$PATH; git diff --name-status 9166de05c42d0c3430d0ebb87c99db378b72420c 0195b139444739e5d4df9fa4a146818182db9fcf`
   prints only `A design_docs/planned/w-defork-pin-redirects-work-repo.md`. The JSON therefore uses
   `0195b139...`, not the sibling sprint's base or the design's earlier measurement commit.
2. **P2 — AC3 in the design is vacuous as an acceptance criterion for the named env-file
   mutation.** It sets `AILANG_DRIVER_PIN=0` directly on the command line, so deleting the env-file
   line cannot turn it red. Resolution: AC-C4 sources the real env file and exits before the helper
   call if the value is absent. The original direct helper probe remains useful premise evidence,
   but it is not M1 acceptance.
3. **P3 — AC4's `grep -c ... >= 1` can pass on history.** No named current mutation can remove an
   old matching log line. Resolution: M1B.0 captures the pre-edit count and AC-C5 requires a strict
   increase after a new fire; AC-C6 separately proves the live cwd, origin, and charter.
4. **P4 — the committed charter at this base contains `D-WORLD-DRIVER-1`, `D-WORLD-32`, and
   `D-WORLD-33`, but no `D-WORLD-34`.** The command
   `export PATH=/opt/homebrew/bin:$PATH; rg -n 'D-WORLD-34' design_docs/world-mission.md` returns no
   hit, while `rg -n 'D-WORLD-DRIVER-1' design_docs/world-mission.md` is the same-file positive
   control. The controller explicitly supplied `D-WORLD-34` and its park/land-no-code disposition
   for this sprint, so this plan treats it as binding external context but does not invent or edit a
   human-decision ledger row.
5. **P5 — queue row 68 is the correct queue artifact; adding a second row would duplicate the
   incident.** At base, `grep -c '^68\. .*w-driver-pin-named-world-points-at-the-wrong-repo'
   design_docs/world-mission.md` returns 1. Its “inert today” premise is now false, so the executor
   refreshes that row in place and preserves its history.
6. **P6 — the controller's three origin/common-dir measurements reproduce exactly.** The command
   in the filing step returns origins `ailang`, `ailang`, `ailang-world` and three different
   absolute common dirs. This confirms origin identity discriminates the de-fork while common-dir
   identity would misclassify the separate motoko clone. These are planner measurements; an
   executor repeat remains **UNINFORMATIVE UNDER SANDBOX**.
7. **P7 — AC-F11's normalization examples reproduce under Bash 3.2 semantics.** Applying the
   design's canonicalizer to `git@github.com:sunholo-data/ailang.git` and
   `https://github.com/sunholo-data/ailang.git` yields the same
   `github.com/sunholo-data/ailang`; the World URL yields
   `github.com/sunholo-data/ailang-world`. No git repository mutation was needed for this planner
   probe.
8. **P8 — the rig-local base is RED exactly where M1 needs it to be.**
   `grep -c '^AILANG_DRIVER_PIN=0$' /Users/voightkampff/.config/ailang/mission-world.env` returned
   0, while the same file's known-positive `grep -c '^MISSION_WORKDIR=' ...` returned 1. The file
   has 109 lines. No planner write was attempted.
9. **P9 — the helper honors the proposed value, but this is premise evidence only.** A direct
   read-only probe with inherited pin variables removed returned
   `PIN_STATUS=disabled MISSION_WORKDIR=/Users/voightkampff/dev/sunholo-data/ailang-world`.
   AC-C4 is stronger because it obtains the value from the file that M1 actually mutates.
10. **P10 — the live wrong-repository discriminator is observable now.** The pin worktree's origin
    is `https://github.com/sunholo-data/ailang.git`; its World charter is absent; the real World
    clone's World charter is present. AC-C6 binds all three facts so a mere disabled log line cannot
    produce a false green.
11. **P11 — current Gate 3b red is inherited and matches the controller premise.** Read-only GitHub
    inspection of run `34033464096` found one successful job and exactly the two named failed jobs.
    Failed logs name missing `tools/launchd/test_mission_routing.sh` and missing
    `tools/launchd/mission-control.sh`. This sprint makes no attempt to repair them.
12. **P12 — the design's Files to Create/Modify list omits the required World-owned M2 hand-off
    artifacts.** The user scope requires a queue row and issue body. Resolution: M1A is limited to
    those two artifacts; it does not widen into a World implementation of M2.

## Estimate and exit

Total active work is **0.25 day / 2.0 hours**: 1.2 executor hours for the two repo-local hand-off
artifacts and their static checks, plus 0.8 controller hour for backup, idempotent edit, direct
checks, next-fire verification, and issue filing. Calendar completion also waits for the next World
fire. Padding this into a code sprint would create false authority over the only operational change.

The sprint exits when M1A is handed off, AC-C1 through AC-C6 pass outside sandbox, and the fleet
issue is filed. M2 implementation and M3 re-enable remain open beyond this sprint.
