package daemon

// Row 243 V3 sprint 3 M9 (design_docs/planned/w-world-address-bar-verbs-v3.md §2.5): the
// confirmation card's read routes, each behind approvalRoute's Host check, served as the decision
// card's are; the data is surfacepage.ConfirmationCard of the row. Nothing here writes (POST: sprint 4).

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/url"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/surface"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

const confirmIDReason = "a verb request id is vr- and 16 lowercase hex"

// handleSurfaceConfirmPage is GET /surface/confirm/{id}: the same bytes for every well-formed id.
func (d *Daemon) handleSurfaceConfirmPage(w http.ResponseWriter, r *http.Request, _ []string) {
	if !store.VerbRequestIDPattern.MatchString(r.PathValue("id")) {
		writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", confirmIDReason)
		return
	}
	setSurfaceHeaders(w)
	_, _ = w.Write(surfacepage.ConfirmPage)
}

// handleSurfaceConfirmScript is the card's script, under the same headers.
func (d *Daemon) handleSurfaceConfirmScript(w http.ResponseWriter, _ *http.Request, _ []string) {
	setSurfaceHeaders(w)
	w.Header().Set("Content-Type", "text/javascript; charset=utf-8")
	_, _ = w.Write(surfacepage.ConfirmScript)
}

// handleSurfaceConfirmation is GET /surface/confirmations/{id}: 400 for a bad id, 503 naming the
// serve restart below schema 9, 404 for an unknown one, else the card's data with the active
// passkeys by origin.
func (d *Daemon) handleSurfaceConfirmation(w http.ResponseWriter, r *http.Request, _ []string) {
	id := r.PathValue("id")
	if !store.VerbRequestIDPattern.MatchString(id) {
		writeApprovalRefusal(w, http.StatusBadRequest, "BadRequest", confirmIDReason)
		return
	}
	ctx, cancel := d.readCtx(r)
	defer cancel()
	ro, err := store.OpenCredentialsReadOnly(d.credentialsPath)
	if err != nil {
		d.approvalInternal(w, r, err)
		return
	}
	defer func() { _ = ro.Close() }()
	row, ok, err := ro.VerbRequest(ctx, id)
	switch {
	case isOlderStore(err):
		writeApprovalRefusal(w, http.StatusServiceUnavailable, "Restart",
			"the credentials store predates verb requests: restart `ailang-worldd serve` with this build once (it upgrades the store)")
		return
	case err != nil:
		d.approvalInternal(w, r, err)
		return
	case !ok:
		writeApprovalRefusal(w, http.StatusNotFound, "NotFound", "no verb request "+id)
		return
	}
	active, err := ro.ActivePasskeys(ctx)
	if err != nil {
		d.approvalInternal(w, r, err)
		return
	}
	var passkeys []surfacepage.CardPasskey
	for _, e := range active {
		passkeys = append(passkeys, surfacepage.CardPasskey{Origin: e.Passkey.Origin, RPID: e.Passkey.RPID, CredentialID: e.Passkey.CredentialID})
	}
	now := d.approveNow()
	card := surfacepage.ConfirmationCard(confirmationInput(row, now), passkeys)
	card.EpisodeHref, card.PRHref = d.confirmationLinks(ctx, ro, row, now)
	writeApprovalJSON(w, http.StatusOK, card)
}

// confirmationInput is row as the card reads it at now.
func confirmationInput(row store.VerbRequestRow, now int64) surfacepage.ConfirmationInput {
	return surfacepage.ConfirmationInput{ID: row.RequestID, Request: row.Request, Digest: row.Digest, Challenge: row.Challenge, Created: row.CreatedAt,
		Expires: row.ExpiresAt, Now: now, Outcome: row.Outcome, ApprovedBy: row.ApprovedBy, Decision: row.DecisionRef}
}

// readConfirmations fills out.Confirmations with the live, pending verb requests on one read-only
// handle: a summary decoded from the stored bytes, or Refused and no summary.
func (d *Daemon) readConfirmations(ctx context.Context, now int64, out *surfacepage.RequestsJSON) error {
	ro, err := store.OpenCredentialsReadOnly(d.credentialsPath)
	if err != nil {
		return err
	}
	defer func() { _ = ro.Close() }()
	rows, err := ro.PendingVerbRequests(ctx, now, store.MaxPendingVerbRequests)
	if err != nil {
		return err
	}
	for _, row := range rows {
		out.Confirmations = append(out.Confirmations, pendingConfirmation(row))
	}
	return nil
}

