# w-gate0-blind-to-its-own-crash-notices — a second, no-authority Gate-0 read for the driver's own crash notices, with the previous week's issue in scope and the notice's reassurance rewritten as a Gate-2 trigger

- Status: **planned** (design only; nothing implemented, no git write performed) · Date: **2026-09-08** (revised once 2026-09-08 — quorum r1: `gemini-3-1-pro` pass, `gpt5-6-sol` reject on D7's watermark one-liner, controller-measured six-state table V16; D7 respecified, D1–D6 and the milestone shape unchanged. **Quorum r2 BLOCKED and CLOSED under the narrow-refinement carve-out**: both reviewers present (`absent_reviewers` empty), `gemini-3-1-pro` **pass**, `gpt5-6-sol` **reject** on the SAME surface — F10's shape-only regex accepting `2026-99-99T99:99:99Z` and `9999-01-01T00:00:00Z`. The objection carries a concrete reviewer-authored `proposed_fix` and disputes no design direction, so the controller applied that fix VERBATIM as new §D7a + AC-M1-16/17 + five mutation rows + V17, after reproducing it first-party in both directions — the finding is BIGGER than filed: a shape-valid far-future watermark turns 4 real crash notices into 0 at rc=0, for ever. Artifacts `…-2026-09-08T06-53-44Z.json` (r1) and `…-2026-09-08T07-07-12Z.json` (r2).) · Designer: `claude:claude-fable-5-1` (iter-173) · Base commit measured at: **`0d316388b435e0455fbda97506ef6b994026816b`** (main checkout, `dev` == `origin/dev`, `git status --short` empty; charter blob `4e98fec339c581520bdc6c5c1ac4ea73f0e7fbf3`) · Owning queue row: **73 — `w-gate-0-cannot-see-the-drivers-own-crash-notices`** (clause-2, ~0.2d, gated on nothing, surfaced iter-151) · Scope class: **HARNESS — completeness of the loop's read of its own channel** · Verify profile: `ailang-code` (`bash ./scripts/verify_ail.sh`, `go vet ./...`, `go test ./... -count=1`; `scripts/verify_go.sh` is NOT a gate, charter row 76; `go build` is not a compile fence for `_test.go`).

---

## §1 Problem — the loop's own death notices sit in the channel it reads every fire, and the read is built to drop them

Row 73's premise, verbatim: the driver announces every failed iteration as a comment on the
bookkeeping issue; Gate 0.6 reads that issue through an author allowlist of `MarkEdmondson1234`;
the allowlist is **correct and must not be widened**, so the notice is invisible by construction.
Iteration 151 found the row-56 slot's notice by accident, at Gate 5, while moving a watermark.

**What is measured now, first-party at `0d31638` (every command in §7):**

- **The writer is content-stable and line-unstable.** The fleet driver posts, on a failed
  iteration, as the loop's own account, a body that BEGINS
  `⚠️ Mission iteration **FAILED to complete** (rc=$RC — timeout or crash) at <date>. Log on the rig: \`$LOG\`. The queue is untouched; the next interval will retry.`
  — at **`:1877`** in the pinned driver tree `48c4a6e49` this rig runs, at **`:1925`** in fleet
  HEAD `f7d04edb3` (the file changed at 08:27 today), and at `:617` in row 73's text, which is
  stale (V1). Nothing in this design cites a line number; the signature is anchored by content.
- **On the PREVIOUS bookkeeping issue `#107` (week of 2026-08-31): 43 comments, all 43 authored by
  `sunholo-voight-kampff`, and FOUR crash notices** — `2026-09-02T12:02:13Z` rc=143,
  `2026-09-02T15:21:07Z` rc=143, `2026-09-02T19:27:22Z` rc=143, `2026-09-07T02:47:28Z` rc=1 (V2).
  On the CURRENT issue `#129` (week of 2026-09-07): 14 comments, all self-authored, **zero** crash
  notices (V4).
- **A substring match over-counts by one, and the extra hit is the loop's own prose.** A search for
  `FAILED to complete` anywhere in the body returns **5** on `#107`; the fifth is the iteration-151
  *correction* comment at `2026-09-03T10:57:25Z`, which quotes the notice text at body offset
  **320**. Anchoring the signature at body offset **0** returns exactly **4** (V3). This is the
  row-57/row-72 lesson (*prose is never a tag*) arriving on a new surface: the loop writes about
  its own notices in the same channel, so the classifier must read a fixed position, never the
  body.
- **The notice's tail has already drifted; its head has not.** The three 09-02 bodies are 200
  characters; the 09-07 body is 284, because the driver gained the sentence *"Further identical
  failures are silent until the rc changes or an iteration completes."* between them (V3). A
  signature that included the tail would have gone dark on 2026-09-07 with no red anywhere. The
  design matches the prefix up to and including `(rc=`, captures the integer, and reads nothing
  after it.
- **Gate 0's read drops all of it, correctly.** `mission_directives.sh` (V1 checkout, absolute
  path — row 69's class) reports `0 directive(s) from [MarkEdmondson1234] … (of 43 comments)` for
  `#107` and `0 … (of 14 comments)` for `#129`; the positive control on `sunholo-data/ailang#972`
  surfaces `1` (`MarkEdmondson1234 @ 2026-09-02T07:17:34Z: D-54 b`) (V5). The zeros are true
  negatives; the allowlist is doing its job; the fix is a second read beside it, not a wider one.
- **Rotation week is the hard case and it has already happened.** The rc=1 notice landed on `#107`
  at `2026-09-07T02:47:28Z`; `#129` was created at `2026-09-07T07:51:05Z`, **5 h 04 m later**.
  `~/.ailang/state/mission-world-gh-issue` reads `129`, `…-gh-issue-prev` reads `107` (V6). A read
  of the current issue alone misses the last death of the previous week exactly when the loop
  rotates. The shared skill's Gate-5 rotation-week catch already says the *Mark-comment* read must
  also check `-prev`; this instrument reads both in one call, by design, not by reminder.
- **The false-reassurance half is the reason to ship.** The notice says *"The queue is untouched;
  the next interval will retry."* At iteration 151 that was true and misleading together: the
  queue was untouched **because** the dead iteration had already merged PR `#113` **61 seconds
  earlier** (mission log lines 984–985), so completed work read as unstarted and nothing retried
  it because nothing was left to retry (V7). The instrument's verdict line therefore says
  **A FIRE DIED — run Gate 2's died-mid-flight traces**, and is forbidden by an arm from containing
  the words `untouched`, `retry` or `nothing`.
- **A local corroborating record exists, but only since 2026-09-05.** The driver's
  `~/.ailang/state/mission-world-slot-verdicts.log` records `2026-09-07T02:47:27Z
  verdict=CRASHED_at=fired rc=1 … controller=pi:ollama/glm-5.3:cloud` — its first line is
  `2026-09-05T20:38:39Z`, so it did not exist for the three 09-02 deaths, and it is capped at 200
  rows (V8). It is a sibling signal, not a substitute (§9).

---

## §2 Design

Two scripts and two generated fixtures mirror the landed pairs `gate1_range_check.sh` /
`test_gate1_range_check.sh` and `queue_census.sh` / `test_queue_census.sh`: a bash-3.2 instrument
with a python3 JSON parser (the house pattern — `verify_ail.sh` and `gate1_range_check.sh` both
parse with python3, "no jq dependency", proven on `ubuntu-latest` by CI), a self-contained arm
suite with `--only <arm>`, one CI step in `go-verify`, and a Repo Profile rule.

### D1 — Interface, and the seam is the gh binary itself, never a bypass flag

```
scripts/gate0_self_notices.sh --issue <n> [--prev-issue <m>] --repo <owner/name> --self <login>
                              ( --since <ISO8601> | --watermark-file <path> [--watermark-file <path> …] )
                              --control <issue>:<count>
                              [--driver-src <file>] [--gh-bin <path>] [--gh-timeout <secs>]
```

**The injectable seam is `--gh-bin`** — the same seam `gate1_range_check.sh` ships and its suite
drives with a SHA-keyed stub. Every arm points `--gh-bin` at a stub script that maps
`issue view <n> --repo <repo> --json comments` to a fixture file and
`api repos/<repo>/issues/<n>` to a second fixture, and **records its argv** to a file the arm can
read. There is deliberately **no `--comments-file` flag**: a file-input bypass is a second code
path, and the live run would be the only caller of the first one — so the suite would prove the
parser and nothing about the read. With the stub seam there is ONE read path, the suite exercises
its exact argv shape (AC-M1-13 asserts the recorded argv), and the live invocation differs from the
tested one only in which binary answers. That is the "the seam exists AND the live path is still
exercised" property the brief demands, made structural rather than promised.

