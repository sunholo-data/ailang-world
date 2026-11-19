package approval_test

// Row 243 V4 sprint 2 M2 (design §3.1, §5.1 the envelope half of TestOperatorRequestLimits): the
// 128 KiB and 8 KiB caps at the exact boundary (one under, at, one over), every Validate rule refused
// on its own from a valid action, and the two digests: the review digest covers the semantic inputs
// and not the id, the signature digest covers both.

import (
	"bytes"
	"errors"
	"fmt"
	"slices"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/verbs"
)

// The envelope canonicalizes host/verbs' admitted tuple: OperatorArgv is the bar's Intent.Argv for every
// admitted line, an action built from the Intent validates, and the two action-id grammars agree.
func TestOperatorArgvIsTheBarsIntent(t *testing.T) {
	yes := verbs.OperatorPolicy{Start: "A", Finish: "A", HandoffProven: true}
	for _, line := range []string{"episode start e1 --preset p --launch --headless", "episode start e1 --launch=true --preset=p",
		"episode finish e1 --preset p", "episode finish e1 --ready --preset p --confirm-on-surface", "episode finish e1 --preset p --ready"} {
		w, _ := verbs.Split(line)
		in := verbs.AdmitOperator(w, yes).Intent
		a := opAction(in.Kind)
		a.Episode, a.Options = in.Episode, approval.OperatorOptions{Preset: in.Preset, Headless: in.Headless, ConfirmOnSurface: in.ConfirmOnSurface, Ready: in.Ready}
		if a.Argv = in.Argv; !slices.Equal(approval.OperatorArgv(a.Kind, a.Episode, a.Options), in.Argv) || a.Validate() != nil {
			t.Errorf("%s: %q, %v", line, approval.OperatorArgv(a.Kind, a.Episode, a.Options), a.Validate())
		}
	}
	for _, id := range []string{"oa-0123456789abcdef0123456789abcdef", "oa-0123456789abcdef0123456789abcde", "oa-0123456789abcdef0123456789abcdef0",
		"OA-0123456789abcdef0123456789abcdef", "oa-0123456789ABCDEF0123456789abcdef", "xoa-0123456789abcdef0123456789abcde", "oa-0123456789abcdef0123456789abcdef\n"} {
		a := opAction("start")
		a.ID = id
		if (a.Validate() == nil) != verbs.ValidOperatorActionID(id) {
			t.Errorf("id %q: the envelope and the bar disagree", id)
		}
	}
}

// opSized is the start action with its review padded so its canonical bytes are exactly n long.
func opSized(t *testing.T, n int) approval.OperatorAction {
	t.Helper()
	a := opAction("start")
	a.Review = "r"
	c, _ := approval.CanonicalOperatorAction(a)
	a.Review = strings.Repeat("r", n-len(c)+1)
	if c, _ = approval.CanonicalOperatorAction(a); len(c) != n {
		t.Fatalf("sized %d, want %d", len(c), n)
	}
	return a
}

