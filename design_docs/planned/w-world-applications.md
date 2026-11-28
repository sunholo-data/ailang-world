# w-world-applications (queue row 274): applications on World — r0

**Status:** r0, **filed attended 2026-11-28** (Mark: *"ok lets take this over as we have better
context"*, then *"ok proceed"*). Not quorum-reviewed. The loop's pick-time designer re-measures
§1 and §9 at its base and runs the quorum as usual. Nothing here waives a gate.
**Upstream companion:** AILANG `design_docs/planned/m-world-os-applications.md` r1
(sunholo-data/ailang PR #2493). It records that World owns this layer. The AILANG repo keeps only
a program-trace reader example and a docs page, plus an `ailang init world-app` scaffold once §2's
manifest is ratified.
**Rulings for Mark:** D-WORLD-123 (external-app authority), D-WORLD-124 (first app), D-WORLD-125
(the record apps read). See §7.

## §0 What Mark asked (2026-11-28, Daneel directive task-988c91fa)

> *"an operating system isn't just files, that we replace with agent tasks. It's the applications
> that use those files. We need an equivalent for AILANG world os. What the agents build and what
> applications they interact with is where a human will actually interact with ai. For example
> stapledons godot traces show interactions with blender to make 3d assets … How do we actually
> build on top of the AILANG world agent traces the same way Unix builds on files to offer
> applications."*

## §1 Problem (measured at `91befdab`)

World already has nearly every part an application needs. What it doesn't have is the
application itself.

| Part | What exists | Where |
|---|---|---|
| Tools | An extension package declares transitions in `transitions.json` (`id`, `title`, `description`, `transitionFnFile`, `inputSchema`). The registry serves them as MCP tools, filtered by session capabilities | `packages/se-tools/transitions.json`, `/v1/registry/{name}`, `POST /mcp/` (rows 107/108) |
| How to start an agent | A **preset**: `schema`, `name`, `repo`, `remote`, `base`, `profile`, `agent`, `model`, `ttl_s`, `grants[]` (`effect`, `scope`, `budget`, `set`), `daemon`, `launch` (`mode`, `deadline_s`), `checks[]` | `cmd/ailang-worldd/preset.go:43-75`; files in `~/.ailang-world/presets/` (ten on the Studio) |
| Running it | An **episode**: a confined agent in a worktree under the preset's grants. Every call commits a log entry, `finish` re-runs the checks and opens a PR, and `why` walks results back | rows 217, 218; browser start and finish in row 243 V4 (S9b landed 2026-11-28) |
| What happened | `GET /v1/episodes/{episode}/report`: start, finish, calls by tool and surface, effects with `why`'s verdicts, budgets, refusals, reads, writes, notes | `host/daemon/episode_report.go:50-75` |
| Showing it | `/surface` (address-bar verbs, Needs-you, approvals, repos and presets), AG-UI stream, agent views (row 151, designed), the OS surface (row 188, approved) | D-WORLD-64, -87, -92, -93 |
| Running other software | `Workspace.Exec`: one profile command per call under `srt`, worktree-only writes, no network, `timeout_ms` ≤ 9000 | QUICKSTART §10, `w-workspace-exec-toolchain-effect` |
| Grants for host capabilities | Skill-declared grants and host-mediated effects (`Host.TailnetPreview` first) | row 272 (designed, NEXT after row 243 V4) |

**Missing:**

1. **No "application".** A person can't pick "Blender asset builder" and launch it. They must
   write a private preset by hand (`preset init`), know which grants and checks it needs, and
   then start a generic episode. Nothing carries a tool's purpose, grants, checks and expected
   outputs together, so nothing is shareable or installable.
2. **No declared outputs.** The episode report counts writes. It does not say "this episode
   produced `exports/commons_kit.glb`, validated". An app's result is the thing a person came for.
3. **External software only fits in 9-second pieces.** Measured 2026-11-28: headless Blender
   5.2.1 builds a torus and exports a GLB in 0.61–1.65 s (outside the sandbox), so small
   operations fit `Workspace.Exec`. Renders and real kit builds do not, and Blender under `srt`
   on macOS has not been tried.

## §2 The model: an app is a package with a manifest

An **app** is an AILANG extension package that adds one file, `app.json`, beside its
`transitions.json` (if it has one):

```json
{
  "schema": "world/app/v1",
  "name": "blender-assets",
  "title": "Blender asset builder",
  "summary": "Build, inspect and validate 3D assets as GLB with headless Blender",
  "repo_kind": "blender-studio",
  "preset": {
    "agent": "claude",
    "launch": { "mode": "headless", "deadline_s": 3600 },
    "ttl_s": 14400,
    "grants": [
      { "effect": "Workspace.Read",  "scope": "worktree", "budget": 2000 },
      { "effect": "Workspace.Write", "scope": "worktree", "budget": 500 },
      { "effect": "Workspace.Exec",  "scope": "worktree", "budget": 300 }
    ],
    "exec_profile": "blender-headless",
    "checks": [ { "name": "validate-glb", "run": "make validate-glb" } ]
  },
  "outputs": [ { "kind": "file", "glob": "exports/**/*.glb", "check": "validate-glb" } ],
  "views": [ "timeline", "table", "object-link" ]
}
```

- **The preset part is a template, not a preset.** It holds only the portable fields: agent,
  launch, TTL, grants, exec profile and checks. `repo`, `remote`, `base`, `profile` and `daemon`
  come from the operator's repo profile at launch (row 269), so an app never carries a machine
  path or a credential.
- **Grants never exceed the template.** Launch can narrow them, never widen them. Widening is a new
  app version, which goes through the same propose → verify → commit as any extension package
  (clause 7).
- **`outputs`** names what the app is for. `finish` resolves each glob in the worktree, runs the
  named check, and records the result in the report (§3).
- **`views`** draws only from row 151's closed catalog. An app can't bring its own UI kit.

**Install** = publish the package through the existing registry path. **Launch** = row 243 V4's
start path, with the preset built from the template plus the chosen repo profile. The review
panel shows the resolved grants exactly as it does for any start. An app adds no new runtime,
no new authority path and no new process model.

## §3 The record apps read

Apps, and anything built on them, read the **episode report**. Today the report is an internal
JSON shape. To make it a contract:

- add `"schema": "world/episode-report/v1"` and version it additively, as the log epoch is;
- add `app` (name and version, when the episode was launched from one) and `outputs[]` (`path`,
  `sha256`, `bytes`, `check`, `check_result`);
- keep `/v1/log`, `/v1/receipts/{id}` and `why` as the deeper layers. Apps link to them rather
  than copying them.

This is the "file format" an application over World history can rely on. Row 188's timeline and
cards, row 151's views, and the training-data exports all read the same record. D-WORLD-125 asks
whether this report is the public contract, or whether apps should read the log directly.

## §4 Where it appears on the surface

One place, not a new screen. The surface's start flow (row 243 V4) gains **"from an app"**
beside "from a preset": pick an app, pick a repo profile, write the task, and review. The episode
page then shows the app's name and, once it finishes, its outputs with their check results,
in the outputs area of the front card that row 188 already defines. Installed apps are listed
by the same read verb family as presets (`app list` / `app show`, read-only).

## §5 External applications (A4)

**The first app wraps existing software.** `~/dev/blender` (`sunholo-data/blender`) is a
Blender studio repo with `make example`, `make inspect`, `make render` and `make validate-glb`,
a `blender.sh` wrapper that keeps Blender's user config inside the repo, and committed scenes and
exports. stapledons-godot consumes its GLBs (`art/commons-backdrop-pilot/kit-source.json` records
their bounds and hashes). So the asset pipeline Mark named already exists. It just runs outside
World.

**Authority options (D-WORLD-123):**

| | How Blender runs | New authority | Fits |
|---|---|---|---|
| **A** | `Workspace.Exec` with a `blender-headless` exec profile: `blender.sh --background --python-exit-code 1 --python <script>` under `srt`, worktree-only writes, no network, ≤ 9 s | none | generate, inspect, validate, small exports (0.6–1.7 s measured) |
| **B** | A long-running exec **job**: the same profile and sandbox, started as a job and polled, with a per-profile time limit above 9 s | a longer time limit only | renders, real kit builds |
| **C** | A host-mediated effect per tool (`Host.Blender`), row 272's pattern: the daemon runs the tool on the host with a fixed script directory | a per-tool host effect | anything the sandbox can't hold (GPU, Metal) |

**Recommendation:** start with A, add B when the probe (S0) shows that a real build exceeds 9 s,
and use C only if Blender can't run under `srt` at all. A and B keep one authority model for
every external tool (Blender, Godot export, ffmpeg). C adds a new host effect per tool.

## §6 Sequence

| Step | What | Who |
|---|---|---|
| **S0 probe** | Run one real asset task as a World episode with no new code: `preset init blender --repo ~/dev/blender --base main`, an exec profile for `blender.sh --background`, then `episode start --preset blender --launch` with a task such as "add a low-poly crate generator and export a validated GLB". Record what it needed, where it fell back, how long each Blender call took under `srt`, and what the report missed. The record goes in `design_docs/verification/world-row274-s0/` | attended (Mark at the terminal for `y`) |
| **M1** | `app.json` schema, pure validation in AILANG (template ⊆ grant law, catalog-only views, globs confined to the worktree), `app list` / `app show` read verbs | loop |
| **M2** | Launch from an app: template + repo profile → resolved preset → row 243 V4's start path; "from an app" in the start flow | loop, after row 243 V4 |
| **M3** | Episode report v1: `schema`, `app`, `outputs[]` with check results at `finish` | loop |
| **M4** | The first app: `blender-assets`, per D-WORLD-123 and D-WORLD-124, with S0's task as its acceptance run | loop + attended exit run |
| **Upstream** | After M1 is ratified: ask `sunholo-data/ailang` for `ailang init world-app` scaffolding the ratified manifest | message to the ailang inbox |

The queue position is NEXT after row 272, because M4's options B and C build on row 272's grant
machinery. S0 is attended and can run sooner, since it changes no code.

## §7 Rulings for Mark

- **D-WORLD-123 (external-app authority):** A, A then B, or C (§5). Recommendation and default: A
  then B, as S0 measures.
- **D-WORLD-124 (first app):** `blender-assets` over `~/dev/blender` (recommended). The
  alternatives are a Godot export app over stapledons-godot, or a non-external first app (a
  document or analysis app) that exercises M1–M3 before any external tool.
- **D-WORLD-125 (the record apps read):** version the episode report as `world/episode-report/v1`
  and make it the app contract (recommended), or keep it internal and have apps read the log.

## §8 Non-goals

- No new runtime, process model or authority path beyond D-WORLD-123's choice.
- No app-specific UI. Views come from row 151's catalog, and the shell is row 188's surface.
- No app store. Distribution is the registry, and installation is a package publish.
- No change to presets as they exist. Apps produce presets and do not replace them.
- Not AILANG core work. Language gaps found on the way route to `sunholo-data/ailang` as issues.

## §9 Measurements (2026-11-28, attended, base `91befdab`)

| # | Measurement | Result |
|---|---|---|
| M-1 | Preset shape | `cmd/ailang-worldd/preset.go:43-75` (`launchPreset`, `presetGrant`, `presetLaunch`) |
| M-2 | Episode report shape | `host/daemon/episode_report.go:50-75`; route `host/daemon/daemon.go:1145` |
| M-3 | Tool descriptor shape | `packages/se-tools/transitions.json` |
| M-4 | Exec effect limits | QUICKSTART §10: `srt` 0.0.78, worktree-only writes, no network, `timeout_ms` ≤ 9000 in a 10 s handler cap |
| M-5 | Presets on the Studio | `ls ~/.ailang-world/presets`: ten files (none for Blender) |
| M-6 | Blender present | `/Applications/Blender.app`, `/opt/homebrew/bin/blender`, `blender --version` → 5.2.1 LTS; `/opt/homebrew/bin/godot` present |
| M-7 | Headless Blender timing (outside `srt`) | `Blender -b --factory-startup --python-expr` torus + glTF export: `real 1.65`, then `real 0.61`; 83,280-byte GLB |
| M-8 | Blender studio repo | `~/dev/blender` (`sunholo-data/blender`, HEAD `e4da963`): `blender.sh`, Makefile targets `doctor init example inspect keyframes render validate-glb open` |
| M-9 | stapledons consumes Blender output | `stapledons-godot/art/commons-backdrop-pilot/` (`commons_backdrop_kit_v1.glb`, `kit-source.json`, `kit-notes.md`); commit `0549182c` on Blender masters and headless Godot import |
| M-10 | Not measured | Blender under `srt` on macOS; Metal/GPU inside the sandbox; a real kit build's duration. S0 measures all three |

## §10 Open questions

1. Does an app need its own transitions? `blender-assets` may need none beyond `se-tools` and
   exec, with the app being just the manifest. That would make M1 smaller.
2. Should `outputs` support kinds other than files (a PR, a published package, a decision)?
   r0 has files only.
3. Should several apps be chainable as one goal (A's GLB → B's Godot import)? That is row 188's
   "change from here" and is deliberately out of r0.