// pendingConfirmation summarises row, or refuses it when its bytes do not hash to its digest, do
// not decode as a valid world/verb-request/v1, or carry no finish approval scope.
func pendingConfirmation(row store.VerbRequestRow) surfacepage.PendingConfirmation {
	p := surfacepage.PendingConfirmation{ID: row.RequestID, Expires: row.ExpiresAt}
	r, s, refused := storedFinishConfirmation(row)
	if refused != "" {
		p.Refused = refused
		return p
	}
	p.Episode, p.Repository, p.Head, p.SHA, p.Draft = r.Episode, s.Repository, s.Head, s.SHA, &s.Draft
	return p
}

func storedFinishConfirmation(row store.VerbRequestRow) (approval.VerbRequest, broker.FinishApprovalScope, string) {
	if sum := sha256.Sum256(row.Request); hex.EncodeToString(sum[:]) != row.Digest {
		return approval.VerbRequest{}, broker.FinishApprovalScope{}, "the stored request does not match its digest"
	}
	r, err := approval.DecodeVerb(row.Request)
	if err == nil {
		err = r.Validate()
	}
	if err != nil {
		return approval.VerbRequest{}, broker.FinishApprovalScope{}, "the stored request is not a valid verb request: " + err.Error()
	}
	s, err := broker.ParseFinishApprovalScope(r.Scope)
	if err != nil {
		return approval.VerbRequest{}, broker.FinishApprovalScope{}, "the stored request's scope is not a finish approval"
	}
	return r, s, ""
}

// confirmationLinks decorates the read card with evidence from the configured
// surface. Neither link contributes to approval or to its signed message.
func (d *Daemon) confirmationLinks(ctx context.Context, ro *store.Credentials, row store.VerbRequestRow, now int64) (episodeHref, prHref string) {
	if d.follow == nil {
		return "", ""
	}
	request, scope, refused := storedFinishConfirmation(row)
	if refused != "" || request.Episode != row.EpisodeID {
		return "", ""
	}
	start, err := ro.EpisodeStart(ctx, request.Episode)
	slug := strings.ToLower(scope.Repository)
	if err != nil || start == nil || start.Episode != request.Episode || !filepath.IsAbs(start.Repo) ||
		start.Branch != scope.Head || start.Base.Branch != scope.Base ||
		start.Remote.URL != scope.Remote || githubOriginSlug(start.Remote.URL) != slug {
		return "", ""
	}

	env := d.followSnapshot().Envelope
	repo := canonicalRepo(start.Repo)
	for _, enrolled := range d.follow.config.Repos {
		if enrolled.Repo != repo {
			continue
		}
		at, err := surface.QualifyID(enrolled.ID, "episode", request.Episode)
		if err != nil {
			continue
		}
		shared, _ := surface.QualifyID(enrolled.ID, "brane", surfacepage.SharedKey)
		stack, err := surfacepage.BuildFollowStack(env, shared, at, time.Unix(now, 0))
		if err != nil || stack.Repo != enrolled.ID || stack.Front.Brane == "" {
			continue
		}
		if episodeHref != "" {
			return "", "" // ambiguous enrollment
		}
		episodeHref = "/surface?" + url.Values{"at": {at}, "brane": {stack.Front.Brane}}.Encode()
	}
	if episodeHref == "" || row.Outcome != store.VerbOutcomeApproved || row.DecisionRef == "" {
		return episodeHref, ""
	}
	return episodeHref, d.confirmedPRHref(ctx, row, request.Scope, scope, slug)
}

func (d *Daemon) confirmedPRHref(ctx context.Context, row store.VerbRequestRow, signedScope string, scope broker.FinishApprovalScope, slug string) string {
	if d.store == nil {
		return ""
	}
	ref, ok, err := d.store.GetRegistryHead(ctx, finishHead(row.EpisodeID))
	if err != nil || !ok {
		return ""
	}
	obj, ok, err := d.store.GetObject(ctx, ref)
	var finish EpisodeFinish
	if err != nil || !ok || obj.SemanticID != EpisodeFinishV1 || obj.Provenance != finishProvenance ||
		json.Unmarshal(obj.Payload, &finish) != nil || finish.Episode != row.EpisodeID ||
		finish.Outcome != "done" || finish.Approval != row.DecisionRef || finish.Scope != signedScope ||
		finish.Commit != scope.SHA || finish.Branch != scope.Head || finish.Remote != scope.Remote ||
		finish.PR.Number <= 0 || finish.PR.Draft != scope.Draft {
		return ""
	}
	href := "https://github.com/" + slug + "/pull/" + strconv.Itoa(finish.PR.Number)
	if !strings.EqualFold(finish.PR.URL, href) {
		return ""
	}
	return href
}
