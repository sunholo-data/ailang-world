package main

// Row 201 A4 sprint 4 (design §2.8, §6 XD1, XD2, XD4): the decision card's browser rig, on the A2/A3
// driver (testdata/passkey_cdp.mjs, unchanged): Chrome headless, a temporary profile and a CDP
// virtual authenticator holding a test key (D-WORLD-109). The page and its script are served by the
// production daemon handler on the loopback; GET /surface/decisions/{ref} and POST
// /surface/decide/{ref} are this rig's, over the same temporary credentials store: the subject is
// served as given (tests tamper with it), and a POST is held to A3a's browser guards, parsed by
// the POST's own grammar (surfacepage.ParseDecideBody), put through the decision law
// (broker.DecisionRefusal) and judged by the daemon's own judges (authority.JudgeDecisionAssertion,
// JudgeDecisionSSH) at the subject's now. The daemon's POST arms over a world store are
// TestDecidePOST*'s (host/daemon); this rig proves what the page sends them.

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"net/http"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/coordinator"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

// cardPost is one POST /surface/decide/{ref} the rig received, and its answer.
type cardPost struct {
	ref    string
	body   surfacepage.DecideBody
	status int
	header http.Header
}

// decideCardRig serves the decision card for the subjects it holds.
type decideCardRig struct {
	t        *testing.T
	creds    string
	origin   string
	subjects map[string]*surfacepage.DecisionSubject
	scopes   map[string]string
	gets     map[string]int
	posts    []cardPost
}

// cardSubject is a subject for request bytes req and packet bytes pkt at now, its window's max
// min(now + 600, deadline) as the daemon's.
func cardSubject(req, pkt []byte, now, deadline int64) *surfacepage.DecisionSubject {
	return &surfacepage.DecisionSubject{V: 1, Request: hashref.SumSHA256(req).String(), RequestBytes: base64.RawURLEncoding.EncodeToString(req),
		Packet: hashref.SumSHA256(pkt).String(), PacketBytes: base64.RawURLEncoding.EncodeToString(pkt), Entry: 7,
		Proposal: hashref.SumSHA256([]byte("plan")).String(), Now: now, Window: surfacepage.DecisionWindow{Min: now + 1, Max: broker.DecisionWindowMax(now, deadline)},
		Passkeys: []surfacepage.CardPasskeys{}, SSH: surfacepage.DecisionSSH{Namespace: approval.DecisionNamespace}}
}

// cardObjects are a Human.Approve request of scope by requester and a packet carrying it with deadline.
func cardObjects(scope, requester string, deadline int64) (req, pkt []byte) {
	req = []byte(fmt.Sprintf(`{"effect":"Human.Approve","scope":%q,"cost":1,"requester":%q,"now":1700}`, scope, requester))
	pkt = coordinator.EncodeDecisionPacket(coordinator.DecisionPacket{ProposalHash: hashref.SumSHA256([]byte("plan")), RequestRef: hashref.SumSHA256(req),
		CreatedAt: 1700, DeadlineAt: deadline, Policy: "Cancel"})
	return req, pkt
}

func newDecideCardRig(t *testing.T) *decideCardRig {
	t.Helper()
	r := &decideCardRig{t: t, creds: filepath.Join(t.TempDir(), "credentials.db"), subjects: map[string]*surfacepage.DecisionSubject{},
		scopes: map[string]string{}, gets: map[string]int{}}
	c, err := store.OpenCredentials(r.creds) // make the store before the daemon opens it
	if err != nil {
		t.Fatal(err)
	}
	_ = c.Close()
	r.origin = cardOrigin(t, r.creds, "", func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			switch {
			case req.Method == "GET" && strings.HasPrefix(req.URL.Path, "/surface/decisions/"):
				r.serveSubject(w, strings.TrimPrefix(req.URL.Path, "/surface/decisions/"))
			case req.Method == "POST" && strings.HasPrefix(req.URL.Path, "/surface/decide/"):
				r.decide(w, req, strings.TrimPrefix(req.URL.Path, "/surface/decide/"))
			default:
				h.ServeHTTP(w, req)
			}
		})
	})
	return r
}

// add holds subject s (for request scope) under its request ref and returns the card's URL.
func (r *decideCardRig) add(s *surfacepage.DecisionSubject, scope string) string {
	r.subjects[s.Request], r.scopes[s.Request] = s, scope
	return r.origin + "/surface/decide/" + s.Request
}

func (r *decideCardRig) serveSubject(w http.ResponseWriter, ref string) {
	r.gets[ref]++
	s, ok := r.subjects[ref]
	w.Header().Set("Content-Type", "application/json")
	if !ok {
		w.WriteHeader(404)
		_, _ = w.Write([]byte(`{"v":1,"error":"NotFound","reason":"no such approval request"}`))
		return
	}
	_ = json.NewEncoder(w).Encode(s)
}

