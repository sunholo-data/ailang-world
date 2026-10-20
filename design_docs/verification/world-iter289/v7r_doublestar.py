#!/usr/bin/env python3
"""V7r-b: V7's pattern-form arms (ONLY ["**/*.test.ts", "tests/**"]) on a given binary,
plus an empty fs_deny_write arm for the multi-link refusal. usage: <bin> <scratch>"""
import json, os, shutil, subprocess, sys
BIN, SCR = sys.argv[1], sys.argv[2]
def run(deny, req, extra_files=(), link=None):
    if os.path.exists(SCR): shutil.rmtree(SCR)
    wt = os.path.join(SCR, "wt")
    for rel in ["src/foo.test.ts", "src/deep/x.test.ts", "free/plain.txt"] + list(extra_files):
        p = os.path.join(wt, rel); os.makedirs(os.path.dirname(p), exist_ok=True); open(p, "w").write("ORIG\n")
    os.makedirs(os.path.join(wt, "cli/tests"), exist_ok=True)
    if link: os.link(os.path.join(wt, link[0]), os.path.join(wt, link[1]))
    wt = os.path.realpath(wt); pol = os.path.join(SCR, "p.toml")
    open(pol, "w").write(f'security_mode = "restricted"\nallowed_caps = ["FS", "IO"]\nfs_sandbox = "{wt}"\n'
        f'timeout_ms = 8000\nfs_deny_write = {json.dumps(deny)}\nentry = "main"\n\n[budgets]\nFS = 0\n')
    r = subprocess.run([BIN, "policy-tool", "--policy", pol], input=json.dumps(req), cwd=wt, capture_output=True, text=True,
                       env={"HOME": SCR, "PATH": "/usr/bin:/bin", "LANG": "C", "LC_ALL": "C"}, timeout=30)
    j = json.loads(r.stdout); p = os.path.join(wt, req["path"])
    host = "absent" if not os.path.exists(p) else ("ORIG" if open(p).read() == "ORIG\n" else "WRITTEN")
    return ("ok" if j.get("ok") else "REFUSED: " + str(j.get("refused"))[:110]), host
print("binary", subprocess.run([BIN, "--version"], capture_output=True, text=True).stdout.splitlines()[0])
D = ["**/*.test.ts", "tests/**"]
for name, deny, req, kw in [
  ("depth0 a.test.ts", D, {"op": "write", "path": "a.test.ts", "content": "P"}, {}),
  ("depth1 src/foo.test.ts", D, {"op": "write", "path": "src/foo.test.ts", "content": "P"}, {}),
  ("depth2 src/deep/x.test.ts", D, {"op": "write", "path": "src/deep/x.test.ts", "content": "P"}, {}),
  ("nested cli/tests/x.py vs tests/**", D, {"op": "write", "path": "cli/tests/x.py", "content": "P"}, {}),
  ("EMPTY deny, single-link free/plain.txt", [], {"op": "write", "path": "free/plain.txt", "content": "P"}, {}),
  ("EMPTY deny, hard-linked free/plain2.txt", [], {"op": "write", "path": "free/plain2.txt", "content": "P"}, {"link": ("free/plain.txt", "free/plain2.txt")}),
]:
    v, h = run(deny, req, **kw); print(f"{name:42s} {v:75s} host={h}")
shutil.rmtree(SCR)
