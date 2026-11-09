package broker

import (
	"errors"
	"strings"
	"testing"
	"time"
)

// Row 256: timestamps describe observations only. Genuine refusals may be
// recorded late; all failures exclude the profile regardless of this flag.
// Mutants: fail ignores the deadline; OutOfTime tests ANY rather than ALL;
// applyProbe serves a failed profile or drops the flag.
func TestExecProbeOutOfTimeNamesTheClockNotAnArm(t *testing.T) {
	past, future := time.Now().Add(-time.Second), time.Now().Add(time.Hour)
	judge := func(deadline time.Time) ExecProbeArmFailure {
		pp := &execProber{deadline: deadline}
		pp.fail(ExecProbeArmWriteInside, "token %q", "")
		return pp.failed[0]
	}
	if !judge(past).spent {
		t.Fatal("an arm judged after the startup deadline is not spent")
	}
	if judge(future).spent {
		t.Fatal("an arm judged with time remaining is spent")
	}
	if judge(time.Time{}).spent {
		t.Fatal("an arm judged with no deadline at all is spent")
	}

	spent := ExecProbeArmFailure{Arm: ExecProbeArmWriteInside, Why: "the startup budget is spent", spent: true}
	refused := ExecProbeArmFailure{Arm: ExecProbeArmWriteOutside, Why: "token \"WORLD-PROBE-WROTE\""}
	// Fixed past deadline; no sleep or scheduler timing assumption. The arm
	// is a genuine-refusal fixture (confinement failed), not a timeout reason.
	late := &execProber{deadline: time.Unix(1, 0)}
	late.fail(ExecProbeArmWriteOutside, "genuine refusal: sibling write succeeded")
	for _, tc := range []struct {
		name   string
		failed []ExecProbeArmFailure
		want   bool
	}{
		{"every arm spent", []ExecProbeArmFailure{spent, spent}, true},
		{"one arm refused in time", []ExecProbeArmFailure{spent, refused}, false},
		{"only a refusal", []ExecProbeArmFailure{refused}, false},
		{"genuine refusal recorded after deadline", late.failed, true},
		{"no failed arm", nil, false},
	} {
		if got := (&ExecProbeError{Project: "p", Failed: tc.failed}).OutOfTime(); got != tc.want {
			t.Errorf("%s: OutOfTime = %t, want %t", tc.name, got, tc.want)
		}
	}

	// applyProbe: the profile whose failures were all recorded late is left out, flagged, and
	// not served; the refused one is left out, not flagged; a non-probe error
	// is not flagged. Only is nil beside them.
	a, b, c, d := &ExecProfile{Project: "a"}, &ExecProfile{Project: "b"}, &ExecProfile{Project: "c"}, &ExecProfile{Project: "d"}
	st := &ExecStartup{Profiles: map[string]*ExecProfile{"a": a, "b": b, "c": c, "d": d},
		Sources: map[string]string{"a": "/d/a.json", "b": "/d/b.json", "c": "/d/c.json", "d": "/d/d.json"}}
	errs := []error{nil,
		&ExecProbeError{Project: "b", Failed: []ExecProbeArmFailure{spent}},
		&ExecProbeError{Project: "c", Failed: []ExecProbeArmFailure{spent, refused}},
		errors.New("broker: exec probe for profile \"d\": scratch episode: exists")}
	if err := st.applyProbe([]*ExecProfile{a, b, c, d}, 0, errs); err != nil {
		t.Fatalf("applyProbe = %v", err)
	}
	if len(st.Profiles) != 1 || st.Profiles["a"] == nil || st.Only != nil {
		t.Fatalf("served = %v (Only %v), want only a, and no default beside left-out profiles", st.Profiles, st.Only)
	}
	got := map[string]ExecProfileExclusion{}
	for _, x := range st.Excluded {
		got[x.Project] = x
	}
	if len(got) != 3 || !got["b"].OutOfTime || got["c"].OutOfTime || got["d"].OutOfTime ||
		!strings.Contains(got["b"].Reason, "the startup budget is spent") {
		t.Fatalf("excluded = %+v, want b flagged (its reason kept), c and d not", st.Excluded)
	}

	// A flagged --exec-profile whose failures were all recorded late still refuses
	// startup, as any failed --exec-profile does.
	st = &ExecStartup{Profiles: map[string]*ExecProfile{"b": b}, Sources: map[string]string{"b": "/f/b.json"}}
	if err := st.applyProbe([]*ExecProfile{b}, 1, errs[1:2]); !errors.As(err, new(*ExecProbeError)) {
		t.Fatalf("a late-recorded --exec-profile = %v, want the probe refusal", err)
	}
}

// A late genuine refusal must retain its reason and exclusion. The shared
// operator text reports observations, never a proven timeout/load cause.
func TestExecProbeGenuineRefusalRecordedAfterDeadline(t *testing.T) {
	pp := &execProber{deadline: time.Unix(1, 0)}
	pp.fail(ExecProbeArmWriteOutside, "genuine refusal: sibling write succeeded")
	e := &ExecProbeError{Project: "late", Failed: pp.failed}
	if !e.OutOfTime() {
		t.Fatal("late genuine refusal should permit OutOfTime=true")
	}
	p := &ExecProfile{Project: "late"}
	st := &ExecStartup{Profiles: map[string]*ExecProfile{"late": p}, Sources: map[string]string{"late": "/p/late.json"}}
	if err := st.applyProbe([]*ExecProfile{p}, 0, []error{e}); err != nil {
		t.Fatal(err)
	}
	if len(st.Profiles) != 0 || st.Only != nil || len(st.Excluded) != 1 || !st.Excluded[0].OutOfTime ||
		!strings.Contains(st.Excluded[0].Reason, "genuine refusal: sibling write succeeded") {
		t.Fatalf("late refusal lost exclusion/reason: %+v", st)
	}
	const want = "Every failed arm was recorded at or after the startup deadline; startup timeout may have contributed."
	if ExecProbeDeadlineObservation != want {
		t.Fatalf("diagnostic must report observations only: %q", ExecProbeDeadlineObservation)
	}
}