func (r *decideCardRig) refuse(w http.ResponseWriter, p cardPost, status int, class, reason string) {
	p.status = status
	r.posts = append(r.posts, p)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(surfacepage.ApproveRefusal{V: 1, Error: class, Reason: reason})
}

func (r *decideCardRig) decide(w http.ResponseWriter, req *http.Request, ref string) {
	p := cardPost{ref: ref, header: req.Header.Clone()}
	mt, _, _ := mime.ParseMediaType(req.Header.Get("Content-Type"))
	if req.Header.Get("Origin") != r.origin || req.Header.Get("Sec-Fetch-Site") != "same-origin" || mt != "application/json" ||
		req.Header.Get("Cookie") != "" || req.Header.Get("Authorization") != "" {
		r.refuse(w, p, 403, "CrossOrigin", "the rig's guards: Origin, Sec-Fetch-Site, Content-Type, no cookie")
		return
	}
	raw, _ := io.ReadAll(io.LimitReader(req.Body, surfacepage.MaxApproveBodyBytes))
	body, packet, request, err := surfacepage.ParseDecideBody(raw)
	p.body = body
	if err != nil {
		r.refuse(w, p, 400, "BadRequest", err.Error())
		return
	}
	s, ok := r.subjects[ref]
	if !ok {
		r.refuse(w, p, 404, "NotFound", "no such approval request "+ref)
		return
	}
	pkt, err := coordinator.DecodeDecisionPacket(mustB64(s.PacketBytes))
	if err != nil {
		r.refuse(w, p, 409, "Conflict", err.Error())
		return
	}
	f := broker.DecisionFacts{Ref: ref, RefOK: true, RequestOK: true, OnChain: true, Discovery: broker.Discovery{Complete: true, Carrying: 1},
		Scope: r.scopes[ref], Decided: s.Decision != "", DecidedVerb: s.Decision, DecidedBy: s.DecidedBy, Now: s.Now, Deadline: pkt.DeadlineAt,
		Packet: refDigest(s.Packet), Request: refDigest(ref), PostedPacket: packet, PostedRequest: request, Verb: body.Verb, Expiry: body.Expiry}
	if law := broker.DecisionRefusal(f); law.Step != "" {
		r.refuse(w, p, law.Status, "Conflict", law.Reason)
		return
	}
	verb, _ := approval.DecisionVerb(body.Verb) // D9 held
	ro, err := store.OpenCredentialsReadOnly(r.creds)
	if err != nil {
		r.t.Error(err)
		return
	}
	defer func() { _ = ro.Close() }()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	in := authority.DecisionJudgeInput{Verb: verb, Packet: packet, Request: request, Expiry: body.Expiry, Now: s.Now}
	var signed *approval.SignedDecision
	if body.WebAuthn != nil {
		signed, err = authority.JudgeDecisionAssertion(ctx, ro, in, *body.WebAuthn)
	} else if body.SSHSIG != nil {
		signed, err = authority.JudgeDecisionSSH(ctx, ro, in, []byte(*body.SSHSIG))
	} else {
		r.refuse(w, p, 403, "Forbidden", "the card never posts typed")
		return
	}
	if err != nil {
		r.refuse(w, p, 422, "NotApproved", "the decision does not decide request "+ref+": "+err.Error())
		return
	}
	s.Decision, s.DecidedBy = body.Verb, signed.DecidedBy()
	s.Refusal = broker.SubjectRefusal(broker.DecisionFacts{Ref: ref, RefOK: true, RequestOK: true, OnChain: true, Discovery: broker.Discovery{Complete: true, Carrying: 1},
		Decided: true, DecidedVerb: s.Decision, DecidedBy: s.DecidedBy, Now: s.Now, Deadline: pkt.DeadlineAt}).Reason
	p.status = 200
	r.posts = append(r.posts, p)
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(surfacepage.DecisionMade{V: 1, Decided: body.Verb, Request: ref, Decision: "sha256:" + strings.Repeat("d", 64), DecidedBy: s.DecidedBy})
}

func mustB64(s string) []byte {
	b, _ := base64.RawURLEncoding.DecodeString(s)
	return b
}

// decisionHash64 is the base64url SHA-256 of the decision message for (verb, s, expiry), as a get's challenge.
func decisionHash64(t *testing.T, verb byte, s *surfacepage.DecisionSubject, expiry int64) string {
	t.Helper()
	h, err := approval.DecisionMessageHash(verb, refDigest(s.Packet), refDigest(s.Request), expiry)
	if err != nil {
		t.Fatal(err)
	}
	return base64.RawURLEncoding.EncodeToString(h[:])
}
