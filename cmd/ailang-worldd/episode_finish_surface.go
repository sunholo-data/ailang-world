package main

// Row 243 V3 sprint 5 M12 (design §2.7): `episode finish … --confirm-on-surface`. After the commit,
// where the typed y sits, finish stores a world/verb-request/v1 (credentials schema 9), prints its
// card's links and waits — polling the credentials store, never reading the terminal — for a key's
// confirmation; before any push it re-verifies the artifact over its own digest, challenge and
// expiry (a) and reads back the world decision (b). Every refusal keeps the commit.

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/store"
)

// verbLandGrace is how long past a request's expiry the wait still polls: above the POST's 3 s
// commitBudget, so a confirmation admitted just before the expiry is seen (design §2.7 step 5).
const verbLandGrace = 10 * time.Second

// surfaceFinish is what finishConfirmOnSurface needs from the finish it confirms: the request's
// fields, the credentials store it is stored in, and the preset daemon's address for the links.
type surfaceFinish struct {
	credentials, addr, ep, sha string
	argv                       []string
	scope                      broker.FinishApprovalScope
	scopeText, title, body     string
	committer                  string
	protected                  []string
}

// verbOwn is this process's own request: its id and episode, the SHA-256 of the bytes it stored,
// the challenge it drew and the expiry it set (design §2.7 step 6(a) re-checks against these).
type verbOwn struct {
	id, ep            string
	digest, challenge [32]byte
	expires           int64
}

// canonicalFinishArgv is the request's argv (pure): ailang-worldd episode finish <ep>, then each flag
// the user set in fs.Visit's lexical order — --name for a true bool, --name value otherwise, and
// --allow-failing once per value in the order given — so interleaved spellings of one command give
// one argv, and a flag left at its default never appears (a bool set false is its default).
func canonicalFinishArgv(fs *flag.FlagSet, allowed []string, ep string) []string {
	argv := []string{"ailang-worldd", "episode", "finish", ep}
	fs.Visit(func(f *flag.Flag) {
		b, isBool := f.Value.(interface{ IsBoolFlag() bool })
		switch {
		case f.Name == "allow-failing":
			for _, v := range allowed {
				argv = append(argv, "--allow-failing", v)
			}
		case isBool && b.IsBoolFlag():
			if f.Value.String() == "true" {
				argv = append(argv, "--"+f.Name)
			}
		default:
			argv = append(argv, "--"+f.Name, f.Value.String())
		}
	})
	return argv
}

// finishVerbRequest is design §2.7 step 1 (pure): the request, validated, and its canonical bytes.
// The prompt is broker.FinishPromptText of the same fields finishPrompt renders.
func finishVerbRequest(in surfaceFinish) (approval.VerbRequest, []byte, error) {
	protected := append([]string{}, in.protected...)
	r := approval.VerbRequest{Schema: approval.VerbRequestSchema, Verb: approval.VerbFinish, Argv: in.argv, Episode: in.ep,
		Prompt: broker.FinishPromptText(in.scope, in.title, in.body, protected, in.committer), Scope: in.scopeText,
		Title: in.title, Body: in.body, Committer: in.committer, Protected: protected}
	if err := r.Validate(); err != nil {
		return r, nil, err
	}
	canon, err := approval.CanonicalVerb(r)
	return r, canon, err
}

