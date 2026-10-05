#!/usr/bin/env python3
"""Row 140 M0 — the srt confinement matrix, portable (macOS Seatbelt and Linux bubblewrap).

Runs the design's startup-probe arms 1-6 plus V3/V8/V9/V11/V12/V13/V41/V42 against ONE pinned srt,
in World's real topology (workspace root + state dir under $HOME), with srt run from an archived
copy of its whole install tree, launched exactly as design §4.5 says:

  <node> <archive>/…/cli.js --settings <rendered> -- /bin/sh -c '/usr/bin/env -i "$@"; exit $?' world-exec K=V… argv…

Host env for srt's node is PATH only (§4.5 r3). Every "refused" arm demands a positive token from
the probe client (World's own node); every arm has a control that proves it can fail.

usage: m0_matrix.py --srt-node-modules DIR --node ABS_NODE [--json OUT]
exit 0 iff every GATING arm passes; INFO arms are recorded, never gating.
"""
import argparse, hashlib, json, os, platform, shutil, signal, socket, subprocess, sys, tempfile, threading, time

ap = argparse.ArgumentParser()
ap.add_argument("--srt-node-modules", required=True, help="a node_modules dir holding @anthropic-ai/sandbox-runtime and its deps")
ap.add_argument("--node", required=True, help="absolute path of the node binary World would verify")
ap.add_argument("--json", help="write the results here")
ap.add_argument("--seccomp-read", choices=["none", "allowread"], default="none",
                help="Linux: how srt's in-sandbox apply-seccomp helper (archived under the read-denied state dir) is reached. "
                     "none = design §4.6 as written; allowread = re-allow ONLY <archive>/…/vendor/seccomp/<arch> for reading")
a = ap.parse_args()

NODE = os.path.realpath(a.node)
HOME = os.path.realpath(os.path.expanduser("~"))
HOST_ENV = {"PATH": "/usr/bin:/bin"}
results = []

def rec(arm, gating, expected, observed, ok):
    results.append({"arm": arm, "gating": gating, "expected": expected, "observed": observed, "pass": bool(ok)})
    print(f"[{'PASS' if ok else 'FAIL'}]{'' if gating else ' (info)'} {arm}: expected {expected!r}; observed {observed!r}", flush=True)

# --- topology (§4.4, V35): everything under a scratch dir inside $HOME, as World's real layout is ---
BASE = tempfile.mkdtemp(prefix=".world-row140-m0-", dir=HOME)
ROOT = os.path.join(BASE, "workspace")          # workspace root (denyRead)
EP1 = os.path.join(ROOT, "ep1")                  # this episode's worktree
EP2 = os.path.join(ROOT, "ep2")                  # a sibling episode
STATE = os.path.join(BASE, "state")              # World state dir (denyRead)
CACHE = os.path.join(STATE, "exec-cache", "ep1", "proj")
for d in (EP1, EP2, os.path.join(CACHE, "home"), os.path.join(CACHE, "tmp"), os.path.join(STATE, "exec-probe"),
          os.path.join(ROOT, ".exec-probe-sibling")):
    os.makedirs(d, exist_ok=True)
with open(os.path.join(EP1, ".git"), "w") as f:  # a worktree pointer FILE, as World's episodes have
    f.write("gitdir: /nonexistent/.git/worktrees/ep1\n")
with open(os.path.join(EP2, "notes.txt"), "w") as f:
    f.write("SIBLING-EPISODE-FILE\n")
DECOYS = [os.path.join(HOME, ".ailang-worldd-exec-probe-decoy-m0"),
          os.path.join(STATE, "exec-probe", "decoy"),
          os.path.join(ROOT, ".exec-probe-sibling", "decoy")]
for p in DECOYS:
    with open(p, "w") as f:
        f.write("DECOY-NOT-A-SECRET\n")

