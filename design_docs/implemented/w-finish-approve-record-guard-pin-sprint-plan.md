# w-finish-approve-record-guard-pin — sprint plan (iteration 335, row 237)

Authority: [design r1](w-finish-approve-record-guard-pin.md), quorum r1 PROCEED 3/3
(oc-kimi-k3, gpt6-1-sol, gemini-4-argon); [coding standards](../coding-standards.md), S6.
Estimate: **0.1 day; one milestone; one new 67-line test file; 0 production lines**.
The controller commits the artifacts and final change. This plan is ready for execution;
planner validation is limited to artifact consistency, formatting and scratch mutation mechanics,
not a claim that the final-file Go drills or gates have run.

## Base and scope

Planner worktree: `/Users/voightkampff/dev/sunholo-data/.planner-wt-iter335` (detached). HEAD:
`936edb59ee3b3a022a6f1bdc1e9565c23c619dfe`.

| Observed file | SHA-256 at base | Lines |
|---|---|---:|
| `host/daemon/episode_finish.go` | `4e320c537e689620480af2bdff9bdcd9d4d6febba43b4a5dd5b3bb09971973bf` | 280 |
| `host/daemon/episode_finish_test.go` | `b872a38cc091e250e570d1bba16afa17044c57ff1d9a69cc63484400429b9d67` | 258 |
| Prototype evidence `.go.txt` | `f5df08b346dd833250ea444ed15427fb4cec1fe818f05273e0fde5187134d2e5` | 61 |

Only implementation file: `host/daemon/episode_finish_approve_guard_test.go` (new). Reuse existing helpers without changing them.
Executor evidence, scratch files and logs may go only under `verification/world-iter335/executor/`;
these are not additional product changes. Preserve designer evidence. No `.ail`, production,
CLI, dependency, CI or shared-driver changes; this route test belongs in the Go host, not an AILANG package.

No git writes: no add, commit, stash, checkout, reset, branch, restore or git apply.
Do not edit `design_docs/world-mission.md`, `design_docs/world-mission-log.md`,
`tools/launchd/`, or the authority. No step writes outside the worktree, including temporary
files, caches, shared logs, `$HOME` or `~/.ailang/state`. External inputs may be read.
Establish every unchanged claim by observing hashes/line counts/diffs; never seed a write to
an unchanged target. Transient mutations below are explicit experiments, restored byte for byte.

Controller's outside-sandbox base observations at this HEAD: vet passes, gofmt output empty,
Finish suite passes **6** top-level tests; M1 survives those six and the prototype kills it.
A gate red at base is a finding to report, never something to work around or quietly repair.

## Quorum decisions and premise correction

1. **oc-kimi-k3:** keep the refusal first, on the same fresh store as the later control.
   Final lines 42–46 state the content-addressing reason, prohibit parallel subtests, and
   observe absence of the ApprovalsV1 head and both approval object kinds **before posting**.
   The three-line assertion is cheap and gives a reordered control a named precondition failure.
   Keep this assertion together with arm 1; no `t.Parallel()` or parallel subtest refactor.
2. **gpt6-1-sol:** six added lines change the prototype's bytes. Re-run **every drill** below
   (M1, M3a, M3c, M2+M3c, M2b, M4+M3a, plus clean M2 and M1-existing controls) against the
   **final file**, and record its full SHA-256 before/after. Prototype evidence cannot discharge
   these checks. Any later byte change requires a new hash and the complete drill set again.
3. **gemini-4-argon:** use the new file named by the design. It gives the entire change an
   independent hash, keeps the 258-line helper/host file untouched, and has no behavior cost.

The design's statement that the prototype is already at the host `_test.go` path is false at
this planner base: it exists only at
`verification/world-iter335/designer/prototype_episode_finish_approve_guard_test.go.txt`.
Its measurement HEAD `77e0dd7c` is historical; use the base above. No behavioral premise was
found false by source/evidence inspection. Six existing Finish tests, rather than five, is
already corrected in design §7. Do not run the designer's `drills.sh`: it uses git checkout.

## M1 — Add the guard pin and verify its three assertions (0.1 day)

1. Observe base HEAD/hashes/line counts, gate evidence and name availability. If base differs,
   report drift before applying. Check no existing `approvalWrites` or test function collision.
2. Create the final file from the exact Go block below (UTF-8, LF, trailing newline, no fence).
   Check `gofmt` output, 67 lines and the pinned hash; never rewrite helpers or production.
