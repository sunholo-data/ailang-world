# w-world-applications (queue row 274): applications on World — r1

**Status:** r1 + M1 built (attended 2026-11-28; §6 and §11 record what M1 decided), **attended 2026-11-28**. r0 was filed that morning (Mark: *"ok lets take this over
as we have better context"*, then *"ok proceed"*). r1 follows the same day's S0 part 1
measurement and Mark's direction on it (§0). It narrows the row to the part that works on any
host OS. External GPU and device programs (r0's §5 and M4) move to the existing post-1.x rows
**210** (`w-device-effects`) and **214** (`w-program-exec-profiles`), sequenced with row **211**
(`w-world-os-image`). Not quorum-reviewed. The loop's pick-time designer re-measures §1 and §9 at
its base and runs the quorum as usual. Nothing here waives a gate.
**Upstream companion:** AILANG `design_docs/planned/m-world-os-applications.md` r1
(sunholo-data/ailang PR #2493). It records that World owns this layer. The AILANG repo keeps only
a program-trace reader example and a docs page, plus an `ailang init world-app` scaffold once §2's
manifest is ratified.
**Rulings:** D-WORLD-123 RESOLVED (external GPU apps deferred to rows 210/214 with row 211).
D-WORLD-124 (first app) and D-WORLD-125 (the record apps read) are open. See §7.

## §0 What Mark asked (2026-11-28)

The directive (Daneel task-988c91fa):

> *"an operating system isn't just files, that we replace with agent tasks. It's the applications
> that use those files. We need an equivalent for AILANG world os. What the agents build and what
> applications they interact with is where a human will actually interact with ai. For example
> stapledons godot traces show interactions with blender to make 3d assets … How do we actually
> build on top of the AILANG world agent traces the same way Unix builds on files to offer
> applications."*

After S0 part 1 showed that Blender cannot run inside World's macOS sandbox (§5):

> *"yeah its not just blender but all apps, and of course all operating systems - but perhaps
> this is something for when we are running ailang world OS as a layer on top of our linux boot
> up image that is on the roadmap."*

Then, on the recommendation to narrow this row and move device programs to rows 210/214 with
row 211: *"yes please proceed"*.

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

**Missing:**

1. **No "application".** A person can't pick "convert these documents" or "check this package"
   and launch it. They must write a private preset by hand (`preset init`), know which grants and
   checks it needs, and then start a generic episode. Nothing carries a tool's purpose, grants,
   checks and expected outputs together, so nothing is shareable or installable.
2. **No declared outputs.** The episode report counts writes. It does not say "this episode
   produced `out/report.md`, checked". An app's result is the thing a person came for.

**Out of scope here, already on the roadmap:** programs that need a GPU or other devices, such as
Blender, Godot, local models and audio. S0 part 1 (§5) measured why they cannot use World's
sandbox on macOS. Rows 210, 214, 215, 216 and 211 hold that work.

## §2 The model: an app is a package with a manifest

An **app** is an AILANG extension package that adds one file, `app.json`, beside its
`transitions.json` (if it has one). Illustrative example, a document app over the docparse repo
(whose preset and `docparse` exec profile exist on the Studio today):

```json
{
  "schema": "world/app/v1",
  "name": "docparse-convert",
  "title": "Convert documents",
  "summary": "Convert the documents in a folder to Markdown with AILANG Parse and check the result",
  "repo_kind": "docparse",
  "preset": {
    "agent": "claude",
    "launch": { "mode": "headless", "deadline_s": 3600 },
    "ttl_s": 14400,
    "grants": [
      { "effect": "Workspace.Read",  "scope": "worktree", "budget": 2000 },
      { "effect": "Workspace.Write", "scope": "worktree", "budget": 500 },
      { "effect": "Ailang.Run",      "scope": "worktree", "budget": 500 },
      { "effect": "Workspace.Exec",  "scope": "worktree", "budget": 100 }
    ],
    "exec_profile": "docparse",
    "checks": [ { "name": "ailang-test", "argv": ["ailang", "test"], "timeout_s": 900 } ]
  },
  "outputs": [ { "kind": "file", "glob": "out/**/*.md", "check": "ailang-test" } ],
  "views": [ "table", "object-link" ]
}
```

- **The preset part is a template, not a preset.** It holds only the portable fields: agent,
  launch, TTL, grants, exec profile name and checks. `repo`, `remote`, `base`, `profile` and
  `daemon` come from the operator's repo profile at launch (row 269), so an app never carries a
  machine path or a credential.
- **Grants never exceed the template.** Launch can narrow them, never widen them. Widening is a new
  app version, which goes through the same propose → verify → commit as any extension package
  (clause 7).
- **Grants are existing effects only.** An app can't declare a new effect. Device grants arrive
  when row 210 defines them, and apps then use them like any other grant.
- **`outputs`** names what the app is for. `finish` resolves each glob in the worktree, runs the
  named check, and records the result in the report (§3).
- **`views`** draws only from row 151's closed catalog. An app can't bring its own UI kit.

**Install** = publish the package through the registry path (M1b; until then, the local catalog in §11). **Launch** = row 243 V4's
start path, with the preset built from the template plus the chosen repo profile. The review
panel shows the resolved grants exactly as it does for any start. An app adds no new runtime,
no new authority path and no new process model. Nothing in it depends on the host OS.

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

## §5 Programs that need devices: measured, and handed to rows 210/214 with row 211

r0 planned a Blender app as this row's M4. S0 part 1 (attended 2026-11-28,
`design_docs/verification/world-row274-s0/readings.md`) measured why that can't work on the
Studio:

- Outside any sandbox, headless Blender 5.2.1 builds a torus and exports a GLB in 0.61–1.65 s.
- A scratch `ailang-worldd` with a `blender` exec profile refuses to start. The startup probe
  exits 139 (SIGSEGV), and the error names the profile, the arm, the exit code and the stderr.
- srt 0.0.78 alone crashes Blender in the same way with World-like settings, an open filesystem,
  and every Mach service allowed. The sandbox log shows `deny(1) iokit-open-user-client
  AGXDeviceUserClient`. srt's macOS profile allows IOKit only for IOSurface and offers no setting
  to extend it. macOS Blender accepts only `--gpu-backend metal`, and it segfaults on the missing
  device instead of falling back.

**This is general, as Mark said.** Any program that needs the GPU, a camera, audio or a display
meets the same wall on macOS, and every OS needs its own way to grant a device. That is the
problem row 210 (`w-device-effects`) already names. Row 214 (`w-program-exec-profiles`) assumed
srt confinement would hold Blender and Godot. On macOS it can't.

**Where it lands: Linux first, with row 211.** On Linux, srt confines through bubblewrap. Granting
one program the GPU should mean binding the render node (`/dev/dri/renderD*`) into its sandbox
for a bounded time, which is row 210's grant shape. Linux Blender should also run headless
scripts and CPU renders with no GPU at all. **Neither claim is measured.** The Studio has no
Linux VM or container tooling. Row 211's NixOS image spike is where to measure both, before any
device design depends on them. A macOS sandbox profile that admits the GPU (World rendering its
own Seatbelt rules, or an upstream srt setting) stays a later option, not a plan.

**Until then**, Blender and Godot work keeps running outside World, as it does today. This row
adds nothing that would run a program unsandboxed to cover the Mac (D-WORLD-123).

## §6 Sequence

| Step | What | Who |
|---|---|---|
| **S0 part 1** | **DONE 2026-11-28**: Blender under World's exec sandbox, measured. Its evidence moves to rows 210/214 | attended |
| **M1** | **BUILT attended 2026-11-28** (§11): `app.json` schema, the law as `design_docs/sketches/appmanifest.ail` (23 Z3-proven laws, 122 named tests) mirrored in `host/appmanifest` under a drift test, a local catalog `~/.ailang-world/apps`, `app list` / `app show` / `app check`, `GET /surface/apps.json` and the `app list` bar verb | attended |
| **M1b** | Install through a registry identity: a `world/app-registry/v1` head (name, version, `app.json` digest, the descriptor ids the app's package publishes), written by an attended `world-publish apps` verb and read by M1's verbs in place of the local catalog. Nothing in World groups descriptors by package today, so this needs its own short design | loop or attended |
| **M2** | Launch from an app: template + repo profile → resolved preset → row 243 V4's start path; "from an app" in the start flow | loop, after row 243 V4 |
| **M3** | Episode report v1: `schema`, `app`, `outputs[]` with check results at `finish` | loop |
| **M4** | The first app, per D-WORLD-124 (recommended: `docparse-convert`), with an attended exit run: install, launch from the surface, finish, outputs checked in the report | loop + attended exit run |
| **Upstream** | After M1 is ratified: ask `sunholo-data/ailang` for `ailang init world-app` scaffolding the ratified manifest | message to the ailang inbox |

The queue position stays NEXT after row 272, as filed. r1 removes the external-program slice,
the only part that needed row 272's grant machinery. M1 and M3 now depend on nothing beyond
landed work, and M2 depends on row 243 V4. Moving the row up is Mark's call.

## §7 Rulings

- **D-WORLD-123 (external-app authority): RESOLVED 2026-11-28, deferred.** Programs that need
  devices move to rows 210/214, sequenced with row 211's Linux image. No unsandboxed host effect
  is built here. Blender and Godot work stays outside World until then.
- **D-WORLD-124 (first app), OPEN:** recommended `docparse-convert`, a document app with no device
  needs (§2). Alternatives: an `ailang-parse` package check, or a repo analysis/report app.
- **D-WORLD-125 (the record apps read), OPEN:** version the episode report as
  `world/episode-report/v1` and make it the app contract (recommended), or keep it internal and
  have apps read the log.

## §8 Non-goals

- No new runtime, process model or authority path, and no new effects.
- No device, GPU or display access. That is rows 210/214/215/216, with row 211.
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
| M-5 | Presets on the Studio | `ls ~/.ailang-world/presets`: ten files, including `docparse.json` (agent `agy`, base `main`, profile `docparse`) |
| M-6 | The docparse exec profile | `~/.ailang-world-exec/profiles/docparse.json`: commands `ailang-check`, `ailang-test` |
| M-7 | Headless Blender timing (outside `srt`) | torus + glTF export: `real 1.65`, then `real 0.61`; 83,280-byte GLB |
| M-8 | Blender under World's exec sandbox (S0 part 1) | scratch daemon probe `exit 139`, startup refused; srt alone (`world`, `open`, all Mach lookups): a crash each time in `MTLBackend::metal_is_supported`; sandbox log `deny(1) iokit-open-user-client AGXDeviceUserClient`; `--gpu-backend` accepts only `metal` |
| M-9 | Linux tooling on the Studio | none of `docker`, `orb`, `limactl`, `colima`, `podman`, `tart`, `utmctl`, `qemu-system-aarch64`, `multipass`; no OrbStack, Docker, UTM or Parallels app |
| M-10 | Not measured | Blender under bubblewrap on Linux, with and without a render-node bind (row 211's spike) |
| M-11 | M1 law and gate | pinned v0.41.0: `ailang verify sketches/appmanifest.ail` → 23 exported functions, 23 verified; `ailang test --format json` → 122 named tests pass; `verify_ail.sh` with an isolated `AILANG_STATE_DIR` → PASSED with Leg 2d "all 122 app manifest sketch named tests pass" and 27 swept modules; mutants killed: a Go budget bound (drift red), an extra Go effect (4 tests red), a body that disagrees with its ensures (Z3 VIOLATION), a flipped row (drift red) |

## §10 Open questions

1. Does an app need its own transitions? `docparse-convert` may need none beyond `se-tools`,
   `Ailang.Run` and exec, with the app being just the manifest. That would make M1 smaller.
2. Should `outputs` support kinds other than files (a PR, a published package, a decision)?
   r1 has files only.
3. Should several apps be chainable as one goal (A's output → B's input)? That is row 188's
   "change from here" and is deliberately out of r1.

## §11 What M1 decided (attended 2026-11-28)

1. **The law is a sketch, not a package** (S3's question). Manifests are validated at the host
   boundary before any episode exists, as the exec-argument law (`sketches/execargs.ail`) and the
   effect-plan law are. AILANG states and proves it, Go runs it, and a drift test replays every
   inline row through the Go mirror. A package root would also have moved the isolated-copy
   census in `host/verifygate`. The app's own runtime parts (transitions, views) stay package code.
2. **Install is a local catalog for now.** Neither registry path has a package identity: only
   `world/core` goes to the AILANG registry, and the transition registry is a flat descriptor set.
   M1 reads `~/.ailang-world/apps/<name>.json` the way presets are read. M1b replaces it with a
   registry head.
3. **Checks use the preset grader shape without `unsandboxed`.** The closed decoder refuses that
   field, so an app cannot ask for a check outside the sandbox. A missing `timeout_s` means the
   preset default of 900 s, and the law caps it at 3600 s.
4. **Effects are the ten episode effects**, bound three ways: the sketch's `ensures`, the Go set,
   and the broker's effect constants (`TestAllowedEffectsAreBrokerEffects`). `Model.Infer`,
   `Human.Approve` and the `FS.*` effects are refused. Adding one is an edit to the law.
5. **The surface never shows a file path.** `/surface/apps.json` rows carry the name, validity,
   template summary, digest and problems. The CLI's `app show` prints the path.
6. **`NarrowGrants` exists now** (same effect and scope, budget 1 to the template's), so M2 can
   build launch on a tested law.

