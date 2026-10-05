#!/bin/bash
# P11: dependency installation inside the fence with a registry-only network allowlist (npm ci; uv sync).
. "$(dirname "$0")/lib.sh"
N=$S/npm-tiny; NE=$S/npm-tiny-ep
rm -rf $N $NE; mkdir -p $N $NE/home $NE/tmp $NE/npm-cache
cat > $N/package.json <<'EOF'
{ "name": "tiny", "version": "0.0.0", "private": true, "dependencies": { "left-pad": "1.3.0" },
  "scripts": { "test": "node -e \"console.log(require('left-pad')('x', 3))\"" } }
EOF
cat > $S/s-npm-reg.json <<EOF
{ "filesystem": { "denyRead": [], "allowRead": [], "allowWrite": ["$N", "$NE"], "denyWrite": [] },
  "network": { "allowedDomains": ["registry.npmjs.org"], "deniedDomains": [] } }
EOF
cat > $S/s-npm-none.json <<EOF
{ "filesystem": { "denyRead": [], "allowRead": [], "allowWrite": ["$N", "$NE"], "denyWrite": [] },
  "network": { "allowedDomains": [], "deniedDomains": [] } }
EOF
NENV="env -i PATH=/usr/bin:/bin:/opt/homebrew/bin HOME=$NE/home npm_config_cache=$NE/npm-cache CLAUDE_CODE_TMPDIR=$NE/tmp"
cd $N
echo "== D1: npm install (lockfile creation) with the registry allowed"
run $NENV $SRT --settings $S/s-npm-reg.json -- npm install --ignore-scripts --no-audit --no-fund
echo "== D2: npm ci from the lockfile, network NONE but the per-episode npm cache now warm (offline)"
rm -rf node_modules
run $NENV $SRT --settings $S/s-npm-none.json -- npm ci --offline --ignore-scripts --no-audit --no-fund
run $NENV $SRT --settings $S/s-npm-none.json -- npm test
echo "== D3: uv sync --frozen for the python cli copy into a fresh venv, pypi allowed, per-episode uv cache"
Y=$S/py-cli2; YE=$S/py-ep2
rm -rf $Y $YE; mkdir -p $YE/home $YE/tmp $YE/uvcache
mkdir -p $Y && cp -c $S/py-cli/pyproject.toml $S/py-cli/uv.lock $Y/ && cp -cR $S/py-cli/aiplatform $S/py-cli/tests $Y/
PYREAL=$(python3 -c "import os;print(os.path.realpath(os.path.expanduser('~/.local/share/uv/python/cpython-3.12-macos-aarch64-none')))")
cat > $S/s-uv.json <<EOF
{ "filesystem": { "denyRead": [], "allowRead": [], "allowWrite": ["$Y", "$YE"], "denyWrite": [] },
  "network": { "allowedDomains": ["pypi.org", "files.pythonhosted.org"], "deniedDomains": [] } }
EOF
cd $Y
UENV="env -i PATH=/usr/bin:/bin LANG=C.UTF-8 HOME=$YE/home UV_CACHE_DIR=$YE/uvcache UV_PYTHON_DOWNLOADS=never UV_PYTHON=$PYREAL/bin/python3.12 CLAUDE_CODE_TMPDIR=$YE/tmp"
t0=$(python3 -c 'import time;print(time.time())')
run $UENV $SRT --settings $S/s-uv.json -- $HOME/.local/bin/uv sync --frozen
t1=$(python3 -c 'import time;print(time.time())'); python3 -c "print('uv sync wall %.1fs' % ($t1-$t0))"
run $UENV $SRT --settings $S/s-uv.json -- .venv/bin/python -m pytest -q -p no:cacheprovider tests/test_cli_access.py
