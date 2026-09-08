# w-defork-pin-redirects-work-repo — stop the pin from moving the World controller into the wrong repository

- **Status:** IMPLEMENTED / CLOSED 2026-09-08 (World iteration 170). M1 (the rig-local
  `AILANG_DRIVER_PIN=0` stopgap) is SUPERSEDED by M3; M2 landed as a FLEET-authored commit,
  `8c56c5863` *"fix(mission): preserve foreign work repository across pinned driver re-exec"*
  (2026-09-07 11:13 +0200, on the fleet's `origin/dev`, and an ANCESTOR of the pinned ref
  `48c4a6e49` this rig actually runs); M3 (re-enable the pin) was deployed attended on 2026-09-07
  and is confirmed HOLDING first-party this fire — pin ACTIVE (`AILANG_DRIVER_PIN=1`,
  `AILANG_DRIVER_PINNED=48c4a6e49`) with `MISSION_WORKDIR` and `pwd` both
  `/Users/voightkampff/dev/sunholo-data/ailang-world`, while the pin worktree
  `~/.ailang-driver-pin/world-ollama-gauge-48c4a6e49` still has origin
  `github.com/sunholo-data/ailang` — the wrong repo the loop is no longer in.
  **AC-F1 through AC-F12: 12/12 PASS**, verified by the controller at iteration 170.
  ⚠ **Run them under `bash`.** The controller's first pass ran them in the login shell (zsh) and
  got `rc=2` ("indeterminate origin identity") on AC-F4/F5/F11 — a pure INSTRUMENT artifact, since
  `pin-root.sh` is bash. Re-run under `bash -c`, every AC matched its stated expectation exactly.
- **Date:** 2026-09-07
- **Owning queue row:** new — surfaced by the controller this fire (iter-166); the de-fork (commit `e92594c`) deleted World's own driver copy and the shared driver's pin now redirects `MISSION_WORKDIR` into a worktree of the **wrong** repository.
- **Base commit measured at:** `9166de0` (== `origin/dev`), worktree `.wt-world-iter166`
- **Fleet source measured at:** `sunholo-data/ailang` HEAD `59b231907` (clean working tree, 0 behind `origin/dev`)

---

## Motivation

The World mission was **de-forked** (commit `e92594c` deleted its own 1,047-line copy of
`tools/launchd/mission-control.sh`). The launchd job `dev.ailang.mission-world` was regenerated
to invoke the **FLEET** driver:

`ProgramArguments = /bin/bash .../ailang/tools/launchd/mission-control.sh`; `MISSION_WORKDIR =
/Users/voightkampff/dev/sunholo-data/ailang-world`; `WorkingDirectory = .../ailang-world`.

The fleet driver deliberately separates two roots (`mission-control.sh:40-48`):

```bash
# MC_DRIVER_ROOT is the repo the DRIVER ships from; REPO is the repo the mission WORKS in.
# For v1/docs/motoko they are the same checkout; for a de-forked mission they are NOT.
MC_DRIVER_ROOT=$(cd "$(dirname "$0")/../.." && pwd)
REPO="${MISSION_WORKDIR:-$MC_DRIVER_ROOT}"
cd "$REPO" || exit 1
```

and its pin block (`mission-control.sh:894-898`) states the intended contract — "the helper
then derives what to pin from its own $0, and MISSION_WORKDIR keeps $REPO pointing at the
mission's work repo across the re-exec". **THAT LAST CLAUSE IS FALSE.**
`tools/launchd/lib/pin-root.sh` unconditionally does (`pin-root.sh:263-267`):

```bash
AILANG_DRIVER_PINNED="$short"; AILANG_DRIVER_SRC="$src"; AILANG_DRIVER_DRIFT="$drift"
AILANG_DRIVER_REF="$ref"; MISSION_WORKDIR="$wt"
export AILANG_DRIVER_PINNED AILANG_DRIVER_SRC AILANG_DRIVER_DRIFT AILANG_DRIVER_REF MISSION_WORKDIR
exec /bin/bash "$wt/tools/launchd/$script" "$@"
```

where `$wt` is a worktree of the **DRIVER's** source clone (`$src`, derived from `$0`), i.e. of
`sunholo-data/ailang`. Its own comment defends the clobber (`pin-root.sh:257-262`):

```
# MISSION_WORKDIR moves too, and that is the point: mission-control.sh:40 reads it AHEAD of
# $0-relative resolution, so pinning only the script would leave motoko and World running a
# pinned driver against a stale charter and skill. That half-fix reports green, which is worse
# than no fix. Sprint worktrees are unaffected (the skill creates them by absolute path).
```

**MEASURED CONSEQUENCE, this fire (2026-09-07 08:44 local, driver log
`/tmp/ailang-mission-world.log`):** `driver pin: running committed origin/dev @ 878939117 from
/Users/voightkampff/.ailang-driver-pin/world (source clone .../ailang was 0 behind)`. And inside
the controller session (verified first-party this fire):

```
MISSION_WORKDIR = /Users/voightkampff/.ailang-driver-pin/world
pwd             = /Users/voightkampff/.ailang-driver-pin/world
git remote -v   -> origin https://github.com/sunholo-data/ailang.git   (NOT ailang-world)
git rev-parse --git-common-dir -> /Users/voightkampff/dev/sunholo-data/ailang/.git
ls design_docs/world-mission.md -> No such file or directory
```

