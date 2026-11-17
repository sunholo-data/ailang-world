package broker

// Row 243 V3 sprint 3 M8 (design_docs/planned/w-world-address-bar-verbs-v3.md §2.3, §6 VT2): the
// daemon's step C10, pure. A stored verb request is confirmable only if its scope is a finish
// approval that outlives the request by at most an hour from now, its title and body are the ones
// the scope approves, and its prompt is the terminal's rendering of them (D-WORLD-118 = A: no
// draft-only rule).

import (
	"crypto/sha256"
	"encoding/hex"
)

// FinishConfirmFacts are the fields of a decoded, validated world/verb-request/v1 that C10 judges,
// with the request's expiry and the daemon's clock.
type FinishConfirmFacts struct {
	Scope, Title, Body, Prompt, Committer string
	Protected                             []string
	RequestExpires, Now                   int64
}

// FinishConfirmRefusal is why f is not confirmable (F1-F5), or "" when it is.
func FinishConfirmRefusal(f FinishConfirmFacts) string {
	scope, err := ParseFinishApprovalScope(f.Scope)
	if err != nil { // F1
		return "the request's scope is not a finish approval"
	}
	if !(f.RequestExpires < scope.ExpiresAt && scope.ExpiresAt <= f.Now+3600) { // F2
		return "the request outlives its approval scope, or the scope is over an hour ahead"
	}
	if title := sha256.Sum256([]byte(f.Title)); hex.EncodeToString(title[:]) != scope.TitleSHA256 { // F3
		return "the title is not the one the scope approves"
	}
	if body := sha256.Sum256([]byte(f.Body)); hex.EncodeToString(body[:]) != scope.BodySHA256 { // F4
		return "the body is not the one the scope approves"
	}
	if f.Prompt != FinishPromptText(scope, f.Title, f.Body, f.Protected, f.Committer) { // F5
		return "the prompt is not the terminal's rendering of what is approved"
	}
	return ""
}
