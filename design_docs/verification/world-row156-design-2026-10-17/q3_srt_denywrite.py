#!/usr/bin/env python3
"""Q3: how srt 0.0.78 (macOS) enforces denyWrite for protected test paths.
Two renderings: GLOB (Seatbelt globs) and LITERAL (render-time expansion of the
files that exist, the only form Linux srt honours for non-/** patterns).
Each arm runs in a FRESH tree; host evidence is read unsandboxed afterwards."""
import hashlib, json, os, shutil, subprocess, sys
SCR = "/Users/voightkampff/dev/sunholo-data/.scratch-iter277-row156/q3"
SRT = os.path.expanduser("~/.ailang-world-exec/srt/node_modules/@anthropic-ai/sandbox-runtime/dist/cli.js")
NODE = shutil.which("node")
FILES = {"pkg/a.go": "package pkg\n", "pkg/a_test.go": "package pkg\n// ORIG\n", "a_test.go": "package root\n// ORIG\n",
         "tests/test_x.py": "# ORIG\n", "cli/conftest.py": "# ORIG\n", "src/foo.test.ts": "// ORIG\n"}
def tree():
    if os.path.exists(SCR): shutil.rmtree(SCR)
    wt = os.path.join(SCR, "wt"); 
    for rel, body in FILES.items():
        p = os.path.join(wt, rel); os.makedirs(os.path.dirname(p), exist_ok=True); open(p, "w").write(body)
    return os.path.realpath(wt)
def settings(wt, mode):
    if mode == "GLOB":
        deny = [f"{wt}/tests", f"{wt}/**/*_test.go", f"{wt}/**/conftest.py", f"{wt}/**/*.test.ts"]
    elif mode == "LITERAL":
        deny = [f"{wt}/tests", f"{wt}/pkg/a_test.go", f"{wt}/a_test.go", f"{wt}/cli/conftest.py", f"{wt}/src/foo.test.ts"]
    else:
        deny = []
    s = {"filesystem": {"denyRead": [], "allowRead": [], "allowWrite": [wt], "denyWrite": deny},
         "network": {"allowedDomains": [], "deniedDomains": [], "allowLocalBinding": False, "allowAllUnixSockets": False},
         "enableWeakerNetworkIsolation": False, "allowAppleEvents": False}
    p = os.path.join(SCR, "s.json"); open(p, "w").write(json.dumps(s)); os.chmod(p, 0o600); return p
ARMS = [
 ("control: write pkg/a.go (non-protected)", "echo X >> pkg/a.go", ["pkg/a.go"]),
 ("modify existing pkg/a_test.go", "echo PWN >> pkg/a_test.go", ["pkg/a_test.go"]),
 ("create NEW pkg/b_test.go", "echo PWN > pkg/b_test.go", ["pkg/b_test.go"]),
 ("create NEW root c_test.go (** = zero dirs)", "echo PWN > c_test.go", ["c_test.go"]),
 ("create NEW tests/new.py (dir deny)", "echo PWN > tests/new.py", ["tests/new.py"]),
 ("create NEW cli/sub/conftest.py", "mkdir -p cli/sub && echo PWN > cli/sub/conftest.py", ["cli/sub/conftest.py"]),
 ("create NEW src/bar.test.ts", "echo PWN > src/bar.test.ts", ["src/bar.test.ts"]),
 ("rm existing pkg/a_test.go", "rm -f pkg/a_test.go", ["pkg/a_test.go"]),
 ("mv pkg/a.go onto NEW pkg/z_test.go", "mv pkg/a.go pkg/z_test.go", ["pkg/z_test.go"]),
 ("mv pkg/a_test.go away (rename source)", "mv pkg/a_test.go pkg/moved.txt", ["pkg/a_test.go", "pkg/moved.txt"]),
 ("rename dir tests -> tests2, write inside", "mv tests tests2 && echo PWN > tests2/test_x.py", ["tests/test_x.py", "tests2/test_x.py"]),
 ("rename dir pkg -> pkg2, write pkg2/a_test.go", "mv pkg pkg2 && echo PWN >> pkg2/a_test.go", ["pkg/a_test.go", "pkg2/a_test.go"]),
 ("symlink lnk->tests, write lnk/x.py", "ln -s tests lnk && echo PWN > lnk/x.py", ["tests/x.py"]),
 ("hardlink hl->pkg/a_test.go, write hl", "ln pkg/a_test.go hl && echo PWN >> hl", ["pkg/a_test.go", "hl"]),
 ("case: append PKG/A_TEST.GO", "echo PWN >> PKG/A_TEST.GO", ["pkg/a_test.go"]),
 ("case: create NEW pkg/B_TEST.GO", "echo PWN > pkg/B_TEST.GO", ["pkg/B_TEST.GO"]),
]
def snap(wt, rels):
    out = []
    for r in rels:
        p = os.path.join(wt, r)
        if os.path.islink(p): out.append(f"{r}=LINK"); continue
        if os.path.exists(p):
            b = open(p, "rb").read(); out.append(f"{r}={'ORIG' if b == FILES.get(r, '').encode() else 'CHANGED:'+b.decode(errors='replace').strip().replace(chr(10),'|')[:30]}")
        else: out.append(f"{r}=ABSENT")
    return " ".join(out)
print(f"srt cli.js sha256 {hashlib.sha256(open(SRT,'rb').read()).hexdigest()}  node {NODE}  platform {sys.platform}")
for mode in ["NONE", "GLOB", "LITERAL"]:
    print(f"\n=== rendering {mode}")
    for name, cmd, rels in ARMS:
        wt = tree(); s = settings(wt, mode)
        r = subprocess.run([NODE, SRT, "--settings", s, "--", "/bin/sh", "-c", cmd], cwd=wt, capture_output=True, text=True, timeout=30)
        err = (r.stderr.strip().splitlines() or [""])[-1][:90]
        print(f"[{mode}] {name:48s} rc={r.returncode} host: {snap(wt, rels)}  stderr: {err}")
shutil.rmtree(SCR)
