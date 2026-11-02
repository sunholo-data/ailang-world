# Executor report — World iteration 334, queue row 235

Authority: the approved sprint plan, pinned JSON, and design r3. Execution started 2026-11-02 13:21:53 UTC in the assigned workspace-write sandbox. No git write operation was performed.

Working directory: `/Users/voightkampff/dev/sunholo-data/.wt-world-iter334-row235`. Scratch directory: `/tmp/world-iter334-executor.desvqzr3`. HEAD: `a0439a2b3cd8ada68d0af42e81c571689cf43beb` (advanced documentation HEAD; both Go base pins matched). Go 1.26.6 darwin/arm64; AILANG v0.41.0; PATH and AILANG_BIN as prescribed.

## Pins and scope

| Item | Measured SHA256 |
|---|---|
| Applied embedded diff | `8dee691695d632d8afa3f775629d213bf4769606072a0225140e9ec9f66c6d51` |
| Fixed test file | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` |
| Unchanged daemon.go | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` |

Base test pin: `53a708a5264fafe53217d59c58dafa7e25eeabe7f5ef9debc554912fa30a0445`. Base daemon.go pin is identical to its final pin. The unique diff and both extracted tool scripts were verified against JSON pins before applying. `patch -p1 --dry-run` and `patch -p1` both returned 0. Scope: 1 test file, 177 insertions(+), 17 deletions(-). All other worktree additions are this report and executor raw outputs.

## Milestone verdicts

**M1: UNVERIFIED.** Formatting, vet, compile fence, target plain/race controls, handshake controls, AC-7, and every race mutation outcome passed. Every plain mutation invocation stopped before go test because the pinned runner uses an empty array with nounset on macOS Bash 3.2. This is an instrument error, not a kill or survival. All mutations restored both files byte-identically. No plain mutation gate is claimed green.

**M2: UNVERIFIED.** See the command table and findings below. Socket-touching suite results inside workspace-write are UNINFORMATIVE UNDER SANDBOX. Controller must rerun all gates outside the sandbox and obtain all three green CI jobs on the final fix head.

## Acceptance summary

| Acceptance | Observed result |
|---|---|
| M1 formatting/vet/compile | rc 0; gofmt output empty; compile-only explicitly discovers no tests by design |
| M1 target controls | Both named tests PASS 3/3 plain and 3/3 race; 6 serialization subtests per mode, no FAIL/SKIP |
| AC-1 | Race base stagger killed 20/20 at measured :270; plain runner instrument error, no tests launched. Empirical demonstration only |
| AC-2 | M1/M1b/M2 race killed 20/20 top-level and 40/40 subtests each at :372; plain instruments blocked |
| AC-3 | M3 race killed 20/20 in late-second at :392; together PASS 20/20. No survival or adjudication rerun. Plain instrument blocked |
| AC-4 | Handshake PASS 200/200 plain and race. H1/H2/H3 race killed 3/3 each at :426/:426/:422. Plain instruments blocked |
| Additional M1 rows | M6 race killed 20/20 at :375; M7 race killed 20/20 at :372; M7b race survives 20/20 as expected (40/40 subtests). Plain instruments blocked |
| AC-7 | Idle top-level 0.80 s, under 1.2 s; together 0.10 s, late-second 0.70 s; rc 0 |
| AC-5 soak-idle | rc 0; PASS locally; 189.297 s |
| AC-5 soak-loaded | rc 7; UNINFORMATIVE UNDER SANDBOX; 0.298 s |
| AC-6 package-vet | rc 0; PASS locally; 0.454 s |
| AC-6 package-plain | rc 1; UNINFORMATIVE UNDER SANDBOX; 9.306 s |
| AC-6 package-race | rc 1; UNINFORMATIVE UNDER SANDBOX; 9.981 s |
| AC-6 full-gate | rc 1; UNINFORMATIVE UNDER SANDBOX; 190.65 s |
| AC-6 verify-ail | rc 0; PASS locally; 11.146 s |

## Every captured command

PASS/FAIL are measured `--- PASS:` / `--- FAIL:` marker counts including subtests; top-level counts are separate. Non-verbose package/soak commands emit no per-test markers when successful; 0 markers does not assert zero discovery or a measured per-test count. Mutation protocol rc is 0 when the expected red test outcome was correctly observed; its inner go-test rc is recorded separately below.

