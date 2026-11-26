# `w-agent-world-context-and-messages`: Agent World Context & Multi-Agent Messages

**Status:** Planned (r0, filed 2026-11-26 from dogfooding interactive episode `world-agy-interactive-2611261006`)
**Related:** [`w-skill-grants.md`](w-skill-grants.md) (PR #683)

---

## 1. Problem

Dogfooding a real interactive `agy` episode inside AILANG World confirmed that file editing and sandboxed execution (`ailang-read`, `ailang-edit`, `workspace-exec`) are crisp and fast, while exposing two major capability gaps compared to unconfined multi-agent workflows:

1. **No multi-agent messaging inside World:** Confined episodes cannot send or receive AILANG collaboration messages (`ailang messages list|read|send`), preventing confined episodes from coordinating with other episodes, host orchestrators, or AILANG core without falling back to an unconfined host session.
2. **Less historical context than raw `git`:** Confined agents have no read-only access to World's own episode/merge/decision history or their own uncommitted worktree diff (`git status`/`git diff` are excluded from `workspace-exec`), even though `ailang-worldd` already indexes all of this for `/surface`.

---

## 2. Multi-Agent Cooperation: `Agent.Messages` Effect

Mediate AILANG collaboration messages through a typed World effect so confined episodes can cooperate safely with full ledger provenance:

- **Effect:** `Agent.Messages` (scope: inbox name or `*`, budget: N operations per episode).
- **MCP Tools:**
  - `messages-list(inbox?, unread_only?, limit?)` — list inbox messages (read-only).
  - `messages-read(id)` — read a specific message payload.
  - `messages-send(to, title, body, correlation_id?)` — send a message attributed to `episode:<episode-id>` (never forging another identity).
  - `messages-ack(id)` — mark a message handled.
- **Host Mediation & Provenance:** World's broker executes the message store operation on the host, stamps `from_agent = "world:" + episodeID`, and records every send/ack in `world.db` so multi-agent handoffs are traceable via `ailang-worldd why`.
- **Skill Integration (`w-skill-grants`):** The `ailang-feedback` skill declares `metadata.world_grants: [{effect: Agent.Messages, scope: "*", budget: 50}]`, so granting `--skill ailang-feedback` automatically exposes both the skill instructions and the `messages-*` MCP tools across `agy`, `claude`, and `codex`.

---

## 3. World-Native Agent Context (Strictly Better Than Raw `git`)

Expose three small read-only MCP tools (governed by existing `Workspace.Read` grants) so confined agents have *richer* context than raw `git log` / `git diff`:

1. **`episode-status` / `workspace-diff`:**
   - Returns the current episode's modified/added/deleted files and unified diff against its base SHA (reusing `/surface`'s `episode-files.json` diff logic), plus remaining grant budgets and session TTL.
   - Lets an agent audit its own uncommitted work before handing back to the operator or calling `episode finish`.
2. **`world-history` (file / `row-NNN` / episode / PR lookup):**
   - Given a worktree-relative file path, `row-NNN`, episode ID, or PR/commit ref, returns recent World episodes and merges that touched it — including episode summaries, finish-check outcomes, and human `Human.Approve` decisions (context raw `git blame` never captures).
3. **`workspace-search` regex clarity:**
   - Support explicit `regex: true|false` (default `false` literal, or Go RE2 when `true`) and return a clear note if a literal query contains regex metacharacters (`|`, `.*`, `^`) with 0 matches, so `a|b` never silently returns empty.

---

## 4. Phased Rollout

- **Phase 0 (Landed alongside this doc):** Immediate harness UX fixes — `agy-world-gate` read-only allowlist for Antigravity `.system_generated/steps/*/output.txt` spill files, `workspace-exec` listing available commands on missing/unknown `command`, `agy` Turn-1 MCP settle delay bump, and `worldToolMap` tips (≤120-line reads, sequential `go-test`/`go-vet`).
- **Phase 1 (Read-only context tools):** `episode-status` / `workspace-diff`, `world-history`, and `workspace-search` regex option.
- **Phase 2 (`Agent.Messages` effect):** Host-mediated `messages-list|read|send|ack` MCP tools with episode-attributed provenance and `ailang-feedback` skill-grant integration.
