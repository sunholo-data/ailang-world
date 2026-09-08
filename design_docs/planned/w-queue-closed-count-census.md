# w-queue-closed-count-census — closure becomes a leading tag the instrument reads, and a census script prints the per-row table with its own controls

- Status: **planned** (design only; nothing implemented, no git write performed) · Date: **2026-09-08** (revised twice 2026-09-08 — quorum r1 BLOCKED N=2 no absentees, designer revision; quorum r2 BLOCKED N=2 no absentees, closed under the Gate-2 **narrow-refinement carve-out** with both reviewer fixes applied VERBATIM by the controller, see §10) · Designer: `claude:claude-fable-5-1` (iter-172) · Base commit measured at: **`8b151533ed701e6648f7037d86579383b04de407`** (worktree `.wt-world-iter172`, branch `sprint/w-queue-census`) · Owning queue row: **72 — `w-queue-closed-count-is-an-increment-chain-no-instrument-reproduces`** (clause-2, ~0.3d, gated on nothing, surfaced iter-151) · Scope class: **HARNESS — auditability of the mission's own progress metric** · Verify profile: `ailang-code` (`bash ./scripts/verify_ail.sh`, `go vet ./...`, `go test ./... -count=1`; `scripts/verify_go.sh` is NOT a gate, charter row 76).

---

## §1 Problem — the closed-row count was an increment chain; now it is nothing, and still nothing can produce one

