package surfacepage

// Row 243 V3 sprint 2 M5 (design §2.5, §6 CP2, CP4): the confirmation card's field table, the golden its
// node tests read (written by this Go test from the real tools, never by hand), and the page's assets.

import (
	"bytes"
	"encoding/hex"
	"encoding/json"
	"os"
	"path/filepath"
	"reflect"
	"regexp"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/verbs"
)

// confirmGolden is the node tests' fixture file; WORLD_UPDATE_CONFIRM_GOLDEN=1 rewrites it.
var confirmGolden = filepath.Join("testdata", "confirm-card.golden.jsonl")

// confirmFields is confirm.js's CONFIRM_FIELDS table.
func confirmFields(t *testing.T) []struct{ Key, Label, Type string } {
	t.Helper()
	_, rest, _ := bytes.Cut(ConfirmScript, []byte("/* CONFIRM_FIELDS */"))
	lit, _, _ := bytes.Cut(rest, []byte("/* /CONFIRM_FIELDS */"))
	var fields []struct{ Key, Label, Type string }
	if err := json.Unmarshal(lit, &fields); err != nil || len(fields) == 0 {
		t.Fatalf("confirm.js's CONFIRM_FIELDS: %v", err)
	}
	return fields
}

// CP2: CONFIRM_FIELDS names exactly approval.VerbRequest's json tags, in order, each with its JSON type
// (string or a list of strings) and a label: a field added to the request without the page fails here.
func TestConfirmPageFieldTable(t *testing.T) {
	var want, got []string
	ty := reflect.TypeOf(approval.VerbRequest{})
	for i := 0; i < ty.NumField(); i++ {
		kind := map[reflect.Kind]string{reflect.String: "string", reflect.Slice: "strings"}[ty.Field(i).Type.Kind()]
		if ty.Field(i).Type.Kind() == reflect.Slice && ty.Field(i).Type.Elem().Kind() != reflect.String {
			kind = ""
		}
		want = append(want, strings.Split(ty.Field(i).Tag.Get("json"), ",")[0]+":"+kind)
	}
	for _, f := range confirmFields(t) {
		if f.Label == "" {
			t.Errorf("CONFIRM_FIELDS %s has no label", f.Key)
		}
		got = append(got, f.Key+":"+f.Type)
	}
	if !reflect.DeepEqual(got, want) || len(want) != 10 { // control: design §2.2's ten fields
		t.Fatalf("CONFIRM_FIELDS %v, approval.VerbRequest %v", got, want)
	}
}

// confirmGoldenLine is one card the node tests open, and what the page must make of it.
type confirmGoldenLine struct {
	Name     string       `json:"name"`
	Card     Confirmation `json:"card"`
	ArgvLine string       `json:"argvLine"` // verbs.Quote(argv)
	Message  string       `json:"message"`  // hex of approval.VerbMessage(digest, challenge, expires)
	Hash     string       `json:"hash"`     // hex of approval.VerbMessageHash
}

