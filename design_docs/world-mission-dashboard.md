# Mission Dashboard — Ailang World

Snapshot: 2026-09-08, iteration 170. History: world-mission-log.md.

- **Row 70 LANDED (reduced)** — `scripts/gate1_range_check.sh` + an 18-arm suite + one additive CI
  step. Gate 1 can now name a zero-check commit that `origin/dev`'s HEAD has moved past. It reports
  `ZERO-CHECK-CANDIDATE` and says *"a controller must look"* — it deliberately makes NO claim about
  push tips, because local git topology cannot establish one.
- **Row 68 CLOSED** — the fleet's de-fork guard (`8c56c5863`) is in and is an ancestor of the pinned
  ref this rig runs; the pin is back ON and confirmed holding first-party (12/12 fleet ACs). Docs →
  `implemented/`.
- **New row 88** — `w-gate1-push-boundary-attribution`, the hard half split out of row 70 at the
  quorum's own signal. Needs an authoritative *persisted* push-boundary record; both the GitHub
  events API and local topology are measured out.

## Queue head (next picks)

1. **69** `w-heartbeat-script-absent` — fleet-gated; predicate re-checked this fire, still false
   (`tools/launchd/mission-heartbeat.sh` absent here, present in the fleet checkout).
2. **71** `w-mission-critical-state-lives-in-a-directory-the-os-wipes-on-boot` — ungated, ~0.2d.
3. **72** `w-queue-closed-count-is-an-increment-chain-no-instrument-reproduces` — ungated, ~0.3d.

Then 73–78, 81–88, then 39. Row 65 DEFERRED per `D-WORLD-33`.

## Loop cadence + routing

- Fires 4-hourly via `dev.ailang.mission-world` (FLEET driver; World's local copy retired at the
  DE-FORK, `D-WORLD-DRIVER-1`). Driver pin ACTIVE at `48c4a6e49`.
- Controller `claude:claude-opus-5`. Designer ROTATION `fable-5-1 → gpt-6-astra → pi:deepseek`.
  Planner per `resolve-role-spawn.sh`/`derive-planner-lane.sh`. Executor `pi:deepseek`.
  Evaluator `sonnet` (generator ≠ judge, enforced on model AND vendor).
- Verify gate: `./scripts/verify_ail.sh` + `go vet ./...` + `go test ./... -count=1`, against the
  pinned `~/.pinned-ailang/ailang` v0.30.0. **Not** `verify_go.sh` — row 76's fleet driver-drift arm
  fatals before `go build` and suspends every Go assertion behind it.

## Standing constraints worth re-reading before a sprint

- `tools/launchd/*` is FROZEN CORE. The shared `mission-control` SKILL.md is the V1 checkout's file
  and World may not edit it — skill gaps are PROPOSALS to Mark + the fleet.
- The pinned v0.30.0 binary's quorum roster does **not** resolve `gpt6-astra` or `oc-glm-5-2`
  (measured iter-170: `unknown-model`). Usable reviewers here: `gpt5-6-sol`, `gemini-3-1-pro`.
- No `timeout(1)` and no `gtimeout` on this rig; `gh api` has no `--timeout`. Login shell is zsh and
  the scripts are bash — invoke them as `bash <script>` or you are measuring the shell.

## Parked on Mark

**Nothing.** Decision ledger: 22 rows, `--check` valid, **ZERO OPEN**.

## Quota posture

Metered this iteration **$0.2322** of $5 (quorum only). Codex/ChatGPT bucket **ration-blocked** on
the last three fires — astra was skipped for that reason (and for a self-review collision). Anthropic
and the `pi:ollama` flat-rate lane both healthy.
