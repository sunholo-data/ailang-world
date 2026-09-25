package projection

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang/serveapi/protocol"

	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/transitionreg"
)

// ---------------------------------------------------------------------------
// Fixtures: the REAL stack wherever possible (in-memory store, real resolver,
// real minted credentials, real StoreReader) — fakes only for call counting,
// failure injection and bounded-wait blocking (AC-DENIAL-JSONRPC / B3 races /
// Decision 6).
// ---------------------------------------------------------------------------

func openStore(t *testing.T) *store.Store {
	t.Helper()
	st, err := store.Open(":memory:")
	if err != nil {
		t.Fatalf("open in-memory store: %v", err)
	}
	t.Cleanup(func() { _ = st.Close() })
	return st
}

// liveGrant returns a capability that is live for the whole test window.
func liveGrant(effect string) broker.Capability {
	return broker.Capability{Effect: effect, Scope: "world", ExpiresAt: time.Now().Unix() + 7200, Budget: 10}
}

func mintToken(t *testing.T, st *store.Store, episode string, grants []broker.Capability) string {
	t.Helper()
	tok, _, _, err := authority.Mint(context.Background(), st, episode, grants, 3600, time.Now().Unix(), io.Discard)
	if err != nil {
		t.Fatalf("mint: %v", err)
	}
	return tok
}

func mintExpiredToken(t *testing.T, st *store.Store) string {
	t.Helper()
	tok, _, _, err := authority.Mint(context.Background(), st, "ep-expired",
		[]broker.Capability{{Effect: "alpha", Scope: "world", ExpiresAt: time.Now().Unix() + 7200, Budget: 10}},
		60, time.Now().Unix()-3600, io.Discard)
	if err != nil {
		t.Fatalf("mint expired: %v", err)
	}
	return tok
}

// descriptor builds a transitionreg.Descriptor whose access requirement names
// exactly one effect (mirrors transitionreg_test's descriptorWithAccess).
func descriptor(id, effect string) transitionreg.Descriptor {
	return transitionreg.Descriptor{
		ID: id, TransitionFn: hashref.SumSHA256([]byte("fn-" + id)), Interpreter: hashref.SumSHA256([]byte("interp")),
		SemanticsEpoch: 1, InputSchema: []byte(`{}`), OutputSchema: []byte(`{}`),
		Access:          transitionreg.EffectRequirement{Effect: effect, Scope: "world", Cost: 1},
		DeclaredEffects: []transitionreg.EffectRequirement{{Effect: effect, Scope: "world", Cost: 1}},
		Title:           "title-" + id, Description: "description-" + id,
	}
}

// seedRegistry stores revision 1 with the given descriptors and points the
// transition-registry head at it (transitionreg_test's seedRevision shape).
// Entries are canonically ordered by ID first — EncodeRevision requires the
// strict bytewise order that Request.Allowed then preserves (F9).
func seedRegistry(t *testing.T, st *store.Store, descs ...transitionreg.Descriptor) hashref.HashRef {
	t.Helper()
	ordered := append([]transitionreg.Descriptor(nil), descs...)
	sort.Slice(ordered, func(i, j int) bool { return ordered[i].ID < ordered[j].ID })
	return publishRevision(t, st, transitionreg.Revision{
		SemanticID: transitionreg.SemanticIDV1, InterfaceHash: transitionreg.InterfaceHashV1,
		Revision: 1, Entries: ordered,
	}, hashref.HashRef{})
}

func publishRevision(t *testing.T, st *store.Store, rev transitionreg.Revision, expected hashref.HashRef) hashref.HashRef {
	t.Helper()
	payload, err := transitionreg.EncodeRevision(rev)
	if err != nil {
		t.Fatalf("encode revision: %v", err)
	}
	obj := store.Object{
		Hash: hashref.SumSHA256(payload), InterfaceHash: transitionreg.InterfaceHashV1,
		SemanticID: transitionreg.SemanticIDV1, Provenance: "projection-test", Payload: payload,
	}
	if err := st.PutObject(obj); err != nil {
		t.Fatalf("put object: %v", err)
	}
	if err := st.CompareAndSetRegistryHead(store.TransitionRegistryV1, expected, obj.Hash); err != nil {
		t.Fatalf("cas head: %v", err)
	}
	return obj.Hash
}

// --- reference denial/error writers -----------------------------------------
// The production writers are the daemon's own (injected at mount; the
// byte-identity proof lives in host/daemon's tests). These test references
// mirror the daemon's writeAPIError envelope shape so the package-level
// denial tests assert status + class + constancy here.
func refDeny(w http.ResponseWriter, k authority.DenialKind) {
	status := http.StatusUnauthorized
	if k == authority.DenialMalformed {
		status = http.StatusBadRequest
	}
	writeRefAPIError(w, k.String(), "ref-denial-"+k.String(), status)
}

func refFail(w http.ResponseWriter, class, message string, status int) {
	writeRefAPIError(w, class, message, status)
}

func writeRefAPIError(w http.ResponseWriter, class, message string, status int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{"error": map[string]any{"class": class, "message": message}})
}

type apiErrorBody struct {
	Error struct {
		Class   string `json:"class"`
		Message string `json:"message"`
	} `json:"error"`
}

// --- counting / blocking seams ------------------------------------------------

type countingHeads struct {
	inner HeadReader
	calls int
}

func (c *countingHeads) GetRegistryHead(ctx context.Context, name string) (hashref.HashRef, bool, error) {
	c.calls++
	return c.inner.GetRegistryHead(ctx, name)
}

type countingReader struct {
	inner transitionreg.Reader
	calls int
}

func (c *countingReader) ReadSnapshot(ctx context.Context) (transitionreg.Snapshot, error) {
	c.calls++
	return c.inner.ReadSnapshot(ctx)
}

type errHeads struct{ err error }

func (e errHeads) GetRegistryHead(context.Context, string) (hashref.HashRef, bool, error) {
	return hashref.HashRef{}, false, e.err
}

// fixedHeads reports a fixed head state (for the B3 race fixtures) and never
// reads any store.
type fixedHeads struct{ ok bool }

func (f fixedHeads) GetRegistryHead(context.Context, string) (hashref.HashRef, bool, error) {
	return hashref.HashRef{}, f.ok, nil
}

type errReader struct{ err error }

func (e errReader) ReadSnapshot(context.Context) (transitionreg.Snapshot, error) {
	return transitionreg.Snapshot{}, e.err
}

// blockingResolver blocks in ResolveContext until the request ctx is done
// (returning ctx.Err()) — the bounded-wait fault injection. Its 2 s escape
// hatch returns a successful binding so a mutation that drops the ctx
// (MUT-DROP-DEADLINE-PROJ) yields a 200 response instead of a hang, and the
// bounded-wait assertions RED on status AND elapsed time.
type blockingResolver struct{}

func (blockingResolver) Resolve(string, int64) authority.ResolveOutcome {
	denied := authority.DenialUnknown
	return authority.ResolveOutcome{Denied: &denied}
}

func (blockingResolver) ResolveContext(ctx context.Context, _ string, _ int64) (authority.ResolveOutcome, error) {
	select {
	case <-ctx.Done():
		return authority.ResolveOutcome{}, ctx.Err()
	case <-time.After(2 * time.Second):
		return authority.ResolveOutcome{Success: &authority.SessionBinding{
			EpisodeID: "ep-blocked", Caps: []broker.Capability{liveGrant("alpha")},
			ExpiresAt: time.Now().Unix() + 3600, CreatedAt: time.Now().Unix(),
		}}, nil
	}
}

