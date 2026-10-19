#!/usr/bin/env python3
"""V23 [r4.1]: standing ignored protected files in the REAL checkouts (read-only: `git ls-files -o -i` writes nothing;
GIT_OPTIONAL_LOCKS=0), with r4 candidate pathspecs, filtered by an ASCII stand-in for matchProtected (all patterns here are
ASCII; fold = upper()). Pre-fix = entries r4 would report on a pristine checkout; post-fix = entries absent from an
episode-start snapshot taken by the same command. Also times a content identity (lstat size+mtime; sha256) of the snapshot."""
import fnmatch, hashlib, os, subprocess, sys, time
HOME = sys.argv[1]
env = dict(os.environ, HOME=HOME, GIT_CONFIG_NOSYSTEM="1", GIT_CONFIG_GLOBAL="/dev/null", GIT_OPTIONAL_LOCKS="0")
def W(s):  # §4.2 widening for the ASCII patterns used here
    out = "".join("*" if c in "sSkK;`?" else c for c in s)
    while "**" in out: out = out.replace("**", "*")
    return out
def specs(pats):
    r = []
    for p in pats:
        if p.endswith("/**"):
            d = "/".join(W(x) for x in p[:-3].split("/")); r += [f":(top,glob,icase){d}", f":(top,glob,icase){d}/**"]
        else: r.append(f":(top,glob,icase)**/{W(p)}")
    return r
def mirror(pats, path):
    k = path.upper(); b = k.rsplit("/", 1)[-1]
    for p in pats:
        if p.endswith("/**"):
            d = p[:-3].upper()
            if k == d or k.startswith(d + "/"): return True
        elif fnmatch.fnmatchcase(b, p.upper()): return True
    return False
REPOS = {
 "sunholo-platform": ("/Users/voightkampff/dev/sunholo-data/sunholo-platform",
   ["test_*.py", "conftest.py", "*.test.ts", "*.test.tsx", "backend/tests/**", "cli/tests/**"]),
 "TwilightGame": ("/Users/voightkampff/dev/markedmondson1234/TwilightGame", ["tests/**", "*.test.ts", "*.test.tsx", "vitest.config.ts"]),
 "ailang": ("/Users/voightkampff/dev/sunholo-data/ailang", ["*_test.go", "*.golden", "tests/**", "examples/tests/**"]),
}
for name, (root, pats) in REPOS.items():
    head = subprocess.run(["git", "rev-parse", "--short", "HEAD"], cwd=root, env=env, capture_output=True, text=True).stdout.strip()
    cmd = ["git", "-c", "core.fsmonitor=false", "ls-files", "-z", "-o", "-i", "--exclude-standard", "--"] + specs(pats)
    ts = []
    for _ in range(3):
        t = time.perf_counter(); p = subprocess.run(cmd, cwd=root, env=env, capture_output=True); ts.append((time.perf_counter() - t) * 1000)
        assert p.returncode == 0, p.stderr
    cand = [x.decode() for x in p.stdout.split(b"\0") if x]
    hit = [c for c in cand if mirror(pats, c)]
    loc = {}
    for c in hit:
        tag = ".venv" if "/.venv/" in "/" + c else ("node_modules" if "node_modules/" in c else "other")
        loc[tag] = loc.get(tag, 0) + 1
    snap = set(c.upper() for c in hit)
    p2 = subprocess.run(cmd, cwd=root, env=env, capture_output=True)
    post = [x.decode() for x in p2.stdout.split(b"\0") if x and mirror(pats, x.decode()) and x.decode().upper() not in snap]
    t = time.perf_counter(); st = [os.lstat(os.path.join(root, c)) for c in hit]; t_stat = (time.perf_counter() - t) * 1000
    t = time.perf_counter(); nbytes = 0
    for c in hit:
        fp = os.path.join(root, c)
        if os.path.isfile(fp) and not os.path.islink(fp):
            with open(fp, "rb") as fh: d = fh.read(); nbytes += len(d); hashlib.sha256(d).digest()
    t_hash = (time.perf_counter() - t) * 1000
    print(f"{name}@{head} | ls-files -o -i ms {min(ts):.0f}-{max(ts):.0f} | candidates={len(cand)} mirror-true(pre-fix reported)={len(hit)} by-location={loc} | post-fix reported={len(post)} | lstat {t_stat:.1f} ms, sha256 {t_hash:.1f} ms over {nbytes} bytes")