func TestOperatorRequestLimits(t *testing.T) {
	max := approval.MaxOperatorActionBytes
	if max != 131072 || approval.MaxOperatorTaskBytes != 8192 {
		t.Fatalf("caps %d %d", max, approval.MaxOperatorTaskBytes)
	}
	at := opSized(t, max)
	c, _ := approval.CanonicalOperatorAction(at)
	if err := at.Validate(); err != nil {
		t.Fatalf("exactly 128 KiB: %v", err)
	}
	if _, err := approval.DecodeOperatorAction(c); err != nil {
		t.Fatalf("decode exactly 128 KiB: %v", err)
	}
	for _, n := range []int{max - 1, max} {
		if err := opSized(t, n).Validate(); err != nil {
			t.Errorf("%d bytes: %v", n, err)
		}
	}
	if err := opSized(t, max+1).Validate(); !errors.Is(err, approval.ErrOperatorTooLarge) {
		t.Errorf("one over: %v", err)
	}
	over, _ := approval.CanonicalOperatorAction(opSized(t, max+1))
	opDecodeRefused(t, "decode one over", over, approval.ErrOperatorTooLarge)
	opDecodeRefused(t, "decode one over by a space, before any parse", append(bytes.Clone(c), ' '), approval.ErrOperatorTooLarge)
	under, _ := approval.CanonicalOperatorAction(opSized(t, max-1))
	opDecodeRefused(t, "at the cap with a space", append(under, ' '), approval.ErrNotCanonical)
	for _, tc := range []struct {
		name, kind, task string
		ok               bool
	}{
		{"8 KiB of two-byte runes", "start", strings.Repeat("é", 4096), true},
		{"one byte", "start", "x", true},
		{"8 KiB + 1", "start", strings.Repeat("é", 4096) + "x", false},
		{"8 KiB - 1 + a two-byte rune", "start", strings.Repeat("x", 8191) + "é", false},
		{"empty start", "start", "", false},
		{"a finish's task", "finish", "x", false},
		{"no finish task", "finish", "", true},
	} {
		a := opAction(tc.kind)
		a.Task = tc.task
		err := a.Validate()
		if tc.ok != (err == nil) || (!tc.ok && !errors.Is(err, approval.ErrOperatorTask)) {
			t.Errorf("task %s: %v", tc.name, err)
		}
	}
}

