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

import re


def protectedScratchRoot(scratch_root, allowed_root=None):
    """Resolve and validate before any fixture creation; reject symlink escapes too."""
    allowed = os.path.realpath(allowed_root if allowed_root is not None else
                               os.environ.get("GITHUB_WORKSPACE", os.getcwd()))
    root = os.path.realpath(scratch_root)
    if root == allowed or os.path.commonpath([root, allowed]) != allowed:
        raise ValueError(f"scratch root {root!r} must be inside allowed root {allowed!r}")
    return root


def validateProtectedRows(rows, census):
    """Return None on success, otherwise a named rejection reason.

    Each row has arm, rendering, outcome (WROTE/REFUSED/PARTIAL), client_token, rc, and
    host_before/host_after maps of host-observed path -> SHA-256 (None = absent).
    Snapshots must be collected independently of client output. rc is diagnostic.
    census is the independent expected collection of (arm, rendering) pairs.
    """
    expected = set(census)
    if not expected:
        return "empty-census"
    if not rows or len(rows) != len(expected) or {(r.get("arm"), r.get("rendering")) for r in rows} != expected:
        return "census-mismatch"
    for row in rows:
        outcome = row.get("outcome")
        if outcome not in ("WROTE", "REFUSED", "PARTIAL"):
            return "invalid-outcome"
        token = row.get("client_token")
        valid = token == "WORLD-PROBE-WROTE" or (isinstance(token, str) and
                re.fullmatch(r"WORLD-PROBE-REFUSED E[A-Z0-9]+", token) is not None)
        if (not valid and (outcome != "PARTIAL" or token is not None)) or (
                outcome == "WROTE" and token != "WORLD-PROBE-WROTE") or (
                outcome == "REFUSED" and token == "WORLD-PROBE-WROTE"):
            return "missing-client"
        if not isinstance(row.get("rc"), int) and row.get("rc") != "timeout":
            return "invalid-rc"
        before, after = row.get("host_before"), row.get("host_after")
        if not isinstance(before, dict) or not before or not isinstance(after, dict) or before.keys() != after.keys():
            return "invalid-host-snapshot"
        for digest in list(before.values()) + list(after.values()):
            if digest is not None and (not isinstance(digest, str) or
                    re.fullmatch(r"(?:symlink:)?[0-9a-f]{64}", digest) is None):
                return "invalid-host-snapshot"
        if outcome == "REFUSED" and before != after:
            return "host-changed"
        if outcome == "WROTE" and before == after:
            return "wrote-unchanged"
    return None


# V5 order; targets include aliases and both sides of every move.
PROTECTED_ARMS = [
    ("control", ["appendFileSync('pkg/a.go','X')"], ["pkg/a.go"]),
    ("modify", ["appendFileSync('pkg/a_test.go','PWN')"], ["pkg/a_test.go"]),
    ("new-nested", ["writeFileSync('pkg/b_test.go','PWN')"], ["pkg/b_test.go"]),
    ("new-root", ["writeFileSync('c_test.go','PWN')"], ["c_test.go"]),
    ("dir-modify", ["writeFileSync('tests/new.py','PWN')"], ["tests/new.py"]),
    ("conftest", ["mkdirSync('cli/sub',{recursive:true})", "writeFileSync('cli/sub/conftest.py','PWN')"], ["cli/sub/conftest.py"]),
    ("ts", ["writeFileSync('src/bar.test.ts','PWN')"], ["src/bar.test.ts"]),
    ("delete", ["unlinkSync('pkg/a_test.go')"], ["pkg/a_test.go"]),
    ("rename-onto", ["renameSync('pkg/a.go','pkg/z_test.go')"], ["pkg/a.go", "pkg/z_test.go"]),
    ("rename-away", ["renameSync('pkg/a_test.go','pkg/moved.txt')"], ["pkg/a_test.go", "pkg/moved.txt"]),
    ("protected-dir-rename", ["renameSync('tests','tests2')", "writeFileSync('tests2/test_x.py','PWN')"], ["tests/test_x.py", "tests2/test_x.py"]),
    ("parent-dir-rename", ["renameSync('pkg','pkg2')", "appendFileSync('pkg2/a_test.go','PWN')"], ["pkg/a_test.go", "pkg2/a_test.go"]),
    ("symlink", ["symlinkSync('tests','lnk')", "writeFileSync('lnk/x.py','PWN')"], ["tests/x.py", "lnk"]),
    ("hardlink", ["linkSync('pkg/a_test.go','hl')", "appendFileSync('hl','PWN')"], ["pkg/a_test.go", "hl"]),
    ("case-existing", ["appendFileSync('PKG/A_TEST.GO','PWN')"], ["PKG/A_TEST.GO"]),
    ("case-new", ["writeFileSync('pkg/B_TEST.GO','PWN')"], ["pkg/B_TEST.GO"]),
]
PROTECTED_RENDERINGS = ("NONE", "DIRS_ONLY", "LITERAL")