// finishConfirmOnSurface is design §2.7 steps 1-6. It returns the world decision a verified key gave
// this finish's request, or the refusal; the caller keeps the commit on a refusal.
func finishConfirmOnSurface(ctx context.Context, c *client, env sessionEnv, in surfaceFinish, tty, stdout io.Writer) (string, string) {
	r, canon, err := finishVerbRequest(in)
	if err != nil {
		return "", "this finish cannot be confirmed on the surface (" + err.Error() + "); confirm at the terminal"
	}
	if env.random == nil {
		return "", "internal error: no randomness source (sessionEnv.random)"
	}
	var idb [8]byte
	var ch [32]byte
	if err := env.random(idb[:]); err != nil {
		return "", "randomness: " + err.Error()
	}
	if err := env.random(ch[:]); err != nil {
		return "", "randomness: " + err.Error()
	}
	now := env.now()
	own := verbOwn{id: "vr-" + hex.EncodeToString(idb[:]), ep: in.ep, digest: sha256.Sum256(canon), challenge: ch,
		expires: now + int64(approval.VerbRequestLifetime/time.Second)}
	links, err := surfaceCardLinks(ctx, in.credentials, "/surface/confirm/"+own.id, in.addr)
	if err != nil {
		return "", fmt.Sprintf("cannot read the credentials store %s: %v", in.credentials, err)
	}
	row := store.VerbRequestRow{RequestID: own.id, Request: canon, Digest: hex.EncodeToString(own.digest[:]),
		Challenge: hex.EncodeToString(own.challenge[:]), EpisodeID: in.ep, CreatedAt: now, ExpiresAt: own.expires}
	if refusal := putFinishVerbRequest(in.credentials, row, now); refusal != "" {
		return "", refusal
	}
	h, err := approval.VerbMessageHash(own.digest, own.challenge, own.expires)
	if err != nil {
		return "", fmt.Sprintf("request %s: %v", own.id, err)
	}
	fmt.Fprint(tty, r.Prompt) // the terminal shows what the card shows; nothing is read from it
	fmt.Fprintf(stdout, "confirm request %s on the surface (expires %s, message hash %s):\n", own.id, approverTime(own.expires), hex.EncodeToString(h[:8]))
	for _, link := range links {
		fmt.Fprintf(stdout, "  open %s\n", link)
	}
	fmt.Fprintf(stdout, "the typed y is not read in this mode; Ctrl-C stops waiting and keeps the commit at %s\n", in.sha)
	got, refusal := awaitVerbConfirmation(env, in.credentials, own)
	if refusal != "" {
		return "", refusal
	}
	if refusal := readBackFinishDecision(worldReader{ctx: ctx, c: c}, got.DecisionRef, got.ApprovedBy, in.scopeText); refusal != "" {
		return "", refusal
	}
	fmt.Fprintf(stdout, "confirmed by %s (decision %s)\n", got.ApprovedBy, got.DecisionRef)
	return got.DecisionRef, ""
}

// surfaceCardLinks is writeAwaitLinks' address rule for a card path: one link per distinct active
// passkey origin, then the preset daemon's loopback unless an origin already is it.
func surfaceCardLinks(ctx context.Context, credPath, path, addr string) ([]string, error) {
	ctx, cancel := context.WithTimeout(ctx, approveReadBudget)
	defer cancel()
	ro, err := store.OpenCredentialsReadOnly(credPath)
	if err != nil {
		return nil, err
	}
	active, err := ro.ActivePasskeys(ctx)
	_ = ro.Close()
	if err != nil {
		return nil, err
	}
	var links []string
	seen := map[string]bool{}
	for _, e := range active {
		if !seen[e.Passkey.Origin] {
			seen[e.Passkey.Origin] = true
			links = append(links, e.Passkey.Origin+path)
		}
	}
	if _, port, err := net.SplitHostPort(addr); err == nil && !seen["http://localhost:"+port] {
		links = append(links, "http://localhost:"+port+path)
	}
	return links, nil
}

// putFinishVerbRequest is design §2.7 step 3: one non-waiting writer under its own 3 s root, the
// put, close. It returns the refusal text, or "".
func putFinishVerbRequest(credPath string, row store.VerbRequestRow, now int64) string {
	ctx, cancel := context.WithTimeout(context.Background(), approveReadBudget)
	defer cancel()
	st, err := store.OpenCredentials(credPath)
	if store.IsWriterAlreadyActive(err) {
		return "the credentials store is busy (another writer holds it); run finish again"
	}
	if err != nil {
		return fmt.Sprintf("cannot open the credentials store %s: %v", credPath, err)
	}
	defer func() { _ = st.Close() }()
	if err := st.PutVerbRequest(ctx, row, now); err != nil {
		return fmt.Sprintf("cannot store request %s: %v", row.RequestID, err)
	}
	return ""
}