// default test config: REAL resolver, REAL reader/heads over the store,
// reference writers, a small finite bound (Decision 6).
func testConfig(st *store.Store) Config {
	return Config{
		Resolver: authority.New(st),
		Reader:   transitionreg.NewReader(st),
		Heads:    st,
		Deny:     refDeny,
		Fail:     refFail,
		Agent:    protocol.AgentInfo{Name: "ailang-worldd", Description: "test projection agent", Version: "0.1.0"},
		MaxWait:  2 * time.Second,
	}
}

func mustHandler(t *testing.T, cfg Config) *Handler {
	t.Helper()
	h, err := New(cfg)
	if err != nil {
		t.Fatalf("projection.New: %v", err)
	}
	return h
}

func getCard(t *testing.T, h *Handler, auth string) *httptest.ResponseRecorder {
	t.Helper()
	req := httptest.NewRequest(http.MethodGet, "/.well-known/agent.json", nil)
	if auth != "" {
		req.Header.Set("Authorization", auth)
	}
	rec := httptest.NewRecorder()
	h.AgentCard(rec, req)
	return rec
}

func postA2A(t *testing.T, h *Handler, auth, body string) *httptest.ResponseRecorder {
	t.Helper()
	var rdr io.Reader
	if body != "" {
		rdr = strings.NewReader(body)
	}
	req := httptest.NewRequest(http.MethodPost, "/a2a/", rdr)
	if auth != "" {
		req.Header.Set("Authorization", auth)
	}
	rec := httptest.NewRecorder()
	h.A2A(rec, req)
	return rec
}

func skillIDs(t *testing.T, body []byte) []string {
	t.Helper()
	var card map[string]any
	if err := json.Unmarshal(body, &card); err != nil {
		t.Fatalf("card body is not JSON: %v (%s)", err, body)
	}
	skills, ok := card["skills"].([]any)
	if !ok {
		t.Fatalf("card.skills is %T, want an array (body %s)", card["skills"], body)
	}
	ids := make([]string, 0, len(skills))
	for _, s := range skills {
		ids = append(ids, s.(map[string]any)["id"].(string))
	}
	return ids
}

// a2aErr decodes a JSON-RPC error response body and asserts the wire shape:
// jsonrpc "2.0", an error OBJECT with EXACTLY {code, message} (no REST
// "class" key — the /a2a/ route never emits the daemon envelope, AC8), and an
// id member. It returns the code and message.
func a2aErr(t *testing.T, body []byte) (int, string, json.RawMessage) {
	t.Helper()
	var resp struct {
		JSONRPC string          `json:"jsonrpc"`
		ID      json.RawMessage `json:"id"`
		Result  json.RawMessage `json:"result"`
		Error   map[string]any  `json:"error"`
	}
	if err := json.Unmarshal(body, &resp); err != nil {
		t.Fatalf("/a2a/ body is not JSON: %v (%s)", err, body)
	}
	if resp.JSONRPC != "2.0" {
		t.Fatalf("jsonrpc = %q, want \"2.0\" (body %s)", resp.JSONRPC, body)
	}
	if resp.Error == nil {
		t.Fatalf("/a2a/ body has no error object (never a success on this surface): %s", body)
	}
	if len(resp.Result) != 0 {
		t.Fatalf("/a2a/ body carries a result member — no success path exists: %s", body)
	}
	for k := range resp.Error {
		if k != "code" && k != "message" {
			t.Fatalf("/a2a/ error object carries unexpected key %q (the REST envelope's 'class' must never appear): %s", k, body)
		}
	}
	codeF, _ := resp.Error["code"].(float64)
	msg, _ := resp.Error["message"].(string)
	return int(codeF), msg, resp.ID
}

func tasksSendBody(skillID string) string {
	return `{"jsonrpc":"2.0","id":7,"method":"tasks/send","params":{"id":"task-1",` +
		`"message":{"role":"user","parts":[{"type":"text","text":"hi"}]},` +
		`"metadata":{"skill_id":` + fmt.Sprintf("%q", skillID) + `}}}`
}

// compile-time witness for F9's CapabilitySource shape: the projection passes
// bindingCaps (over the authority binding's immutable grants) into NewRequest.
var _ transitionreg.CapabilitySource = bindingCaps{}

// ---------------------------------------------------------------------------
// AC2 / AC3 / AC4 / AC8-card — the card body
// ---------------------------------------------------------------------------

// TestAgentCard_ExactSkillSetPerSession is AC2 (with the F6 verbatim-ID proof
// that kills MUT-CARDSURFACE and the exactness that kills MUT-SESSION-UNION):
// two sessions with unequal capability sets see EXACTLY their own authorized
// transition-ID sets — verbatim, in registry bytewise order, with no extras.
// The registry IDs deliberately use REAL World grammar (dots and slashes) that
// protocol.CallerSurface's MCP name regex rejects (F6), so any routing of the
// card through CallerSurface fails here.
func TestAgentCard_ExactSkillSetPerSession(t *testing.T) {
	st := openStore(t)
	dAlpha := descriptor("tools.echo", "alpha")                 // dot grammar
	dBeta := descriptor("world/recovery-transition/v1", "beta") // slash grammar
	dGamma := descriptor("tools.hidden", "gamma")               // nobody's capability
	seedRegistry(t, st, dAlpha, dBeta, dGamma)

	h := mustHandler(t, testConfig(st))
	tokA := mintToken(t, st, "ep-a", []broker.Capability{liveGrant("alpha")})
	tokB := mintToken(t, st, "ep-b", []broker.Capability{liveGrant("beta")})
	tokBoth := mintToken(t, st, "ep-both", []broker.Capability{liveGrant("alpha"), liveGrant("beta")})

	recA := getCard(t, h, "Bearer "+tokA)
	if recA.Code != http.StatusOK {
		t.Fatalf("session A card status=%d body=%s", recA.Code, recA.Body)
	}
	if got, want := skillIDs(t, recA.Body.Bytes()), []string{"tools.echo"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("session A skills = %v, want exactly %v", got, want)
	}
	recB := getCard(t, h, "Bearer "+tokB)
	if got, want := skillIDs(t, recB.Body.Bytes()), []string{"world/recovery-transition/v1"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("session B skills = %v, want exactly %v", got, want)
	}
	// Registry bytewise order is the card order (tools.* sorts before world/*).
	recBoth := getCard(t, h, "Bearer "+tokBoth)
	if got, want := skillIDs(t, recBoth.Body.Bytes()), []string{"tools.echo", "world/recovery-transition/v1"}; !reflect.DeepEqual(got, want) {
		t.Fatalf("both-cap session skills = %v, want registry bytewise order %v", got, want)
	}
	// Verbatim: the emitted IDs are the registry IDs byte-for-byte, separators
	// and all — never mangled toward the MCP name grammar.
	for _, id := range skillIDs(t, recBoth.Body.Bytes()) {
		if id != "tools.echo" && id != "world/recovery-transition/v1" {
			t.Fatalf("skill id %q is not a verbatim registry ID", id)
		}
	}
}

