# Sprint plan: w-check-no-personal-email

**Iteration:** World 175
**Design:** `design_docs/planned/w-check-no-personal-email.md` (252 lines; authored by
`pi:ollama/deepseek-v4-flash:0731-cloud`, revised once by the designer after quorum r1, revised once
more by the controller under the narrow-refinement carve-out after r2 — both recorded in its header)
**Owning queue row:** 74, `w-the-personal-email-gate-the-shared-skill-cites-as-enforcement-does-not-exist`
(clause-2, ~0.3d, gated on nothing, surfaced iter-152; charter line 5698)
**Planning base:** `a295291bea069d793960e08f055e501c6f9661f4` (== `origin/dev`; the sprint branch
`sprint/w-check-no-personal-email` sits at `f53b353`, three design-doc commits on top, and differs
from base ONLY in `design_docs/planned/w-check-no-personal-email.md` — out of the gate's scope)
**Executor worktree:** `/Users/voightkampff/dev/sunholo-data/.wt-world-iter175` (measured: exists,
at `f53b353`, `git status --porcelain` empty)
**Executor lane (M1):** `pi:ollama/deepseek-v4-flash:0731-cloud` via V1's `mission_pi_run.sh` — a
flat-rate, UNSANDBOXED pi lane that cannot read this repo's skills. **Everything M1 needs is written
into this plan verbatim; nothing is implicit.** It performs NO git write and snapshots its work
into `.snap/M1/`; the controller reconstructs the commit from the snapshot by explicit file list.
**Controller-run milestone:** M2 (charter edits only), after M1 lands on `dev`.
**Fleet precedent:** `sunholo-data/ailang` `f0d44915d` — `scripts/check_no_personal_email.sh` +
`scripts/test_check_no_personal_email.sh` (read from `origin/dev` of the fleet checkout; the
instrument text below is that precedent adapted per D1/D2/D5/D6 and then RUN, not transcribed).
**Status:** planned; nothing implemented, no git write performed by the planner (one permitted
`git worktree add --detach` of a sibling base tree, removed again).

**Address hygiene in this file:** no real personal address appears anywhere below. The current
default at `scripts/mission_answer.sh:69` is referred to as `<the line-69 personal default>`; fixtures
use `someone@realdomain.com`-style placeholders that the plan's own text never assembles into a
PAT-matching token inside a scanned file (see P1/P2 — that trap fired twice while drafting).

---

## Outcome and scope

One bash-3.2 instrument, an 8-arm self-test that drives it as a process, a two-line change to
`scripts/mission_answer.sh`, one new arm in its suite, two CI steps, and two charter edits:

```text
scripts/check_no_personal_email.sh        # D1/D2/D6: git ls-files × World SCOPE_RE × PAT × anchored
                                          # exclusion list; ✓/✗/floor lines; exit 0 / 1 / 2
scripts/test_check_no_personal_email.sh   # D5: arms A–H, one probe string per arm, fixtures split at "@"
scripts/mission_answer.sh                 # D3: line 69 default → GitHub noreply; +1 dry-run identity note
scripts/test_mission_answer.sh            # D3: ARM 6 (6a, 6b) — 16 → 18 assertions
.github/workflows/ci.yml                  # D4: two steps appended to go-verify (16 steps after)
design_docs/world-mission.md              # M2: row 74 LANDED tag, one Repo Profile bullet, new row 91
```

Not in scope: `tools/launchd/*` (frozen), any `.ail`, the shared skill, a Makefile, an allow-list
mechanism (D3 declares it NOT built), widening the scan to `design_docs/**` (M2 files it as row 91).

---

## Baseline table — every acceptance command re-run by the planner at `a295291`

Measured in a detached sibling worktree `/Users/voightkampff/dev/sunholo-data/.base-world-iter175-planner`
at `a295291` (`git status --porcelain` empty), `PATH=/opt/homebrew/bin:$PATH`,
`AILANG_BIN=$HOME/.pinned-ailang/ailang` (v0.30.0), every script run as `bash <script>`, every grep
`/usr/bin/grep`. Worktree removed afterwards. A gate already red at base measures the repo, not the
change (rule 3e) — none of the profile gates is.

| AC | Command (exactly as run) | Observed at base | Doc's post-change expectation | Verdict |
|---|---|---|---|---|
| AC-M1-1 | `test -x scripts/check_no_personal_email.sh` | file ABSENT (rc=1) | true | consistent (created in M1; 17 of 28 tracked `scripts/` files are `100755`, the two suites among them) |
| AC-M1-2 | `bash scripts/check_no_personal_email.sh; echo rc=$?` | ABSENT (`No such file`, rc=127) | `✓ … loop-written surface`, rc=0 | consistent — the planner's DRAFT instrument (M1 text below) run at base reads rc=1 naming `scripts/mission_answer.sh`, `1 occurrence(s) in 34 in-scope files` |
| AC-M1-3 | the doc's sibling-worktree procedure, with the fleet's `origin/dev:scripts/check_no_personal_email.sh` as the stand-in | rc=1; ✗ line names `scripts/mission_answer.sh -> <redacted>`; 6 output lines; `git ls-files \| grep -c check_no_personal_email` → 0 (the copy is untracked, not scanned); porcelain shows only `?? scripts/check_no_personal_email.sh` | `rc=1`, ✗ names `scripts/mission_answer.sh` | consistent (P5); the fleet's `.claude/skills/.*` alternative matched nothing |
| AC-M1-4 | `bash scripts/test_check_no_personal_email.sh` | ABSENT | `8 passed, 0 failed` | consistent — the planner's draft suite reads exactly `8 passed, 0 failed` in a scratch repo (P4) |
| AC-M1-5 | `bash scripts/test_mission_answer.sh` | rc=0, **`16 passed, 0 failed`**; arm ids 1a–1g, 2a–2c, 3a, 4a–4d, 5a | "all arms pass including ARM 6" | consistent; post-change line is **`18 passed, 0 failed`** (P9, measured with ARM 6 spliced in) |
| AC-M1-6 | `/usr/bin/grep -c 'check_no_personal_email' .github/workflows/ci.yml` | **0**; positive control `/usr/bin/grep -c 'test_gate0_self_notices' .github/workflows/ci.yml` → **1** | 2 | consistent (the two `run:` lines; the `name:` lines use hyphens, not underscores) |
| AC-M1-7 | `git diff --stat a295291 -- tools/launchd/` | empty (0 lines) | empty | consistent |
| AC-M1-8a | `go vet ./...` | rc=0, 0 lines | rc 0 | GREEN at base |
| AC-M1-8b | `go test ./... -count=1 -timeout 600s` with `AILANG_BIN` | rc=0, **19 packages `ok`**, 0 `no test files` | rc 0 | GREEN at base |
| AC-M1-8c | `bash ./scripts/verify_ail.sh` | **rc=1 WITHOUT `AILANG_BIN`** (`✗ AILANG_BIN refused [DEV_BUILD]: … v0.38.5-13-ga11ffb7b5-dirty`); **rc=0 WITH it**: `✓ verify gate PASSED: 11 required identities verified, 40 named tests pass`, `9/9 steps performed non-zero work` | rc 0, 11 / 40 | GREEN at base with the pin exported (P10) |
| AC-M1-9 | `git ls-files \| /usr/bin/grep -cE '^(design_docs/[^/]*mission[^/]*\.md\|scripts/.*)$'` | **34** (6 mission docs + 28 `scripts/`); root-level `^[^/]*mission[^/]*\.md$` → 0; binary-extension probe over the 34 → 0 | 36 on the sprint tree | consistent (34 + the 2 new tracked scripts) |
| AC-M1-10 | fleet-style in-scope scan (PAT × anchored exclusion) over the 34 files | **exactly one file hits: `scripts/mission_answer.sh: 1 <redacted>`**, at line 69; `/usr/bin/grep -c 'users\.noreply\.github\.com' scripts/mission_answer.sh` → **0** | 0 hits, count 1 | consistent (P6) |
| V15 | full-tree scan, 357 tracked files | 8 files hit: `scripts/mission_answer.sh` + 7 out-of-scope files (4 `toolchain@v0.0.1-go1.2x.y.darwin` in `design_docs/implemented/w-defork-pin-redirects-work-repo*` ×3 + `w-setup-go-pin-unguarded-sprint-plan.md`; 3 `git@github.com` in `design_docs/planned/w-daemon-timeout-test-flake*` ×2 + `sprint_w-daemon-timeout-test-flake.json`) | 8 = 1 + 7 | consistent (P11) |
| V13 | `env -u MISSION_ATTENDED_EMAIL MISSION_ATTENDED_NAME="Test Human" scripts/mission_answer.sh --id D-2 --answer x --file <fixture> --dry-run` | rc=0; last line `• --dry-run: <fixture> not modified`; the string `attended identity` occurs ONCE in the output — inside the diff body (line 103's prose), NOT as a note | identity printed after D3 | consistent (P12) |
| F2 | `ls Makefile`; `ls -d .claude/skills` | both `No such file or directory` | — | consistent |
| census | `bash scripts/queue_census.sh --doc design_docs/world-mission.md --control-closed 1 --control-open 79` | `43 closed / 90 rows`, controls 1/79 ok | — | after M2: `44 closed / 91 rows` (P17) |

Rig facts re-measured: `/bin/bash` 3.2.57; the interactive `grep` is a ugrep shell function that a
child `bash <script>` does NOT inherit (`type grep` inside bash → `/usr/bin/grep`); `cp` of a 0755
file yields `-rwxr-xr-x`; `cat > existing-file` keeps the inode and mode; BSD `cat` has no `-A`.

---

## Planner findings

### P1 — the instrument's OWN comments are in scope, and a header that NAMES a bad address is a hit *(act on this: the header below is address-free; step 2 runs the gate on itself)*

A first draft whose header said "`git@github.com` remotes" and "`a@example.com.evil.net` is a HIT"
was flagged by itself the moment it was tracked: `✗ … scripts/check_no_personal_email.sh -> a@example.com.evil.net`
and `-> git@github.com`, rc=1 (measured in a scratch repo). `scripts/.*` is a live alternative, so
the instrument and its suite police their own text. The M1 texts below say "SSH remote strings" and
"an address at example.com.evil.net" (no `@`), and the executor's first act after writing each file
is to run the gate over the tree (AC-M1-2) BEFORE running the suite.

### P2 — a fixture split at the DOT is still address-shaped; split at the `@` *(act on this: verbatim fixture lines)*

The fleet's technique (`FIX="someone@realdomain"; FIX="$FIX.com"`) works for that one token but NOT
for the noreply fixture: `NOREPLY="3155884+SomeUser@users.noreply"` matches PAT on its own
(`users` + `.noreply`), is not excluded (no `github.com` tail), and the suite was flagged by the gate
(measured: `-> 3155884+SomeUser@users.noreply`, rc=1). Every fixture in the suite below is
`LOCAL="x"; LOCAL="${LOCAL}@domain.tld"` — neither half matches PAT. Post-fix PAT probe over both
new files → 0 tokens.

### P3 — the fleet's ✗ line has a colour reset BETWEEN `file:` and the filename *(act on this: colour only on a TTY)*

Fleet output, captured: `ESC[0;31m✗ personal email in loop-written file:ESC[0m scripts/mission_answer.sh  ->  …`.
A "verbatim" assertion on `loop-written file: scripts/mission_answer.sh` cannot match through the
reset. The instrument below sets the colour variables only when `[ -t 1 ]`, so captured output
(suite, CI, `$(...)`) is the bare D6 line and the AC-M1-3 assertion is
`/usr/bin/grep -q 'loop-written file: scripts/mission_answer.sh'`. This is a second, cosmetic
deviation from fleet parity beside the `$` anchor; record it in the report and in M2's bullet.

### P4 — the §5 "sole killer" column, MEASURED against the drafts (rule 3i) *(act on this: the drill expects THESE red sets)*

Run on the scratch copies of the exact M1 texts below (`cp` backup, `sed` mutant, `shasum -a 256 -c`
+ `test -x` restore, all OK):

| § | mutation | doc says | measured red set | note |
|---|---|---|---|---|
| M1 | hit-path `exit 1` → `exit 0` | B | **{B, G, H}** — probe `rc=0 hit=1 ok=0` | every rc-1 arm reds; B is the named killer, G/H are the same mechanism |
| M2 | exclusion `grep -vE` line deleted | C or D | **{A, C, D}** (D reads `hit=4`) | A reds because the REAL tree carries allowed tokens (the noreply default, `human@example.com` ×4 and the bot noreply ×2 in `test_mission_answer.sh`) |
| M3 | `SCOPE_RE` narrowed to mission docs only | G | **{G}** — *after* arm E was given a clean mission doc; before that {E, G} because E's throwaway had ZERO mission docs and the floor fired | E now writes `# charter (clean)` so the scope stays non-empty under every mutation |
| M4 | floor block removed | F | **{F}** — probe `rc=0 hit=0 ok=1 floor=0` | sole |
| M5 | `hits=$((hits+1))` → `:` | B | **{B, G, H}** — probe `rc=0 hit=1 ok=1` | same arms as M1 but a DIFFERENT signature (the ✓ line prints after the ✗ lines); the probe string localises it |
| M6 | `$` dropped from `@example\.(com\|org\|net)$` | H | **{H}** | sole |
| M7 | line-69 default reverted to a raw address | AC-M1-2 | **{arm A, AC-M1-2 (rc=1), AC-M1-10 (count 0), arm 6b}** | the live gate kills it — and so do three other readers; all four are downstream of the same byte |
| M8 | comment reword (control) | none | **{}** — `8 passed, 0 failed` | control holds |
| M9 | dry-run identity note deleted | 6b | **{6b}** — `17 passed, 1 failed` | sole; 6a and 5a stay green |

The doc's "yes" in the sole-killer column is true in the sense that every mutation HAS a named
killer that reds; three rows (M1, M2, M5) also red arms sharing the mechanism. The executor records
the measured set, never adjusts an arm to shrink it (row 90's lesson: detection is right, locality
is what the column is for — and here the probe signature localises M1 vs M5).

### P5 — AC-M1-3's sibling-worktree procedure scans the base surface *(record only)*

Confirmed with the fleet script as stand-in and again with the draft: the copied instrument is
untracked in the base tree (`git ls-files | grep -c` → 0), the scan covers the 34 base files, the
verdict is rc=1 naming `scripts/mission_answer.sh`. The doc's `BASEWT` path `…/.base-world-iter175`
is a sibling of the repo, never `/tmp`. One caveat: the doc's snippet runs `git worktree add` /
`remove` — a git WRITE. The executor may NOT run it; AC-M1-3 is a **controller** step (below).

### P6 — AC-M1-10's `users\.noreply\.github\.com` count is 1 only because the note line uses the variables *(act on this: verbatim line)*

With line 69 rewritten and `note "--dry-run: attended identity $ATT_NAME <$ATT_EMAIL>"` inserted at
line 124, `/usr/bin/grep -c 'users\.noreply\.github\.com' scripts/mission_answer.sh` reads **1**
(measured on a spliced copy). A note line that spelled the literal would read 2 and red the AC.

### P7 — PAT's `-` makes the extracted token `-3155884+…` on a `${VAR:-default}` line *(act on this: assert file names, never exact addresses)*

`echo 'ATT_EMAIL="${MISSION_ATTENDED_EMAIL:-3155884+…@users.noreply.github.com}"' | grep -oE "$PAT"`
→ `-3155884+MarkEdmondson1234@users.noreply.github.com` (still excluded — the domain tail decides).
The base line-69 token likewise begins with `-` (measured: first char `-`, length 20). So the
AC-M1-3 ✗-line assertion greps the FILE NAME; no AC compares an address string byte-for-byte.

### P8 — `scripts/test_mission_answer.sh` is not in CI *(record only)*

`/usr/bin/grep -n mission_answer .github/workflows/ci.yml` → nothing. ARM 6 is proven on the rig by
AC-M1-5 and by the evaluator, not by CI. The doc's D4 adds exactly two steps (AC-M1-6 pins `2`), so
this plan does not add a third; M2's new row 91 text mentions it as a candidate follow-up.

### P9 — the suite's `N passed` line counts assertions, not arms *(act on this: AC-M1-5 expects `18 passed, 0 failed`)*

Base reads `16 passed, 0 failed` (16 assertion ids across arms 1–5). ARM 6 adds 6a and 6b → 18.

### P10 — `verify_ail.sh` ALSO needs `AILANG_BIN` *(act on this: export it before every gate, not only `go test`)*

Without the pin it refuses the dev build on PATH (`DEV_BUILD … -dirty`, rc=1). With it, rc=0 and the
11/40 line. Same variable, same export, both gates.

### P11 — V15's "`git@github.com` ×3" is a FILE count; raw occurrences are 5 *(record only)*

Per-file `sort -u` (the instrument's shape) gives 3 files; `uniq -c` over raw tokens gives 5. The
doc's 7 false positives = 7 FILES (3 + 4), which is the unit the instrument reports. Row 91's text
uses the file count.

### P12 — `attended identity` already appears once in the dry-run output at base, inside the diff body *(record only)*

Line 103's evidence prose ("stamps a fixed attended identity for EVERY caller …") is echoed by the
`diff -u`. 6b's needle carries the `<…users.noreply.github.com>` tail, which that prose lacks, so 6b
reds on base and under M9 (both measured) — not vacuous.

### P13 — nothing in `.gitignore` covers `.snap/` or `.cnpe_scratch.*` *(record only)*

`git status --porcelain` will show `?? .snap/` after the snapshot; the controller reconstructs by
explicit file list (never `git add -A`). The suite removes its scratch on EXIT via `trap`.

### P14 — the exec-bit hazard is real and `cp` alone is safe here *(act on this: never `mv`, always `test -x`)*

Iteration 174's executor lost the mode by `sed > /tmp/x && mv x file`. Measured: `cp` of a 0755
file → `-rwxr-xr-x`; `cat > existing` keeps the inode. Every write in M1 is `cp`/`cat >` onto the
existing path; every restore is followed by `shasum -a 256 -c` AND `test -x`.

### P15 — the suite runs arm F LAST, in a second throwaway *(record only)*

F needs a repo whose scope is EMPTY, so the instrument cannot sit in `scripts/` there; it is copied
to `tools/` (not in `SCOPE_RE`) beside a tracked `README.md`. Doc order A–H is preserved in the ids;
the execution order is A, B, C, D, E, G, H, F. AC-M1-4's line is unaffected.

### P16 — the D3 dry-run note prints for EVERY `--dry-run`, not only the default *(record only)*

Arms 5a (explicit `human@example.com`) still pass — 5a asserts only that the file is unchanged. The
note is one line more on the dry-run path; nothing parses it.

### P17 — the census moves twice under M2 *(act on this in M2)*

Row 74's LANDED tag → `44 closed / 90 rows`; the new row 91 → `44 closed / 91 rows`. The CI census
step's controls (`--control-closed 1 --control-open 79`) are row NUMBERS, unaffected.

---

## Cross-checks

### Rule 3b(vi) — Verification Log (§7) vs acceptance criteria (§4)

| V-row | supports | planner re-measurement |
|---|---|---|
| V1/V2 | premise (no gate; grep control) | `check-no-personal-email` 0 outside `design_docs`; `verify_ail.sh` in `.github` 3 — AGREES |
| V3 | D4 (no Makefile) | AGREES |
| V4/V16 | D1/D5 (fleet precedent text) | read from the fleet's `origin/dev`; regex, exclusion, arms A–E AGREE; the ✗ line's colour placement is P3 |
| V5 | AC-M1-3, AC-M1-10 | one in-scope hit, `scripts/mission_answer.sh` line 69 — AGREES |
| V6 | D3 | not re-read (fleet's line 108); the noreply form is the one D3 prescribes and the gate excludes it — consistent |
| V7 | D3 / AC-M1-5 | arms 1–5 all set `MISSION_ATTENDED_EMAIL`; 16 assertions — AGREES |
| V8/V14 | D2 / AC-M1-9 | 34 = 6 + 28; 0 root-level mission docs — AGREES |
| V9 | AC-M1-8 | rc=0 ×3 with the pin; 19 packages; 11/40 — AGREES (P10: the pin gates `verify_ail.sh` too) |
| V13 | D3 / M9 | identity not printed at base; `attended identity` appears once in the diff body — AGREES, P12 |
| V15 | D2 / row 91 | 8 = 1 + 7 files — AGREES (P11) |

### Rule 3i — for each §5 row the OBSERVABLE the arm reads is downstream of the mechanism removed

| § | mechanism removed | observable read | downstream? |
|---|---|---|---|
| M1 | the `exit 1` on the hit path | `rc=` field of the probe | yes — rc is that statement's value |
| M2 | the exclusion filter | `hit=`/`rc=` on allowed tokens (C, D) and the real tree (A) | yes — every hit line passes through the filter |
| M3 | `scripts/.*` in `SCOPE_RE` | `hit=` on a `scripts/` fixture (G) | yes — the file list is the regex's output |
| M4 | the floor block | `rc=2`, `floor=1`, `ok=0` (F) | yes — only the floor prints that line / exits 2 |
| M5 | `hits` increment | `rc=`/`ok=` (the ✓ path is taken when hits==0) | yes |
| M6 | the `$` anchor | `hit=` on `a@example.com.evil.net` (H) | yes |
| M7 | the noreply default byte | AC-M1-2 rc, AC-M1-10 count, arm A, 6b | yes — all read the same line |
| M9 | the note line | 6b's grep for `--dry-run: attended identity … <…>` | yes |

### pi false-green (5) — writes outside the worktree

The executor's writes are confined to `/Users/voightkampff/dev/sunholo-data/.wt-world-iter175/`:
`scripts/` (4 files), `.github/workflows/ci.yml`, `.snap/M1/`, the suite's transient
`.cnpe_scratch.*`, and the mutation backups in the SIBLING `…/.wt-world-iter175-mutbak/`. After the
run the controller reads `git -C /Users/voightkampff/dev/sunholo-data/ailang-world status --short`
and requires it EMPTY, and `git -C …/.wt-world-iter175 status --porcelain` to list exactly:
`?? scripts/check_no_personal_email.sh`, `?? scripts/test_check_no_personal_email.sh`,
` M scripts/mission_answer.sh`, ` M scripts/test_mission_answer.sh`, ` M .github/workflows/ci.yml`,
`?? .snap/` — nothing else (no `.cnpe_scratch.*`, no `*.bak`).

---

## Rig constraints that bind every command in this plan

| # | Constraint | Consequence |
|---|---|---|
| 1 | interactive `grep` is a ugrep wrapper; child `bash` gets `/usr/bin/grep` | the instrument and suite call `/usr/bin/grep` explicitly (present on macOS and `ubuntu-latest`); commands the executor records use `/usr/bin/grep` |
| 2 | login shell zsh; `/bin/bash` is 3.2.57 | every script runs as `bash <script>`; no `mapfile`, `${var,,}`, associative arrays, `local -n`; `<<<` here-strings are fine |
| 3 | no `timeout(1)`, no `gdate`; BSD `sed`/`date`/`cat` | mutants are made with `sed '…' backup > file` (never `sed -i`, never `mv`); no `cat -A` |
| 4 | `export PATH=/opt/homebrew/bin:$PATH` | required for `go`, `gh`, `jq` |
| 5 | `export AILANG_BIN=/Users/voightkampff/.pinned-ailang/ailang` | required for `go test` AND `verify_ail.sh` (P10) |
| 6 | `scripts/verify_go.sh` is NOT a gate (row 76) | never run as one; `go build` is not a compile fence for `_test.go` — `go vet ./...` is |
| 7 | `tools/launchd/*` FROZEN | `git diff --stat a295291 -- tools/launchd/` empty at every gate |
| 8 | worktrees/state never under `/tmp` (row 71) | suite scratch inside the worktree (`.cnpe_scratch.XXXXXX`, trap-removed); backups in the sibling `-mutbak` dir; the base tree for AC-M1-3 is a sibling; gate logs may go to `/tmp` (transient) |
| 9 | row 82: never `git checkout --` on uncommitted work | every restore is `cp` from the `-mutbak` backup, asserted by `shasum -a 256 -c` AND `test -x` |
| 10 | the executor is unsandboxed pi | NO git write of any kind (`add`, `commit`, `stash`, `checkout`, `reset`, `branch`, `worktree`, `mv` under git); AC-M1-3 and AC-M1-7's `git diff` are READ-only and allowed; AC-M1-3's worktree add/remove is the CONTROLLER's |
| 11 | CI has no TTY | the instrument prints colour only when `[ -t 1 ]` (P3) |

---

## Sequencing

**M1 (executor) → controller reconstructs the commit from `.snap/M1`, runs AC-M1-3 with the real
worktree procedure, opens the PR, lands → M2 (controller, charter only).** M2's Repo Profile bullet
names `scripts/check_no_personal_email.sh`; writing it before the script is on `dev` would prescribe
a `No such file`. Row 91 is filed in the same M2 commit.

---

## M1 — the instrument, its suite, the default fix, ARM 6 and the two CI steps (executor, ~0.25d)

### Files

| Path | Action | Mode | LOC | Content |
|---|---|---|---|---|
| `scripts/check_no_personal_email.sh` | **create** | 0755 | 76 | the text in §M1-A, byte-for-byte |
| `scripts/test_check_no_personal_email.sh` | **create** | 0755 | 100 | the text in §M1-B, byte-for-byte |
| `scripts/mission_answer.sh` | **modify** | keep 0755 | 146 → 147 | line 69 replaced; one line inserted after line 123 (§M1-C) |
| `scripts/test_mission_answer.sh` | **modify** | keep 0755 | 83 → 92 | ARM 6 block inserted after line 79 (§M1-D) |
| `.github/workflows/ci.yml` | **modify** | — | 216 → 224 | two steps appended at EOF (§M1-E) |

### Order of work inside M1 — do not reorder

1. `cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter175; export PATH=/opt/homebrew/bin:$PATH;
   export AILANG_BIN=/Users/voightkampff/.pinned-ailang/ailang`. Confirm `git rev-parse HEAD` =
   `f53b353…` and `git status --porcelain` is empty. Record the base readings you will compare
   against: `bash scripts/test_mission_answer.sh | tail -1` (`16 passed, 0 failed`),
   `git ls-files | /usr/bin/grep -cE '^(design_docs/[^/]*mission[^/]*\.md|scripts/.*)$'` (34),
   `/usr/bin/grep -c 'check_no_personal_email' .github/workflows/ci.yml` (0).
2. Write `scripts/check_no_personal_email.sh` (§M1-A) with a quoted heredoc (`cat > … <<'EOF'`);
   `chmod +x`; `bash -n`. Run it: `bash scripts/check_no_personal_email.sh; echo rc=$?` — expect
   **rc=1** naming `scripts/mission_answer.sh` (the before-state; the instrument is untracked so it
   is not yet scanned). The FAILED line must read `1 occurrence(s) in 34 in-scope files`.
3. Apply §M1-C to `scripts/mission_answer.sh`. Re-run step 2's command — expect **rc=0** and the ✓
   line (AC-M1-2). `test -x scripts/mission_answer.sh` must still hold.
4. Write `scripts/test_check_no_personal_email.sh` (§M1-B); `chmod +x`; `bash -n`. Run
   `bash scripts/check_no_personal_email.sh` AGAIN (P1/P2: the suite's own text is in scope once
   tracked; the gate must stay ✓ — it scans tracked files only, so ALSO run the PAT probe in AC-M1-10b
   over both new files and require 0 tokens). Then `bash scripts/test_check_no_personal_email.sh`
   → `8 passed, 0 failed`.
5. Apply §M1-D to `scripts/test_mission_answer.sh`; run it → `18 passed, 0 failed`.
6. Append §M1-E to `.github/workflows/ci.yml`; `/usr/bin/grep -c 'check_no_personal_email'` → 2.
7. Run every AC-M1 command (table below) and record rc + output.
8. Run the mutation drill (protocol below) for ALL of M1–M9; restore and re-verify after each.
9. Run the M1 boundary gate. Snapshot into `.snap/M1/` and write `.snap/M1/REPORT.md`. **No git write.**

### §M1-A — `scripts/check_no_personal_email.sh` (complete; 76 lines; planner-run: clean tree ✓ rc=0, base tree ✗ rc=1, empty scope rc=2)

```bash
#!/usr/bin/env bash
# check_no_personal_email.sh — refuse a personal email address in the surface the LOOP writes.
#
# WHY (the shared mission-control skill's ATTENDED LEDGER EDITS contract, Mark, attended
# 2026-09-02): ledger rows are pasted VERBATIM into the public bookkeeping issue by every mission
# report, so an address written into an evidence cell is published. The rule says compare, never
# record; this gate is what makes that stick, because the generator is automated and prose review
# is not. World carried a personal address in 7 doc locations at iter-152 (row 74); the redaction
# was a one-time sweep with nothing behind it until this gate.
#
# SCOPE IS DELIBERATELY NARROW: the loop-written surface only (World's own SCOPE_RE, row 74 D2):
#   * design_docs/*mission*.md — the six mission docs the loop writes (charter, log, archives,
#     dashboard, index)
#   * scripts/**               — World-owned instruments the loop writes and CI runs
# It does NOT police design_docs/planned|implemented|verification (sprint artifacts; a full-tree
# scan reads 7 non-address false positives there today — SSH remote strings, Go toolchain
# module ids), Go/.ail sources, workflows, go.mod/go.sum, docs/**, README.md or CLAUDE.md
# (deliberate public contact points — a governance decision, not a lint). Widening this without
# ruling on those first would just fail the build; it is a separate queue row.
#
# Allowed: GitHub noreply addresses (attributable, expose nothing), example.com/org/net and
# .invalid/.test/.localhost placeholders (RFC 2606/6761 reserved, cannot be real), and machine
# identities such as GCP service accounts — none of which is a person. The reserved-TLD exclusion
# is ANCHORED (`@example\.(com|org|net)$`): an address at example.com.evil.net is a HIT. World's one
# deliberate deviation from the fleet precedent (sunholo-data/ailang f0d44915d), pinned by arm H.
#
# Exit codes: 0 clean; 1 at least one address hit; 2 instrument failure (zero in-scope files —
# a gate that scans nothing must never print ✓).
# bash 3.2 safe: set -u, no associative arrays, no mapfile, no ${var,,}. Colour only on a TTY so
# captured output is the bare contract lines.

set -u
if [ -t 1 ]; then RED=$'\033[0;31m'; GREEN=$'\033[0;32m'; RESET=$'\033[0m'; else RED=''; GREEN=''; RESET=''; fi
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT" || exit 1

# The loop writes here (World's own scope, row 74 D2 — no .claude/skills/: World has no such tree).
SCOPE_RE='^(design_docs/[^/]*mission[^/]*\.md|scripts/.*)$'
PAT='[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}'

files="$(git ls-files | /usr/bin/grep -E "$SCOPE_RE")"
n_scope=$(printf '%s\n' "$files" | /usr/bin/grep -c .)
if [ "$n_scope" -eq 0 ]; then
    echo "${RED}✗ check-no-personal-email: zero in-scope files scanned — the gate must never print ✓ over an empty surface${RESET}"
    exit 2
fi

hits=0
while IFS= read -r f; do
    [ -n "$f" ] || continue
    case "$f" in
        *.png|*.jpg|*.gif|*.pdf|*.zip) continue ;;
    esac
    [ -f "$f" ] || continue
    found=$(LC_ALL=C /usr/bin/grep -oE "$PAT" "$f" 2>/dev/null \
        | /usr/bin/grep -vE 'users\.noreply\.github\.com|noreply@|@example\.(com|org|net)$|\.(invalid|test|localhost)$|gserviceaccount\.com|@sentry\.io' \
        | sort -u)
    if [ -n "$found" ]; then
        while IFS= read -r addr; do
            [ -n "$addr" ] || continue
            echo "${RED}✗ personal email in loop-written file:${RESET} $f  ->  $addr"
            hits=$((hits+1))
        done <<< "$found"
    fi
done <<< "$files"

if [ "$hits" -gt 0 ]; then
    echo ""
    echo "${RED}check-no-personal-email FAILED${RESET} ($hits occurrence(s) in $n_scope in-scope files)."
    echo "Ledger rows are pasted verbatim into the public bookkeeping issue — an address here is published."
    echo "Record the provenance VERDICT (\"attended\" / \"fleet\"), never the address."
    echo "Use a GitHub noreply address where an identity is genuinely needed."
    exit 1
fi
echo "${GREEN}✓ check-no-personal-email: no personal addresses in the loop-written surface${RESET}"
exit 0
```

Contract notes the executor must not "improve": the ✓, ✗ and floor lines are D6's verbatim contract
and the suite greps ASCII substrings of them; `n_scope` is counted BEFORE the loop so the floor fires
on an empty list; the header names no address-shaped token (P1); colour is TTY-only (P3);
`/usr/bin/grep` everywhere (rig row 1); `LC_ALL=C` on the PAT grep as in the fleet.

### §M1-B — `scripts/test_check_no_personal_email.sh` (complete; 100 lines; planner-run: `8 passed, 0 failed`, self-scan clean, matrix as in P4)

```bash
#!/usr/bin/env bash
# test_check_no_personal_email.sh — 8 arms (A–H) for scripts/check_no_personal_email.sh (row 74).
#
# Rule 3k: every arm runs the REAL instrument as a process — `bash "$SCRIPT_UT"` from inside a
# throwaway `git init` repo (arms B–H) or the real repo (arm A) — and asserts rc AND the three
# contract lines (D6) via one `probe` string per arm: "rc=N hit=N ok=N floor=N", where hit = count
# of per-hit lines, ok = count of the clean line, floor = count of the empty-scope floor line.
# Every fixture address is ASSEMBLED AT RUNTIME by concatenation: a literal address in this file
# would be flagged by the very gate it tests, because this file lives in scripts/ and is in scope.
# bash 3.2 safe. Scratch lives inside the worktree (never /tmp, row 71) and is removed on exit.
set -uo pipefail
cd "$(dirname "$0")/.." || exit 1
ROOT="$(pwd)"
SCRIPT_UT="$ROOT/scripts/check_no_personal_email.sh"
pass=0; fail=0
ok()   { echo "ok $*";      pass=$((pass+1)); }
notok(){ echo "not ok $*";  fail=$((fail+1)); }
ck() { if [ "$2" = "$3" ]; then ok "$1"; else notok "$1 (got '$2' want '$3')"; fi; }

SCRATCH="$(mktemp -d "$ROOT/.cnpe_scratch.XXXXXX")"
trap 'cd "$ROOT"; rm -rf "$SCRATCH"' EXIT

# probe <instrument-path>: run it from the CURRENT directory, summarise rc + the contract lines.
probe() {
    local out rc
    out="$(bash "$1" 2>&1)"; rc=$?
    printf 'rc=%s hit=%s ok=%s floor=%s\n' "$rc" \
        "$(printf '%s\n' "$out" | /usr/bin/grep -c 'personal email in loop-written file:')" \
        "$(printf '%s\n' "$out" | /usr/bin/grep -c 'check-no-personal-email: no personal addresses in the loop-written surface')" \
        "$(printf '%s\n' "$out" | /usr/bin/grep -c 'zero in-scope files scanned')"
}
stage() { git add -A >/dev/null 2>&1; }

# Fixture addresses, assembled at runtime by splitting at the "@" (see header): neither half is
# address-shaped on its own, so the gate reads this file as clean.
PERSONAL="someone";      PERSONAL="${PERSONAL}@realdomain.com"
NOREPLY="3155884+SomeUser"; NOREPLY="${NOREPLY}@users.noreply.github.com"
RESERVED_COM="you";      RESERVED_COM="${RESERVED_COM}@example.com"
RESERVED_ORG="you";      RESERVED_ORG="${RESERVED_ORG}@example.org"
RESERVED_NET="you";      RESERVED_NET="${RESERVED_NET}@example.net"
RESERVED_INV="test";     RESERVED_INV="${RESERVED_INV}@example.invalid"
PREFIX_ONLY="a";         PREFIX_ONLY="${PREFIX_ONLY}@example.com.evil.net"

# A — the real repo is clean (the gate's live assertion; this is what CI's first new step runs).
ck "A real repo clean" "$(probe "$SCRIPT_UT")" "rc=0 hit=0 ok=1 floor=0"

# B–E, G, H run in ONE throwaway repo whose scripts/ holds a copy of the instrument (in scope,
# clean by construction — the instrument's own text carries no address-shaped token).
W="$SCRATCH/repo"; mkdir -p "$W/design_docs" "$W/scripts"
cp "$SCRIPT_UT" "$W/scripts/check_no_personal_email.sh"
cd "$W" || exit 1
git init -q .
git -c user.email=t@t -c user.name=t commit -q --allow-empty -m init >/dev/null 2>&1
stage

# B — a personal address in a mission doc must FAIL (the whole point).
echo "provenance is the commit author $PERSONAL" > design_docs/world-mission.md; stage
ck "B personal addr in mission doc -> rc 1" "$(probe scripts/check_no_personal_email.sh)" "rc=1 hit=1 ok=0 floor=0"

# C — a GitHub noreply is ALLOWED (identities may still be recorded).
echo "author $NOREPLY" > design_docs/world-mission.md; stage
ck "C noreply allowed" "$(probe scripts/check_no_personal_email.sh)" "rc=0 hit=0 ok=1 floor=0"

# D — reserved-TLD placeholders are allowed (RFC 2606/6761): .com/.org/.net/.invalid.
echo "to: $RESERVED_COM $RESERVED_ORG $RESERVED_NET and $RESERVED_INV" > design_docs/world-mission.md; stage
ck "D reserved placeholders allowed" "$(probe scripts/check_no_personal_email.sh)" "rc=0 hit=0 ok=1 floor=0"

# E — OUT OF SCOPE by design: governance/code/planned-doc files are not policed. A clean mission
# doc stays in place so the scope is non-empty under every mutation (the floor is arm F's, not E's).
echo "# charter (clean)" > design_docs/world-mission.md
mkdir -p design_docs/planned
echo "maintainer $PERSONAL" > SECURITY.md
echo "const owner = \"$PERSONAL\"" > access_control.go
echo "reviewer $PERSONAL" > design_docs/planned/some-design.md
stage
ck "E governance/code/planned out of scope" "$(probe scripts/check_no_personal_email.sh)" "rc=0 hit=0 ok=1 floor=0"

# G — a hit in scripts/ (not only in a mission doc) is caught: scripts/.* is a live alternative.
rm -f SECURITY.md access_control.go design_docs/planned/some-design.md
echo "ATT_EMAIL=\"\${MISSION_ATTENDED_EMAIL:-$PERSONAL}\"" > scripts/some_tool.sh; stage
ck "G hit in scripts/ caught" "$(probe scripts/check_no_personal_email.sh)" "rc=1 hit=1 ok=0 floor=0"

# H — the reserved-TLD exclusion is ANCHORED: example.com as a mere PREFIX of the domain is a HIT.
rm -f scripts/some_tool.sh
echo "contact $PREFIX_ONLY" > design_docs/world-mission.md; stage
ck "H example.com-as-prefix is a HIT (anchored)" "$(probe scripts/check_no_personal_email.sh)" "rc=1 hit=1 ok=0 floor=0"

# F — the anti-vacuity floor: a repo with ZERO in-scope tracked files reads rc 2, no ✓ line.
# The instrument lives OUTSIDE scripts/ here (tools/ is not in SCOPE_RE) so the scope is empty.
W2="$SCRATCH/empty"; mkdir -p "$W2/tools"
cp "$SCRIPT_UT" "$W2/tools/check_no_personal_email.sh"
cd "$W2" || exit 1
git init -q .
echo "# nothing the loop writes" > README.md; stage
ck "F empty scope -> rc 2, floor line, no clean line" "$(probe tools/check_no_personal_email.sh)" "rc=2 hit=0 ok=0 floor=1"

cd "$ROOT" || exit 1
echo "---"
echo "$pass passed, $fail failed"
[ "$fail" -eq 0 ]
```

### §M1-C — `scripts/mission_answer.sh`: the two-line change (D3)

Line 69 today is `ATT_EMAIL="${MISSION_ATTENDED_EMAIL:-<the line-69 personal default>}"`; line 123
is `	diff -u "$DOC" "$TMP" | head -40` and line 124 is `	note "--dry-run: $DOC not modified"` (both
TAB-indented, inside `if [ "$DRYRUN" -eq 1 ]; then`). Splice with `sed -n` ranges into a sibling
backup dir and write back with `cat >` so the inode and 0755 mode survive (P14):

```bash
BAK=/Users/voightkampff/dev/sunholo-data/.wt-world-iter175-mutbak; mkdir -p "$BAK"
cp scripts/mission_answer.sh "$BAK/mission_answer.sh.orig"
{ sed -n '1,68p' "$BAK/mission_answer.sh.orig"
  printf '%s\n' 'ATT_EMAIL="${MISSION_ATTENDED_EMAIL:-3155884+MarkEdmondson1234@users.noreply.github.com}"'
  sed -n '70,123p' "$BAK/mission_answer.sh.orig"
  printf '\t%s\n' 'note "--dry-run: attended identity $ATT_NAME <$ATT_EMAIL>"'
  sed -n '124,$p' "$BAK/mission_answer.sh.orig"
} > "$BAK/mission_answer.sh.new"
cat "$BAK/mission_answer.sh.new" > scripts/mission_answer.sh
diff "$BAK/mission_answer.sh.orig" scripts/mission_answer.sh | /usr/bin/grep -c '^[<>]'   # exactly 3 (69 replaced = 2, 124 inserted = 1)
wc -l scripts/mission_answer.sh                                                          # 147
test -x scripts/mission_answer.sh; echo "rc=$?"                                          # 0
sed -n '69p;122,127p' scripts/mission_answer.sh
```

Expected lines 122–127 after the splice:

```bash
if [ "$DRYRUN" -eq 1 ]; then
	diff -u "$DOC" "$TMP" | head -40
	note "--dry-run: attended identity $ATT_NAME <$ATT_EMAIL>"
	note "--dry-run: $DOC not modified"
	exit 0
fi
```

The note uses the VARIABLES, never the literal (P6): `/usr/bin/grep -c 'users\.noreply\.github\.com'
scripts/mission_answer.sh` → 1. The `.orig` copy in `$BAK` holds the personal default — it lives
OUTSIDE the repo, is never copied under the worktree, and the executor deletes `$BAK` as its last act.

### §M1-D — `scripts/test_mission_answer.sh`: ARM 6 (D3), inserted after line 79, before the `echo "---"` at line 81

The suite's conventions (read at base): `SCRIPT_UT=scripts/mission_answer.sh` (invoked directly —
its exec bit matters), `FIX="$(mktemp)"` + `write_fixture` (a two-row ledger with `D-2` OPEN),
`ok "<id> …"` / `notok "<id> …"` print `ok`/`not ok` and count; the summary is `---` then
`N passed, M failed`; the last line `[ "$fail" -eq 0 ]` makes any `not ok` a non-zero exit.
Line 79 is `[ "$before" = "$after" ] && ok "5a --dry-run does not write" || notok "5a --dry-run wrote"`;
line 80 is blank; line 81 is `echo "---"`. Insert this block (the continuation line begins with ONE
TAB, as lines 32 and 46 do) so that the file reads 79, blank, ARM 6, blank, `echo "---"`:

```bash
# ARM 6 — with no MISSION_ATTENDED_EMAIL override, the DEFAULT identity is the GitHub noreply
# form (row 74): a personal address as the default is exactly what the privacy gate forbids.
write_fixture
out=$(env -u MISSION_ATTENDED_EMAIL MISSION_ATTENDED_NAME="Test Human" \
	"$SCRIPT_UT" --id D-2 --answer "x" --file "$FIX" --dry-run 2>&1); rc=$?
[ $rc -eq 0 ] && ok "6a default identity accepted on --dry-run" || notok "6a default identity refused: $out"
echo "$out" | grep -q 'attended identity .*<3155884+MarkEdmondson1234@users.noreply.github.com>' && ok "6b default identity is the noreply form" || notok "6b default identity is not the noreply form"
```

Splice (same `cat >` discipline):

```bash
cp scripts/test_mission_answer.sh "$BAK/test_mission_answer.sh.orig"
cat > "$BAK/arm6.txt" <<'ARM6'
<the 8 lines above, verbatim — including the TAB on the "$SCRIPT_UT" line>
ARM6
{ sed -n '1,80p' "$BAK/test_mission_answer.sh.orig"; cat "$BAK/arm6.txt"; echo; sed -n '81,$p' "$BAK/test_mission_answer.sh.orig"; } > "$BAK/test_mission_answer.sh.new"
cat "$BAK/test_mission_answer.sh.new" > scripts/test_mission_answer.sh
wc -l scripts/test_mission_answer.sh        # 92
bash scripts/test_mission_answer.sh | tail -3   # ok 6b …, ---, 18 passed, 0 failed
```

`env -u` is BSD- and GNU-safe (measured on the rig). 6b's needle carries the `<…>` tail so the diff
body's own "attended identity" prose cannot satisfy it (P12).

### §M1-E — `.github/workflows/ci.yml`: two steps appended at EOF (after line 216, inside `go-verify`)

The job's last step today (lines 214–216, six spaces before `- name:`, eight before `timeout-minutes:`/`run:`):

```yaml
      - name: Gate-0 self-notice instrument suite
        timeout-minutes: 2
        run: bash scripts/test_gate0_self_notices.sh
```

Append exactly (a blank separator line first; the file already ends with a newline):

```yaml

      - name: Refuse a personal email in the loop-written surface
        timeout-minutes: 2
        run: bash scripts/check_no_personal_email.sh

      - name: Personal-email gate self-test
        timeout-minutes: 2
        run: bash scripts/test_check_no_personal_email.sh
```

Planner-validated: the appended file parses (ruby `YAML.load_file`) with `go-verify` at 16 steps, the
last two named as above; `wc -l` 224; `/usr/bin/grep -c 'check_no_personal_email'` → 2.

### Acceptance criteria — quoted from §4 by AC id (executor runs each, records rc)

```bash
export PATH=/opt/homebrew/bin:$PATH; export AILANG_BIN=/Users/voightkampff/.pinned-ailang/ailang
test -x scripts/check_no_personal_email.sh; echo "AC-M1-1 rc=$?"                                   # 0
test -x scripts/test_check_no_personal_email.sh; echo "AC-M1-1b rc=$?"                             # 0 (both new scripts are 0755)
bash scripts/check_no_personal_email.sh; echo "AC-M1-2 rc=$?"                                      # ✓ line; 0
# AC-M1-3: CONTROLLER step (needs `git worktree add`) — see below; the executor records "deferred to controller"
bash scripts/test_check_no_personal_email.sh | tail -1; echo "AC-M1-4 rc=${PIPESTATUS[0]}"         # 8 passed, 0 failed; 0
bash scripts/test_mission_answer.sh | tail -1; echo "AC-M1-5 rc=${PIPESTATUS[0]}"                  # 18 passed, 0 failed; 0
/usr/bin/grep -c 'check_no_personal_email' .github/workflows/ci.yml                               # AC-M1-6: 2
git diff --stat a295291 -- tools/launchd/ | wc -l                                                  # AC-M1-7: 0
go vet ./... ; echo "AC-M1-8a rc=$?"                                                               # 0
go test ./... -count=1 -timeout 600s > /tmp/gotest175.log 2>&1; echo "AC-M1-8b rc=$?"; /usr/bin/grep -c '^ok' /tmp/gotest175.log   # 0; 19
bash ./scripts/verify_ail.sh > /tmp/va175.log 2>&1; echo "AC-M1-8c rc=$?"; tail -1 /tmp/va175.log  # 0; "11 required identities verified, 40 named tests pass"
git ls-files | /usr/bin/grep -cE '^(design_docs/[^/]*mission[^/]*\.md|scripts/.*)$'              # AC-M1-9: 34 UNTIL the controller tracks the two scripts; the executor ALSO runs:
{ git ls-files; echo scripts/check_no_personal_email.sh; echo scripts/test_check_no_personal_email.sh; } | /usr/bin/grep -cE '^(design_docs/[^/]*mission[^/]*\.md|scripts/.*)$'   # 36 (the post-commit count, simulated without a git write)
/usr/bin/grep -c 'users\.noreply\.github\.com' scripts/mission_answer.sh                          # AC-M1-10a: 1
# AC-M1-10b — the in-scope scan INCLUDING the two untracked new files (they become tracked at the commit):
PAT='[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}'
{ git ls-files; echo scripts/check_no_personal_email.sh; echo scripts/test_check_no_personal_email.sh; } | /usr/bin/grep -E '^(design_docs/[^/]*mission[^/]*\.md|scripts/.*)$' \
  | while IFS= read -r f; do LC_ALL=C /usr/bin/grep -oE "$PAT" "$f" | /usr/bin/grep -vE 'users\.noreply\.github\.com|noreply@|@example\.(com|org|net)$|\.(invalid|test|localhost)$|gserviceaccount\.com|@sentry\.io' | sed "s|^|$f: |"; done | wc -l   # 0
```

AC-M1-9 is honest only after the commit; the executor records BOTH numbers (34 tracked-now, 36
simulated) and the controller re-reads 36 after reconstruction.

**AC-M1-3 — controller, after reconstruction, before the PR** (the doc's snippet, with the sibling
path and the file-name assertion of P3/P7):

```bash
cd /Users/voightkampff/dev/sunholo-data/.wt-world-iter175
BASEWT="$(cd "$(git rev-parse --show-toplevel)/.." && pwd)/.base-world-iter175"
git worktree add --detach "$BASEWT" a295291bea069d793960e08f055e501c6f9661f4
cp scripts/check_no_personal_email.sh "$BASEWT/scripts/"
( cd "$BASEWT" && bash scripts/check_no_personal_email.sh > /tmp/ac3.log 2>&1; echo "rc=$?" )        # rc=1
/usr/bin/grep -c 'loop-written file: scripts/mission_answer.sh' /tmp/ac3.log                         # 1
/usr/bin/grep -c 'in 34 in-scope files' /tmp/ac3.log                                                  # 1
( cd "$BASEWT" && git ls-files | /usr/bin/grep -c check_no_personal_email )                          # 0 — the copy is untracked
git worktree remove --force "$BASEWT"
```

Never print `/tmp/ac3.log` into a report: its ✗ line carries the base address. Record the four
numbers only.

### Mutation drill — protocol and the nine rows (executor runs ALL; measured red sets from P4)

```bash
BAK=/Users/voightkampff/dev/sunholo-data/.wt-world-iter175-mutbak; mkdir -p "$BAK"
cp scripts/check_no_personal_email.sh "$BAK/cnpe.bak"; shasum -a 256 scripts/check_no_personal_email.sh > "$BAK/cnpe.sha"
cp scripts/mission_answer.sh "$BAK/ma.bak";               shasum -a 256 scripts/mission_answer.sh > "$BAK/ma.sha"
drill() {   # drill <label> — runs the gate suite, prints the red set, restores, asserts digest AND mode
  bash scripts/test_check_no_personal_email.sh > "$BAK/$1.log" 2>&1; echo "$1 suite rc=$? :: $(tail -1 "$BAK/$1.log") :: red={ $(/usr/bin/grep '^not ok' "$BAK/$1.log" | cut -d' ' -f3 | tr '\n' ' ')}"
  cp "$BAK/cnpe.bak" scripts/check_no_personal_email.sh; cp "$BAK/ma.bak" scripts/mission_answer.sh
  (cd "$(git rev-parse --show-toplevel)" && shasum -a 256 -c "$BAK/cnpe.sha" "$BAK/ma.sha")        # both "OK" — else STOP
  test -x scripts/check_no_personal_email.sh && test -x scripts/mission_answer.sh && echo "mode ok"  # else STOP
}
B="$BAK/cnpe.bak"; T=scripts/check_no_personal_email.sh
sed 's/^    exit 1$/    exit 0/' "$B" > "$T";                                                        drill M1   # red={ B G H }
/usr/bin/grep -v "grep -vE 'users" "$B" > "$T";                                                      drill M2   # red={ A C D }
sed "s|^SCOPE_RE=.*|SCOPE_RE='^(design_docs/[^/]*mission[^/]*\\\\.md)\$'|" "$B" > "$T";              drill M3   # red={ G }
sed '/^if \[ "\$n_scope" -eq 0 \]; then$/,/^fi$/d' "$B" > "$T";                                       drill M4   # red={ F }
sed 's/^            hits=\$((hits+1))$/            :/' "$B" > "$T";                                    drill M5   # red={ B G H }  (probe shows ok=1 — distinct from M1)
sed 's/@example\\\.(com|org|net)\$|/@example\\.(com|org|net)|/' "$B" > "$T";                          drill M6   # red={ H }
# M7 — the live-gate kill: revert the default to a RAW (placeholder) address; four readers red
sed '69s|3155884+MarkEdmondson1234@users.noreply.github.com|someone@realdomain.com|' "$BAK/ma.bak" > scripts/mission_answer.sh
bash scripts/check_no_personal_email.sh | /usr/bin/grep -c 'loop-written file: scripts/mission_answer.sh'; echo "M7 live rc=${PIPESTATUS[0]}"   # 1; rc=1 (AC-M1-2 RED)
/usr/bin/grep -c 'users\.noreply\.github\.com' scripts/mission_answer.sh                                                                          # 0 (AC-M1-10a RED)
bash scripts/test_mission_answer.sh | /usr/bin/grep -E '^not ok|passed'                                                                          # not ok 6b …; 17 passed, 1 failed
drill M7                                                                                                                                         # red={ A }
sed 's/^# Allowed: GitHub noreply addresses/# Permitted: GitHub noreply addresses/' "$B" > "$T";      drill M8   # red={ }  (control)
# M9 — drop the dry-run identity note
sed '/--dry-run: attended identity/d' "$BAK/ma.bak" > scripts/mission_answer.sh
bash scripts/test_mission_answer.sh | /usr/bin/grep -E '^not ok|passed'                                                                          # not ok 6b …; 17 passed, 1 failed
drill M9                                                                                                                                         # red={ } (the gate suite is untouched by M9)
bash scripts/test_check_no_personal_email.sh | tail -1; bash scripts/test_mission_answer.sh | tail -1                                            # 8 passed, 0 failed; 18 passed, 0 failed — pristine
```

Each `sed` above was run by the planner against the §M1-A/§M1-C texts and produced the listed
set. A drill that reds MORE or FEWER arms than listed is a finding for `.snap/M1/REPORT.md`, never
a number to adjust. M7's `someone@realdomain.com` is a placeholder, not a person; it exists in the
tree only between the `sed` and the `cp` restore, and `drill` restores `scripts/mission_answer.sh`
too.

### M1 boundary gate — all green before the snapshot

```bash
export PATH=/opt/homebrew/bin:$PATH; export AILANG_BIN=/Users/voightkampff/.pinned-ailang/ailang
bash -n scripts/check_no_personal_email.sh; echo "rc=$?"                                              # 0
bash -n scripts/test_check_no_personal_email.sh; echo "rc=$?"                                         # 0
bash -n scripts/mission_answer.sh; bash -n scripts/test_mission_answer.sh; echo "rc=$?"               # 0
bash scripts/check_no_personal_email.sh; echo "rc=$?"                                                 # ✓; 0
bash scripts/test_check_no_personal_email.sh | tail -1                                                # 8 passed, 0 failed
bash scripts/test_mission_answer.sh | tail -1                                                         # 18 passed, 0 failed
bash scripts/test_gate1_range_check.sh | tail -1                                                      # 18 passed, 0 failed
bash scripts/test_queue_census.sh | tail -1                                                           # 68 passed, 0 failed
bash scripts/test_gate0_self_notices.sh | tail -1                                                     # 88 passed, 0 failed
bash ./scripts/verify_ail.sh > /tmp/va175.log 2>&1; echo "rc=$?"; tail -1 /tmp/va175.log              # 0; 11 identities / 40 named tests
go vet ./...; echo "rc=$?"                                                                            # 0
go test ./... -count=1 -timeout 600s > /tmp/gotest175.log 2>&1; echo "rc=$?"                          # 0 (19 packages)
git diff --stat a295291 -- tools/launchd/ | wc -l                                                     # 0
git diff --quiet a295291 -- design_docs/world-mission.md; echo "rc=$?"                                # 0 — M1 touches no charter byte
git status --porcelain                                                                                # exactly the 5 files + ?? .snap/ (see the false-green section); no .cnpe_scratch.*
ls -d .cnpe_scratch.* 2>/dev/null | wc -l                                                             # 0
```

### Snapshot and report — the executor's last step

```bash
mkdir -p .snap/M1/scripts .snap/M1/.github/workflows
cp scripts/check_no_personal_email.sh      .snap/M1/scripts/check_no_personal_email.sh
cp scripts/test_check_no_personal_email.sh .snap/M1/scripts/test_check_no_personal_email.sh
cp scripts/mission_answer.sh               .snap/M1/scripts/mission_answer.sh
cp scripts/test_mission_answer.sh          .snap/M1/scripts/test_mission_answer.sh
cp .github/workflows/ci.yml                .snap/M1/.github/workflows/ci.yml
(cd .snap/M1 && shasum -a 256 scripts/* .github/workflows/ci.yml) > .snap/M1/SHA256
ls -l .snap/M1/scripts/                                                                               # all four scripts -rwxr-xr-x
rm -rf /Users/voightkampff/dev/sunholo-data/.wt-world-iter175-mutbak                                  # holds the .orig with the personal default — delete it
```

`.snap/M1/REPORT.md` records, in this order: (1) `git rev-parse HEAD` (`f53b353…`); (2) every
boundary-gate command with rc and last line; (3) every AC-M1 command with observed rc/output
(AC-M1-3 "deferred to controller"; AC-M1-9 both numbers); (4) the two suite lines `8 passed, 0 failed`
and `18 passed, 0 failed`; (5) every mutation M1–M9 with its OBSERVED red set beside the expected
one, and the `shasum -c` OK + `mode ok` lines after each restore; (6) the PAT probe over both new
files (0 tokens); (7) any finding. It must NOT contain the base address, the `.orig` contents, or
`/tmp/ac3.log`. **No git write: no `add`, `commit`, `stash`, `checkout`, `reset`, `branch`,
`worktree`, `mv` under git.**

### Controller reconstruction (between M1 and M2)

In the sprint worktree: verify `.snap/M1/SHA256` against the tree (`shasum -a 256 -c`), `test -x` on
the four scripts, `git status --porcelain` matches the false-green list; then `git add` the FIVE
paths by name (never `-A`, never `.snap/`), commit, run AC-M1-3 (above), re-read AC-M1-9 → 36, open
the PR against `dev`, wait for 2/2 CI on the head (the ubuntu leg is the first CI run of both new
steps — `ubuntu-latest` has `/usr/bin/grep`, GNU sed is not used), squash-merge, Gate 3b on the merge.

---

## M2 — charter edits (controller, ~0.05d, after M1 is on `dev`)

Runs in the MAIN checkout `/Users/voightkampff/dev/sunholo-data/ailang-world` on `dev` as the
controller's own commit. Three edits, `≤12` lines each; `<PR>`/`<merge-sha>` filled at landing.

**(a) Row 74's leading LANDED tag** — replace the row's first line
`74. **w-the-personal-email-gate-the-shared-skill-cites-as-enforcement-does-not-exist** · clause-2 ·`
with:

```markdown
74. **[LANDED 2026-09-14 (iter-175), CI green on the merge `<merge-sha>` (PR #<PR>), evaluator <verdict>; the gate is `scripts/check_no_personal_email.sh` + `scripts/test_check_no_personal_email.sh` (arms A–H), two `go-verify` steps, the line-69 default is now the GitHub noreply form (ARM 6), the "allow-list for functional defaults" was deliberately NOT built (D3), and widening to `design_docs/**` is row 91]** **w-the-personal-email-gate-the-shared-skill-cites-as-enforcement-does-not-exist** · clause-2 ·
```

**(b) One Repo Profile bullet**, inserted immediately after the row-73 bullet that begins
`- **GATE 0 HAS A SECOND, NO-AUTHORITY READ FOR THE DRIVER'S OWN CRASH NOTICES` (line 137 at base;
insert after that bullet's last line, before the next `- **` bullet):

```markdown
- **THE PERSONAL-EMAIL GATE THE SHARED SKILL CITES IS REAL HERE (row 74, iter-175).**
  `scripts/check_no_personal_email.sh` scans `git ls-files` × `^(design_docs/[^/]*mission[^/]*\.md|scripts/.*)$`
  (the loop-written surface — six mission docs + `scripts/`; NOT `design_docs/planned|implemented`,
  see row 91) for `[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}` minus GitHub noreply,
  `@example.(com|org|net)` ANCHORED at the end, `.invalid/.test/.localhost`, service accounts.
  Exit 0 clean / 1 hit / **2 zero in-scope files** (the floor). CI runs it directly in `go-verify`
  (`Refuse a personal email in the loop-written surface`, then the 8-arm `Personal-email gate
  self-test`); there is no Makefile. Mirror of the fleet's `sunholo-data/ailang` `f0d44915d` with
  two deliberate deviations: the `$` anchor (arm H) and TTY-only colour. Fixtures in any suite
  under `scripts/` are assembled at runtime split at the `@` — a literal address there is a hit.
  `scripts/mission_answer.sh`'s default identity is the noreply form; set `MISSION_ATTENDED_EMAIL`
  for another human. Run `bash scripts/check_no_personal_email.sh` before any ledger edit lands.
```

**(c) New queue row 91**, appended after row 90's text (charter line 5611 block), ≤12 lines:

```markdown
91. **w-personal-email-gate-scope-stops-at-the-loop-written-surface** · clause-2 · **THE ROW-74
    GATE SCANS SIX MISSION DOCS + `scripts/`; `design_docs/planned|implemented|verification/**` AND
    THE ROOT `sprint_*.json` ARE UNSCANNED, AND A FULL-TREE SCAN READS 7 NON-ADDRESS FALSE
    POSITIVES TODAY.** Measured iter-175 at `a295291` (design V15, planner-reproduced): 357 tracked
    files, 8 hits = the one personal default row 74 removed + **7 files** that are not people —
    `git AT github.com` SSH-remote strings ×3 (written without the `@`: the charter is in scope) (`design_docs/planned/w-daemon-timeout-test-flake*` ×2,
    `sprint_w-daemon-timeout-test-flake.json`) and Go toolchain ids `toolchain AT v0.0.1-go1.2x.y.darwin`
    ×4 (`design_docs/implemented/w-defork-pin-redirects-work-repo*` ×3, `w-setup-go-pin-unguarded-sprint-plan.md`).
    **THE ITEM:** rule on each class (exclude `git@` and `toolchain@` by pattern, or by path), widen
    `SCOPE_RE`, add arms for every new exclusion, and consider a CI step for `scripts/test_mission_answer.sh`
    (not in CI today, P8). Sprint artifacts are loop-written too; the fleet's header warns widening
    without ruling first "would just fail the build". · ~0.2d · gated on nothing · surfaced iter-175.
```

**(d) Acceptance — run each in the main checkout, record output:**

```bash
/usr/bin/grep -c '^74\. \*\*\[LANDED 2026-09-14 (iter-175)' design_docs/world-mission.md            # 1
/usr/bin/grep -c 'THE PERSONAL-EMAIL GATE THE SHARED SKILL CITES IS REAL HERE' design_docs/world-mission.md   # 1
/usr/bin/grep -c '^91\. \*\*w-personal-email-gate-scope-stops-at-the-loop-written-surface' design_docs/world-mission.md   # 1
bash scripts/queue_census.sh --doc design_docs/world-mission.md --control-closed 1 --control-open 79 | tail -1   # 44 closed / 91 rows (P17)
bash scripts/mission_decisions.sh --check --file design_docs/world-mission.md; echo "rc=$?"          # 0 (ledger untouched: 22 rows)
bash scripts/check_no_personal_email.sh; echo "rc=$?"                                                # ✓; 0 — the charter is in scope; the new bullet/row must not plant a token
bash scripts/test_check_no_personal_email.sh | tail -1                                               # 8 passed, 0 failed
```

Both false-positive classes are PAT matches (`git@github.com`, `toolchain@v0.0.1-…darwin` — the
planner probed each), so the row-91 text above spells them with ` AT ` — the charter, the log and
the STATUS stamp are all in scope, and a literal there would red the very gate this row lands.
Run the AC-M2 gate line after EVERY charter edit.

---

## Verify gate — the exact list at every commit and on the merge

`bash ./scripts/verify_ail.sh` (rc 0, 11/40, with `AILANG_BIN`) · `go vet ./...` (rc 0) ·
`go test ./... -count=1` (rc 0, 19 packages, with `AILANG_BIN`) · `bash scripts/check_no_personal_email.sh`
(rc 0) · the four suites (`8`, `18`, `18`, `68`, `88` passed / 0 failed) · `tools/launchd/` untouched ·
CI 2/2 on the PR head and SHA-pinned on the merge. `scripts/verify_go.sh` is never run as a gate.

## Commit policy

M1: ONE commit reconstructed by the controller from `.snap/M1` by explicit file list — message
`row 74: a real check-no-personal-email gate over the loop-written surface, wired into go-verify`.
M2: ONE controller commit on `dev` — `docs(mission): row 74 LANDED — the personal-email gate is
real; row 91 files the scope widening (M2)`. Neither commit message may contain an address.

## Velocity

Doc estimate 0.3d. Planner measurement: the instrument and suite are already run-verified here
(76 + 100 lines), the splices are mechanical, the drill is nine `sed` lines. M1 ≈ 0.25d for the pi
lane including drills and report; M2 ≈ 0.05d. Risk: LOW — the only places the pi lane can drift are
the fixture-splitting rule (P2) and the exec bit (P14), both of which the boundary gate catches.

## Start here

Executor: read "Rig constraints", then "Order of work inside M1", then write §M1-A verbatim. Run the
gate after EVERY file you write. Your report is `.snap/M1/REPORT.md`; your last command deletes the
`-mutbak` sibling. You perform no git write. Controller: reconstruct, AC-M1-3, PR, then M2.
