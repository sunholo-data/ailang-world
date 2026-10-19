#!/usr/bin/env python3
"""V20/V21: cost of the r3 vs r4 (widened superset) vs unfiltered detection invocations on a
`git clone --shared` of ailang in scratch, then with a 30 000-file ignored tree added. Read-only to the source."""
import os, subprocess, sys, time, shutil
S = sys.argv[1]; SRC = sys.argv[2]
env = dict(os.environ, HOME=S + "/home", GIT_CONFIG_NOSYSTEM="1", GIT_CONFIG_GLOBAL="/dev/null", GIT_OPTIONAL_LOCKS="0")
R = S + "/ailang-clone"
shutil.rmtree(R, ignore_errors=True)
subprocess.run(["git", "clone", "-q", "--shared", SRC, R], check=True, env=env)
G = ["git", "-c", "core.fsmonitor=false", "-c", "core.quotepath=false"]
r3 = [":(top,glob,icase)**/*_test.go", ":(top,glob,icase)**/*.golden", ":(top,glob,icase)tests", ":(top,glob,icase)tests/**",
      ":(top,glob,icase)examples/tests", ":(top,glob,icase)examples/tests/**"]
r4 = [":(top,glob,icase)**/*_te*t.go", ":(top,glob,icase)**/*.golden", ":(top,glob,icase)te*t*", ":(top,glob,icase)te*t*/**",
      ":(top,glob,icase)example*/te*t*", ":(top,glob,icase)example*/te*t*/**"]
ST = ["status", "--porcelain=v2", "-z", "--branch", "--untracked-files=all", "--no-renames", "--ignore-submodules=all"]
cmds = {
 "r3 status --ignored=matching": ST + ["--ignored=matching", "--"] + r3,
 "r4 status --ignored=no": ST + ["--ignored=no", "--"] + r4,
 "r4 ls-files -o -i": ["ls-files", "-z", "-o", "-i", "--exclude-standard", "--"] + r4,
 "r4 ls-files -v": ["ls-files", "-v", "-z", "--"] + r4,
 "none status --ignored=no": ST + ["--ignored=no"],
 "none ls-files -o -i": ["ls-files", "-z", "-o", "-i", "--exclude-standard"],
 "none ls-files -v": ["ls-files", "-v", "-z"],
}
def run(label):
    for name, c in cmds.items():
        ts, n = [], 0
        for _ in range(4):
            t = time.perf_counter(); p = subprocess.run(G + c, cwd=R, env=env, capture_output=True); ts.append((time.perf_counter() - t) * 1000)
            assert p.returncode == 0, (name, p.stderr)
            n = len([x for x in p.stdout.split(b"\0") if x]); b = len(p.stdout)
        ts = sorted(ts[1:])  # first run warms
        print(f"{label} | {name} | entries={n} bytes={b} | warm ms {ts[0]:.0f}-{ts[-1]:.0f}")
print("files tracked:", len(subprocess.run(["git", "ls-files", "-z"], cwd=R, env=env, capture_output=True).stdout.split(b"\0")) - 1)
run("clean")
# ignored explosion arm: a node_modules-like tree, 30 000 files in 300 dirs, two of them *_test.go
nm = R + "/node_modules"
for d in range(300):
    os.makedirs(f"{nm}/p{d}", exist_ok=True)
    for f in range(100):
        open(f"{nm}/p{d}/f{f}.js", "w").close()
open(f"{nm}/p0/zz_test.go", "w").close(); open(f"{nm}/p1/yy_test.go", "w").close()
with open(R + "/.git/info/exclude", "a") as fh: fh.write("node_modules/\n")
run("ignored-30k")
shutil.rmtree(R)
print("removed clone")
