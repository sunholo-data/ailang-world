# Executor report — iteration 268, row150 M2a

Implementation ready for independent evaluation; this is not an evaluator verdict, full-suite verdict or completed demo/release claim.

Worktree: `/Users/voightkampff/dev/sunholo-data/.wt-world-iter268-run2`.
Branch: `sprint/world-iter268-row150-m2a-run2`; base `39066a46c1ac13c994bacb162b8d4b754d5ec083`.
Actual declared executor model: native GPT-6.1 Sol (`codex:gpt-6.1-sol`); role token usage not reported. Operator-authorized native Agent transport replaces the recipe; path deviation FLAGGED. No subagents, inbox operations or git writes. Controller owns plan/bookkeeping edits, commits, full candidate gates and the fresh independent evaluator.

## Scope and behavior

Only executor implementation files: `cmd/ailang-worldd/demo.go`, `demo_seed.go`, `demo_test.go`, `main.go`, and `docs/QUICKSTART.md`. The two changed sprint-plan artifacts are controller-owned clarification edits, not executor edits. Source copies with relative paths are banked under this report directory.

Adds `demo seed`, global/local --addr and help, all designed flags, normal session-file resolution, preflight and S1 only. It commits five real objects through session-authenticated sole-writer REST: goal, two distinct candidate descriptions, unchanged production logepoch source and explicitly demo-authored claim. It emits valid frozen six-field headers with real archived pin, payload/interface hashes, carrying entries, observed-head CAS and stored worlds. It never mints, publishes, creates a worktree, writes a competing store, produces forged proof/approval objects or writes a success manifest. The incremental command returns exit1 with the explicit executed S1 / S2-S6 incomplete boundary.

Before writes: absent manifest; daemon/db identity; pin hash+version; known live demo session; all four required grants with exact scope, live expiry and sufficient remaining budget after existing EffectSpend; empty public log/head and no unselected World rows. The controller approved and documented the narrow read-only clarification in both plan artifacts: after health identity, existing OpenReadOnly+authority resolver+EffectSpend+bounded ScanUnreadableWorlds are permitted. All public writes remain authenticated REST. Scan.Scanned tests any World row, not only corruption; no generic all-object emptiness is claimed.

QUICKSTART gives explicit disposable worktree and combined descriptor preparation, separate pins, attended publish then typed-y session mint, sole-writer serve and incremental seed invocation. Attended publish/mint/release were not executed or impersonated. The help/preparation gate runs actual help forms and checks required runbook preparation. The real S1 happy path is executed in test-only authority setup before the daemon starts; no full attended recipe or release capture is claimed.

## Test-first and corrections

`test-first.log` is a compiling named assertion red: the valid-session real-daemon test could not reach the expected incomplete boundary because demo was unknown. An earlier test setup incorrectly used an unbounded Mint context and produced `store: context has no deadline`; it was corrected to the existing bounded test context before recording the real CLI assertion red. Initial implementation then exposed the epoch registry nomination's actual string (`AILANG v0.41.0`, not `0.41.0`); evidence in first-implementation.log preceded correction.

Final-focused r1 exposed a test-instrument bug in the new unknown-subcommand arm: usage text truthfully includes 'S2-S6 incomplete', which the generic negative check mistook for executed seeding. Its complete failure is preserved in preflight-final-r1-failed.log/focused-gates-r1.json. The test now distinguishes the executed boundary phrase, retaining zero real-log-change and zero seed-request assertions. The unknown-subcommand compiling inventory mutation was re-killed afterward (post-correction-unknown-mutant.json).

## Named boundary gates

All gate subprocesses preserve pins and remove only the four first-party established tracing keys: GOOGLE_CLOUD_PROJECT, OTEL_EXPORTER_OTLP_ENDPOINT, OTEL_RESOURCE_ATTRIBUTES, OTLP_GOOGLE_CLOUD_PROJECT. Expected released pins measured in versions.json: AILANG v0.41.0/24ee108 and tool v0.52.1/c68ded4. No AILANG source was authored or changed.