So the World mission's controller now executes inside a throwaway worktree of the **WRONG
REPOSITORY**. Its charter, decision ledger, mission log, iteration index and dashboard are all
in `ailang-world` and are unreachable from `$REPO`. The pin reports **FULL SUCCESS** while doing
this: `PIN_STATUS=pinned`, drift 0, no warning on any channel. The environment of this very
design session confirms it — `MISSION_WORKDIR=/Users/voightkampff/.ailang-driver-pin/world`,
`AILANG_DRIVER_PINNED=878939117`, `AILANG_DRIVER_SRC=.../ailang`.

**The second layer.** `~/.config/ailang/mission-world.env` (a RIG-LOCAL file, not in any repo)
contains `MISSION_WORKDIR="${MISSION_WORKDIR:-/Users/voightkampff/dev/sunholo-data/ailang-world}"`,
changed from a bare assignment on 2026-08-18 so that "the pin wins" — right pre-de-fork, now
defers to the wrong value.

---

## Premises (hard constraints)

1. **P1 — D-WORLD-DRIVER-1 (RESOLVED, ratified attended by Mark 2026-08-17):** the driver stays
   FLEET-owned; the sync is a COMMIT, never working-tree dirt. Updates to `tools/launchd/*` land
   in this repo only as fleet-authored commits. World's controller never edits them. The durable
   fix to `tools/launchd/lib/pin-root.sh` is therefore **FLEET work**, not World work.
2. **P2 — World owns and may change:** its own repo `sunholo-data/ailang-world` (charter, log,
   queue, `scripts/`, CI workflow) and the RIG-LOCAL file `~/.config/ailang/mission-world.env`.
   Precedent: the iter-23 PATH fix and the D-WORLD-20 executor-fallback pin both landed there
   EXPRESSLY because `tools/launchd/*` is frozen core. That file is sourced at driver line 71,
   before the pin call at line ~894, so `AILANG_DRIVER_PIN` set there DOES take effect.
3. **P3 — World must not edit `tools/launchd/*`.** Any design that requires it goes in the FLEET
   half and says so.
4. **P4 — The env file cannot repair `REPO`.** It is sourced at line 71, after `REPO` is computed
   at line 48 and after `cd "$REPO"`. Forcing `MISSION_WORKDIR` there changes the variable but
   not the already-computed `REPO` or the already-executed `cd`. The only env-file lever that
   works is one that prevents the pin from running at all (`AILANG_DRIVER_PIN=0`).