The instrument does not `cd` and derives nothing it is not handed: **`--issue` (digits), `--repo`,
`--self`, exactly one of `--since` | `--watermark-file`, and `--control` are all mandatory** (floor
F0/F6). The window has **no epoch default** — an epoch default would make every historical notice
"in window", exit 1 on every fire, and train the channel to be ignored (the exact failure Mark
named for the driver's own notices, 2026-08-31). `--since` is the literal form the suite drives;
**`--watermark-file`, repeatable, is the form the Repo Profile rule uses** — the instrument reads
the watermark files itself, validates each, takes the older, and reports every file's state on
stdout (D7), so the derivation is code the suite drives over fixtures rather than a shell line in a
charter bullet (quorum r1: both reviewers landed on that line, §7 V16).

### D2 — The signature is an anchored prefix over self-authored comments, and no body is ever printed

A comment is a **crash notice** iff:

1. `author.login`, lowercased, equals `--self` lowercased (GitHub logins are case-insensitive;
   `mission_directives.sh` lowercases for the same reason); and
2. `body` **starts at offset 0** with the literal
   `⚠️ Mission iteration **FAILED to complete** (rc=` followed by one or more ASCII digits.

The digits are captured as `rc=`. Nothing after them is read. Why offset 0 and not "contains":
V3 — the loop quotes its own notices in prose, at offset 320 on `#107`, and a substring rule
classifies that correction as a death. Why the prefix stops at `(rc=`: V3 again — the tail has
already drifted once (200 → 284 characters) with no fleet notice, and World cannot edit the writer.

**Foreign authors are not read at all** for classification: a stranger posting the exact
signature on the public issue produces `foreign` in the per-issue count and no `crash:` line. This
is not because a foreign hit would grant anything — a hit grants nothing (D6) — but because the
instrument's claim is *"the driver announced a death"*, and the driver posts as `--self`.

**The instrument prints no comment body, ever.** Its `crash:` lines carry `issue=`, `at=`, `rc=`,
`in-window=` and `url=`; the URL is how a controller reaches the text. Nothing an issue commenter
writes can therefore appear in the instrument's output and be mistaken, by a model reading a
transcript, for an instruction. AC-M1-7 plants a sentinel in a fixture body and asserts its
absence from stdout and stderr.

### D3 — Scope decision: crash notices are the ONLY classified signal; everything else self-authored is COUNTED, not classified

The evidence package names three more driver-authored notice classes on `#129` right now —
`Executor/planner lane degraded on this fire` (×3), `🔁 Controller model: X → Y` (×2),
`Driver ran UNPINNED on this fire` (×1) — the same blindness (self-authored, allowlist-hidden), a
different signal (conditions, not death). **This design does not classify them**, and the reasons
are argued rather than assumed:

- **They describe the fire the controller is IN; the crash notice describes a fire that VANISHED.**
  All three are posted *"before the iteration ran"*, and each has a **local, durable, already-read
  source** the controller does not need the issue for: the lane state in
  `~/.ailang/state/mission-world-lane-degraded.episode`, the model switch in
  `…/mission-world-model-last`, the pin failure in the notice spool
  (`…/mission-world-notice-spool.tsv`, which holds all three classes verbatim, V9), and the
  controller's own env. The crash notice's only other record was the `/tmp` driver log, which row
  71 measured a reboot destroying twice — the issue comment is the one durable copy, which is why
  its invisibility is a defect and theirs is a convenience.
- **A class vocabulary is a maintenance surface against a frozen writer.** Every notice class the
  fleet adds or rewords would need a World edit and a fixture, and the instrument would go quietly
  stale one class at a time. One anchored literal, defended by two controls (D5), is a surface this
  mission can keep honest at ~0.3d; four literals is a summary nobody re-measures.
- **The widening stays reviewable because the count is printed.** Every `issue:` line carries
  `other=N` — self-authored comments that are not crash notices. On `#129` today that is `14`
  (reports, corrections and the six condition notices together). A future row that wants
  conditions classified starts from a measured `other=` count and a stated reason, which is how
  scope should widen; it does not widen here by default.

So: `crash` is a class with a verdict; `other` is a bucket with a number; `foreign` is
`comments − self` and is never opened.

### D4 — Rotation week: the previous issue is read in the same call, and the arm proves the miss without it

`--prev-issue <m>` is optional at the interface and **mandatory in the Repo Profile rule** (D7): the
rule passes `$(cat ~/.ailang/state/mission-world-gh-issue-prev)` on every fire, not only "the first
fire after a rotation", because (a) the `-prev` file always names a real issue, (b) a second read
costs one bounded gh call, and (c) "first fire after rotation" is a condition a controller has to
notice, which is the class of thing this row exists to stop relying on. Both issues are read under
the same `--since`; `crash:` lines carry `issue=` so the reader knows which thread to open.

AC-M1-6's arm is two runs over the same stub: the rc=1 shape on issue `107` with `--since` before
it, once with `--issue 129` alone (exit **0** — the miss, demonstrated) and once with
`--issue 129 --prev-issue 107` (exit **1**, one `crash: issue=107 … in-window=yes` line). The
miss is asserted, not just the hit, so the arm reds if a future edit makes the current-issue read
accidentally cover the previous one (e.g. by reading a hard-coded list).

### D5 — Two controls and eight floors, one arm each, and how each arm reaches its floor

A signature is a literal against a writer this repo may not edit; a literal that silently stops
matching is a green gate over a dead read. Two controls make that loud:

- **`--control <issue>:<count>` (mandatory, F6/F7).** After the table is printed, the instrument
  asserts the named issue has exactly `<count>` crash notices **all-time** (not in-window) and
  prints `control: issue=<n> expect=<k> got=<g> ok` or `✗ CONTROL MISMATCH …` (exit 2). The
  published live control is **`107:4`**: `#107` is closed, the driver posts only to the current
  issue, and the four notices are the loop's own account's — so the count is fixed unless the
  loop itself begins a comment on `#107` with the signature (a controller must not, and the
  mismatch would be the control working). A control run through the SAME read path proves, on
  real data, that the read and the signature still agree with what the driver actually posted.
  The control issue is fetched by the same path (deduplicated if it is also `--prev-issue`).
- **`--driver-src <file>` (optional at the interface, mandatory in the rule, F8).** When given, the
  instrument asserts the literal prefix (D2) occurs in that file and prints
  `signature-source: <path> ok`; absent or unreadable → `✗ signature not found in driver source:
  <path>`, exit 2. When NOT given it prints the loud token **`signature-source: unchecked`** on
  its own line, so a transcript that omitted it is visibly weaker. The rule passes the pinned
  driver tree's `tools/launchd/mission-control.sh` by absolute path (V1: the line is at `:1877`
  there today; the check is a content grep, not a line). The two controls fail in opposite cases:
  `107:4` still passes after a fleet rewording (old notices keep the old text) while
  `--driver-src` reds — which is precisely the fire a controller must learn the wording changed.

Floors protect the instrument, not the input; each has its own message and its own arm, and each
floor arm asserts its **specific** message so that a floor reached by falling through to a
different one reds on message mismatch (the iter-108 defence, as `test_queue_census.sh` does it):

| # | floor | message (stderr) | exit |
|---|---|---|---|
| F0 | a mandatory argument is absent or malformed (`--issue`/`--prev-issue` not digits; both or neither of `--since`/`--watermark-file`) | `✗ usage: --issue <digits>, --repo, --self and exactly one of --since | --watermark-file are required` | 2 |
| F1 | gh binary unreachable (child rc 127) | `✗ gh binary unreachable: <path>` | 2 |
| F2 | gh call exceeded `--gh-timeout` (default 5 s; `run_bounded`, the gate1 wrapper — no `timeout(1)` on this rig) | `✗ gh timed out after <t>s for issue #<n>` | 2 |
| F3 | gh exited non-zero | `✗ gh failed (rc=<rc>) for issue #<n>` | 2 |
| F4 | payload is not JSON, has no `comments` array, or an element lacks `author.login`/`createdAt`/`body` | `✗ malformed comments payload for issue #<n>: <reason>` | 2 |
| F5 | `comments` array length ≠ the REST issue's `comments` count | `✗ comment array truncated for issue #<n>: got <x> of <y>` | 2 |
| F6 | `--control` absent | `✗ control required: --control <issue>:<count>` | 2 |
| F7 | control issue has a different all-time crash count | `✗ CONTROL MISMATCH: issue <n> expect=<k> got=<g>` (table printed first) | 2 |
| F8 | `--driver-src` given and the signature prefix is not in it | `✗ signature not found in driver source: <path>` | 2 |
| F9 | every `--watermark-file` is absent (no readable window at all) | `✗ watermark channel unavailable: no readable watermark among <n> file(s): <paths>` | 2 |
| F10 | a present watermark (or the `--since` literal) is empty, unreadable, or fails **strict UTC timestamp validation** — not `YYYY-MM-DDTHH:MM:SSZ` in shape, OR not round-trip canonical, OR later than now+skew | `✗ watermark invalid: <path or --since> (<empty \| unreadable \| not-ISO8601 '<value>' \| not-canonical '<value>' → '<canonical>' \| future '<value>' > now+<skew>>)` | 2 |

F0, F6, F8, F9 and F10 fire **before any network call** (they are argument/file checks, in that
order); F1–F5 fire per issue read; F7 fires after the table. F9/F10 are D7's watermark
dispositions, specified there. **F5 exists because the `gh issue view --json comments`
comment cap is UNVERIFIED** (§9): the largest issue on this rig has 84 comments and `view` agreed
with REST on every one tried (V10), so the cap could not be measured, and the design converts the
unknown into a fail-loud consistency check — one extra bounded call per issue
(`gh api repos/<repo>/issues/<n>`, parsed for `.comments`).

The suite has **no helper that synthesises a comment, an author or a body**: the only fixture
builder is `make_json <path>` reading literal lines from a heredoc, and each floor arm writes its
broken shape by hand (a payload with no `comments` key; a `.meta.json` whose count disagrees; a
driver file without the signature).

### D6 — Output contract, exit codes, and the sentence that replaces "the queue is untouched"

```
gate0_self_notices: repo=sunholo-data/ailang-world self=sunholo-voight-kampff issues=129,107
watermark: /Users/voightkampff/.ailang/state/mission-world-last-seen = 2026-09-08T06:19:00Z ok
watermark: /Users/voightkampff/.ailang/state/mission-129-last-seen = 2026-09-08T06:19:00Z ok
since: 2026-09-08T06:19:00Z (older of 2 readable watermark(s) of 2)
crash: issue=107 at=2026-09-02T12:02:13Z rc=143 in-window=no url=https://github.com/sunholo-data/ailang-world/issues/107#issuecomment-<id>
crash: issue=107 at=2026-09-02T15:21:07Z rc=143 in-window=no url=…
crash: issue=107 at=2026-09-02T19:27:22Z rc=143 in-window=no url=…
crash: issue=107 at=2026-09-07T02:47:28Z rc=1 in-window=no url=…
issue: 129 comments=14 self=14 crash=0 other=14
issue: 107 comments=43 self=43 crash=4 other=39
control: issue=107 expect=4 got=4 ok
signature-source: /Users/voightkampff/dev/sunholo-data/ailang/tools/launchd/mission-control.sh ok
verdict: 0 crash notice(s) since 2026-09-08T06:19:00Z on #129,#107 — no death signal from the driver; Gate 2's traces (a)-(c) still run [gate0_self_notices.sh, control 107:4, watermarks 2/2]
```

The block above is the **expected live reading at landing** given V2/V4/V6 and today's watermarks
(both `2026-09-08T06:19:00Z`, V6); the controller records the actual one in §7 (AC-M2-3). One
`watermark:` line per `--watermark-file` (or none under `--since`), then one `since:` line naming
the value used and how many files it was the older of. `crash:` lines print **all-time** notices
per issue, sorted by `at`, with `in-window=yes|no` decided by `createdAt > since` (strict, as
`mission_directives.sh` compares) — so the week's death history is visible on every fire even when
the window is empty, and the count that drives the exit code is the `yes` rows. `issue:` lines
print in the order the issues were passed. On a hit the last line is:

```
verdict: 1 crash notice(s) since <since> on #129,#107 — A FIRE DIED (newest rc=143 at 2026-09-02T19:27:22Z on #107): run Gate 2's died-mid-flight traces (a) open PRs (b) worktrees (c) uncommitted state BEFORE picking, and credit the orphan in the log; "the queue is untouched" describes the charter, not the work [gate0_self_notices.sh, control 107:4, watermarks 2/2]
```

