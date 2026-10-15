exec(open('/tmp/world268-evaluator-r1/run.py').read().split("for name in ['AuthenticatedCommits'")[0])
results=[]
run('independent-subprocess-audit-red','^TestEverySubprocessSiteIsDrivenAndScrubsTheRegistryCredential$','./host/broker')
f=root/'cmd/ailang-worldd/evaluator_r1_probe_test.go'
s='''package main
import (
 "encoding/json"
 "net/http"
 "net/http/httptest"
 "os"
 "path/filepath"
 "testing"
 "github.com/sunholo-data/ailang-world/host/broker"
 "github.com/sunholo-data/ailang-world/host/daemon"
 "github.com/sunholo-data/ailang-world/host/hashref"
)
func TestEvaluatorR1VersionProbeInheritsCredential(t *testing.T) {
 t.Setenv(broker.RegistryCredentialVariable,"EVALUATOR_R1_SENTINEL_ONLY")
 dir:=t.TempDir(); pin:=filepath.Join(dir,"pin"); dump:=filepath.Join(dir,"sentinel.txt")
 script:=[]byte("#!/bin/sh\\nprintf '%s' \\\"$AILANG_REGISTRY_API_KEY\\\" > "+dump+"\\nprintf 'AILANG v0.41.0\\\\n'\\n")
 if err:=os.WriteFile(pin,script,0700);err!=nil {t.Fatal(err)}
 db:=filepath.Join(dir,"missing.db")
 h:=daemon.HealthResponse{DBPath:db,InterpreterRef:hashref.SumSHA256(script).String(),InterpreterVersion:"AILANG v0.41.0\\n"}
 server:=httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter,r *http.Request){_ = json.NewEncoder(w).Encode(h)}));defer server.Close()
 ctx:=boundedTestContext(t)
 _,_,_ = demoPreflight(ctx,demoOptions{db:db,pin:pin,manifest:filepath.Join(dir,"new.json")},"unused",worldReader{ctx:ctx,c:newClient(server.URL)})
 b,err:=os.ReadFile(dump);if err!=nil {t.Fatalf("probe did not run: %v",err)}
 if string(b)!="EVALUATOR_R1_SENTINEL_ONLY" {t.Fatalf("sentinel not inherited: %q",b)}
 t.Log("CONFIRMED actual demoPreflight version child inherited the sentinel registry credential; no real credential used")
}
'''
try:
 f.write_text(s);run('independent-sentinel-leak-observation','^TestEvaluatorR1VersionProbeInheritsCredential$')
finally:f.unlink(missing_ok=True)
(out/'subprocess-ledger.json').write_text(json.dumps(results,indent=2));print('SUBPROCESS TERMINAL RESTORED',flush=True)
