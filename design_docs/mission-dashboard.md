# AILANG World — mission dashboard

*Snapshot, overwritten every Gate 4. History: `world-mission.md` STATUS + `world-mission-log.md`.*

**As of** 2026-08-10, **attended session with Mark** · dev @ `7817d8e` · issue **#53**.

## In flight

- **All three asks DISCHARGED; nothing is parked on Mark.** `9/OD-10` ratified **ACCEPT**; the
  driver export **applied**; `8/OD-1` needed no ruling — answered 2026-08-05.
- **THE QUEUE IS NOT EMPTY.** `SM.D`'s entrypoint and item 9's three pieces are all routable.

## Latest — a milestone whose final action is attended is not an attended milestone

Mark attended to run `SM.D`, the irreversible first publish. **It could not be run, and not for
want of a decision: the code that performs it does not exist.** Measured four ways — no
`Publish`/`Approve` in `cmd/ailang-worldd` outside tests (`registry` is a read-only GET); no
publishing script; the runbook's three ```bash blocks are **all in Stage A**, so Stage B is prose;
and the publish machinery is a library whose `RegistryOrigin` has **no non-test caller**.

**The absence is the safety property.** `registry_publish.go:396-399` demands `https` and refuses
loopback while every caller is an `httptest`. Nothing here can publish, so nothing headless can
trip it. Writing the entrypoint is the deliberate relaxation — headless may **build** it, never
**run** it; `pi` stays barred. **The bookkeeping defect it exposes:** iter-65/66 called `SM.D`
"attended-only, blocked on `8/OD-1`" and declared the queue empty, then picked filler while item 8
held routable work. Rule: *name the ARTIFACT a human owes, or the park cannot be discharged.*

**Stage A ran green anyway** — projection reproduces byte-identically, readiness gate 9/9 equal to
the golden, `verify_ail.sh` 4 identities / 14 tests, `verify_go.sh` 28 `ok` / 0 FAIL / race control
2. Packet reviewed: `world/core@0.1.0`, 4 exports, effects `[]`, 5773 bytes. **No publish occurred;
`world/core@0.1.0` remains unclaimed.** Carry closed: the publisher's success marker is now
**observed** (`⚠ Dry run complete`), not read from upstream source. **Two runbook defects, both
`SM.D` inputs:** step 4 asks the human to confirm digests the gate never emits (the dry-run shows
68 bits, not 256; the real check is mechanical at steps 7/9), and Stage B has no commands at all.

## Loop · cost · asks

- Attended; controller `claude-opus-5`, no executor/evaluator, `metered=$0.00`; **v0.30.0** pinned.
- **Parked on Mark: NOTHING.** Next: `SM.D`'s entrypoint (build-only) · item 9's is-a-release
  assertion + `9/CF-A-1` shim + `9/CF-A-2` (under ACCEPT the DRIFT warning now fires every run).