// TestAgentCard_TracksSessionSnapshot is AC3: from ONE handler (the same
// daemon-level surface), two live sessions with different capability snapshots
// observe DIFFERENT cards, and re-fetching a session's card does not change it
// without a capability or registry change (stable per snapshot).
func TestAgentCard_TracksSessionSnapshot(t *testing.T) {
	st := openStore(t)
	seedRegistry(t, st, descriptor("tools.echo", "alpha"), descriptor("tools.probe", "beta"))
	h := mustHandler(t, testConfig(st))

	tokA := mintToken(t, st, "ep-a", []broker.Capability{liveGrant("alpha")})
	tokB := mintToken(t, st, "ep-b", []broker.Capability{liveGrant("beta")})
	cardA := skillIDs(t, getCard(t, h, "Bearer "+tokA).Body.Bytes())
	cardB := skillIDs(t, getCard(t, h, "Bearer "+tokB).Body.Bytes())
	if reflect.DeepEqual(cardA, cardB) {
		t.Fatalf("sessions with different snapshots observed EQUAL cards %v", cardA)
	}
	cardAAgain := skillIDs(t, getCard(t, h, "Bearer "+tokA).Body.Bytes())
	if !reflect.DeepEqual(cardA, cardAAgain) {
		t.Fatalf("same session, same snapshot: card changed between requests: %v -> %v", cardA, cardAAgain)
	}
}

// TestAgentCard_AmbientExportsAbsent is AC4: ambient interpreter names are
// never projected. The banned set below names them; the control assertion
// (card non-empty for a capable session) keeps the scan non-vacuous. Note the
// grammar already refuses most ambient names at registration (validateID is
// lowercase-only) — the projection must additionally never INJECT them
// (MUT-UNFILTERED-PROJECTION adds a hardcoded "exit" skill: this test and the
// exact-set test both red).
func TestAgentCard_AmbientExportsAbsent(t *testing.T) {
	st := openStore(t)
	seedRegistry(t, st, descriptor("tools.echo", "alpha"))
	h := mustHandler(t, testConfig(st))
	tok := mintToken(t, st, "ep-wide", []broker.Capability{liveGrant("alpha"), liveGrant("beta"), liveGrant("gamma")})

	ids := skillIDs(t, getCard(t, h, "Bearer "+tok).Body.Bytes())
	if len(ids) == 0 {
		t.Fatal("control failed: a capable session's card is empty, so the banned-name scan below is vacuous")
	}
	banned := []string{"exit", "writeBytes", "readBytes", "writeLine", "readLine", "flush", "submit_feedback", "std/io.writeBytes", "std/io.readBytes"}
	for _, id := range ids {
		for _, b := range banned {
			if id == b {
				t.Fatalf("ambient export %q must never appear in the card (skills %v)", b, ids)
			}
		}
		if strings.HasPrefix(id, "std/io.") || strings.HasPrefix(id, "std/io/") {
			t.Fatalf("std/io ambient export %q must never appear in the card", id)
		}
	}
}

// TestAgentCard_UpstreamKeySet is the card half of AC8 + the AC1 shape gate:
// the card is a map[string]any literal with EXACTLY the upstream
// serveapi/a2a_handler key set (F5) and the exact skill-entry key set.
func TestAgentCard_UpstreamKeySet(t *testing.T) {
	st := openStore(t)
	seedRegistry(t, st, descriptor("tools.echo", "alpha"))
	h := mustHandler(t, testConfig(st))
	tok := mintToken(t, st, "ep-a", []broker.Capability{liveGrant("alpha")})

	rec := getCard(t, h, "Bearer "+tok)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body)
	}
	var card map[string]any
	if err := json.Unmarshal(rec.Body.Bytes(), &card); err != nil {
		t.Fatalf("card not JSON: %v", err)
	}
	want := []string{"capabilities", "defaultInputModes", "defaultOutputModes", "description", "name", "skills", "url", "version"}
	got := make([]string, 0, len(card))
	for k := range card {
		got = append(got, k)
	}
	sortStrings(got)
	sortStrings(want)
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("card key set = %v, want the upstream key set %v", got, want)
	}
	caps := card["capabilities"].(map[string]any)
	for _, k := range []string{"streaming", "pushNotifications", "stateTransitionHistory"} {
		if v, ok := caps[k].(bool); !ok || v {
			t.Fatalf("capabilities[%q] = %v, want false", k, caps[k])
		}
	}
	if len(caps) != 3 {
		t.Fatalf("capabilities key count = %d, want exactly 3", len(caps))
	}
	if !reflect.DeepEqual(card["defaultInputModes"], []any{"application/json"}) ||
		!reflect.DeepEqual(card["defaultOutputModes"], []any{"application/json"}) {
		t.Fatalf("input/output modes = %v/%v, want [application/json]", card["defaultInputModes"], card["defaultOutputModes"])
	}
	if card["name"] != "ailang-worldd" || card["version"] != "0.1.0" || card["description"] != "test projection agent" {
		t.Fatalf("agent identity = %v/%v/%v, want the injected protocol.AgentInfo", card["name"], card["version"], card["description"])
	}
	skills := card["skills"].([]any)
	if len(skills) != 1 {
		t.Fatalf("skills len = %d, want 1 (control for the entry-key assertion)", len(skills))
	}
	entry := skills[0].(map[string]any)
	wantEntry := []string{"description", "examples", "id", "name", "tags"}
	gotEntry := make([]string, 0, len(entry))
	for k := range entry {
		gotEntry = append(gotEntry, k)
	}
	sortStrings(gotEntry)
	if !reflect.DeepEqual(gotEntry, wantEntry) {
		t.Fatalf("skill entry key set = %v, want %v", gotEntry, wantEntry)
	}
	if entry["id"] != "tools.echo" || entry["name"] != "title-tools.echo" || entry["description"] != "description-tools.echo" {
		t.Fatalf("skill entry = %v — id must be the verbatim stable ID, name/description the descriptor's Title/Description", entry)
	}
}

func sortStrings(s []string) {
	for i := 1; i < len(s); i++ {
		for j := i; j > 0 && s[j] < s[j-1]; j-- {
			s[j], s[j-1] = s[j-1], s[j]
		}
	}
}

// ---------------------------------------------------------------------------
// AC-CARD-DENIALS / AC6 — card-route denial matrix
// ---------------------------------------------------------------------------

// TestAgentCard_DenialMatrix is the card-route half of AC-DENIAL-JSONRPC and
// all of AC-CARD-DENIALS: each of the four denial kinds yields its F2 HTTP
// status (401/401/401/400) and the daemon-envelope class
// (SessionAbsent/SessionUnknown/SessionExpired/InvalidSession), the body comes
// through the INJECTED writer, and denial messages are CONSTANT per kind (the
// same bad credential twice yields byte-identical bodies — nothing leaks).
// MUT-CARD-DENIAL-TABLE and MUT-CARD-ENVELOPE die here.
func TestAgentCard_DenialMatrix(t *testing.T) {
	st := openStore(t)
	h := mustHandler(t, testConfig(st))
	expiredTok := mintExpiredToken(t, st)
	cases := []struct {
		name       string
		auth       string
		wantStatus int
		wantClass  string
	}{
		{"absent", "", http.StatusUnauthorized, "SessionAbsent"},
		{"unknown", "Bearer " + strings.Repeat("b", 64), http.StatusUnauthorized, "SessionUnknown"},
		{"expired", "Bearer " + expiredTok, http.StatusUnauthorized, "SessionExpired"},
		{"malformed", "Token abc", http.StatusBadRequest, "InvalidSession"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := getCard(t, h, tc.auth)
			if rec.Code != tc.wantStatus {
				t.Fatalf("status=%d, want %d; body=%s", rec.Code, tc.wantStatus, rec.Body)
			}
			var body apiErrorBody
			if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
				t.Fatalf("card denial body is not the APIError envelope shape: %v (%s)", err, rec.Body)
			}
			if body.Error.Class != tc.wantClass {
				t.Fatalf("class=%q, want %q (body %s)", body.Error.Class, tc.wantClass, rec.Body)
			}
			if body.Error.Message == "" {
				t.Fatal("denial message is empty — every kind has a constant message")
			}
			// Constant per kind: repeat the exact request, expect identical bytes.
			rec2 := getCard(t, h, tc.auth)
			if !reflect.DeepEqual(rec.Body.Bytes(), rec2.Body.Bytes()) {
				t.Fatalf("denial bodies differ between identical requests:\n%s\n%s", rec.Body, rec2.Body)
			}
		})
	}
}

