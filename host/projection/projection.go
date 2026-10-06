// Package projection exposes a session-scoped agent card, A2A admission,
// and MCP tools over the daemon's shared authority resolver. MCP wire and
// callback execution belong to the pinned protocol seams; each invocation
// receives fresh World admission and a fresh coordinator task ID.
//
// The surface holds no credential path or store handle of its own. The card is
// read-only; an admitted A2A call can invoke the injected coordinator. Both
// routes resolve the session with Resolver.ResolveContext over one bounded
// request context.
//
// Authorization on this surface is conceptually
//
//	visible(session, transition) = session is live AND
//	                               transition capability ∈ session capabilities
//
// (explanatory prose, not source): the implementation reuses the landed
// capability predicate via transitionreg.Request.Allowed over broker.Allows —
// this package is NOT a second policy engine.
//
// Wire ownership (AC1): the /a2a/ route parses and emits EXCLUSIVELY through
// the pinned github.com/sunholo-data/ailang/serveapi/protocol helpers
// (A2ARequest in, A2AError/A2AResult out); the card is a map[string]any literal with the
// SAME keys the upstream serveapi handler emits; the REST APIError envelope
// used for card-route denials is DAEMON-OWNED — injected as function values at
// mount — and never hand-formatted here. Skill IDs are World stable
// transition IDs emitted VERBATIM (F6): they never pass through
// protocol.CallerSurface / ValidateMCPName, whose MCP name grammar rejects
// them.
package projection

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/sunholo-data/ailang/serveapi/protocol"

	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/coordinator"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/transitionreg"
)

// JSON-RPC 2.0 codes this surface emits, per the design's authoritative
// Denial matrix and Decision 3 (B4). protocol.A2AError ALWAYS writes HTTP 200
// (F5b), so every /a2a/ failure is a JSON-RPC error body on a 200 carrier.
const (
	// codeSessionDenied (-32001) is the implementation-defined server-error
	// code (JSON-RPC 2.0 reserved range -32000..-32099) for an absent,
	// unknown or expired session credential on /a2a/.
	codeSessionDenied = -32001
	// codeInvalidRequest (-32600) covers a malformed session credential, an
	// undecodable body, and a jsonrpc value other than "2.0".
	codeInvalidRequest = -32600
	// codeMethodNotFound (-32601) covers any method other than "tasks/send".
	codeMethodNotFound = -32601
	// codeInvalidParams (-32602) covers undecodable params and a skill_id not
	// in this session's Allowed set (unlisted/guessed/stale names).
	codeInvalidParams = -32602
	// codeInternal (-32603) covers resolution, registry, and invocation failures.
	codeInternal = -32603
)

// Constant wire messages. None of these is ever interpolated with a skill ID,
// credential material, request content or store detail (constant-message
// rules in the Denial matrix and B4).
const (
	msgAbsent        = "session credential is absent: send Authorization: Bearer <session-credential>"
	msgUnknown       = "unknown session credential"
	msgExpired       = "session credential has expired"
	msgMalformed     = "malformed Authorization header: expected Bearer <64-hex-credential>"
	msgInvalidReq    = "invalid JSON-RPC request"
	msgMethodUnknown = "method not found"
	msgInvalidParams = "invalid params"
	// msgNotAuthorized is Decision 3 (B4)'s exact "not authorized" refusal for
	// a skill_id outside the session's Allowed set. Stale and guessed names
	// get the same refusal; the message carries no name back to the caller.
	msgNotAuthorized = "not authorized"
	// notAvailableMessage covers unavailable dispatch, including a daemon
	// without an archived interpreter, and resolution/registry failures.
	notAvailableMessage = "transition invocation is not available in this daemon"
)

// Card-route failure envelope classes for the non-denial failures, written
// through the daemon-injected ErrorWriter so the APIError envelope itself
// stays daemon-owned (B1). These classes are NEW and exist only on the NEW
// card route — no frozen /v1/ vocabulary changes (P6).
const (
	classUnavailable      = "ProjectionUnavailable"
	msgUnavailable        = "the projection read failed; retry the request"
	classDeadlineExceeded = "ProjectionDeadlineExceeded"
	msgDeadlineExceeded   = "the projection read exceeded its bounded deadline"
)

// maxA2ARequestBytes bounds the /a2a/ request body (1 MiB — the same cap the
// upstream serveapi A2A handler applies).
const maxA2ARequestBytes = 1 << 20

