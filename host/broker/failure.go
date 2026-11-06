package broker

// Row 255 (design_docs/planned/w-exec-failure-carries-a-reason.md): a handler
// failure that carries a bounded, recordable reason. The broker records
// Line() in the failed effect record's `failure` field; the coordinator
// returns the same line in the caller's world block, and `why` prints it.
// Only a *HandlerFailure has a line: any other handler error records as
// before, byte for byte.

import (
	"errors"
	"path/filepath"
	"strconv"
	"strings"
	"time"
	"unicode"
	"unicode/utf8"
)

// FailureLineMax bounds a failure line in bytes, ellipsis included.
const FailureLineMax = 200

// The Workspace.Exec failure classes (§3.2). The set is closed: FailureLine
// answers "" for any other class, so nothing unclassified is recorded.
const (
	FailExecMisrouted     = "exec-misrouted"
	FailExecTools         = "exec-tools"
	FailExecStackDrift    = "exec-stack-drift"
	FailExecSettingsDrift = "exec-settings-drift"
	FailProjectRoot       = "project-root"
	FailExecCache         = "exec-cache"
	FailSeedCache         = "seed-cache"
	FailTimeBudgetSpent   = "time-budget-spent"
	FailCleanupIncomplete = "cleanup-incomplete"
	FailSandboxStart      = "sandbox-start"
)

var failureClasses = map[string]bool{
	FailExecMisrouted: true, FailExecTools: true, FailExecStackDrift: true, FailExecSettingsDrift: true,
	FailProjectRoot: true, FailExecCache: true, FailSeedCache: true, FailTimeBudgetSpent: true,
	FailCleanupIncomplete: true, FailSandboxStart: true,
}

// HandlerFailure is a handler error with a class and a detail built from
// World's own words and disclosed paths. Error and Unwrap are the cause's, so
// every caller that matches on the cause is unchanged; only Line is recorded.
type HandlerFailure struct {
	Class  string
	Detail string
	Err    error
}

func (f *HandlerFailure) Error() string {
	if f.Err == nil {
		return f.Line()
	}
	return f.Err.Error()
}

func (f *HandlerFailure) Unwrap() error { return f.Err }

// Line is the recorded, returned and logged reason.
func (f *HandlerFailure) Line() string { return FailureLine(f.Class, f.Detail) }

// FailureLine is the pure bounding rule: a known class, then ": " and the
// detail on one line (control and format characters, such as a bidi
// override, and runs of spaces become one space; invalid UTF-8 becomes "?"),
// cut at a rune boundary to
// FailureLineMax bytes with a trailing "…". An unknown class gives "".
func FailureLine(class, detail string) string {
	if !failureClasses[class] {
		return ""
	}
	line := class
	clean := strings.Map(func(r rune) rune {
		if unicode.IsControl(r) || unicode.Is(unicode.Cf, r) {
			return ' '
		}
		return r
	}, strings.ToValidUTF8(detail, "?"))
	if d := strings.Join(strings.Fields(clean), " "); d != "" {
		line += ": " + d
	}
	if len(line) <= FailureLineMax {
		return line
	}
	const ellipsis = "…"
	cut := FailureLineMax - len(ellipsis)
	for cut > 0 && !utf8.RuneStart(line[cut]) {
		cut--
	}
	return strings.TrimRight(line[:cut], " ") + ellipsis
}

// ValidFailureLine is true when s is a line FailureLine can produce: it is a
// fixed point of the rule. A record whose line is not refuses to decode.
func ValidFailureLine(s string) bool {
	class, detail, _ := strings.Cut(s, ": ")
	return s != "" && FailureLine(class, detail) == s
}

// execDisclosed is what an exec failure line may name verbatim: the paths a
// successful Workspace.Exec result already returns to the agent (§4.5:
// profile.argv0 and argv0_target, sandbox.node and node_target). A path under
// the srt archive is named relative to it; any other path by its base name.
type execDisclosed struct {
	archive string
	paths   []string
}

func (h *ExecHandler) disclosed(sb *ExecSandbox, cmd ExecProfileCommand) execDisclosed {
	if sb == nil {
		sb = h.sandbox
	}
	return execDisclosed{archive: sb.Dir, paths: []string{sb.Node, sb.NodeConfigured, cmd.Argv0, cmd.ConfiguredArgv0}}
}

func (d execDisclosed) name(p string) string {
	for _, k := range d.paths {
		if k != "" && p == k {
			return p
		}
	}
	if d.archive != "" {
		if p == d.archive {
			return "srt archive"
		}
		if rel, err := filepath.Rel(d.archive, p); err == nil && rel != ".." && !strings.HasPrefix(rel, "../") {
			return "srt archive " + filepath.ToSlash(rel)
		}
	}
	return filepath.Base(p)
}

// execFailure is Execute's one classification point. class and detail name
// the step that stopped; a typed cause overrides them with its own class and a
// detail from fixed words and disclosed paths. err's text never reaches the line.
func execFailure(class, detail string, err error, d execDisclosed) *HandlerFailure {
	var hf *HandlerFailure
	var drift *ExecStackDriftError
	var settings *ExecSettingsError
	switch {
	case errors.As(err, &hf):
		return hf
	case errors.As(err, &drift):
		return &HandlerFailure{Class: FailExecStackDrift, Detail: d.name(drift.Path) + ": " + driftWhat(drift.Why), Err: err}
	case errors.As(err, &settings):
		return &HandlerFailure{Class: FailExecSettingsDrift, Detail: filepath.Base(settings.Path) + ": " + settingsWhat(settings.Why), Err: err}
	}
	return &HandlerFailure{Class: class, Detail: detail, Err: err}
}

// driftWhat maps an ExecStackDriftError's Why to fixed words: verifyArchive's
// three own phrases pass, a digest mismatch is "digest changed", and anything
// else (an OS error, which names host paths) is "unreadable".
func driftWhat(why string) string {
	switch {
	case why == "changed since it was archived", why == "missing from the archive", why == "not in the archived manifest":
		return why
	case strings.HasPrefix(why, "sha256 "):
		return "digest changed"
	}
	return "unreadable"
}

func settingsWhat(why string) string {
	switch {
	case strings.HasPrefix(why, "mode "):
		return "not a regular 0600 file"
	case why == "its bytes differ from the canonical rendering":
		return "bytes differ from World's rendering"
	}
	return "unreadable"
}

// seedWhat names how a seed copy failed, from the run error's type only.
func seedWhat(err error) string {
	var timeout *HandlerTimeoutError
	code, exited := handlerExitStatus(err)
	switch {
	case errors.As(err, &timeout):
		return "timed out after " + timeout.Timeout.Round(time.Millisecond).String()
	case cleanupIncomplete(err):
		return "the copy was not reaped"
	case exited:
		return "the copy exited " + strconv.Itoa(code)
	}
	return "the copy could not run"
}
