#!/usr/bin/env python3
"""P14 (quorum round 2 premises):
 (a) gemini: a probe client missing from the sandbox exits 126/127 — a 'fails' arm must demand a positive
     refusal token. Measure World's own node as the probe client: connect token under the sandbox, and
     what a missing client looks like.
 (b) glm: cost of re-verifying the confinement stack before every call — the srt package tree digest and
     the node binary (+ its libnode dylib) digest."""
import glob, hashlib, json, os, socket, subprocess, threading, time
S = os.path.dirname(os.path.abspath(__file__))
HOME = os.path.expanduser("~")
NODE = "/opt/homebrew/bin/node"
SRTDIR = HOME + "/dev/sunholo-data/ailang/tools/pi-extensions/sandbox/node_modules/@anthropic-ai/sandbox-runtime"
SRT = [NODE, SRTDIR + "/dist/cli.js"]
os.chdir(os.path.join(S, "ws"))
cfg = os.path.join(S, "s-b.json")

srv = socket.socket(); srv.bind(("127.0.0.1", 0)); srv.listen(8); port = srv.getsockname()[1]
accepts = []
def acceptor():
    while True:
        try:
            c, a = srv.accept(); accepts.append(a); c.close()
        except OSError:
            return
threading.Thread(target=acceptor, daemon=True).start()
js = ("const n=require('net');const s=n.connect(%d,'127.0.0.1');"
      "s.on('connect',()=>{console.log('WORLD-PROBE-CONNECTED');s.destroy()});"
      "s.on('error',e=>{console.log('WORLD-PROBE-REFUSED '+e.code)})" % port)
def sbx(argv):
    r = subprocess.run(SRT + ["--settings", cfg, "--"] + argv, capture_output=True, text=True)
    return r.returncode, (r.stdout + r.stderr).strip()[:160]
print("== (a) probe client = World's own node (absolute path)")
r = subprocess.run([NODE, "-e", js], capture_output=True, text=True); time.sleep(0.2)
print("unsandboxed control:", r.stdout.strip(), "accepts =", len(accepts))
print("sandboxed node     :", sbx([NODE, "-e", js])); time.sleep(0.2)
print("accepts after      =", len(accepts))
print("missing client     :", sbx(["/nonexistent/curl", "http://127.0.0.1:%d/" % port]))
srv.close()

print("== (b) per-call re-verification cost")
files = sorted(p for p in glob.glob(SRTDIR + "/**", recursive=True) if os.path.isfile(p))
t = time.perf_counter()
h = hashlib.sha256()
for p in files:
    h.update(os.path.relpath(p, SRTDIR).encode() + b"\0" + hashlib.sha256(open(p, "rb").read()).digest())
print(f"srt package tree: {len(files)} files, {sum(os.path.getsize(p) for p in files)} bytes, digest in {1000*(time.perf_counter()-t):.1f} ms")
real = os.path.realpath(NODE)
libs = glob.glob(os.path.dirname(os.path.dirname(real)) + "/lib/libnode*.dylib")
for p in [real] + libs:
    t = time.perf_counter(); hashlib.sha256(open(p, "rb").read()).hexdigest()
    print(f"{os.path.basename(p)}: {os.path.getsize(p)} bytes, sha256 in {1000*(time.perf_counter()-t):.1f} ms")
