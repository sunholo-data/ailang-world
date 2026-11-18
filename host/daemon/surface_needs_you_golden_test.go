package daemon

// Row 243 V3 sprint 6 M15 (design §2.5, §6 CP7): the fixtures of app.js's Needs-you confirmations are
// GET /surface/requests.json as this handler answers it — a draft and a ready confirmation and a refused
// one on a schema-9 store, then the same read on a schema-8 store — never hand-typed JSON. The page's
// node harness (surfacepage TestNeedsYouListsConfirmations) renders them.

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"
)

// confirmGolden is the two answers: regenerate with WORLD_UPDATE_CONFIRMATIONS_GOLDEN=1; otherwise they
// must equal it byte for byte.
var confirmGolden = filepath.Join("..", "surfacepage", "testdata", "requests-confirmations.golden.json")

func TestRequestsJSONConfirmationsGolden(t *testing.T) {
	a := newApprovalRig(t)
	a.putVerb("vr-00000000000000c2", "w2", apT0-50, false)
	a.putVerb("vr-00000000000000c4", "w1", apT0-40, true)
	edited := a.putVerb("vr-00000000000000c6", "w6", apT0-20, true)
	setVerbColumn(t, a.d.credentialsPath, edited.RequestID, "request", append(edited.Request, ' '))
	get := func() json.RawMessage {
		rec := a.serve("GET", "/surface/requests.json", apLoopback, nil)
		if rec.Code != 200 || !json.Valid(rec.Body.Bytes()) {
			t.Fatalf("requests.json = %d %s", rec.Code, rec.Body.String())
		}
		return bytes.TrimSpace(rec.Body.Bytes())
	}
	v9 := get()
	setVersion(t, a.d.credentialsPath, `DROP TABLE verb_requests`, `PRAGMA user_version = 8`)
	want, err := json.Marshal(map[string]json.RawMessage{"v9": v9, "v8": get()})
	if err != nil {
		t.Fatal(err)
	}
	want = append(want, '\n')
	if os.Getenv("WORLD_UPDATE_CONFIRMATIONS_GOLDEN") == "1" {
		if err := os.WriteFile(confirmGolden, want, 0o644); err != nil {
			t.Fatal(err)
		}
	}
	if got, err := os.ReadFile(confirmGolden); err != nil || !bytes.Equal(got, want) {
		t.Errorf("%s is not the handler's answers (WORLD_UPDATE_CONFIRMATIONS_GOLDEN=1 regenerates it): %v", confirmGolden, err)
	}
}
