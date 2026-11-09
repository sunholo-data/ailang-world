"""Row 256: unmodified base and candidate net6 runs aligned to the start page's
periodic GET, sweeping the start offset. Witness listener records each scan."""
import os, socket, subprocess, sys, threading, time

S = "$SCRATCH"
OUT = S + "/m256a"; os.makedirs(OUT, exist_ok=True)
ARMS = {"base": S + "/wt-base256/host/broker", "cand": S + "/wt-row256/host/broker"}
BIN = {"base": S + "/m256/broker-base.test", "cand": S + "/m256/broker-cand.test"}
T = "^TestExecProbeNet6ForeignDialIsRetriedNotRefused$/^raw-twice-proxy$"
env = dict(os.environ, AILANG_BIN=os.path.expanduser("~/.pinned-ailang/ailang"))

scans = []
def witness():
    s = socket.socket(); s.bind(("127.0.0.1", 0)); s.listen(16)
    while True:
        c, _ = s.accept(); scans.append(time.time())
        try: c.settimeout(0.3); c.recv(200)
        except Exception: pass
        c.close()
threading.Thread(target=witness, daemon=True).start()

log = open(OUT + "/aligned.tsv", "w")
while len(scans) < 2: time.sleep(0.2)
offsets = [x / 10 for x in range(2, 20, 2)]  # seconds before the predicted scan
for off in offsets:
    for arm in ("base", "cand"):
        period = scans[-1] - scans[-2]
        nxt = scans[-1] + period
        while nxt - off < time.time() + 0.5: nxt += period
        time.sleep(max(0, nxt - off - time.time()))
        t0 = time.time()
        p = subprocess.run([BIN[arm], "-test.run", T, "-test.count=1", "-test.v"], cwd=ARMS[arm], env=env, capture_output=True, text=True)
        t1 = time.time()
        time.sleep(1.5)
        hit = [round(s - t0, 2) for s in scans if t0 - 1 <= s <= t1 + 1]
        sig = "foreign 1, 2, 1" in p.stdout
        log.write(f"{off}\t{arm}\t{p.returncode}\t{t1 - t0:.2f}\t{hit}\t{sig}\n"); log.flush()
        if p.returncode: open(f"{OUT}/fail-{arm}-{off}.log", "w").write(p.stdout)
log.write("DONE\n")