The trailing `watermarks a/b` is `readable/given` under `--watermark-file` and the literal
`watermarks literal` under `--since`, so a stamp that quotes the verdict also quotes whether the
window came from two files, one, or a typed value. Exit codes: **0** no in-window crash notice,
controls ok; **1** one or more in-window crash notices (*a controller must look* — the gate1
convention; never "abort"); **2** any floor. **The verdict line is the only line a STATUS stamp
may quote**, and exactly one arm (`report`) owns its literal shape — including that on exit 1 it
contains `A FIRE DIED` and `died-mid-flight` and contains none of `untouched`, `retry`, `nothing
happened`.

**No authority, stated as mechanism rather than as a wish:** the instrument reads only
self-authored comments for classification, prints no body, cannot change any file, and its
positive outcome is an exit code whose prescribed consequence (D7) is *to run three read-only
Gate-2 traces and write a log credit*. It cannot name a queue row, cannot unpark, and the Repo
Profile rule says in words that a hit is evidence a fire died and never a pick.

### D7 — Where it runs: a Repo Profile rule at Gate 0, beside the directive read, quoting the exact command

World may not edit the shared skill (`~/.claude/skills/mission-control` is a symlink into the V1
checkout, V11), so the Gate-0 rule lands as a Repo Profile bullet, in the same section and shape
as the TWO-WATERMARKS bullet (iter-52) and the ABSOLUTE-PATH bullet (iter-72). The bullet opens
with the bold literal **`GATE 0 HAS A SECOND, NO-AUTHORITY READ FOR THE DRIVER'S OWN CRASH
NOTICES (process fix, iter-173)`** — AC-M2-1 scopes its no-redirect/no-epoch grep to the bullet by
that opening. Its content, which M2 pastes and AC-M2-1/2 pin:

- **Immediately after Gate 0.6's `mission_directives.sh` call, run the second read.** The rule
  types no watermark logic: it names the two files the iter-52 rule names, and the instrument
  reads, validates and reports them. No stderr redirect and no epoch literal anywhere in the
  command — nothing needs suppressing, because every state of every file is a line on stdout or
  a named floor:
  ```bash
  ISSUE="$(cat ~/.ailang/state/mission-world-gh-issue)"; PREV="$(cat ~/.ailang/state/mission-world-gh-issue-prev)"
  bash scripts/gate0_self_notices.sh --issue "$ISSUE" --prev-issue "$PREV" --repo sunholo-data/ailang-world \
    --self sunholo-voight-kampff \
    --watermark-file ~/.ailang/state/mission-world-last-seen \
    --watermark-file "$HOME/.ailang/state/mission-${ISSUE}-last-seen" \
    --control 107:4 --driver-src /Users/voightkampff/dev/sunholo-data/ailang/tools/launchd/mission-control.sh; echo "rc=$?"
  ```
  (an absent issue file makes `cat` print its error AND hands the instrument an empty `--issue`,
  which F0 refuses by name — two loud signals, no fallback).

**The watermark disposition — every state the controller measured (V16) has a defined, loud
outcome.** The instrument classifies each `--watermark-file` as `ok` (present, readable, non-empty,
passing **strict UTC timestamp validation**, D7a), `ABSENT` (no such path), or `invalid`
(present but empty, unreadable, or failing that validation). It prints one `watermark:` line per
file, then decides:

| state of the two files | stdout / stderr | since | rc |
|---|---|---|---|
| both `ok` | two `watermark: … ok`; `since: <older> (older of 2 readable watermark(s) of 2)` | older of two | unchanged (0/1 by notices) |
| exactly one `ok`, the other `ABSENT` | `watermark: <path> ABSENT — DEGRADED: single-watermark read (iter-52 rule wants two)`; `since: <value> (older of 1 readable watermark(s) of 2)`; verdict suffix `watermarks 1/2` | the one value | **unchanged, DEGRADED named** — see below |
| both `ABSENT` | F9 `✗ watermark channel unavailable: no readable watermark among 2 file(s): <paths>` | none | **2** |
| any file present but **empty** | F10 `✗ watermark invalid: <path> (empty)` | none | **2** |
| any file present but **not ISO8601** (`not-a-date`) | F10 `✗ watermark invalid: <path> (not-ISO8601 'not-a-date')` | none | **2** |
| any file present but **calendar-invalid** (`2026-02-30T00:00:00Z`) | F10 `✗ watermark invalid: <path> (not-canonical '2026-02-30T00:00:00Z' → '2026-03-02T00:00:00Z')` | none | **2** |
| any file present but **shape-valid garbage** (`2026-99-99T99:99:99Z`) | F10 `✗ watermark invalid: <path> (not-canonical '2026-99-99T99:99:99Z')` | none | **2** |
| any file present but **far-future** (`9999-01-01T00:00:00Z`) | F10 `✗ watermark invalid: <path> (future '9999-01-01T00:00:00Z' > now+300s)` | none | **2** |
| any file present but **unreadable** | F10 `✗ watermark invalid: <path> (unreadable)` | none | **2** |

"Older" is the **lexicographic minimum, and it is used ONLY after strict UTC timestamp
validation** — which is exactly the step the r0 one-liner lacked, so a malformed value there lost
the sort and vanished (V16, state 6). D7a is evaluated for every present file before any minimum is
taken; a malformed file is never outvoted. The same validation applies to the `--since` literal.

### D7a — Strict UTC timestamp validation (round-2 quorum fix, applied VERBATIM from the reviewer)

Round 2 blocked on this surface. `gpt5-6-sol`'s objection, verbatim: *"F10 does not actually
validate ISO-8601 timestamps: it accepts any fixed-width shape, including `2026-99-99T99:99:99Z`
and valid-but-implausible future values such as `9999-01-01T00:00:00Z`. Such a watermark can pass
as `ok`, become the selected `since`, suppress every crash notice, and yield exit 0 indefinitely.
This violates the fail-loud/no-silent-fallback requirement."* Its `catch`, verbatim: *"The claim
that the regex has 'forced every survivor into' a valid chronological timestamp is unverified and
false. AC-M1-15 tests only `not-a-date`; it does not test calendar-invalid, future, or
clock-skewed values."* Both are correct; see V17 for the controller's first-party reproduction,
including the consequence measurement (4 real notices → 0, at rc=0).

Its `proposed_fix`, applied VERBATIM as this section's requirement: *"Replace shape-only validation
with strict parsing and round-trip canonicalization, then reject timestamps later than current UTC
plus a documented small skew allowance. Add F10 cases for `2026-02-30T00:00:00Z`,
`2026-99-99T99:99:99Z`, and a far-future canonical timestamp, each requiring rc=2 before any gh
call. Revise D7 to say lexicographic comparison is used only after strict UTC timestamp
validation, and add a verification-log row demonstrating all three failures."*

So validation is **three ordered checks**, all before any network call, each with its own F10 reason:

1. **Shape** — `^[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}Z$`, written expanded for
   bash 3.2 (no `{n}`). Reason `not-ISO8601 '<value>'`.
2. **Round-trip canonicalization** — parse as UTC and re-format with the same layout; the result
   MUST be byte-identical to the input. Reason `not-canonical '<value>' → '<canonical>'` (the arrow
   half is omitted when the parse itself fails). This is what rejects calendar-invalid values,
   which shape alone accepts: measured on the rig, `2026-02-30T00:00:00Z` normalises to
   `2026-03-02T00:00:00Z` and `2026-02-29T00:00:00Z` (2026 is not a leap year) to
   `2026-03-01T00:00:00Z`, while `2026-99-99T99:99:99Z` fails the parse outright (V17).
3. **Future ceiling** — reject any value later than `now(UTC) + SKEW`, with **SKEW = 300 s**, a
   documented small allowance for rig/GitHub clock drift. Reason `future '<value>' > now+300s`.
   Round-tripping alone cannot catch this: `9999-01-01T00:00:00Z` is perfectly canonical (V17).
   The ceiling is the only check that stops the failure the objection names — a watermark in the
   future suppresses every notice for ever, silently, at rc=0.

**Portability is an implementation constraint, not a design change (controller note, labelled as
such).** The rig is macOS/BSD `date` and CI is ubuntu/GNU `date`, and their strict-parse invocations
differ: BSD needs `TZ=UTC date -j -f '%Y-%m-%dT%H:%M:%SZ' "$v" +'%Y-%m-%dT%H:%M:%SZ'`, GNU needs
`date -u -d "$v" +'%Y-%m-%dT%H:%M:%SZ'`. The BSD arm is measured working for all five cases above
(V17); the GNU arm is **UNMEASURED on this rig** (no `gdate` installed) and is therefore proven by
CI, not by assertion — which is exactly what AC-M1-16 is for. Implement it as a two-arm helper that
tries one form and falls back to the other, and refuse loudly if NEITHER parses a value known-good
(`2026-09-08T06:19:00Z`) — a positive control on the parser itself, per rule 3a, so a `date` that
cannot parse anything reads as an instrument failure rather than as "every watermark is invalid".