func TestOperatorActionValidate(t *testing.T) {
	for _, kind := range []string{"start", "finish"} {
		if err := opAction(kind).Validate(); err != nil {
			t.Fatalf("control %s: %v", kind, err)
		}
	}
	res := func(names ...string) []approval.OperatorResolved {
		r := []approval.OperatorResolved{}
		for _, n := range names {
			r = append(r, approval.OperatorResolved{Name: n, Value: "v"})
		}
		return r
	}
	many := []string{}
	for i := 0; i < 65; i++ {
		many = append(many, fmt.Sprintf("n%02d", i))
	}
	argv := func(a *approval.OperatorAction) { a.Argv = approval.OperatorArgv(a.Kind, a.Episode, a.Options) }
	for _, tc := range []struct {
		name, kind string
		edit       func(*approval.OperatorAction)
		want       error
	}{
		{"schema", "start", func(a *approval.OperatorAction) { a.Schema = approval.VerbRequestSchema }, approval.ErrSchema},
		{"NUL in a leaf", "start", func(a *approval.OperatorAction) { a.Resolution[0].Value = "a\x00" }, approval.ErrString},
		{"invalid UTF-8 in a leaf", "start", func(a *approval.OperatorAction) { a.Review = "\xff" }, approval.ErrString},
		{"no id", "start", func(a *approval.OperatorAction) { a.ID = "" }, approval.ErrOperatorAction},
		{"id with a prefix", "start", func(a *approval.OperatorAction) { a.ID = "x" + a.ID }, approval.ErrOperatorAction},
		{"id with a suffix", "start", func(a *approval.OperatorAction) { a.ID += "0" }, approval.ErrOperatorAction},
		{"id with a newline", "start", func(a *approval.OperatorAction) { a.ID += "\n" }, approval.ErrOperatorAction},
		{"id upper-case", "start", func(a *approval.OperatorAction) { a.ID = strings.ToUpper(a.ID) }, approval.ErrOperatorAction},
		{"id 31 hex", "start", func(a *approval.OperatorAction) { a.ID = a.ID[:34] }, approval.ErrOperatorAction},
		{"id mr-", "start", func(a *approval.OperatorAction) { a.ID = "mr" + a.ID[2:] }, approval.ErrOperatorAction},
		{"kind launch", "start", func(a *approval.OperatorAction) { a.Kind = "launch"; argv(a) }, approval.ErrOperatorAction},
		{"kind Start", "start", func(a *approval.OperatorAction) { a.Kind = "Start"; argv(a) }, approval.ErrOperatorAction},
		{"episode", "start", func(a *approval.OperatorAction) { a.Episode = "-ep"; argv(a) }, approval.ErrEpisode},
		{"store relative", "start", func(a *approval.OperatorAction) { a.Store = "state/credentials.db" }, approval.ErrPath},
		{"store unclean", "start", func(a *approval.OperatorAction) { a.Store = "/ws/x/../credentials.db" }, approval.ErrPath},
		{"store empty", "start", func(a *approval.OperatorAction) { a.Store = "" }, approval.ErrPath},
		{"preset empty", "start", func(a *approval.OperatorAction) { a.Options.Preset = ""; argv(a) }, approval.ErrOperatorAction},
		{"preset se-tools", "start", func(a *approval.OperatorAction) { a.Options.Preset = "se-tools"; argv(a) }, approval.ErrOperatorAction},
		{"preset 52", "start", func(a *approval.OperatorAction) { a.Options.Preset = strings.Repeat("p", 52); argv(a) }, approval.ErrOperatorAction},
		{"preset upper-case", "start", func(a *approval.OperatorAction) { a.Options.Preset = "Sol"; argv(a) }, approval.ErrOperatorAction},
		{"start ready", "start", func(a *approval.OperatorAction) { a.Options.Ready = true; argv(a) }, approval.ErrOperatorAction},
		{"start confirm-on-surface", "start", func(a *approval.OperatorAction) { a.Options.ConfirmOnSurface = true; argv(a) }, approval.ErrOperatorAction},
		{"finish headless", "finish", func(a *approval.OperatorAction) { a.Options.Headless = true; argv(a) }, approval.ErrOperatorAction},
		{"argv a word dropped", "start", func(a *approval.OperatorAction) { a.Argv = a.Argv[:7] }, approval.ErrOperatorAction},
		{"argv a word added", "start", func(a *approval.OperatorAction) { a.Argv = append(a.Argv, "--grant") }, approval.ErrOperatorAction},
		{"argv reordered", "start", func(a *approval.OperatorAction) { a.Argv[6], a.Argv[7] = a.Argv[7], a.Argv[6] }, approval.ErrOperatorAction},
		{"argv another episode", "start", func(a *approval.OperatorAction) { a.Argv[3] = "other-ep" }, approval.ErrOperatorAction},
		{"argv without --headless", "start", func(a *approval.OperatorAction) { a.Options.Headless = false }, approval.ErrOperatorAction},
		{"argv without --confirm-on-surface", "finish", func(a *approval.OperatorAction) { a.Options.ConfirmOnSurface = false }, approval.ErrOperatorAction},
		{"argv with --ready", "finish", func(a *approval.OperatorAction) { a.Argv = append(a.Argv, "--ready") }, approval.ErrOperatorAction},
		{"resolution nil", "start", func(a *approval.OperatorAction) { a.Resolution = nil }, approval.ErrOperatorAction},
		{"resolution empty", "start", func(a *approval.OperatorAction) { a.Resolution = res() }, approval.ErrOperatorAction},
		{"resolution 65", "start", func(a *approval.OperatorAction) { a.Resolution = res(many...) }, approval.ErrOperatorAction},
		{"resolution unsorted", "start", func(a *approval.OperatorAction) { a.Resolution = res("b", "a") }, approval.ErrOperatorAction},
		{"resolution repeated", "start", func(a *approval.OperatorAction) { a.Resolution = res("a", "a") }, approval.ErrOperatorAction},
		{"resolution name", "start", func(a *approval.OperatorAction) { a.Resolution = res("Agent") }, approval.ErrOperatorAction},
		{"review empty", "start", func(a *approval.OperatorAction) { a.Review = "" }, approval.ErrOperatorAction},
	} {
		a := opAction(tc.kind)
		a.Resolution = append([]approval.OperatorResolved(nil), a.Resolution...)
		tc.edit(&a)
		if err := a.Validate(); !errors.Is(err, tc.want) {
			t.Errorf("%s: %v, want %v", tc.name, err, tc.want)
		}
		if _, err := approval.OperatorSignatureDigest(a); !errors.Is(err, tc.want) {
			t.Errorf("%s: the signature digest of an invalid action: %v", tc.name, err)
		}
	}
	for _, ok := range [][]string{{"a"}, many[:64]} {
		a := opAction("start")
		a.Resolution = res(ok...)
		if err := a.Validate(); err != nil {
			t.Errorf("resolution of %d: %v", len(a.Resolution), err)
		}
	}
	a := opAction("start")
	a.Options.Preset = strings.Repeat("p", 51)
	a.Argv = approval.OperatorArgv(a.Kind, a.Episode, a.Options)
	if err := a.Validate(); err != nil {
		t.Errorf("preset 51: %v", err)
	}
}