Row 72's premise, verbatim: every STATUS stamp published an `N of M rows closed` headline that was
*carried* (previous iteration's number plus one) and never re-derived; iteration 151 tried to
measure it twice and got three disagreeing numbers (carried `58 of 70`, attempt 1 `36 of 72`,
attempt 2 `53 of 70`), with attempt 2 demonstrably over-counting because prose inside an open row
matched a closure keyword.

**What is measured now, first-party at `8b15153` (every command in §7):**

- The queue heading is line **1610** (V1). After it there are **89** numbered rows, numbers exactly
  `1..89`, **no duplicates, no gaps** (V2).
- Rows are **not in file order**: row 82 sits at line 5340, row 68 at line 5543 (V3). Any instrument
  that assumes ascending order is wrong by construction.
- There is **one unnumbered closed block** between the heading and row 1 — the
  `**[LANDED 2026-07-24 (iter-13) …] w-m1-ailang-hardening` preamble at line 1615 (V4). A census
  must say, explicitly, whether it is in the denominator.
- **Tag position is honoured by a minority.** A closure/parking keyword in brackets appears on
  **45** of 89 first lines; a bracket *immediately* after the row number on **30** (V5). So the
  fixed position the heading itself promises (`## Queue (top = next; tags: [NEXT] [IN-SPRINT]
  [PARKED] [LANDED] [RULED OUT])`) holds for a third of the queue.
- **No census instrument exists.** `scripts/` contains nothing that reads a closure tag; the
  positive control (the same needle over the charter) hits (V6).
- **The row's over-count mechanism reproduces.** Open row 57's body contains the word *landed*
  exactly once, lowercase, as an English verb — "`-title` and `-from` landed correctly in both"
  (V8). A body-wide, case-insensitive keyword search classes that row closed. That is the defect a
  prose-guessing classifier has, and it is what this design refuses to do.
- **Under the anchored rule this doc proposes (§2 D1), the live charter reads `41 closed / 89 rows
  (tagged-open 4, untagged 44)`** — prototyped in `awk` against the real file, per-class row lists
  recorded (V9). Five *closed* rows are among the 44 UNTAGGED because their tag is not in the
  leading position (16, 18, 66, 68, 70 — V10); the instrument reports them by number instead of
  guessing them closed.

**Honest scope note — half the row's premise has rotted, and the doc must say so.** The
`**Progress:** N of M queue rows closed` headline appears in the mission log at iterations 146–151
and **nowhere after 151**; entries 152+ publish `**Progress:** **goal unmoved.**` with no count
(V7). So the loop already complied with the row's interim instruction ("say the count is carried,
not measured") by dropping the number entirely. **The surviving defect is the other half: no
verified census exists and nothing on this rig can produce one.** This document designs that
instrument. It does not resurrect a carried number; it makes a *reading* available and says how a
reading, and only a reading, may be published (D6).

---

## §2 Design

Two files mirror the landed precedent pair `scripts/gate1_range_check.sh` +
`scripts/test_gate1_range_check.sh` (a small bash-3.2 instrument, a self-contained arm suite with
`--only <arm>`, one CI step in `go-verify`): **`scripts/queue_census.sh`** and
**`scripts/test_queue_census.sh`**, plus one generated fixture under `scripts/testdata/`.

### D1 — The classifier reads ONE fixed position and never guesses from prose

**Row start.** After the line matching `^## Queue`, every line matching `^[0-9]+\. ` is a queue-row
start; the integer is the row number; the line number is recorded. Nothing else is a row. The scan
runs to end-of-file — **it does not stop at a `## ` line**, because a row body contains one
(`## Premise Verification Log …` at line 5512 sits inside row 89, V11); that was attempt 2's
enumeration bug.

**Tag.** Take the row-start line, drop `^[0-9]+\. `, then drop the leading run of markdown
decoration characters — exactly `[` and `*`, nothing else, no whitespace. The tag must now be at
position 1 and be one of the closed vocabulary, followed by a non-letter or end of line:

| tag | class |
|---|---|
| `LANDED`, `ROUTED`, `RULED OUT` | **closed** |
| `PARKED`, `NEXT`, `IN-SPRINT` | **open** (tagged) |
| anything else | **UNTAGGED** — the instrument makes no claim |

Decoration tolerance is what lets `1. [LANDED …`, `3. [**LANDED …`, `4. [**[LANDED …` and
`41. **[LANDED …` all classify: the tag *word* is the first token of the row in every one of those
shapes. It is not a prose search: `65. **w-slug** · **[PARKED …` is UNTAGGED because the first
token is `w-slug`; `57. **[AMENDED …` is UNTAGGED because `AMENDED` is not in the vocabulary;
`16. [**[COMPLETE …` is UNTAGGED because `COMPLETE` is not in the vocabulary. **The instrument does
not read the row body at all**, so the tag literals listed in row 72's own body (line 5591, V12)
and the lowercase *landed* in row 57's body (V8) cannot match anything.

**Why one position, not two.** Four recent rows (65, 66, 68, 70) put the tag after the slug, and
accepting that as a second position would classify them today. It is rejected: every additional
accepted position is a regex somebody argues about, and two positions is the first step back
towards "wherever the keyword is". The row's own words are *a leading tag in a fixed position*.
A slug-first row that closes moves its tag to the front — a one-line edit the census then verifies.

**Why UNTAGGED is a bucket and not a default.** Open rows are filed without a tag by house habit,
so UNTAGGED is *consistent with* open — but 5 of today's 44 UNTAGGED rows are closed (V10). Defaulting
UNTAGGED to open would publish a false denominator split; defaulting it to closed is attempt 2. The
honest reading is `41 closed / 89 rows` with the 44 named. The per-row table is the worklist.

**ROUTED is closed.** Row 53 is `[ROUTED — WORLD-SIDE COMPLETE …; TRACKING sunholo-data/ailang#…]`
and row 72's item lists `[ROUTED]` as a closure tag. The total line prints per-tag counts so a
reader who wants LANDED-only can subtract. The heading's vocabulary list lacks `[ROUTED]`; M2 adds
it (safe because the heading anchor is the `^## Queue` prefix — D5 arm `heading-tail`; the shared
skill's own read is also prefix-anchored, `grep -A2 "^## Queue"`, V16).

### D2 — Normalising the 89 existing rows is OUT of scope; the instrument lands against the charter as-is

A first-line rewrite of up to 44 rows in a 5,722-line, 904 KB ratified charter (V13) is not a
~0.3d item, and bundling it with the instrument would mix the measurement with the measured — the
instrument's *first* published reading must be of an unmodified charter, or the reading proves
nothing about the instrument. So:

- **This item ships the instrument, its suite, and the rule.** The charter's queue body is
  untouched by M1; M2 touches the heading line (append `[ROUTED]`) and the Repo Profile only.
- **In the meantime the census reports** `41 closed / 89 rows (tagged-open 4, untagged 44)` at
  `8b15153`, rising by one when the controller's landing commit tags row 72 in the leading position
  — which the live CI step (M2) then verifies.
- **Recipe for the follow-on tag moves, if a controller chooses to do them** (a docs commit, not a
  milestone here; rows 65, 66, 68, 70 are pure position moves; rows 7, 16, 18, 57 need a *word*
  change and are content edits to a ratified row, so they are recorded here and left to an
  attended or explicitly-logged decision): (a) `wc -l` before == after; (b) every changed line in
  `git diff -U0` matches `^[-+][0-9]+\. ` — no body line moves; (c) the census's `rows` total is
  unchanged and the closed count rises by exactly the number of rows moved; (d) the per-row table
  differs *only* on the moved row numbers (`diff` of the two tables restricted to `^row=`). The
  census is what proves the edit lost nothing.

### D3 — Controls in the same call, asserted by NUMBER and by CLASS, so recording them does not spend them

`--control-closed <n>` and `--control-open <n>` are **both mandatory**. The instrument prints the
full table, then asserts row `<n>` classified **closed** / **open** respectively, printing
`control: row=<n> expect=<class> got=<TAG> ok` or `… MISMATCH`, and exits 2 on any mismatch (the
table is still printed first so the reader can see what it did). A run without both controls is
`✗ controls required` exit 2 — a reading with no calibration is not published.

The published values for the live call are **closed = row 1** (`[LANDED 2026-07-24]`, will never
reopen) and **open = row 79** (`[PARKED — DESIGN REVIEW AND USER APPROVAL] w-evidence-applicability`,
parked on an attended decision). *A control you record is a control you spend* — the trap is a
census that greps the charter for a literal, which the charter then contains because the census
doc was written into it. This design does not have that hole: a control is a row **number** and a
**class**, matched against the anchored first-token classifier, and the classifier never reads
prose. Writing "`--control-open 79`" into a STATUS stamp adds a number to prose the instrument does
not scan. When row 79 eventually closes, the census exits 2 `CONTROL MISMATCH` and the closing
edit re-points the flag — the mismatch is the control working, not a defect.

### D4 — Anti-vacuity floors, ordered, one arm each, and how each arm REACHES its floor

The charter's iter-108 bullet (V15) records why floors go unpinned: the harness helper that builds a
gate's input emits the very field the floor checks for, so the floor's input shape cannot be
produced. This suite therefore has **no helper that synthesises a heading, a row number, or a
tag**. Its only fixture builder is `make_doc <path>` reading literal lines from a heredoc; each floor
arm writes the broken shape *by hand*. Floors fire in this order, each with its own message, and
every floor arm asserts its **specific** message — so if a different floor fires first (e.g. the
`zero-rows` fixture forgot its heading and `no-heading` fired), the arm reds on message mismatch,
which is the proof it reached the floor it names.

| # | floor (protects the instrument, not the input) | message (stderr) | exit |
|---|---|---|---|
| F0 | mission doc missing/unreadable | `✗ mission doc unreadable: <path>` | 2 |
| F1 | no line matches `^## Queue` | `✗ no "## Queue" heading in <path>` | 2 |
| F2 | zero rows enumerated after the heading | `✗ zero queue rows enumerated after line <n>` | 2 |
| F3 | a row number appears twice | `✗ duplicate row number: <n> (lines <a>, <b>)` | 2 |
| F4 | numbering is not `1..max` | `✗ numbering gap: expected 1..<max>, missing <list>` | 2 |
| F5 | a control flag is absent | `✗ controls required: --control-closed <n> --control-open <n>` | 2 |
| F6 | a control row is not enumerated or has the wrong class | `✗ CONTROL MISMATCH: row <n> expect=<class> got=<TAG or ABSENT>` | 2 |

F3/F4 are deliberately loud rather than tolerant: a column-0 numbered list inside a row body would
collide with a row number, and the right response is "a human looks", not a guess. (Today none
exists: the 89 numbered lines after the heading are exactly the 89 rows, V2. The two column-0
numbered lines *before* the heading — lines 866, 869 — are why the scan starts at the heading, V14.)

Exit codes are **0** (census printed, both controls ok) and **2** (any floor). There is no exit 1:
UNTAGGED > 0 is the queue's normal state, not an alarm.

### D5 — Report format

```
queue_census: doc=design_docs/world-mission.md heading_line=1610
row=1 tag=LANDED class=closed line=1637
row=2 tag=LANDED class=closed line=1643
…
row=39 tag=- class=UNTAGGED line=4396
…
preamble: 1 unnumbered closed block (line 1615) excluded from rows
control: row=1 expect=closed got=LANDED ok
control: row=79 expect=open got=PARKED ok
census: 41 closed / 89 rows (tagged-open 4, untagged 44) [LANDED 40, ROUTED 1, RULED OUT 0; PARKED 3, IN-SPRINT 1, NEXT 0] [queue_census.sh, charter-blob 2dda66a36a6bc60ab280a73b3f0b3c6be89ed399, controls 1/79]
```

**The block above is the reading at the BASE charter `8b15153`** (`heading_line=1610`, blob
`2dda66a3…`). At the commit this item lands, M2's three Repo Profile bullets shift the heading to
`1613` and the blob to `52aaf9e8…`, and the count is unchanged at `41 closed / 89 rows` — which is
the point of citing the blob rather than a line number. Every blob literal in this document is
labelled with the tree it was read from; none is a claim about the current tree.

Rows print **sorted by row number** regardless of file order; `line=` is the row-start line in the
doc. The `preamble:` line reports any column-0 `**[` block after the heading that is not a numbered
row (exactly one today, V4) — counted **out** of `rows`, always printed (even `0`), so the
denominator decision is visible in every reading. The last line is the only line a stamp may quote.

**The total line identifies the exact bytes it was read from** (quorum r1, `gpt5-6-sol`'s fix applied
verbatim, §10): `charter-blob <oid>` is `git hash-object <doc>` of the input path, computed by the
script itself (`--git-bin` overridable, default `git`). `git hash-object` needs no repository — it
hashes a file — but it does need a git binary and a readable path (a missing path is `rc=128`,
V21). If the OID cannot be produced the script prints the distinct token **`charter-blob
unavailable`** on the total line and a stderr warning, and still exits 0 — the table is valid, the
*reading* is not publishable (D6). It never silently omits the field. The `controls a/b` suffix
names the two control rows passed.

**Implementation constraints (bash 3.2, portable awk).** The rig's `/bin/bash` is 3.2.57 and its
`awk` is BWK 20200816 (V17); CI is `ubuntu-latest` (mawk). Enumeration, classification, dup/gap
detection and sorting live in ONE awk program invoked by the bash wrapper — no bash associative
arrays, no `{n,m}` intervals, no `\b`, no `gensub`, no `asort`; `sort -n` does the ordering. The
wrapper owns argument parsing, F0, F5, F6 and exit codes. It does **not** `cd` to the repo root
(the precedent's note): `--doc <path>` defaults to `design_docs/world-mission.md` relative to cwd.

### D6 — Where the number gets published: a reading cites its instrument, or it is not published

- **Mission log `**Progress:**` and the STATUS stamp may carry a closed-row count again, ONLY as a
  census citation**, in EXACTLY the form the script prints (quorum r2, `gemini-3-1-pro`'s fix
  applied verbatim — the previous, truncated form omitted `tagged-open` and the per-tag array and
  so contradicted D5's *"the last line is the only line a stamp may quote"*):
  `census: N closed / M rows (tagged-open T, untagged U) [LANDED X, ROUTED Y, RULED OUT Z; PARKED P, IN-SPRINT Q, NEXT R] [queue_census.sh, charter-blob <oid>, controls a/b]`.
  The live reading at `8b15153` is therefore quoted as
  `census: 41 closed / 89 rows (tagged-open 4, untagged 44) [LANDED 40, ROUTED 1, RULED OUT 0; PARKED 3, IN-SPRINT 1, NEXT 0] [queue_census.sh, charter-blob 2dda66a36a6bc60ab280a73b3f0b3c6be89ed399, controls 1/79]`.
  The `charter-blob <oid>`
  is the blob OID of the exact charter bytes the reading was taken from, printed by the script (D5)
  and checkable with `git hash-object design_docs/world-mission.md`; the `controls a/b` are the two
  row numbers passed. A stamp with no fresh reading prints **no number** (the current `goal
  unmoved` practice), never a carried one; a line reading `charter-blob unavailable` is not
  publishable. This is a Repo Profile bullet (M2); World may not edit the shared skill, and the
  shared skill's Gate-4 text names "queue tags" without a position rule (V16), so nothing upstream
  conflicts.
- **Why a blob OID and not a commit SHA** (quorum r1, `gpt5-6-sol`, adjudicated TRUE by the
  controller): a commit cannot contain its own SHA — the message and every tracked file are inputs
  to it — so a reading that cited the commit it lands in was impossible as written. A blob OID is
  fixed the moment the file's bytes are, before any commit exists, and identifies the bytes directly
  where a commit SHA does so only transitively. **Reachability caveat, stated honestly:** the
  reading is taken against a charter state, and the STATUS stamp lives *inside* the charter, so a
  stamp that quotes a reading necessarily names the blob of the charter *before the stamp line was
  written*. That blob is reachable in git only if it was committed — which the house workflow
  already does: the sprint PR's squash commit carries the row's `[LANDED …]` tag (e.g. `1ee9ea8`
  tagged row 71), and the separate iteration-record commit (e.g. `8b15153`) writes the stamp and the
  log. So: **take the reading against the committed tag state, publish it in the record commit's
  log entry and stamp; the cited blob is then `<tag-commit>:design_docs/world-mission.md`.** A
  reading published in the mission log or a commit message needs no caveat — those are different
  objects from the charter.
- **A row tagged closed is re-read**: the commit that adds a `[LANDED …]` tag runs the census
  (AC-M2-5 proves the printed blob is the file's blob) and the next record commit publishes the
  reading. The live CI step (M2) is what makes a mis-positioned tag visible — it does not move the
  count — and a control flip reds. **What no CI step can own, and this doc therefore does not
  claim:** that the count "moved by exactly one" across two commits. CI sees one tree; a delta is a
  controller reading of two census lines, and it is left as exactly that.

---

## §3 Files to create / modify

- `scripts/queue_census.sh` — **new** (M1). The instrument, D1–D5, including the `charter-blob`
  field (`git hash-object` of `--doc`, `--git-bin` overridable, loud `unavailable` token).
- `scripts/test_queue_census.sh` — **new** (M1). Arms, `--only <arm>`, `ok`/`not ok` lines, final
  `N passed, M failed`, exit non-zero on any failure — the precedent's shape. The blob arms compute
  `git hash-object` on their own fixture file, so they depend on a git *binary*, never on this
  repository's git state (both CI jobs check out shallow).
- `scripts/testdata/queue_census_rowstarts_8b15153.md` — **new** (M1). Frozen snapshot of the
  charter's queue *row starts* at `8b15153`, **generated by this command, checked in as produced**:
  `awk '/^## Queue/{h=NR; print substr($0,1,120); next} h && (/^[0-9]+\. / || /^\*\*\[/) {print substr($0,1,120)}' design_docs/world-mission.md`
  → 91 lines, 10,038 bytes (V18): heading + preamble + 89 first lines truncated to 120 bytes. The
  classifier reads only the leading token, so truncation cannot change a class; bodies are absent
  by construction, so the body traps (`## ` line, prose keyword) get **synthetic** arms instead.
  Its expected reading is the prototype's: `41 closed / 89 rows (tagged-open 4, untagged 44)` — two
  independent implementations, one number.
- `.github/workflows/ci.yml` — **modify** (M1: suite step; M2: live-charter step), both in
  `go-verify` immediately after `Gate-1 range instrument suite` (line 202–204, V19). Additive.
- `design_docs/world-mission.md` — **modify** (M2 only): heading line 1610 gains `[ROUTED]`; three
  Repo Profile bullets (D1 position rule, D6 citation rule, D3 control re-point rule). No queue row
  body is touched by this item.
- `design_docs/planned/w-queue-closed-count-census.md` — this document (controller moves it to
  `implemented/` at landing).

---

## §4 Acceptance criteria

Every criterion is a command run from the repo root **under `bash`**, with its expected output.
**All were baselined on the pristine tree at `8b15153` and are RED there** (§7 B1–B9); a criterion
green at base would measure the repo, not the change. `grep` is `/usr/bin/grep` throughout.

**M1 — instrument + suite + fixture + suite CI step**

- **AC-M1-1** — both scripts exist and are executable.
  `test -x scripts/queue_census.sh && test -x scripts/test_queue_census.sh && echo EXISTS` → `EXISTS`. (Base: `rc=1`, B1/B2.)
- **AC-M1-2** — enumeration is order-independent and line-accurate. Arm `enumerate` writes rows
  `3` (`[LANDED …`), `1` (`[LANDED …`), `2` (`[PARKED …`) in that file order, asserts the table
  prints `row=1`, `row=2`, `row=3` in that order, each with the correct `line=`, and with controls
  `--control-closed 1 --control-open 2`, `2 closed / 3 rows (tagged-open 1, untagged 0)`.
  `bash scripts/test_queue_census.sh --only enumerate; echo "rc=$?"` → `ok …` lines, `rc=0`. (Base: `rc=127`, B3.)
- **AC-M1-3** — the closed vocabulary and every decoration shape classify closed. Arm
  `closed-vocab` writes rows 1–6 as `[LANDED …`, `[**LANDED …`, `[**[LANDED …`, `**[LANDED …`,
  `**[ROUTED …`, `**[RULED OUT …`, **plus row 7 `[PARKED] control`** (quorum r2, `gpt5-6-sol`'s
  fix verbatim: mandatory controls mean every success fixture must contain a row of BOTH classes)
  → six `class=closed`, and with controls `--control-closed 1 --control-open 7`
  `6 closed / 7 rows (tagged-open 1, untagged 0)`.
  `bash scripts/test_queue_census.sh --only closed-vocab; echo "rc=$?"` → `rc=0`.
- **AC-M1-4** — the open vocabulary classifies open, not closed. Arm `open-vocab`: `[PARKED …`,
  `**[NEXT] …`, `[**[IN-SPRINT] …` as rows 1–3, **plus row 4 `[LANDED] control`** (quorum r2,
  `gpt5-6-sol` verbatim) → three `class=open`, and with controls `--control-closed 4 --control-open 1`
  `1 closed / 4 rows (tagged-open 3, untagged 0)`.
  `bash scripts/test_queue_census.sh --only open-vocab; echo "rc=$?"` → `rc=0`.
- **AC-M1-5** — an untagged row is UNTAGGED, counted in rows, not in closed. Arm `untagged`: one
  LANDED row (1), one PARKED row (2) and one `**w-slug** · clause-2 · prose` row (3) — the
  three-class shape `gpt5-6-sol` specified in quorum r2 — → row 3 prints `tag=- class=UNTAGGED`,
  and with controls `--control-closed 1 --control-open 2`
  `1 closed / 3 rows (tagged-open 1, untagged 1)`.
  `bash scripts/test_queue_census.sh --only untagged; echo "rc=$?"` → `rc=0`.
- **AC-M1-6** — **prose is never a tag (the row-57 defect).** Arm `prose-not-tag` writes an
  untagged row whose body contains the lines `[LANDED] [ROUTED] [RULED OUT] [PARKED] [NEXT]` and
  `it landed correctly` and `LANDED 2026-09-08 (iter-1)`, and a second row `**w-slug** · **[LANDED
  2026-09-08]** …` (slug first), **plus one `[LANDED] control` row and one `[PARKED] control`
  row** (quorum r2, `gpt5-6-sol` verbatim: "adds one LANDED and one PARKED control row and expects
  the two target rows to remain UNTAGGED within a four-row census"). Numbering: 1 `[LANDED] control`,
  2 `[PARKED] control`, 3 and 4 the two target rows. Both targets → `class=UNTAGGED`; with controls
  `--control-closed 1 --control-open 2`, `1 closed / 4 rows (tagged-open 1, untagged 2)`.
  `bash scripts/test_queue_census.sh --only prose-not-tag; echo "rc=$?"` → `rc=0`.
- **AC-M1-7** — a `## ` line inside a row body does not end the enumeration (attempt 2's bug). Arm
  `body-heading`: row 1 (`[LANDED …`), then a body line `## Premise Verification Log`, then row 2
  (`[PARKED …`) and row 3 (`[LANDED …`) → with controls `--control-closed 1 --control-open 2`,
  `2 closed / 3 rows (tagged-open 1, untagged 0)`.
  `bash scripts/test_queue_census.sh --only body-heading; echo "rc=$?"` → `rc=0`.
- **AC-M1-8** — the unnumbered preamble block is reported and excluded. Arm `preamble`: heading,
  `**[LANDED …] w-m1** …`, row 1 (`[LANDED …`) and row 2 (`[PARKED …`) → with controls
  `--control-closed 1 --control-open 2`, `preamble: 1 unnumbered closed block (line 2) excluded
  from rows` and `1 closed / 2 rows (tagged-open 1, untagged 0)`; a second fixture with the same
  two rows and no preamble prints `preamble: 0 …` under the same controls.
  `bash scripts/test_queue_census.sh --only preamble; echo "rc=$?"` → `rc=0`.
- **AC-M1-9** — floors F0–F4, one arm each, each asserting its own message and `rc=2`:
  `bash scripts/test_queue_census.sh --only no-doc; bash scripts/test_queue_census.sh --only no-heading; bash scripts/test_queue_census.sh --only zero-rows; bash scripts/test_queue_census.sh --only dup-number; bash scripts/test_queue_census.sh --only gap; echo "rc=$?"`
  → five `ok …` blocks, final `rc=0`. Fixture shapes: `no-doc` passes `--doc /nonexistent/x.md`;
  `no-heading` has rows and no `## Queue` line; `zero-rows` has the heading and prose only;
  `dup-number` has rows `1, 2, 2`; `gap` has rows `1, 3`. **These five arms pass no control flags
  and must not**: the floors are ORDERED (D4), F0–F4 all precede F5, and each arm asserts its own
  floor's message — which is exactly what proves it reached the floor it names rather than falling
  through to `controls required`. (Controller clarification at the r2 carve-out, not a reviewer
  fix: `gpt5-6-sol`'s objection is scoped to *success* arms, and this states why the failure arms
  are outside it.)
- **AC-M1-10** — controls are mandatory (F5) and adjudicated (F6). Arm `controls-required` runs a
  valid doc with no flags → `controls required`, `rc=2`. Arm `control-mismatch` uses a three-row
  fixture (1 `[LANDED …`, 2 `[PARKED …`, 3 untagged) and runs it twice with BOTH flags present:
  `--control-closed 2 --control-open 2` (closed control aimed at an open row) and
  `--control-closed 1 --control-open 3` (open control aimed at an UNTAGGED row) → each
  `CONTROL MISMATCH: row=<n> expect=… got=…`, `rc=2`, **and the per-row table is still printed
  above it**.
  `bash scripts/test_queue_census.sh --only controls-required; bash scripts/test_queue_census.sh --only control-mismatch; echo "rc=$?"` → `rc=0`.
- **AC-M1-11** — the report format (D5). Arm `report`: a 4-row fixture (2 closed, 1 open, 1
  untagged) with controls `--control-closed 1 --control-open 3` → last line matches exactly
  `census: 2 closed / 4 rows (tagged-open 1, untagged 1) [LANDED 1, ROUTED 1, RULED OUT 0; PARKED 1, IN-SPRINT 0, NEXT 0] [queue_census.sh, charter-blob <40 hex>, controls 1/3]`
  (the OID asserted as `[0-9a-f]{40}` via `/usr/bin/grep -E`; its *value* is AC-M1-15's job),
  both `control: … ok` lines present, `rc=0`, and the output contains neither `carried` nor `of 4
  rows closed` (the carried-prose form is not reproduced).
  `bash scripts/test_queue_census.sh --only report; echo "rc=$?"` → `rc=0`.
- **AC-M1-15** — the printed blob OID is the input file's blob, and its absence is loud (quorum
  r1 arm). Arm `blob-oid` writes a two-row fixture (1 `[LANDED …`, 2 `[PARKED …`), runs the
  instrument with `--control-closed 1 --control-open 2`, and asserts the `charter-blob`
  value equals `git hash-object <that fixture>`; arm `blob-unavailable` uses the same fixture and
  controls and runs with
  `--git-bin /nonexistent/git` and asserts the total line carries the literal token
  `charter-blob unavailable`, a stderr line names it, and `rc=0`.
  `bash scripts/test_queue_census.sh --only blob-oid; bash scripts/test_queue_census.sh --only blob-unavailable; echo "rc=$?"` → `rc=0`.
- **AC-M1-12** — the heading anchor is the `^## Queue` prefix. Arm `heading-tail`: a doc whose
  heading is `## Queue (top = next; tags: [NEXT] [IN-SPRINT] [PARKED] [LANDED] [RULED OUT] [ROUTED])`
  enumerates its two rows (1 `[LANDED …`, 2 `[PARKED …`); `heading_line=` is reported and, with
  controls `--control-closed 1 --control-open 2`, `1 closed / 2 rows (tagged-open 1, untagged 0)`.
  `bash scripts/test_queue_census.sh --only heading-tail; echo "rc=$?"` → `rc=0`.
- **AC-M1-13** — the frozen charter snapshot reproduces the prototype. Arm `snapshot` runs the
  instrument on `scripts/testdata/queue_census_rowstarts_8b15153.md` with `--control-closed 1
  --control-open 79` → `census: 41 closed / 89 rows (tagged-open 4, untagged 44)` prefix,
  `preamble: 1 …`, both controls ok, and `row=16`, `row=18`, `row=66`, `row=68`, `row=70` each
  `class=UNTAGGED` (the mis-positioned closed rows are *reported*, not guessed).
  `bash scripts/test_queue_census.sh --only snapshot; echo "rc=$?"` → `rc=0`. (Base: fixture absent, B9.)
- **AC-M1-14** — the whole suite is green and wired into CI.
  `bash scripts/test_queue_census.sh; echo "rc=$?"` → `… passed, 0 failed`, `rc=0`;
  `/usr/bin/grep -n 'test_queue_census' .github/workflows/ci.yml` → one `run: bash scripts/test_queue_census.sh` line inside `go-verify`. (Base: `rc=1`, B4.)

**M2 — live-charter CI step, heading vocabulary, Repo Profile rule**

- **AC-M2-1** — the live charter census runs clean with the published controls.
  `bash scripts/queue_census.sh --doc design_docs/world-mission.md --control-closed 1 --control-open 79; echo "rc=$?"`
  → a `row=` line per row, `preamble: 1 …`, two `control: … ok`, a `census: N closed / M rows …`
  last line, `rc=0`. (Base: `rc=127`, B8.) The reading at landing is recorded in §7 by the
  controller, and the row count cross-checks the independent enumeration:
  `test "$(bash scripts/queue_census.sh --doc design_docs/world-mission.md --control-closed 1 --control-open 79 | /usr/bin/grep -c '^row=')" = "$(awk '/^## Queue/{h=1;next} h && /^[0-9]+\. /' design_docs/world-mission.md | wc -l | tr -d ' ')" && echo AGREE` → `AGREE`.
- **AC-M2-2** — the live step is in CI:
  `/usr/bin/grep -n 'queue_census.sh --doc design_docs/world-mission.md --control-closed 1 --control-open 79' .github/workflows/ci.yml` → one hit in `go-verify`. (Base: `rc=1`, B5.)
- **AC-M2-3** — the heading vocabulary names `[ROUTED]` and the anchor still holds:
  `/usr/bin/grep -c '^## Queue .*\[ROUTED\]' design_docs/world-mission.md` → `1` (base `0`, B6);
  `/usr/bin/grep -c '^## Queue' design_docs/world-mission.md` → `1` (unchanged).
- **AC-M2-4** — the Repo Profile carries the rule, before the queue heading (so a future row
  mentioning the script cannot satisfy it):
  `awk '/^## Queue/{exit} /queue_census\.sh/{c++} END{print c+0}' design_docs/world-mission.md` → `3` or more (base `0`, B7): one bullet each for tag position, citation form, control re-point.
- **AC-M2-5** — the live reading's blob OID is the charter's blob (the reviewer-named arm):
  `test "$(bash scripts/queue_census.sh --doc design_docs/world-mission.md --control-closed 1 --control-open 79 | /usr/bin/grep '^census: ' | sed -E 's/.*charter-blob ([0-9a-f]{40}).*/\1/')" = "$(git hash-object design_docs/world-mission.md)" && echo BLOB-MATCH` → `BLOB-MATCH`. (Base: `rc=127`, B8; the charter's blob at base is `2dda66a36a6bc60ab280a73b3f0b3c6be89ed399`, V21.)
- **AC-M2-6** — the D6 citation protocol is pinned at both ends: the Repo Profile bullet quotes the
  exact form, and the instrument prints that form. `awk '/^## Queue/{exit} /charter-blob/{c++} END{print c+0}' design_docs/world-mission.md` → `1` or more (base `0`, V22); and the `report` arm (AC-M1-11) holds, so the printed line and the quoted rule share the literal `[queue_census.sh, charter-blob …, controls a/b]` shape. The remaining half of the protocol — a controller *publishing* a reading in a log entry — is a procedure step verified at landing by re-running AC-M2-5 against the tag commit's charter, not by CI.

---

## §5 Test plan / mutation matrix

One row per arm; the mutation is applied to `scripts/queue_census.sh` (or the named file), the
expected red set is **exactly** the listed arm(s) with every other arm green. Restore from a
`cp`-captured backup and assert its sha256 — never `git checkout --` (row 82). The enumeration is
anchored to the diff this item ships (instrument, suite, fixture, CI steps, heading, Repo Profile),
not only to the defect it fixes.

| arm | kills this mutation | expected red |
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
| AC-M2-2 | live step's `--doc` pointed at the snapshot fixture (table still plausible, blob wrong) | **AC-M2-2**, not AC-M2-5 — corrected post-evaluation (judge finding F-NB-5, measured): AC-M2-5 runs the *instrument* against the live charter path and is independent of what the CI step's `--doc` says, so this mutation is caught by AC-M2-2's grep on the run line, not by AC-M2-5. |
| AC-M2-6 | citation bullet written after the heading, or quoting a `@ <sha>` form | AC-M2-6 |
| AC-M1-14 grep | CI suite step removed | AC-M1-14 |
| AC-M2-1/2 | live step removed, or control row re-pointed to an UNTAGGED row | AC-M2-2 / AC-M2-1 |
| AC-M2-3 | `[ROUTED]` not appended, or appended to a *new* `## Queue…` line (count becomes 2) | AC-M2-3 |
| AC-M2-4 | rule written after the heading (inside a row) | AC-M2-4 |

**Post-evaluation correction — this table's "expected red" column systematically UNDER-states the
coupling, and the independent judge measured it** (findings F-NB-1, F-NB-2, F-NB-3, and the
`heading-tail` row). Four rows predict a one- or two-arm red set where the real one is much wider,
because almost every success arm greps the census total line, so any mutation that changes the
line's shape or a class count reds all of them. The direction is safe — every mutation is caught by
**more** arms than predicted, never fewer, and the judge confirmed **no mutation escapes the suite**
— but a prediction that is wrong in the generous direction is still a prediction nobody can use to
localise a regression.

**Reproduced first-party by the controller before it was recorded** (not banked from the judge):
mutating the total line's wording (`closed / ` → `closed of `) with a `cp`-captured backup and a
sha256-asserted restore gives **`58 passed, 10 failed`** — ten failing assertions across **nine**
arms (`body-heading`, `closed-vocab`, `enumerate`, `heading-tail`, `open-vocab`, `preamble`,
`prose-not-tag`, `report`, `snapshot`, `untagged`), against this table's predicted `report` only.
Restore verified sha-identical and the suite returned to `68 passed, 0 failed`. The judge reported
11 arms for the same mutation; the difference is arms-vs-assertions counting, and the class is the
same. **A future revision should either widen these predictions to what is measured, or decouple
the arms from the total line so a localised mutation reds a localised set** — the second is the
better fix and is out of scope here.

**Floor reachability is itself asserted:** each floor arm's `ok` requires the floor's own message,
so a harness change that makes a floor unreachable (e.g. a helper that starts auto-inserting the
heading) reds `no-heading` immediately — the iter-108 failure mode is caught by construction.

---

## §6 Conflict Surface

- **`design_docs/world-mission.md` readers/writers (repo-wide grep, V6/V20):**
  `scripts/verify_go.sh` (`check_mission_config`, lines 112–125) runs
  `scripts/mission_decisions.sh --check --file design_docs/world-mission.md`, which parses only the
  `<!-- decision-ledger:start/end -->` block at lines 827–852 — **before** the queue heading and
  disjoint from anything this item touches; `host/verifygate/mission_config_gate_test.go` writes
  synthetic charters containing only a ledger and never invokes the census — unaffected;
  `README.md`/`CLAUDE.md` mention the path only. Nothing in `scripts/`, `.github/`, `host/` or
  `Makefile` reads a queue tag today (V6). The census adds one more *reader*; the only *writer* is
  M2's two edits (heading tail, Repo Profile bullets), both outside the row bodies.
- **The shared skill (V1 checkout, read-only for World, V16):** reads the queue head with
  `grep -A2 "^## Queue" … | tail -2` — prefix-anchored, so M2's heading tail is safe; Gate-3b/Gate-4
  prose mentions `[LANDED]`/`[PARKED]` tags without a position rule — no conflict, and the D1
  position rule is a World Repo Profile bullet, not a skill edit. `tools/launchd/*` reads
  `MISSION_DOC` from the env files only (V16) — frozen core untouched.
- **CI (`.github/workflows/ci.yml`):** two additive steps in `go-verify` after line 204. The live
  step (M2) makes CI red on a structurally broken queue (duplicate/gap) or a flipped control — the
  same shape as the existing ledger `--check` already run by `verify_go.sh` in CI, so a docs-only
  commit reddening CI for a *charter-structure* reason is precedented, not new. Re-point rule for
  the open control lives in the Repo Profile (AC-M2-4).
- **Adjacent queue rows:** 85 (STATUS rotation rule stranded) and 86 (log headings) are housekeeping
  on the same files but different regions; 83 (a row's proposed remedy is a claim too) is honoured
  here by prototyping the remedy against the live charter before proposing it (V9). Row 72's own
  landing tag goes in the D1 position; the live step then verifies it.
- **Row bodies:** none edited. Rows 65/66/68/70's slug-first tags stay where they are until a
  separate docs commit moves them under D2's recipe.

---

## §7 Verification Log

Every claim above that the codebase/rig currently does X has a row here with its COMMAND and
OBSERVED OUTPUT. All commands run under `bash` from `/Users/voightkampff/dev/sunholo-data/.wt-world-iter172`
at `8b151533ed701e6648f7037d86579383b04de407` (`git rev-parse HEAD`, `git branch --show-current` →
`sprint/w-queue-census`; `git status --short` → empty). Negative/empty results carry a positive
control in the same row.

| # | Claim | COMMAND | OBSERVED OUTPUT |
|---|---|---|---|
| V1 | Queue heading is line 1610, and there is exactly one | `/usr/bin/grep -n "^## Queue" design_docs/world-mission.md`; `/usr/bin/grep -c '^## Queue' …` | `1610:## Queue (top = next; tags: [NEXT] [IN-SPRINT] [PARKED] [LANDED] [RULED OUT])`; `1` |
| V2 | 89 numbered rows after it, numbers exactly 1..89, no dups, no gaps | `awk 'NR>1610 && /^[0-9]+\. /' … \| wc -l`; `… \| sed -E 's/^([0-9]+)\..*/\1/' \| sort -n \| uniq -d`; `comm -13 <(…sorted) <(seq 1 89 \| sort -n)` | `89`; empty; empty |
| V3 | Rows are not in file order | `/usr/bin/grep -nE '^(82\|68)\. ' design_docs/world-mission.md` | `5340:82. **w-mutation-drill-…`; `5543:68. **w-driver-pin-…` |
| V4 | Exactly one unnumbered closed block after the heading, at 1615 | `/usr/bin/grep -c '^\*\*\[LANDED' …`; `awk 'NR>1610 && /^\*\*\[/ {print NR}' …` | `1`; `1615` (`**[LANDED 2026-07-24 (iter-13), CI green on dev \`d0009c8\`] w…`) |
| V5 | Tag keyword on 45 first lines; bracket right after the number on 30 | `awk 'NR>1610 && /^[0-9]+\. /' … \| /usr/bin/grep -cE '\[(LANDED\|RULED OUT\|PARKED\|NEXT\|IN-SPRINT)'`; `… \| /usr/bin/grep -cE '^[0-9]+\. \*\*\['` | `45`; `30` |
| V6 | No script reads a closure tag; positive control hits | `ls scripts/`; `/usr/bin/grep -rlnE 'LANDED\|RULED OUT\|closed' scripts/; echo rc=$?`; control `/usr/bin/grep -rlE 'LANDED\|RULED OUT\|closed' design_docs/world-mission.md` | `bench_worldd.sh build_world_package.sh gate1_range_check.sh hooks mission_answer.sh mission_decisions.sh test_gate1_range_check.sh test_mission_answer.sh testdata verify_ail.sh verify_go.sh verify_world_package.sh world_package_ready_packet.golden.json` (the brief's transcription omitted `test_mission_answer.sh`; `git ls-files` confirms it is tracked — no bearing on this design); `rc=1`, zero files; control → `design_docs/world-mission.md` |
| V7 | The `N of M` headline stops after iteration 151 | `/usr/bin/grep -n '^\*\*Progress:\*\*' design_docs/world-mission-log.md \| tail -30` | lines 232, 378, 514, 680, 816 (`54 of 61` … `58 of 70`), 900 (`the carried count is **59 of 72 queue rows closed** … It is carried, not measured`), then 1032/1247/1417… `**Progress:** **goal unmoved.**` — no count |
| V8 | Open row 57's body contains a closure keyword only as lowercase prose | `for w in landed 'ruled out' parked routed next in-sprint complete; do sed -n '5035,5138p' … \| /usr/bin/grep -ci "$w"; done`; `sed -n '5035,5138p' … \| /usr/bin/grep -in landed` | `landed: 1`, all others `0`; `89: \`-title\` and \`-from\` landed correctly in both.` — case-sensitive `/usr/bin/grep -c LANDED` over the same range → `0` |
| V9 | Under D1 the live charter reads 41/89/4/44 — **executable prototype, re-run by the designer after the controller's independent run agreed** (quorum r1, `gemini-3-1-pro`, §10). This program is a prototype for this log, NOT the shipped script: no floors, no controls, no preamble line, no blob field, and it takes the FIRST `^## Queue` line by `!h`. | `PROG='/^## Queue/ && !h { h=NR; next }`<br>`h && /^[0-9]+\. / {`<br>`  n=$0; sub(/^[0-9]+\. /,"",n)`<br>`  sub(/^[\[\*]+/,"",n)`<br>`  tag="-"; cls="UNTAGGED"`<br>`  if (n ~ /^RULED OUT([^A-Za-z]\|$)/)      { tag="RULED OUT"; cls="closed" }`<br>`  else if (n ~ /^LANDED([^A-Za-z]\|$)/)    { tag="LANDED";    cls="closed" }`<br>`  else if (n ~ /^ROUTED([^A-Za-z]\|$)/)    { tag="ROUTED";    cls="closed" }`<br>`  else if (n ~ /^PARKED([^A-Za-z]\|$)/)    { tag="PARKED";    cls="open" }`<br>`  else if (n ~ /^IN-SPRINT([^A-Za-z]\|$)/) { tag="IN-SPRINT"; cls="open" }`<br>`  else if (n ~ /^NEXT([^A-Za-z]\|$)/)      { tag="NEXT";      cls="open" }`<br>`  num=$0; sub(/\..*$/,"",num)`<br>`  printf "%d\t%s\t%s\t%d\n", num, tag, cls, NR`<br>`}'`<br>then: `awk "$PROG" design_docs/world-mission.md \| sort -n \| wc -l`; `awk "$PROG" … \| sort -n \| awk -F'\t' '{c[$3]++; t[$2]++} END{printf "closed=%d open=%d untagged=%d\n", c["closed"], c["open"], c["UNTAGGED"]; for(k in t) printf "  %s=%d\n", k, t[k]}'`; `awk "$PROG" … \| sort -n \| awk -F'\t' '$3=="UNTAGGED"{printf "%s ",$1}'`; `awk "$PROG" … \| sort -n \| /usr/bin/grep -E '^(1\|79)\b'`; and `diff <(awk -f /tmp/proto_census_iter172.awk … \| sort -n) <(awk "$PROG" … \| sort -n) && echo IDENTICAL` against the controller's file | `89`; `closed=41 open=4 untagged=44` / `ROUTED=1` / `PARKED=3` / `LANDED=40` / `-=44` / `IN-SPRINT=1` (RULED OUT and NEXT absent = 0); UNTAGGED = `7 16 18 22 23 24 25 26 27 28 29 30 31 32 33 34 35 36 37 38 39 40 57 65 66 68 69 70 72 73 74 75 76 77 78 81 82 83 84 85 86 87 88 89` (44 numbers, listed literally so a re-run can be diffed); `1 LANDED closed 1637`, `79 PARKED open 5711`; `IDENTICAL`. Per-class lists from the designer's first prototype (same rule, `substr` loop instead of `sub`): `LANDED: 40 -> 1 2 3 4 5 9 10 11 12 13 14 15 17 19 20 21 41 42 43 44 45 46 47 48 49 50 51 52 54 55 56 58 59 60 61 62 63 64 67 71`; `ROUTED: 1 -> 53`; `PARKED: 3 -> 6 79 80`; `IN-SPRINT: 1 -> 8` |
| V10 | Five closed rows are UNTAGGED by position/word | `for n in 7 16 18 40 57 65 66 68 70; do /usr/bin/grep -nE "^$n\. " … \| cut -c1-110; done` | `16. [**[COMPLETE 2026-08-13 (iter-80)] — M1 LANDED …`; `18. [**ITEM COMPLETE — M3 LANDED …`; `66. **w-flow-key-…** · **[LANDED 2026-09-07 …`; `68. **w-driver-pin-…** · **[LANDED 2026-09-08 …`; `70. **w-gate-1-…** · **[LANDED 2026-09-08 …`; also `7. [**PARK CONDITION CORRECTED …`, `40. **w-a2a-…** · … **[BLOCKED on row 39]**`, `57. **[AMENDED …`, `65. **w-go-build-…** · **[PARKED — …` |
| V11 | A `## ` heading lives inside a row body after the queue heading | `awk 'NR>1610 && /^## /{print NR": "$0}' …` | `5512: ## Premise Verification Log (quorum objection #1 — …)` (between row 89 at 5493 and row 67 at 5541) |
| V12 | Row 72's own body lists the tag literals (the spent-control trap) | `sed -n '5573,5596p' … \| /usr/bin/grep -c '\[LANDED\]\|\[RULED OUT\]\|\[PARKED\]\|\[NEXT\]\|\[ROUTED\]'`; `/usr/bin/grep -c '\[ROUTED\]' …` | `1`; `1` (the only `[ROUTED]` in the charter is row 72's prose — the heading has none, B6) |
| V13 | Charter size | `wc -l design_docs/world-mission.md`; `wc -c …` | `5722`; `904369` |
| V14 | Two column-0 numbered lines exist BEFORE the heading | `awk 'NR<=1610 && /^[0-9]+\. / {print NR": "substr($0,1,60)}' …` | `866: 1. ~~Iteration 0: ratify the bar~~ …`; `869: 2. **NOW**: work the queue …` |
| V15 | The iter-108 floors bullet exists | `/usr/bin/grep -n 'iter-108' design_docs/world-mission.md \| head -1` | `115:- **AN ANTI-VACUITY FLOOR IS A BRANCH TOO, AND RULE 3j's ENUMERATION STOPS AT *REFUSAL* BRANCHES — SO THE CHECKS THAT PROTECT THE INSTRUMENT ITSELF ARE SYSTEMATICALLY THE LAST THING ANYONE PINS …` |
| V16 | Shared skill reads the heading prefix-anchored; driver reads only `MISSION_DOC`; no position rule upstream | in `/Users/voightkampff/dev/sunholo-data/ailang`: `/usr/bin/grep -rn '## Queue' tools/launchd .claude/skills/mission-control`; `/usr/bin/grep -rnE '\[(LANDED\|RULED OUT\|IN-SPRINT)\]' …`; control `/usr/bin/grep -rln MISSION_DOC tools/launchd`; `/usr/bin/grep -rniE 'leading tag\|fixed position' …/SKILL.md; echo rc=$?` | `SKILL.md:591: … grep -A2 "^## Queue" "${MISSION_DOC:-…}" \| tail -2`; `gate-3b-ci-green.md:261 … upgrades the queue tag to [LANDED]`, `gate-4-record.md:80 … queue tags ([LANDED], [PARKED], etc.)`; control → 3 env files; `rc=1` (no position rule) |
| V17 | Rig shell/awk versions | `/bin/bash --version \| head -1`; `awk --version` | `GNU bash, version 3.2.57(1)-release (arm64-apple-darwin25)`; `awk version 20200816` |
| V18 | Snapshot fixture size when generated by the named command | `awk '/^## Queue/{h=NR; print substr($0,1,120); next} h && (/^[0-9]+\. (\*\*\|\[)/ \|\| /^\*\*\[/) {print substr($0,1,120)}' … \| wc -lc` (all 89 row starts match `^[0-9]+\. (\*\*\|\[)` — `awk 'NR>1610 && /^[0-9]+\. (\*\*\|\[)/' … \| wc -l` → `89`) | `91   10038` |
| V19 | CI insertion point and the precedent step | `/usr/bin/grep -n 'Gate-1 range' -A 2 .github/workflows/ci.yml`; `/usr/bin/grep -nE '^  [a-z-]+:$' .github/workflows/ci.yml` | `202: - name: Gate-1 range instrument suite / 203: timeout-minutes: 2 / 204: run: bash scripts/test_gate1_range_check.sh`; jobs `ailang-verify` (18), `go-verify` (99) |
| V20 | Repo readers of the charter path outside `design_docs/` | `/usr/bin/grep -rln 'world-mission\.md' --exclude-dir=.git --exclude-dir=design_docs .`; `/usr/bin/grep -n 'check_mission_config' scripts/verify_go.sh`; `/usr/bin/grep -n 'decision-ledger:start\|decision-ledger:end' design_docs/world-mission.md` | `host/verifygate/mission_config_gate_test.go`, `README.md`, `scripts/verify_go.sh`, `CLAUDE.md`; `112`, `142`, `218`; `827`, `852` |
| V21 | The charter's blob OID at base, identical in the worktree, the main checkout and the commit; a missing path is loud | `git hash-object design_docs/world-mission.md` (worktree); `git -C /Users/voightkampff/dev/sunholo-data/ailang-world hash-object design_docs/world-mission.md`; `git rev-parse 8b15153:design_docs/world-mission.md`; `git hash-object /nonexistent/x.md; echo rc=$?` | `2dda66a36a6bc60ab280a73b3f0b3c6be89ed399` ×3; `fatal: could not open '/nonexistent/x.md' for reading: No such file or directory`, `rc=128` |
| V22 | AC-M2-6 red at base | `awk '/^## Queue/{exit} /charter-blob/{c++} END{print c+0}' design_docs/world-mission.md` | `0` |
| B1–B2 | AC-M1-1 red at base | `test -x scripts/queue_census.sh; echo rc=$?`; same for the suite | `rc=1`; `rc=1` |
| B3 | AC-M1-2…13 red at base; the precedent suite is the positive control | `bash scripts/test_queue_census.sh; echo rc=$?`; control `bash scripts/test_gate1_range_check.sh --only report \| tail -1` | `bash: scripts/test_queue_census.sh: No such file or directory`, `rc=127`; control `4 passed, 0 failed` rc=0 |
| B4 | AC-M1-14 grep red at base; control hits | `/usr/bin/grep -n 'test_queue_census' .github/workflows/ci.yml; echo rc=$?`; control `… 'test_gate1_range_check' …` | `rc=1`; `204:        run: bash scripts/test_gate1_range_check.sh` |
| B5 | Nothing named `queue_census` exists anywhere; control hits | `/usr/bin/grep -rn 'queue_census' --exclude-dir=.git .; echo rc=$?`; control `/usr/bin/grep -rln 'gate1_range_check' --exclude-dir=.git . \| head -3` | `rc=1`; `design_docs/world-mission-log.md`, `-index.md`, `-dashboard.md` |
| B6 | AC-M2-3 red at base | `/usr/bin/grep -c '^## Queue .*\[ROUTED\]' design_docs/world-mission.md` | `0` |
| B7 | AC-M2-4 red at base | `awk '/^## Queue/{exit} /queue_census\.sh/{c++} END{print c+0}' design_docs/world-mission.md` | `0` |
| B8 | AC-M2-1 red at base | `bash scripts/queue_census.sh --doc design_docs/world-mission.md --control-closed 1 --control-open 79; echo rc=$?` | `No such file or directory`, `rc=127` |
| B9 | Fixture absent at base | `ls scripts/testdata/queue_census_rowstarts_8b15153.md; echo rc=$?` | `No such file or directory`, `rc=1` |

---

## §8 Milestones

Total ~0.3d. Two milestones, each independently CI-green and committable; neither touches
`tools/launchd/*`, the shared skill, or any queue row body.

- **M1 — the instrument, its suite, the frozen snapshot, and the suite's CI step (~0.2d).**
  `scripts/queue_census.sh` (D1 classifier, D3 controls, D4 floors, D5 report; bash 3.2 + one
  portable awk program), `scripts/test_queue_census.sh` (19 arms, `--only`, literal-line fixtures
  only), `scripts/testdata/queue_census_rowstarts_8b15153.md` (generated by the §3 command, checked
  in as produced), and a `Queue census instrument suite` step in `go-verify` after the Gate-1 step.
  Green when AC-M1-1…15 hold and the three profile gates are unchanged-green. Useful alone: from
  this commit on, `bash scripts/queue_census.sh --control-closed 1 --control-open 79` produces a
  reading anyone can reproduce, stamped with the blob it was read from.
- **M2 — the live-charter CI step, the heading vocabulary, and the Repo Profile rule (~0.1d).**
  A `Queue census (live charter, controlled)` step in `go-verify` running the AC-M2-1 command;
  heading line 1610 gains `[ROUTED]`; three Repo Profile bullets (tag position = first token after
  the row number, decoration-tolerant; a stamp/log count is published only as
  `census: N closed / M rows (tagged-open T, untagged U) [LANDED X, ROUTED Y, RULED OUT Z; PARKED P, IN-SPRINT Q, NEXT R] [queue_census.sh, charter-blob <oid>, controls a/b]`
  — byte-for-byte the script's own total line (quorum r2, `gemini-3-1-pro`) — never carried, never
  with `charter-blob unavailable`; when a control row changes class, the same
  commit re-points the flag in `ci.yml`). Green when AC-M2-1…6 hold. The controller's landing
  commit tags row 72 `[LANDED …]` in the D1 position; the census is run against that committed
  state (AC-M2-5) and the reading is published in the following record commit's log entry and
  stamp, citing that blob (D6). The former claim that "the live step verifies the count moved by
  exactly one" is **deleted**: it was a two-commit delta no CI step can own (quorum r1, §10).

---

## §9 Deferred scope and declared residuals

- **Normalising the 44 UNTAGGED rows** — out (D2), with the safe recipe recorded. The four
  position-only moves (65, 66, 68, 70) are a candidate docs commit for any later fire; the four
  word changes (7, 16, 18, 57) edit ratified prose and need an explicit decision.
- **UNTAGGED is not split into "plain" vs "bracketed-but-unrecognised".** Any such split reads the
  rest of the first line, which is the road back to prose matching. If the mission wants
  `[BLOCKED on …]` or `[AMENDED …]` recognised, that is a vocabulary change to the heading and the
  script's table, made together and measured by the `snapshot` arm.
- **The instrument does not read the STATUS rotation rule or the log.** Row 85 (rotation-rule body
  stranded) and row 86 (log headings) stay separate.
- **CI portability of the awk program** is asserted by CI itself (the suite runs on `ubuntu-latest`
  mawk and on the rig's BWK awk); no `gawk`-only construct is permitted (D5 constraints).

---

## §10 Quorum verification log

**Round 1 — verdict `blocked`, N=2, both reviewers present, `absent_reviewers = []`.** Artifact:
`.ailang/state/mission-quorum/w-queue-closed-count-census-2026-09-08T04-10-03Z.json` (main
checkout state dir; verdict/absentees/costs read back from the JSON by the designer). Costs:
`gpt5-6-sol` **$0.073785**, `gemini-3-1-pro` **$0.030432**. Both objections were narrow, carried a
reviewer-authored fix, and neither disputed the design direction; this is the one protocol-mandated
revision pass.

### Objection 1 — `gpt5-6-sol` (reject) — adjudicated TRUE by the controller; fix applied verbatim

> **strongest_objection**: "D6 requires a reading to cite the SHA of the charter commit while also
> requiring the commit that changes the closure tag to publish that citation in its own commit
> message or stamp. A Git commit cannot reliably contain its own SHA because the message and tracked
> stamp are inputs to that SHA; amending the citation changes the SHA again. The landing protocol is
> therefore impossible as written and violates the deterministic-behavior axiom."
>
> **catch**: "The contradiction recurs in §8: the landing commit must tag row 72 and quote the
> census, but the required citation identifies that same not-yet-stable commit. AC-M2-1 through
> AC-M2-4 also never test the D6 citation protocol or the claim that the count moved by exactly one."
>
> **proposed_fix**: "Replace the self-referential SHA with an identifier available before commit
> creation, such as the charter blob OID: `blob=$(git hash-object design_docs/world-mission.md)`, and
> require `census: N closed / M rows (untagged U) [queue_census.sh, charter-blob <oid>, controls
> a/b]`. Have the script print that blob OID and add an acceptance arm verifying it matches `git
> hash-object` for the input document. Alternatively, define a two-commit protocol: commit A changes
> the tag; commit B publishes a reading citing commit A. Remove all claims that the same commit can
> cite its own SHA."

**Controller adjudication:** TRUE by the deterministic-SHA argument. **Applied (blob-OID form, not
the two-commit alternative):** D6's citation form is now `[queue_census.sh, charter-blob <oid>,
controls a/b]`; every `@ <sha>` occurrence and every sentence claiming a commit can cite its own
SHA is removed from D6 and §8; D5 specifies the script prints `charter-blob <oid>` via `git
hash-object` (`--git-bin` overridable) and degrades LOUDLY to the literal `charter-blob
unavailable` plus a stderr line, never omitting the field; AC-M1-15 adds the `blob-oid` and
`blob-unavailable` arms (fixture-local `git hash-object`, no dependence on the repo's git state);
AC-M2-5 is the reviewer-named arm (printed OID == `git hash-object design_docs/world-mission.md`);
AC-M2-6 pins the citation protocol at both ends (Repo Profile quotes the form; the `report` arm
prints it). **The "count moved by exactly one" claim is deleted, not tested:** it is a delta
between two commits that no CI step can own, and §8 now says so. **Recorded honestly, not
smoothed:** a STATUS stamp lives inside the charter, so a stamp quoting a reading names the blob of
the charter *before* the stamp line — D6 states this and resolves it through the house workflow's
existing two commits (tag in the PR squash, reading in the record commit), citing
`<tag-commit>:design_docs/world-mission.md`. Not applied: the two-commit *protocol* as a new rule
(out of scope for a ~0.3d item and unnecessary given the OID).

### Objection 2 — `gemini-3-1-pro` (reject) — adjudicated TRUE by the controller; the number independently CONFIRMED

> **strongest_objection**: "The verification log entry V9 uses English pseudo-code ('the prototype
> awk in §2 D1... | sort -n | awk...') instead of an executable command. This hides the exact regex
> and text-munging logic used to prove the 41/89/4/44 classification baseline, leaving the core
> premise of the D1 classifier's feasibility unverified and irreproducible."
>
> **proposed_fix**: "Replace the pseudo-code in V9's COMMAND column with the exact, executable
> awk/bash one-liner that was actually run against the charter to produce the observed per-class
> tally."

**Controller adjudication:** TRUE; and per the obligation on an unverified-premise objection the
controller wrote a SECOND, INDEPENDENT implementation of the D1 rule
(`/tmp/proto_census_iter172.awk`, `sub(/^[\[\*]+/,"")` instead of the designer's `substr` loop),
ran it from the main checkout at `8b15153` under `bash`, and got `89` / `closed=41 open=4
untagged=44` / the same 44-number untagged list / `1 LANDED closed 1637`, `79 PARKED open 5711`.
**Applied:** V9's COMMAND column is now the executable program and its four invocations, **re-run by
the designer in the worktree before recording** (not transcribed), with a fifth invocation diffing
its output against the controller's file → `IDENTICAL`. The untagged list is written as 44 literal
numbers (the earlier `22–40` range notation is gone) so a re-run can be diffed. V9 states the
program is a prototype for the log, not the shipped script (no floors, controls, preamble line or
blob field; first `^## Queue` by `!h`). The agreement of two implementations is recorded as
corroboration of the **number**, not of the prose.

**Not changed:** nothing else. No section a reviewer did not object to was restructured; scope is
unchanged (normalisation still OUT, one tag position, UNTAGGED bucket, controls 1/79).

---

### Round 2 — `blocked`, N=2, `absent_reviewers = []`, closed under the narrow-refinement carve-out

Artifact `.ailang/state/mission-quorum/w-queue-closed-count-census-2026-09-08T04-17-38Z.json`
(rc=3). Costs: `gpt5-6-sol` **$0.095650**, `gemini-3-1-pro` **$0.041312**. Both reviewers PRESENT —
`jq -r '.synthesis.absent_reviewers'` → `[]`, cross-checked with
`[.reviewers[] | select(.present==false) | .model]` → `[]`. So this is a verdict with no hole.

**Disposition.** Both remaining objections (a) carry a concrete reviewer-authored `proposed_fix`
and (b) dispute **completeness/internal consistency**, not the design DIRECTION — the classifier
rule, the UNTAGGED bucket, the scoped-out normalisation and the floor set drew no objection in
either round. That is exactly the Gate-2 **narrow-refinement carve-out**: after the one
protocol-mandated re-quorum, the controller makes a bounded 2nd revision applying the reviewers'
**verbatim** fixes and routes to the planner. It is not a force-pass — both objections are
SATISFIED, not overridden — and it is a controller routing call, not `needs-human-review`
(standing rule 8). This mission has used the carve-out before (iter-169, iter-171), so no
first-use ratification is owed.

**Where the objections landed, per round** (the Gate-2 localisation rule): r1 = D6's citation
identity + §7's V9 record; r2 = D3's control contract vs the §4 fixtures + the D5/D6 format
mismatch. Four objections, four different surfaces, and r2's gemini objection sits on the surface
r1's sol objection had just moved — i.e. a defect introduced by the previous round's fix, not a
surface holding the doc hostage. No localisation, so no SPLIT.

#### Objection 1 — `gpt5-6-sol` (reject)

> **strongest_objection**: "D3 makes both a closed and an open control mandatory for every
> successful invocation, but several acceptance fixtures contain no row of one or both required
> classes while still requiring rc=0. For example, AC-M1-3 contains six closed rows and no open
> row, AC-M1-4 contains only open rows, and AC-M1-6 contains only UNTAGGED rows. Those criteria
> cannot pass under the specified control-mismatch floor, so the design is internally
> unimplementable as written."
>
> **catch**: "Reconcile mandatory controls with AC-M1-3 through AC-M1-8 and any other fixture-based
> success arms. Do not add an undocumented test-only control bypass, because that would weaken the
> no-silent-fallback calibration rule."
>
> **proposed_fix**: "Add explicit control rows of the missing classes to every successful fixture
> and revise expected totals accordingly. For example: AC-M1-3 adds row 7 `[PARKED] control` and
> expects `6 closed / 7 rows (tagged-open 1, untagged 0)` with controls 1/7; AC-M1-4 adds row 4
> `[LANDED] control` and expects `1 closed / 4 rows (tagged-open 3, untagged 0)` with controls
> 4/1; AC-M1-5 uses LANDED, PARKED, and UNTAGGED rows and expects `1 closed / 3 rows (tagged-open
> 1, untagged 1)`; AC-M1-6 adds one LANDED and one PARKED control row and expects the two target
> rows to remain UNTAGGED within a four-row census. State the control numbers and revised totals
> for every remaining success arm."

**Controller adjudication: TRUE, verified first-party by reading D3 against §4.** D3 makes BOTH
`--control-closed` and `--control-open` mandatory (F5) and adjudicated (F6, exit 2 on mismatch),
while AC-M1-3's fixture is six closed rows with no open row, AC-M1-4's is three open rows with no
closed row, AC-M1-5's is one closed + one untagged, and AC-M1-6's is two untagged — each of which
makes one mandatory control unsatisfiable while the criterion demands `rc=0`. Unimplementable as
written, and the reviewer's own instruction ("do not add a test-only bypass") is the right
constraint: the fix is to the fixtures, never to the control contract.

**Applied VERBATIM**, all four named fixes exactly as written (AC-M1-3 row 7 `[PARKED] control`,
controls 1/7, `6 closed / 7 rows (tagged-open 1, untagged 0)`; AC-M1-4 row 4 `[LANDED] control`,
controls 4/1, `1 closed / 4 rows (tagged-open 3, untagged 0)`; AC-M1-5 LANDED+PARKED+UNTAGGED,
controls 1/2, `1 closed / 3 rows (tagged-open 1, untagged 1)`; AC-M1-6 plus one LANDED and one
PARKED control row, controls 1/2, `1 closed / 4 rows (tagged-open 1, untagged 2)`), **plus the
sweep the objection's last sentence demands** — control numbers and revised totals now stated for
every remaining success arm: AC-M1-2 `enumerate` (controls 1/2, `2 closed / 3 rows`), AC-M1-7
`body-heading` (1/2, `2 closed / 3 rows`), AC-M1-8 `preamble` (1/2, `1 closed / 2 rows`, both
fixtures), AC-M1-12 `heading-tail` (1/2, `1 closed / 2 rows`), AC-M1-15 `blob-oid` and
`blob-unavailable` (1/2 on a two-row fixture), and AC-M1-10's `control-mismatch` respecified to
pass BOTH flags on a three-row fixture in two runs. AC-M1-11 (`report`, 1/3) and AC-M1-13
(`snapshot`, 1/79) already carried satisfiable controls and are unchanged. No bypass was added.

**One controller clarification, labelled as such rather than smuggled in:** the five *failure* arms
of AC-M1-9 (`no-doc`, `no-heading`, `zero-rows`, `dup-number`, `gap`) pass no control flags and must
not — D4's floors are ordered and F0–F4 all precede F5, and each arm asserts its own floor's
message, which is what proves it reached that floor rather than falling through to `controls
required`. The objection is scoped to *success* arms; this states why the failure arms sit outside
it. It is the controller's sentence, not the reviewer's, and it resolves no objection.

#### Objection 2 — `gemini-3-1-pro` (reject)

> **strongest_objection**: "The citation format mandated by the Repo Profile rule in D6 and §8 M2
> structurally contradicts the actual output generated by the script in D5 and AC-M1-11. D6
> dictates the reading may be published 'ONLY as a census citation in this form: `census: 41
> closed / 89 rows (untagged 44) [queue_census.sh...`', which entirely omits the `tagged-open`
> count and the per-tag breakdown array. However, D5 specifies the script prints the rich format
> `census: 41 closed / 89 rows (tagged-open 4, untagged 44) [LANDED 40...] [...]` and instructs
> 'The last line is the only line a stamp may quote.' This forces an agent to either violate the
> strict Repo Profile rule by quoting the literal output, or manually mutilate the instrument's
> output to satisfy the rule."
>
> **catch**: "D6 and §8 M2 truncate the citation format, missing the `tagged-open` count and
> per-tag array `[LANDED X, ...]` that the script actually prints."
>
> **proposed_fix**: "Update the citation protocol string in D6 and §8 M2 to exactly match the
> script's output format: `census: N closed / M rows (tagged-open T, untagged U) [LANDED X, ROUTED
> Y, RULED OUT Z; PARKED P, IN-SPRINT Q, NEXT R] [queue_census.sh, charter-blob <oid>, controls
> a/b]`."

**Controller adjudication: TRUE, verified first-party by reading D5 and D6 side by side.** D5's
printed total line carries `(tagged-open 4, untagged 44)` and the per-tag array and closes with
*"The last line is the only line a stamp may quote"*; D6's mandated citation carried
`(untagged 44)` and no array. A stamp cannot satisfy both. Note the shape, because it is this
loop's own: **the r1 fix to D6 is what introduced it** — the citation string was rewritten for the
blob OID and its prefix was not re-checked against the line the script prints.

**Applied VERBATIM**: D6's and §8 M2's citation strings are now byte-for-byte the reviewer's
string, with the `8b15153` reading quoted underneath in that exact shape. AC-M2-6 is unchanged and
still holds — its `awk … /charter-blob/` needle is a substring of the new form — and AC-M1-11
already pins the printed line, so the rule and the instrument now share one literal.

**Not changed:** nothing else. No section either reviewer left unobjected was restructured; scope
is unchanged (normalisation still OUT, one tag position, UNTAGGED bucket, live controls 1/79).
Metered quorum spend across both rounds: **$0.241179**.
