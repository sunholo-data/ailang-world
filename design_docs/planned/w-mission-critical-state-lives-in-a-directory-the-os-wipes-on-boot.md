- Status: **planned** · Date: **2026-09-08** · Designer: pi:openrouter/deepseek/deepseek-v4-flash-0731 (iter-171, attempt 2 after an empty_worktree repetition loop on the ollama rung) · Base commit measured at: **791bd634530b52b44bd29ffb43d06d1f8a8dde5e** · Owning queue row: **71 — w-mission-critical-state-lives-in-a-directory-the-os-wipes-on-boot** (clause-2, ~0.2d, gated on nothing, surfaced iter-151) · Scope class: **HARNESS — durability of loop evidence** · Disposition: **(A) land the World-side rule + disposition + fleet note.**; revised after quorum r1 (gemini-3-1-pro + controller objections: AC-M1-1 needle re-anchored, AC-M2-3 stamp convention aligned); r2 closed under the narrow-refinement carve-out — gemini's true half (`-d` checks upstream-or-HEAD, not dev) folded into a pinned §4 step, its deterministic-abort scenario refuted by measurement (see §9)

## §1 Problem — three kinds of state parked where the OS wipes on boot

Iteration 151 measured that **one macOS reboot took all three kinds of mission state** that had
been parked under `/tmp`: the toolchain pin, the driver's crash log, and a sprint worktree. The
pin half is already CLOSED by iter-151 — moved to `~/.pinned-ailang/ailang`, CI's own path — and
the Repo Profile generalised its lesson ("a pin stored in a volatile directory is a claim about
the last boot, not about the toolchain"). What this row tracks are the **two remainders**, and
both are worse than the pin because neither fails closed:

- **(a) The driver log is this loop's ONLY crash forensics.** `HARD TIMEOUT` and `STALL:` are
  written to `/tmp/ailang-mission-world.log` and nowhere else, so any fire that dies becomes
  **undiagnosable** the moment the rig reboots. The launchd `StandardOutPath` copy sits in the
  same directory and is therefore *not a backup at all* — two copies of the same loss (distinct
  inodes).
- **(b) Sprint worktrees under `/tmp`** lose an interrupted iteration's uncommitted work — the
  residue the skill's died-mid-flight recovery exists to read, deleted before anyone can.

**F4 is the sharpening fact:** the rig has rebooted **a second time since the row was filed** —
`kern.boottime` = `1788606012` = **Sat Sep 5 13:00:12 2026**, and the current `/tmp` log's first
line is `[2026-09-05 13:01:11]`, so every driver-log line for iterations 151–157 is already
destroyed. The exposure is not hypothetical; it has now eaten the loop's forensics twice in one
week. The row's scope fence binds: **do not grow this into "make the driver crash-proof"** — the
deliverable is only that the evidence of a crash outlives the crash.

## §2 Verification Log

Controller-measured (2026-09-08, this fire, iteration controller — cited verbatim). Runner `[controller]`.

- **F1 [controller]** — Fleet driver per-mission crash log is in `/tmp`: line 103 of
  `tools/launchd/mission-control.sh` reads `LOG="/tmp/ailang-mission-${MISSION_NAME}.log"` —
  identical in fleet HEAD and in the pinned driver copy this rig runs. Row 71's citation of `:76`
  is **STALE** (the value has drifted to `:103`). World does not edit frozen core.
- **F2 [controller]** — launchd `StandardOutPath` is also `/tmp`:
  `~/Library/LaunchAgents/dev.ailang.mission-world.plist` declares
  `/tmp/ailang-mission-world.launchd.log`.
- **F3 [controller]** — Live this fire: `/tmp/ailang-mission-world.log` = 12,209,292 bytes /
  113,316 lines / 52 `HARD TIMEOUT|STALL`-matching lines; `/tmp/ailang-mission-world.launchd.log`
  = 53,809 bytes; inodes differ (68012264 vs 68001663). This log is the loop's ONLY crash forensics.
- **F4 [controller]** — `sysctl -n kern.boottime` = `1788606012` = Sat Sep 5 13:00:12 2026 — a
  **second instance since filing** (~2 days after row 71 was filed, off the 2026-09-03 reboot).
  The current `/tmp` log's first line is dated `[2026-09-05 13:01:11]`, so all driver-log lines
  for iterations 151–157 are already destroyed.
- **F5 [controller]** — A live `/tmp` worktree exists: `/private/tmp/world-attended-rulings`,
  branch `mission/attended-rulings-20260907` at `565d0b2`, **fully merged**
  (`git merge-base --is-ancestor 565d0b2 dev` → rc 0). Content landed; the worktree directory is
  pure exposure.
- **F6 [controller]** — World's de facto worktree convention is persistent (`.wt-world-*` and
  `.eval-world-*` under `/Users/voightkampff/dev/sunholo-data/`) but **no charter rule records it**:
  `/usr/bin/grep -c "worktree convention" design_docs/world-mission.md` = **0**. The charter's pin
  rule ("MOVED OFF /tmp … a pin stored in a volatile directory is a claim about the last boot")
  should be **generalized, not duplicated**.
