# Iteration 314 — executor report (haiku, Agent tool)

Banked by the controller from the executor's returned report (the executor declined to write
report files under its own instructions). No git writes were made by the executor; the controller
rebuilt one commit per milestone from `.snap/M1..M3`, ran the compile fence and the probe tests at
every boundary (rc 0 each), and proved the reconstruction faithful with `shasum -a 256 -c` against
the executor's final tree (4/4 OK).

## Base readings (pristine `59a3e9d6`; Go tree = `2769cff9`)

G1 gofmt empty · G2 vet rc 0 · G3 compile fence rc 0 · G4 `go test ./host/broker/` rc 0 (167.1 s) ·
G5 daemon `TestExecStartup|TestExecSrt` PASS 2 SKIP 2 · G6 probe tests FAIL 0 · G7 race rc 0 ·
P0 reproduction (scratch, deleted): `chatty` FAIL 3/3 and `silent` FAIL 3/3 with the CI signature,
`control` PASS 3/3. G10 at base not re-run (planner's reading: rc 1, an unrelated verifygate load
timeout).

## Milestones

- **M1** `exec_probe.go` +167/−29, test +31: G1–G3, T1, G4–G7 green.
- **M2** test +219 (250 cumulative): G1–G3, G8 node-less PASS 5 SKIP 0, `-count=10`, G4–G7 green.
- **M3** test +71 (321), `ci.yml` +1, `QUICKSTART.md` +5: T4/T6 PASS; G6 PASS 18 SKIP 1 (Arm8);
  G9 20/20 for all 7 names; G10 `verify_go.sh` plain PASSED (304.5 s); G11 actionlint clean;
  G12 local step-8 copy 7/7 PASS; G13 four scope files only.

## Mutation table (executor-reported; 21/21 killed, every restore sha-verified)

M1 FOREIGN-IS-REACH · M2 NO-RETRY · M3 ANY-BYTES · M4 PREFIX-NONCE · M5 EXHAUST-PASSES ·
M6 UNBOUNDED · M7 EXHAUST-AS-REACH · M8 EADDRINUSE-PASSES · M9 NEVER-REAL · M10 SAME-PORT ·
M11 FOREIGN-DOMINATES · M12 NOTOKEN-RETRIED · M13 ALWAYS-RETRY · M14 SEAM-ONLY · M15 HELLO-DROPPED ·
M16 SILENT-RETRIED · M17 SILENT-NOT-STICKY · M18 EOF-IS-SILENT · M19 RESET-IS-FOREIGN ·
M20 SILENT-AS-REACH-TEXT · M21 SILENT-DOMINATES-CHILD — each with its named killer in the red set.

## Deviations

1. M18/M19 as written fail to compile (`declared and not used: werr`); `_ = werr` added, meaning kept.
2. G10 base not re-measured (planner's reading used).
3. Report file not written by the executor (banked here by the controller).
4. `exec_probe.go` +167/−29 vs the plan's ~+125/−30 (plan's own D8 sizing note).
5. Plan header says `e4dbac36`; worktree HEAD was `59a3e9d6` (Go tree identical).
