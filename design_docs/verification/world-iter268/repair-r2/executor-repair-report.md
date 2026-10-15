# Executor portability repair — World268 M2a / R2-B1

Terminal source-only local commit **e8f916c51285856f8823b9251876aa49b99b94dc**, branch repair/world-iter268-m2a-r2, worktree `/Users/voightkampff/dev/sunholo-data/.repair-world-iter268-r2`, base7ba691b732d327f6b38a16223152b876e7d197df. **Not pushed.** Only cmd/ailang-worldd/demo_test.go committed; no production/gate/workflow/docs change. Previous source and evaluator trees untouched. Declared actual role model native GPT-6.1 Sol; tokens unknown; no subagents.

Read terminal evaluator-r2.md and parsed its JSON from `.eval-world-iter268-r2/design_docs/verification/world-iter268/`: FAIL70/100, sole blocking R2-B1; R1-B1/B2/B3 closed. The ordinary-CI missing tool-variable finding is repaired with a unique regular byte file under t.TempDir. Its content hash differs from the real archived pin, so the same production hash guard must refuse before version/manifest lookup. Existing exact refusal and zero-inventory-change assertions remain, with no skip or CI provisioning change. Real AILANG_BIN remains required; this M2a fixture no longer requires the later workspace-tool binary.

Independently reproduced original scoped unset-tool failure before edits (before-unset-tool.log/json). Repaired scoped command passes unset; reverting exactly the old test fixture compiles and dies on required-tool-pin assertion; byte-exact restoration recorded in fixture-mutation.json and raw log. This is a test-fixture portability mutation, not relabelled as a new production-law mutation.

Final immutable boundary: **18 commands, all required named top-level PASS, zero SKIP**, covering nine commands / eleven named identities separately with WORLD_TOOL_AILANG_BIN set and unset. Included all original6 M2a gates, TestQuickstartDecisionPacketSection, all3 earlier regression names and the nine-site subprocess audit. Every command checked all five implementation/document SHA256s unchanged. Literal argv, bounds, environment arm and wall times are in portability-gates.json; each log is banked. Test children remove only GOOGLE_CLOUD_PROJECT, OTEL_EXPORTER_OTLP_ENDPOINT, OTEL_RESOURCE_ATTRIBUTES and OTLP_GOOGLE_CLOUD_PROJECT. No environment values or credentials were dumped.

| Tool variable | Required names | rc | Wall seconds |
|---|---|---:|---:|
| set | TestDemoSeedAuthenticatedCommits | 0 | 3.42 |
| set | TestDemoSeedPreflightRefusals | 0 | 51.048 |
| set | TestDemoSeedHelpPreparation | 0 | 1.001 |
| set | TestSessionMint_RequiresControllingTerminal | 0 | 1.078 |
| set | TestRESTCommitRefusesDecisionPacketObjects | 0 | 2.018 |
| set | TestDemoSeedCommitFailures | 0 | 5.985 |
| set | TestQuickstartDecisionPacketSection | 0 | 1.072 |
| set | TestQuickstartSection11FlagsMatchTheCLI,TestWorldBoundaryDependencyAllowlist,TestWorldBoundaryRecordingWriter | 0 | 1.876 |
| set | TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential | 0 | 19.978 |
| unset | TestDemoSeedAuthenticatedCommits | 0 | 2.58 |
| unset | TestDemoSeedPreflightRefusals | 0 | 61.135 |
| unset | TestDemoSeedHelpPreparation | 0 | 1.655 |
| unset | TestSessionMint_RequiresControllingTerminal | 0 | 1.716 |
| unset | TestRESTCommitRefusesDecisionPacketObjects | 0 | 1.73 |
| unset | TestDemoSeedCommitFailures | 0 | 3.629 |
| unset | TestQuickstartDecisionPacketSection | 0 | 2.301 |
| unset | TestQuickstartSection11FlagsMatchTheCLI,TestWorldBoundaryDependencyAllowlist,TestWorldBoundaryRecordingWriter | 0 | 4.122 |
| unset | TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential | 0 | 57.475 |

All **12 protected paths** retain exact before/after SHA256 and newline equality, including recursive launchd testdata. Protected diff, diffcheck and formatting clean; observations.json and protected-before/after.json bank this. Committed source bytes equal the final gate subject; source-commit.json records the full SHA and one-file commit scope.

R2's nonblocking budget90s→24h, StateRoot and entryHash coverage survivors remain separate controller-owned3n residuals; no scope expansion or normative canonical-hash claim. No broad duplicate suite, production mint/publish, attended confirmation, inbox call or push. Full repaired-head gates, report/evidence commits, fresh R3 independent evaluation and remote CI remain controller-owned. This report is an executor handoff, not an independent PASS or release claim.
