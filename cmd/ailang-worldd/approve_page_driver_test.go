package main

// Row 201 A3 sprint 3 M8a (design §2.7, §6 XA1-XA5): the request card's browser harness, on A2's
// driver (testdata/passkey_cdp.mjs, extended): Chrome headless, a temporary profile and a CDP virtual
// authenticator holding a test key (D-WORLD-109). The card is served by the production daemon
// handler on the loopback, over the request's temporary credentials store.

import (
	"bytes"
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
	"github.com/sunholo-data/ailang-world/host/daemon"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

// cardView is what the card page shows (the driver's card command).
type cardView struct {
	State, What, Status, Again, SaveHref, SaveName, SSHLine, Note      string
	StatusError, ButtonHidden, ButtonDisabled, AgainHidden, SaveHidden bool
	Fields                                                             map[string]string
	TTLLabel                                                           string
	Gets                                                               []cardGet
	Violations, Errors                                                 []string
	Op                                                                 *struct {
		Signal, Aborted             bool
		AbortedAfter, ResolvedAfter *float64
	}
}

// cardGet is one navigator.credentials.get call the getRecorder injection saw.
type cardGet struct {
	Challenge, RPID, UserVerification string
	Allow                             []string
	Timeout                           int
	Signal                            bool
}

// cardOrigin serves the production daemon handler (wrapped by wrap, if any) over credentials store
// creds, with surface index index (or none), on the loopback, and returns its origin.
func cardOrigin(t *testing.T, creds, index string, wrap func(http.Handler) http.Handler) string {
	t.Helper()
	d, err := daemon.New(boundedTestContext(t), daemon.Config{DBPath: filepath.Join(t.TempDir(), "world.db"), CredentialsPath: creds,
		SurfaceIndex: index, BindHost: daemon.DefaultBindHost, ErrorLog: io.Discard})
	if err != nil {
		t.Fatal(err)
	}
	t.Cleanup(func() { _ = d.Close() })
	h := d.Handler()
	if wrap != nil {
		h = wrap(h)
	}
	srv := httptest.NewServer(h)
	t.Cleanup(srv.Close)
	_, port, _ := net.SplitHostPort(srv.Listener.Addr().String())
	return "http://localhost:" + port
}

// virtualPasskey enrols approvaltest's ES256 passkey seed as name at origin (its rpID: the origin's
// host) in creds (at `at`) and gives its private key to the tab's virtual authenticator, so the
// browser signs as that test key.
func virtualPasskey(t *testing.T, b *passkeyBrowser, creds, origin, name string, seed byte, at int64) {
	t.Helper()
	pk := approvaltest.NewPasskey(approval.AlgES256, seed)
	k, err := ecdsa.ParseRawPrivateKey(elliptic.P256(), bytes.Repeat([]byte{seed}, 32))
	if err != nil {
		t.Fatal(err)
	}
	if pub, _ := k.PublicKey.Bytes(); !bytes.Equal(pub[1:33], pk.Key().X) { // control: the same key
		t.Fatal("the virtual credential is not approvaltest's key")
	}
	der, err := x509.MarshalPKCS8PrivateKey(k)
	if err != nil {
		t.Fatal(err)
	}
	rp, err := approval.ParseCeremonyOrigin(origin)
	if err != nil {
		t.Fatal(err)
	}
	st, err := store.OpenCredentials(creds)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = st.Close() }()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := st.EnrolPasskey(ctx, store.ApproverRow{Name: name, Format: approval.FormatWebAuthn, PublicKey: pk.Key().Canonical(), EnrolledAt: at},
		store.PasskeyRow{Name: name, RPID: rp, Origin: origin, CredentialID: pk.CredentialID(), AAGUID: make([]byte, 16), Alg: approval.AlgES256, SignCount: 1}); err != nil {
		t.Fatal(err)
	}
	b.do(t, nil, "credential", map[string]any{"id": base64.StdEncoding.EncodeToString(pk.CredentialID()), "rpId": rp,
		"privateKey": base64.StdEncoding.EncodeToString(der)})
}

// cardFields is what the card must show for canonical request bytes req: every leaf, by its
// concrete key path, as text (a boolean as FIELDS' word for it).
func cardFields(t *testing.T, req []byte) map[string]string {
	t.Helper()
	_, rest, _ := bytes.Cut(surfacepage.ApproveScript, []byte("/* FIELDS */"))
	lit, _, _ := bytes.Cut(rest, []byte("/* /FIELDS */"))
	var fields []struct{ Path, Type, No, Yes string }
	if err := json.Unmarshal(lit, &fields); err != nil || len(fields) == 0 {
		t.Fatalf("approve.js's FIELDS: %v", err)
	}
	words := map[string][2]string{}
	for _, f := range fields {
		words[f.Path] = [2]string{f.No, f.Yes}
	}
	dec := json.NewDecoder(bytes.NewReader(req))
	dec.UseNumber()
	var v any
	if err := dec.Decode(&v); err != nil {
		t.Fatal(err)
	}
	out := map[string]string{}
	var walk func(v any, pattern, concrete string)
	walk = func(v any, pattern, concrete string) {
		dot := func(p string) string {
			if p == "" {
				return p
			}
			return p + "."
		}
		switch x := v.(type) {
		case map[string]any:
			for k, e := range x {
				walk(e, dot(pattern)+k, dot(concrete)+k)
			}
		case []any:
			for i, e := range x {
				walk(e, pattern+"[]", fmt.Sprintf("%s.%d", concrete, i))
			}
		case bool:
			w := words[pattern][0]
			if x {
				w = words[pattern][1]
			}
			out[concrete] = w
		default:
			out[concrete] = fmt.Sprint(x)
		}
	}
	walk(v, "", "")
	return out
}

// surfaceRecorded is the surface approval recorded for request id in creds, if any.
func surfaceRecorded(t *testing.T, creds, id string) (store.SurfaceApprovalRow, bool) {
	t.Helper()
	ro, err := store.OpenCredentialsReadOnly(creds)
	if err != nil {
		t.Fatal(err)
	}
	defer func() { _ = ro.Close() }()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	row, ok, err := ro.SurfaceApproval(ctx, id)
	if err != nil {
		t.Fatal(err)
	}
	return row, ok
}