// HeadReader is the B3 absent-head pre-check seam. The daemon wires its
// readStore (daemon.go:331-337), so the projection's GetRegistryHead call goes
// through the same read seam every /v1 GET route uses.
type HeadReader interface {
	GetRegistryHead(ctx context.Context, name string) (hashref.HashRef, bool, error)
}

// DenyWriter renders one typed session denial in the caller's wire envelope.
// The daemon injects its writeAPIError-backed writeSessionDenial, so a
// card-route denial is byte-identical to the /v1/commit middleware's (B1):
// the envelope stays daemon-owned and this package never formats it (AC1).
type DenyWriter func(w http.ResponseWriter, kind authority.DenialKind)

// ErrorWriter renders a plain REST error envelope for the card route's
// non-denial failures (503/504). The daemon injects writeAPIError.
type ErrorWriter func(w http.ResponseWriter, class, message string, status int)

// Config is the mount-time injection surface the daemon fills from its ONE
// existing handle set (F3): no new store, no second resolver, no second
// credential path.
type Config struct {
	// Resolver is the daemon's own authority.Resolver instance.
	Resolver authority.Resolver
	// Reader reads transition-registry snapshots (prod:
	// transitionreg.NewReader over the daemon's store).
	Reader transitionreg.Reader
	// Heads is the GetRegistryHead pre-check seam (prod: the daemon's
	// readStore).
	Heads HeadReader
	// Deny renders session denials on the card route (daemon envelope).
	Deny DenyWriter
	// Fail renders the card route's non-denial 503/504 failures.
	Fail ErrorWriter
	// ErrorLog is the required operator sink for sanitized A2A refusals.
	ErrorLog io.Writer
	// Agent identities the card's name/description/version fields.
	Agent protocol.AgentInfo
	// MaxWait is the ONE bounded wait per projection request (Decision 6):
	// the deadline is the earlier of the client cancellation/deadline and
	// this finite, positive server maximum, passed WITHOUT replacement
	// through session resolution and the registry/capability snapshot reads.
	// Zero/negative is rejected at startup, never silently "unlimited".
	MaxWait         time.Duration
	InvokeWait      time.Duration
	Coordinator     *coordinator.Coordinator
	CredentialWait  time.Duration
	CallbackTimeout time.Duration
	MaxCallbacks    int
	WriteWait       time.Duration
}

// Handler serves the two projection routes. Both handlers run resolution and
// snapshot reads on the CALLING goroutine (no worker is ever spawned — that is
// what makes the bounded ctx the only wait), and both are safe for concurrent
// requests: every request builds its OWN fresh session, snapshot and request
// set.
type Handler struct {
	resolver       authority.Resolver
	reader         transitionreg.Reader
	heads          HeadReader
	deny           DenyWriter
	fail           ErrorWriter
	errorLog       io.Writer
	agent          protocol.AgentInfo
	maxWait        time.Duration
	invokeWait     time.Duration
	coord          *coordinator.Coordinator
	credentialWait time.Duration
	mcp            http.Handler
	mintTask       func() (string, error)
}

// New validates the config at startup (Decision 6: a zero/negative/omitted
// bound is a construction error, never "unlimited") and returns the handler.
func New(cfg Config) (*Handler, error) {
	switch {
	case cfg.Resolver == nil:
		return nil, errors.New("projection: Resolver is required")
	case cfg.Reader == nil:
		return nil, errors.New("projection: Reader is required")
	case cfg.Heads == nil:
		return nil, errors.New("projection: Heads is required")
	case cfg.Deny == nil:
		return nil, errors.New("projection: Deny is required")
	case cfg.Fail == nil:
		return nil, errors.New("projection: Fail is required")
	case cfg.ErrorLog == nil:
		return nil, errors.New("projection: ErrorLog is required")
	case cfg.MaxWait <= 0:
		return nil, fmt.Errorf("projection: MaxWait must be finite and positive, got %v", cfg.MaxWait)
	case cfg.CredentialWait <= 0 || cfg.CallbackTimeout <= 0 || cfg.MaxCallbacks <= 0 || cfg.WriteWait <= 0 || cfg.InvokeWait >= cfg.WriteWait:
		return nil, errors.New("projection: MCP bounds must be positive and invocation below write timeout")
	case cfg.InvokeWait <= 0:
		return nil, fmt.Errorf("projection: InvokeWait must be finite and positive, got %v", cfg.InvokeWait)
	}
	h := &Handler{
		resolver:       cfg.Resolver,
		reader:         cfg.Reader,
		heads:          cfg.Heads,
		deny:           cfg.Deny,
		fail:           cfg.Fail,
		errorLog:       cfg.ErrorLog,
		agent:          cfg.Agent,
		maxWait:        cfg.MaxWait,
		invokeWait:     cfg.InvokeWait,
		coord:          cfg.Coordinator,
		credentialWait: cfg.CredentialWait, mintTask: mintMCPTask,
	}
	var err error
	h.mcp, err = newMCPHandler(h, cfg)
	if err != nil {
		return nil, err
	}
	return h, nil
}

