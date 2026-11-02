#!/usr/bin/env bash
# Row 237 mutation drills, iteration 335 designer. Run from the worktree root.
# Each drill: apply one substitution, run, restore, print sha256 before/edited/restored.
set -u
cd "$(git rev-parse --show-toplevel)"
PROD=host/daemon/episode_finish.go
TEST=host/daemon/episode_finish_approve_guard_test.go
BAK=verification/world-iter335/designer/test.bak
sha() { shasum -a 256 "$1" | cut -c1-12; }
NEW='^TestEpisodeFinishApproveRefusesRecord$'
run() { # label, go test args...
  local label=$1; shift
  echo "---- $label: go test -count=1 $*"
  go test -count=1 "$@" 2>&1 | grep -E '^(---|ok|FAIL|PASS)|_test.go:[0-9]+:' ; echo "rc=${PIPESTATUS[0]}"
}
echo "== base: prod $(sha $PROD) test $(sha $TEST) =="
run CONTROL-new -v -run "$NEW" ./host/daemon/
run CONTROL-finish-all -run Finish ./host/daemon/

prod_mut() { # id, perl expr
  local id=$1 expr=$2 before; before=$(sha $PROD)
  perl -0pi -e "$expr" $PROD
  echo "== $id: prod $before -> $(sha $PROD); diff lines: $(git diff --numstat -- $PROD)"
  go vet ./host/daemon/ >/dev/null 2>&1 && echo "vet ok" || echo "VET FAILED"
}
prod_restore() { git checkout -- $PROD; echo "restored prod $(sha $PROD)"; }

# M1: the row's mutant.
prod_mut M1 's/ \|\| input\.Record != nil \{/ {/'
run M1-new -v -run "$NEW" ./host/daemon/
run M1-existing-finish -run Finish -skip "$NEW" ./host/daemon/
prod_restore

# M3a: write-then-refuse. Guard arm dropped; the refusal moves after RecordOperatorApproval.
prod_mut M3a 's/ \|\| input\.Record != nil \{/ {/; s/(request, decision, err := broker\.RecordOperatorApproval\([^\n]*\)\n)/$1\t\tif input.Record != nil {\n\t\t\trefuse("operator route: invalid or expired finish scope")\n\t\t\treturn\n\t\t}\n/'
run M3a-new -v -run "$NEW" ./host/daemon/
prod_restore

# M3c: the approve guard refuses everything (a different early refusal).
prod_mut M3c 's/if err != nil \|\| scope\.ExpiresAt <= now \|\| scope\.ExpiresAt > now\+3600 \|\| input\.Record != nil \{/if true || err != nil || scope.ExpiresAt <= now || scope.ExpiresAt > now+3600 || input.Record != nil {/'
run M3c-new -v -run "$NEW" ./host/daemon/
prod_restore

# M2 (test-side): the control arm removed (test lines 50-60 deleted), run against clean prod and M3c.
cp $TEST $BAK
perl -0pi -e 's/\n\tproof = finishProofFile.*?\n\t\}\n(\t\}\n)?\}\n$/\n}\n/s' $TEST
perl -0pi -e 's/\t"encoding\/json"\n//' $TEST
echo "== M2: test $(sha $BAK) -> $(sha $TEST); lines $(wc -l < $TEST)"
go vet ./host/daemon/ >/dev/null 2>&1 && echo "vet ok" || echo "VET FAILED"
run M2-clean -v -run "$NEW" ./host/daemon/
prod_mut M2+M3c 's/if err != nil \|\| scope\.ExpiresAt <= now \|\| scope\.ExpiresAt > now\+3600 \|\| input\.Record != nil \{/if true || err != nil || scope.ExpiresAt <= now || scope.ExpiresAt > now+3600 || input.Record != nil {/'
run M2+M3c-new -v -run "$NEW" ./host/daemon/
prod_restore
cp $BAK $TEST

# M2b (test-side): the control's expectation inverted (want 403 instead of 200) on clean prod.
perl -0pi -e 's/if w\.Code != 200 \|\| json\.Unmarshal/if w.Code != 403 || json.Unmarshal/' $TEST
echo "== M2b: test $(sha $BAK) -> $(sha $TEST)"
run M2b-clean -v -run "$NEW" ./host/daemon/
cp $BAK $TEST

# M4 (test-side): the no-approval assertion removed, run against M3a.
perl -0pi -e 's/\tif head, requests, decisions := approvalWrites\(t, d\); head .*?\n\t\}\n//s' $TEST
echo "== M4: test $(sha $BAK) -> $(sha $TEST)"
go vet ./host/daemon/ >/dev/null 2>&1 && echo "vet ok" || echo "VET FAILED"
prod_mut M4+M3a 's/ \|\| input\.Record != nil \{/ {/; s/(request, decision, err := broker\.RecordOperatorApproval\([^\n]*\)\n)/$1\t\tif input.Record != nil {\n\t\t\trefuse("operator route: invalid or expired finish scope")\n\t\t\treturn\n\t\t}\n/'
run M4+M3a-new -v -run "$NEW" ./host/daemon/
prod_restore
cp $BAK $TEST; rm -f $BAK

echo "== restored: prod $(sha $PROD) test $(sha $TEST) =="
git status --porcelain
