# AILANG World — mission dashboard

*Snapshot, overwritten every Gate 4. History lives in `world-mission.md` STATUS + `world-mission-log.md`.*

**As of** 2026-08-05, iteration 53 · dev @ `13315da` · CI green both jobs (SHA-addressed, step-log verified)

## In flight

- **Item 8 `w-self-mod-vertical` — milestone `SM.A` LANDED** (PR #41 → squash `13315da`, 1,547
  insertions). Projected `world/core@0.1.0`, `host/pkgproj` hash re-implementation, a nine-step
  package gate as `verify_ail.sh`'s THIRD LEG, and AC12's boundary guard.
- **`[NEXT]` is `SM.B1`** — the durability kernel (`approval_claims`, schema 1→2). Gated on
  nothing. Must be **ONE commit**: splitting it lands a red DDL gate. `DD-2`'s ~3× blast radius
  and `DD-3` (`store.go:354`'s bare `return nil`) are binding on it.
- **Item 5 `w-mcp-projection` — still BLOCKED** on one prerequisite (transition registry absent at
  HEAD). Unchanged this iteration.

## Latest

- **The sprint's one gating unknown came back GREEN, on two platforms.** `AC6`'s cross-check agrees
  on `content`/`interface`/`tarball` — World's `go1.25.6` gzip+tar output reproduces the
  `go1.26.5`-built pinned CLI byte-for-byte on darwin/arm64 **and** linux/amd64. Proven a
  measurement, not a co-occurrence: two mutations red one arm each, naming both values, restore
  byte-identical.
- **Queue item 9 went from latent to ACTIVE — a day before anyone looked.** `releases/latest` moved
  to **v0.33.0** on 2026-08-04, so CI job 1 has been verifying `.ail` against an unpinned compiler.
  Measured in the step log at `af0c3b4`: job 1 `v0.33.0`, job 2 `v0.30.0`, same run. **Item 9's own
  row predicted exactly this** and graded itself "latent, not active". A prediction is not a monitor.
- **`DD-7` (found at landing, named by no designer/planner/reviewer):** a byte-exact compiler pin is
  platform-specific, so the single-constant version would have redded CI 100%. Now a per-platform
  table; `compilerSHA256` is machine provenance and is kept OUT of the artifact golden.
- **`AC12` landed with honest limits** rather than a clean claim: its "network confined to
  `host/broker`" control is **vacuous today** (zero `net/http` there — network arrives in SM.B2a).

## Loop · cost · asks

- launchd `mission-world`; controller `claude-opus-5`. Executor **`codex:gpt-5.6-sol`** (4 bounded
  30-min runs); evaluator **`sonnet`** (generator≠judge). Designer/planner **not fired** — doc and
  plan both already landed; designer rotation unchanged.
- Verify profile `ailang-code`; AILANG pinned **v0.30.0** at `/tmp/ailang-v0300/ailang`; the package
  leg now carries its own pinned install in CI via `WORLD_PKG_AILANG_BIN`. Issue **#32**.
- **`metered=$0.00`** vs the $5 ceiling — every role on a quota bucket.
- **Parked on Mark: NONE.** `8/OD-1` ratified; `8/OD-2` open but non-blocking. Worth his attention,
  not blocking: item 9's human-gated half (pin CI job 1 vs keep tracking `latest`) — the cheap
  observability half (`verify_ail.sh` announcing its resolved binary) is recommended first.
