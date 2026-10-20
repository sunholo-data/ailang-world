#!/usr/bin/env python3
"""V7r (iteration 289): re-measure design V7's alias arms on a given policy-tool
binary. Same fixture and patterns as q4_policytool_denywrite.py (V7), plus the
hard-link arms V7 never ran (D-156-4 names symlinks AND hard links) and an
unprotected hard-link control (does the binary refuse every multiply-linked
file, or only aliases of protected ones?).

usage: v7r_policytool_aliases.py <ailang-binary> <scratch-dir>
Host evidence is read back from disk after every arm; the tool's JSON alone is
never the verdict."""
import json, os, shutil, subprocess, sys
BIN, SCR = sys.argv[1], sys.argv[2]
FILES = {"pkg/a.go": "package pkg\n", "pkg/a_test.go": "// ORIG\n", "tests/test_x.py": "# ORIG\n",
         "cli/conftest.py": "# ORIG\n", "src/foo.test.ts": "// ORIG\n", "src/deep/x.test.ts": "// ORIG\n",
         "free/plain.txt": "ORIG\n"}
PATTERNS = ["tests/**", "*_test.go", "*.test.ts", "conftest.py"]

def setup(patterns, hardlinks):
    if os.path.exists(SCR): shutil.rmtree(SCR)
    wt = os.path.join(SCR, "wt")
    for rel, body in FILES.items():
        p = os.path.join(wt, rel); os.makedirs(os.path.dirname(p), exist_ok=True); open(p, "w").write(body)
    os.symlink("tests", os.path.join(wt, "lnk"))                      # dir alias of a protected dir
    os.symlink("pkg/a_test.go", os.path.join(wt, "alias.txt"))        # file symlink alias of a protected file
    if hardlinks:  # only for hard-link arms, so symlink arms' targets stay single-linked (no confound)
        os.link(os.path.join(wt, "pkg/a_test.go"), os.path.join(wt, "hl.txt"))          # hard link to a protected file
        os.link(os.path.join(wt, "tests/test_x.py"), os.path.join(wt, "pkg/hl2.go"))    # hard link, other dir
        os.link(os.path.join(wt, "free/plain.txt"), os.path.join(wt, "free/plain2.txt"))  # hard link, NOTHING protected
    wt = os.path.realpath(wt)
    pol = os.path.join(SCR, "policy.toml")
    deny = [".github/**"] + patterns
    open(pol, "w").write(f'security_mode = "restricted"\nallowed_caps = ["FS", "IO"]\nfs_sandbox = "{wt}"\n'
        f'timeout_ms = 8000\nfs_deny_write = {json.dumps(deny)}\nentry = "main"\n\n[budgets]\nFS = 0\n')
    return wt, pol

def call(wt, pol, req):
    r = subprocess.run([BIN, "policy-tool", "--policy", pol], input=json.dumps(req), cwd=wt, capture_output=True,
                       text=True, env={"HOME": SCR, "PATH": "/usr/bin:/bin", "LANG": "C", "LC_ALL": "C"}, timeout=30)
    try:
        j = json.loads(r.stdout)
        verdict = "ok" if j.get("ok") else "REFUSED: " + str(j.get("refused") or j.get("error"))[:150]
    except Exception:
        verdict = "UNPARSEABLE: " + r.stdout.strip()[:120] + " | " + r.stderr.strip()[:120]
    return r.returncode, verdict

def host_changes(wt):
    ev = []
    for root, dirs, files in os.walk(wt):
        for f in files:
            p = os.path.join(root, f); rel = os.path.relpath(p, wt)
            if rel in ("lnk", "alias.txt") or os.path.islink(p): continue
            b = open(p).read()
            orig = FILES.get(rel)
            if rel == "hl.txt": orig = FILES["pkg/a_test.go"]
            if rel == "pkg/hl2.go": orig = FILES["tests/test_x.py"]
            if rel == "free/plain2.txt": orig = FILES["free/plain.txt"]
            if orig is None: ev.append(f"{rel}:NEW")
            elif b != orig: ev.append(f"{rel}:CHANGED")
    # an alias write that replaced the link (rename-over) instead of writing through it
    for rel in ("hl.txt", "pkg/hl2.go"):
        p = os.path.join(wt, rel)
        if os.path.exists(p) and os.stat(p).st_nlink == 1: ev.append(f"{rel}:LINK-BROKEN(nlink=1)")
        if not os.path.exists(p): continue
    for rel in ("lnk", "alias.txt"):
        if not os.path.islink(os.path.join(wt, rel)): ev.append(f"{rel}:SYMLINK-REPLACED")
    return sorted(ev) or ["none"]

ARMS = [
 ("control write pkg/a.go", {"op": "write", "path": "pkg/a.go", "content": "X\n"}),
 ("lexical write pkg/a_test.go", {"op": "write", "path": "pkg/a_test.go", "content": "PWN\n"}),
 ("lexical write NEW tests/new.py", {"op": "write", "path": "tests/new.py", "content": "PWN\n"}),
 ("symlink-dir write NEW lnk/new5.py", {"op": "write", "path": "lnk/new5.py", "content": "PWN\n"}),
 ("symlink-dir write lnk/test_x.py", {"op": "write", "path": "lnk/test_x.py", "content": "PWN\n"}),
 ("symlink-dir edit lnk/test_x.py", {"op": "edit", "path": "lnk/test_x.py", "old_text": "ORIG", "new_text": "PWN"}),
 ("symlink-file write alias.txt", {"op": "write", "path": "alias.txt", "content": "PWN\n"}),
 ("symlink-file edit alias.txt", {"op": "edit", "path": "alias.txt", "old_text": "ORIG", "new_text": "PWN"}),
 ("hardlink write hl.txt (=pkg/a_test.go)", {"op": "write", "path": "hl.txt", "content": "PWN\n"}),
 ("hardlink edit hl.txt (=pkg/a_test.go)", {"op": "edit", "path": "hl.txt", "old_text": "ORIG", "new_text": "PWN"}),
 ("hardlink write pkg/hl2.go (=tests/test_x.py)", {"op": "write", "path": "pkg/hl2.go", "content": "PWN\n"}),
 ("hardlink edit pkg/hl2.go (=tests/test_x.py)", {"op": "edit", "path": "pkg/hl2.go", "old_text": "ORIG", "new_text": "PWN"}),
 ("UNPROTECTED hardlink write free/plain2.txt", {"op": "write", "path": "free/plain2.txt", "content": "X\n"}),
]
v = subprocess.run([BIN, "--version"], capture_output=True, text=True).stdout.splitlines()
print(f"binary {BIN}: {v[0] if v else '?'} {v[1] if len(v) > 1 else ''}")
n = 0
for label, pats in [("WITHOUT protected patterns (control)", []), ("WITH " + json.dumps(PATTERNS), PATTERNS)]:
    print(f"\n=== fs_deny_write {label}")
    for name, req in ARMS:
        wt, pol = setup(pats, hardlinks=("hardlink" in name))
        rc, verdict = call(wt, pol, req)
        n += 1
        print(f"{name:46s} rc={rc} {verdict:60s} host={','.join(host_changes(wt))}")
print(f"\narms_run={n}")
shutil.rmtree(SCR)
