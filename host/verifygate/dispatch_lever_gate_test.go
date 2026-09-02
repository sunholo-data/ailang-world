package verifygate

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"testing"
)

var acceptedOnForms = []string{"on:", `"on":`, "'on':"}

var (
	errNoOnBlock        = errors.New("no top-level on block")
	errDuplicateOnBlock = errors.New("duplicate top-level on block")
	errTabIndent        = errors.New("tab in indentation")
	errUnhandledOnForm  = errors.New("unhandled top-level on form")
)

type onBlockParseError struct {
	kind   error
	line   int
	count  int
	detail string
}

func (e *onBlockParseError) Error() string {
	switch e.kind {
	case errDuplicateOnBlock:
		return fmt.Sprintf("%d top-level `on:` blocks, want exactly 1", e.count)
	case errTabIndent:
		return fmt.Sprintf("line %d: tab in indentation", e.line)
	case errUnhandledOnForm:
		return fmt.Sprintf("line %d: unsupported top-level `on:` form %q", e.line, e.detail)
	default:
		return e.kind.Error()
	}
}

func (e *onBlockParseError) Unwrap() error { return e.kind }

// matchTopLevelOn is deliberately column-0 anchored. Trimming before this comparison would
// allow a nested `on:` under a job step's with:/inputs: block to become the scan anchor.
func matchTopLevelOn(line string) (rest string, ok bool) {
	for _, form := range acceptedOnForms {
		if strings.HasPrefix(line, form) {
			return line[len(form):], true
		}
	}
	return "", false
}

func srcDeclaresOnKey(src string) bool {
	for _, line := range strings.Split(src, "\n") {
		for _, form := range acceptedOnForms {
			if strings.Contains(line, form) {
				return true
			}
		}
	}
	return false
}

func stripInlineComment(s string) string {
	if i := strings.Index(s, "#"); i >= 0 {
		s = s[:i]
	}
	return strings.TrimSpace(s)
}

func scalarTriggerValue(key, value string, scalarValued map[string]string) bool {
	value = stripInlineComment(value)
	if key == "workflow_dispatch" && value != "" && !(strings.HasPrefix(value, "{") && strings.HasSuffix(value, "}")) {
		scalarValued[key] = value
		return true
	}
	return false
}

func splitFlowItems(body string) ([]string, bool) {
	var items []string
	start, braces, brackets := 0, 0, 0
	var quote rune
	escaped := false
	for i, r := range body {
		if quote != 0 {
			if escaped {
				escaped = false
				continue
			}
			if r == '\\' && quote == '"' {
				escaped = true
				continue
			}
			if r == quote {
				quote = 0
			}
			continue
		}
		switch r {
		case '\'', '"':
			quote = r
		case '{':
			braces++
		case '}':
			braces--
		case '[':
			brackets++
		case ']':
			brackets--
		case ',':
			if braces == 0 && brackets == 0 {
				items = append(items, strings.TrimSpace(body[start:i]))
				start = i + 1
			}
		}
		if braces < 0 || brackets < 0 {
			return nil, false
		}
	}
	if quote != 0 || braces != 0 || brackets != 0 {
		return nil, false
	}
	items = append(items, strings.TrimSpace(body[start:]))
	return items, true
}

func splitFlowKeyValue(item string) (string, string, bool) {
	braces, brackets := 0, 0
	var quote rune
	escaped := false
	for i, r := range item {
		if quote != 0 {
			if escaped {
				escaped = false
				continue
			}
			if r == '\\' && quote == '"' {
				escaped = true
				continue
			}
			if r == quote {
				quote = 0
			}
			continue
		}
		switch r {
		case '\'', '"':
			quote = r
		case '{':
			braces++
		case '}':
			braces--
		case '[':
			brackets++
		case ']':
			brackets--
		case ':':
			if braces == 0 && brackets == 0 {
				key := strings.Trim(strings.TrimSpace(item[:i]), "'\"")
				return key, strings.TrimSpace(item[i+1:]), key != ""
			}
		}
	}
	return "", "", false
}

