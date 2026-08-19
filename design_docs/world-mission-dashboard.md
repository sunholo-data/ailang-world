# Mission Dashboard — AILANG World

*Snapshot, overwritten every iteration. History: `world-mission.md` (STATUS),
`world-mission-status-archive.md`, `world-mission-log.md`.*

**As of** 2026-08-19, iteration **94** · `dev` @ `6c2a537` · CI green (`checks=2`, both `success`
on the MERGE commit, SHA-addressed)

## Just landed

**Row 19 `w-daemon-timeout-test-flake` — COMPLETE** (PR #72 → `6c2a537`; evaluator `sonnet`
**97/100**, zero blocking). The daemon read-deadline tests said "already expired" and spelled it
`1 * time.Nanosecond` — which in Go is a **future** deadline that arms a timer, so a fast read
finished first and the route answered 200. Now one constant, `-1 * time.Nanosecond`, which
cancels synchronously at construction. Base 6/1000 and 3/500 → head 0/2000 and 0/1000.

**The find**: the queue row named ONE test; there are **two** on the one stimulus. And the design
doc's own mutation table declared a two-test red set with "any red outside that set fails the arm"
— the measured set is **four**, so following the doc would have scored a *correct* mutant as a
failed arm. Corrected by measurement, reproduced four times across three roles.

**The reviewer earned their fee**: `gpt5-6-sol` blocked twice on premises nobody had measured.
Running its own proposed fix produced row 22 below.

## Next picks (all unblocked, none needing a human)

1. **20 `w-capsule-output-cap-load-flake`** — two caps race in one fixture; `go test ./...`
   non-deterministic under load. ~0.25 d.
2. **21 `w-archive-stderr-in-manifest`** — `CombinedOutput()` writes an `Observatory:` stderr line
   into the content-addressed archive manifest, served as `interpreter_version`. ~0.5 d.
3. **22 `w-daemon-lock-wait-not-deadline-bound`** *(new)* — a lock-blocked read is bounded by
   `busy_timeout` (~2 s), **not** by the request deadline (measured 2.043 s under a 300 ms
   deadline). Safe today only because 2 s < 10 s, an ordering nothing asserts. ~0.5 d, needs a doc.

## Loop cadence + routing

Controller `claude-opus-5`. Designer **rotation**: `claude:claude-fable-5` used (codex probed
`rc=1`, exhausted until 2026-08-20 05:34). Planner **opus** (`derive-planner-lane.sh` →
`opus fail-closed:env-pin`). Executor **opus** — the chain end; `pi:deepseek-v4-flash-0731` is
now **5-for-5** dead. Evaluator **sonnet** (generator≠judge holds).

## Quota / cost posture

`metered=$0.2342` of the $5 ceiling — two quorum rounds ($0.0955 + $0.1185) and one dead `pi` run
($0.0202). fable/opus/sonnet are subscription buckets.

## Parked on Mark — TWO open asks, both one-word

- **`D-WORLD-19`** — item 17's scope: may tranche 1 extend `host/store` with a bounded object
  read? **A** = yes (adopt the reviewer's fix verbatim) · **B** = no (declare the residual).
- **`D-WORLD-20`** — does `pi:openrouter/deepseek/deepseek-v4-flash-0731` stay in the ratified
  executor chain? **A** = suspend until re-qualified · **B** = keep as ratified. Now **5-for-5**
  zero-completion; this run changed bytes for the first time and still died mid-turn.
