# Sprint plan — w-interface-hash-covers-the-interface (iteration 188, queue row 27)

**Design doc**: `design_docs/planned/w-interface-hash-covers-the-interface.md`. Quorum r1 and r2 closed it under the narrow-refinement carve-out (§13); T19 and T20 were CLAIMS when it closed.
**Base commit**: `903c17a`. It differs from the designer's prototype base `1a1d929` only by the design doc (`git diff --stat 1a1d929 903c17a` shows 1 file, +378).
**Planner**: planner-role, iter-188. **Verify profile**: go-host. No `.ail` change and no Z3 change.

**PROTOTYPE FIRST.** The whole design, including the r2 additions T19 and T20 and the M3 prose, was built in a disposable copy before this plan was written. Every LOC figure, gate reading and KILLED/SURVIVED verdict below was **measured on that copy**.

The copy is `~/.ailang/state/iter188-planner/proto`. It is a git repo with three commits:
- `base`: `git archive 903c17a`;
- `r1-designer-proto`: the designer's r1 prototype rsynced over the base;
- `planner: T19 T20 generator M3`: this plan's changes.

Evidence:
- `~/.ailang/state/iter188-planner/planner-prototype.diff` (base → final, 1829 lines) and `SHA256SUMS` (31 files);
- the runner `mutate.py`, with per-mutation logs in `mut/`;
- `mut_all.out`;
- the logs `base_test.log`, `final_test.log`, `final_race.log`, `final_gate.log`, `final_ail.log`.

The designer's copy `~/.ailang/state/iter188-designer/proto` was only copied, never modified.

**Executor rule.** You may copy files from the planner prototype. `git -C ~/.ailang/state/iter188-planner/proto diff 411e16e HEAD -- <file>` shows each file's change. You must then **re-run every mutation in §4 against the real worktree**. The prototype's verdicts are evidence, not a substitute.

---

## 1. Baselines and premises (measured this run)

`B=$HOME/.pinned-ailang/ailang` reports `AILANG v0.41.0 / Commit: 24ee108`. The toolchain is `go1.26.6 darwin/arm64`. PATH `ailang` was never used.

| # | Premise | Command | Observed |
|---|---|---|---|
| P1 | Base suite is green | `cd ~/.ailang/state/iter188-planner/base && AILANG_BIN=$B go vet ./... && AILANG_BIN=$B go test ./... -count=1 -timeout 15m` | vet rc=0; **21 `ok`**, 0 FAIL, test rc=0 |
| P2 | `--no-run` output: schema, gates, size, exit | `cp -R packages/world-core $T/p; cd $T/p; $B pkg quality --json --no-run .` twice, plus once in run mode | `--no-run`: rc=2, **3267 B**, schema `ailang.package-quality/v1`, gates `[{"code":"PUB015","level":"gate","msg":"_smoke.ail failed"}]`, v2 `b25fe031…`, 53 signatures. Run mode: rc=0, 3285–3286 B, gates `[]`, same v2. **The output is not byte-deterministic:** two identical `--no-run` runs differ only at `.contracts.wall_seconds`. Fixtures are therefore "as produced", never byte-reproducible. Parsers ignore the field. |
| P3 | Generated fixtures (§7 discipline) | `AILANG_BIN=$B host/pkgproj/testdata/gen_quality_fixtures.sh` (prototype) | pristine rc=2 3267 B, gates `[PUB015]`; pristine_run rc=0, `[]`; probe_ctor rc=2, v2 `0a46f29a…`, 54; v2_unbuilt rc=2, gates `[PUB000,PUB015]`, **no** `hash_v2`, 0 signatures. It embeds only the `mktemp` path (`/var/folders/…`: 2 hits, `$HOME`: 0). Derived fixtures come from five named jq transforms. |
| P4 | The registry still serves the captured bytes | `curl -sS -o … https://storage.googleapis.com/ailang-registry/packages/{world/core/0.1.0,sunholo/auth/0.4.1}/metadata.json; cmp` against the fixtures | `200 8143` / `200 1289`. **Both byte-identical** to the checked-in fixtures. world/core serves `interface_hash_v2 = sha256:ifacev2:b25fe031…e8d8`, schema v2. auth: `has_v2:false`, schema v1. |
| P5 | Live reconcile (AC5), prototype vs base, read-only | `go run ./cmd/world-publish reconcile --store <scratch>.db --registry-origin https://storage.googleapis.com/ailang-registry --probe` | prototype: `state=succeeded-reconciled … detail="served metadata matches all four expected digests"`. Base (control): `… all three expected digests`. |
| P6 | AC8 grep plus control | `grep -c 'interface hash changes because a public ADT changes' <proto doc>`; `git show 1a1d929:design_docs/implemented/w-validated-proven-evidence-boundary.md \| grep -c …` | `0` / `1` |
| P7 | gofmt baseline | `gofmt -l host cmd` in base and in the prototype | Both list the same 2 pre-existing files (`host/store/schema_version_test.go`, `host/store/store.go`). The prototype adds none. |
| P8 | M1 is independently green (bisectable) | base + M1 files only (`~/.ailang/state/iter188-planner/m1only`): `go vet ./...`; `go test ./host/pkgproj/ ./host/broker/ -run 'InterfaceV2\|QueryInterface\|ParseQuality\|EverySubprocessSite'` | vet rc=0; both `ok` |
| P9 | `go build` is not a compile fence for `_test.go` | memory rule, applied | every gate below uses `go vet ./...` |

