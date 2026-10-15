exec(open('/tmp/world268-evaluator-r1/run.py').read().split("for name in ['AuthenticatedCommits'")[0])
results=[]
run('independent-quickstart-red','^TestQuickstartSection11FlagsMatchTheCLI$')
run('independent-boundary-red','^TestWorldBoundary(DependencyAllowlist|RecordingWriter)$','./host/boundary')
try:
 doc='docs/QUICKSTART.md';(root/doc).write_bytes(subprocess.check_output(['git','show','39066a46c1ac13c994bacb162b8d4b754d5ec083:'+doc],cwd=root));run('attribution-quickstart-base-doc-green','^TestQuickstartSection11FlagsMatchTheCLI$');(root/doc).write_bytes(original[doc])
 f='cmd/ailang-worldd/demo_seed.go';s=original[f].decode();s=s.replace('\n\t"github.com/sunholo-data/ailang-world/host/registry"','');s=s.replace('reg, err := registry.Decode(obj.Payload)','var reg struct { Epochs []struct { Epoch int64 `json:"epoch"`; Candidates []string `json:"candidates"` } `json:"epochs"` }\n\terr = json.Unmarshal(obj.Payload, &reg)');(root/f).write_text(s);run('attribution-local-wire-decode-boundary-green','^TestWorldBoundary(DependencyAllowlist|RecordingWriter)$','./host/boundary');(root/f).write_bytes(original[f])
 # Own M2a production diff only reverted; earlier milestones remain unchanged.
 for f in ['cmd/ailang-worldd/demo.go','cmd/ailang-worldd/demo_seed.go']:(root/f).unlink()
 f='cmd/ailang-worldd/main.go';(root/f).write_bytes(subprocess.check_output(['git','show','39066a46c1ac13c994bacb162b8d4b754d5ec083:'+f],cwd=root))
 run('M2a-own-production-revert','^TestDemoSeed(AuthenticatedCommits|PreflightRefusals|HelpPreparation)$')
finally:
 for f,b in original.items():(root/f).write_bytes(b)
# Correctly named response assertion, rather than the executor's earlier no-retry expectation.
f='cmd/ailang-worldd/demo_seed.go';s=original[f].decode();old='if json.Unmarshal(body, &result) != nil || result.SelectedHead != commit.NextWorld.Ref {';assert s.count(old)==1
try:
 (root/f).write_text(s.replace(old,'if false && (json.Unmarshal(body, &result) != nil || result.SelectedHead != commit.NextWorld.Ref) {'));run('mutant-response-correct-assertion','^TestDemoSeedCommitFailures$/^commit-response$',expect='commit refusal boundary')
finally:(root/f).write_bytes(original[f])
run('final-restored-happy-help','^TestDemoSeed(AuthenticatedCommits|HelpPreparation|CommitFailures)$')
commands=[['gofmt','-l']+files[:4],['go','vet','./cmd/ailang-worldd'],['git','diff','--check','39066a46c1ac13c994bacb162b8d4b754d5ec083','67ee2efe5b2fdca3dcab4eaeda4849a8caea1c8c','--']+files]
for i,cmd in enumerate(commands):
 t=time.monotonic();p=subprocess.run(cmd,cwd=root,env=env,stdout=subprocess.PIPE,stderr=subprocess.STDOUT,timeout=120);(out/('quality-'+str(i)+'.log')).write_bytes(p.stdout);results.append(dict(name='quality-'+str(i),command=cmd,rc=p.returncode,seconds=round(time.monotonic()-t,3)))
(out/'attack-ledger.json').write_text(json.dumps(results,indent=2));assert all((root/f).read_bytes()==b for f,b in original.items());print('ATTACKS TERMINAL RESTORED',flush=True)
