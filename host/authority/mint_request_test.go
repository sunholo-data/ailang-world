package authority

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/store"
)

// au1Request is a valid request for episode "ep" (one grant, 600 s, profile proj) and its stored
// row, created at 1000 and expiring at 1600.
func au1Request(t *testing.T, st *store.Credentials) (approval.MintRequest, store.MintRequestRow) {
	t.Helper()
	mr := approval.MintRequest{
		Schema: approval.RequestSchema, Episode: "ep",
		WorldDB: "/abs/world.db", WorkspaceRoot: "/abs/ws", Worktree: "/abs/ws/ep",
		Repo: "/abs/repo", Git: "/usr/bin/git", Branch: "world/ep",
		Grants:      []approval.Grant{{Effect: "Workspace.Read", Scope: "worktree", Budget: 5}},
		TTLSeconds:  600,
		Profile:     approval.ProfileRef{Project: "proj", Digest: "sha256:" + strings.Repeat("c", 64)},
		ProfilesDir: "/abs/profiles", CredentialsDB: "/abs/credentials.db", SessionFile: "/abs/sessions/ep.session",
		Agent: "codex", Model: "gpt-test",
		Daemon: approval.DaemonRef{Addr: "127.0.0.1:1"},
	}
	canon, err := approval.Canonical(mr)
	if err != nil {
		t.Fatal(err)
	}
	sum := sha256.Sum256(canon)
	row := store.MintRequestRow{
		RequestID: "mr-0123456789abcdef", Request: canon, Digest: hex.EncodeToString(sum[:]),
		Challenge: strings.Repeat("ab", 32), EpisodeID: "ep", CreatedAt: 1000, ExpiresAt: 1600,
	}
	if err := st.PutMintRequest(boundedTestContext(t), row, 1000); err != nil {
		t.Fatal(err)
	}
	return mr, row
}

// au1Challenge is au1Request's row challenge, strings.Repeat("ab", 32), as bytes.
var au1Challenge = [32]byte{0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab,
	0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab, 0xab}

// au1Grants is the grant list provisionEpisode builds: ExpiresAt is now + ttl.
func au1Grants(budget, expiresAt int64) []broker.Capability {
	return []broker.Capability{{Effect: "Workspace.Read", Scope: "worktree", Budget: budget, ExpiresAt: expiresAt}}
}

// TestMintForRequestChecksTheRowInsideTheTransaction (AU1): MintForRequest checks the row it is
// about to write, inside the consume's transaction, and writes nothing when a column differs.
func TestMintForRequestChecksTheRowInsideTheTransaction(t *testing.T) {
	good := ProfileBinding{Project: "proj", Digest: "sha256:" + strings.Repeat("c", 64)}
	consumeAt := func(st *store.Credentials, mr approval.MintRequest, row store.MintRequestRow, at int64) store.RequestConsume {
		return store.RequestConsume{ID: row.RequestID, Digest: row.Digest, Clock: func() int64 { return at }}
	}
	t.Run("good", func(t *testing.T) {
		st, err := store.OpenCredentials(filepath.Join(t.TempDir(), store.CredentialsFileName))
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = st.Close() }()
		mr, row := au1Request(t, st)
		ctx := boundedTestContext(t)
		raw, credID, sess, err := MintForRequest(ctx, st, "ep", good, au1Grants(5, 1610), 600, 1010, nil,
			RequestMint{Request: mr, RequestCreatedAt: 1000, Challenge: au1Challenge, Consume: consumeAt(st, mr, row, 1020)})
		if err != nil {
			t.Fatal(err)
		}
		if credID != hashHexToken(raw) || sess.CredentialID != credID {
			t.Errorf("credential id %q, token hash %q", credID, hashHexToken(raw))
		}
		if sess.CreatedAt != 1010 || sess.ExpiresAt != 1610 {
			t.Errorf("session created %d expires %d, want 1010 and 1610", sess.CreatedAt, sess.ExpiresAt)
		}
		got, ok, err := st.MintRequest(ctx, row.RequestID)
		if err != nil || !ok || !got.Consumed || got.ApprovedBy != "tty" || got.CredentialID != credID {
			t.Errorf("request after the mint: %+v ok %v err %v", got, ok, err)
		}
	})
	mismatches := []struct {
		name    string
		episode string
		binding ProfileBinding
		grants  []broker.Capability
		ttl     int64
		now     int64
		clockAt int64
		column  string
	}{
		{"budget", "ep", good, au1Grants(6, 1610), 600, 1010, 1020, "grants_json"},
		{"episode", "other", good, au1Grants(5, 1610), 600, 1010, 1020, "episode_id"},
		{"project", "ep", ProfileBinding{}, au1Grants(5, 1610), 600, 1010, 1020, "project"},
		{"ttl", "ep", good, au1Grants(5, 1611), 601, 1010, 1020, "expires_at"},
		{"created-after-consume", "ep", good, au1Grants(5, 1621), 600, 1021, 1020, "created_at"},
		{"created-before-request", "ep", good, au1Grants(5, 1599), 600, 999, 1020, "created_at"},
	}
	for _, m := range mismatches {
		t.Run(m.name, func(t *testing.T) {
			st, err := store.OpenCredentials(filepath.Join(t.TempDir(), store.CredentialsFileName))
			if err != nil {
				t.Fatal(err)
			}
			defer func() { _ = st.Close() }()
			mr, row := au1Request(t, st)
			ctx := boundedTestContext(t)
			_, _, _, err = MintForRequest(ctx, st, m.episode, m.binding, m.grants, m.ttl, m.now, nil,
				RequestMint{Request: mr, RequestCreatedAt: 1000, Challenge: au1Challenge, Consume: consumeAt(st, mr, row, m.clockAt)})
			var mm *approval.MintMismatchError
			if !errors.As(err, &mm) || mm.Column != m.column {
				t.Fatalf("err %v, want a MintMismatchError on %s", err, m.column)
			}
			got, _, gerr := st.MintRequest(ctx, row.RequestID)
			if gerr != nil || got.Consumed {
				t.Errorf("request consumed after a refused mint: %+v err %v", got, gerr)
			}
			sessions, err := st.ListSessions(ctx, "", store.MaxListSessions)
			if err != nil || len(sessions) != 0 {
				t.Errorf("%d session rows after a refused mint (err %v)", len(sessions), err)
			}
		})
	}
	t.Run("expired-at-the-consume", func(t *testing.T) {
		st, err := store.OpenCredentials(filepath.Join(t.TempDir(), store.CredentialsFileName))
		if err != nil {
			t.Fatal(err)
		}
		defer func() { _ = st.Close() }()
		mr, row := au1Request(t, st)
		_, _, _, err = MintForRequest(boundedTestContext(t), st, "ep", good, au1Grants(5, 1610), 600, 1010, nil,
			RequestMint{Request: mr, RequestCreatedAt: 1000, Challenge: au1Challenge, Consume: consumeAt(st, mr, row, 1600)})
		if !errors.Is(err, store.ErrRequestExpired) {
			t.Fatalf("err %v, want ErrRequestExpired", err)
		}
	})
}
