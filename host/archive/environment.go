package archive

import (
	"context"
	"errors"
	"fmt"
)

// EnvironmentFailureLabel prefixes every message AttributeFailure builds for an
// instrument failure. It is the single string a reader (or a grep over a CI log)
// keys on to separate "the rig could not run the measurement" from "the code
// under test is wrong", so it is a constant rather than an inline literal.
const EnvironmentFailureLabel = "ENVIRONMENT FAILURE (rig, not a replay defect)"

// environmentGuidance is the actionable half of an instrument-failure message:
// the mechanism, the numbers that pin it, and the one command that discriminates.
// It is deliberately verbose. The failure it explains is rare, arrives inside an
// unrelated test's output, and its whole cost is the hour someone spends
// attributing it to the diff in front of them.
//
// MEASURED on the rig 2026-09-03 (mission iteration 153, row 58), pinned
// AILANG v0.30.0 on Darwin/arm64:
//
//   - one exec of the pinned binary at its ESTABLISHED path: 47-52 ms (n=5)
//   - one FIRST exec of a freshly-written copy at a new path: 1211-1294 ms (n=5)
//   - 8 such first execs started CONCURRENTLY: 1289, 2332, 3396, 4436, 5531,
//     6747, 7788, 8871 ms
//   - 12 concurrently: 1255 ... 13691 ms, still linear at ~1.13 s per exec
//
// The series is linear, not parallel: macOS assesses a never-before-executed
// binary's provenance under a global serializing lock, so N concurrent first
// execs cost N x 1.13 s for the LAST one to return. probeTimeout is a
// PER-PROBE wall-clock bound, so it is crossed at N >= 9 regardless of how fast
// any single probe is -- which is why the same suite passes when a package runs
// alone and reds when `go test ./...` schedules several archiving packages into
// one window.
//
// Two candidate causes are REFUTED by the same measurements rather than argued
// away. CPU contention is not it: three first execs under 16 busy spinners took
// 1322/1373/1385 ms, barely above the unloaded figure. Per-invocation Observatory
// retention cleanup over the (then 553 MB) database is not it either: the warm
// arm above pays that cleanup and still returns in ~50 ms.
const environmentGuidance = "the archived interpreter's `--version` probe exceeded its wall-clock bound.\n" +
	"On macOS the FIRST exec of a freshly-written path pays a provenance assessment that is\n" +
	"globally SERIALIZED (~1.13 s each, measured), so parallel test packages that each archive\n" +
	"the pinned interpreter queue behind one another and the Nth crosses the bound at N >= 9.\n" +
	"This is a property of the rig, not of the code under test. DISCRIMINATE BEFORE ATTRIBUTING:\n" +
	"re-run this package ALONE (`go test ./<pkg> -count=1 -p 1`); if it passes, the bound was\n" +
	"queued behind sibling packages and nothing here is broken."

// EnvironmentFailure reports whether err is an INSTRUMENT failure -- the
// interpreter version probe ran out of wall-clock -- as opposed to a replay or
// correctness defect.
//
// Both conjuncts are load-bearing and neither is sufficient. KindExecFailure
// alone covers a genuinely broken interpreter (not executable, wrong platform,
// non-zero exit), which is a REAL defect and must keep its own attribution.
// context.DeadlineExceeded alone can arrive from any other bounded operation in
// a caller's stack and says nothing about the probe. Only their conjunction
// identifies the probe deadline, which is the one archive failure this repo has
// measured to be caused by the machine rather than by the tree.
//
// It never suppresses a failure. Callers use it to LABEL one: a test that hits
// this still fails, it just fails saying which of the two things went wrong.
func EnvironmentFailure(err error) bool {
	re, ok := IsReplayError(err)
	if !ok || re.Kind != KindExecFailure {
		return false
	}
	return errors.Is(err, context.DeadlineExceeded)
}

// AttributeFailure renders err as a failure message for step, prefixed with
// EnvironmentFailureLabel and the discriminating guidance when
// EnvironmentFailure(err) holds, and unchanged otherwise.
//
// The default branch is the point: an ordinary failure must read exactly as it
// did before this floor existed, so the floor can only ever ADD attribution and
// can never quietly relabel a defect as somebody else's problem.
func AttributeFailure(step string, err error) string {
	if EnvironmentFailure(err) {
		return fmt.Sprintf("%s\n%s: %v\n%s", EnvironmentFailureLabel, step, err, environmentGuidance)
	}
	return fmt.Sprintf("%s: %v", step, err)
}
