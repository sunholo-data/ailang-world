# Mission Dashboard — Ailang World

*Snapshot, overwritten every iteration. History lives in `world-mission.md` (STATUS),
`world-mission-status-archive.md` and `world-mission-log.md`.*

**As of**: 2026-09-06 · iteration **159** · `dev` @ [`2115172`](https://github.com/sunholo-data/ailang-world/commit/2115172) · CI **GREEN 3/3**

## Latest landing

- **Row 63** — `w-locator-derivation-refusals-are-unpinned-and-undeclared` ·
  PR [#121](https://github.com/sunholo-data/ailang-world/pull/121) → squash `2115172`.
  The row-52 step locator's two derivation refusals were deletable with the whole suite green.
  Both are now armed via an extracted `stepBlockAnchors`, and the branch the row left unfired
  turns out to be reachable from a **re-style that parses deep-equal to the pristine `ci.yml`**.
- Rows **58–63** have landed on six consecutive iterations.

## Next picks

| # | item | size |
|---|---|---|
| **64** | `w-fleet-residual-net-shares-phase-1-pathspec` | ~0.1d |
| 65 | `w-go-build-is-not-a-compile-fence-for-a-test-file` | ~0.1d |
| 66 | `w-flow-key-quote-trim-is-uncovered` | ~0.1d |
| 68–78, 81, 82, 83 | remaining clause-2 gate-hardening rows | small |
| 39 | next non-gate-hardening item | — |

Rows **79 / 80** are `[PARKED — DESIGN REVIEW]` by their own text.

## Loop cadence & routing

- Controller: `claude:` CLI, opus, quota bucket. Metered budget **$5/iteration**;
  recent iterations spend **$0.00** (small rows are controller-authored direct fixes).
- Designer `claude:claude-fable-5` · planner `opus` · executor `codex:gpt-5.6-sol` ·
  evaluator `sonnet` — **not spawned** for ~0.1d rows; when that happens the record states
  plainly that generator == judge.
- Verify gate: `scripts/verify_ail.sh` + `go build`/`go vet`/`go test ./...` with `AILANG_BIN`
  set to the pinned **v0.30.0** at `~/.pinned-ailang/ailang`.
  `scripts/verify_go.sh` is **rc=1 at base** on the FLEET-owned driver-drift arm (row 76) —
  only a fleet commit clears it, never a World edit.

## Parked on Mark

**Nothing.** Decision ledger: **18 rows, `--check` valid, ZERO OPEN.**

## Known standing conditions

- `scripts/mission_directives.sh` and `tools/launchd/mission-heartbeat.sh` are absent from this
  repo (row 69) — directives are read via the V1 checkout's copy by absolute path, and no
  per-gate heartbeat stamp fires here.
- `tools/launchd/mission-control.sh.tmp.astra` sits untracked in the shared main checkout: a
  fleet artifact, frozen core, left alone.
