package broker

// Row 140 M1 (w-workspace-exec-toolchain-effect §4.3, §6 AC2.4's grammar
// half): MatchExecArgs against the design's worked `go test` profile command.
// AC2.4's refusals are pure, so they are pinned here at the law; M2 adds the
// handler arm (refused with an exec counter of 0, no srt spawn).

import (
	"errors"
	"reflect"
	"strconv"
	"strings"
	"testing"
)

// goTestCommand is §4.3's worked "test" command.
var goTestCommand = ExecCommand{
	Argv:       []string{"go", "test", "-count=1"},
	Flags:      map[string]string{"-run": ExecClassRegex, "-v": ExecClassBool, "-short": ExecClassBool},
	Positional: ExecPositionalPkgPattern,
	MaxArgs:    8,
}

// pytestFileCommand declares passthrough and the testfile class.
var pytestFileCommand = ExecCommand{
	Argv:        []string{"python", "-m", "pytest", "-q"},
	Flags:       map[string]string{"-k": ExecClassRegex, "-x": ExecClassBool, "--maxfail": ExecClassInt},
	Positional:  ExecPositionalTestfile,
	Suffixes:    []string{".py"},
	MaxArgs:     6,
	Passthrough: ExecPassthrough,
}

func TestMatchExecArgsAdmitsAndNormalizes(t *testing.T) {
	for _, tc := range []struct {
		name string
		cmd  ExecCommand
		args []string
		want []string
	}{
		{"-f v and -f=v mean one thing", goTestCommand, []string{"-run", "TestLex", "./internal/lexer/"}, []string{"-run=TestLex", "./internal/lexer/"}},
		{"already normalized", goTestCommand, []string{"-run=TestLex", "./internal/lexer/"}, []string{"-run=TestLex", "./internal/lexer/"}},
		{"bool flag consumes nothing", goTestCommand, []string{"-v", "./..."}, []string{"-v", "./..."}},
		{"no args", goTestCommand, []string{}, []string{}},
		{"value that looks like a flag", goTestCommand, []string{"-run", "-v"}, []string{"-run=-v"}},
		{"passthrough then a test file", pytestFileCommand, []string{"--maxfail", "2", "--", "tests/test_cli.py"},
			[]string{"--maxfail=2", "--", "tests/test_cli.py"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := MatchExecArgs(tc.cmd, tc.args)
			if err != nil || !reflect.DeepEqual(got, tc.want) {
				t.Fatalf("MatchExecArgs(%q) = (%q, %v), want %q", tc.args, got, err, tc.want)
			}
		})
	}
}

// TestMatchExecArgsRefusesAC24Grammar is AC2.4's grammar arms: each is refused
// naming the first bad arg. -c and -s=x.json are srt's own options and
// --settings its settings flag (V12): none is a flag the command lists, so no
// agent byte can reach srt's option parser.
func TestMatchExecArgsRefusesAC24Grammar(t *testing.T) {
	for _, tc := range []struct {
		name   string
		args   []string
		index  int
		reason string
	}{
		{"unprofiled flag", []string{"-race", "./..."}, 0, ExecArgFlag},
		{"-exec=sh", []string{"./...", "-exec=sh"}, 1, ExecArgFlag},
		{"-exec sh", []string{"-exec", "sh"}, 0, ExecArgFlag},
		{"-- without passthrough", []string{"--", "./..."}, 0, ExecArgDashDash},
		{"-c", []string{"-c"}, 0, ExecArgFlag},
		{"-s=x.json", []string{"-s=x.json"}, 0, ExecArgFlag},
		{"../x", []string{"../x"}, 0, ExecArgPositional},
		{"/abs", []string{"/abs"}, 0, ExecArgPositional},
		{"flag value with a newline", []string{"-run=a\nb"}, 0, ExecArgValue},
		{"split flag value with a newline", []string{"-run", "a\nb"}, 1, ExecArgValue},
		{"exactly --settings", []string{"--settings"}, 0, ExecArgFlag},
		{"--settings with a value", []string{"--settings=/tmp/s.json"}, 0, ExecArgFlag},
		{"prefix of a listed flag", []string{"-ru=X"}, 0, ExecArgFlag},
		{"listed flag as a prefix", []string{"-runX"}, 0, ExecArgFlag},
		{"bool flag given a value", []string{"-v=false"}, 0, ExecArgValue},
		{"valued flag at the end", []string{"./...", "-run"}, 1, ExecArgMissingValue},
		{"max_args", []string{"-v", "a/", "b/", "c/", "d/", "e/", "f/", "g/", "h/"}, 8, ExecArgMaxArgs},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, err := MatchExecArgs(goTestCommand, tc.args)
			var bad *ExecArgError
			if !errors.As(err, &bad) || bad.Index != tc.index || bad.Reason != tc.reason || bad.Arg != tc.args[tc.index] || got != nil {
				t.Fatalf("MatchExecArgs(%q) = (%q, %v), want arg %d refused for %q", tc.args, got, err, tc.index, tc.reason)
			}
			// The arg is named quoted, so a newline in it never reaches a log raw.
			if !strings.Contains(err.Error(), strconv.Quote(tc.args[tc.index])) {
				t.Fatalf("refusal %q does not name the bad arg %q", err, tc.args[tc.index])
			}
		})
	}
}

// TestMatchExecArgsPassthroughIsPositionalOnly: after the declared "--" a
// flag-shaped arg is not a flag, and a positional must still be a test file.
func TestMatchExecArgsPassthroughIsPositionalOnly(t *testing.T) {
	for _, tc := range []struct {
		args   []string
		index  int
		reason string
	}{
		{[]string{"--", "-k"}, 1, ExecArgPositional},
		{[]string{"--", "tests/conftest.txt"}, 1, ExecArgPositional},
		{[]string{"--", "../t.py"}, 1, ExecArgPositional},
		{[]string{"--maxfail=x"}, 0, ExecArgValue},
		{[]string{"--maxfail", "1000001"}, 1, ExecArgValue},
	} {
		_, err := MatchExecArgs(pytestFileCommand, tc.args)
		var bad *ExecArgError
		if !errors.As(err, &bad) || bad.Index != tc.index || bad.Reason != tc.reason {
			t.Fatalf("MatchExecArgs(%q) = %v, want arg %d refused for %q", tc.args, err, tc.index, tc.reason)
		}
	}
}
