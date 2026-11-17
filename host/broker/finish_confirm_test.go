package broker

// Row 243 V3 sprint 3 M8 (design §2.3, §6 VT2): FinishConfirmRefusal's F1-F5, each from a fixture
// whose scope is FormatFinishApprovalScope's and whose prompt is FinishPromptText's (never typed).

import (
	"crypto/sha256"
	"encoding/hex"
	"testing"
)

func vtHex(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

// vtFacts is an admissible request at now 1000 expiring 1600, its scope expiring 1900.
func vtFacts() (FinishConfirmFacts, FinishApprovalScope) {
	title, body, protected, committer := "Fix the bar", "## Summary\nline 2\n", []string{".github/workflows/ci.yml", "go.mod"}, ""
	s := FinishApprovalScope{Repository: "sunholo-data/ailang-world", Remote: "origin", Head: "world/w1", SHA: "0123456789abcdef0123456789abcdef01234567", Base: "dev",
		TitleSHA256: vtHex(title), BodySHA256: vtHex(body), Draft: true, ExpiresAt: 1900}
	return FinishConfirmFacts{Scope: FormatFinishApprovalScope(s), Title: title, Body: body, Prompt: FinishPromptText(s, title, body, protected, committer),
		Committer: committer, Protected: protected, RequestExpires: 1600, Now: 1000}, s
}

func TestFinishConfirmRefusal(t *testing.T) {
	const f1, f2, f3, f4, f5 = "the request's scope is not a finish approval", "the request outlives its approval scope, or the scope is over an hour ahead",
		"the title is not the one the scope approves", "the body is not the one the scope approves", "the prompt is not the terminal's rendering of what is approved"
	rescope := func(edit func(*FinishApprovalScope)) func(*FinishConfirmFacts) {
		return func(f *FinishConfirmFacts) {
			_, s := vtFacts()
			edit(&s)
			f.Scope = FormatFinishApprovalScope(s)
			f.Prompt = FinishPromptText(s, f.Title, f.Body, f.Protected, f.Committer)
		}
	}
	for _, tc := range []struct {
		name string
		edit func(*FinishConfirmFacts)
		want string
	}{
		{"admissible", func(*FinishConfirmFacts) {}, ""},
		{"ready", rescope(func(s *FinishApprovalScope) { s.Draft = false }), ""}, // D-WORLD-118 = A: no draft-only rule
		{"F1-not-a-scope", func(f *FinishConfirmFacts) { f.Scope = "github-pr:o/r#x" }, f1},
		{"F1-trailing-byte", func(f *FinishConfirmFacts) { f.Scope += " " }, f1},
		{"F2-scope-equals-request", rescope(func(s *FinishApprovalScope) { s.ExpiresAt = 1600 }), f2},
		{"F2-scope-before-request", rescope(func(s *FinishApprovalScope) { s.ExpiresAt = 1599 }), f2},
		{"F2-scope-an-hour-ahead", rescope(func(s *FinishApprovalScope) { s.ExpiresAt = 4600 }), ""},
		{"F2-scope-over-an-hour-ahead", rescope(func(s *FinishApprovalScope) { s.ExpiresAt = 4601 }), f2},
		{"F3-title", func(f *FinishConfirmFacts) { f.Title = "Fix the bar." }, f3},
		{"F4-body", func(f *FinishConfirmFacts) { f.Body += "x" }, f4},
		{"F5-prompt-edited", func(f *FinishConfirmFacts) { f.Prompt += "\n" }, f5},
		{"F5-protected-not-shown", func(f *FinishConfirmFacts) { f.Protected = f.Protected[:1] }, f5},
		{"F5-committer-not-shown", func(f *FinishConfirmFacts) { f.Committer = "Mark <m@x>" }, f5},
		{"F5-prompt-of-another-scope", func(f *FinishConfirmFacts) {
			_, s := vtFacts()
			s.Base = "main"
			f.Prompt = FinishPromptText(s, f.Title, f.Body, f.Protected, f.Committer)
		}, f5},
	} {
		f, _ := vtFacts()
		tc.edit(&f)
		if got := FinishConfirmRefusal(f); got != tc.want {
			t.Errorf("%s: %q, want %q", tc.name, got, tc.want)
		}
	}
}
