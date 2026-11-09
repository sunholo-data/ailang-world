#!/usr/bin/env bash
# Row 256 attended attribution measurement (D-WORLD-111 = B).
set -u
S=$SCRATCH
CAND=$S/wt-row256; BASE=$S/wt-base256; OUT=$S/m256; mkdir -p $OUT
export AILANG_BIN=$HOME/.pinned-ailang/ailang
"$AILANG_BIN" --version > $OUT/ailang-version.txt 2>&1
N=${N:-100}

# 1. full suites, sequential, same env
for arm in cand base; do
  d=$CAND; [ $arm = base ] && d=$BASE
  ( cd $d && uptime && git rev-parse HEAD && go test -count=1 ./... ) > $OUT/full-$arm.log 2>&1
  echo "full-$arm rc=$?" >> $OUT/summary.txt
  grep -E '^(--- FAIL|FAIL|ok )' $OUT/full-$arm.log | grep -E 'FAIL' > $OUT/full-$arm.fails
done

# 2. compiled test binaries for host/broker
for arm in cand base; do
  d=$CAND; [ $arm = base ] && d=$BASE
  ( cd $d/host/broker && go test -c -o $OUT/broker-$arm.test . ) || echo "compile $arm failed" >> $OUT/summary.txt
done

# 3. paired, alternating arms
for t in TestExecProbeNet6ForeignDialIsRetriedNotRefused TestAttendedPublishMintsThroughTheLandedTraversalAndSpendsExactlyOnce TestWorktreeSnapshotRefusesAConfiguredFilter; do
  fc=0; fb=0
  for i in $(seq 1 $N); do
    for arm in base cand; do
      d=$CAND; [ $arm = base ] && d=$BASE
      ( cd $d/host/broker && $OUT/broker-$arm.test -test.run "^$t\$" -test.count=1 ) > $OUT/run.tmp 2>&1
      rc=$?
      if [ $rc -ne 0 ]; then
        cp $OUT/run.tmp $OUT/fail-$t-$arm-$i.log
        if [ $arm = cand ]; then fc=$((fc+1)); else fb=$((fb+1)); fi
      fi
    done
  done
  echo "$t base_fail=$fb/$N cand_fail=$fc/$N" >> $OUT/summary.txt
done
echo DONE >> $OUT/summary.txt
