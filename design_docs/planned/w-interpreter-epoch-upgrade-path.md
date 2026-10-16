# w-interpreter-epoch-upgrade-path (queue row 139, widened) — r1

**Status:** r1 draft, attended 2026-10-16 (Mark attending). Row widened by Mark: "world needs to be able to cope
with ailangs versioning - we release very often"; sequencing "widen 139, do it next". Ratification is Mark's.

## §1 Problem (measured 2026-10-16)

AILANG shipped 8 releases in 6 days (v0.54.4 → v0.54.11). World hard-pins two AILANG binaries and cannot move
either as a normal operation.

| id | measurement |
|---|---|
| V1 | Restarting an existing store with `serve --ailang-bin <v0.54.11>` (copy of the row-182 rig store) exits at startup: `daemon startup failed at registry-bootstrap: … registry: existing head "sha256:c560…" diverges from bootstrap revision "sha256:a4ae…"`. `registry.bootstrap` (`host/registry/registry.go:125-171`) only builds `{epoch 1, candidates:[served release]}`; `firstEpoch` is fixed (`:34`); nothing appends an epoch or a candidate. |
| V2 | World's AILANG code is compatible with v0.54.11: `verify_ail.sh` under v0.54.11 → all 23 modules ai-check + Z3, 16/16 required identities, 40 kernel + 138 execargs + 444 se-tools named tests pass. Only leg 3 fails: `verify_world_package.sh:213-226` parses the dry-run `Tarball:` line, which v0.54.11 indents. Leg 1 prints `UNRECOGNISED RELEASE` (notice, from `scripts/testdata/ailang_release_observed.txt`). |
| V3 | The coordinator runs the interpreter named by the bound transition DESCRIPTOR, never the served pin (`coordinator.go:322-325,341-343`; `plan.go:109,129,139-145` copies it into the entry); `capsule.go:198-202` resolves `entry.Interpreter` via `archive.Resolve` and re-hashes it. |
| V4 | `world-publish transitions` binds a revision to the epoch whose `Candidates` contain the archive pin's release (`transitionreg/epoch.go:62-84`); refuses a release no epoch nominates (`transitions.go:236-239`), and needs `semanticsEpoch` when several match (`:246-249`). |
| V5 | Replay resolves each entry's OWN interpreter hash from `<db>.artifacts/interpreters/<algo>/<digest>/` (`replay.go:188-199`, `archive.go:184-204`); there is no GC, so archived binaries persist across a pin change. Goldens re-verified on v0.41.0; a two-interpreter replay test exists (`replay_test.go:25-28,514,574`). |
| V6 | The tool binary is an exact-match pin: `ToolBinaryRelease = "AILANG v0.52.1"` (`host/daemon/workspace.go:44`), refused otherwise (`:299-300`). Confinement behaviour is measured on v0.52.1 (restricted defaults, `.git` protection, fs_deny_write, `pkg/...` resolution; `handlers_ailang_run.go`, `handlers_ailang.go`, `workspace_layout.go`; matrix V32–V80). |
| V7 | Release literals are scattered: `pins.go:30,39-52` (release + per-platform digests; used by `setup`/`doctor` only), `fences.go:91` (`frozenCompilerVersion`, ready-packet recompute `:364`), `demo_proof.go:35,104-107`, `demo_seed.go:90,132`, `verify_world_package.sh:25-35,278-279,286-288`, `ci.yml:77-106`, `ailang_release_observed.txt`, `workspace.go:44`, plus many tests. `TestPinTableCoversBothReleases` (`setup_test.go:231-275`) cross-checks only some. |
| V8 | No code reads a project's declared AILANG range (`ailang.toml` `ailang = ">=0.50.0"` in ailang-parse; World's own `packages/world-core` says `>=0.30.0`); `pkgproj.go:47` keeps it only as an interface-hash input. |

## §2 Design

**D1 — one release manifest.** `pins/ailang-releases.json` (new) is the single source of truth: per release,
the version line, per-platform tarball + binary sha256, and roles (`interpreter` candidate, `tool`, `proof`),
each role with the date and evidence of its qualification. `pins.go`, the demo/proof literals, `fences.go`'s
frozen compiler, `verify_world_package.sh`'s table, `ci.yml`'s install step and `ailang_release_observed.txt`
are generated from it or checked against it by one test that covers EVERY literal in V7 (a release string in
production code that the manifest does not explain fails the test). A release bump is one manifest edit plus
`make pins` (regenerate), in one PR.

**D2 — nomination (the registry can grow).** A new operation appends an immutable registry revision (CAS on
the head, same object mechanism as bootstrap), in one of two shapes:
- **candidate** — add release R to the CURRENT epoch's `Candidates`. Allowed only when D4's compatibility gate
  passed for R, including every replay golden reproducing **byte-identically** under R (same semantics).
- **new epoch** — append epoch N+1 with `Candidates:[R]`, when the goldens differ (semantics moved). Old
  entries keep their epoch and their own interpreter (V5).
