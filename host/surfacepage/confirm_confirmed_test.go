package surfacepage

// Row 243 V3 sprint 4 M10 (design §2.6 "ok"): the POST's success body, the design's wire.

import (
	"encoding/json"
	"strings"
	"testing"
)

func TestConfirmedWire(t *testing.T) {
	ref := "sha256:" + strings.Repeat("5e", 32)
	got, err := json.Marshal(Confirmed("vr-0123456789abcdef", "webauthn:mark", ref))
	want := `{"v":1,"confirmed":true,"request":"vr-0123456789abcdef","approvedBy":"webauthn:mark","decision":"` + ref + `","next":"the waiting episode finish opens the pull request"}`
	if err != nil || string(got) != want {
		t.Fatalf("Confirmed = %s (%v), want %s", got, err, want)
	}
}
