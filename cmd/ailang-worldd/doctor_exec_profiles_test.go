package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/broker"
	"github.com/sunholo-data/ailang-world/host/daemon"
)

// Row 256: doctor was all green while every exec profile was left out.
// doctor (and --quick) now names each left-out profile from /v1/health as a ✗
// in the exec class (exit 13), with observation-only advice for failures
// recorded at/after the startup deadline. Control: the same daemon
// with nothing left out adds no line and fails nothing.
// Mutants: doctorExecProfiles never called; the line a warning; "exec
// profiles" outside the exec class; out_of_time ignored (one fix for both).
func TestDoctorNamesLeftOutExecProfiles(t *testing.T) {
	t.Setenv("HOME", t.TempDir())
	db := filepath.Join(t.TempDir(), "world.db")
	excluded := []daemon.HealthExecExclusion{
		{Project: "ailang-world", File: "/p/ailang-world.json", OutOfTime: true,
			Reason: `exec startup probe failed for profile "ailang-world" (7 of 7 arms): arm1-write-inside: … (the startup budget is spent; stderr "")`},
		{Project: "broken", File: "/p/broken.json",
			Reason: `exec startup probe failed for profile "broken" (1 of 7 arms): arm7-toolchain: probe "cargo --version": exit 127`},
	}
	serve := func(x []daemon.HealthExecExclusion) *httptest.Server {
		return httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if r.URL.Path != "/v1/health" {
				http.NotFound(w, r)
				return
			}
			_ = json.NewEncoder(w).Encode(daemon.HealthResponse{Status: "ok", DaemonVersion: "0.1.0", DBPath: db,
				ExecProfiles: []daemon.HealthExecProfile{}, ExecProfilesExcluded: x})
		}))
	}

	srv := serve(excluded)
	defer srv.Close()
	code, out, _ := runDoctorCLI(t, srv.URL, "--quick", "--db", db)
	if code != 13 {
		t.Fatalf("rc %d, want 13 (exec class):\n%s", code, out)
	}
	for _, want := range []string{
		"✓ daemon     ailang-worldd 0.1.0",
		"✗ exec profiles ailang-world (/p/ailang-world.json) was left out at startup: " + broker.ExecProbeDeadlineObservation,
		"fix: Review the failure reasons and restart World if appropriate",
		"✗ exec profiles broken (/p/broken.json) was left out at startup: exec startup probe failed for profile \"broken\" (1 of 7 arms): arm7-toolchain",
		"fix: review /p/broken.json and the failure reasons; correct the profile if needed, then restart World if appropriate",
		"2 check(s) failed",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("doctor lacks %q:\n%s", want, out)
		}
	}

	ok := serve(nil)
	defer ok.Close()
	code, out, _ = runDoctorCLI(t, ok.URL, "--quick", "--db", db)
	if code != exitOK || strings.Contains(out, "exec profiles") || !strings.Contains(out, "no check failed") {
		t.Fatalf("control (nothing left out): rc %d\n%s", code, out)
	}
}

// Direct report capture avoids sockets and drives the production report
// helper that doctorDaemon calls after decoding health. A genuine refusal with OutOfTime=true
// must still advise review, retaining the reason; the reason is capped at
// 240 bytes by oneLine (project/file/fixed prose are outside that cap).
func TestDoctorLateRefusalObservationAndReasonCap(t *testing.T) {
	for _, late := range []bool{false, true} {
		var out strings.Builder
		r := &doctorReport{w: &out, asJSON: true}
		reason := "genuine refusal: sibling write succeeded\n" + strings.Repeat("x", 300)
		doctorExecProfiles(r, []daemon.HealthExecExclusion{{Project: "late", File: "/p/late.json", Reason: reason, OutOfTime: late}})
		var got struct{ Detail, Fix, Mark, Class string }
		if err := json.Unmarshal([]byte(out.String()), &got); err != nil {
			t.Fatal(err)
		}
		if r.classCode != 13 || r.fails != 1 || got.Mark != "fail" || got.Class != "exec" {
			t.Fatalf("report = %s", out.String())
		}
		capped := oneLine(reason, 240-len("…"))
		if len(capped) > 240 || !strings.HasSuffix(got.Detail, capped) || strings.Contains(got.Detail, strings.Repeat("x", 241)) || strings.Contains(got.Detail, "\n") {
			t.Fatalf("reason must be one line, capped at 240 bytes: %q", got.Detail)
		}
		if late && (!strings.Contains(got.Detail, broker.ExecProbeDeadlineObservation) ||
			!strings.Contains(got.Fix, "Review the failure reasons and restart World if appropriate")) {
			t.Fatalf("late refusal needs observation-only advice: %s", out.String())
		}
		for _, forbidden := range []string{"host was too busy", "nothing to fix", "host load caused", "when the host is quieter"} {
			if strings.Contains(got.Detail+got.Fix, forbidden) {
				t.Fatalf("causal diagnosis %q: %s", forbidden, out.String())
			}
		}
	}
}
