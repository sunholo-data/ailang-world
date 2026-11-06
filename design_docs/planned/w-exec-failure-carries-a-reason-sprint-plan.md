# Sprint plan — w-exec-failure-carries-a-reason (iteration 360, row 255)

Approved design: [w-exec-failure-carries-a-reason.md](w-exec-failure-carries-a-reason.md), **r3**, including §3.7, Claim F and §11's controller narrow-refinement carve-out. Planner base: `adc03d4223a2ff89412d49d1499420684f1b2ccb` (detached HEAD). Exactly one milestone, **M1**, status **planned**, approximately 0.5 day. The controller commits and pushes; this planner performs no Git writes to tracked implementation files or the index.

M1 carries one bounded classified exec failure line through the stored record, live/replayed world output, operator log, `why` and `call`. It lazily marks the store schema 6 before writing a lined record. If marking fails, return `nil, hashref.HashRef{}, &OutcomeWriteError{Effect: req.Effect, Scope: req.Scope, Err: fmt.Errorf("broker: mark store for failure lines: %w", markErr)}`; write no effect record, append no outcome, preserve the existing debit and execute the handler exactly once. Keep the mark-before-record and plain-failure-at-5 TV rows. This is Go host-boundary error classification and persistence, which requires the Go errors at their source; it is not an AILANG package/kernel change (§3.6).

## M1 — exact implementation scope and measured bytes

`ABSENT` means the file does not exist at HEAD. All post hashes below were measured on a fresh `git archive HEAD` with deliverable 3 applied, after the manual silent-fallback mutation was restored and `cmp` verified. The four planner artifacts are additional planning outputs, not implementation paths.

| File | Required change | Pre-change sha256 | Post-change sha256 | Versus applied designer r2 |
|---|---|---|---|---|
| `cmd/ailang-worldd/call_failure_line_test.go` | TR; with/without reason rendering | `ABSENT` | `52651e7e6387d07dffa4d9b7450246058a3c3669742f187b0e0ffdebc9fe89fc` | IDENTICAL |
| `cmd/ailang-worldd/mcpclient.go` | Print the optional reason | `91be874ff41288331eaa078ef432011afac93f0382a1a824f48b609a9a56ddaa` | `291703f0184eb0ef21a12e597aaa5a46f48fa58487a32fb4abc36459fbb6c96d` | IDENTICAL |
| `cmd/ailang-worldd/why.go` | Print record failure on the effect link | `47080a095198489b3acaddce150f4095b3bdfd2a24f17393cb308f271252cdd1` | `d5b2abb7732ce3355b95c9b2a07e959aeefc26aba3a1b531313cad698f918844` | IDENTICAL |
| `cmd/ailang-worldd/why_test.go` | TW and failure-capable walk rig | `48b7089c74259df09ab206dc48e5f1630464651e544546ac09a55348ceef31ef` | `79730a25230b5ab404348aed626066ac8c5305f6ff5b679ed35d7300fc2ec3ac` | IDENTICAL |
| `docs/QUICKSTART.md` | Failure-line usage and compatibility/rollback docs | `48f09b67a0e8cb6504fb5c23a3ca69826471e69c54d98361f7db8c43baaaa7a4` | `715d845827cac15fbdb6f0f7dd80146b28ddda9fc342087133683e0da8833387` | IDENTICAL |
| `host/broker/broker.go` | Reason carrier, replay, validation, r3 mark-before-write error | `788f6eaef667182abc4f671ca3abc57a78afce56981d4706340f19601b8dcc7f` | `4684030cfd3d438ca3fa4031b7744e2290db303c1e207943665eacafd46f43ee` | DIFFERS (r3) |
| `host/broker/broker_test.go` | Test store delegates marker | `85a8e1cc0688370aee575292d56bfe04122966c956058cf8545ceda77f7ac1ac` | `6090016bd96d4e311f16af7f6cd2fa3f2770db809c559205e86d51348ad11c39` | IDENTICAL |
| `host/broker/failure.go` | Closed classes; bounded sanitiser; classifier and disclosure rules | `ABSENT` | `fe5a82c7840f2143dfdb27e235c4728b33d214868a5681d77cf404b9ac38c92c` | IDENTICAL |
| `host/broker/failure_test.go` | TL/TG/TK/TX/TV; r3 fault injection and accounting assertions | `ABSENT` | `6263ffc560dafc552aecb6d51cac4f3f4af2ee41b25e928840fc19797634c7c2` | DIFFERS (r3) |
| `host/broker/handler_error_repro_test.go` | Test store delegates marker | `14eb30e43b49db7e46cb1e53c157ea1acd059c87c009ce46202a6a6dc8dcd67e` | `984a307ab1c2a3493fa3c01ab35e2287a6909a5fa8171ac7b7948c2f87b50af6` | IDENTICAL |
| `host/broker/handlers_exec.go` | Classify Execute and cache failures | `d85457256d78add86ac475fc0eb77bcf8cc373082c27139ba3d9f7fd023b2ad0` | `996c5ee11daed4b94f982c3eccebe11ca458f41c2d5fa1ff041887ebf96c9697` | IDENTICAL |
| `host/broker/handlers_test.go` | Test store delegates marker | `fbf7ddc4805d14b530c7c5a2e98dac076952b570ebeaf4db27a40cedbd943b17` | `ecbf80d29ada65310e9e9fe7bdf47f4db286e020ea27dd76b85050b657ce902e` | IDENTICAL |
| `host/broker/record.go` | Optional last failure field and consistency checks | `57a0ec0e8b09c0ea47223f2465681590773bc9d75a2076cda6349898202cd187` | `184009f5fbcebb04e97bdfb70cd28c9dbfeab670f75e93f3de6476cadb4203b8` | IDENTICAL |
| `host/broker/registry_publish_test.go` | Test store delegates marker | `cdc43738affc360676267bc396005be797d9c57e2193d88982e54f20fe29d6ce` | `e09bc1d24d48d191b50c12af35759f8cd523655c372a0478bb9de3495b188ea5` | IDENTICAL |
| `host/coordinator/effect_failure_reason_test.go` | TC; live and replay output bytes | `ABSENT` | `5cbc93c6c94efff8574eb09e6d910dac269428051db8e3877eb4f9e2acad8132` | IDENTICAL |
| `host/coordinator/effectful.go` | Copy reason without changing finish input | `f003b0de143c9dae175e07378ad659361dc9cd819476bc1d61a127c2d0a8e70b` | `b6422dfeae92ec0cac6b31f228d93c4f73bd154f5b43e75231bc3f6e9a4c09d4` | IDENTICAL |
| `host/projection/effect_failure_log_test.go` | TP; MCP/A2A output and log-once | `ABSENT` | `4e77608a9d9c88154976121f0cef9aa70fa654ed2c7bfbcf75eb4b409cee6a13` | IDENTICAL |
| `host/projection/mcp.go` | Log classified MCP failures | `e456f9656b5c49fa04c38dca13a922c7bbf9a90da824f9bea43bbb11a7a9289b` | `f6ffddc79c7a2c0850d6db449cfde495a8921ba3066aba486627bc6ea4908301` | IDENTICAL |
| `host/projection/projection.go` | Log classified A2A failures; suppress reconciled log | `9aa265f9bd9a33635fa53aa23fd326f91b3bb50b387126de5f1d3c2878dbfad1` | `8dbe67ec63254d42e8bbafe038dba5d0aaaa8e11bf968cefd144f7a9e0dc4533` | IDENTICAL |
| `host/store/deadline_guard_test.go` | Marker guard fixture | `9219a4aa634b01d514ecef747d45f0d931329df484297a9327d9a16fe3a4b24c` | `6803bd2044eeaa9118d96c90797022512aa0b367f89b1aad37d9a6e83997356d` | IDENTICAL |
| `host/store/schema_version_test.go` | TS and future-schema fixtures | `d9413412818a85c4c987e2d28fcd2e306772304304203a9c5464839d611bbd73` | `6aa056dc27ddfce6e29cd2d4c143e25afffecf9233a1ee853b10d3662f20c23f` | IDENTICAL |
| `host/store/store.go` | Lazy schema 5-to-6 marker; accept 6/refuse above 6 | `050f46a6643204272be9f4542db2174c13c9fea4192381214320c6eb6ff0cee3` | `e0e10bcd97a80aa8213a8a63ce03bf76d3eef12a91bd6156570063db49dcc340` | IDENTICAL |

