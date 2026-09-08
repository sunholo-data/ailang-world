# Sprint plan: w-queue-closed-count-census

**Iteration:** World 172
**Design:** `design_docs/planned/w-queue-closed-count-census.md` (736 lines; two quorum rounds, r2
closed under the Gate-2 narrow-refinement carve-out with both reviewer fixes applied VERBATIM — §10)
**Owning queue row:** 72, `w-queue-closed-count-is-an-increment-chain-no-instrument-reproduces`
(clause-2, ~0.3d, gated on nothing; charter line 5573)
**Planning base:** `8b151533ed701e6648f7037d86579383b04de407`
**Executor branch:** `sprint/w-queue-census`
**Executor worktree:** `/Users/voightkampff/dev/sunholo-data/.wt-world-iter172`
**Machine plan:** `.ailang/state/sprints/w-queue-closed-count-census.plan.json`
**Status:** planned; nothing implemented, no git write performed

---

## Outcome and scope

One runnable instrument in this repo, the suite that proves it, a frozen snapshot of the charter it
reads, and the rule that says how its number may be published:

```text
scripts/queue_census.sh          # anchored leading-tag classifier; mandatory closed+open controls;
                                 # ordered floors F0-F6; report stamped with the charter's blob OID
scripts/test_queue_census.sh     # 19 self-contained arms, --only <arm>, "N passed, M failed"
scripts/testdata/queue_census_rowstarts_8b15153.md   # generated, 91 lines / 10038 bytes
```

**The instrument does not read the row body at all.** That is the whole design. Row 72 was filed
because iteration 151's attempt 2 classed open row 57 as closed on the strength of the word *landed*
appearing once, lowercase, as an English verb in its body. A prose-guessing classifier has that
defect by construction; a first-token classifier cannot.

