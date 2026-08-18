package daemon

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strconv"

	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
)

// APIError is the single error envelope shared by every M2.B route and the
// M2.C client. Class names and status codes mirror ApiError/httpStatus in the
// frozen, checked design_docs/sketches/worlddapi.ail exactly.
type APIError struct {
	Error APIErrorDetail `json:"error"`
}

type APIErrorDetail struct {
	Class        string `json:"class"`
	Message      string `json:"message"`
	ObservedHead string `json:"observedHead,omitempty"`
	SelectedHead string `json:"selectedHead,omitempty"`
}

type worldResponse struct {
	Ref       string `json:"ref"`
	Revision  int64  `json:"revision"`
	StateRoot string `json:"stateRoot"`
	LogHead   string `json:"logHead"`
}

type objectResponse struct {
	Hash          string  `json:"hash"`
	InterfaceHash string  `json:"interfaceHash"`
	SemanticID    string  `json:"semanticId"`
	Provenance    string  `json:"provenance"`
	Payload       *[]byte `json:"payload,omitempty"`
}

type logHeaderResponse struct {
	EntryIndex     int64  `json:"entryIndex"`
	SemanticsEpoch int64  `json:"semanticsEpoch"`
	TransitionFn   string `json:"transitionFn"`
	Interpreter    string `json:"interpreter"`
	PrevEntryHash  string `json:"prevEntryHash"`
	WrittenBy      string `json:"writtenBy"`
}

type logEntryResponse struct {
	Header        logHeaderResponse `json:"header"`
	EntryHash     string            `json:"entryHash"`
	TransitionRef string            `json:"transitionRef"`
}

type logRangeResponse struct {
	Items []logEntryResponse `json:"items"`
}

type registryResponse struct {
	Name string `json:"name"`
	Head string `json:"head"`
}

type commitRequest struct {
	ObservedHead string          `json:"observedHead"`
	Objects      []objectRequest `json:"objects"`
	NextWorld    worldRequest    `json:"nextWorld"`
	Entry        logEntryRequest `json:"entry"`
}

type objectRequest struct {
	Hash          string `json:"hash"`
	InterfaceHash string `json:"interfaceHash"`
	SemanticID    string `json:"semanticId"`
	Provenance    string `json:"provenance"`
	Payload       []byte `json:"payload"`
}

type worldRequest struct {
	Ref       string `json:"ref"`
	Revision  int64  `json:"revision"`
	StateRoot string `json:"stateRoot"`
	LogHead   string `json:"logHead"`
}

type logHeaderRequest struct {
	EntryIndex     int64  `json:"entryIndex"`
	SemanticsEpoch int64  `json:"semanticsEpoch"`
	TransitionFn   string `json:"transitionFn"`
	Interpreter    string `json:"interpreter"`
	PrevEntryHash  string `json:"prevEntryHash"`
	WrittenBy      string `json:"writtenBy"`
}

type logEntryRequest struct {
	Header        logHeaderRequest `json:"header"`
	EntryHash     string           `json:"entryHash"`
	TransitionRef string           `json:"transitionRef"`
}

type commitResponse struct {
	SelectedHead string `json:"selectedHead"`
}

func writeJSON(w http.ResponseWriter, status int, value any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(value)
}

func writeAPIError(w http.ResponseWriter, class, message string, status int) {
	writeJSON(w, status, APIError{Error: APIErrorDetail{Class: class, Message: message}})
}

func writeConflict(w http.ResponseWriter, conflict *store.ConflictError) {
	writeJSON(w, http.StatusConflict, APIError{Error: APIErrorDetail{
		Class:        "HeadConflict",
		Message:      fmt.Sprintf("conflict: observed %s selected %s", conflict.ObservedHead, conflict.SelectedHead),
		ObservedHead: conflict.ObservedHead.String(),
		SelectedHead: conflict.SelectedHead.String(),
	}})
}

func parseRef(text, field string) (hashref.HashRef, error) {
	ref, err := hashref.Parse(text)
	if err != nil {
		return hashref.HashRef{}, fmt.Errorf("%s: %w", field, err)
	}
	return ref, nil
}