Byte bank: `design_docs/verification/world-iter360/plan-M1.diff.txt` (1578 lines, 22 files), sha256 `83ef9f9db480ed5f9488fbcf7e043fde29daea1622605b45fc541b5896c1fdb0`. Designer r2 input sha256: `02259c3bdc1c942bf94b62b5e527a6ee95492430fd59878d449b98152c050dd4`. Apply only this combined bank; do not apply r2 and then a second patch.

Measured per-file `diff` between the designer-applied archive and bank-applied archive: `host/broker/broker.go` differs because a mark error now takes the existing storage-error path rather than clearing `reason`; its fallback comment is updated. `host/broker/failure_test.go` differs because TV now injects a stable mark error, checks type and wrapping, checks nil result/zero refs, zero recorded objects and outcome calls, one handler execution and budget 4 after cost 1 from 5. Test-only counters support those assertions. Every other file is byte-identical, as the final column states. A scan of all r2 diff additions for fallback claims found these two files only; no other implementation/doc change is justified.

Scratch proof used only `.plan-scratch/` below this worktree:

```bash
mkdir -p .plan-scratch/base .plan-scratch/designer .plan-scratch/final .plan-scratch/applied
git archive HEAD | tar -x -C .plan-scratch/base
git archive HEAD | tar -x -C .plan-scratch/designer
git archive HEAD | tar -x -C .plan-scratch/final
git archive HEAD | tar -x -C .plan-scratch/applied
git apply --directory=.plan-scratch/designer design_docs/planned/w-exec-failure-carries-a-reason.patch
# r3 was authored only in .plan-scratch/final after scratch-only r2 application.
git apply --check design_docs/verification/world-iter360/plan-M1.diff.txt
git apply --check --directory=.plan-scratch/applied design_docs/verification/world-iter360/plan-M1.diff.txt
git apply --directory=.plan-scratch/applied design_docs/verification/world-iter360/plan-M1.diff.txt
# For each file in the table: diff -u .plan-scratch/designer/<file> .plan-scratch/applied/<file>
# rc 1 only for broker.go and failure_test.go, rc 0 for the other 20 files.
# sha256 was measured from .plan-scratch/base/<file> and .plan-scratch/applied/<file>.
# .plan-scratch/ was removed after verification.
```

The archive directories are below a Git worktree. Run `git apply --directory=...` from the root: merely using `git -C` inside an archive still discovers the parent Git repository and can skip archive paths. The planner used scratch-prefix applications only; root application appears below solely as an executor command.

## Acceptance commands — M1 stays planned until A8

Run from the implementation tree root, in Bash, sequentially. Leave the controller's `TMPDIR` unchanged; never relocate it under `~/dev`. Do not pin `GOTOOLCHAIN=local`: `go.mod` requires 1.26.6. A3–A8 and the mutation runner run **outside the planner sandbox**. No skipped required identity satisfies a gate. A listener result inside this sandbox is **UNINFORMATIVE UNDER SANDBOX**. Full package suites and the full mutation runner were not run by the planner.

```bash
set -euo pipefail
export AILANG_BIN="$HOME/.pinned-ailang/ailang"
export EVIDENCE="$PWD/design_docs/verification/world-iter360/controller-M1"
mkdir -p "$EVIDENCE"
go version
"$AILANG_BIN" --version
```

**A0 — base, prototype and byte bank before application.**

```bash
test "$(git rev-parse HEAD)" = adc03d4223a2ff89412d49d1499420684f1b2ccb
shasum -a 256 -c <<'SHA'
02259c3bdc1c942bf94b62b5e527a6ee95492430fd59878d449b98152c050dd4  design_docs/planned/w-exec-failure-carries-a-reason.patch
83ef9f9db480ed5f9488fbcf7e043fde29daea1622605b45fc541b5896c1fdb0  design_docs/verification/world-iter360/plan-M1.diff.txt
SHA
python3 -I - <<'PRE'
from pathlib import Path
import hashlib
expected={'cmd/ailang-worldd/call_failure_line_test.go': 'ABSENT', 'cmd/ailang-worldd/mcpclient.go': '91be874ff41288331eaa078ef432011afac93f0382a1a824f48b609a9a56ddaa', 'cmd/ailang-worldd/why.go': '47080a095198489b3acaddce150f4095b3bdfd2a24f17393cb308f271252cdd1', 'cmd/ailang-worldd/why_test.go': '48b7089c74259df09ab206dc48e5f1630464651e544546ac09a55348ceef31ef', 'docs/QUICKSTART.md': '48f09b67a0e8cb6504fb5c23a3ca69826471e69c54d98361f7db8c43baaaa7a4', 'host/broker/broker.go': '788f6eaef667182abc4f671ca3abc57a78afce56981d4706340f19601b8dcc7f', 'host/broker/broker_test.go': '85a8e1cc0688370aee575292d56bfe04122966c956058cf8545ceda77f7ac1ac', 'host/broker/failure.go': 'ABSENT', 'host/broker/failure_test.go': 'ABSENT', 'host/broker/handler_error_repro_test.go': '14eb30e43b49db7e46cb1e53c157ea1acd059c87c009ce46202a6a6dc8dcd67e', 'host/broker/handlers_exec.go': 'd85457256d78add86ac475fc0eb77bcf8cc373082c27139ba3d9f7fd023b2ad0', 'host/broker/handlers_test.go': 'fbf7ddc4805d14b530c7c5a2e98dac076952b570ebeaf4db27a40cedbd943b17', 'host/broker/record.go': '57a0ec0e8b09c0ea47223f2465681590773bc9d75a2076cda6349898202cd187', 'host/broker/registry_publish_test.go': 'cdc43738affc360676267bc396005be797d9c57e2193d88982e54f20fe29d6ce', 'host/coordinator/effect_failure_reason_test.go': 'ABSENT', 'host/coordinator/effectful.go': 'f003b0de143c9dae175e07378ad659361dc9cd819476bc1d61a127c2d0a8e70b', 'host/projection/effect_failure_log_test.go': 'ABSENT', 'host/projection/mcp.go': 'e456f9656b5c49fa04c38dca13a922c7bbf9a90da824f9bea43bbb11a7a9289b', 'host/projection/projection.go': '9aa265f9bd9a33635fa53aa23fd326f91b3bb50b387126de5f1d3c2878dbfad1', 'host/store/deadline_guard_test.go': '9219a4aa634b01d514ecef747d45f0d931329df484297a9327d9a16fe3a4b24c', 'host/store/schema_version_test.go': 'd9413412818a85c4c987e2d28fcd2e306772304304203a9c5464839d611bbd73', 'host/store/store.go': '050f46a6643204272be9f4542db2174c13c9fea4192381214320c6eb6ff0cee3'}
for f,want in expected.items():
    p=Path(f)
    assert (not p.exists()) if want=='ABSENT' else (p.is_file() and hashlib.sha256(p.read_bytes()).hexdigest()==want),f
PRE
git apply --check design_docs/verification/world-iter360/plan-M1.diff.txt
```

