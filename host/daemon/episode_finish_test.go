package daemon

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/hashref"
)

func finishProofFile(t *testing.T, d *Daemon, phase string, expires int64, mode os.FileMode) string {
	t.Helper()
	proof := strings.Repeat("a", 64)
	dir := filepath.Join(filepath.Dir(d.credentialsPath), "operator-proofs")
	if err := os.MkdirAll(dir, 0o700); err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(dir, "ep1.json")
	data := mustMarshal(t, operatorProof{SHA256: hashref.SumSHA256([]byte(proof)).Digest(), Phase: phase, ExpiresAt: expires})
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.Chmod(path, mode); err != nil {
		t.Fatal(err)
	}
	return proof
}

func finishRoute(t *testing.T, d *Daemon, input FinishRequest, proof, authorization string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodPost, "/v1/episodes/ep1/finish", bytes.NewReader(mustMarshal(t, input)))
	if proof != "" {
		req.Header.Set("X-World-Operator-Proof", proof)
	}
	if authorization != "" {
		req.Header.Set("Authorization", authorization)
	}
	if d.isProtected(req) {
		t.Fatal("operator route is session gated")
	}
	w := httptest.NewRecorder()
	d.Handler().ServeHTTP(w, req)
	return w
}

func testFinishScope() string {
	return broker.FormatFinishApprovalScope(broker.FinishApprovalScope{
		Repository: "o/r", Remote: "https://github.com/o/r.git", Head: "world/ep1", SHA: strings.Repeat("a", 40),
		Base: "dev", TitleSHA256: strings.Repeat("b", 64), BodySHA256: strings.Repeat("c", 64),
		Draft: true, ExpiresAt: time.Now().Unix() + 900,
	})
}