5. **P5 — The pin worktree is always a worktree of `$src` (the driver's source clone).** For a
   de-forked mission `$src` is `ailang`, so `$wt` is a worktree of `ailang` — never of
   `ailang-world`. Redirecting `MISSION_WORKDIR` to it is therefore always wrong for world.
6. **P6 — The pin's original purpose is already broken for world.** The pin was designed to move
   script + skill + charter together in one move. For world those now live in two different
   repositories, so the pin can only ever pin the DRIVER; it cannot move the charter/skill. The
   durable fix must therefore keep the driver-pin benefit (committed driver code) while dropping
   the `MISSION_WORKDIR` redirect for world.

---

## Decision 1 — Immediate mitigation: `AILANG_DRIVER_PIN=0` in the rig-local env file (WORLD)

**Decision:** add `AILANG_DRIVER_PIN=0` to `~/.config/ailang/mission-world.env`. This disables the
pin for the World mission only. The driver then skips `pin_root_to_committed_ref`'s pin path
(`PIN_STATUS=disabled`, verified this fire), `MISSION_WORKDIR` stays at the plist value
`/Users/voightkampff/dev/sunholo-data/ailang-world`, and `REPO` resolves to `ailang-world` again.
The controller runs in its own repository.

**Why this is the right immediate move:**
- It is the **only** env-file lever that works (P4): no other `AILANG_DRIVER_*` variable controls
  the `MISSION_WORKDIR` redirect specifically; `AILANG_DRIVER_PIN_DIR` only changes WHERE the
  worktree lives, not WHICH repo it is a worktree of (P5).
- It is **World-landable** (P2): a rig-local file, no git write, no `tools/launchd/*` edit.
- It is **trivially reversible**: comment out or delete one line. Reversing is `# AILANG_DRIVER_PIN=0`.
- It does **not** make the fleet fix harder to land or verify: it is orthogonal to `pin-root.sh`.
  When the fleet fix lands, World un-comments the line and the durable fix takes over.

**Rejected alternative — detect-and-relocate at Gate 1 (controller discipline):** keep the pin,
and have the controller, at Gate 1, detect it is in the wrong repo and `cd` to `ailang-world`.
Rejected because (a) the controller session is spawned with cwd = `$REPO` = the wrong repo, and
its first actions read the charter, which is unreachable — relocation must happen BEFORE any
read, the fragile ordering the mission has repeatedly found unreliable (row 69); (b) any write
before relocation lands in the wrong repo.

**Rejected alternative — do nothing, park every fire:** this is the current state (iters 162-164
parked de-fork repair). It keeps the mission stalled and leaves the pin reporting false success
with no warning on any channel. Not a fix.

**Acknowledged cost:** `AILANG_DRIVER_PIN=0` re-opens the unpinned-driver class #558 for this
mission — the driver runs the working tree of `ailang` as-is rather than a pinned commit. This
is acceptable as a stopgap because (a) the driver is now the FLEET driver, which the fleet
actively maintains and syncs by commit (P1), and (b) measured this fire: `ailang` is clean and 0
behind `origin/dev`.

---

## Decision 2 — Durable fix (FLEET): pin-root.sh must not redirect a mission whose work repo is a DIFFERENT repository from the driver's source clone

**Decision:** hand to the fleet a precise change to `tools/launchd/lib/pin-root.sh`: keep the
driver pin (re-exec from `$wt`, committed driver code) but make the `MISSION_WORKDIR="$wt"`
clobber conditional on the mission's work repo being the **same repository** as `$src`.

**The discriminator — an explicit flag, authoritative when set, with a NORMALISED origin-URL
inference as the default.** The fleet fix adds one env flag, `AILANG_DRIVER_MISSION_IS_DE_FORKED`
(fleet namespace `AILANG_DRIVER_*`), and makes the `MISSION_WORKDIR="$wt"` clobber conditional on
the mission's work repo being the **same repository** as `$src`. The flag is authoritative WHEN
SET; when unset, the helper infers from the normalised `origin` URL. No existing mission must be
reconfigured.

**Why the flag is authoritative-when-set, not mandatory.** A flag every mission MUST set is a flag
that silently breaks the mission that forgets it. Authoritative-when-set with inference as the
default means a mission that forgets the flag still gets the correct answer from the normalised
URL, and a mission whose URL is ambiguous (or whose operator wants to override a mis-inference)
can force the answer. Setting the flag also removes the indeterminate case.

**Why `git-common-dir` is the WRONG discriminator (measured this fire):** v1/docs/motoko's
`MISSION_WORKDIR` is a **separate CLONE of the same repo** (`ailang`), so it has its OWN `.git`.
Resolved to realpath, `ailang-motoko`'s common dir is `.../ailang-motoko/.git` while `ailang`'s
is `.../ailang/.git` — they DIFFER even though the repo is the same. A `git-common-dir`
comparison would therefore classify motoko/docs as "different repository" and wrongly drop the
redirect, breaking them. (Verified: AC-F1/F2/F3 below.)

**Why the clobber must be KEPT for v1/docs/motoko and DROPPED for world:** for motoko/docs the
pin worktree (a worktree of `$src`) is the same repository, so redirecting `MISSION_WORKDIR` to
`$wt` moves the charter and skill together with the driver — the whole point of the pin. For
world, `MISSION_WORKDIR` is a DIFFERENT repository (`ailang-world`); redirecting it to `$wt` (a
worktree of `ailang`) strands the charter, log, iteration index and dashboard in an unreachable
repo while the pin reports full success.

**The normalisation (answers the URL-fragility objection).** Exact string equality across
`git@github.com:org/repo.git` vs `https://github.com/org/repo.git` vs a trailing slash is unsafe,
so the inference compares a CANONICAL form, not the raw URL. `_pin_origin_canon` applies, in
order: (1) strip a trailing `.git`; (2) strip a trailing `/`; (3) rewrite
`scheme://[user@]host[:port]/path` and `[user@]host:path` to `host/path` (drop scheme, user and
`:port`); (4) lowercase the HOST but NOT the path. Applied to the three real rig values:
- `https://github.com/sunholo-data/ailang.git` → `github.com/sunholo-data/ailang`
- `git@github.com:sunholo-data/ailang-world.git` → `github.com/sunholo-data/ailang-world`
- `https://github.com/sunholo-data/ailang-world.git` → `github.com/sunholo-data/ailang-world`

So world (`.../ailang-world`) and the driver (`.../ailang`) differ, while SSH and HTTPS clones of
the SAME repo canonicalise to the same `host/path` and are correctly "same" (AC-F11).

**The indeterminate case is LOUD, not defaulted (answers the silent-fallback objection).** When
the flag is unset and the inference cannot decide — `origin` unset or empty on either side, `git`
failing, or a non-repo path — `_pin_same_repo` returns 2 (indeterminate) and `_set_pin_workdir`
calls `_pin_stale "<reason>"`: `PIN_STATUS=STALE`, rc=1, and the DRIVER posts `driver ran
UNPINNED (code provenance unknown)` on the human channel (this mission received one such notice
this very fire). We choose `_pin_stale` over a stderr WARN because the harms are asymmetric: a
wrong redirect strands world's charter silently, while a refused pin runs uncommitted driver code
loudly.

**Exact fleet spec (factor the decision into testable pure functions):**

```bash
# --- DE-FORK GUARD (fleet fix for world; D-WORLD-DRIVER-1) ---
# A de-forked mission (world) runs this shared driver out of the ailang repo while its
# charter, log and worktrees live in a DIFFERENT repository (ailang-world). Redirecting
# MISSION_WORKDIR to the pin worktree — a worktree of $src (ailang) — would move the
# controller into the wrong repo and strand its charter, log and skill. For v1/docs/motoko,
# MISSION_WORKDIR is a separate CLONE of the SAME repo, so redirecting is correct and is
# the whole point (the charter moves with the driver).
#
# Discriminator: AILANG_DRIVER_MISSION_IS_DE_FORKED is authoritative when set; otherwise
# infer from the NORMALISED origin URL. NOT git-common-dir: a separate clone of the same
# repo (motoko/docs) has its OWN .git, so git-common-dir differs from $src even though the
# repo is the same — that would wrongly drop the redirect for motoko/docs.
_pin_origin_canon() {  # $1 = dir -> canonical "host/path" on stdout, rc=0; rc=1 if unset/empty/not a repo
  local url host path
  url=$(git -C "$1" config --get remote.origin.url 2>/dev/null) || return 1
  [ -n "$url" ] || return 1
  url=${url%.git}; url=${url%/}          # strip trailing .git and trailing /
  case "$url" in
    *://*) path=${url#*://}; path=${path#*@}; host=${path%%/*}; path=${path#*/}; host=${host%%:*} ;;
    *@*:*) host=${url%%:*}; host=${host#*@}; path=${url#*:} ;;
    *) return 1 ;;
  esac
  host=$(printf '%s' "$host" | tr '[:upper:]' '[:lower:]')   # lowercase host, NOT path
  printf '%s/%s' "$host" "$path"
}

_pin_same_repo() {  # $1 = dir A, $2 = dir B -> 0 same, 1 different, 2 indeterminate
  local a b
  a=$(_pin_origin_canon "$1") || return 2
  b=$(_pin_origin_canon "$2") || return 2
  [ "$a" = "$b" ]
}

# Sets MISSION_WORKDIR IN THE CURRENT SHELL: the pin worktree for a same-repo mission, the
# mission's own repo for a de-forked (different-repo) mission. Indeterminate -> STALE, rc=1.
#
# MUTATES rather than echoes, DELIBERATELY (quorum round 2, gpt6-astra + gemini-3-1-pro, both
# REJECT on the same mechanism): the round-2 draft called this helper inside a command
# substitution, which runs in a SUBSHELL — so `_pin_stale`'s `PIN_STATUS=STALE` was discarded on
# return and the driver saw an unchanged PIN_STATUS, destroying the very loud-failure guarantee
# this design rests on. `MISSION_WORKDIR="$(...)"` also assigns the captured stdout BEFORE
# `|| return 1` runs, so a failed call would first blank the work directory. Assigning in the
# current shell removes BOTH hazards by construction: there is no subshell and no capture.
_set_pin_workdir() {  # $1 = wt, $2 = src -> sets MISSION_WORKDIR; rc 0 ok, 1 indeterminate (STALE)
  local wd rc
  wd="${MISSION_WORKDIR:-$2}"
  if [ -n "${AILANG_DRIVER_MISSION_IS_DE_FORKED:-}" ]; then
    if [ "$AILANG_DRIVER_MISSION_IS_DE_FORKED" = "1" ]; then
      MISSION_WORKDIR="$wd"; return 0        # de-forked: keep the mission's own repo
    fi
    MISSION_WORKDIR="$1"; return 0           # same repo: redirect to the pin worktree
  fi
  _pin_same_repo "$wd" "$2"; rc=$?
  case "$rc" in
    0) MISSION_WORKDIR="$1"; return 0 ;;
    1) MISSION_WORKDIR="$wd"; return 0 ;;
    2) printf 'WARN: pin de-fork guard: cannot canonicalise the origin URL of %s or %s\n' "$wd" "$2" >&2
       _pin_stale "cannot decide whether MISSION_WORKDIR ($wd) is the same repository as the driver source ($2): origin unset/empty or not a git repo on one side — refusing to redirect silently"
       return 1 ;;
  esac
}
```

and replace the unconditional clobber at `pin-root.sh:267`:

```bash
  _set_pin_workdir "$wt" "$src" || return 1
```

The call is a PLAIN function call, not a command substitution, so `_pin_stale` runs in the
CURRENT shell and its `PIN_STATUS=STALE` reaches `pin_root_to_committed_ref` and the driver;
`|| return 1` then propagates the refusal, and no value is written to `MISSION_WORKDIR` on the
failure path. The `AILANG_DRIVER_PINNED/SRC/DRIFT/REF` exports and the
`exec /bin/bash "$wt/tools/launchd/$script"` are unchanged — the driver is still pinned to
committed code; only the work-repo redirect becomes conditional. `_pin_origin_canon` and
`_pin_same_repo` stay pure and testable in isolation; `_set_pin_workdir` is testable by calling
it and reading `$MISSION_WORKDIR` afterwards (no exec, no fire).

**Edge cases the discriminator handles:**
- `MISSION_WORKDIR` unset (v1 plist sets none; the env file sets it to `ailang`): inference compares
  `$src` to `$src` → same → redirect, correct.
- `MISSION_WORKDIR` a worktree or separate clone of the same repo (motoko/docs): origin matches →
  redirect, correct.
- `MISSION_WORKDIR` a different repo (world): origin differs → no redirect, correct.
- `MISSION_WORKDIR` has no origin (local-only), flag unset: `_pin_same_repo` returns 2 →
  `_pin_stale` → STALE, LOUD, no silent branch. A mission that sets the flag avoids this entirely.

---

## Milestones

- **M1 (WORLD, this iteration):** land `AILANG_DRIVER_PIN=0` in `~/.config/ailang/mission-world.env`.
  Restores `REPO=ailang-world` immediately. Reversible by commenting one line.
- **M2 (FLEET):** land the `pin-root.sh` de-fork guard (Decision 2) as a fleet-authored commit.
- **M3 (WORLD, after M2 lands):** un-comment `AILANG_DRIVER_PIN` (re-enable the pin) and verify
  the durable fix holds: driver pinned, controller in `ailang-world`. This is the hand-back gate.

---

## Files to Create/Modify

- `~/.config/ailang/mission-world.env` — **WORLD** — add `AILANG_DRIVER_PIN=0` (rig-local, no git).
- `tools/launchd/lib/pin-root.sh` — **FLEET** — add `_pin_same_repo` / `_set_pin_workdir` and make
  the `MISSION_WORKDIR` clobber conditional (Decision 2). World does NOT touch this.
- `design_docs/planned/w-defork-pin-redirects-work-repo.md` — **WORLD** — this document.

---

## Acceptance Criteria

Each AC is a command with its expected output. AC1–AC3 verify the WORLD mitigation and are
runnable on this rig now (no launchd fire). AC-F1–F3 verify the discriminator logic and are
runnable now. AC-F4–F11 verify the fleet fix and require the fixed `pin-root.sh` to be present
(they source it and call its helpers directly — no fire).

**Mitigation (WORLD):**

- **AC1** — the env file carries the pin opt-out.
  `grep -n '^AILANG_DRIVER_PIN=0' ~/.config/ailang/mission-world.env`
  → a line `AILANG_DRIVER_PIN=0`. (Red at base: the line is absent.)

- **AC2** — sourcing the env file with the plist's `MISSION_WORKDIR` yields the opt-out and keeps
  the work repo.
  `MISSION_WORKDIR=/Users/voightkampff/dev/sunholo-data/ailang-world bash -c '. ~/.config/ailang/mission-world.env; echo "PIN=${AILANG_DRIVER_PIN:-<unset>} WD=$MISSION_WORKDIR"'`
  → `PIN=0 WD=/Users/voightkampff/dev/sunholo-data/ailang-world`. (Red at base: `PIN=<unset>`.)

- **AC3** — `pin-root.sh` honors the opt-out and does NOT re-exec (returns `disabled`).
  ```
  env -u AILANG_DRIVER_PINNED -u AILANG_DRIVER_SRC -u AILANG_DRIVER_DRIFT -u AILANG_DRIVER_REF -u AILANG_DRIVER_PIN_GATE_REFRESHED \
    AILANG_DRIVER_PIN=0 MISSION_WORKDIR=/Users/voightkampff/dev/sunholo-data/ailang-world \
    bash -c '. /Users/voightkampff/dev/sunholo-data/ailang/tools/launchd/lib/pin-root.sh; pin_root_to_committed_ref; echo "PIN_STATUS=$PIN_STATUS"'
  ```
  → `PIN_STATUS=disabled` and the process returns to the `echo` (no re-exec). The `env -u` list is
  required because this rig's controller session inherits `AILANG_DRIVER_PINNED` etc. from the
  pinned driver (measured this fire). (Red at base: with `AILANG_DRIVER_PIN` unset the function
  proceeds to pin and re-execs.)

- **AC4** — next-fire confirmation (the one fire-dependent AC; checked at the next launchd fire).
  `grep -c 'pin disabled' /tmp/ailang-mission-world.log` → `>= 1`, and the controller session's
  `pwd` is `/Users/voightkampff/dev/sunholo-data/ailang-world`.

**Fleet discriminator (FLEET; AC-F1–F3 runnable now, AC-F4–F11 require the fix):**

- **AC-F1** — world's work repo is a DIFFERENT repository from the driver's source clone.
  ```
  [ "$(git -C /Users/voightkampff/dev/sunholo-data/ailang-world config --get remote.origin.url)" != "$(git -C /Users/voightkampff/dev/sunholo-data/ailang config --get remote.origin.url)" ] && echo DIFFERENT-REPO || echo SAME-REPO
  ```
  → `DIFFERENT-REPO`.

- **AC-F2** — motoko's work repo is the SAME repository as the driver's source clone.
  ```
  [ "$(git -C /Users/voightkampff/dev/sunholo-data/ailang-motoko config --get remote.origin.url)" = "$(git -C /Users/voightkampff/dev/sunholo-data/ailang config --get remote.origin.url)" ] && echo SAME-REPO || echo DIFFERENT-REPO
  ```
  → `SAME-REPO`.

- **AC-F3** — `git-common-dir` (resolved to realpath) DIFFERS for motoko vs ailang, proving it is
  the wrong discriminator (a separate clone of the same repo has its own `.git`).
  ```
  resolve_common() { local d; d=$(git -C "$1" rev-parse --git-common-dir) || return 1; case "$d" in /*) echo "$d";; *) echo "$(cd "$1" && pwd -P)/$d";; esac; }
  [ "$(resolve_common /Users/voightkampff/dev/sunholo-data/ailang-motoko)" = "$(resolve_common /Users/voightkampff/dev/sunholo-data/ailang)" ] && echo SAME-REPO || echo DIFFERENT-REPO
  ```
  → `DIFFERENT-REPO`. This is the known-positive control that justifies rejecting `git-common-dir`.

- **AC-F4** — the fixed helper classifies world as a different repo (requires the fix present).
  ```
  . /Users/voightkampff/dev/sunholo-data/ailang/tools/launchd/lib/pin-root.sh; _pin_same_repo /Users/voightkampff/dev/sunholo-data/ailang-world /Users/voightkampff/dev/sunholo-data/ailang; echo "rc=$?"
  ```
  → `rc=1` (different → no redirect).

- **AC-F5** — the fixed helper classifies motoko as the same repo (requires the fix present).
  ```
  . /Users/voightkampff/dev/sunholo-data/ailang/tools/launchd/lib/pin-root.sh; _pin_same_repo /Users/voightkampff/dev/sunholo-data/ailang-motoko /Users/voightkampff/dev/sunholo-data/ailang; echo "rc=$?"
  ```
  → `rc=0` (same → redirect). **Reds if the fleet uses `git-common-dir`:** motoko would return `rc=1`.

- **AC-F6** — the fixed `_set_pin_workdir` keeps world's `MISSION_WORKDIR` (requires the fix).
  ```
  . /Users/voightkampff/dev/sunholo-data/ailang/tools/launchd/lib/pin-root.sh
  MISSION_WORKDIR=/Users/voightkampff/dev/sunholo-data/ailang-world; _set_pin_workdir /tmp/pinwt /Users/voightkampff/dev/sunholo-data/ailang; echo "rc=$? MISSION_WORKDIR=$MISSION_WORKDIR"
  ```
  → `rc=0 MISSION_WORKDIR=/Users/voightkampff/dev/sunholo-data/ailang-world` (NOT `/tmp/pinwt`).

- **AC-F7** — the fixed `_set_pin_workdir` redirects motoko's `MISSION_WORKDIR` (requires the fix).
  ```
  . /Users/voightkampff/dev/sunholo-data/ailang/tools/launchd/lib/pin-root.sh
  MISSION_WORKDIR=/Users/voightkampff/dev/sunholo-data/ailang-motoko; _set_pin_workdir /tmp/pinwt /Users/voightkampff/dev/sunholo-data/ailang; echo "rc=$? MISSION_WORKDIR=$MISSION_WORKDIR"
  ```
  → `rc=0 MISSION_WORKDIR=/tmp/pinwt` (redirect preserved for same-repo missions).

- **AC-F8** — the flag set to `1` forces a de-forked mission to keep its own repo even if the inference would say otherwise.
  ```
  . /Users/voightkampff/dev/sunholo-data/ailang/tools/launchd/lib/pin-root.sh
  AILANG_DRIVER_MISSION_IS_DE_FORKED=1; MISSION_WORKDIR=/Users/voightkampff/dev/sunholo-data/ailang-world; _set_pin_workdir /tmp/pinwt /Users/voightkampff/dev/sunholo-data/ailang; echo "rc=$? MISSION_WORKDIR=$MISSION_WORKDIR"
  ```
  → `rc=0 MISSION_WORKDIR=/Users/voightkampff/dev/sunholo-data/ailang-world` (flag wins over any inference).

- **AC-F9** — the flag unset falls back to inference and still keeps world's repo.
  ```
  . /Users/voightkampff/dev/sunholo-data/ailang/tools/launchd/lib/pin-root.sh
  unset AILANG_DRIVER_MISSION_IS_DE_FORKED; MISSION_WORKDIR=/Users/voightkampff/dev/sunholo-data/ailang-world; _set_pin_workdir /tmp/pinwt /Users/voightkampff/dev/sunholo-data/ailang; echo "rc=$? MISSION_WORKDIR=$MISSION_WORKDIR"
  ```
  → `rc=0 MISSION_WORKDIR=/Users/voightkampff/dev/sunholo-data/ailang-world` (inference: different repo).

- **AC-F10** — an empty/absent-origin work repo is LOUD: `_set_pin_workdir` returns 1, sets `PIN_STATUS=STALE`, and emits a WARN on stderr.
  ```
  . /Users/voightkampff/dev/sunholo-data/ailang/tools/launchd/lib/pin-root.sh
  d=$(mktemp -d); git -C "$d" init -q; git -C "$d" config --unset-all remote.origin 2>/dev/null
  MISSION_WORKDIR="$d"; _set_pin_workdir /tmp/pinwt /Users/voightkampff/dev/sunholo-data/ailang 2>/tmp/f10err.txt; echo "rc=$? PIN_STATUS=$PIN_STATUS MISSION_WORKDIR=$MISSION_WORKDIR"; cat /tmp/f10err.txt
  ```
  → `rc=1 PIN_STATUS=STALE` with `MISSION_WORKDIR` UNCHANGED (still `$d`, never blanked) and a WARN
  line naming the indeterminate reason. Reading `PIN_STATUS` in the caller is the load-bearing half:
  it is only observable because the call is not a command substitution.

- **AC-F11** — an SSH clone and an HTTPS clone of the SAME repo are the same repo after normalisation.
  ```
  . /Users/voightkampff/dev/sunholo-data/ailang/tools/launchd/lib/pin-root.sh
  a=$(mktemp -d); b=$(mktemp -d); git -C "$a" init -q; git -C "$b" init -q; git -C "$a" remote add origin git@github.com:sunholo-data/ailang.git; git -C "$b" remote add origin https://github.com/sunholo-data/ailang.git
  _pin_same_repo "$a" "$b"; echo "rc=$?"
  ```
  → `rc=0` (same repo despite SSH-vs-HTTPS URL formats).

- **AC-F12** — the STALE refusal REACHES THE CALLER (the round-2 quorum's central objection).
  The production call site must be a plain call, not a command substitution, so that
  `_pin_stale`'s `PIN_STATUS=STALE` is observable in `pin_root_to_committed_ref`.
  ```
  . /Users/voightkampff/dev/sunholo-data/ailang/tools/launchd/lib/pin-root.sh
  grep -n 'MISSION_WORKDIR="\$(_set_pin_workdir' /Users/voightkampff/dev/sunholo-data/ailang/tools/launchd/lib/pin-root.sh; echo "substitution-hits-rc=$?"
  grep -nc '^  _set_pin_workdir "\$wt" "\$src" || return 1$' /Users/voightkampff/dev/sunholo-data/ailang/tools/launchd/lib/pin-root.sh
  ```
  → `substitution-hits-rc=1` (grep found NO command-substitution form) and a count of `1` for the
  plain-call form. Pair it with the runtime reading from AC-F10: `PIN_STATUS=STALE` visible in the
  caller's shell is the property; the greps are what stop the form regressing.

---

## Non-Vacuity — Named RED Mutation for Every Gate

- **AC1** — mutation: comment out / change `AILANG_DRIVER_PIN=0` → `grep '^AILANG_DRIVER_PIN=0'`
  finds nothing → reds.
- **AC2** — mutation: remove the `AILANG_DRIVER_PIN=0` line → sourcing yields `PIN=<unset>` → reds.
- **AC3** — mutation: set `AILANG_DRIVER_PIN=1` (or unset) → function pins and re-execs, the `echo`
  is never reached → reds.
- **AC4** — mutation: remove the `AILANG_DRIVER_PIN=0` line → next fire logs `running committed ...`
  and the controller's `pwd` is the ailang worktree → reds.
- **AC-F1** — mutation: `git -C .../ailang-world remote set-url origin https://github.com/sunholo-data/ailang.git`
  → URLs match → `SAME-REPO` → reds.
- **AC-F2** — mutation: `git -C .../ailang-motoko remote set-url origin https://github.com/sunholo-data/ailang-world.git`
  → URLs differ → `DIFFERENT-REPO` → reds.
- **AC-F3** — mutation: make `ailang-motoko` a WORKTREE of `ailang` → common dirs match →
  `SAME-REPO` → reds. (Proves AC-F3 measures clone-vs-worktree, not a tautology.)
- **AC-F4** — mutation: `_pin_same_repo` always returns 0 → world classified same → `rc=0` → reds.
- **AC-F5** — mutation: `_pin_same_repo` uses `git-common-dir` → motoko `rc=1` → reds. **Catches
  the wrong-discriminator bug; without it a `git-common-dir` impl passes every other AC.**
- **AC-F6** — mutation: keep unconditional `MISSION_WORKDIR="$wt"` (no `_set_pin_workdir`) → the
  command fails; or always echo `$1` → world gets `/tmp/pinwt` → reds.
- **AC-F7** — mutation: `_set_pin_workdir` always assigns the incoming `$MISSION_WORKDIR` → motoko gets its own
  repo → reds.
- **AC-F8** — mutation: `_set_pin_workdir` ignores the flag (always infer) → with flag `1` and a
  mis-read URL, world gets `/tmp/pinwt` → reds.
- **AC-F9** — mutation: `_set_pin_workdir` requires the flag (no inference) → flag unset, world
  gets `/tmp/pinwt` or errors → reds.
- **AC-F10** — mutation: default the indeterminate case to "same repo" instead of `_pin_stale` →
  empty-origin repo gets `/tmp/pinwt`, `rc=0`, no WARN, `PIN_STATUS` not `STALE` → reds. **Catches
  the silent-fallback bug the quorum rejected; without it a quietly-defaulting impl passes every
  other AC.**
- **AC-F11** — mutation: compare raw URLs instead of the canonical form → SSH and HTTPS differ as
  strings → `rc=1` → reds. **Catches the URL-fragility bug; without it a raw-string comparison
  passes every other AC.**
- **AC-F12** — mutation: restore the command-substitution call site,
  `MISSION_WORKDIR="$(_set_pin_workdir "$wt" "$src")" || return 1`. `_pin_stale` then runs in a
  SUBSHELL, `PIN_STATUS` stays unchanged in the caller, and `MISSION_WORKDIR` is assigned the
  captured stdout (empty on the failure path) BEFORE `|| return 1` executes → the grep for the
  substitution form hits and AC-F10's `PIN_STATUS=STALE` reading reds. **This is the mutation
  that pins the round-2 objection itself: every other AC passes under it, which is exactly why
  two independent reviewers had to catch it by reading the code rather than by running it.**

---

## Risks and What Could Still Break

- **#558 re-opened for world (the mitigation's cost).** With `AILANG_DRIVER_PIN=0`, the driver
  runs the working tree of `ailang` as-is. If the fleet ever leaves uncommitted or stale changes
  in `ailang`, world runs them. Mitigated by P1 (fleet syncs by commit) and by the measured clean
  state this fire, but it is a real residual until M2 lands. This is the strongest reason to push
  the fleet fix promptly.
- **The fleet fix could regress v1/docs/motoko.** If the fleet implements the discriminator with
  `git-common-dir` instead of `origin`, motoko/docs lose their redirect and run against a stale
  charter — the exact half-fix the pin's comment warns about. AC-F5 is the guard; it must be part
  of the fleet's own verification before landing.
- **A mission whose work repo has no `origin`, flag unset** is now LOUD, not silent: the inference
  returns indeterminate and `_set_pin_workdir` calls `_pin_stale`, so the driver posts
  `driver ran UNPINNED (code provenance unknown)` on the human channel and the fire runs the
  working tree. No current mission hits this; a mission that sets the flag avoids it entirely.
  This is the observable runtime diagnostic the round-1 objection demanded — a commit-message
  note is no longer the only witness.
- **The stale pin worktree `~/.ailang-driver-pin/world`.** Under the mitigation it is unused; under
  the durable fix it is the driver's pinned home (a worktree of `ailang`). Its name is now
  slightly misleading (it is the driver's pinned worktree, not the mission's work repo). No action
  needed; do not delete it while the pin is enabled.
- **The env file is sourced twice** (lines 71 and 74-75). `AILANG_DRIVER_PIN=0` is a bare
  assignment, so double-sourcing is idempotent. No risk.
- **A future de-fork of another mission** (docs/motoko) would hit the same bug. The durable fix
  generalizes to any de-forked mission; the mitigation is world-specific and would need a
  per-mission `AILANG_DRIVER_PIN=0` in that mission's env file until the fleet fix lands.

---

## Deferred Scope

- **Cleanup / rename of `~/.ailang-driver-pin/world`** — cosmetic; defer until the durable fix is
  confirmed, then decide whether to keep it as the driver's pinned home.
- **A general "de-forked mission" flag in the driver** — RESOLVED by Decision 2, not deferred.
  The explicit flag `AILANG_DRIVER_MISSION_IS_DE_FORKED` is now part of the durable fix
  (authoritative when set, normalised-origin inference as the default), so this bullet is no
  longer open.
- **A World-side regression test** that asserts `REPO` resolves to `ailang-world` under the
  mitigation. The env file is rig-local and not in the repo, so a repo test cannot reach it; a
  `scripts/` check that greps the env file is possible but low-value. Defer.

---

## Quorum Round 1 — objections and how this revision answers them

**gemini-3-1-pro (REJECT) — URL comparison fragility.** The objection: exact string equality of
`remote.origin.url` is non-deterministic across SSH-vs-HTTPS and trailing `.git`, so the match
fails and silently drops the redirect for same-repo missions. Answered on the merits: the
inference now compares a CANONICAL form (`_pin_origin_canon`: strip trailing `.git` and `/`,
rewrite both URL schemes to `host/path`, lowercase host but not path), and the explicit flag
`AILANG_DRIVER_MISSION_IS_DE_FORKED` is adopted as authoritative-when-set, with the normalised
inference as the default so no existing mission must be reconfigured. AC-F11 (SSH-vs-HTTPS same
repo) and AC-F8/F9 (flag set/unset) pin both mechanisms down.

**oc-glm-5-2 (REJECT) — the indeterminate case is unobservable.** The objection: when `origin` is
unset or empty on either side, the old function silently returned non-zero with no warning, log,
`PIN_STATUS` change, or AC — a false-success bug of the same class the doc was fixing. Answered
on the merits: the indeterminate case is now LOUD by construction. `_pin_same_repo` returns 2
(indeterminate) and `_set_pin_workdir` calls `_pin_stale "<reason>"`: `PIN_STATUS=STALE`, rc=1,
and the DRIVER posts `driver ran UNPINNED (code provenance unknown)` on the human channel — the
exact primitive `pin-root.sh` already uses for every other refusal, and the exact channel this
mission saw fire this very fire. We chose `_pin_stale` over a stderr WARN alone because the harms
are asymmetric (a wrong redirect strands the charter silently; a refused pin runs uncommitted
driver code loudly) and the STALE path is observable by construction. AC-F10 (below) pins this
down with a named RED mutation.

**gpt6-astra (ABSENT — budget refusal, not a pass).** No objection to answer; the quorum degraded
to N-1. This revision stays at or below ~500 lines so astra's raised cap can review it in round 2.

---

## Quorum Round 2 — objections, and the narrow-refinement carve-out applied

Round 2 was **BLOCKED**, and the block was correct. `gpt6-astra` and `gemini-3-1-pro` independently
rejected on the SAME mechanism: the round-2 draft invoked the work-dir helper inside a command
substitution — `MISSION_WORKDIR="$(_pin_workdir_for "$wt" "$src")" || return 1` — which runs in a
**subshell**, so `_pin_stale`'s `PIN_STATUS=STALE` was discarded on return and the driver would have
seen an unchanged `PIN_STATUS`. That defeats the loud-failure guarantee round 2 was written to
provide. astra added a second hazard in the same line: the substitution assigns captured stdout
(empty on the failure path) BEFORE `|| return 1` runs, so a refused call would first blank the work
directory. `oc-glm-5-2` was **ABSENT** — `ollama error: 429 … session usage limit` — so the round-2
quorum degraded to N−1; its round-1 objection was answered by the round-2 revision, but that answer
was **not re-reviewed by glm**, and this doc says so rather than counting it as a pass.

Disposition: the skill's **narrow-refinement carve-out** (Gate 2, added V1 iter-95) applies and was
taken. Both remaining blocking objections (a) carry a concrete reviewer-authored `proposed_fix` and
(b) dispute no part of the design DIRECTION — both reviewers accept the de-fork guard, the
canonicalised-origin discriminator, the authoritative-when-set flag and the loud-indeterminate
requirement, and object only to the shell mechanism that implements the last one. The controller
therefore applied the reviewers' verbatim fix and routed to the sprint-planner rather than parking.

What was applied, and why one fix answers both: `gemini-3-1-pro`'s prescription — *"Rewrite the
helper to mutate the variable in the current shell rather than echoing it from a subshell. Rename
`_pin_workdir_for` to `_set_pin_workdir`, replace its `echo` statements with direct assignments, and
replace the command substitution at `pin-root.sh:267` with a direct call: `_set_pin_workdir "$wt"
"$src" || return 1`"* — was applied verbatim. It **subsumes** `gpt6-astra`'s objection by
construction rather than by controller judgement: astra's second hazard is a property of assigning a
command substitution's stdout, and after gemini's fix there is no substitution and no capture, so
nothing can be assigned on the failure path. astra's own alternative (a temp variable plus an
explicit `if`) reaches the same two guarantees by a different route; the controller did not invent a
third. **AC-F12** and its named RED mutation were added so the form cannot regress silently, since —
as the reviewers demonstrated — every other AC passes under the defective form.

---

## Open Decisions

None requiring a human ruling. The mitigation is World-landable under P2; the durable fix is
fleet work under P1 and is specified precisely enough to implement without re-derivation. The
only "wait" is the fleet's cadence for landing M2 — a scheduling fact, not a decision.
