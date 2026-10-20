package broker

import (
	"context"
	"crypto/sha256"
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"math"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"slices"
	"strings"
	"testing"
	"time"
)

type v24Sample struct {
	MS                              map[string]float64 `json:"ms"`
	Bytes, Candidates, Mirror       map[string]int
	Intact                          *bool             `json:"intact"`
	Error                           string            `json:"error"`
	Changes                         []protectedChange `json:"changes"`
	Base, BaseDigest, IgnoredDigest string
	BaseFiles, BaseIgnoredFiles     int
}
type v24Stats struct {
	N             int
	P50, P95, Max float64
}
type v24Summary struct {
	Samples       []v24Sample         `json:"samples"`
	Stats         map[string]v24Stats `json:"per_leg_ms"`
	IgnoredCensus *int                `json:"ignored_census"`
}

func v24Summarize(samples []v24Sample, census *int) v24Summary {
	s, values := v24Summary{samples, map[string]v24Stats{}, census}, map[string][]float64{}
	for _, x := range samples {
		for k, v := range x.MS {
			values[k] = append(values[k], v)
		}
	}
	for k, v := range values {
		slices.Sort(v)
		n := len(v)
		s.Stats[k] = v24Stats{n, v[int(math.Ceil(.5*float64(n)))-1], v[int(math.Ceil(.95*float64(n)))-1], v[n-1]}
	}
	return s
}
func v24Accept(s v24Summary, phaseBoundMS, constructBoundMS int) error {
	for _, leg := range []string{"status", "ignored", "flags", "phase", "construction"} {
		if s.Stats[leg].N == 0 || s.Stats[leg].N != len(s.Samples) {
			return fmt.Errorf("missing %s timing/sample count", leg)
		}
	}
	for _, x := range s.Samples {
		switch {
		case x.Intact == nil:
			return fmt.Errorf("null verdict: %s", x.Error)
		case x.Intact != nil && (!*x.Intact || x.Error != ""):
			return fmt.Errorf("changed/error verdict: %s", x.Error)
		}
	}
	switch {
	case s.IgnoredCensus == nil:
		return fmt.Errorf("missing ignored census")
	case s.Stats["phase"].P95 >= float64(phaseBoundMS):
		return fmt.Errorf("phase p95 >= bound")
	case s.Stats["construction"].Max >= float64(constructBoundMS):
		return fmt.Errorf("construction max >= bound")
	}
	return nil
}
func v24Need(t *testing.T, ok bool, format string, args ...any) {
	if t.Helper(); !ok {
		t.Fatalf(format, args...)
	}
}
func v24Must(t *testing.T, err error) { t.Helper(); v24Need(t, err == nil, "%v", err) }
func TestProtectedV24RealWorktree(t *testing.T) {
	yes, n := true, 1
	good := func() v24Summary {
		return v24Summarize([]v24Sample{{MS: map[string]float64{"status": 1, "ignored": 1, "flags": 1, "phase": 9, "construction": 9}, Intact: &yes}}, &n)
	}
	expect := func(t *testing.T, s v24Summary, want string) {
		err := v24Accept(s, 10, 10)
		v24Need(t, (want == "" && err == nil) || (want != "" && err != nil && strings.Contains(err.Error(), want)), "want %q; got %v", want, err)
	}
	t.Run("reject-null", func(t *testing.T) {
		s := good()
		s.Samples[0].Intact, s.Samples[0].Error = nil, "deadline"
		expect(t, s, "null")
		no := false
		s.Samples[0].Intact, s.Samples[0].Error = &no, ""
		expect(t, s, "changed/error verdict")
		s.Samples[0].Intact, s.Samples[0].Error = &yes, "partial"
		expect(t, s, "changed/error verdict")
	})
	t.Run("reject-over-budget", func(t *testing.T) {
		expect(t, good(), "")
		for _, leg := range []string{"phase", "construction"} {
			for _, ms := range []float64{10, 11} {
				s := good()
				s.Samples[0].MS[leg], s.Stats[leg] = ms, v24Stats{1, ms, ms, ms}
				expect(t, s, leg+map[string]string{"phase": " p95", "construction": " max"}[leg]+" >= bound")
			}
		}
	})
	t.Run("required-legs", func(t *testing.T) {
		expect(t, good(), "")
		for _, leg := range []string{"ignored", "status", "flags"} {
			s := good()
			delete(s.Stats, leg)
			expect(t, s, "missing "+leg+" timing/sample count")
			s = good()
			s.Stats[leg] = v24Stats{}
			expect(t, s, "missing "+leg+" timing/sample count")
		}
		s := good()
		s.IgnoredCensus = nil
		expect(t, s, "missing ignored census")
	})
	t.Run("measurement", v24Measurement)
}

