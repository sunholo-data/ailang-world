#!/usr/bin/env python3
# iter-361 designer: row 256 mutation table. Usage: python3 designer-mutate.py.txt <tree> <out>
# For each mutant: every edit is an exact-string replace whose anchor must occur
# exactly once in its file; run the named tests once (-count=1) across the four
# packages; parse --- PASS/--- FAIL of the TOP-LEVEL tests; restore the files
# byte-for-byte (sha256 checked). CONTROL (no edit) runs first and must be green
# before any mutant verdict is read. Environment: TMPDIR/GOTMPDIR untouched,
# GOTOOLCHAIN untouched, AILANG_BIN pinned.
import hashlib, os, re, subprocess, sys, time

tree, out = sys.argv[1], sys.argv[2]
TESTS = {
    "Serve": "TestServeStartsWhenTheProbeSpendsTheStartupClock",
    "Phase": "TestRegistryPhaseHasItsOwnClock",
    "Bound": "TestRegistryPhaseClockIsBounded",
    "Pin": "TestServeStartupClocksArePinned",
    "Spent": "TestExecProbeOutOfTimeNamesTheClockNotAnArm",
    "Doctor": "TestDoctorNamesLeftOutExecProfiles",
    # existing neighbours, expected green under every mutant
    "Quar": "TestStartupBudgetQuarantinesPostCutoffBootstrap",
    "Held": "TestStartupBudgetHeldConnectionSettles",
    "Health": "TestHealthReportsBuildToolCredentialsAndProfiles",
    "ProfileFor": "TestProfileForResolutionOrder",
    "Isolate": "TestExecProfilesDirProbeFailureIsolatesOneProfile",
    "DocDaemon": "TestDoctorDaemonVsForeignListener",
    "DocHelp": "TestDoctorHelpAndAPIKeyGuard",
    "DocM3": "TestDoctorExecToolsWithAnM3Daemon",
    "Roots": "TestProductionContextRoots",
}
PKGS = ["./host/daemon/", "./host/broker/", "./cmd/ailang-worldd/", "./host/store/"]
D, W, B, S, DOC, DH = ("host/daemon/daemon.go", "host/daemon/workspace.go", "host/broker/exec_probe.go",
                       "host/broker/exec_startup.go", "cmd/ailang-worldd/doctor.go", "cmd/ailang-worldd/doctor_hosts.go")

M = [
    ("CONTROL", [], set()),
    ("MUT-RUN-HEAD-SHAPE", [(D, "d, err := startDaemon(ctx, cfg, startupBudget, registryPhaseBudget)",
                             "d, err := startBounded(ctx, cfg, startupBudget, New)")], {"Serve"}),
    ("MUT-SHARED-PARENT", [(D, "return newDaemon(startupCtx, ctx, registry, cfg)",
                            "return newDaemon(startupCtx, startupCtx, registry, cfg)")], {"Serve", "Phase"}),
    ("MUT-REGISTRY-UNBOUNDED", [(D, "regCtx, cancelReg := context.WithTimeout(phaseParent, registryBudget)",
                                 "regCtx, cancelReg := context.WithCancel(phaseParent)")], {"Bound"}),
    ("MUT-SWEEP-ON-STARTUP-CLOCK", [(D, "d.integrity = d.scanIntegrity(regCtx)", "d.integrity = d.scanIntegrity(ctx)")], {"Phase"}),
    ("MUT-TEMPLATES-ON-STARTUP-CLOCK", [(D, "d.reportTemplates(regCtx)", "d.reportTemplates(ctx)")], set()),  # predicted survivor
    ("MUT-SERVE-SPENT", [(S, "\t\tif errs[i] == nil {\n\t\t\tcontinue\n\t\t}\n\t\tp := order[i]",
                          "\t\tif pe := (*ExecProbeError)(nil); errs[i] == nil || errors.As(errs[i], &pe) && pe.OutOfTime() {\n\t\t\tcontinue\n\t\t}\n\t\tp := order[i]")],
     {"Spent", "Serve", "Phase"}),
    ("MUT-SPENT-STILL-SERVED", [(S, "\t\tdelete(st.Profiles, p.Project)\n\t\tdelete(st.Sources, p.Project)",
                                 "\t\tif !spent {\n\t\t\tdelete(st.Profiles, p.Project)\n\t\t}\n\t\tdelete(st.Sources, p.Project)")],
     {"Spent", "Serve", "Phase"}),
    ("MUT-FAIL-NEVER-SPENT", [(B, "spent := !pp.deadline.IsZero() && !time.Now().Before(pp.deadline)",
                               "spent := false && !pp.deadline.IsZero() && !time.Now().Before(pp.deadline)")],
     {"Spent", "Serve", "Phase"}),
    ("MUT-DEADLINE-UNSET", [(B, "\t\tpp.deadline = dl\n", "\t\t_ = dl\n")], {"Serve", "Phase"}),
    ("MUT-SPENT-ANY", [(B, "\t\tif !f.spent {\n\t\t\treturn false\n\t\t}\n\t}\n\treturn len(e.Failed) > 0",
                        "\t\tif f.spent {\n\t\t\treturn true\n\t\t}\n\t}\n\treturn false")], {"Spent"}),
    ("MUT-FLAG-NOT-CARRIED", [(S, "OutOfTime: spent})", "OutOfTime: false && spent})")], {"Spent", "Serve", "Phase"}),
    ("MUT-HEALTH-DROPS-FLAG", [(W, "Reason: x.Reason, OutOfTime: x.OutOfTime})", "Reason: x.Reason})")], {"Serve", "Phase"}),
    ("MUT-NO-LOG-LINE", [(W, "\t\tif x.OutOfTime {", "\t\tif false && x.OutOfTime {")], {"Serve"}),
    ("MUT-DOCTOR-NOT-CALLED", [(DOC, "\tdoctorExecProfiles(r, h.ExecExcluded)\n", "\t_ = h.ExecExcluded\n")], {"Doctor"}),
    ("MUT-DOCTOR-WARN", [(DOC, 'r.add(markFail, "exec profiles", detail, fix)', 'r.add(markWarn, "exec profiles", detail, fix)')], {"Doctor"}),
    ("MUT-DOCTOR-CLASS", [(DH, 'case "exec tools", "exec profiles":', 'case "exec tools":')], {"Doctor"}),
    ("MUT-DOCTOR-ONE-FIX", [(DOC, "\t\tif x.OutOfTime {", "\t\tif false && x.OutOfTime {")], {"Doctor"}),
    ("MUT-REGISTRY-BUDGET-30S", [(D, "registryPhaseBudget = 9 * time.Second", "registryPhaseBudget = 30 * time.Second")], {"Pin"}),
    ("MUT-STARTUP-BUDGET-5S", [(D, "startupBudget       = 9 * time.Second", "startupBudget       = 5 * time.Second")], {"Pin"}),
]

