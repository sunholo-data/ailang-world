package broker

// Row 201 A4 sprint 1 (design §2.5, §2.6, §6 DL1–DL12, DL4a, DL4b, R4a): the decision law, one
// test per refusal (the unit is the branch), its order, its must-live control, and L1.

import (
	"crypto/sha256"
	"strings"
	"testing"
)

const lawRef = "sha256:00112233445566778899aabbccddeeff00112233445566778899aabbccddeeff"

// lawFacts is a decidable approve: every fact holds.
func lawFacts() DecisionFacts {
	p, q := sha256.Sum256([]byte("p")), sha256.Sum256([]byte("q"))
	return DecisionFacts{Ref: lawRef, RefOK: true, RequestOK: true, OnChain: true,
		Discovery: Discovery{Complete: true, Scanned: 3, Carrying: 1}, Scope: "worktree",
		Now: 1000, Deadline: 5000, Packet: p, Request: q, PostedPacket: p, PostedRequest: q, Verb: "approve", Expiry: 1600}
}

// wantLaw fails unless the law refuses f at step with status and a reason containing each of want.
func wantLaw(t *testing.T, what string, f DecisionFacts, step string, status int, want ...string) {
	t.Helper()
	r := DecisionRefusal(f)
	if r.Step != step || r.Status != status {
		t.Fatalf("%s: %+v, want %s %d", what, r, step, status)
	}
	for _, w := range want {
		if !strings.Contains(r.Reason, w) {
			t.Errorf("%s: reason %q lacks %q", what, r.Reason, w)
		}
	}
}

func TestDecisionRefusalD1(t *testing.T) {
	f := lawFacts()
	f.RefOK, f.Ref = false, "sha256:XYZ"
	wantLaw(t, "D1", f, "D1", 400, `"sha256:XYZ" is not a request ref`)
}

func TestDecisionRefusalD2(t *testing.T) {
	f := lawFacts()
	f.RequestOK = false
	wantLaw(t, "D2", f, "D2", 404, "no such approval request "+lawRef)
}

func TestDecisionRefusalD3(t *testing.T) {
	f := lawFacts()
	f.OnChain = false
	wantLaw(t, "D3", f, "D3", 404, "not on the approvals chain")
}

// DL4a: a scan that hit its sentinel, or saw an undecodable packet, is refused even when it saw
// exactly one carrying packet (one visible match is never exactly one match).
func TestDecisionRefusalD4a(t *testing.T) {
	f := lawFacts()
	f.Discovery = Discovery{Complete: false, Scanned: 2000, Carrying: 1}
	wantLaw(t, "incomplete", f, "D4a", 409, "cannot be counted completely", "more than 2000")
	f.Discovery = Discovery{Complete: true, Scanned: 4, Carrying: 1, Undecodable: 1}
	wantLaw(t, "undecodable", f, "D4a", 409, "cannot be counted completely")
	f.Discovery = Discovery{Complete: false, Scanned: 2000, Carrying: 2}
	wantLaw(t, "incomplete before the count", f, "D4a", 409)
}

// DL4b: a complete count of 0 or 2, or one carrying packet that does not validate.
func TestDecisionRefusalD4b(t *testing.T) {
	f := lawFacts()
	f.Discovery.Carrying = 0
	wantLaw(t, "zero", f, "D4b", 409, "no committed decision packet carries request "+lawRef)
	f.Discovery.Carrying = 2
	wantLaw(t, "two", f, "D4b", 409, "is carried by 2 decision packets")
	f = lawFacts()
	f.PacketCode, f.PacketSentence = 5, "its approval request is not among the newest 40 approval heads, so its approvals-chain membership is not established"
	wantLaw(t, "code 5", f, "D4b", 409, "is refused: ", "not established")
}

func TestDecisionRefusalD5(t *testing.T) {
	f := lawFacts()
	f.Scope = "registry:https://r.example/package:v/n/version:1.0.0"
	wantLaw(t, "publish", f, "D5", 403, "world-publish approve", "(R4)")
	f.Scope = "github-pr:o/r"
	wantLaw(t, "finish", f, "D5", 403, "ailang-worldd episode finish")
}

func TestDecisionRefusalD6(t *testing.T) {
	f := lawFacts()
	f.Decided, f.DecidedVerb, f.DecidedBy = true, "deny", "operator:terminal"
	wantLaw(t, "D6", f, "D6", 409, "request "+lawRef+" was already decided deny by operator:terminal")
}

func TestDecisionRefusalD7(t *testing.T) {
	f := lawFacts()
	f.Now = f.Deadline
	wantLaw(t, "at the deadline", f, "D7", 409, "deadline passed at 1970-01-01T01:23:20Z", "Cancel")
	f.Now = f.Deadline - 1
	f.Expiry = f.Deadline
	if r := DecisionRefusal(f); r.Step != "" {
		t.Errorf("one second before the deadline (control): %+v", r)
	}
}

func TestDecisionRefusalD8(t *testing.T) {
	f := lawFacts()
	f.PostedPacket[0] ^= 1
	wantLaw(t, "P", f, "D8", 409, "reload it")
	f = lawFacts()
	f.PostedRequest[31] ^= 1
	wantLaw(t, "Q", f, "D8", 409, "reload it")
}

func TestDecisionRefusalD9(t *testing.T) {
	for _, verb := range []string{"", "Approve", "yes", "approve "} {
		f := lawFacts()
		f.Verb = verb
		wantLaw(t, "verb "+verb, f, "D9", 400, "a decision is approve or deny")
	}
}

