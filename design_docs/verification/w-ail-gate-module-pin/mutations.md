# Module-manifest gate mutation transcript

Run 2026-08-12 with `AILANG_BIN` set to the released v0.30.0 binary and
`GOTOOLCHAIN=go1.25.6`. Every source mutant was backed up with `cp`, proved landed by a changed
SHA-256, vetted when it touched Go, restored with `cp`, checked byte-identical, and followed by a
green rerun of its scoped arm. The inverse source arm was
`go test ./host/verifygate/ -skip 'TestModuleManifest' -count=1`.

| ID | target and SHA-256 movement | scoped result | inverse | restored arm |
|---|---|---|---|---|
| M1 MUT-STRAY | isolated probe: ABSENT → `5cb14f…b555` | `TestModuleManifestRejectsStrayModule` PASS; its internal gate rc=1 and named `+world/_stray_manifest_probe.ail` | rc=0 | PASS |
| M2 MUT-DEL-LEAF | isolated leaf: `adf076…5079` → ABSENT | `TestModuleManifestRejectsDeletedModule` PASS; its internal gate rc=1 and named `-design_docs/sketches/storejournal.ail` | rc=0 | PASS |
| M3 MUT-SWAP | both M1 and M2 movements in one isolated root | direct copied gate rc=1; one diff named both paths while the real count stayed 11 | rc=0 | pristine copied controls in M1/M2 PASS |
| M4 MUT-NEUTER-CMP | script `a23fd0…7f8` → `44f821…8b31` | rc=1: stray and delete tests FAIL (`want rc=1, got 127`) because the mutant passed membership and reached missing isolated Leg 3 | rc=0 | rc=0 |
| M5 MUT-EMPTY-ALLOWLIST | script `a23fd0…7f8` → `c4f686…b0f3` | rc=1: empty-allowlist test FAIL; output was `LEG1_MODULES[@]: unbound variable`, not its named guard | rc=0 | rc=0 |
| M6 MUT-ENUM-EMPTY | script `a23fd0…7f8` → `b730ff…fb0` | rc=1: empty-enumeration test FAIL; output was `mods[@]: unbound variable`, not its named guard | rc=0 | rc=0 |
| M7 MUT-DIAG-SILENT | script `a23fd0…7f8` → `7208a0…b7a` | rc=1: stray test FAIL at `stray refusal omits "+world/_stray_manifest_probe.ail"` | rc=0 | rc=0 |
| M8 MUT-ISO-INCOMPLETE | Go `6a8eca…b38a` → `94a90e…2d45` | vet rc=0; stray test rc=1 at `pristine isolated control missing`, diff named the omitted leaf | rc=0 | rc=0 |
| M9 MUT-ARM-RC-ONLY + M7 | script → `7208a0…b7a`; Go → `b4ebb7…a4efb` | vet rc=0; composed stray arm rc=0 (GREEN), proving the path/label assertions kill M7 | rc=0; M9 alone rc=0 | rc=0 |
| M10 MUT-ARM-CONTROL-DEAD + M8 | Go `6a8eca…b38a` → `fb9357…e6384` | vet rc=0; composed stray arm rc=0 (GREEN), proving the pristine control observes copy completeness | rc=0; M10 alone (`783efe…8c2f1`) rc=0 | rc=0 |

## Branch inventory from the landed diff

The shell refusal branches are B1 empty actual enumeration, B2 empty expected allowlist, B3 set
inequality, and B4 the path-naming diagnostic inside B3. No additional shell refusal branch was
found. The cut instrument counted five added `exit 1` lines because moving the existing timeout
and JSON-parse exits into the indexed consume loop makes them appear as additions; reading the
diff shows those two are pre-existing behavior, not new membership refusals. The untracked Go
file instrument counted 30 fatal/error calls (positive control 87, negative control 0); these are
harness assertions rather than new production-gate refusal branches.

## Checkable execution deviations

- M8 must also adjust the mandated copy-count expectation from 13/11 to 12/10. Without that
  adjustment, `newIsolatedGateRoot` fails first with `isolated copy landed 12 files / 10 .ail
  files`, so the plan's claimed `pristine isolated control missing` observable is unreachable.
- M10 must neutralize both the pristine-control call and its immediately dependent 11-line count
  assertion. Removing only the call either does not compile (`control` undefined) or fails on the
  empty replacement string before the intended composed experiment.
- Two malformed M9 attempts were discarded before verdict because `go vet` returned rc=1. The
  recorded M9 is the corrected, vet-green run; no parse failure was credited as a kill.