type v24Harness struct {
	t   *testing.T
	git string
	ctx context.Context
}

func (h v24Harness) run(root string, args ...string) []byte {
	ctx, cancel := context.WithTimeout(h.ctx, 30*time.Second)
	defer cancel()
	argv, stderr := append([]string{"-c", "core.fsmonitor=false"}, args...), newBoundedSink(protectedCheckMaxOutput)
	out, err := runBounded(ctx, handlerBounds{execTimeout: 30 * time.Second, maxOutputBytes: 64 << 20}, handlerCommand{path: h.git, args: argv, dir: root, env: gitEnv(), stderr: stderr})
	v24Need(h.t, err == nil, "git %q in %s: %v: %s", argv, root, err, stderr.Bytes())
	return out
}
func v24Filter(canonical []string, leg string, out []byte) (candidates int, kept []string) {
	for _, r := range strings.Split(strings.TrimSuffix(string(out), "\x00"), "\x00") {
		if r == "" || (leg == "status" && strings.HasPrefix(r, "# ")) {
			continue
		}
		n := map[string]int{"ignored": 1, "flags": 2, "ls_tree": 2, "status": map[byte]int{'?': 2, '1': 9, 'u': 11}[r[0]]}[leg]
		p := strings.SplitN(r, map[string]string{"ignored": " ", "flags": " ", "status": " ", "ls_tree": "\t"}[leg], n)[n-1]
		candidates++
		if matchProtected(canonical, p) {
			kept = append(kept, r)
		}
	}
	return
}
func (h v24Harness) measure(root string, canonical []string, argv [][]string) (v24Sample, map[string]struct{}) {
	s := v24Sample{MS: map[string]float64{}, Bytes: map[string]int{}, Candidates: map[string]int{}, Mirror: map[string]int{}}
	tick := func(leg string, f func()) {
		if leg == "construction" || leg == "phase" {
			var cancel context.CancelFunc
			h.ctx, cancel = context.WithTimeout(context.Background(), 30*time.Second)
			defer cancel()
		}
		start := time.Now()
		f()
		s.MS[leg] = float64(time.Since(start)) / float64(time.Millisecond)
	}
	preflight := func() { v24Must(h.t, protectedGitPreflight(h.ctx, h.git, root)) }
	snapshot := map[string]struct{}{}
	tick("construction", func() {
		tick("construct_preflight", preflight)
		tick("rev_parse", func() { s.Base = strings.TrimSpace(string(h.run(root, "rev-parse", "HEAD"))) })
		tick("ls_tree", func() {
			out := h.run(root, "ls-tree", "-r", "-z", "--full-tree", s.Base)
			var kept []string
			s.Candidates["ls_tree"], kept = v24Filter(canonical, "ls_tree", out)
			s.Bytes["ls_tree"], s.Mirror["ls_tree"], s.BaseFiles = len(out), len(kept), len(kept)
			s.BaseDigest = fmt.Sprintf("%x", sha256.Sum256([]byte(strings.Join(append(kept, ""), "\x00"))))
		})
		tick("construct_ignored", func() {
			out := h.run(root, argv[1]...)
			s.Bytes["construct_ignored"] = len(out)
			var kept []string
			s.Candidates["construct_ignored"], kept = v24Filter(canonical, "ignored", out)
			s.Mirror["construct_ignored"] = len(kept)
			for _, p := range kept {
				snapshot[foldProtectedKey(p)] = struct{}{}
			}
			keys := slices.Sorted(maps.Keys(snapshot))
			s.IgnoredDigest, s.BaseIgnoredFiles = fmt.Sprintf("%x", sha256.Sum256([]byte(strings.Join(append(keys, ""), "\n")))), len(snapshot)
		})
	})
	var streams [3][]byte
	tick("phase", func() {
		tick("preflight", preflight)
		for i, leg := range []string{"status", "ignored", "flags"} {
			tick(leg, func() { streams[i] = h.run(root, argv[i]...); s.Bytes[leg] = len(streams[i]) })
		}
		tick("parse", func() {
			head, _, err := parseProtectedCandidates(canonical, streams[0], streams[1], streams[2], snapshot)
			v24Need(h.t, err == nil && head == s.Base, "parse: HEAD=%q error=%v", head, err)
		})
	})
	for i, leg := range []string{"status", "ignored", "flags"} {
		count, kept := v24Filter(canonical, leg, streams[i])
		s.Candidates[leg], s.Mirror[leg] = count, len(kept)
	}
	r := h.check(root, s.Base, canonical, snapshot)
	s.MS["checker_ms"], s.Intact, s.Changes, s.Error = r.MS["checker_ms"], r.Intact, r.Changes, r.Error
	return s, snapshot
}
func (h v24Harness) check(root, base string, canonical []string, snapshot map[string]struct{}) v24Sample {
	start := time.Now()
	r := checkProtected(context.Background(), h.git, root, base, canonical, snapshot)
	return v24Sample{MS: map[string]float64{"checker_ms": float64(time.Since(start)) / float64(time.Millisecond)}, Intact: r.intact, Changes: r.changes, Error: map[bool]string{true: fmt.Sprint(r.err)}[r.err != nil]}
}
func (h v24Harness) sampleInWorktree(root, source, cold string, deps, canonical []string, argv [][]string) (v24Sample, map[string]struct{}) {
	if cold == "" {
		return h.measure(root, canonical, argv)
	}
	_, err := os.Lstat(cold)
	v24Need(h.t, os.IsNotExist(err), "cold path exists or inaccessible: %s: %v", cold, err)
	moved := []string{}
	defer func() {
		var restoreErr error
		for i := len(moved) - 1; i >= 0; i-- {
			restoreErr = errors.Join(restoreErr, os.Rename(filepath.Join(cold, moved[i]), filepath.Join(root, moved[i])))
		}
		v24Need(h.t, restoreErr == nil, "RESTORE: %v", restoreErr)
		h.run(source, "worktree", "remove", "--force", cold)
	}()
	h.run(source, "worktree", "add", "--detach", cold, "HEAD")
	for _, p := range deps {
		dest := filepath.Join(cold, p)
		mustMkdir(h.t, filepath.Dir(dest))
		v24Must(h.t, os.Rename(filepath.Join(root, p), dest))
		moved = append(moved, p)
	}
	return h.measure(cold, canonical, argv)
}
func (h v24Harness) arm(root, probe string, canonical []string, base string, snapshot map[string]struct{}) []v24Sample {
	dir, path := filepath.Dir(filepath.Join(root, probe)), filepath.Join(root, probe)
	v24Must(h.t, os.Mkdir(dir, 0755))
	defer func() { v24Must(h.t, os.RemoveAll(dir)) }()
	mustWrite(h.t, path, "V24 post-snapshot dependency\n")
	v24Need(h.t, strings.TrimSpace(string(h.run(root, "check-ignore", "--", probe))) == probe, "probe is not ignored")
	before := h.check(root, base, canonical, snapshot)
	named := slices.ContainsFunc(before.Changes, func(c protectedChange) bool { return c.Path == probe })
	v24Need(h.t, before.Error == "" && before.Intact != nil && !*before.Intact && named, "post-snapshot created: %+v", before)
	v24Must(h.t, os.Remove(path))
	after := h.check(root, base, canonical, snapshot)
	v24Need(h.t, after.Error == "" && after.Intact != nil && *after.Intact, "post-snapshot removed: %+v", after)
	return []v24Sample{before, after}
}
func v24Measurement(t *testing.T) {
	root := os.Getenv("WORLD_PROTECTED_V24_ROOT")
	if root == "" {
		t.Skip("V24 UNMEASURED: WORLD_PROTECTED_V24_ROOT unset — a skip never satisfies V24")
	}
	git, err := exec.LookPath("git")
	v24Need(t, err == nil, "real git required: %v", err)
	h := v24Harness{t, git, context.Background()}
	report := map[string]any{"git_path": git, "go_version": runtime.Version(), "cold_label": "cold = first Git access on a freshly added worktree (fresh index, no untracked cache); OS page cache NOT purged; construction precedes phase", "install_order": "Dependencies installed before snapshot here; M4-wide ordering unproven; post-snapshot installation reports CHANGED (R-156-6)."}
	defer func() {
		if out := os.Getenv("WORLD_PROTECTED_V24_OUT"); out != "" {
			b, err := json.MarshalIndent(report, "", "  ")
			v24Need(t, err == nil, "JSON: %v", err)
			v24Must(t, os.WriteFile(out, append(b, '\n'), 0644))
		}
	}()
	groups, phaseBound, constructBound := map[string]v24Summary{}, 2500, 3200
	report["groups"], report["git_version"] = groups, strings.TrimSpace(string(h.run(root, "--version")))
	for _, project := range []struct{ name, patterns, probe string }{
		{"sunholo-platform", "test_*.py conftest.py *.test.ts *.test.tsx backend/tests/** cli/tests/**", "backend/.venv/v24_probe/test_v24_probe.py"},
		{"ailang", "*_test.go *.golden tests/** examples/tests/**", "docs/node_modules/v24_probe/v24_probe_test.go"},
	} {
		warm, source := filepath.Join(root, project.name), filepath.Join(root, "src", project.name)
		v24Need(t, strings.Count(string(h.run(source, "worktree", "list", "--porcelain")), "worktree ") == 2, "expected source clone and warm worktree: %s", source)
		deps := strings.Split(strings.TrimSpace(string(mustRead(t, warm+".deps"))), "\n")
		v24Need(t, len(deps) > 0 && deps[0] != "", "empty deps: %s", warm)
		checkDeps := func() {
			for _, p := range deps {
				_, err := os.Lstat(filepath.Join(warm, p))
				v24Need(t, err == nil, "missing dependency %s: %v", p, err)
			}
		}
		checkDeps()
		v24Need(t, len(h.run(warm, "status", "--porcelain")) == 0, "warm worktree dirty: %s", warm)
		canonical, err := ParseProtectedPatterns(strings.Fields(project.patterns))
		v24Need(t, err == nil, "patterns: %v", err)
		argv := [][]string{strings.Fields("status --porcelain=v2 -z --branch --untracked-files=all --ignored=no --no-renames --ignore-submodules=all --"), strings.Fields(protectedIgnoredCommand), strings.Fields("ls-files -v -z --")}
		for i := range argv {
			argv[i] = append(argv[i], gitCandidatePathspecs(canonical)...)
			t.Logf("%s argv: %q", project.name, append([]string{git, "-c", "core.fsmonitor=false"}, argv[i]...))
		}
		census := strings.Count(string(h.run(warm, "ls-files", "-z", "-o", "-i", "--exclude-standard")), "\x00")
		report[project.name] = map[string]any{"canonical": canonical, "argv": argv, "argv_prefix": []string{git, "-c", "core.fsmonitor=false"}, "deps_entries": len(deps), "ignored_files": census, "revision": strings.TrimSpace(string(h.run(warm, "rev-parse", "HEAD")))}
		h.measure(warm, canonical, argv) // One discarded warm-up.
		for _, class := range []string{"warm", "cold"} {
			samples := []v24Sample{}
			for k := 1; k <= 20; k++ {
				cold := map[bool]string{true: fmt.Sprintf("%s-cold-%02d", warm, k)}[class == "cold"]
				s, snapshot := h.sampleInWorktree(warm, source, cold, deps, canonical, argv)
				samples = append(samples, s)
				groups[project.name+"/"+class] = v24Summarize(samples, &census)
				phaseBound, constructBound = max(phaseBound, int(math.Ceil(1.5*s.MS["phase"]))), max(constructBound, int(math.Ceil(1.5*s.MS["construction"])))
				if class == "warm" && k == 20 {
					report[project.name+"/post_snapshot"] = h.arm(warm, project.probe, canonical, s.Base, snapshot)
				}
			}
		}
		checkDeps()
		v24Need(t, strings.Count(string(h.run(source, "worktree", "list", "--porcelain")), "worktree ") == 2, "expected exactly two worktrees for %s", project.name)
		v24Need(t, len(h.run(warm, "status", "--porcelain")) == 0, "warm worktree dirty after restoration: %s", warm)
	}
	report["phase_bound_ms"], report["construct_bound_ms"] = phaseBound, constructBound
	counts := map[string]int{}
	var acceptErr error
	for _, key := range []string{"sunholo-platform/warm", "sunholo-platform/cold", "ailang/warm", "ailang/cold"} {
		s := groups[key]
		for _, x := range s.Samples {
			counts["null"] += map[bool]int{true: 1}[x.Intact == nil]
			counts["error"] += map[bool]int{true: 1}[x.Error != ""]
			counts["changed"] += map[bool]int{true: 1}[x.Intact != nil && !*x.Intact]
		}
		t.Logf("%s n=%d phase(p50/p95/max)=%+v construction=%+v bounds=%d/%d ms", key, len(s.Samples), s.Stats["phase"], s.Stats["construction"], phaseBound, constructBound)
		acceptErr = errors.Join(acceptErr, v24Accept(s, phaseBound, constructBound))
		v24Need(t, len(s.Samples) == 20, "%s: need 20 samples, got %d", key, len(s.Samples))
	}
	report["checker_null"], report["checker_error"], report["checker_changed"], report["checker_bound_ms"] = counts["null"], counts["error"], counts["changed"], float64(protectedCheckTimeout)/float64(time.Millisecond)
	t.Logf("production checker: null=%d error=%d changed=%d", counts["null"], counts["error"], counts["changed"])
	v24Must(t, acceptErr)
}