Verb: `world-publish epoch nominate --store <db> --release-manifest pins/ailang-releases.json --release R
(--candidate | --new-epoch) --gate-verdict <file>`. It refuses without a passing verdict for R whose binary
digest matches the manifest. Who may run it is **decision D-139-1** (§4).

**D3 — startup accepts any nominated release.** `registry.bootstrap` keeps its fresh-store path (epoch 1).
On an existing head, startup loads the head and requires the served interpreter's release to be a candidate
of some epoch; otherwise it refuses with the exact `epoch nominate` command to run (still fail-closed). The
selected epoch for publishing is the one nominating the served release (V4 already handles several).

**D4 — the per-release compatibility gate.** `scripts/ailang_release_check.sh <release>` (new) downloads the
release, verifies its sha256 against the published `.sha256`, writes a manifest candidate row, then runs, with
the candidate as both binaries where applicable: `verify_ail.sh` (all legs), the Go tests that execute the
interpreter or the tool binary (replay goldens on the candidate; the tool confinement matrix V32–V80; fences
ready-packet recompute), and a two-interpreter replay of an existing store copy (old entries on their own
archived interpreter, new entries on the candidate). Output: a verdict JSON (per check: pass/fail, evidence),
`golden_identical: true|false` deciding candidate vs new epoch. The loop runs it on each new AILANG release
(fleet release notice or a daily `gh release list` check) and opens a PR with the manifest row + regenerated
pins when it passes; a failing release is filed as an ailang issue with the verdict, never worked around.

**D5 — tool binary by role and range.** `ToolBinaryRelease` becomes "the manifest's releases with role
`tool`" (exact digests). `serve --tool-ailang-bin` stays one binary per daemon in r1; at session/episode
start the workspace reads the project's `ailang.toml` `ailang` range (if any) and **refuses the episode** when
the served tool release is outside it, naming the release that would satisfy it. Several tool binaries per
daemon, selected by range, is a later step (stated, not built here).

**D6 — tolerant gate parsers.** `verify_world_package.sh` leg-3 parsing accepts the label with any leading
whitespace; each parser has a fixture per qualified release.

**D7 — the first move.** Run D4 on v0.54.11 (or the latest at the time), fix what it finds (D6 at least),
nominate, republish se-tools transitions on the new epoch/candidate, move both pins via the manifest, re-run
the FINAL-shaped checks that depend on the pins (verify gate, replay goldens, a floor smoke cell). The v0.41.0
and v0.52.1 binaries stay archived and in the manifest for history.

**Why not a package (S3)?** The registry, startup and gate are host/kernel-boundary operations; no `world/` change.

## §3 Acceptance criteria

| AC | criterion |
|---|---|
| AC1 | One test fails if any production release literal is not explained by the manifest (covers every V7 site). |
| AC2 | `epoch nominate` appends a candidate or a new epoch as an immutable revision; refuses without a passing, digest-matching verdict; CAS conflict refused. |
| AC3 | An existing store restarts on a nominated release (V1 now passes); on an un-nominated release it refuses, printing the nominate command. |
| AC4 | Old entries replay byte-for-byte on their own archived interpreter after the move; new entries replay on the new one (two-interpreter replay on a real store copy). |
| AC5 | `ailang_release_check.sh` produces a verdict for v0.54.11 with every check named; a deliberately broken candidate (wrong digest, failing golden) is refused. |
| AC6 | An episode whose project declares a range excluding the served tool release is refused before any tool call, naming the satisfying release. |
| AC7 | The first move (D7) lands: both pins on the new release, se-tools republished, verify gate + replay goldens + one floor smoke cell green. |

## §4 Decisions for Mark

- **D-139-1 — who may nominate a release?** (a) **the loop**, when D4's verdict passes and the release digest
  matches GitHub's `.sha256` — needed for a daily cadence; or (b) human-typed only, like `publish` and session
  mint today (QUICKSTART §0) — safest, but one typed phrase per AILANG release. Recommendation: (a) for
  **candidate** nominations (same semantics, goldens identical), (b) for **new epochs**.
- **D-139-2 — first target release:** v0.54.11 (measured compatible apart from D6), or the latest on the day.

## §5 Risks

| risk | mitigation |
|---|---|
| A release changes semantics silently | goldens byte-identical or it is a new epoch, never a candidate |
| Tool confinement regresses on a new release | the V32–V80 matrix runs on the candidate in D4; a failure blocks the tool role |
| Daily churn of PRs | one manifest PR per qualified release; failing releases filed upstream, not bumped |
| Literal drift | AC1's exhaustive test |
| History unreplayable | binaries never removed from the archive; AC4 on a real store |

## §6 Unknown

- Replay goldens' behaviour under v0.54.11 (D4 measures it: identical → candidate, else new epoch).
- How many Go tests hard-code v0.41.0/v0.52.1 behaviour beyond the literal (D4's first run lists them).
