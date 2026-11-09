#!/usr/bin/env bash
# Row 256 net6 attribution, phase-randomized (the start page dials every 32 s).
set -u
S=$SCRATCH
OUT=$S/m256j; mkdir -p $OUT; : > $OUT/runs.tsv
export AILANG_BIN=$HOME/.pinned-ailang/ailang
N=${N:-100}
T='^TestExecProbeNet6ForeignDialIsRetriedNotRefused$/^raw-twice-proxy$'
# scanner witness: a listener that logs every unsolicited accept, for the whole run
python3 -I $S/listen_probe.py 3000 > $OUT/scanner.log 2>&1 &
W=$!
for i in $(seq 1 $N); do
  if [ $((RANDOM % 2)) -eq 0 ]; then order="base cand"; else order="cand base"; fi
  for arm in $order; do
    d=$S/wt-row256; [ $arm = base ] && d=$S/wt-base256
    t0=$(python3 -c 'import time;print(f"{time.time():.3f}")')
    ( cd $d/host/broker && $S/m256/broker-$arm.test -test.run "$T" -test.count=1 -test.v ) > $OUT/run.tmp 2>&1
    rc=$?
    t1=$(python3 -c 'import time;print(f"{time.time():.3f}")')
    printf '%s\t%s\t%s\t%s\t%s\n' $i $arm $rc $t0 $t1 >> $OUT/runs.tsv
    [ $rc -ne 0 ] && cp $OUT/run.tmp $OUT/fail-$arm-$i.log
    sleep 0.$((RANDOM % 10))$((RANDOM % 10)); sleep $((RANDOM % 5))
  done
done
kill $W 2>/dev/null
awk -F'\t' '{n[$2]++; if($3!=0) f[$2]++} END{for(a in n) printf "%s fail=%d/%d\n", a, f[a]+0, n[a]}' $OUT/runs.tsv > $OUT/summary.txt
echo DONE >> $OUT/summary.txt
