# AILANG World — mission dashboard

*Snapshot, overwritten every iteration. History: charter STATUS + `world-mission-log.md`.*
**As of** 2026-08-14 (iter-86) · **dev** `6fd26f0` · **CI** green, SHA-addressed, `checks=2` both
`success`, control fires on a known-runs SHA. Nothing landed to dev this iteration.

## Item 18 `w-daemon-read-cancellation` — DESIGNED, then PARKED `needs-human-review`

673-line doc, 9 ACs, 27+2 verification rows, 12 mutation rows. **Two quorum rounds, both
`blocked` — and both with BOTH external reviewers PRESENT** (`absent_reviewers: []` each round,
so neither verdict is a pass-with-a-hole; the per-reviewer cap was raised to $0.25/$0.30 against a
$0.10 default precisely to pre-empt the budget-degrade trap). `metered=$0.2299`.

Round 2 left **one** live objection. gemini's was a PREMISE objection — measured by the controller,
**both premises TRUE** (`handleLogRange` writes once at the end; `defaultClientTimeout = 30s` at
`daemon.go:110`), recorded as V28/V29, design unchanged, objection **discharged by measurement**.
`gpt5-6-sol`'s is a **DIRECTION dispute on the scope boundary**, which forecloses the carve-out.

## The spine: one enumeration, wrong at three different scopes, in one iteration

The queue row said **four** context-free read getters. The designer found a **fifth**
(`GetRegistryHead:628`) — then stated it as a store-wide universal, and the store has **six**
(`GetVerifyResult:773`, off the daemon path). Meanwhile the controller's own directive, under a
**VERIFIED-BY-ME** heading, put the DSN builders in `store.go` — that file has **zero** definitions;
they are at `writer_lock.go:120/176/187`. I had read **call sites as definitions**. The designer
refuted it. Under-count, over-claim and mis-location are the same defect: **a count is only true
inside the scope it was taken in, and the scope is the part nobody writes down.**

## What each item needs now

- **18** — **PARKED**: one-word A/B (A = ship the 1.5 d scoped item · B = re-size to repo-wide).
- **17** — *revision round*: MAC seam, V27 repair (`verify.results[]`, never the int), neg control.
- **14** — blocked on 18. **5** — still blocked.
- New follow-on filed by the doc: `w-bounded-waits-operator-and-write-paths` (item 18's residue).

## Loop · carry-forward

launchd, ~6h watchdog. Controller `opus` · designer `claude:claude-fable-5` (rotation advanced
`codex` → `claude`; **FLAGGED**: two bounded Fable runs, design + revision, vs the one-run
discipline). Cap $5; spent **$0.2299**. **ONE open ask.**