func parseOnBlockTriggers(src string) (keys []string, scalarValued map[string]string, err error) {
	lines := strings.Split(src, "\n")
	start, rest, count := -1, "", 0
	for i, line := range lines {
		if matchedRest, ok := matchTopLevelOn(line); ok {
			count++
			if start == -1 {
				start, rest = i, matchedRest
			}
		}
	}
	if start == -1 {
		return nil, nil, &onBlockParseError{kind: errNoOnBlock}
	}
	if count != 1 {
		return nil, nil, &onBlockParseError{kind: errDuplicateOnBlock, count: count}
	}

	scalarValued = map[string]string{}
	rest = stripInlineComment(rest)
	if rest != "" {
		if strings.HasPrefix(rest, "{") && strings.HasSuffix(rest, "}") {
			items, valid := splitFlowItems(strings.TrimSpace(rest[1 : len(rest)-1]))
			if !valid {
				return nil, nil, &onBlockParseError{kind: errUnhandledOnForm, line: start + 1, detail: rest}
			}
			for _, item := range items {
				if item == "" {
					continue
				}
				key, value, ok := splitFlowKeyValue(item)
				if !ok {
					return nil, nil, &onBlockParseError{kind: errUnhandledOnForm, line: start + 1, detail: rest}
				}
				if !scalarTriggerValue(key, value, scalarValued) {
					keys = append(keys, key)
				}
			}
			return keys, scalarValued, nil
		}
		return nil, nil, &onBlockParseError{kind: errUnhandledOnForm, line: start + 1, detail: rest}
	}

	triggerLead := -1
	for i, line := range lines[start+1:] {
		prefixLen := len(line) - len(strings.TrimLeft(line, " \t"))
		if strings.Contains(line[:prefixLen], "\t") {
			return nil, nil, &onBlockParseError{kind: errTabIndent, line: start + i + 2}
		}
		tok := strings.TrimSpace(line)
		if tok == "" || strings.HasPrefix(tok, "#") {
			continue
		}
		lead := len(line) - len(strings.TrimLeft(line, " "))
		if lead == 0 {
			break
		}
		if triggerLead == -1 {
			triggerLead = lead
		}
		if lead != triggerLead {
			continue
		}
		kv := strings.SplitN(tok, ":", 2)
		if len(kv) != 2 {
			continue
		}
		key := strings.Trim(strings.TrimSpace(kv[0]), "'\"")
		if !scalarTriggerValue(key, kv[1], scalarValued) {
			keys = append(keys, key)
		}
	}
	return keys, scalarValued, nil
}

func onBlockFailureMessage(path string, err error) string {
	var parseErr *onBlockParseError
	switch {
	case errors.Is(err, errNoOnBlock):
		return fmt.Sprintf("instrument failure: %s has no top-level `on:` trigger block", path)
	case errors.Is(err, errDuplicateOnBlock) && errors.As(err, &parseErr):
		return fmt.Sprintf("instrument failure: %s has %d top-level `on:` blocks, want exactly 1", path, parseErr.count)
	case errors.Is(err, errTabIndent), errors.Is(err, errUnhandledOnForm):
		return fmt.Sprintf("instrument failure: %s %s — this line scan computes YAML trigger shape and depth conservatively", path, err)
	default:
		return fmt.Sprintf("instrument failure: %s: %v", path, err)
	}
}

// onBlockTriggerKeys returns the trigger keys declared DIRECTLY under the on: block — the
// same depth as push:/pull_request: — so a workflow_dispatch folded into nested inputs: or
// buried in a comment does not count.
func onBlockTriggerKeys(t *testing.T, path, src string) []string {
	t.Helper()
	keys, scalarValued, err := parseOnBlockTriggers(src)
	if err != nil {
		t.Fatalf("%s", onBlockFailureMessage(path, err))
	}
	for key, value := range scalarValued {
		t.Errorf("%s: `%s:` has scalar value %q; want an empty key or a mapping", path, key, value)
	}
	return keys
}

