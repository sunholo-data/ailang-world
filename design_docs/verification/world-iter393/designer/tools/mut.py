import json, subprocess, sys, shutil
V604 = "/tmp/ailang-v0604-i393/ailang"
MUTS = [
 ("MUT-1 drop deriving(Eq) on EvidenceGrade", "world/types.ail", "  | CLAIMED\n  deriving (Eq)\n", "  | CLAIMED\n", 1),
 ("MUT-2 drop deriving(Eq) on TimeoutOutcome", "world/types.ail", "  | ExecuteUnderPriorAuthority\n  deriving (Eq)\n", "  | ExecuteUnderPriorAuthority\n", 1),
 ("MUT-3 plan ensures back to record ==", "world/transitions.ail", "result.inputWorld.algo == w.stateRoot.algo && result.inputWorld.digest == w.stateRoot.digest", "result.inputWorld == w.stateRoot", 1),
 ("MUT-4 verify ensures back to record ==", "world/transitions.ail", "result.proposalHash.algo == p.proposalHash.algo && result.proposalHash.digest == p.proposalHash.digest &&", "result.proposalHash == p.proposalHash &&", 1),
 ("MUT-5 plan body inputWorld: w.logHead", "world/transitions.ail", "inputWorld: w.stateRoot,", "inputWorld: w.logHead,", 1),
 ("MUT-6 plan body proposalHash: transitionFn", "world/transitions.ail", "    proposalHash: proposalHash,\n    inputWorld", "    proposalHash: transitionFn,\n    inputWorld", 1),
 ("MUT-7 verify rejected branch accepted: true", "world/transitions.ail", 'accepted: false, reason: "input world mismatch"', 'accepted: true, reason: "input world mismatch"', 1),
 ("MUT-8 plan ensures drops the algo conjuncts (weakening)", "world/transitions.ail", "result.inputWorld.algo == w.stateRoot.algo && ", "", 1),
]
for name, f, old, new, n in MUTS:
    orig = open(f).read(); assert orig.count(old) >= 1, name
    open(f, "w").write(orig.replace(old, new, n))
    try:
        out = subprocess.run([V604, "ai-check", "-timeout", "5s", f], capture_output=True, text=True).stdout
        d = json.loads(out); c = d["check"]; v = d.get("verify") or {}
        errs = [(e["code"], e.get("line")) for e in c.get("errors", [])]
        st = sorted((r["function"], r["status"]) for r in v.get("results", []) or [] if r["status"] != "verified")
        nver = sum(1 for r in v.get("results", []) or [] if r["status"] == "verified")
        print("%-58s check.passed=%s errs=%s non-verified=%s verified=%d" % (name, c["passed"], errs, st, nver))
    finally:
        open(f, "w").write(orig)