// bindingCaps adapts the authority binding's grant set to
// transitionreg.CapabilitySource. The binding (row 39) already carries an
// IMMUTABLE per-session capability snapshot (P2), so the capability filter
// needs no live broker session — and TR.C's dispatch-binding gate
// (host/broker) forbids constructing one outside host/broker anyway. The
// snapshot reading itself is built by broker.NewCapabilitySnapshot: this type
// adds nothing but the interface method, so the capability policy still lives
// exactly in broker.Allows (no second policy engine here).
type bindingCaps struct{ caps []broker.Capability }

// CapabilitySnapshot satisfies transitionreg.CapabilitySource.
func (c bindingCaps) CapabilitySnapshot(now int64) broker.CapabilitySnapshot {
	return broker.NewCapabilitySnapshot(c.caps, now)
}

// AgentCard serves GET /.well-known/agent.json (B1/B2/B3): resolve the
// session, then emit the upstream-keyshaped card whose skills are EXACTLY
// this session's allowed transition IDs, verbatim, in registry bytewise
// order. A denied session gets the daemon's REST APIError envelope with the
// F2 status mapping; an absent registry head is a legitimate ZERO-skills card
// at HTTP 200 (still authenticated); a read failure is 503/504 fail-closed.
func (h *Handler) AgentCard(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.maxWait)
	defer cancel()
	out, err := h.resolver.ResolveContext(ctx, r.Header.Get("Authorization"), time.Now().Unix())
	if err != nil {
		h.writeUnavailable(w, err)
		return
	}
	if out.Denied != nil {
		h.deny(w, *out.Denied)
		return
	}
	_, allowed, serr := h.allowedDescriptors(ctx, out.Success)
	if serr != nil {
		h.writeUnavailable(w, serr)
		return
	}
	skills := make([]map[string]any, 0, len(allowed))
	for _, d := range allowed {
		skills = append(skills, map[string]any{
			// The stable transition ID is the tool identity, emitted VERBATIM
			// (F6: never through CallerSurface / ValidateMCPName).
			"id": d.ID, "name": d.Title, "description": d.Description,
			"tags": []string{}, "examples": []string{},
		})
	}
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	card := map[string]any{
		"name": h.agent.Name, "description": h.agent.Description,
		"url": scheme + "://" + r.Host, "version": h.agent.Version,
		"capabilities":      map[string]any{"streaming": false, "pushNotifications": false, "stateTransitionHistory": false},
		"defaultInputModes": []string{"application/json"}, "defaultOutputModes": []string{"application/json"},
		"skills": skills,
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_ = json.NewEncoder(w).Encode(card)
}