# --- archive the WHOLE install tree under the read-denied state dir (§4.6 r3, V42) ---
SRC_NM = os.path.realpath(a.srt_node_modules)
ARCH_NM = os.path.join(STATE, "archive", "exec-sandbox", "node_modules")
shutil.copytree(SRC_NM, ARCH_NM, symlinks=True)
CLI = os.path.join(ARCH_NM, "@anthropic-ai", "sandbox-runtime", "dist", "cli.js")
PKG_ONLY = os.path.join(STATE, "archive", "pkg-only", "sandbox-runtime")
shutil.copytree(os.path.join(SRC_NM, "@anthropic-ai", "sandbox-runtime"), PKG_ONLY, symlinks=True)

def sha(p):
    return hashlib.sha256(open(p, "rb").read()).hexdigest()

def ver(cmd):
    try:
        r = subprocess.run(cmd, capture_output=True, text=True, timeout=10)
        return (r.stdout + r.stderr).strip().splitlines()[0] if (r.stdout + r.stderr).strip() else f"rc={r.returncode}"
    except Exception as e:
        return f"absent ({type(e).__name__})"

meta = {
    "platform": platform.platform(), "machine": platform.machine(),
    "srt_version": json.load(open(os.path.join(ARCH_NM, "@anthropic-ai", "sandbox-runtime", "package.json")))["version"],
    "cli_sha256": sha(CLI), "node": NODE, "node_version": ver([NODE, "--version"]),
    "bwrap": ver(["bwrap", "--version"]), "socat": ver(["socat", "-V"]) if sys.platform != "darwin" else "n/a",
    "rg": ver(["rg", "--version"]),
    "userns_sysctl": open("/proc/sys/kernel/apparmor_restrict_unprivileged_userns").read().strip()
        if os.path.exists("/proc/sys/kernel/apparmor_restrict_unprivileged_userns") else "n/a",
}
print("META", json.dumps(meta), flush=True)

DENY_WRITE = [os.path.join(EP1, n) for n in (".git", ".github", ".pi", ".claude", ".ailang", ".gitmodules", ".gitattributes")] + \
             ["/tmp/claude", "/private/tmp/claude", os.path.join(HOME, ".npm", "_logs"), os.path.join(HOME, ".claude", "debug")]
if sys.platform != "darwin":
    DENY_WRITE = [p for p in DENY_WRITE if not p.startswith("/private/")]

SECCOMP_DIR = os.path.join(ARCH_NM, "@anthropic-ai", "sandbox-runtime", "vendor", "seccomp",
                           {"x86_64": "x64", "aarch64": "arm64", "arm64": "arm64"}.get(platform.machine(), platform.machine()))
EXTRA_READ = [SECCOMP_DIR] if (a.seccomp_read == "allowread" and sys.platform != "darwin") else []
meta["seccomp_read"] = a.seccomp_read
meta["extra_allow_read"] = EXTRA_READ

def settings(name, **over):
    s = {
        "filesystem": {"denyRead": [HOME, STATE, ROOT], "allowRead": [EP1, CACHE] + EXTRA_READ,
                       "allowWrite": [EP1, CACHE], "denyWrite": DENY_WRITE},
        "network": {"allowedDomains": [], "deniedDomains": [], "allowLocalBinding": False, "allowAllUnixSockets": False},
        "enableWeakerNetworkIsolation": False, "allowAppleEvents": False,
    }
    for k, v in over.items():
        sec, key = k.split("__")
        s[sec][key] = v
    p = os.path.join(BASE, name + ".srt.json")
    with open(p, "w") as f:
        json.dump(s, f)
    os.chmod(p, 0o600)
    return p

SET = settings("rendered")
SET_NODENYW = settings("no-denywrite", filesystem__denyWrite=[])
KV = [f"HOME={CACHE}/home", f"TMPDIR={CACHE}/tmp", "PATH=/usr/bin:/bin", "LANG=C.UTF-8", "LC_ALL=C.UTF-8", "CI=1"]
TRAMP = ["/bin/sh", "-c", '/usr/bin/env -i "$@"; exit $?', "world-exec"]