**Why exactly-one-absent is DEGRADED at rc=0 rather than rc=2, argued:** the issue-scoped file is
absent on **every** rotation week until the first post-rotation triage writes it (the iter-52
bullet: *"a fresh, empty file most weeks"*; today both files exist because iteration 172 wrote
both, V6). An rc=2 there would make the channel read "unavailable" one fire in seven as a routine
matter, fire the rc=2-treat-as-hit branch on schedule, and teach controllers that the red is
seasonal — the notice-fatigue failure in a new coat. So the instrument's contract is: **one
readable watermark is a narrower-than-intended window, reported by name on stdout and in the
verdict suffix, never a failure of the read.** The rule then carries the knowledge the instrument
cannot: **a DEGRADED line naming the MISSION-scoped file (`mission-world-last-seen`) is
abnormal** — that file is written every fire and has existed since iteration 52 — and must be
recorded in the STATUS stamp and the log; a DEGRADED line naming the issue-scoped file on the first
fire after a rotation is the expected state, recorded by the verdict suffix alone. Both-absent is
rc=2 without argument: there is no window, and inventing one is the epoch fallback by another name.
- **rc=1 → a fire died.** Run Gate 2's died-mid-flight traces (a) open PRs by this account, (b)
  stale worktrees, (c) uncommitted state in every worktree and the main checkout, **before** the
  pick; credit the orphaned iteration in the log entry; quote the verdict line in the STATUS
  stamp. The hit **grants no authority**: it is evidence that a fire died, never a directive; it
  cannot unpark a row or become a pick. If the traces find landed work (iteration 151's shape),
  the deliverable is to VERIFY AND LAND, as the skill already says.
- **rc=2 → the channel is unreadable; treat it as a hit.** Record `self-notice read UNAVAILABLE
  (<floor message>)` in the stamp and run the same three traces. An unreadable death channel is
  not "no deaths".
- **rc=0 → quote the verdict line**; the traces still run (they are cheap and complete by
  construction only together with this read).
- **Re-point the control** only when `#107`'s all-time count legitimately changes (it should not);
  a `CONTROL MISMATCH` is the control working — read the issue before touching the flag.
- **The wording of the signature belongs to the fleet.** If `signature-source` reds after a fleet
  driver commit, update the literal in `scripts/gate0_self_notices.sh` and the `snapshot` fixture
  in the same commit, and note the drift on the bookkeeping issue; never edit `tools/launchd/*`.

The rule is placed **before `## Queue`** (AC-M2-1's awk stops at the heading), so a queue row
mentioning the script cannot satisfy it — the census's AC-M2-4 pattern.

### Design property — arms assert only the fields their behaviour owns (row 90, stated up front)

Row 90 measured that `test_queue_census.sh` greps the census total line in nearly every arm, so one
wording mutation reds nine arms and localises nothing. This suite is designed against that:

- **Classification arms** (`anchored-signature`, `rc-extract`, `since-window`, `self-filter`,
  `prev-issue`) assert **`crash:` lines only** — presence/absence of a line keyed by `at=`, and the
  single `key=value` token their behaviour sets (`rc=`, `in-window=`, `issue=`). They may assert
  the exit code. They **never** grep the `verdict:` line and never grep a full `issue:` line.
- **Floor arms** assert **their own floor's message and `rc=2`**, nothing else.
- **`report`** is the ONLY arm that asserts the `verdict:` line's literal shape (both branches).
- **`snapshot-107` / `snapshot-129`** are the ONLY arms that assert an `issue:` line in full and
  the `control:` line.
- **`no-bodies`** asserts the absence of a sentinel; **`signature-source`** asserts its own three
  lines.
- **Watermark arms** (`watermark-older`, `watermark-degraded`, `watermark-unavailable`,
  `watermark-invalid`) assert only **`watermark:` lines, the `since:` line and their floor's
  message**; the `since:` line is the field they jointly own. They never grep `crash:`, `issue:`
  or `verdict:` lines — `report` owns the `watermarks a/b` suffix along with the rest of that line.

§5's "kills which arm" column is written to this property, and AC-M1-14's landing step is a
first-party drill of the `report` mutation confirming the red set is **`report` alone**.

---

## §3 Files to create / modify

- `scripts/gate0_self_notices.sh` — **new** (M1). D1–D6. bash 3.2 (`set -uo pipefail`, no
  associative arrays), `run_bounded` copied in shape from `gate1_range_check.sh`, ONE python3
  parser invoked per issue payload that emits `CRASH <at> <rc> <url>` / `COUNT <comments> <self>`
  lines the bash wrapper assembles; the wrapper owns argument parsing, floors, the control, the
  verdict and exit codes.
- `scripts/test_gate0_self_notices.sh` — **new** (M1). 22 arms, `--only <arm>`, `ok`/`not ok`
  lines, final `N passed, M failed`, non-zero on any failure — the precedent shape. Scratch lives
  in the worktree (`mktemp -d "$(pwd)/.g0_scratch.XXXXXX"`, row 71), cleaned on exit; watermark
  fixtures are literal files written there by each arm (`printf` of one value, `: >` for empty,
  `chmod 000` for unreadable). Stub gh
  built by `make_stub_gh <out> <dir> <argv-log>`: serves `<dir>/<n>.json` for `issue view <n>` and
  `<dir>/<n>.meta.json` for `api repos/*/issues/<n>`, appends its argv to `<argv-log>`, and
  `exit 1` with `stub gh: no fixture` for anything else. Every instrument invocation is an explicit
  `bash "$SCRIPT_UT" …` — never the login shell (zsh), never `sh`.
- `scripts/testdata/gate0_self_notices_107.json` and `…_129.json` — **new** (M1). Frozen, cutoff-
  pinned, body-truncated snapshots of the two issues, **generated by this command, checked in as
  produced**:
  ```bash
  gen() { gh issue view "$1" --repo sunholo-data/ailang-world --json comments \
    | jq --arg cut "$2" '{comments: [.comments[] | select(.createdAt <= $cut) | {author: {login: .author.login}, createdAt, url, body: .body[0:700]}]}'; }
  gen 107 2026-09-07T23:59:59Z > scripts/testdata/gate0_self_notices_107.json
  gen 129 2026-09-08T06:28:52Z > scripts/testdata/gate0_self_notices_129.json
  ```
  Dry-run at design time (V12): `107` → **36,718 bytes**, sha256
  `79ad0ad78ea7be1548c4fe4049888a60e6b1a354b96a3d94573f41f994406e2f`, 43 comments, last
  `2026-09-07T07:51:27Z`; `129` → **12,624 bytes**, sha256
  `d61af56d9826e85d7f3c9c95a03dd6dbb1663e95df06534cd012b3383446cf14`, 14 comments, last
  `2026-09-08T06:28:52Z`. Three parts are load-bearing, each bought by a prior red: the
  **`createdAt <= $cut`** filter makes the file reproducible after the live issue grows (`#129` is
  live; `#107` is closed but closed issues accept comments); **`.body[0:700]`** truncates by
  **codepoint** — jq strings are Unicode, measured on `a—b—c` (V12) — so no multi-byte character is
  cut in half (row 72's `substr`-counts-bytes red); and 700 is chosen because the iter-151
  correction carries the quoted signature at offset **320** (V3), so the truncated fixture still
  contains the negative control that a substring classifier fails. The four crash bodies (200,
  200, 200, 284 chars) survive whole. Alongside each, a hand-written
  `scripts/testdata/gate0_self_notices_<n>.meta.json` of the form `{"comments": 43}` /
  `{"comments": 14}` stands in for the REST count (F5); `gh api repos/sunholo-data/ailang-world/issues/107 --jq .comments` → `43` is the value to write (V10's method).
- `.github/workflows/ci.yml` — **modify** (M1): one additive step `Gate-0 self-notice instrument
  suite` in `go-verify` immediately after `Queue census (live charter, controlled)` (line 210–212,
  V13). **No live step**: CI has no `gh` auth and the live read is a controller-run Gate-0 step.
- `design_docs/world-mission.md` — **modify** (M2 only): one Repo Profile bullet (D7), placed
  after the `mission_directives.sh` ABSOLUTE-PATH bullet (line 122, V14). No queue row body is
  touched by this item; the controller's landing commit tags row 73 in the D1-of-row-72 leading
  position.
- `design_docs/planned/w-gate0-blind-to-its-own-crash-notices.md` — this document (controller
  moves it to `implemented/` at landing).
- **Nothing under `tools/launchd/`** (frozen core, `D-WORLD-DRIVER-1`); nothing in the shared skill.

---

## §4 Acceptance criteria

Every criterion is a command run from the repo root **under `bash`**, with its expected exit code
and output. All were baselined on the pristine tree at `0d31638` and are RED there (§7 B1–B6);
`grep` is `/usr/bin/grep` throughout (the login shell's `grep` is ugrep).

**M1 — instrument + suite + fixtures + suite CI step**

- **AC-M1-1** — both scripts exist and are executable.
  `test -x scripts/gate0_self_notices.sh && test -x scripts/test_gate0_self_notices.sh && echo EXISTS` → `EXISTS`, rc=0. (Base: rc=1, B1.)
- **AC-M1-2** — the signature is anchored at offset 0 and matches the prefix only. Arm
  `anchored-signature`: five self-authored bodies — (i) the exact 09-02 body, (ii) the 09-07 body
  with the drifted tail, (iii) one space then the signature, (iv) the signature quoted at offset 40
  inside prose, (v) `…(rc=abc — …` (non-digit rc) — `--since 2000-01-01T00:00:00Z`, matching
  control → exactly **two** `crash:` lines, keyed by (i)'s and (ii)'s `at=`, none by (iii)–(v); rc=1.
  `bash scripts/test_gate0_self_notices.sh --only anchored-signature; echo "rc=$?"` → `ok …` lines, `rc=0`. (Base: rc=127, B2.)
- **AC-M1-3** — `rc=` is the captured integer. Arm `rc-extract`: notices `(rc=143 —` and `(rc=1 —`
  → `crash: … rc=143` and `crash: … rc=1` keyed by their `at=`, neither carrying the other's value.
  `bash scripts/test_gate0_self_notices.sh --only rc-extract; echo "rc=$?"` → `rc=0`.
- **AC-M1-4** — the window is `createdAt > since`, strict, and only in-window rows drive the exit
  code. Arm `since-window`: notices at `T1 < T2 < T3`; run A `--since T2` → `T1`,`T2` `in-window=no`,
  `T3` `in-window=yes`, rc=1; run B `--since T3` → all `no`, rc=0; all three `crash:` lines printed
  in both runs.
  `bash scripts/test_gate0_self_notices.sh --only since-window; echo "rc=$?"` → `rc=0`.
- **AC-M1-5** — only `--self` comments classify, case-insensitively. Arm `self-filter`: the exact
  signature by `SUNHOLO-VOIGHT-KAMPFF`, by `MarkEdmondson1234` and by `stranger-42`, `--self
  sunholo-voight-kampff` → exactly one `crash:` line, keyed by the uppercase author's `at=`; rc=1.
  `bash scripts/test_gate0_self_notices.sh --only self-filter; echo "rc=$?"` → `rc=0`.
- **AC-M1-6** — rotation week: the previous issue is read in the same call, and the miss without
  it is demonstrated. Arm `prev-issue`: stub serves `129` (no notices) and `107` (one notice at
  `2026-09-07T02:47:28Z`), `--since 2026-09-07T00:00:00Z`, control `107:1`. The control read
  fetches `107` in both runs, but **classification is scoped to `--issue`/`--prev-issue` only** —
  a control issue that is not also named as an issue contributes to `control:` and to nothing
  else. Run A `--issue 129` → **rc=0**, no `crash: issue=107 … in-window=yes` line, and the
  verdict names `on #129` only. Run B `--issue 129 --prev-issue 107` → exactly one line
  `crash: issue=107 at=2026-09-07T02:47:28Z rc=1 in-window=yes`, **rc=1**. Both runs are
  asserted — the miss as well as the hit — so the arm reds if a future edit makes the
  current-issue read cover the previous one by accident (a hard-coded list, or the control read
  leaking into classification).
  `bash scripts/test_gate0_self_notices.sh --only prev-issue; echo "rc=$?"` → `rc=0`.
- **AC-M1-7** — no comment body is ever printed. Arm `no-bodies`: a notice whose tail carries
  `SENTINEL-DO-NOT-PRINT-7f3a` and an `other` comment that is only that sentinel → combined
  stdout+stderr contains the sentinel **zero** times and does contain the notice's `url=`.
  `bash scripts/test_gate0_self_notices.sh --only no-bodies; echo "rc=$?"` → `rc=0`.
- **AC-M1-8** — floors F0–F5, one arm each, each asserting its own message and rc=2:
  `for a in usage gh-missing gh-timeout gh-failed malformed truncated; do bash scripts/test_gate0_self_notices.sh --only $a || exit 1; done; echo "rc=$?"` → six `ok …` blocks, `rc=0`.
  Shapes: `usage` is three runs — omit both `--since` and `--watermark-file`; pass both; pass
  `--issue ""` — each asserting F0's message; `gh-missing` passes `--gh-bin /nonexistent/gh`; `gh-timeout`
  uses a stub that sleeps 3 s with `--gh-timeout 0.5`; `gh-failed` uses a stub that `exit 7`;
  `malformed` serves `{"nope": []}`; `truncated` serves a 2-comment payload with a `.meta.json` of
  `{"comments": 5}`. These arms pass a `--control` where the floor is downstream of F6 and omit it
  only in `usage` (F0 precedes F6 in the ordered floors, and the arm asserts F0's message, which
  is what proves it reached F0 rather than F6).
- **AC-M1-9** — the control is mandatory (F6) and adjudicated (F7), table first. Arm
  `control-required`: valid stub, no `--control` → `control required`, rc=2. Arm
  `control-mismatch`: a stub with two notices on `107`, `--control 107:4` → `✗ CONTROL MISMATCH:
  issue 107 expect=4 got=2`, rc=2, **and both `crash:` lines are printed above it**.
  `bash scripts/test_gate0_self_notices.sh --only control-required; bash scripts/test_gate0_self_notices.sh --only control-mismatch; echo "rc=$?"` → `rc=0`.
- **AC-M1-10** — the driver-source control (F8) in all three states. Arm `signature-source`: run A
  with `--driver-src` pointing at a heredoc file containing the signature line as the driver
  writes it (with `$RC` unexpanded) → `signature-source: <path> ok`, rc unchanged by it; run B
  pointing at a file that lacks it → `✗ signature not found in driver source: <path>`, rc=2,
  **before any gh call** (the stub's argv log is empty); run C without the flag →
  `signature-source: unchecked` present.
  `bash scripts/test_gate0_self_notices.sh --only signature-source; echo "rc=$?"` → `rc=0`.
- **AC-M1-11** — the report format (D6), both branches, and the forbidden words. Arm `report`:
  run A (no in-window notice) → last line matches exactly
  `^verdict: 0 crash notice\(s\) since <since> on #129,#107 — no death signal from the driver; Gate 2's traces \(a\)-\(c\) still run \[gate0_self_notices\.sh, control 107:4, watermarks literal\]$`, rc=0;
  run B (one in-window) → last line matches
  `^verdict: 1 crash notice\(s\) since <since> on #129,#107 — A FIRE DIED \(newest rc=143 at <at> on #107\): run Gate 2's died-mid-flight traces \(a\) open PRs \(b\) worktrees \(c\) uncommitted state BEFORE picking, and credit the orphan in the log; "the queue is untouched" describes the charter, not the work \[gate0_self_notices\.sh, control 107:4, watermarks literal\]$`, rc=1;
  run C (as A but with two `ok` `--watermark-file`s) → last line ends `, watermarks 2/2\]$`;
  and run B's output contains none of `untouched;`, `will retry`, `nothing happened` (the quoted
  phrase *"the queue is untouched" describes the charter* is the only permitted occurrence of
  `untouched`, hence the `;`-suffixed needle).
  `bash scripts/test_gate0_self_notices.sh --only report; echo "rc=$?"` → `rc=0`.
- **AC-M1-12** — the frozen snapshots reproduce the live readings. Arm `snapshot-107`: the checked-in
  `_107.json` + `_107.meta.json`, `--issue 107 --since 2026-09-08T00:00:00Z --control 107:4` →
  exactly four `crash:` lines with `at=` `2026-09-02T12:02:13Z`/`rc=143`,
  `2026-09-02T15:21:07Z`/`rc=143`, `2026-09-02T19:27:22Z`/`rc=143`, `2026-09-07T02:47:28Z`/`rc=1`,
  all `in-window=no`; **no `crash:` line with `at=2026-09-03T10:57:25Z`** (the quoting
  correction); `issue: 107 comments=43 self=43 crash=4 other=39`; `control: issue=107 expect=4
  got=4 ok`; rc=0. Arm `snapshot-129`: `_129.json` + `_129.meta.json`, `--issue 129 --prev-issue
  107 --since 2026-09-08T00:00:00Z --control 107:4` → `issue: 129 comments=14 self=14 crash=0
  other=14`, no `crash: issue=129` line, rc=0.
  `bash scripts/test_gate0_self_notices.sh --only snapshot-107; bash scripts/test_gate0_self_notices.sh --only snapshot-129; echo "rc=$?"` → `rc=0`. (Base: fixtures absent, B3.)
- **AC-M1-13** — the seam is the ONLY read path, and its argv shape is what the live run sends.
  Instrument-health control (labelled, not load-bearing): `/usr/bin/grep -cE '(^|[^"$])gh (issue|api) ' scripts/gate0_self_notices.sh` → `0`. Load-bearing: `snapshot-107` reads the stub's argv log and asserts exactly `issue view 107 --repo sunholo-data/ailang-world --json comments` and `api repos/sunholo-data/ailang-world/issues/107` (order-insensitive).
  `bash scripts/test_gate0_self_notices.sh --only snapshot-107 | /usr/bin/grep -c 'argv'` → `2`.
- **AC-M1-14** — the whole suite is green, wired into CI, and the row-90 property holds under a
  first-party drill. `bash scripts/test_gate0_self_notices.sh; echo "rc=$?"` → `… passed, 0
  failed`, rc=0; `/usr/bin/grep -n 'test_gate0_self_notices' .github/workflows/ci.yml` → one
  `run:` line inside `go-verify` (base rc=1, B4). Drill at landing (§7): `cp` backup, mutate the
  verdict wording (`crash notice(s)` → `crash notices`), suite → red set **`report` only** (distinct
  arm names from `^not ok` lines → 1); restore by `cp` + `shasum -a 256 -c`, never `git checkout --`
  (row 82).
- **AC-M1-15** — the watermark dispositions (D7 table), four arms, each asserting only its
  `watermark:`/`since:` lines or floor message and the exit code. `watermark-older`: two `ok` files
  `2026-09-08T06:19:00Z` and `2026-09-05T00:00:00Z` → `since: 2026-09-05T00:00:00Z (older of 2
  readable watermark(s) of 2)`, both `watermark: … ok`, rc=0 (fixture has no in-window notice).
  `watermark-degraded`: one `ok` file, one absent path → `watermark: <absent> ABSENT — DEGRADED:
  single-watermark read`, `since: <value> (older of 1 readable watermark(s) of 2)`, rc=0.
  `watermark-unavailable`: two absent paths → F9's message naming both paths, rc=2, stub argv log
  empty. `watermark-invalid`: four runs — an empty file (`(empty)`), a file holding `not-a-date`
  beside an `ok` file (`(not-ISO8601 'not-a-date')` — the malformed file must NOT be outvoted),
  a `chmod 000` file (`(unreadable)`; the arm first asserts `[ ! -r ]` holds after the chmod and
  reds with `precondition: still readable — running as root?` if not, so the branch is never
  silently skipped), and `--since not-a-date` (`--since (not-ISO8601 'not-a-date')`) — each rc=2
  with its message.
  `for a in watermark-older watermark-degraded watermark-unavailable watermark-invalid; do bash scripts/test_gate0_self_notices.sh --only $a || exit 1; done; echo "rc=$?"` → four `ok …` blocks, `rc=0`.
- **AC-M1-16** — strict UTC timestamp validation (D7a; the round-2 quorum fix), **one arm,
  `watermark-strict`**, asserting only its own F10 messages and exit codes. The reviewer named
  three cases and each gets a run, every one refused BEFORE any gh call (asserted by the stub's
  argv log being empty): `2026-02-30T00:00:00Z` → `(not-canonical '2026-02-30T00:00:00Z' →
  '2026-03-02T00:00:00Z')`, rc=2; `2026-99-99T99:99:99Z` → `(not-canonical
  '2026-99-99T99:99:99Z')`, rc=2; `9999-01-01T00:00:00Z` → `(future '9999-01-01T00:00:00Z' >
  now+300s)`, rc=2. Two controls in the same arm, because three refusals prove nothing on their
  own (rule 3a): a known-good value `2026-09-08T06:19:00Z` must be classified `ok` and reach the
  gh stub, and `now(UTC) - 60s` — computed at run time, never a literal — must ALSO be `ok`, which
  is what proves the future ceiling is a ceiling and not a blanket refusal. The arm asserts no
  `crash:`, `issue:` or `verdict:` line (row-90 property).
  `bash scripts/test_gate0_self_notices.sh --only watermark-strict; echo "rc=$?"` → `ok
  watermark-strict`, `rc=0`.
- **AC-M1-17** — the parser is portable and says so when it is not. The two-arm `date` helper
  (BSD `-j -f`, GNU `-u -d`) resolves on the runner it is executing on, and refuses loudly if
  NEITHER arm can parse the known-good control: a `watermark-parser-control` arm forces both arms
  to fail (a `PATH`-shadowing stub `date` that always exits 1) and asserts the instrument prints
  `✗ timestamp parser unusable: neither BSD nor GNU date parsed the control value` at rc=2 —
  never "every watermark is invalid". This AC is the one that PROVES the GNU arm, because CI is
  ubuntu and the rig is macOS: the suite passing in the `go host build + test gate` job **is** the
  GNU measurement (V17 records the BSD arm as measured on the rig and the GNU arm as UNMEASURED
  locally, deliberately).
  `bash scripts/test_gate0_self_notices.sh --only watermark-parser-control; echo "rc=$?"` → `ok
  watermark-parser-control`, `rc=0`; and the same suite green in CI on ubuntu.

**M2 — the Repo Profile rule and the live reading**

- **AC-M2-1** — the rule exists, before the queue heading, names the script, both issue files and
  both watermark files, and suppresses nothing: `awk '/^## Queue/{exit} /gate0_self_notices\.sh/{c++} END{print c+0}' design_docs/world-mission.md` → `2` or more (base `0`, B5); `awk '/^## Queue/{exit} /mission-world-gh-issue-prev/{c++} /--watermark-file/{w++} END{print c+0, w+0}' design_docs/world-mission.md` → `≥1 2` (base `0 0`, B5); and, scoped to this bullet alone (it opens with the literal `SECOND, NO-AUTHORITY READ`), `awk '/^## Queue/{exit} /^- .*SECOND, NO-AUTHORITY READ/{f=1;next} f&&/^- /{f=0} f' design_docs/world-mission.md | /usr/bin/grep -c '2>/dev/null\|1970-01-01'` → `0` (no redirect, no epoch anywhere in the rule).
- **AC-M2-2** — the rule pins the no-authority clause, the hit consequence and the unreadable-
  channel branch, as literals: `awk '/^## Queue/{exit} /grants no authority/{a++} /A FIRE DIED/{b++} /self-notice read UNAVAILABLE/{c++} END{print a+0, b+0, c+0}' design_docs/world-mission.md` → three integers each `≥ 1` (base `0 0 0`, B5).
- **AC-M2-3** — the live reading at landing, controller-run on the rig (network; not CI). The D7
  command verbatim → rc=**0**, two `watermark: … ok` lines (or one `ok` + one `ABSENT — DEGRADED`
  naming the issue-scoped file if landing falls on a rotation week — recorded either way, and a
  DEGRADED line naming `mission-world-last-seen` is a landing blocker), a `since:` line,
  `issue: 107 comments=43 self=43 crash=4 other=39`, `control: issue=107 expect=4 got=4 ok`,
  `signature-source: …/mission-control.sh ok`, and the verdict line quoted in the STATUS stamp. The `issue: 129 …` line's counts are
  whatever the live issue holds at that moment and are recorded, not predicted. (Base: rc=127, B2.)

**M3 — the fleet half, filed (controller-run)**

- **AC-M3-1** — a proposal issue exists on `sunholo-data/ailang`, authored by this loop, and row
  73's landing tag names it: `gh issue view <N> --repo sunholo-data/ailang --json author,title --jq '"\(.author.login) \(.title)"'` → `sunholo-voight-kampff …Gate 0…`; `awk '/^73\. /,/^74\. /' design_docs/world-mission.md | /usr/bin/grep -c 'sunholo-data/ailang#<N>'` → `≥ 1`.
- **AC-M3-2** — the cross-mission note is sent and read back: `ailang messages send mission-control "<§10 body>" -title "…" -from world` (body POSITIONAL, single-dash flags — the memory-file rule), then the artifact id read back and recorded in the log entry; `awk '/^73\. /,/^74\. /' design_docs/world-mission.md | /usr/bin/grep -c 'msg_'` → `≥ 1`.

---

## §5 Test plan / mutation matrix

One row per arm or criterion; the mutation is applied to `scripts/gate0_self_notices.sh` unless
named otherwise; **expected red is the single named arm, every other arm green** — that is the
row-90 property, and AC-M1-14's drill measures one row of it at landing. Restore from a
`cp`-captured backup and assert its sha256; never `git checkout --` (row 82). The enumeration is
anchored to the diff this item ships (instrument, suite, fixtures, CI step, Repo Profile rule).

| arm / criterion | kills this mutation | kills which arm (expected red) |
|---|---|---|
| `anchored-signature` | `startswith` → `contains` (substring); or prefix shortened to `⚠️ Mission iteration`; or `(rc=` digit check dropped | `anchored-signature` only |
| `rc-extract` | rc captured as the whole `(…)` group, or hard-coded `143` | `rc-extract` only |
| `since-window` | comparison `>` → `>=`; or in-window computed but exit code from all-time count | `since-window` only |
| `self-filter` | author filter removed; or case-sensitive compare | `self-filter` only |
| `prev-issue` | `--prev-issue` ignored (only `--issue` classified); or `--issue` silently also reads `-prev` | `prev-issue` only |
| `no-bodies` | `crash:` line gains a `body=` field, or a debug `echo "$body"` | `no-bodies` only |
| `usage` | `--since` given an epoch default | `usage` only |
| `gh-missing` | rc 127 mapped to "0 comments" | `gh-missing` only |
| `gh-timeout` | `run_bounded` replaced by a bare call | `gh-timeout` only |
| `gh-failed` | non-zero gh rc swallowed | `gh-failed` only |
| `malformed` | parser defaults a missing `comments` key to `[]` | `malformed` only |
| `truncated` | REST cross-check removed (F5 gone) | `truncated` only |
| `control-required` | `--control` made optional | `control-required` only |
| `control-mismatch` | mismatch downgraded to a warning (exit 0); or table suppressed on mismatch | `control-mismatch` only |
| `signature-source` | F8 check removed; or `unchecked` token dropped when the flag is absent | `signature-source` only |
| `watermark-older` | minimum → maximum (`head -1` → `tail -1`); or only the first `--watermark-file` read | `watermark-older` only |
| `watermark-degraded` | an absent file treated as the epoch (r0's silent collapse); or the `DEGRADED` line dropped; or one-absent escalated to rc=2 | `watermark-degraded` only |
| `watermark-unavailable` | both-absent falls back to the epoch (sol's headline); or F9 fires only when the FIRST file is absent | `watermark-unavailable` only |
| `watermark-invalid` (empty) | empty file treated as absent (→ DEGRADED rc=0 instead of F10) | `watermark-invalid` only |
| `watermark-invalid` (not-ISO8601) | the regex removed, so `sort` discards the malformed value (V16 state 6); or validation applied only to the chosen minimum | `watermark-invalid` only |
| `watermark-invalid` (unreadable) | `-e` used without `-r`, so an unreadable file reads as empty | `watermark-invalid` only |
| `watermark-invalid` (`--since`) | the literal form exempted from the regex | `watermark-invalid` only |
| `report` | verdict wording changed (`crash notice(s)` → `crash notices`; `A FIRE DIED` removed; the `untouched; … retry` sentence copied in) | `report` only — **measured at landing (AC-M1-14)** |
| `snapshot-107` | fixture regenerated without the cutoff or with byte-truncation; `other=` computed as `comments − crash` (ignoring `self`); classifier drift that flips the 09-03 correction | `snapshot-107` only |
| `snapshot-129` | `issue:` line for a zero-crash issue omitted | `snapshot-129` only |
| AC-M1-13 | a second, literal `gh issue view` call added beside `"$GH_BIN"`; argv shape changed (`--json comments` dropped) | AC-M1-13's grep → `1`; `snapshot-107`'s argv assertion |
| AC-M1-14 grep | CI suite step removed | AC-M1-14 |
| AC-M2-1 | rule written after `## Queue` (inside a row); `-prev` file not named; a `2>/dev/null` or an epoch literal reintroduced into the command; one `--watermark-file` dropped | AC-M2-1 |
| AC-M2-2 | "grants no authority" softened; unreadable-channel branch dropped | AC-M2-2 |
| AC-M2-3 | rule's command points `--driver-src` at a World path or omits it (`unchecked`); control re-pointed to `129:0` | AC-M2-3 (the recorded live line differs) |
| AC-M3-1/2 | proposal not filed; note not read back | AC-M3-1 / AC-M3-2 |

**Floor reachability is itself asserted:** each floor arm's `ok` requires the floor's own message,
so a harness change that makes a floor unreachable (a stub that starts answering every argv, a
`make_json` helper that inserts `comments`) reds that arm immediately.

**Honest note on shared assertions:** `snapshot-107` also asserts the `control:` line and the argv
log, so a mutation of the control *printer* reds `control-mismatch` (message) and `snapshot-107`
(the `ok` line) — two arms, both owning a field of that behaviour. This is the one intentional
overlap and it is recorded here rather than discovered by a judge.

---

| `watermark-strict` (canonicalization) | the round-trip canonicalization check removed, leaving shape-only validation | `watermark-strict` only |
| `watermark-strict` (future ceiling) | the `now+300s` ceiling removed | `watermark-strict` only |
| `watermark-strict` (ceiling inverted) | the ceiling applied as a FLOOR (reject anything OLDER than now-300s) — kills the blanket-refusal reading | `watermark-strict` only |
| `watermark-strict` (order) | validation moved to AFTER the first gh call | `watermark-strict` only (its empty-argv-log assertion) |
| `watermark-parser-control` | the neither-arm-parsed refusal replaced by treating an unparseable control as `invalid` | `watermark-parser-control` only |


## §6 Conflict surface

- **The shared skill (V1 checkout, read-only for World).** Gate 0.6 prescribes
  `mission_directives.sh` and forbids hand-rolling the `gh | jq` pipeline *for directives*. This
  instrument is not a directive reader: it never reads `MISSION_DIRECTIVE_AUTHORS`, cannot widen
  the allowlist, and filters to the loop's OWN account — the account `mission_directives.sh`
  refuses to accept as a principal. The two reads are complementary by construction. Gate 5's
  rotation step already writes `-gh-issue-prev` (`gate-5-retro.md:152`, V6) — this design consumes
  it, changes nothing about it. Gate 2's traces (a)–(c) (`gate-2-pick.md:184–240`) are what a hit
  points at; no trace is redefined.
- **`tools/launchd/*`** — frozen core, untouched. The signature is a World literal *about* the
  driver; the driver-src control reads the pinned tree by absolute path and writes nothing. If the
  fleet rewords the notice, the World edit is one literal + one fixture (D7's last bullet), and the
  fleet proposal (§10) asks for the wording to be stable and machine-anchored.
- **`scripts/gate1_range_check.sh` / `queue_census.sh`** — siblings, not dependencies; the
  `run_bounded` wrapper and the `--gh-bin` stub convention are copied in shape, not sourced.
- **CI (`.github/workflows/ci.yml`)** — one additive suite step; no network, no `gh` auth, no live
  step. The `ailang-verify` job is untouched.
- **State files** — read by the Repo Profile rule only (`-gh-issue`, `-gh-issue-prev`, the two
  watermarks); nothing is written. The instrument itself reads no state.
- **Adjacent rows** — 69 (heartbeat by absolute path: same class, same remedy shape), 71 (the
  `/tmp` log this notice is the durable twin of), 88 (push-boundary attribution — a different
  Gate-1 surface), 90 (the arm-coupling lesson this suite is built to).

---

## §7 Verification Log

Every claim above that the rig or the codebase currently does X has a row here with its COMMAND
and OBSERVED OUTPUT. All commands run under `bash` from
`/Users/voightkampff/dev/sunholo-data/ailang-world` at `0d316388b435e0455fbda97506ef6b994026816b`
(`git rev-parse HEAD` == `git rev-parse origin/dev`; `git status --short` empty). **Provenance
column:** `first-party` = run by the designer this iteration; `inherited` = the controller's
evidence package, cited but not re-run; `UNVERIFIED` = a claim neither measured. `gh` is
`/opt/homebrew/bin/gh` 2.99.0, authenticated as `sunholo-voight-kampff` (`gh api user --jq .login`).

| # | Claim | COMMAND | OBSERVED OUTPUT | Provenance |
|---|---|---|---|---|
| V1 | The writer's line number differs across the three trees; content is the anchor | `grep -n 'FAILED to complete' /Users/voightkampff/dev/sunholo-data/ailang/tools/launchd/mission-control.sh`; `git -C …/ailang show 48c4a6e49:tools/launchd/mission-control.sh \| grep -n 'FAILED to complete'`; `git -C …/ailang rev-parse --short HEAD`; `git -C …/ailang log -1 --format='%H %ci' -- tools/launchd/mission-control.sh` | `1925: --body "⚠️ Mission iteration **FAILED to complete** (rc=$RC — timeout or crash) at $(date '+%F %H:%M %Z'). Log on the rig: \`$LOG\`. The queue is untouched; the next interval will retry. Further identical failures are silent until the rc changes or an iteration completes."`; pinned `1877:`; fleet HEAD `f7d04edb3`; `f090cfec… 2026-09-08 08:27:30 +0200`. Row 73's `:617` matches neither. The writer is episode-gated (`rcfail.episode`: identical rc suppressed until it changes or an iteration completes) | first-party (the controller's `:1877` is the pinned tree — agrees) |
| V2 | `#107`: 43 comments, all self-authored, four anchored crash notices | `gh issue view 107 --repo sunholo-data/ailang-world --json comments,createdAt,title --jq '{createdAt, n: (.comments\|length), authors: (.comments\|map(.author.login)\|group_by(.)\|map({(.[0]): length})\|add)}'`; `… --jq '[.comments[] \| select(.body \| startswith("⚠️ Mission iteration **FAILED to complete** (rc=")) \| {at: .createdAt, rc: (.body\|capture("rc=(?<rc>[0-9]+)").rc)}]'` | `{"createdAt":"2026-08-31T09:26:51Z","n":43,"authors":{"sunholo-voight-kampff":43}}`; `[{"at":"2026-09-02T12:02:13Z","rc":"143"},{"at":"2026-09-02T15:21:07Z","rc":"143"},{"at":"2026-09-02T19:27:22Z","rc":"143"},{"at":"2026-09-07T02:47:28Z","rc":"1"}]` | first-party (agrees with the controller's four) |
| V3 | A substring search over-counts by one; the fifth is the iter-151 correction quoting the notice at offset 320; the notice tail drifted 200 → 284 chars | `… --jq '[.comments[] \| select(.body\|test("FAILED to complete"))] \| length'`; `… --jq '.comments[] \| select(.createdAt=="2026-09-03T10:57:25Z") \| {len: (.body\|length), idx: (.body\|index("FAILED to complete")), idx2: (.body\|index("⚠️ Mission iteration"))}'`; `… --jq '[.comments[] \| select(.body \| startswith("⚠️ …(rc=")) \| (.body\|length)]'`; `… --jq '.comments[] \| select(.body\|test("FAILED to complete")) \| {at: .createdAt, head: (.body\|.[0:260])}'` | `5`; `{"len":2580,"idx":341,"idx2":320}`; `[200,200,200,284]`; the fifth head: `**Iteration 151 — correction, against my own record above**\n\nI reported slot 1's cause … as **unrecoverable** …` | first-party (NOT in the evidence package — designer's finding) |
| V4 | `#129`: 14 comments, all self-authored, zero crash notices, and the three condition classes | `gh issue view 129 … --jq '{createdAt, n, authors, crash: [… test("FAILED to complete") …], degraded: [… "lane degraded on this fire" …], model: [… "Controller model:" …], unpinned: [… "ran UNPINNED on this fire" …]}'` | `{"createdAt":"2026-09-07T07:51:05Z","n":14,"authors":{"sunholo-voight-kampff":14},"crash":[],"degraded":["2026-09-07T14:26:29Z","2026-09-08T02:45:06Z","2026-09-08T06:28:52Z"],"model":["2026-09-08T02:43:12Z","2026-09-08T06:28:34Z"],"unpinned":["2026-09-08T02:47:00Z"]}` | first-party |
| V5 | Gate 0's read drops every self-authored comment; the positive control hits | `bash /Users/voightkampff/dev/sunholo-data/ailang/scripts/mission_directives.sh --issue 107 --repo sunholo-data/ailang-world`; same `--issue 129`; control `--issue 972 --repo sunholo-data/ailang` | `#107 — 0 directive(s) from [MarkEdmondson1234] since 1970-01-01T00:00:00Z (of 43 comments …)`; `#129 — 0 … (of 14 comments …)`; control `#972 — 1 directive(s) … (of 59 comments …)` with `MarkEdmondson1234 @ 2026-09-02T07:17:34Z: D-54 b` on stdout | first-party (the controller's `#852`/`#850` → `2`/`1` are inherited, not re-run) |
| V6 | Rotation: `-prev` = 107, current = 129, and the rc=1 notice predates `#129` by 5 h 04 m | `cat ~/.ailang/state/mission-world-gh-issue ~/.ailang/state/mission-world-gh-issue-prev`; `#129` `createdAt` (V4) vs V2's last `at` | `129` / `107`; `2026-09-07T07:51:05Z` − `2026-09-07T02:47:28Z` = 5 h 03 m 37 s | first-party |
| V7 | The false-reassurance instance is in the mission log | `grep -n 'FAILED to complete\|crash notice' design_docs/world-mission-log.md`; `sed -n '984,985p;1010p' …-log.md` | `984: already posted its own crash notice** — \`#107\` comment \`2026-09-02T19:27:22Z\` …`; `985: … **61 seconds after the merge**.`; `1010: … the **driver** posts its crash notices to that same issue, so they are` | first-party |
| V8 | The local slot-verdict log records the 09-07 crash but begins 2026-09-05 | `cat ~/.ailang/state/mission-world-slot-verdicts.log` | first line `2026-09-05T20:38:39Z verdict=COMPLETED …`; `2026-09-07T02:47:27Z verdict=CRASHED_at=fired rc=1 attempt=1/3 elapsed_s=168 stamps=1 controller=pi:ollama/glm-5.3:cloud`; 15 lines; the driver keeps `tail -n 200` | first-party |
| V9 | The condition-class notices have local durable copies | `cut -c1-200 ~/.ailang/state/mission-world-notice-spool.tsv`; `ls ~/.ailang/state/ \| grep mission-world` | four spool rows (`lane degraded` ×3, `driver ran UNPINNED` ×1, 2026-09-07/08); `mission-world-lane-degraded.episode`, `mission-world-model-last` present | first-party |
| V10 | `gh issue view --json comments` agreed with the REST count on every issue tried; none ≥ 100, so the cap is unmeasured | `for n in 972 852 1089 979; do gh issue view $n --repo sunholo-data/ailang --json comments --jq '.comments\|length'; gh api repos/sunholo-data/ailang/issues/$n --jq .comments; done` | `59/59`, `84/84`, `3/3`, `23/23` | first-party; **the cap itself is UNVERIFIED** → F5 |
| V11 | World cannot edit the skill; the driver env lives under `~/.config/ailang` | `readlink -f ~/.claude/skills/mission-control`; `grep -n '\.env' …/ailang/tools/launchd/mission-control.sh \| head -3` | `/Users/voightkampff/dev/sunholo-data/ailang/.claude/skills/mission-control`; `4: # sources ~/.config/ailang/mission-<name>.env`, `81–82: … "$HOME/.config/ailang/mission-${MISSION_NAME}.env"` | first-party |
| V12 | The fixture generator is reproducible and codepoint-safe | the §3 `gen` command into a scratch dir, then `wc -c`, `shasum -a 256`, `jq -r '.comments \| "n=\(length) last=\(map(.createdAt)\|max)"'`; `printf '%s' '"a—b—c"' \| jq -r '.[0:3]' \| od -An -c`; `grep -c -axv '.*' 107.json`; python3 re-count of `startswith` vs substring on the truncated file | `107: bytes=36718 sha256=79ad0ad7… n=43 last=2026-09-07T07:51:27Z`; `129: bytes=12624 sha256=d61af56d… n=14 last=2026-09-08T06:28:52Z`; slice → `a 342 200 224 b` (the em dash intact, 3 bytes); `0` invalid lines; `comments 43 crash 4 substr 5` / `comments 14 crash 0 substr 0` | first-party |
| V13 | CI insertion point; python3 is the house parser and is present in CI | `grep -nE '^  [a-z-]+:$\|name: Queue census' .github/workflows/ci.yml`; `sed -n '20,25p' scripts/verify_ail.sh`; `grep -n python3 scripts/gate1_range_check.sh \| head -2` | jobs `ailang-verify` (18), `go-verify` (99); steps at 206 and 210; `python3 is REQUIRED (JSON parse; no jq dependency) … present on macOS dev machines and ubuntu-latest CI runners`; `69: # ── python3 total_count parser (house pattern …)` | first-party |
| V14 | Repo Profile anchors: the heading, the watermark bullet, the absolute-path bullet, no existing crash-notice rule | `grep -n '^## Repo Profile\|^## Queue' design_docs/world-mission.md`; `grep -n 'TWO\* MARK-COMMENT' …`; `grep -n 'mission_directives' … \| awk -F: '$1<833'`; `awk '/^## Queue/{exit} /crash notice\|crash-notice/{c++} END{print c+0}' …` | `36`, `1636`; `44`; `122,123,128,133`; `0` | first-party |
| V15 | Rig shell/tools | `/bin/bash --version \| head -1`; `awk --version`; `python3 --version`; `jq --version` | `GNU bash, version 3.2.57(1)-release (arm64-apple-darwin25)`; `awk version 20200816`; `Python 3.14.5`; `jq-1.7.1-apple` | first-party |
| V16 | **The r0 D7 one-liner degrades silently in four of six watermark states** — the r1 objection, measured by the CONTROLLER on 2026-09-08 before it was forwarded (rule 3f), and the reason D7 is now instrument code with F9/F10 | the r0 snippet `SINCE="$(sort <(cat ~/.ailang/state/mission-world-last-seen) <(cat "$HOME/.ailang/state/mission-${ISSUE}-last-seen" 2>/dev/null) \| head -1)"` run against six states of the two files | both present (`2026-09-08T06:19:00Z`, `2026-09-05T00:00:00Z`) → clean, `2026-09-05T00:00:00Z`, correct; **mission-scoped ABSENT** → `cat: …: No such file or directory` on stderr, `2026-09-05T00:00:00Z` — the take-the-older rule silently collapses to ONE watermark; mission-scoped present / issue-scoped absent → clean, `2026-09-08T06:19:00Z`, same collapse with no stderr at all; **both absent** → leaks, SINCE empty → silent epoch; **mission-scoped present but EMPTY** (a real state: rotation creates the file) → clean, `2026-09-05T00:00:00Z`, indistinguishable from healthy; **mission-scoped `not-a-date`** → clean, `2026-09-05T00:00:00Z` — `sort` puts the ISO string first, so the malformed value is silently discarded and loses to a valid one. Four distinct silent-degradation paths; r0 named one | **controller-measured (iter-173, 2026-09-08)**; designer re-derived the D7 table from it, did not re-run |
| B1 | AC-M1-1 red at base | `for f in scripts/gate0_self_notices.sh scripts/test_gate0_self_notices.sh; do test -x $f; echo rc=$?; done` | `rc=1`, `rc=1` | first-party |
| B2 | AC-M1-2…12, AC-M2-3 red at base; precedent suite is the control | `bash scripts/test_gate0_self_notices.sh; echo rc=$?`; control `bash scripts/test_gate1_range_check.sh --only report \| tail -1` | `No such file or directory`, `rc=127`; control `4 passed, 0 failed` | first-party (control inherited from the census doc's B3 shape; re-run) |
| B3 | Fixtures absent at base | `for f in scripts/testdata/gate0_self_notices_107.json scripts/testdata/gate0_self_notices_129.json; do test -e $f; echo rc=$?; done` | `rc=1`, `rc=1` | first-party |
| B4 | AC-M1-14 grep red at base; control hits | `grep -c 'gate0_self_notices' .github/workflows/ci.yml`; control `grep -c 'test_gate1_range_check' …` | `0`; `1` | first-party |
| B5 | AC-M2-1/2 red at base | `awk '/^## Queue/{exit} /gate0_self_notices/{c++} END{print c+0}' design_docs/world-mission.md`; `grep -rn 'gate0_self_notices\|self_notices' --exclude-dir=.git . \| wc -l` | `0`; `0` | first-party |
| B6 | Charter blob at base | `git hash-object design_docs/world-mission.md` | `4e98fec339c581520bdc6c5c1ac4ea73f0e7fbf3` | first-party |

**Inherited and not re-run:** the controller's directive-reader controls on `sunholo-data/ailang`
`#852`/`#850` (`2`/`1`). Everything else in §1 was re-measured here, and one inherited number was
*sharpened* rather than contradicted: the evidence package's "four crash notices" is correct under
the anchored rule and would read five under a substring rule (V3).

---

| V17 | **Shape-only validation accepts calendar-invalid and far-future watermarks, and a far-future watermark silently zeroes the instrument** — the round-2 objection, measured by the CONTROLLER on 2026-09-08 before the fix was applied (rule 3f; reproduced in BOTH directions per the judge-findings rule, and it is bigger than filed) | (a) the doc's r1 regex `^[0-9]{4}-[0-9]{2}-[0-9]{2}T[0-9]{2}:[0-9]{2}:[0-9]{2}Z$` run against six values in bash; (b) `TZ=UTC date -j -f '%Y-%m-%dT%H:%M:%SZ' <v> +'%Y-%m-%dT%H:%M:%SZ'` round-trip on the rig; (c) `jq --arg s <since> '[.comments[]\|select((.body\|startswith("⚠️ Mission iteration **FAILED to complete** (rc=")) and (.createdAt > $s))]\|length'` over a saved `#107` payload | (a) ACCEPTED: `2026-09-08T06:19:00Z`, `2026-99-99T99:99:99Z`, `2026-02-30T00:00:00Z`, `9999-01-01T00:00:00Z`; rejected: `not-a-date`, `2026-9-8T6:19:0Z` — so shape admits three invalid values. (b) `2026-09-08T06:19:00Z` → SAME (rc=0); `2026-02-30T00:00:00Z` → `2026-03-02T00:00:00Z` (DIFFER); `2026-02-29T00:00:00Z` → `2026-03-01T00:00:00Z` (DIFFER, 2026 is not a leap year); `2026-99-99T99:99:99Z` → rc=1 `illegal time format`; `9999-01-01T00:00:00Z` → **SAME (rc=0)** — round-tripping alone does NOT catch the future case, which is why the ceiling is a separate check. GNU `date` arm UNMEASURED (no `gdate` on the rig) → proven by CI, AC-M1-17. (c) the consequence: `since=1970-01-01T00:00:00Z` → **4** notices, `since=2026-09-01T00:00:00Z` → **4**, `since=9999-01-01T00:00:00Z` → **0**, `since=2026-99-99T99:99:99Z` → **0**. A shape-valid bad watermark turns four real deaths into a clean rc=0, for ever | **controller-measured (iter-173, 2026-09-08)**; the fix is `gpt5-6-sol`'s `proposed_fix` applied VERBATIM under the narrow-refinement carve-out |

## §8 Milestones

Total **~0.3d**. Three milestones, each independently committable; none touches `tools/launchd/*`,
the shared skill, or any queue row body.

- **M1 — the instrument, its suite, the two snapshots, and the suite's CI step (~0.2d).**
  `scripts/gate0_self_notices.sh` (D1–D6; bash 3.2 + one python3 parser + `run_bounded`),
  `scripts/test_gate0_self_notices.sh` (22 arms, `--only`, literal-heredoc fixtures, argv-logging
  stub), `scripts/testdata/gate0_self_notices_{107,129}.json` + `.meta.json` (generated by §3's
  command, checked in as produced), one `Gate-0 self-notice instrument suite` step in `go-verify`
  after the live census step. Green when AC-M1-1…14 hold and the three profile gates are
  unchanged-green. **Useful alone:** from this commit on, any controller can run the D7 command by
  hand and see the week's deaths on both issues, with controls — this mission's Gate 0 is no
  longer blind even if M2 were never written and the fleet never acted.
- **M2 — the Repo Profile rule and the live reading (~0.05d, controller-run at landing).** One
  bullet after the ABSOLUTE-PATH bullet (D7 verbatim); the live D7 command run on the rig and its
  verdict line quoted in the iteration's STATUS stamp; row 73 tagged `[LANDED …]` in the leading
  position. Green when AC-M2-1…3 hold.
- **M3 — the fleet half, filed (~0.05d, controller-run).** A GitHub issue on `sunholo-data/ailang`
  and an `ailang messages send mission-control` note (§10), both ids recorded in row 73's tag and
  the log entry. Green when AC-M3-1…2 hold. M3 is not a stub that M1 waits on; M1 is complete
  without it, and M3 is the mission-independent half being handed to its owner.

---

## §9 Risks, declared residuals and honest limitations

- **The count is a lower bound, and a death is in-window for one fire.** The driver suppresses
  identical-rc repeats (*"Further identical failures are silent until the rc changes …"*, V1), so
  two same-rc deaths post ONE notice and the instrument cannot see what was never posted; and with
  the window at the watermark, a death is `in-window` only for the first fire after it. The
  all-time `crash:` lines mitigate both (every fire sees the week's history) and the exit-1
  consequence is idempotent. The local slot-verdict log (V8) records every slot and is the natural
  corroborating source, but it is rig-local, 200-row capped and three days old — reading it at
  Gate 0 is a candidate follow-on row, not this item.
- **The signature is a literal against a frozen writer.** Mitigated by two opposite-direction
  controls (D5) rather than by pretending the wording is stable. A fleet rewording produces a red
  on the very next fire (`signature-source`), and the fix is one literal + one fixture.
- **`gh issue view --json comments` may cap at 100** — UNVERIFIED (V10). Converted to F5 (REST
  cross-check, exit 2). If the cap is real and an issue reaches it, the instrument refuses loudly
  rather than reporting a partial week as complete; the remedy then is pagination via
  `gh api …/comments --paginate`, out of scope until measured.
- **`other=N` is a bucket, not a classification** (D3). The three condition classes are visible
  as a number and reachable by URL; classifying them is declared out, with the argument recorded so
  a later row can overturn it with a measurement rather than re-litigate it.
- **The live step is controller discipline, not CI.** Like the directive read, the D7 command needs
  `gh` auth. CI proves the instrument on fixtures whose argv shape equals the live one (AC-M1-13);
  the rule plus the stamp quote is what makes an omitted run visible in the record.
- **Exit 1 is "look", not "abort".** A hit on a fire that was already correctly credited (a
  controller who read `#107` at Gate 5 last week) costs three read-only traces. That is the price
  of a channel that is now read every fire; the alternative was the channel not being read.
- **The instrument treats its two watermark files symmetrically; the rule carries which one's
  absence is abnormal.** An instrument that knew "the mission-scoped file must always exist" would
  need a `--watermark-file-required` flag and a fifth arm for ~0.05d of value; instead the DEGRADED
  line names the path and the rule says which name is a blocker (D7). If a second instance of a
  missing mission-scoped file ever lands, that is the row that adds the flag.
- **Row 73's row text cites `mission-control.sh:617`** — stale (V1); recorded here, not edited into
  the ratified row.

---

## §10 The fleet half — proposal to Mark + V1 (World cannot edit the shared skill)

**Bar:** the shared skill's ≥2-friction bar for a proposal is **met**: instance 1 at iteration 151
(the row-56 slot's notice, found by accident at Gate 5, controller-measured), instance 2 at
iteration 173 (re-measured here first-party: four notices on `#107`, zero of them visible to the
prescribed Gate-0 read, positive control on `sunholo-data/ailang#972` surfacing a real directive —
V2/V5). The defect is mission-independent: every mission on this rig posts the same notice as its
own account to its own issue and filters it with the same allowlist.

**The ask, one paragraph (the controller sends this; the doc only drafts it):**

> Gate 0.6 reads the bookkeeping issue through `mission_directives.sh`'s author allowlist, which is
> correct and must not be widened — and which, by construction, drops the driver's own
> `⚠️ Mission iteration **FAILED to complete** (rc=…)` notices, because the driver posts them as the
> loop's account. Measured on World: four such notices on `#107` (three rc=143 on 2026-09-02, one
> rc=1 on 2026-09-07), 0 visible to the Gate-0 read on two consecutive fires, the last one landing
> five hours before the weekly rotation created `#129`. Ask: add to Gate 0, beside the directive
> read, a **second, separate, no-authority read** — self-authored comments on the current AND the
> `-prev` issue, matched by the notice's anchored prefix, since the watermark — whose only
> consequence is to trigger Gate 2's died-mid-flight traces (a)–(c) and a log credit; it can
> unpark nothing and pick nothing. Two supporting asks for the driver, fleet-owned: (1) keep the
> notice's **prefix** stable and machine-anchored (the tail already drifted once, 200 → 284 chars,
> between 2026-09-02 and 2026-09-07); (2) reword the reassurance — *"The queue is untouched; the
> next interval will retry"* read as "nothing happened" on the one occasion it mattered, when the
> dead fire had merged its PR 61 seconds earlier and the queue was untouched **because** the work
> was done. World ships its own read (`scripts/gate0_self_notices.sh` + Repo Profile rule) so this
> mission is covered regardless; this asks for the same second read in the shared Gate 0 so the
> other missions are too.

Filed as a GitHub issue on `sunholo-data/ailang` and an `ailang messages send mission-control`
note (AC-M3-1/2). Not proposed: any change to the allowlist, any local driver edit, any vendored
copy of `mission_directives.sh`.
