package main

import (
	"bytes"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"testing"
)

// ---------------------------------------------------------------------------
// AC21 — all 12 refusal branches, driven, with a per-row positive control
//
// 12, not 14: R-TTY-CHARDEV and R-TTY-SAMEFILE were REMOVED on 2026-10-06
// (attended-steps ergonomics). They refused an operator whose stdin was not the
// controlling terminal; the fence now reads the typed line from /dev/tty itself
// in that case, so stdin's identity no longer decides anything and the only
// terminal refusal left is R-TTY-OPEN. TestRedirectedStdinIsNeverTheConfirmation
// pins what the two rows used to guard: a phrase on a redirected stdin is never
// read.
//
// SAFETY: every row below runs IN PROCESS with an injected ttyProbe and an
// injected getenv. No row performs a network request of any kind: the only row
// that reaches broker.NewRegistryPublishHandler hands it a LOOPBACK registry
// origin, which the PRODUCTION constructor refuses before any transport exists.
// AILANG_REGISTRY_API_KEY is neither read, set, nor named anywhere in this
// package — TestCommandNeverNamesTheCredentialVariable measures that.
//
// THE POSITIVE CONTROL IS THE POINT. Asserting only that a bad invocation exits
// 3 would pass identically if the command refused EVERYTHING — a panic, a
// missing file, a broken build. Each row therefore also runs the IDENTICAL
// invocation with only that row's triggering condition removed, and requires it
// to reach a STRICTLY LATER stage. A refusal that cannot be removed is not a
// fence, it is a wall.
// ---------------------------------------------------------------------------

// fenceStages is the total order of refusal stages a `publish` invocation walks.
// "Later" in AC21 is defined against THIS list, and the list is finer-grained
// than the fence names: the TTY refusal and the two confirmation refusals are
// distinct stages, so a control that merely reached a different reason on
// the same fence cannot be mistaken for progress.
var fenceStages = []string{
	"mode",
	"store",
	"approval",
	"credential",
	"packet",
	"ci",
	"tty/no-controlling-terminal",
	"confirmation/eof",
	"confirmation/mismatch",
	"handler",
}

// multiStageFences are the fences whose REASONS are sequential rather than
// alternative: the TTY refusal and the two confirmation refusals are
// checked in order, so a control that reached a different reason on the same
// fence HAS made progress and must be scored as such. Every other fence's
// reasons are alternatives reached at one position.
var multiStageFences = map[string]bool{"tty": true, "confirmation": true}

// stageOf maps an observed STOP line to its index in fenceStages. It parses the
// line the command actually printed rather than being told what to expect, so a
// row whose command printed nothing scores -1 and fails loudly.
func stageOf(t *testing.T, stderr string) (int, string) {
	t.Helper()
	line := stopLine(stderr)
	if line == "" {
		return -1, ""
	}
	fence, reason := parseStopLine(t, line)
	key := fence
	if multiStageFences[fence] {
		key = fence + "/" + reason
	}
	for i, stage := range fenceStages {
		if stage == key {
			return i, line
		}
	}
	t.Fatalf("observed STOP line %q maps to no known stage (key %q). Either a fence was added "+
		"without a stage, or the STOP contract changed", line, key)
	return -1, line
}

func stopLine(stderr string) string {
	for _, line := range strings.Split(stderr, "\n") {
		if strings.HasPrefix(line, "STOP fence=") {
			return line
		}
	}
	return ""
}

func parseStopLine(t *testing.T, line string) (fence, reason string) {
	t.Helper()
	fields := strings.Fields(line)
	if len(fields) < 2 || fields[0] != "STOP" {
		t.Fatalf("malformed STOP line %q", line)
	}
	fence = strings.TrimPrefix(fields[1], "fence=")
	if len(fields) >= 3 {
		reason = strings.TrimPrefix(fields[2], "reason=")
	}
	return fence, reason
}

// ---------------------------------------------------------------------------
// probes and fixtures
// ---------------------------------------------------------------------------

func commandRepoRoot(t *testing.T) string {
	t.Helper()
	_, file, _, ok := runtime.Caller(0)
	if !ok {
		t.Fatal("cannot locate repository root")
	}
	return filepath.Clean(filepath.Join(filepath.Dir(file), "..", ".."))
}

func statOrFail(t *testing.T, path string) fs.FileInfo {
	t.Helper()
	info, err := os.Stat(path)
	if err != nil {
		t.Fatalf("stat %s: %v", path, err)
	}
	return info
}