def srt(argv, st=SET, tramp=True, dashdash=True, cli=CLI, timeout=60):
    cmd = [NODE, cli, "--settings", st] + (["--"] if dashdash else []) + ((TRAMP + KV) if tramp else []) + argv
    try:
        r = subprocess.run(cmd, cwd=EP1, env=HOST_ENV, capture_output=True, text=True, timeout=timeout, stdin=subprocess.DEVNULL)
        return r.returncode, (r.stdout + r.stderr).strip()
    except subprocess.TimeoutExpired:
        return "timeout", ""

def js(code):
    return [NODE, "-e", code]

def readp(p):
    return js(f"try{{require('fs').readFileSync({json.dumps(p)});console.log('WORLD-PROBE-READ')}}catch(e){{console.log('WORLD-PROBE-REFUSED '+e.code)}}")

def writep(p):
    return js(f"try{{require('fs').writeFileSync({json.dumps(p)},'x');console.log('WORLD-PROBE-WROTE')}}catch(e){{console.log('WORLD-PROBE-REFUSED '+e.code)}}")

def tok(out):
    for line in out.splitlines():
        if line.startswith("WORLD-PROBE-"):
            return line.strip()
    return None

def unsbx(argv):
    r = subprocess.run(argv, capture_output=True, text=True, timeout=30)
    return tok(r.stdout)

REFUSED = ("WORLD-PROBE-REFUSED EPERM", "WORLD-PROBE-REFUSED EACCES", "WORLD-PROBE-REFUSED EROFS", "WORLD-PROBE-REFUSED ENOENT")
# ENOENT counts as refused ONLY for reads/writes into a denied region whose control proves the path exists
# (bwrap may mask a denied dir with an empty tmpfs rather than return EPERM); recorded verbatim either way.

