import subprocess,json,time,sys,pathlib
sha=sys.argv[1]; target=pathlib.Path(sys.argv[2]); deadline=time.time()+1800
assert len(sha)==40 and all(c in "0123456789abcdef" for c in sha)
expected={"ailang-code verify gate","go host build + test gate"}
while time.time()<deadline:
    r=subprocess.run(["gh","api",f"repos/sunholo-data/ailang-world/commits/{sha}/check-runs?per_page=100"],capture_output=True,text=True,timeout=30)
    if r.returncode:
        print("INSTRUMENT FAILURE",r.returncode,flush=True);time.sleep(20);continue
    data=json.loads(r.stdout); checks=data.get("check_runs",[]); by={c["name"]:c for c in reversed(checks)}
    target.write_text(json.dumps({"sha":sha,"read_epoch":time.time(),"data":data},indent=2)+"\n")
    missing=expected-set(by); bad=[(n,by[n]["conclusion"]) for n in expected&set(by) if by[n]["status"]=="completed" and by[n]["conclusion"]!="success"]
    pending=[n for n in by if by[n]["status"]!="completed"]
    bad += [(n,by[n]["conclusion"]) for n in set(by)-expected if by[n]["status"]=="completed" and by[n]["conclusion"] not in ("success","skipped","neutral")]
    print(sha,"expected=2", "present="+str(len(expected&set(by))),"missing="+str(sorted(missing)),"pending="+str(pending),"bad="+str(bad),flush=True)
    if bad: sys.exit(1)
    if not missing and not pending:
        rr=subprocess.run(["gh","api",f"repos/sunholo-data/ailang-world/actions/runs?head_sha={sha}&per_page=100"],capture_output=True,text=True,timeout=30)
        assert rr.returncode==0
        runs=json.loads(rr.stdout); target.with_suffix(".runs.json").write_text(json.dumps(runs,indent=2)+"\n")
        ci=[x for x in runs["workflow_runs"] if x["name"]=="CI" and x["status"]=="completed" and x["conclusion"]=="success"]
        if ci:
            print("ALL COMPLETE SUCCESS for",sha,"CI runs",[x["id"] for x in ci],flush=True);sys.exit(0)
    time.sleep(30)
print("TIMEOUT no verdict for",sha,flush=True);sys.exit(124)