def protectedRowGate(row):
    if row["rendering"] == "NONE":
        return row["outcome"] == "WROTE"
    if row["arm"] == "dir-modify" or (row["rendering"] == "LITERAL" and row["arm"] == "modify"):
        return row["outcome"] == "REFUSED" and row["host_before"] == row["host_after"]
    return None


def runProtectedMatrix(args):
    root = protectedScratchRoot(args.scratch_root)  # Before any mkdir, including metadata cwd.
    os.makedirs(root, exist_ok=True)
    node = os.path.realpath(args.node)
    package = os.path.join(os.path.realpath(args.srt_node_modules), "@anthropic-ai", "sandbox-runtime")
    cli = os.path.join(package, "dist", "cli.js")
    host_env = {"PATH": "/usr/bin:/bin"}

    def run(cmd, cwd=root):
        try:
            r = subprocess.run(cmd, cwd=cwd, env=host_env, stdin=subprocess.DEVNULL,
                               capture_output=True, text=True, timeout=30)
            return r.returncode, r.stdout, r.stderr
        except subprocess.TimeoutExpired as exc:
            return "timeout", (exc.stdout or b"").decode(), (exc.stderr or b"").decode()

    def snapshot(wt, targets):
        result = {}
        for rel in targets:
            path = os.path.join(wt, rel)
            if os.path.islink(path):
                result[rel] = "symlink:" + hashlib.sha256(os.fsencode(os.readlink(path))).hexdigest()
            elif os.path.exists(path):
                with open(path, "rb") as f:
                    result[rel] = hashlib.sha256(f.read()).hexdigest()
            else:
                result[rel] = None
        return result

    with open(os.path.join(package, "package.json")) as f:
        version = json.load(f)["version"]
    with open(cli, "rb") as f:
        cli_hash = hashlib.sha256(f.read()).hexdigest()
    meta = {"platform": platform.platform(), "srt_version": version, "cli_sha256": cli_hash,
            "node_version": run([node, "--version"])[1].strip(),
            "bwrap_version": run(["bwrap", "--version"])[1].strip() if sys.platform == "linux" else "n/a"}
    arch = {"x86_64": "x64", "aarch64": "arm64"}.get(platform.machine(), platform.machine())
    seccomp = os.path.join(package, "vendor", "seccomp", arch)
    reads = [seccomp] if sys.platform == "linux" and args.seccomp_read == "allowread" and os.path.isdir(seccomp) else []
    files = ("pkg/a.go", "pkg/a_test.go", "a_test.go", "tests/test_x.py", "cli/conftest.py",
             "src/foo.test.ts", "PKG/A_TEST.GO")
    census = {(arm, rendering) for arm, _, _ in PROTECTED_ARMS for rendering in PROTECTED_RENDERINGS}
    rows, gates = [], {}
    for rendering in PROTECTED_RENDERINGS:
        for arm, steps, targets in PROTECTED_ARMS:
            rowroot = tempfile.mkdtemp(prefix=rendering + "-" + arm + "-", dir=root)
            wt = os.path.join(rowroot, "wt")
            for rel in files:
                path = os.path.join(wt, rel)
                os.makedirs(os.path.dirname(path), exist_ok=True)
                with open(path, "w") as f:
                    f.write("// ORIG\n")
            for rel in ("cli/sub", ".home", ".tmp"):
                os.makedirs(os.path.join(wt, rel), exist_ok=True)
            literals = sorted(os.path.join(wt, rel) for rel in files if
                              re.fullmatch(r".*_test\.go|conftest\.py|.*\.test\.ts", os.path.basename(rel), re.I))
            deny = [] if rendering == "NONE" else [os.path.join(wt, "tests")]
            if rendering == "LITERAL":
                deny += literals
            settings = {"filesystem": {"allowWrite": [wt], "denyWrite": deny, "denyRead": [], "allowRead": reads},
                        "network": {"allowedDomains": [], "deniedDomains": [], "allowLocalBinding": False,
                                    "allowAllUnixSockets": False}, "enableWeakerNetworkIsolation": False, "allowAppleEvents": False}
            settings_path = os.path.join(rowroot, "settings.json")
            with open(settings_path, "w", opener=lambda p, flags: os.open(p, flags, 0o600)) as f:
                json.dump(settings, f)
            script = "process.env.HOME=process.cwd()+'/.home';process.env.TMPDIR=process.cwd()+'/.tmp';const fs=require('fs');try{"
            script += "".join("fs." + step + ";console.log('WORLD-PROBE-STEP " + str(i + 1) + "');" for i, step in enumerate(steps))
            script += "console.log('WORLD-PROBE-WROTE')}catch(e){console.log('WORLD-PROBE-REFUSED '+e.code);process.exitCode=1}"
            before = snapshot(wt, targets)
            rc, stdout, stderr = run([node, cli, "--settings", settings_path, "--", node, "-e", script], wt)
            after = snapshot(wt, targets)
            # Remove the row's fixture tree at once: it holds .go files, and repo-wide Go scans
            # later in the CI job must never see them (the JSON keeps the evidence).
            shutil.rmtree(rowroot)
            tokens = [line for line in stdout.splitlines() if line.startswith(("WORLD-PROBE-WROTE", "WORLD-PROBE-REFUSED"))]
            token = tokens[0] if len(tokens) == 1 else None
            outcome = "WROTE" if before != after and token == "WORLD-PROBE-WROTE" else "PARTIAL"
            if before == after and token and re.fullmatch(r"WORLD-PROBE-REFUSED E[A-Z0-9]+", token):
                outcome = "REFUSED"
            row = {"arm": arm, "rendering": rendering, "outcome": outcome, "client_token": token, "rc": rc, "worktree": wt,
                   "host_before": before, "host_after": after, "deny_write": deny, "literals": literals if rendering == "LITERAL" else [],
                   "steps_passed": re.findall(r"^WORLD-PROBE-STEP [0-9]+$", stdout, re.M), "stdout": stdout, "stderr": stderr}
            rows.append(row)
            gate = protectedRowGate(row)
            if gate is not None:
                gates[rendering + "/" + arm] = gate
            print(f"[{'INFO' if gate is None else 'GATE PASS' if gate else 'GATE FAIL'}] {rendering}/{arm} outcome={outcome} token={token} host={'changed' if before != after else 'unchanged'}", flush=True)
    leftover = sorted(os.listdir(root))
    gates["fixtures-removed"] = not leftover
    if not leftover:
        os.rmdir(root)
    rejection = validateProtectedRows(rows, census)
    gates["validator"] = rejection is None
    result = {"meta": meta, "rows": rows, "census_size": len(census), "gating": gates, "validation_error": rejection}
    with open(args.json, "w") as f:
        json.dump(result, f, indent=2)
    print("PROTECTED-JSON " + json.dumps(result, separators=(",", ":")), flush=True)
    return 0 if all(gates.values()) else 1


