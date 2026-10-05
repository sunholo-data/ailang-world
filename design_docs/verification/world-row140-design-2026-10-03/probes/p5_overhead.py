#!/usr/bin/env python3
"""P5: per-call overhead of the srt wrapper (CLI) vs a bare spawn and a bare sandbox-exec."""
import os, statistics, subprocess, sys, time
S = os.path.dirname(os.path.abspath(__file__))
SRT = os.path.expanduser("~/dev/sunholo-data/ailang/tools/pi-extensions/sandbox/node_modules/@anthropic-ai/sandbox-runtime/dist/cli.js")
N = int(sys.argv[1]) if len(sys.argv) > 1 else 20
os.chdir(os.path.join(S, "ws"))

def timeit(label, argv):
    ts = []
    for _ in range(N):
        t = time.perf_counter()
        r = subprocess.run(argv, stdout=subprocess.DEVNULL, stderr=subprocess.DEVNULL)
        ts.append((time.perf_counter() - t) * 1000)
        assert r.returncode == 0, (label, r.returncode)
    ts.sort()
    print(f"{label:48s} n={N} median={statistics.median(ts):7.1f}ms p90={ts[int(0.9*N)-1]:7.1f}ms min={ts[0]:7.1f}ms")

timeit("bare /usr/bin/true", ["/usr/bin/true"])
timeit("node -e 0 (node startup alone)", ["node", "-e", "0"])
timeit("sandbox-exec minimal profile /usr/bin/true", ["/usr/bin/sandbox-exec", "-p", "(version 1)(allow default)", "/usr/bin/true"])
timeit("srt (no network) /usr/bin/true", ["node", SRT, "--settings", os.path.join(S, "s-basic.json"), "--", "/usr/bin/true"])
timeit("srt (1 allowed domain) /usr/bin/true", ["node", SRT, "--settings", os.path.join(S, "s-net.json"), "--", "/usr/bin/true"])
timeit("srt (read fence) /usr/bin/true", ["node", SRT, "--settings", os.path.join(S, "s-readfence.json"), "--", "/usr/bin/true"])
