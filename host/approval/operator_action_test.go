package approval_test

// Row 243 V4 sprint 2 M1 (design §3.1, §5.1 TestOperatorStrictJSON): the canonical bytes and both digests,
// pinned by a golden the real encoder writes (WORLD_UPDATE_OPERATOR_ACTION_GOLDEN=1 go test ./host/approval/
// -run '^TestOperatorActionGolden$'); every invalid input class, each the golden's bytes with one change,
// refused before canonicalization by its own error.

import (
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"encoding/json"
	"errors"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
)

// opAction is a valid action of kind start or finish.
func opAction(kind string) approval.OperatorAction {
	a := approval.OperatorAction{Schema: approval.OperatorActionSchema, ID: "oa-0123456789abcdef0123456789abcdef", Kind: kind,
		Episode: "fixture-ep", Store: "/ws/state/credentials.db", Options: approval.OperatorOptions{Preset: "sol-headless"},
		Resolution: []approval.OperatorResolved{{Name: "agent", Value: "codex"}, {Name: "base_oid", Value: "8db22168689f31e19349b075e0199c8d9c940e61"},
			{Name: "model", Value: "gpt-6.2-sol<&>"}, {Name: "repo", Value: "/src/ailang-world"}},
		Task: "Fix the flaky test; $(rm -rf ~) `x` <b>é</b>\n", Review: "Start fixture-ep with sol-headless, headless.\n"}
	if kind == "start" {
		a.Options.Headless = true
	} else {
		a.Options.ConfirmOnSurface, a.Task, a.Review = true, "", "Finish fixture-ep: check, commit, then confirm the PR.\n"
	}
	a.Argv = approval.OperatorArgv(a.Kind, a.Episode, a.Options)
	return a
}

var opChallenge = [32]byte{0x6f, 0x61, 1, 2, 3}

const opExpiry = int64(1767225600)

type opGolden struct {
	Canonical, ReviewDigest, SignatureDigest, Message, MessageHash string
}

func opGoldenOf(t *testing.T, a approval.OperatorAction) opGolden {
	t.Helper()
	c, err := approval.CanonicalOperatorAction(a)
	rd, err2 := approval.OperatorReviewDigest(a)
	sd, err3 := approval.OperatorSignatureDigest(a)
	m, err4 := approval.OperatorActionMessage(sd, opChallenge, opExpiry)
	h, err5 := approval.OperatorActionMessageHash(sd, opChallenge, opExpiry)
	if err := errors.Join(err, err2, err3, err4, err5); err != nil {
		t.Fatal(err)
	}
	return opGolden{string(c), hex.EncodeToString(rd[:]), hex.EncodeToString(sd[:]), hex.EncodeToString(m), hex.EncodeToString(h[:])}
}

// opGoldenBytes is the start action's canonical bytes as the golden file holds them.
func opGoldenBytes(t *testing.T) []byte {
	t.Helper()
	b, err := os.ReadFile(filepath.Join("testdata", "operator-action.golden.json"))
	if err != nil {
		t.Fatal(err)
	}
	var g map[string]opGolden
	if err := json.Unmarshal(b, &g); err != nil {
		t.Fatal(err)
	}
	return []byte(g["start"].Canonical)
}

func TestOperatorActionGolden(t *testing.T) {
	got := map[string]opGolden{"start": opGoldenOf(t, opAction("start")), "finish": opGoldenOf(t, opAction("finish"))}
	file := filepath.Join("testdata", "operator-action.golden.json")
	if os.Getenv("WORLD_UPDATE_OPERATOR_ACTION_GOLDEN") == "1" {
		b, _ := json.MarshalIndent(got, "", " ")
		if err := os.WriteFile(file, append(b, '\n'), 0o644); err != nil {
			t.Fatal(err)
		}
	}
	b, err := os.ReadFile(file)
	if err != nil {
		t.Fatal(err)
	}
	var want map[string]opGolden
	if err := json.Unmarshal(b, &want); err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("the encoder's bytes or digests moved:\n got %+v\nwant %+v", got, want)
	}
	for kind, g := range want {
		a, err := approval.DecodeOperatorAction([]byte(g.Canonical))
		if err != nil || !reflect.DeepEqual(a, opAction(kind)) {
			t.Errorf("%s: the golden does not decode to its action: %v", kind, err)
		}
	}
}

// The message's layout, field by field: prefix, NUL, digest, challenge, big-endian expiry.
func TestOperatorActionMessage(t *testing.T) {
	d := [32]byte{0xd1, 0xd2}
	m, err := approval.OperatorActionMessage(d, opChallenge, opExpiry)
	if err != nil || approval.OperatorActionMessageLen != 97 || len(m) != 97 {
		t.Fatalf("len %d, %v; want 97", len(m), err)
	}
	if string(m[:24]) != "world-operator-action/v1" || m[24] != 0 || !bytes.Equal(m[25:57], d[:]) || !bytes.Equal(m[57:89], opChallenge[:]) ||
		binary.BigEndian.Uint64(m[89:]) != uint64(opExpiry) {
		t.Fatalf("layout: %x", m)
	}
	if _, err := approval.OperatorActionMessage(d, opChallenge, -1); !errors.Is(err, approval.ErrNegativeExpiry) {
		t.Fatalf("negative expiry: %v", err)
	}
	if _, err := approval.OperatorActionMessageHash(d, opChallenge, -1); !errors.Is(err, approval.ErrNegativeExpiry) {
		t.Fatalf("negative expiry hash: %v", err)
	}
}

