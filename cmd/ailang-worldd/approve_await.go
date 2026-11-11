package main

// Row 201 A3 sprint 1 M3 (design_docs/planned/w-world-approvals-by-key-a3.md §2.6): `approve
// <request-id> --await [--await-for <duration>]`, the minter. It opens no terminal: it waits for an
// approval the daemon's surface recorded (credentials schema 7, surface_approvals), judges that
// artifact again with the same law (authority.JudgeAssertion, JudgeSSHSignature) against the
// request it reads in the same poll, and only a verified one goes on to approve's steps 2-6. The
// row's approver and digest columns are never trusted.

import (
	"context"
	"encoding/hex"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/store"
)

// awaitPollInterval is the minter's pause between polls (M1), through env.after.
const awaitPollInterval = time.Second

// writeAwaitLinks prints the minter's start (design §2.6): the request, its message hash and
// expiry, and the card's address for each distinct active passkey origin and the request's own
// loopback daemon. Nothing it prints is secret. It returns the refusal text, or "".
func writeAwaitLinks(ctx context.Context, credPath string, stdout io.Writer, row store.MintRequestRow, mr approval.MintRequest) string {
	ro, err := store.OpenCredentialsReadOnly(credPath)
	if err != nil {
		return fmt.Sprintf("cannot read the credentials store %s: %v", credPath, err)
	}
	active, err := ro.ActivePasskeys(ctx)
	_ = ro.Close()
	if err != nil { // below schema 6 names the serve restart
		return err.Error()
	}
	h, err := requestHash(row)
	if err != nil {
		return fmt.Sprintf("request %s: %v", row.RequestID, err)
	}
	fmt.Fprintf(stdout, "waiting for an approval of request %s (expires %s, message hash %s) from the surface; open its card and approve it there:\n",
		row.RequestID, approverTime(row.ExpiresAt), hex.EncodeToString(h[:8]))
	seen := map[string]bool{}
	for _, e := range active {
		if !seen[e.Passkey.Origin] {
			seen[e.Passkey.Origin] = true
			fmt.Fprintf(stdout, "  open %s/surface/approve/%s\n", e.Passkey.Origin, row.RequestID)
		}
	}
	if _, port, err := net.SplitHostPort(mr.Daemon.Addr); err == nil && !seen["http://localhost:"+port] {
		fmt.Fprintf(stdout, "  open http://localhost:%s/surface/approve/%s\n", port, row.RequestID)
	}
	return ""
}

// awaitPoll is one poll's reading: the request row, its recorded approval if any, and that
// approval's judgement (a verified approval, or the refusal text).
type awaitPoll struct {
	row      store.MintRequestRow
	rec      store.SurfaceApprovalRow
	recorded bool
	signed   *approval.SignedApproval
	refusal  string
}

// pollSurfaceApproval is one poll (M1, M3): one read-only open under its own 3 s root, the request
// row and its recorded approval, and the approval's artifact judged against that row at now.
func pollSurfaceApproval(credPath, id string, now int64) (awaitPoll, error) {
	ctx, cancel := context.WithTimeout(context.Background(), approveReadBudget)
	defer cancel()
	ro, err := store.OpenCredentialsReadOnly(credPath)
	if err != nil {
		return awaitPoll{}, err
	}
	defer func() { _ = ro.Close() }()
	var p awaitPoll
	var ok bool
	if p.row, ok, err = ro.MintRequest(ctx, id); err != nil {
		return awaitPoll{}, err
	}
	if !ok {
		return awaitPoll{}, fmt.Errorf("request %s is no longer in the store", id)
	}
	if p.rec, p.recorded, err = ro.SurfaceApproval(ctx, id); err != nil || !p.recorded {
		return p, err
	}
	var jerr error
	if p.rec.Format == "ssh" { // M3: the artifact, judged again; never the approver column
		if p.signed, jerr = authority.JudgeSSHSignature(ctx, ro, p.row, p.rec.Artifact, now); jerr != nil {
			p.refusal = sshRefusal(jerr, "the surface's approval", p.row)
		}
	} else if p.signed, jerr = authority.JudgeAssertion(ctx, ro, p.row, p.rec.Artifact, now); jerr != nil {
		p.refusal = passkeyRefusal(jerr, "the surface's approval", p.row)
	}
	return p, nil
}

// awaitSurfaceApproval is the minter's wait (M1-M6), after admission and F4: a poll every second
// until the request has a recorded approval that the law verifies (returned), or the request is
// used elsewhere or expires (admission's own text), or the wait reaches now + awaitFor (0: the
// expiry). A refused approval is printed once per received_at, and the wait goes on.
func awaitSurfaceApproval(env sessionEnv, credPath string, row store.MintRequestRow, awaitFor time.Duration, stderr io.Writer) (*approval.SignedApproval, string) {
	deadline := row.ExpiresAt
	if awaitFor > 0 {
		deadline = min(deadline, env.now()+int64(awaitFor/time.Second))
	}
	shown := int64(-1)
	for {
		now := env.now()
		p, err := pollSurfaceApproval(credPath, row.RequestID, now)
		if err != nil {
			return nil, fmt.Sprintf("cannot read request %s: %v", row.RequestID, err)
		}
		if _, refusal := approveAdmission(p.row, now); refusal != "" { // M2: used elsewhere, expired
			return nil, refusal
		}
		if p.signed != nil { // M5
			return p.signed, ""
		}
		if p.recorded && p.rec.ReceivedAt != shown { // M4
			fmt.Fprintf(stderr, "ailang-worldd approve: refused the surface's approval: %s; waiting for another\n", p.refusal)
			shown = p.rec.ReceivedAt
		}
		if now >= deadline { // M6
			return nil, fmt.Sprintf("no approval for request %s by %s; it stays pending until %s",
				row.RequestID, approverTime(now), approverTime(row.ExpiresAt))
		}
		<-env.after(awaitPollInterval)
	}
}
