package broker

// Row 189 M2 (D-189-2 A): `serve --exec-profiles-dir` loads every
// <project>.json, skips *.proposed.json (D-189-5: inert until renamed), and
// isolates a profile that fails to load or probe, so one broken profile
// never takes down every other project's episodes.

import (
	"errors"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// writeDirProfile writes <dir>/<file> holding a minimal loadable profile for
// project (no read_roots or cache seeds, which need host paths).
func writeDirProfile(t *testing.T, dir, file, project, bin string) string {
	t.Helper()
	data := profileJSON(t, bin, func(m map[string]any) {
		m["project"] = project
		delete(m, "read_roots")
		delete(m, "caches")
	})
	path := filepath.Join(dir, file)
	if err := os.WriteFile(path, data, 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

func dirFixture(t *testing.T) (dir, bin string, h ExecHostPaths, forbidden []namedRoot) {
	t.Helper()
	base, h := hostPathsFixture(t)
	dir = filepath.Join(base, "profiles")
	mustMkdir(t, dir)
	forbidden = []namedRoot{{"the state dir", h.StateDir}, {"the workspace root", h.WorkspaceRoot}}
	return dir, fakeToolDir(t, "go"), h, forbidden
}

func excludedNames(st *ExecStartup) map[string]string {
	out := map[string]string{}
	for _, x := range st.Excluded {
		out[x.Project] = x.Reason
	}
	return out
}

// The loader keeps good.json, never reads ready.proposed.json, and leaves
// out broken.json (does not parse) and misnamed.json (its project is not its
// file name), each with a reason. Mutants: drop the *.proposed.json skip and
// "ready" loads; return the first load error and good never loads.
func TestExecProfilesDirSkipsProposedAndIsolatesABrokenProfile(t *testing.T) {
	dir, bin, h, forbidden := dirFixture(t)
	writeDirProfile(t, dir, "good.json", "good", bin)
	writeDirProfile(t, dir, "ready.proposed.json", "ready", bin)
	writeDirProfile(t, dir, "misnamed.json", "other", bin)
	if err := os.WriteFile(filepath.Join(dir, "broken.json"), []byte(`{"profile": "world/exec-profile/v1", "project": `), 0o600); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "notes.txt"), []byte("not a profile"), 0o600); err != nil {
		t.Fatal(err)
	}
	st := &ExecStartup{Profiles: map[string]*ExecProfile{}, Sources: map[string]string{}}
	loaded, err := loadExecProfilesDir(dir, h, forbidden, st, nil)
	if err != nil {
		t.Fatalf("one broken profile refused the whole directory: %v", err)
	}
	if len(loaded) != 1 || loaded[0].Project != "good" || st.Profiles["good"] == nil {
		t.Fatalf("loaded %v, want only good", loaded)
	}
	if st.Sources["good"] != filepath.Join(dir, "good.json") && !strings.HasSuffix(st.Sources["good"], "good.json") {
		t.Fatalf("source of good = %q", st.Sources["good"])
	}
	ex := excludedNames(st)
	if len(ex) != 2 || ex["broken"] == "" || !strings.Contains(ex["misnamed"], "not its file name") {
		t.Fatalf("excluded = %v, want broken and misnamed with reasons", ex)
	}
	if _, ok := st.Profiles["ready"]; ok || ex["ready"] != "" {
		t.Fatal("a *.proposed.json was read: it is inert until the operator renames it")
	}

	// Control: once renamed, the proposal loads.
	if err := os.Rename(filepath.Join(dir, "ready.proposed.json"), filepath.Join(dir, "ready.json")); err != nil {
		t.Fatal(err)
	}
	st = &ExecStartup{Profiles: map[string]*ExecProfile{}, Sources: map[string]string{}}
	if _, err := loadExecProfilesDir(dir, h, forbidden, st, nil); err != nil || st.Profiles["ready"] == nil {
		t.Fatalf("the renamed profile did not load: %v %v", err, st.Profiles)
	}
}

// A directory file that is also a --exec-profile is the same profile, not a
// duplicate; a different file claiming a flagged project is left out.
func TestExecProfilesDirBesideFlaggedProfiles(t *testing.T) {
	dir, bin, h, forbidden := dirFixture(t)
	same := writeDirProfile(t, dir, "good.json", "good", bin)
	other := t.TempDir()
	flaggedOther := writeDirProfile(t, other, "dup.json", "dup", bin)
	writeDirProfile(t, dir, "dup.json", "dup", bin)
	st := &ExecStartup{Profiles: map[string]*ExecProfile{}, Sources: map[string]string{}}
	for _, f := range []string{same, flaggedOther} {
		p, err := loadExecProfileFile(f, h, forbidden)
		if err != nil {
			t.Fatal(err)
		}
		st.Profiles[p.Project], st.Sources[p.Project] = p, realpathOr(f)
	}
	loaded, err := loadExecProfilesDir(dir, h, forbidden, st, map[string]string{realpathOr(same): same, realpathOr(flaggedOther): flaggedOther})
	if err != nil || len(loaded) != 0 {
		t.Fatalf("loaded %v, %v; want nothing new", loaded, err)
	}
	ex := excludedNames(st)
	if len(ex) != 1 || !strings.Contains(ex["dup"], "already configured") {
		t.Fatalf("excluded = %v, want dup only, as already configured", ex)
	}
}

// An explicit directory that does not exist, or lies inside the workspace
// root, refuses startup: that is the operator's configuration, not one
// profile's fault.
func TestExecProfilesDirItselfIsChecked(t *testing.T) {
	_, _, h, forbidden := dirFixture(t)
	st := &ExecStartup{Profiles: map[string]*ExecProfile{}, Sources: map[string]string{}}
	if _, err := loadExecProfilesDir(filepath.Join(h.StateDir, "nope"), h, forbidden, st, nil); err == nil {
		t.Fatal("a missing --exec-profiles-dir was accepted")
	}
	inside := filepath.Join(h.WorkspaceRoot, "profiles")
	mustMkdir(t, inside)
	var serr *ExecStartupError
	if _, err := loadExecProfilesDir(inside, h, forbidden, st, nil); !errors.As(err, &serr) || serr.Flag != "--exec-profiles-dir" {
		t.Fatalf("a profiles dir inside the workspace root: %v", err)
	}
}

// A directory profile that fails the startup probe is left out with the
// probe's reason and the rest serve; a failed --exec-profile still refuses
// startup. With a sibling left out, the survivor is NOT every episode's
// default. Mutant: make applyProbe return the joined error for every
// profile and the first arm fails.
func TestExecProfilesDirProbeFailureIsolatesOneProfile(t *testing.T) {
	a, b, c := &ExecProfile{Project: "a"}, &ExecProfile{Project: "b"}, &ExecProfile{Project: "c"}
	newSt := func() *ExecStartup {
		return &ExecStartup{Profiles: map[string]*ExecProfile{"a": a, "b": b, "c": c},
			Sources: map[string]string{"a": "/f/a.json", "b": "/d/b.json", "c": "/d/c.json"}}
	}
	probeErr := &ExecProbeError{Project: "b", Failed: []ExecProbeArmFailure{{Arm: ExecProbeArmToolchain, Why: "exit 3"}}}

	st := newSt()
	if err := st.applyProbe([]*ExecProfile{a, b, c}, 1, []error{nil, probeErr, nil}); err != nil {
		t.Fatalf("a failed directory profile refused startup: %v", err)
	}
	if st.Profiles["a"] == nil || st.Profiles["c"] == nil || st.Profiles["b"] != nil {
		t.Fatalf("profiles after the probe = %v, want a and c", st.Profiles)
	}
	if ex := excludedNames(st); !strings.Contains(ex["b"], ExecProbeArmToolchain) || len(st.Excluded) != 1 || st.Excluded[0].File != "/d/b.json" {
		t.Fatalf("excluded = %+v, want b with the probe's arm", st.Excluded)
	}
	// Control: the same failure on a --exec-profile refuses startup.
	st = newSt()
	if err := st.applyProbe([]*ExecProfile{b, a, c}, 1, []error{probeErr, nil, nil}); !errors.As(err, new(*ExecProbeError)) {
		t.Fatalf("a failed --exec-profile: %v, want the probe refusal", err)
	}
	// One survivor beside a left-out sibling is no default.
	st = &ExecStartup{Profiles: map[string]*ExecProfile{"a": a, "b": b}, Sources: map[string]string{}}
	if err := st.applyProbe([]*ExecProfile{a, b}, 0, []error{nil, probeErr}); err != nil || st.Only != nil {
		t.Fatalf("Only = %v (%v), want nil beside a left-out profile", st.Only, err)
	}
	st = &ExecStartup{Profiles: map[string]*ExecProfile{"a": a}, Sources: map[string]string{}}
	if err := st.applyProbe([]*ExecProfile{a}, 0, []error{nil}); err != nil || st.Only != a {
		t.Fatalf("Only = %v (%v), want the one profile", st.Only, err)
	}
}
