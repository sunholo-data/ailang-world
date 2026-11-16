import json, subprocess, sys
BIN = sys.argv[1]
F = "world/transitions.ail"
MUTS = [
 ("MUT-5  plan inputWorld: w.logHead", "inputWorld: w.stateRoot,", "inputWorld: w.logHead,"),
 ("MUT-5a plan inputWorld algo-only differs", "inputWorld: w.stateRoot,", 'inputWorld: { algo: "x", digest: w.stateRoot.digest },'),
 ("MUT-5b plan inputWorld digest-only differs", "inputWorld: w.stateRoot,", 'inputWorld: { algo: w.stateRoot.algo, digest: "x" },'),
 ("MUT-6  plan proposalHash: transitionFn", "    proposalHash: proposalHash,\n    inputWorld", "    proposalHash: transitionFn,\n    inputWorld"),
 ("MUT-6a plan transitionFn: proposalHash", "    transitionFn: transitionFn,\n    goal", "    transitionFn: proposalHash,\n    goal"),
 ("MUT-7  verify rejected accepted: true", 'accepted: false, reason: "input world mismatch"', 'accepted: true, reason: "input world mismatch"'),
 ("MUT-7a verify accepted-branch proposalHash: p.inputWorld", 'then { proposalHash: p.proposalHash, accepted: true', 'then { proposalHash: p.inputWorld, accepted: true'),
 ("CTRL   unmutated", None, None),
]
orig = open(F).read()
ver = subprocess.run([BIN, "--version"], capture_output=True, text=True).stdout.splitlines()[0]
for name, old, new in MUTS:
    if old is not None:
        assert orig.count(old) == 1, name
        open(F, "w").write(orig.replace(old, new))
    try:
        d = json.loads(subprocess.run([BIN, "ai-check", "-timeout", "5s", F], capture_output=True, text=True).stdout)
        st = {r["function"]: r["status"] for r in (d.get("verify") or {}).get("results", []) or []}
        print("%s | %-55s | passed=%s | plan=%s verify=%s applyRevision=%s" % (ver, name, d["check"]["passed"], st.get("plan"), st.get("verify"), st.get("applyRevision")))
    finally:
        open(F, "w").write(orig)
