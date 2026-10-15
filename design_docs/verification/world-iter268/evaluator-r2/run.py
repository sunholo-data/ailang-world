import pathlib,subprocess,json,time,os,hashlib
W=pathlib.Path('/Users/voightkampff/dev/sunholo-data/.eval-world-iter268-r2'); E=W/'design_docs/verification/world-iter268/evaluator-r2'
ENV=os.environ.copy();ENV.update(PATH='/opt/homebrew/bin:'+ENV['PATH'],AILANG_BIN=str(pathlib.Path.home()/'.pinned-ailang/ailang'),WORLD_TOOL_AILANG_BIN=str(pathlib.Path.home()/'.pinned-ailang-tools/v0.52.1/ailang'),GOTOOLCHAIN='go1.26.6')
for k in ['GOOGLE_CLOUD_PROJECT','OTEL_EXPORTER_OTLP_ENDPOINT','OTEL_RESOURCE_ATTRIBUTES','OTLP_GOOGLE_CLOUD_PROJECT']:ENV.pop(k,None)
ledger=[]
def run(name,argv):
 start=time.time()
 with (E/(name+'.log')).open('wb') as f:
  try: rc=subprocess.run(argv,cwd=W,env=ENV,stdout=f,stderr=subprocess.STDOUT,timeout=150).returncode
  except subprocess.TimeoutExpired:rc=124
 d={'name':name,'argv':argv,'cwd':str(W),'rc':rc,'seconds':round(time.time()-start,3)};ledger.append(d);(E/'gates.json').write_text(json.dumps(ledger,indent=2));print(json.dumps(d),flush=True)
if __name__=='__main__':
 run('regressions',['go','test','./cmd/ailang-worldd','./host/boundary','-run','^(TestQuickstartSection11FlagsMatchTheCLI|TestWorldBoundaryDependencyAllowlist|TestWorldBoundaryRecordingWriter)$','-count=1','-timeout=90s','-v'])
 run('subprocess-audit',['go','test','./host/broker','-run','^TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential$','-count=1','-timeout=90s','-v'])
 for n,p in [('TestDemoSeedAuthenticatedCommits','./cmd/ailang-worldd'),('TestDemoSeedPreflightRefusals','./cmd/ailang-worldd'),('TestDemoSeedHelpPreparation','./cmd/ailang-worldd'),('TestDemoSeedCommitFailures','./cmd/ailang-worldd'),('TestSessionMint_RequiresControllingTerminal','./cmd/ailang-worldd'),('TestRESTCommitRefusesDecisionPacketObjects','./host/daemon')]:
  run(n,['go','test',p,'-run','^'+n+'$','-count=1','-timeout=120s','-v'])