| Milestone / command log | Exact command | rc | PASS / FAIL / SKIP | Top PASS / FAIL | Wall seconds |
|---|---|---:|---|---|---:|
| M1 [environment](executor-environment.txt) | `pwd<br>git rev-parse HEAD<br>git status --porcelain<br>go version<br>"$AILANG_BIN" --version` | 0 | 0 / 0 / 0 | 0 / 0 | 0.193 |
| M1 [base-hashes](executor-base-hashes.txt) | `shasum -a 256 host/daemon/decisionpacket_daemon_test.go host/daemon/daemon.go > "$WX/base.sha"<br>cat "$WX/base.sha"<br>shasum -a 256 "$WX/plan.diff"` | 0 | 0 / 0 / 0 | 0 / 0 | 0.026 |
| M1 [archive-base](executor-archive-base.txt) | `git archive HEAD \| tar -x -C "$WX/base"` | 0 | 0 / 0 / 0 | 0 / 0 | 1.523 |
| M1 [A-plain-protocol](executor-A-plain-protocol.txt) | `"$WX/tools/runmut.sh" "/tmp/world-iter334-executor.desvqzr3/base" A plain 20` | 1 | 0 / 0 / 0 | 0 / 0 | 0.317 |
| M1 [A-race-protocol](executor-A-race-protocol.txt) | `"$WX/tools/runmut.sh" "/tmp/world-iter334-executor.desvqzr3/base" A race 20` | 0 | 0 / 20 / 0 | 0 / 20 | 27.689 |
| M1 [patch-dry-run](executor-patch-dry-run.txt) | `patch -p1 --dry-run < "$WX/plan.diff"` | 0 | 0 / 0 / 0 | 0 / 0 | 0.006 |
| M1 [patch-apply](executor-patch-apply.txt) | `patch -p1 < "$WX/plan.diff"` | 0 | 0 / 0 / 0 | 0 / 0 | 0.005 |
| M1 [fixed-hashes](executor-fixed-hashes.txt) | `shasum -a 256 host/daemon/decisionpacket_daemon_test.go host/daemon/daemon.go > "$WX/fixed.sha"<br>cat "$WX/fixed.sha"` | 0 | 0 / 0 / 0 | 0 / 0 | 0.011 |
| M1 [gofmt](executor-gofmt.txt) | `perl -e 'alarm 180; exec @ARGV' gofmt -l host/daemon/` | 0 | 0 / 0 / 0 | 0 / 0 | 0.034 |
| M1 [vet](executor-vet.txt) | `perl -e 'alarm 300; exec @ARGV' go vet ./host/daemon/` | 0 | 0 / 0 / 0 | 0 / 0 | 0.404 |
| M1 [compile](executor-compile.txt) | `go test -timeout 180s -run '^$' ./host/daemon/` | 0 | 0 / 0 / 0 | 0 / 0 | 2.067 |
| M1 [target-plain](executor-target-plain.txt) | `go test -timeout 600s ./host/daemon -run '^(TestHumanApproveIsSerializedAcrossEpisodes\|TestConcurrencyProbeNamesAMissedArrival)$' -count=3 -v` | 0 | 12 / 0 / 0 | 6 / 0 | 3.093 |
| M1 [target-race](executor-target-race.txt) | `go test -race -timeout 600s ./host/daemon -run '^(TestHumanApproveIsSerializedAcrossEpisodes\|TestConcurrencyProbeNamesAMissedArrival)$' -count=3 -v` | 0 | 12 / 0 / 0 | 6 / 0 | 7.97 |
| M1 [archive-fixed](executor-archive-fixed.txt) | `git archive HEAD \| tar -x -C "$WX/mutants"<br>patch -d "$WX/mutants" -p1 < "$WX/plan.diff"` | 0 | 0 / 0 / 0 | 0 / 0 | 0.558 |
| M1 [M1-plain-protocol](executor-M1-plain-protocol.txt) | `"$WX/tools/runmut.sh" "/tmp/world-iter334-executor.desvqzr3/mutants" M1 plain 20` | 1 | 0 / 0 / 0 | 0 / 0 | 0.052 |
| M1 [M1-race-protocol](executor-M1-race-protocol.txt) | `"$WX/tools/runmut.sh" "/tmp/world-iter334-executor.desvqzr3/mutants" M1 race 20` | 0 | 0 / 60 / 0 | 0 / 20 | 12.751 |
| M1 [M1b-plain-protocol](executor-M1b-plain-protocol.txt) | `"$WX/tools/runmut.sh" "/tmp/world-iter334-executor.desvqzr3/mutants" M1b plain 20` | 1 | 0 / 0 / 0 | 0 / 0 | 0.052 |
| M1 [M1b-race-protocol](executor-M1b-race-protocol.txt) | `"$WX/tools/runmut.sh" "/tmp/world-iter334-executor.desvqzr3/mutants" M1b race 20` | 0 | 0 / 60 / 0 | 0 / 20 | 11.657 |
| M1 [M2-plain-protocol](executor-M2-plain-protocol.txt) | `"$WX/tools/runmut.sh" "/tmp/world-iter334-executor.desvqzr3/mutants" M2 plain 20` | 1 | 0 / 0 / 0 | 0 / 0 | 0.052 |
| M1 [M2-race-protocol](executor-M2-race-protocol.txt) | `"$WX/tools/runmut.sh" "/tmp/world-iter334-executor.desvqzr3/mutants" M2 race 20` | 0 | 0 / 60 / 0 | 0 / 20 | 11.63 |
| M1 [M3-plain-protocol](executor-M3-plain-protocol.txt) | `"$WX/tools/runmut.sh" "/tmp/world-iter334-executor.desvqzr3/mutants" M3 plain 20` | 1 | 0 / 0 / 0 | 0 / 0 | 0.051 |
| M1 [M3-race-protocol](executor-M3-race-protocol.txt) | `"$WX/tools/runmut.sh" "/tmp/world-iter334-executor.desvqzr3/mutants" M3 race 20` | 0 | 20 / 40 / 0 | 0 / 20 | 23.66 |
| M1 [handshake-plain](executor-handshake-plain.txt) | `go test -timeout 600s ./host/daemon -run '^TestConcurrencyProbeNamesAMissedArrival$' -count=200 -v` | 0 | 200 / 0 / 0 | 200 / 0 | 1.658 |
| M1 [handshake-race](executor-handshake-race.txt) | `go test -race -timeout 600s ./host/daemon -run '^TestConcurrencyProbeNamesAMissedArrival$' -count=200 -v` | 0 | 200 / 0 / 0 | 200 / 0 | 2.677 |
| M1 [H1-plain-protocol](executor-H1-plain-protocol.txt) | `"$WX/tools/runmut.sh" "/tmp/world-iter334-executor.desvqzr3/mutants" H1 plain 3` | 1 | 0 / 0 / 0 | 0 / 0 | 0.071 |
| M1 [H1-race-protocol](executor-H1-race-protocol.txt) | `"$WX/tools/runmut.sh" "/tmp/world-iter334-executor.desvqzr3/mutants" H1 race 3` | 0 | 0 / 3 / 0 | 0 / 3 | 2.922 |
| M1 [H2-plain-protocol](executor-H2-plain-protocol.txt) | `"$WX/tools/runmut.sh" "/tmp/world-iter334-executor.desvqzr3/mutants" H2 plain 3` | 1 | 0 / 0 / 0 | 0 / 0 | 0.051 |
| M1 [H2-race-protocol](executor-H2-race-protocol.txt) | `"$WX/tools/runmut.sh" "/tmp/world-iter334-executor.desvqzr3/mutants" H2 race 3` | 0 | 0 / 3 / 0 | 0 / 3 | 2.783 |
| M1 [H3-plain-protocol](executor-H3-plain-protocol.txt) | `"$WX/tools/runmut.sh" "/tmp/world-iter334-executor.desvqzr3/mutants" H3 plain 3` | 1 | 0 / 0 / 0 | 0 / 0 | 0.052 |
| M1 [H3-race-protocol](executor-H3-race-protocol.txt) | `"$WX/tools/runmut.sh" "/tmp/world-iter334-executor.desvqzr3/mutants" H3 race 3` | 0 | 0 / 3 / 0 | 0 / 3 | 17.823 |
| M1 [M6-plain-protocol](executor-M6-plain-protocol.txt) | `"$WX/tools/runmut.sh" "/tmp/world-iter334-executor.desvqzr3/mutants" M6 plain 20` | 1 | 0 / 0 / 0 | 0 / 0 | 0.055 |
| M1 [M6-race-protocol](executor-M6-race-protocol.txt) | `"$WX/tools/runmut.sh" "/tmp/world-iter334-executor.desvqzr3/mutants" M6 race 20` | 0 | 0 / 60 / 0 | 0 / 20 | 11.691 |
| M1 [M7-plain-protocol](executor-M7-plain-protocol.txt) | `"$WX/tools/runmut.sh" "/tmp/world-iter334-executor.desvqzr3/mutants" M7 plain 20` | 1 | 0 / 0 / 0 | 0 / 0 | 0.051 |
| M1 [M7-race-protocol](executor-M7-race-protocol.txt) | `"$WX/tools/runmut.sh" "/tmp/world-iter334-executor.desvqzr3/mutants" M7 race 20` | 0 | 0 / 60 / 0 | 0 / 20 | 14.047 |
| M1 [M7b-plain-protocol](executor-M7b-plain-protocol.txt) | `"$WX/tools/runmut.sh" "/tmp/world-iter334-executor.desvqzr3/mutants" M7b plain 20` | 1 | 0 / 0 / 0 | 0 / 0 | 0.052 |
| M1 [M7b-race-protocol](executor-M7b-race-protocol.txt) | `"$WX/tools/runmut.sh" "/tmp/world-iter334-executor.desvqzr3/mutants" M7b race 20` | 0 | 60 / 0 / 0 | 20 / 0 | 22.679 |
| M1 [scope-m1](executor-scope-m1.txt) | `shasum -a 256 -c "$WX/fixed.sha"<br>git diff --name-only<br>git diff --stat` | 0 | 0 / 0 / 0 | 0 / 0 | 0.049 |
| M1 [wall-time](executor-wall-time.txt) | `go test -timeout 600s ./host/daemon -run '^TestHumanApproveIsSerializedAcrossEpisodes$' -count=1 -v` | 0 | 3 / 0 / 0 | 1 / 0 | 1.487 |
| M1 [bash-version](executor-bash-version.txt) | `/bin/bash --version` | 0 | 0 / 0 / 0 | 0 / 0 | 0.009 |
| M1 [script-hashes](executor-script-hashes.txt) | `shasum -a 256 "$WX/tools/runmut.sh" "$WX/tools/loaded-soak.sh"` | 0 | 0 / 0 / 0 | 0 / 0 | 0.029 |
| M2 [soak-idle](executor-soak-idle.txt) | `go test -race -timeout 900s ./host/daemon -run '^TestHumanApproveIsSerializedAcrossEpisodes$' -count=200` | 0 | 0 / 0 / 0 | 0 / 0 | 189.297 |
| M2 [soak-loaded](executor-soak-loaded.txt) | `"$WX/tools/loaded-soak.sh" "$PWD"` | 7 | 0 / 0 / 0 | 0 / 0 | 0.298 |
| M2 [package-vet](executor-package-vet.txt) | `perl -e 'alarm 300; exec @ARGV' go vet ./host/daemon/` | 0 | 0 / 0 / 0 | 0 / 0 | 0.454 |
| M2 [package-plain](executor-package-plain.txt) | `go test -timeout 900s ./host/daemon/ -count=1` | 1 | 0 / 4 / 0 | 0 / 2 | 9.306 |
| M2 [package-race](executor-package-race.txt) | `go test -race -timeout 1500s ./host/daemon/ -count=1` | 1 | 0 / 4 / 0 | 0 / 2 | 9.981 |
| M2 [head-provenance](executor-head-provenance.txt) | `git diff --name-only cf79a7a92199eb65ddda4bf724167ba9c8c092e1 HEAD` | 0 | 0 / 0 / 0 | 0 / 0 | 0.016 |
| M2 [full-gate](executor-full-gate.txt) | `perl -e 'alarm 2400; exec @ARGV' bash <<'FULL'<br>go build ./... && go test -timeout 1800s ./...<br>FULL` | 1 | 0 / 44 / 0 | 0 / 23 | 190.65 |
| M2 [verify-ail](executor-verify-ail.txt) | `(<br>    unset AILANG_BIN<br>    export PATH=$HOME/.ailang/state/ailang-release-v0.55.0/x:$PATH<br>    export WORLD_PKG_AILANG_BIN=$HOME/.ailang/state/ailang-release-v0.41.0/x/ailang<br>    perl -e 'alarm 600; exec @ARGV' ./scripts/verify_ail.sh<br>)` | 0 | 0 / 0 / 0 | 0 / 0 | 11.146 |
| M2 [base-compare-package-plain](executor-base-compare-package-plain.txt) | `go test -timeout 900s ./host/daemon/ -count=1` | 1 | 0 / 4 / 0 | 0 / 2 | 10.082 |
| M2 [base-compare-package-race](executor-base-compare-package-race.txt) | `go test -race -timeout 1500s ./host/daemon/ -count=1` | 1 | 0 / 4 / 0 | 0 / 2 | 12.736 |
| M2 [final-pin-check](executor-final-pin-check.txt) | `shasum -a 256 -c "$WX/fixed.sha"` | 0 | 0 / 0 / 0 | 0 / 0 | 0.026 |
| M2 [process-enumeration](executor-process-enumeration.txt) | `scan_rc=0<br>pgrep -x yes > "$WX/logs/final-yes.txt" 2> "$WX/logs/final-yes-error.txt" \|\| scan_rc=$?<br>printf 'pgrep -x yes rc=%s\n' "$scan_rc"<br>if [ "$scan_rc" -le 1 ]; then<br>    cat "$WX/logs/final-yes.txt"<br>    wc -l < "$WX/logs/final-yes.txt"<br>else<br>    cat "$WX/logs/final-yes-error.txt"<br>    echo 'UNINFORMATIVE UNDER SANDBOX: process count unavailable; no burners started by executor'<br>fi` | 0 | 0 / 0 / 0 | 0 / 0 | 0.006 |
| M2 [base-compare-full-gate](executor-base-compare-full-gate.txt) | `perl -e 'alarm 2400; exec @ARGV' bash <<'FULL'<br>go build ./... && go test -timeout 1800s ./...<br>FULL` | 1 | 0 / 45 / 0 | 0 / 24 | 184.903 |
| M2 [final-scope](executor-final-scope.txt) | `set -e<br>shasum -a 256 -c "$WX/fixed.sha"<br>git diff --name-only<br>git diff --stat` | 0 | 0 / 0 / 0 | 0 / 0 | 0.058 |
| M2 [git-status](executor-git-status.txt) | `git status --porcelain` | 0 | 0 / 0 / 0 | 0 / 0 | 0.027 |