def ProtectedMatrixSelfTest():
    census = {("modify", "DIRS_ONLY"), ("control", "NONE")}
    digest = hashlib.sha256(b"original host bytes").hexdigest()
    refused = {"arm": "modify", "rendering": "DIRS_ONLY", "outcome": "REFUSED",
               "client_token": "WORLD-PROBE-REFUSED EPERM", "rc": 1,
               "host_before": {"tests/test.py": digest}, "host_after": {"tests/test.py": digest}}
    wrote = {"arm": "control", "rendering": "NONE", "outcome": "WROTE",
             "client_token": "WORLD-PROBE-WROTE", "rc": 0,
             "host_before": {"ordinary.txt": None}, "host_after": {"ordinary.txt": digest}}
    cases = []

    def check(name, ok):
        cases.append(bool(ok))
        print(f"[{'PASS' if ok else 'FAIL'}] {name}", flush=True)

    partial = dict(refused, outcome="PARTIAL", host_after={"tests/test.py": None})
    linked = dict(refused, host_before={"link": "symlink:" + digest}, host_after={"link": "symlink:" + digest})
    check("ok", all(validateProtectedRows([r, wrote], census) is None for r in (
        refused, partial, dict(partial, client_token=None, rc="timeout"), linked)))
    check("empty", validateProtectedRows([], census) == "census-mismatch" and
          validateProtectedRows([], set()) == "empty-census" and
          validateProtectedRows([refused], census) == "census-mismatch" and
          validateProtectedRows([refused, refused], census) == "census-mismatch")
    check("missing-client", validateProtectedRows([dict(refused, client_token=None), wrote], census) == "missing-client")
    changed = dict(refused, host_after={"tests/test.py": hashlib.sha256(b"forged refusal").hexdigest()})
    check("host-changed", validateProtectedRows([changed, wrote], census) == "host-changed")
    none_refused = dict(wrote, outcome="REFUSED", client_token="WORLD-PROBE-REFUSED EBUSY", host_before=wrote["host_after"])
    check("none-not-written", validateProtectedRows([refused, none_refused], census) is None and protectedRowGate(none_refused) is False)
    check("wrote-unchanged", validateProtectedRows([refused, dict(wrote, host_before=wrote["host_after"])], census) == "wrote-unchanged")
    check("invalid-outcome", validateProtectedRows([dict(refused, outcome="BOGUS"), wrote], census) == "invalid-outcome")
    check("bad-digest", all(validateProtectedRows([dict(refused, host_before={"tests/test.py": d}), wrote], census) == "invalid-host-snapshot"
                            for d in ("g" * 64, "a" * 63, "symlink:" + "a" * 63)))
    check("bad-token", all(validateProtectedRows([dict(refused, client_token=t), wrote], census) == "missing-client"
                           for t in ("WORLD-PROBE-REFUSED eperm", "WORLD-PROBE-REFUSED")))
    check("invalid-rc", all(validateProtectedRows([dict(refused, rc=r), wrote], census) == "invalid-rc"
                            for r in ("bogus", None)) and
          validateProtectedRows([dict(refused, rc="timeout"), wrote], census) is None)
    check("snapshot-keys", validateProtectedRows([dict(refused, host_after={"other.txt": digest}), wrote], census)
          == "invalid-host-snapshot" and
          validateProtectedRows([dict(refused, host_before={}, host_after={}), wrote], census) == "invalid-host-snapshot")
    check("partial-token", validateProtectedRows([dict(partial, client_token="garbage"), wrote], census) == "missing-client")
    check("wrote-refused-token", validateProtectedRows([dict(wrote, client_token="WORLD-PROBE-REFUSED EPERM"), refused],
                                                       census) == "missing-client" and
          validateProtectedRows([dict(refused, client_token="WORLD-PROBE-WROTE"), wrote], census) == "missing-client")
    # The candidate is outside a synthetic allowed root but still under cwd,
    # so even a broken containment guard cannot cause off-worktree writes.
    with tempfile.TemporaryDirectory(prefix=".protected-self-test-", dir=os.getcwd()) as tmp:
        allowed = os.path.join(tmp, "allowed")
        candidate = os.path.join(tmp, "allowed-sibling", "must-not-exist")
        os.mkdir(allowed)
        outside = os.path.join(tmp, "outside")
        os.mkdir(outside)
        link = os.path.join(allowed, "link")
        os.symlink(outside, link)
        for name, path in (("root-equals-allowed", allowed), ("symlink-escape", link)):
            rejected = False
            try:
                protectedScratchRoot(path, allowed)
            except ValueError:
                rejected = True
            check(name, rejected)
        before = os.listdir(tmp)
        rejected = False
        try:
            protectedScratchRoot(candidate, allowed)
        except ValueError:
            rejected = True
        check("outside-root", rejected and not os.path.exists(candidate) and os.listdir(tmp) == before)
    print(f"self-test: {len(cases)} cases, {sum(cases)} passed", flush=True)
    return 0 if all(cases) else 1


