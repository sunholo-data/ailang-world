#!/usr/bin/env python3
"""Q4: does the pinned tool binary's policy-tool (v0.52.1, the engine behind
World's ailang-write/ailang-edit) refuse protected-path writes when the
patterns are appended to fs_deny_write? Arms: lexical forms, `..`, case, and a
symlink alias that already exists in the worktree."""
import json, os, shutil, subprocess
BIN = "/Users/voightkampff/.pinned-ailang-tools/v0.52.1/ailang"
SCR = "/Users/voightkampff/dev/sunholo-data/.scratch-iter277-row156/q4"
FILES = {"pkg/a.go": "package pkg\n", "pkg/a_test.go": "// ORIG\n", "tests/test_x.py": "# ORIG\n",
         "cli/conftest.py": "# ORIG\n", "src/foo.test.ts": "// ORIG\n", "src/deep/x.test.ts": "// ORIG\n"}
def setup(patterns):
    if os.path.exists(SCR): shutil.rmtree(SCR)
    wt = os.path.join(SCR, "wt")
    for rel, body in FILES.items():
        p = os.path.join(wt, rel); os.makedirs(os.path.dirname(p), exist_ok=True); open(p, "w").write(body)
    os.makedirs(os.path.join(wt, "cli/tests"), exist_ok=True); os.symlink("tests", os.path.join(wt, "lnk"))          # dir alias of a protected dir
    os.symlink("pkg/a_test.go", os.path.join(wt, "alias.txt"))  # file alias of a protected file
    wt = os.path.realpath(wt)
    pol = os.path.join(SCR, "policy.toml")
    deny = [".github/**"] + patterns
    open(pol, "w").write(f'security_mode = "restricted"\nallowed_caps = ["FS", "IO"]\nfs_sandbox = "{wt}"\n'
        f'timeout_ms = 8000\nfs_deny_write = {json.dumps(deny)}\nentry = "main"\n\n[budgets]\nFS = 0\n')
    return wt, pol
def call(wt, pol, req):
    r = subprocess.run([BIN, "policy-tool", "--policy", pol], input=json.dumps(req), cwd=wt, capture_output=True, text=True,
                       env={"HOME": SCR, "PATH": "/usr/bin:/bin", "LANG": "C", "LC_ALL": "C"}, timeout=30)
    return r.returncode, r.stdout.strip()[:160], r.stderr.strip()[:120]
PATTERNS = ["tests/**", "*_test.go", "*.test.ts", "conftest.py", "**/*.test.ts"]
ARMS = [
 ("control write pkg/a.go", {"op": "write", "path": "pkg/a.go", "content": "X\n"}),
 ("write existing pkg/a_test.go", {"op": "write", "path": "pkg/a_test.go", "content": "PWN\n"}),
 ("edit existing pkg/a_test.go", {"op": "edit", "path": "pkg/a_test.go", "old_text": "ORIG", "new_text": "PWN"}),
 ("write NEW pkg/b_test.go", {"op": "write", "path": "pkg/b_test.go", "content": "PWN\n"}),
 ("write NEW tests/new.py", {"op": "write", "path": "tests/new.py", "content": "PWN\n"}),
 ("write NEW cli/sub/conftest.py", {"op": "write", "path": "cli/sub/conftest.py", "content": "PWN\n"}),
 ("write src/deep/x.test.ts", {"op": "write", "path": "src/deep/x.test.ts", "content": "PWN\n"}),
 ("dotdot pkg/../tests/new2.py", {"op": "write", "path": "pkg/../tests/new2.py", "content": "PWN\n"}),
 ("abs path <wt>/tests/new3.py", {"op": "write", "path": "@WT@/tests/new3.py", "content": "PWN\n"}),
 ("case TESTS/new4.py", {"op": "write", "path": "TESTS/new4.py", "content": "PWN\n"}),
 ("case pkg/A_TEST.GO", {"op": "write", "path": "pkg/A_TEST.GO", "content": "PWN\n"}),
 ("symlink dir lnk/new5.py (lnk->tests)", {"op": "write", "path": "lnk/new5.py", "content": "PWN\n"}),
 ("symlink dir lnk/test_x.py (existing)", {"op": "write", "path": "lnk/test_x.py", "content": "PWN\n"}),
 ("symlink file alias.txt -> pkg/a_test.go", {"op": "write", "path": "alias.txt", "content": "PWN\n"}),
]
print(f"binary {BIN}: {subprocess.run([BIN,'--version'],capture_output=True,text=True).stdout.splitlines()[0]}")
ARMS2 = [("depth0 a.test.ts", {"op": "write", "path": "a.test.ts", "content": "P"}), ("depth1 src/foo.test.ts", {"op": "write", "path": "src/foo.test.ts", "content": "P"}), ("depth2 src/deep/x.test.ts", {"op": "write", "path": "src/deep/x.test.ts", "content": "P"}), ("nested dir pattern: cli/tests/x.py vs tests/**", {"op": "write", "path": "cli/tests/x.py", "content": "P"})]
for label, pats, arms in [("WITHOUT protected patterns (control)", [], ARMS), ("WITH " + json.dumps(PATTERNS), PATTERNS, ARMS), ("ONLY [\"**/*.test.ts\", \"tests/**\"] (pattern-form semantics)", ["**/*.test.ts", "tests/**"], ARMS2)]:
    print(f"\n=== fs_deny_write {label}")
    for name, req in arms:
        wt, pol = setup(pats)
        if "path" in req: req = dict(req, path=req["path"].replace("@WT@", wt))
        rc, out, err = call(wt, pol, req)
        # host evidence (APFS is case-insensitive: case arms land on the lowercase name)
        ev = []
        for rel in ["pkg/a.go", "pkg/a_test.go", "pkg/b_test.go", "tests/new.py", "cli/sub/conftest.py", "src/deep/x.test.ts",
                    "tests/new2.py", "tests/new3.py", "tests/new4.py", "tests/new5.py", "tests/test_x.py", "a.test.ts", "cli/tests/x.py"]:
            p = os.path.join(wt, rel)
            if os.path.exists(p):
                b = open(p).read()
                if b != FILES.get(rel): ev.append(f"{rel}:{'NEW' if rel not in FILES else 'CHANGED'}")
        print(f"{name:42s} rc={rc} {out}  host-changes={ev or 'none'}")
shutil.rmtree(SCR)