// ---------------------------------------------------------------------------
// AC-ABSENT-HEAD (B3) — absent head vs read failure vs the head race
// ---------------------------------------------------------------------------

// TestAgentCard_AbsentHeadZeroSkills is AC-ABSENT-HEAD's absent half: no
// transition-registry head at all is a LEGITIMATE state — the card is HTTP
// 200 with an EMPTY (non-null) skills array, still session-gated. This is the
// zero-skills-at-200 anchor that kills MUT-ABSENT-HEAD (absent treated as a
// read error) while TestAgentCard_HeadReadFailuresAre5xx stays green as the
// control.
func TestAgentCard_AbsentHeadZeroSkills(t *testing.T) {
	st := openStore(t) // never seeded — the absent-head state (F8)
	h := mustHandler(t, testConfig(st))
	tok := mintToken(t, st, "ep-a", []broker.Capability{liveGrant("alpha")})

	rec := getCard(t, h, "Bearer "+tok)
	if rec.Code != http.StatusOK {
		t.Fatalf("absent-head card status=%d, want 200 (absent head is not a failure); body=%s", rec.Code, rec.Body)
	}
	if ids := skillIDs(t, rec.Body.Bytes()); len(ids) != 0 {
		t.Fatalf("absent-head skills = %v, want EMPTY", ids)
	}
	if !strings.Contains(string(rec.Body.Bytes()), `"skills":[]`) {
		t.Fatalf("absent-head card body must carry an explicit empty skills array: %s", rec.Body)
	}
	// Still authenticated: the same request WITHOUT a session denies (401), so
	// the zero-skills 200 is a property of the SESSION's surface, not a public leak.
	if rec := getCard(t, h, ""); rec.Code != http.StatusUnauthorized {
		t.Fatalf("unauthenticated absent-head card status=%d, want 401 — the route never degrades", rec.Code)
	}
}

// TestAgentCard_HeadReadFailuresAre5xx is AC-ABSENT-HEAD's failure half: a
// head-check error is 5xx (503), and a snapshot failure AFTER the check saw a
// PRESENT head is 5xx (503) — a present-but-unreadable registry is a failure,
// never the absent state (B3's race rule, second half).
func TestAgentCard_HeadReadFailuresAre5xx(t *testing.T) {
	st := openStore(t)
	tok := mintToken(t, st, "ep-a", []broker.Capability{liveGrant("alpha")})

	// Head check itself fails (e.g. store contention surfaced as an error).
	cfgErrHeads := testConfig(st)
	cfgErrHeads.Heads = errHeads{err: errors.New("head check blew up")}
	h := mustHandler(t, cfgErrHeads)
	rec := getCard(t, h, "Bearer "+tok)
	if rec.Code != http.StatusServiceUnavailable {
		t.Fatalf("head-check error status=%d, want 503; body=%s", rec.Code, rec.Body)
	}
	var body apiErrorBody
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil || body.Error.Class != "ProjectionUnavailable" {
		t.Fatalf("head-check failure body = %s, want the daemon envelope with class ProjectionUnavailable", rec.Body)
	}

	// Check sees a head, then the snapshot read fails (present-but-unreadable).
	cfgRace := testConfig(st)
	cfgRace.Heads = fixedHeads{ok: true}
	cfgRace.Reader = errReader{err: errors.New("snapshot read blew up")}
	h2 := mustHandler(t, cfgRace)
	rec2 := getCard(t, h2, "Bearer "+tok)
	if rec2.Code != http.StatusServiceUnavailable {
		t.Fatalf("present-head snapshot failure status=%d, want 503; body=%s", rec2.Code, rec2.Body)
	}
}

// TestAgentCard_HeadRaceAbsentThenSucceeds is B3's first race half: the head
// check said ABSENT, but by snapshot time the head has been published — the
// card must USE NewRequest's result (the head arrived; succeed), not the
// stale absent reading.
func TestAgentCard_HeadRaceAbsentThenSucceeds(t *testing.T) {
	st := openStore(t)
	// The real reader CAN read it (the head has raced in)…
	seedRegistry(t, st, descriptor("tools.echo", "alpha"))
	tok := mintToken(t, st, "ep-a", []broker.Capability{liveGrant("alpha")})
	cfg := testConfig(st)
	// …but the check observes the PRE-publish world.
	cfg.Heads = fixedHeads{ok: false}
	h := mustHandler(t, cfg)

	rec := getCard(t, h, "Bearer "+tok)
	if rec.Code != http.StatusOK {
		t.Fatalf("raced-head card status=%d, want 200; body=%s", rec.Code, rec.Body)
	}
	if ids := skillIDs(t, rec.Body.Bytes()); !reflect.DeepEqual(ids, []string{"tools.echo"}) {
		t.Fatalf("head published between check and snapshot must be USED (NewRequest's result): skills = %v, want [tools.echo]", ids)
	}
}

// ---------------------------------------------------------------------------
// AC-A2A-CODES / AC-DENIAL-JSONRPC /a2a/ half — the JSON-RPC surface
// ---------------------------------------------------------------------------

// TestA2A_DenialMatrix is the /a2a/ half of AC-DENIAL-JSONRPC (the r2
// carve-out): all four denial kinds, each a protocol.A2AError on an HTTP 200
// carrier — -32001 for absent/unknown/expired, -32600 for malformed — with
// constant messages and never the REST envelope (no "class" key). It kills
// MUT-DENIAL-401 (status assertion) and MUT-CARD-ENVELOPE's /a2a/ twin.
func TestA2A_DenialMatrix(t *testing.T) {
	st := openStore(t)
	h := mustHandler(t, testConfig(st))
	expiredTok := mintExpiredToken(t, st)
	cases := []struct {
		name     string
		auth     string
		wantCode int
	}{
		{"absent", "", -32001},
		{"unknown", "Bearer " + strings.Repeat("b", 64), -32001},
		{"expired", "Bearer " + expiredTok, -32001},
		{"malformed", "Token abc", -32600},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := postA2A(t, h, tc.auth, tasksSendBody("tools.echo"))
			if rec.Code != http.StatusOK {
				t.Fatalf("/a2a/ denial status=%d, want HTTP 200 (F5b: A2AError always writes 200); body=%s", rec.Code, rec.Body)
			}
			code, msg, id := a2aErr(t, rec.Body.Bytes())
			if code != tc.wantCode {
				t.Fatalf("code=%d, want %d (denial matrix)", code, tc.wantCode)
			}
			if msg == "" {
				t.Fatal("denial message is empty — every kind has a constant message")
			}
			if string(id) != "null" && len(id) != 0 {
				t.Fatalf("denial error id = %s, want null/omitted (the request ID is not trusted before parse)", id)
			}
			rec2 := postA2A(t, h, tc.auth, tasksSendBody("tools.echo"))
			if !reflect.DeepEqual(rec.Body.Bytes(), rec2.Body.Bytes()) {
				t.Fatalf("denial bodies differ between identical requests (messages must be constant):\n%s\n%s", rec.Body, rec2.Body)
			}
			// Messages must not echo request content (the guessed skill id).
			if strings.Contains(msg, "tools.echo") {
				t.Fatalf("denial message interpolates request content: %q", msg)
			}
		})
	}
}