// A2A serves POST /a2a/: admission and invocation share one bounded context
// and one captured registry request. Denials precede body parsing, so a denied
// request never acquires a registry snapshot.
func (h *Handler) A2A(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), h.invokeWait)
	defer cancel()
	out, err := h.resolver.ResolveContext(ctx, r.Header.Get("Authorization"), time.Now().Unix())
	if err != nil {
		protocol.A2AError(w, nil, codeInternal, notAvailableMessage)
		return
	}
	if out.Denied != nil {
		switch *out.Denied {
		case authority.DenialAbsent:
			protocol.A2AError(w, nil, codeSessionDenied, msgAbsent)
		case authority.DenialUnknown:
			protocol.A2AError(w, nil, codeSessionDenied, msgUnknown)
		case authority.DenialExpired:
			protocol.A2AError(w, nil, codeSessionDenied, msgExpired)
		case authority.DenialMalformed:
			protocol.A2AError(w, nil, codeInvalidRequest, msgMalformed)
		}
		return
	}
	var req protocol.A2ARequest
	if err := json.NewDecoder(http.MaxBytesReader(w, r.Body, maxA2ARequestBytes)).Decode(&req); err != nil {
		protocol.A2AError(w, nil, codeInvalidRequest, msgInvalidReq)
		return
	}
	if req.JSONRPC != "2.0" {
		protocol.A2AError(w, req.ID, codeInvalidRequest, msgInvalidReq)
		return
	}
	if req.Method != "tasks/send" {
		protocol.A2AError(w, req.ID, codeMethodNotFound, msgMethodUnknown)
		return
	}
	var params protocol.A2ATaskSendParams
	if err := json.Unmarshal(req.Params, &params); err != nil {
		protocol.A2AError(w, req.ID, codeInvalidParams, msgInvalidParams)
		return
	}
	// The skill name arrives exactly where the upstream handler reads it
	// (F5b): params.Metadata["skill_id"].(string).
	skillID, _ := params.Metadata["skill_id"].(string)
	// The admission check builds its OWN fresh per-request snapshot set:
	// possession of a previously listed skill name conveys no authority (B2's
	// one-snapshot rule, applied at /a2a/ admission time).
	admitted, allowed, serr := h.allowedDescriptors(ctx, out.Success)
	if serr != nil {
		h.logRefusal("a2a", req.Method, "-", serr)
		protocol.A2AError(w, req.ID, codeInternal, notAvailableMessage)
		return
	}
	for _, d := range allowed {
		if d.ID == skillID {
			if h.coord == nil {
				h.logRefusal("a2a", req.Method, "-", errors.New("projection: invocation coordinator is unavailable"))
				protocol.A2AError(w, req.ID, codeInternal, notAvailableMessage)
				return
			}
			if len(params.Message.Parts) != 1 || params.Message.Parts[0].Type != "data" || params.Message.Parts[0].Data == nil {
				protocol.A2AError(w, req.ID, codeInvalidParams, msgInvalidParams)
				return
			}
			var pin *hashref.HashRef
			if raw, ok := params.Metadata["transition_fn"]; ok {
				s, valid := raw.(string)
				if !valid {
					protocol.A2AError(w, req.ID, codeInvalidParams, msgInvalidParams)
					return
				}
				parsed, err := hashref.Parse(s)
				if err != nil {
					protocol.A2AError(w, req.ID, codeInvalidParams, msgInvalidParams)
					return
				}
				pin = &parsed
			}
			result, err := h.coord.Dispatch(ctx, coordinator.Call{Request: admitted, EpisodeID: out.Success.EpisodeID,
				Grants: out.Success.Caps, SkillID: d.ID, TaskID: params.ID, Input: params.Message.Parts[0].Data, PinnedFn: pin})
			if err != nil {
				// A timed-out invocation must release its transport as well as its
				// capsule. The client can then retry with the same task ID after
				// reconciling an uncertain durable outcome.
				if ctx.Err() != nil {
					w.Header().Set("Connection", "close")
				}
				code, msg := dispatchError(err)
				if msg == notAvailableMessage {
					h.logRefusal("a2a", req.Method, coordinator.InvocationID(out.Success.EpisodeID, params.ID), err)
				}
				protocol.A2AError(w, req.ID, code, msg)
				return
			}
			task := map[string]any{"id": params.ID, "status": map[string]any{"state": "completed"},
				"artifacts": []any{map[string]any{"parts": []any{map[string]any{"type": "data", "data": result.Output}}}},
				"metadata":  map[string]any{"invocation_id": result.InvocationID, "world_ref": result.WorldRef, "entry_index": result.EntryIndex}}
			protocol.A2AResult(w, req.ID, task)
			return
		}
	}
	// Unlisted, guessed or stale: never reach dispatch, never name the gap.
	protocol.A2AError(w, req.ID, codeInvalidParams, msgNotAuthorized)
}

// allowedDescriptors captures ONE registry snapshot + ONE capability snapshot
// for this request (transitionreg.NewRequest over the freshly constructed
// broker session — a pure constructor whose CapabilitySnapshot is the
// CapabilitySource, F9) and returns the admitted descriptors in registry
// bytewise order.
//
// An empty success requires absence in BOTH the precheck and the snapshot's
// typed result. A raced-in head's store, integrity or context failure remains
// a failure; present-precheck then absent-snapshot remains unavailable. A
// successful fresh snapshot always wins over the earlier precheck.
func (h *Handler) allowedDescriptors(ctx context.Context, b *authority.SessionBinding) (transitionreg.Request, []transitionreg.Descriptor, error) {
	_, hasHead, err := h.heads.GetRegistryHead(ctx, store.TransitionRegistryV1)
	if err != nil {
		return transitionreg.Request{}, nil, fmt.Errorf("projection: registry head check: %w", err)
	}
	req, rerr := transitionreg.NewRequest(ctx, h.reader, bindingCaps{b.Caps}, time.Now().Unix())
	if rerr != nil {
		var target transitionreg.RegistryHeadAbsentError
		if !hasHead && errors.As(rerr, &target) {
			return transitionreg.Request{}, nil, nil
		}
		return transitionreg.Request{}, nil, fmt.Errorf("projection: registry snapshot: %w", rerr)
	}
	return req, req.Allowed(), nil
}