// parseGenesisRef parses observedHead, and is the reason a genesis commit is
// expressible over REST at all.
//
// EXACTLY ONE field gets this lenience, and the boundary is drawn where the
// KERNEL draws it. store.Commit documents that "a nil selected head (genesis) is
// treated as a match only when c.ObservedHead is also the zero HashRef" — so the
// zero observed head is a value the kernel defines, not a hole in validation.
// Every other ref field stays strict through parseRef;
// TestGenesisRefLenienceIsExactlyOneField pins that blast radius, field by field.
//
// The empty string is NOT a second hash encoding invented at the transport
// boundary — it is the FIRST one, read in reverse. HashRef.String() already
// renders the zero value as "" (host/hashref/hashref.go), so "" is precisely the
// canonical text of the zero ref, and hashref.Parse rejects it only because
// Parse's job is to reject malformed CONTENT ADDRESSES, which the zero value is
// not. Decision 3's "one hash encoding everywhere" is satisfied by round-tripping
// it here rather than violated.
//
// PrevEntryHash deliberately does NOT get this lenience even though it is also
// "absent" at genesis. M1's own genesis convention seeds entry 0's PrevEntryHash
// from the genesis world's LogHead (host/store/store_test.go:103), i.e. always a
// real content address; and store.Commit will WRITE a zero PrevEntryHash that
// store.GetLogEntry then cannot READ BACK ("store: log entry 0 prevEntryHash:
// hashref: empty hashref text"). ARM V1's landed kernel validation now refuses
// that write. This boundary check remains as defense-in-depth and does not widen
// the daemon's carefully bounded genesis exception.
//
// Found by the M2.B sprint-evaluator (round 1, BLOCKING) and independently by the
// controller: before this, POST /v1/commit answered a genesis commit that the
// embedded store ACCEPTS with 400 "observedHead: hashref: empty hashref text", so
// the REST surface could not express a commit its own kernel supports and the
// acceptance check "a genesis+commit episode driven entirely over REST" was
// unreachable. TestRESTGenesisAndCommitAreByteEquivalent is the gate.
func parseGenesisRef(text, field string) (hashref.HashRef, error) {
	if text == "" {
		return hashref.HashRef{}, nil
	}
	return parseRef(text, field)
}

func worldJSON(w store.World) worldResponse {
	return worldResponse{Ref: w.Ref.String(), Revision: w.Revision, StateRoot: w.StateRoot.String(), LogHead: w.LogHead.String()}
}

func logJSON(entry store.LogEntry) logEntryResponse {
	return logEntryResponse{
		Header: logHeaderResponse{
			EntryIndex: entry.Header.EntryIndex, SemanticsEpoch: entry.Header.SemanticsEpoch,
			TransitionFn: entry.Header.TransitionFn.String(), Interpreter: entry.Header.Interpreter.String(),
			PrevEntryHash: entry.Header.PrevEntryHash.String(), WrittenBy: entry.Header.WrittenBy,
		},
		EntryHash: entry.EntryHash.String(), TransitionRef: entry.TransitionRef.String(),
	}
}

// clampLimit mirrors the Z3-proven sketch predicate exactly. In particular,
// non-positive requests DEFAULT to 100; they do not clamp upward to one.
func clampLimit(requested int) int {
	if requested < 1 {
		return 100
	}
	if requested > 500 {
		return 500
	}
	return requested
}

// readCtx derives the context every store read a GET handler performs runs
// under. It is DERIVED FROM THE REQUEST, not from context.Background(): a
// client that disconnects releases the store's single connection immediately
// instead of holding it for the residual deadline.
//
// In M1 it adds no deadline of its own (a cancel-only context derived from the
// request is exactly the daemon's behaviour before this item), so M1 is
// behaviour-identical and independently landable. M2 changes only this body to
// context.WithTimeout(r.Context(), d.readDeadline).
//
// EVERY caller must `defer cancel()` immediately after calling it — a
// WithCancel/WithTimeout context whose CancelFunc is never called leaks its
// context subtree (and, from M2, its timer) once per store-reading request.
func (d *Daemon) readCtx(r *http.Request) (context.Context, context.CancelFunc) {
	return context.WithCancel(r.Context())
}

