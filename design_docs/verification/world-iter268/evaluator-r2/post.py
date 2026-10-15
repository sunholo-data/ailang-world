from run import *
B=pathlib.Path.home()/'.ailang/state/world-iter268/evaluator-r2-backups';orig={f:(B/f).read_bytes() for f in json.loads((E/'start.json').read_text())['files']}
for f,b in orig.items():assert (W/f).read_bytes()==b
ENV.pop('WORLD_TOOL_AILANG_BIN',None)
run('ci-env-unset-toolpin',['go','test','./cmd/ailang-worldd','-run','^TestDemoSeedPreflightRefusals$/^pin-mismatch$','-count=1','-timeout=90s','-v'])
ENV['WORLD_TOOL_AILANG_BIN']=str(pathlib.Path.home()/'.pinned-ailang-tools/v0.52.1/ailang')
run('decision-packet-doc',['go','test','./host/daemon','-run','^TestQuickstartDecisionPacketSection$','-count=1','-timeout=90s','-v'])
# Diff-derived supporting helper predicates: measure coverage, not product validity.
for name,old,new in [('state-root-coverage','[]byte(observed + "\\n" + o.Hash)','[]byte(o.Hash)'),('entry-hash-coverage','entryHash := hashref.SumSHA256(encoded).String()','entryHash := hashref.SumSHA256(append(encoded, 0)).String()')]:
 f='cmd/ailang-worldd/demo_seed.go';s=orig[f].decode();assert s.count(old)==1,(name,s.count(old))
 try:
  (W/f).write_text(s.replace(old,new))
  run(name+'-build',['go','test','-c','-o',str(B/(name+'.test')),'./cmd/ailang-worldd'])
  run(name,['go','test','./cmd/ailang-worldd','-run','^TestDemoSeedAuthenticatedCommits$','-count=1','-timeout=90s','-v'])
 finally:(W/f).write_bytes(orig[f])
# Correct the anticipated response assertion, retaining initial raw attack log.
f='cmd/ailang-worldd/demo_seed.go';s=orig[f].decode();old='if json.Unmarshal(body, &result) != nil || result.SelectedHead != commit.NextWorld.Ref {'
try:
 (W/f).write_text(s.replace(old,'if false && (json.Unmarshal(body, &result) != nil || result.SelectedHead != commit.NextWorld.Ref) {'))
 run('response-mutant-build',['go','test','-c','-o',str(B/'response.test'),'./cmd/ailang-worldd'])
 run('response-mutant-assertion',['go','test','./cmd/ailang-worldd','-run','^TestDemoSeedCommitFailures$/^commit-response$','-count=1','-timeout=90s','-v'])
finally:(W/f).write_bytes(orig[f])
try:
 for f in ['cmd/ailang-worldd/demo.go','cmd/ailang-worldd/demo_seed.go']:(W/f).unlink()
 f='cmd/ailang-worldd/main.go';(W/f).write_bytes(subprocess.check_output(['git','show','39066a46c1ac13c994bacb162b8d4b754d5ec083:'+f],cwd=W))
 run('M2a-own-revert-build',['go','test','-c','-o',str(B/'revert.test'),'./cmd/ailang-worldd'])
 run('M2a-own-revert',['go','test','./cmd/ailang-worldd','-run','^(TestDemoSeedAuthenticatedCommits|TestDemoSeedPreflightRefusals|TestDemoSeedHelpPreparation|TestDemoSeedCommitFailures|TestSessionMint_RequiresControllingTerminal)$','-count=1','-timeout=120s','-v'])
finally:
 for f,b in orig.items():(W/f).write_bytes(b)
for n in ['TestDemoSeedAuthenticatedCommits','TestDemoSeedPreflightRefusals','TestDemoSeedHelpPreparation','TestDemoSeedCommitFailures']:
 run('restored-'+n,['go','test','./cmd/ailang-worldd','-run','^'+n+'$','-count=1','-timeout=120s','-v'])
run('gofmt',['gofmt','-l']+list(orig)[:4])
run('vet',['go','vet','./cmd/ailang-worldd'])
run('code-doc-diffcheck',['git','diff','--check','39066a46c1ac13c994bacb162b8d4b754d5ec083','HEAD','--']+list(orig))
(E/'post.json').write_text(json.dumps(ledger,indent=2))
for f,b in orig.items():assert (W/f).read_bytes()==b
(E/'restoration.json').write_text(json.dumps({f:hashlib.sha256((W/f).read_bytes()).hexdigest() for f in orig},indent=2))
print('POST TERMINAL RESTORED',flush=True)