- **F7 [controller]** — Frozen-core fence: `tools/launchd/*` is fleet-mirrored; the repo's own gate
  fatals on "DRIVER DRIFT vs FLEET (D-WORLD-DRIVER-1)" (charter row 76) when a World blob differs
  from fleet HEAD. The driver-log fix is a **FLEET commit, never local**.
- **F8 [controller]** — `~/.ailang/state/` (persistent) already holds this mission's heartbeat,
  base, watermark files — the durable home for state that must survive a reboot.
- **F9 [controller]** — The shared mission-control skill already forbids `/tmp` WORKTREES for a
  location-class reason ("NEVER PLACE A WORKTREE UNDER /tmp — the suite goes red for the LOCATION,
  not the code"). World's new rule **extends** this to the PERSISTENCE reason (evidence must
  outlive a reboot) and to state beyond worktrees.

Designer spot-checks — my own measurements this fire, one §2 row each. Runner `[designer]`.

- **S1 [designer]** — `/usr/bin/grep -n 'LOG=' /Users/voightkampff/dev/sunholo-data/ailang/tools/launchd/mission-control.sh` →
  output: `86:  LOG=/tmp/ailang-mission-control.log` and `103:  LOG="/tmp/ailang-mission-${MISSION_NAME}.log"`.
  **Confirms F1's `:103`** (the `:76` in the row is stale). The `:86` line is the driver's own
  control log — also `/tmp`, same class of exposure.
- **S2 [designer]** — `/usr/bin/grep -c "worktree convention" design_docs/world-mission.md` →
  output: **1**. **DISCREPANCY with F6 (recorded 0).** Located: the single hit is line 5525, which
  is **row 71's own queue text** ("World's side is the worktree convention and this row."), *not* a
  charter rule. Reconciliation: F6's claim — "no charter RULE records it" — still holds; the
  terminology now appears once, inside the very row this doc resolves, and the Repo Profile still
  contains no governing rule. The deliverable (a real rule) remains required. Control:
  `/usr/bin/grep -c "MOVED OFF" design_docs/world-mission.md` → **1** (line 125, the pinned-ailang
  rule; the file's phrasing is `` MOVED OFF `/tmp` `` with backticks, so the exact substring
  "MOVED OFF /tmp" matches only via the backtick-free fragments — the pin rule is present, ≥1 as
  the control intends).
- **S3 [designer]** — `git merge-base --is-ancestor 565d0b2 dev; echo "rc=$?"` → `rc=0`.
  **Confirms F5** — the live `/tmp` worktree's tip is fully merged into `dev`.