// opDecodeRefused decodes b and wants err class want and the zero action.
func opDecodeRefused(t *testing.T, name string, b []byte, want error) {
	t.Helper()
	a, err := approval.DecodeOperatorAction(b)
	if !errors.Is(err, want) || !reflect.DeepEqual(a, approval.OperatorAction{}) {
		t.Errorf("%s: %v (%+v); want %v and no action", name, err, a.ID, want)
	}
}

func TestOperatorStrictJSON(t *testing.T) {
	g := opGoldenBytes(t)
	if _, err := approval.DecodeOperatorAction(g); err != nil {
		t.Fatalf("control: %v", err)
	}
	s := string(g)
	rep := func(old, new string) []byte {
		if strings.Count(s, old) != 1 {
			t.Fatalf("fixture: %q not once in %s", old, s)
		}
		return []byte(strings.Replace(s, old, new, 1))
	}
	for _, tc := range []struct {
		name string
		b    []byte
		want error
	}{
		{"duplicate top-level key", rep(`"kind":"start",`, `"kind":"start","kind":"start",`), approval.ErrOperatorDuplicateKey},
		{"duplicate nested key", rep(`"headless":true,`, `"headless":true,"headless":true,`), approval.ErrOperatorDuplicateKey},
		{"duplicate key in an array element", rep(`{"name":"agent",`, `{"name":"agent","name":"agent",`), approval.ErrOperatorDuplicateKey},
		{"unknown top-level key", rep(`"kind":"start",`, `"kind":"start","shell":"x",`), approval.ErrOperatorUnknownKey},
		{"unknown nested key", rep(`"headless":true,`, `"headless":true,"env":"x",`), approval.ErrOperatorUnknownKey},
		{"trailing object", append(bytes.Clone(g), `{}`...), approval.ErrOperatorTrailing},
		{"trailing word", append(bytes.Clone(g), ` x`...), approval.ErrOperatorTrailing},
		{"trailing newline", append(bytes.Clone(g), '\n'), approval.ErrNotCanonical},
		{"invalid UTF-8", rep(`é`, "\xc3"), approval.ErrOperatorUTF8},
		{"overlong UTF-8", rep(`é`, "\xc0\xa9"), approval.ErrOperatorUTF8},
		{"raw NUL", rep(`<b>`, "<\x00b>"), approval.ErrOperatorNUL},
		{"escaped NUL", rep(`<b>`, `<\u0000b>`), approval.ErrString},
		{"null", []byte(`null`), approval.ErrOperatorShape},
		{"array", []byte(`[` + s + `]`), approval.ErrOperatorShape},
		{"string", []byte(`"x"`), approval.ErrOperatorShape},
		{"number", []byte(`1`), approval.ErrOperatorShape},
		{"empty", nil, approval.ErrOperatorShape},
		{"unclosed", g[:len(g)-1], approval.ErrDecode},
		{"argv a string", rep(`"argv":["ailang-worldd","episode","start","fixture-ep","--preset","sol-headless","--launch","--headless"]`, `"argv":"ailang-worldd"`), approval.ErrDecode},
		{"argv an int element", rep(`"--headless"]`, `7]`), approval.ErrDecode},
		{"options an array", rep(`"options":{"preset":"sol-headless","headless":true,"confirm_on_surface":false,"ready":false}`, `"options":[]`), approval.ErrDecode},
		{"headless a string", rep(`"headless":true`, `"headless":"true"`), approval.ErrDecode},
		{"argv null", rep(`"argv":["ailang-worldd","episode","start","fixture-ep","--preset","sol-headless","--launch","--headless"]`, `"argv":null`), approval.ErrOperatorAction},
		{"id null", rep(`"id":"oa-0123456789abcdef0123456789abcdef"`, `"id":null`), approval.ErrNotCanonical},
		{"case-folded key", rep(`"kind":`, `"Kind":`), approval.ErrNotCanonical},
		{"reordered keys", rep(`"schema":"world/operator-action/v1","id":"oa-0123456789abcdef0123456789abcdef"`, `"id":"oa-0123456789abcdef0123456789abcdef","schema":"world/operator-action/v1"`), approval.ErrNotCanonical},
		{"spaced", rep(`"kind":"start"`, `"kind": "start"`), approval.ErrNotCanonical},
		{"HTML-escaped", rep(`<b>`, `\u003cb>`), approval.ErrNotCanonical},
		{"canonical but invalid", rep(`"kind":"start"`, `"kind":"launch"`), approval.ErrOperatorAction},
	} {
		opDecodeRefused(t, tc.name, tc.b, tc.want)
	}
}