// devNullInfo is the MEASURED fact that makes R-TTY-SAMEFILE necessary:
// /dev/null carries the character-device bit, so a naive isatty check admits
// `--live < /dev/null`.
func devNullInfo(t *testing.T) fs.FileInfo {
	t.Helper()
	info := statOrFail(t, os.DevNull)
	if info.Mode()&os.ModeCharDevice == 0 {
		t.Fatalf("instrument failure: %s is not a character device on this platform (mode %v); "+
			"the R-TTY-SAMEFILE row is not exercising what it claims", os.DevNull, info.Mode())
	}
	return info
}

func regularFileInfo(t *testing.T) fs.FileInfo {
	t.Helper()
	path := filepath.Join(t.TempDir(), "not-a-terminal")
	if err := os.WriteFile(path, []byte("x"), 0o600); err != nil {
		t.Fatal(err)
	}
	info := statOrFail(t, path)
	if info.Mode()&os.ModeCharDevice != 0 {
		t.Fatalf("instrument failure: a regular file reports the character-device bit (mode %v)", info.Mode())
	}
	return info
}

// satisfiedProbe is the ONLY probe in this repository that passes the fence. It
// is /dev/null paired with ITSELF, which os.SameFile reports as the same file —
// measured, and the reason the passing branch is reachable with no pty and no
// new dependency. In production the ctty FileInfo comes only from
// os.Open("/dev/tty"), which never resolves to /dev/null.
func satisfiedProbe(t *testing.T) ttyProbe {
	t.Helper()
	info := devNullInfo(t)
	probe := ttyProbe{stdin: info, ctty: info}
	if err := requireControllingTerminal(probe); err != nil {
		t.Fatalf("instrument failure: the satisfied probe does not satisfy the fence: %v", err)
	}
	return probe
}

func noEnv(string) string { return "" }

func ciEnv(name string) string {
	if name == "CI" {
		return "true"
	}
	return ""
}