3. Run the mutation script against these exact bytes. Inspect every return code **and named
   assertion**; compile errors, timeouts and infrastructure errors are not mutation kills.
   Restore both files after each drill and on early exit; require the full base production hash.
4. Run clean gates, observe both existing files unchanged, and write the executor report with
   final hash/line count, raw logs, AC mapping and limitations. Controller runs package/full
   gates and records CI status for the fix head. Do not invent a CI result before it exists.

| Design AC | Acceptance in M1 | Evidence |
|---|---|---|
| AC-1 | M1 fails only at final :51, `approve with record: 200`; existing six survive | M1, M1-existing logs |
| AC-2 | M3a writes then returns correct refusal; fails at :54, `refused approve recorded an approval: head=true` | M3a log; M4+M3a survives after deleting only :53–55 |
| AC-3 | M3c always refuses; fails at :61, `control approve without record: 403` | M3c log; M2-clean and M2+M3c survive; M2b fails at :61 with `200` |
| AC-4 | Clean new test passes once, 20 times and race 5 times; Finish has 7 named top-level PASS; vet/gofmt clean | Clean gate logs, full stdout and exit codes |
| AC-5 | Controller's daemon/broker packages, full Go build/tests and verify_ail pass with released v0.41.0; CI green on fix head | Outside-sandbox raw logs + CI reference, or explicit pending/finding |
| AC-6 | Production SHA remains `4e320c537e689620480af2bdff9bdcd9d4d6febba43b4a5dd5b3bb09971973bf` / 280 lines; existing host test SHA remains `b872a38cc091e250e570d1bba16afa17044c57ff1d9a69cc63484400429b9d67` / 258 lines; only new Go test is added | Before/after observations, clean existing-file diff, status |

S6 is discharged by M1/M3a/M3c's distinct failing assertions and the M2/M4 deletion
controls. Counting assertion text cannot substitute for those behavioral drills.

## Exact final test

Path: `host/daemon/episode_finish_approve_guard_test.go`. **SHA-256 `c9e279de7c7ce87edf5fbc0e43e690253616d320b5585625062cda3d29cdec57`; 67 lines.**
The fenced block is the full file, including a newline after its final brace.

```go
package daemon

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/broker"
)

// approvalWrites reports what broker.RecordOperatorApproval leaves behind: the
// approvals head and the request and decision objects it puts.
func approvalWrites(t *testing.T, d *Daemon) (head bool, requests, decisions []string) {
	t.Helper()
	ctx := boundedTestContext(t)
	_, head, err := d.store.GetRegistryHead(ctx, broker.ApprovalsV1)
	if err != nil {
		t.Fatal(err)
	}
	for id, out := range map[string]*[]string{broker.ApprovalRequestV1: &requests, broker.ApprovalDecisionV1: &decisions} {
		objs, err := d.store.ObjectsBySemanticID(ctx, id, "", 10)
		if err != nil {
			t.Fatal(err)
		}
		for _, o := range objs {
			*out = append(*out, o.Hash.String())
		}
	}
	return head, requests, decisions
}

// TestEpisodeFinishApproveRefusesRecord pins the approve guard's
// `|| input.Record != nil` arm (row 237): an approve request that also carries a
// record is refused and records no approval. The control is the identical
// request without the record, on the same clock and scope, which is approved.
// Each arm mints a fresh proof: the route consumes it before the guard runs.
func TestEpisodeFinishApproveRefusesRecord(t *testing.T) {
	d := newHandlerDaemon(t)
	pinFinishClock(d, finishClockAt)
	scope := finishScopeAt(finishClockAt + 900)

	// Arm 1 must run first on a pristine store: under M1/M3a both arms write
	// byte-identical approval objects. Do not split these into parallel subtests.
	if head, requests, decisions := approvalWrites(t, d); head || len(requests) != 0 || len(decisions) != 0 {
		t.Fatalf("arm 1 precondition: pristine approval store required: head=%v requests=%v decisions=%v", head, requests, decisions)
	}

	proof := finishProofFile(t, d, "approve", finishClockAt+100, 0o600)
	w := finishRoute(t, d, FinishRequest{Phase: "approve", Scope: scope, Record: &EpisodeFinish{Episode: "ep1"}}, proof, "")
	if w.Code != 403 || !strings.Contains(w.Body.String(), "operator route: invalid or expired finish scope") {
		t.Fatalf("approve with record: %d %s", w.Code, w.Body.String())
	}
	if head, requests, decisions := approvalWrites(t, d); head || len(requests) != 0 || len(decisions) != 0 {
		t.Fatalf("refused approve recorded an approval: head=%v requests=%v decisions=%v", head, requests, decisions)
	}

	proof = finishProofFile(t, d, "approve", finishClockAt+100, 0o600)
	w = finishRoute(t, d, FinishRequest{Phase: "approve", Scope: scope}, proof, "")
	var approval FinishApproval
	if w.Code != 200 || json.Unmarshal(w.Body.Bytes(), &approval) != nil {
		t.Fatalf("control approve without record: %d %s", w.Code, w.Body.String())
	}
	head, requests, decisions := approvalWrites(t, d)
	if !head || len(requests) != 1 || requests[0] != approval.Request || len(decisions) != 1 || decisions[0] != approval.Approval {
		t.Fatalf("control recorded head=%v requests=%v decisions=%v, want %+v", head, requests, decisions, approval)
	}
}
```

