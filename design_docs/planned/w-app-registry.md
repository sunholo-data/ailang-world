# w-app-registry (queue row 274 M1b): install an app by publishing it — r0

**Status:** r0, **attended 2026-11-28** (Mark: *"yep"*, to starting M1b while row 274 M1, PR #712,
was in CI). Not quorum-reviewed. The loop's pick-time designer re-measures §1 at its base and runs
the quorum as usual. Nothing here waives a gate.
**Parent:** [`w-world-applications`](w-world-applications.md) §6 (M1b) and §11 (M1's decisions).
**Rulings:** one, D-WORLD-129 (§6).

## §0 Why

M1 (row 274, PR #712) defined an app as an extension package with an `app.json` manifest, proved
the manifest law in AILANG (`design_docs/sketches/appmanifest.ail`) and read apps from a local
catalog, `~/.ailang-world/apps/<name>.json`. A local file is not an install. Nobody approved it,
it has no version, World's history does not record it, and nothing stops it changing under a
running episode. The parent design says **install = publish the package through the registry
path**. M1 found there is no package identity to publish into. M1b adds one.

## §1 What exists (measured at `ffd70358`)

| Part | What it does | Where |
|---|---|---|
| Registry heads | A generic `name → object ref` table with compare-and-set (`GetRegistryHead`, `CompareAndSetRegistryHead`) | `host/store/store.go:802-857`, table `epoch_registry_heads` |
| Content-addressed objects | `store.Object{Hash, InterfaceHash, SemanticID, Provenance, Payload}`, hash = sha256 of the payload | `host/store/store.go:593-629` |
| The transition registry | A head `world/transition-registry/v1` over revisions `{SemanticID, InterfaceHash, Revision, Parent, Entries}`. `Publish` checks revision = parent + 1 and parent = the captured head, puts the object, then CAS-advances the head | `host/transitionreg/transitionreg.go:30-36, 237-281` |
| The attended write fence | Refuse an automation environment, require a controlling terminal, then require a typed phrase that names *this* write | `cmd/world-publish/fences.go:275-288`, used by `world-publish transitions` (`transitions.go:89-100`) |
| Reading a head | `GET /v1/registry/{name...}` returns `{name, head}` | `host/daemon/handlers.go:664-692` |
| The app law | `appmanifest.Parse` (closed) and `Validate`, mirrored from the Z3-proven sketch | `host/appmanifest` (PR #712) |

**Missing:** nothing groups descriptors by package, and no head records an app.

## §2 The model

A new registry head **`world/app-registry/v1`** points at an immutable **revision object**:

```json
{
  "schema": "world/app-registry/v1",
  "revision": 3,
  "parent": "sha256:…",
  "apps": [
    {
      "name": "docparse-convert",
      "version": 2,
      "manifest": "sha256:…",
      "transitions": [],
      "published_at": "2026-11-29T10:12:00Z"
    }
  ]
}
```

- **`manifest`** is the ref of a separate object (semantic ID `world/app-manifest/v1`) whose
  payload is the **exact `app.json` bytes** that passed `appmanifest.Validate`. A launch resolves
  those bytes, never a file on disk, so a published app cannot change under an episode.
- **`version`** is a positive integer that increases by one for each publish of the same name.
  The manifest gains no field for it. The registry owns versions, as the transition registry owns
  revisions.
- **`transitions`** lists the descriptor IDs this app's package publishes. Each must exist in the
  transition registry's current snapshot when the app is published. M1's apps have none. The field
  is there so a later app with its own tools is one record.
- **One entry per name.** The latest version replaces the older one in the next revision. Older
  versions stay readable through the parent chain, which is World's roll-forward rollback: to undo
  version 3, publish version 2's manifest again as version 4.
- **Retire** is a revision without the entry. Its history stays in the parent chain.

## §3 Writing it: `world-publish apps` (attended)

```bash
world-publish apps --manifest docs/examples/apps/docparse-convert.json --store ~/.ailang/world/world.db
world-publish apps --retire docparse-convert --store ~/.ailang/world/world.db
```

1. Open the store, read the head, and decode the current revision (absent means revision 0).
2. Read the manifest. `appmanifest.Parse` and `Validate` must pass, and `name` must equal the file's
   basename. Every listed transition ID must be in the transition registry snapshot.
3. Build the next revision: `version` is the entry's previous version + 1 (or 1).
4. **The attended fence**, identical to `world-publish transitions`. The phrase names this write:
   `publish app docparse-convert version 2 to <store>`. A headless loop cannot publish an app.
5. Put the manifest object and the revision object, then CAS the head from the captured ref. A
   lost race refuses with the current head and changes nothing.

Every refusal says what failed, why, and the next command, as M1's `app check` does.

## §4 Reading it

- `ailang-worldd app list|show` read the registry head first and mark each row `published vN`.
  Local catalog entries that are not published show as `local (not published)`. `app check`
  stays a local, file-level tool.
- `GET /surface/apps.json` adds `source` (`registry` or `local`), `version` and `manifest` (the
  object ref) to each row. Still no file paths.
- `GET /v1/registry/world/app-registry/v1` already serves the head with no new code.

## §5 What M2 may launch

This is D-WORLD-129 (§6). The recommendation is **published apps only**. M2 resolves the
manifest object the registry entry names, and the start review shows `docparse-convert v2
(sha256:…)`. The episode start records that ref, so the episode report's `app` field
(D-WORLD-125) names an exact, immutable manifest. Local catalog files stay a draft area for
`app check`.

## §6 Ruling for Mark

- **D-WORLD-129 (what M2 may launch):**
  - **A:** only published apps (registry entries). The start pins the manifest object ref.
  - **B:** published apps, plus local catalog files labelled `local` in the review.
  - **Recommendation:** A. A launch should name something World recorded and an operator
    approved at a terminal, and a mutable file is neither.
  - **Default while unanswered:** A.

## §7 Sequence

| Step | What | Size |
|---|---|---|
| **M1b-1** | `host/appregistry`: revision codec, `BuildNext` (publish, retire), `Publish` with parent/revision/CAS checks, and `ReadSnapshot`, reusing `transitionreg`'s ObjectStore interface | ~400 lines with tests |
| **M1b-2** | `world-publish apps`: the verb, the attended fence and the typed phrase, transition-ID existence checks, and the same refusal table as `transitions` (`R-STORE`, automation env, no TTY, wrong phrase, CAS race) | ~350 lines with tests |
| **M1b-3** | `app list|show` and `/surface/apps.json` read the registry first (`source`, `version`, `manifest`), plus a QUICKSTART §22 publish step run verbatim (attended) | ~300 lines with tests |

No AILANG law changes. The manifest law stays M1's. The revision rules (version + 1, parent equals
the captured head) are small enough for Go table tests, as `transitionreg`'s are.

## §8 Non-goals

- No network publish. The AILANG registry (`world-publish publish`) is for `world/core`, and apps
  are local-first.
- No signatures beyond the attended fence. A typed phrase at a terminal is the authority, as it is
  for transition revisions.
- No automatic install from the local catalog, and no loop-run publish.

## §9 Measurements (attended 2026-11-28, base `ffd70358`)

| # | Measurement | Result |
|---|---|---|
| R-1 | Registry head API | `store.go:802` `SetRegistryHead`, `:832` `GetRegistryHead` (table `epoch_registry_heads`, keyed by `registry_name`), CAS at `:857` |
| R-2 | Revision publish rules | `transitionreg.go:237-281`: revision 1 needs a zero parent; otherwise revision = current + 1 and parent = expected head; object then CAS |
| R-3 | Attended fence | `fences.go:275-288` `requireAttendedLocalWrite` = automation-env refusal + controlling TTY + typed phrase |
| R-4 | Package identity today | only `world/core` has an `ailang.toml` and a ready-packet; `packages/se-tools` and `world-demo` are flat `transitions.json` sets; the transition registry groups nothing by package |
