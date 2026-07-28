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
  mission-control` note — no local workarounds, no vendored forks.
- **Never touch** `~/.ailang/state/mission-v1*` or the V1 checkout from work in this repo.

## Operating the daemon

Build/run/commit/read walkthrough (executed-verbatim, S7-maintained): [docs/QUICKSTART.md](docs/QUICKSTART.md).
