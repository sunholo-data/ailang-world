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
       ailang-worldd approve <request-id> --sign-with <key> [--db <path>] [--credentials-db <path>]
       ailang-worldd approve <request-id> --write-message <file> [--db <path>] [--credentials-db <path>]
       ailang-worldd approve <request-id> --webauthn-link [--db <path>] [--credentials-db <path>]
       ailang-worldd approve <request-id> --webauthn-assertion <file> [--db <path>] [--credentials-db <path>]
       ailang-worldd approve <request-id> --await [--await-for <duration>] [--db <path>] [--credentials-db <path>]
       ailang-worldd approve --list [--decided] [--db <path>] [--credentials-db <path>]
       ailang-worldd approve <request-ref> [--deny] [--daemon <host:port>] [--db <path>] [--credentials-db <path>]
       ailang-worldd approve <request-ref> [--deny] --sign-with <key> [--daemon <host:port>] [--db <path>] [--credentials-db <path>]
       ailang-worldd approve <request-ref> [--deny] --write-message <file> | --webauthn-link [--daemon <host:port>] [--db <path>] [--credentials-db <path>]
       ailang-worldd approve <request-ref> [--deny] --signature <file> | --webauthn-assertion <file> --expiry <t> [--daemon <host:port>] [...]

ATTENDED, or SIGNED. Mints the session a stored mint request describes (row
201), after a y/N typed at /dev/tty, or, with --signature or --sign-with,
after an enrolled approver's SSH signature over the request verifies, or, with
--webauthn-assertion, after an enrolled passkey's assertion verifies, with no
terminal.
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
--sign-with <key> signs that message by running ssh-keygen -Y sign -f <key>
  (a private key file, or a .pub whose key ssh-agent holds; ssh-keygen may ask
  for a passphrase itself) and checks the signature as --signature does. It
  waits at most until the request expires. World never reads the private key.
  A .pub path does not prove the agent signed: ssh-keygen uses the private
  file beside it when there is one.
--write-message <file> writes that message (90 bytes) to a NEW file and prints
  the ssh-keygen -Y sign line and the --signature step. It opens no terminal
  and writes no store row: the request stays pending.
--webauthn-link prints, for each origin of the enrolled passkeys, a link to
  the daemon's page /surface/passkey that asks the browser's passkey (with
  Face ID, Touch ID, a PIN or a security key) to sign the request's message
  hash, and the --webauthn-assertion step. Check the page shows the request,
  expiry and hash it prints. It opens no terminal and writes no store row.
--webauthn-assertion <file> replaces step 3's y/N with the page's saved file
  (at most 16 KiB): the assertion is checked against the passkey enrolled for
  its credential before any network call, its sign counter must move forward,
  and the request records approved_by=webauthn:<name>.
--await replaces step 3's y/N with an approval made on the daemon's surface
  (the request's card, /surface/approve/<id>): it prints the card's address,
  then waits, with no terminal, until the surface has recorded a passkey
  assertion or SSH signature for the request. That artifact is judged again
  here, as --webauthn-assertion and --signature judge theirs; a refused one is
  shown and the wait goes on. It ends when the request is approved, used,
  expired, or --await-for (at least 1s) has passed; the request stays pending.
--list prints the pending requests, read-only. A request expires 10 minutes
after it is made; at most 16 are pending per credentials store.
--list --decided prints the used requests instead, newest first (at most 64):
  who approved each (tty, or ssh:<name>), when, and its credential's first 16
  hex. 'ailang-worldd session revoke' takes the credential id.

  --db <path>               the world store the request was made for
                            (default ~/.ailang/world/world.db)
  --credentials-db <path>   default credentials.db beside --db
The global --addr is refused: the daemon's address is part of the request.

approve <request-ref> (sha256: and 64 hex, row 201 A4) decides a pending
Human.Approve request an agent raised in an episode, whose decision packet the
daemon committed: approve, or with --deny deny. It reads the request and its
packet from the daemon serving --db (--daemon <host:port>, default the
installed service's address), checks that both hash to their refs, shows every
field, asks y at /dev/tty, and has the daemon record the decision (decided by
operator:terminal). A request is decided once. --await is refused (the daemon
decides at once); publish and finish approvals stay with world-publish
approve and episode finish. --sign-with <key> opens no terminal: ssh-keygen
-Y sign -n world-decision@sunholo signs the 91-byte decision message, and the
decision is by ssh:<name>, an approver enrolled with 'approver enroll'.
--write-message <file> writes that message (expiry: the latest the daemon
allows) to a NEW file and prints the signing line and the --signature <file>
--expiry <t> step that posts it. --webauthn-link prints, per passkey origin, a
link to the passkey page over the message's hash (the page says "Approve
request decide-<verb>-…" for either verb: check the hash it shows), and the
--webauthn-assertion <file> --expiry <t> step that posts the saved file; that
decision is by webauthn:<name>. The daemon judges every artifact.
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
		return approval.MintRequest{}, expiredRefusal(row)
	case approval.RefuseDigestMismatch:
		return approval.MintRequest{}, digestRefusal
	default:
		return approval.MintRequest{}, fmt.Sprintf("request %s: unexpected admission verdict %s", row.RequestID, verdict)
	}
}