## Mutation results (measured, one mutant at a time)

| ID | Mode / requested count | Protocol rc / Go test rc | Verdict | Top PASS / FAIL | Subtest PASS / FAIL | Sole killer measured from raw log |
|---|---|---|---|---|---|---|
| A | plain / 20 | 1 / not started | INSTRUMENT ERROR | 0 / 0 | 0 / 0 | none measured; no test launched |
| A | race / 20 | 0 / 1 | KILLED | 0 / 20 | 0 / 0 | :270 unwrapped probe max = 1, want 2 |
| M1 | plain / 20 | 1 / not started | INSTRUMENT ERROR | 0 / 0 | 0 / 0 | none measured; no test launched |
| M1 | race / 20 | 0 / 1 | KILLED | 0 / 20 | 0 / 40 | :372 serialized probe ran 2 of 2 executions without the daemon mutex held |
| M1b | plain / 20 | 1 / not started | INSTRUMENT ERROR | 0 / 0 | 0 / 0 | none measured; no test launched |
| M1b | race / 20 | 0 / 1 | KILLED | 0 / 20 | 0 / 40 | :372 serialized probe ran 2 of 2 executions without the daemon mutex held |
| M2 | plain / 20 | 1 / not started | INSTRUMENT ERROR | 0 / 0 | 0 / 0 | none measured; no test launched |
| M2 | race / 20 | 0 / 1 | KILLED | 0 / 20 | 0 / 40 | :372 serialized probe ran 2 of 2 executions without the daemon mutex held |
| M3 | plain / 20 | 1 / not started | INSTRUMENT ERROR | 0 / 0 | 0 / 0 | none measured; no test launched |
| M3 | race / 20 | 0 / 1 | KILLED | 0 / 20 | 20 / 20 | :392 unwrapped probe max = 1, want 2 |
| M6 | plain / 20 | 1 / not started | INSTRUMENT ERROR | 0 / 0 | 0 / 0 | none measured; no test launched |
| M6 | race / 20 | 0 / 1 | KILLED | 0 / 20 | 0 / 40 | :375 serialized probe max = 2, want 1 |
| M7 | plain / 20 | 1 / not started | INSTRUMENT ERROR | 0 / 0 | 0 / 0 | none measured; no test launched |
| M7 | race / 20 | 0 / 1 | KILLED | 0 / 20 | 0 / 40 | :372 serialized probe ran 2 of 2 executions without the daemon mutex held |
| M7b | plain / 20 | 1 / not started | INSTRUMENT ERROR | 0 / 0 | 0 / 0 | none measured; no test launched |
| M7b | race / 20 | 0 / 0 | SURVIVES (expected PASS boundary) | 20 / 0 | 40 / 0 | none |
| H1 | plain / 3 | 1 / not started | INSTRUMENT ERROR | 0 / 0 | 0 / 0 | none measured; no test launched |
| H1 | race / 3 | 0 / 1 | KILLED | 0 / 3 | 0 / 0 | :426 missed = "", want "second Human.Approve dispatch did not arrive within 10s" |
| H2 | plain / 3 | 1 / not started | INSTRUMENT ERROR | 0 / 0 | 0 / 0 | none measured; no test launched |
| H2 | race / 3 | 0 / 1 | KILLED | 0 / 3 | 0 / 0 | :426 missed = "second Human.Approve dispatch did not arrive", want "second Human.Approve dispatch did not arrive within 10s" |
| H3 | plain / 3 | 1 / not started | INSTRUMENT ERROR | 0 / 0 | 0 / 0 | none measured; no test launched |
| H3 | race / 3 | 0 / 1 | KILLED | 0 / 3 | 0 / 0 | :422 TestConcurrencyProbeNamesAMissedArrival: Execute did not return after the injected timeout within 5s |

