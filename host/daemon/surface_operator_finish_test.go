package daemon

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"os"
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
}
