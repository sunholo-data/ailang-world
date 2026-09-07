# pin-root: preserve a de-forked mission work repository across driver pin re-exec

World hand-off: `w-defork-pin-redirects-work-repo`.

## Problem and first-party evidence

World was de-forked by commit `e92594c`: its launchd job now invokes the FLEET driver at
`/Users/voightkampff/dev/sunholo-data/ailang/tools/launchd/mission-control.sh`, while the mission's
charter and work remain in the separate `sunholo-data/ailang-world` repository. The driver pin
re-executes committed driver code from `/Users/voightkampff/.ailang-driver-pin/world` and also
unconditionally redirects `MISSION_WORKDIR` there. That path is a worktree of `sunholo-data/ailang`,
so the World controller starts in the wrong repository and cannot resolve
`design_docs/world-mission.md`.

Observed in the live controller session on 2026-09-07:

```text
MISSION_WORKDIR = /Users/voightkampff/.ailang-driver-pin/world
pwd             = /Users/voightkampff/.ailang-driver-pin/world
origin          = https://github.com/sunholo-data/ailang.git
ls design_docs/world-mission.md -> No such file or directory
PIN_STATUS      = pinned
drift           = 0
```

The pin therefore logged full success while moving the controller into the wrong repository:

```text
[2026-09-07 08:44:42] driver pin: running committed origin/dev @ 878939117 from /Users/voightkampff/.ailang-driver-pin/world (source clone /Users/voightkampff/dev/sunholo-data/ailang was 0 behind)
```

The real origin measurements are:

```text
git -C ~/dev/sunholo-data/ailang config --get remote.origin.url
  -> https://github.com/sunholo-data/ailang.git
git -C ~/dev/sunholo-data/ailang-motoko config --get remote.origin.url
  -> https://github.com/sunholo-data/ailang.git
git -C ~/dev/sunholo-data/ailang-world config --get remote.origin.url
  -> https://github.com/sunholo-data/ailang-world.git
```

`git-common-dir` is the wrong discriminator. `ailang-motoko` is a separate clone of the same
repository as `ailang`, yet the three absolute common directories are all different:

```text
git -C <each> rev-parse --path-format=absolute --git-common-dir
  ailang        -> /Users/voightkampff/dev/sunholo-data/ailang/.git
  ailang-motoko -> /Users/voightkampff/dev/sunholo-data/ailang-motoko/.git
  ailang-world  -> /Users/voightkampff/dev/sunholo-data/ailang-world/.git
```

A common-dir comparison would misclassify the same-repository `ailang-motoko` clone as de-forked
and break the redirect that v1/docs/motoko require.

## Requested M2 semantics

Please change `tools/launchd/lib/pin-root.sh` so the committed driver remains pinned, but its pin
worktree replaces `MISSION_WORKDIR` only for a same-repository mission:

- `AILANG_DRIVER_MISSION_IS_DE_FORKED` is authoritative when set.
- With the flag unset, infer same versus different repository from normalized `origin` URLs. Strip
  trailing `.git` and `/`, normalize SSH and scheme URLs to `host/path`, and lowercase the host but
  not the path.
- If either origin is absent, empty, unreadable, or otherwise indeterminate, fail LOUD through
  `_pin_stale`: return 1, set `PIN_STATUS=STALE`, retain the incoming `MISSION_WORKDIR`, and emit the
  reason. Do not silently choose either redirect behavior.
- Preserve the incoming mission work repository for a different-repository mission; redirect to
  the pin worktree for a same-repository mission.

The round-2 quorum binding is load-bearing: `_set_pin_workdir` must mutate `MISSION_WORKDIR` in the
current shell, and the production call site must be exactly the plain call

```bash
_set_pin_workdir "$wt" "$src" || return 1
```

Never use a command substitution. A form such as
`MISSION_WORKDIR="$(_set_pin_workdir "$wt" "$src")" || return 1` runs `_pin_stale` in a subshell,
so `PIN_STATUS=STALE` never reaches the driver; it also assigns captured stdout before the failure
is propagated and can blank `MISSION_WORKDIR`.

## Acceptance commands and expected results

Run these against the fleet implementation. AC-F1 through AC-F3 establish the real topology;
AC-F4 through AC-F12 require the fixed helper.

### AC-F1 — World is a different repository

```bash
[ "$(git -C /Users/voightkampff/dev/sunholo-data/ailang-world config --get remote.origin.url)" != "$(git -C /Users/voightkampff/dev/sunholo-data/ailang config --get remote.origin.url)" ] && echo DIFFERENT-REPO || echo SAME-REPO
```

Expected: `DIFFERENT-REPO`.

### AC-F2 — Motoko is the same repository