// awaitVerbConfirmation is design §2.7 step 5: a poll every awaitPollInterval through env.after
// until the row is approved and re-checked (returned), superseded, unreadable, or the wait reaches
// expires_at + verbLandGrace. It never reads the terminal.
func awaitVerbConfirmation(env sessionEnv, credPath string, own verbOwn) (store.VerbRequestRow, string) {
	deadline := own.expires + int64(verbLandGrace/time.Second)
	for {
		now := env.now()
		row, refusal, err := pollVerbRequest(credPath, own, now)
		switch {
		case err != nil:
			return store.VerbRequestRow{}, fmt.Sprintf("cannot read request %s: %v", own.id, err)
		case refusal != "":
			return store.VerbRequestRow{}, refusal
		case row.Outcome == store.VerbOutcomeApproved:
			return row, ""
		case row.Outcome == store.VerbOutcomeSuperseded:
			return store.VerbRequestRow{}, fmt.Sprintf("request %s was superseded by a newer finish of %s", own.id, own.ep)
		case now >= deadline:
			return store.VerbRequestRow{}, fmt.Sprintf("no confirmation of request %s by %s", own.id, approverTime(now))
		}
		<-env.after(awaitPollInterval)
	}
}

// pollVerbRequest is one poll (design §2.7 steps 5 and 6(a)): one read-only open under its own 3 s
// root, the row, and for an approved row the re-check in the same open — the row's digest,
// challenge and expiry must be this process's own, authority.RecheckVerb must verify its artifact
// over this process's own values against the enrolled key, and that key's approver must be the
// row's. It returns the row and a refusal ("" when none).
func pollVerbRequest(credPath string, own verbOwn, now int64) (store.VerbRequestRow, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), approveReadBudget)
	defer cancel()
	ro, err := store.OpenCredentialsReadOnly(credPath)
	if err != nil {
		return store.VerbRequestRow{}, "", err
	}
	defer func() { _ = ro.Close() }()
	row, ok, err := ro.VerbRequest(ctx, own.id)
	if err != nil {
		return row, "", err
	}
	if !ok {
		return row, "", fmt.Errorf("request %s is no longer in the store", own.id)
	}
	if row.Outcome != store.VerbOutcomeApproved {
		return row, "", nil
	}
	bad := func(reason string) (store.VerbRequestRow, string, error) {
		return row, fmt.Sprintf("the confirmation of request %s does not verify: %s", own.id, reason), nil
	}
	if row.Digest != hex.EncodeToString(own.digest[:]) || row.Challenge != hex.EncodeToString(own.challenge[:]) || row.ExpiresAt != own.expires {
		return bad("the stored request is not the one this finish made")
	}
	in := authority.VerbJudgeInput{Digest: own.digest, Challenge: own.challenge, ExpiresAt: own.expires, Now: now}
	signed, err := authority.RecheckVerb(ctx, ro, in, row.Artifact)
	if err != nil {
		return bad(err.Error())
	}
	if signed.ApprovedBy() != row.ApprovedBy {
		return bad(fmt.Sprintf("the key is %s's, the row says %s", signed.ApprovedBy(), row.ApprovedBy))
	}
	return row, "", nil
}

// readBackFinishDecision is design §2.7 step 6(b): through worldReader.object (GET /v1/objects), the
// decision the row names must approve, be decided by the row's approver, and its request must be
// this finish scope's Human.Approve for operator:episode-finish. It returns the refusal, or "".
func readBackFinishDecision(r worldReader, ref, approvedBy, scopeText string) string {
	bad := func(why string) string {
		return "the world decision " + ref + " does not approve this finish (F6): " + why
	}
	obj, ok, err := r.object(ref)
	if err != nil || !ok || obj.SemanticID != broker.ApprovalDecisionV1 {
		return bad(fmt.Sprintf("not a finish approval decision (found %v, %v)", ok, err))
	}
	var d struct {
		RequestRef string `json:"requestRef"`
		Decision   string `json:"decision"`
		DecidedBy  string `json:"decidedBy"`
	}
	if err := json.Unmarshal(obj.Payload, &d); err != nil {
		return bad(err.Error())
	}
	if d.Decision != "approve" || d.DecidedBy != approvedBy {
		return bad(fmt.Sprintf("decision %q by %q", d.Decision, d.DecidedBy))
	}
	qobj, ok, err := r.object(d.RequestRef)
	if err != nil || !ok || qobj.SemanticID != broker.ApprovalRequestV1 {
		return bad(fmt.Sprintf("its request %s is not readable (found %v, %v)", d.RequestRef, ok, err))
	}
	q, err := broker.DecodeApprovalRequest(qobj.Payload)
	if err != nil {
		return bad(err.Error())
	}
	if q.Effect != broker.EffectHumanApprove || q.Requester != "operator:episode-finish" || q.Scope != scopeText {
		return bad(fmt.Sprintf("it approves %s for %s over another scope", q.Effect, q.Requester))
	}
	return ""
}
