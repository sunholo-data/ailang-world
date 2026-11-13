package main

// Row 201 A4 sprint 4 (design §2.8, §6 XD2): the decision card fails closed, in a live Chrome, gated
// as XB1 (WORLD_CHROME, WORLD_REQUIRE_PASSKEY_BROWSER=1; skipped and counted without them).

import (
	"encoding/base64"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/coordinator"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

// XD2: the card fails closed. A subject whose bytes do not hash to the refs, a packet carrying
// another request, a key the table does not list in either object, a listed key missing, a value of
// another type, bytes that are not base64url and a malformed window: refused with the reason, no
// fields, no decide control, no get. A reserved, a refused and an already decided subject render
// their fields but offer no decide control. Control: the untampered subject is ready.
func TestDecidePageFailsClosed(t *testing.T) {
	b := startPasskeyBrowser(t)
	b.do(t, nil, "tab", map[string]any{"authenticator": "uv", "inject": "getRecorder"})
	r := newDecideCardRig(t)
	other, _ := cardObjects("worktree:other", "agent", 2400)
	packetFor := func(req []byte, extra string) []byte {
		p := coordinator.EncodeDecisionPacket(coordinator.DecisionPacket{ProposalHash: hashref.SumSHA256([]byte("plan")), RequestRef: hashref.SumSHA256(req),
			CreatedAt: 1700, DeadlineAt: 2400, Policy: "Cancel"})
		if extra != "" {
			p = append(p[:len(p)-1], []byte(","+extra+"}")...)
		}
		return p
	}
	// subject is the card URL of request bytes req and the subject served for it, after edit.
	subject := func(req []byte, edit func(s *surfacepage.DecisionSubject)) cardCase {
		s := cardSubject(req, packetFor(req, ""), 1710, 2400)
		s.Passkeys = cardPasskeys(r.origin, 92)
		ref := s.Request
		if edit != nil {
			edit(s)
		}
		r.subjects[ref] = s
		return cardCase{r.origin + "/surface/decide/" + ref, s}
	}
	const unshown = "This request has a field the page cannot show ("
	for _, tc := range []struct {
		name, state, what string
		c                 cardCase
	}{
		{"control", "ready", "Decide request ", subject([]byte(`{"effect":"Human.Approve","scope":"worktree:1","cost":1,"requester":"a","now":1700}`), nil)},
		{"request bytes swapped", "refused", "The daemon's objects do not match their hashes; nothing to decide.",
			subject([]byte(`{"effect":"Human.Approve","scope":"worktree:2","cost":1,"requester":"a","now":1700}`), func(s *surfacepage.DecisionSubject) {
				s.RequestBytes = base64.RawURLEncoding.EncodeToString(other)
			})},
		{"packet hash", "refused", "The daemon's objects do not match their hashes; nothing to decide.",
			subject([]byte(`{"effect":"Human.Approve","scope":"worktree:3","cost":1,"requester":"a","now":1700}`), func(s *surfacepage.DecisionSubject) {
				s.Packet = hashref.SumSHA256([]byte("another packet")).String()
			})},
		{"another request's subject", "refused", "The daemon's objects do not match their hashes; nothing to decide.",
			subject([]byte(`{"effect":"Human.Approve","scope":"worktree:4","cost":1,"requester":"a","now":1700}`), func(s *surfacepage.DecisionSubject) {
				s.RequestBytes = base64.RawURLEncoding.EncodeToString(other) // its bytes and ref agree, but it is not this card's
				s.Request = hashref.SumSHA256(other).String()
			})},
		{"packet for another request", "refused", "The packet does not carry this request; nothing to decide.",
			subject([]byte(`{"effect":"Human.Approve","scope":"worktree:5","cost":1,"requester":"a","now":1700}`), func(s *surfacepage.DecisionSubject) {
				p := packetFor(other, "")
				s.Packet, s.PacketBytes = hashref.SumSHA256(p).String(), base64.RawURLEncoding.EncodeToString(p)
			})},
		{"unlisted request key", "refused", unshown + "request.extra)",
			subject([]byte(`{"effect":"Human.Approve","scope":"worktree:6","cost":1,"requester":"a","now":1700,"extra":"x"}`), nil)},
		{"unlisted packet key", "refused", unshown + "packet.extra)",
			subject([]byte(`{"effect":"Human.Approve","scope":"worktree:7","cost":1,"requester":"a","now":1700}`), func(s *surfacepage.DecisionSubject) {
				p := packetFor(mustB64(s.RequestBytes), `"extra":1`)
				s.Packet, s.PacketBytes = hashref.SumSHA256(p).String(), base64.RawURLEncoding.EncodeToString(p)
			})},
		{"missing key", "refused", unshown + "request.cost)",
			subject([]byte(`{"effect":"Human.Approve","scope":"worktree:8","requester":"a","now":1700}`), nil)},
		{"string for int", "refused", unshown + "request.cost)",
			subject([]byte(`{"effect":"Human.Approve","scope":"worktree:9","cost":"1","requester":"a","now":1700}`), nil)},
		{"fraction", "refused", unshown + "request.now)",
			subject([]byte(`{"effect":"Human.Approve","scope":"worktree:10","cost":1,"requester":"a","now":1700.5}`), nil)},
		{"int for string", "refused", unshown + "request.scope)",
			subject([]byte(`{"effect":"Human.Approve","scope":11,"cost":1,"requester":"a","now":1700}`), nil)},
		{"not an object", "refused", unshown + "request)", subject([]byte(`["Human.Approve"]`), nil)},
		{"not base64url", "refused", "The daemon's objects are not base64url; nothing to decide.",
			subject([]byte(`{"effect":"Human.Approve","scope":"worktree:12","cost":1,"requester":"a","now":1700}`), func(s *surfacepage.DecisionSubject) { s.PacketBytes = "*" + s.PacketBytes })},
		{"window", "refused", "The daemon's decision window is malformed; nothing to decide.",
			subject([]byte(`{"effect":"Human.Approve","scope":"worktree:13","cost":1,"requester":"a","now":1700}`), func(s *surfacepage.DecisionSubject) { s.Window.Max = -1 })},
		{"reserved", "reserved", "Decide request ",
			subject([]byte(`{"effect":"Human.Approve","scope":"registry:x","cost":1,"requester":"a","now":1700}`), func(s *surfacepage.DecisionSubject) {
				s.Reserved, s.Refusal = "publish", "publish approvals are minted by `world-publish approve` at a terminal (R4); this request cannot be decided here"
			})},
		{"refused", "refused", "Decide request ",
			subject([]byte(`{"effect":"Human.Approve","scope":"worktree:14","cost":1,"requester":"a","now":1700}`), func(s *surfacepage.DecisionSubject) {
				s.Refusal = "the packet's deadline passed at 1970-01-01T00:40:00Z; its policy is Cancel"
			})},
		{"decided", "decided", "Decide request ",
			subject([]byte(`{"effect":"Human.Approve","scope":"worktree:15","cost":1,"requester":"a","now":1700}`), func(s *surfacepage.DecisionSubject) {
				s.Decision, s.DecidedBy, s.Refusal = "deny", "ssh:mark", "request was already decided deny by ssh:mark"
			})},
	} {
		var view cardView
		b.do(t, &view, "openCard", map[string]any{"url": tc.c.url})
		var controls string
		b.do(t, &controls, "waitFor", map[string]any{"expression": "['go','deny','ssh','save','save-deny'].map((id) => document.getElementById(id).hidden ? '-' : id).join(' ') || 'x'"})
		if view.State != tc.state || !strings.HasPrefix(view.What, tc.what) || len(view.Gets) != 0 || len(view.Violations) != 0 {
			t.Fatalf("%s: state %q, what %q, %d gets; want %q, %q", tc.name, view.State, view.What, len(view.Gets), tc.state, tc.what)
		}
		switch tc.state {
		case "ready":
			if controls != "go deny ssh save save-deny" || len(view.Fields) != 13 {
				t.Fatalf("%s: controls %q, %d fields", tc.name, controls, len(view.Fields))
			}
		case "refused", "reserved", "decided":
			if controls != "- - - - -" {
				t.Fatalf("%s: a %s card offers %q", tc.name, tc.state, controls)
			}
			wantFields := 0
			if strings.HasPrefix(tc.what, "Decide request ") {
				wantFields = 13
			}
			if len(view.Fields) != wantFields || (wantFields != 0 && !strings.Contains(view.Status, map[string]string{"reserved": "A publish approval: use world-publish approve at a terminal.",
				"refused": "This request cannot be decided here: the packet's deadline passed", "decided": "Already decided: deny by ssh:mark."}[tc.state])) {
				t.Fatalf("%s: %d fields, status %q", tc.name, len(view.Fields), view.Status)
			}
		}
	}
	if len(r.posts) != 0 {
		t.Fatalf("%d POSTs from cards that offer nothing", len(r.posts))
	}
}

// cardCase is one card's URL and the subject the rig serves for it.
type cardCase struct {
	url string
	s   *surfacepage.DecisionSubject
}
