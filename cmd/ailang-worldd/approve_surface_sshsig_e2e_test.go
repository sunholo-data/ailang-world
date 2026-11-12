package main

// Row 201 A3 sprint 3 M5d (design §2.4, §6 AW6's route half; the sprint-2 plan's Appendix C): an
// SSHSIG POSTed to the surface is recorded, and approve --await mints it, as E2E1 does a passkey's.

import (
	"bytes"
	"context"
	"encoding/hex"
	"encoding/json"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/daemon"
)

// TestSurfaceSSHSIGThenAwaitMints: episode start --request; the card GET; the 90-byte message
// from the card's digest, nonce and expiry; sigSigner's SSHSIG POSTed (200, recorded); approve
// --await mints it with no terminal; requests.json lists it under decided, by ssh:mark-mac.
func TestSurfaceSSHSIGThenAwaitMints(t *testing.T) {
	now := time.Now().Unix() // the daemon's routes read the wall clock
	a := newApproveRig(t, "ep1", now)
	seedApprover(t, a.credPath, "mark-mac", sigSigner.PublicKey(), now-100)
	d, err := daemon.New(boundedTestContext(t), daemon.Config{DBPath: filepath.Join(t.TempDir(), "world.db"), CredentialsPath: a.credPath,
		BindHost: daemon.DefaultBindHost, ErrorLog: io.Discard})
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = d.Close() }()
	serve := func(method, target string, body []byte, out any) int {
		req := httptest.NewRequest(method, target, bytes.NewReader(body))
		req = req.WithContext(context.WithValue(req.Context(), http.LocalAddrContextKey, net.Addr(&net.TCPAddr{IP: net.IPv4(127, 0, 0, 1), Port: 7644})))
		req.Host = "localhost:7644"
		for k, v := range map[string]string{"Origin": "http://localhost:7644", "Sec-Fetch-Site": "same-origin", "Content-Type": "application/json"} {
			req.Header.Set(k, v)
		}
		rec := httptest.NewRecorder()
		d.Handler().ServeHTTP(rec, req)
		if err := json.Unmarshal(rec.Body.Bytes(), out); err != nil {
			t.Fatalf("%s %s: %d %s", method, target, rec.Code, rec.Body.String())
		}
		return rec.Code
	}
	var card struct {
		Digest, Challenge string
		Expires           int64
	}
	if code := serve("GET", "/surface/requests/"+a.id, nil, &card); code != 200 || card.Digest != a.row(t).Digest {
		t.Fatalf("card = %d %+v", code, card)
	}
	dg, _ := hex.DecodeString(card.Digest)
	nonce, _ := hex.DecodeString(card.Challenge)
	msg, err := approval.SignedMessage([32]byte(dg), [32]byte(nonce), card.Expires)
	if err != nil {
		t.Fatal(err)
	}
	body, _ := json.Marshal(map[string]any{"v": 1, "digest": card.Digest, "sshsig": string(sigSigner.Sign(approval.Namespace, msg))})
	var recorded struct {
		Recorded   bool
		ApprovedBy string
	}
	if code := serve("POST", "/surface/approve/"+a.id, body, &recorded); code != 200 || !recorded.Recorded || recorded.ApprovedBy != "ssh:mark-mac" {
		t.Fatalf("POST = %d %+v", code, recorded)
	}
	var opened bool
	env, _ := awaitEnv(time.Now().Unix(), &opened, nil)
	var stdout, stderr bytes.Buffer
	if code := runApprove(a.args("--await", "--await-for", "5s"), &stdout, &stderr, env); code != exitOK || opened {
		t.Fatalf("approve --await: rc %d, terminal opened %v\nstderr %s", code, opened, stderr.String())
	}
	requireMintedBy(t, a, "ssh:mark-mac", stdout.String())
	var list struct {
		Decided []struct{ ID, ApprovedBy string }
	}
	if code := serve("GET", "/surface/requests.json", nil, &list); code != 200 || len(list.Decided) != 1 || list.Decided[0].ID != a.id || list.Decided[0].ApprovedBy != "ssh:mark-mac" {
		t.Fatalf("requests.json = %d %+v", code, list)
	}
}
