# AILANG World — mission dashboard

*Snapshot, overwritten every iteration. History: charter STATUS + `world-mission-log.md`.*
**As of** 2026-08-15 (iter-88) · **dev** `e53876a` · **CI** green, SHA-addressed, `checks=2` both
`success`. Charter/record only; no code changed, nothing routed. `metered=$0.00` of $5.

## ✅ CORRECTION — THE QUEUE IS **NOT** FULLY BLOCKED. Item 5's blockers have rotted.

Iteration 87's headline is wrong one iteration after it was written. **All three of item 5
`w-mcp-projection`'s P6.B prerequisites are discharged**, measured first-party with firing controls.
Queue is **24 rows** (18 numbered + `4b,4c,4d,4e,4f,6b`), not the 19 last recorded.

## The spine — a blocker's STATE is not its PURPOSE

Iter-87 re-checked `ailang#498` and recorded *"still OPEN, last updated 2026-08-04"*. True about the
issue, and the wrong instrument — Gate 2 says `OPEN` + long-untouched is evidence **toward**
superseded. **Lane B landed in full**: `f5ebcc0b5` M1 (#585) · `6166adab8` M2 (#592) · `b8c038647`
M3 (#601), shipping the **module-root public package `serveapi/`**, released in **v0.33.1**.

It answers the row's four named requirements **four-for-four** — caller-owned mux (`Mount`),
principal resolved before discovery *or* invocation (`SessionResolver` in both the `Tools` and
`Invoke` closures of both handlers), caller-supplied exact descriptors (`ToolSource.Tools`), and
**no built-in tool unless supplied**: `submit_feedback` = **0** across the four seam files while the
same grep in the same directory returns **4** on `feedback_tool.go` (control fires).

Prereq 2 was discharged at iter-75. Prereq 3's basis *"No landed API exposes these"* is now **FALSE**
— `JournalIntent`+`bindCommitIntentTx`, `InvocationID`, `GetReceipt`/`GetEffectReceipt` are all
landed public APIs. **Residual, one word:** *"a **VERIFIED** commit-boundary contract"* — none of the
**10** pinned Z3 identities is a commit-boundary law, so the surface exists and the proof does not.
That is `world/*.ail` work, not more `host/`.

## Why it is still not routable — prescription rot, 3rd instance

The doc chose **path (c)** ("a narrow seam over `internal/apiserver`") *because* upstream exposed
nothing public. Upstream now ships it, so path (c) may reduce to "import `serveapi`". **The clearing
of the blocker IS the falsifier** — invisible to any "has the blocker cleared?" sweep. Item 5 needs a
**design revision** before any sprint.

## ⚠ THREE open asks — item 5's unblocks the most

- **Item 5 (NEW)** — World's `go.mod` requires **only** `modernc.org/sqlite`; there is no ailang Go
  dep to bump, so this is World's **first**. **A** = import `serveapi` at v0.33.1 (smallest, tracks
  upstream, couples the host + re-opens the `.ail` verifier pin `v0.30.0` as a second axis) ·
  **B** = stay dependency-free, build path (c) informed by upstream's contract. Frozen-core cuts for
  **A**, slim-kernel for **B** — which is why it is not the loop's call.
- **Item 17** — one word (evidence-boundary architecture, re-parked after round 4).
- **Item 18** — one word (daemon read cancellation, scope A/B).

## Loop

launchd, ~6h watchdog. Controller `opus`; **no designer/planner/executor/evaluator spawned** this
iteration — nothing routable without the asks. Rotation pointer untouched at `codex:gpt-5.6-sol`.
Upstream `#498` annotated and **deliberately left OPEN**: World's blocker is discharged, its CLI
`serve-api` half is not.