The sole-killer counts are checked by the unchanged pinned runner in every race invocation. Logs contain no build error, race finding, Go timeout, SKIP, or socket error in those mutation runs. M1/M1b/M2/M7 fail in both subtests at unheld before max. M6 is empirical. M7b consistently rewires the test mutex and is expected to survive; it makes no claim about production binding. H3 fails by its own 5 s bound, not the Go timeout. M5 and historical M6b are not required drill rows.

### Exact substitutions

A:

```sh
perl -0pi -e 's/(\t\t\t\t_, errs\[i\] = c\.Dispatch\(ctx, calls\[i\]\))/\t\t\t\tif i == 1 {\n\t\t\t\t\ttime.Sleep(300 * time.Millisecond)\n\t\t\t\t}\n$1/' "$T"
```

M1:

```sh
perl -0pi -e 's/\ts\.mu\.Lock\(\)         \/\/ MD3 target\n\tdefer s\.mu\.Unlock\(\) \/\/ MD3 target\n//' "$G"
```

M1b:

```sh
perl -0pi -e 's/defer s\.mu\.Unlock\(\) \/\/ MD3 target/s.mu.Unlock() \/\/ MD3 target/' "$G"
```

M2:

```sh
perl -0pi -e 's/gate\.h = serializedHandler\{mu: sh\.mu, h: probe\}/gate.h = probe/' "$T"
```

