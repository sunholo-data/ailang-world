#!/usr/bin/env python3
"""P8: srt exit-status trampoline, process-group kill of the whole sandboxed tree, output volume."""
import os, signal, subprocess, time
S = os.path.dirname(os.path.abspath(__file__))
SRT = ["/opt/homebrew/bin/node", os.path.expanduser("~/dev/sunholo-data/ailang/tools/pi-extensions/sandbox/node_modules/@anthropic-ai/sandbox-runtime/dist/cli.js"),
       "--settings", os.path.join(S, "s-basic.json"), "--"]
TRAMP = ["/bin/sh", "-c", '"$@"; exit $?', "world-exec"]
os.chdir(os.path.join(S, "ws"))

def rc(argv):
    p = subprocess.run(argv, capture_output=True, text=True)
    return p.returncode, p.stderr.strip()[:120]

print("== exit status: bare srt vs the sh trampoline (child killed by its own signal)")
for sig in ["TERM", "INT", "KILL", "SEGV"]:
    child = ["/bin/sh", "-c", f"kill -{sig} $$"]
    print(f"SIG{sig:5s} bare srt rc={rc(SRT + child)}   trampoline rc={rc(SRT + TRAMP + child)}")
print("plain exit 3 via trampoline:", rc(SRT + TRAMP + ["/bin/sh", "-c", "exit 3"]))
print("argv survives the trampoline:", subprocess.run(SRT + TRAMP + ["printf", "[%s]", "a b", "$(id)", "-c", "--"], capture_output=True, text=True).stdout)

print("== process-group kill of node + sandbox + a background grandchild")
p = subprocess.Popen(SRT + ["/bin/sh", "-c", "sleep 300 & sleep 300 & wait"], start_new_session=True,
                     stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
time.sleep(1.5)
pgid = os.getpgid(p.pid)
members = subprocess.run(["pgrep", "-g", str(pgid)], capture_output=True, text=True).stdout.split()
print(f"pgid={pgid} members before kill: {len(members)}")
t = time.perf_counter(); os.killpg(pgid, signal.SIGKILL); p.wait()
time.sleep(0.3)
left = subprocess.run(["pgrep", "-g", str(pgid)], capture_output=True, text=True).stdout.split()
print(f"members after killpg(SIGKILL): {len(left)}; node wait status {p.returncode}; kill->reap {1000*(time.perf_counter()-t):.0f}ms")

print("== output volume: 20 MB on stdout through srt")
t = time.perf_counter()
out = subprocess.run(SRT + ["/bin/sh", "-c", "head -c 20000000 /dev/zero | tr '\\0' x"], capture_output=True).stdout
print(f"bytes={len(out)} wall={1000*(time.perf_counter()-t):.0f}ms")
