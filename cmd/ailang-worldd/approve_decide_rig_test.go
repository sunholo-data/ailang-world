package main

// Row 201 A4 sprint 2 M4b (design §2.10, §6 CD2, CD3): approve <request-ref> [--deny], typed,
// against a fake daemon (httptest on the loopback) whose subject is built with the daemon's own
// encoders and JSON types (coordinator.EncodeDecisionPacket, surfacepage.DecisionSubject); the
// real daemon side of the same wire is host/daemon's RD2-RD9 (D-WORLD-109: temp dirs only, no
// real store, no /dev/tty).

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/coordinator"
	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

// decideTerm is a terminal answering answer and keeping what was written to it.
type decideTerm struct {
	r   io.Reader
	out bytes.Buffer
}

func (d *decideTerm) Read(p []byte) (int, error)  { return d.r.Read(p) }
func (d *decideTerm) Write(p []byte) (int, error) { return d.out.Write(p) }
func (d *decideTerm) Close() error                { return nil }

// fakeDecider is a daemon serving /v1/health, GET /surface/decisions/{ref} and POST
// /surface/decide/{ref} for one request; tests tamper with subject, health and the answers.
type fakeDecider struct {
	srv        *httptest.Server
	dir, db    string
	creds      string
	ref        string
	subject    surfacepage.DecisionSubject
	getStatus  int
	getBody    string
	health     daemon.HealthResponse
	gets       int
	posts      []*http.Request
	postBodies []surfacepage.DecideBody
	postStatus int
	postBody   string
}

// decisionObjects are a Human.Approve request of scope and a packet carrying it, as stored bytes.
func decisionObjects(scope string) (req, pkt []byte, q hashref.HashRef) {
	req = []byte(fmt.Sprintf(`{"effect":"Human.Approve","scope":%q,"cost":1,"requester":"agent","now":1700}`, scope))
	q = hashref.SumSHA256(req)
	pkt = coordinator.EncodeDecisionPacket(coordinator.DecisionPacket{ProposalHash: hashref.SumSHA256([]byte("plan " + scope)), RequestRef: q,
		CreatedAt: 1700, DeadlineAt: 1760, Policy: "Cancel"})
	return req, pkt, q
}

