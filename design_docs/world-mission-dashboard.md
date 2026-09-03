# Mission Dashboard — AILANG World

_Snapshot, overwritten every iteration. History: `world-mission.md` (STATUS), `world-mission-log.md`._

**As of:** 2026-09-03 · iter 152 · `dev` = [`6874a98`](https://github.com/sunholo-data/ailang-world/commit/6874a98) · CI green (3/3) · local verify gate green both legs

## Last iteration
**No PR to this repo's code — the fix is upstream.** **REFUTATION** · metered **$0.00** (controller-authored; no sub-agent spawned).

**Row 57's causal claim died, and the issue it kills is one this mission filed itself.** The row said `ailang messages send --type` is dropped, so a "typed sub-query" finds zero `approval_request` rows and `coordinator pending` prints a green all-clear under a live ask. Measured at ailang `origin/dev`: **there is no typed sub-query and one is not expressible** — `printApprovalsInboxPending` applies no type filter (proven empirically: it *listed* a row typed `notification`), and `InboxListOptions` has 10 filter fields, none for message type. The green comes from a **different store**: `FROM approval_requests WHERE status='pending'` (local SQLite) vs the inbox (Firestore). **Two-arm control: at 0 and at 1 unread inbox rows the green line is byte-identical** — provably invariant to the inbox. So fixing the type would change nothing.

**`--type` is real but MISFILED, not ignored:** `messages_send.go:42` binds it to `Category`; `:132` hardcodes `MessageType`. All 18 `approvals` rows split totally by author class (12 mission `notification`/`approval_request`, 6 coordinator `approval_request`/empty). Where it actually hurts: `messages activity` reports **zero** `approval_request` (29 notification, 1 completion) while listing `1 approvals`; plus template routing and sweep classification.

**Settled in the safe direction — the row's own open question:** the Discord push does **not** filter on type (`messageNotification` references `MessageType` 0 times; switches on `ToInbox`; live from `daemon.go:204`). **Mission approvals do reach Discord.** The ask is not lost.

Upstream: correction on [ailang#984](https://github.com/sunholo-data/ailang/issues/984) (0→1 comments asserted) · real cause filed as [ailang#1036](https://github.com/sunholo-data/ailang/issues/1036) · cross-mission note sent.

**Second finding, remediated:** the shared skill cites `make check-no-personal-email` as enforcement — **no such target exists** in either repo, while this PUBLIC repo carried a personal address in 7 doc locations. Redacted (7→0, balanced 7+/7− diff). The missing gate is **new row 74**.

## Goal distance
**Goal unmoved** (no product surface changed). Row census remains **carried, not measured** — row 72 tracks that; three re-derivations disagreed. Row 57 is now **tracking-only** on upstream. Row 50 parked on `D-WORLD-31`.

## Next picks
1. **Row 58** `w-verify-go-is-red-at-pristine-base` — AMENDED: the gate is **flaky**, not deterministically red; wants an instrument-failure floor so a rig cost can't wear a correctness defect's clothes.
2. **Row 59** `w-static-grep-cannot-prove-an-assertion-is-live` — an AC proved "load-bearing" by `grep -c`, which cannot tell a live assertion from a compiled-and-unreached one.
3. **Row 74** `w-the-personal-email-gate-...-does-not-exist` — build the gate the rulebook already claims exists, with a non-vacuity mutation arm. Then 60–66, 68–73, then 39.

## Routing / cadence
Controller `claude:claude-opus-5`. Designer rotation pointer: `claude:claude-fable-5` (not advanced — no designer ran). Verify profile `ailang-code`; pin **v0.30.0** at `~/.pinned-ailang/ailang` (PATH's `ailang` is `-dirty` and is never used for gates).

## Parked on Mark
**`D-WORLD-31`** — ONE WORD: ship rule A as ratified, or hold row 50 for the fixture migration? Re-asked unchanged; **no new ask this iteration**. Re-posted to the approvals spine (prior rows had all been marked read, so the ask was invisible there).

## Quota posture
metered **$0.00** of $5 this iteration. Billing tripwire CLEAN.
