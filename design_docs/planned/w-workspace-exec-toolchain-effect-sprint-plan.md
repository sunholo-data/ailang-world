# Sprint plan — w-workspace-exec-toolchain-effect (row 140, ATTENDED)

**Authority:** the approved design [w-workspace-exec-toolchain-effect.md](w-workspace-exec-toolchain-effect.md) (rulings D-WORLD-62; round-3 changes approved by Mark 2026-10-05; PR #197 `fccaae9`). Queue row 140 is tagged **IN-SPRINT — ATTENDED**, so the loop does not pick it. Critical path after this row: 136 → 141 → 93 (D-WORLD-61 = A).

**Mode:** attended sessions (Mark + Claude). One PR per milestone, each merged on green CI. Each milestone's mutants (design §7) are run and recorded in its PR before merge.

## Gates (every milestone)

```sh
export AILANG_BIN=$HOME/.pinned-ailang/ailang     # v0.41.0, the gate binary; never a -dirty build
./scripts/verify_ail.sh                           # needs AILANG_BIN too
go vet ./... && go test -count=1 ./...            # go build is not a compile fence for _test.go
go test -race -count=1 ./host/broker/ ./host/daemon/   # CI runs -race
```

The `-race` run covers only the packages a milestone touches. Real-srt tests skip with a stated reason when srt is absent, and every skip has an assertion naming it. CI must run them for real (see M2).

## Milestones

| M | Scope | Design ACs | Mutants it owns | Est. |
|---|---|---|---|---|
| **M0** | Pin srt and measure Linux. No production code. | AC0.1–0.3 | — | 0.5 d |
| **M1** | `exec.ail` transition, `execargs.ail` sketch with contracts, 9th tool descriptor | AC1.1–1.3 | MUT-PREFIX-FROM-AGENT (plan half) | 1 d |
| **M2** | `handlers_exec.go`, `captureHeadTail` on `runBounded`, settings rendering, caches, trampoline, per-call re-verification | AC2.1–2.10 | NO-DASHDASH, NO-TRAMPOLINE, SHELL-JOIN, GRAMMAR-SKIP, PREFIX-FROM-AGENT, ENV-INHERIT, NO-DENYWRITE-DEFAULTS, NO-WORLD-DENY, READ-FENCE-OFF, ROOT-NOT-DENIED, STRICT-REGRESSED, SECOND-LIFECYCLE, SEED-WRITABLE, OVERFLOW-ERROR, NO-GROUP-KILL, STACK-UNVERIFIED, HOST-ENV, TRAMPOLINE-EXEC | 1–1.5 d |
| **M3** | `serve` flags, profile load refusals, placement, startup probe, `/mcp/` end to end, replay, QUICKSTART | AC3.1–3.5 | NET-ON, PROBE-NET-VACUOUS, READROOT-ANCESTOR, PROBE-RC-ONLY, PLACEMENT, ARCHIVE-PKG-ONLY, PROBE-SKIPPED | 0.75–1 d |
| **M4** | Attended (Mark): republish `se-tools`, then three real projects through a real MCP client | M4 list | — | 0.5 d |

Estimate: about 4 days of sessions. The design said 3 + 0.5. The extra half day is buffer for M2, the largest milestone, which owns 18 of the 25 mutants.

### M0 — pin and Linux (PR #206)

- **AC0.1, done on macOS:** `design_docs/verification/world-row140-m0/m0_matrix.py` is one portable matrix. It runs the startup-probe arms 1–6 plus V3/V8/V9/V11/V12/V13/V41/V42 in World's real topology:
  - the workspace root and state dir sit under `$HOME`;
  - srt runs from an archived whole-tree copy, with a PATH-only host env and the `env -i` trampoline.

  Results:
  - **srt 0.0.78: 33/33 gating arms pass.**
  - **Control:** a pass-through (unsandboxed) srt fails all 21 confinement arms.
  - **Pin:** `0.0.78`, `cli.js` sha256 `3c3092bd…d96e` (`pin.json`).
- **AC0.2:** a step in `ci.yml`'s go job runs the same matrix on `ubuntu-latest`, as design §10 planned. It installs `bubblewrap socat ripgrep`, sets `kernel.apparmor_restrict_unprivileged_userns=0`, checks the pinned digest and runs the gating mode. A separate workflow was tried first; `TestGoToolchainPinsAgreeAndMatchJobList` refuses any second workflow file.
- **AC0.3 — DONE.** Linux differed, and the design now records it: V43–V49, P19/P20, an amended §4.4 `allowRead` and §4.6 probe semantics. Run as written, the design fails every Linux arm with rc 127, because srt runs its seccomp helper inside the sandbox. The fix re-allows only `<archive>/…/vendor/seccomp/<arch>` for reading. With it, Linux passes 32/32 gating arms and macOS 33/33. R-140-8 is closed.

### M1 — transition and sketch

- Load the AILANG reference (the `ailang-docs` MCP) before writing any `.ail`. The checks run against the pinned binary only.
- `packages/se-tools/se_tools/exec.ail` follows calling convention v2. It uses `commandIdOk`, `argOk` and `argsOk` as contracted predicates. It emits one effect or a zero-effect refusal (L7). Named tests cover AC1.1's rows.
- `design_docs/sketches/execargs.ail` holds the argument-matching law under Z3. A Go-mirror drift test keeps it in step with the handler (the `effectplan.ail` pattern).
- In `transitions.json`, the descriptor gets 1 triple and `additionalProperties:false`. The workspace registry binds 9 names, and `verify_ail.sh` gets the new identities.

### M2 — handler

- `runBounded` gains a capture strategy. Its zero value is today's strict capture. No existing `runBounded` or handler test may be edited; new arms only (MUT-STRICT-REGRESSED).
- `handlers_exec.go` holds no `exec.Cmd`, `Setpgid`, kill or `procbound` identifier, enforced by a source-scan arm (MUT-SECOND-LIFECYCLE).
- **CI:** the go job in `ci.yml` gains the srt + bubblewrap install, so the real-srt ACs (2.2, 2.3, 2.5, 2.7, 2.8, 2.10) run on Linux at merge. The step is pinned to `pin.json`. A pass-line check, like the fork-lock step's, fails the job if those tests skip.

### M3 — wiring

- There are 4 `serve` flags. The refusal table is bound to QUICKSTART. Each probe arm must be driven red by the pass-through shim from M0's control (`out/macos-passthrough-control.txt` shows the shape).
- `/mcp/` end to end; grant/budget denial; byte-equal replay with a nil registry; the QUICKSTART section.

### M4 — attended (Mark)

1. Republish `se-tools` with `exec.ail`.
2. Write three profiles and seed their caches:
   - the ailang compiler (Go): `test` on one package;
   - TwilightGame (TS): `typecheck`;
   - sunholo-platform/cli (Python): `test-file`.
3. Run each through a real MCP client and record the measured result. An unprofiled command must be refused, and so must a decoy read inside a test.
4. Re-tag row 140 LANDED. That frees the loop to pick 136.

## Registry reuse

`none` for every milestone. srt is an npm runtime dependency the operator installs (pinned + hashed), not an AILANG registry package. `se-tools` is World's own package and gains a transition.

## Risks carried from the design

- **R-140-2:** srt is beta. Re-run the M0 matrix on any pin bump.
- **R-140-8:** Linux. M0 measures it; see AC0.3.
- **R-140-9:** a cold cache can still time out. The result is typed `timed_out`.
- A long-run lane (full suites past 10 s) is a follow-up row, per D-140-4.