func TestEpisodeFinishOperatorProofAC11(t *testing.T) {
	d := newHandlerDaemon(t)
	token := mintSessionGrants(t, d, "ep1", "Workspace.Write")
	input := FinishRequest{Phase: "approve", Scope: testFinishScope()}
	for _, name := range []string{"session", "missing", "wrong", "expired", "future", "readable", "world-readable", "symlink", "phase"} {
		t.Run(name, func(t *testing.T) {
			proof, auth := "", ""
			path := filepath.Join(filepath.Dir(d.credentialsPath), "operator-proofs", "ep1.json")
			os.Remove(path)
			switch name {
			case "session":
				auth = "Bearer " + token
			case "wrong":
				finishProofFile(t, d, "approve", time.Now().Unix()+100, 0o600)
				proof = strings.Repeat("b", 64)
			case "expired":
				proof = finishProofFile(t, d, "approve", time.Now().Unix()-1, 0o600)
			case "future":
				proof = finishProofFile(t, d, "approve", time.Now().Unix()+301, 0o600)
			case "readable":
				proof = finishProofFile(t, d, "approve", time.Now().Unix()+100, 0o640)
			case "world-readable":
				proof = finishProofFile(t, d, "approve", time.Now().Unix()+100, 0o644)
			case "symlink":
				proof = finishProofFile(t, d, "approve", time.Now().Unix()+100, 0o600)
				os.Rename(path, path+".target")
				os.Symlink(path+".target", path)
			case "phase":
				proof = finishProofFile(t, d, "record", time.Now().Unix()+100, 0o600)
			}
			w := finishRoute(t, d, input, proof, auth)
			if w.Code != 403 || !strings.Contains(w.Body.String(), "why:") || !strings.Contains(w.Body.String(), "next:") {
				t.Fatalf("%d %s", w.Code, w.Body.String())
			}
			for _, head := range []string{broker.ApprovalsV1, finishHead("ep1")} {
				if _, exists, err := d.store.GetRegistryHead(boundedTestContext(t), head); err != nil || exists {
					t.Fatalf("refusal wrote %s: %v %v", head, exists, err)
				}
			}
			os.Remove(path)
		})
	}
	proof := finishProofFile(t, d, "approve", time.Now().Unix()+100, 0o600)
	w := finishRoute(t, d, input, proof, "")
	if w.Code != 200 {
		t.Fatalf("%d %s", w.Code, w.Body.String())
	}
	var approval FinishApproval
	if err := json.Unmarshal(w.Body.Bytes(), &approval); err != nil {
		t.Fatal(err)
	}
	if w := finishRoute(t, d, input, proof, ""); w.Code != 403 {
		t.Fatalf("replay: %d", w.Code)
	}
	record := &EpisodeFinish{Episode: "ep1", Branch: "world/ep1", Commit: strings.Repeat("a", 40),
		Remote: "https://github.com/o/r.git", Approval: approval.Approval, Outcome: "done",
		PR: FinishPR{URL: "https://github.com/o/r/pull/1", Number: 1, Draft: true}}
	input = FinishRequest{Phase: "record", Record: record}
	for _, wrong := range []string{"branch", "commit", "remote", "episode", "approval"} {
		copy := *record
		switch wrong {
		case "branch":
			copy.Branch = "other"
		case "commit":
			copy.Commit = strings.Repeat("f", 40)
		case "remote":
			copy.Remote = "https://github.com/other/repo"
		case "episode":
			copy.Episode = "ep2"
		case "approval":
			copy.Approval = approval.Request
		}
		proof = finishProofFile(t, d, "record", time.Now().Unix()+100, 0o600)
		w = finishRoute(t, d, FinishRequest{Phase: "record", Record: &copy}, proof, "")
		if w.Code != 403 {
			t.Fatalf("%s: %d %s", wrong, w.Code, w.Body.String())
		}
		if _, exists, _ := d.store.GetRegistryHead(boundedTestContext(t), finishHead("ep1")); exists {
			t.Fatal("invalid record moved head")
		}
	}
	for _, class := range []string{"push-failed", "pr-failed", "ambiguous", "done"} {
		record.Outcome = class
		proof = finishProofFile(t, d, "record", time.Now().Unix()+100, 0o600)
		w = finishRoute(t, d, input, proof, "")
		if w.Code != 200 {
			t.Fatalf("%s: %d %s", class, w.Code, w.Body.String())
		}
	}
	proof = finishProofFile(t, d, "record", time.Now().Unix()+100, 0o600)
	if w := finishRoute(t, d, input, proof, ""); w.Code != 409 {
		t.Fatalf("second done: %d %s", w.Code, w.Body.String())
	}
	w = requestRecorder(t, d, "GET", "/v1/episodes/ep1/report", nil)
	var report EpisodeReport
	if w.Code != 200 || json.Unmarshal(w.Body.Bytes(), &report) != nil || report.Finish == nil ||
		report.Finish.PR.URL != record.PR.URL || report.Finish.Previous == "" {
		t.Fatalf("report: %d %s", w.Code, w.Body.String())
	}
}

func TestEpisodeFinishScopeExpiry(t *testing.T) {
	d := newHandlerDaemon(t)
	for _, expiry := range []int64{time.Now().Unix() - 1, time.Now().Unix() + 3601} {
		scope, _ := broker.ParseFinishApprovalScope(testFinishScope())
		scope.ExpiresAt = expiry
		proof := finishProofFile(t, d, "approve", time.Now().Unix()+100, 0o600)
		w := finishRoute(t, d, FinishRequest{Phase: "approve", Scope: broker.FormatFinishApprovalScope(scope)}, proof, "")
		if w.Code != 403 {
			t.Fatalf("expiry %d: %d", expiry, w.Code)
		}
	}
}

type unreadFinishBody struct {
	t *testing.T
}

func (b unreadFinishBody) Read([]byte) (int, error) {
	b.t.Fatal("operator route read a session-authorized body")
	return 0, nil
}

func TestEpisodeFinishAuthorizationPrecedesBody(t *testing.T) {
	d := newHandlerDaemon(t)
	token := mintSessionGrants(t, d, "ep1", "Workspace.Write")
	for _, header := range []string{"Bearer " + token, ""} {
		r := httptest.NewRequest("POST", "/v1/episodes/ep1/finish", unreadFinishBody{t})
		r.Header.Set("Authorization", header)
		w := httptest.NewRecorder()
		d.Handler().ServeHTTP(w, r)
		if w.Code != 403 || !strings.Contains(w.Body.String(), "operator route: session credentials are refused") {
			t.Fatalf("%d %s", w.Code, w.Body.String())
		}
	}
}