# Dispatch BEFORE legacy required arguments and all HOME/decoy/archive setup.
mode_parser = argparse.ArgumentParser(add_help=False, allow_abbrev=False)
mode_group = mode_parser.add_mutually_exclusive_group()
mode_group.add_argument("--protected-self-test", action="store_true")
mode_group.add_argument("--protected-only", action="store_true")
mode, _ = mode_parser.parse_known_args()
if mode.protected_self_test or mode.protected_only or "--help" in sys.argv or "-h" in sys.argv:
    protected_parser = argparse.ArgumentParser(parents=[mode_parser], allow_abbrev=False)
    if mode.protected_only or "--help" in sys.argv or "-h" in sys.argv:
        protected_parser.add_argument("--scratch-root", required=True)
        protected_parser.add_argument("--srt-node-modules", required=True)
        protected_parser.add_argument("--node", required=True)
        protected_parser.add_argument("--json", required=True)
        protected_parser.add_argument("--seccomp-read", choices=["none", "allowread"], default="none")
    protected_args = protected_parser.parse_args()
    if mode.protected_self_test:
        sys.exit(ProtectedMatrixSelfTest())
    try:
        sys.exit(runProtectedMatrix(protected_args))
    except ValueError as exc:
        protected_parser.error(str(exc))

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