**A1 — executor applies once, then verifies every resulting file.**

```bash
git apply design_docs/verification/world-iter360/plan-M1.diff.txt
shasum -a 256 -c <<'SHA'
52651e7e6387d07dffa4d9b7450246058a3c3669742f187b0e0ffdebc9fe89fc  cmd/ailang-worldd/call_failure_line_test.go
291703f0184eb0ef21a12e597aaa5a46f48fa58487a32fb4abc36459fbb6c96d  cmd/ailang-worldd/mcpclient.go
d5b2abb7732ce3355b95c9b2a07e959aeefc26aba3a1b531313cad698f918844  cmd/ailang-worldd/why.go
79730a25230b5ab404348aed626066ac8c5305f6ff5b679ed35d7300fc2ec3ac  cmd/ailang-worldd/why_test.go
715d845827cac15fbdb6f0f7dd80146b28ddda9fc342087133683e0da8833387  docs/QUICKSTART.md
4684030cfd3d438ca3fa4031b7744e2290db303c1e207943665eacafd46f43ee  host/broker/broker.go
6090016bd96d4e311f16af7f6cd2fa3f2770db809c559205e86d51348ad11c39  host/broker/broker_test.go
fe5a82c7840f2143dfdb27e235c4728b33d214868a5681d77cf404b9ac38c92c  host/broker/failure.go
6263ffc560dafc552aecb6d51cac4f3f4af2ee41b25e928840fc19797634c7c2  host/broker/failure_test.go
984a307ab1c2a3493fa3c01ab35e2287a6909a5fa8171ac7b7948c2f87b50af6  host/broker/handler_error_repro_test.go
996c5ee11daed4b94f982c3eccebe11ca458f41c2d5fa1ff041887ebf96c9697  host/broker/handlers_exec.go
ecbf80d29ada65310e9e9fe7bdf47f4db286e020ea27dd76b85050b657ce902e  host/broker/handlers_test.go
184009f5fbcebb04e97bdfb70cd28c9dbfeab670f75e93f3de6476cadb4203b8  host/broker/record.go
e09bc1d24d48d191b50c12af35759f8cd523655c372a0478bb9de3495b188ea5  host/broker/registry_publish_test.go
5cbc93c6c94efff8574eb09e6d910dac269428051db8e3877eb4f9e2acad8132  host/coordinator/effect_failure_reason_test.go
b6422dfeae92ec0cac6b31f228d93c4f73bd154f5b43e75231bc3f6e9a4c09d4  host/coordinator/effectful.go
4e77608a9d9c88154976121f0cef9aa70fa654ed2c7bfbcf75eb4b409cee6a13  host/projection/effect_failure_log_test.go
f6ffddc79c7a2c0850d6db449cfde495a8921ba3066aba486627bc6ea4908301  host/projection/mcp.go
8dbe67ec63254d42e8bbafe038dba5d0aaaa8e11bf968cefd144f7a9e0dc4533  host/projection/projection.go
6803bd2044eeaa9118d96c90797022512aa0b367f89b1aad37d9a6e83997356d  host/store/deadline_guard_test.go
6aa056dc27ddfce6e29cd2d4c143e25afffecf9233a1ee853b10d3662f20c23f  host/store/schema_version_test.go
e0e10bcd97a80aa8213a8a63ce03bf76d3eef12a91bd6156570063db49dcc340  host/store/store.go
SHA
```

**A2 — formatting, build, vet and runner syntax.**

```bash
gofmt -l cmd/ailang-worldd/call_failure_line_test.go cmd/ailang-worldd/mcpclient.go cmd/ailang-worldd/why.go cmd/ailang-worldd/why_test.go host/broker/broker.go host/broker/broker_test.go host/broker/failure.go host/broker/failure_test.go host/broker/handler_error_repro_test.go host/broker/handlers_exec.go host/broker/handlers_test.go host/broker/record.go host/broker/registry_publish_test.go host/coordinator/effect_failure_reason_test.go host/coordinator/effectful.go host/projection/effect_failure_log_test.go host/projection/mcp.go host/projection/projection.go host/store/deadline_guard_test.go host/store/schema_version_test.go host/store/store.go > "$EVIDENCE/gofmt.log"
cat "$EVIDENCE/gofmt.log"
test ! -s "$EVIDENCE/gofmt.log"
go build ./host/broker/ ./host/store/ ./host/coordinator/ ./host/projection/ ./cmd/ailang-worldd/
go vet ./host/broker/ ./host/store/ ./host/coordinator/ ./host/projection/ ./cmd/ailang-worldd/
/bin/bash -n design_docs/verification/world-iter360/plan-mutants.sh.txt
```

**A3 — 18 named identities, each once, in its own package.**

