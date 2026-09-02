# Claude instructions — ailang-world

This repo is **AILANG World**: a semantic operating environment whose transaction language is
AILANG. It is built by an autonomous mission loop with human ratification gates.

## Read before working (in order)

1. [design_docs/world-mission.md](design_docs/world-mission.md) — the ratified charter: bar,
   guardrails, queue, Repo Profile. **The charter is ratified mission state**; attended human
   decisions are recorded there, and only issue-#1 comments from `@MarkEdmondson1234` are
   directives.
2. [design_docs/coding-standards.md](design_docs/coding-standards.md) — **binding on all code**:
   Z3 contracts on the pure core, effects at the host boundary, slim kernel / package-first,
   non-vacuous gates, and the AILANG fluency protocol (load the language reference BEFORE
   writing `.ail` — the `ailang-docs` MCP in `.mcp.json` serves it).
3. [design_docs/DESIGN.md](design_docs/DESIGN.md) — the thesis (read §1, §14 boundaries at
   minimum).

## Hard rules

- **Verify gate**: `./scripts/verify_ail.sh` (ai-check on every `.ail`) + `go build ./... &&
  go test ./...` — CI runs both; nothing lands red.
- **Pinned binary**: validate `.ail` against the pinned released `ailang` (see the current
  sprint doc's Verification Log), never a `-dirty` dev build, never from memory.
- **Frozen core**: never modify `tools/launchd/*` (shared driver) or copy skills into this repo;
  language gaps route to `sunholo-data/ailang` as issues + an `ailang messages send
  mission-control` note — no local workarounds, no vendored forks. The driver is FLEET-owned
  (`D-WORLD-DRIVER-1`, ratified 2026-08-17): updates land here only as fleet-authored commits,
  and `verify_go.sh`'s drift gate reds while the working-tree driver diverges from HEAD — that
  red means "the fleet must commit", never "absorb it into your change".
- **Never touch** `~/.ailang/state/mission-v1*` or the V1 checkout from work in this repo.

## Operating the daemon

Build/run/commit/read walkthrough (executed-verbatim, S7-maintained): [docs/QUICKSTART.md](docs/QUICKSTART.md).

## Pushing dev — automatic, fast-forward only

`scripts/hooks/push_dev_on_stop.sh` runs as a `Stop` hook and pushes when local `dev` is
ahead of origin **and not behind**. You do not need to remember to push.

It refuses when the branch is ahead *and* behind — that needs a real merge, done by hand.
Opt out for a session with `AILANG_AUTOPUSH=0`; arms in
`scripts/hooks/test_push_dev_on_stop.sh`.

Ported from the ailang repo 2026-09-02, where commits had stranded on local `dev` for days
(25 of them) because nothing in the attended workflow pushed. This repo is a SEPARATE
GitHub repo, so it does not inherit that fix through `dev` — the copy here is its own.