// TestA2A_CodeMatrix is AC-A2A-CODES + AC5 + the AC8 code restriction: every
// /a2a/ outcome is a protocol.A2AError body at HTTP 200 with EXACTLY one of
// the five sanctioned codes; an authorized skill_id receives the CONSTANT
// -32603 not-available message (never a success — MUT-A2A-FAKE-SUCCESS,
// MUT-A2A-MESSAGE-INTERP die here); unauthorized/guessed/stale skill_ids
// receive -32602 "not authorized" BEFORE any other processing
// (MUT-A2A-STALE); code mappings are exact (MUT-A2A-CODEFLIP).
func TestA2A_CodeMatrix(t *testing.T) {
	st := openStore(t)
	seedRegistry(t, st, descriptor("tools.echo", "alpha"))
	h := mustHandler(t, testConfig(st))
	tok := mintToken(t, st, "ep-a", []broker.Capability{liveGrant("alpha")})
	auth := "Bearer " + tok

	cases := []struct {
		name     string
		body     string
		wantCode int
		wantMsg  string // exact
		wantID   string // exact raw id member; "" = skip id assertion
	}{
		{"garbage JSON", "this is not json", -32600, "invalid JSON-RPC request", ""},
		{"wrong version", `{"jsonrpc":"1.9","id":3,"method":"tasks/send","params":{}}`, -32600, "invalid JSON-RPC request", "3"},
		{"missing version", `{"id":4,"method":"tasks/send","params":{}}`, -32600, "invalid JSON-RPC request", "4"},
		{"wrong method", `{"jsonrpc":"2.0","id":5,"method":"tasks/list","params":{}}`, -32601, "method not found", "5"},
		{"params not object", `{"jsonrpc":"2.0","id":6,"method":"tasks/send","params":"junk"}`, -32602, "invalid params", "6"},
		{"no skill_id", `{"jsonrpc":"2.0","id":8,"method":"tasks/send","params":{"metadata":{}}}`, -32602, "not authorized", "8"},
		{"guessed skill", tasksSendBody("tools.guessed"), -32602, "not authorized", "7"},
		{"stale skill", tasksSendBody("tools.retired"), -32602, "not authorized", "7"},
		{"ambient-name skill", tasksSendBody("writeBytes"), -32602, "not authorized", "7"},
		{"authorized skill", tasksSendBody("tools.echo"), -32603, "transition invocation is not available in this daemon", "7"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			rec := postA2A(t, h, auth, tc.body)
			if rec.Code != http.StatusOK {
				t.Fatalf("status=%d, want 200; body=%s", rec.Code, rec.Body)
			}
			code, msg, id := a2aErr(t, rec.Body.Bytes())
			if code != tc.wantCode {
				t.Fatalf("code=%d, want %d; body=%s", code, tc.wantCode, rec.Body)
			}
			if msg != tc.wantMsg {
				t.Fatalf("message=%q, want exactly %q (constant, never interpolated)", msg, tc.wantMsg)
			}
			if tc.wantID != "" && string(id) != tc.wantID {
				t.Fatalf("id=%s, want echoed id %s", id, tc.wantID)
			}
		})
	}

	// The not-available answer is byte-constant: two authorized calls produce
	// identical bodies (id aside, which is echoed).
	a := postA2A(t, h, auth, tasksSendBody("tools.echo"))
	b := postA2A(t, h, auth, tasksSendBody("tools.echo"))
	if !reflect.DeepEqual(a.Body.Bytes(), b.Body.Bytes()) {
		t.Fatalf("authorized-skill bodies differ:\n%s\n%s", a.Body, b.Body)
	}
}

// TestA2A_NeverWritesStore is the "produces no store/log change" half of AC5
// and the no-store-write clause of AC8: a full denial+admission battery leaves
// the world head AND the registry head byte-identical.
func TestA2A_NeverWritesStore(t *testing.T) {
	st := openStore(t)
	seeded := seedRegistry(t, st, descriptor("tools.echo", "alpha"))
	h := mustHandler(t, testConfig(st))
	tok := mintToken(t, st, "ep-a", []broker.Capability{liveGrant("alpha")})

	postA2A(t, h, "", tasksSendBody("tools.echo"))             // denied
	postA2A(t, h, "Bearer "+tok, "garbage")                    // parse error
	postA2A(t, h, "Bearer "+tok, tasksSendBody("tools.guess")) // unauthorized
	postA2A(t, h, "Bearer "+tok, tasksSendBody("tools.echo"))  // authorized -> constant refusal
	getCard(t, h, "Bearer "+tok)                               // card route too

	// The store's every committed mutation address is content-hashed, so the
	// transition-registry head (the only store mutation the project surface
	// could even conceptually touch) must be byte-identical after the battery.
	// The world log gains rows only through POST /v1/commit, which this suite
	// never calls.
	regHeadAfter, ok, err := st.GetRegistryHead(context.Background(), store.TransitionRegistryV1)
	if err != nil || !ok {
		t.Fatalf("GetRegistryHead after the battery: ok=%v err=%v — the probe must still read", ok, err)
	}
	if regHeadAfter != seeded {
		t.Fatalf("projection writes leaked: registry head %q -> %q across the battery (want unchanged)", seeded, regHeadAfter)
	}
	// Control: the probe distinguishes a write when one actually happens —
	// re-seeding moves the head, so an unchanged head above is meaningful.
	d2 := descriptor("tools.probe", "beta")
	next, err := transitionreg.BuildNext(transitionreg.Revision{
		SemanticID: transitionreg.SemanticIDV1, InterfaceHash: transitionreg.InterfaceHashV1,
		Revision: 1, Entries: []transitionreg.Descriptor{descriptor("tools.echo", "alpha")},
	}, []transitionreg.Change{{ID: d2.ID, Descriptor: &d2}})
	if err != nil {
		t.Fatalf("BuildNext control: %v", err)
	}
	head2, err := transitionreg.NewReader(st).Publish(context.Background(), seeded, next)
	if err != nil {
		t.Fatalf("Publish control: %v", err)
	}
	if head2 == seeded {
		t.Fatal("control failed: a real publish did not move the registry head — the no-write probe is vacuous")
	}
}

// ---------------------------------------------------------------------------
// AC-DENIAL-JSONRPC cross-route — a denied request never acquires a snapshot
// ---------------------------------------------------------------------------

