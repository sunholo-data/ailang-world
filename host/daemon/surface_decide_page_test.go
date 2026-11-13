package daemon

// Row 201 A4 sprint 4 (design §2.8, §6 PD1, PD2): the decision card page, its script and its field
// table.

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/surfacepage"
)

// decisionFields is decide.js's DECISION_FIELDS table.
func decisionFields(t *testing.T) []struct{ Object, Path, Label, Type string } {
	t.Helper()
	_, rest, _ := bytes.Cut(surfacepage.DecideScript, []byte("/* DECISION_FIELDS */"))
	lit, _, _ := bytes.Cut(rest, []byte("/* /DECISION_FIELDS */"))
	var fields []struct{ Object, Path, Label, Type string }
	if err := json.Unmarshal(lit, &fields); err != nil || len(fields) == 0 {
		t.Fatalf("decide.js's DECISION_FIELDS: %v", err)
	}
	return fields
}

// PD1: DECISION_FIELDS names exactly the keys, and their JSON types, of a stored
// world/approval-request/v1 and world/decision-packet/v1 as the subject serves them (the broker's
// and the coordinator's own encoders wrote them): a wire field added without the page fails here.
func TestDecidePageFieldTable(t *testing.T) {
	r := newDecideRig(t)
	q, _ := r.packet("fields", broker.WorkspaceScope, "agent")
	rec := r.serve("GET", "/surface/decisions/"+q.String(), apLoopback, nil)
	var s surfacepage.DecisionSubject
	if err := json.Unmarshal(rec.Body.Bytes(), &s); err != nil || rec.Code != 200 {
		t.Fatalf("subject = %d %s", rec.Code, rec.Body.String())
	}
	wire := map[string]string{}
	for object, b64 := range map[string]string{"request": s.RequestBytes, "packet": s.PacketBytes} {
		b, err := base64.RawURLEncoding.DecodeString(b64)
		if err != nil {
			t.Fatal(err)
		}
		var m map[string]any
		if err := json.Unmarshal(b, &m); err != nil {
			t.Fatal(err)
		}
		for k, v := range m {
			kind := "number"
			if _, ok := v.(string); ok {
				kind = "string"
			}
			wire[object+"."+k] = kind
		}
	}
	table := map[string]string{}
	for _, f := range decisionFields(t) {
		kind := map[string]string{"string": "string", "int": "number", "time": "number"}[f.Type]
		if kind == "" || f.Label == "" {
			t.Errorf("DECISION_FIELDS %s.%s: type %q, label %q", f.Object, f.Path, f.Type, f.Label)
		}
		table[f.Object+"."+f.Path] = kind
	}
	if !reflect.DeepEqual(table, wire) {
		t.Fatalf("DECISION_FIELDS %v, the stored wires %v", table, wire)
	}
	if len(wire) != 12 { // control: request 5 + packet 7 (design §2.8)
		t.Fatalf("the wires have %d fields", len(wire))
	}
}

// PD2: GET /surface/decide/{ref} and /surface/decide.js serve the embedded bytes, byte-equal to
// host/surfacepage/decide/, under the unchanged surfaceCSP and the surface headers; a malformed ref
// is 400 and a foreign Host 421; the A2 and A3 pages' bytes are unchanged; the page loads exactly
// its own script.
func TestDecidePageAssets(t *testing.T) {
	const csp = "default-src 'none'; style-src 'unsafe-inline'; script-src 'self'; connect-src 'self'; font-src 'self'; base-uri 'none'; form-action 'self'; frame-ancestors 'none'"
	if surfaceCSP != csp {
		t.Fatalf("surfaceCSP changed: %q", surfaceCSP)
	}
	a := newApprovalRig(t)
	ref := "sha256:" + strings.Repeat("ab", 32)
	for _, tc := range []struct{ target, file, contentType string }{
		{"/surface/decide/" + ref, "index.html", "text/html; charset=utf-8"},
		{"/surface/decide/sha256:" + strings.Repeat("0", 64), "index.html", "text/html; charset=utf-8"},
		{"/surface/decide.js", "decide.js", "text/javascript; charset=utf-8"},
	} {
		want, err := os.ReadFile(filepath.Join("..", "surfacepage", "decide", tc.file))
		if err != nil {
			t.Fatal(err)
		}
		rec := a.serve("GET", tc.target, apLoopback, nil)
		if rec.Code != 200 || len(want) == 0 || !bytes.Equal(rec.Body.Bytes(), want) {
			t.Errorf("GET %s: %d, or not the bytes of host/surfacepage/decide/%s", tc.target, rec.Code, tc.file)
		}
		for k, v := range map[string]string{"Content-Security-Policy": csp, "Content-Type": tc.contentType,
			"Cache-Control": "no-store", "X-Content-Type-Options": "nosniff", "Referrer-Policy": "no-referrer"} {
			if got := rec.Header().Get(k); got != v {
				t.Errorf("GET %s: %s = %q, want %q", tc.target, k, got, v)
			}
		}
		if rec.Header().Get("Set-Cookie") != "" || rec.Header().Get("Access-Control-Allow-Origin") != "" {
			t.Errorf("GET %s sets a cookie or a CORS header", tc.target)
		}
		requireRefusal(t, "foreign host", a.serve("GET", tc.target, "evil.example", nil), 421, "MisdirectedRequest")
	}
	for _, bad := range []string{"sha256:" + strings.Repeat("AB", 32), "sha256:abc", strings.Repeat("ab", 32), "mr-0123456789abcdef"} {
		requireRefusal(t, bad, a.serve("GET", "/surface/decide/"+bad, apLoopback, nil), 400, "BadRequest")
	}
	for file, sum := range map[string]string{"passkey/index.html": "1695ebc2ae20f303dbfa3b617b3bec4c76035447ec290deba91e695fec9cc06d",
		"passkey/passkey.js": "0064c6c90099a3565b178cf170618b80a67bbf1f3444a1018cc4ee29f51ed53f",
		"approve/index.html": "e353429ab873633b9aae73e26cd43cd8697f7f2af647ecf7be842cad27f4950e",
		"approve/approve.js": "9c7fadd4e8a5985f53c2e95981e6746cf422c03f07a2b459d13ef67893ad4b25"} {
		b, err := os.ReadFile(filepath.Join("..", "surfacepage", file))
		if h := sha256.Sum256(b); err != nil || hex.EncodeToString(h[:]) != sum {
			t.Errorf("the A2/A3 page's %s changed (%v)", file, err)
		}
	}
	page := string(surfacepage.DecidePage)
	if scripts := regexp.MustCompile(`(?i)<script\b[^>]*>`).FindAllString(page, -1); len(scripts) != 1 || !strings.Contains(page, `<script src="/surface/decide.js"></script>`) {
		t.Errorf("scripts %q", scripts)
	}
	if m := regexp.MustCompile(`(?i)\son[a-z]+\s*=|<form|<img|<link|<iframe|https?://`).FindString(page); m != "" {
		t.Errorf("the page has %q", m)
	}
	// Labelled instrument-health controls (not load-bearing; XD1, XD2 and XD4 are).
	for _, banned := range []string{"innerHTML", "outerHTML", "insertAdjacentHTML", "document.write", "eval(", "new Function", "document.cookie", "XMLHttpRequest", "sendBeacon", "WebSocket"} {
		if bytes.Contains(surfacepage.DecideScript, []byte(banned)) {
			t.Errorf("decide.js uses %s", banned)
		}
	}
}
