#!/usr/bin/env python3
"""P15c: the revised trampoline `/bin/sh -c '/usr/bin/env -i "$@"; exit $?' world-exec K=V... argv0 args...`
keeps the signal mapping (no `exec`), sets the child env exactly, and passes argv verbatim; host node env = PATH only."""
import os, subprocess
S = os.path.dirname(os.path.abspath(__file__))
SRT = ["/opt/homebrew/bin/node", os.path.expanduser("~/dev/sunholo-data/ailang/tools/pi-extensions/sandbox/node_modules/@anthropic-ai/sandbox-runtime/dist/cli.js"),
       "--settings", os.path.join(S, "s-basic.json"), "--"]
TRAMP = ["/bin/sh", "-c", '/usr/bin/env -i "$@"; exit $?', "world-exec", "HOME=/nonexistent-home", "TMPDIR=" + os.path.join(S, "ws"), "LANG=C.UTF-8"]
os.chdir(os.path.join(S, "ws"))
def run(argv):
    r = subprocess.run(SRT + TRAMP + argv, capture_output=True, text=True, env={"PATH": "/usr/bin:/bin"})
    return r.returncode, (r.stdout + r.stderr).strip().replace("\n", " | ")[:200]
for sig in ["TERM", "INT", "KILL"]:
    print(f"SIG{sig}:", run(["/bin/sh", "-c", f"kill -{sig} $$"]))
print("exit 3:", run(["/bin/sh", "-c", "exit 3"]))
print("env seen:", run(["/usr/bin/env"]))
print("argv:", run(["/usr/bin/printf", "[%s]", "a b", "$(id)", "-c", "--", "X=Y"]))