func newFakeDecider(t *testing.T) *fakeDecider {
	t.Helper()
	f := &fakeDecider{dir: t.TempDir(), getStatus: 200, postStatus: 200}
	f.db, f.creds = filepath.Join(f.dir, "world.db"), filepath.Join(f.dir, "credentials.db")
	for _, p := range []string{f.db, f.creds} {
		if err := os.WriteFile(p, nil, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	req, pkt, q := decisionObjects("worktree")
	f.ref = q.String()
	f.subject = surfacepage.DecisionSubject{V: 1, Request: f.ref, RequestBytes: base64.RawURLEncoding.EncodeToString(req),
		Packet: hashref.SumSHA256(pkt).String(), PacketBytes: base64.RawURLEncoding.EncodeToString(pkt), Entry: 7,
		Proposal: hashref.SumSHA256([]byte("plan worktree")).String(), Now: 1710, Window: surfacepage.DecisionWindow{Min: 1711, Max: 1760},
		SSH: surfacepage.DecisionSSH{Namespace: "world-decision@sunholo"}}
	f.health = daemon.HealthResponse{DBPath: f.db, CredentialsDB: f.creds}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /v1/health", func(w http.ResponseWriter, _ *http.Request) { _ = json.NewEncoder(w).Encode(f.health) })
	mux.HandleFunc("GET /surface/decisions/{ref}", func(w http.ResponseWriter, r *http.Request) {
		f.gets++
		w.WriteHeader(f.getStatus)
		if f.getStatus != 200 {
			_, _ = w.Write([]byte(`{"v":1,"error":"Conflict","reason":"the packet's deadline passed"}`))
			return
		}
		if f.getBody != "" {
			_, _ = w.Write([]byte(f.getBody))
			return
		}
		_ = json.NewEncoder(w).Encode(f.subject)
	})
	mux.HandleFunc("POST /surface/decide/{ref}", func(w http.ResponseWriter, r *http.Request) {
		var b surfacepage.DecideBody
		raw, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(raw, &b)
		f.posts, f.postBodies = append(f.posts, r), append(f.postBodies, b)
		w.WriteHeader(f.postStatus)
		if f.postBody != "" {
			_, _ = w.Write([]byte(f.postBody))
			return
		}
		_ = json.NewEncoder(w).Encode(surfacepage.DecisionMade{V: 1, Decided: b.Verb, Request: "sha256:" + b.Request, Decision: "sha256:" + strings.Repeat("d", 64), DecidedBy: "operator:terminal"})
	})
	f.srv = httptest.NewServer(mux)
	t.Cleanup(f.srv.Close)
	return f
}

func (f *fakeDecider) addr() string { return strings.TrimPrefix(f.srv.URL, "http://") }

// run is approve <ref> --db --daemon extra... with term (nil: no terminal); it returns the code,
// stdout and stderr, and whether a terminal was opened.
func (f *fakeDecider) run(t *testing.T, term io.ReadWriteCloser, extra ...string) (int, string, string, bool) {
	t.Helper()
	opened := false
	env := approveEnv(term, &opened, atTime(1710))
	if term == nil {
		env.openTerminal = func() (io.ReadWriteCloser, error) { opened = true; return nil, errors.New("no tty") }
	}
	var stdout, stderr bytes.Buffer
	code := runApprove(append([]string{f.ref, "--db", f.db, "--daemon", f.addr()}, extra...), &stdout, &stderr, env)
	return code, stdout.String(), stderr.String(), opened
}

// proofFile is the operator proof the CLI wrote for the request, or ok false.
func (f *fakeDecider) proofFile(t *testing.T) (p struct {
	SHA256    string `json:"sha256"`
	Phase     string `json:"phase"`
	ExpiresAt int64  `json:"expires_at"`
}, ok bool) {
	t.Helper()
	data, err := os.ReadFile(filepath.Join(f.dir, "operator-proofs", "decision-"+strings.TrimPrefix(f.ref, "sha256:")+".json"))
	if err != nil {
		return p, false
	}
	if err := json.Unmarshal(data, &p); err != nil {
		t.Fatal(err)
	}
	return p, true
}

// CD2: no terminal → refused naming the key flags, nothing posted, no proof written.
func TestApproveDecideNoTerminal(t *testing.T) {
	f := newFakeDecider(t)
	code, _, stderr, opened := f.run(t, nil)
	if _, wrote := f.proofFile(t); code != exitUsage || !opened || len(f.posts) != 0 || wrote || !strings.Contains(stderr, "--sign-with <key>") {
		t.Fatalf("no terminal: rc %d, opened %v, %d POSTs, proof %v\nstderr %s", code, opened, len(f.posts), wrote, stderr)
	}
}

// CD3: --deny and --daemon need a request ref (an mr- id, --list); --await, and an assertion
// without --expiry (CD9), are refused with a request ref; two modes are F1. Each before any daemon call.
func TestApproveDecideFlagRules(t *testing.T) {
	f := newFakeDecider(t)
	for _, tc := range []struct {
		args []string
		want string
	}{
		{[]string{"mr-0123456789abcdef", "--deny"}, "a mint request has no deny"},
		{[]string{f.ref + "0"}, "is not a mint request id"}, // the decision mode's ref is anchored
		{[]string{"mr-0123456789abcdef", "--daemon", f.addr()}, "a mint request has no deny"},
		{[]string{"--list", "--deny"}, "a mint request has no deny"},
		{[]string{f.ref, "--await"}, "there is nothing to await"},
		{[]string{f.ref, "--webauthn-assertion", "x"}, "go with --expiry <t>"},
		{[]string{f.ref, "--signature", "x", "--sign-with", "y"}, "usage: ailang-worldd approve"},
	} {
		var stdout, stderr bytes.Buffer
		opened := false
		code := runApprove(append(tc.args, "--db", f.db), &stdout, &stderr, approveEnv(nil, &opened, atTime(1710)))
		if code != exitUsage || opened || f.gets != 0 || !strings.Contains(stderr.String(), tc.want) {
			t.Errorf("%v: rc %d, opened %v, %d GETs\nstderr %s", tc.args, code, opened, f.gets, stderr.String())
		}
	}
}