```bash
[ "$(git -C /Users/voightkampff/dev/sunholo-data/ailang-motoko config --get remote.origin.url)" = "$(git -C /Users/voightkampff/dev/sunholo-data/ailang config --get remote.origin.url)" ] && echo SAME-REPO || echo DIFFERENT-REPO
```

Expected: `SAME-REPO`.

### AC-F3 — the common-dir counterexample

```bash
resolve_common() { local d; d=$(git -C "$1" rev-parse --git-common-dir) || return 1; case "$d" in /*) echo "$d";; *) echo "$(cd "$1" && pwd -P)/$d";; esac; }
[ "$(resolve_common /Users/voightkampff/dev/sunholo-data/ailang-motoko)" = "$(resolve_common /Users/voightkampff/dev/sunholo-data/ailang)" ] && echo SAME-REPO || echo DIFFERENT-REPO
```

Expected: `DIFFERENT-REPO`. This proves common-dir identity cannot represent repository identity
across separate clones.

### AC-F4 — classify World as different

```bash
. /Users/voightkampff/dev/sunholo-data/ailang/tools/launchd/lib/pin-root.sh; _pin_same_repo /Users/voightkampff/dev/sunholo-data/ailang-world /Users/voightkampff/dev/sunholo-data/ailang; echo "rc=$?"
```

Expected: `rc=1` (different, so do not redirect).

### AC-F5 — wrong-discriminator regression killer

```bash
. /Users/voightkampff/dev/sunholo-data/ailang/tools/launchd/lib/pin-root.sh; _pin_same_repo /Users/voightkampff/dev/sunholo-data/ailang-motoko /Users/voightkampff/dev/sunholo-data/ailang; echo "rc=$?"
```

Expected: `rc=0` (same, so preserve the redirect). This is the wrong-discriminator regression
killer: a `git-common-dir` implementation returns `rc=1` for the separate Motoko clone.

### AC-F6 — keep World's work repository

```bash
. /Users/voightkampff/dev/sunholo-data/ailang/tools/launchd/lib/pin-root.sh
MISSION_WORKDIR=/Users/voightkampff/dev/sunholo-data/ailang-world; _set_pin_workdir /tmp/pinwt /Users/voightkampff/dev/sunholo-data/ailang; echo "rc=$? MISSION_WORKDIR=$MISSION_WORKDIR"
```

Expected: `rc=0 MISSION_WORKDIR=/Users/voightkampff/dev/sunholo-data/ailang-world` (not
`/tmp/pinwt`).

### AC-F7 — redirect Motoko to the pin worktree

```bash
. /Users/voightkampff/dev/sunholo-data/ailang/tools/launchd/lib/pin-root.sh
MISSION_WORKDIR=/Users/voightkampff/dev/sunholo-data/ailang-motoko; _set_pin_workdir /tmp/pinwt /Users/voightkampff/dev/sunholo-data/ailang; echo "rc=$? MISSION_WORKDIR=$MISSION_WORKDIR"
```

Expected: `rc=0 MISSION_WORKDIR=/tmp/pinwt`.

### AC-F8 — authoritative de-fork flag

```bash
. /Users/voightkampff/dev/sunholo-data/ailang/tools/launchd/lib/pin-root.sh
AILANG_DRIVER_MISSION_IS_DE_FORKED=1; MISSION_WORKDIR=/Users/voightkampff/dev/sunholo-data/ailang-world; _set_pin_workdir /tmp/pinwt /Users/voightkampff/dev/sunholo-data/ailang; echo "rc=$? MISSION_WORKDIR=$MISSION_WORKDIR"
```

Expected: `rc=0 MISSION_WORKDIR=/Users/voightkampff/dev/sunholo-data/ailang-world`; the flag wins
over inference.

### AC-F9 — unset flag uses origin inference

```bash
. /Users/voightkampff/dev/sunholo-data/ailang/tools/launchd/lib/pin-root.sh
unset AILANG_DRIVER_MISSION_IS_DE_FORKED; MISSION_WORKDIR=/Users/voightkampff/dev/sunholo-data/ailang-world; _set_pin_workdir /tmp/pinwt /Users/voightkampff/dev/sunholo-data/ailang; echo "rc=$? MISSION_WORKDIR=$MISSION_WORKDIR"
```

Expected: `rc=0 MISSION_WORKDIR=/Users/voightkampff/dev/sunholo-data/ailang-world` from the
different-origin inference.

### AC-F10 — loud-indeterminate regression killer