// TestProjection_DeniedNeverTouchesRegistry is the zero-snapshot half of
// AC-DENIAL-JSONRPC (MUT-DENIAL-SNAPSHOT-FIRST's killer): for ALL FOUR denial
// kinds on BOTH routes, a counting GetRegistryHead observer and a counting
// snapshot reader must read ZERO calls — denial happens before any registry
// acquisition. A valid-session request that fails JSON-RPC parsing must also
// stay at zero (parse precedes admission).
func TestProjection_DeniedNeverTouchesRegistry(t *testing.T) {
	mk := func(t *testing.T) (*countingHeads, *countingReader, *Handler, string) {
		st := openStore(t)
		seedRegistry(t, st, descriptor("tools.echo", "alpha"))
		heads := &countingHeads{inner: st}
		reader := &countingReader{inner: transitionreg.NewReader(st)}
		cfg := testConfig(st)
		cfg.Heads, cfg.Reader = heads, reader
		h := mustHandler(t, cfg)
		tok := mintToken(t, st, "ep-a", []broker.Capability{liveGrant("alpha")})
		return heads, reader, h, tok
	}

	t.Run("card route, four denial kinds", func(t *testing.T) {
		heads, reader, h, tok := mk(t)
		_ = tok
		expiredTok := mintExpiredToken(t, openStore(t)) // shape-only: unknown on THIS store would also be fine; use a real expired below
		_ = expiredTok
		for _, auth := range []string{"", "Token abc", "Bearer " + strings.Repeat("b", 64)} {
			getCard(t, h, auth)
		}
		if heads.calls != 0 || reader.calls != 0 {
			t.Fatalf("denied card requests touched the registry: heads.calls=%d reader.calls=%d, want 0/0", heads.calls, reader.calls)
		}
	})

	t.Run("/a2a/ route, four denial kinds", func(t *testing.T) {
		heads, reader, h, _ := mk(t)
		for _, auth := range []string{"", "Token abc", "Bearer " + strings.Repeat("b", 64)} {
			postA2A(t, h, auth, tasksSendBody("tools.echo"))
		}
		if heads.calls != 0 || reader.calls != 0 {
			t.Fatalf("denied /a2a/ requests touched the registry: heads.calls=%d reader.calls=%d, want 0/0", heads.calls, reader.calls)
		}
	})

	t.Run("expired session, both routes", func(t *testing.T) {
		st := openStore(t)
		seedRegistry(t, st, descriptor("tools.echo", "alpha"))
		heads := &countingHeads{inner: st}
		reader := &countingReader{inner: transitionreg.NewReader(st)}
		cfg := testConfig(st)
		cfg.Heads, cfg.Reader = heads, reader
		h := mustHandler(t, cfg)
		expiredTok := mintExpiredToken(t, st)
		getCard(t, h, "Bearer "+expiredTok)
		postA2A(t, h, "Bearer "+expiredTok, tasksSendBody("tools.echo"))
		if heads.calls != 0 || reader.calls != 0 {
			t.Fatalf("expired-session requests touched the registry: heads.calls=%d reader.calls=%d, want 0/0", heads.calls, reader.calls)
		}
	})

	t.Run("valid session + malformed JSON never reaches admission", func(t *testing.T) {
		heads, reader, h, tok := mk(t)
		postA2A(t, h, "Bearer "+tok, "this is not json")
		postA2A(t, h, "Bearer "+tok, `{"jsonrpc":"2.0","id":1,"method":"tasks/list"}`)
		if heads.calls != 0 || reader.calls != 0 {
			t.Fatalf("parse/method failures touched the registry: heads.calls=%d reader.calls=%d, want 0/0", heads.calls, reader.calls)
		}
	})

	t.Run("control: a real admission check reads exactly once each", func(t *testing.T) {
		heads, reader, h, tok := mk(t)
		postA2A(t, h, "Bearer "+tok, tasksSendBody("tools.echo"))
		if heads.calls != 1 || reader.calls != 1 {
			t.Fatalf("control admission read counts = %d/%d, want exactly 1/1 (zero would make the assertions above vacuous)", heads.calls, reader.calls)
		}
	})
}

// ---------------------------------------------------------------------------
// AC7 — one snapshot per request
// ---------------------------------------------------------------------------

// TestProjection_OneSnapshotPerRequest is AC7 (MUT-SPLIT-SNAPSHOT's killer):
// one card request performs EXACTLY ONE head pre-check and EXACTLY ONE
// registry snapshot read, and the card derives from that one captured
// snapshot. A mutation that re-reads (a second NewRequest) pushes the reader
// count to 2 and REDs.
func TestProjection_OneSnapshotPerRequest(t *testing.T) {
	st := openStore(t)
	seedRegistry(t, st, descriptor("tools.echo", "alpha"))
	heads := &countingHeads{inner: st}
	reader := &countingReader{inner: transitionreg.NewReader(st)}
	cfg := testConfig(st)
	cfg.Heads, cfg.Reader = heads, reader
	h := mustHandler(t, cfg)
	tok := mintToken(t, st, "ep-a", []broker.Capability{liveGrant("alpha")})

	rec := getCard(t, h, "Bearer "+tok)
	if rec.Code != http.StatusOK {
		t.Fatalf("status=%d body=%s", rec.Code, rec.Body)
	}
	if heads.calls != 1 || reader.calls != 1 {
		t.Fatalf("one card request caused heads.calls=%d reader.calls=%d, want exactly 1/1 (never mixes epochs)", heads.calls, reader.calls)
	}

	heads2 := &countingHeads{inner: st}
	reader2 := &countingReader{inner: transitionreg.NewReader(st)}
	cfg2 := testConfig(st)
	cfg2.Heads, cfg2.Reader = heads2, reader2
	h2 := mustHandler(t, cfg2)
	postA2A(t, h2, "Bearer "+tok, tasksSendBody("tools.echo"))
	if heads2.calls != 1 || reader2.calls != 1 {
		t.Fatalf("one /a2a/ admission caused heads.calls=%d reader.calls=%d, want exactly 1/1", heads2.calls, reader2.calls)
	}
}

// ---------------------------------------------------------------------------
// AC12 — dynamic source
// ---------------------------------------------------------------------------

// TestAgentCard_HeadChangeWithoutRestart is AC12 (MUT-STARTUP-CACHE's
// killer): publishing a new transition-registry head changes the NEXT
// authorized card with no restart, no rebuild and no .ail edit.
func TestAgentCard_HeadChangeWithoutRestart(t *testing.T) {
	st := openStore(t)
	d := descriptor("tools.echo", "alpha")
	head1 := seedRegistry(t, st, d)
	h := mustHandler(t, testConfig(st))
	tok := mintToken(t, st, "ep-a", []broker.Capability{liveGrant("alpha"), liveGrant("beta")})

	if ids := skillIDs(t, getCard(t, h, "Bearer "+tok).Body.Bytes()); !reflect.DeepEqual(ids, []string{"tools.echo"}) {
		t.Fatalf("rev-1 card = %v, want [tools.echo]", ids)
	}

	d2 := descriptor("tools.probe", "beta")
	next, err := transitionreg.BuildNext(transitionreg.Revision{
		SemanticID: transitionreg.SemanticIDV1, InterfaceHash: transitionreg.InterfaceHashV1,
		Revision: 1, Entries: []transitionreg.Descriptor{d},
	}, []transitionreg.Change{{ID: d2.ID, Descriptor: &d2}})
	if err != nil {
		t.Fatalf("BuildNext: %v", err)
	}
	if _, err := transitionreg.NewReader(st).Publish(context.Background(), head1, next); err != nil {
		t.Fatalf("Publish rev 2: %v", err)
	}

	if ids := skillIDs(t, getCard(t, h, "Bearer "+tok).Body.Bytes()); !reflect.DeepEqual(ids, []string{"tools.echo", "tools.probe"}) {
		t.Fatalf("rev-2 card = %v, want [tools.echo tools.probe] without restart (the registry head is read per request, never cached at startup)", ids)
	}
}