```bash
go test -count=1 -run '^(TestFailureLineBoundsAndSanitises|TestRecordFailureGoldenBytes|TestExecFailureClassifies|TestExecFailureClassesReachTheRecord|TestRecordGoldenBytes|TestExecHandlerReverifiesTheStackBeforeEachCall|TestCFJ2HandlerErrorKeepsDebitAndWritesOneFailureRecord|TestSketchRows|TestFailureLineMarksTheStoreFirst)$' -v -timeout 300s ./host/broker/ | tee "$EVIDENCE/host-broker-named.log"
go test -count=1 -run '^(TestMarkFailureLinesRaisesFiveToSixOnce|TestFutureVersionIsRejected|TestOpenReadOnlyEnforcesVersionFive|TestStoreDeadlineGuardDynamic)$' -v -timeout 300s ./host/store/ | tee "$EVIDENCE/host-store-named.log"
go test -count=1 -run '^(TestEffectFailureReasonReachesTheWorldBlockAndReplays|TestEffectReplay)$' -v -timeout 300s ./host/coordinator/ | tee "$EVIDENCE/host-coordinator-named.log"
go test -count=1 -run '^(TestMCPEffectFailureLineReturnedAndLoggedOnce)$' -v -timeout 300s ./host/projection/ | tee "$EVIDENCE/host-projection-named.log"
go test -count=1 -run '^(TestWhyShowsTheFailureLine|TestCallPrintsTheFailureLine)$' -v -timeout 300s ./cmd/ailang-worldd/ | tee "$EVIDENCE/cmd-ailang-worldd-named.log"
python3 -I - "$EVIDENCE" <<'NAMED'
from pathlib import Path
import collections,re,sys
E=Path(sys.argv[1])
TESTS={'host/broker': ['TestFailureLineBoundsAndSanitises', 'TestRecordFailureGoldenBytes', 'TestExecFailureClassifies', 'TestExecFailureClassesReachTheRecord', 'TestRecordGoldenBytes', 'TestExecHandlerReverifiesTheStackBeforeEachCall', 'TestCFJ2HandlerErrorKeepsDebitAndWritesOneFailureRecord', 'TestSketchRows', 'TestFailureLineMarksTheStoreFirst'], 'host/store': ['TestMarkFailureLinesRaisesFiveToSixOnce', 'TestFutureVersionIsRejected', 'TestOpenReadOnlyEnforcesVersionFive', 'TestStoreDeadlineGuardDynamic'], 'host/coordinator': ['TestEffectFailureReasonReachesTheWorldBlockAndReplays', 'TestEffectReplay'], 'host/projection': ['TestMCPEffectFailureLineReturnedAndLoggedOnce'], 'cmd/ailang-worldd': ['TestWhyShowsTheFailureLine', 'TestCallPrintsTheFailureLine']}
for pkg,names in TESTS.items():
    text=(E/(pkg.replace('/','-')+'-named.log')).read_text()
    passed=collections.Counter(re.findall(r'^--- PASS: (\w+) ',text,re.M))
    assert passed==collections.Counter(names),(pkg,passed)
    assert not re.search(r'^\s*--- (SKIP|FAIL):',text,re.M),pkg
NAMED
```

**A4 — stability: all ten new named tests, three times.**

```bash
go test -count=3 -run '^(TestFailureLineBoundsAndSanitises|TestRecordFailureGoldenBytes|TestExecFailureClassifies|TestExecFailureClassesReachTheRecord|TestFailureLineMarksTheStoreFirst)$' -v -timeout 300s ./host/broker/ | tee "$EVIDENCE/host-broker-stability.log"
go test -count=3 -run '^(TestMarkFailureLinesRaisesFiveToSixOnce)$' -v -timeout 300s ./host/store/ | tee "$EVIDENCE/host-store-stability.log"
go test -count=3 -run '^(TestEffectFailureReasonReachesTheWorldBlockAndReplays)$' -v -timeout 300s ./host/coordinator/ | tee "$EVIDENCE/host-coordinator-stability.log"
go test -count=3 -run '^(TestMCPEffectFailureLineReturnedAndLoggedOnce)$' -v -timeout 300s ./host/projection/ | tee "$EVIDENCE/host-projection-stability.log"
go test -count=3 -run '^(TestWhyShowsTheFailureLine|TestCallPrintsTheFailureLine)$' -v -timeout 300s ./cmd/ailang-worldd/ | tee "$EVIDENCE/cmd-ailang-worldd-stability.log"
python3 -I - "$EVIDENCE" <<'STABILITY'
from pathlib import Path
import collections,re,sys
E=Path(sys.argv[1])
TESTS={'host/broker': ['TestFailureLineBoundsAndSanitises', 'TestRecordFailureGoldenBytes', 'TestExecFailureClassifies', 'TestExecFailureClassesReachTheRecord', 'TestFailureLineMarksTheStoreFirst'], 'host/store': ['TestMarkFailureLinesRaisesFiveToSixOnce'], 'host/coordinator': ['TestEffectFailureReasonReachesTheWorldBlockAndReplays'], 'host/projection': ['TestMCPEffectFailureLineReturnedAndLoggedOnce'], 'cmd/ailang-worldd': ['TestWhyShowsTheFailureLine', 'TestCallPrintsTheFailureLine']}
for pkg,names in TESTS.items():
    text=(E/(pkg.replace('/','-')+'-stability.log')).read_text()
    passed=collections.Counter(re.findall(r'^--- PASS: (\w+) ',text,re.M))
    assert passed==collections.Counter({n:3 for n in names}),(pkg,passed)
    assert not re.search(r'^\s*--- (SKIP|FAIL):',text,re.M),pkg
STABILITY
```

**A5 — controller mutation matrix outside sandbox.**

