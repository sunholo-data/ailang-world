package projection

import (
	"bytes"
	"context"
	"errors"
	"github.com/sunholo-data/ailang-world/host/authority"
	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/coordinator"
	"github.com/sunholo-data/ailang-world/host/hashref"
	"github.com/sunholo-data/ailang-world/host/store"
	"github.com/sunholo-data/ailang-world/host/transitionreg"
	"net/http"
	"strings"
	"testing"
	"time"
)

type absenceReader struct {
	err   error
	calls int
}

func (r *absenceReader) ReadSnapshot(context.Context) (transitionreg.Snapshot, error) {
	r.calls++
	return transitionreg.Snapshot{}, r.err
}

type racedObjectStore struct {
	transitionreg.ObjectStore
	heads, objects int
	err            error
}

func (r *racedObjectStore) GetRegistryHead(context.Context, string) (hashref.HashRef, bool, error) {
	r.heads++
	return hashref.SumSHA256([]byte("raced head")), true, nil
}
func (r *racedObjectStore) GetObject(context.Context, hashref.HashRef) (store.Object, bool, error) {
	r.objects++
	return store.Object{}, false, r.err
}

func TestAgentCard_AbsentPrecheckSnapshotFailure(t *testing.T) {
	for _, name := range []string{"store_error", "real_reader", "deadline"} {
		t.Run(name, func(t *testing.T) {
			st := openStore(t)
			tok := mintToken(t, st, "ep-a", []broker.Capability{liveGrant("alpha")})
			cfg := testConfig(st)
			cfg.Heads = fixedHeads{ok: false}
			injected := error(errors.New("private snapshot fault"))
			if name == "deadline" {
				injected = context.DeadlineExceeded
			}
			direct := &absenceReader{err: injected}
			real := &racedObjectStore{ObjectStore: st, err: injected}
			cfg.Reader = direct
			if name == "real_reader" {
				cfg.Reader = transitionreg.NewReader(real)
			}
			rec := getCard(t, mustHandler(t, cfg), "Bearer "+tok)
			want := 503
			class := classUnavailable
			msg := msgUnavailable
			if name == "deadline" {
				want = 504
				class = classDeadlineExceeded
				msg = msgDeadlineExceeded
			}
			if rec.Code != want {
				t.Fatalf("card status=%d want %d: %s", rec.Code, want, rec.Body)
			}
			wantBody := `{"error":{"class":"` + class + `","message":"` + msg + `"}}` + "\n"
			if rec.Body.String() != wantBody || strings.Contains(rec.Body.String(), injected.Error()) {
				t.Fatalf("existing sanitized envelope: %q", rec.Body.String())
			}
			if name == "real_reader" {
				if real.heads != 1 || real.objects != 1 {
					t.Fatalf("real fault not reached: %d/%d", real.heads, real.objects)
				}
			} else if direct.calls != 1 {
				t.Fatalf("snapshot calls=%d", direct.calls)
			}
		})
	}
}
func TestProjectionAbsenceRaceControls(t *testing.T) {
	for _, tc := range []struct {
		name           string
		present        bool
		headErr, error error
		empty          bool
		calls          int
	}{
		{name: "true_absence", error: transitionreg.RegistryHeadAbsentError{}, empty: true, calls: 1},
		{name: "present_then_absent", present: true, error: transitionreg.RegistryHeadAbsentError{}, calls: 1},
		{name: "same_text_genuine", error: errors.New("read transition registry: head is absent"), calls: 1},
		{name: "precheck_error", headErr: errors.New("head read failed"), calls: 0},
	} {
		t.Run(tc.name, func(t *testing.T) {
			st := openStore(t)
			cfg := testConfig(st)
			r := &absenceReader{err: tc.error}
			cfg.Reader = r
			cfg.Heads = fixedHeads{ok: tc.present}
			if tc.headErr != nil {
				cfg.Heads = errHeads{err: tc.headErr}
			}
			h := mustHandler(t, cfg)
			_, ds, err := h.allowedDescriptors(boundedTestContext(t), &authority.SessionBinding{Caps: []broker.Capability{liveGrant("alpha")}})
			if (err == nil) != tc.empty || len(ds) != 0 || r.calls != tc.calls {
				t.Fatalf("classification err=%v empty=%t calls=%d want %d", err, tc.empty, r.calls, tc.calls)
			}
			tok := mintToken(t, st, "ep-a", []broker.Capability{liveGrant("alpha")})
			rec := getCard(t, h, "Bearer "+tok)
			want := 503
			if tc.empty {
				want = 200
			}
			if rec.Code != want {
				t.Fatalf("card status=%d want %d", rec.Code, want)
			}
		})
	}
}
func TestA2A_AbsentPrecheckSnapshotFailure(t *testing.T) {
	for _, realReader := range []bool{false, true} {
		t.Run(map[bool]string{false: "direct", true: "real_reader"}[realReader], func(t *testing.T) {
			st := openStore(t)
			runner := &countingRunner{}
			coord, err := coordinator.New(coordinator.Config{Store: st, Runner: runner, Binder: func(ep string, caps []broker.Capability) transitionreg.Binder {
				return broker.OpenBinder(st, ep, caps, nil)
			}, Now: func() int64 { return time.Now().Unix() }, MaxInput: 1 << 20, MaxOutput: 1 << 20})
			if err != nil {
				t.Fatal(err)
			}
			cfg := testConfig(st)
			cfg.Coordinator = coord
			cfg.Heads = fixedHeads{ok: false}
			injected := errors.New("private reached snapshot fault")
			r := &absenceReader{err: injected}
			real := &racedObjectStore{ObjectStore: st, err: injected}
			cfg.Reader = r
			if realReader {
				cfg.Reader = transitionreg.NewReader(real)
			}
			var sink bytes.Buffer
			cfg.ErrorLog = &sink
			tok := mintToken(t, st, "ep-a", []broker.Capability{liveGrant("alpha")})
			rec := postA2A(t, mustHandler(t, cfg), "Bearer "+tok, `{"jsonrpc":"2.0","id":7,"method":"tasks/send","params":{"id":"race","metadata":{"skill_id":"tools.echo"},"message":{"parts":[{"type":"data","data":{"x":1}}]}}}`)
			want := `{"error":{"code":-32603,"message":"` + notAvailableMessage + `"},"id":7,"jsonrpc":"2.0"}` + "\n"
			if rec.Code != http.StatusOK || rec.Body.String() != want {
				t.Fatalf("A2A wire=%d %q want %q", rec.Code, rec.Body.String(), want)
			}
			if strings.Count(sink.String(), injected.Error()) != 1 || runner.runs != 0 {
				t.Fatalf("genuine refusal/runner=%q/%d", sink.String(), runner.runs)
			}
			if realReader {
				if real.objects != 1 || real.heads != 1 {
					t.Fatalf("real fault counts %d/%d", real.heads, real.objects)
				}
			} else if r.calls != 1 {
				t.Fatalf("reader calls=%d", r.calls)
			}
			if _, seen, err := st.GetReceipt(boundedTestContext(t), coordinator.InvocationID("ep-a", "race")); err != nil || seen {
				t.Fatalf("unexpected receipt seen=%t err=%v", seen, err)
			}
			if _, seen, err := st.SelectedHead(boundedTestContext(t)); err != nil || seen {
				t.Fatalf("unexpected world head seen=%t err=%v", seen, err)
			}
		})
	}
}