```bash
. /Users/voightkampff/dev/sunholo-data/ailang/tools/launchd/lib/pin-root.sh
d=$(mktemp -d); git -C "$d" init -q; git -C "$d" config --unset-all remote.origin 2>/dev/null
MISSION_WORKDIR="$d"; _set_pin_workdir /tmp/pinwt /Users/voightkampff/dev/sunholo-data/ailang 2>/tmp/f10err.txt; echo "rc=$? PIN_STATUS=$PIN_STATUS MISSION_WORKDIR=$MISSION_WORKDIR"; cat /tmp/f10err.txt
```

Expected: `rc=1 PIN_STATUS=STALE`, `MISSION_WORKDIR` unchanged at `$d` and never blank, plus a WARN
line naming the indeterminate reason. This is the loud-indeterminate regression killer: silently
defaulting to same-repo would instead redirect to `/tmp/pinwt` with rc 0 and no STALE status.

### AC-F11 — SSH and HTTPS normalize to the same repository

```bash
. /Users/voightkampff/dev/sunholo-data/ailang/tools/launchd/lib/pin-root.sh
a=$(mktemp -d); b=$(mktemp -d); git -C "$a" init -q; git -C "$b" init -q; git -C "$a" remote add origin git@github.com:sunholo-data/ailang.git; git -C "$b" remote add origin https://github.com/sunholo-data/ailang.git
_pin_same_repo "$a" "$b"; echo "rc=$?"
```

Expected: `rc=0` despite the SSH-versus-HTTPS forms.

### AC-F12 — current-shell regression killer

```bash
. /Users/voightkampff/dev/sunholo-data/ailang/tools/launchd/lib/pin-root.sh
grep -n 'MISSION_WORKDIR="\$(_set_pin_workdir' /Users/voightkampff/dev/sunholo-data/ailang/tools/launchd/lib/pin-root.sh; echo "substitution-hits-rc=$?"
grep -nc '^  _set_pin_workdir "\$wt" "\$src" || return 1$' /Users/voightkampff/dev/sunholo-data/ailang/tools/launchd/lib/pin-root.sh
```

Expected: `substitution-hits-rc=1` and a count of `1` for the plain-call form. Pair this with
AC-F10's runtime observation that `PIN_STATUS=STALE` is visible in the caller. AC-F12 is the
current-shell regression killer: command substitution makes the first grep hit and loses the
caller's STALE status.

## Ownership and hand-back

The fleet owns, implements, tests, and lands this M2 change in `sunholo-data/ailang`. World supplies
the first-party evidence and this acceptance contract; under ratified `D-WORLD-DRIVER-1`, World
does not edit, temporarily mutate, vendor, or sync frozen core under `tools/launchd/*`. Until M2 is
landed and handed back, the controller owns the rig-local M1 mitigation `AILANG_DRIVER_PIN=0`.
World's later M3 re-enables the pin only after the fleet fix passes these checks and a live fire
shows a pinned driver with the controller still in `sunholo-data/ailang-world`.

External clone, launchd, live-controller, and rig-env verdicts made by the World executor are
**UNINFORMATIVE UNDER SANDBOX**; the controller must reproduce those facts before filing and owns
the live proof.

---

## Addendum — a gap the independent judge found in AC-F8, before you implement

World's iteration-166 evaluator (an independent Anthropic lane, distinct from the OpenAI executor
that wrote this body) implemented BOTH the specified `_set_pin_workdir` and AC-F8's own named RED
mutant ("ignores the flag entirely, always infers") and ran AC-F8's literal command against each:

```
AILANG_DRIVER_MISSION_IS_DE_FORKED=1; MISSION_WORKDIR=<ailang-world>; _set_pin_workdir /tmp/pinwt <ailang>
  correct impl -> rc=0 MISSION_WORKDIR=/Users/voightkampff/dev/sunholo-data/ailang-world
  mutant  impl -> rc=0 MISSION_WORKDIR=/Users/voightkampff/dev/sunholo-data/ailang-world   (IDENTICAL)
```

**AC-F8 is therefore vacuous as written.** Because world's real origin genuinely differs from
ailang's, ignoring the flag and falling back to inference produces the same answer as honouring it,
so no implementation that drops the flag can be caught by this criterion. The doc's prose smuggles
in an unstated extra condition ("with flag `1` and a **mis-read URL**") that the command never
exercises.

Fix it before implementing, not after: exercise the flag against a case where the flag and the
inference **disagree** — e.g. set `AILANG_DRIVER_MISSION_IS_DE_FORKED=1` with `MISSION_WORKDIR`
pointed at a clone whose origin MATCHES `$src` (a temp clone of `ailang`, or `ailang-motoko`), and
assert the work repo is KEPT rather than redirected. Under a flag-ignoring implementation that case
redirects to the pin worktree and the criterion reds.

AC-F5, AC-F10 and AC-F12 were checked by the same judge and are sound; AC-F8 is the one that needs
rewriting.
