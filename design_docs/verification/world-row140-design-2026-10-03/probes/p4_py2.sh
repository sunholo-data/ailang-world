#!/bin/bash
# P4-Py2: the read fence with the python store's REAL path re-allowed; project code (a test) trying to escape.
. "$(dirname "$0")/lib.sh"
Y=$S/py-cli; EP=$S/py-ep
PYLINK=$HOME/.local/share/uv/python/cpython-3.12-macos-aarch64-none
PYREAL=$(python3 -c "import os,sys;print(os.path.realpath(sys.argv[1]))" $PYLINK)
echo "python store link: $PYLINK -> $PYREAL"
cat > $S/s-py-readfence3.json <<EOF
{ "filesystem": { "denyRead": ["$HOME"], "allowRead": ["$Y", "$EP", "$PYLINK", "$PYREAL"], "allowWrite": ["$Y", "$EP"], "denyWrite": [] },
  "network": { "allowedDomains": [], "deniedDomains": [] } }
EOF
cd $Y
ENVV="env -i PATH=/usr/bin:/bin LANG=C.UTF-8 HOME=$EP/home CLAUDE_CODE_TMPDIR=$EP/tmp PYTHONPYCACHEPREFIX=$EP/pycache"
echo "== PY3c: read fence with the python store's link and real path re-allowed"
run $ENVV $SRT --settings $S/s-py-readfence3.json -- .venv/bin/python -m pytest -q -p no:cacheprovider tests/test_cli_access.py
echo "== PY5: a test (project code, run by definition) tries to write outside and read the decoy"
cat > $Y/tests/test_zz_probe_escape.py <<EOF
def test_escape():
    for p in ["/tmp/srt-row140-pyescape", "$HOME/.srt-row140-pyescape", "$S/other/pyescape"]:
        try:
            open(p, "w").write("x"); r = "WROTE"
        except OSError as e:
            r = "REFUSED " + str(e.strerror)
        print("WRITE", p, r)
    try:
        r = open("$HOME/.srt-row140-decoy").read().strip()
    except OSError as e:
        r = "REFUSED " + str(e.strerror)
    print("READ decoy:", r)
EOF
run $ENVV $SRT --settings $S/s-py.json -- .venv/bin/python -m pytest -q -s -p no:cacheprovider tests/test_zz_probe_escape.py
run $ENVV $SRT --settings $S/s-py-readfence3.json -- .venv/bin/python -m pytest -q -s -p no:cacheprovider tests/test_zz_probe_escape.py
rm -f $Y/tests/test_zz_probe_escape.py
ls -la /tmp/srt-row140-pyescape $HOME/.srt-row140-pyescape $S/other/pyescape 2>&1