```bash
bash design_docs/verification/world-iter360/plan-mutants.sh.txt "$PWD" | tee "$EVIDENCE/mutants.log"
test "$(grep -c '^EXPECTATION MATCH:' "$EVIDENCE/mutants.log")" -eq 212
if grep -q '^EXPECTATION MISMATCH:' "$EVIDENCE/mutants.log"; then exit 1; fi
test "$(grep -c '^RESTORE cmp ok: .*tests=' "$EVIDENCE/mutants.log")" -eq 212
grep -Fq 'GO FILE CENSUS unchanged; no extra .go files remain' "$EVIDENCE/mutants.log"
shasum -a 256 -c <<'SHA'
52651e7e6387d07dffa4d9b7450246058a3c3669742f187b0e0ffdebc9fe89fc  cmd/ailang-worldd/call_failure_line_test.go
291703f0184eb0ef21a12e597aaa5a46f48fa58487a32fb4abc36459fbb6c96d  cmd/ailang-worldd/mcpclient.go
d5b2abb7732ce3355b95c9b2a07e959aeefc26aba3a1b531313cad698f918844  cmd/ailang-worldd/why.go
79730a25230b5ab404348aed626066ac8c5305f6ff5b679ed35d7300fc2ec3ac  cmd/ailang-worldd/why_test.go
715d845827cac15fbdb6f0f7dd80146b28ddda9fc342087133683e0da8833387  docs/QUICKSTART.md
4684030cfd3d438ca3fa4031b7744e2290db303c1e207943665eacafd46f43ee  host/broker/broker.go
6090016bd96d4e311f16af7f6cd2fa3f2770db809c559205e86d51348ad11c39  host/broker/broker_test.go
fe5a82c7840f2143dfdb27e235c4728b33d214868a5681d77cf404b9ac38c92c  host/broker/failure.go
6263ffc560dafc552aecb6d51cac4f3f4af2ee41b25e928840fc19797634c7c2  host/broker/failure_test.go
984a307ab1c2a3493fa3c01ab35e2287a6909a5fa8171ac7b7948c2f87b50af6  host/broker/handler_error_repro_test.go
996c5ee11daed4b94f982c3eccebe11ca458f41c2d5fa1ff041887ebf96c9697  host/broker/handlers_exec.go
ecbf80d29ada65310e9e9fe7bdf47f4db286e020ea27dd76b85050b657ce902e  host/broker/handlers_test.go
184009f5fbcebb04e97bdfb70cd28c9dbfeab670f75e93f3de6476cadb4203b8  host/broker/record.go
e09bc1d24d48d191b50c12af35759f8cd523655c372a0478bb9de3495b188ea5  host/broker/registry_publish_test.go
5cbc93c6c94efff8574eb09e6d910dac269428051db8e3877eb4f9e2acad8132  host/coordinator/effect_failure_reason_test.go
b6422dfeae92ec0cac6b31f228d93c4f73bd154f5b43e75231bc3f6e9a4c09d4  host/coordinator/effectful.go
4e77608a9d9c88154976121f0cef9aa70fa654ed2c7bfbcf75eb4b409cee6a13  host/projection/effect_failure_log_test.go
f6ffddc79c7a2c0850d6db449cfde495a8921ba3066aba486627bc6ea4908301  host/projection/mcp.go
8dbe67ec63254d42e8bbafe038dba5d0aaaa8e11bf968cefd144f7a9e0dc4533  host/projection/projection.go
6803bd2044eeaa9118d96c90797022512aa0b367f89b1aad37d9a6e83997356d  host/store/deadline_guard_test.go
6aa056dc27ddfce6e29cd2d4c143e25afffecf9233a1ee853b10d3662f20c23f  host/store/schema_version_test.go
e0e10bcd97a80aa8213a8a63ce03bf76d3eef12a91bd6156570063db49dcc340  host/store/store.go
SHA
```

**A6 — reproduce V20 with the planner base as the old binary.**

```bash
compat=$(mktemp -d "$PWD/.plan-compat.XXXXXX")
trap 'rm -rf "$compat"' EXIT
mkdir "$compat/old"
git archive adc03d4223a2ff89412d49d1499420684f1b2ccb | tar -x -C "$compat/old"
cat > "$compat/mark.go" <<'MARK'
package main
import (
 "context"
 "errors"
 "fmt"
 "os"
 "time"
 "github.com/sunholo-data/ailang-world/host/broker"
 "github.com/sunholo-data/ailang-world/host/store"
)
func main() {
 ctx,cancel:=context.WithTimeout(context.Background(),10*time.Second); defer cancel()
 path:=os.Args[1]
 st,err:=store.Open(path); if err!=nil { panic(err) }; defer st.Close()
 h:=broker.HandlerFunc(func(context.Context,broker.EffectRequest,[]byte)([]byte,error){return nil,&broker.HandlerFailure{Class:broker.FailSeedCache,Detail:"GOCACHE: timed out after 9s",Err:errors.New("injected handler error")}})
 session:=broker.NewSession(st,"compat",[]broker.Capability{{Effect:broker.EffectWorkspaceExec,Scope:broker.WorkspaceScope,ExpiresAt:100,Budget:5}},broker.Registry{broker.EffectWorkspaceExec:h})
 _,ref,err:=session.Invoke(ctx,broker.EffectRequest{Effect:broker.EffectWorkspaceExec,Scope:broker.WorkspaceScope,Cost:1,Now:1},[]byte(`{}`))
 var failed *broker.EffectFailedError
 if !errors.As(err,&failed) || ref.IsZero() { panic(fmt.Sprintf("Invoke: %v",err)) }
 obj,ok,err:=st.GetObject(ctx,ref); if err!=nil || !ok { panic("missing effect record") }
 rec,err:=broker.DecodeRecord(obj.Payload); if err!=nil || rec.Failure!=failed.Reason || rec.Failure=="" { panic("missing failure line") }
 v,err:=store.WorldSchemaVersion(path); if err!=nil || v!=6 { panic("store not marked 6") }
 fmt.Printf("lined record %s; store user_version = %d\n",ref,v)
}
MARK
cat > "$compat/old-open.go" <<'OLD'
package main
import (
 "errors"
 "fmt"
 "os"
 "strings"
 "github.com/sunholo-data/ailang-world/host/store"
)
func main() {
 for _,op:=range []struct{name string;open func(string)(*store.Store,error)}{{"Open",store.Open},{"OpenReadOnly",store.OpenReadOnly}} {
  st,err:=op.open(os.Args[1]); if st!=nil {st.Close();panic("older binary opened marked store")}
  var future *store.FutureSchemaVersionError
  if !errors.As(err,&future) || future.Found!=6 || future.Current!=5 || !strings.Contains(err.Error(),"use a compatible binary") {panic(fmt.Sprintf("wrong refusal: %v",err))}
  fmt.Printf("%s: %v\n",op.name,err)
 }
}
OLD
go run "$compat/mark.go" "$compat/marked.db" | tee "$EVIDENCE/compat-mark.log"
(cd "$compat/old" && go run "$compat/old-open.go" "$compat/marked.db") | tee "$EVIDENCE/compat-old-open.log"
rm -rf "$compat"
trap - EXIT
```

**A7 — full touched-package suites and repository gates, outside sandbox.**