// ---------------------------------------------------------------------------
// AC6 — carrier fixed (D-WORLD-26 = ARM A)
// ---------------------------------------------------------------------------

// TestProjection_CarrierFixed_OnlyBearerHeader is AC6 on both routes: the
// resolver reads ONLY the standard Authorization: Bearer value. The rejected
// Arm-B header X-World-Session is never read, even as a fallback
// (MUT-ALT-HEADER's killer), and a well-shaped but non-session Bearer value —
// the semantic stand-in for the static serve-api key — resolves to nothing
// (MUT-KEY-AS-SESSION / MUT-DEFAULT-CAPS die here): constraint (i), one
// credential path, fail closed.
func TestProjection_CarrierFixed_OnlyBearerHeader(t *testing.T) {
	st := openStore(t)
	seedRegistry(t, st, descriptor("tools.echo", "alpha"))
	h := mustHandler(t, testConfig(st))
	tok := mintToken(t, st, "ep-a", []broker.Capability{liveGrant("alpha")})

	// Arm B carried in ONLY X-World-Session — with a VALID session token —
	// must fail closed as ABSENT on both routes (the header is never read;
	// even a real credential smuggled on the rejected carrier gains nothing).
	req := httptest.NewRequest(http.MethodGet, "/.well-known/agent.json", nil)
	req.Header.Set("X-World-Session", tok)
	rec := httptest.NewRecorder()
	h.AgentCard(rec, req)
	if rec.Code != http.StatusUnauthorized {
		t.Fatalf("X-World-Session-only card request status=%d, want 401 — Arm B is REJECTED even as a fallback", rec.Code)
	}
	var body apiErrorBody
	_ = json.Unmarshal(rec.Body.Bytes(), &body)
	if body.Error.Class != "SessionAbsent" {
		t.Fatalf("X-World-Session-only card class=%q, want SessionAbsent (the only Authorization signal is the standard header)", body.Error.Class)
	}

	reqA := httptest.NewRequest(http.MethodPost, "/a2a/", strings.NewReader(tasksSendBody("tools.echo")))
	reqA.Header.Set("X-World-Session", tok)
	recA := httptest.NewRecorder()
	h.A2A(recA, reqA)
	code, _, _ := a2aErr(t, recA.Body.Bytes())
	if recA.Code != http.StatusOK || code != -32001 {
		t.Fatalf("X-World-Session-only /a2a/ request = (status %d, code %d), want (200, -32001)", recA.Code, code)
	}

	// Key-shaped Bearer (well-formed 64-hex that matches NO session row — the
	// semantic content of "the static serve-api key must never be accepted as
	// a session"): unknown, 401 / -32001, on both routes.
	keyBearer := "Bearer " + strings.Repeat("c", 64)
	if rec := getCard(t, h, keyBearer); rec.Code != http.StatusUnauthorized {
		t.Fatalf("key-shaped Bearer card status=%d, want 401", rec.Code)
	}
	recA2 := postA2A(t, h, keyBearer, tasksSendBody("tools.echo"))
	code2, _, _ := a2aErr(t, recA2.Body.Bytes())
	if code2 != -32001 {
		t.Fatalf("key-shaped Bearer /a2a/ code=%d, want -32001", code2)
	}

	// Control: the same token over the STANDARD carrier resolves the card.
	if rec := getCard(t, h, "Bearer "+tok); rec.Code != http.StatusOK {
		t.Fatalf("standard-carrier control status=%d, want 200 (the X-World-Session assertions are vacuous without this)", rec.Code)
	}
}

// ---------------------------------------------------------------------------
// Decision 6 / AC13(retained half) — the bounded wait
// ---------------------------------------------------------------------------

// blockingHeads blocks in GetRegistryHead until ctx is done — the snapshot
// half of the bounded-wait fault injection.
type blockingHeads struct{}

func (blockingHeads) GetRegistryHead(ctx context.Context, _ string) (hashref.HashRef, bool, error) {
	select {
	case <-ctx.Done():
		return hashref.HashRef{}, false, ctx.Err()
	case <-time.After(2 * time.Second):
		return hashref.HashRef{}, false, nil
	}
}

// TestProjection_BoundedWait is the retained Decision-6/AC13 half
// (MUT-DROP-DEADLINE-PROJ's killer): with a fault-injected resolver (and,
// separately, a fault-injected head check) that blocks until the context is
// done, each route answers within the configured finite bound — the card with
// 504, /a2a/ with the constant -32603 — and never hangs. Under the mutation
// (a context.Background() substituted before resolution or snapshot reads)
// the fault-injected seams hit their 2 s escape hatch, the test exceeds the
// bound, and the status/code assertions RED as well.
func TestProjection_BoundedWait(t *testing.T) {
	const maxWait = 100 * time.Millisecond
	const slack = 1500 * time.Millisecond

	t.Run("card route, resolution blocked", func(t *testing.T) {
		st := openStore(t)
		cfg := testConfig(st)
		cfg.Resolver = blockingResolver{}
		cfg.MaxWait = maxWait
		h := mustHandler(t, cfg)
		start := time.Now()
		rec := getCard(t, h, "Bearer "+strings.Repeat("d", 64))
		elapsed := time.Since(start)
		if rec.Code != http.StatusGatewayTimeout {
			t.Fatalf("blocked-resolution card status=%d, want 504; body=%s", rec.Code, rec.Body)
		}
		var body apiErrorBody
		if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil || body.Error.Class != "ProjectionDeadlineExceeded" {
			t.Fatalf("504 body = %s, want envelope class ProjectionDeadlineExceeded", rec.Body)
		}
		if elapsed > slack {
			t.Fatalf("blocked resolution returned in %v — beyond the configured bound (%v + slack): the request context was dropped or replaced", elapsed, maxWait)
		}
	})

	t.Run("card route, head check blocked", func(t *testing.T) {
		st := openStore(t)
		tok := mintToken(t, st, "ep-a", []broker.Capability{liveGrant("alpha")})
		cfg := testConfig(st)
		cfg.Heads = blockingHeads{}
		cfg.MaxWait = maxWait
		h := mustHandler(t, cfg)
		start := time.Now()
		rec := getCard(t, h, "Bearer "+tok)
		elapsed := time.Since(start)
		if rec.Code != http.StatusGatewayTimeout {
			t.Fatalf("blocked head-check card status=%d, want 504; body=%s", rec.Code, rec.Body)
		}
		if elapsed > slack {
			t.Fatalf("blocked head check returned in %v — beyond the configured bound", elapsed)
		}
	})

	t.Run("/a2a/ route, resolution blocked", func(t *testing.T) {
		st := openStore(t)
		cfg := testConfig(st)
		cfg.Resolver = blockingResolver{}
		cfg.MaxWait = maxWait
		h := mustHandler(t, cfg)
		start := time.Now()
		rec := postA2A(t, h, "Bearer "+strings.Repeat("d", 64), tasksSendBody("tools.echo"))
		elapsed := time.Since(start)
		code, msg, _ := a2aErr(t, rec.Body.Bytes())
		if rec.Code != http.StatusOK || code != -32603 || msg != "transition invocation is not available in this daemon" {
			t.Fatalf("/a2a/ blocked resolution = (status %d, code %d, msg %q), want (200, -32603, constant not-available message)",
				rec.Code, code, msg)
		}
		if elapsed > slack {
			t.Fatalf("/a2a/ blocked resolution returned in %v — beyond the configured bound", elapsed)
		}
	})
}

