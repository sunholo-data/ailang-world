# Mission Dashboard — Ailang World

*Snapshot, overwritten every iteration. History: `world-mission.md` (STATUS), `-status-archive.md`, `-log.md`.*
**Iteration 105** · 2026-08-21 · `dev` @ `bd48f68` · CI green (both jobs, SHA-addressed, run confirmed `event=push`)

> **PE.C's anti-vacuity pin was hollow and its drill passed anyway** — the M27 kill depended on an
> UNREACHABLE branch. 2nd consecutive iteration whose real defect was in its own verification machinery.
## In flight
- **Item 17 `w-validated-proven-evidence-boundary`** `[IN-SPRINT]` — `PE.A`–`PE.F` (4.70 d), **three landed**.
- **`PE.C` LANDED** — PR #78 → `bd48f68`. New `host/evidence`: strict canonical `ProofReportV1` (nine
  fields in order) + envelope (`report`, `mac`) codecs, `DecodeProposal` with its 256 KiB pre-parse cap,
  byte caps, AC19 depth pin. Judge `sonnet` **88/100, ZERO blocking**, in its own worktree.
- **NEXT: `PE.D`** (0.92 d) — the largest: validator, sealed mint authority, resolved grade, three
  constructor refusals, **15 mutations**. Then `PE.E` → `PE.F`; `PE.F` last without exception.

## The findings worth carrying
Both ways on the identical tree: as delivered M27 → arm rc=1; with the unreachable `if err == nil` branch
removed the unmutated suite is still rc=0 and the **same mutant survives**. The arm's observable ("some
typed `malformed` refusal") is satisfied by the **trailing-JSON** bystander guard too — repaired by pinning
the scanner's own `exceeded max depth`, not by keeping dead code. The judge's finding then **split**:
`report_codec.go`'s arity guard is genuinely unpinned (mutant builds, suite rc=0, and a tail-truncated
report **panics**, `index out of range [8] with length 8`, in code whose §3.3 mandate is "malformed input
→ typed refusal, never a panic") — now pinned, sole killer; its claim about the **envelope** guard is
**REFUTED** (neutering the whole condition reds `TestEnvelopeStrictRefusals/unknown` on a real assertion).
One undeclared unreachable branch deleted.

## Queue / parked
Rows **22**/**23** headless-routable · **24**–**27** designed-pending · **28**/**29** from iter-104 · **30**
new (iter-105 judge finding 3). Row **5** blocked — `sunholo-data/ailang#764` re-measured today `OPEN`, 0
comments, untouched since 2026-08-17, control answering. **Row 14 is UNBLOCKED** (blocker item 18 COMPLETE).
Parked on Mark: **NONE**; decision ledger **11 rows, 0 OPEN** (`scripts/mission_decisions.sh --check`).

## Loop / routing / cost
Controller `claude:claude-opus-5` · no planner/designer (plan existed) · executor `codex:gpt-5.6-sol`
(probe rc=0) · judge `sonnet`, own worktree (generator≠judge). Fable unspent a **6th** iteration.
`metered=$0.00` of $5 · quota `opus` ×1 / `codex` ×2 / `sonnet` ×1.
## Gates
`AILANG_BIN=/tmp/ailang-v0300/ailang ./scripts/verify_ail.sh` · `AILANG_BIN=… GOTOOLCHAIN=go1.25.6
./scripts/verify_go.sh` — **both exports mandatory**. Pinned `v0.30.0` (`e37b370`). Baseline measured rc=0
on the pristine tree BEFORE the change. `TestHandlerTimeoutKillsTheWholeProcessGroup` load-flaky 2/5.