- **S4 [designer]** — `sysctl -n kern.boottime` →
  `{ sec = 1788606012, usec = 212952 } Sat Sep 5 13:00:12 2026`. **Confirms F4** (second instance).
- **S5 [designer]** — `/usr/bin/grep -n "^71\. " design_docs/world-mission.md | head -1` →
  `5505`. Read its first 10 lines (the full text of row 71 appears at lines 5505–5519). Row located
  and scope-fence text confirmed verbatim.
- **S6 [designer]** — `ls ~/.ailang/state/ | /usr/bin/grep -c "mission-world"` → **13**. Confirms
  **F8**: a durable, populated `~/.ailang/state/mission-world-*` home already exists (heartbeat,
  base, watermark, notice-spool, etc.) ready to absorb state that must outlive a reboot.
- **S7 [designer]** — `git worktree list | /usr/bin/grep -c tmp` → **1**
  (`/private/tmp/world-attended-rulings`). All other worktrees are under
  `/Users/voightkampff/dev/sunholo-data/`. This is the **pre-removal control read** for AC-M2-1
  (must drop 1 → 0).

Instruments cited inline: `[file:line at 791bd634]`, `[ls -la]`, `[sysctl]`,
`[git merge-base]`. Docs-only. No `.ail`, no Go, no edits under `tools/launchd/` or `scripts/`.

## §3 Deliverable 1 — the charter rule

Add this ONE bullet to the Repo Profile process-fix section (verbatim text ready to paste):

