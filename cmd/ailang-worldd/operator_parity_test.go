package main

// Row 243 V4 S1 (design §4.1, W23, W29): host/verbs' admitted operator subset is joined to the CLI's
// own flag sets and identifier grammars. The flag sets are the ones runEpisodeStart and
// runEpisodeFinishWith build, seen through episodeFlagSetSeen on a --help run, which returns before
// any effect.

import (
	"bytes"
	"context"
	"flag"
	"reflect"
	"sort"
	"strings"
	"testing"

	"github.com/sunholo-data/ailang-world/host/verbs"
)

// cliEpisodeFlagSet runs the CLI's own `episode <verb> --help` and returns the flag set it built.
func cliEpisodeFlagSet(t *testing.T, verb string) *flag.FlagSet {
	t.Helper()
	t.Setenv("HOME", t.TempDir())
	var seen []*flag.FlagSet
	episodeFlagSetSeen = func(fs *flag.FlagSet) { seen = append(seen, fs) }
	defer func() { episodeFlagSetSeen = nil }()
	var out, errOut bytes.Buffer
	rc := -1
	switch verb {
	case "start":
		rc = runEpisodeStart([]string{"--help"}, &out, &errOut, sessionEnv{})
	case "finish":
		rc = runEpisodeFinishWith([]string{"--help"}, &out, &errOut, sessionEnv{}, nil)
	}
	if rc != exitOK || len(seen) != 1 || out.Len() == 0 || errOut.Len() != 0 {
		t.Fatalf("episode %s --help: rc %d, %d flag sets seen, stderr %q", verb, rc, len(seen), errOut.String())
	}
	return seen[0]
}

// flagKinds is a flag set's names, each true when the flag is a boolean.
func flagKinds(fs *flag.FlagSet) map[string]bool {
	m := map[string]bool{}
	fs.VisitAll(func(f *flag.Flag) {
		b, ok := f.Value.(interface{ IsBoolFlag() bool })
		m[f.Name] = ok && b.IsBoolFlag()
	})
	return m
}

func mirrorKinds(fl []verbs.OperatorFlag) (m map[string]bool, admitted []string) {
	m = map[string]bool{}
	for _, f := range fl {
		m[f.Name] = f.Bool
		if f.Admit {
			admitted = append(admitted, f.Name)
		}
	}
	sort.Strings(admitted)
	return m, admitted
}

// TestOperatorFlagParity: host/verbs' mirror is exactly the CLI's flag set — every name, each kind —
// so the CLI adding, renaming or removing a flag, or the subset naming one the CLI lacks, fails here.
func TestOperatorFlagParity(t *testing.T) {
	for _, c := range []struct {
		verb   string
		mirror []verbs.OperatorFlag
		n      int
		admit  []string
	}{
		{"start", verbs.EpisodeStartFlags, 23, []string{"headless", "launch", "preset"}},
		{"finish", verbs.EpisodeFinishFlags, 9, []string{"confirm-on-surface", "preset", "ready"}},
	} {
		cli := flagKinds(cliEpisodeFlagSet(t, c.verb))
		mir, admitted := mirrorKinds(c.mirror)
		if len(cli) != c.n || len(c.mirror) != len(mir) || !reflect.DeepEqual(cli, mir) {
			t.Errorf("episode %s: the CLI's %d flags %v; host/verbs' %d (%d entries) %v", c.verb, len(cli), cli, len(mir), len(c.mirror), mir)
		}
		if !reflect.DeepEqual(admitted, c.admit) {
			t.Errorf("episode %s: the subset admits %v, want %v (design §4.1)", c.verb, admitted, c.admit)
		}
		// Controls: one flag more, one fewer, one renamed, one kind flipped — each differs.
		for i, edit := range []func(m map[string]bool){
			func(m map[string]bool) { m["zzz"] = true },
			func(m map[string]bool) { delete(m, "preset") },
			func(m map[string]bool) { delete(m, "preset"); m["presets"] = false },
			func(m map[string]bool) { m["preset"] = true },
		} {
			m := map[string]bool{}
			for k, v := range cli {
				m[k] = v
			}
			if edit(m); reflect.DeepEqual(m, mir) {
				t.Errorf("episode %s: control %d passed the parity check", c.verb, i)
			}
		}
	}
}

