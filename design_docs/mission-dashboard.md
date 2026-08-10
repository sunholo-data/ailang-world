# AILANG World — mission dashboard

*Snapshot, overwritten every Gate 4. History: `world-mission.md` STATUS + `world-mission-log.md`.*

**As of** 2026-08-10, **iteration 67** · dev @ `a4452d1` · CI green both jobs · issue **#53**.

## In flight

- **`SM.D0` LANDED** (PR #55 → `a4452d1`) — the attended-publish entrypoint exists. **`SM.D` is now
  a real procedure and is ATTENDED-ONLY**: never headless, never CI. Item 8 has no headless
  milestone left.
- **`[NEXT]` is item 9's three pieces**, routable today under the `9/OD-10` ACCEPT ruling.
- **Parked on Mark: NOTHING. Zero open asks.**

## Latest — the deliverable was a fence, not a feature

Until this commit the *absence* of a caller was the safety property: `grep -rhoE 'Publish|Approve'
cmd/` = **0** (control: 27 `func `), and the production constructor refuses loopback while every
caller was an `httptest`. `SM.D0` deliberately relaxes that, so its centre is the fence.

**The fence is a controlling-terminal check** — chosen because it is the one thing this loop is
structurally unable to satisfy (stdin is a socket; `open(/dev/tty)` → *device not configured*).
`/dev/null` **is** a character device, so a naive isatty would admit `--live < /dev/null`; hence the
`os.SameFile` branch. `R-CI` is a declared tripwire, **not** the fence. 14 refusal branches, 22
mutations, 22 killed. Stage B of the runbook now carries commands — it was prose, which is how its
step-4 defect (confirm digests the gate never emits) survived a milestone.

**The finding, and it came from the judge.** `MUT-D0-FENCE-ORDER`'s non-vacuity claim was **false**
in three places — *"every AC21 row still passes; killed only by AC22"*. Reproduced first-party:
**6 of 15 rows red** (`R-CI`, `R-TTY-OPEN/CHARDEV/SAMEFILE`, `R-PHRASE-EOF/PHRASE`), because those
rows' fixture uses a **loopback** origin so a hoisted constructor refuses first. AC22 still earns
its place; it is not the unique killer. Corrected in code and commit message. **Spine: a
non-vacuity claim never run as literally described is a vacuous pass — and it hides best when the
code it describes is correct.**

**Also:** two landed defects repaired — `TestRunbookStageAPerformsNoPublicWrite` was structurally
vacuous (its detection loop had never executed its body), and `protectedGoGroups` let any new
`cmd/` package escape the boundary gate. **No publish occurred; `world/core@0.1.0` unclaimed.**

## Loop · cost · asks

- Planner+executor `opus` (**`pi` barred** — publish-capable), evaluator `sonnet` **88/100 zero
  blocking**, `metered=$0.00`; **v0.30.0** pinned. Carry `8/CF-D0-1`: the fence guards the call
  site, not `host/broker` — a direct broker call from other Go code skips all 14 branches.