M3:

```sh
perl -0pi -e 's/\t\tp\.awaitSecond\(ctx\)\n/\t\tselect {\n\t\tcase <-p.inside2:\n\t\tcase <-time.After(200 * time.Millisecond):\n\t\tcase <-ctx.Done():\n\t\t}\n/' "$T"
```

M6:

```sh
perl -0pi -e 's/\ts\.mu\.Lock\(\)         \/\/ MD3 target\n\tdefer s\.mu\.Unlock\(\) \/\/ MD3 target\n//' "$G"
perl -0pi -e 's/unheld := p\.held != nil && p\.held\.TryLock\(\)/unheld := false \&\& p.held != nil \&\& p.held.TryLock()/' "$T"
```

M7:

```sh
perl -0pi -e 's/func \(s serializedHandler\) Execute/var otherMu sync.Mutex\n\nfunc (s serializedHandler) Execute/; s/\ts\.mu\.Lock\(\)         \/\/ MD3 target\n\tdefer s\.mu\.Unlock\(\) \/\/ MD3 target/\totherMu.Lock()         \/\/ MD3 target\n\tdefer otherMu.Unlock() \/\/ MD3 target/' "$G"
```

M7b:

```sh
perl -0pi -e 's/\tsh := d\.human\.\(serializedHandler\)\n/\tsh := d.human.(serializedHandler)\n\t_ = sh\n\totherMu := \&sync.Mutex{}\n/; s/held: sh\.mu/held: otherMu/; s/gate\.h = serializedHandler\{mu: sh\.mu, h: probe\}/gate.h = serializedHandler{mu: otherMu, h: probe}/' "$T"
```

H1:

```sh
perl -0pi -e 's/(case <-timeout:\n\t\t)why = fmt\.Sprintf\("second Human\.Approve dispatch did not arrive within %s", probeWaitBound\)/$1why = ""/' "$T"
```

