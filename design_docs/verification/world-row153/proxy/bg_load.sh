#!/bin/bash
# Row 153 AC3.8 / AC4.1 throttled-load proxy (the design's bg_load.sh, parameterised).
# usage: bg_load.sh <burners> <count> <out> <test-binary> <package-dir>
#   burners      number of background-QoS CPU burners (`taskpolicy -b yes`)
#   count        -test.count for the two e2e tests
#   out          output log
#   test-binary  a `go test -race -c` binary of host/daemon (probe-instrumented)
#   package-dir  the host/daemon directory the binary runs in (its testdata paths are relative)
# The test binary itself runs under background QoS (E-cores), which is what makes a CI-like
# slow host reproducible on an otherwise idle Apple-silicon machine. Overridable: RUNPAT.
n=$1; count=$2; out=$3; bin=$4; pkg=$5
export AILANG_BIN=$HOME/.pinned-ailang/ailang WORLD_EXEC_SRT_NODE_MODULES=$HOME/.ailang/state/mission-world-iter241/srt/node_modules WORLD_PROBE_PHASE_TIMING=1 PATH=/opt/homebrew/bin:$PATH
pids=()
for i in $(seq 1 $n); do taskpolicy -b yes > /dev/null & pids+=($!); done
cd "$pkg" || exit 2
taskpolicy -b "$bin" -test.count=$count -test.v -test.run "${RUNPAT:-^(TestExecSrtMCPEndToEnd|TestExecSrtMCPGrantAndBudget)$}" > "$out" 2>&1
echo "RC=$?" >> "$out"
kill "${pids[@]}" 2>/dev/null; wait 2>/dev/null
echo DONE >> "$out"