```bash
go test -count=1 -timeout 600s ./host/broker/ ./host/store/ ./host/coordinator/ ./host/projection/ ./cmd/ailang-worldd/ | tee "$EVIDENCE/touched-packages.log"
AILANG_BIN="$HOME/.pinned-ailang/ailang" ./scripts/verify_ail.sh | tee "$EVIDENCE/verify-ail.log"
VERIFY_GO_LEGS=plain AILANG_BIN="$HOME/.pinned-ailang/ailang" ./scripts/verify_go.sh | tee "$EVIDENCE/verify-go-plain.log"
shasum -a 256 -c <<'SHA'
52651e7e6387d07dffa4d9b7450246058a3c3669742f187b0e0ffdebc9fe89fc  cmd/ailang-worldd/call_failure_line_test.go
291703f0184eb0ef21a12e597aaa5a46f48fa58487a32fb4abc36459fbb6c96d  cmd/ailang-worldd/mcpclient.go
d5b2abb7732ce3355b95c9b2a07e959aeefc26aba3a1b531313cad698f918844  cmd/ailang-worldd/why.go
79730a25230b5ab404348aed626066ac8c5305f6ff5b679ed35d7300fc2ec3ac  cmd/ailang-worldd/why_test.go
715d845827cac15fbdb6f0f7dd80146b28ddda9fc342087133683e0da8833387  docs/QUICKSTART.md
4684030cfd3d438ca3fa4031b7744e2290db303c1e207943665eacafd46f43ee  host/broker/broker.go
6090016bd96d4e311f16af7f6cd2fa3f2770db809c559205e86d51348ad11c39  host/broker/broker_test.go
fe5a82c7840f2143dfdb27e235c4728b33d214868a5681d77cf404b9ac38c92c  host/broker/failure.go
6263ffc560dafc552aecb6d51cac4f3f4af2ee41b25e928840fc19797634c7c2  host/broker/failure_test.go
984a307ab1c2a3493fa3c01ab35e2287a6909a5fa8171ac7b7948c2f87b50af6  host/broker/handler_error_repro_test.go
996c5ee11daed4b94f982c3eccebe11ca458f41c2d5fa1ff041887ebf96c9697  host/broker/handlers_exec.go
ecbf80d29ada65310e9e9fe7bdf47f4db286e020ea27dd76b85050b657ce902e  host/broker/handlers_test.go
184009f5fbcebb04e97bdfb70cd28c9dbfeab670f75e93f3de6476cadb4203b8  host/broker/record.go
e09bc1d24d48d191b50c12af35759f8cd523655c372a0478bb9de3495b188ea5  host/broker/registry_publish_test.go
5cbc93c6c94efff8574eb09e6d910dac269428051db8e3877eb4f9e2acad8132  host/coordinator/effect_failure_reason_test.go
b6422dfeae92ec0cac6b31f228d93c4f73bd154f5b43e75231bc3f6e9a4c09d4  host/coordinator/effectful.go
4e77608a9d9c88154976121f0cef9aa70fa654ed2c7bfbcf75eb4b409cee6a13  host/projection/effect_failure_log_test.go
f6ffddc79c7a2c0850d6db449cfde495a8921ba3066aba486627bc6ea4908301  host/projection/mcp.go
8dbe67ec63254d42e8bbafe038dba5d0aaaa8e11bf968cefd144f7a9e0dc4533  host/projection/projection.go
6803bd2044eeaa9118d96c90797022512aa0b367f89b1aad37d9a6e83997356d  host/store/deadline_guard_test.go
6aa056dc27ddfce6e29cd2d4c143e25afffecf9233a1ee853b10d3662f20c23f  host/store/schema_version_test.go
e0e10bcd97a80aa8213a8a63ce03bf76d3eef12a91bd6156570063db49dcc340  host/store/store.go
SHA
```

**A8 — final acceptance: controller verifies the PR head and all three CI checks after push.**

```bash
: "${PR_NUMBER:?controller selects implementation PR}"
: "${CI_RUN_ID:?controller selects CI run on pushed implementation head}"
gh pr view "$PR_NUMBER" --json headRefOid > "$EVIDENCE/pr-head.json"
gh pr checks "$PR_NUMBER" --json name,state,link | tee "$EVIDENCE/pr-checks.json"
gh run view "$CI_RUN_ID" --json headSha,conclusion,jobs,url > "$EVIDENCE/ci.json"
gh run view "$CI_RUN_ID" --log > "$EVIDENCE/ci.log"
python3 -I - "$EVIDENCE" <<'CI'
from pathlib import Path
import json,sys
E=Path(sys.argv[1]); head=json.loads((E/'pr-head.json').read_text())['headRefOid']
run=json.loads((E/'ci.json').read_text()); checks=json.loads((E/'pr-checks.json').read_text())
required={'ailang-code verify gate','go host build + test gate','go race gate (go test -race, own job)'}
assert run['headSha']==head and run['conclusion']=='success',run
jobs={j['name']:j for j in run['jobs']}
assert set(jobs)==required and all(j['conclusion']=='success' for j in jobs.values()),jobs
observed={c['name']:c['state'] for c in checks}
assert set(observed)==required and all(v=='SUCCESS' for v in observed.values()),observed
plain='go build + test gate (replay tests run against pinned AILANG_BIN)'
steps=[x for x in jobs['go host build + test gate']['steps'] if x['name']==plain]
assert len(steps)==1 and steps[0]['conclusion']=='success',steps
print('Accepted PR head '+head+'; three CI checks successful; Linux plain gate successful: '+run['url'])
CI
```

`scripts/verify_go.sh` exists at this base; A7's plain leg is required. Its successful Linux CI counterpart is required again at A8, including TX's GNU cp exit-code and node EACCES rows. The CI workflow has exactly three checks, named in A8. All three must be green on the pushed implementation PR head; a local green does not accept M1. No CI workflow edit belongs to M1.

## Named identities and mutation expectations (design §5)

| Alias | Package | Exact test |
|---|---|---|
| TL | `host/broker` | `TestFailureLineBoundsAndSanitises` |
| TG | `host/broker` | `TestRecordFailureGoldenBytes` |
| TK | `host/broker` | `TestExecFailureClassifies` |
| TX | `host/broker` | `TestExecFailureClassesReachTheRecord` |
| Gold | `host/broker` | `TestRecordGoldenBytes` |
| Reverify | `host/broker` | `TestExecHandlerReverifiesTheStackBeforeEachCall` |
| CFJ2 | `host/broker` | `TestCFJ2HandlerErrorKeepsDebitAndWritesOneFailureRecord` |
| Sketch | `host/broker` | `TestSketchRows` |
| TV | `host/broker` | `TestFailureLineMarksTheStoreFirst` |
| TS | `host/store` | `TestMarkFailureLinesRaisesFiveToSixOnce` |
| Future | `host/store` | `TestFutureVersionIsRejected` |
| ROFive | `host/store` | `TestOpenReadOnlyEnforcesVersionFive` |
| Guard | `host/store` | `TestStoreDeadlineGuardDynamic` |
| TC | `host/coordinator` | `TestEffectFailureReasonReachesTheWorldBlockAndReplays` |
| Replay | `host/coordinator` | `TestEffectReplay` |
| TP | `host/projection` | `TestMCPEffectFailureLineReturnedAndLoggedOnce` |
| TW | `cmd/ailang-worldd` | `TestWhyShowsTheFailureLine` |
| TR | `cmd/ailang-worldd` | `TestCallPrintsTheFailureLine` |

Each mutant is applied alone. RED means rc 1 plus a named terminal FAIL, GREEN means rc 0 plus a named terminal PASS. The runner rejects compilation/setup errors, missing identities and skips; it restores every mutated source after every row and checks with `cmp`. Every alias not listed in a mutant's RED column must remain GREEN. Thus CFJ2, Sketch and Replay are GREEN for every mutant. Expected outcomes reproduce the designer's r2 matrix, refined for TV's r3 branch; only the manual silent-fallback kill below is newly measured here.

