package archive

import (
	"context"
	"errors"
	"fmt"
	"os"
	"runtime"
	"strings"
	"testing"
	"time"
)

// deadlineExecFailure is the exact error shape probeVersion builds when the
// bound is crossed: a KindExecFailure ReplayError whose Err chain reaches
// context.DeadlineExceeded. Constructed here rather than reused from a live
// probe so the UNIT arms below stay fast; the SHIPPED-path arms at the bottom
// of this file are what prove the shape is the real one (rule 3k -- a test that
// rebuilds the artifact by a second route verifies your arithmetic, never your
// artifact).
func deadlineExecFailure() error {
	return &ReplayError{
		Kind:   KindExecFailure,
		Path:   "/some/archived/ailang",
		Detail: "cannot obtain --version from archived interpreter",
		Err:    fmt.Errorf("--version: timed out after 10s: %w", context.DeadlineExceeded),
	}
}

// TestEnvironmentFailureRequiresBothConjuncts pins the classifier's whole
// contract: KindExecFailure AND a deadline. Each negative row neuters exactly
// one conjunct, so a classifier that dropped either check would turn that row
// green -- which is the point, because the cost of a false positive here is a
// REAL exec defect wearing an "it's just the rig" label.
func TestEnvironmentFailureRequiresBothConjuncts(t *testing.T) {
	for _, tc := range []struct {
		name string
		err  error
		want bool
	}{
		{
			name: "exec failure that timed out is the instrument",
			err:  deadlineExecFailure(),
			want: true,
		},
		{
			name: "wrapped by a caller, still the instrument",
			err:  fmt.Errorf("newFixtureEnv: %w", deadlineExecFailure()),
			want: true,
		},
		{
			// The load-bearing negative: a genuinely broken interpreter. Same
			// Kind, no deadline. This MUST keep its own attribution.
			name: "exec failure without a deadline is a real defect",
			err: &ReplayError{
				Kind:   KindExecFailure,
				Detail: "cannot obtain --version from archived interpreter",
				Err:    os.ErrPermission,
			},
			want: false,
		},
		{
			name: "a deadline under a different Kind is not the probe",
			err: &ReplayError{
				Kind:   KindHashMismatch,
				Detail: "existing sidecar records a different hash",
				Err:    fmt.Errorf("wrapped: %w", context.DeadlineExceeded),
			},
			want: false,
		},
		{
			name: "a bare deadline from somewhere else in the stack",
			err:  context.DeadlineExceeded,
			want: false,
		},
		{
			name: "an exec failure wrapping nothing",
			err:  &ReplayError{Kind: KindExecFailure, Detail: "no cause recorded"},
			want: false,
		},
		{
			name: "no error at all",
			err:  nil,
			want: false,
		},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := EnvironmentFailure(tc.err); got != tc.want {
				t.Errorf("EnvironmentFailure(%v) = %v, want %v", tc.err, got, tc.want)
			}
		})
	}
}

// TestAttributeFailureLabelsOnlyTheInstrument asserts both branches. The
// DEFAULT branch is the one worth pinning: an ordinary failure must render
// exactly as `step: err` and must not acquire the label, so this floor can only
// add attribution and can never relabel a defect as the rig's fault.
func TestAttributeFailureLabelsOnlyTheInstrument(t *testing.T) {
	env := AttributeFailure("archive pinned interpreter", deadlineExecFailure())
	if !strings.Contains(env, EnvironmentFailureLabel) {
		t.Errorf("instrument message does not carry the label:\n%s", env)
	}
	if !strings.Contains(env, "archive pinned interpreter") {
		t.Errorf("instrument message dropped the step:\n%s", env)
	}
	if !strings.Contains(env, "re-run this package ALONE") {
		t.Errorf("instrument message carries no discriminating command:\n%s", env)
	}

	defect := &ReplayError{Kind: KindExecFailure, Detail: "boom", Err: os.ErrPermission}
	got := AttributeFailure("archive pinned interpreter", defect)
	want := fmt.Sprintf("archive pinned interpreter: %v", defect)
	if got != want {
		t.Errorf("defect message = %q, want %q", got, want)
	}
	if strings.Contains(got, EnvironmentFailureLabel) {
		t.Errorf("a real defect acquired the ENVIRONMENT label:\n%s", got)
	}
}

// TestEnvironmentFailureClassifiesTheSHIPPEDProbeDeadline runs the real
// Archive() against a blocking interpreter under a shrunk probeTimeout and
// classifies the error IT returns. Without this arm every assertion above is
// about an error literal typed in this file; with it, the classifier is pinned
// to the shape probeVersion actually produces, so a change to that wrapping
// (dropping the %w, replacing the Kind) reds here rather than silently
// disabling the floor at every call site.
func TestEnvironmentFailureClassifiesTheSHIPPEDProbeDeadline(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell-script fake interpreter is POSIX-only")
	}
	dir := t.TempDir()
	execPath, _ := fakeInterpreterBlocking(t, dir, "environment-floor")

	a := New(storeDBPath(t))
	a.probeTimeout = 1 * time.Second

	_, err := a.Archive(execPath)
	if err == nil {
		t.Fatal("Archive against a blocking interpreter returned nil error")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("fixture did not produce a deadline: %v", err)
	}
	if !EnvironmentFailure(err) {
		t.Errorf("EnvironmentFailure = false for the shipped probe deadline: %v", err)
	}
	if msg := AttributeFailure("archive pinned interpreter", err); !strings.Contains(msg, EnvironmentFailureLabel) {
		t.Errorf("shipped deadline was not labelled:\n%s", msg)
	}
}

// TestEnvironmentFailureRejectsTheSHIPPEDExecDefect is the negative control for
// the arm above, and it is the one that makes this floor safe to install: a
// REAL exec failure returned by the same shipped code path -- an archived
// interpreter stripped of its execute bit -- must NOT be classified as the rig.
// Both arms run through Archive(); only the cause differs.
func TestEnvironmentFailureRejectsTheSHIPPEDExecDefect(t *testing.T) {
	if runtime.GOOS == "windows" {
		t.Skip("shell-script fake interpreter is POSIX-only")
	}
	if os.Geteuid() == 0 {
		t.Skip("root ignores the execute permission bit, so the fixture cannot fail")
	}
	dir := t.TempDir()
	// "v1\n" fails versionPrefix, so the idempotent re-archive probes again.
	execPath, _ := fakeInterpreter(t, dir, "environment-floor-defect", "v1\n")

	a := New(storeDBPath(t))
	ref, err := a.Archive(execPath)
	if err != nil {
		t.Fatalf("first Archive: %v", err)
	}
	archived := a.pathFor(ref)
	if err := os.Chmod(archived, 0o444); err != nil {
		t.Fatalf("chmod archived interpreter: %v", err)
	}
	t.Cleanup(func() { _ = os.Chmod(archived, 0o555) })

	_, err = a.Archive(execPath)
	re, ok := IsReplayError(err)
	if !ok {
		t.Fatalf("expected *ReplayError from the healing probe, got %v", err)
	}
	if re.Kind != KindExecFailure {
		t.Fatalf("Kind = %q, want %q -- the arm is not exercising the shared branch", re.Kind, KindExecFailure)
	}
	if EnvironmentFailure(err) {
		t.Errorf("a real exec defect was classified as the rig: %v", err)
	}
	if msg := AttributeFailure("archive pinned interpreter", err); strings.Contains(msg, EnvironmentFailureLabel) {
		t.Errorf("a real exec defect acquired the ENVIRONMENT label:\n%s", msg)
	}
}