// expiredRefusal is the refusal of a request that expired before it was approved.
func expiredRefusal(row store.MintRequestRow) string {
	return fmt.Sprintf("request %s expired at %s; fix: make a new request with episode start --request",
		row.RequestID, time.Unix(row.ExpiresAt, 0).UTC().Format(time.RFC3339))
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
	signWith := fs.String("sign-with", "", "approve by signing the request with ssh-keygen -Y sign -f <key> (no terminal)")
	passkeyLink := fs.Bool("webauthn-link", false, "print the ceremony page's links for the request (no terminal, no write)")
	assertion := fs.String("webauthn-assertion", "", "approve with an enrolled passkey's assertion file (no terminal)")
	await := fs.Bool("await", false, "wait for an approval recorded by the daemon's surface, judge it, and mint (no terminal)")
	awaitFor := fs.Duration("await-for", 0, "with --await: stop waiting after this long (default: the request's expiry)")
	deny := fs.Bool("deny", false, "with a request ref: deny it instead of approving it")
	daemonAddr := fs.String("daemon", "", "with a request ref: the daemon serving --db (default: the installed service's address)")
	expiry := fs.Int64("expiry", 0, "with a request ref and --signature or --webauthn-assertion: the expiry the signed decision message carries")
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
	modes := 0
	for _, on := range []bool{*sigPath != "", *signWith != "", *msgPath != "", *passkeyLink, *assertion != "", *await} {
		if on {
			modes++
		}
	}
	if modes > 1 { // F1: at most one of --signature, --sign-with, --write-message, --webauthn-link, --webauthn-assertion and --await
		fmt.Fprint(stderr, approveHelp)
		return exitUsage
	}
	if *list && modes != 0 { // F2: --list is read-only
		fmt.Fprint(stderr, approveHelp)
		return exitUsage
	}
	if *decided && !*list { // F3: --decided is a --list
		fmt.Fprint(stderr, approveHelp)
		return exitUsage
	}
	awaitForSet := false
	fs.Visit(func(f *flag.Flag) { awaitForSet = awaitForSet || f.Name == "await-for" })
	if awaitForSet && (!*await || *awaitFor < time.Second) { // AW5: --await-for is an --await, of at least 1s
		fmt.Fprint(stderr, approveHelp)
		return exitUsage
	}
	if (*deny || *daemonAddr != "" || *expiry != 0) && (*list || !decisionRefPattern.MatchString(positional[0])) { // CD3: they decide a request ref
		return refuse("--deny, --daemon and --expiry take a Human.Approve request ref (sha256: and 64 hex); a mint request has no deny: it expires")
	}
	id := ""
	if !*list {
		id = positional[0]
		if decisionRefPattern.MatchString(id) { // row 201 A4: a request ref selects the decision mode (design §2.10)
			if *await { // CD3
				return refuse("a request ref is decided at once by the daemon: there is nothing to await")
			}
			if (*sigPath != "" || *assertion != "") != (*expiry != 0) { // CD7, CD9: an artifact is posted with the expiry it signed
				return refuse("--signature and --webauthn-assertion go with --expiry <t> on a request ref, and --expiry with one of them: the expiry is part of the signed message")
			}
			verb := "approve"
			if *deny {
				verb = "deny"
			}
			return runApproveDecide(decideOptions{ref: id, verb: verb, db: *db, credsDB: *credsDB, daemon: *daemonAddr, signWith: *signWith,
				sigPath: *sigPath, msgPath: *msgPath, assertion: *assertion, passkeyLink: *passkeyLink, expiry: *expiry}, stdout, stderr, env, refuse)
		}
		if !store.MintRequestIDPattern.MatchString(id) {
			return refuse("%q is not a mint request id (mr- and 16 lowercase hex)", id)
		}
		if env.after == nil {
			fmt.Fprintf(stderr, "%sinternal error: no prompt timer (sessionEnv.after)\n", prefix)
			return exitFatal
		}
		if *signWith != "" && env.sign == nil { // F7
			fmt.Fprintf(stderr, "%sinternal error: no signer (sessionEnv.sign)\n", prefix)
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
		}
		for _, r := range rows {
			fmt.Fprintf(stdout, "%s  episode %s  digest %s  expires %s\n", r.RequestID, r.EpisodeID, r.Digest[:16],
				time.Unix(r.ExpiresAt, 0).UTC().Format(time.RFC3339))
		}
		approveListDecisions(ctx, *db, env.now(), stdout, stderr) // row 201 A4 (design Q2)
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
	if *passkeyLink { // WF1: the ceremony links, and nothing else (row 201 A2 §2.13)
		if refusal := writePasskeyLinks(ctx, credPath, stdout, row, mr); refusal != "" {
			return refuse("%s", refusal)
		}
		return exitOK
	}
	var signed *approval.SignedApproval
	if *await { // the minter: a surface approval, judged again (row 201 A3 §2.6)
		if refusal := writeAwaitLinks(ctx, credPath, stdout, row, mr); refusal != "" {
			return refuse("%s", refusal)
		}
		var refusal string
		if signed, refusal = awaitSurfaceApproval(env, credPath, row, *awaitFor, stderr); refusal != "" {
			return refuse("%s", refusal)
		}
	}
	if *sigPath != "" { // the signature, judged before any network call (design §2.7 F9-F14)
		var refusal string
		if signed, refusal = approveSignature(ctx, credPath, *sigPath, row, now); refusal != "" {
			return refuse("%s", refusal)
		}
	}
	if *assertion != "" { // the passkey's assertion, judged before any network call (design §2.13 WF2-WF6)
		var refusal string
		if signed, refusal = approveAssertion(ctx, credPath, *assertion, row, now); refusal != "" {
			return refuse("%s", refusal)
		}
	}
	if *signWith != "" { // the signer's output, judged as --signature's (design §2.7 step 3, F7a, F8)
		armored, refusal := signForApproval(env, *signWith, row, stderr)
		if refusal != "" {
			return refuse("%s", refusal)
		}
		if signed, refusal = approveArmored(ctx, credPath, "ssh-keygen's output for --sign-with "+*signWith, armored, row, now); refusal != "" {
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