# .git/hooks and .git/config come FIRST (design §4.4 [M2-linux], V50): when an episode's .git is a
# DIRECTORY that lacks them, srt (Linux) stubs each absent one with a bind mount, and a stub emitted
# after the .git read-only bind makes bwrap abort ("Can't create file at …/.git/hooks: Read-only file
# system"). Listed before .git, the stubs are mounted while .git is still writable.
DENY_WRITE = [os.path.join(EP1, n) for n in (".git/hooks", ".git/config", ".git", ".github", ".pi", ".claude", ".ailang",
                                             ".gitmodules", ".gitattributes")] + \
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

def srt(argv, st=SET, tramp=True, dashdash=True, cli=CLI, timeout=60, cwd=EP1):
    cmd = [NODE, cli, "--settings", st] + (["--"] if dashdash else []) + ((TRAMP + KV) if tramp else []) + argv
    try:
        r = subprocess.run(cmd, cwd=cwd, env=HOST_ENV, capture_output=True, text=True, timeout=timeout, stdin=subprocess.DEVNULL)
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
        # macOS refuses (EPERM). Linux/bwrap masks a read-denied dir with an empty tmpfs, so the write can
        # land in that throwaway mount: the gating property is "host bytes unchanged", the token is recorded.
        ok = (not os.path.exists(p)) and (tok(out) in REFUSED if sys.platform == "darwin" else tok(out) is not None)
        rec(f"arm2 write {name}: host unchanged", True, "host absent (macOS: REFUSED; Linux: REFUSED or WROTE-into-tmpfs)",
            f"{tok(out)}; host exists={os.path.exists(p)}", ok)
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
    rec("V8 ls $HOME leaks nothing", True, "ran; BASE not listed (macOS rc!=0; Linux masked: empty)", f"rc {rc}: {out[:160]}",
        os.path.basename(BASE) not in out and (rc != 0 if sys.platform == "darwin" else rc == 0))
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

    # ---- .git as a DIRECTORY (M2-linux, V50): an episode whose .git holds HEAD, run from that episode ----
    # srt adds <cwd>/.git/hooks and <cwd>/.git/config to its own deny only when .git is a directory, and
    # leaves .git/HEAD writable (V5): only World's rendered .git deny closes it. Each arm has its own
    # episode so one arm's write cannot feed the next.
    def git_dir_episode(name, with_hooks_config=False):
        ep = os.path.join(ROOT, name)
        os.makedirs(os.path.join(ep, ".git"), exist_ok=True)
        with open(os.path.join(ep, ".git", "HEAD"), "w") as f:
            f.write("ref: refs/heads/dev\n")
        if with_hooks_config:
            os.makedirs(os.path.join(ep, ".git", "hooks"), exist_ok=True)
            open(os.path.join(ep, ".git", "config"), "w").close()
        return ep
    def ep_settings(name, ep, deny_write=None):
        dw = [os.path.join(ep, os.path.relpath(p, EP1)) if p.startswith(EP1 + os.sep) else p for p in DENY_WRITE]
        return settings(name, filesystem__allowRead=[ep, CACHE] + EXTRA_READ, filesystem__allowWrite=[ep, CACHE],
                        filesystem__denyWrite=deny_write(dw) if deny_write else dw)
    def git_head(label, ep, st, gating, want_refused=True):
        head = os.path.join(ep, ".git", "HEAD")
        rc, out = srt(writep(head), st=st, cwd=ep)
        after = open(head).read()
        left = sorted(n for n in os.listdir(os.path.join(ep, ".git")) if n != "HEAD")
        observed = f"rc {rc}: {tok(out)}; HEAD={after!r}; host .git also holds {left}; out={out[-300:]!r}"
        if want_refused:
            rec(label, gating, "REFUSED token; HEAD bytes unchanged", observed,
                tok(out) in REFUSED and after == "ref: refs/heads/dev\n")
        else:
            rec(label, gating, "WROTE (srt alone leaves .git/HEAD open, V5)", observed, tok(out) == "WORLD-PROBE-WROTE")
    epg = git_dir_episode("epg")
    git_head("V50 write .git/HEAD refused (.git a dir WITHOUT hooks/ or config)", epg, ep_settings("epg", epg), True)
    eph = git_dir_episode("eph", with_hooks_config=True)
    git_head("V50 write .git/HEAD refused (.git a dir WITH hooks/ and config)", eph, ep_settings("eph", eph), True)
    epn = git_dir_episode("epn")
    git_head("V50 control: .git/HEAD WRITTEN without World's .git denies", epn,
             ep_settings("epn", epn, deny_write=lambda dw: [p for p in dw if os.path.relpath(p, epn) not in
                                                           (".git", os.path.join(".git", "hooks"), os.path.join(".git", "config"))]),
             True, want_refused=False)
    epo = git_dir_episode("epo")
    rc, out = srt(writep(os.path.join(epo, ".git", "HEAD")), cwd=epo,
                  st=ep_settings("epo", epo, deny_write=lambda dw: [p for p in dw if not p.endswith(("/.git/hooks", "/.git/config"))]))
    rec("V50 hazard: .git deny WITHOUT the leading .git/hooks + .git/config entries", False,
        "Linux: bwrap aborts (Can't create file at …/.git/hooks); macOS: REFUSED", f"rc {rc}: {tok(out)}; out={out[-300:]!r}", True)

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
    rec("V9 bind 127.0.0.1 (allowLocalBinding=false)", False, "macOS REFUSED; Linux binds inside its own netns", tok(out), tok(out) is not None)
    # the gating property: a listener the child opens is NOT reachable from the host
    lst = js("const n=require('net');const s=n.createServer(c=>c.end());s.on('error',e=>console.log('WORLD-PROBE-REFUSED '+e.code));"
             "s.listen(0,'127.0.0.1',()=>{console.log('WORLD-PROBE-BOUND '+s.address().port);setTimeout(()=>s.close(),4000)})")
    sp = subprocess.Popen([NODE, CLI, "--settings", SET, "--"] + TRAMP + KV + lst, cwd=EP1, env=HOST_ENV,
                          stdout=subprocess.PIPE, stderr=subprocess.STDOUT, text=True, stdin=subprocess.DEVNULL)
    line = sp.stdout.readline().strip()
    reached = None
    if line.startswith("WORLD-PROBE-BOUND "):
        bport = int(line.split()[1])
        c = socket.socket(); c.settimeout(2)
        try:
            c.connect(("127.0.0.1", bport)); reached = True
        except OSError as e:
            reached = False
        finally:
            c.close()
    sp.wait(timeout=20)
    rec("V9 child listener unreachable from host", True, "BOUND-and-unreachable or REFUSED",
        f"{line}; host reached={reached}", line.startswith("WORLD-PROBE-REFUSED") or (line.startswith("WORLD-PROBE-BOUND ") and reached is False))
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
