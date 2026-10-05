package broker

// Row 140 M1 (design_docs/planned/w-workspace-exec-toolchain-effect.md §4.2
// gate 2, §4.3): the Workspace.Exec argument-matching law. Before anything
// is spawned, the handler matches the agent's args against the one profile
// command the call names; only flags the command lists, values of their
// declared class, positionals of its declared class and "--" where it
// declares passthrough get through, normalized so "-f v" and "-f=v" record
// the same argv. The law is authored as the contracted sketch
// design_docs/sketches/execargs.ail; every function below is the Go mirror of
// the sketch law of the same name, bound by execargs_drift_test.go, which
// evaluates the sketch's own test rows here.
//
// The grammar is not an execution fence (§4.3, R-140-3): running a project's
// tests runs project code. Its job is to keep the agent off srt's options and
// off any shell, keep the executable the operator's choice, and make every
// call's meaning an exact, recorded argv.

import (
	"fmt"
	"strconv"
	"strings"
)

// Flag value classes (§4.3). An enum class is spelled "enum:[a,b,...]".
const (
	ExecClassBool       = "bool"
	ExecClassRegex      = "regex"
	ExecClassInt        = "int"
	ExecClassRelpath    = "relpath"
	execClassEnumPrefix = "enum:"
)

// Positional classes (§4.3); a command with no class admits no positional.
const (
	ExecPositionalRelpath    = "relpath"
	ExecPositionalPkgPattern = "pkgpattern"
	ExecPositionalTestfile   = "testfile"
)

// ExecPassthrough is the only passthrough a command may declare.
const ExecPassthrough = "--"

// ExecCommand is one command of an operator exec profile (§4.3). Argv is the
// fixed prefix; the agent never supplies it, so MatchExecArgs does not read
// it. Flags maps each listed flag's exact name to its value class.
type ExecCommand struct {
	Argv        []string          `json:"argv"`
	Flags       map[string]string `json:"flags,omitempty"`
	Positional  string            `json:"positional,omitempty"`
	Suffixes    []string          `json:"suffixes,omitempty"`
	MaxArgs     int               `json:"max_args"`
	Passthrough string            `json:"passthrough,omitempty"`
	// Forms maps a valued flag to how the EMITTED argv spells it (row 140
	// M2, EmitExecArgs); matching never reads it.
	Forms map[string]string `json:"-"`
}

// Why an agent arg was refused (the sketch's reason words).
const (
	ExecArgFlag         = "flag"
	ExecArgValue        = "value"
	ExecArgMissingValue = "missing-value"
	ExecArgDashDash     = "dash-dash"
	ExecArgPositional   = "positional"
	ExecArgMaxArgs      = "max-args"
)

var execArgReasonText = map[string]string{
	ExecArgFlag:         "the command does not list this flag",
	ExecArgValue:        "the value does not match the flag's class",
	ExecArgMissingValue: "the flag takes a value and none follows",
	ExecArgDashDash:     `"--" is admitted only where the command declares passthrough`,
	ExecArgPositional:   "not a positional of the command's class",
	ExecArgMaxArgs:      "more arguments than the command's max_args",
}

// ExecArgError names the first agent arg the grammar refused.
type ExecArgError struct {
	Index  int    // position in the agent's args
	Arg    string // the arg itself
	Reason string // one of the ExecArg* words
}

func (e *ExecArgError) Error() string {
	return fmt.Sprintf("argument %d %q refused: %s", e.Index, e.Arg, execArgReasonText[e.Reason])
}

// MatchExecArgs applies the §4.3 law (sketch matchArgs) to the agent's args
// and returns them normalized, or the first bad arg.
func MatchExecArgs(cmd ExecCommand, args []string) ([]string, error) {
	out := make([]string, 0, len(args))
	passthrough := cmd.Passthrough == ExecPassthrough
	afterDashDash := false
	bad := func(i int, reason string) ([]string, error) {
		return nil, &ExecArgError{Index: i, Arg: args[i], Reason: reason}
	}
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case afterDashDash:
			if !positionalOk(cmd.Positional, a, suffixIn(cmd.Suffixes, a)) {
				return bad(i, ExecArgPositional)
			}
		case a == "--":
			if !dashDashOk(a, passthrough) {
				return bad(i, ExecArgDashDash)
			}
			afterDashDash = true
		case strings.HasPrefix(a, "-"):
			cls := flagClassOf(cmd.Flags, a)
			hasEq := strings.Contains(a, "=")
			if cls == "" {
				return bad(i, ExecArgFlag)
			}
			if hasEq {
				if !execValueOk(cls, flagValueOf(a)) {
					return bad(i, ExecArgValue)
				}
			} else if consumesNext(cls, hasEq) {
				if i+1 >= len(args) {
					return bad(i, ExecArgMissingValue)
				}
				if !execValueOk(cls, args[i+1]) {
					return bad(i+1, ExecArgValue)
				}
				if !argCountOk(len(out)+1, cmd.MaxArgs) {
					return bad(i, ExecArgMaxArgs)
				}
				out = append(out, a+"="+args[i+1])
				i++
				continue
			}
		default:
			if !positionalOk(cmd.Positional, a, suffixIn(cmd.Suffixes, a)) {
				return bad(i, ExecArgPositional)
			}
		}
		if !argCountOk(len(out)+1, cmd.MaxArgs) {
			return bad(i, ExecArgMaxArgs)
		}
		out = append(out, a)
	}
	return out, nil
}