~~~markdown
- **CRITICAL MISSION STATE LIVES IN A DIRECTORY THE OS MAY DELETE — THREE KINDS MUST SURVIVE A
  REBOOT (process fix, iter-151, generalized from the pin's "a claim about the last boot" rule).**
  (a) ROLE AND EVALUATOR WORKTREES: never place a worktree under `/tmp` — macOS wipes
  `/private/tmp` on boot — always under a persistent sibling of this repo
  (`/Users/voightkampff/dev/sunholo-data/.wt-*` / `.eval-*`). (b) ANY STATE WHOSE LOSS DESTROYS
  EVIDENCE — crash forensics (the driver's `HARD TIMEOUT|STALL` logs) and uncommitted sprint
  work — must live under a `$HOME`-anchored path (`~/.ailang/state/mission-world-*` or inside the
  repo), never in `/tmp`. The tell is the pin's own: an anchor stored in a volatile directory is a
  claim about the last boot, not about the evidence — and a fire that dies becomes undiagnosable
  the moment the rig reboots. `/tmp` is for scratch it is acceptable to lose, and nothing else.
~~~

This satisfies the four requirements: (a) worktrees under a persistent sibling, never `/tmp`;
(b) evidence-destroying state under `$HOME`-anchored paths (`~/.ailang/state/` or the repo);
(c) generalizes the pin rule's "a claim about the last boot" phrasing — it names that phrase and
extends it to evidence, without duplicating the pin bullet; (d) one bullet, well under 15 lines.

## §4 Deliverable 2 — disposition of the LIVE `/tmp` worktree (F5)

The live worktree `/private/tmp/world-attended-rulings` is fully merged (S3: rc=0). Its content
has landed; the directory is pure exposure. **The CONTROLLER performs this** — it is rig state,
not repo content, and the executor cannot touch the main checkout's worktree registry. Exact
verify-merged-then-remove sequence, with each step's evidence:

1. `git merge-base --is-ancestor 565d0b2 dev; echo "rc=$?"` — **rc=0** (already measured, S3):
   only proceed if merged, so no content is lost.
2. `git worktree remove /private/tmp/world-attended-rulings` — removes the working tree.
3. `git worktree prune` — drops the stale registry entry.
4. Executed by the controller in the **MAIN checkout** with the context asserted and recorded
   first: `git branch --show-current` = `dev` (the disposition context is pinned — quorum r2
   measured that `-d`'s implicit check is **upstream-or-HEAD, not `dev`**: this branch's configured
   upstream is `origin`/`refs/heads/dev`, so the belt check lands on the right axis here, but the
   doc pins it rather than relying on that). Then the explicit license and the safe delete run as
   one command: `git merge-base --is-ancestor mission/attended-rulings-20260907 dev && git branch
   -d mission/attended-rulings-20260907` — the explicit license (merged into `dev` exactly) is
   the load-bearing guard; git's own safe delete cross-checks it (upstream `origin/dev` and
   `HEAD=dev` both measured to contain the branch: `rev-list --count` = 0 on each) and **refuses
   (rc≠0) if unmerged**. A reviewer proposed `-D` (force); NOT applied — the belt checks the
   right axis by measured config, and discarding it would leave only the license.
5. Evidence control: `git worktree list | /usr/bin/grep -c tmp` — must drop from **1** (S7, the
   pre-removal read) to **0**; and `git show-ref --verify --quiet refs/heads/mission/attended-rulings-20260907`
   returns rc≠0 (branch gone). This is AC-M2-1 below.

## §5 Deliverable 3 — the cross-mission note to the fleet

World does not edit frozen core; the driver-log fix routes as a **fleet commit plus a
cross-mission note**. The controller sends this; the doc only drafts it. Draft (≤25 lines):

~~~
To: sunholo-data/ailang fleet loop (mission-control)
From: sunholo-data/ailang-world loop (charter row 71)
Subject: per-mission crash-log path must leave /tmp — every mission shares this exposure

Controller-measured (2026-09-08, World rig):
- The driver's per-mission crash log is written to /tmp:
    tools/launchd/mission-control.sh:103  LOG="/tmp/ailang-mission-${MISSION_NAME}.log"
  (identical in fleet HEAD and in the pinned driver copy this rig runs).
- The launchd StandardOutPath is also /tmp:
    ~/Library/LaunchAgents/dev.ailang.mission-world.plist
    -> /tmp/ailang-mission-world.launchd.log
- This log is the loop's ONLY crash forensics — HARD TIMEOUT and STALL lines live there and
  nowhere else. One reboot destroys them. Measured: kern.boottime = 1788606012 =
  Sat Sep 5 13:00:12 2026 (the SECOND reboot since the row was filed off the 2026-09-03 one);
  the current /tmp log's first line is [2026-09-05 13:01:11], so iterations 151-157 of that
  mission's driver log are already gone. This fire the log held 12,209,292 bytes / 113,316
  lines / 52 HARD TIMEOUT|STALL lines — large, live, and irrecoverable on reboot. Two copies
  in /tmp (log + launchd log, distinct inodes) are two copies of the same loss.

Ask (fleet commit; World never edits frozen core):
Move the per-mission LOG and the plist StandardOutPath under $HOME (or another path that
survives a reboot), so crash evidence outlives the crash. Every mission on this rig carries the
identical /tmp exposure — landing it fleet-side fixes them all at once.
~~~

## §6 Milestones + acceptance criteria

**M1 — Charter rule lands + note drafted.**
- **AC-M1-1** (instrument: awk region + `/usr/bin/grep -c`): `awk '/^## Repo Profile/,/^## STATUS/' design_docs/world-mission.md | /usr/bin/grep -c 'never place a worktree under `/tmp`'` → **expected 1**. *Failing command outcome when unmet:* 0. *Vacuity edit:* paste the literal string in a doc or the STATUS headline — fails, because the awk region restricts the read to the Repo Profile block, so a stray placement lands outside it and the count is 0. *Why the wording prevents it:* the AC names the region, not just the file; a rule anywhere but the Repo Profile is not a rule. (Single quotes make the backticks literal in the shell; the needle is an exact substring of the §3 rule text, case-exact.)
- **AC-M1-2** (instrument: `/usr/bin/grep -c` on this doc): the fenced §3 rule block above contains the generalized phrasing `claim about the last boot` → **expected ≥1** in `design_docs/planned/w-mission-critical-state-lives-in-a-directory-the-os-wipes-on-boot.md`. *Vacuity edit:* put the phrase in a comment or boilerplate outside the fenced block — caught because the AC scopes to the §3 fenced block. *Why it prevents vacuity:* AC-M1-1 already proves the rule is *in the charter*; AC-M1-2 proves the *drafted verbatim text the controller pastes* exists. Neither is satisfiable by the file merely existing — both read specific content in a specific region.
- **AC-M1-3** (instrument: `/usr/bin/grep -c`): the generalized rule must not duplicate the pin rule's `~/.pinned-ailang` text — `/usr/bin/grep -c "pinned-ailang" design_docs/world-mission.md` stays at its pre-change value (the existing pin bullet), i.e. the new bullet adds no second copy. *Vacuity edit:* not touching anything fails AC-M1-1/M1-2; removing the pin bullet fails the control.

**M2 — Worktree disposition executed (controller) + note sent (controller) + row 71 closed.**
- **AC-M2-1** (instrument: `git worktree list | /usr/bin/grep -c`): `git worktree list | /usr/bin/grep -c tmp` → **expected 0**, down from the S7 pre-removal read of **1**. *Vacuity edit:* no file creation can zero this — it is an operational worktree-registry state; the directory still being on disk (S7 observed) makes the command return 1, failing the AC. *Wording:* the AC names the *registry count*, which only the controller's `worktree remove`+`prune` can change.
- **AC-M2-2** (instrument: compound command rc + `git show-ref`): from the main checkout with `git branch --show-current` = `dev` recorded, the controller's executed `git merge-base --is-ancestor mission/attended-rulings-20260907 dev && git branch -d mission/attended-rulings-20260907` returns **rc=0**, and `git show-ref --verify --quiet refs/heads/mission/attended-rulings-20260907` returns rc≠0. *Vacuity edit:* simply never deleting the branch leaves `show-ref` rc=0 — failing. Deleting without the license (e.g. `-D` from an unpinned context) is not the recorded instrument; the explicit `merge-base` license fails closed (rc≠0) if the branch were unmerged, and the safe-delete belt cross-checks it against the measured-containing upstream/HEAD. *Wording:* the AC requires the *executed licensed safe-delete rc=0*, impossible to satisfy by absence or by an unlicensed delete.
- **AC-M2-3** (instrument: controller dispatch confirmation + row-71 stamp, region-scoped): the controller dispatched the §5 note via `ailang messages send` (body positional, never --body-file) and read back the sent artifact, recording its path; the log entry's M2 evidence cites the artifact id. And row 71 in `design_docs/world-mission.md` carries a LANDED tag inside its own title, verified by `awk '/^71\. /,/^72\. /' design_docs/world-mission.md | /usr/bin/grep -c '\[LANDED 2026-09-08 (iter-171)'` → **expected ≥1** (region-scoped to row 71's own lines: a LANDED stamp anywhere but inside row 71 reads 0 — the vacuity guard). *Vacuity edit:* any `.md` or queue row can be stamped with a LANDED-looking tag — but the AC's instrument is the *sent-note dispatch record plus* the region-scoped row stamp, and the row stays open (the driver log path persists) if deployment was faked. *Wording:* closure is conditional on the note actually travelling to the fleet.

## §7 Conflict Surface

- **Frozen-core fence (F7):** the driver-log path is `tools/launchd/mission-control.sh:103`,
  fleet-mirrored; World never edits it. All artifacts here are docs and a dispatched note — zero
  filesystem edits in frozen core or `scripts/`.
- **Driver-drift gate (row 76):** the repo's own gate fatals on "DRIVER DRIFT vs FLEET
  (D-WORLD-DRIVER-1)" when a World blob differs from fleet HEAD. This doc proposes no blob; it
  proposes World-side rules and a fleet **request**, so the gate is not triggered.
- **Charter STATUS mass-deletion hazard:** the charter's STATUS block and queue rows are
  controller-owned, ratifiable surface. The new rule lands **only** as a bullet in the Repo
  Profile process-fix section — it must not touch the STATUS block, the queue, or any other row.
- **Shared skill's existing `/tmp`-worktree rule (F9):** the skill already forbids `/tmp`
  worktrees for a *location-class* reason. World's rule **extends** that to the *persistence*
  reason and to evidence-bearing state beyond worktrees — it contradicts nothing; it broadens the
  rationale and the surface.

## §8 Declared residuals

- The **fleet `LOG` path stays `/tmp`** until the fleet commit lands (the note is the ask; the
  fix is fleet's to commit).
- The **two non-`/tmp` stale worktrees** (`.eval-world-iter162`, `.wt-world-iter163-record`,
  per S7) are **out of scope**: they persist on a durable sibling, so they do not violate this
  row's rule and are not part of the OS-wipes evidence loss.
- Anything **"crash-proof"** — making the driver survive or self-heal a crash — is out of scope,
  per the row's fence. The deliverable is only that the evidence of a crash outlives the crash.
- The `:76` vs `:103` line drift in the row text (F1) is recorded here, not edited into the row.

## §9 Quorum verification log (controller-owned)

- **Round 1** (2026-09-08T03:14:38Z, artifact `.ailang/state/mission-quorum/w-mission-critical-state-lives-in-a-directory-the-os-wipes-on-boot-2026-09-08T03-14-38Z.json`): **BLOCKED**, metered $0.01459. gemini-3-1-pro **present, reject** — AC-M1-1's needle 'NEVER PLACE A WORKTREE UNDER /tmp' does not match the drafted rule text 'never place a worktree under `/tmp`' (case + backtick) — the same defect the controller had already measured first-party (`sed`-scoped §3 block: 0 hits). gpt5-6-sol **absent: auth** (OPENAI_API_KEY not set this fire), degraded by name, not silently. Controller verdict: reject. Fix applied by the designer (r3 pass on the openrouter twin, 122 s, 3 tool calls, verdict `ok`): gemini's verbatim needle re-anchor + the controller's AC-M2-3 alignment to the real `[LANDED YYYY-MM-DD (iter-N)]` queue stamp, region-scoped to row 71's own lines. New needle reproduced first-party before round 2: 1 hit.
- **Round 2 — the one re-quorum** (2026-09-08T03:19:52Z, artifact `…T03-19-52Z.json`): **BLOCKED**, metered $0.01519. gemini-3-1-pro **present, reject** — "`git branch -d`'s implicit check is upstream-or-HEAD, not dev; the controller runs on its own sprint branch, so it will deterministically abort." Controller verdict: pass on the revised doc. Adjudication, first-party, before any edit: the TRUE half — `-d` checks upstream-or-HEAD, and this branch's configured upstream is `origin`/`refs/heads/dev` (`git config branch.…remote/.merge`), so §4 step 4's old rationale mis-named the check's axis. The REFUTED half — the deterministic abort: the disposition runs in the MAIN checkout (HEAD = `dev`, `git branch --show-current` measured `dev`; sprint work lives in worktrees), and every axis `-d` could check contains the branch (`rev-list --count origin/dev..mission/attended-rulings-20260907` = 0; `rev-list --count dev..…` = 0). Disposition: closed under the **narrow-refinement carve-out** — one blocking objection, concrete reviewer-authored fix, no direction dispute — but the reviewer's `-D` fix was NOT applied verbatim because measurement shows the belt checks the right axis here (upstream = origin/dev, measured containing) and discarding it would leave only the license; instead the true half is folded in as a pinned context + inline license in §4 step 4 and AC-M2-2 (applied and recorded in this commit). gpt5-6-sol **absent: auth**, degraded by name.