try:
    # ---- sanity: the stack runs at all, from the archive, with PATH-only host env ----
    rc, out = srt(["/bin/echo", "M0-ALIVE"])
    rec("sanity: srt from whole-tree archive runs (V42)", True, "rc 0 + M0-ALIVE", f"rc {rc}: {out[:300]}", rc == 0 and "M0-ALIVE" in out)
    rc, out = srt(["/bin/echo", "x"], cli=os.path.join(PKG_ONLY, "dist", "cli.js"))
    rec("V42 package-dir-only archive fails to load", False, "non-zero, module not found", f"rc {rc}: {out[:200]}", rc != 0)

    # ---- arm 1: write inside ----
    target = os.path.join(EP1, "in.txt")
    rc, out = srt(writep(target))
    back = open(target).read() if os.path.exists(target) else None
    rec("arm1 write inside worktree", True, "WORLD-PROBE-WROTE + readback 'x'", f"{tok(out)} readback={back!r}", tok(out) == "WORLD-PROBE-WROTE" and back == "x")
    rc, out = srt(writep(os.path.join(CACHE, "tmp", "c.txt")))
    rec("arm1b write inside episode exec cache", True, "WORLD-PROBE-WROTE", tok(out), tok(out) == "WORLD-PROBE-WROTE")

    # ---- arm 2: write outside (sibling episode), with control ----
    ctl = unsbx(writep(os.path.join(EP2, "control.txt")))
    tgt = os.path.join(EP2, "pwn.txt")
    rc, out = srt(writep(tgt))
    rec("arm2 write sibling episode refused", True, "control WROTE; sandboxed REFUSED; file absent",
        f"control={ctl}; sbx={tok(out)}; exists={os.path.exists(tgt)}",
        ctl == "WORLD-PROBE-WROTE" and tok(out) in REFUSED and not os.path.exists(tgt))
    for name, p in (("$HOME", os.path.join(HOME, ".world-row140-m0-home-write")), ("state dir", os.path.join(STATE, "w.txt"))):
        rc, out = srt(writep(p))
        rec(f"arm2 write {name} refused", True, "REFUSED; absent", f"{tok(out)}; exists={os.path.exists(p)}", tok(out) in REFUSED and not os.path.exists(p))
        if os.path.exists(p):
            os.remove(p)

    # ---- arm 3: read fence, fixed decoys, each with an unsandboxed control ----
    for p in DECOYS:
        ctl = unsbx(readp(p))
        rc, out = srt(readp(p))
        rec(f"arm3 read decoy {p.replace(HOME, '$HOME')} refused", True, "control READ; sandboxed REFUSED", f"control={ctl}; sbx={tok(out)}",
            ctl == "WORLD-PROBE-READ" and tok(out) in REFUSED)
    rc, out = srt(readp(os.path.join(EP2, "notes.txt")))
    rec("arm3b read sibling episode file refused (V35)", True, "REFUSED", tok(out), tok(out) in REFUSED)
    rc, out = srt(readp(CLI))
    rec("V42 child cannot read the archived cli.js", True, "REFUSED", tok(out), tok(out) in REFUSED)
    rc, out = srt(["/bin/ls", HOME])
    rec("V8 ls $HOME refused", True, "non-zero, no listing of BASE", f"rc {rc}: {out[:160]}", rc != 0 and os.path.basename(BASE) not in out)
    rc, out = srt(readp("/etc/hosts"))
    rec("V8 system path /etc/hosts stays readable (R-140-1)", False, "READ", tok(out), tok(out) == "WORLD-PROBE-READ")
    rc, out = srt(readp(os.path.join(EP1, "in.txt")))
    rec("V8 own worktree readable under denied $HOME", True, "READ", tok(out), tok(out) == "WORLD-PROBE-READ")

    # ---- arm 4: srt's always-writable defaults closed by denyWrite ----
    rc, out = srt(writep("/tmp/claude/m0-probe"))
    rec("arm4 /tmp/claude write refused (rendered denyWrite)", True, "REFUSED", tok(out), tok(out) in REFUSED)
    rc, out = srt(writep("/tmp/claude/m0-probe-nodeny"), st=SET_NODENYW)
    rec("V4 /tmp/claude writable WITHOUT World's denyWrite", False, "WROTE (srt default)", tok(out), tok(out) == "WORLD-PROBE-WROTE")
    # the dirs exist (made unsandboxed) so a refusal cannot be a vacuous ENOENT; on a case-insensitive
    # fs .GITHUB IS .github, on Linux it is a different, non-special dir (so not gating there)
    os.makedirs(os.path.join(EP1, ".github"), exist_ok=True)
    os.makedirs(os.path.join(EP1, ".GITHUB"), exist_ok=True)
    for p in (os.path.join(EP1, ".git"), os.path.join(EP1, ".GITHUB", "x"), os.path.join(EP1, ".github", "x")):
        before = open(p).read() if os.path.isfile(p) else None
        rc, out = srt(writep(p))
        after = open(p).read() if os.path.isfile(p) else None
        changed = (before != after) if before is not None else (after is not None)
        rec(f"V6/V35 write {os.path.relpath(p, EP1)} refused", sys.platform == "darwin" or not p.endswith("GITHUB/x"),
            "REFUSED; bytes unchanged", f"{tok(out)}; changed={changed}", tok(out) in REFUSED and not changed)

    # ---- arm 5: exit status through the trampoline (V11, V41) ----
    for sig, want in (("TERM", 143), ("INT", 130), ("KILL", 137)):
        rc, out = srt(["/bin/sh", "-c", f"kill -{sig} $$"])
        rec(f"arm5 SIG{sig} self-kill -> {want} via trampoline", True, want, rc, rc == want)
    rc, out = srt(["/bin/sh", "-c", "exit 3"])
    rec("arm5 exit 3 -> 3", True, 3, rc, rc == 3)
    rc, out = srt(["/nonexistent/binary"])
    rec("V11 missing binary -> 127", True, 127, rc, rc == 127)
    rc, out = srt(["/bin/sh", "-c", "kill -TERM $$"], tramp=False)
    rec("V11 bare srt (no trampoline) SIGTERM rc", False, "0 on macOS 0.0.71/0.0.78 (the hazard)", rc, True)

    # ---- child env exactly the K=V set (V41) ----
    rc, out = srt(["/usr/bin/env"])
    got = sorted(l for l in out.splitlines() if "=" in l)
    rec("V41 child env == rendered K=V set (no proxy vars, no TMPDIR=/tmp/claude)", True, sorted(KV), got, got == sorted(KV))

    # ---- argv verbatim and the `--` rule (V12) ----
    rc, out = srt(["/usr/bin/printf", "[%s]", "a b", "$(id)", "-c", "--", "X=Y", "-s", "/nonexistent.json"])
    want = "[a b][$(id)][-c][--][X=Y][-s][/nonexistent.json]"
    rec("V12 argv verbatim after -- (incl. -c, -s, X=Y)", True, want, out, out == want)
    rc, out = srt(["/usr/bin/printf", "[%s]", "-c", "echo OPTION-INJECTED"], tramp=False, dashdash=False)
    rec("V12 WITHOUT -- srt takes -c (the hazard)", False, "OPTION-INJECTED executed", f"rc {rc}: {out[:160]}", True)

    # ---- V3: grandchild, orphan, symlink writes outside ----
    g = os.path.join(EP2, "grandchild.txt")
    rc, out = srt(["/bin/sh", "-c", f'echo WORLD-PROBE-RAN; sh -c "sh -c \\"echo z > {g}\\""'])
    rec("V3 grandchild write outside refused", True, "ran; absent", f"{tok(out)}; exists={os.path.exists(g)}", tok(out) == "WORLD-PROBE-RAN" and not os.path.exists(g))
    o = os.path.join(EP2, "orphan.txt")
    rc, out = srt(["/bin/sh", "-c", f'( nohup sh -c "sleep 2; echo late > {o}" >/dev/null 2>&1 & ); echo WORLD-PROBE-RAN'])
    time.sleep(3)
    rec("V3 orphan write (2 s after srt exits) refused", True, "ran; absent", f"{tok(out)}; exists={os.path.exists(o)}", tok(out) == "WORLD-PROBE-RAN" and not os.path.exists(o))
    os.symlink(EP2, os.path.join(EP1, "lnk"))
    s = os.path.join(EP2, "viaslink.txt")
    rc, out = srt(["/bin/sh", "-c", "echo WORLD-PROBE-RAN; echo y > lnk/viaslink.txt"])
    rec("V3 write through symlink to sibling refused", True, "ran; absent", f"{tok(out)}; exists={os.path.exists(s)}", tok(out) == "WORLD-PROBE-RAN" and not os.path.exists(s))

    # ---- arm 6: network against a LIVE listener with accept count (V37, V38) ----
    srv = socket.socket(); srv.bind(("127.0.0.1", 0)); srv.listen(16); port = srv.getsockname()[1]
    accepts = []
    def acceptor():
        while True:
            try:
                c, _ = srv.accept(); accepts.append(1); c.close()
            except OSError:
                return
    threading.Thread(target=acceptor, daemon=True).start()
    conn = js(f"const n=require('net');const s=n.connect({port},'127.0.0.1');s.on('connect',()=>{{console.log('WORLD-PROBE-CONNECTED');s.destroy()}});s.on('error',e=>{{console.log('WORLD-PROBE-REFUSED '+e.code)}})")
    ctl = unsbx(conn); time.sleep(0.3); n0 = len(accepts)
    rc, out = srt(conn); time.sleep(0.3); n1 = len(accepts)
    rec("arm6 raw loopback connect refused (live listener)", True, "control CONNECTED count 1; sbx not CONNECTED; count stays 1",
        f"control={ctl} n0={n0}; sbx={tok(out)} n1={n1}", ctl == "WORLD-PROBE-CONNECTED" and n0 == 1 and tok(out) != "WORLD-PROBE-CONNECTED" and tok(out) is not None and n1 == n0)
    rc, out = srt(["/bin/sh", "-c", f"echo WORLD-PROBE-RAN; curl -sS --max-time 5 -o /dev/null -w '%{{http_code}}' http://127.0.0.1:{port}/; echo \" rc=$?\""], tramp=False)
    time.sleep(0.3); n2 = len(accepts)
    rec("arm6 loopback via srt proxy env (no trampoline) refused", True, "ran; count unchanged", f"{out[:160]} n2={n2}", "WORLD-PROBE-RAN" in out and n2 == n0)
    rc, out = srt(["/bin/sh", "-c", f"echo WORLD-PROBE-RAN; curl -sS --noproxy '*' --max-time 5 -o /dev/null http://127.0.0.1:{port}/; echo \" rc=$?\""], tramp=False)
    time.sleep(0.3); n3 = len(accepts)
    rec("arm6 loopback with --noproxy refused", True, "ran; count unchanged", f"{out[:160]} n3={n3}", "WORLD-PROBE-RAN" in out and n3 == n0)
    srv.close()
    rc, out = srt(["/bin/sh", "-c", "echo WORLD-PROBE-RAN; curl -sS --max-time 8 -o /dev/null -w '%{http_code}' https://example.com/; echo \" rc=$?\""], tramp=False)
    rec("V9 external https refused (allowedDomains=[])", True, "ran; no 200", out[:160], "WORLD-PROBE-RAN" in out and "200" not in out.split(" rc=")[0])
    bind = js("const n=require('net');const s=n.createServer();s.on('error',e=>console.log('WORLD-PROBE-REFUSED '+e.code));s.listen(0,'127.0.0.1',()=>{console.log('WORLD-PROBE-BOUND');s.close()})")
    rc, out = srt(bind)
    rec("V9 bind 127.0.0.1 refused (allowLocalBinding=false)", True, "not BOUND", tok(out), tok(out) is not None and tok(out) != "WORLD-PROBE-BOUND")
    rc, out = srt(js("require('dns').lookup('example.com',(e,a)=>console.log(e?'WORLD-PROBE-REFUSED '+e.code:'WORLD-PROBE-RESOLVED '+a))"))
    rec("V9 direct DNS", False, "REFUSED (macOS)", tok(out), True)

    # ---- V13: group kill leaves no member ----
    p = subprocess.Popen([NODE, CLI, "--settings", SET, "--"] + TRAMP + KV + ["/bin/sh", "-c", "sleep 3017 & sleep 3017 & wait"],
                         cwd=EP1, env=HOST_ENV, start_new_session=True, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
    time.sleep(2.0)
    def count():
        r = subprocess.run(["pgrep", "-f", "^sleep 3017$"], capture_output=True, text=True)
        return len(r.stdout.split())
    before = count()
    os.killpg(p.pid, signal.SIGKILL); p.wait(); time.sleep(1.0)
    after = count()
    rec("V13 killpg(srt group) leaves no sleeper", True, "before 2, after 0", f"before {before}, after {after}", before == 2 and after == 0)
    if after:
        subprocess.run(["pkill", "-9", "-f", "sleep 3017"])
finally:
    for p in DECOYS[:1]:
        try: os.remove(p)
        except OSError: pass
    shutil.rmtree(BASE, ignore_errors=True)
    for p in ("/tmp/claude/m0-probe", "/tmp/claude/m0-probe-nodeny"):
        try: os.remove(p)
        except OSError: pass

gating = [r for r in results if r["gating"]]
failed = [r["arm"] for r in gating if not r["pass"]]
summary = {"meta": meta, "results": results, "gating_total": len(gating), "gating_failed": failed}
print(f"SUMMARY gating {len(gating) - len(failed)}/{len(gating)} pass; failed: {failed}", flush=True)
if a.json:
    with open(a.json, "w") as f:
        json.dump(summary, f, indent=1)
sys.exit(1 if failed else 0)