// execNoDotDotSegment mirrors the sketch's noDotDotSegment.
func execNoDotDotSegment(p string) bool {
	return p != ".." && !strings.HasPrefix(p, "../") && !strings.HasSuffix(p, "/..") && !strings.Contains(p, "/../")
}

// execPathOk mirrors the sketch's pathOk (relpath).
func execPathOk(p string) bool {
	return len(p) >= 1 && !strings.HasPrefix(p, "/") && !strings.HasPrefix(p, "-") && execNoDotDotSegment(p)
}

func flagNameOf(a string) string {
	if i := strings.Index(a, "="); i >= 0 {
		return a[:i]
	}
	return a
}

func flagValueOf(a string) string {
	if i := strings.Index(a, "="); i >= 0 {
		return a[i+1:]
	}
	return ""
}

func flagRuleMatches(name, a string) bool { return name == flagNameOf(a) }

func consumesNext(cls string, hasEq bool) bool { return !hasEq && cls != ExecClassBool }

func digitByteOk(c int) bool { return c >= '0' && c <= '9' }

func intValueOk(n int) bool { return n >= 0 && n <= 1000000 }

func regexValueOk(v string, n int) bool {
	return n >= 1 && n <= 256 && !strings.Contains(v, "\n") && !strings.Contains(v, "\x00")
}

func flagValueOk(cls, v string, n int, intOk, enumOk bool) bool {
	return (cls == ExecClassRegex && regexValueOk(v, n)) || (cls == ExecClassRelpath && execPathOk(v)) ||
		(cls == ExecClassInt && intOk) || (strings.HasPrefix(cls, execClassEnumPrefix) && enumOk)
}

func pkgpatternOk(p string) bool {
	if p == "./..." || (execPathOk(p) && !strings.Contains(p, "...")) {
		return true
	}
	if len(p) < 5 || !strings.HasSuffix(p, "/...") {
		return false
	}
	prefix := p[:len(p)-4]
	return execPathOk(prefix) && !strings.Contains(prefix, "...")
}

func testfileOk(p, suffix string) bool {
	return execPathOk(p) && len(suffix) >= 1 && strings.HasSuffix(p, suffix) && len(p) > len(suffix)
}

func positionalOk(cls, p string, suffixMatched bool) bool {
	return (cls == ExecPositionalRelpath && execPathOk(p)) || (cls == ExecPositionalPkgPattern && pkgpatternOk(p)) ||
		(cls == ExecPositionalTestfile && suffixMatched)
}

func dashDashOk(a string, passthrough bool) bool { return a != "--" || passthrough }

func argCountOk(n, maxArgs int) bool { return n >= 0 && n <= maxArgs }

// intValueParses: 1 to 7 decimal digits whose value passes intValueOk.
func intValueParses(v string) bool {
	if len(v) == 0 || len(v) > 7 {
		return false
	}
	for i := 0; i < len(v); i++ {
		if !digitByteOk(int(v[i])) {
			return false
		}
	}
	n, err := strconv.Atoi(v)
	return err == nil && intValueOk(n)
}

// enumHas: v is one of the words of "enum:[a,b,...]".
func enumHas(cls, v string) bool {
	if !strings.HasPrefix(cls, "enum:[") || !strings.HasSuffix(cls, "]") || v == "" {
		return false
	}
	for _, w := range strings.Split(cls[len("enum:["):len(cls)-1], ",") {
		if w == v {
			return true
		}
	}
	return false
}

func suffixIn(suffixes []string, p string) bool {
	for _, s := range suffixes {
		if testfileOk(p, s) {
			return true
		}
	}
	return false
}

// flagClassOf is the class of the listed flag whose name EXACTLY equals the
// arg's flag name, or "" when the command does not list it. The map lookup
// is flagRuleMatches over every listed name at once: a key equal to the
// whole flag name, never a prefix.
func flagClassOf(flags map[string]string, a string) string {
	return flags[flagNameOf(a)]
}

func execValueOk(cls, v string) bool {
	return flagValueOk(cls, v, len(v), intValueParses(v), enumHas(cls, v))
}