only = set(sys.argv[3].split(",")) if len(sys.argv) > 3 else None
regex = "^(" + "|".join(TESTS.values()) + ")$"
env = dict(os.environ, AILANG_BIN=os.path.expanduser("~/.pinned-ailang/ailang"))
name_of = {v: k for k, v in TESTS.items()}
log = open(out, "a")


def sha(p):
    return hashlib.sha256(open(os.path.join(tree, p), "rb").read()).hexdigest()


def say(s):
    print(s, flush=True)
    log.write(s + "\n")
    log.flush()


say("== run %s  tree=%s  load=%s" % (time.strftime("%Y-%m-%dT%H:%M:%SZ", time.gmtime()), tree, os.getloadavg()))
matches = mismatches = 0
for label, edits, want in M:
    if only and label not in only and label != "CONTROL":
        continue
    files = sorted({f for f, _, _ in edits})
    before = {f: open(os.path.join(tree, f), "rb").read() for f in files}
    shas = {f: sha(f) for f in files}
    ok = True
    for f, old, new in edits:
        p = os.path.join(tree, f)
        s = open(p).read()
        n = s.count(old)
        if n != 1:
            say("%s: ANCHOR-COUNT %d in %s: %r" % (label, n, f, old[:60]))
            ok = False
            break
        open(p, "w").write(s.replace(old, new))
    if ok:
        t0 = time.time()
        build_rc = subprocess.run(["go", "build", "./..."], cwd=tree, env=env, capture_output=True, text=True, timeout=120).returncode
        say("%s: explicit go build rc=%s" % (label, build_rc))
        if build_rc != 0: raise RuntimeError("mutant build failed, no behavioral verdict")
        r = subprocess.run(["go", "test", "-count=1", "-timeout", "120s", "-v", "-run", regex] + PKGS,
                           cwd=tree, env=env, capture_output=True, text=True, timeout=180)
        o = r.stdout + r.stderr
        passed = set(re.findall(r"^--- PASS: (\S+)", o, re.M))
        failed = set(re.findall(r"^--- FAIL: (\S+)", o, re.M))
        build = "BUILD-FAIL" if ("[build failed]" in o or "[setup failed]" in o) else ""
        red = sorted(name_of[t] for t in failed if t in name_of)
        green = len([t for t in passed if t in name_of])
        missing = sorted(k for k, v in TESTS.items() if v not in passed and v not in failed)
        verdict = "MATCH" if set(red) == want and not missing and not build else "MISMATCH"
        if verdict == "MATCH":
            matches += 1
        else:
            mismatches += 1
        say("%s: rc=%d pass=%d red=%s expected=%s missing=%s %s %s (%.0fs)" % (
            label, r.returncode, green, ",".join(red) or "-", ",".join(sorted(want)) or "-",
            ",".join(missing) or "-", build, verdict, time.time() - t0))
        if verdict != "MATCH" or label == "CONTROL":
            for line in o.splitlines():
                if re.match(r"^\s+\S+_test\.go:\d+:", line) or "FAIL" in line or "build failed" in line:
                    log.write("    | " + line[:400] + "\n")
    for f in files:
        open(os.path.join(tree, f), "wb").write(before[f])
        if sha(f) != shas[f]:
            say("%s: RESTORE FAILED for %s" % (label, f))
            sys.exit(2)
say("== done: %d MATCH, %d MISMATCH" % (matches, mismatches))