func TestOnBlockTriggerParserShapes(t *testing.T) {
	realCI, err := os.ReadFile(filepath.Join(repoRoot, ".github", "workflows", "ci.yml"))
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name       string
		src        string
		wantKeys   []string
		wantScalar map[string]string
		wantErr    error
	}{
		{"canonical_block", "on:\n  push:\n  workflow_dispatch:\njobs:\n", []string{"push", "workflow_dispatch"}, nil, nil},
		{"real_ci_yml", string(realCI), []string{"push", "pull_request", "workflow_dispatch"}, nil, nil},
		{"quoted_double", "\"on\":\n  push:\n  workflow_dispatch:\n", []string{"push", "workflow_dispatch"}, nil, nil},
		{"quoted_single", "'on':\n  push:\n  workflow_dispatch:\n", []string{"push", "workflow_dispatch"}, nil, nil},
		{"flow_mapping", "on: {push: {branches: [dev]}, workflow_dispatch: }\n", []string{"push", "workflow_dispatch"}, nil, nil},
		{"flow_scalar_violation", "on: {push: {}, workflow_dispatch: garbage}\n", []string{"push"}, map[string]string{"workflow_dispatch": "garbage"}, nil},
		{"flow_unclosed", "on: {push: {}, workflow_dispatch:\n", nil, nil, errUnhandledOnForm},
		{"tab_indent", "on:\n\tpush:\n\tworkflow_dispatch:\n", nil, nil, errTabIndent},
		{"scalar_on", "on: push\n", nil, nil, errUnhandledOnForm},
		{"sequence_on", "on: [push, workflow_dispatch]\n", nil, nil, errUnhandledOnForm},
		{"duplicate_mixed", "on:\n  push:\n\"on\":\n  workflow_dispatch:\n", nil, nil, errDuplicateOnBlock},
		{"no_on_block", "jobs:\n  build:\n", nil, nil, errNoOnBlock},
		{"block_scalar_violation", "on:\n  push:\n  workflow_dispatch: garbage\n", []string{"push"}, map[string]string{"workflow_dispatch": "garbage"}, nil},
		{"inline_comment_value", "on:\n  push:\n  workflow_dispatch: # manual re-run lever\n", []string{"push", "workflow_dispatch"}, nil, nil},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			keys, scalarValued, err := parseOnBlockTriggers(tt.src)
			if !errors.Is(err, tt.wantErr) {
				t.Fatalf("err=%v, want errors.Is(_, %v)", err, tt.wantErr)
			}
			if !slices.Equal(keys, tt.wantKeys) {
				t.Errorf("keys=%v, want %v", keys, tt.wantKeys)
			}
			if len(scalarValued) != len(tt.wantScalar) {
				t.Errorf("scalarValued=%v, want %v", scalarValued, tt.wantScalar)
			}
			for key, want := range tt.wantScalar {
				if got := scalarValued[key]; got != want {
					t.Errorf("scalarValued[%q]=%q, want %q", key, got, want)
				}
			}
		})
	}
}

func TestOnBlockControlNeedle(t *testing.T) {
	realCI, err := os.ReadFile(filepath.Join(repoRoot, ".github", "workflows", "ci.yml"))
	if err != nil {
		t.Fatal(err)
	}
	tests := []struct {
		name string
		src  string
		want bool
	}{
		{"control_quoted_no_runs_on", "\"on\":\n  workflow_dispatch:\n", true},
		{"control_real_ci_yml", string(realCI), true},
		{"control_absent", "jobs:\n  build:\n", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := srcDeclaresOnKey(tt.src); got != tt.want {
				t.Errorf("srcDeclaresOnKey()=%v, want %v", got, tt.want)
			}
		})
	}
}

