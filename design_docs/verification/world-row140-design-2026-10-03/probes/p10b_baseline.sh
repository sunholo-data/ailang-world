#!/bin/bash
# P10b: the same full suites UNSANDBOXED with the same env — are P10's failures caused by the sandbox?
. "$(dirname "$0")/lib.sh"
T=$S/ts-twilight; Y=$S/py-cli
cd $T
env -i PATH=/usr/bin:/bin:/opt/homebrew/bin HOME=$S/ts-ep/home npm_config_cache=$S/ts-ep/npm-cache TMPDIR=$S/ts-ep/tmp CI=1 NO_COLOR=1 ./node_modules/.bin/vitest run > $S/p10b_vitest_full.log 2>&1
echo "vitest full UNSANDBOXED rc=$?"; tail -5 $S/p10b_vitest_full.log | head -3
grep -h "FAIL " $S/p10_vitest_full.log | sort -u > $S/p10_fail_sbx.txt
grep -h "FAIL " $S/p10b_vitest_full.log | sort -u > $S/p10_fail_bare.txt
echo "vitest failing lines sandboxed=$(wc -l < $S/p10_fail_sbx.txt) bare=$(wc -l < $S/p10_fail_bare.txt); diff:"; diff $S/p10_fail_sbx.txt $S/p10_fail_bare.txt | head -10
cd $Y
env -i PATH=/usr/bin:/bin LANG=C.UTF-8 HOME=$S/py-ep/home TMPDIR=$S/py-ep/tmp PYTHONPYCACHEPREFIX=$S/py-ep/pycache .venv/bin/python -m pytest -q -p no:cacheprovider > $S/p10b_pytest_full.log 2>&1
echo "pytest full UNSANDBOXED rc=$?"; tail -1 $S/p10b_pytest_full.log
grep -h "^FAILED" $S/p10_pytest_full.log | sort > $S/p10_pyfail_sbx.txt; grep -h "^FAILED" $S/p10b_pytest_full.log | sort > $S/p10_pyfail_bare.txt
echo "pytest FAILED sandboxed=$(wc -l < $S/p10_pyfail_sbx.txt) bare=$(wc -l < $S/p10_pyfail_bare.txt); diff:"; diff $S/p10_pyfail_sbx.txt $S/p10_pyfail_bare.txt | head
echo "pytest failures with the operator's real HOME (unsandboxed):"
env -i PATH=/usr/bin:/bin LANG=C.UTF-8 HOME=$HOME .venv/bin/python -m pytest -q -p no:cacheprovider 2>&1 | tail -1
