#!/bin/bash
# P4-Py: pytest on a scratch APFS clone of sunholo-platform/cli (with its .venv) under srt.
. "$(dirname "$0")/lib.sh"
SRC=$HOME/dev/sunholo-data/sunholo-platform/cli
Y=$S/py-cli; EP=$S/py-ep
if [ ! -d $Y ]; then cp -cR $SRC $Y; find $Y -name __pycache__ -type d -prune -exec rm -rf {} +; rm -rf $Y/.pytest_cache; fi
rm -rf $EP; mkdir -p $EP/home $EP/tmp $EP/pycache $EP/uvcache
PYHOME=$HOME/.local/share/uv/python/cpython-3.12-macos-aarch64-none
cat > $S/s-py.json <<EOF
{ "filesystem": { "denyRead": [], "allowRead": [], "allowWrite": ["$Y", "$EP"], "denyWrite": [] },
  "network": { "allowedDomains": [], "deniedDomains": [] } }
EOF
cat > $S/s-py-readfence.json <<EOF
{ "filesystem": { "denyRead": ["$HOME"], "allowRead": ["$Y", "$EP", "$PYHOME"], "allowWrite": ["$Y", "$EP"], "denyWrite": [] },
  "network": { "allowedDomains": [], "deniedDomains": [] } }
EOF
cd $Y
ENVV="env -i PATH=/usr/bin:/bin LANG=C.UTF-8 HOME=$EP/home CLAUDE_CODE_TMPDIR=$EP/tmp"
echo "== PY1: the copied venv's python (symlink into the uv python store) -m pytest, default env"
t0=$(python3 -c 'import time;print(time.time())')
run $ENVV $SRT --settings $S/s-py.json -- .venv/bin/python -m pytest -q -p no:cacheprovider tests/test_cli_access.py
t1=$(python3 -c 'import time;print(time.time())'); python3 -c "print('PY1 wall %.1fs' % ($t1-$t0))"
echo "== PY1-where: did pytest/python write bytecode into the worktree?"
run find $Y -name __pycache__ -not -path '*/.venv/*' -type d
find $Y -name __pycache__ -type d -prune -exec rm -rf {} +
echo "== PY2: PYTHONPYCACHEPREFIX + PYTHONDONTWRITEBYTECODE alternatives; cache provider on"
run $ENVV PYTHONPYCACHEPREFIX=$EP/pycache $SRT --settings $S/s-py.json -- .venv/bin/python -m pytest -q tests/test_cli_access.py
run find $Y -name __pycache__ -not -path '*/.venv/*' -type d
run find $Y -maxdepth 1 -name .pytest_cache
run du -sh $EP/pycache
echo "== PY3: the read fence (HOME denied; worktree, episode cache and the uv python store re-allowed)"
run $ENVV PYTHONPYCACHEPREFIX=$EP/pycache $SRT --settings $S/s-py-readfence.json -- .venv/bin/python -m pytest -q -p no:cacheprovider tests/test_cli_access.py
echo "== PY3b: the read fence WITHOUT the python store re-allowed"
cat > $S/s-py-readfence2.json <<EOF
{ "filesystem": { "denyRead": ["$HOME"], "allowRead": ["$Y", "$EP"], "allowWrite": ["$Y", "$EP"], "denyWrite": [] },
  "network": { "allowedDomains": [], "deniedDomains": [] } }
EOF
run $ENVV $SRT --settings $S/s-py-readfence2.json -- .venv/bin/python -c 'print(1)'
echo "== PY4: uv run --offline (per-episode UV_CACHE_DIR, cold) vs the venv"
run $ENVV UV_CACHE_DIR=$EP/uvcache $SRT --settings $S/s-py.json -- $HOME/.local/bin/uv run --offline --frozen pytest -q -p no:cacheprovider tests/test_cli_access.py
echo "== PY5: conftest.py executes project code (inherent): a conftest writing outside the worktree is fenced"
cat > $Y/tests/zz_probe_conftest_check.py <<'EOF'
import os
def test_escape_write():
    try:
        open(os.path.expanduser("~") + "/../../../tmp/srt-row140-pyescape", "w").write("x")
        r = "WROTE"
    except OSError as e:
        r = "REFUSED " + e.strerror
    print("ESCAPE:", r)
EOF
run $ENVV $SRT --settings $S/s-py.json -- .venv/bin/python -m pytest -q -s -p no:cacheprovider tests/zz_probe_conftest_check.py
rm -f $Y/tests/zz_probe_conftest_check.py
