# `w-skill-grants`: Skill-Declared World Grants (Harness-Agnostic)

**Status:** Planned (r0, filed 2026-11-26 from dogfooding interactive episode `world-agy-interactive-2611261006`)
**Applies to:** All agent harnesses (`agy`, `claude`, `codex`, and future adapters)

---

## 1. Problem

Agents running inside confined AILANG World episodes hit two friction points with Agent Skills (`~/.agents/skills/<name>/SKILL.md`):

1. **False skill advertising & blocked reads:** Harnesses like `agy` inject the user's global `~/.agents/skills/` catalog into the system prompt, even though World's confinement blocks `read_file` outside the episode worktree. As a result:
   - Purely instructional coding skills (`adk-cheatsheet`, `blueprint-spec`, `ailang-packages`) cannot have their `SKILL.md` read inside World.
   - Host-side capability skills (`studio-tailnet`, `onepassword-vault`, `herdr-ops`) are advertised to confined agents that lack the underlying host capability, causing wasted turns hitting pre-tool denials.
2. **No ergonomic way to grant a skill's host capability to a specific episode:** When an operator *does* want a specific UI/UX episode to publish a live Tailscale preview (`studio-tailnet`), there is no typed World effect or skill-to-grant binding to authorize that narrow capability without dropping confinement.

---

## 2. Core Design: Typed Effects as Authority, Skills as Ergonomic Bundles

### 2.1 Skill Frontmatter (`metadata.world_grants`)

Skills declare the exact World effects they require in standard YAML frontmatter:

```yaml
---
name: studio-tailnet
description: Publish a web app running on the Mac Studio rig over Tailscale...
metadata:
  world_grants:
    - effect: Host.TailnetPreview
      scope: worktree
      budget: 5
---
```

- **Pure knowledge/coding skills** omit `metadata.world_grants` (or set `[]`). They require no host-side grants beyond read-only skill access.
- **Capability skills** list every World effect they require (`effect`, default `scope`, default `budget`).
- **Authority stays strictly on World Effects:** the World kernel continues to enforce typed effects, scopes, budgets, and ledger provenance. Skills never grant raw shell, raw network, or arbitrary filesystem access.

### 2.2 Harness-Agnostic Skill Delivery & Filtering

Instead of per-harness hook exceptions (which would diverge across `agy`, `claude`, and `codex`), **World itself** mediates skills:

1. **Skill Catalog Indexing (`host/skills`):** At `episode start`, World scans configured skill roots (default `~/.agents/skills/`) and parses each `SKILL.md` frontmatter.
2. **Grant-Filtered Visibility:** Only skills whose `metadata.world_grants` are satisfied by the episode's minted grants are exposed to that episode:
   - Via read-only World MCP tools (`skill-list`, `skill-read`), **and/or**
   - By staging read-only copies of admitted skills into the episode's isolated launch directory (`.agents/skills/` for `agy`, `.claude/skills/` for `claude`).
   - Unadmitted host-ops skills (`onepassword-vault`, `studio-tailnet` when not granted) are never exposed to the episode.
3. **Skill-to-Grant Expansion at Episode Start / Presets:**
   - Launch presets may list `"skills": ["studio-tailnet"]` (pinned in the preset SHA-256 and standing allowance).
   - CLI `ailang-worldd episode start --skill <name>` / `ailang-worldd run --skill <name>` expands `<name>`'s declared `world_grants` into the session's grant list and displays them clearly at the approval prompt.
4. **Mid-Session Escalation:**
   - If a confined agent needs a skill not granted at start, it can request the skill's effect(s) through `Human.Approve` on `/surface`, so the operator can grant it with one tap.

---

## 3. First Host-Mediated Skill Effect: `Host.TailnetPreview`

As the motivating vertical slice for `studio-tailnet`:

- **Effect:** `Host.TailnetPreview` (scope: `worktree`, budget: N registrations per episode).
- **MCP Tool:** `tailnet-preview-serve` (args: `entry` relative to worktree, `name`, `path`, optional `description`).
- **Host Mediation:** World's broker validates `entry` is inside the episode worktree, starts the loopback-only preview process under bounded lifetime (torn down automatically when the episode session expires or is revoked), calls `~/.config/sunholo-platform/homepage/register.py add`, records the provenance entry in `world.db`, and returns the tailnet HTTPS URL to the agent.
- **Security Invariant:** The confined agent never receives raw `tailscale`, raw host shell, or `0.0.0.0` bind authority.

---

## 4. Phased Rollout

- **Phase 1 (Read-only skills, all harnesses):** Add `host/skills` frontmatter parser + `skill-list` / `skill-read` MCP tools (or launch-dir staging) exposing skills with `world_grants: []` uniformly to `agy`, `claude`, and `codex`.
- **Phase 2 (Preset/CLI `--skill` grant expansion):** Expand `metadata.world_grants` in `preset` and `episode start --skill` / `run --skill`, filtering exposed skills to satisfied grants.
- **Phase 3 (`Host.TailnetPreview` vertical slice):** Implement the host-mediated Tailscale preview effect + automatic teardown on session expiry/revocation.