// goldenCopy writes a copy of the committed golden with one transformation
// applied, so the packet rows drive a MUTATED artifact without touching the
// working tree.
func goldenCopy(t *testing.T, transform func(string) string) string {
	t.Helper()
	src := filepath.Join(commandRepoRoot(t), defaultGolden)
	data, err := os.ReadFile(src)
	if err != nil {
		t.Fatalf("read committed golden: %v", err)
	}
	mutated := transform(string(data))
	if mutated == string(data) {
		t.Fatal("instrument failure: the golden transformation changed nothing, so the row " +
			"would be driving the UNMODIFIED artifact")
	}
	path := filepath.Join(t.TempDir(), "golden.json")
	if err := os.WriteFile(path, []byte(mutated), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

// baseArgs is a fully-formed live publish invocation with every fence
// satisfiable. Each row REMOVES or CORRUPTS exactly one thing.
type invocation struct {
	verb  string
	flags map[string]string
	bools []string
}

func (inv invocation) args(t *testing.T) []string {
	t.Helper()
	args := []string{inv.verb}
	names := make([]string, 0, len(inv.flags))
	for name := range inv.flags {
		names = append(names, name)
	}
	// Deterministic order: a flaky argv would make a flaky fence.
	for _, name := range sortedStrings(names) {
		args = append(args, "--"+name, inv.flags[name])
	}
	for _, name := range inv.bools {
		args = append(args, "--"+name)
	}
	return args
}

func (inv invocation) with(mutate func(*invocation)) invocation {
	clone := invocation{verb: inv.verb, flags: map[string]string{}, bools: append([]string(nil), inv.bools...)}
	for k, v := range inv.flags {
		clone.flags[k] = v
	}
	mutate(&clone)
	return clone
}

func sortedStrings(in []string) []string {
	out := append([]string(nil), in...)
	for i := 1; i < len(out); i++ {
		for j := i; j > 0 && out[j] < out[j-1]; j-- {
			out[j], out[j-1] = out[j-1], out[j]
		}
	}
	return out
}

// liveInvocation is the base every publish row starts from.
//
// --registry-origin is LOOPBACK on purpose: if the whole fence stack failed,
// the PRODUCTION constructor still refuses it and no request can be built.
//
// The scheme is `https` deliberately, and it was MEASURED rather than guessed.
// validatePublishOrigin checks `scheme != "https"` BEFORE it checks loopback,
// so `http://127.0.0.1:1` is refused for its SCHEME and never reaches the
// loopback branch — the row would then be green while proving nothing about
// loopback refusal. `https://127.0.0.1:1` reaches the branch this milestone
// actually depends on, and it is no less loopback for being https.
func liveInvocation(t *testing.T) invocation {
	t.Helper()
	root := commandRepoRoot(t)
	return invocation{
		verb: "publish",
		flags: map[string]string{
			"store":           filepath.Join(t.TempDir(), "world.db"),
			"package-dir":     filepath.Join(root, defaultPackageDir),
			"golden":          filepath.Join(root, defaultGolden),
			"registry-origin": "https://127.0.0.1:1",
			"publisher":       filepath.Join(root, defaultGolden), // any readable file; never executed
			"credential-file": filepath.Join(t.TempDir(), "registry.key"),
			"approval-ref":    "sha256:" + strings.Repeat("ab", 32),
			"now":             "10",
			"expires":         "100",
		},
		bools: []string{"live"},
	}
}

type armResult struct {
	code   int
	stdout string
	stderr string
}

func drive(t *testing.T, inv invocation, stdin string, getenv func(string) string, probe ttyProbe) armResult {
	t.Helper()
	var out, errw bytes.Buffer
	code := run(inv.args(t), strings.NewReader(stdin), &out, &errw,
		environment{getenv: getenv, probe: func() ttyProbe { return probe }})
	return armResult{code: code, stdout: out.String(), stderr: errw.String()}
}

// ---------------------------------------------------------------------------
// the table
// ---------------------------------------------------------------------------

type fenceRow struct {
	// branch is the enumerated refusal branch this row drives.
	branch string
	// wantLine is the EXACT STOP line the row must produce.
	wantLine string
	// trigger builds the invocation, stdin, env and probe that FIRE the branch.
	trigger func(t *testing.T) (invocation, string, func(string) string, ttyProbe)
	// control removes ONLY the triggering condition.
	control func(t *testing.T) (invocation, string, func(string) string, ttyProbe)
	// controlMayExitZero is set for the two READ-ONLY verbs, whose control
	// legitimately completes: that a headless loop CAN run them is the whole
	// asymmetry R-RECONCILE-LIVE-FLAG exists to protect.
	controlMayExitZero bool
}

func fenceRows() []fenceRow {
	openProbe := func(t *testing.T) ttyProbe {
		return ttyProbe{stdin: regularFileInfo(t), cttyErr: fmt.Errorf("device not configured")}
	}
	phrase := attendedPhrase + "\n"

	return []fenceRow{
		{
			branch:   "R-MODE-NONE",
			wantLine: "STOP fence=mode reason=none",
			trigger: func(t *testing.T) (invocation, string, func(string) string, ttyProbe) {
				return liveInvocation(t).with(func(i *invocation) { i.bools = nil }), phrase, noEnv, satisfiedProbe(t)
			},
			control: func(t *testing.T) (invocation, string, func(string) string, ttyProbe) {
				// Only the missing mode is restored; everything else stays.
				return liveInvocation(t).with(func(i *invocation) { i.flags["store"] = "" }), phrase, noEnv, satisfiedProbe(t)
			},
		},
		{
			branch:   "R-MODE-BOTH",
			wantLine: "STOP fence=mode reason=both",
			trigger: func(t *testing.T) (invocation, string, func(string) string, ttyProbe) {
				return liveInvocation(t).with(func(i *invocation) { i.bools = []string{"live", "dry-run"} }),
					phrase, noEnv, satisfiedProbe(t)
			},
			control: func(t *testing.T) (invocation, string, func(string) string, ttyProbe) {
				return liveInvocation(t).with(func(i *invocation) { i.flags["store"] = "" }), phrase, noEnv, satisfiedProbe(t)
			},
		},
		{
			branch:   "R-STORE",
			wantLine: "STOP fence=store reason=absent",
			trigger: func(t *testing.T) (invocation, string, func(string) string, ttyProbe) {
				return liveInvocation(t).with(func(i *invocation) { i.flags["store"] = "" }), phrase, noEnv, satisfiedProbe(t)
			},
			control: func(t *testing.T) (invocation, string, func(string) string, ttyProbe) {
				return liveInvocation(t).with(func(i *invocation) { i.flags["approval-ref"] = "" }),
					phrase, noEnv, satisfiedProbe(t)
			},
		},
		{
			branch:   "R-APPROVAL-ABSENT",
			wantLine: "STOP fence=approval reason=absent",
			trigger: func(t *testing.T) (invocation, string, func(string) string, ttyProbe) {
				return liveInvocation(t).with(func(i *invocation) { i.flags["approval-ref"] = "" }),
					phrase, noEnv, satisfiedProbe(t)
			},
			control: func(t *testing.T) (invocation, string, func(string) string, ttyProbe) {
				return liveInvocation(t).with(func(i *invocation) { i.flags["credential-file"] = "" }),
					phrase, noEnv, satisfiedProbe(t)
			},
		},
		{
			branch:   "R-CRED-FLAG",
			wantLine: "STOP fence=credential reason=absent",
			trigger: func(t *testing.T) (invocation, string, func(string) string, ttyProbe) {
				return liveInvocation(t).with(func(i *invocation) { i.flags["credential-file"] = "" }),
					phrase, noEnv, satisfiedProbe(t)
			},
			control: func(t *testing.T) (invocation, string, func(string) string, ttyProbe) {
				return liveInvocation(t).with(func(i *invocation) {
					i.flags["golden"] = goldenCopy(t, flipOneNibble)
				}), phrase, noEnv, satisfiedProbe(t)
			},
		},
		{
			branch:   "R-PACKET-DRIFT",
			wantLine: "STOP fence=packet reason=drift",
			trigger: func(t *testing.T) (invocation, string, func(string) string, ttyProbe) {
				return liveInvocation(t).with(func(i *invocation) {
					i.flags["golden"] = goldenCopy(t, flipOneNibble)
				}), phrase, noEnv, satisfiedProbe(t)
			},
			control: func(t *testing.T) (invocation, string, func(string) string, ttyProbe) {
				// Golden restored; the CI tripwire is the next stage.
				return liveInvocation(t), phrase, ciEnv, satisfiedProbe(t)
			},
		},
		{
			branch:   "R-PACKET-VERSION",
			wantLine: "STOP fence=packet reason=version",
			trigger: func(t *testing.T) (invocation, string, func(string) string, ttyProbe) {
				return liveInvocation(t).with(func(i *invocation) {
					i.flags["golden"] = goldenCopy(t, func(s string) string {
						return strings.Replace(s, `"version":"`+frozenPackageVersion+`"`, `"version":"9.9.9"`, 1)
					})
				}), phrase, noEnv, satisfiedProbe(t)
			},
			control: func(t *testing.T) (invocation, string, func(string) string, ttyProbe) {
				return liveInvocation(t), phrase, ciEnv, satisfiedProbe(t)
			},
		},
		{
			branch:   "R-CI",
			wantLine: "STOP fence=ci reason=ci",
			trigger: func(t *testing.T) (invocation, string, func(string) string, ttyProbe) {
				return liveInvocation(t), phrase, ciEnv, satisfiedProbe(t)
			},
			control: func(t *testing.T) (invocation, string, func(string) string, ttyProbe) {
				return liveInvocation(t), phrase, noEnv, openProbe(t)
			},
		},
		{
			branch:   "R-TTY-OPEN",
			wantLine: "STOP fence=tty reason=no-controlling-terminal",
			trigger: func(t *testing.T) (invocation, string, func(string) string, ttyProbe) {
				return liveInvocation(t), phrase, noEnv, openProbe(t)
			},
			control: func(t *testing.T) (invocation, string, func(string) string, ttyProbe) {
				// /dev/tty opens; stdin is still a regular file carrying the
				// phrase, which is NOT read: the terminal is, and nothing is
				// typed there, so the next stage (confirmation/eof) is reached.
				return liveInvocation(t), phrase, noEnv, terminalProbe(t, regularFileInfo(t), "")
			},
		},
		{
			branch:   "R-PHRASE-EOF",
			wantLine: "STOP fence=confirmation reason=eof",
			trigger: func(t *testing.T) (invocation, string, func(string) string, ttyProbe) {
				return liveInvocation(t), "", noEnv, satisfiedProbe(t)
			},
			control: func(t *testing.T) (invocation, string, func(string) string, ttyProbe) {
				return liveInvocation(t), "not the phrase\n", noEnv, satisfiedProbe(t)
			},
		},
		{
			branch:   "R-PHRASE",
			wantLine: "STOP fence=confirmation reason=mismatch",
			trigger: func(t *testing.T) (invocation, string, func(string) string, ttyProbe) {
				return liveInvocation(t), "not the phrase\n", noEnv, satisfiedProbe(t)
			},
			control: func(t *testing.T) (invocation, string, func(string) string, ttyProbe) {
				// THE POSITIVE CONTROL OF THE WHOLE FENCE STACK: the exact phrase
				// passes every fence, and the invocation is then stopped by the
				// PRODUCTION constructor refusing a loopback registry origin.
				return liveInvocation(t), phrase, noEnv, satisfiedProbe(t)
			},
		},
		{
			branch:   "R-RECONCILE-LIVE-FLAG",
			wantLine: "STOP fence=mode reason=reconcile-is-read-only",
			trigger: func(t *testing.T) (invocation, string, func(string) string, ttyProbe) {
				inv := liveInvocation(t)
				inv.verb = "reconcile"
				return inv, "", noEnv, satisfiedProbe(t)
			},
			control: func(t *testing.T) (invocation, string, func(string) string, ttyProbe) {
				inv := liveInvocation(t).with(func(i *invocation) { i.bools = nil })
				inv.verb = "reconcile"
				return inv, "", noEnv, openProbe(t)
			},
			controlMayExitZero: true,
		},
		{
			branch:   "R-RECONCILE-LIVE-FLAG (packet arm — the same branch, the other read-only verb)",
			wantLine: "STOP fence=mode reason=packet-is-read-only",
			trigger: func(t *testing.T) (invocation, string, func(string) string, ttyProbe) {
				inv := liveInvocation(t)
				inv.verb = "packet"
				return inv, "", noEnv, satisfiedProbe(t)
			},
			control: func(t *testing.T) (invocation, string, func(string) string, ttyProbe) {
				inv := liveInvocation(t).with(func(i *invocation) { i.bools = nil })
				inv.verb = "packet"
				return inv, "", noEnv, openProbe(t)
			},
			controlMayExitZero: true,
		},
	}
}

func flipOneNibble(golden string) string {
	const marker = `"contentHash":"sha256:`
	idx := strings.Index(golden, marker)
	if idx < 0 {
		return golden
	}
	at := idx + len(marker)
	replacement := byte('0')
	if golden[at] == '0' {
		replacement = '1'
	}
	return golden[:at] + string(replacement) + golden[at+1:]
}

// enumeratedBranches is the frozen list of refusal branches this milestone adds.
// It is stated INDEPENDENTLY of the table so a row deleted from the table reds
// rather than silently reducing coverage.
var enumeratedBranches = []string{
	"R-MODE-NONE", "R-MODE-BOTH", "R-CI",
	"R-TTY-OPEN",
	"R-PHRASE-EOF", "R-PHRASE",
	"R-STORE", "R-APPROVAL-ABSENT", "R-CRED-FLAG",
	"R-PACKET-DRIFT", "R-PACKET-VERSION", "R-RECONCILE-LIVE-FLAG",
}

func TestEveryRefusalBranchStopsWithItsExactLineAndHasAPositiveControl(t *testing.T) {
	rows := fenceRows()

	// ANTI-VACUITY: the table must cover every enumerated branch exactly, or a
	// green here means "the rows that remain still pass".
	covered := map[string]bool{}
	for _, row := range rows {
		covered[strings.Fields(row.branch)[0]] = true
	}
	if len(covered) != len(enumeratedBranches) {
		t.Fatalf("the table drives %d distinct branches, but %d are enumerated: %v",
			len(covered), len(enumeratedBranches), covered)
	}
	for _, branch := range enumeratedBranches {
		if !covered[branch] {
			t.Fatalf("enumerated branch %s has no row", branch)
		}
	}
	t.Logf("AC21: %d rows driving %d enumerated refusal branches", len(rows), len(enumeratedBranches))

	for _, row := range rows {
		row := row
		t.Run(row.branch, func(t *testing.T) {
			inv, stdin, getenv, probe := row.trigger(t)
			got := drive(t, inv, stdin, getenv, probe)

			if got.code != exitStop {
				t.Fatalf("exit code = %d, want %d\nstderr:\n%s", got.code, exitStop, got.stderr)
			}
			line := stopLine(got.stderr)
			if line != row.wantLine {
				t.Fatalf("STOP line = %q, want %q\nstderr:\n%s", line, row.wantLine, got.stderr)
			}
			triggerStage, _ := stageOf(t, got.stderr)

			// THE POSITIVE CONTROL.
			cInv, cStdin, cGetenv, cProbe := row.control(t)
			control := drive(t, cInv, cStdin, cGetenv, cProbe)
			if control.code == exitOK {
				if !row.controlMayExitZero {
					t.Fatalf("the positive control EXITED 0. A row whose only removable condition "+
						"leads straight to success is not a fence in a stack.\nstdout:\n%s", control.stdout)
				}
				t.Logf("control (read-only verb) completed: exit 0")
				return
			}
			controlStage, controlLine := stageOf(t, control.stderr)
			if controlStage < 0 {
				t.Fatalf("the positive control produced NO STOP line (exit %d). Either it crashed or "+
					"the fence stack is not being walked.\nstderr:\n%s", control.code, control.stderr)
			}
			if controlStage <= triggerStage {
				t.Fatalf("the positive control reached stage %q (index %d), which is not strictly later "+
					"than this row's %q (index %d): removing the trigger did not make progress, so the "+
					"row may be passing for an unrelated reason",
					controlLine, controlStage, row.wantLine, triggerStage)
			}
			t.Logf("trigger -> %q (stage %d); control -> %q (stage %d)",
				row.wantLine, triggerStage, controlLine, controlStage)
		})
	}
}

// TestTheFenceStackPassesOnlyForAnAttendedOperator drives requireAttendedOperator
// in BOTH directions at the unit level. Without the passing direction, every
// refusal above would be consistent with a fence that can never be satisfied —
// and a fence that can never be satisfied is indistinguishable from a broken
// command.
func TestTheFenceStackPassesOnlyForAnAttendedOperator(t *testing.T) {
	var out bytes.Buffer
	if err := requireAttendedOperator(strings.NewReader(attendedPhrase+"\n"), &out,
		noEnv, satisfiedProbe(t)); err != nil {
		t.Fatalf("the fence stack refused a satisfied operator: %v", err)
	}
	if !strings.Contains(out.String(), attendedPhrase) {
		t.Fatal("the prompt did not show the operator the phrase they must type")
	}
	// And the phrase without a trailing newline (a terminal that sends the line
	// and closes) must still pass: R-PHRASE-EOF is about an EMPTY read, not
	// about the absence of '\n'.
	out.Reset()
	if err := requireAttendedOperator(strings.NewReader(attendedPhrase), &out,
		noEnv, satisfiedProbe(t)); err != nil {
		t.Fatalf("the fence stack refused a phrase with no trailing newline: %v", err)
	}
}

// TestGithubActionsAlsoTripsTheDeclaredTripwire covers the second CI variable.
// The tripwire is DECLARED as a tripwire, not as the fence: `env -u CI` defeats
// it, which is exactly why the load-bearing layer is the controlling terminal.
func TestGithubActionsAlsoTripsTheDeclaredTripwire(t *testing.T) {
	if len(ciVariables) != 2 {
		t.Fatalf("ciVariables = %v, want exactly the two runner markers", ciVariables)
	}
	for _, name := range ciVariables {
		err := refuseAutomationEnvironment(func(query string) string {
			if query == name {
				return "1"
			}
			return ""
		})
		if err == nil {
			t.Fatalf("%s set did not trip the tripwire", name)
		}
		if err.Fence != fenceCI {
			t.Fatalf("%s tripped fence %q, want %q", name, err.Fence, fenceCI)
		}
	}
	// NEGATIVE CONTROL in the same test: with neither set, the tripwire is silent.
	if err := refuseAutomationEnvironment(noEnv); err != nil {
		t.Fatalf("the tripwire fired with no CI variable set: %v", err)
	}
}

// scriptedTerminal stands in for an opened /dev/tty: Read yields what the
// operator typed, Write records the prompt, Close is observed.
type scriptedTerminal struct {
	typed  *strings.Reader
	shown  bytes.Buffer
	closed bool
}

func (s *scriptedTerminal) Read(p []byte) (int, error)  { return s.typed.Read(p) }
func (s *scriptedTerminal) Write(p []byte) (int, error) { return s.shown.Write(p) }
func (s *scriptedTerminal) Close() error                { s.closed = true; return nil }

// terminalProbe is a process whose /dev/tty opened but whose stdin is NOT that
// terminal (stdin is the given file): an IDE pane, `< /dev/null`, a pipe.
func terminalProbe(t *testing.T, stdin fs.FileInfo, typed string) ttyProbe {
	t.Helper()
	return ttyProbe{stdin: stdin, ctty: regularFileInfo(t), term: &scriptedTerminal{typed: strings.NewReader(typed)}}
}

// TestRedirectedStdinIsNeverTheConfirmation is what R-TTY-CHARDEV and
// R-TTY-SAMEFILE used to guard, restated for the fence that reads from the
// controlling terminal. For a stdin that is a pipe-like regular file and for
// /dev/null (a character device, the case a naive isatty check admits):
//
//   - the exact phrase on stdin, nothing typed at the terminal -> STOP
//     confirmation/eof. `echo phrase |` and `< /dev/null` are not an operator;
//   - nothing on stdin, the phrase typed at the terminal -> the fence passes,
//     and the prompt went to the terminal, not to stdout.
//
// Each arm is the other's control: the same probe shape, differing only in
// which device carries the line.
func TestRedirectedStdinIsNeverTheConfirmation(t *testing.T) {
	stdins := map[string]fs.FileInfo{"regular file": regularFileInfo(t), os.DevNull: devNullInfo(t)}
	for name, stdin := range stdins {
		t.Run(name, func(t *testing.T) {
			var out bytes.Buffer
			piped := terminalProbe(t, stdin, "")
			err := requireAttendedOperator(strings.NewReader(attendedPhrase+"\n"), &out, noEnv, piped)
			if err == nil || err.Line() != "STOP fence=confirmation reason=eof" {
				t.Fatalf("a phrase on a redirected stdin was accepted or refused oddly: %v", err)
			}
			if !piped.term.(*scriptedTerminal).closed {
				t.Fatal("the terminal was left open after the refusal")
			}

			out.Reset()
			typed := terminalProbe(t, stdin, attendedPhrase+"\n")
			if err := requireAttendedOperator(strings.NewReader(""), &out, noEnv, typed); err != nil {
				t.Fatalf("the phrase typed at the controlling terminal was refused: %v", err)
			}
			term := typed.term.(*scriptedTerminal)
			if !strings.Contains(term.shown.String(), attendedPhrase) {
				t.Fatalf("the prompt did not reach the terminal: %q", term.shown.String())
			}
			if out.Len() != 0 {
				t.Fatalf("the prompt went to stdout %q, not to the terminal the line is read from", out.String())
			}
			if !term.closed {
				t.Fatal("the terminal was left open after the confirmation")
			}
		})
	}
}

// TestEveryAttendedRefusalEndsWithItsFix: an operator who reads a STOP line
// gets the next thing to do, as the LAST thing printed.
func TestEveryAttendedRefusalEndsWithItsFix(t *testing.T) {
	cases := map[string]*stopError{
		"tty":      requireControllingTerminal(ttyProbe{cttyErr: fmt.Errorf("device not configured")}),
		"ci":       refuseAutomationEnvironment(ciEnv),
		"eof":      requireTypedPhrase(strings.NewReader(""), io.Discard, "p", "phrase"),
		"mismatch": requireTypedPhrase(strings.NewReader("nope\n"), io.Discard, "p", "phrase"),
	}
	want := map[string]string{"tty": ttyFix, "ci": ciFix, "eof": phraseFix, "mismatch": phraseFix}
	for name, err := range cases {
		if err == nil {
			t.Fatalf("%s: no refusal", name)
		}
		var buf bytes.Buffer
		report(&buf, err)
		if !strings.HasSuffix(strings.TrimSpace(buf.String()), want[name]) {
			t.Errorf("%s refusal does not end with its fix %q:\n%s", name, want[name], buf.String())
		}
		if !strings.HasPrefix(want[name], "fix: ") {
			t.Errorf("%s fix %q does not start with \"fix: \"", name, want[name])
		}
	}
}

// TestSameFileFailsClosedOnAMissingObservation pins the nil handling. It matters
// because MUT-D0-04 neuters R-TTY-OPEN and leaves ctty nil; os.SameFile would
// panic there, and a mutant that panics reds as a crash rather than as the
// refusal the row names.
func TestSameFileFailsClosedOnAMissingObservation(t *testing.T) {
	info := devNullInfo(t)
	if sameFile(nil, info) || sameFile(info, nil) || sameFile(nil, nil) {
		t.Fatal("sameFile reported a match for an observation that was never made")
	}
	if !sameFile(info, info) {
		t.Fatal("instrument failure: sameFile does not match a file with itself")
	}
}