// The review digest covers the semantic inputs and not the id; the signature digest covers both.
func TestOperatorDigestsAreDistinct(t *testing.T) {
	a := opAction("start")
	rd, _ := approval.OperatorReviewDigest(a)
	sd, _ := approval.OperatorSignatureDigest(a)
	if rd == sd || rd == [32]byte{} || sd == [32]byte{} {
		t.Fatalf("review %x signature %x", rd, sd)
	}
	other := a
	other.ID = "oa-ffffffffffffffffffffffffffffffff"
	rd2, _ := approval.OperatorReviewDigest(other)
	sd2, _ := approval.OperatorSignatureDigest(other)
	if rd2 != rd || sd2 == sd {
		t.Fatalf("another id: review moved %v, signature kept %v", rd2 != rd, sd2 == sd)
	}
	noID := a
	noID.ID = ""
	if d, err := approval.OperatorReviewDigest(noID); err != nil || d != rd {
		t.Fatalf("the review digest needs no id: %v", err)
	}
	if _, err := approval.OperatorSignatureDigest(noID); !errors.Is(err, approval.ErrOperatorAction) {
		t.Fatalf("the signature digest needs an id: %v", err)
	}
	bad := a
	bad.Kind = "launch"
	if _, err := approval.OperatorReviewDigest(bad); !errors.Is(err, approval.ErrOperatorAction) {
		t.Fatalf("the review digest of an invalid action: %v", err)
	}
	for name, edit := range opFieldEdits() {
		b := opAction("start")
		edit(&b)
		r, err := approval.OperatorReviewDigest(b)
		s, err2 := approval.OperatorSignatureDigest(b)
		if err != nil || err2 != nil || r == rd || s == sd {
			t.Errorf("%s: review moved %v, signature moved %v (%v, %v)", name, r != rd, s != sd, err, err2)
		}
	}
}

// opFieldEdits changes one semantic field of a valid start action, keeping it valid.
func opFieldEdits() map[string]func(*approval.OperatorAction) {
	re := func(a *approval.OperatorAction) { a.Argv = approval.OperatorArgv(a.Kind, a.Episode, a.Options) }
	return map[string]func(*approval.OperatorAction){
		"episode":  func(a *approval.OperatorAction) { a.Episode = "other-ep"; re(a) },
		"store":    func(a *approval.OperatorAction) { a.Store = "/other/credentials.db" },
		"preset":   func(a *approval.OperatorAction) { a.Options.Preset = "sol-two"; re(a) },
		"headless": func(a *approval.OperatorAction) { a.Options.Headless = false; re(a) },
		"resolution": func(a *approval.OperatorAction) {
			a.Resolution = []approval.OperatorResolved{{Name: "agent", Value: "claude"}}
		},
		"task":   func(a *approval.OperatorAction) { a.Task += "!" },
		"review": func(a *approval.OperatorAction) { a.Review += "!" },
	}
}