| Gate | rc | Wall seconds | Result |
|---|---:|---:|---|
| `TestDemoSeedAuthenticatedCommits` | 0 | 1.7 | named top-level PASS; no SKIP |
| `TestDemoSeedPreflightRefusals` | 0 | 28.408 | named top-level PASS; no SKIP |
| `TestDemoSeedHelpPreparation` | 0 | 0.627 | named top-level PASS; no SKIP |
| `TestSessionMint_RequiresControllingTerminal` | 0 | 0.621 | named top-level PASS; no SKIP |
| `TestRESTCommitRefusesDecisionPacketObjects` | 0 | 0.671 | named top-level PASS; no SKIP |
| `TestDemoSeedCommitFailures` | 0 | 3.317 | named top-level PASS; no SKIP |

The required five commands plus one targeted commit-failure command are literal in focused-gates.json. The last six-command boundary asserted SHA256 equality for all four Go source files before/after every command, with no mutators active. An earlier overlapping loaded-run/confirmation-mutant boundary was conservatively repeated and is not used as the final evidence. The last additional read-only-open mutation ran only after final gates were terminal and restored byte-exact source, then its named restored subcase passed.

Preflight has 32 refusal cases, including real missing world.apply while other grants remain valid, expired grant/session, insufficient/spent budget, wrong episode, wrong db/interpreter/version, real interloper CAS change, existing manifest, selected-world/log and unselected-world-only cases. Synthetic malformed/absent/down HTTP controls are refusal instrumentation only; the acceptance path uses the unmodified real daemon. Commit-outcome cases prove refusal/no retry for denied response, mismatched success response and a real committed request whose connection closes before the answer.

## Mutations and non-vacuity

mutation-ledger.json: **29/29 compiling assertion kills**, no build-red, panic, SKIP, zero-test or timeout counted. Each temporary edit restored exact source bytes. Raw logs identify actual named assertions.

Required M2a-bearer: removing the production Authorization header makes the successful-stage test red with SessionAbsent. Required M2a-nonempty: fixture creates a legal unselected World before the daemon starts, with empty selected head/log, valid grants and fresh output; removing the single composite emptiness guard produces five actual log rows and the public-inventory equality assertion reds. The older selected-head/log arm alone was masked by CAS; no inventory kill is claimed from that reason-only arm. Separate world-scan suppression, missing grants, spend debit, version/hash, episode, CAS, epoch integrity/decode/nomination, missing source, manifest, flags, read-only-open, HTTP error/absence, commit status/response/uncertainty and frozen-header/carrying support drills are recorded.

refusal-and-support-inventory.json enumerates deliberate guards and exact drills. Not every low-level OS/SQLite operational error was independently forced: disappearing-CWD filepath.Abs, SQLite resolver/EffectSpend/world-scan query failures and every inherited worldReader call-site transport failure remain unprovoked. Primitive-only fixed commit JSON cannot represent a Marshal failure; resolver's nil/nil defensive state violates its production contract. These are named limits, not invented all-error-path coverage. No all-platform behavior is claimed; gates ran on this darwin/arm64 rig.

## Protected observations and remaining ownership

All **12 protected paths** have identical before/after SHA256 and newline counts, including recursive tools/launchd/testdata. Protected git diff rc0; git diff --check rc0; gofmt -l empty. final-observations.json and protected-before/after.json contain measurements and source hashes. No protected file was restored or written.

No executor full-suite/race claim: controller owns pristine/candidate repository gates. No independent evaluation has run inside this role, and no milestone status or authority policy was changed by the executor. Source is stable for controller commit and fresh-context independent judgment. All evidence and source snapshots remain outside the checkout here.

Protected-before to report elapsed: 844.8 seconds (initial skill/design reading precedes this timestamp); implementation stayed within the30-minute role bound.