H2:

```sh
perl -0pi -e 's/(case <-timeout:\n\t\t)why = fmt\.Sprintf\("second Human\.Approve dispatch did not arrive within %s", probeWaitBound\)/$1why = "second Human.Approve dispatch did not arrive"/' "$T"
```

H3:

```sh
perl -0pi -e 's/timeout := p\.timeout/var timeout <-chan time.Time/' "$T"
```

### Full before / edited / restored SHA256 evidence

Both files are included for every invocation, including unchanged nontargeted files. Each edited digest was checked against the JSON (or the unchanged pin for nontargets), every restored digest equals its before digest, and both archive files equal their base/fixed pins after each invocation.

| ID / mode | File | Before SHA256 | Edited SHA256 | Restored SHA256 |
|---|---|---|---|---|
| A / plain | `host/daemon/decisionpacket_daemon_test.go` | `53a708a5264fafe53217d59c58dafa7e25eeabe7f5ef9debc554912fa30a0445` | `2c4d47bc3e89e7b488e7998d1708173cbfef6b83573eb285dc79ea44ed2f8722` | `53a708a5264fafe53217d59c58dafa7e25eeabe7f5ef9debc554912fa30a0445` |
| A / plain | `host/daemon/daemon.go` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` |
| A / race | `host/daemon/decisionpacket_daemon_test.go` | `53a708a5264fafe53217d59c58dafa7e25eeabe7f5ef9debc554912fa30a0445` | `2c4d47bc3e89e7b488e7998d1708173cbfef6b83573eb285dc79ea44ed2f8722` | `53a708a5264fafe53217d59c58dafa7e25eeabe7f5ef9debc554912fa30a0445` |
| A / race | `host/daemon/daemon.go` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` |
| M1 / plain | `host/daemon/decisionpacket_daemon_test.go` | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` |
| M1 / plain | `host/daemon/daemon.go` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` | `3737b1ed8187fdb92efee6ba71f36761068739736a226104324ffc73b92a1be3` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` |
| M1 / race | `host/daemon/decisionpacket_daemon_test.go` | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` |
| M1 / race | `host/daemon/daemon.go` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` | `3737b1ed8187fdb92efee6ba71f36761068739736a226104324ffc73b92a1be3` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` |
| M1b / plain | `host/daemon/decisionpacket_daemon_test.go` | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` |
| M1b / plain | `host/daemon/daemon.go` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` | `66076a39d936863f431cc74e814d14395e59e8c3869656129432383c8b4b5ffd` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` |
| M1b / race | `host/daemon/decisionpacket_daemon_test.go` | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` |
| M1b / race | `host/daemon/daemon.go` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` | `66076a39d936863f431cc74e814d14395e59e8c3869656129432383c8b4b5ffd` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` |
| M2 / plain | `host/daemon/decisionpacket_daemon_test.go` | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` | `3dc59ab26e2a7e56d89a175602d84e3bc1f043508e275da352dd4c7ce4d1f041` | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` |
| M2 / plain | `host/daemon/daemon.go` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` |
| M2 / race | `host/daemon/decisionpacket_daemon_test.go` | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` | `3dc59ab26e2a7e56d89a175602d84e3bc1f043508e275da352dd4c7ce4d1f041` | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` |
| M2 / race | `host/daemon/daemon.go` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` |
| M3 / plain | `host/daemon/decisionpacket_daemon_test.go` | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` | `092a2b0bb9879673d38ae2870e84299836ef50bb23c1fe35881e2d9ce40b9c80` | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` |
| M3 / plain | `host/daemon/daemon.go` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` |
| M3 / race | `host/daemon/decisionpacket_daemon_test.go` | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` | `092a2b0bb9879673d38ae2870e84299836ef50bb23c1fe35881e2d9ce40b9c80` | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` |
| M3 / race | `host/daemon/daemon.go` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` |
| M6 / plain | `host/daemon/decisionpacket_daemon_test.go` | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` | `0914d8a6d150d48a11ddcbce1eeb91132453f0cc8a3a94774d6a10a9f8336b2a` | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` |
| M6 / plain | `host/daemon/daemon.go` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` | `3737b1ed8187fdb92efee6ba71f36761068739736a226104324ffc73b92a1be3` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` |
| M6 / race | `host/daemon/decisionpacket_daemon_test.go` | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` | `0914d8a6d150d48a11ddcbce1eeb91132453f0cc8a3a94774d6a10a9f8336b2a` | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` |
| M6 / race | `host/daemon/daemon.go` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` | `3737b1ed8187fdb92efee6ba71f36761068739736a226104324ffc73b92a1be3` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` |
| M7 / plain | `host/daemon/decisionpacket_daemon_test.go` | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` |
| M7 / plain | `host/daemon/daemon.go` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` | `682bf08e8e8d240835fcb96626d70d505a40a8299172d82252eda5afee5f059b` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` |
| M7 / race | `host/daemon/decisionpacket_daemon_test.go` | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` |
| M7 / race | `host/daemon/daemon.go` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` | `682bf08e8e8d240835fcb96626d70d505a40a8299172d82252eda5afee5f059b` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` |
| M7b / plain | `host/daemon/decisionpacket_daemon_test.go` | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` | `2311eafcf5b56b71460006d2b5de9c81ee5347b1944895bc4de73b7da406ee8b` | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` |
| M7b / plain | `host/daemon/daemon.go` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` |
| M7b / race | `host/daemon/decisionpacket_daemon_test.go` | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` | `2311eafcf5b56b71460006d2b5de9c81ee5347b1944895bc4de73b7da406ee8b` | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` |
| M7b / race | `host/daemon/daemon.go` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` |
| H1 / plain | `host/daemon/decisionpacket_daemon_test.go` | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` | `5c836e58172ca1896728af59b371d347a4547e1626af42198f30fd975d3ef7ff` | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` |
| H1 / plain | `host/daemon/daemon.go` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` |
| H1 / race | `host/daemon/decisionpacket_daemon_test.go` | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` | `5c836e58172ca1896728af59b371d347a4547e1626af42198f30fd975d3ef7ff` | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` |
| H1 / race | `host/daemon/daemon.go` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` |
| H2 / plain | `host/daemon/decisionpacket_daemon_test.go` | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` | `d16c35e5f270f12c36c52b23763742efe73e393e6e930292f661fd4f13b98153` | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` |
| H2 / plain | `host/daemon/daemon.go` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` |
| H2 / race | `host/daemon/decisionpacket_daemon_test.go` | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` | `d16c35e5f270f12c36c52b23763742efe73e393e6e930292f661fd4f13b98153` | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` |
| H2 / race | `host/daemon/daemon.go` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` |
| H3 / plain | `host/daemon/decisionpacket_daemon_test.go` | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` | `7f98f9f8e890b7da6c9e00dc03346bf1d72593505b0fd1e81d14142632b0b604` | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` |
| H3 / plain | `host/daemon/daemon.go` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` |
| H3 / race | `host/daemon/decisionpacket_daemon_test.go` | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` | `7f98f9f8e890b7da6c9e00dc03346bf1d72593505b0fd1e81d14142632b0b604` | `877013b89be5344b39402b07dd498dbf404ffe2f9e5640c4853f355f3a450383` |
| H3 / race | `host/daemon/daemon.go` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` | `e6b70d7e3f3fb20ef6a89d52b82a87bbdff0a05915c6afaab8956cb05be6c668` |

## Findings, soak, baseline comparisons, and limits

The execution is not an unqualified green result. M1 is unverified because all 11 plain mutation Go commands could not launch. M2 is unverified because the loaded soak could not launch, socket/TTY gates are sandbox-uninformative, other full-gate findings remain red at base, and process enumeration is unavailable.

### Soak and load

- Idle: prescribed `-race -timeout 900s -count=200` command rc 0; wrapper wall 189.297 s. Package output is `ok`; 0 FAIL markers and no DATA RACE. The command is non-verbose: 200 is the requested repetition count, and named PASS marker count is 0, not a separately measured 200 named PASS lines.
- Loaded: approved script rc 7 at its enumeration preflight, before starting burners or Go. **UNINFORMATIVE UNDER SANDBOX**. Exactly 0 burners were started by this executor; no load processes were created or left by it.
- Final `pgrep -x yes` rc 3: `sysmond service not found` / `Cannot get process list`. Stdout contains no PIDs, but a global count of zero is NOT established. No piped `wc` count is claimed as cleanup evidence. See executor-process-enumeration.txt, executor-final-yes.txt and executor-final-yes-error.txt.

### Package baseline

Fixed package plain rc 1 (9.306 s), race rc 1 (9.981 s); untouched base plain rc 1 (10.082 s), race rc 1 (12.736 s). Each emits 2 top-level FAIL and 2 subtest FAIL markers: `TestAGUIShutdownEndsRuns` (false/true) cannot bind `127.0.0.1:0` at agui_test.go:583; `TestAGUISeesDirectStoreCommit` panics on httptest listen `[::1]:0` at agui_test.go:684. **UNINFORMATIVE UNDER SANDBOX**, neither acceptance pass nor fail. The named sprint test is not claimed discovered in these early-aborted package runs.

### Full baseline

Fixed full command rc 1 in 190.650 s; 23 top-level and 21 subtest FAIL markers (44 total), 0 PASS/SKIP markers (non-verbose). Base exact command, same binaries and environment, rc 1 in 184.903 s; 24 top-level and 21 subtest FAIL markers (45 total). No build failure: the && reached go test in both commands.

Package summaries:

| Tree | ok packages | FAIL packages | no-test packages |
|---|---:|---:|---:|
| Fixed | 27 | 8 | 2 |
| Base archive | 26 | 9 | 2 |

The failing fixed packages are cmd/ailang-worldd, cmd/world-publish, host/broker, host/capsule, host/daemon, host/pinfetch, host/ptytest, and scripts/marketingcheck. Loopback bind denials, `/dev/tty` denial, confinement decoy write denial and dependent assertions are **UNINFORMATIVE UNDER SANDBOX**. Their raw nonzero rc and failure markers remain preserved.

Other findings are reported without working around or attributing them to the fix: `TestOrdinaryOverflowCarriesNoESRCH` cannot observe an unreaped zombie in both broker and capsule; `TestExecHandlerTimeoutKillsTheGroup` fails at its specific-child-readiness observation; `TestTheTwoShapesOfControllingTerminal` does not observe the PTY prompt within 10 s. All of these are also red at base. The logs do not establish an unsandboxed outcome for these findings.

Base-only failure: `TestReleaseLiteralsExplainedByManifest` in host/pins, scan_test.go:256, `git ls-files: exit status 128`. This is the documented missing-.git archive artifact; fixed worktree host/pins passed. No git repository was initialized and no workaround was attempted. Every named top-level fixed full failure is also present in the base run; this is a baseline comparison, not a claim that an unsandboxed gate passes.

| Measured top-level failing test | Fixed | Base |
|---|---|---|
| `TestAGUISeesDirectStoreCommit` | red | red |
| `TestAGUIShutdownEndsRuns` | red | red |
| `TestAttendedPublishMintsThroughTheLandedTraversalAndSpendsExactlyOnce` | red | red |
| `TestCLIRealSubprocessEpisode` | red | red |
| `TestClientDeadlineAgainstAcceptingServer` | red | red |
| `TestDemoProofHTTPProducerRoundTrip` | red | red |
| `TestExecHandlerTimeoutKillsTheGroup` | red | red |
| `TestExecProbeForeignListenerOnTheReversePortDoesNotRefuse` | red | red |
| `TestExecProbePassThroughShimFailsEveryConfinementArm` | red | red |
| `TestExecProbeReverseDialErrIsSilent` | red | red |
| `TestExecProbeReverseLegChattyForeignHolderIsRetried` | red | red |
| `TestExecProbeReverseLegExhaustionRefusesDistinctly` | red | red |
| `TestExecProbeReverseLegNoTokenIsNotRetried` | red | red |
| `TestExecProbeReverseLegRealReachSurvivesAForeignAttempt` | red | red |
| `TestExecProbeReverseLegSilentHolderRefusesSticky` | red | red |
| `TestExecProbeReversePollRecordsConnectReset` | red | red |
| `TestExecSrtConfinementMatrix` | red | red |
| `TestInstallGoodRelease` | red | red |
| `TestMarketingCapturePassFakeChrome` | red | red |
| `TestOrdinaryOverflowCarriesNoESRCH` | red | red |
| `TestReleaseLiteralsExplainedByManifest` | not red | red |
| `TestTheTwoShapesOfControllingTerminal` | red | red |
| `TestTransitionsReadsTheConfirmationFromThePtyWhenStdinIsDevNull` | red | red |

### AILANG gate

`verify_ail.sh` rc 0 in 11.146 s under the exact prescribed release environment: 16 required identities verified, 40 named tests pass, 503 se-tools named tests pass, world package 9/9 steps performed nonzero work. The world-package output reports exactly 9 named no-generator skips, which are the documented cross-file derive gap; there are no skips caused by missing binaries. `UNRECOGNISED RELEASE` for legs 1–2 v0.55.0 is the prescribed environment’s informational warning; package compiler bytes are pinned to v0.41.0. No binary or manifest was changed.

### Not reached and controller handoff

Every planned acceptance command was attempted. The actual Go tests within all plain mutation invocations were not reached because the unchanged runner errored first; the loaded soak Go test and load creation were not reached because the unchanged preflight refused to proceed. No acceptance item was omitted for time. M5/M6b are explicitly not required. Controller-only commit/PR/CI, all three green CI jobs on the fix head, and unsandboxed reruns remain outstanding. All waits were bounded; execution and final evidence collection completed within the 30-minute cap.


## Deviations and executor instrument findings

Initial orchestration invocation exited rc=1 after A-plain returned rc=1. The evidence helper then raised KeyError: mutations because the JSON key is mutation_commands. Corrected only the scratch evidence helper key; did not rerun or overwrite A-plain. All before/edited/restored digests for that invocation were subsequently verified against the JSON.
Approved runmut.sh is byte-identical to its extracted pin. On /bin/bash 3.2.57, an empty flags array under set -u errors at line 93: flags[@]: unbound variable. No go test is started in any plain mutation invocation. These are instrument errors, never mutation kills or survivals. Race flags are nonempty and race invocations run without this shell error.
Deviation: independent mutation invocations continue after recording and inspecting the known plain-run instrument error, rather than abandoning all later independent gates at the first || exit. No script, test, binary, timeout, or assertion is changed to work around the error.
During report regeneration, a scratch-only quoting syntax error interrupted the report helper. Corrected the report string and regenerated successfully; no acceptance command or mutation was rerun or modified.


Fenced blocks were extracted with Python rather than the planned awk helper; the resulting diff and both scripts matched every pinned SHA256 exactly before use. The acceptance Go commands and runner contents were unchanged. No approved code, script, substitution, binary, timeout, stagger, or assertion was adjusted. The scratch Python evidence helper was used only to invoke exact acceptance commands, capture outputs/rc/timings, compare pins, and produce this report. Long commands retain the plan bounds and every go test has -timeout. Package-scale gates run sequentially.

## Final scope and handoff

Final pin and scope checks are recorded in executor-final-scope.txt and exact `git status --porcelain` in executor-git-status.txt. No git writes, mission edits, AILANG edits, CI edits, or production changes were performed. The controller owns commits, PR, CI and reruns outside the sandbox.

Report generated at 2026-11-02T13:38:55.375707+00:00; total elapsed since 13:21:53 UTC: 1022.4 s (30-minute cap).