// D10: now < expiry <= min(now + 600, deadline), each edge on both sides.
func TestDecisionRefusalD10(t *testing.T) {
	for _, tc := range []struct {
		now, deadline, expiry int64
		ok                    bool
	}{
		{1000, 5000, 1000, false}, {1000, 5000, 1001, true}, {1000, 5000, 1600, true}, {1000, 5000, 1601, false},
		{1000, 1300, 1300, true}, {1000, 1300, 1301, false}, {1000, 5000, 999, false},
	} {
		f := lawFacts()
		f.Now, f.Deadline, f.Expiry = tc.now, tc.deadline, tc.expiry
		r := DecisionRefusal(f)
		if tc.ok != (r.Step == "") || !tc.ok && (r.Step != "D10" || r.Status != 409 || !strings.Contains(r.Reason, "at most 10 minutes")) {
			t.Errorf("now %d deadline %d expiry %d: %+v, want ok=%v", tc.now, tc.deadline, tc.expiry, r, tc.ok)
		}
	}
	if got := DecisionWindowMax(1000, 5000); got != 1600 {
		t.Errorf("window max = %d, want 1600", got)
	}
	if got := DecisionWindowMax(1000, 1300); got != 1300 {
		t.Errorf("window max = %d, want 1300", got)
	}
}

// DL11: with two facts false, the lower step wins (each adjacent pair, D1 ... D10).
func TestDecisionRefusalOrder(t *testing.T) {
	breaks := []func(*DecisionFacts){
		func(f *DecisionFacts) { f.RefOK = false },
		func(f *DecisionFacts) { f.RequestOK = false },
		func(f *DecisionFacts) { f.OnChain = false },
		func(f *DecisionFacts) { f.Discovery.Complete = false },
		func(f *DecisionFacts) { f.Discovery.Carrying = 2 },
		func(f *DecisionFacts) { f.Scope = "github-pr:o/r" },
		func(f *DecisionFacts) { f.Decided = true },
		func(f *DecisionFacts) { f.Now = f.Deadline + 1; f.Expiry = f.Deadline + 2 },
		func(f *DecisionFacts) { f.PostedPacket[0] ^= 1 },
		func(f *DecisionFacts) { f.Verb = "maybe" },
		func(f *DecisionFacts) { f.Expiry = 1 },
	}
	steps := []string{"D1", "D2", "D3", "D4a", "D4b", "D5", "D6", "D7", "D8", "D9", "D10"}
	for i := 0; i+1 < len(breaks); i++ {
		for j := i + 1; j < len(breaks); j++ {
			f := lawFacts()
			breaks[j](&f)
			breaks[i](&f)
			if r := DecisionRefusal(f); r.Step != steps[i] {
				t.Errorf("%s and %s false: %s won", steps[i], steps[j], r.Step)
			}
		}
	}
	// The subject law is the decision law's first seven steps, and stops there.
	f := lawFacts()
	f.Verb, f.PostedPacket = "maybe", [32]byte{}
	if r := SubjectRefusal(f); r.Step != "" {
		t.Errorf("SubjectRefusal reads a posted fact: %+v", r)
	}
	f.Decided = true
	if r := SubjectRefusal(f); r.Step != "D6" {
		t.Errorf("SubjectRefusal D6: %+v", r)
	}
}

// DL12: a well-formed approve and a well-formed deny are accepted (the control).
func TestDecisionRefusalMustLive(t *testing.T) {
	for _, verb := range []string{"approve", "deny"} {
		f := lawFacts()
		f.Verb = verb
		if r := DecisionRefusal(f); r != (LawRefusal{}) {
			t.Errorf("%s: %+v, want accepted", verb, r)
		}
		if r := SubjectRefusal(f); r != (LawRefusal{}) {
			t.Errorf("subject %s: %+v, want accepted", verb, r)
		}
	}
}

// R4a: L1's publish and finish shapes, their near misses, and unreserved controls.
func TestReservedApprovalScope(t *testing.T) {
	id := PublishIdentity{RegistryOrigin: "https://registry.example", Vendor: "v", Name: "n", Version: "1.0.0"}
	for _, tc := range []struct{ scope, owner string }{
		{PublishApprovalScopeFor(id, PublishHashes{TarballSHA256: "t", ContentHash: "c", InterfaceHash: "i"}, 9), "publish"},
		{"pkg#effect=Other&manifest=m&tarball=t&content=c&interface=i&expires=9", "publish"}, // the parse arm alone
		{"registry:x", "publish"},                           // the prefix arm alone
		{"worktree#effect=Registry.Publish", "publish"},     // the effect arm alone
		{"worktree#effect=Registry.Publish&x=1", "publish"}, // the effect arm, not at the end (sprint-1 judge J07)
		{FormatFinishApprovalScope(FinishApprovalScope{Repository: "o/r", Remote: "https://github.com/o/r.git", Head: "world/e", SHA: strings.Repeat("a", 40),
			Base: "dev", TitleSHA256: strings.Repeat("b", 64), BodySHA256: strings.Repeat("c", 64), ExpiresAt: 9}), "finish"},
		{"github-pr:x", "finish"}, // the finish prefix arm alone
		{"worktree", ""}, {"world", ""}, {"", ""}, {"worktree#x", ""}, {"Registry.Publish", ""},
	} {
		owner, reserved := ReservedApprovalScope(tc.scope)
		if owner != tc.owner || reserved != (tc.owner != "") {
			t.Errorf("%q: %q %v, want %q", tc.scope, owner, reserved, tc.owner)
		}
	}
}
