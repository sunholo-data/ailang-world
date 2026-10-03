#!/usr/bin/env python3
"""P13 (quorum round 1 premises):
 (a) an allowRead entry that is an ANCESTOR of a denyRead entry re-opens it (glm's premise);
 (b) a live-listener network probe: an unsandboxed control connect is accepted, the sandboxed connect is
     refused, and the listener logs exactly the control accept (kimi's premise)."""
import json, os, socket, subprocess, threading
S = os.path.dirname(os.path.abspath(__file__))
HOME = os.path.expanduser("~")
SRT = ["/opt/homebrew/bin/node", HOME + "/dev/sunholo-data/ailang/tools/pi-extensions/sandbox/node_modules/@anthropic-ai/sandbox-runtime/dist/cli.js"]
H = HOME + "/.srt-row140-scratch3"
subprocess.run(["rm", "-rf", H]); os.makedirs(H + "/ws"); os.chdir(H + "/ws")

def settings(name, deny, allow):
    p = os.path.join(S, name)
    json.dump({"filesystem": {"denyRead": deny, "allowRead": allow, "allowWrite": [H + "/ws"], "denyWrite": []},
               "network": {"allowedDomains": [], "deniedDomains": []}}, open(p, "w"))
    return p

def run(cfg, argv):
    r = subprocess.run(SRT + ["--settings", cfg, "--"] + argv, capture_output=True, text=True)
    return r.returncode, (r.stdout + r.stderr).strip()[:160]

print("== (a) allowRead ancestor of a denied root")
decoy = HOME + "/.srt-row140-decoy"
print("narrow  allow=[ws]        :", run(settings("s-a1.json", [HOME], [H + "/ws"]), ["cat", decoy]))
print("ancestor allow=[ws, HOME] :", run(settings("s-a2.json", [HOME], [H + "/ws", HOME]), ["cat", decoy]))
print("ancestor allow=[ws, /Users]:", run(settings("s-a3.json", [HOME], [H + "/ws", "/Users"]), ["cat", decoy]))
link = H + "/homelink"; os.symlink(HOME, link)
print("symlink  allow=[ws, link->HOME]:", run(settings("s-a4.json", [HOME], [H + "/ws", link]), ["cat", decoy]))

print("== (b) live-listener network probe")
srv = socket.socket(); srv.bind(("127.0.0.1", 0)); srv.listen(8); port = srv.getsockname()[1]
accepts = []
def acceptor():
    while True:
        try:
            c, a = srv.accept(); accepts.append(a); c.close()
        except OSError:
            return
threading.Thread(target=acceptor, daemon=True).start()
ctl = socket.create_connection(("127.0.0.1", port), 3); ctl.close()
import time; time.sleep(0.2)
print(f"listener 127.0.0.1:{port}; unsandboxed control accepts={len(accepts)}")
cfg = settings("s-b.json", [], [])
code = f"import socket\ntry:\n socket.create_connection(('127.0.0.1',{port}),3); print('CONNECTED')\nexcept OSError as e:\n print('REFUSED', e.errno, e.strerror)"
print("sandboxed python connect:", run(cfg, ["python3", "-c", code]))
print("sandboxed curl (proxy env):", run(cfg, ["curl", "-sS", "-o", "/dev/null", "-w", "%{http_code}", "--max-time", "3", f"http://127.0.0.1:{port}/"]))
time.sleep(0.3)
print(f"listener accepts after the sandboxed arms = {len(accepts)} (1 = only the control)")
srv.close()
subprocess.run(["rm", "-rf", H])
