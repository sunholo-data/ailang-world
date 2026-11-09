import subprocess,sys,os
R="/Users/voightkampff/.ailang/state/world-iter372-judge"
os.environ["AILANG_BIN"]=os.path.expanduser("~/.pinned-ailang/ailang")
B=("./host/broker/","Seed|LegacyPartialCache")
D=("./host/daemon/","ExecCache|BindSeed|ReportRaces|EpisodeReportShowsBindSeed")
C=("./cmd/ailang-worldd/","PrintsExecCache")
M=[
("START-NOOP","host/broker/exec_seed.go","func (h *ExecHandler) StartSeed() {\n","func (h *ExecHandler) StartSeed() {\n\treturn\n",B),
("SKIP-CHECK","host/broker/handlers_exec.go","return h.awaitSeed(ctx)","return nil",B),
("NO-CAP","host/broker/exec_seed.go","case s.attempts < execSeedMaxAttempts:","case true:",B),
("SWAP","host/broker/exec_seed.go","case s.running:\n\t\treturn seedJoin\n\tcase s.attempts < execSeedMaxAttempts:\n\t\treturn seedStart","case s.running:\n\t\treturn seedStart\n\tcase s.attempts < execSeedMaxAttempts:\n\t\treturn seedJoin",B),
("CALLER-CTX","host/broker/exec_seed.go","context.WithTimeout(context.Background(), h.seedBudget)","context.WithTimeout(context.Background(), h.seedBudget)\n\t_ = 0",B),  # placeholder replaced below
("BUDGET","host/broker/exec_seed.go","handlerBounds{execTimeout: h.seedBudget}","handlerBounds{execTimeout: handlerExecTimeout}",B),
("MARKER-PREFIX","host/broker/exec_seed.go","bytes.Equal(data, seedMarkerBytes(name))","bytes.HasPrefix(data, []byte(\"world-exec-seed v\"))",B),
("LSTAT-READY","host/broker/exec_seed.go","func seedCacheVerified(dst, name string) bool {\n","func seedCacheVerified(dst, name string) bool {\n\t_, err0 := os.Lstat(dst)\n\treturn err0 == nil\n",B),
("NO-SIZECAP","host/broker/exec_seed.go","len(data) <= seedMarkerMax && ","",B),
("FIXED-STAGING","host/broker/exec_seed.go","return seedStagingPrefix(project, name) + incarnation + \"-\" + strconv.Itoa(k)","return strings.TrimSuffix(seedStagingPrefix(project, name), \"-\")",B),
("NO-CLEANUP","host/broker/exec_seed.go","func (h *ExecHandler) cleanupStaging(epDir, name, own string) {\n","func (h *ExecHandler) cleanupStaging(epDir, name, own string) {\n\treturn\n",B),
("UNBOUNDED","host/broker/exec_seed.go","\t\tif removed == seedCleanupMax {\n\t\t\tbreak\n\t\t}\n","",B),
("NO-STEP1","host/broker/exec_seed.go","if _, err := os.Lstat(filepath.Join(seed, seedMarkerName)); err == nil {","if false {",B),
("REMOVE-OK","host/broker/exec_seed.go","return fail(\"an unverified copy cannot be removed\", wrap(err))","return nil",B),
("NO-REPORT","host/daemon/episode_report.go","rep.ExecCache = d.workspace.execCacheStatus(episode)\n","",D),
("NO-BIND","host/daemon/workspace.go","\teh.StartSeed()\n","",D),
("NO-STEP0","host/broker/exec_seed.go","if err := os.MkdirAll(h.cache, 0o700); err != nil {\n\t\treturn fail(\"cache dir cannot be created\"","if err := error(nil); err != nil {\n\t\treturn fail(\"cache dir cannot be created\"",D),
("NO-LINE","cmd/ailang-worldd/episode_report.go","for _, e := range r.ExecCache {\n\t\tfmt.Fprintln(w, execCacheLine(e))\n\t}\n","",C),
("UNLOCKED","host/daemon/exec_seed.go","\tw.execMu.Lock()\n\tt, ok := w.execEntryLocked(episodeID)\n\tw.execMu.Unlock()","\tt, ok := w.execEntryLocked(episodeID)",("./host/daemon/","ReportRacesFirstDispatch -race")),
# own
("X-MARKER-BEFORE-CP","host/broker/exec_seed.go","if err := os.WriteFile(filepath.Join(staging, \"copy\", seedMarkerName)","if err := os.WriteFile(filepath.Join(staging, seedMarkerName)",B),
("X-GIVEUP-NOTGUARD","host/broker/exec_seed.go","func (h *ExecHandler) giveUpLocked(name string) *HandlerFailure {\n","func (h *ExecHandler) giveUpLocked(name string) *HandlerFailure {\n\tif true {\n\t\treturn nil\n\t}\n",B),
("X-READY-ON-FAIL","host/broker/exec_seed.go","h.seed.took, h.seed.last, h.seed.ok, h.seed.running = took, f, f == nil, nil","h.seed.took, h.seed.last, h.seed.ok, h.seed.running = took, f, true, nil",B),
("X-ATTEMPT-NOINC","host/broker/exec_seed.go","\th.seed.attempts++\n\tk := h.seed.attempts","\tk := h.seed.attempts + 1",B),
("X-NOT-IDEMPOTENT","host/broker/exec_seed.go","if h.seed.running != nil || h.seed.attempts > 0 || h.seed.ok {","if false {",B),
("X-RM-STAGING-OTHERS-OWN","host/broker/exec_seed.go","|| e.Name() == own {","{",B),
("X-LAST-NOTRETURNED","host/broker/exec_seed.go","\t\tif last != nil {\n\t\t\tc := *last\n\t\t\treturn &c","\t\tif false {\n\t\t\tc := *last\n\t\t\treturn &c",B),
]
only=sys.argv[1:]
res=[]
for id,f,old,new,(pkg,rx) in M:
    if only and id not in only: continue
    if id=="CALLER-CTX": continue
    p=os.path.join(R,f); s=open(p).read()
    if old not in s: print(id,"PATTERN-MISSING"); continue
    open(p,"w").write(s.replace(old,new,1))
    race=[]
    if rx.endswith(" -race"): rx=rx[:-6]; race=["-race"]
    try:
        r=subprocess.run(["go","test"]+race+[pkg,"-run",rx,"-count=1","-timeout","120s"],cwd=R,capture_output=True,text=True,timeout=200)
        out=r.stdout+r.stderr
        killed=r.returncode!=0
        why="BUILDFAIL" if "build failed" in out or "[build failed]" in out or "cannot use" in out or "declared and not used" in out else ""
        fails=[l.split()[2] for l in out.splitlines() if l.startswith("--- FAIL")][:6]
    except Exception as e:
        killed=True;why="TIMEOUT";fails=[]
    subprocess.run(["git","checkout","--","f"[0:0] or f],cwd=R)
    print(id,"KILLED" if killed else "SURVIVED",why,fails,flush=True)
st=subprocess.run(["git","status","--short"],cwd=R,capture_output=True,text=True).stdout
print("STATUS:",repr(st))
