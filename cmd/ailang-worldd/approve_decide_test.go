package main

// Row 201 A4 sprint 2 M4c (design §2.10, §6 CD1, CD4): the typed decision end to end against the
// fake daemon (approve_decide_rig_test.go), and the subject checks before any prompt.

import (
	"encoding/base64"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/coordinator"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

// CD1: y at the terminal → one POST: typed, the verb, both digests, the window's max as expiry,
// the daemon's loopback Origin, the proof whose file (phase decide-<verb>, named for the request)
// hashes to it; every field was shown with the verb in capitals; --deny denies; n posts nothing.
func TestApproveDecideTyped(t *testing.T) {
	for _, tc := range []struct{ verb, answer string }{{"approve", "y\n"}, {"deny", "y\n"}, {"approve", "n\n"}} {
		f := newFakeDecider(t)
		term := &decideTerm{r: strings.NewReader(tc.answer)}
		var extra []string
		if tc.verb == "deny" {
			extra = []string{"--deny"}
		}
		code, stdout, stderr, _ := f.run(t, term, extra...)
		prompt := term.out.String()
		for _, want := range []string{f.ref, "effect               Human.Approve", "scope                worktree", "cost                 1",
			"requester            agent", "now                  1700", f.subject.Packet + " (log entry 7)", "proposalHash         " + f.subject.Proposal,
			"requestRef           " + f.ref, "createdAt            1700", "deadlineAt           1760 (1970-01-01T00:29:20Z)", "escalationsRemaining 0",
			"policy               Cancel", strings.ToUpper(tc.verb) + " this request?"} {
			if !strings.Contains(prompt, want) {
				t.Errorf("%s %q: the prompt lacks %q:\n%s", tc.verb, tc.answer, want, prompt)
			}
		}
		proof, wrote := f.proofFile(t)
		if tc.answer == "n\n" {
			if code != exitUsage || len(f.posts) != 0 || wrote {
				t.Fatalf("n: rc %d, %d POSTs, proof written %v", code, len(f.posts), wrote)
			}
			continue
		}
		if code != exitOK || len(f.posts) != 1 {
			t.Fatalf("%s: rc %d, %d POSTs\nstderr %s", tc.verb, code, len(f.posts), stderr)
		}
		b, r := f.postBodies[0], f.posts[0]
		want := surfacepage.DecideBody{V: 1, Verb: tc.verb, Packet: strings.TrimPrefix(f.subject.Packet, "sha256:"), Request: strings.TrimPrefix(f.ref, "sha256:"), Expiry: 1760, Typed: true}
		if b != want || r.Header.Get("Origin") != f.srv.URL || r.Header.Get("Content-Type") != "application/json" {
			t.Fatalf("%s: POST %+v (Origin %q), want %+v from %s", tc.verb, b, r.Header.Get("Origin"), want, f.srv.URL)
		}
		if !wrote || proof.Phase != "decide-"+tc.verb || proof.SHA256 != finishSHA([]byte(r.Header.Get("X-World-Operator-Proof"))) {
			t.Fatalf("%s: proof %+v (written %v) does not prove the POST's header", tc.verb, proof, wrote)
		}
		if !strings.Contains(stdout, "decided "+tc.verb+": request "+f.ref) || !strings.Contains(stdout, "by operator:terminal") {
			t.Errorf("%s: stdout %q", tc.verb, stdout)
		}
	}
	f := newFakeDecider(t) // a refused POST is reported with the daemon's reason
	f.postStatus, f.postBody = 409, `{"v":1,"error":"Conflict","reason":"request was already decided deny by operator:terminal"}`
	if code, _, stderr, _ := f.run(t, &decideTerm{r: strings.NewReader("y\n")}); code != exitUsage || !strings.Contains(stderr, "(HTTP 409): request was already decided deny") {
		t.Fatalf("refused POST: rc %d, stderr %s", code, stderr)
	}
}

// CD4: the CLI checks the daemon's answer before any prompt: request bytes of another request,
// packet bytes that do not hash to the packet ref, a packet carrying another request, a request
// that does not decode, an unknown field, a daemon refusal, a display refusal, a daemon serving
// another store or credentials store → refused, no terminal opened, nothing posted.
func TestApproveDecideRefusesAMismatchedSubject(t *testing.T) {
	otherReq, otherPkt, _ := decisionObjects("world")
	_, _, q := decisionObjects("worktree")
	samePkt := coordinator.EncodeDecisionPacket(coordinator.DecisionPacket{ProposalHash: hashref.SumSHA256([]byte("another plan")), RequestRef: q, CreatedAt: 1700, DeadlineAt: 1760, Policy: "Cancel"})
	garbage := []byte(`{"effect":"Human.Approve","extra":1}`)
	garbagePkt := coordinator.EncodeDecisionPacket(coordinator.DecisionPacket{ProposalHash: q, RequestRef: hashref.SumSHA256(garbage), CreatedAt: 1700, DeadlineAt: 1760, Policy: "Cancel"})
	for _, tc := range []struct {
		name string
		edit func(f *fakeDecider)
		want string
	}{
		{"request bytes", func(f *fakeDecider) { f.subject.RequestBytes = base64.RawURLEncoding.EncodeToString(otherReq) }, "do not match their refs"},
		{"packet bytes", func(f *fakeDecider) { f.subject.PacketBytes = base64.RawURLEncoding.EncodeToString(samePkt) }, "do not match their refs"},
		{"another request's packet", func(f *fakeDecider) {
			f.subject.PacketBytes, f.subject.Packet = base64.RawURLEncoding.EncodeToString(otherPkt), hashref.SumSHA256(otherPkt).String()
		}, "do not match their refs"},
		{"undecodable request", func(f *fakeDecider) {
			f.ref = hashref.SumSHA256(garbage).String()
			f.subject.Request, f.subject.RequestBytes = f.ref, base64.RawURLEncoding.EncodeToString(garbage)
			f.subject.PacketBytes, f.subject.Packet = base64.RawURLEncoding.EncodeToString(garbagePkt), hashref.SumSHA256(garbagePkt).String()
		}, "do not match their refs"},
		{"undecodable packet", func(f *fakeDecider) {
			f.subject.PacketBytes, f.subject.Packet = base64.RawURLEncoding.EncodeToString(garbage), hashref.SumSHA256(garbage).String()
		}, "do not match their refs"},
		{"bad base64", func(f *fakeDecider) { f.subject.RequestBytes = "!" }, "do not match their refs"},
		{"unknown field", func(f *fakeDecider) { f.getBody = `{"v":1,"message":"AAAA"}` }, "is not a decision subject"},
		{"daemon refusal", func(f *fakeDecider) { f.getStatus = 409 }, "will not show request"},
		{"display refusal", func(f *fakeDecider) { f.subject.Refusal = "the packet's deadline passed" }, "cannot be decided now: the packet's deadline passed"},
		{"another store", func(f *fakeDecider) { f.health.DBPath = filepath.Join(f.dir, "other.db") }, "serves store"},
		{"another credentials store", func(f *fakeDecider) { f.health.CredentialsDB = filepath.Join(f.dir, "other-credentials.db") }, "serves credentials store"},
	} {
		f := newFakeDecider(t)
		tc.edit(f)
		code, _, stderr, opened := f.run(t, &decideTerm{r: strings.NewReader("y\n")})
		if code != exitUsage || opened || len(f.posts) != 0 || !strings.Contains(stderr, tc.want) {
			t.Errorf("%s: rc %d, opened %v, %d POSTs\nstderr %s", tc.name, code, opened, len(f.posts), stderr)
		}
	}
}