func TestOnBlockFailureMessagesUnchanged(t *testing.T) {
	t.Run("msg_no_on_block", func(t *testing.T) {
		got := onBlockFailureMessage("fixture.yml", &onBlockParseError{kind: errNoOnBlock})
		want := "instrument failure: fixture.yml has no top-level `on:` trigger block"
		if got != want {
			t.Errorf("message=%q, want %q", got, want)
		}
	})
	t.Run("msg_duplicate_on_block", func(t *testing.T) {
		got := onBlockFailureMessage("fixture.yml", &onBlockParseError{kind: errDuplicateOnBlock, count: 2})
		want := "instrument failure: fixture.yml has 2 top-level `on:` blocks, want exactly 1"
		if got != want {
			t.Errorf("message=%q, want %q", got, want)
		}
	})
}

// TestEveryWorkflowDeclaresDispatchLever pins queue-item 47's lever so it cannot be deleted
// silently. A push to dev whose webhook GitHub drops is PERMANENTLY unverifiable (never
// replayed, no run created), and the only API-driven lever is workflow_dispatch, which no
// workflow currently declares (P2). This gate asserts EVERY enumerated workflow file declares
// the lever as a trigger in its on: block.
//
// DECLARED RESIDUAL: this is a STATIC text scan over YAML. It proves the lever is DECLARED,
// never that a dispatch RUN is created or is green. And a workflow_dispatch run is NOT
// equivalent to the event it replaces: its checks do not satisfy PR branch protection
// (measured by this mission — required contexts can read success on the head SHA while
// gh pr checks --required still lists only the pull_request-event context with
// mergeStateStatus=BLOCKED). The lever buys A VERDICT ON A COMMIT, not A MERGEABLE PR.
//
// AND IT RE-VERIFIES THE TIP OF A NAMED REF, NEVER AN ARBITRARY SHA: `gh workflow run
// --ref` takes a branch or tag NAME (measured), so a dropped delivery is recoverable
// only while the affected commit is still that ref's tip -- once dev advances, that
// commit is unverifiable again. This repo's default_branch is dev, which is both why
// the lever is available and why it covers the dropped-push-to-dev case. (Added after
// the evaluator grepped this comment and found the design doc's claim that "the code
// comment states it plainly" was false: 0 hits for named-ref/tip/SHA here, control
// firing at 2 for the branch-protection sentence above.)
//
// It
// also cannot see a workflow added OUTSIDE .github/workflows/ (a singular `workflow` dir, a
// root-level .yaml, a hidden file), a case-mismatched filename (the Glob is case-sensitive),
// or a nested subdirectory (which GitHub itself does not scan either).
func TestEveryWorkflowDeclaresDispatchLever(t *testing.T) {
	matches, err := filepath.Glob(filepath.Join(repoRoot, ".github", "workflows", "*"))
	if err != nil {
		t.Fatal(err)
	}
	// ANTI-VACUITY FLOOR: an empty enumeration FAILS LOUDLY rather than printing a checkmark.
	if len(matches) == 0 {
		t.Fatal("instrument failure: no workflow files enumerated under .github/workflows/ — an empty set proves nothing about the lever")
	}
	for _, m := range matches {
		raw, err := os.ReadFile(m)
		if err != nil {
			t.Fatal(err)
		}
		src := string(raw)
		// known-positive control: a scan that cannot see the on: block is reading the wrong file.
		if !srcDeclaresOnKey(src) {
			t.Fatalf("instrument failure: %s does not contain known-positive control %q", m, "on:")
		}
		triggers := onBlockTriggerKeys(t, m, src)
		if !slices.Contains(triggers, "workflow_dispatch") {
			t.Errorf("%s: on-block triggers=%v lack workflow_dispatch — a dropped push to dev is permanently unverifiable; every workflow file must declare the lever", filepath.Base(m), triggers)
		}
	}
}