## 2. Milestones, files, LOC (measured on the prototype, `git diff --numstat` base → final)

**Commit plan: THREE commits, one per milestone.** Each builds and passes on its own: P8 measured M1 alone, and M1+M2+M3 is the full prototype. M2 depends on M1, because `QueryInterface` is used by the step-7 helper and the wiring test.

### M1 — `pkgproj` v2 query: bounded, gate-checked, parsed (≈0.5 d)

**Files**

| File | LOC | Kind |
|---|---|---|
| `host/pkgproj/iface.go` (new) | 211 | production |
| `host/pkgproj/iface_test.go` (new) | 344 | test (15 tests) |
| `host/pkgproj/testdata/gen_quality_fixtures.sh` (new, `chmod +x`) | 54 | generator |
| `host/pkgproj/testdata/quality_world_core_{pristine,pristine_run,probe_ctor,v2_unbuilt}.json` | 84+79+84+83 | **generated**, checked in as produced |
| `host/pkgproj/testdata/derived_quality_{v2_holds_v1,v1_absent,schema_v2,extra_gate,pub015_other_msg}.json` | 1 each | **derived** by the generator's named jq transforms |
| `host/broker/registry_publish_test.go` | +7 | AC10(a) driver `drivePkgprojQueryInterface` |

**Surface** (as prototyped):

```go
type InterfaceIdentity struct{ V1, V2 string; Signatures int }
func ParseQualityInterface(out []byte) (InterfaceIdentity, error)            // pure; refuses, never defaults
func QueryInterface(ctx context.Context, packageDir string, m Manifest, ailangBin string) (InterfaceIdentity, error)
func queryInterface(ctx, dir, m, bin string, b queryBounds) (InterfaceIdentity, error) // tests pass small bounds
func checkQualityGates(out []byte, exitCode int) error                        // THE one exit-code/gate rule (T19)
type QueryTimeoutError struct{ Timeout time.Duration }
var ErrQualityOutputOverflow                                                  // "stdout or stderr exceeded its output cap"
const queryInterfaceTimeout = 30 * time.Second; queryInterfaceWaitDelay = 2 * time.Second; maxQualityOutputBytes = 1 << 20
const spuriousNoRunGateCode, spuriousNoRunGateMsg = "PUB015", "_smoke.ail failed"
type cappedBuffer struct{ buf bytes.Buffer; limit int; over bool; onOverflow func() } // onOverflow = runCtx's cancel (T20)
```

**Ordered tasks**
1. Write `gen_quality_fixtures.sh` first and run it. Commit its output unmodified. Hand-editing a fixture is a defect. See DV1 and DV2.
2. Write `ParseQualityInterface`, then T5 and T6 against the fixtures.
3. Write `checkQualityGates`. Its rule, in order:
   1. exit ∉ {0,2} → refuse `exit N`;
   2. parse `gates` (missing array → refuse);
   3. exit 0 requires `len(gates)==0`;
   4. exit 2 requires exactly one gate with code `PUB015` **and** msg `_smoke.ail failed`.

   The refusal names the codes, e.g. `exit 2 with gates [PUB000 PUB015]`. It runs **before** `ParseQualityInterface`, so the broken export is refused by the gate rule and not by the v2-absent guard.
4. Write `queryInterface`, with the D1a result order:
   1. overflow (`stdout.over || stderr.over`);
   2. caller `ctx.Err()` → wrapped `cancelled by caller`;
   3. `runCtx` `DeadlineExceeded` → `*QueryTimeoutError`;
   4. a non-`ExitError` run error;
   5. `checkQualityGates`;
   6. parse;
   7. v1 cross-check against `InterfaceHash(manifest)`.

   Both caps use `onOverflow: cancel`.
5. Write the tests T1–T6 and T14–T20, then the AC10(a) driver.