// TestProjection_ConfigValidation is Decision 6's startup half: the bound is
// validated at construction — zero/negative/omitted is an ERROR, never
// silently "unlimited" — and every injected seam is required.
func TestProjection_ConfigValidation(t *testing.T) {
	st := openStore(t)
	base := testConfig(st)

	bad := base
	bad.MaxWait = 0
	if _, err := New(bad); err == nil {
		t.Fatal("New with MaxWait=0: want an error (zero/negative/omitted is never 'unlimited')")
	}
	bad2 := base
	bad2.MaxWait = -time.Second
	if _, err := New(bad2); err == nil {
		t.Fatal("New with MaxWait<0: want an error")
	}

	strip := []struct {
		name string
		fn   func(*Config)
	}{
		{"Resolver", func(c *Config) { c.Resolver = nil }},
		{"Reader", func(c *Config) { c.Reader = nil }},
		{"Heads", func(c *Config) { c.Heads = nil }},
		{"Deny", func(c *Config) { c.Deny = nil }},
		{"Fail", func(c *Config) { c.Fail = nil }},
	}
	for _, tc := range strip {
		c := base
		tc.fn(&c)
		if _, err := New(c); err == nil {
			t.Fatalf("New with nil %s: want an error (the projection never invents its own seams)", tc.name)
		}
	}
	if _, err := New(base); err != nil {
		t.Fatalf("New with the valid base config: %v (control)", err)
	}
}

// ---------------------------------------------------------------------------
// Source-level gates (AC1 wire ownership, AC14 deadline tampering, zero-skip)
// ---------------------------------------------------------------------------

// projectionSources returns the package's non-test Go sources.
func projectionSources(t *testing.T) map[string]string {
	t.Helper()
	files, err := filepath.Glob(filepath.Join(".", "*.go"))
	if err != nil {
		t.Fatalf("glob projection files: %v", err)
	}
	out := map[string]string{}
	for _, f := range files {
		if strings.HasSuffix(f, "_test.go") {
			continue
		}
		data, err := os.ReadFile(f)
		if err != nil {
			t.Fatalf("read %s: %v", f, err)
		}
		out[f] = string(data)
	}
	if len(out) == 0 {
		t.Fatal("source scan checked zero non-test projection files — non-vacuous gate vacuous")
	}
	return out
}

// TestProjection_WireOwnershipSource is the source half of AC1: all A2A wire
// forms on this surface come from the pinned protocol package — the source
// must import it, must hand-format NO JSON-RPC bytes (no "jsonrpc" literal),
// must declare NO parallel wire struct (no json-tagged `jsonrpc` field), and
// must never route skill IDs through the MCP name grammar (no CallerSurface /
// ValidateMCPName reference — F6) nor read the rejected Arm-B header. Kills
// MUT-PROTO-OWNER and the source half of MUT-CARDSURFACE / MUT-ALT-HEADER.
func TestProjection_WireOwnershipSource(t *testing.T) {
	srcs := projectionSources(t)

	sawProtocolImport := false
	for f, src := range srcs {
		if strings.Contains(src, `"github.com/sunholo-data/ailang/serveapi/protocol"`) {
			sawProtocolImport = true
		}
		for i, line := range strings.Split(src, "\n") {
			if strings.HasPrefix(trimCommentSafe(line), "//") {
				continue
			}
			if strings.Contains(line, `"jsonrpc"`) || strings.Contains(line, "`json:\"jsonrpc\"`") {
				t.Errorf("%s:%d: hand-formatted JSON-RPC wire material (a parallel wire shape is forbidden — AC1): %q", f, i+1, strings.TrimSpace(line))
			}
			if strings.Contains(line, "CallerSurface") || strings.Contains(line, "ValidateMCPName") {
				t.Errorf("%s:%d: reference to the MCP name-grammar gate (F6 — card IDs are verbatim, never CallerSurface): %q", f, i+1, strings.TrimSpace(line))
			}
			if strings.Contains(line, "X-World-Session") {
				t.Errorf("%s:%d: reference to the REJECTED Arm-B header (D-WORLD-26, never a fallback): %q", f, i+1, strings.TrimSpace(line))
			}
		}
	}
	if !sawProtocolImport {
		t.Error("no projection source imports the pinned serveapi/protocol package — AC1 wire reuse would be vacuous")
	}
	if !t.Failed() {
		t.Logf("wire ownership pinned over %d non-test source file(s)", len(srcs))
	}
}

// trimCommentSafe returns the line unchanged; it exists as the comment-strip
// hook so a future AST pass can replace line scanning without touching the
// assertion table. (Keep line-based: a `go`-style full-line comment is the
// only comment form the projection sources use inside function bodies.)
func trimCommentSafe(line string) string { return strings.TrimSpace(line) }

// TestProjection_NoDeadlineTampering is AC14's source assertion
// (MUT-DEADLINE-RELAX's killer): no projection source calls
// ResponseController / SetWriteDeadline / SetReadDeadline or otherwise
// relaxes the frozen D7 deadlines. (The REST-side control that D7 constants
// did not move is the same-run TestBoundedWaitsAndBodyLimit in host/daemon,
// plus read_deadline_test.go's D7 assertions.)
func TestProjection_NoDeadlineTampering(t *testing.T) {
	banned := []string{
		"Response" + "Controller",
		"SetWrite" + "Deadline",
		"SetRead" + "Deadline",
	}
	for f, src := range projectionSources(t) {
		for i, line := range strings.Split(src, "\n") {
			for _, b := range banned {
				if strings.Contains(line, b) {
					t.Errorf("%s:%d: deadline-tampering call %q — projection inherits the frozen D7 deadlines unchanged (AC14)", f, i+1, strings.TrimSpace(line))
				}
			}
		}
	}
}

// TestProjection_ZeroSkipSource is the zero-skip clause (travelling with the
// descoped AC10) as a source gate over THIS package's tests
// (MUT-SKIP-SOCKET's re-shaped target — a t.Skip anywhere in
// host/projection's tests REDs): a skipped gate is a false green.
func TestProjection_ZeroSkipSource(t *testing.T) {
	files, err := filepath.Glob(filepath.Join(".", "*_test.go"))
	if err != nil {
		t.Fatalf("glob: %v", err)
	}
	if len(files) == 0 {
		t.Fatal("zero test files scanned — the gate would pass vacuously")
	}
	needle := "t." + "Skip"
	for _, f := range files {
		data, err := os.ReadFile(f)
		if err != nil {
			t.Fatalf("read %s: %v", f, err)
		}
		for i, line := range strings.Split(string(data), "\n") {
			if strings.Contains(line, needle+"(") || strings.Contains(line, needle+"f(") || strings.Contains(line, needle+"Now(") {
				t.Errorf("%s:%d: skip call in a projection test: %q — the projection gates assert, they never skip", f, i+1, strings.TrimSpace(line))
			}
		}
	}
}