// handleWorld serves GET /v1/worlds/{ref} (Decision 3) over store.GetWorld.
//
// It is the first of the four read routes that share one shape, and the shape
// is the point: parse the ref at the boundary (malformed -> 400), ask the store
// (failure -> 500), distinguish absent from malformed (absent -> 404). The two
// failure classes are kept apart deliberately — collapsing them would tell a
// client "you asked wrongly" when the truth is "it is not here yet", which is
// exactly the distinction a re-planning caller needs.
func (d *Daemon) handleWorld(w http.ResponseWriter, r *http.Request) {
	ref, err := parseRef(r.PathValue("ref"), "world ref")
	if err != nil {
		writeAPIError(w, "BadRequest", err.Error(), http.StatusBadRequest)
		return
	}
	ctx, cancel := d.readCtx(r)
	defer cancel()
	world, ok, err := d.store.GetWorld(ctx, ref)
	if err != nil {
		writeAPIError(w, "Internal", err.Error(), http.StatusInternalServerError)
		return
	}
	if !ok {
		writeAPIError(w, "NotFound", "world not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, worldJSON(world))
}

// handleObject serves GET /v1/objects/{ref} (Decision 3) over store.GetObject.
//
// The envelope is ALWAYS returned; the payload bytes are returned only for
// ?payload=true. That default is a bounded-response rule, not a preference: an
// object payload has no size bound of its own, so an unconditional payload
// would put an unbounded allocation on every object read and quietly break
// Decision 7. TestReadRoutesAndPayloadGate asserts the default response carries
// no "payload" key at all rather than an empty one, so the gate cannot decay
// into "always present, sometimes empty".
func (d *Daemon) handleObject(w http.ResponseWriter, r *http.Request) {
	ref, err := parseRef(r.PathValue("ref"), "object ref")
	if err != nil {
		writeAPIError(w, "BadRequest", err.Error(), http.StatusBadRequest)
		return
	}
	ctx, cancel := d.readCtx(r)
	defer cancel()
	object, ok, err := d.store.GetObject(ctx, ref)
	if err != nil {
		writeAPIError(w, "Internal", err.Error(), http.StatusInternalServerError)
		return
	}
	if !ok {
		writeAPIError(w, "NotFound", "object not found", http.StatusNotFound)
		return
	}
	body := objectResponse{
		Hash: object.Hash.String(), InterfaceHash: object.InterfaceHash.String(),
		SemanticID: object.SemanticID, Provenance: object.Provenance,
	}
	if r.URL.Query().Get("payload") == "true" {
		payload := object.Payload
		body.Payload = &payload
	}
	writeJSON(w, http.StatusOK, body)
}

// handleLogEntry serves GET /v1/log/{index} (Decision 3) over store.GetLogEntry.
//
// The index is an integer, not a HashRef, so the 400 arm is a parse failure or
// a negative value rather than a malformed content address.
func (d *Daemon) handleLogEntry(w http.ResponseWriter, r *http.Request) {
	index, err := strconv.ParseInt(r.PathValue("index"), 10, 64)
	if err != nil || index < 0 {
		writeAPIError(w, "BadRequest", "log index must be a non-negative integer", http.StatusBadRequest)
		return
	}
	ctx, cancel := d.readCtx(r)
	defer cancel()
	entry, ok, err := d.store.GetLogEntry(ctx, index)
	if err != nil {
		writeAPIError(w, "Internal", err.Error(), http.StatusInternalServerError)
		return
	}
	if !ok {
		writeAPIError(w, "NotFound", "log entry not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, logJSON(entry))
}

// handleLogRange serves GET /v1/log?from=N&limit=M (Decision 3): the surface's
// ONLY deliberate N+1.
//
// It is a bounded loop over the EXISTING store.GetLogEntry and adds no store
// method. That is a recorded trade, not an oversight — the alternative is a new
// range query in the kernel, and M2 declines to grow the kernel for a read the
// daemon can compose. What makes the trade honest is that it is MEASURED:
// BenchmarkLogRange runs at limit=100 (the default page) and limit=500 (the
// clamp max) and both rows are in bench/BASELINE.md, so if the N+1 ever becomes
// the reason to add a range method, the numbers will say so first.
//
// The loop is bounded twice over: by clampLimit (<= 500, the Z3-proven ceiling)
// and by the first absent index, so a `from` past the end returns an empty item
// list rather than scanning.
func (d *Daemon) handleLogRange(w http.ResponseWriter, r *http.Request) {
	from := int64(0)
	if text := r.URL.Query().Get("from"); text != "" {
		parsed, err := strconv.ParseInt(text, 10, 64)
		if err != nil || parsed < 0 {
			writeAPIError(w, "BadRequest", "from must be a non-negative integer", http.StatusBadRequest)
			return
		}
		from = parsed
	}
	requested := 0
	if text := r.URL.Query().Get("limit"); text != "" {
		parsed, err := strconv.Atoi(text)
		if err != nil {
			writeAPIError(w, "BadRequest", "limit must be an integer", http.StatusBadRequest)
			return
		}
		requested = parsed
	}
	limit := clampLimit(requested)
	// One context for the whole bounded loop: a page that goes slow mid-loop
	// times out AS A WHOLE. Nothing has been written yet (writeJSON runs once
	// at the end), so the error response replaces the page cleanly.
	ctx, cancel := d.readCtx(r)
	defer cancel()
	items := make([]logEntryResponse, 0, limit)
	for offset := 0; offset < limit; offset++ {
		entry, ok, err := d.store.GetLogEntry(ctx, from+int64(offset))
		if err != nil {
			writeAPIError(w, "Internal", err.Error(), http.StatusInternalServerError)
			return
		}
		if !ok {
			break
		}
		items = append(items, logJSON(entry))
	}
	writeJSON(w, http.StatusOK, logRangeResponse{Items: items})
}

// handleRegistry serves GET /v1/registry/{name...} (Decision 3) over
// store.GetRegistryHead.
//
// The multi-segment wildcard is load-bearing: registry semantic IDs contain
// slashes ("world/epoch-registry/v1"), and a single-segment {name} would make
// the epoch registry — the one registry that exists — unreachable. r.PathValue
// returns the UNESCAPED segment, which is what the store keys on.
func (d *Daemon) handleRegistry(w http.ResponseWriter, r *http.Request) {
	name := r.PathValue("name")
	if name == "" {
		writeAPIError(w, "BadRequest", "registry name is empty", http.StatusBadRequest)
		return
	}
	ctx, cancel := d.readCtx(r)
	defer cancel()
	ref, ok, err := d.store.GetRegistryHead(ctx, name)
	if err != nil {
		writeAPIError(w, "Internal", err.Error(), http.StatusInternalServerError)
		return
	}
	if !ok {
		writeAPIError(w, "NotFound", "registry head not found", http.StatusNotFound)
		return
	}
	writeJSON(w, http.StatusOK, registryResponse{Name: name, Head: ref.String()})
}

// handleCommit serves POST /v1/commit (Decision 3): the surface's only
// mutation, and the only route that reads a request body.
//
// Two bounds and one mapping carry the weight here:
//
//   - The body is wrapped in http.MaxBytesReader at maxCommitBytes BEFORE it is
//     read, so an oversized body is refused rather than buffered (Decision 7).
//     The cap is the Z3-proven withinCommitBytes bound; oversized renders as
//     413 PayloadTooLarge per the sketch's httpStatus.
//   - Unknown fields are rejected (DisallowUnknownFields) and trailing JSON
//     values are rejected, so a client that misspells a field learns it at the
//     boundary instead of silently committing a zero value into the kernel.
//   - store.ConflictError maps to 409 with BOTH heads in machine-readable form
//     (writeConflict), because the store's compare-and-append contract is only
//     usable if a stale caller can re-plan against the selected head. A 409 that
//     carries only prose is a dead end, which is why the test re-plans from the
//     body and commits again rather than asserting on the status code.
//
// The daemon adds NO semantics: every ref is parsed at the boundary and the
// assembled store.Commit is handed to the kernel unchanged.
func (d *Daemon) handleCommit(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, maxCommitBytes)
	defer func() { _ = r.Body.Close() }()

	raw, err := io.ReadAll(r.Body)
	if err != nil {
		var tooLarge *http.MaxBytesError
		if errors.As(err, &tooLarge) {
			writeAPIError(w, "PayloadTooLarge", fmt.Sprintf("commit body exceeds %d bytes", maxCommitBytes), http.StatusRequestEntityTooLarge)
			return
		}
		writeAPIError(w, "BadRequest", "cannot read commit body: "+err.Error(), http.StatusBadRequest)
		return
	}
	var request commitRequest
	decoder := json.NewDecoder(bytes.NewReader(raw))
	decoder.DisallowUnknownFields()
	if err := decoder.Decode(&request); err != nil {
		writeAPIError(w, "BadRequest", "invalid commit JSON: "+err.Error(), http.StatusBadRequest)
		return
	}
	var extra any
	if err := decoder.Decode(&extra); !errors.Is(err, io.EOF) {
		writeAPIError(w, "BadRequest", "invalid commit JSON: multiple values", http.StatusBadRequest)
		return
	}

	commit, err := decodeCommit(request)
	if err != nil {
		writeAPIError(w, "BadRequest", err.Error(), http.StatusBadRequest)
		return
	}
	if err := d.store.Commit(commit); err != nil {
		var conflict *store.ConflictError
		if errors.As(err, &conflict) {
			writeConflict(w, conflict)
			return
		}
		writeAPIError(w, "Internal", err.Error(), http.StatusInternalServerError)
		return
	}
	writeJSON(w, http.StatusOK, commitResponse{SelectedHead: commit.NextWorld.Ref.String()})
}

// decodeCommit maps the wire request onto store.Commit, parsing EVERY digest
// field from canonical HashRef text so no unvalidated string reaches the kernel.
//
// Field-named errors ("objects[2].interfaceHash: ...") are deliberate: a commit
// carries a dozen refs, and a bare "malformed hashref" would make a 400
// undebuggable. The only fields permitted to be empty are the two the kernel
// itself defines a zero value for — see parseGenesisRef.
func decodeCommit(r commitRequest) (store.Commit, error) {
	observed, err := parseGenesisRef(r.ObservedHead, "observedHead")
	if err != nil {
		return store.Commit{}, err
	}
	nextRef, err := parseRef(r.NextWorld.Ref, "nextWorld.ref")
	if err != nil {
		return store.Commit{}, err
	}
	stateRoot, err := parseRef(r.NextWorld.StateRoot, "nextWorld.stateRoot")
	if err != nil {
		return store.Commit{}, err
	}
	logHead, err := parseRef(r.NextWorld.LogHead, "nextWorld.logHead")
	if err != nil {
		return store.Commit{}, err
	}
	entryHash, err := parseRef(r.Entry.EntryHash, "entry.entryHash")
	if err != nil {
		return store.Commit{}, err
	}
	transitionRef, err := parseRef(r.Entry.TransitionRef, "entry.transitionRef")
	if err != nil {
		return store.Commit{}, err
	}
	transitionFn, err := parseRef(r.Entry.Header.TransitionFn, "entry.header.transitionFn")
	if err != nil {
		return store.Commit{}, err
	}
	interpreter, err := parseRef(r.Entry.Header.Interpreter, "entry.header.interpreter")
	if err != nil {
		return store.Commit{}, err
	}
	prev, err := parseRef(r.Entry.Header.PrevEntryHash, "entry.header.prevEntryHash")
	if err != nil {
		return store.Commit{}, err
	}
	objects := make([]store.Object, len(r.Objects))
	for i, object := range r.Objects {
		hash, err := parseRef(object.Hash, fmt.Sprintf("objects[%d].hash", i))
		if err != nil {
			return store.Commit{}, err
		}
		iface, err := parseRef(object.InterfaceHash, fmt.Sprintf("objects[%d].interfaceHash", i))
		if err != nil {
			return store.Commit{}, err
		}
		objects[i] = store.Object{
			Hash: hash, InterfaceHash: iface, SemanticID: object.SemanticID,
			Provenance: object.Provenance, Payload: object.Payload,
		}
	}
	return store.Commit{
		ObservedHead: observed, Objects: objects,
		NextWorld: store.World{Ref: nextRef, Revision: r.NextWorld.Revision, StateRoot: stateRoot, LogHead: logHead},
		Entry: store.LogEntry{
			Header: store.LogHeader{
				EntryIndex: r.Entry.Header.EntryIndex, SemanticsEpoch: r.Entry.Header.SemanticsEpoch,
				TransitionFn: transitionFn, Interpreter: interpreter, PrevEntryHash: prev,
				WrittenBy: r.Entry.Header.WrittenBy,
			},
			EntryHash: entryHash, TransitionRef: transitionRef,
		},
	}, nil
}
