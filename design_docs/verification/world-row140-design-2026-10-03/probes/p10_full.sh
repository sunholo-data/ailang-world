#!/bin/bash
# P10: does a FULL suite fit the frozen 10 s handler cap? (TwilightGame vitest run; sunholo-platform cli pytest)
. "$(dirname "$0")/lib.sh"
T=$S/ts-twilight; Y=$S/py-cli
cd $T
t0=$(python3 -c 'import time;print(time.time())')
env -i PATH=/usr/bin:/bin:/opt/homebrew/bin HOME=$S/ts-ep/home npm_config_cache=$S/ts-ep/npm-cache CLAUDE_CODE_TMPDIR=$S/ts-ep/tmp CI=1 NO_COLOR=1 $SRT --settings $S/s-ts.json -- ./node_modules/.bin/vitest run > $S/p10_vitest_full.log 2>&1
echo "vitest full rc=$?"
t1=$(python3 -c 'import time;print(time.time())'); python3 -c "print('TS full vitest wall %.1fs' % ($t1-$t0))"
tail -6 $S/p10_vitest_full.log
echo "vitest full stdout+stderr bytes: $(wc -c < $S/p10_vitest_full.log)"
cd $Y
t0=$(python3 -c 'import time;print(time.time())')
env -i PATH=/usr/bin:/bin LANG=C.UTF-8 HOME=$S/py-ep/home CLAUDE_CODE_TMPDIR=$S/py-ep/tmp PYTHONPYCACHEPREFIX=$S/py-ep/pycache $SRT --settings $S/s-py.json -- .venv/bin/python -m pytest -q -p no:cacheprovider > $S/p10_pytest_full.log 2>&1
echo "pytest full rc=$?"
t1=$(python3 -c 'import time;print(time.time())'); python3 -c "print('PY full pytest wall %.1fs' % ($t1-$t0))"
tail -3 $S/p10_pytest_full.log
echo "pytest full output bytes: $(wc -c < $S/p10_pytest_full.log)"