// TestOperatorParsesLikeTheCLI: each admitted spelling, parsed by the CLI's own flag set with the same
// ParseInterleaved, gives the episode and flag values of the admitted tuple, and so does the tuple's
// canonical line; a line the CLI's parse rejects is refused with the CLI's own parse error.
func TestOperatorParsesLikeTheCLI(t *testing.T) {
	yes := verbs.OperatorPolicy{Start: "A", Finish: "A", HandoffProven: true}
	cliParse := func(verb string, words []string) (map[string]string, []string, error) {
		fs := cliEpisodeFlagSet(t, verb)
		pos, err := verbs.ParseInterleaved(fs, words)
		vals := map[string]string{}
		for _, n := range []string{"preset", "launch", "headless", "ready", "confirm-on-surface"} {
			if f := fs.Lookup(n); f != nil {
				vals[n] = f.Value.String()
			}
		}
		return vals, pos, err
	}
	for _, line := range []string{
		"episode start w1 --preset p --launch --headless", "episode start --headless --launch w1 --preset=p",
		`episode start -preset "p" -launch=true -headless=1 w1`, "episode start w1 --preset p --launch",
		"episode finish w1 --preset p", "episode finish --confirm-on-surface --preset=p w1 --ready",
		"episode finish w1 --ready=T --preset 'p'",
	} {
		argv, err := verbs.Split(line)
		if err != nil {
			t.Fatal(err)
		}
		a := verbs.AdmitOperator(argv, yes)
		if a.Class != verbs.OperatorReview {
			t.Fatalf("%q: %q %q", line, a.Class, a.Reason)
		}
		in := a.Intent
		want := map[string]string{"preset": in.Preset}
		if in.Kind == "start" {
			want["launch"], want["headless"] = "true", map[bool]string{true: "true", false: "false"}[in.Headless]
		} else {
			want["ready"], want["confirm-on-surface"] = map[bool]string{true: "true", false: "false"}[in.Ready], map[bool]string{true: "true", false: "false"}[in.ConfirmOnSurface]
		}
		for _, words := range [][]string{argv[2:], in.Argv[3:]} {
			vals, pos, err := cliParse(in.Kind, words)
			if err != nil || !reflect.DeepEqual(pos, []string{in.Episode}) || !reflect.DeepEqual(vals, want) {
				t.Errorf("%q: the CLI's parse of %q is %v %q (%v), the subset's tuple %v %q", line, words, vals, pos, err, want, in.Episode)
			}
		}
	}
	for _, line := range []string{"episode start w1 --Preset p --launch", "episode start w1 --launch --preset", "episode start w1 --preset p --launch=maybe",
		"episode finish w1 --preset p --ready=x", "episode finish w1 --dryrun --preset p"} {
		argv, _ := verbs.Split(line)
		_, _, cliErr := cliParse(argv[1], argv[2:])
		a := verbs.AdmitOperator(argv, yes)
		if cliErr == nil || a.Class != verbs.OperatorUsage || a.Intent != nil || !strings.HasSuffix(a.Reason, ": "+cliErr.Error()) {
			t.Errorf("%q: the CLI's parse says %v; the subset %q %q", line, cliErr, a.Class, a.Reason)
		}
	}
}

// TestOperatorIdentifierGrammarsMatchTheCLI: the subset's episode grammar decides exactly as the CLI's
// on every boundary class, and its preset cap is the CLI's — where loadLaunchPreset itself draws it.
func TestOperatorIdentifierGrammarsMatchTheCLI(t *testing.T) {
	n := 0
	for _, id := range []string{"", "a", "0", "-a", "a-", "A", "aB", "a_b", "a.b", "é", "a b", "a\n", "\na", "a\nb", strings.Repeat("a", 51),
		strings.Repeat("a", 52), strings.Repeat("a", 64), strings.Repeat("a", 65), "a" + strings.Repeat("-", 63), "a" + strings.Repeat("-", 64)} {
		if verbs.EpisodeIDPattern.MatchString(id) != workspaceEpisodePattern.MatchString(id) {
			t.Errorf("%q: host/verbs %v, the CLI %v", id, verbs.EpisodeIDPattern.MatchString(id), workspaceEpisodePattern.MatchString(id))
		}
		if workspaceEpisodePattern.MatchString(id) {
			n++
		}
	}
	if n != 7 {
		t.Errorf("%d of the corpus match the CLI's grammar, want 7", n)
	}
	if verbs.MaxPresetName != maxPresetNameLen {
		t.Errorf("host/verbs' preset cap %d, the CLI's %d", verbs.MaxPresetName, maxPresetNameLen)
	}
	dir := t.TempDir()
	yes := verbs.OperatorPolicy{Start: "A", Finish: "A"}
	for _, c := range []struct {
		n       int
		invalid bool
	}{{maxPresetNameLen, false}, {maxPresetNameLen + 1, true}, {51, false}, {52, true}} {
		name := strings.Repeat("p", c.n)
		_, err := loadLaunchPreset(context.Background(), name, dir, dir, dir)
		a := verbs.AdmitOperator([]string{"episode", "finish", "w1", "--preset", name}, yes)
		if err == nil || strings.Contains(err.Error(), "invalid preset name") != c.invalid || (a.Class == verbs.OperatorReview) == c.invalid {
			t.Errorf("a %d-character preset: the CLI %v; the subset %q %q", c.n, err, a.Class, a.Reason)
		}
	}
}