Out of scope, explicitly and by ruling: **normalising the 44 UNTAGGED rows** (D2 — the four
position-only moves 65/66/68/70 and the four word changes 7/16/18/57 are a separate docs commit
under D2's recipe), splitting the UNTAGGED bucket, recognising `[BLOCKED …]`/`[AMENDED …]`, and any
claim that the count "moved by exactly one" across two commits (deleted at quorum r1 — no CI step
sees two trees). Re-introducing any of these is scope creep against a settled adjudication.

---

## Baseline table — every acceptance command re-run by the planner at `8b15153`

Charter rule 3e: an AC green at base measures the repo, not the change. The doc's §7 rows B1–B9,
V21 and V22 are **the designer's claims**. Every one below was re-run first-party from
`/Users/voightkampff/dev/sunholo-data/.wt-world-iter172`, under `bash`, with `/usr/bin/grep`.

| # | What I ran | What I saw | Verdict |
|---|---|---|---|
| **B-a** | `bash ./scripts/verify_ail.sh` (AILANG_BIN exported) | `rc=0` — `✓ world package gate PASSED: 9/9 steps performed non-zero work`; `✓ verify gate PASSED: 11 required identities verified, 40 named tests pass` | **GREEN at base** |
| **B-a** | `go vet ./...` | `rc=0`, no output | **GREEN at base** |
| **B-a** | `go test ./... -count=1` **without** `AILANG_BIN` | `rc=1` — **13 FAILs in `host/verifygate`**, each `AILANG_BIN is unset — the shim arms need the pinned released delegate to run the real gate. Never skip.` | **RED — see Finding P1** |
| **B-a** | `AILANG_BIN=/Users/voightkampff/.pinned-ailang/ailang go test ./... -count=1` | `rc=0`, 12 host packages `ok` (`host/verifygate 116.148s`, `host/replay 24.318s`) | **GREEN at base** |
| **B-b** | `test -x scripts/queue_census.sh` / `test -x scripts/test_queue_census.sh` | `rc=1` / `rc=1` | **AGREES** doc B1/B2 |
| **B-c** | `/usr/bin/grep -rn 'queue_census' --exclude-dir=.git .` + control `/usr/bin/grep -rln 'gate1_range_check' --exclude-dir=.git . \| head -3` | **`rc=0`, 50 hits** — every one in `design_docs/planned/w-queue-closed-count-census.md`. Control **fires**: `world-mission-log.md`, `world-mission-index.md`, `world-mission-dashboard.md` | **DISAGREES** with doc B5 — see Finding P2 |
| **B-d** | `/usr/bin/grep -c '^## Queue .*\[ROUTED\]' design_docs/world-mission.md`; control `/usr/bin/grep -c '^## Queue' …` | `0`; control **fires** → `1` | **AGREES** doc B6 + V1 |
| **B-d** | `awk '/^## Queue/{exit} /queue_census\.sh/{c++} END{print c+0}' design_docs/world-mission.md`; control same awk with `/verify_ail\.sh/` | `0`; control **fires** → `9` | **AGREES** doc B7 |
| **B-d** | `awk '/^## Queue/{exit} /charter-blob/{c++} END{print c+0}' design_docs/world-mission.md`; control same awk with `/pinned-ailang/` | `0`; control **fires** → `5` | **AGREES** doc V22 |
| **B-e** | `git hash-object design_docs/world-mission.md` and `git rev-parse 8b15153:design_docs/world-mission.md` | `2dda66a36a6bc60ab280a73b3f0b3c6be89ed399` — **identical by both routes** | **CONFIRMS** doc V21 |
| **B-f** | `bash scripts/test_gate1_range_check.sh` (rc via redirect, not a pipe) | `rc=0`, `18 passed, 0 failed`. Doc's own B3 control `--only report` → `rc=0`, `4 passed, 0 failed` | **AGREES** — the harness shape is live and green |
| **B-g** | The §7 **V9 prototype awk program verbatim**, against `design_docs/world-mission.md` | `89` rows; `closed=41 open=4 untagged=44`; per-tag `LANDED=40 ROUTED=1 PARKED=3 IN-SPRINT=1 -=44`; controls `1 LANDED closed 1637` and `79 PARKED open 5711`; untagged list = the **same 44 numbers** V9 records, in order; rows 16/18/66/68/70 all `- UNTAGGED` | **AGREES** — 41/89/4/44 confirmed |
| **B-h** | `/usr/bin/grep -n 'test_queue_census' .github/workflows/ci.yml`; control `test_gate1_range_check` | `rc=1`; control **fires** → `204:        run: bash scripts/test_gate1_range_check.sh` | **AGREES** doc B4 + V19 |
| **B-h** | `bash scripts/test_queue_census.sh`; `bash scripts/queue_census.sh --doc … --control-closed 1 --control-open 79`; `ls scripts/testdata/queue_census_rowstarts_8b15153.md` | `rc=127`; `rc=127`; `rc=1` | **AGREES** doc B3/B8/B9 |
| **B-i** | §3's fixture command **and** V18's variant, both run and `diff`ed | Both `91` lines / `10038` bytes; `diff` → **IDENTICAL**; sha256 `42fe74686d105871046c84c6c081780ca117d1715f9dd77f4a008ac717fd54c1`. Prototype **on the produced fixture** → `89` rows, `closed=41 open=4 untagged=44`, same 44-number list, numbers exactly `1..89`, zero dups, controls ok | **AGREES** doc V18 — but see Findings P3, P4 |
| **B-j** | `sed -n '200,400p' .github/workflows/ci.yml`; `wc -l` | File **ends at line 204**; the Gate-1 step (202–204) is its last content. Jobs: `ailang-verify` (18), `go-verify` (99) | "immediately after" = **append at EOF** |

**The doc's headline number is confirmed by a third independent run.** The designer's `substr`-loop
prototype, the controller's `sub()`-based `/tmp/proto_census_iter172.awk`, and now the planner's
verbatim re-run of V9's recorded program all read the same bytes and produce
`41 closed / 89 rows (tagged-open 4, untagged 44)` with the same 44-number untagged list. Nothing in
this plan rests on a transcription.

---

## Planner findings

### P1 — `go test ./... -count=1` is RED at base without `AILANG_BIN` *(act on this)*

The role brief names the third gate as bare `go test ./... -count=1`. Run that way at `8b15153` it
is **rc=1** with 13 failures in `host/verifygate`:

```
TestPlainDirtyRefused  TestGitDescribeRefused  TestNoVersionToken  TestNotARelease
TestInScriptControl  TestReleaseChangeNotice  TestEmptyExpectedReleaseSetFailsLoudly
TestFixtureDiscrimination  TestModuleManifestRejectsStrayModule
TestModuleManifestRejectsCaseVariantExtension  TestModuleManifestRejectsDeletedModule
TestModuleManifestEmptyAllowlistFailsLoudly  TestModuleManifestEmptyEnumerationFailsLoudly
```

each saying `AILANG_BIN is unset — the shim arms need the pinned released delegate to run the real
gate. Never skip: a silent skip here is the false-green class this milestone closes.`

With `export AILANG_BIN=/Users/voightkampff/.pinned-ailang/ailang` the same command is **rc=0**.
This is a deliberate anti-false-green guard, not a base red — `ci.yml:124` says so in as many words
(*"Without AILANG_BIN, pinnedBinary(t) t.Skip()s and CI is false-green"*). **Every boundary-gate
invocation in this plan carries the export.** An executor reporting this gate red without checking
the export has mis-measured, and the 13 test names above are the signature.

### P2 — the doc's baseline B5 is refuted *(record only, no AC affected)*

Doc B5 claims `/usr/bin/grep -rn 'queue_census' --exclude-dir=.git .` → `rc=1`, zero hits. It is now
**rc=0 with 50 hits**, all inside `design_docs/planned/w-queue-closed-count-census.md` — the design
doc itself, present-untracked in the worktree, written into the tree *after* B5 was measured. The
positive control in the same call fires, so the empty result was never a fact about the instrument.

This is the doc's own D3 **"a control you record is a control you spend"** trap, realised one level
up: at the *file* level rather than the *row* level. It is also the reason D3's controls are a row
**number** and a **class** rather than a literal — that design decision is vindicated by its own
document breaking a coarser baseline within hours of being written.

**No acceptance consequence.** Every AC needle is scoped either to `ci.yml` (B-h: still rc=1) or to a
pre-`## Queue` region of `world-mission.md` (B-d: still 0). The executor must **never** use a
repo-wide `queue_census` grep as a red/green signal at any point in this sprint.

### P3 — the generated fixture is not valid UTF-8 *(record only; do not "fix" it)*

`iconv -f UTF-8 -t UTF-8 <fixture> >/dev/null` → **rc=1**; `file` → `Non-ISO extended-ASCII text`;
a per-line check counts **2 invalid lines of 91**. Cause: BWK awk 20200816's `substr($0,1,120)` is
**byte**-based and the charter's row lines contain `·` (U+00B7, 2 bytes), so two truncations land
mid-character. No NUL bytes, so git treats the file as text and diffs it normally.

**Measured harmless.** Running the classifier prototype against the produced fixture reproduces the
live reading exactly — 89 rows, `closed=41 open=4 untagged=44`, the same 44-number untagged list,
numbers exactly `1..89` with zero duplicates, both controls correct. The classifier reads only the
leading token, which is pure ASCII in every row.

The doc mandates **"checked in as produced"** and this plan keeps that: **do not hand-repair the two
lines.** The sha256 is pinned so any drift is detectable. GitHub may render those two lines with
U+FFFD; that is cosmetic.

### P4 — §3 and V18 name different generating regexes *(record only)*

§3: `/^[0-9]+\. /`. V18: `/^[0-9]+\. (\*\*|\[)/`. I ran both and diffed them: **IDENTICAL**, 91
lines / 10038 bytes each. Supporting counts: `awk 'NR>1610 && /^[0-9]+\. (\*\*|\[)/' … | wc -l` → 89
and `awk 'NR>1610 && /^[0-9]+\. /' … | wc -l` → 89 — no body line matches the broad form at this
base. The plan prescribes **§3's command verbatim** (§3 is the normative file list). Recorded so a
future regeneration against a changed charter knows the two are only *coincidentally* equal.

### P5 — D1 does not say which `^## Queue` match wins *(act on this)*

D1 says the scan starts "after the line matching `^## Queue`". V9's prototype takes the **first**
(`!h`). I measured `/usr/bin/grep -n '^## Queue' design_docs/world-mission.md` → exactly one hit,
line 1610, so the ambiguity is unreachable today. **The shipped script must take the FIRST match**,
matching the prototype the number was derived from, and say so in its header comment. AC-M2-3's
second half (`/usr/bin/grep -c '^## Queue' …` → 1) is the standing guard that keeps the choice
unobservable.

### P6 — §7's "`git status --short` → empty" is stale *(record only)*

It is `?? design_docs/planned/w-queue-closed-count-census.md`. Expected — the doc is written into
the worktree after §7 was measured. Noted so the executor is not surprised and does not "clean" it.

---

## Sequencing

### M1 strictly before M2 — not a preference

M2's `Queue census (live charter, controlled)` CI step invokes `bash scripts/queue_census.sh`. That
file does not exist until M1 lands. Landing M2 first turns `go-verify` red on **rc=127** for every
commit until M1 arrives. M2's other two edits (the heading `[ROUTED]`, the Repo Profile bullets) are
docs-only and could technically precede M1, but splitting them from the CI step buys nothing and
costs a commit.

### Order of work inside M1 — do not reorder

1. Record the three boundary gates **green before touching anything** (they are, per B-a — confirm,
   don't assume).
2. Write `scripts/queue_census.sh`; `chmod +x`; fence with `bash -n scripts/queue_census.sh`.
3. **Generate** the fixture with the §3 command; verify 91 lines / 10038 bytes / the pinned sha256.
4. Write `scripts/test_queue_census.sh` (19 arms); `chmod +x`; `bash -n` fence.
5. Run every `AC-M1-*` command, then the full suite.
6. Run the M1 mutation rows under the backup/restore discipline below.
7. Append the suite CI step at EOF of `.github/workflows/ci.yml`.
8. Re-run the three boundary gates. Commit.

### Order of work inside M2

1. Append the live-charter CI step at EOF (after M1's step).
2. Edit charter **line 1610 only**: `[ROUTED]` goes *inside* the existing parenthesis, before `)`.
3. Add **three** Repo Profile bullets (charter lines 36–806).
4. Run every `AC-M2-*` command.
5. Run the M2 mutation rows.
6. Re-run the three boundary gates. Commit.

---

## The M2 live-charter CI step is a NEW way for CI to go red on a docs-only commit

From M2 on, `go-verify` runs

```bash
bash scripts/queue_census.sh --doc design_docs/world-mission.md --control-closed 1 --control-open 79
```

on every push. A charter edit that breaks queue **structure** now reddens CI. The doc (§6) argues
this is precedented — `verify_go.sh` already runs the decision-ledger `--check` in CI — but a
precedent is not a licence to leave the recovery unwritten. Here it is.

**Control rows, confirmed first-party by the planner (not taken from the doc):**

| flag | row | charter line | first line, as measured | class |
|---|---|---|---|---|
| `--control-closed 1` | 1 | 1637 | `1. [LANDED 2026-07-24] **w-log-epoch-decision** · clause-1 · …` | **closed** (LANDED) |
| `--control-open 79` | 79 | 5711 | `79. **[PARKED — DESIGN REVIEW AND USER APPROVAL] w-evidence-applicability** · …` | **open** (PARKED) |

Row 1's decoration is bare `[`; row 79's is `**[` — between them the two controls exercise both
decoration shapes. Row 1 is a landed 2026-07-24 decision row that will never reopen; row 79 is
parked on an **attended** decision, so it cannot close without a human in the loop — precisely when
re-pointing the flag is affordable. There is also a *second* `1. ` line at charter line 866
(`1. ~~Iteration 0: ratify the bar~~`), **before** the heading and therefore outside the scan. That
is exactly why D1 anchors at the heading (V14), and I confirmed the scan ignores it.

**What reddens it, and what the recovery is:**

| Trigger | Signal | Recovery |
|---|---|---|
| **Duplicate row number** (F3) — a real duplicate, or a column-0 numbered line inside a row *body* colliding with a row number | `✗ duplicate row number: <n> (lines <a>, <b>)`, rc=2 | Real duplicate → renumber in the **same commit**. Body-line collision → **indent the offending body line by 2 spaces** so it is no longer column-0. **Never widen the classifier to tolerate it** — that is the road back to prose guessing (§9). |
| **Numbering gap** (F4) — a row deleted or renumbered without closing the range | `✗ numbering gap: expected 1..<max>, missing <list>`, rc=2 | Renumber to close the gap in the same commit. If a removal was genuinely intended, the queue's numbering contract is what changed — a controller decision, not a CI workaround. |
| **A control row changes class** (F6) — realistically, row 79 gets its attended decision and closes | `✗ CONTROL MISMATCH: row 79 expect=open got=LANDED`, rc=2, **with the full table still printed above it** | The **same commit** that closes row 79 re-points `--control-open` in `ci.yml` to another open tagged row. Today's open tagged rows, measured: **6, 79, 80** (PARKED) and **8** (IN-SPRINT). The mismatch is the control working, not a defect — and this recovery *is* Repo Profile bullet 3 (AC-M2-4), so the rule ships with the trap. |
| **Mis-generated heading** — the `[ROUTED]` append breaks the `^## Queue` prefix | `✗ no "## Queue" heading in design_docs/world-mission.md`, rc=2 | Restore line 1610 to a single line whose prefix is exactly `## Queue `. |
| **Mis-generated heading, other half** — `[ROUTED]` written on a *new* `## Queue…` line | `/usr/bin/grep -c '^## Queue'` becomes 2; AC-M2-3's second half fails and the instrument silently reads the **first** heading (P5) | Delete the new line, edit 1610 in place. AC-M2-3 is a two-part criterion for exactly this reason: the `[ROUTED]` count goes 0→1 **and** the `^## Queue` count stays at 1. |
| **A STATUS stamp introduces a column-0 `## Queue` line before 1610** | Under the first-match rule the census would enumerate the whole STATUS block | AC-M2-3's `^## Queue` count == 1 catches it. If a stamp must quote the heading, indent or fence it. |

**What M2 does NOT claim:** that the count "moved by exactly one". That is a delta between two
commits and no CI step sees two trees. Deleted at quorum r1 — do not re-introduce it as an AC, a
comment, or a log line.

---

## The snapshot fixture is generated, never typed

**Path:** `scripts/testdata/queue_census_rowstarts_8b15153.md`

**Generating command — §3's, verbatim, run from the repo root inside `bash -c`:**

```bash
awk '/^## Queue/{h=NR; print substr($0,1,120); next} h && (/^[0-9]+\. / || /^\*\*\[/) {print substr($0,1,120)}' \
  design_docs/world-mission.md > scripts/testdata/queue_census_rowstarts_8b15153.md
```

**Check in the bytes as produced.** Do not reformat, do not re-wrap, do not repair the two
invalid-UTF-8 lines (P3). Post-generation assertions, all measured by the planner:

```text
wc -lc  → "      91   10038"                                                    (matches V18)
shasum -a 256 → 42fe74686d105871046c84c6c081780ca117d1715f9dd77f4a008ac717fd54c1
line 1 = the queue heading; line 2 = the "**[LANDED 2026-07-24 (iter-13) …] w-m1-ailang-hardening"
preamble; line 3 = row 1; line 89 = row 79
its own census = 41 closed / 89 rows (tagged-open 4, untagged 44); numbers exactly 1..89, zero dups
```

**Expected reading: `41 closed / 89 rows (tagged-open 4, untagged 44)`.**

**Why generated:** it is the second of two independent implementations agreeing on one number. A
hand-typed fixture would be a transcription of the answer, and the `snapshot` arm would then be
testing typing rather than classification. **Never regenerate it** as part of routine work — it is
frozen at `8b15153` by name, and the §5 `snapshot` mutation row exists precisely to catch a fixture
regenerated from a different commit.

---

## M1 — the instrument, its suite, the frozen snapshot, and the suite's CI step (~0.2d)

### Files

| Path | Action | Mode | ~LOC | Note |
|---|---|---|---|---|
| `scripts/queue_census.sh` | **create** | 0755 | ~190 | D1 classifier (**first** `^## Queue` match, P5; scan to EOF, never stopping at a `## ` line), D3 mandatory controls, D4 ordered floors F0–F6, D5 report with `charter-blob`. Flags: `--doc` (default `design_docs/world-mission.md`, relative to cwd), `--control-closed`, `--control-open`, `--git-bin` (default `git`), `--help`. **Does NOT `cd` to the repo root** (D5; the precedent `gate1_range_check.sh` states the same rule in its header). Enumeration/classification/dup/gap live in **ONE** awk program; the bash wrapper owns arg parsing, F0, F5, F6 and exit codes. |
| `scripts/test_queue_census.sh` | **create** | 0755 | ~300 | 19 arms, `--only <arm>` with a loud refusal on an unknown arm, `ok`/`not ok`, final `N passed, M failed`. Copies the precedent's shape: `#!/usr/bin/env bash`, `set -uo pipefail`, `cd "$(dirname "$0")/.."`, `pass=0; fail=0`, `ok()`/`notok()`, `run_arm()`, `--only` dispatch. Every instrument invocation goes through an explicit `bash "$SCRIPT_UT" …`. |
| `scripts/testdata/queue_census_rowstarts_8b15153.md` | **create** | 0644 | 91 lines | Generated — see above. |
| `.github/workflows/ci.yml` | **modify** | — | +3 | Append at **EOF** (B-j: the file is 204 lines and the Gate-1 step is its last content), in `go-verify`. |

CI step shape, matching the precedent exactly:

```yaml
      - name: Queue census instrument suite
        timeout-minutes: 2
        run: bash scripts/test_queue_census.sh
```

### The one thing the suite must not do

**No helper may synthesise a heading, a row number, or a tag.** The only fixture builder is
`make_doc <path>` reading literal lines from a heredoc; each floor arm writes its broken shape by
hand. This is D4's defence against the iter-108 failure mode (V15): the harness helper that builds a
gate's input emits the very field the floor checks for, so the floor's input shape becomes
unproducible and the floor is never actually reached. Every floor arm asserts its **own specific
message**, so if a different floor fires first the arm reds on message mismatch — which is the proof
it reached the floor it names.

### Acceptance criteria — quoted from §4 by AC id

The 19 arms: `enumerate`, `closed-vocab`, `open-vocab`, `untagged`, `prose-not-tag`, `body-heading`,
`preamble`, `no-doc`, `no-heading`, `zero-rows`, `dup-number`, `gap`, `controls-required`,
`control-mismatch`, `report`, `blob-oid`, `blob-unavailable`, `heading-tail`, `snapshot`.

| AC | Command | Expected | Base (planner) |
|---|---|---|---|
| **AC-M1-1** — *"both scripts exist and are executable"* | `test -x scripts/queue_census.sh && test -x scripts/test_queue_census.sh && echo EXISTS` | `EXISTS` | `rc=1` (B-b) |
| **AC-M1-2** — *"enumeration is order-independent and line-accurate"*; arm writes rows `3`,`1`,`2` in that file order, asserts `row=1`,`row=2`,`row=3` in that order with correct `line=` | `bash scripts/test_queue_census.sh --only enumerate; echo "rc=$?"` | `ok …` lines, `rc=0`; controls **1/2**, `2 closed / 3 rows (tagged-open 1, untagged 0)` | `rc=127` (B-h) |
| **AC-M1-3** — *"the closed vocabulary and every decoration shape classify closed"*; rows 1–6 as `[LANDED …`, `[**LANDED …`, `[**[LANDED …`, `**[LANDED …`, `**[ROUTED …`, `**[RULED OUT …`, **plus row 7 `[PARKED] control`** | `bash scripts/test_queue_census.sh --only closed-vocab; echo "rc=$?"` | `rc=0`; controls **1/7**, `6 closed / 7 rows (tagged-open 1, untagged 0)` | `rc=127` |
| **AC-M1-4** — *"the open vocabulary classifies open, not closed"*; `[PARKED …`, `**[NEXT] …`, `[**[IN-SPRINT] …` as rows 1–3, **plus row 4 `[LANDED] control`** | `bash scripts/test_queue_census.sh --only open-vocab; echo "rc=$?"` | `rc=0`; controls **4/1** (not 1/4 — getting these backwards is an F6 mismatch), `1 closed / 4 rows (tagged-open 3, untagged 0)` | `rc=127` |
| **AC-M1-5** — *"an untagged row is UNTAGGED, counted in rows, not in closed"*; LANDED (1), PARKED (2), `**w-slug** · clause-2 · prose` (3) | `bash scripts/test_queue_census.sh --only untagged; echo "rc=$?"` | `rc=0`; row 3 prints `tag=- class=UNTAGGED`; controls **1/2**, `1 closed / 3 rows (tagged-open 1, untagged 1)` | `rc=127` |
| **AC-M1-6** — ***"prose is never a tag (the row-57 defect)"***; an untagged row whose body carries `[LANDED] [ROUTED] [RULED OUT] [PARKED] [NEXT]`, `it landed correctly`, `LANDED 2026-09-08 (iter-1)`; a slug-first row `**w-slug** · **[LANDED 2026-09-08]** …`; plus one `[LANDED] control` and one `[PARKED] control`. Numbering: 1 = LANDED control, 2 = PARKED control, 3 and 4 = the targets | `bash scripts/test_queue_census.sh --only prose-not-tag; echo "rc=$?"` | `rc=0`; **both targets `class=UNTAGGED`**; controls **1/2**, `1 closed / 4 rows (tagged-open 1, untagged 2)` | `rc=127` |
| **AC-M1-7** — *"a `## ` line inside a row body does not end the enumeration"* (attempt 2's bug); row 1, body line `## Premise Verification Log`, rows 2 and 3 | `bash scripts/test_queue_census.sh --only body-heading; echo "rc=$?"` | `rc=0`; controls **1/2**, `2 closed / 3 rows (tagged-open 1, untagged 0)` | `rc=127` |
| **AC-M1-8** — *"the unnumbered preamble block is reported and excluded"*; **two** fixtures | `bash scripts/test_queue_census.sh --only preamble; echo "rc=$?"` | `rc=0`; fixture A → `preamble: 1 unnumbered closed block (line 2) excluded from rows` + `1 closed / 2 rows (tagged-open 1, untagged 0)`; fixture B (no preamble, same rows, same controls **1/2**) → `preamble: 0 …` | `rc=127` |
| **AC-M1-9** — *"floors F0–F4, one arm each, each asserting its own message and `rc=2`"* | `bash scripts/test_queue_census.sh --only no-doc; … --only no-heading; … --only zero-rows; … --only dup-number; … --only gap; echo "rc=$?"` | five `ok …` blocks, final `rc=0`. Shapes: `no-doc` → `--doc /nonexistent/x.md`; `no-heading` → rows, no `## Queue`; `zero-rows` → heading + prose only; `dup-number` → `1, 2, 2`; `gap` → `1, 3` | `rc=127` |
| **AC-M1-10** — *"controls are mandatory (F5) and adjudicated (F6)"* | `bash scripts/test_queue_census.sh --only controls-required; … --only control-mismatch; echo "rc=$?"` | `rc=0`. `controls-required`: valid doc, no flags → `controls required`, rc=2. `control-mismatch`: 3-row fixture (1 LANDED, 2 PARKED, 3 untagged) run **twice with both flags** — `--control-closed 2 --control-open 2`, then `--control-closed 1 --control-open 3` → each `CONTROL MISMATCH: row=<n> expect=… got=…`, rc=2, **table still printed above it** | `rc=127` |
| **AC-M1-11** — *"the report format (D5)"*; 4-row fixture (2 closed incl. one ROUTED, 1 open, 1 untagged), controls **1/3** | `bash scripts/test_queue_census.sh --only report; echo "rc=$?"` | `rc=0`; last line exactly `census: 2 closed / 4 rows (tagged-open 1, untagged 1) [LANDED 1, ROUTED 1, RULED OUT 0; PARKED 1, IN-SPRINT 0, NEXT 0] [queue_census.sh, charter-blob <40 hex>, controls 1/3]`; both `control: … ok` present; output contains **neither** `carried` **nor** `of 4 rows closed` | `rc=127` |
| **AC-M1-15** — *"the printed blob OID is the input file's blob, and its absence is loud"* | `bash scripts/test_queue_census.sh --only blob-oid; … --only blob-unavailable; echo "rc=$?"` | `rc=0`. `blob-oid`: 2-row fixture, controls **1/2**, printed `charter-blob` == `git hash-object <that fixture>`. `blob-unavailable`: same fixture and controls, `--git-bin /nonexistent/git` → total line carries the literal `charter-blob unavailable`, a stderr line names it, **rc=0** | `rc=127` |
| **AC-M1-12** — *"the heading anchor is the `^## Queue` prefix"*; heading is the **post-M2** text `## Queue (top = next; tags: [NEXT] [IN-SPRINT] [PARKED] [LANDED] [RULED OUT] [ROUTED])` | `bash scripts/test_queue_census.sh --only heading-tail; echo "rc=$?"` | `rc=0`; `heading_line=` reported; controls **1/2**, `1 closed / 2 rows (tagged-open 1, untagged 0)` | `rc=127` |
| **AC-M1-13** — *"the frozen charter snapshot reproduces the prototype"*; controls **1/79** | `bash scripts/test_queue_census.sh --only snapshot; echo "rc=$?"` | `rc=0`; `census: 41 closed / 89 rows (tagged-open 4, untagged 44)` prefix, `preamble: 1 …`, both controls ok, and `row=16`, `row=18`, `row=66`, `row=68`, `row=70` each `class=UNTAGGED` | fixture absent, `rc=1` (B-h). **Planner pre-checked this arm is satisfiable** — see B-i |
| **AC-M1-14** — *"the whole suite is green and wired into CI"* | `bash scripts/test_queue_census.sh; echo "rc=$?"` then `/usr/bin/grep -n 'test_queue_census' .github/workflows/ci.yml` | `19 passed, 0 failed`, `rc=0`; one `run: bash scripts/test_queue_census.sh` line inside `go-verify` | grep `rc=1` (B-h) |

Two constraints worth calling out inside these ACs:

- **AC-M1-9's five failure arms pass no control flags and must not.** D4's floors are **ordered** and
  F0–F4 all precede F5. Each arm asserting its **own** floor's message is what proves it reached that
  floor rather than falling through to `controls required`. Asserting only `rc=2` would make all five
  vacuously equivalent. (This is the controller's r2 clarification, labelled as such in §4 — the
  gpt5-6-sol objection about mandatory controls is scoped to *success* arms.)
- **AC-M1-11's `[0-9a-f]{40}`** is a `{n,m}` interval. It is allowed **there** because it is a
  `/usr/bin/grep -E` pattern in the *test* script. It must **not** appear in the awk program — D5
  bans intervals for mawk/BWK portability.

### Mutation matrix (from §5)

| Arm | Mutation applied to `scripts/queue_census.sh` unless noted | Expected red set |
|---|---|---|
| `enumerate` | sort dropped (print in file order) / `line=` printed as row index | `enumerate` only |
| `closed-vocab` | `RULED OUT` removed from the vocabulary (two-word tag); decoration strip limited to one char | `closed-vocab`, `snapshot` |
| `open-vocab` | `PARKED` moved to the closed set | `open-vocab`, `snapshot`, `report` |
| `untagged` | UNTAGGED defaulted to open (counted in tagged-open) | `untagged`, `report`, `snapshot` |
| `prose-not-tag` | classifier searches the whole first line (`~ /\[LANDED/`) or the body; slug-first accepted as a second position | `prose-not-tag` (+ `snapshot`: rows 66/68/70 flip) |
| `body-heading` | enumeration stops at the next `^## ` | `body-heading` only |
| `preamble` | preamble counted as a row / `preamble:` line dropped | `preamble`, `snapshot` |
| `no-doc` | F0 check removed (awk reads nothing, prints `0 rows`) | `no-doc` (and `zero-rows` if F2 also gone) |
| `no-heading` | scan starts at line 1 instead of the heading | `no-heading`, and `snapshot`/`enumerate` if pre-heading numbered lines exist in the fixture |
| `zero-rows` | `[ rows -gt 0 ]` floor removed → `census: 0 closed / 0 rows` exit 0 | `zero-rows` only |
| `dup-number` | duplicate detection removed (last wins) | `dup-number` only |
| `gap` | contiguity check removed | `gap` only |
| `controls-required` | controls made optional | `controls-required` only |
| `control-mismatch` | mismatch downgraded to a warning (exit 0) / table suppressed on mismatch | `control-mismatch` only |
| `report` | total line reworded (`N of M rows closed`) or per-tag breakdown dropped | `report` only |
| `heading-tail` | heading matched as the full literal string | `heading-tail` only |
| `snapshot` | fixture regenerated from a different commit / any classifier drift | `snapshot` only |
| `blob-oid` | OID computed on the wrong path (e.g. `$0`), or hashed via `sha256sum` instead of `git hash-object` | `blob-oid` only |
| `blob-unavailable` | missing git swallowed (field omitted, no token) or turned into exit 2 | `blob-unavailable` only |
| AC-M1-14 grep | **`.github/workflows/ci.yml`**: CI suite step removed | AC-M1-14 |

**Floor reachability is itself asserted**: each floor arm's `ok` requires the floor's own message, so
a harness change that makes a floor unreachable reds `no-heading` immediately.

### M1 boundary gate — which of the three profile gates must be green at that commit

**All three.** They are green at base (B-a), so these are non-regression gates, not gates to turn
green.

```bash
export PATH=/opt/homebrew/bin:$PATH
export AILANG_BIN=/Users/voightkampff/.pinned-ailang/ailang
bash ./scripts/verify_ail.sh          # rc=0, "✓ verify gate PASSED: 11 required identities verified, 40 named tests pass"
go vet ./...                          # rc=0
go test ./... -count=1                # rc=0  (~2 min; RED without AILANG_BIN — Finding P1)
```

Plus, at the M1 commit:

- `bash scripts/test_gate1_range_check.sh` → `rc=0`, `18 passed, 0 failed` (the precedent suite
  shares `scripts/testdata/`; it must not regress)
- `bash -n scripts/queue_census.sh` and `bash -n scripts/test_queue_census.sh` → `rc=0` each
- `git diff --quiet 8b15153 -- tools/launchd` → `rc=0` (frozen core untouched)
- `git diff --quiet 8b15153 -- design_docs/world-mission.md` → `rc=0` — **M1 touches no charter
  byte**; the charter blob must still be `2dda66a36a6bc60ab280a73b3f0b3c6be89ed399`

The `11 identities / 40 named tests` numbers must not move — M1 adds no `.ail`. A change in either is
a red flag, not a success.

---

## M2 — the live-charter CI step, the heading vocabulary, and the Repo Profile rule (~0.1d)

### Files

**`.github/workflows/ci.yml`** — append a second step at EOF, after M1's:

```yaml
      - name: Queue census (live charter, controlled)
        timeout-minutes: 2
        run: bash scripts/queue_census.sh --doc design_docs/world-mission.md --control-closed 1 --control-open 79
```

The `run:` line must contain the AC-M2-2 needle **verbatim**.

**`design_docs/world-mission.md`** — two surgical edits, nothing else.

*(a) Line 1610 only.* Current bytes, measured:

```
## Queue (top = next; tags: [NEXT] [IN-SPRINT] [PARKED] [LANDED] [RULED OUT])
```

becomes

```
## Queue (top = next; tags: [NEXT] [IN-SPRINT] [PARKED] [LANDED] [RULED OUT] [ROUTED])
```

`[ROUTED]` goes **inside** the existing parenthesis. Do not add a line. Do not reflow.

*(b) Three bullets in the Repo Profile section.* Measured bounds: `## Repo Profile (…)` at line
**36**; the next top-level heading `## Human Decision Ledger (authoritative current state)` at line
**807**. The bullets go inside 36–806 — 800+ lines before the queue heading, so AC-M2-4's
`/^## Queue/{exit}` region trivially contains them.

1. **D1 tag position** — a queue row's closure/parking tag is the FIRST token after `^<n>. `,
   tolerant only of a leading run of `[` and `*`. A tag after the slug is UNTAGGED and the census
   reports it by number rather than guessing. Closed: `LANDED`, `ROUTED`, `RULED OUT`. Open (tagged):
   `PARKED`, `NEXT`, `IN-SPRINT`. Read by `scripts/queue_census.sh`.
2. **D6 citation form** — a closed-row count is published in a STATUS stamp or a mission-log
   `**Progress:**` line ONLY as the script's own total line, byte-for-byte:
   `census: N closed / M rows (tagged-open T, untagged U) [LANDED X, ROUTED Y, RULED OUT Z; PARKED P, IN-SPRINT Q, NEXT R] [queue_census.sh, charter-blob <oid>, controls a/b]`.
   Never carried. A line reading `charter-blob unavailable` is not publishable. No fresh reading →
   **no number**.
3. **D3 control re-point** — the live CI step passes `--control-closed 1 --control-open 79`; when a
   control row changes class the census exits 2 `CONTROL MISMATCH` and the SAME commit re-points the
   flag in `.github/workflows/ci.yml`.

Bullets 2 and 3 must each contain the literal `queue_census.sh` (AC-M2-4 needs ≥3 occurrences —
one per bullet) and bullet 2 must contain the literal `charter-blob` (AC-M2-6).

**Must not touch:** any queue row body (D2), **row 72's own line at 5573** (the *controller* tags it
at Gate-4), the STATUS block (`## STATUS (rotation rule)` at 856, stamps from 858), the Human
Decision Ledger block (807–855, including the `decision-ledger:start/end` markers at 827/852 that
`scripts/mission_decisions.sh --check` parses), or `tools/launchd/*`.

### Acceptance criteria — quoted from §4 by AC id

| AC | Command | Expected | Base (planner) |
|---|---|---|---|
| **AC-M2-1** — *"the live charter census runs clean with the published controls"* | `bash scripts/queue_census.sh --doc design_docs/world-mission.md --control-closed 1 --control-open 79; echo "rc=$?"` | a `row=` line per row, `preamble: 1 …`, two `control: … ok`, a `census: N closed / M rows …` last line, `rc=0` | `rc=127` (B-h) |
| **AC-M2-1** cross-check | `test "$(bash scripts/queue_census.sh --doc design_docs/world-mission.md --control-closed 1 --control-open 79 \| /usr/bin/grep -c '^row=')" = "$(awk '/^## Queue/{h=1;next} h && /^[0-9]+\. /' design_docs/world-mission.md \| wc -l \| tr -d ' ')" && echo AGREE` | `AGREE` | — |
| **AC-M2-2** — *"the live step is in CI"* | `/usr/bin/grep -n 'queue_census.sh --doc design_docs/world-mission.md --control-closed 1 --control-open 79' .github/workflows/ci.yml` | one hit in `go-verify` | `rc=1` (B-h) |
| **AC-M2-3** — *"the heading vocabulary names `[ROUTED]` and the anchor still holds"* | `/usr/bin/grep -c '^## Queue .*\[ROUTED\]' design_docs/world-mission.md` **then** `/usr/bin/grep -c '^## Queue' design_docs/world-mission.md` | `1`; then `1` (unchanged) | `0`; `1` (B-d) |
| **AC-M2-4** — *"the Repo Profile carries the rule, before the queue heading"* | `awk '/^## Queue/{exit} /queue_census\.sh/{c++} END{print c+0}' design_docs/world-mission.md` | `3` or more | `0` (B-d; control `verify_ail\.sh` → 9) |
| **AC-M2-5** — *"the live reading's blob OID is the charter's blob"* | `test "$(bash scripts/queue_census.sh --doc design_docs/world-mission.md --control-closed 1 --control-open 79 \| /usr/bin/grep '^census: ' \| sed -E 's/.*charter-blob ([0-9a-f]{40}).*/\1/')" = "$(git hash-object design_docs/world-mission.md)" && echo BLOB-MATCH` | `BLOB-MATCH` | `rc=127` (B-h). Base blob `2dda66a3…` confirmed by two routes (B-e) |
| **AC-M2-6** — *"the D6 citation protocol is pinned at both ends"* | `awk '/^## Queue/{exit} /charter-blob/{c++} END{print c+0}' design_docs/world-mission.md` | `1` or more; **and** the `report` arm (AC-M1-11) holds, so the printed line and the quoted rule share the literal `[queue_census.sh, charter-blob …, controls a/b]` shape | `0` (B-d; control `pinned-ailang` → 5) |

**AC-M2-5 note.** After M2's edits the charter's blob OID **changes** — `2dda66a3…` was the *base*
charter. The AC compares the printed value to a live `git hash-object`, not to a constant, so it
holds. **Do not hard-code `2dda66a3…` anywhere.**

**Expected reading at the M2 commit:** the count is **unchanged** at `41 closed / 89 rows
(tagged-open 4, untagged 44)` — M2 edits only the heading and the Repo Profile, both outside the row
set.

**AC-M2-6's uncoverable half.** A controller *publishing* a reading in a log entry is a procedure
step, verified at landing by re-running AC-M2-5 against the tag commit's charter. Not an executor
deliverable and not a CI assertion.

### Mutation matrix (from §5)

| Arm | Mutation | Subject | Expected red |
|---|---|---|---|
| AC-M2-5 | live step's `--doc` pointed at the snapshot fixture (table still plausible, blob wrong) | `.github/workflows/ci.yml` | AC-M2-5 |
| AC-M2-6 | citation bullet written after the heading, or quoting a `@ <sha>` form | `design_docs/world-mission.md` | AC-M2-6 |
| AC-M2-1/2 | live step removed, or control row re-pointed to an UNTAGGED row (e.g. `--control-open 72` → `✗ CONTROL MISMATCH: row 72 expect=open got=UNTAGGED`, rc=2) | `.github/workflows/ci.yml` | AC-M2-2 / AC-M2-1 |
| AC-M2-3 | `[ROUTED]` not appended, or appended to a *new* `## Queue…` line (count becomes 2) | `design_docs/world-mission.md` | AC-M2-3 |
| AC-M2-4 | rule written after the heading (inside a row) | `design_docs/world-mission.md` | AC-M2-4 |

The charter is 5722 lines / 904369 bytes of **ratified mission state**. Back it up and sha256-verify
the restore before and after every charter mutation.

### M2 boundary gate

**All three profile gates green**, same commands as M1's. Plus, at the M2 commit:

- `bash scripts/test_queue_census.sh` → `rc=0`, `19 passed, 0 failed` — in particular the
  `heading-tail` arm now describes the **live** heading
- `bash scripts/test_gate1_range_check.sh` → `rc=0`, `18 passed, 0 failed`
- `bash scripts/mission_decisions.sh --check --file design_docs/world-mission.md` → `rc=0` (the
  ledger block at 827–852 is untouched; this is the existing charter reader, §6/V20)
- `git diff 8b15153 -- design_docs/world-mission.md | /usr/bin/grep -E '^@@'` — **every hunk header
  must start inside 36–806 or at 1610.** Any hunk touching 807–1609 or 1611+ is out of scope and
  must be reverted from the backup.

---

## Mutation restore discipline — the row-82 rule

**Restore from a `cp`-captured backup and assert its sha256. NEVER `git checkout --`.**

`git checkout -- <path>` restores the file to **HEAD**. Every subject under mutation in this sprint
— `queue_census.sh`, `test_queue_census.sh`, `ci.yml`'s new steps, the fixture, the charter's new
bullets — is **uncommitted work** at the moment it is mutated. `git checkout --` would delete the
milestone, not undo the mutation. Charter row 82 has measured this class twice, including a
sub-agent reverting the *controller's* staged git state in a shared worktree.

```bash
export PATH=/opt/homebrew/bin:$PATH
BAK=/Users/voightkampff/dev/sunholo-data/.wt-world-iter172-mutbak   # OUTSIDE the repo, NOT under /tmp
mkdir -p "$BAK"

cp scripts/queue_census.sh "$BAK/queue_census.sh.bak"
shasum -a 256 "$BAK/queue_census.sh.bak"          # RECORD this digest

# … apply the mutation with a targeted sed/edit …
bash scripts/test_queue_census.sh > /dev/null 2>&1; echo "rc=$?"   # capture the failing arm names

cp "$BAK/queue_census.sh.bak" scripts/queue_census.sh
shasum -a 256 scripts/queue_census.sh             # MUST equal the recorded digest — if not, STOP
bash scripts/test_queue_census.sh > /dev/null 2>&1; echo "rc=$?"   # MUST be 0 before the next mutation
```

The backup directory is a **sibling of the worktree, outside the repo** (so a backup can never
pollute the diff or be committed) and **not under `/tmp`** (row 71's ratified rule: critical state
never lives where the OS wipes).

**Expected red set = exactly the listed arm(s), every other arm green.** A mutation that reds *more*
arms than listed means an arm is coupled to the wrong behaviour; *fewer* means an arm is vacuous.
Either way it is a finding to record, not a number to adjust.

**No git write of any kind during mutation work** — no `add`, `commit`, `stash`, `checkout`, `reset`,
`mv`, `branch`, `worktree`.

---

## Can M2's charter edit conflict with the controller's Gate-4 record edits?

**Measured, not assumed.**

| Who | Region touched |
|---|---|
| **M2 (executor)** | line **1610** (heading) + one insertion inside the Repo Profile block, lines **36–806** |
| **Controller (Gate-4)** | the STATUS block — `## STATUS (rotation rule)` at **856**, stamps from **858** — and **row 72's line at 5573** (append `[LANDED 2026-09-08 (iter-172) …]` in the D1 leading position) |

**Disjoint, with a 49-line gap** between the Repo Profile's end (806) and the STATUS block's start
(856), and 3900+ lines to row 72. Git's 3-way merge needs 3 lines of context and has 49. **There is
no textual conflict.**

**The real risk is not a merge conflict — it is a shared-worktree state clobber** (charter row 82,
measured twice). The executor works **only** in `/Users/voightkampff/dev/sunholo-data/.wt-world-iter172`
on `sprint/w-queue-census`. The controller's Gate-4 edits happen in the **main checkout after the
branch merges**. The two never touch the same working tree at the same time.

**Staging discipline that makes it survivable:**

1. M2's charter edit is **one commit** whose diff hunks are provably confined to line 1610 and the
   Repo Profile block. Assert it with the `git diff … | grep -E '^@@'` check in the M2 boundary gate.
2. The executor **never** edits row 72, the STATUS block, or the ledger block. A Gate-4 stamp is the
   controller's commit, not this sprint's.
3. Do not rebase or amend the M2 commit after the controller has begun Gate-4.

**Downstream count change the controller must expect.** Row 72 is **UNTAGGED today** — measured:
`72. **w-queue-closed-count-is-an-increment-chain-no-instrument-reproduces** · clause-2 ·` at line
5573, slug-first. When the controller tags it `[LANDED …]` in the **leading** position the census
moves to `42 closed / 89 rows (tagged-open 4, untagged 43)` with `LANDED 41`, and the live CI step
stays green (controls 1/79 unaffected). If the controller instead tags it slug-first, **the count
does not move and the row stays UNTAGGED** — which is D1 working as designed and is a visible,
reportable outcome, not a CI failure. Repo Profile bullet 1 exists to make the leading position the
house habit.

---

## Rig constraints that bind every command in this plan

| # | Constraint | Verified at base |
|---|---|---|
| 1 | `grep` in this login shell is **ugrep** | Every grep in this plan and every grep the executor runs is `/usr/bin/grep`. |
| 2 | The login shell is **zsh**; `/bin/bash` is **3.2.57** | Every bash script is invoked as `bash <script>` or inside `bash -c '…'`. Running a bash script under zsh on this rig has manufactured `rc=2` against a positive control. |
| 3 | **`${PIPESTATUS[0]}` expands EMPTY in zsh** | Planner-measured this fire — zsh spells it `pipestatus`, 1-indexed. Never capture an rc through a pipe; redirect to a file and read `$?` inside `bash -c`. **An empty rc reading is not `rc=0`.** |
| 4 | `awk` is **BWK 20200816**; CI is `ubuntu-latest` (**mawk**) | Banned everywhere: `gensub`, `asort`, `{n,m}` intervals, `\b`, bash-4 associative arrays, `${x^^}`. `sort -n` does the ordering. `bash -n` is a syntax fence but **not** an awk fence — the only real awk fence is running the suite, which CI does on mawk. |
| 5 | `substr` in BWK awk is **byte**-based | Planner-measured — this is why the fixture carries 2 truncated multi-byte characters (P3). |
| 6 | There is **no `timeout(1)`** on this rig | Never write one into a command. `timeout-minutes:` in the CI YAML is a GitHub Actions key, not the binary. |
| 7 | `scripts/verify_go.sh` is **NOT a gate** | Charter row 76: its fleet driver-drift arm fatals early and suspends every Go assertion behind it. An AC naming it is unsatisfiable by any World change. It is never named as a gate or an AC here. |
| 8 | `tools/launchd/*` is **FROZEN CORE** (`D-WORLD-DRIVER-1`) | No milestone touches it; `git diff --quiet 8b15153 -- tools/launchd` is a boundary check at both milestones. |
| 9 | The shared `mission-control` SKILL.md is **not in this repo** | It lives in the V1 checkout. Nothing here edits it. §6/V16 confirms it reads the heading prefix-anchored (`grep -A2 "^## Queue"`), so M2's heading tail is safe upstream. |
| 10 | `export PATH=/opt/homebrew/bin:$PATH` | Required in every Bash call needing `go`/`gh`/`node`. |
| 11 | `export AILANG_BIN=/Users/voightkampff/.pinned-ailang/ailang` | Required — and load-bearing for `go test` (Finding P1). Measured present, mode 0755, 91826738 bytes. |

---

## Commit policy

Two commits on `sprint/w-queue-census`, in order:

1. `feat(scripts): queue census instrument — anchored leading-tag classifier, mandatory controls,
   ordered floors, blob-stamped report + 19-arm suite + frozen 8b15153 snapshot + CI step`
2. `docs(mission): row 72 — live-charter census CI step, [ROUTED] in the queue heading, three Repo
   Profile bullets (tag position, citation form, control re-point)`

The executor performs exactly these two `add`+`commit` operations. **No push, no PR, no rebase, no
stash, no checkout, no reset.** The controller builds every other commit. A pi-runner
`empty_worktree` verdict on a committed sprint is a **false negative** (charter row 84) — the
controller adjudicates by reading `git log`, never discards on that verdict alone.

---

## Velocity

~500 new lines (instrument ~190, suite ~300, fixture 91 generated lines) + 8 CI lines + ~4 charter
lines. Two commits, one iteration.

**Basis:** the precedent pair — `gate1_range_check.sh` (9661 bytes) + `test_gate1_range_check.sh`
(10790 bytes, 18 arms) — landed in one iteration (iter-170, row 70). This sprint is the same shape
with 19 arms and a *smaller* instrument surface: no `gh`, no network, no throwaway git repos, every
fixture a heredoc. **The planner does not dispute the doc's ~0.3d.**

---

## Start here

**M1.** Record the three boundary gates green (they are — re-confirm, do not assume), then write
`scripts/queue_census.sh` and fence it with `bash -n`.

**The one way to fail:** letting the classifier read anything but the first token after `^<n>. `.
Every widening — a second accepted position, a body search, a case-insensitive keyword — is attempt
2's over-count returning. `prose-not-tag` and `snapshot` are the two arms that exist to stop it.
