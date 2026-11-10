package main

// Row 201 A0 (design §2.6): `ailang-worldd approve <request-id>` mints the session a stored mint
// request describes after the typed y, through provisionEpisode; `approve --list` shows the
// pending requests. The request is consumed in the mint's own transaction.

import (
	"bytes"
	"context"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"flag"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/store"
)

const approveHelp = `usage: ailang-worldd approve <request-id> [--db <path>] [--credentials-db <path>]
       ailang-worldd approve <request-id> --signature <file> [--db <path>] [--credentials-db <path>]
       ailang-worldd approve <request-id> --write-message <file> [--db <path>] [--credentials-db <path>]
       ailang-worldd approve --list [--decided] [--db <path>] [--credentials-db <path>]

ATTENDED, or SIGNED. Mints the session a stored mint request describes (row
201), after a y/N typed at /dev/tty, or, with --signature, after an enrolled
approver's SSH signature over the request verifies, with no terminal.
'ailang-worldd episode start <episode> ... --request' makes the request without
a terminal; nothing is minted until approve.

approve <request-id> (mr- and 16 hex):
  1. reads the request and refuses, before any terminal opens, one that is
     unknown, already used, expired, or whose stored bytes do not match their
     digest;
  2. re-derives the request from the world as it is now (the exec profile, the
     serving daemon's build, the worktree, the session and MCP files) and
     refuses, naming the first field, if anything differs;
  3. runs episode start's own path: the /dev/tty fence, a y/N showing the whole
     request, then the worktree and the mint. The y/N waits at most until the
     request expires. The request is consumed in the mint's own transaction, so
     it mints at most once;
  4. prints the launch line, as episode start does. The token is never printed.
--signature <file> replaces step 3's y/N. The file (at most 16 KiB) is an
  SSHSIG over the request's exact message, in namespace world-approval@sunholo,
  by a key enrolled with 'ailang-worldd approver enroll' and not revoked. It is
  checked before any network call, the approver is re-checked in the mint's own
  transaction, and the request records approved_by=ssh:<name>. A signature
  proves that something holding the key signed, not that you were there.
--write-message <file> writes that message (90 bytes) to a NEW file and prints
  the ssh-keygen -Y sign line and the --signature step. It opens no terminal
  and writes no store row: the request stays pending.
--list prints the pending requests, read-only. A request expires 10 minutes
after it is made; at most 16 are pending per credentials store.
--list --decided prints the used requests instead, newest first (at most 64):
  who approved each (tty, or ssh:<name>), when, and its credential's first 16
  hex. 'ailang-worldd session revoke' takes the credential id.

  --db <path>               the world store the request was made for
                            (default ~/.ailang/world/world.db)
  --credentials-db <path>   default credentials.db beside --db
The global --addr is refused: the daemon's address is part of the request.
`

// approveAdmission is approve's admission of a stored row, before any terminal: the decoded
// request, or the refusal text ("" on admission).
func approveAdmission(row store.MintRequestRow, now int64) (approval.MintRequest, string) {
	mr, err := approval.Decode(row.Request)
	if err != nil {
		return approval.MintRequest{}, fmt.Sprintf("request %s is not a canonical mint request (%v); fix: make a new request with episode start --request", row.RequestID, err)
	}
	digestRefusal := fmt.Sprintf("request %s's stored bytes do not match its digest; refusing; fix: make a new request with episode start --request", row.RequestID)
	stored, err := hex.DecodeString(row.Digest)
	if err != nil || len(stored) != 32 {
		return approval.MintRequest{}, digestRefusal
	}
	verdict := approval.Admit(approval.AdmitInput{Now: now, CreatedAt: row.CreatedAt, ExpiresAt: row.ExpiresAt,
		Consumed: row.Consumed, StoredDigest: [32]byte(stored), BytesDigest: sha256.Sum256(row.Request)})
	switch verdict {
	case approval.Admitted:
		return mr, ""
	case approval.RefuseConsumed:
		return approval.MintRequest{}, fmt.Sprintf("request %s was already used at %s by %s (credential %s); fix: make a new request with episode start --request",
			row.RequestID, time.Unix(row.ConsumedAt, 0).UTC().Format(time.RFC3339), row.ApprovedBy, row.CredentialID)
	case approval.RefuseClockBehind:
		return approval.MintRequest{}, fmt.Sprintf("the clock is behind request %s's creation time (%s); refusing",
			row.RequestID, time.Unix(row.CreatedAt, 0).UTC().Format(time.RFC3339))
	case approval.RefuseExpired:
		return approval.MintRequest{}, fmt.Sprintf("request %s expired at %s; fix: make a new request with episode start --request",
			row.RequestID, time.Unix(row.ExpiresAt, 0).UTC().Format(time.RFC3339))
	case approval.RefuseDigestMismatch:
		return approval.MintRequest{}, digestRefusal
	default:
		return approval.MintRequest{}, fmt.Sprintf("request %s: unexpected admission verdict %s", row.RequestID, verdict)
	}
}