// CP4, CP6: the golden is the fixtures' cards (ConfirmationCard over approval.CanonicalVerb's bytes) in
// each state, with Go's own message, its hash and the quoted argv; the node tests check the page's
// against it. WORLD_UPDATE_CONFIRM_GOLDEN=1 go test ./host/surfacepage/ -run '^TestConfirmCardGolden$'
// writes it.
func TestConfirmCardGolden(t *testing.T) {
	var out bytes.Buffer
	fx := confirmFixtures(t)
	for _, c := range []struct {
		name    string
		f       confirmFixture
		outcome string
		now     int64
	}{
		{"draft", fx[0], "", 0}, {"hostile", fx[1], "", 0}, {"draft-approved", fx[0], ConfirmationApproved, 0},
		{"draft-superseded", fx[0], ConfirmationSuperseded, 0}, {"draft-expired", fx[0], "", fx[0].in.Expires},
	} {
		in := c.f.in
		if c.outcome != "" {
			in.Outcome = c.outcome
		}
		if c.now != 0 {
			in.Now = c.now
		}
		if c.outcome == ConfirmationApproved {
			in.ApprovedBy, in.Decision = "webauthn:mark-mac", "sha256:"+strings.Repeat("5e", 32)
		}
		var d, ch [32]byte
		if _, err := hex.Decode(d[:], []byte(in.Digest)); err != nil {
			t.Fatal(err)
		}
		if _, err := hex.Decode(ch[:], []byte(in.Challenge)); err != nil {
			t.Fatal(err)
		}
		m, err := approval.VerbMessage(d, ch, in.Expires)
		if err != nil {
			t.Fatal(err)
		}
		h, _ := approval.VerbMessageHash(d, ch, in.Expires)
		line, err := json.Marshal(confirmGoldenLine{Name: c.name, Card: ConfirmationCard(in, []CardPasskey{{Origin: "http://localhost:1", RPID: "localhost", CredentialID: []byte{7, 7}}}),
			ArgvLine: verbs.Quote(c.f.req.Argv), Message: hex.EncodeToString(m), Hash: hex.EncodeToString(h[:])})
		if err != nil {
			t.Fatal(err)
		}
		out.Write(append(line, '\n'))
	}
	if os.Getenv("WORLD_UPDATE_CONFIRM_GOLDEN") == "1" {
		if err := os.WriteFile(confirmGolden, out.Bytes(), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	got, err := os.ReadFile(confirmGolden)
	if err != nil || !bytes.Equal(got, out.Bytes()) {
		t.Fatalf("%s is not the cards the fixtures make (%v); regenerate it with WORLD_UPDATE_CONFIRM_GOLDEN=1", confirmGolden, err)
	}
}

// The page loads exactly its own script and has no inline handler, form, image, link, frame or remote
// URL; the ids confirm.js reads exist; the embedded bytes are the files'. Labelled instrument-health
// controls (not load-bearing; the node tests and XC1 are): the script uses no HTML sink, eval, cookie
// or other transport.
func TestConfirmPageAssets(t *testing.T) {
	for file, b := range map[string][]byte{"index.html": ConfirmPage, "confirm.js": ConfirmScript} {
		want, err := os.ReadFile(filepath.Join("confirm", file))
		if err != nil || len(want) == 0 || !bytes.Equal(b, want) {
			t.Fatalf("confirm/%s is not embedded (%v)", file, err)
		}
	}
	page := string(ConfirmPage)
	if scripts := regexp.MustCompile(`(?i)<script\b[^>]*>`).FindAllString(page, -1); len(scripts) != 1 || !strings.Contains(page, `<script src="/surface/confirm.js"></script>`) {
		t.Errorf("scripts %q", scripts)
	}
	if m := regexp.MustCompile(`(?i)\son[a-z]+\s*=|<form|<img|<link|<iframe|https?://`).FindString(page); m != "" {
		t.Errorf("the page has %q", m)
	}
	if !strings.Contains(page, `<header class="world-header">`) || !strings.Contains(page, `class="world-logo"><svg`) ||
		!strings.Contains(page, `role="group" aria-label="Pull request branch direction"`) {
		t.Error("the page needs the World brand and pull request direction")
	}
	if !strings.Contains(page, `<a id="back" href="/surface">`) ||
		!strings.Contains(page, `role="status" aria-atomic="true">`) ||
		!strings.Contains(page, `id="status" tabindex="-1"`) ||
		!strings.Contains(page, `id="answer" role="status"`) {
		t.Error("the page needs a World link and focusable, announced status")
	}
	for _, pair := range [][2]string{{"overview", "shown"}, {"status", "review"}, {"review", "go"}, {"go", "signed"}} {
		first := strings.Index(page, `id="`+pair[0]+`"`)
		second := strings.Index(page, `id="`+pair[1]+`"`)
		if first < 0 || first >= second {
			t.Errorf("#%s should precede #%s", pair[0], pair[1])
		}
	}
	verification := regexp.MustCompile(`(?s)<details id="verification">\s*<summary>Verification details</summary>.*?</details>`).FindString(page)
	for _, id := range []string{"prompt", "argv", "signed"} {
		if !strings.Contains(verification, `id="`+id+`"`) {
			t.Errorf("#%s should be in the collapsed verification group", id)
		}
	}
	for _, id := range []string{"f-body", "save", "prompt", "argv", "signed"} {
		found := false
		for _, block := range regexp.MustCompile(`(?s)<details(?: id="(?:ssh-fallback|verification)")?>.*?</details>`).FindAllString(page, -1) {
			if strings.Contains(block, `id="`+id+`"`) {
				found = true
			}
		}
		if !found {
			t.Errorf("#%s should be in a collapsed disclosure", id)
		}
	}
	for _, m := range regexp.MustCompile(`\$\('([a-z-]+)'\)`).FindAllStringSubmatch(string(ConfirmScript), -1) {
		if !strings.Contains(page, ` id="`+m[1]+`"`) {
			t.Errorf("confirm.js reads #%s, which the page lacks", m[1])
		}
	}
	for _, banned := range []string{"innerHTML", "outerHTML", "insertAdjacentHTML", "document.write", "eval(", "new Function", "document.cookie", "XMLHttpRequest", "sendBeacon", "WebSocket"} {
		if bytes.Contains(ConfirmScript, []byte(banned)) {
			t.Errorf("confirm.js uses %s", banned)
		}
	}
}