func dispatchError(err error) (int, string) {
	var invalid *coordinator.InvalidCallError
	var absent *transitionreg.TransitionAbsentError
	var denied *transitionreg.AccessDeniedError
	var mismatch *transitionreg.ProposalMismatchError
	var duplicate *store.DuplicateInvocationError
	var inFlight *coordinator.InFlightError
	var noWorld *coordinator.WorldAbsentError
	var effects *coordinator.EffectsUnsupportedError
	var incompatible *coordinator.IncompatibleError
	var execution *coordinator.ExecutionError
	var deadline *coordinator.UnconfirmedError
	var output *coordinator.OutputError
	var notCommitted *coordinator.NotCommittedError
	var conflict *store.ConflictError
	var integrity *coordinator.IntegrityError
	var unrecorded *coordinator.EffectsUnrecordedError
	switch {
	// First: an effect already ran, whatever the cause (which it wraps — a
	// conflict, a deadline, an unconfirmed commit) says. The caller gets the
	// effect-record refs so the executed effect stays addressable (row 134).
	case errors.As(err, &unrecorded):
		return codeInternal, effectsUnrecordedMessage(unrecorded)
	case errors.As(err, &invalid):
		return codeInvalidParams, msgInvalidParams
	case errors.As(err, &absent), errors.As(err, &denied):
		return codeInvalidParams, msgNotAuthorized
	case errors.As(err, &mismatch):
		return codeInvalidParams, "proposal does not match the registered transition"
	case errors.As(err, &integrity):
		return codeInternal, notAvailableMessage
	case errors.As(err, &notCommitted):
		return codeInternal, "invocation was not committed; send a new task id"
	case errors.As(err, &deadline):
		return codeInternal, "invocation outcome is not confirmed; resend the same task id"
	case errors.As(err, &duplicate), errors.As(err, &inFlight):
		return codeInvalidParams, "task id already used in this session"
	case errors.As(err, &noWorld):
		return codeInternal, "no world is selected; commit a genesis world first"
	case errors.As(err, &effects):
		return codeInternal, "transition declares an effect this daemon has no handler for"
	case errors.As(err, &incompatible):
		return codeInternal, "transition does not implement the invocation calling convention"
	case errors.As(err, &execution):
		return codeInternal, "transition execution failed"
	case errors.Is(err, context.Canceled), errors.Is(err, context.DeadlineExceeded):
		return codeInternal, "invocation exceeded its deadline"
	case errors.As(err, &output):
		return codeInternal, "transition output is not a JSON object"
	case errors.As(err, &conflict):
		return codeInternal, "world head moved during invocation; not committed; send a new task id"
	default:
		return codeInternal, notAvailableMessage
	}
}

// EffectsUnrecordedPrefix opens the message dispatchError gives an
// *coordinator.EffectsUnrecordedError; the effect-record refs follow it.
const EffectsUnrecordedPrefix = "effects were requested but the invocation was not confirmed committed; effect records:"

func effectsUnrecordedMessage(e *coordinator.EffectsUnrecordedError) string {
	msg := EffectsUnrecordedPrefix
	if len(e.EffectRecords) == 0 {
		return msg + " none"
	}
	for _, r := range e.EffectRecords {
		msg += " " + r.String()
	}
	return msg
}

// writeUnavailable maps a resolution/registry failure on the card route to the
// bounded-wait status pair (Decision 6): 504 when the bounded deadline was
// hit waiting on the store's single pooled connection, 503 otherwise. The
// envelope comes from the daemon-injected ErrorWriter.
func (h *Handler) writeUnavailable(w http.ResponseWriter, err error) {
	if errors.Is(err, context.DeadlineExceeded) {
		h.fail(w, classDeadlineExceeded, msgDeadlineExceeded, http.StatusGatewayTimeout)
		return
	}
	h.fail(w, classUnavailable, msgUnavailable, http.StatusServiceUnavailable)
}

// logRefusal keeps arbitrary cause text on one physical operator-only line.
func (h *Handler) logRefusal(surface, method, id string, err error) {
	fmt.Fprintf(h.errorLog, "ailang-worldd: %s refusal: %s %s: %q\n", surface, method, id, err.Error())
}