// runApprove is `ailang-worldd approve <request-id>` and `approve --list` (row 201 §2.6). Every
// refusal before the mint names its fix and exits exitUsage; a failure at or after the write is
// provisionEpisode's exitFatal.
func runApprove(args []string, stdout, stderr io.Writer, env sessionEnv) int {
	const prefix = "ailang-worldd approve: "
	refuse := func(format string, a ...any) int {
		fmt.Fprintf(stderr, "%s"+format+"\n", append([]any{prefix}, a...)...)
		return exitUsage
	}
	def := defaultWorldPaths()
	fs := flag.NewFlagSet("ailang-worldd approve", flag.ContinueOnError)
	fs.SetOutput(stderr)
	fs.Usage = func() {}
	db := fs.String("db", def.db, "the world store the request was made for (default ~/.ailang/world/world.db)")
	credsDB := fs.String("credentials-db", "", credentialsDBUsage)
	list := fs.Bool("list", false, "print the pending requests, read-only")
	decided := fs.Bool("decided", false, "with --list: print the used requests and who approved them")
	sigPath := fs.String("signature", "", "approve with an enrolled approver's SSHSIG over the request (no terminal)")
	msgPath := fs.String("write-message", "", "write the request's approval message to this new file, to sign elsewhere")
	positional, err := parseInterleaved(fs, args)
	if errors.Is(err, flag.ErrHelp) {
		fmt.Fprint(stdout, approveHelp)
		return exitOK
	}
	if err != nil {
		fmt.Fprint(stderr, approveHelp)
		return exitUsage
	}
	if *list && len(positional) != 0 || !*list && len(positional) != 1 {
		fmt.Fprint(stderr, approveHelp)
		return exitUsage
	}
	if *sigPath != "" && *msgPath != "" { // F1: one of --signature and --write-message
		fmt.Fprint(stderr, approveHelp)
		return exitUsage
	}
	if *list && (*sigPath != "" || *msgPath != "") { // F2: --list is read-only
		fmt.Fprint(stderr, approveHelp)
		return exitUsage
	}
	if *decided && !*list { // F3: --decided is a --list
		fmt.Fprint(stderr, approveHelp)
		return exitUsage
	}
	id := ""
	if !*list {
		id = positional[0]
		if !store.MintRequestIDPattern.MatchString(id) {
			return refuse("%q is not a mint request id (mr- and 16 lowercase hex)", id)
		}
		if env.after == nil {
			fmt.Fprintf(stderr, "%sinternal error: no prompt timer (sessionEnv.after)\n", prefix)
			return exitFatal
		}
	}
	credPath, err := credentialsFor(*db, *credsDB)
	if err != nil {
		return refuse("%v", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), approveReadBudget)
	defer cancel()
	ro, err := store.OpenCredentialsReadOnly(credPath)
	if err != nil {
		return refuse("cannot read the credentials store %s: %v", credPath, err)
	}
	if *list && *decided {
		return approveListDecided(ctx, ro, credPath, stdout, refuse)
	}
	if *list {
		rows, err := ro.PendingMintRequests(ctx, env.now(), 0)
		_ = ro.Close()
		if err != nil {
			return refuse("cannot list the requests in %s: %v", credPath, err)
		}
		if len(rows) == 0 {
			fmt.Fprintln(stdout, "no pending mint requests")
			return exitOK
		}
		for _, r := range rows {
			fmt.Fprintf(stdout, "%s  episode %s  digest %s  expires %s\n", r.RequestID, r.EpisodeID, r.Digest[:16],
				time.Unix(r.ExpiresAt, 0).UTC().Format(time.RFC3339))
		}
		return exitOK
	}
	row, ok, err := ro.MintRequest(ctx, id)
	_ = ro.Close()
	if err != nil {
		return refuse("cannot read request %s from %s: %v", id, credPath, err)
	}
	if !ok {
		return refuse("no mint request %s in %s; fix: ailang-worldd approve --list shows the pending ones, or make one with episode start --request", id, credPath)
	}

	// 1. Admission, before any terminal (design §2.3).
	now := env.now()
	mr, refusal := approveAdmission(row, now)
	if refusal != "" {
		return refuse("%s", refusal)
	}
	challenge, err := hex.DecodeString(row.Challenge)
	if err != nil || len(challenge) != 32 { // F4: the consume binds the row's challenge (row 201 A1b §2.7)
		return refuse("request %s's stored challenge is not 64 hex; refusing; fix: make a new request with episode start --request", row.RequestID)
	}
	if *msgPath != "" { // F5, F6: the message to sign elsewhere, and nothing else
		if refusal := writeApprovalMessage(stdout, *msgPath, row, mr); refusal != "" {
			return refuse("%s", refusal)
		}
		return exitOK
	}
	var signed *approval.SignedApproval
	if *sigPath != "" { // the signature, judged before any network call (design §2.7 F9-F14)
		var refusal string
		if signed, refusal = approveSignature(ctx, credPath, *sigPath, row, now); refusal != "" {
			return refuse("%s", refusal)
		}
	}

	// 2. The same store: the request was made for this world store and credentials file.
	w, werr := store.CanonicalPath(*db)
	if werr != nil || w != mr.WorldDB || credPath != mr.CredentialsDB {
		shown := *db
		if werr == nil {
			shown = w
		}
		return refuse("--db/--credentials-db name %s and %s, but request %s was made for %s and %s; fix: pass --db %s --credentials-db %s",
			shown, credPath, id, mr.WorldDB, mr.CredentialsDB, mr.WorldDB, mr.CredentialsDB)
	}

	// 3. The world as it is now: the profile, then the MCP config file, which must not exist yet.
	req, err := prepareMint(fromRequest(mr))
	if err != nil {
		return refuse("the world changed since request %s was made: %v; it stays pending", id, err)
	}
	if mr.MCPConfig != "" {
		if _, err := os.Lstat(mr.MCPConfig); err == nil {
			return refuse("the world changed since request %s was made: %s already exists; it stays pending", id, mr.MCPConfig)
		}
	}

	// 4. The serving daemon is this build, serves this build's pin and reads this store.
	health, err := fetchDaemonHealth(mr.Daemon.Addr)
	if err != nil {
		return refuse("refusing: %v", err)
	}
	if err := checkDaemon(health, req, mr.Daemon.Addr, episodeSelfBuild()); err != nil {
		return refuse("refusing: %v", err)
	}

	// 5. The request, re-derived from the world now, is the stored bytes exactly.
	fresh, err := buildMintRequest(req, health.Build, mr.Daemon.Addr, mr.Agent, mr.Model, mr.Task, mr.MCPConfig)
	if err != nil {
		return refuse("the world changed since request %s was made: %v; it stays pending", id, err)
	}
	fb, _ := approval.Canonical(fresh)
	if !bytes.Equal(fb, row.Request) {
		return refuse("the world changed since request %s was made: %s differs; it stays pending", id, approval.Diff(mr, fresh))
	}

	// 6. The mint: episode start's own path, bounded by the request's expiry. With a verified
	// signature, provisionEpisode skips the terminal (design §2.8).
	req.presetPrompt = renderRequest(row.RequestID, mr, row.Digest, row.ExpiresAt, now)
	req.afterToken = mcpAfterToken(mr.MCPConfig, mr.Daemon.Addr)
	req.answerBy = row.ExpiresAt
	req.consume = &authority.RequestMint{Request: mr, RequestCreatedAt: row.CreatedAt, Challenge: [32]byte(challenge),
		ExpiresAt: row.ExpiresAt, Signed: signed, Consume: store.RequestConsume{ID: row.RequestID, Digest: row.Digest, Clock: env.now}}
	if code := provisionEpisode(req, stdout, stderr, env); code != exitOK {
		return code
	}
	printLaunch(stdout, launchInfo{episode: mr.Episode, agent: mr.Agent, model: mr.Model, addr: mr.Daemon.Addr,
		worktree: req.dir, branch: mr.Branch, reused: req.reuse, sessionFile: mr.SessionFile, mcpFile: mr.MCPConfig,
		task: mr.Task, grants: toCapabilities(mr.Grants)})
	return exitOK
}

// approveListDecided is approve --list --decided (design §2.7 F15, F16): the used requests,
// newest first, each with who approved it, when, and its credential id's first 16 hex. It closes ro.
func approveListDecided(ctx context.Context, ro *store.Credentials, credPath string, stdout io.Writer, refuse func(string, ...any) int) int {
	rows, err := ro.DecidedMintRequests(ctx, "", 0)
	_ = ro.Close()
	if err != nil {
		return refuse("cannot list the used requests in %s: %v", credPath, err)
	}
	if len(rows) == 0 { // F15
		fmt.Fprintln(stdout, "no decided mint requests")
		return exitOK
	}
	for _, r := range rows { // F16
		fmt.Fprintf(stdout, "%s  episode %s  approved by %s  at %s  credential %s\n", r.RequestID, r.EpisodeID, r.ApprovedBy,
			approverTime(r.ConsumedAt), credentialPrefix(r.CredentialID))
	}
	return exitOK
}

// credentialPrefix is a credential id's first 16 characters (all of a shorter one).
func credentialPrefix(id string) string { return id[:min(16, len(id))] }
