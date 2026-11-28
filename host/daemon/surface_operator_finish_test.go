package daemon

import (
	"context"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/sunholo-data/ailang-world/host/approval"
	"github.com/sunholo-data/ailang-world/host/approval/approvaltest"
	"github.com/sunholo-data/ailang-world/host/store"
)

func TestOperatorFinishAdmission(t *testing.T) {
	f := newOperatorHTTP(t, 3*time.Minute)
	reader, err := os.ReadFile("../../docs/QUICKSTART.md")
	if err != nil {
		t.Fatal(err)
	}
	_, usage, ok := strings.Cut(string(reader), "### Prepare a local finish in the World app")
	_, usage, found := strings.Cut(usage, "```text\n")
	command, _, closed := strings.Cut(usage, "\n```")
	if !ok || !found || !closed {
		t.Fatal("QUICKSTART finish fixture is absent")
	}
	query := "/surface/verb.json?q=" + url.QueryEscape(command)
	path, err := store.CanonicalPath(f.d.credentialsPath)
	if err != nil {
		t.Fatal(err)
	}
	manifest := store.OperatorFinishReview{Store: path, Episode: "ep", Preset: "p",
		Report: "private-manifest-fixture-marker", Commit: store.OperatorLocalCommit{Message: "signed local finish"}}
	catalogue := map[string]json.RawMessage{}
	if err := json.Unmarshal(f.h.Catalogue, &catalogue); err != nil {
		t.Fatal(err)
	}
	catalogue[store.OperatorFinishCatalogueKey] = opJSON(map[string]store.OperatorFinishReview{"ep": manifest})
	f.h.Catalogue = opJSON(catalogue)
	f.publish()
	get := func(path string) *verbAnswer {
		t.Helper()
		response := f.do("GET", path, nil, nil)
		var answer verbAnswer
		if response.Code != http.StatusOK || json.Unmarshal(response.Body.Bytes(), &answer) != nil {
			t.Fatalf("read %d %s", response.Code, response.Body.String())
		}
		return &answer
	}
	before := f.evidence()
	if got := get(query); got.Outcome != "terminal" || !strings.Contains(got.Reason, "blocked: github-auth") || f.evidence() != before {
		t.Fatalf("finish readiness refusal: %+v", got)
	}
	f.h.Readiness = "both"
	f.publish()
	before = f.evidence()
	review := get(query)
	if review.Outcome != "review" || !strings.Contains(review.Text, manifest.Report) || f.evidence() != before {
		t.Fatalf("review was unavailable or effected a write: %+v", review)
	}

	post := func(path string, data []byte) *httptest.ResponseRecorder { return f.do("POST", path, data, nil) }
	req := operatorRequestBody{Argv: review.Argv, OperationID: opRequest(35).OperationID, CatalogueGeneration: review.CatalogueGeneration}
	wrong := req
	wrong.Task = "no task belongs to prep"
	if got := post(operatorPath, opJSON(wrong)); got.Code != 400 || !strings.Contains(got.Body.String(), "BadRequest") || f.evidence() != before {
		t.Fatalf("task refusal changed state: %+v", got)
	}
	wrong = req
	wrong.CatalogueGeneration--
	if got := post(operatorPath, opJSON(wrong)); got.Code != 503 || !strings.Contains(got.Body.String(), "RunnerUnavailable") || f.evidence() != before {
		t.Fatalf("stale catalogue changed state: %+v", got)
	}

	for _, tc := range []struct {
		name, class string
		code        int
		edit        func(*store.OperatorFinishReview)
	}{
		{"other-store", "RunnerUnavailable", http.StatusServiceUnavailable, func(m *store.OperatorFinishReview) {
			m.Store = filepath.Join(filepath.Dir(path), "other-fixture.db")
		}},
		{"protected", "RunnerUnavailable", http.StatusConflict, func(m *store.OperatorFinishReview) { m.Protected = []string{"protected.txt"} }},
		{"unattributed", "RunnerUnavailable", http.StatusConflict, func(m *store.OperatorFinishReview) { m.Unattributed = []string{"unattributed.txt"} }},
	} {
		t.Run(tc.name, func(t *testing.T) {
			m := manifest
			tc.edit(&m)
			current := f.h.Catalogue
			catalogue[store.OperatorFinishCatalogueKey] = opJSON(map[string]store.OperatorFinishReview{"ep": m})
			f.h.Catalogue = opJSON(catalogue)
			f.publish()
			t.Cleanup(func() {
				f.h.Catalogue = current
				f.publish()
			})
			ro, err := store.OpenCredentialsReadOnly(path)
			if err != nil {
				t.Fatal(err)
			}
			defer ro.Close()
			h, err := ro.OperatorHeartbeat(boundedTestContext(t))
			if err != nil {
				t.Fatal(err)
			}
			attempt := req
			attempt.CatalogueGeneration = h.CatalogueGeneration
			before := f.evidence()
			got := post(operatorPath, opJSON(attempt))
			if got.Code != tc.code || !strings.Contains(got.Body.String(), tc.class) || before != f.evidence() {
				t.Fatalf("request refusal lost class or changed durable state: HTTP %d %s", got.Code, got.Body.String())
			}
		})
	}
	review = get(query)
	req.CatalogueGeneration = review.CatalogueGeneration
	key := approvaltest.NewPasskey(approval.AlgEdDSA, 43)
	f.enrolPasskey("finish-key", key, opOrigin)
	var requested operatorStatus
	if got := post(operatorPath, opJSON(req)); got.Code != 200 || json.Unmarshal(got.Body.Bytes(), &requested) != nil ||
		requested.Kind != "finish" || requested.State != "requested" {
		t.Fatalf("request %+v", got)
	}
	row := f.row(requested.ID)
	card := f.do("GET", requested.Href, nil, nil)
	if card.Code != 200 || !strings.Contains(card.Body.String(), manifest.Report) || strings.Contains(card.Body.String(), `id="operator-task"`) {
		t.Fatal("signed card did not display exact finish review without a task")
	}
	for _, path := range []string{"/surface/requests.json", "/surface/home-summary.json"} {
		routine := f.do("GET", path, nil, nil)
		if routine.Code != 200 || strings.Contains(routine.Body.String(), manifest.Report) || strings.Contains(routine.Body.String(), "artifact") ||
			strings.Contains(routine.Body.String(), opTask) {
			t.Fatalf("routine action list exposed protected payloads: %s HTTP %d", path, routine.Code)
		}
	}
	if got := post(operatorPath, opJSON(req)); got.Code != 200 || !strings.Contains(got.Body.String(), requested.ID) {
		t.Fatalf("retry lost same request: %+v", got)
	}
	t.Run("existing-operation/other-generation", func(t *testing.T) {
		other := req
		other.CatalogueGeneration++
		before := f.evidence()
		got := post(operatorPath, opJSON(other))
		if got.Code != http.StatusConflict || !strings.Contains(got.Body.String(), `"Conflict"`) ||
			f.evidence() != before || !reflect.DeepEqual(row, f.row(row.ID)) {
			t.Fatalf("retry at another catalogue generation changed state or lacked conflict: HTTP %d %s", got.Code, got.Body.String())
		}
	})
	t.Run("status/start-only-runner", func(t *testing.T) {
		f.h.Readiness = "start"
		f.publish()
		t.Cleanup(func() { f.h.Readiness = "both"; f.publish() })
		before := f.evidence()
		resp := f.do("GET", requested.StatusHref, nil, nil)
		var read operatorStatus
		const reason = "finish blocked: github-auth; repair PATH, GitHub login or Git credential helper in attended service-runner setup and re-run readiness"
		if resp.Code != http.StatusOK || json.Unmarshal(resp.Body.Bytes(), &read) != nil ||
			read.Reason != reason || read.State != "requested" || read.Kind != "finish" || read.ID != row.ID ||
			read.RunnerSeenAt != f.h.SeenAt || read.LocalCommit != "" || len(read.Checks) != 0 ||
			!strings.Contains(f.do("GET", requested.Href, nil, nil).Body.String(), reason) || f.evidence() != before ||
			!reflect.DeepEqual(row, f.row(row.ID)) {
			t.Fatalf("start-only runner lost finish remedy or changed evidence: %+v", read)
		}
	})
	artifact := approveBody(hex.EncodeToString(row.Digest),
		key.Assert(opOrigin, "localhost", sha256.Sum256(opMessage(t, [32]byte(row.Digest), [32]byte(row.Challenge), row.ExpiresAt)), 2))
	before = f.evidence()
	if got := post(requested.Href+"/confirm", sshBody(hex.EncodeToString(row.Digest), []byte("y"))); got.Code != 422 || !strings.Contains(got.Body.String(), "NotApproved") || f.evidence() != before {
		t.Fatalf("plain y changed state: %+v", got)
	}
	if got := post(requested.Href+"/confirm", artifact); got.Code != 200 || !strings.Contains(got.Body.String(), "confirmed") {
		t.Fatalf("key refused: %+v", got)
	}
	first := f.evidence()
	if got := post(requested.Href+"/confirm", artifact); got.Code != 200 || !strings.Contains(got.Body.String(), "already-confirmed") || f.evidence() != first {
		t.Fatalf("replay changed durable state: %+v", got)
	}
	if e := f.evidence(); e.actions != 1 || e.operations != 1 || e.sessions != before.sessions || !reflect.DeepEqual(e.world, before.world) {
		t.Fatalf("unexpected prep-route effects %+v", e)
	}
	for _, output := range []string{f.do("GET", requested.StatusHref, nil, nil).Body.String(), f.log.String()} {
		if strings.Contains(output, manifest.Report) || strings.Contains(output, "authenticatorData") || f.row(row.ID).V3ID != "" {
			t.Fatal("confirmed status exposed prep data or issued a V3 prompt")
		}
	}

	t.Run("confirm/without-surface-option", func(t *testing.T) {
		other := manifest
		other.Episode = "finish-elsewhere"
		a, err := other.Action(req.CatalogueGeneration, false)
		if err != nil {
			t.Fatal(err)
		}
		a.ID = "oa-" + strings.Repeat("b", 32)
		a.Options.ConfirmOnSurface = false
		a.Argv = approval.OperatorArgv(a.Kind, a.Episode, a.Options)
		challenge := sha256.Sum256([]byte(t.Name()))
		var pending store.OperatorActionRow
		f.write(func(ctx context.Context, c *store.Credentials) error {
			var err error
			pending, err = c.RequestOperatorAction(ctx, store.OperatorRequest{
				Action: a, OperationID: opRequest(55).OperationID, Challenge: challenge,
				CreatedAt: apT0, ExpiresAt: apT0 + 600,
			}, apT0)
			return err
		})
		digest := hex.EncodeToString(pending.Digest)
		msg := opMessage(t, [32]byte(pending.Digest), challenge, pending.ExpiresAt)
		proof := approveBody(digest, key.Assert(opOrigin, "localhost", sha256.Sum256(msg), 3))
		before := f.evidence()
		got := post(operatorPath+"/"+pending.ID+"/confirm", proof)
		if got.Code != http.StatusConflict || !strings.Contains(got.Body.String(), `"Conflict"`) ||
			f.evidence() != before || !reflect.DeepEqual(pending, f.row(pending.ID)) {
			t.Fatalf("signed non-surface finish crossed confirm: HTTP %d %s", got.Code, got.Body.String())
		}
	})
}
