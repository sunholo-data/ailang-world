# Mission Dashboard — Ailang World

*Snapshot, overwritten every iteration. History lives in the charter STATUS + `world-mission-log.md`.*

**Updated**: 2026-08-25 (iteration 124) · `dev` @ `612828b` · CI **GREEN** (`checks=2`, both success)

## Where we are

- **Queue head is row 5 `w-mcp-projection`** — live since item 14 completed and `D-WORLD-25`'s
  `Finish 14` was discharged. It is the **sole remaining blocker on M4**, the value gate.
- **The design doc was REVISED this iteration** (Fable designer; 641 → 974 lines) because upstream
  shipping `serveapi/protocol` at `v0.33.2` **falsified its central premise**. The revision is
  sound and re-derived rather than patched.
- **Re-quorum BLOCKED at round 3** — both reviewers present, `absent_reviewers` empty. The blocking
  objection was **confirmed first-party**, not inherited, and it is a scope finding:

  > `serveapi/protocol` carries the whole **A2A** surface and the MCP **envelope** helpers, but
  > **no MCP JSON-RPC dispatch**. Dispatch lives in `mcp_handler.go`, which delegates to the MCP
  > SDK. Importing that SDK costs **+34 packages / 5 module roots / 28 allowlist violations**,
  > including `golang.org/x/oauth2` — an outbound-credential stack in the daemon core, breaching
  > clause 2 and clause 3. Writing our own dispatcher is forbidden by §3.7. Both routes closed.

- **Disposition: SPLIT** (a controller routing call, not a human ask). Upstream ask filed as
  [`ailang#885`](https://github.com/sunholo-data/ailang/issues/885) — `D-WORLD-5`'s own prescribed
  default, the same route that produced `#764` → `v0.33.2`.

## What is executable right now — none of it blocked

| Milestone | Scope | Blocked on |
|---|---|---|
| **`P6.T`** | toolchain floor `go1.25.6 → go1.26.6` (~0.1d, independently mergeable) | nothing |
| `P6.D` | pinned dep + **one** package-path allowlist line + narrowness test (~0.15d) | `P6.T` |
| `P6.V` | the `"verified"` residual — a commit-boundary law in `world/*.ail` (~0.3d) | nothing |
| `P6.B` A2A half | World-authored A2A handler over `protocol` | `D-WORLD-26` |
| `P6.B` MCP half | — | **`ailang#885`** |

## Loop cadence + routing

Controller `opus` ×1 · designer **`fable` ×1** (rotation collapsed onto Fable and FLAGGED —
`codex` is this doc's own quorum reviewer, `gemini` cannot author; diet ceiling of one doc hit
exactly) · planner/executor/evaluator unspent. `metered=$0.1658` of $5 (re-quorum only).

## Parked on Mark

- **`D-WORLD-26` — ONE WORD.** Session credential carrier: **A** = `Authorization: Bearer`
  (recommended) or **B** = `X-World-Session`. Gates only `P6.B`; **the queue does not stall on it.**

## Known drift (not ours to fix)

The **running shared skill is 27 lines behind `origin/dev`** (3,757 vs 3,784; V1 commit
`065a4f16c`) — same as iteration 123, still unrepaired. World cannot edit it (frozen core).
