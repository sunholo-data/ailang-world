#!/bin/bash
# Q2b: run one targeted test command per M4 project in a THROWAWAY APFS clone
# (cp -c / git clone --shared) with World's child env shape (CI=1, HOME/TMPDIR
# in a scratch cache, PYTHONPYCACHEPREFIX, npm_config_cache), then list every
# path under the clone newer than a marker. Unsandboxed: this measures what the
# TOOLCHAIN writes, not what srt allows. The real checkouts are never written.
set -u
SC=/Users/voightkampff/dev/sunholo-data/.scratch-iter277-row156/q2b; rm -rf $SC; mkdir -p $SC/cache/{home,tmp,pyc,npm}
ENV=(env -i PATH=/opt/homebrew/bin:/usr/bin:/bin HOME=$SC/cache/home TMPDIR=$SC/cache/tmp LANG=C.UTF-8 LC_ALL=C.UTF-8
     PYTHONPYCACHEPREFIX=$SC/cache/pyc npm_config_cache=$SC/cache/npm CI=1 NO_COLOR=1)
newer() { find "$1" -newer $SC/marker -not -path '*/.git/*' \( -type f -o -type l \) 2>/dev/null | sed "s#^$1/##" | grep -v '^node_modules/' ; echo "  (+ node_modules entries newer: $(find "$1/node_modules" -newer $SC/marker -type f 2>/dev/null | wc -l | tr -d ' '))"; }
echo "== TwilightGame: vitest run <one file>"
cp -c -R /Users/voightkampff/dev/markedmondson1234/TwilightGame $SC/tg 2>/dev/null
F=$(cd $SC/tg && ls tests/*.test.ts | head -1); sleep 1; touch $SC/marker; sleep 1
(cd $SC/tg && "${ENV[@]}" node node_modules/vitest/vitest.mjs run "$F" 2>&1 | tail -4)
echo "-- paths written under the clone:"; newer $SC/tg
echo "== sunholo-platform/cli: pytest -p no:cacheprovider <one file>, then WITHOUT -p no:cacheprovider"
cp -c -R /Users/voightkampff/dev/sunholo-data/sunholo-platform/cli $SC/cli
touch $SC/marker; sleep 1
(cd $SC/cli && "${ENV[@]}" .venv/bin/python -m pytest -q -p no:cacheprovider tests/test_cli_bq.py 2>&1 | tail -2)
echo "-- paths written (with no:cacheprovider):"; find $SC/cli -newer $SC/marker -type f -not -path '*/.venv/*' | sed "s#^$SC/cli/##"
touch $SC/marker; sleep 1
(cd $SC/cli && "${ENV[@]}" .venv/bin/python -m pytest -q tests/test_cli_bq.py 2>&1 | tail -1)
echo "-- paths written (default cacheprovider):"; find $SC/cli -newer $SC/marker -type f -not -path '*/.venv/*' | sed "s#^$SC/cli/##" | head
echo "== ailang: go test -count=1 ./internal/lexer/ (GOFLAGS=-mod=readonly GOPROXY=off, shared modcache read-only by convention)"
git clone -q --shared /Users/voightkampff/dev/sunholo-data/ailang $SC/ail
touch $SC/marker; sleep 1
(cd $SC/ail && "${ENV[@]}" GOCACHE=$SC/cache/gocache GOMODCACHE=/Users/voightkampff/go/pkg/mod GOFLAGS=-mod=readonly GOPROXY=off GOTELEMETRY=off GOTOOLCHAIN=auto go test -count=1 ./internal/lexer/ 2>&1 | tail -2)
echo "-- paths written under the clone:"; find $SC/ail -newer $SC/marker -not -path '*/.git/*' -type f | sed "s#^$SC/ail/##" | head
rm -rf $SC
