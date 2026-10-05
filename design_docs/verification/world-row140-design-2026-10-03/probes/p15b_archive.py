#!/usr/bin/env python3
"""P15b (quorum round 3, kimi): srt run from an archived copy of its WHOLE install tree (srt + its runtime
deps) placed under a read-DENIED state dir; and the per-call digest cost of that whole tree."""
import glob, hashlib, json, os, shutil, subprocess, time
S = os.path.dirname(os.path.abspath(__file__))
HOME = os.path.expanduser("~")
NM = HOME + "/dev/sunholo-data/ailang/tools/pi-extensions/sandbox/node_modules"
H = HOME + "/.srt-row140-scratch5"
shutil.rmtree(H, ignore_errors=True)
os.makedirs(H + "/root/ep1"); os.makedirs(H + "/state/archive")
shutil.copytree(NM, H + "/state/archive/node_modules", symlinks=True)
cfg = os.path.join(S, "s-r3b.json")
json.dump({"filesystem": {"denyRead": [HOME, H + "/state", H + "/root"], "allowRead": [H + "/root/ep1"],
                          "allowWrite": [H + "/root/ep1"], "denyWrite": []},
           "network": {"allowedDomains": [], "deniedDomains": []}}, open(cfg, "w"))
os.chdir(H + "/root/ep1")
cli = H + "/state/archive/node_modules/@anthropic-ai/sandbox-runtime/dist/cli.js"
r = subprocess.run(["/opt/homebrew/bin/node", cli, "--settings", cfg, "--", "/bin/sh", "-c",
                    "echo ran-from-archive; cat " + cli + " >/dev/null && echo CHILD-CAN-READ-ARCHIVE || echo CHILD-CANNOT-READ-ARCHIVE"],
                   capture_output=True, text=True, env={"PATH": "/usr/bin:/bin"})
print("run from archive under denied state dir: rc", r.returncode, "|", (r.stdout + r.stderr).strip().replace("\n", " | ")[:300])
files = sorted(p for p in glob.glob(NM + "/**", recursive=True) if os.path.isfile(p))
pkgs = sorted(os.listdir(NM))
t = time.perf_counter(); h = hashlib.sha256()
for p in files:
    h.update(os.path.relpath(p, NM).encode() + b"\0" + hashlib.sha256(open(p, "rb").read()).digest())
print(f"whole install tree: top-level entries {pkgs}; {len(files)} files, {sum(os.path.getsize(p) for p in files)} bytes; digest {1000*(time.perf_counter()-t):.1f} ms")
shutil.rmtree(H)