**Named tests** (all PASS on the prototype, `-v` timings in parentheses)

| T | Test | Asserts |
|---|---|---|
| T1 | `TestInterfaceV2MovesWhenAnExportedADTGainsAConstructor` (0.50 s) | real binary; `ProbeArm(HashRef)` → v2 moves, v1 does not, signatures +1 |
| T2 | `TestInterfaceV2IgnoresACommentOnlyEdit` (0.46 s) | identity unchanged |
| T3 | `TestQueryInterfaceRefusesAV1Disagreement` | `Edition:"2"` → "disagrees with pkgproj.InterfaceHash" |
| T4 | `TestQueryInterfaceRefusesAnInfraExitEvenWithValidJSON` | fake prints the real fixture and exits 1 → refused |
| T5 | `TestParseQualityInterfaceFixtures` | pristine `b25fe031…`/53; ctor ≠ pristine and v1 equal; v2_unbuilt → "interface identity v2 absent" |
| T6 | `TestParseQualityInterfaceDerivedNegatives` | each derived fixture reds its own guard's text |
| T14 | `TestQueryInterfaceTimesOutOnAHangingBinary` (3.01 s) | `exec sleep 30` → `*QueryTimeoutError{3s}` |
| T15 | `TestQueryInterfaceReturnsWhileADescendantHoldsStdout` (4.04 s) | returns the named timeout at timeout + WaitDelay. It logs `descendant … alive after return: true` (the declared row-24 residual) and reaps the pid |
| T16 | `TestQueryInterfaceCapsStdout` (0.17 s) | `head -c 100000 /dev/zero`, cap 4096, **60 s internal timeout** → `ErrQualityOutputOverflow` |
| T17 | `TestQueryInterfaceHonoursCallerCancellation` (1.51 s) | caller 1.5 s, internal 6 s → the caller's `DeadlineExceeded`, never `*QueryTimeoutError` |
| T18 | `TestQueryInterfaceUsesFiniteProductionBounds` | 0 < timeout < 120 s (the literal `run_bounded 120`), WaitDelay > 0, cap > 0 |
| **T19** | `TestQueryInterfaceRefusesExit2WithANonSpuriousGate` (0.95 s) | **Controls first:** real `--no-run` pristine at exit 2 is ACCEPTED (53); real run-mode pristine at exit 0 is ACCEPTED. **Refusals:** the generated v2_unbuilt at exit 2 → "exit 2 with gates [PUB000 PUB015]"; derived extra_gate (keeps `hash_v2`) at exit 2 → "…[PUB015 PUB000]"; `--no-run` pristine at **exit 0** → "exit 0 with gates [PUB015]"; derived pub015_other_msg at exit 2 → refused |
| **T20** | `TestQueryInterfaceOverflowCancelsTheChild` (0.34 s) | For **both** stdout and stderr (stderr's fixed 64 KiB cap): the double runs `head -c 100000 /dev/zero[>&2]; exec sleep 30`. Internal timeout **60 s**, return guard **10 s**. Result: `ErrQualityOutputOverflow` at 0.11 s (stdout) and 0.23 s (stderr) |

The shared `guard = 10 * time.Second` is the `runWithin` limit for T14–T16 and T20, so a mutant **reds instead of hanging**. T17 uses 5 s.

**Acceptance (M1)**, with `<WT>` = the executor's absolute worktree path, e.g. `/Users/voightkampff/dev/sunholo-data/ailang-world`:

```
cd <WT> && AILANG_BIN=$HOME/.pinned-ailang/ailang host/pkgproj/testdata/gen_quality_fixtures.sh   # prints "gen_quality_fixtures: OK"
cd <WT> && go vet ./...
cd <WT> && AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./host/pkgproj/ ./host/broker/ -run 'InterfaceV2|QueryInterface|ParseQuality|EverySubprocessSite' -count=1 -timeout 300s -v
cd <WT> && AILANG_BIN=$HOME/.pinned-ailang/ailang go test -race ./host/pkgproj/ -count=1 -timeout 300s
pgrep -fl 'sleep 30' || echo none      # must print "none" once the run has finished (T15 reaps its descendant)
```

Then run the M1 mutations in §4. Each must be KILLED by its named test.

### M2 — packet, gate, fence, reconcile (≈0.3 d)

**Files**

| File | LOC | Change |
|---|---|---|
| `host/pkgproj/readypacket.go` | +6/−2 | field `InterfaceHashV2` (`interfaceHashV2`, sorted between `interfaceHash` and `package`), `ReadyPacketFields`, `Field()`, `RecomputeReadyPacket(dir, m, compilerVersion, interfaceHashV2)` |
| `host/pkgproj/readypacket_test.go` | +6/−5 | T7: Equal-table row |
| `host/broker/registry_reconcile.go` | +21/−3 | D4: `ExpectedInterfaceV2`, `RegistryMetadata.InterfaceHashV2 *string`, R6 "four", P4 (nil → conflict, mismatch → conflict, match → "all four") |
| `host/broker/registry_reconcile_test.go` | +15/−9 | `reconcileExpectedV2`, fixtures, R6 text |
| `host/broker/registry_reconcile_v2_test.go` (new) | 91 | T8–T11 (the auth-fixture read is now `t.Fatal` on error, DV9) |
| `host/broker/testdata/metadata_world_core_0.1.0.json`, `metadata_sunholo_auth_0.4.1.json`, `CAPTURED.txt` (new) | 169+33+3 | captured as served (P4); the capture commands are in `CAPTURED.txt` (DV3) |
| `cmd/world-publish/fences.go` | +5/−1 | `frozenInterfaceHashV2`; `RecomputeReadyPacket(…, frozenInterfaceHashV2)` |
| `cmd/world-publish/main.go` | +1 | `ExpectedInterfaceV2: packet.InterfaceHashV2` |
| `cmd/world-publish/wiring_test.go` | +19/−1 | T12; the golden test recomputes v2 **live** via `QueryInterface` (`t.Fatal` if `AILANG_BIN` is unset) |
| `host/runbook/runbook_stageb_test.go` | +21 | T13 (`ifaceV2Digest` regex, disjoint from `fullDigest`) |
| `scripts/verify_world_package.sh` | +6/−4 | step-7 helper calls `QueryInterface` and emits `interfaceHashV2=`; the key loop and step 9's python require and write it |
| `scripts/world_package_ready_packet.golden.json` | +1/−1 | **regenerated by the gate's red diff** (S8 recipe), never hand-typed. The new line is in §8. |
| `docs/SELF_MOD_PUBLISH.md` | +2/−1 of the file's +8/−5 | the relabelled `interfaceHash` row and the new `interfaceHashV2` row (T13 binds them) |

Production (Go + shell) is about 39 LOC plus the golden. Tests are about 152.

**Ordered tasks**
1. Change `readypacket.go`. `go vet ./...` then reds every `RecomputeReadyPacket` caller; fix `fences.go` and `wiring_test.go`.
2. Update the gate script.
3. Run the gate. Step 9 reds with a diff; copy its new line into the golden and re-run until green.
4. Make the reconcile changes and add the captured testdata.
5. Add `main.go` wiring, T12, the runbook row and T13.

**Acceptance (M2)**

```
cd <WT> && AILANG_BIN=$HOME/.pinned-ailang/ailang WORLD_PKG_AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_world_package.sh   # AC3: "9/9", rc=0
cd <WT> && go run ./cmd/world-publish reconcile --store "$(mktemp -d)/r.db" --registry-origin https://storage.googleapis.com/ailang-registry --probe   # AC5: "… all four expected digests"
cd <WT> && git diff --stat 903c17a -- host/broker/registry_publish.go host/broker/approve.go   # AC7: empty output
cd <WT> && go vet ./... && AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./... -count=1 -timeout 15m   # 21 ok / 0 FAIL
```

AC4 is `MUT-GATE-PY-OMITS-V2` in §4, and it must show gate rc=1. Then run the M2 mutations.

### M3 — honesty prose (≈0.15 d)

**Files** (measured):
- `host/pkgproj/pkgproj.go`, +6/−3 (the `Package` comment: v1 covers the manifest only; v2 comes from `QueryInterface`);
- `scripts/verify_ail.sh`, +2/−1 (the comment at line 49);
- `design_docs/implemented/w-validated-proven-evidence-boundary.md`, +8/−1 (§8.3's new first sentence, verbatim from D5, plus a dated correction blockquote);
- `docs/SELF_MOD_PUBLISH.md`, the remaining prose:
  - §4's field list gains `interfaceHashV2`;
  - "These three digests" → "four";
  - the reconcile state table "all three expected digests" → "all four";
  - **the approval-scope sentence stays three and names them** (DV7).

Do **not** edit `coding-standards.md` (§9).

**Acceptance (M3)**

```
cd <WT> && grep -c 'interface hash changes because a public ADT changes' design_docs/implemented/w-validated-proven-evidence-boundary.md   # 0
cd <WT> && git show 1a1d929:design_docs/implemented/w-validated-proven-evidence-boundary.md | grep -c 'interface hash changes because a public ADT changes'   # 1 (control)
cd <WT> && AILANG_BIN=$HOME/.pinned-ailang/ailang ./scripts/verify_ail.sh   # "11 required identities verified, 40 named tests pass"
cd <WT> && AILANG_BIN=$HOME/.pinned-ailang/ailang go test ./host/runbook/ -count=1   # binds the runbook to the golden after the prose edits
```

The controller files §11 O1 (upstream PUB015 under `--no-run`) and O5 (the `CrossCheck` row candidate). The executor does not.

## 3. Full gate list, derived from `.github/workflows/ci.yml`

There are two jobs. Both install the pinned **v0.41.0**:
- `ailang-verify` exports it as `WORLD_PKG_AILANG_BIN` (PATH `ailang` there is `releases/latest`);
- `go-verify` exports it as `AILANG_BIN`.

Both install z3 4.16.0.

| CI job | CI step | Local command (pinned binary) | Prototype reading |
|---|---|---|---|
| ailang-verify | `./scripts/verify_ail.sh` | `AILANG_BIN=$B WORLD_PKG_AILANG_BIN=$B ./scripts/verify_ail.sh` | rc=0, 11 identities / 40 tests (it runs the world-package leg, 9/9) |
| go-verify | `./scripts/verify_go.sh` (25 min) | Do **not** run the script locally (memory rule). Run its legs directly: `go vet ./...` (compile fence, P9); `go build ./...`; `go test -json ./host/evidence -count=1` (37-test manifest, untouched here); `go test ./... -count=1`; `go test ./... -count=1 -race -timeout 8m` | vet rc=0; plain run **21 ok / 0 FAIL, 54 s** (and 2 more green runs, 57 s); **-race 21 ok / 0 FAIL, 91 s** |
| go-verify | reproducer `design_docs/verification/w-race-gate-blindspot/run.sh` | CI only (platform-gated) | not run: untouched surface |
| go-verify | `./scripts/bench_worldd.sh --smoke` / `--check-claims` | same | rc=0 / rc=0 |
| go-verify | recorder refuses off-rig | CI only (needs a linux runner to refuse at `hw.ncpu`) | not run |
| go-verify | `bash scripts/test_gate1_range_check.sh` | same | 18 passed |
| go-verify | `bash scripts/test_queue_census.sh` | same | 68 passed |
| go-verify | `bash scripts/queue_census.sh --doc design_docs/world-mission.md --control-closed 1 --control-open 79` | same | rc=0 |
| go-verify | `bash scripts/test_gate0_self_notices.sh` | same | 88 passed |
| go-verify | `bash scripts/check_no_personal_email.sh` | same | rc=0 |
| go-verify | `bash scripts/test_check_no_personal_email.sh` | same | 10 passed |

Plus these checks, which CI does not run:
- `gofmt -l host cmd`: only the 2 pre-existing files;
- `./scripts/verify_world_package.sh` with both env vars: 9/9;
- AC5 live;
- AC7 empty;
- AC8 `0`/`1`.

## 4. Mutation matrix, run against the prototype

**Method** (`~/.ailang/state/iter188-planner/mutate.py`):
- apply an exact-string anchor, and assert its occurrence count;
- run the named test(s) with the pinned `AILANG_BIN`;
- restore the file and verify the restore by sha256;
- a build failure is reported as COMPILE-FAIL and **never counted as a kill**.

**Tally: 33 mutations run → 32 KILLED, 1 SURVIVED (an equivalent mutant; see the last row).**
- The design's full matrix is 24 r1 mutations plus `MUT-EXIT2-ANY-GATE` and `MUT-OVERFLOW-NO-CANCEL`. **All 26 are KILLED.**
- The two r2 CLAIMS (T19/T20) are now measured.
- 6 extra mutations added by this plan are KILLED.

All M1 rows edit `host/pkgproj/iface.go` unless another file is named.

| ID | M | Anchor → mutant | Must red | Result |
|---|---|---|---|---|
| MUT-V2-IS-V1 | 1 | `V2: *in.HashV2,` → `*in.HashV1` | T1, T5 | KILLED (T1+T5) |
| MUT-V2-IS-SOURCE | 1 | before the final `return id, nil`: `id.V2 = ContentHash(packageDir)` | T2 | KILLED |
| MUT-ABSENT-V2-OK | 1 | the v2 guard: nil → `new(string)` (absent = "" success) | T5 | KILLED |
| MUT-V2-SHAPE-UNCHECKED | 1 | `if in.HashV2 == nil \|\| !v2Shape…` → `if in.HashV2 == nil {` | T6 | KILLED |
| MUT-V1-GUARD-OFF | 1 | the v1 guard: nil → `new(string)`, shape check `false` | T6 | KILLED |
| MUT-SCHEMA-UNCHECKED | 1 | `if doc.Schema != "ailang.package-quality/v1" {` → `if false {` | T6 | KILLED |
| MUT-NO-V1-CROSSCHECK | 1 | `id.V1 != local {` → `false && …` | T3 | KILLED |
| MUT-RC1-ACCEPTED | 1 | `if exitCode != 0 && exitCode != 2 {` → `… && exitCode != 1 {` | T4 | KILLED |
| MUT-RC2-REFUSED | 1 | same anchor → `if exitCode != 0 {` | T1 | KILLED |
| **MUT-EXIT2-ANY-GATE** (r2) | 1 | `if len(gates) != 1 \|\| gates[0].Code != … \|\| gates[0].Msg != … {` → `if false {` | T19 | **KILLED** |
| MUT-EXIT2-CODE-ONLY (+) | 1 | drop ` \|\| gates[0].Msg != spuriousNoRunGateMsg` | T19 (pub015_other_msg leg) | KILLED |
| MUT-EXIT0-GATES-UNCHECKED (+) | 1 | `\t\tif len(gates) != 0 {` → `if false {` | T19 (exit-0 leg) | KILLED |
| MUT-GATE-RULE-REFUSES-ALL (+) | 1 | the exit-0 `return nil` → `return errors.New("refuse")` | T19 (positive control) | KILLED |
| MUT-NO-DEADLINE | 1 | `context.WithTimeout(ctx, b.timeout)` → `context.WithCancel(ctx)` | T14, T15 | KILLED (10 s guard, 20.5 s) |
| MUT-TIMEOUT-UNNAMED | 1 | delete the `QueryTimeoutError` branch | T14 | KILLED |
| MUT-NO-WAITDELAY | 1 | `cmd.WaitDelay = b.waitDelay` → `= 0` | T15 | KILLED (guard, 10.5 s) |
| MUT-NO-CAP | 1 | `if c.buf.Len()+len(p) > c.limit {` → `if false {` | T16 | KILLED |
| MUT-PROD-TIMEOUT-HUGE | 1 | `30 * time.Second` → `3600 * time.Second` | T18 | KILLED |
| MUT-TIMEOUT-DROPS-CALLER-CTX | 1 | `WithTimeout(ctx, …)` → `WithTimeout(context.Background(), …)` | T17 | KILLED |
| MUT-CALLER-CHECK-OFF | 1 | `if cerr := ctx.Err(); cerr != nil {` → `…; false && cerr != nil {` | T17 | KILLED |
| **MUT-OVERFLOW-NO-CANCEL** (r2) | 1 | `stdout := &cappedBuffer{limit: b.maxOutput, onOverflow: cancel}` → without `onOverflow` | T20 | **KILLED** (10 s guard vs a 60 s internal timeout) |
| MUT-STDERR-OVERFLOW-NO-CANCEL (+) | 1 | the same for `stderr := &cappedBuffer{limit: 64 << 10, onOverflow: cancel}` | T20 (stderr leg) | KILLED |
| MUT-AC10-DRIVER-DROPPED (+) | 1 | `registry_publish_test.go`: delete the `"host/pkgproj/iface.go": drivePkgprojQueryInterface,` line | `TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential` | KILLED |
| MUT-PACKET-FIELD-ALIAS | 2 | `readypacket.go` `case "interfaceHashV2": return p.InterfaceHashV2` → `p.InterfaceHash` | T7 (`TestReadyPacketEqualNamesTheFirstDifferingField/interfaceHashV2`) | KILLED |
| MUT-RECON-NO-V2-COMPARE | 2 | `if *served.InterfaceHashV2 != cfg.ExpectedInterfaceV2 {` → `if false {` | T9 | KILLED |
| MUT-RECON-ABSENT-OK | 2 | the nil branch returns `succeeded-reconciled` | T10 | KILLED |
| MUT-RECON-R6-V2 | 2 | drop ` \|\| cfg.ExpectedInterfaceV2 == ""` | T11 | KILLED |
| MUT-WP-DROP-EXPECTED-V2 | 2 | `main.go`: delete `ExpectedInterfaceV2: packet.InterfaceHashV2,` | T12 | KILLED |
| MUT-FENCE-CONST | 2 | `fences.go` `…a693d621e8d8"` → `…e8d9"` | `TestEveryRefusalBranchStopsWithItsExactLineAndHasAPositiveControl`, `TestTheEntrypointReachesTheProductionRefusal` | KILLED |
| MUT-RUNBOOK-ROW-DROPPED | 2 | delete the `interfaceHashV2` table row | T13 | KILLED |
| MUT-GATE-PY-OMITS-V2 (AC4) | 2 | step 9 python: drop `"interfaceHashV2":values["interfaceHashV2"], ` | package gate | KILLED: rc=1, `✗ ready packet differs byte-for-byte from golden` |
| MUT-GATE-HELPER-EMITS-V1 (+) | 2 | step-7 helper `r.Local.Interface, id.V2, r.Local.Tarball` → `id.V1` | package gate | KILLED: rc=1, same line |
| MUT-GOLDEN-LIVE-V2-CONST (+) | 2 | `wiring_test.go` `frozenCompilerVersion, iface.V2)` → `frozenInterfaceHashV2)` (+ `_ = iface`; the literal form was COMPILE-FAIL and was reshaped) | `TestWorldCoreManifestMatchesTheCommittedGolden` | **SURVIVED: equivalent mutant.** The constant equals the live v2 on today's tree, so no test can separate the two until the exported interface changes. This is not a gap in any criterion the design claims; it is recorded so nobody cites that test as the kill for "live". |

Removed by the design and not re-run: `MUT-SIGS-UNCHECKED` (the guard was deleted, D1).

For the executor's re-run, `mutate.py` is reusable: point its `P` at the worktree. Every anchor above has count 1 at the prototype's final state.

## 5. Deviations from the design doc (each forced by a measurement)

- **DV1 — fixture mode.** The r1 prototype's pristine and ctor fixtures were **run-mode** captures (3286 B, gates `[]`, rc 0). Its v2_unbuilt was `--no-run`. §7's generator line also omits `--no-run`. Production parses `--no-run`, so the plan's generator captures every fixture with `--no-run` (P3), plus **one** run-mode capture (`quality_world_core_pristine_run.json`) as the exit-0 positive control for T19. Size figures: `--no-run` 3267 B (V24 confirmed); run 3285–3286 B, varying with `wall_seconds` (P2).
- **DV2 — the generator did not exist, and the prototype's v2_unbuilt fixture embedded a home path.** `grep -c /Users/` → 2. The error text named `~/.ailang/state/iter188-designer/broken/…`, which violates §7 ("temp paths, never a home directory"). The plan's generator uses `mktemp -d` and refuses if any fixture contains `$HOME`.
- **DV3 — `host/broker/testdata/CAPTURED.txt` was missing** from the prototype (§7 requires it). It has been added with the two `curl` commands (P4).
- **DV4 — the r1 test bounds FAIL under full-suite load.** At 1 s / 0.3 s:
  - T15 failed **3 of 3** full `go test ./...` runs with "descendant pid not recorded";
  - T16 failed 3 of 3 with "exceeded its 1s bound";
  - both passed in isolation. V22's 3/3 green was isolation-only.

  Cause: the fake shell had not produced output or its pid within 1 s while about 20 packages ran in parallel.

  Plan values:
  - `tinyBounds` 3 s / 1 s;
  - T16 and T20 use a 60 s internal timeout, so the cap alone must produce the error;
  - T17 is caller 1.5 s / internal 6 s;
  - the guard is 10 s (5 s for T17).

  Measured: full suite green **3/3** afterwards (2 before the final T19 refactor, 1 after), `-race` green, and every bounded-wait mutation still KILLED. **AC9's "within timeout + 0.1 s" slack was never asserted by the prototype and is not asserted here.** Boundedness is enforced by the guard. A sub-second slack assertion would reintroduce the load flake. *Controller: accept, or require a timing assertion with measured slack.*
- **DV5 — one exit-code rule.** The r1 prototype accepted exit 2 in `queryInterface`. T19 adds a gate rule, so exit-code acceptance would have lived in two places. `MUT-RC1-ACCEPTED` against the outer check would then be masked by the inner one. That is reasoned from the code shape, not measured. The plan puts exit-code acceptance only in `checkQualityGates`, and `MUT-RC1-ACCEPTED` and `MUT-RC2-REFUSED` target that single anchor (both KILLED).
- **DV6 — T19/T20 are broader than the doc's text, to kill the obvious neighbours.** T19 adds:
  - positive controls at exit 2 and exit 0 (killing `MUT-GATE-RULE-REFUSES-ALL`);
  - an exit-0-with-gate leg (`MUT-EXIT0-GATES-UNCHECKED`);
  - a fifth derived fixture, `derived_quality_pub015_other_msg.json` (`.gates[0].msg = "_smoke.ail failed: derived other message"`), because the design names the msg and a code-only check survived nothing without it (`MUT-EXIT2-CODE-ONLY`).

  T20 covers stderr too, matching D1a's "either output cap". `ErrQualityOutputOverflow`'s text is now "stdout or stderr exceeded its output cap".
- **DV7 — a SELF_MOD_PUBLISH.md site that D5 would have made false.** D5 says "These three digests → four". The same page's approval paragraph (line 182, "The scope binds the three digests above") must **not** become four: D6 keeps v2 out of the scope grammar. The plan rewrites it to name `contentHash`, `interfaceHash` and `tarballSHA256` explicitly. D5 also missed the reconcile state table (line 219, "all three expected digests"), which now contradicts the tool's output (P5); the plan changes it to "four".
- **DV8 — stale comments in the prototype's `iface.go`.** "3286 bytes … ~300x" and "0.28-0.65 s … ~50x" are corrected to 3267 B / ~320× (V24) and 0.28–0.85 s / ~35× (V17/V22), matching D1a.
- **DV9 — T10's instrument could pass vacuously.** The prototype read the auth fixture with `ctl, _ := os.ReadFile` and `_ = json.Unmarshal`. A missing file gave a nil map, and "no v2 key" held vacuously. It now calls `t.Fatal` on a read or decode error or an empty map.

## 6. AC → test map

| AC | Discharged by |
|---|---|
| AC1 re-break | T1 + `MUT-V2-IS-V1` |
| AC2 negative control | T2 + `MUT-V2-IS-SOURCE` |
| AC3 gate 9/9 | `verify_world_package.sh` rc=0 (prototype reading) |
| AC4 | `MUT-GATE-PY-OMITS-V2` → gate rc=1 |
| AC5 live | P5 (four vs base three) + T8 offline |
| AC6 | §3 gate list |
| AC7 | `git diff --stat` empty for `registry_publish.go` and `approve.go` (the prototype does not touch them) |
| AC8 | P6 `0`/`1` |
| AC9 bounded waits | T14–T18, T20 + the D1a mutations (all KILLED) |
| r2 exit-2 rule | T19 + `MUT-EXIT2-ANY-GATE` (+3 neighbours) |
| r2 overflow cancels | T20 + `MUT-OVERFLOW-NO-CANCEL` (+ the stderr twin) |

## 7. Residuals, risks, controller decisions

1. **D3 named consequence (controller's call; the doc recommends accept).** A future exported-interface change to `packages/world-core` reds CI at `STOP fence=packet reason=drift` until a version-bump design exists.
2. **DV4 timing slack (controller's call).** Accept guard-only boundedness, or ask for a measured-slack assertion.
3. **Descendant survives** (T15 logs `alive after return: true`). This is declared and owned by row 24.
4. **`CrossCheck` stays unbounded in-process** (O5). File it as a new row.
5. **Upstream O1** (`--no-run` reports "not run" as PUB015 "failed"). The controller files it on `sunholo-data/ailang` with a mission-control note. If upstream fixes it, `--no-run` exits 0 with `[]`, which T19's exit-0 arm already accepts.
6. **Fixtures are not byte-reproducible** (`wall_seconds`). Regeneration produces a diff on that line. This is expected and parsers ignore it.
7. **Registry drift risk** for the captured docs: P4 re-verified them byte-identical this run.

## 8. Prototype reference

- Root: `~/.ailang/state/iter188-planner/proto`. Diff: `git -C … diff 411e16e HEAD`, also saved as `~/.ailang/state/iter188-planner/planner-prototype.diff`.
- Hashes: `~/.ailang/state/iter188-planner/SHA256SUMS`. Leading bytes:

| File | sha256 (prefix) |
|---|---|
| `iface.go` | `6a401e022d99c73b` |
| `iface_test.go` | `7f0ca4d2a085a279` |
| `gen_quality_fixtures.sh` | `6933005afe6058ad` |
| `readypacket.go` | `5fb638cf10769c04` |
| `registry_reconcile.go` | `883cdc724db99a1a` |
| `registry_reconcile_v2_test.go` | `900dd94266208b9d` |

New golden line (produced by the gate, P3/AC3):

```
{"compilerVersion":"AILANG v0.41.0","contentHash":"sha256:0c8c60616e592dc01891e8bbb59350786f242a2f79a9eb2c587ae8b0ca2e00b9","effects":[],"exports":["world/types","world/contracts","world/transitions","world/logepoch"],"interfaceHash":"sha256:d16cc88270ff4c4eaaa583e644d3ea30e2e4b2e36f95fd7108d920046cdb4083","interfaceHashV2":"sha256:ifacev2:b25fe03155db0c7bf595cf730295b945d6ac64ec415a1998fae8a693d621e8d8","package":"world/core","tarballBytes":9785,"tarballSHA256":"sha256:44fc9fab7be710f09b84274f744445db41a79df77c71cc43c0e952d75d97c27f","version":"0.1.0"}
```