| Mutant | Edit / load-bearing assertion | Expected RED aliases | Expected GREEN aliases | Triggering assertion (one sufficient named kill) |
|---|---|---|---|---|
| CONTROL (one per package) | Unmutated applied tree; all 18 identities execute | — | ALL | All 18 named terminal identities pass exactly once |
| MUT-NO-OMITEMPTY | `json:"failure,omitempty"` → `json:"failure"` | Gold | ALL except RED | Gold: encoded unclassified/success bytes equal their frozen goldens |
| MUT-RECORD-DROPS-LINE | the failed record omits `Failure: reason` | TC,TV,TX | ALL except RED | TX: decoded rec.Failure equals the class row want |
| MUT-REPLAY-DROPS-REASON | `invokeReplay` omits `Reason: rec.Failure` | TC | ALL except RED | TC: Replay succeeds and equals committed output bytes |
| MUT-COORD-DROPS-REASON | `requestEffect` does not copy `failed.Reason` | TC,TP | ALL except RED | TC: classified output contains the exact reason key/value |
| MUT-WORLD-NO-OMITEMPTY | `json:"reason,omitempty"` → `json:"reason"` | TC | ALL except RED | TC/unclassified: output contains no reason key |
| MUT-ERR-TEXT | the untyped fallback's detail is `err.Error()` | TK,TX | ALL except RED | TK/untyped cause: got.Line() equals fixed sandbox-start text |
| MUT-RAW-PATH | `execDisclosed.name` returns every path verbatim | TK,TX | ALL except RED | TK/archive file: got.Line() names srt archive relative path |
| MUT-DRIFT-WHY-RAW | `driftWhat` returns `Why` verbatim | TK,TX | ALL except RED | TK/undisclosed path: got.Line() equals basename plus unreadable |
| MUT-NO-SANITISE | no character mapping | TL | ALL except RED | TL: control-character fixture equals its cleaned expected string |
| MUT-CONTROL-ONLY | maps control characters but not format characters | TL | ALL except RED | TL: U+202E/U+2028 fixture equals evil gnp.txt next |
| MUT-NO-BOUND | no length cut | TL | ALL except RED | TL: long line length is at most FailureLineMax and ends with ellipsis |
| MUT-CLASS-OPEN | an unknown class is accepted | TG,TL | ALL except RED | TL: unknown class yields empty string |
| MUT-DECODE-NO-CHECK | `DecodeRecord` skips `FailureConsistent` | TG | ALL except RED | TG: malformed failure fixture must return a failure-line refusal |
| MUT-FAILURE-ANY-ARM | a line is allowed on any arm | TG | ALL except RED | TG/on a success: DecodeRecord must return a failure-line refusal |
| MUT-PUT-NO-CHECK **SURVIVOR** | `putRecord` skips `FailureConsistent` | — | ALL | No reachable assertion; by-construction survivor |
| MUT-HF-ERROR-IS-LINE | `HandlerFailure.Error()` returns the line, not the cause | Reverify,TX | ALL except RED | Reverify: err.Error() contains original drift detail want |
| MUT-SEED-UNCLASSIFIED | the seed copy returns a plain error | TX | ALL except RED | TX/seed: EffectFailedError.Reason equals the expected seed-cache line |
| MUT-SEED-TIMEOUT-AS-EXIT | `seedWhat` skips the timeout case | TK,TX | ALL except RED | TK: seedWhat(timeout) equals timed out after duration |
| MUT-LOG-MCP-NONE | no log call on MCP | TP | ALL except RED | TP: log matches the anchored exactly-one MCP-line regexp |
| MUT-LOG-A2A-NONE | no log call on A2A | TP | ALL except RED | TP: log matches the anchored exactly-one A2A-line regexp |
| MUT-LOG-RECONCILED | a reconciled result also logs | TP | ALL except RED | TP: repeated A2A task leaves exactly one A2A log line |
| MUT-WHY-NO-LINE | `why` ignores `Failure` | TW | ALL except RED | TW: stdout contains no result colon plus exact line |
| MUT-CALL-NO-LINE | `call` ignores `reason` | TR | ALL except RED | TR: output contains failed record, em dash and exact line |
| MUT-CLEANUP-UNWIRED **SURVIVOR** | `:356`'s return goes back to the bare `runErr` | — | ALL | No Execute-driven assertion; explicit R3 survivor |
| MUT-NO-MARK | the broker never calls `MarkFailureLines` | TV | ALL except RED | TV/line marks: store user_version must be 6 |
| MUT-MARK-ERROR-IGNORED | a failed mark still writes the line | TV | ALL except RED | TV/mark refused: error must be OutcomeWriteError wrapping markErr |
| MUT-MARK-AFTER-RECORD | the mark moves after `putRecord` (a crash between them would leave a lined record in a version-5 store) | TV | ALL except RED | TV/line marks: joined call order must equal mark,record |
| MUT-STORE-REJECTS-6 | `enforceSchemaVersion` does not accept 6 | TS | ALL except RED | TS: writer/reader reopen of marked store must succeed |
| MUT-FUTURE-BOUND-5 **SURVIVOR** | the future check reads `> currentSchemaVersion` | — | ALL | No differentiating assertion; equivalent survivor |
| MUT-MARK-NOOP | `MarkFailureLines` returns nil without writing | TS,TV | ALL except RED | TS: MarkFailureLines must leave user_version at 6 |
| MUT-MARK-ANY-VERSION | `MarkFailureLines` tries to mark any version | TS | ALL except RED | TS/schema-4: mark must refuse with is not 5 or 6 |
| MUT-HEAD-SHAPE | the record drops the line **and** the coordinator drops the reason: today's behaviour inside the new structure | TC,TP,TV,TX | ALL except RED | TX: rec.Failure equals the expected class line |
| MUT-MARK-FAIL-SILENT | a failed mark writes the record without the line and returns the ordinary `*EffectFailedError` (r2's fallback) | TV | ALL except RED | TV/mark refused: error must be OutcomeWriteError wrapping markErr |
| MUT-MARK-REFUNDS-GRANT | Restore grant budget on mark error; TV must report budget != 4 | TV | ALL except RED | TV/mark refused: grant budget must equal 4 |
| MUT-MARK-RERUNS-HANDLER | Call Execute again on mark error; TV must report handler calls != 1 | TV | ALL except RED | TV/mark refused: handler.calls must equal 1 |

The r2 bank contains 33 scenarios including CONTROL (32 mutants), although the brief calls them "33 mutants". This runner contains **36 scenarios: CONTROL + 35 mutants**, or **212 expectation rows**: five package controls, all 32 banked mutants, r3 MUT-MARK-FAIL-SILENT and two accounting mutants. Rows split each package's tests into expected-red and expected-green groups; each named identity is checked, so one failure cannot mask an incorrectly green identity. Each survivor has five expected-green package rows explicitly labelled SURVIVOR. All 18 identities execute in each scenario (648 terminal named outcomes total).

The three survivors remain explicit residuals, never claimed kills:

- MUT-PUT-NO-CHECK: the sole writer supplies `Line()`, a fixed point by construction; the put guard protects a future writer. TG pins decode; TL pins fixed points.
- MUT-CLEANUP-UNWIRED: no process-kill fault injector drives `cleanup-incomplete` through Execute. TK pins its text, not that wiring.
- MUT-FUTURE-BOUND-5: equivalent, because schema 6 exits from an earlier accepted-version branch; `>5` and `>6` both refuse 7 and above.

The runner takes only `$1`, resolves it before use, and keeps its state, logs, private caches and pristine copies under that root. Pristine copies end in `.pristine`, never `.go`; it removes pristine/runtime caches even on failure and checks a before/after Go-file census. It does not set TMPDIR or GOTOOLCHAIN. It reuses designer exact edit texts wherever r3 did not change the anchor. It has been syntax-checked and all anchors statically checked on the pinned applied tree; the full runner is intentionally unexecuted in sandbox. F1 below limits the literal subprocess-write confinement claim.

## Residuals, findings and exclusions

Design §6 residuals carried into M1:

- R1: after a line may exist, a pre-change binary refuses the schema-6 store. Automatic upgrade rollback across this boundary can stop the service with both builds kept. An older CLI over HTTP is still permitted; `why` shows an unknown-field link failure and `call` omits the reason. Q1/Q2 retain the design's defaults; this plan adds no policy.
- R2: the operator sees the bounded line, not the full Go cause; some failures still need reproduction. Q3 remains deferred.
- R3: cleanup-incomplete is not Execute-driven; MUT-CLEANUP-UNWIRED survives.
- R4: a non-disclosed host file can disclose its base name, never its full path.
- R5: ExecUnconfiguredHandler and ExecRefusalHandler retain unreachable bare misroute errors and refuse real calls.
- R6: Linux TX seed-copy exit 1 and sandbox-start EACCES remain UNMEASURED locally; A8's Linux plain CI is binding.
- R7: `why --json` carries the line in detail text rather than a structured field.

Planner findings (reported, not scope-expanded fixes):

- **F1 — runner confinement versus unchanged TMPDIR.** Existing required tests use `t.TempDir()` (for example TV and TS) and canonicalTempDir, so keeping TMPDIR unchanged allows those subprocesses to create files in the default OS temp directory outside `$1`. No portable Bash/Go-test flag redirects t.TempDir without changing TMPDIR or the tests. Runner-owned artifacts and Go caches are confined; literal confinement of every descendant write requires external filesystem isolation/mapping, or a clarified requirement. The runner does not silently set TMPDIR to satisfy this. For the planner's small socket-free tests only, TMPDIR was explicitly set to the canonical `.plan-scratch/runtime` to obey the separate sandbox no-outside-writes rule. Those results are limited evidence, not the controller's default-TMPDIR acceptance gates.
- **F2 — stale design wording.** §8 still says "a failed mark drops the line"; r3 §3.7, Claim F and §11 supersede it. No approved design file was edited. Likewise §7's sample command applies r2; A0/A1 here apply the combined bank.
- **F3 — artifact counts.** The banked mutation listing has 33 scenarios including CONTROL, rather than 33 actual mutants. All are preserved. `designer-applycheck.txt.txt` lists only 15 paths; it does not prove the 22-file r2/r3 apply. This plan's fresh 22-file apply checks discharge that independently.
- **F4 — marker meaning.** Version 6 means a lined record *may* exist: marking can succeed before putRecord fails. Compatibility table parentheticals "a line written" are not a converse guarantee. Claim F's one-way safety argument remains valid; no record/debit/retry policy was changed.

M1 must NOT touch:

- Row 256 scope: startup probe behavior/budget, registry bootstrap, load-handling design, or left-out-profile reporting.
- Row 257 scope: bind-time seed preparation, cache seeding schedule, cache placement or its separate visibility tests.
- `tools/launchd/*`, shared or copied skills, `.agents/`, `.codex/`, mission/fleet machinery, `~/.ailang/state/*`, the V1 checkout or HOST_CHANGELOG.md.
- Any `.ail` file, including effectbroker/effectplan sketches, `world/` and package projections; RecordConsistent, the effect-record semantic id, budgets, debit law, retry/backoff behavior, or an added production seam.
- Unrelated handlers, episode reporting, build-skew checks for read verbs, CI workflow, compatibility policy beyond approved r3, or any implementation file outside the exact 22-path table.

## Planner verification receipt

- Toolchain: `go version go1.26.6 darwin/arm64`.
- r2 SHA check: exact full hash above; combined root apply-check rc 0; scratch-prefixed apply-check/application rc 0 on the final bank. Fresh archive bytes compared equal to the authored scratch for all 22 files. Per-file comparison to r2: only the two r3 files differ.
- On `.plan-scratch/applied`: touched-package `go build` rc 0; touched-package `go vet` rc 0; `gofmt -l` over all 21 touched Go files rc 0, empty output. Runner `/bin/bash -n` rc 0. All mutation anchors unique, checked statically; no full runner execution.
- Socket-free named tests, `-count=1`: broker TL, TG, TK, TV rc 0 (four terminal passes); coordinator TC and Replay rc 0 (two); projection TP rc 0 (one, in-process ResponseRecorder); CLI TR rc 0 (one); store TS, Future, ROFive and Guard rc 0 (four). Total 12 named terminal passes. TMPDIR was the canonical scratch/runtime directory solely to contain planner test writes.
- Broker TL/TG/TK/TV `-count=3`: rc 0, exactly three terminal passes per name. The remaining stability gates are controller-owned.
- Manual isolated MUT-MARK-FAIL-SILENT (not the runner): restored r2 fallback in scratch broker.go; `go test -count=1 -run '^(TestFailureLineMarksTheStoreFirst)$' ./host/broker/` rc 1. TV's mark-refused row reported ordinary EffectFailedError, wanted OutcomeWriteError wrapping the injected mark error (failure_test.go:370). Restored with `cp` from a non-.go pristine copy; `cmp` rc 0 and all post hashes rechecked. The two other TV rows passed in that run.
- TW opens httptest.NewServer; not run. Its sandbox result would be **UNINFORMATIVE UNDER SANDBOX**. TX/Reverify, full package suites, AIL gate, full plain/race gates, Linux and old-binary probe remain controller-owned; no sandbox result is claimed for them.
- Discarded instrumentation attempt: an initial generated new-file diff lacked `new file mode`, so apply-check returned rc 1 (`dev/null: No such file`). It was corrected before final apply-check and all accepted measurements. Two test commands accidentally proceeded on the unpatched archive: Replay rc 0 without TC discovered, and broker controls rc 1 in Reverify because the planner's noncanonical scratch TMPDIR contained `/../`. This is a base-tree test finding under that explicit invalid temp-path setup, not a product regression, and is not counted as acceptance. The corrected bank was freshly applied; only the twelve identities listed above were run for accepted evidence. No base implementation was repaired and the controller's default-TMPDIR tests are still required.
- Scratch directories and all helper/generated `.go` files were removed at completion. Tracked implementation `git diff` is empty; only the four requested planner deliverables were authored outside scratch. Existing untracked designer evidence was preserved.