## Execution environment and evidence policy

All pinned shell blocks target **`/bin/bash` 3.2 on macOS**, even when the interactive shell
is zsh. Run/source the blocks in order in one Bash process. No script uses arrays; never
expand a possibly-empty array under `set -u` (iteration 334's failure).
Run foreground commands; if a tool returns a session ID, poll it to completion with a bounded
loop (at most 60 one-second polls for a small command, or a documented finite bound for a gate).
Do not leave background work or end a turn announcing a wait.

Use the following environment before any Go gate. It copies only pinned dependency inputs
from the existing module cache into the worktree; all cache/temp/state writes stay local.
If an input/toolchain is missing, report it; do not fetch/install under a shared path.
Keep this environment for the controller's outside-sandbox runs as well. Do not repurpose HOME.

```bash
#!/bin/bash
# Run from the executor worktree root; all writes stay under this root.
set -euo pipefail
ROOT=$(pwd -P)
E="$ROOT/verification/world-iter335/executor"
export GOTOOLCHAIN=local
module_seed=$(go env GOMODCACHE)
export GOCACHE="$E/go-build" GOMODCACHE="$E/modules" GOPATH="$E/go"
export TMPDIR="$E/tmp" GOTMPDIR="$E/tmp"
export XDG_CACHE_HOME="$E/cache" XDG_CONFIG_HOME="$E/config" XDG_STATE_HOME="$E/state"
export AILANG_STATE_DIR="$E/ailang-state" AILANG_CACHE_DIR="$E/ailang-cache"
export GOPROXY=off GOSUMDB=off GOFLAGS=-mod=readonly
mkdir -p "$E" "$GOCACHE" "$GOMODCACHE" "$GOPATH" "$TMPDIR"
mkdir -p "$XDG_CACHE_HOME" "$XDG_CONFIG_HOME" "$XDG_STATE_HOME" "$AILANG_STATE_DIR" "$AILANG_CACHE_DIR"
# Clone graph/download metadata, then only the source modules in go.mod.
# The source cache is READ ONLY; cp -cR uses macOS copy-on-write clones.
if [ ! -d "$GOMODCACHE/cache" ]; then
  cp -cR "$module_seed/cache" "$GOMODCACHE/cache"
fi
awk '$1 ~ /^(github[.]com|golang[.]org|modernc[.]org)\// {print $1, $2}' go.mod |
while read -r mod ver; do
  mkdir -p "$GOMODCACHE/$(dirname "$mod")"
  if [ ! -d "$GOMODCACHE/$mod@$ver" ]; then
    cp -cR "$module_seed/$mod@$ver" "$GOMODCACHE/$mod@$ver"
  fi
done
```

Capture each gate's stdout/stderr and actual exit status under `$E`, without a pipeline masking
that status. **Any gate result produced inside the Codex sandbox that fails on a socket bind
or a path outside the worktree is labelled `UNINFORMATIVE UNDER SANDBOX`, never pass or fail.**
This applies to base, clean and mutant runs; it cannot satisfy an expected-failure drill.
Inspect logs for paths even when the script's coarse permission detector does not match them.
Report the exact gate/error and leave its outside-sandbox result to the controller. Do not change
code, weaken a test, skip a gate or alter its expected result to work around such an error.

## Exact mutation drills and restores

Execute this block in the configured Bash shell after creating the final file. `m1`, `m3a`
and `m3c` below are the exact sed edits; M2/M2b/M4 use final-file line numbers/anchored text.
Every row begins from the pinned original except the explicitly paired mutations. Every
`restore` copies the two saved originals and checks **full SHA-256 byte identity**, including
production `4e320c537e68…`; no git mutation is involved. No gofmt on mutants.
The EXIT trap also restores on a failed apply, vet, test or interrupted gate.

| Drill | Applied edit | Expected outcome |
|---|---|---|
| M1 | Drop only the record guard clause | FAIL with `approve with record: 200` |
| M1-existing | Same M1, skip new test | PASS, six original tests |
| M3a | M1 plus same refusal after RecordOperatorApproval | FAIL with `refused approve recorded an approval: head=true` |
| M3c | Force approve guard true | FAIL with `control approve without record: 403` |
| M2-clean | Delete final :57–66 and json import | PASS on clean production |
| M2+M3c | M2 against always-refusing guard | PASS: the deleted control was necessary |
| M2b | Change control's expected 200 to 403 | FAIL with `control approve without record: 200` |
| M4+M3a | Delete only post-refusal :53–55; write then refuse | PASS: the deleted no-write assertion was necessary |

```bash
#!/bin/bash
# Source the environment block first, then run this block in the same shell.
set -euo pipefail
P=host/daemon/episode_finish.go
T=host/daemon/episode_finish_approve_guard_test.go
N='^TestEpisodeFinishApproveRefusesRecord$'
sha() { shasum -a 256 "$1" | awk '{print $1}'; }
[ "$(sha "$P")" = "4e320c537e689620480af2bdff9bdcd9d4d6febba43b4a5dd5b3bb09971973bf" ]
[ "$(sha "$T")" = "c9e279de7c7ce87edf5fbc0e43e690253616d320b5585625062cda3d29cdec57" ]
cp "$P" "$E/production.base"
cp "$T" "$E/final-test.base"
restore() {
  cp "$E/production.base" "$P"
  cp "$E/final-test.base" "$T"
  [ "$(sha "$P")" = "4e320c537e689620480af2bdff9bdcd9d4d6febba43b4a5dd5b3bb09971973bf" ]
  [ "$(sha "$T")" = "c9e279de7c7ce87edf5fbc0e43e690253616d320b5585625062cda3d29cdec57" ]
  shasum -a 256 "$P" "$T"
}
trap restore EXIT
trap 'exit 130' INT
trap 'exit 143' TERM
# run LABEL EXPECTED_RC REQUIRED_LITERAL COMMAND [ARG ...]
# "$@" is positional arguments, never a possibly-empty Bash array.
run() {
  local label=$1 want=$2 marker=$3 rc=0
  shift 3
  "$@" > "$E/$label.txt" 2>&1 || rc=$?
  cat "$E/$label.txt"
  printf '%s rc=%s expected=%s\n' "$label" "$rc" "$want"
  # Inspect ALL other sandbox/path errors manually under the rule below.
  if grep -Eiq '(bind:.*(not permitted|permission denied)|sandbox.*deny|operation not permitted)' "$E/$label.txt"; then
    printf '%s\n' 'UNINFORMATIVE UNDER SANDBOX'
    return 1
  fi
  [ "$rc" -eq "$want" ]
  if [ -n "$marker" ]; then grep -Fq -- "$marker" "$E/$label.txt"; fi
}
probe() {
  shasum -a 256 "$P" "$T" > "$E/$1.hashes.txt"
  run "$1-vet" 0 '' go vet ./host/daemon/
  run "$1" "$2" "$3" go test -count=1 -v -run "$N" ./host/daemon/
}
m1() {
  sed 's/ || input.Record != nil {/ {/' "$P" > "$E/edited"
  cp "$E/edited" "$P"
}
m3c() {
  sed 's/if err != nil || scope.ExpiresAt <= now || scope.ExpiresAt > now+3600 || input.Record != nil {/if true || err != nil || scope.ExpiresAt <= now || scope.ExpiresAt > now+3600 || input.Record != nil {/' "$P" > "$E/edited"
  cp "$E/edited" "$P"
}
m3a() {
  m1
  sed '/^[[:space:]]*request, decision, err := broker.RecordOperatorApproval(ctx, d.store, input.Scope, "operator:episode-finish", "operator:terminal", now)$/a\
		if input.Record != nil {\
			refuse("operator route: invalid or expired finish scope")\
			return\
		}
' "$P" > "$E/edited"
  cp "$E/edited" "$P"
}
probe clean 0 '--- PASS: TestEpisodeFinishApproveRefusesRecord'
m1
probe M1 1 'approve with record: 200'
run M1-existing 0 'ok' go test -count=1 -run Finish -skip "$N" ./host/daemon/
restore
m3a
probe M3a 1 'refused approve recorded an approval: head=true'
restore
m3c
probe M3c 1 'control approve without record: 403'
restore
# M2: remove final-file lines 57..66 (control), and line 4 (json import).
sed '4d;57,66d' "$T" > "$E/edited"
cp "$E/edited" "$T"
probe M2-clean 0 '--- PASS: TestEpisodeFinishApproveRefusesRecord'
m3c
probe M2+M3c 0 '--- PASS: TestEpisodeFinishApproveRefusesRecord'
restore
# M2b: invert only the control's expected status.
sed 's/if w.Code != 200 || json.Unmarshal/if w.Code != 403 || json.Unmarshal/' "$T" > "$E/edited"
cp "$E/edited" "$T"
probe M2b 1 'control approve without record: 200'
restore
# M4: delete ONLY lines 53..55 (post-refusal no-write check).
# Keep lines 44..46 (the new pristine-store precondition).
sed '53,55d' "$T" > "$E/edited"
cp "$E/edited" "$T"
m3a
probe M4+M3a 0 '--- PASS: TestEpisodeFinishApproveRefusesRecord'
restore
trap - EXIT INT TERM
```

## Clean gates and controller handoff

Run these only after successful hash restoration. The additional verbose Finish invocation
observes the seven top-level names and PASS count; it does not replace checking the exit status.
Record each command individually even if an earlier gate is red, then report the red finding.
If running as one script with `set -e`, resume the remaining clean gates individually after
recording the error. Do not classify sandbox errors as ordinary reds.

```bash
#!/bin/bash
# Same shell/environment as above, AFTER all mutations have been restored.
go vet ./host/daemon/
gofmt -l host/daemon/ > "$E/gofmt.txt"
cat "$E/gofmt.txt"
test ! -s "$E/gofmt.txt"
go test -count=1 -run '^TestEpisodeFinishApproveRefusesRecord$' ./host/daemon/
go test -count=20 -run '^TestEpisodeFinishApproveRefusesRecord$' ./host/daemon/
go test -race -count=5 -run '^TestEpisodeFinishApproveRefusesRecord$' ./host/daemon/
go test -count=1 -run Finish ./host/daemon/
go test -count=1 -v -run Finish ./host/daemon/ > "$E/finish-verbose.txt" 2>&1
cat "$E/finish-verbose.txt"
test "$(grep -c '^--- PASS:' "$E/finish-verbose.txt")" -eq 7
shasum -a 256 host/daemon/episode_finish.go host/daemon/episode_finish_test.go host/daemon/episode_finish_approve_guard_test.go
wc -l host/daemon/episode_finish.go host/daemon/episode_finish_test.go host/daemon/episode_finish_approve_guard_test.go
git diff --exit-code -- host/daemon/episode_finish.go host/daemon/episode_finish_test.go
git status --short
```

The controller re-runs the package and full gates **outside the sandbox**, including broker
because both whole packages need `AILANG_BIN=$HOME/.pinned-ailang/ailang` (AILANG v0.41.0).
The binary is a read-only input; verify its version without changing it. Use the same local
cache/temp/state paths and existing containment of writes to the worktree. If an older binary
ignores a state/cache override and tries an external path, report the restriction; do not seed
or modify that path. No package/full success is claimed by the planner.

```bash
#!/bin/bash
# Controller only, outside the Codex sandbox, with the workspace-local
# cache/temp/state environment above and its existing write containment.
export AILANG_BIN="$HOME/.pinned-ailang/ailang"
"$AILANG_BIN" --version
# Require released AILANG v0.41.0; do not install/replace anything under HOME.
go test -count=1 ./host/daemon/
go test -count=1 ./host/broker/
go build ./... && go test ./...
./scripts/verify_ail.sh
```

Report any unrelated base red, sandbox limitation or missing prerequisite with its raw evidence.
Do not enlarge the single-file implementation to make an unrelated gate green. Before handoff,
observe hashes and line counts once more, and ensure the product diff adds exactly the pinned
67-line file, with no production changes. The controller owns committing and final CI evidence.
